package store

import (
	"database/sql"

	"github.com/levionstudio/fintech/internal/models"
)

type CreditCardPaymentStore interface {
}

type PostgresCreditCardPaymentStore struct {
	db *sql.DB
}

func NewPostgresCreditCardPaymentStore(db *sql.DB) *PostgresCreditCardPaymentStore {
	return &PostgresCreditCardPaymentStore{
		db,
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

func (cc *PostgresCreditCardPaymentStore) GetBeneficiaryByBeneficiaryID(beneficiaryId string) (*models.GetCreditCardBeneficiaryDetailsResponseModel, error) {
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

func (cc *PostgresCreditCardPaymentStore) CreateCreditCardPaymentTransaction() {
	
}
