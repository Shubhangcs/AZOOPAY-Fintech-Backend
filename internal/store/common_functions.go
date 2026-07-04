package store

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/levionstudio/fintech/internal/models"
)

type userTableInfo struct {
	TableName               string
	IDColumnName            string
	WalletBalanceColumnName string
	HoldAmountColumnName    string
	MpinColumnName          string
	KYCColumnName           string
}

type transaction struct {
	UserID      string  `json:"user_id"`
	ReferenceID string  `json:"reference_id"`
	Amount      float64 `json:"amount"`
	Reason      string  `json:"reason"`
	Remarks     string  `json:"remarks"`
	userTableInfo
}

func getUserTableInfo(id string) (*userTableInfo, error) {
	if id == "" {
		return nil, errors.New("invalid user id")
	}
	switch string(id[0]) {
	case "A":
		return &userTableInfo{TableName: "admins", IDColumnName: "admin_id", WalletBalanceColumnName: "admin_wallet_balance"}, nil
	case "M":
		return &userTableInfo{TableName: "master_distributors", IDColumnName: "master_distributor_id", WalletBalanceColumnName: "master_distributor_wallet_balance", HoldAmountColumnName: "hold_amount", MpinColumnName: "master_distributor_mpin", KYCColumnName: "master_distributor_kyc_status"}, nil
	case "D":
		return &userTableInfo{TableName: "distributors", IDColumnName: "distributor_id", WalletBalanceColumnName: "distributor_wallet_balance", HoldAmountColumnName: "hold_amount", MpinColumnName: "distributor_mpin", KYCColumnName: "distributor_kyc_status"}, nil
	case "R":
		return &userTableInfo{TableName: "retailers", IDColumnName: "retailer_id", WalletBalanceColumnName: "retailer_wallet_balance", HoldAmountColumnName: "hold_amount", MpinColumnName: "retailer_mpin", KYCColumnName: "retailer_kyc_status"}, nil
	default:
		return nil, errors.New("invalid user id")
	}
}

func checkRowsAffected(res sql.Result) error {
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func scanDropdown(db *sql.DB, query string, args ...any) ([]models.DropdownItem, error) {
	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []models.DropdownItem
	for rows.Next() {
		var item models.DropdownItem
		if err := rows.Scan(&item.ID, &item.Name); err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return items, rows.Err()
}

func debitTx(tx *sql.Tx, txn transaction, wts WalletTransactionStore) error {
	if txn.Amount <= 0 {
		return nil
	}
	var q string
	if txn.HoldAmountColumnName != "" {
		q = fmt.Sprintf(
			`UPDATE %s SET %s = %s - $1, updated_at = CURRENT_TIMESTAMP
			 WHERE %s = $2 AND %s - $1 >= %s
			 RETURNING %s + $1, %s;`,
			txn.TableName, txn.WalletBalanceColumnName,
			txn.WalletBalanceColumnName, txn.IDColumnName,
			txn.WalletBalanceColumnName, txn.HoldAmountColumnName,
			txn.WalletBalanceColumnName, txn.WalletBalanceColumnName,
		)
	} else {
		q = fmt.Sprintf(
			`UPDATE %s SET %s = %s - $1, updated_at = CURRENT_TIMESTAMP
			 WHERE %s = $2 AND %s >= $1
			 RETURNING %s + $1, %s;`,
			txn.TableName, txn.WalletBalanceColumnName,
			txn.WalletBalanceColumnName, txn.IDColumnName,
			txn.WalletBalanceColumnName, txn.WalletBalanceColumnName,
			txn.WalletBalanceColumnName,
		)
	}
	var before, after float64
	err := tx.QueryRow(q, txn.Amount, txn.UserID).Scan(&before, &after)

	if err != nil {
		return err
	}

	err = wts.CreateWalletTransactionTx(tx, &models.WalletTransactionModel{
		UserID: txn.UserID, ReferenceID: txn.ReferenceID,
		DebitAmount: &txn.Amount, BeforeBalance: before, AfterBalance: after,
		TransactionReason: txn.Reason, Remarks: txn.Remarks,
	})
	return err
}

func creditTx(tx *sql.Tx, txn transaction, wts WalletTransactionStore) error {
	q := fmt.Sprintf(
		`UPDATE %s SET %s = %s + $1, updated_at = CURRENT_TIMESTAMP
		 WHERE %s = $2
		 RETURNING %s - $1, %s;`,
		txn.TableName, txn.WalletBalanceColumnName, txn.WalletBalanceColumnName,
		txn.IDColumnName, txn.WalletBalanceColumnName, txn.WalletBalanceColumnName,
	)
	var before, after float64
	err := tx.QueryRow(q, txn.Amount, txn.UserID).Scan(&before, &after)
	if err != nil {
		return err
	}

	err = wts.CreateWalletTransactionTx(tx, &models.WalletTransactionModel{
		UserID: txn.UserID, ReferenceID: txn.ReferenceID, CreditAmount: &txn.Amount,
		BeforeBalance: before, AfterBalance: after, TransactionReason: txn.Reason, Remarks: txn.Remarks,
	})
	return err
}

func verifyMpin(db *sql.DB, userID string, mpin int) error {
	info, err := getUserTableInfo(userID)
	if err != nil {
		return err
	}
	if info.MpinColumnName == "" {
		return nil
	}
	q := fmt.Sprintf(`SELECT EXISTS(SELECT 1 FROM %s WHERE %s = $1 AND %s = $2)`,
		info.TableName, info.IDColumnName, info.MpinColumnName)
	var valid bool
	if err := db.QueryRow(q, userID, mpin).Scan(&valid); err != nil {
		return err
	}
	if !valid {
		return errors.New("invalid mpin")
	}
	return nil
}

func verifyKYC(db *sql.DB, userID string) error {
	info, err := getUserTableInfo(userID)
	if err != nil {
		return err
	}
	if info.KYCColumnName == "" {
		return nil
	}
	q := fmt.Sprintf(`SELECT %s FROM %s WHERE %s = $1`, info.KYCColumnName, info.TableName, info.IDColumnName)
	var verified bool
	if err := db.QueryRow(q, userID).Scan(&verified); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("user not found")
		}
		return err
	}
	if !verified {
		return errors.New("KYC is not verified")
	}
	return nil
}

func checkExistsTx(tx *sql.Tx, table, idCol, id, role string) error {
	q := fmt.Sprintf(`SELECT EXISTS(SELECT 1 FROM %s WHERE %s = $1);`, table, idCol)
	var exists bool
	_ = tx.QueryRow(q, id).Scan(&exists)
	if !exists {
		return fmt.Errorf("%s not found", role)
	}
	return errors.New("insufficient balance")
}

func (ps *PostgresPayoutTransactionStore) resolveCommision(
	retailerID, distributorID, mdID, adminID, service string,
	amount float64,
) *models.CommisionModel {
	for _, userID := range []string{retailerID, distributorID, mdID, adminID} {
		c, err := ps.commisionStore.GetCommisionByUserIDServiceAndAmount(userID, service, amount)
		if err == nil && c != nil {
			return c
		}
	}
	return ps.commisionStore.GetDefaultCommision(amount)
}

func (ps *PostgresPayoutTransactionStore) getRetailerTransactionLimit(
	retailerID, service string,
) (float64, error) {
	limit, _, err := ps.transactionLimitStore.GetTransactionLimitByRetailerIDAndService(&models.TransactionLimitModel{RetailerID: retailerID, Service: service})
	if err != nil {
		return 0, err
	}
	return limit, nil
}

func getRetailerDetails(
	db *sql.DB, retailerID string,
) (retailerChain, error) {
	const q = `
	SELECT
		r.retailer_kyc_status,
		r.is_retailer_blocked,
		r.is_payout_blocked,
		r.retailer_wallet_balance,
		r.distributor_id,
		r.address,
		r.email,
		d.master_distributor_id,
		md.admin_id
	FROM retailers r
	JOIN distributors d         ON r.distributor_id         = d.distributor_id
	JOIN master_distributors md ON d.master_distributor_id  = md.master_distributor_id
	WHERE r.retailer_id = $1;
	`
	var rc retailerChain
	err := db.QueryRow(q, retailerID).Scan(
		&rc.kyc, &rc.blocked, &rc.payoutBlocked, &rc.balance,
		&rc.distributorID, &rc.address, &rc.email, &rc.mdID, &rc.adminID,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return retailerChain{}, errors.New("retailer not found")
		}
		return retailerChain{}, err
	}
	return rc, nil
}
