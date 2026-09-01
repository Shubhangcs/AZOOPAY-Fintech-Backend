package handlers

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
	"github.com/levionstudio/fintech/internal/models"
	"github.com/levionstudio/fintech/internal/store"
	"github.com/levionstudio/fintech/internal/utils"
)

type UPIATMHandler struct {
	aepsStore   store.AEPSStore
	upiAtmStore store.UPIATMStore
	logger      *slog.Logger
}

func NewUPIATMHandler(logger *slog.Logger, upiAtmStore store.UPIATMStore, aepsStore store.AEPSStore) *UPIATMHandler {
	return &UPIATMHandler{
		aepsStore,
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

	details, err := ua.aepsStore.GetAEPSDetailsByRetailerID(retailerId)
	if err != nil {
		utils.ServerError(w, ua.logger, "create upi qr", err)
		return
	}
	var apiReq models.UPIATMCreateQRAPIRequestModel
	apiReq.RequestData = &req
	apiReq.RetailerID = retailerId
	apiReq.RequestID = uuid.NewString()
	apiReq.Latitude = details.Latitude
	apiReq.Longitude = details.Longitude
	apiReq.OutletID = details.OutletID

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

	if res.Status == "FAILED" {
		return nil, errors.New(res.Message)
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

	if apiRes.Status != "PENDING" {
		if err := ua.upiAtmStore.FinilizeQRStatus(apiRes); err != nil {
			utils.ServerError(w, ua.logger, "check upi qr status", err)
			return
		}
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

	if res.Status == "FAILED" && res.QRStatus == "" {
		return nil, errors.New(res.Message)
	}

	return &res, nil
}

func (ua *UPIATMHandler) HandleGetUPIATMTransactionsByRetailerID(w http.ResponseWriter, r *http.Request) {
	retailerId, err := utils.ReadParamID(r)
	if err != nil {
		utils.BadRequest(w, ua.logger, "get upi atm transactions by retailer id", err)
		return
	}

	res, err := ua.upiAtmStore.GetUPIATMTransactionsByRetailerID(retailerId)
	if err != nil {
		utils.ServerError(w, ua.logger, "get upi atm transactions by retailer id", err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"message": "upi atm transactions fetched successfully", "response": res})
}

func (ua *UPIATMHandler) HandleGetAllUPIATMTransactions(w http.ResponseWriter, r *http.Request) {
	res, err := ua.upiAtmStore.GetALLUPIATMTransactions()
	if err != nil {
		utils.ServerError(w, ua.logger, "get upi atm transactions by retailer id", err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"message": "upi atm transactions fetched successfully", "response": res})
}