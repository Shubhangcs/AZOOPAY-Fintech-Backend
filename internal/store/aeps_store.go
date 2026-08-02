package store

import (
	"database/sql"

	"github.com/levionstudio/fintech/internal/models"
)

type AEPSStore interface {
	GetAEPSDetailsByRetailerID(retailerId string) (*models.AEPSDetailsModel, error)
}

type PostgresAEPSStore struct {
	db *sql.DB
}

func NewPostgresAEPSStore(db *sql.DB) *PostgresAEPSStore {
	return &PostgresAEPSStore{
		db,
	}
}

func (as *PostgresAEPSStore) GetAEPSDetailsByRetailerID(retailerId string) (*models.AEPSDetailsModel, error) {
	query := `
		SELECT
			r.retailer_aadhaar_number,
			m.outlet_id
		FROM retailers r
		JOIN aeps_merchant_details m ON m.retailer_id = r.retailer_id
		WHERE r.retailer_id = $1;
	`
	var res models.AEPSDetailsModel
	if err := as.db.QueryRow(
		query,
		retailerId,
	).Scan(
		&res.AadhaarNumber,
		&res.OutletID,
	); err != nil {
		return nil, err
	}

	return &res, nil
}
