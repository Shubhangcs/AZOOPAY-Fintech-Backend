package store

import (
	"database/sql"
	"fmt"

	"github.com/levionstudio/fintech/internal/models"
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
