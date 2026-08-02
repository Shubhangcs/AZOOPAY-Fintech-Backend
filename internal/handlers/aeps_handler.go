package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
	"github.com/levionstudio/fintech/internal/models"
	"github.com/levionstudio/fintech/internal/store"
	"github.com/levionstudio/fintech/internal/utils"
)

type AEPSHandler struct {
	AEPSStore store.AEPSStore
	logger    *slog.Logger
}

func NewAEPSHandler(AEPSStore store.AEPSStore, logger *slog.Logger) *AEPSHandler {
	return &AEPSHandler{
		AEPSStore,
		logger,
	}
}

func (ah *AEPSHandler) GetAEPSBanks(w http.ResponseWriter, r *http.Request) {
	retailerId, err := utils.ReadParamID(r)
	if err != nil {
		utils.BadRequest(w, ah.logger, "get aeps banks", err)
		return
	}

	data, err := ah.AEPSStore.GetAEPSDetailsByRetailerID(retailerId)
	if err != nil {
		utils.ServerError(w, ah.logger, "get aeps banks", err)
		return
	}
	var res []models.AEPSBankModel

	if err := utils.GetRequest2(
		utils.PayntricAPI+utils.AEPSBankList+`?outletId=`+data.OutletID,
		"token",
		utils.PayntricAPIToken,
		"username",
		utils.PayntricUsername,
		&res,
	); err != nil {
		utils.BadRequest(w, ah.logger, "get aeps banks", err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"message": "aeps banks fetched successfully", "banks": res})
}

func (ah *AEPSHandler) DailyLogin(w http.ResponseWriter, r *http.Request) {
	retailerId, err := utils.ReadParamID(r)
	if err != nil {
		utils.BadRequest(w, ah.logger, "aeps daily login", err)
		return
	}

	var req models.AEPSDailyLoginRequestModel
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.BadRequest(w, ah.logger, "aeps daily login", err)
		return
	}

	data, err := ah.AEPSStore.GetAEPSDetailsByRetailerID(retailerId)
	if err != nil {
		utils.ServerError(w, ah.logger, "aeps daily login", err)
		return
	}

	res, err := aepsDailyLogin(data, &req)
	if err != nil {
		utils.BadRequest(w, ah.logger, "aeps daily login", err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"message": "daily login successfull", "details": res})

}

func aepsDailyLogin(req *models.AEPSDetailsModel, bd *models.AEPSDailyLoginRequestModel) (*models.AEPSDailyLoginResponseModel, error) {
	var res models.AEPSDailyLoginResponseModel

	if err := utils.PostRequest2(
		utils.PayntricAPI+utils.AEPSOutletLogin,
		"token",
		utils.PayntricAPIToken,
		"username",
		utils.PayntricUsername,
		map[string]any{
			"type":          "DAILY_LOGIN",
			"externalRef":   uuid.NewString(),
			"outletId":      req.OutletID,
			"latitude":      bd.Latitude,
			"longitude":     bd.Longitude,
			"captureType":   "FINGER",
			"aadhaar":       req.AadhaarNumber,
			"biometricData": bd.BiometricData,
		},
		&res,
	); err != nil {
		return nil, err
	}

	if res.Status == "FAILED" || res.Status == "FAILURE" || res.Status == "Failure" {
		return nil, errors.New(res.Message)
	}

	return &res, nil
}

func (ah *AEPSHandler) CheckAEPSDailyLoginStatus(w http.ResponseWriter, r *http.Request) {
	retailerId, err := utils.ReadParamID(r)
	if err != nil {
		utils.BadRequest(w, ah.logger, "check aeps daily login", err)
		return
	}

	data, err := ah.AEPSStore.GetAEPSDetailsByRetailerID(retailerId)
	if err != nil {
		utils.ServerError(w, ah.logger, "check aeps daily login", err)
		return
	}

	res, err := checkAEPSDailyLoginStatus(data)
	if err != nil {
		utils.BadRequest(w, ah.logger, "check aeps daily login", err)
		return
	}

	fmt.Println(res)

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"message": "daily login check successfull"})
}

func checkAEPSDailyLoginStatus(req *models.AEPSDetailsModel) (*models.AEPSDailyLoginResponseModel, error) {
	var res models.AEPSDailyLoginResponseModel

	if err := utils.PostRequest2(
		utils.PayntricAPI+utils.AEPSOutletLoginStatus+"?outletId="+req.OutletID,
		"token",
		utils.PayntricAPIToken,
		"username",
		utils.PayntricUsername,
		map[string]any{},
		&res,
	); err != nil {
		return nil, err
	}

	if res.Status == "FAILED" || res.Status == "FAILURE" || res.Status == "Failure" {
		return nil, errors.New(res.Message)
	}

	return &res, nil
}

func (ah *AEPSHandler) GetMiniStatement(w http.ResponseWriter, r *http.Request) {
	var req models.AEPSMiniStatementRequestModel
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.BadRequest(w, ah.logger, "aeps mini statement", err)
		return
	}
	req.RequestID = uuid.NewString()

	data, err := ah.AEPSStore.GetAEPSDetailsByRetailerID(req.RetailerID)
	if err != nil {
		utils.ServerError(w, ah.logger, "aeps mini statement", err)
		return
	}

	res, err := getMiniStatement(data, &req)
	if err != nil {
		utils.BadRequest(w, ah.logger, "aeps mini statement", err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"message": "mini statement fetched successfully", "mini_statement": res})
}

func getMiniStatement(retailerData *models.AEPSDetailsModel, req *models.AEPSMiniStatementRequestModel) (*models.AEPSMiniStatementResponseModel, error) {
	var res models.AEPSMiniStatementResponseModel

	if err := utils.PostRequest2(
		utils.PayntricAPI+utils.AEPSMiniStatement,
		"token",
		utils.PayntricAPIToken,
		"username",
		utils.PayntricUsername,
		map[string]any{
			"requestId":     req.RequestID,
			"outletId":      retailerData.OutletID,
			"bankin":        req.BankID,
			"mobile":        req.MobileNumber,
			"amount":        0,
			"latitude":      req.Latitude,
			"longitude":     req.Longitude,
			"customerName":  req.CustomerName,
			"captureType":   "finger",
			"aadhaar":       req.AadhaarNumber,
			"biometricData": req.BiometricData,
		},
		&res,
	); err != nil {
		return nil, err
	}

	if res.Status == "FAILED" || res.Status == "FAILURE" || res.Status == "Failure" {
		return nil, errors.New(res.Message)
	}

	return &res, nil
}

func (ah *AEPSHandler) GetBalance(w http.ResponseWriter, r *http.Request) {
	var req models.AEPSBalanceEnquiryRequestModel
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.BadRequest(w, ah.logger, "aeps balance enquiry", err)
		return
	}
	req.RequestID = uuid.NewString()

	data, err := ah.AEPSStore.GetAEPSDetailsByRetailerID(req.RetailerID)
	if err != nil {
		utils.ServerError(w, ah.logger, "aeps balance enquiry", err)
		return
	}

	res, err := getBalance(data, &req)
	if err != nil {
		utils.BadRequest(w, ah.logger, "aeps balance enquiry", err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"message": "mini statement fetched successfully", "mini_statement": res})
}

func getBalance(retailerData *models.AEPSDetailsModel, req *models.AEPSBalanceEnquiryRequestModel) (*models.AEPSMiniStatementResponseModel, error) {
	var res models.AEPSMiniStatementResponseModel

	if err := utils.PostRequest2(
		utils.PayntricAPI+utils.AEPSBalanceEnquiry,
		"token",
		utils.PayntricAPIToken,
		"username",
		utils.PayntricUsername,
		map[string]any{
			"requestId":     req.RequestID,
			"outletId":      retailerData.OutletID,
			"bankin":        req.BankID,
			"mobile":        req.MobileNumber,
			"amount":        0,
			"latitude":      req.Latitude,
			"longitude":     req.Longitude,
			"customerName":  req.CustomerName,
			"captureType":   "finger",
			"aadhaar":       req.AadhaarNumber,
			"biometricData": req.BiometricData,
		},
		&res,
	); err != nil {
		return nil, err
	}

	if res.Status == "FAILED" || res.Status == "FAILURE" || res.Status == "Failure" {
		return nil, errors.New(res.Message)
	}

	return &res, nil
}
