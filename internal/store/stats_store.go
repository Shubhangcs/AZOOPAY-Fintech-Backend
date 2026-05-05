package store

import (
	"database/sql"
	"fmt"

	"github.com/levionstudio/fintech/internal/models"
)

type PostgresStatsStore struct {
	db *sql.DB
}

func NewPostgresStatsStore(db *sql.DB) *PostgresStatsStore {
	return &PostgresStatsStore{db: db}
}

type StatsStore interface {
	GetRetailerStats(retailerID string) (*models.StatsResponse, error)
	GetDistributorStats(distributorID string) (*models.StatsResponse, error)
	GetMasterDistributorStats(masterDistributorID string) (*models.StatsResponse, error)
}

const statsQuery = `
SELECT
	COALESCE(SUM(CASE WHEN created_at >= CURRENT_DATE THEN amount END), 0),
	COUNT(CASE WHEN created_at >= CURRENT_DATE THEN 1 END),
	COALESCE(SUM(CASE WHEN created_at >= date_trunc('week', CURRENT_DATE) THEN amount END), 0),
	COUNT(CASE WHEN created_at >= date_trunc('week', CURRENT_DATE) THEN 1 END),
	COALESCE(SUM(CASE WHEN created_at >= date_trunc('month', CURRENT_DATE) THEN amount END), 0),
	COUNT(CASE WHEN created_at >= date_trunc('month', CURRENT_DATE) THEN 1 END),
	COALESCE(SUM(amount), 0),
	COUNT(*)
FROM (%s) t
`

func scanStats(row *sql.Row) (*models.StatsResponse, error) {
	var s models.StatsResponse
	err := row.Scan(
		&s.Today.TotalAmount, &s.Today.TotalTransactions,
		&s.Weekly.TotalAmount, &s.Weekly.TotalTransactions,
		&s.Monthly.TotalAmount, &s.Monthly.TotalTransactions,
		&s.Yearly.TotalAmount, &s.Yearly.TotalTransactions,
	)
	return &s, err
}

func (ss *PostgresStatsStore) GetRetailerStats(retailerID string) (*models.StatsResponse, error) {
	union := `
		SELECT amount, created_at FROM payout_transactions
		WHERE retailer_id = $1 AND payout_transaction_status = 'SUCCESS'
		AND created_at >= date_trunc('year', CURRENT_DATE)
		UNION ALL
		SELECT amount, created_at FROM mobile_recharge
		WHERE retailer_id = $1 AND recharge_status = 'SUCCESS'
		AND created_at >= date_trunc('year', CURRENT_DATE)
		UNION ALL
		SELECT amount, created_at FROM dth_recharge
		WHERE retailer_id = $1 AND status = 'SUCCESS'
		AND created_at >= date_trunc('year', CURRENT_DATE)
		UNION ALL
		SELECT amount, created_at FROM electricity_bill_payments
		WHERE retailer_id = $1 AND transaction_status = 'SUCCESS'
		AND created_at >= date_trunc('year', CURRENT_DATE)
	`
	query := fmt.Sprintf(statsQuery, union)
	return scanStats(ss.db.QueryRow(query, retailerID))
}

func (ss *PostgresStatsStore) GetDistributorStats(distributorID string) (*models.StatsResponse, error) {
	union := `
		SELECT pt.amount, pt.created_at FROM payout_transactions pt
		JOIN retailers re ON pt.retailer_id = re.retailer_id
		WHERE re.distributor_id = $1 AND pt.payout_transaction_status = 'SUCCESS'
		AND pt.created_at >= date_trunc('year', CURRENT_DATE)
		UNION ALL
		SELECT mr.amount, mr.created_at FROM mobile_recharge mr
		JOIN retailers re ON mr.retailer_id = re.retailer_id
		WHERE re.distributor_id = $1 AND mr.recharge_status = 'SUCCESS'
		AND mr.created_at >= date_trunc('year', CURRENT_DATE)
		UNION ALL
		SELECT dr.amount, dr.created_at FROM dth_recharge dr
		JOIN retailers re ON dr.retailer_id = re.retailer_id
		WHERE re.distributor_id = $1 AND dr.status = 'SUCCESS'
		AND dr.created_at >= date_trunc('year', CURRENT_DATE)
		UNION ALL
		SELECT eb.amount, eb.created_at FROM electricity_bill_payments eb
		JOIN retailers re ON eb.retailer_id = re.retailer_id
		WHERE re.distributor_id = $1 AND eb.transaction_status = 'SUCCESS'
		AND eb.created_at >= date_trunc('year', CURRENT_DATE)
	`
	query := fmt.Sprintf(statsQuery, union)
	return scanStats(ss.db.QueryRow(query, distributorID))
}

func (ss *PostgresStatsStore) GetMasterDistributorStats(masterDistributorID string) (*models.StatsResponse, error) {
	union := `
		SELECT pt.amount, pt.created_at FROM payout_transactions pt
		JOIN retailers re ON pt.retailer_id = re.retailer_id
		JOIN distributors d ON re.distributor_id = d.distributor_id
		WHERE d.master_distributor_id = $1 AND pt.payout_transaction_status = 'SUCCESS'
		AND pt.created_at >= date_trunc('year', CURRENT_DATE)
		UNION ALL
		SELECT mr.amount, mr.created_at FROM mobile_recharge mr
		JOIN retailers re ON mr.retailer_id = re.retailer_id
		JOIN distributors d ON re.distributor_id = d.distributor_id
		WHERE d.master_distributor_id = $1 AND mr.recharge_status = 'SUCCESS'
		AND mr.created_at >= date_trunc('year', CURRENT_DATE)
		UNION ALL
		SELECT dr.amount, dr.created_at FROM dth_recharge dr
		JOIN retailers re ON dr.retailer_id = re.retailer_id
		JOIN distributors d ON re.distributor_id = d.distributor_id
		WHERE d.master_distributor_id = $1 AND dr.status = 'SUCCESS'
		AND dr.created_at >= date_trunc('year', CURRENT_DATE)
		UNION ALL
		SELECT eb.amount, eb.created_at FROM electricity_bill_payments eb
		JOIN retailers re ON eb.retailer_id = re.retailer_id
		JOIN distributors d ON re.distributor_id = d.distributor_id
		WHERE d.master_distributor_id = $1 AND eb.transaction_status = 'SUCCESS'
		AND eb.created_at >= date_trunc('year', CURRENT_DATE)
	`
	query := fmt.Sprintf(statsQuery, union)
	return scanStats(ss.db.QueryRow(query, masterDistributorID))
}
