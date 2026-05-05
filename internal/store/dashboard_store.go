package store

import (
	"database/sql"

	"github.com/levionstudio/fintech/internal/models"
)

type PostgresDashboardStore struct {
	db *sql.DB
}

func NewPostgresDashboardStore(db *sql.DB) *PostgresDashboardStore {
	return &PostgresDashboardStore{db: db}
}

type DashboardStore interface {
	GetRetailerDashboard(retailerID string) (*models.RetailerDashboard, error)
}

func (ds *PostgresDashboardStore) GetRetailerDashboard(retailerID string) (*models.RetailerDashboard, error) {
	var d models.RetailerDashboard

	err := ds.db.QueryRow(`
		SELECT
			COALESCE(SUM(amount), 0), COUNT(*)
		FROM payout_transactions
		WHERE retailer_id = $1
		AND payout_transaction_status = 'SUCCESS'
		AND created_at >= CURRENT_DATE
	`, retailerID).Scan(&d.Payout.TotalAmount, &d.Payout.TotalTransactions)
	if err != nil {
		return nil, err
	}

	err = ds.db.QueryRow(`
		SELECT
			COALESCE(SUM(amount), 0), COUNT(*)
		FROM mobile_recharge
		WHERE retailer_id = $1
		AND recharge_status = 'SUCCESS'
		AND created_at >= CURRENT_DATE
	`, retailerID).Scan(&d.MobileRecharge.TotalAmount, &d.MobileRecharge.TotalTransactions)
	if err != nil {
		return nil, err
	}

	err = ds.db.QueryRow(`
		SELECT
			COALESCE(SUM(amount), 0), COUNT(*)
		FROM dth_recharge
		WHERE retailer_id = $1
		AND status = 'SUCCESS'
		AND created_at >= CURRENT_DATE
	`, retailerID).Scan(&d.DTHRecharge.TotalAmount, &d.DTHRecharge.TotalTransactions)
	if err != nil {
		return nil, err
	}

	err = ds.db.QueryRow(`
		SELECT
			COALESCE(SUM(amount), 0), COUNT(*)
		FROM electricity_bill_payments
		WHERE retailer_id = $1
		AND transaction_status = 'SUCCESS'
		AND created_at >= CURRENT_DATE
	`, retailerID).Scan(&d.ElectricityBill.TotalAmount, &d.ElectricityBill.TotalTransactions)
	if err != nil {
		return nil, err
	}

	err = ds.db.QueryRow(`
		SELECT
			COALESCE(SUM(amount), 0), COUNT(*)
		FROM fund_requests
		WHERE requester_id = $1
		AND created_at >= CURRENT_DATE
	`, retailerID).Scan(&d.FundRequests.TotalAmount, &d.FundRequests.TotalTransactions)
	if err != nil {
		return nil, err
	}

	err = ds.db.QueryRow(`
		SELECT
			COALESCE(SUM(credit_amount), 0),
			COALESCE(SUM(debit_amount), 0)
		FROM wallet_transactions
		WHERE user_id = $1
		AND created_at >= CURRENT_DATE
	`, retailerID).Scan(&d.WalletCredited, &d.WalletDebited)
	if err != nil {
		return nil, err
	}

	return &d, nil
}
