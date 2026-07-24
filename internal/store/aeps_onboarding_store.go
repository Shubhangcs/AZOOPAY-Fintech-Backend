package store

import (
	"database/sql"

	"github.com/levionstudio/fintech/internal/models"
)

type AEPSOnboardingStore interface {
	ApplyForAEPS(data *models.ApplyForAEPSRequestModel) error
	ChangeAEPSApplicationStatus(data *models.ApplyForAEPSRequestModel) error
	GetAEPSApplications() ([]models.AEPSApplicationResponseModel, error)
	GetAEPSApplication(userId string) (*models.AEPSApplicationResponseModel, error)
	CheckAEPSApplicationStatus(userId string) (*models.AEPSApplicationResponseModel, error)
	GetUserDetailsForAEPSSignup(userId string) (*models.AEPSOnboardingSubMerchantSignupRequestModel, error)
	SubMerchantAEPSSignup(userId string, data *models.AEPSOnboardingSubMerchantSignupSuccessResponseModel) error
	GetSubMerchantIDForAEPSEKYCCheck(userId string) (string, error)
	MerchantAEPSEKYCCheck(userId string, data *models.AEPSOnboardingeKycCheckResponseModel) error
	GetMerchantDetailsForBiometricKYC(userId string, req *models.AEPSOnboardingBiometricKYCRequestModel) error
}

type PostgresAEPSOnboardingStore struct {
	db *sql.DB
}

func NewAEPSOnboardingStore(db *sql.DB) *PostgresAEPSOnboardingStore {
	return &PostgresAEPSOnboardingStore{
		db,
	}
}

func (pa *PostgresAEPSOnboardingStore) ApplyForAEPS(data *models.ApplyForAEPSRequestModel) error {
	query := `
		INSERT INTO aeps_applications(
			retailer_id,
			retailer_name,
			aeps_application_status,
			remarks,
			latitude,
			longitude
		) VALUES(
			$1, $2, $3, $4, $5, $6
		);
	`

	res, err := pa.db.Exec(
		query,
		data.RetailerID,
		data.RetailerName,
		"PENDING",
		"Applied to AEPS",
		data.Latitude,
		data.Longitude,
	)
	if err != nil {
		return err
	}

	return checkRowsAffected(res)
}

func (pa *PostgresAEPSOnboardingStore) ChangeAEPSApplicationStatus(data *models.ApplyForAEPSRequestModel) error {
	query := `
		UPDATE aeps_applications
		SET aeps_application_status = $1,
			remarks = $2,
			updated_at = NOW()
		WHERE aeps_application_id = $3;
	`

	res, err := pa.db.Exec(
		query,
		data.AEPSApplicationStatus,
		data.Remarks,
		data.AEPSApplicationID,
	)

	if err != nil {
		return err
	}

	return checkRowsAffected(res)
}

func (pa *PostgresAEPSOnboardingStore) GetAEPSApplications() ([]models.AEPSApplicationResponseModel, error) {
	query := `
		SELECT 
			a.aeps_application_id,
			a.retailer_id,
			a.retailer_name,
			a.aeps_application_status,
			a.remarks,
			a.latitude,
			a.longitude,
			a.created_at,
			a.updated_at,
			u.retailer_aadhar_number,
			u.retailer_pan_number
			u.retailer_gender,
			u.retailer_email,
			u.retailer_address,
			u.retailer_city,
			u.retailer_pincode,
			u.retailer_date_of_birth,
			u.retailer_phone
		FROM aeps_applications a
		JOIN retailers u
			ON a.retailer_id = u.retailer_id;
	`

	res, err := pa.db.Query(query)
	if err != nil {
		return nil, err
	}

	var application models.AEPSApplicationResponseModel
	var applications []models.AEPSApplicationResponseModel
	for res.Next() {
		if err := res.Scan(
			&application.AEPSApplicationID,
			&application.RetailerID,
			&application.RetailerName,
			&application.AEPSApplicationStatus,
			&application.Remarks,
			&application.Latitude,
			&application.Longitude,
			&application.CreatedAT,
			&application.UpdatedAT,
			&application.RetailerAadhaar,
			&application.RetailerPAN,
			&application.RetailerGender,
			&application.RetailerEmail,
			&application.RetailerAddress,
			&application.RetailerCity,
			&application.RetailerPincode,
			&application.RetailerDateOfBirth,
			&application.RetailerPhone,
		); err != nil {
			return nil, err
		}

		applications = append(applications, application)
	}

	if res.Err() != nil {
		return nil, res.Err()
	}

	return applications, nil
}

func (pa *PostgresAEPSOnboardingStore) GetAEPSApplication(userId string) (*models.AEPSApplicationResponseModel, error) {
	query := `
		SELECT 
			a.aeps_application_id,
			a.retailer_id,
			a.retailer_name,
			a.aeps_application_status,
			a.remarks,
			a.latitude,
			a.longitude,
			a.created_at,
			a.updated_at,
			u.retailer_aadhar_number,
			u.retailer_pan_number,
			u.retailer_gender,
			u.retailer_email,
			u.retailer_address,
			u.retailer_city,
			u.retailer_pincode,
			u.retailer_date_of_birth,
			u.retailer_phone
		FROM aeps_applications a
		JOIN retailers u
			ON a.retailer_id = u.retailer_id
		WHERE a.retailer_id = $1;
	`

	var application models.AEPSApplicationResponseModel
	if err := pa.db.QueryRow(
		query,
		userId,
	).Scan(
		&application.AEPSApplicationID,
		&application.RetailerID,
		&application.RetailerName,
		&application.AEPSApplicationStatus,
		&application.Remarks,
		&application.Latitude,
		&application.Longitude,
		&application.CreatedAT,
		&application.UpdatedAT,
		&application.RetailerAadhaar,
		&application.RetailerPAN,
		&application.RetailerGender,
		&application.RetailerEmail,
		&application.RetailerAddress,
		&application.RetailerCity,
		&application.RetailerPincode,
		&application.RetailerDateOfBirth,
		&application.RetailerPhone,
	); err != nil {
		return nil, err
	}

	return &application, nil
}

func (pa *PostgresAEPSOnboardingStore) CheckAEPSApplicationStatus(userId string) (*models.AEPSApplicationResponseModel, error) {
	query := `
		SELECT 
			aeps_application_id,
			aeps_application_status
		FROM aeps_applications
		WHERE user_id = $1;
	`

	var status models.AEPSApplicationResponseModel
	if err := pa.db.QueryRow(
		query,
	).Scan(
		&status.AEPSApplicationID,
		&status.AEPSApplicationStatus,
	); err != nil {
		return nil, err
	}

	return &status, nil
}

func (pa *PostgresAEPSOnboardingStore) GetUserDetailsForAEPSSignup(userId string) (*models.AEPSOnboardingSubMerchantSignupRequestModel, error) {
	getUserDetailsQuery := `
		SELECT 
			r.retailer_name,
			r.retailer_phone,
			r.retailer_email,
			r.retailer_aadhar_number,
			r.retailer_pan_number,
			r.retailer_address,
			r.retailer_city,
			r.retailer_pincode,
			r.retailer_date_of_birth,
			r.retailer_gender,
			u.latitude,
			u.longitude,
		FROM aeps_applications u
		JOIN retailers r ON u.user_id = r.retailer_id
		WHERE u.retailer_id = $1;
	`

	var res models.AEPSOnboardingSubMerchantSignupRequestModel
	if err := pa.db.QueryRow(
		getUserDetailsQuery,
		userId,
	).Scan(
		&res.Name,
		&res.Mobile,
		&res.Email,
		&res.Aadhaar,
		&res.Pan,
		&res.Address.Full,
		&res.Address.City,
		&res.Address.Pincode,
		&res.DateOfBirth,
		&res.Gender,
		&res.Latitude,
		&res.Longitude,
	); err != nil {
		return nil, err
	}

	return &res, nil
}

func (pa *PostgresAEPSOnboardingStore) SubMerchantAEPSSignup(userId string, data *models.AEPSOnboardingSubMerchantSignupSuccessResponseModel) error {
	query := `
		INSERT INTO aeps_merchant_details(
			retailer_id,
			sub_merchant_id,
			parent_merchant_id,
			outlet_id,
			min_kyc_status,
			ekyc_status,
			mobile_change_state,
			i_pay_uuid,
			timestamp
		) VALUES(
			$1, $2, $3, $4, $5, $6, $7, $8, $9
		);
	`

	res, err := pa.db.Exec(
		query,
		userId,
		data.SubMerchantID,
		data.ParentMerchantID,
		data.OutletID,
		data.MinKYCStatus,
		data.EKYCStatus,
		data.MobileChangeState,
		data.IPayUUID,
		data.Timestamp,
	)

	if err != nil {
		return err
	}

	return checkRowsAffected(res)
}

func (pa *PostgresAEPSOnboardingStore) GetSubMerchantIDForAEPSEKYCCheck(userId string) (string, error) {
	query := `
		SELECT 
			sub_merchant_id
		FROM aeps_merchant_details
		WHERE retailer_id = $1;
	`

	var subMerchantId string
	if err := pa.db.QueryRow(query, userId).Scan(&subMerchantId); err != nil {
		return "", nil
	}
	return subMerchantId, nil
}

func (pa *PostgresAEPSOnboardingStore) MerchantAEPSEKYCCheck(userId string, data *models.AEPSOnboardingeKycCheckResponseModel) error {
	query := `
		UPDATE aeps_merchant_details
		SET ekyc_action = $1,
			ekyc_status = $4,
			reference_key = COALESCE(NULLIF($2, ''), reference_key),
		 	updated_at = NOW()
		WHERE retailer_id = $3;
	`

	res, err := pa.db.Exec(query, data.EKYCAction, data.ReferenceKey, userId, data.EKYCStatus)
	if err != nil {
		return err
	}

	return checkRowsAffected(res)
}

func (pa *PostgresAEPSOnboardingStore) GetMerchantDetailsForBiometricKYC(userId string, req *models.AEPSOnboardingBiometricKYCRequestModel) error {
	query := `
		SELECT
			sub_merchant_id,
			reference_key
		FROM aeps_merchant_details
		WHERE retailer_id = $1;
	`

	if err := pa.db.QueryRow(query, userId).Scan(&req.SubMerchantID, &req.ReferenceKey); err != nil {
		return err
	}

	return nil
}
