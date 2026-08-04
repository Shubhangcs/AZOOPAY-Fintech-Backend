package store

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/levionstudio/fintech/internal/models"
	"github.com/levionstudio/fintech/internal/utils"
)

type AEPSStore interface {
	GetAEPSDetailsByRetailerID(retailerId string) (*models.AEPSDetailsModel, error)
	InitilizeCashWithdrawal(retailerId string, merchantData *models.AEPSDetailsModel, transactionData *models.AEPSCashWithdrawalRequestModel) (int64, error)
	FinilizeCashWithdrawal(aepsTransactionId int64, res *models.AEPSCashWithdrawalResponseModel) error
	GetAllAEPSTransactions(p utils.QueryParams) ([]models.AepsTransactionResponse, error)
	GetAEPSTransactionsByRetailerID(retailerId string, p utils.QueryParams) ([]models.AepsTransactionResponse, error)
	GetAllAEPSTDSDeductions(p utils.QueryParams) ([]models.AepsTdsDeductionResponse, error)
	GetAEPSTDSDeductionsByRetailerID(retailerId string, p utils.QueryParams) ([]models.AepsTdsDeductionResponse, error)
	GetAEPSTDSDeductionsByDistributorID(distributorId string, p utils.QueryParams) ([]models.AepsTdsDeductionResponse, error)
	GetAEPSTDSDeductionsByMDID(mdId string, p utils.QueryParams) ([]models.AepsTdsDeductionResponse, error)
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

	rtTableInfo, err := getUserTableInfo(retailerId)
	if err != nil {
		return 0, err
	}

	disTableInfo, err := getUserTableInfo(rtds.distributorID)
	if err != nil {
		return 0, err
	}

	mdTableInfo, err := getUserTableInfo(rtds.mdID)
	if err != nil {
		return 0, err
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
		commision.MasterDistributorCommision*0.98,
		commision.DistributorCommision*0.98,
		commision.RetailerCommision*0.98,
		"PENDING",
	).Scan(
		&aepsTransactionId,
	); err != nil {
		return 0, err
	}

	if err := creditTx(
		tx,
		transaction{
			UserID:        retailerId,
			ReferenceID:   fmt.Sprintf("%d", aepsTransactionId),
			Amount:        transactionData.Amount,
			Reason:        "AEPS",
			Remarks:       "AEPS Amount Credit To: " + retailerId,
			userTableInfo: *rtTableInfo,
		},
		as.walletStore,
	); err != nil {
		return 0, err
	}

	if commision.MasterDistributorCommision != 0 {
		if err := creditTx(
			tx,
			transaction{
				UserID:        rtds.mdID,
				ReferenceID:   fmt.Sprintf("%d", aepsTransactionId),
				Amount:        commision.MasterDistributorCommision * 0.98,
				Reason:        "AEPS_COMMISSION",
				Remarks:       "AEPS Commision for MD: " + rtds.mdID,
				userTableInfo: *mdTableInfo,
			},
			as.walletStore,
		); err != nil {
			return 0, err
		}

		if err := tdsAepsTx(tx, &aepsTdsTransaction{
			TransactionID: fmt.Sprintf("%d", aepsTransactionId),
			TDSAmount:     commision.MasterDistributorCommision * 0.02,
			UserID:        rtds.mdID,
			UserType:      "MD",
		}); err != nil {
			return 0, err
		}
	}

	if commision.DistributorCommision != 0 {
		if err := creditTx(
			tx,
			transaction{
				UserID:        rtds.distributorID,
				ReferenceID:   fmt.Sprintf("%d", aepsTransactionId),
				Amount:        commision.DistributorCommision * 0.98,
				Reason:        "AEPS_COMMISSION",
				Remarks:       "AEPS Commision for Distributor: " + rtds.distributorID,
				userTableInfo: *disTableInfo,
			},
			as.walletStore,
		); err != nil {
			return 0, err
		}

		if err := tdsAepsTx(tx, &aepsTdsTransaction{
			TransactionID: fmt.Sprintf("%d", aepsTransactionId),
			TDSAmount:     commision.DistributorCommision * 0.02,
			UserID:        rtds.distributorID,
			UserType:      "DIS",
		}); err != nil {
			return 0, err
		}
	}

	if commision.RetailerCommision != 0 {
		if err := creditTx(
			tx,
			transaction{
				UserID:        retailerId,
				ReferenceID:   fmt.Sprintf("%d", aepsTransactionId),
				Amount:        commision.RetailerCommision * 0.98,
				Reason:        "AEPS_COMMISSION",
				Remarks:       "AEPS Commision for Retailer: " + retailerId,
				userTableInfo: *rtTableInfo,
			},
			as.walletStore,
		); err != nil {
			return 0, err
		}

		if err := tdsAepsTx(tx, &aepsTdsTransaction{
			TransactionID: fmt.Sprintf("%d", aepsTransactionId),
			TDSAmount:     commision.RetailerCommision * 0.02,
			UserID:        retailerId,
			UserType:      "RT",
		}); err != nil {
			return 0, err
		}
	}

	adminCreditQuery := `
		UPDATE admins
		SET aeps_wallet = aeps_wallet + $1,
			updated_at = NOW()
		WHERE admin_id = $2
		RETURNING aeps_wallet;
	`

	var adminAepsWalletAfterBalance float64
	if err := tx.QueryRow(
		adminCreditQuery,
		transactionData.Amount,
		rtds.adminID,
	).Scan(
		&adminAepsWalletAfterBalance,
	); err != nil {
		return 0, err
	}

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

const aepsSelectBase = `
	SELECT 
		at.aeps_transaction_id,
		at.reference_id,
		at.transaction_id,
		at.order_id,
		at.customer_name,
		at.customer_phone,
		at.customer_aadhaar,
		at.amount,
		at.md_commision,
		at.dis_commision,
		at.retailer_commision,
		at.transaction_status,
		at.created_at,
		at.updated_at,
		r.retailer_name,
		w.before_balance,
		w.after_balance,
		w.transaction_reason,
		w.remarks
	FROM retailers r
	JOIN aeps_transactions at ON at.retailer_id = r.retailer_id
	LEFT JOIN wallet_transactions w ON w.user_id = r.retailer_id AND w.reference_id = at.aeps_transaction_id::VARCHAR AND w.transaction_reason = 'AEPS'
`

func (as *PostgresAEPSStore) GetAllAEPSTransactions(p utils.QueryParams) ([]models.AepsTransactionResponse, error) {
	q := aepsSelectBase + `
	WHERE at.created_at >= COALESCE($3, '-infinity'::TIMESTAMPTZ)
	AND at.created_at <= COALESCE($4, 'infinity'::TIMESTAMPTZ)
	AND ($5::TEXT IS NULL OR at.transaction_status = $5)
	AND ($6::TEXT IS NULL OR (
		at.aeps_transaction_id::TEXT ILIKE '%'||$6||'%' OR
		at.reference_id ILIKE '%'||$6||'%' OR
		at.transaction_id ILIKE '%'||$6||'%' OR
		at.order_id ILIKE '%'||$6||'%' OR
		at.customer_phone ILIKE '%'||$6||'%'
	))
	ORDER BY at.created_at DESC
	LIMIT $1 OFFSET $2;
	`
	return scanAepsTransactions(as.db, q, p.Limit, p.Offset, p.StartDate, p.EndDate, p.Status, p.Search)
}

func (as *PostgresAEPSStore) GetAEPSTransactionsByRetailerID(retailerId string, p utils.QueryParams) ([]models.AepsTransactionResponse, error) {
	q := aepsSelectBase + `
	WHERE r.retailer_id = $7
	AND at.created_at >= COALESCE($3, '-infinity'::TIMESTAMPTZ)
	AND at.created_at <= COALESCE($4, 'infinity'::TIMESTAMPTZ)
	AND ($5::TEXT IS NULL OR at.transaction_status = $5)
	AND ($6::TEXT IS NULL OR (
		at.aeps_transaction_id::TEXT ILIKE '%'||$6||'%' OR
		at.reference_id ILIKE '%'||$6||'%' OR
		at.transaction_id ILIKE '%'||$6||'%' OR
		at.order_id ILIKE '%'||$6||'%' OR
		at.customer_phone ILIKE '%'||$6||'%'
	))
	ORDER BY at.created_at DESC
	LIMIT $1 OFFSET $2;
	`
	return scanAepsTransactions(as.db, q, p.Limit, p.Offset, p.StartDate, p.EndDate, p.Status, p.Search, retailerId)
}

func scanAepsTransactions(db *sql.DB, query string, args ...any) ([]models.AepsTransactionResponse, error) {
	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query aeps transactions: %w", err)
	}
	defer rows.Close()

	var transactions []models.AepsTransactionResponse

	for rows.Next() {
		var t models.AepsTransactionResponse

		if err := rows.Scan(
			&t.AepsTransactionID,
			&t.ReferenceID,
			&t.TransactionID,
			&t.OrderID,
			&t.CustomerName,
			&t.CustomerPhone,
			&t.CustomerAadhaar,
			&t.Amount,
			&t.MdCommission,
			&t.DisCommission,
			&t.RetailerCommission,
			&t.TransactionStatus,
			&t.CreatedAt,
			&t.UpdatedAt,
			&t.RetailerName,
			&t.BeforeBalance,
			&t.AfterBalance,
			&t.Reason,
			&t.Remarks,
		); err != nil {
			return nil, fmt.Errorf("failed to scan aeps transaction row: %w", err)
		}

		transactions = append(transactions, t)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating aeps transaction rows: %w", err)
	}

	return transactions, nil
}

const aepsTdsSelectAllBase = `
	SELECT 
		t.tds_id,
		t.aeps_transaction_id,
		t.user_id,
		t.user_type,
		t.created_at,
		t.updated_at,
		at.customer_name,
		COALESCE(r.retailer_name, d.distributor_name, md.master_distributor_name) AS user_name,
		t.tds_amount
	FROM aeps_tds t
	JOIN aeps_transactions at ON at.aeps_transaction_id = t.aeps_transaction_id
	LEFT JOIN retailers r ON t.user_type = 'RT' AND r.retailer_id = t.user_id
	LEFT JOIN distributors d ON t.user_type = 'DIS' AND d.distributor_id = t.user_id
	LEFT JOIN master_distributors md ON t.user_type = 'MD' AND md.master_distributor_id = t.user_id
`

func (as *PostgresAEPSStore) GetAllAEPSTDSDeductions(p utils.QueryParams) ([]models.AepsTdsDeductionResponse, error) {
	q := aepsTdsSelectAllBase + `
	WHERE t.created_at >= COALESCE($3, '-infinity'::TIMESTAMPTZ)
	AND t.created_at <= COALESCE($4, 'infinity'::TIMESTAMPTZ)
	AND ($5::TEXT IS NULL OR t.user_type = $5)
	AND ($6::TEXT IS NULL OR (
		t.tds_id::TEXT ILIKE '%'||$6||'%' OR
		t.user_id ILIKE '%'||$6||'%' OR
		at.customer_name ILIKE '%'||$6||'%'
	))
	ORDER BY t.created_at DESC
	LIMIT $1 OFFSET $2;
	`
	return scanAepsTdsDeductions(as.db, q, p.Limit, p.Offset, p.StartDate, p.EndDate, p.Status, p.Search)
}

const aepsTdsSelectRetailerBase = `
	SELECT 
		t.tds_id,
		t.aeps_transaction_id,
		t.user_id,
		t.user_type,
		t.created_at,
		t.updated_at,
		at.customer_name,
		r.retailer_name AS user_name,
		t.tds_amount
	FROM aeps_tds t
	JOIN aeps_transactions at ON at.aeps_transaction_id = t.aeps_transaction_id
	JOIN retailers r ON r.retailer_id = t.user_id
`

func (as *PostgresAEPSStore) GetAEPSTDSDeductionsByRetailerID(retailerId string, p utils.QueryParams) ([]models.AepsTdsDeductionResponse, error) {
	q := aepsTdsSelectRetailerBase + `
	WHERE t.user_type = 'RT' AND t.user_id = $7
	AND t.created_at >= COALESCE($3, '-infinity'::TIMESTAMPTZ)
	AND t.created_at <= COALESCE($4, 'infinity'::TIMESTAMPTZ)
	AND ($5::TEXT IS NULL OR t.user_type = $5)
	AND ($6::TEXT IS NULL OR (
		t.tds_id::TEXT ILIKE '%'||$6||'%' OR
		at.customer_name ILIKE '%'||$6||'%'
	))
	ORDER BY t.created_at DESC
	LIMIT $1 OFFSET $2;
	`
	return scanAepsTdsDeductions(as.db, q, p.Limit, p.Offset, p.StartDate, p.EndDate, p.Status, p.Search, retailerId)
}

const aepsTdsSelectDistributorBase = `
	SELECT 
		t.tds_id,
		t.aeps_transaction_id,
		t.user_id,
		t.user_type,
		t.created_at,
		t.updated_at,
		at.customer_name,
		d.distributor_name AS user_name,
		t.tds_amount
	FROM aeps_tds t
	JOIN aeps_transactions at ON at.aeps_transaction_id = t.aeps_transaction_id
	JOIN distributors d ON d.distributor_id = t.user_id
`

func (as *PostgresAEPSStore) GetAEPSTDSDeductionsByDistributorID(distributorId string, p utils.QueryParams) ([]models.AepsTdsDeductionResponse, error) {
	q := aepsTdsSelectDistributorBase + `
	WHERE t.user_type = 'DIS' AND t.user_id = $7
	AND t.created_at >= COALESCE($3, '-infinity'::TIMESTAMPTZ)
	AND t.created_at <= COALESCE($4, 'infinity'::TIMESTAMPTZ)
	AND ($5::TEXT IS NULL OR t.user_type = $5)
	AND ($6::TEXT IS NULL OR (
		t.tds_id::TEXT ILIKE '%'||$6||'%' OR
		at.customer_name ILIKE '%'||$6||'%'
	))
	ORDER BY t.created_at DESC
	LIMIT $1 OFFSET $2;
	`
	return scanAepsTdsDeductions(as.db, q, p.Limit, p.Offset, p.StartDate, p.EndDate, p.Status, p.Search, distributorId)
}

const aepsTdsSelectMdBase = `
	SELECT 
		t.tds_id,
		t.aeps_transaction_id,
		t.user_id,
		t.user_type,
		t.created_at,
		t.updated_at,
		at.customer_name,
		md.master_distributor_name AS user_name,
		t.tds_amount
	FROM aeps_tds t
	JOIN aeps_transactions at ON at.aeps_transaction_id = t.aeps_transaction_id
	JOIN master_distributors md ON md.master_distributor_id = t.user_id
`

func (as *PostgresAEPSStore) GetAEPSTDSDeductionsByMDID(mdId string, p utils.QueryParams) ([]models.AepsTdsDeductionResponse, error) {
	q := aepsTdsSelectMdBase + `
	WHERE t.user_type = 'MD' AND t.user_id = $7
	AND t.created_at >= COALESCE($3, '-infinity'::TIMESTAMPTZ)
	AND t.created_at <= COALESCE($4, 'infinity'::TIMESTAMPTZ)
	AND ($5::TEXT IS NULL OR t.user_type = $5)
	AND ($6::TEXT IS NULL OR (
		t.tds_id::TEXT ILIKE '%'||$6||'%' OR
		at.customer_name ILIKE '%'||$6||'%'
	))
	ORDER BY t.created_at DESC
	LIMIT $1 OFFSET $2;
	`
	return scanAepsTdsDeductions(as.db, q, p.Limit, p.Offset, p.StartDate, p.EndDate, p.Status, p.Search, mdId)
}

func scanAepsTdsDeductions(db *sql.DB, query string, args ...any) ([]models.AepsTdsDeductionResponse, error) {
	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query aeps tds deductions: %w", err)
	}
	defer rows.Close()

	var deductions []models.AepsTdsDeductionResponse

	for rows.Next() {
		var d models.AepsTdsDeductionResponse

		if err := rows.Scan(
			&d.TdsID,
			&d.AepsTransactionID,
			&d.UserID,
			&d.UserType,
			&d.CreatedAt,
			&d.UpdatedAt,
			&d.CustomerName,
			&d.UserName,
			&d.TdsAmount,
		); err != nil {
			return nil, fmt.Errorf("failed to scan aeps tds deduction row: %w", err)
		}

		deductions = append(deductions, d)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating aeps tds deduction rows: %w", err)
	}

	return deductions, nil
}
