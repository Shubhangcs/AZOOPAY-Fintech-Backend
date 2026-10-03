package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/levionstudio/fintech/internal/models"
	"github.com/levionstudio/fintech/internal/store"
	"github.com/levionstudio/fintech/internal/utils"
)

type UPIATMHandler struct {
	upiAtmStore store.UPIATMStore
	logger      *slog.Logger
}

func NewUPIATMHandler(logger *slog.Logger, upiAtmStore store.UPIATMStore) *UPIATMHandler {
	return &UPIATMHandler{
		upiAtmStore,
		logger,
	}
}

func (ua *UPIATMHandler) HandleCreateUPIQR(w http.ResponseWriter, r *http.Request) {
	retailerId, err := utils.ReadParamID(r)
	if err != nil {
		utils.BadRequest(w, ua.logger, "create upi qr", err)
		return
	}

	var req models.UPIATMCreateQRRequestModel
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.BadRequest(w, ua.logger, "create upi qr", err)
		return
	}

	if req.Amount < 50 || req.Amount > 50000 {
		utils.BadRequest(w, ua.logger, "create upi qr", errors.New("amount must be between 50 and 50000"))
		return
	}

	if len(req.Mobile) != 10 {
		utils.BadRequest(w, ua.logger, "create upi qr", errors.New("mobile must be a 10-digit number"))
		return
	}

	details, err := ua.upiAtmStore.GetUPIATMDetailsByRetailerID(retailerId)
	if err != nil {
		utils.ServerError(w, ua.logger, "create upi qr", err)
		return
	}

	if details.IsMerchantBlocked {
		utils.BadRequest(w, ua.logger, "create upi qr", errors.New("upi atm merchant is blocked"))
		return
	}

	if details.EKYCStatus != "APPROVED" {
		utils.BadRequest(w, ua.logger, "create upi qr", errors.New("upi atm merchant ekyc is not approved"))
		return
	}

	// The amount is debited from the wallet when the payment succeeds.
	if err := ua.upiAtmStore.ValidateRetailerForQR(retailerId, req.Amount); err != nil {
		utils.BadRequest(w, ua.logger, "create upi qr", err)
		return
	}

	var apiReq models.UPIATMCreateQRAPIRequestModel
	apiReq.RequestData = &req
	apiReq.RetailerID = retailerId
	apiReq.RequestID = uuid.NewString()
	apiReq.Latitude = details.Latitude
	apiReq.Longitude = details.Longitude
	apiReq.OutletID = details.OutletID

	if req.Latitude != "" && req.Longitude != "" {
		apiReq.Latitude = req.Latitude
		apiReq.Longitude = req.Longitude
	}

	apiRes, err := createUpiQr(&apiReq)
	if err != nil {
		utils.BadRequest(w, ua.logger, "create upi qr", err)
		return
	}

	if err := ua.upiAtmStore.CreateQR(retailerId, apiRes); err != nil {
		utils.ServerError(w, ua.logger, "create upi qr", err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"message": apiRes.Message, "response": apiRes})
}

func createUpiQr(req *models.UPIATMCreateQRAPIRequestModel) (*models.UPIATMCreateQRAPIResponseModel, error) {
	var res models.UPIATMCreateQRAPIResponseModel
	if err := utils.PostRequest2(
		utils.PayntricAPI+utils.UPIATMCreateQR,
		"username",
		utils.PayntricUsername,
		"token",
		utils.PayntricAPIToken,
		map[string]any{
			"requestId": req.RequestID,
			"amount":    req.RequestData.Amount,
			"mobile":    req.RequestData.Mobile,
			"latitude":  req.Latitude,
			"longitude": req.Longitude,
			"outletId":  req.OutletID,
		},
		&res,
	); err != nil {
		return nil, err
	}

	if res.Status != "SUCCESS" || res.TransactionID == "" {
		if res.Message == "" {
			return nil, errors.New("upi atm qr generation failed")
		}
		return nil, errors.New(res.Message)
	}

	// Some providers do not echo requestId back; it is our key for status checks.
	if res.RequestID == "" {
		res.RequestID = req.RequestID
	}

	return &res, nil
}

func (ua *UPIATMHandler) HandleCheckQRTransactionStatus(w http.ResponseWriter, r *http.Request) {
	var req struct {
		RequestID string `json:"requestId"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.BadRequest(w, ua.logger, "check upi qr status", err)
		return
	}

	apiReq, err := ua.upiAtmStore.GetQRDetailsForStatusCheck(req.RequestID)
	if err != nil {
		utils.ServerError(w, ua.logger, "check upi qr status", err)
		return
	}

	apiRes, err := checkQRTransactionStatus(apiReq)
	if err != nil {
		utils.BadRequest(w, ua.logger, "check upi qr status", err)
		return
	}

	// Always key the update on our stored request id.
	apiRes.RequestID = apiReq.RequestID

	if err := ua.upiAtmStore.FinilizeQRStatus(apiRes); err != nil {
		if errors.Is(err, store.ErrUPIATMDebitFailed) {
			utils.BadRequest(w, ua.logger, "check upi qr status", err)
			return
		}
		utils.ServerError(w, ua.logger, "check upi qr status", err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"message": apiRes.Message, "response": apiRes})
}

func checkQRTransactionStatus(req *models.UPIATMCheckQRTransactionStatusAPIRequestModel) (*models.UPIATMCheckQRTransactionStatusAPIResponseModel, error) {
	var res models.UPIATMCheckQRTransactionStatusAPIResponseModel
	if err := utils.PostRequest2(
		utils.PayntricAPI+utils.UPIATMCheckQRTransactionStatus,
		"username",
		utils.PayntricUsername,
		"token",
		utils.PayntricAPIToken,
		map[string]any{
			"transactionId": req.TransactionID,
			"ipayId":        req.IpayID,
			"requestId":     req.RequestID,
			"outletId":      req.OutletID,
		},
		&res,
	); err != nil {
		return nil, err
	}

	// A FAILED status without a qrStatus is an API error (not found, auth, system),
	// not a transaction outcome.
	if res.QRStatus == "" {
		if res.Message == "" {
			return nil, errors.New("upi atm qr status check failed")
		}
		return nil, errors.New(res.Message)
	}

	return &res, nil
}

// StartQRStatusPoller checks and updates every non-final UPI ATM QR at the given
// interval until ctx is cancelled. Runs never overlap: a slow run simply delays
// the next tick.
func (ua *UPIATMHandler) StartQRStatusPoller(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	ua.logger.Info("upi atm qr status poller started", "interval", interval.String())
	for {
		select {
		case <-ctx.Done():
			ua.logger.Info("upi atm qr status poller stopped")
			return
		case <-ticker.C:
			ua.pollPendingQRs(ctx)
		}
	}
}

func (ua *UPIATMHandler) pollPendingQRs(ctx context.Context) {
	pending, err := ua.upiAtmStore.GetPendingQRsForStatusCheck(200)
	if err != nil {
		ua.logger.Error("upi atm poller: fetch pending qrs", "error", err)
		return
	}

	for i := range pending {
		if ctx.Err() != nil {
			return
		}

		qr := &pending[i]
		apiRes, err := checkQRTransactionStatus(qr)
		if err != nil {
			ua.logger.Warn("upi atm poller: status check", "request_id", qr.RequestID, "error", err)
			continue
		}

		apiRes.RequestID = qr.RequestID
		if err := ua.upiAtmStore.FinilizeQRStatus(apiRes); err != nil {
			ua.logger.Error("upi atm poller: finalize status", "request_id", qr.RequestID, "error", err)
		}
	}
}

func (ua *UPIATMHandler) HandleGetUPIATMTransactionsByRetailerID(w http.ResponseWriter, r *http.Request) {
	retailerId, err := utils.ReadParamID(r)
	if err != nil {
		utils.BadRequest(w, ua.logger, "get upi atm transactions by retailer id", err)
		return
	}

	qp := utils.ReadQueryParams(r)

	res, err := ua.upiAtmStore.GetUPIATMTransactionsByRetailerID(retailerId, qp)
	if err != nil {
		utils.ServerError(w, ua.logger, "get upi atm transactions by retailer id", err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"message": "upi atm transactions fetched successfully", "response": res})
}

func (ua *UPIATMHandler) HandleGetAllUPIATMTransactions(w http.ResponseWriter, r *http.Request) {
	qp := utils.ReadQueryParams(r)

	res, err := ua.upiAtmStore.GetALLUPIATMTransactions(qp)
	if err != nil {
		utils.ServerError(w, ua.logger, "get all upi atm transactions", err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"message": "upi atm transactions fetched successfully", "response": res})
}
