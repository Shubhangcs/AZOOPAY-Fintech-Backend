package store

import (
	"database/sql"

	"github.com/levionstudio/fintech/internal/models"
)

type AEPSOnboardingStore interface {
	CreateAEPSApplication(data *models.CreateAEPSApplicationRequestModel) error
	ChangeAEPSApplicationStatus(data *models.ChangeAEPSApplicationStatusRequestModel) error
	GetAEPSApplicationByRetailerID(retailerId string) (*models.AEPSApplicationResponseModel, error)
	GetAllAEPSApplications() ([]models.AEPSApplicationResponseModel, error)
	CreateAEPSMerchant(retailerId string, data *models.CreateAEPSMerchantResponseModel) error
	UpdateAEPSMerchant(retailerId string, data *models.UpdateAEPSMerchantResponseModel) error
	GetAEPSMerchantDetailsByRetailerID(retailerId string) (*models.AEPSMerchantDetailsResponseModel, error)
}

type PostgresAEPSOnboardingStore struct {
	db *sql.DB
}

func NewPostgresAEPSOnboardingStore(db *sql.DB) *PostgresAEPSOnboardingStore {
	return &PostgresAEPSOnboardingStore{
		db,
	}
}

func (pa *PostgresAEPSOnboardingStore) CreateAEPSApplication(data *models.CreateAEPSApplicationRequestModel) error {
	query := `
		INSERT INTO aeps_applications(
			retailer_id,
			retailer_remarks,
			latitude,
			longitude
		) VALUES (
			$1 , $2 , $3 , $4
		);
	`

	res, err := pa.db.Exec(
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

func (pa *PostgresAEPSOnboardingStore) ChangeAEPSApplicationStatus(data *models.ChangeAEPSApplicationStatusRequestModel) error {
	query := `
		UPDATE aeps_applications
		SET aeps_application_status = COALESCE(NULLIF($1, '') , aeps_application_status),
			admin_remarks = COALESCE(NULLIF($2 , '') , admin_remarks),
			retailer_remarks = COALESCE(NULLIF($3 , '') , retailer_remarks),
			updated_at = NOW()
		WHERE retailer_id = $4;
	`

	res, err := pa.db.Exec(
		query,
		data.AEPSApplicationStatus,
		data.AdminRemarks,
		data.RetailerRemarks,
		data.RetailerID,
	)
	if err != nil {
		return err
	}

	return checkRowsAffected(res)
}

func (pa *PostgresAEPSOnboardingStore) GetAEPSApplicationByRetailerID(retailerId string) (*models.AEPSApplicationResponseModel, error) {
	query := `
		SELECT
			a.aeps_application_id,
			a.retailer_id,
			a.aeps_application_status,
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
		FROM aeps_applications a
		JOIN retailers r ON r.retailer_id = a.retailer_id
		WHERE a.retailer_id = $1;
	`

	var res models.AEPSApplicationResponseModel
	if err := pa.db.QueryRow(
		query,
		retailerId,
	).Scan(
		&res.AEPSApplicationID,
		&res.RetailerID,
		&res.AEPSApplicationStatus,
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
	); err != nil {
		return nil, err
	}

	return &res, nil
}

func (pa *PostgresAEPSOnboardingStore) GetAllAEPSApplications() ([]models.AEPSApplicationResponseModel, error) {
	query := `
		SELECT
			a.aeps_application_id,
			a.retailer_id,
			a.aeps_application_status,
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
		FROM aeps_applications a
		JOIN retailers r ON r.retailer_id = a.retailer_id;
	`

	dbres, err := pa.db.Query(query)
	if err != nil {
		return nil, err
	}

	var res models.AEPSApplicationResponseModel
	var resList []models.AEPSApplicationResponseModel
	for dbres.Next() {
		if err := dbres.Scan(
			&res.AEPSApplicationID,
			&res.RetailerID,
			&res.AEPSApplicationStatus,
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
		); err != nil {
			return nil, err
		}

		resList = append(resList, res)
	}

	if dbres.Err() != nil {
		return nil, dbres.Err()
	}

	return resList, nil
}

func (pa *PostgresAEPSOnboardingStore) CreateAEPSMerchant(retailerId string, data *models.CreateAEPSMerchantResponseModel) error {
	query := `
		INSERT INTO aeps_merchant_details(
			retailer_id,
			sub_merchant_id,
			parent_merchant_id,
			outlet_id,
			min_kyc_status,
			ekyc_status,
			mobile_charge_state,
			i_pay_uuid
		) VALUES (
			$1 , $2 , $3 , $4 , $5 , $6 , $7 , $8 
		);
	`

	res, err := pa.db.Exec(
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

func (pa *PostgresAEPSOnboardingStore) UpdateAEPSMerchant(retailerId string, data *models.UpdateAEPSMerchantResponseModel) error {
	query := `
		UPDATE aeps_merchant_details
		SET ekyc_status = COALESCE(NULLIF($1, ''), ekyc_status),
    		ekyc_action = COALESCE(NULLIF($2, ''), ekyc_action),
    		reference_key = COALESCE(NULLIF($3, ''), reference_key),
    		status = COALESCE(NULLIF($4, ''), status),
    		is_face_auth_available = COALESCE($5, is_face_auth_available),
    		is_biometric_kyc_manditory = COALESCE($6, is_biometric_kyc_manditory),
    		bank_name = COALESCE(NULLIF($7, ''), bank_name)
		WHERE retailer_id = $8;
	`

	res, err := pa.db.Exec(
		query,
		data.EKYCStatus,
		data.EKYCAction,
		data.ReferenceKey,
		data.Data.Status,
		data.Data.IsFaceAuthAvailable,
		data.Data.IsBiometricKycManditory,
		data.Data.BankName,
		retailerId,
	)

	if err != nil {
		return err
	}

	return checkRowsAffected(res)
}

func (pa *PostgresAEPSOnboardingStore) GetAEPSMerchantDetailsByRetailerID(retailerId string) (*models.AEPSMerchantDetailsResponseModel, error) {
	query := `
		SELECT
			aeps_merchant_id,
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
			is_biometric_kyc_manditory,
			bank_name,
			is_merchant_blocked
		FROM aeps_merchant_details
		WHERE retailer_id=$1;
	`

	var res models.AEPSMerchantDetailsResponseModel
	if err := pa.db.QueryRow(
		query,
		retailerId,
	).Scan(
		&res.AEPSMerchantID,
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
		&res.IsBiometricKycManditory,
		&res.BankName,
		&res.IsMerchantBlocked,
	); err != nil {
		return nil, err
	}

	return &res, nil
}
