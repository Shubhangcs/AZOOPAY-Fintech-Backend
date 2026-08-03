package store

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/levionstudio/fintech/internal/models"
)

type AEPSStore interface {
	GetAEPSDetailsByRetailerID(retailerId string) (*models.AEPSDetailsModel, error)
	InitilizeCashWithdrawal(retailerId string, merchantData *models.AEPSDetailsModel, transactionData *models.AEPSCashWithdrawalRequestModel) (int64, error)
	FinilizeCashWithdrawal(aepsTransactionId int64, res *models.AEPSCashWithdrawalResponseModel) error
}

type PostgresAEPSStore struct {
	db             *sql.DB
	commisionStore CommisionStore
	walletStore    WalletTransactionStore
}

func NewPostgresAEPSStore(db *sql.DB, commisionStore CommisionStore, walletStore WalletTransactionStore) *PostgresAEPSStore {
	return &PostgresAEPSStore{
		db,
		commisionStore,
		walletStore,
	}
}

func (as *PostgresAEPSStore) GetAEPSDetailsByRetailerID(retailerId string) (*models.AEPSDetailsModel, error) {
	query := `
		SELECT
			r.retailer_aadhar_number,
			m.outlet_id,
			a.latitude,
			a.longitude
		FROM retailers r
		JOIN aeps_merchant_details m ON m.retailer_id = r.retailer_id
		JOIN aeps_applications a ON a.retailer_id = r.retailer_id
		WHERE r.retailer_id = $1;
	`
	var res models.AEPSDetailsModel
	if err := as.db.QueryRow(
		query,
		retailerId,
	).Scan(
		&res.AadhaarNumber,
		&res.OutletID,
		&res.Latitude,
		&res.Longitude,
	); err != nil {
		return nil, err
	}

	return &res, nil
}

func (as *PostgresAEPSStore) InitilizeCashWithdrawal(retailerId string, merchantData *models.AEPSDetailsModel, transactionData *models.AEPSCashWithdrawalRequestModel) (int64, error) {

	rtds, err := getRetailerDetails(as.db, retailerId)
	if err != nil {
		return 0, err
	}

	commision := as.commisionStore.GetAEPSCommision(transactionData.Amount)
	if commision == nil {
		return 0, errors.New("invalid amount")
	}

	query := `
		INSERT INTO aeps_transactions(
			retailer_id,
			reference_id,
			customer_name,
			customer_phone,
			customer_aadhaar,
			amount,
			md_commision,
			dis_commision,
			retailer_commision,
			transaction_status
		) VALUES (
			$1 , $2 , $3 , $4 , $5 , $6 , $7 , $8 , $9 , $10 
		) 
		RETURNING aeps_transaction_id;
	`

	tx, err := as.db.Begin()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	var aepsTransactionId int64
	if err := tx.QueryRow(
		query,
		retailerId,
		transactionData.RequestID,
		transactionData.CustomerName,
		transactionData.Mobile,
		transactionData.Aadhaar,
		transactionData.Amount,
		commision.MasterDistributorCommision,
		commision.DistributorCommision,
		commision.RetailerCommision,
		"PENDING",
	).Scan(
		&aepsTransactionId,
	); err != nil {
		return 0, err
	}

	// if err := debitTx(
	// 	tx,
	// 	transaction{
	// 		UserID:      retailerId,
	// 		ReferenceID: fmt.Sprintf("%d", aepsTransactionId),
	// 		Amount:      transactionData.Amount,
	// 		Reason:      "AEPS",
	// 		Remarks:     "AEPS transaction for Customer Name: " + transactionData.CustomerName + " Customer Phone: " + transactionData.Mobile + " Retailer id: " + retailerId,
	// 	},
	// 	as.walletStore,
	// ); err != nil {
	// 	return 0, err
	// }

	if commision.MasterDistributorCommision != 0 {
		if err := creditTx(
			tx,
			transaction{
				UserID:      rtds.mdID,
				ReferenceID: fmt.Sprintf("%d", aepsTransactionId),
				Amount:      commision.MasterDistributorCommision * 0.98,
				Reason:      "AEPS_COMMISION",
				Remarks:     "AEPS Commision for MD: " + rtds.mdID,
			},
			as.walletStore,
		); err != nil {
			return 0, err
		}

		// if err := tdsAepsTx(tx, &aepsTdsTransaction{
		// 	TransactionID: fmt.Sprintf("%d", aepsTransactionId),
		// 	TDSAmount:     commision.MasterDistributorCommision * 0.02,
		// 	UserID:        rtds.mdID,
		// 	UserType:      "MD",
		// }); err != nil {
		// 	return 0, err
		// }
	}

	if commision.DistributorCommision != 0 {
		if err := creditTx(
			tx,
			transaction{
				UserID:      rtds.distributorID,
				ReferenceID: fmt.Sprintf("%d", aepsTransactionId),
				Amount:      commision.DistributorCommision * 0.98,
				Reason:      "AEPS_COMMISION",
				Remarks:     "AEPS Commision for Distributor: " + rtds.mdID,
			},
			as.walletStore,
		); err != nil {
			return 0, err
		}
		// if err := tdsAepsTx(tx, &aepsTdsTransaction{
		// 	TransactionID: fmt.Sprintf("%d", aepsTransactionId),
		// 	TDSAmount:     commision.DistributorCommision * 0.02,
		// 	UserID:        rtds.distributorID,
		// 	UserType:      "DIS",
		// }); err != nil {
		// 	return 0, err
		// }
	}

	if commision.RetailerCommision != 0 {
		if err := creditTx(
			tx,
			transaction{
				UserID:      retailerId,
				ReferenceID: fmt.Sprintf("%d", aepsTransactionId),
				Amount:      commision.RetailerCommision * 0.98,
				Reason:      "AEPS_COMMISION",
				Remarks:     "AEPS Commision for Retailer: " + rtds.mdID,
			},
			as.walletStore,
		); err != nil {
			return 0, err
		}

		// if err := tdsAepsTx(tx, &aepsTdsTransaction{
		// 	TransactionID: fmt.Sprintf("%d", aepsTransactionId),
		// 	TDSAmount:     commision.RetailerCommision * 0.02,
		// 	UserID:        retailerId,
		// 	UserType:      "RT",
		// }); err != nil {
		// 	return 0, err
		// }
	}

	// adminCreditQuery := `
	// 	UPDATE admins
	// 	SET aeps_wallet = aeps_wallet + $1,
	// 		updated_at = NOW()
	// 	WHERE admin_id = $2
	// 	RETURNING aeps_wallet;
	// `

	// var adminAepsWalletAfterBalance float64
	// if err := tx.QueryRow(
	// 	adminCreditQuery,
	// 	transactionData.Amount,
	// 	rtds.adminID,
	// ).Scan(
	// 	&adminAepsWalletAfterBalance,
	// ); err != nil {
	// 	return 0, err
	// }

	// if err := as.walletStore.CreateWalletTransactionTx(
	// 	tx,
	// 	&models.WalletTransactionModel{
	// 		UserID:            rtds.adminID,
	// 		ReferenceID:       fmt.Sprintf("%d", aepsTransactionId),
	// 		CreditAmount:      &transactionData.Amount,
	// 		BeforeBalance:     adminAepsWalletAfterBalance - transactionData.Amount,
	// 		AfterBalance:      adminAepsWalletAfterBalance,
	// 		TransactionReason: "AEPS",
	// 		Remarks:           "AEPS from retailer " + retailerId,
	// 	},
	// ); err != nil {
	// 	return 0, err
	// }

	return aepsTransactionId, tx.Commit()
}

func (as *PostgresAEPSStore) FinilizeCashWithdrawal(aepsTransactionId int64, res *models.AEPSCashWithdrawalResponseModel) error {
	query := `
		UPDATE aeps_transactions
		SET transaction_status = COALESCE(NULLIF($1 , '') , transaction_status),
			transaction_id = COALESCE(NULLIF($2 , '') , transaction_id),
			order_id = COALESCE(NULLIF($3, '') , order_id)
		WHERE aeps_transaction_id = $4;
	`

	dbres, err := as.db.Exec(
		query,
		res.TransactionStatus,
		res.TransactionID,
		res.OrderID,
		aepsTransactionId,
	)

	if err != nil {
		return err
	}

	return checkRowsAffected(dbres)
}
