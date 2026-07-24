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

type AEPSOnboardingHandler struct {
	logger              *slog.Logger
	AEPSOnboardingStore store.AEPSOnboardingStore
}

func NewAEPSOnboardingHandler(logger *slog.Logger, AEPSOnboardingStore store.AEPSOnboardingStore) *AEPSOnboardingHandler {
	return &AEPSOnboardingHandler{
		logger,
		AEPSOnboardingStore,
	}
}

func (ah *AEPSOnboardingHandler) HandleApplyForAEPS(w http.ResponseWriter, r *http.Request) {
	var req models.ApplyForAEPSRequestModel
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.BadRequest(w, ah.logger, "apply for aeps", err)
		return
	}

	if err := req.Validate(); err != nil {
		utils.BadRequest(w, ah.logger, "apply for aeps", err)
		return
	}

	if err := ah.AEPSOnboardingStore.ApplyForAEPS(&req); err != nil {
		utils.ServerError(w, ah.logger, "apply for aeps", err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"message": "AEPS application submitted successfully"})
}

func (ah *AEPSOnboardingHandler) HandleCheckAEPSApplicationStatus(w http.ResponseWriter, r *http.Request) {
	userId, err := utils.ReadParamID(r)
	if err != nil {
		utils.BadRequest(w, ah.logger, "check aeps application status", err)
		return
	}

	res, err := ah.AEPSOnboardingStore.CheckAEPSApplicationStatus(userId)
	if err != nil {
		utils.ServerError(w, ah.logger, "check aeps application status", err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"message": "AEPS application status fetched successfully", "application": res})
}

func (ah *AEPSOnboardingHandler) HandleChangeAEPSApplicationStatus(w http.ResponseWriter, r *http.Request) {
	var req models.ApplyForAEPSRequestModel
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.BadRequest(w, ah.logger, "change aeps application status", err)
		return
	}

	if err := ah.AEPSOnboardingStore.ChangeAEPSApplicationStatus(&req); err != nil {
		utils.ServerError(w, ah.logger, "change aeps application status", err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"message": "AEPS application status updated successfully"})
}

func (ah *AEPSOnboardingHandler) HandleAEPSSignupMerchant(w http.ResponseWriter, r *http.Request) {
	userId, err := utils.ReadParamID(r)
	if err != nil {
		utils.BadRequest(w, ah.logger, "aeps signup merchant", err)
		return
	}

	res, err := ah.AEPSOnboardingStore.GetAEPSApplication(userId)
	if err != nil {
		fmt.Println("error here")
		utils.ServerError(w, ah.logger, "aeps signup merchant", err)
		return
	}

	apiRes, err := aepsMerchantSignup(res)
	if err != nil {
		utils.BadRequest(w, ah.logger, "aeps signup merchant", err)
		return
	}

	if err := ah.AEPSOnboardingStore.SubMerchantAEPSSignup(userId, apiRes); err != nil {
		utils.BadRequest(w, ah.logger, "aeps signup merchant", err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"message": "AEPS signup successfull", "api_response": apiRes})
}

func aepsMerchantSignup(data *models.AEPSApplicationResponseModel) (*models.AEPSOnboardingSubMerchantSignupSuccessResponseModel, error) {
	var res models.AEPSOnboardingSubMerchantSignupSuccessResponseModel

	switch data.RetailerGender {
	case "MALE":
		data.RetailerGender = "M"
	case "FEMALE":
		data.RetailerGender = "F"
	default:
		data.RetailerGender = "T"
	}

	fmt.Println(map[string]any{
		"mobile":      data.RetailerPhone,
		"name":        data.RetailerName,
		"gender":      data.RetailerGender,
		"email":       data.RetailerEmail,
		"pan":         data.RetailerPAN,
		"aadhaar":     data.RetailerAadhaar,
		"dateOfBirth": data.RetailerDateOfBirth,
		"address": map[string]string{
			"full":    data.RetailerAddress,
			"city":    data.RetailerCity,
			"pincode": data.RetailerPincode,
		},
		"latitude":  data.Latitude,
		"longitude": data.Longitude,
	})

	if err := utils.PostRequest2(
		utils.PayntricAPI+utils.AEPSSubMerchantSignup,
		"token",
		utils.PayntricAPIToken,
		"username",
		utils.PayntricUsername,
		map[string]any{
			"mobile":      data.RetailerPhone,
			"name":        data.RetailerName,
			"gender":      data.RetailerGender,
			"email":       data.RetailerEmail,
			"pan":         data.RetailerPAN,
			"aadhaar":     data.RetailerAadhaar,
			"dateOfBirth": data.RetailerDateOfBirth,
			"address": map[string]string{
				"full":    data.RetailerAddress,
				"city":    data.RetailerCity,
				"pincode": data.RetailerPincode,
			},
			"latitude":  data.Latitude,
			"longitude": data.Longitude,
		},
		&res,
	); err != nil {
		fmt.Println(res)
		return nil, err
	}

	if res.Status == "FAILED" || res.Status == "FAILURE" || res.Status == "Failure" {
		fmt.Println(res)
		return nil, errors.New(res.Message)
	}

	return &res, nil
}

func (ah *AEPSOnboardingHandler) HandleAEPSCheckMerchantEKYC(w http.ResponseWriter, r *http.Request) {
	userId, err := utils.ReadParamID(r)
	if err != nil {
		utils.BadRequest(w, ah.logger, "aeps check merchant ekyc", err)
		return
	}

	var BankDetails struct {
		BankCode string `json:"bank_code"`
	}

	if err := json.NewDecoder(r.Body).Decode(&BankDetails); err != nil {
		utils.BadRequest(w, ah.logger, "aeps check merchant ekyc", err)
		return
	}

	subMerchantId, err := ah.AEPSOnboardingStore.GetSubMerchantIDForAEPSEKYCCheck(userId)
	if err != nil {
		utils.ServerError(w, ah.logger, "aeps check merchant ekyc", err)
		return
	}

	res, err := aepsCheckMerchantEKYC(subMerchantId, BankDetails.BankCode)
	if err != nil {
		utils.BadRequest(w, ah.logger, "aeps check merchant ekyc", err)
		return
	}

	if err := ah.AEPSOnboardingStore.MerchantAEPSEKYCCheck(userId, res); err != nil {
		utils.ServerError(w, ah.logger, "aeps check merchant ekyc", err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"message": "AEPS merchant eKyc check completed successfully", "api_response": res})
}

func aepsCheckMerchantEKYC(subMerchantId string, bankCode string) (*models.AEPSOnboardingeKycCheckResponseModel, error) {
	var reqJson = make(map[string]any)
	reqJson["subMerchantId"] = subMerchantId
	reqJson["spKey"] = "WAP"
	if bankCode != "" {
		reqJson["gw"] = bankCode
	}
	var res models.AEPSOnboardingeKycCheckResponseModel
	if err := utils.PostRequest2(
		utils.PayntricAPI+utils.AEPSMerchantEKYCStatusCheck,
		"token",
		utils.PayntricAPIToken,
		"username",
		utils.PayntricUsername,
		reqJson,
		res,
	); err != nil {
		return nil, err
	}

	if res.Status == "FAILED" || res.Status == "FAILURE" || res.Status == "Failure" {
		return nil, errors.New(res.Message)
	}

	return &res, nil
}

func (ah *AEPSOnboardingHandler) HandleBiometricKYC(w http.ResponseWriter, r *http.Request) {
	userId, err := utils.ReadParamID(r)
	if err != nil {
		utils.BadRequest(w, ah.logger, "aeps biometric kyc", err)
		return
	}

	var biometricData models.AEPSBiometricDataModel
	if err := json.NewDecoder(r.Body).Decode(&biometricData); err != nil {
		utils.BadRequest(w, ah.logger, "aeps biometric kyc", err)
		return
	}

	var res models.AEPSOnboardingBiometricKYCRequestModel
	res.BiometricData = biometricData
	if err := ah.AEPSOnboardingStore.GetMerchantDetailsForBiometricKYC(userId, &res); err != nil {
		utils.ServerError(w, ah.logger, "aeps biometric kyc", err)
		return
	}

	apiRes, err := aepsBiometricKYC(&res)
	if err != nil {
		utils.BadRequest(w, ah.logger, "aeps biometric kyc", err)
		return
	}

	if err := ah.AEPSOnboardingStore.MerchantAEPSEKYCCheck(userId, &models.AEPSOnboardingeKycCheckResponseModel{
		EKYCStatus: apiRes.EKYCStatus,
		EKYCAction: apiRes.EKYCAction,
	}); err != nil {
		utils.BadRequest(w, ah.logger, "aeps biometric kyc", err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"message": "AEPS biometric verification completed", "api_response": apiRes})
}

func aepsBiometricKYC(data *models.AEPSOnboardingBiometricKYCRequestModel) (*models.AEPSOnboardingBiometricKYCResponseModel, error) {
	var res models.AEPSOnboardingBiometricKYCResponseModel
	if err := utils.PostRequest2(
		utils.PayntricAPI+utils.AEPSBiometricKYC,
		"token",
		utils.PayntricAPIToken,
		"username",
		utils.PayntricUsername,
		map[string]any{
			"subMerchantId": data.SubMerchantID,
			"referenceKey":  data.ReferenceKey,
			"latitude":      data.Latitude,
			"longitude":     data.Longitude,
			"externalRef":   uuid.NewString(),
			"captureType":   "FMR",
			"biometricData": data.BiometricData,
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
