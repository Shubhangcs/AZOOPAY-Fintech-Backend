package store

import (
	"database/sql"
	"errors"

	"github.com/levionstudio/fintech/internal/models"
)

type PostgresDistributorStore struct {
	db *sql.DB
}

func NewPostgresDistributorStore(db *sql.DB) *PostgresDistributorStore {
	return &PostgresDistributorStore{db: db}
}

type DistributorStore interface {
	CreateDistributor(d *models.DistributorModel) error
	UpdateDistributorDetails(d *models.DistributorModel) error
	UpdateDistributorPassword(d *models.DistributorModel) error
	UpdateDistributorMpin(d *models.DistributorModel) error
	UpdateDistributorKYCStatus(d *models.DistributorModel) error
	UpdateDistributorBlockStatus(d *models.DistributorModel) error
	GetDistributorByID(id string) (*models.DistributorModel, error)
	GetDistributorsByMasterDistributorID(masterDistributorID string, limit, offset int) ([]models.DistributorModel, error)
	GetDistributorsByAdminID(adminID string, limit, offset int) ([]models.DistributorModel, error)
	GetDistributorDetailsForLogin(d *models.DistributorModel) error
	GetDistributorsByMasterDistributorIDForDropdown(mdID string) ([]models.DropdownItem, error)
	GetDistributorsByAdminIDForDropdown(adminID string) ([]models.DropdownItem, error)
	ChangeDistributorsMasterDistributor(distributorID, masterDistributorID string) error
	DeleteDistributor(id string) error
	UpdateDistributorAadharFrontImage(path, id string) error
	UpdateDistributorAadharBackImage(path, id string) error
	UpdateDistributorPanImage(path, id string) error
	UpdateDistributorPanWithAgentImage(path, id string) error
	UpdateDistributorShopImage(path, id string) error
	UpdateDistributorSignatureImage(path, id string) error
	UpdateDistributorSelfieImage(path, id string) error
	GetDistributorWalletBalance(id string) (float64, error)
	UpdateDistributorHoldAmount(id string, amount float64) error
}

// Create Distributor
func (ds *PostgresDistributorStore) CreateDistributor(d *models.DistributorModel) error {
	query := `
	INSERT INTO distributors (
		master_distributor_id,
		distributor_name,
		distributor_phone,
		distributor_email,
		distributor_password,
		distributor_aadhar_number,
		distributor_pan_number,
		distributor_date_of_birth,
		distributor_gender,
		distributor_city,
		distributor_state,
		distributor_address,
		distributor_pincode,
		distributor_business_name,
		distributor_business_type,
		distributor_gst_number
	) VALUES (
		$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16
	)
	RETURNING distributor_id, distributor_mpin, distributor_wallet_balance, created_at, updated_at;
	`

	return ds.db.QueryRow(
		query,
		d.MasterDistributorID,
		d.DistributorName,
		d.DistributorPhone,
		d.DistributorEmail,
		d.DistributorPassword,
		d.DistributorAadharNumber,
		d.DistributorPanNumber,
		d.DistributorDateOfBirth,
		d.DistributorGender,
		d.DistributorCity,
		d.DistributorState,
		d.DistributorAddress,
		d.DistributorPincode,
		d.DistributorBusinessName,
		d.DistributorBusinessType,
		d.DistributorGSTNumber,
	).Scan(
		&d.DistributorID,
		&d.DistributorMpin,
		&d.DistributorWalletBalance,
		&d.CreatedAT,
		&d.UpdatedAT,
	)
}

// Update Distributor
func (ds *PostgresDistributorStore) UpdateDistributorDetails(d *models.DistributorModel) error {
	query := `
	UPDATE distributors
	SET
		distributor_name         = COALESCE(NULLIF($1, ''), distributor_name),
		distributor_phone        = COALESCE(NULLIF($2, ''), distributor_phone),
		distributor_email        = COALESCE(NULLIF($3, ''), distributor_email),
		distributor_city         = COALESCE(NULLIF($4, ''), distributor_city),
		distributor_state        = COALESCE(NULLIF($5, ''), distributor_state),
		distributor_address      = COALESCE(NULLIF($6, ''), distributor_address),
		distributor_pincode      = COALESCE(NULLIF($7, ''), distributor_pincode),
		distributor_business_name = COALESCE(NULLIF($8, ''), distributor_business_name),
		distributor_business_type = COALESCE(NULLIF($9, ''), distributor_business_type),
		distributor_gst_number   = COALESCE($10, distributor_gst_number),
		updated_at               = CURRENT_TIMESTAMP
	WHERE distributor_id = $11;
	`

	res, err := ds.db.Exec(
		query,
		d.DistributorName,
		d.DistributorPhone,
		d.DistributorEmail,
		d.DistributorCity,
		d.DistributorState,
		d.DistributorAddress,
		d.DistributorPincode,
		d.DistributorBusinessName,
		d.DistributorBusinessType,
		d.DistributorGSTNumber,
		d.DistributorID,
	)
	if err != nil {
		return err
	}

	return checkRowsAffected(res)
}

// Update Distributor Password
func (ds *PostgresDistributorStore) UpdateDistributorPassword(d *models.DistributorModel) error {
	query := `
	UPDATE distributors
	SET distributor_password = $1,
		updated_at           = CURRENT_TIMESTAMP
	WHERE distributor_id     = $2;
	`

	res, err := ds.db.Exec(query, d.DistributorPassword, d.DistributorID)
	if err != nil {
		return err
	}

	return checkRowsAffected(res)
}

// Update Distributor MPIN
func (ds *PostgresDistributorStore) UpdateDistributorMpin(d *models.DistributorModel) error {
	query := `
	UPDATE distributors
	SET distributor_mpin = $1,
		updated_at       = CURRENT_TIMESTAMP
	WHERE distributor_id = $2;
	`

	res, err := ds.db.Exec(query, d.DistributorMpin, d.DistributorID)
	if err != nil {
		return err
	}

	return checkRowsAffected(res)
}

// Update Distributor KYC Status
func (ds *PostgresDistributorStore) UpdateDistributorKYCStatus(d *models.DistributorModel) error {
	query := `
	UPDATE distributors
	SET distributor_kyc_status = $1,
		updated_at             = CURRENT_TIMESTAMP
	WHERE distributor_id       = $2;
	`

	res, err := ds.db.Exec(query, d.DistributorKYCStatus, d.DistributorID)
	if err != nil {
		return err
	}

	return checkRowsAffected(res)
}

// Update Distributor Block Status
func (ds *PostgresDistributorStore) UpdateDistributorBlockStatus(d *models.DistributorModel) error {
	query := `
	UPDATE distributors
	SET is_distributor_blocked = $1,
		updated_at             = CURRENT_TIMESTAMP
	WHERE distributor_id       = $2;
	`

	res, err := ds.db.Exec(query, d.IsDistributorBlocked, d.DistributorID)
	if err != nil {
		return err
	}

	return checkRowsAffected(res)
}

// Get Distributor By ID
func (ds *PostgresDistributorStore) GetDistributorByID(id string) (*models.DistributorModel, error) {
	query := `
	SELECT
		distributor_id,
		master_distributor_id,
		distributor_name,
		distributor_phone,
		distributor_email,
		distributor_password,
		distributor_mpin,
		distributor_aadhar_number,
		distributor_pan_number,
		distributor_date_of_birth,
		distributor_gender,
		distributor_city,
		distributor_state,
		distributor_address,
		distributor_pincode,
		distributor_business_name,
		distributor_business_type,
		distributor_gst_number,
		distributor_kyc_status,
		distributor_wallet_balance,
		hold_amount,
		is_distributor_blocked,
		distributor_aadhar_front_image,
		distributor_aadhar_back_image,
		distributor_pan_image,
		distributor_pan_with_agent_image,
		distributor_signature_image,
		distributor_shop_image,
		distributor_selfie_image,
		created_at,
		updated_at
	FROM distributors
	WHERE distributor_id = $1;
	`

	var d models.DistributorModel
	err := ds.db.QueryRow(query, id).Scan(
		&d.DistributorID,
		&d.MasterDistributorID,
		&d.DistributorName,
		&d.DistributorPhone,
		&d.DistributorEmail,
		&d.DistributorPassword,
		&d.DistributorMpin,
		&d.DistributorAadharNumber,
		&d.DistributorPanNumber,
		&d.DistributorDateOfBirth,
		&d.DistributorGender,
		&d.DistributorCity,
		&d.DistributorState,
		&d.DistributorAddress,
		&d.DistributorPincode,
		&d.DistributorBusinessName,
		&d.DistributorBusinessType,
		&d.DistributorGSTNumber,
		&d.DistributorKYCStatus,
		&d.DistributorWalletBalance,
		&d.HoldAmount,
		&d.IsDistributorBlocked,
		&d.DistributorAadharFrontImage,
		&d.DistributorAadharBackImage,
		&d.DistributorPanImage,
		&d.DistributorPanWithAgentImage,
		&d.DistributorSignatureImage,
		&d.DistributorShopImage,
		&d.DistributorSelfieImage,
		&d.CreatedAT,
		&d.UpdatedAT,
	)

	return &d, err
}

// Get Distributors By Master Distributor ID
func (ds *PostgresDistributorStore) GetDistributorsByMasterDistributorID(masterDistributorID string, limit, offset int) ([]models.DistributorModel, error) {
	query := `
	SELECT
		distributor_id,
		master_distributor_id,
		distributor_name,
		distributor_phone,
		distributor_email,
		distributor_password,
		distributor_mpin,
		distributor_aadhar_number,
		distributor_pan_number,
		distributor_date_of_birth,
		distributor_gender,
		distributor_city,
		distributor_state,
		distributor_address,
		distributor_pincode,
		distributor_business_name,
		distributor_business_type,
		distributor_gst_number,
		distributor_kyc_status,
		distributor_wallet_balance,
		hold_amount,
		is_distributor_blocked,
		distributor_aadhar_front_image,
		distributor_aadhar_back_image,
		distributor_pan_image,
		distributor_pan_with_agent_image,
		distributor_signature_image,
		distributor_shop_image,
		distributor_selfie_image,
		created_at,
		updated_at
	FROM distributors
	WHERE master_distributor_id = $1
	ORDER BY created_at DESC
	LIMIT $2 OFFSET $3;
	`

	return scanDistributors(ds.db, query, masterDistributorID, limit, offset)
}

// Get Distributors By Admin ID
func (ds *PostgresDistributorStore) GetDistributorsByAdminID(adminID string, limit, offset int) ([]models.DistributorModel, error) {
	query := `
	SELECT
		d.distributor_id,
		d.master_distributor_id,
		d.distributor_name,
		d.distributor_phone,
		d.distributor_email,
		d.distributor_password,
		d.distributor_mpin,
		d.distributor_aadhar_number,
		d.distributor_pan_number,
		d.distributor_date_of_birth,
		d.distributor_gender,
		d.distributor_city,
		d.distributor_state,
		d.distributor_address,
		d.distributor_pincode,
		d.distributor_business_name,
		d.distributor_business_type,
		d.distributor_gst_number,
		d.distributor_kyc_status,
		d.distributor_wallet_balance,
		d.hold_amount,
		d.is_distributor_blocked,
		d.distributor_aadhar_front_image,
		d.distributor_aadhar_back_image,
		d.distributor_pan_image,
		d.distributor_pan_with_agent_image,
		d.distributor_signature_image,
		d.distributor_shop_image,
		d.distributor_selfie_image,
		d.created_at,
		d.updated_at
	FROM distributors d
	JOIN master_distributors md ON d.master_distributor_id = md.master_distributor_id
	WHERE md.admin_id = $1
	ORDER BY d.created_at DESC
	LIMIT $2 OFFSET $3;
	`

	return scanDistributors(ds.db, query, adminID, limit, offset)
}

// Get Distributors By Master Distributor ID For Dropdown
func (ds *PostgresDistributorStore) GetDistributorsByMasterDistributorIDForDropdown(mdID string) ([]models.DropdownItem, error) {
	query := `
	SELECT distributor_id, distributor_name
	FROM distributors
	WHERE master_distributor_id = $1
	ORDER BY distributor_name;
	`
	return scanDropdown(ds.db, query, mdID)
}

// Get Distributors By Admin ID For Dropdown
func (ds *PostgresDistributorStore) GetDistributorsByAdminIDForDropdown(adminID string) ([]models.DropdownItem, error) {
	query := `
	SELECT d.distributor_id, d.distributor_name
	FROM distributors d
	JOIN master_distributors md ON d.master_distributor_id = md.master_distributor_id
	WHERE md.admin_id = $1
	ORDER BY d.distributor_name;
	`
	return scanDropdown(ds.db, query, adminID)
}

// Get Distributor Details For Login
func (ds *PostgresDistributorStore) GetDistributorDetailsForLogin(d *models.DistributorModel) error {
	query := `
	SELECT
		dist.distributor_id,
		dist.distributor_name,
		md.admin_id
	FROM distributors dist
	JOIN master_distributors md ON dist.master_distributor_id = md.master_distributor_id
	WHERE dist.distributor_phone = $1
	AND dist.distributor_password = $2
	AND dist.is_distributor_blocked = FALSE;
	`

	err := ds.db.QueryRow(query, d.DistributorPhone, d.DistributorPassword).Scan(
		&d.DistributorID,
		&d.DistributorName,
		&d.AdminID,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return errors.New("invalid credentials")
	}

	return err
}

// Delete Distributor
func (ds *PostgresDistributorStore) DeleteDistributor(id string) error {
	query := `
	DELETE FROM distributors
	WHERE distributor_id = $1;
	`

	res, err := ds.db.Exec(query, id)
	if err != nil {
		return err
	}

	return checkRowsAffected(res)
}

func scanDistributors(db *sql.DB, query string, args ...any) ([]models.DistributorModel, error) {
	rows, err := db.Query(query, args...)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return []models.DistributorModel{}, nil
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var distributors []models.DistributorModel
	for rows.Next() {
		var d models.DistributorModel
		err = rows.Scan(
			&d.DistributorID,
			&d.MasterDistributorID,
			&d.DistributorName,
			&d.DistributorPhone,
			&d.DistributorEmail,
			&d.DistributorPassword,
			&d.DistributorMpin,
			&d.DistributorAadharNumber,
			&d.DistributorPanNumber,
			&d.DistributorDateOfBirth,
			&d.DistributorGender,
			&d.DistributorCity,
			&d.DistributorState,
			&d.DistributorAddress,
			&d.DistributorPincode,
			&d.DistributorBusinessName,
			&d.DistributorBusinessType,
			&d.DistributorGSTNumber,
			&d.DistributorKYCStatus,
			&d.DistributorWalletBalance,
			&d.HoldAmount,
			&d.IsDistributorBlocked,
			&d.DistributorAadharFrontImage,
			&d.DistributorAadharBackImage,
			&d.DistributorPanImage,
			&d.DistributorPanWithAgentImage,
			&d.DistributorSignatureImage,
			&d.DistributorShopImage,
			&d.DistributorSelfieImage,
			&d.CreatedAT,
			&d.UpdatedAT,
		)
		if err != nil {
			return nil, err
		}
		distributors = append(distributors, d)
	}

	return distributors, rows.Err()
}

// Change Distributors Master Distributor
func (ds *PostgresDistributorStore) ChangeDistributorsMasterDistributor(distributorID, masterDistributorID string) error {
	res, err := ds.db.Exec(`
		UPDATE distributors
		SET master_distributor_id = $1, updated_at = CURRENT_TIMESTAMP
		WHERE distributor_id = $2
	`, masterDistributorID, distributorID)
	if err != nil {
		return err
	}
	return checkRowsAffected(res)
}

// Update Distributor Aadhar Front Image
func (ds *PostgresDistributorStore) UpdateDistributorAadharFrontImage(path, id string) error {
	query := `
		UPDATE distributors
		SET distributor_aadhar_front_image = $1,
		updated_at = CURRENT_TIMESTAMP
		WHERE distributor_id = $2;
	`

	res, err := ds.db.Exec(query, path, id)
	if err != nil {
		return err
	}

	return checkRowsAffected(res)
}

// Update Distributor Aadhar Back Image
func (ds *PostgresDistributorStore) UpdateDistributorAadharBackImage(path, id string) error {
	query := `
		UPDATE distributors
		SET distributor_aadhar_back_image = $1,
		updated_at = CURRENT_TIMESTAMP
		WHERE distributor_id = $2;
	`

	res, err := ds.db.Exec(query, path, id)
	if err != nil {
		return err
	}

	return checkRowsAffected(res)
}

// Update Distributor Pan Image
func (ds *PostgresDistributorStore) UpdateDistributorPanImage(path, id string) error {
	query := `
		UPDATE distributors
		SET distributor_pan_image = $1,
		updated_at = CURRENT_TIMESTAMP
		WHERE distributor_id = $2;
	`

	res, err := ds.db.Exec(query, path, id)
	if err != nil {
		return err
	}

	return checkRowsAffected(res)
}

// Update Distributor Pan With Agent Image
func (ds *PostgresDistributorStore) UpdateDistributorPanWithAgentImage(path, id string) error {
	query := `
		UPDATE distributors
		SET distributor_pan_with_agent_image = $1,
		updated_at = CURRENT_TIMESTAMP
		WHERE distributor_id = $2;
	`

	res, err := ds.db.Exec(query, path, id)
	if err != nil {
		return err
	}

	return checkRowsAffected(res)
}

// Update Distributor Shop Image
func (ds *PostgresDistributorStore) UpdateDistributorShopImage(path, id string) error {
	query := `
		UPDATE distributors
		SET distributor_shop_image = $1,
		updated_at = CURRENT_TIMESTAMP
		WHERE distributor_id = $2;
	`

	res, err := ds.db.Exec(query, path, id)
	if err != nil {
		return err
	}

	return checkRowsAffected(res)
}

// Update Distributor Signature Image
func (ds *PostgresDistributorStore) UpdateDistributorSignatureImage(path, id string) error {
	query := `
		UPDATE distributors
		SET distributor_signature_image = $1,
		updated_at = CURRENT_TIMESTAMP
		WHERE distributor_id = $2;
	`

	res, err := ds.db.Exec(query, path, id)
	if err != nil {
		return err
	}

	return checkRowsAffected(res)
}

// Update Distributor Selfie Image
func (ds *PostgresDistributorStore) UpdateDistributorSelfieImage(path, id string) error {
	query := `
		UPDATE distributors
		SET distributor_selfie_image = $1,
		updated_at = CURRENT_TIMESTAMP
		WHERE distributor_id = $2;
	`

	res, err := ds.db.Exec(query, path, id)
	if err != nil {
		return err
	}

	return checkRowsAffected(res)
}

// Get Distributor Wallet Balance
func (ds *PostgresDistributorStore) GetDistributorWalletBalance(id string) (float64, error) {
	query := `
		SELECT
			distributor_wallet_balance
		FROM distributors
		WHERE distributor_id = $1;
	`
	var balance float64
	err := ds.db.QueryRow(
		query,
		id,
	).Scan(
		&balance,
	)

	return balance, err
}

// Update Distributor Hold Amount
func (ds *PostgresDistributorStore) UpdateDistributorHoldAmount(id string, amount float64) error {
	res, err := ds.db.Exec(
		`UPDATE distributors SET hold_amount = $1, updated_at = CURRENT_TIMESTAMP WHERE distributor_id = $2`,
		amount, id,
	)
	if err != nil {
		return err
	}
	return checkRowsAffected(res)
}
