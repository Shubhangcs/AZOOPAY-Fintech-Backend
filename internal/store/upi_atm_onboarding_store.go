package store

import (
	"database/sql"

	"github.com/levionstudio/fintech/internal/models"
)

type UPIATMOnboardingStore interface {
	CreateUPIATMApplication(data *models.CreateUPIATMApplicationRequestModel) error
	ChangeUPIATMApplicationStatus(data *models.ChangeUPIATMApplicationStatusRequestModel) error
	GetUPIATMApplicationByRetailerID(retailerId string) (*models.UPIATMApplicationResponseModel, error)
	GetAllUPIATMApplications() ([]models.UPIATMApplicationResponseModel, error)
	CreateUPIATMMerchant(retailerId string, data *models.CreateUPIATMMerchantResponseModel) error
	UpdateUPIATMMerchant(retailerId string, data *models.UpdateUPIATMMerchantResponseModel) error
	GetUPIATMMerchantDetailsByRetailerID(retailerId string) (*models.UPIATMMerchantDetailsResponseModel, error)
	GetAllUPIATMMerchants() ([]models.UPIATMMerchantDetailsResponseModel, error)
	ChangeUPIATMMerchantBlockStatus(retailerId string, blocked bool) error
}

type PostgresUPIATMOnboardingStore struct {
	db *sql.DB
}

func NewPostgresUPIATMOnboardingStore(db *sql.DB) *PostgresUPIATMOnboardingStore {
	return &PostgresUPIATMOnboardingStore{
		db,
	}
}

func (us *PostgresUPIATMOnboardingStore) CreateUPIATMApplication(data *models.CreateUPIATMApplicationRequestModel) error {
	query := `
		INSERT INTO upi_atm_applications(
			retailer_id,
			retailer_remarks,
			latitude,
			longitude
		) VALUES (
			$1 , $2 , $3 , $4
		);
	`

	res, err := us.db.Exec(
		query,
		data.RetailerID,
		data.RetailerRemarks,
		data.Latitude,
		data.Longitude,
	)
	if err != nil {
		return err
	}

	return checkRowsAffected(res)
}

func (us *PostgresUPIATMOnboardingStore) ChangeUPIATMApplicationStatus(data *models.ChangeUPIATMApplicationStatusRequestModel) error {
	query := `
		UPDATE upi_atm_applications
		SET upi_atm_application_status = COALESCE(NULLIF($1, '') , upi_atm_application_status),
			admin_remarks = COALESCE(NULLIF($2 , '') , admin_remarks),
			retailer_remarks = COALESCE(NULLIF($3 , '') , retailer_remarks),
			updated_at = NOW()
		WHERE retailer_id = $4;
	`

	res, err := us.db.Exec(
		query,
		data.UPIATMApplicationStatus,
		data.AdminRemarks,
		data.RetailerRemarks,
		data.RetailerID,
	)
	if err != nil {
		return err
	}

	return checkRowsAffected(res)
}

const upiAtmApplicationSelectBase = `
	SELECT
		a.upi_atm_application_id,
		a.retailer_id,
		a.upi_atm_application_status,
		a.retailer_remarks,
		a.admin_remarks,
		a.latitude,
		a.longitude,
		a.created_at,
		a.updated_at,
		r.retailer_name,
		r.retailer_email,
		r.retailer_phone,
		r.retailer_aadhar_number,
		r.retailer_pan_number,
		r.retailer_address,
		r.retailer_city,
		r.retailer_pincode,
		r.retailer_date_of_birth,
		r.retailer_gender
	FROM upi_atm_applications a
	JOIN retailers r ON r.retailer_id = a.retailer_id
`

func scanUPIATMApplication(row interface{ Scan(dest ...any) error }, res *models.UPIATMApplicationResponseModel) error {
	return row.Scan(
		&res.UPIATMApplicationID,
		&res.RetailerID,
		&res.UPIATMApplicationStatus,
		&res.RetailerRemarks,
		&res.AdminRemarks,
		&res.Latitude,
		&res.Longitude,
		&res.CreatedAT,
		&res.UpdatedAT,
		&res.RetailerDetails.RetailerName,
		&res.RetailerDetails.RetailerEmail,
		&res.RetailerDetails.RetailerPhone,
		&res.RetailerDetails.RetailerAadhaarNumber,
		&res.RetailerDetails.RetailerPanNumber,
		&res.RetailerDetails.RetailerFullAddress,
		&res.RetailerDetails.RetailerCity,
		&res.RetailerDetails.RetailerPincode,
		&res.RetailerDetails.RetailerDateOfBirth,
		&res.RetailerDetails.RetailerGender,
	)
}

func (us *PostgresUPIATMOnboardingStore) GetUPIATMApplicationByRetailerID(retailerId string) (*models.UPIATMApplicationResponseModel, error) {
	query := upiAtmApplicationSelectBase + `
		WHERE a.retailer_id = $1;
	`

	var res models.UPIATMApplicationResponseModel
	if err := scanUPIATMApplication(us.db.QueryRow(query, retailerId), &res); err != nil {
		return nil, err
	}

	return &res, nil
}

func (us *PostgresUPIATMOnboardingStore) GetAllUPIATMApplications() ([]models.UPIATMApplicationResponseModel, error) {
	query := upiAtmApplicationSelectBase + `
		ORDER BY a.created_at DESC;
	`

	dbres, err := us.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer dbres.Close()

	var resList []models.UPIATMApplicationResponseModel
	for dbres.Next() {
		var res models.UPIATMApplicationResponseModel
		if err := scanUPIATMApplication(dbres, &res); err != nil {
			return nil, err
		}

		resList = append(resList, res)
	}

	if err := dbres.Err(); err != nil {
		return nil, err
	}

	return resList, nil
}

func (us *PostgresUPIATMOnboardingStore) CreateUPIATMMerchant(retailerId string, data *models.CreateUPIATMMerchantResponseModel) error {
	query := `
		INSERT INTO upi_atm_merchant_details(
			retailer_id,
			sub_merchant_id,
			parent_merchant_id,
			outlet_id,
			min_kyc_status,
			ekyc_status,
			mobile_change_state,
			ipay_uuid
		) VALUES (
			$1 , $2 , $3 , $4 , $5 , $6 , $7 , $8
		);
	`

	res, err := us.db.Exec(
		query,
		retailerId,
		data.SubMerchantID,
		data.ParentMerchantID,
		data.OutletID,
		data.MinKYCStatus,
		data.EKYCStatus,
		data.MobileChangeState,
		data.IPayUUID,
	)
	if err != nil {
		return err
	}

	return checkRowsAffected(res)
}

func (us *PostgresUPIATMOnboardingStore) UpdateUPIATMMerchant(retailerId string, data *models.UpdateUPIATMMerchantResponseModel) error {
	query := `
		UPDATE upi_atm_merchant_details
		SET min_kyc_status = COALESCE(NULLIF($1, ''), min_kyc_status),
			ekyc_status = COALESCE(NULLIF($2, ''), ekyc_status),
			ekyc_action = COALESCE(NULLIF($3, ''), ekyc_action),
			reference_key = COALESCE(NULLIF($4, ''), reference_key),
			status = COALESCE(NULLIF($5, ''), status),
			is_face_auth_available = COALESCE($6, is_face_auth_available),
			is_biometric_kyc_mandatory = COALESCE($7, is_biometric_kyc_mandatory),
			bank_name = COALESCE(NULLIF($8, ''), bank_name),
			updated_at = NOW()
		WHERE retailer_id = $9;
	`

	res, err := us.db.Exec(
		query,
		data.MinKYCStatus,
		data.EKYCStatus,
		data.EKYCAction,
		data.ReferenceKey,
		data.Data.Status,
		data.Data.IsFaceAuthAvailable,
		data.Data.IsBiometricKycMandatory,
		data.Data.BankName,
		retailerId,
	)
	if err != nil {
		return err
	}

	return checkRowsAffected(res)
}

const upiAtmMerchantSelectBase = `
	SELECT
		upi_atm_merchant_id,
		retailer_id,
		sub_merchant_id,
		parent_merchant_id,
		outlet_id,
		min_kyc_status,
		ekyc_status,
		ekyc_action,
		reference_key,
		status,
		is_face_auth_available,
		is_biometric_kyc_mandatory,
		bank_name,
		mobile_change_state,
		ipay_uuid,
		is_merchant_blocked,
		created_at,
		updated_at
	FROM upi_atm_merchant_details
`

func scanUPIATMMerchant(row interface{ Scan(dest ...any) error }, res *models.UPIATMMerchantDetailsResponseModel) error {
	return row.Scan(
		&res.UPIATMMerchantID,
		&res.RetailerID,
		&res.SubMerchantID,
		&res.ParentMerchantID,
		&res.OutletID,
		&res.MinKYCStatus,
		&res.EKYCStatus,
		&res.EKYCAction,
		&res.ReferenceKey,
		&res.Status,
		&res.IsFaceAuthAvailable,
		&res.IsBiometricKycMandatory,
		&res.BankName,
		&res.MobileChangeState,
		&res.IPayUUID,
		&res.IsMerchantBlocked,
		&res.CreatedAT,
		&res.UpdatedAT,
	)
}

func (us *PostgresUPIATMOnboardingStore) GetUPIATMMerchantDetailsByRetailerID(retailerId string) (*models.UPIATMMerchantDetailsResponseModel, error) {
	query := upiAtmMerchantSelectBase + `
		WHERE retailer_id = $1;
	`

	var res models.UPIATMMerchantDetailsResponseModel
	if err := scanUPIATMMerchant(us.db.QueryRow(query, retailerId), &res); err != nil {
		return nil, err
	}

	return &res, nil
}

func (us *PostgresUPIATMOnboardingStore) GetAllUPIATMMerchants() ([]models.UPIATMMerchantDetailsResponseModel, error) {
	query := upiAtmMerchantSelectBase + `
		ORDER BY created_at DESC;
	`

	dbres, err := us.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer dbres.Close()

	var resList []models.UPIATMMerchantDetailsResponseModel
	for dbres.Next() {
		var res models.UPIATMMerchantDetailsResponseModel
		if err := scanUPIATMMerchant(dbres, &res); err != nil {
			return nil, err
		}

		resList = append(resList, res)
	}

	if err := dbres.Err(); err != nil {
		return nil, err
	}

	return resList, nil
}

func (us *PostgresUPIATMOnboardingStore) ChangeUPIATMMerchantBlockStatus(retailerId string, blocked bool) error {
	query := `
		UPDATE upi_atm_merchant_details
		SET is_merchant_blocked = $1,
			updated_at = NOW()
		WHERE retailer_id = $2;
	`

	res, err := us.db.Exec(query, blocked, retailerId)
	if err != nil {
		return err
	}

	return checkRowsAffected(res)
}
