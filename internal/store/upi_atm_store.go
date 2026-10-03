package store

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/levionstudio/fintech/internal/models"
	"github.com/levionstudio/fintech/internal/utils"
)

type UPIATMStore interface {
	GetUPIATMDetailsByRetailerID(retailerId string) (*models.UPIATMDetailsModel, error)
	ValidateRetailerForQR(retailerId string, amount float64) error
	CreateQR(retailerId string, req *models.UPIATMCreateQRAPIResponseModel) error
	GetQRDetailsForStatusCheck(requestId string) (*models.UPIATMCheckQRTransactionStatusAPIRequestModel, error)
	FinilizeQRStatus(req *models.UPIATMCheckQRTransactionStatusAPIResponseModel) error
	GetPendingQRsForStatusCheck(limit int) ([]models.UPIATMCheckQRTransactionStatusAPIRequestModel, error)
	GetUPIATMTransactionsByRetailerID(retailerId string, p utils.QueryParams) ([]models.UPIATMTransactionResponseModel, error)
	GetALLUPIATMTransactions(p utils.QueryParams) ([]models.UPIATMTransactionResponseModel, error)
}

// ErrUPIATMDebitFailed wraps the checkExistsTx reason ("insufficient balance" /
// "retailer not found") when a QR payment succeeded but debitTx could not debit
// the retailer. The QR is left non-final so the status poller retries the debit.
var ErrUPIATMDebitFailed = errors.New("upi atm payment received but wallet debit failed")

type PostgresUPIATMStore struct {
	db          *sql.DB
	walletStore WalletTransactionStore
}

func NewPostgresUPIATMStore(db *sql.DB, walletStore WalletTransactionStore) *PostgresUPIATMStore {
	return &PostgresUPIATMStore{
		db,
		walletStore,
	}
}

func (us *PostgresUPIATMStore) GetUPIATMDetailsByRetailerID(retailerId string) (*models.UPIATMDetailsModel, error) {
	query := `
		SELECT
			m.outlet_id,
			m.ekyc_status,
			m.is_merchant_blocked,
			a.latitude,
			a.longitude
		FROM upi_atm_merchant_details m
		JOIN upi_atm_applications a ON a.retailer_id = m.retailer_id
		WHERE m.retailer_id = $1;
	`

	var res models.UPIATMDetailsModel
	if err := us.db.QueryRow(query, retailerId).Scan(
		&res.OutletID,
		&res.EKYCStatus,
		&res.IsMerchantBlocked,
		&res.Latitude,
		&res.Longitude,
	); err != nil {
		return nil, err
	}

	return &res, nil
}

// ValidateRetailerForQR applies the same pre-checks the other wallet services
// run before debiting (KYC, blocked, retailer_wallet_balance), since the QR
// amount is debited from the wallet once the payment succeeds.
func (us *PostgresUPIATMStore) ValidateRetailerForQR(retailerId string, amount float64) error {
	rc, err := getRetailerDetails(us.db, retailerId)
	if err != nil {
		return err
	}
	if !rc.kyc {
		return errors.New("retailer KYC is not verified")
	}
	if rc.blocked {
		return errors.New("retailer is blocked")
	}
	if rc.balance < amount {
		return errors.New("insufficient wallet balance")
	}
	return nil
}

func (us *PostgresUPIATMStore) CreateQR(retailerId string, req *models.UPIATMCreateQRAPIResponseModel) error {
	qrStatus, ok := NormalizeUPIATMQRStatus(req.QRStatus)
	if !ok {
		qrStatus = "INITIATED"
	}
	settlementStatus, ok := NormalizeUPIATMSettlementStatus(req.SettlementStatus)
	if !ok {
		settlementStatus = "PENDING"
	}

	query := `
		INSERT INTO upi_atm (
			retailer_id,
			request_id,
			transaction_id,
			ipay_id,
			amount,
			payable_value,
			transaction_value,
			commission_amount,
			tds_amount,
			net_amount,
			settlement_status,
			qr_status
		) VALUES (
			$1 , $2 , $3 , $4 , $5 , $6 , $7 , $8 , $9 , $10, $11 , $12
		);
	`

	res, err := us.db.Exec(
		query,
		retailerId,
		req.RequestID,
		req.TransactionID,
		req.IpayID,
		req.Amount,
		req.PayableValue,
		req.TransactionValue,
		req.CommissionAmount,
		req.TDSAmount,
		req.NETAmount,
		settlementStatus,
		qrStatus,
	)
	if err != nil {
		return err
	}

	return checkRowsAffected(res)
}

func (us *PostgresUPIATMStore) GetQRDetailsForStatusCheck(requestId string) (*models.UPIATMCheckQRTransactionStatusAPIRequestModel, error) {
	query := `
		SELECT
			o.request_id,
			o.transaction_id,
			o.ipay_id,
			n.outlet_id
		FROM upi_atm o
		JOIN upi_atm_merchant_details n ON n.retailer_id = o.retailer_id
		WHERE o.request_id = $1;
	`
	var res models.UPIATMCheckQRTransactionStatusAPIRequestModel
	if err := us.db.QueryRow(query, requestId).Scan(
		&res.RequestID,
		&res.TransactionID,
		&res.IpayID,
		&res.OutletID,
	); err != nil {
		return nil, err
	}

	return &res, nil
}

func isUPIATMFinalQRStatus(status string) bool {
	return status == "SUCCESS" || status == "FAILED"
}

// NormalizeUPIATMQRStatus maps a Payntric qrStatus onto the values allowed by
// the upi_atm_qr_status_check constraint (INITIATED, PENDING, SUCCESS, FAILED).
// ok is false for values it does not recognise.
func NormalizeUPIATMQRStatus(status string) (string, bool) {
	switch strings.ToUpper(strings.TrimSpace(status)) {
	case "INITIATED", "CREATED", "INITIALIZED":
		return "INITIATED", true
	case "PENDING", "PROCESSING", "IN_PROGRESS", "INPROGRESS":
		return "PENDING", true
	case "SUCCESS", "SUCCESSFUL", "COMPLETED", "PAID":
		return "SUCCESS", true
	case "FAILED", "FAILURE", "EXPIRED", "TIMEOUT", "CANCELLED", "CANCELED", "REJECTED", "DECLINED":
		return "FAILED", true
	default:
		return "", false
	}
}

// NormalizeUPIATMSettlementStatus maps a Payntric settlementStatus onto the
// values allowed by the upi_atm_settlement_status_check constraint
// (PENDING, SUCCESS, FAILED). Payntric documents SETTLED for a completed settlement.
func NormalizeUPIATMSettlementStatus(status string) (string, bool) {
	switch strings.ToUpper(strings.TrimSpace(status)) {
	case "PENDING", "PROCESSING", "INITIATED":
		return "PENDING", true
	case "SETTLED", "SUCCESS", "SUCCESSFUL", "COMPLETED":
		return "SUCCESS", true
	case "FAILED", "FAILURE", "REJECTED":
		return "FAILED", true
	default:
		return "", false
	}
}

// FinilizeQRStatus records the latest provider status. The retailer wallet is
// debited with the QR amount exactly once, on the first transition into SUCCESS.
func (us *PostgresUPIATMStore) FinilizeQRStatus(req *models.UPIATMCheckQRTransactionStatusAPIResponseModel) error {
	tx, err := us.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	lockQuery := `
		SELECT
			retailer_id,
			amount,
			qr_status
		FROM upi_atm
		WHERE request_id = $1
		FOR UPDATE;
	`

	var retailerId, currentQRStatus string
	var amount float64
	if err := tx.QueryRow(lockQuery, req.RequestID).Scan(
		&retailerId,
		&amount,
		&currentQRStatus,
	); err != nil {
		return err
	}

	// Unrecognised provider values are not written (they would violate the
	// table's CHECK constraints); the current value is kept and the QR stays
	// in the poller until a known status arrives.
	newQRStatus, ok := NormalizeUPIATMQRStatus(req.QRStatus)
	if !ok || isUPIATMFinalQRStatus(currentQRStatus) {
		// QR outcome is already final; only the settlement status may still move.
		newQRStatus = currentQRStatus
	}

	newSettlementStatus, _ := NormalizeUPIATMSettlementStatus(req.SettlementStatus)

	if !isUPIATMFinalQRStatus(currentQRStatus) && newQRStatus == "SUCCESS" {
		rtTableInfo, err := getUserTableInfo(retailerId)
		if err != nil {
			return err
		}

		if err := debitTx(
			tx,
			transaction{
				UserID:        retailerId,
				ReferenceID:   req.RequestID,
				Amount:        amount,
				Reason:        "UPI_ATM",
				Remarks:       "UPI ATM Amount Debit From: " + retailerId,
				userTableInfo: *rtTableInfo,
			},
			us.walletStore,
		); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return fmt.Errorf("%w: %w", ErrUPIATMDebitFailed, checkExistsTx(tx, rtTableInfo.TableName, rtTableInfo.IDColumnName, retailerId, "retailer"))
			}
			return err
		}
	}

	updateQuery := `
		UPDATE upi_atm
		SET qr_status = COALESCE(NULLIF($1, ''), qr_status),
			settlement_status = COALESCE(NULLIF($2, ''), settlement_status),
			updated_at = NOW()
		WHERE request_id = $3;
	`

	res, err := tx.Exec(updateQuery, newQRStatus, newSettlementStatus, req.RequestID)
	if err != nil {
		return err
	}

	if err := checkRowsAffected(res); err != nil {
		return err
	}

	return tx.Commit()
}

// GetPendingQRsForStatusCheck returns QRs whose outcome is not final yet,
// oldest first. QRs older than a day are skipped so unknown/abandoned
// transactions are not polled forever.
func (us *PostgresUPIATMStore) GetPendingQRsForStatusCheck(limit int) ([]models.UPIATMCheckQRTransactionStatusAPIRequestModel, error) {
	query := `
		SELECT
			o.request_id,
			o.transaction_id,
			o.ipay_id,
			n.outlet_id
		FROM upi_atm o
		JOIN upi_atm_merchant_details n ON n.retailer_id = o.retailer_id
		WHERE o.qr_status NOT IN ('SUCCESS', 'FAILED')
		AND o.created_at >= NOW() - INTERVAL '1 day'
		ORDER BY o.created_at ASC
		LIMIT $1;
	`

	rows, err := us.db.Query(query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []models.UPIATMCheckQRTransactionStatusAPIRequestModel
	for rows.Next() {
		var r models.UPIATMCheckQRTransactionStatusAPIRequestModel
		if err := rows.Scan(
			&r.RequestID,
			&r.TransactionID,
			&r.IpayID,
			&r.OutletID,
		); err != nil {
			return nil, err
		}
		res = append(res, r)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return res, nil
}

const upiAtmSelectBase = `
	SELECT
		t.retailer_id,
		t.request_id,
		t.transaction_id,
		t.ipay_id,
		t.amount,
		t.payable_value,
		t.transaction_value,
		t.commission_amount,
		t.tds_amount,
		t.net_amount,
		t.settlement_status,
		t.qr_status,
		t.created_at,
		t.updated_at,
		r.retailer_name,
		w.before_balance,
		w.after_balance,
		w.transaction_reason,
		w.remarks
	FROM upi_atm t
	JOIN retailers r ON r.retailer_id = t.retailer_id
	LEFT JOIN wallet_transactions w ON w.user_id = t.retailer_id AND w.reference_id = t.request_id AND w.transaction_reason = 'UPI_ATM'
`

func (us *PostgresUPIATMStore) GetUPIATMTransactionsByRetailerID(retailerId string, p utils.QueryParams) ([]models.UPIATMTransactionResponseModel, error) {
	q := upiAtmSelectBase + `
	WHERE t.retailer_id = $7
	AND t.created_at >= COALESCE($3, '-infinity'::TIMESTAMPTZ)
	AND t.created_at <= COALESCE($4, 'infinity'::TIMESTAMPTZ)
	AND ($5::TEXT IS NULL OR t.qr_status = $5)
	AND ($6::TEXT IS NULL OR (
		t.request_id ILIKE '%'||$6||'%' OR
		t.transaction_id ILIKE '%'||$6||'%' OR
		t.ipay_id ILIKE '%'||$6||'%'
	))
	ORDER BY t.created_at DESC
	LIMIT $1 OFFSET $2;
	`
	return scanUPIATMTransactions(us.db, q, p.Limit, p.Offset, p.StartDate, p.EndDate, p.Status, p.Search, retailerId)
}

func (us *PostgresUPIATMStore) GetALLUPIATMTransactions(p utils.QueryParams) ([]models.UPIATMTransactionResponseModel, error) {
	q := upiAtmSelectBase + `
	WHERE t.created_at >= COALESCE($3, '-infinity'::TIMESTAMPTZ)
	AND t.created_at <= COALESCE($4, 'infinity'::TIMESTAMPTZ)
	AND ($5::TEXT IS NULL OR t.qr_status = $5)
	AND ($6::TEXT IS NULL OR (
		t.retailer_id ILIKE '%'||$6||'%' OR
		t.request_id ILIKE '%'||$6||'%' OR
		t.transaction_id ILIKE '%'||$6||'%' OR
		t.ipay_id ILIKE '%'||$6||'%'
	))
	ORDER BY t.created_at DESC
	LIMIT $1 OFFSET $2;
	`
	return scanUPIATMTransactions(us.db, q, p.Limit, p.Offset, p.StartDate, p.EndDate, p.Status, p.Search)
}

func scanUPIATMTransactions(db *sql.DB, query string, args ...any) ([]models.UPIATMTransactionResponseModel, error) {
	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query upi atm transactions: %w", err)
	}
	defer rows.Close()

	var txns []models.UPIATMTransactionResponseModel
	for rows.Next() {
		var txn models.UPIATMTransactionResponseModel
		if err := rows.Scan(
			&txn.RetailerID,
			&txn.RequestID,
			&txn.TransactionID,
			&txn.IpayID,
			&txn.Amount,
			&txn.PayableValue,
			&txn.TransactionValue,
			&txn.CommissionAmount,
			&txn.TDSAmount,
			&txn.NETAmount,
			&txn.SettlementStatus,
			&txn.QRStatus,
			&txn.CreatedAT,
			&txn.UpdatedAT,
			&txn.RetailerName,
			&txn.BeforeBalance,
			&txn.AfterBalance,
			&txn.Reason,
			&txn.Remarks,
		); err != nil {
			return nil, fmt.Errorf("failed to scan upi atm transaction row: %w", err)
		}

		txns = append(txns, txn)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating upi atm transaction rows: %w", err)
	}

	return txns, nil
}
