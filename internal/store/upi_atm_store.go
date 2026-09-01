package store

import (
	"database/sql"

	"github.com/levionstudio/fintech/internal/models"
)

type UPIATMStore interface {
	CreateQR(retailerId string, req *models.UPIATMCreateQRAPIResponseModel) error
	GetQRDetailsForStatusCheck(requestId string) (*models.UPIATMCheckQRTransactionStatusAPIRequestModel, error)
	FinilizeQRStatus(req *models.UPIATMCheckQRTransactionStatusAPIResponseModel) error
	GetUPIATMTransactionsByRetailerID(retailerId string) ([]models.UPIATMTransactionResponseModel, error)
	GetALLUPIATMTransactions() ([]models.UPIATMTransactionResponseModel, error)
}

type PostgresUPIATMStore struct {
	db *sql.DB
}

func NewPostgresUPIATMStore(db *sql.DB) *PostgresUPIATMStore {
	return &PostgresUPIATMStore{
		db,
	}
}

func (us *PostgresUPIATMStore) CreateQR(retailerId string, req *models.UPIATMCreateQRAPIResponseModel) error {
	query := `
		INSERT INTO upi_atm (
			retailer_id,
			request_id,
			transaction_id,
			ipay_id,
			amount,
			payable_value,
			transaction_value,
			commission_amount,
			tds_amount,
			net_amount,
			settlement_status,
			qr_status
		) VALUES (
			$1 , $2 , $3 , $4 , $5 , $6 , $7 , $8 , $9 , $10, $11 , $12
		);
	`

	res, err := us.db.Exec(
		query,
		retailerId,
		req.RequestID,
		req.TransactionID,
		req.IpayID,
		req.Amount,
		req.PayableValue,
		req.TransactionValue,
		req.CommissionAmount,
		req.TDSAmount,
		req.NETAmount,
		req.SettlementStatus,
		req.QRStatus,
	)
	if err != nil {
		return err
	}

	return checkRowsAffected(res)
}

func (us *PostgresUPIATMStore) GetQRDetailsForStatusCheck(requestId string) (*models.UPIATMCheckQRTransactionStatusAPIRequestModel, error) {
	query := `
		SELECT 
			o.request_id,
			o.transaction_id,
			o.ipay_id,
			n.outlet_id
		FROM upi_atm o
		JOIN aeps_merchant_details n ON n.retailer_id = o.retailer_id
		WHERE o.request_id = $1;
	`
	var res models.UPIATMCheckQRTransactionStatusAPIRequestModel
	if err := us.db.QueryRow(query, requestId).Scan(
		&res.RequestID,
		&res.TransactionID,
		&res.IpayID,
		&res.OutletID,
	); err != nil {
		return nil, err
	}

	return &res, nil
}

func (us *PostgresUPIATMStore) FinilizeQRStatus(req *models.UPIATMCheckQRTransactionStatusAPIResponseModel) error {
	query := `
		UPDATE upi_atm
		SET qr_status = $1,
			settlement_status = $2,
			updated_at = NOW()
		WHERE request_id = $3;
	`

	res, err := us.db.Exec(query, req.QRStatus, req.SettlementStatus, req.RequestID)
	if err != nil {
		return err
	}

	return checkRowsAffected(res)
}

func (us *PostgresUPIATMStore) GetUPIATMTransactionsByRetailerID(retailerId string) ([]models.UPIATMTransactionResponseModel, error) {
	query := `
		SELECT 
			t.retailer_id,
			t.request_id,
			t.transaction_id,
			t.ipay_id,
			t.amount,
			t.payable_value,
			t.transaction_value,
			t.commission_amount,
			t.tds_amount,
			t.net_amount,
			t.settlement_status,
			t.qr_status,
			t.created_at,
			t.updated_at,
			r.retailer_name
		FROM upi_atm t
		JOIN retailers r ON r.retailer_id = t.retailer_id
		WHERE t.retailer_id = $1;
	`

	res, err := us.db.Query(query, retailerId)
	if err != nil {
		return nil, err
	}
	defer res.Close()

	var txn models.UPIATMTransactionResponseModel
	var txns []models.UPIATMTransactionResponseModel
	for res.Next() {
		if err := res.Scan(
			&txn.RetailerID,
			&txn.RequestID,
			&txn.TransactionID,
			&txn.IpayID,
			&txn.Amount,
			&txn.PayableValue,
			&txn.TransactionValue,
			&txn.CommissionAmount,
			&txn.TDSAmount,
			&txn.NETAmount,
			&txn.SettlementStatus,
			&txn.QRStatus,
			&txn.CreatedAT,
			&txn.UpdatedAT,
			&txn.RetailerName,
		); err != nil {
			return nil, err
		}

		txns = append(txns, txn)
	}

	if err := res.Err(); err != nil {
		return nil, err
	}

	return txns, nil
}

func (us *PostgresUPIATMStore) GetALLUPIATMTransactions() ([]models.UPIATMTransactionResponseModel, error) {
	query := `
		SELECT 
			t.retailer_id,
			t.request_id,
			t.transaction_id,
			t.ipay_id,
			t.amount,
			t.payable_value,
			t.transaction_value,
			t.commission_amount,
			t.tds_amount,
			t.net_amount,
			t.settlement_status,
			t.qr_status,
			t.created_at,
			t.updated_at,
			r.retailer_name
		FROM upi_atm t
		JOIN retailers r ON r.retailer_id = t.retailer_id;
	`

	res, err := us.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer res.Close()

	var txn models.UPIATMTransactionResponseModel
	var txns []models.UPIATMTransactionResponseModel
	for res.Next() {
		if err := res.Scan(
			&txn.RetailerID,
			&txn.RequestID,
			&txn.TransactionID,
			&txn.IpayID,
			&txn.Amount,
			&txn.PayableValue,
			&txn.TransactionValue,
			&txn.CommissionAmount,
			&txn.TDSAmount,
			&txn.NETAmount,
			&txn.SettlementStatus,
			&txn.QRStatus,
			&txn.CreatedAT,
			&txn.UpdatedAT,
			&txn.RetailerName,
		); err != nil {
			return nil, err
		}

		txns = append(txns, txn)
	}

	if err := res.Err(); err != nil {
		return nil, err
	}

	return txns, nil
}
