package store

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/levionstudio/fintech/internal/models"
	"github.com/levionstudio/fintech/internal/utils"
)

type CreditCardPaymentStore interface {
	CreateCreditCardBeneficiary(data *models.CreateCreditCardBeneficiaryRequestModel) error
	UpdateCreditCardBeneficiary(data *models.UpdateCreditCardBeneficiaryRequestModel) error
	DeleteCreditCardBeneficiary(beneficiaryId int64) error
	GetBeneficiariesByRetailerID(retailerId string) ([]models.GetCreditCardBeneficiaryDetailsResponseModel, error)
	GetBeneficiaryByBeneficiaryID(beneficiaryId int64) (*models.GetCreditCardBeneficiaryDetailsResponseModel, error)
	InitilizeCreateCreditCardPaymentTransaction(data *models.CreateCreditCardPaymentTransactionRequestModel) (int64, error)
	FinalizeCreateCreditCardPaymentTransaction(
		transactionID int64,
		res *models.CreditCardBillPaymentAPIResponse,
	) error
	GetCreditCardPaymentTransactionByID(transactionID int64) (*models.CreditCardPaymentTransactionModel, error)
	GetAllCreditCardPaymentTransactions(p utils.QueryParams) ([]models.CreditCardPaymentTransactionModel, error)
	GetCreditCardPaymentTransactionsByRetailerID(retailerID string, p utils.QueryParams) ([]models.CreditCardPaymentTransactionModel, error)
	GetCreditCardPaymentTransactionsByDistributorID(distributorID string, p utils.QueryParams) ([]models.CreditCardPaymentTransactionModel, error)
	GetCreditCardPaymentTransactionsByMasterDistributorID(mdID string, p utils.QueryParams) ([]models.CreditCardPaymentTransactionModel, error)
	RefundCreditCardPaymentTransaction(transactionID int64) error
}

type PostgresCreditCardPaymentStore struct {
	db  *sql.DB
	wts *PostgresWalletTransactionStore
}

func NewPostgresCreditCardPaymentStore(db *sql.DB, wts *PostgresWalletTransactionStore) *PostgresCreditCardPaymentStore {
	return &PostgresCreditCardPaymentStore{
		db,
		wts,
	}
}

func (cc *PostgresCreditCardPaymentStore) CreateCreditCardBeneficiary(data *models.CreateCreditCardBeneficiaryRequestModel) error {
	query := `
		INSERT INTO credit_card_beneficiaries(
			retailer_id,
			retailer_name,
			beneficiary_name,
			beneficiary_phone,
			beneficiary_bank_name,
			beneficiary_account_number,
			beneficiary_ifsc_code,
			operator_name,
			operator_code
		) VALUES (
			$1 , $2 , $3 , $4 , $5 , $6 , $7 , $8 , $9 
		);
	`

	res, err := cc.db.Exec(
		query,
		data.RetailerID,
		data.RetailerName,
		data.BeneficiaryName,
		data.PhoneNumber,
		data.BankName,
		data.AccountNumber,
		data.IFSCCode,
		data.OperatorName,
		data.OperatorCode,
	)

	if err != nil {
		return err
	}

	return checkRowsAffected(res)
}

func (cc *PostgresCreditCardPaymentStore) UpdateCreditCardBeneficiary(data *models.UpdateCreditCardBeneficiaryRequestModel) error {
	query := `
		UPDATE credit_card_beneficiaries
		SET beneficiary_name = COALESCE(NULLIF($1 , '') , beneficiary_name),
			beneficiary_phone = COALESCE(NULLIF($2 , '') , beneficiary_phone),
			beneficiary_account_number = COALESCE(NULLIF($3 , '') , beneficiary_account_number),
			beneficiary_ifsc_code = COALESCE(NULLIF($4 , '') , beneficiary_ifsc_code),
			beneficiary_bank_name = COALESCE(NULLIF($5 , '') , beneficiary_bank_name),
			operator_name = COALESCE(NULLIF($6 , '') , operator_name),
			operator_code = COALESCE(NULLIF($7 , '') , operator_code)
		WHERE beneficiary_id = $8;
	`

	res, err := cc.db.Exec(
		query,
		data.BeneficiaryName,
		data.PhoneNumber,
		data.AccountNumber,
		data.IFSCCode,
		data.BankName,
		data.OperatorName,
		data.OperatorCode,
		data.BeneficiaryID,
	)

	if err != nil {
		return err
	}

	return checkRowsAffected(res)
}

func (cc *PostgresCreditCardPaymentStore) DeleteCreditCardBeneficiary(beneficiaryId int64) error {
	query := `
		DELETE FROM credit_card_beneficiaries WHERE beneficiary_id=$1;
	`

	res, err := cc.db.Exec(
		query,
		beneficiaryId,
	)

	if err != nil {
		return err
	}

	return checkRowsAffected(res)
}

func (cc *PostgresCreditCardPaymentStore) GetBeneficiariesByRetailerID(retailerId string) ([]models.GetCreditCardBeneficiaryDetailsResponseModel, error) {
	query := `
		SELECT
			beneficiary_id,
			retailer_id,
			retailer_name,
			beneficiary_name,
			beneficiary_phone,
			beneficiary_account_number,
			beneficiary_ifsc_code,
			beneficiary_bank_name,
			operator_name,
			operator_code,
			created_at,
			updated_at
		FROM credit_card_beneficiaries
		WHERE retailer_id = $1;
	`

	res, err := cc.db.Query(query, retailerId)
	if err != nil {
		return nil, err
	}
	defer res.Close()

	var bene models.GetCreditCardBeneficiaryDetailsResponseModel
	var benes []models.GetCreditCardBeneficiaryDetailsResponseModel
	for res.Next() {
		if err := res.Scan(
			&bene.BeneficiaryID,
			&bene.RetailerID,
			&bene.RetailerName,
			&bene.BeneficiaryName,
			&bene.PhoneNumber,
			&bene.AccountNumber,
			&bene.IFSCCode,
			&bene.BankName,
			&bene.OperatorName,
			&bene.OperatorCode,
			&bene.CreatedAT,
			&bene.UpdatedAT,
		); err != nil {
			return nil, err
		}

		benes = append(benes, bene)
	}

	if res.Err() != nil {
		return nil, res.Err()
	}

	return benes, nil
}

func (cc *PostgresCreditCardPaymentStore) GetBeneficiaryByBeneficiaryID(beneficiaryId int64) (*models.GetCreditCardBeneficiaryDetailsResponseModel, error) {
	query := `
		SELECT
			beneficiary_id,
			retailer_id,
			retailer_name,
			beneficiary_name,
			beneficiary_phone,
			beneficiary_account_number,
			beneficiary_ifsc_code,
			beneficiary_bank_name,
			operator_name,
			operator_code,
			created_at,
			updated_at
		FROM credit_card_beneficiaries
		WHERE beneficiary_id = $1;
	`

	var bene models.GetCreditCardBeneficiaryDetailsResponseModel
	if err := cc.db.QueryRow(
		query,
		beneficiaryId,
	).Scan(
		&bene.BeneficiaryID,
		&bene.RetailerID,
		&bene.RetailerName,
		&bene.BeneficiaryName,
		&bene.PhoneNumber,
		&bene.AccountNumber,
		&bene.IFSCCode,
		&bene.BankName,
		&bene.OperatorName,
		&bene.OperatorCode,
		&bene.CreatedAT,
		&bene.UpdatedAT,
	); err != nil {
		return nil, err
	}

	return &bene, nil
}

func (cc *PostgresCreditCardPaymentStore) InitilizeCreateCreditCardPaymentTransaction(data *models.CreateCreditCardPaymentTransactionRequestModel) (int64, error) {
	if err := verifyMpin(cc.db, data.BeneDetails.RetailerID, data.Mpin); err != nil {
		return 0, err
	}

	rc, err := getRetailerDetails(cc.db, data.BeneDetails.RetailerID)
	if err != nil {
		return 0, err
	}
	if !rc.kyc {
		return 0, errors.New("retailer KYC is not verified")
	}
	if rc.blocked {
		return 0, errors.New("retailer is blocked")
	}
	if rc.balance < data.Amount {
		return 0, errors.New("insufficient wallet balance")
	}

	tx, err := cc.db.Begin()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	query := `
		INSERT INTO credit_card_payment_transactions(
			retailer_id,
			beneficiary_id,
			amount,
			status,
			partner_request_id
		) VALUES (
			$1 , $2 , $3 , $4 , $5
		) RETURNING transaction_id;
	`

	var ccTransactionId int64
	if err := tx.QueryRow(
		query,
		data.BeneDetails.RetailerID,
		data.BeneDetails.BeneficiaryID,
		data.Amount,
		"PENDING",
		data.PartnerRequestID,
	).Scan(&ccTransactionId); err != nil {
		return 0, err
	}

	userTableInfo, err := getUserTableInfo(data.BeneDetails.RetailerID)
	if err != nil {
		return 0, err
	}

	if err := debitTx(tx, transaction{
		UserID:        data.BeneDetails.RetailerID,
		ReferenceID:   fmt.Sprintf("%d", ccTransactionId),
		Amount:        data.Amount,
		Reason:        "CC_BILL_PAYMENT",
		Remarks:       fmt.Sprintf("Credit Card Bill Payment By Retailer: %s For Beneficiary: %s", data.BeneDetails.RetailerID, data.BeneDetails.BeneficiaryName),
		userTableInfo: *userTableInfo,
	}, cc.wts); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, checkExistsTx(tx, userTableInfo.TableName, userTableInfo.IDColumnName, data.BeneDetails.RetailerID, "retailer")
		}
		return 0, err
	}

	return ccTransactionId, tx.Commit()
}

func (cc *PostgresCreditCardPaymentStore) FinalizeCreateCreditCardPaymentTransaction(
	transactionID int64,
	res *models.CreditCardBillPaymentAPIResponse,
) error {

	var status string

	switch res.Status {
	case 1:
		status = "SUCCESS"
	case 2:
		status = "PENDING"
	case 3:
		status = "FAILED"
	default:
		return fmt.Errorf("invalid payment status: %d", res.Status)
	}

	query := `
		UPDATE credit_card_payment_transactions
		SET
			status = $1,
			order_id = $2,
			updated_at = NOW()
		WHERE transaction_id = $3
	`

	args := []any{
		status,
		res.OrderID,
		transactionID,
	}

	if res.Status == 1 {
		query = `
			UPDATE credit_card_payment_transactions
			SET
				status = $1,
				order_id = $2,
				operator_transaction_id = $3,
				updated_at = NOW()
			WHERE transaction_id = $4
		`

		args = []any{
			status,
			res.OrderID,
			res.OperatorTransactionID,
			transactionID,
		}
	}

	dbres, err := cc.db.Exec(query, args...)
	if err != nil {
		return err
	}

	return checkRowsAffected(dbres)
}

const ccPaymentSelectBase = `
SELECT
	t.transaction_id,
	t.retailer_id,
	COALESCE(r.retailer_name, '') AS retailer_name,
	r.retailer_business_name,
	t.beneficiary_id,
	b.beneficiary_name,
	b.beneficiary_phone,
	b.beneficiary_account_number,
	b.beneficiary_ifsc_code,
	b.beneficiary_bank_name,
	b.operator_name,
	b.operator_code,
	t.amount,
	t.status,
	t.partner_request_id,
	t.order_id,
	t.operator_transaction_id,
	COALESCE(wt.before_balance, 0) AS before_balance,
	COALESCE(wt.after_balance, 0) AS after_balance,
	t.created_at,
	t.updated_at
FROM credit_card_payment_transactions t
JOIN retailers r ON r.retailer_id = t.retailer_id
JOIN credit_card_beneficiaries b ON b.beneficiary_id = t.beneficiary_id
LEFT JOIN wallet_transactions wt ON wt.reference_id = t.transaction_id::TEXT
	AND wt.user_id = t.retailer_id AND wt.debit_amount IS NOT NULL
	AND wt.transaction_reason = 'CC_BILL_PAYMENT'
`

const ccPaymentFilters = `
	AND t.created_at >= COALESCE($4, '-infinity'::TIMESTAMPTZ)
	AND t.created_at <= COALESCE($5, 'infinity'::TIMESTAMPTZ)
	AND ($6::TEXT IS NULL OR t.status = $6)
	AND ($7::TEXT IS NULL OR (
		t.transaction_id::TEXT ILIKE '%'||$7||'%' OR
		t.partner_request_id ILIKE '%'||$7||'%' OR
		t.order_id ILIKE '%'||$7||'%' OR
		t.operator_transaction_id ILIKE '%'||$7||'%' OR
		b.beneficiary_name ILIKE '%'||$7||'%' OR
		b.beneficiary_phone ILIKE '%'||$7||'%'
	))
	ORDER BY t.created_at DESC
	LIMIT $2 OFFSET $3;
`

func (cc *PostgresCreditCardPaymentStore) GetCreditCardPaymentTransactionByID(transactionID int64) (*models.CreditCardPaymentTransactionModel, error) {
	q := ccPaymentSelectBase + `WHERE t.transaction_id = $1;`
	results, err := scanCreditCardPaymentTransactions(cc.db, q, transactionID)
	if err != nil {
		return nil, err
	}
	if len(results) == 0 {
		return nil, errors.New("cc transaction not found")
	}
	return &results[0], nil
}

func (cc *PostgresCreditCardPaymentStore) GetAllCreditCardPaymentTransactions(p utils.QueryParams) ([]models.CreditCardPaymentTransactionModel, error) {
	// $1 is a no-op placeholder so every list query shares the same filter numbering.
	q := ccPaymentSelectBase + `WHERE ($1::TEXT IS NULL OR TRUE)` + ccPaymentFilters
	return scanCreditCardPaymentTransactions(cc.db, q, nil, p.Limit, p.Offset, p.StartDate, p.EndDate, p.Status, p.Search)
}

func (cc *PostgresCreditCardPaymentStore) GetCreditCardPaymentTransactionsByRetailerID(retailerID string, p utils.QueryParams) ([]models.CreditCardPaymentTransactionModel, error) {
	q := ccPaymentSelectBase + `WHERE t.retailer_id = $1` + ccPaymentFilters
	return scanCreditCardPaymentTransactions(cc.db, q, retailerID, p.Limit, p.Offset, p.StartDate, p.EndDate, p.Status, p.Search)
}

func (cc *PostgresCreditCardPaymentStore) GetCreditCardPaymentTransactionsByDistributorID(distributorID string, p utils.QueryParams) ([]models.CreditCardPaymentTransactionModel, error) {
	q := ccPaymentSelectBase + `WHERE r.distributor_id = $1` + ccPaymentFilters
	return scanCreditCardPaymentTransactions(cc.db, q, distributorID, p.Limit, p.Offset, p.StartDate, p.EndDate, p.Status, p.Search)
}

func (cc *PostgresCreditCardPaymentStore) GetCreditCardPaymentTransactionsByMasterDistributorID(mdID string, p utils.QueryParams) ([]models.CreditCardPaymentTransactionModel, error) {
	q := ccPaymentSelectBase + `
	JOIN distributors d ON d.distributor_id = r.distributor_id
	WHERE d.master_distributor_id = $1` + ccPaymentFilters
	return scanCreditCardPaymentTransactions(cc.db, q, mdID, p.Limit, p.Offset, p.StartDate, p.EndDate, p.Status, p.Search)
}

func scanCreditCardPaymentTransactions(db *sql.DB, q string, args ...any) ([]models.CreditCardPaymentTransactionModel, error) {
	rows, err := db.Query(q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []models.CreditCardPaymentTransactionModel
	for rows.Next() {
		var t models.CreditCardPaymentTransactionModel
		if err := rows.Scan(
			&t.TransactionID,
			&t.RetailerID,
			&t.RetailerName,
			&t.RetailerBusinessName,
			&t.BeneficiaryID,
			&t.BeneficiaryName,
			&t.PhoneNumber,
			&t.AccountNumber,
			&t.IFSCCode,
			&t.BankName,
			&t.OperatorName,
			&t.OperatorCode,
			&t.Amount,
			&t.Status,
			&t.PartnerRequestID,
			&t.OrderID,
			&t.OperatorTransactionID,
			&t.BeforeBalance,
			&t.AfterBalance,
			&t.CreatedAT,
			&t.UpdatedAT,
		); err != nil {
			return nil, err
		}
		results = append(results, t)
	}
	return results, rows.Err()
}

// RefundCreditCardPaymentTransaction refunds a FAILED payment to the retailer's
// refund_wallet (same as DTH / mobile recharge / electricity refunds) and marks
// it REFUNDED. The status guard prevents a double refund.
func (cc *PostgresCreditCardPaymentStore) RefundCreditCardPaymentTransaction(transactionID int64) error {
	t, err := cc.GetCreditCardPaymentTransactionByID(transactionID)
	if err != nil {
		return err
	}
	if t.Status != "FAILED" {
		return errors.New("only FAILED cc transactions can be refunded")
	}

	tx, err := cc.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	res, err := tx.Exec(`
		UPDATE credit_card_payment_transactions
		SET status = 'REFUNDED',
			updated_at = NOW()
		WHERE transaction_id = $1 AND status = 'FAILED'
	`, transactionID)
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return errors.New("cc transaction not found or already refunded")
	}

	refID := fmt.Sprintf("%d", transactionID)
	remarks := fmt.Sprintf("Credit Card Bill Payment refund | Ref: %s", refID)

	var refundAfterBalance float64
	if err := tx.QueryRow(`
		UPDATE retailers
		SET refund_wallet = refund_wallet + $1,
			updated_at = NOW()
		WHERE retailer_id = $2
		RETURNING refund_wallet;
	`, t.Amount, t.RetailerID).Scan(&refundAfterBalance); err != nil {
		return err
	}

	if err := cc.wts.CreateWalletTransactionTx(tx, &models.WalletTransactionModel{
		UserID:            t.RetailerID,
		ReferenceID:       refID,
		CreditAmount:      &t.Amount,
		BeforeBalance:     refundAfterBalance - t.Amount,
		AfterBalance:      refundAfterBalance,
		TransactionReason: "CC_BILL_PAYMENT_REFUND",
		Remarks:           remarks,
	}); err != nil {
		return err
	}

	return tx.Commit()
}
