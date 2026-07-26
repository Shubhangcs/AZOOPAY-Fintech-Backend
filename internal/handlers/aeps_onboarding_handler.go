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

func (ah *AEPSOnboardingHandler) HandleCreateAEPSApplication(w http.ResponseWriter, r *http.Request) {
	var req models.CreateAEPSApplicationRequestModel
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.BadRequest(w, ah.logger, "create aeps application", err)
		return
	}

	if err := ah.AEPSOnboardingStore.CreateAEPSApplication(&req); err != nil {
		utils.ServerError(w, ah.logger, "create aeps application", err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"message": "aeps application created successfully"})
}

func (ah *AEPSOnboardingHandler) HandleChangeAEPSApplicationStatus(w http.ResponseWriter, r *http.Request) {
	var req models.ChangeAEPSApplicationStatusRequestModel
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.BadRequest(w, ah.logger, "change aeps application status", err)
		return
	}

	if err := ah.AEPSOnboardingStore.ChangeAEPSApplicationStatus(&req); err != nil {
		utils.ServerError(w, ah.logger, "change aeps application status", err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"message": "aeps application status updated successfully"})
}

func (ah *AEPSOnboardingHandler) HandleSignupAEPSMerchant(w http.ResponseWriter, r *http.Request) {
	retailerId, err := utils.ReadParamID(r)
	if err != nil {
		utils.BadRequest(w, ah.logger, "signup aeps merchant", err)
		return
	}

	details, err := ah.AEPSOnboardingStore.GetAEPSApplicationByRetailerID(retailerId)
	if err != nil {
		utils.ServerError(w, ah.logger, "signup aeps merchant", err)
		return
	}

	apiRes, err := aepsMerchantSignup(details)
	if err != nil {
		utils.BadRequest(w, ah.logger, "signup aeps merchant", err)
		return
	}

	if err := ah.AEPSOnboardingStore.CreateAEPSMerchant(retailerId, apiRes); err != nil {
		utils.ServerError(w, ah.logger, "signup aeps merchant", err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"message": "aeps merchant signup successfull", "api_response": apiRes})
}

func (ah *AEPSOnboardingHandler) HandleCheckEKYCRequired(w http.ResponseWriter, r *http.Request) {
	retailerId, err := utils.ReadParamID(r)
	if err != nil {
		utils.BadRequest(w, ah.logger, "check ekyc required", err)
		return
	}

	details, err := ah.AEPSOnboardingStore.GetAEPSMerchantDetailsByRetailerID(retailerId)
	if err != nil {
		utils.ServerError(w, ah.logger, "check ekyc required", err)
		return
	}

	res, err := aepsCheckMerchantEKYC(details.SubMerchantID)
	if err != nil {
		utils.BadRequest(w, ah.logger, "check ekyc required", err)
		return
	}

	if err := ah.AEPSOnboardingStore.UpdateAEPSMerchant(retailerId, res); err != nil {
		utils.ServerError(w, ah.logger, "check ekyc required", err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"message": "ekyc check successfull", "api_response": res})
}

func (ah *AEPSOnboardingHandler) HandleBiometricKYC(w http.ResponseWriter, r *http.Request) {
	retailerId, err := utils.ReadParamID(r)
	if err != nil {
		utils.BadRequest(w, ah.logger, "biometric kyc", err)
		return
	}

	var req struct {
		Latitude      string                        `json:"latitude"`
		Longitude     string                        `json:"longitude"`
		BiometricData models.AEPSBiometricDataModel `json:"biometric_data"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.BadRequest(w, ah.logger, "biometric kyc", err)
		return
	}

	details, err := ah.AEPSOnboardingStore.GetAEPSMerchantDetailsByRetailerID(retailerId)
	if err != nil {
		utils.ServerError(w, ah.logger, "biometric kyc", err)
		return
	}

	res, err := aepsBiometricKYC(&req.BiometricData, details, req.Latitude, req.Longitude)
	if err != nil {
		utils.BadRequest(w, ah.logger, "biometric kyc", err)
		return
	}

	if err := ah.AEPSOnboardingStore.UpdateAEPSMerchant(retailerId, res); err != nil {
		utils.ServerError(w, ah.logger, "biometric kyc", err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"message": "biometric kyc updated successfully", "api_response": res})
}

func (ah *AEPSOnboardingHandler) HandleGetAllAEPSApplications(w http.ResponseWriter, r *http.Request) {
	res, err := ah.AEPSOnboardingStore.GetAllAEPSApplications()
	if err != nil {
		utils.ServerError(w, ah.logger, "get all aeps applications", err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"message": "all aeps application fetched successfully", "applications": res})
}

func (ah *AEPSOnboardingHandler) HandleGetAEPSApplicationByRetailerID(w http.ResponseWriter, r *http.Request) {
	retailerId, err := utils.ReadParamID(r)
	if err != nil {
		utils.BadRequest(w, ah.logger, "get aeps application", err)
		return
	}

	res, err := ah.AEPSOnboardingStore.GetAEPSApplicationByRetailerID(retailerId)
	if err != nil {
		utils.ServerError(w, ah.logger, "get aeps application", err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"message": "aeps application fetched successfully", "application": res})
}

func (ah *AEPSOnboardingHandler) HandleGetAEPSMerchantDetails(w http.ResponseWriter, r *http.Request) {
	retailerId, err := utils.ReadParamID(r)
	if err != nil {
		utils.BadRequest(w, ah.logger, "get aeps merchant details", err)
		return
	}

	res, err := ah.AEPSOnboardingStore.GetAEPSMerchantDetailsByRetailerID(retailerId)
	if err != nil {
		utils.ServerError(w, ah.logger, "get aeps application", err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"message": "aeps merchant details fetched successfully", "merchants": res})
}

func aepsMerchantSignup(data *models.AEPSApplicationResponseModel) (*models.CreateAEPSMerchantResponseModel, error) {
	var res models.CreateAEPSMerchantResponseModel

	switch data.RetailerDetails.RetailerGender {
	case "MALE":
		data.RetailerDetails.RetailerGender = "M"
	case "FEMALE":
		data.RetailerDetails.RetailerGender = "F"
	default:
		data.RetailerDetails.RetailerGender = "T"
	}

	fmt.Println(map[string]any{
		"mobile":      data.RetailerDetails.RetailerPhone,
		"name":        data.RetailerDetails.RetailerName,
		"gender":      data.RetailerDetails.RetailerGender,
		"email":       data.RetailerDetails.RetailerEmail,
		"pan":         data.RetailerDetails.RetailerPanNumber,
		"aadhaar":     data.RetailerDetails.RetailerAadhaarNumber,
		"dateOfBirth": data.RetailerDetails.RetailerDateOfBirth,
		"address": map[string]string{
			"full":    data.RetailerDetails.RetailerFullAddress,
			"city":    data.RetailerDetails.RetailerCity,
			"pincode": data.RetailerDetails.RetailerPincode,
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
			"mobile":      data.RetailerDetails.RetailerPhone,
			"name":        data.RetailerDetails.RetailerName,
			"gender":      data.RetailerDetails.RetailerGender,
			"email":       data.RetailerDetails.RetailerEmail,
			"pan":         data.RetailerDetails.RetailerPanNumber,
			"aadhaar":     data.RetailerDetails.RetailerAadhaarNumber,
			"dateOfBirth": data.RetailerDetails.RetailerDateOfBirth,
			"address": map[string]string{
				"full":    data.RetailerDetails.RetailerFullAddress,
				"city":    data.RetailerDetails.RetailerCity,
				"pincode": data.RetailerDetails.RetailerPincode,
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

func aepsCheckMerchantEKYC(subMerchantId string) (*models.UpdateAEPSMerchantResponseModel, error) {
	var reqJson = make(map[string]any)
	reqJson["subMerchantId"] = subMerchantId
	reqJson["spKey"] = "WAP"
	var res models.UpdateAEPSMerchantResponseModel
	if err := utils.PostRequest2(
		utils.PayntricAPI+utils.AEPSMerchantEKYCStatusCheck,
		"token",
		utils.PayntricAPIToken,
		"username",
		utils.PayntricUsername,
		reqJson,
		&res,
	); err != nil {
		return nil, err
	}

	if res.Status == "FAILED" || res.Status == "FAILURE" || res.Status == "Failure" {
		return nil, errors.New(res.Message)
	}

	return &res, nil
}

func aepsBiometricKYC(bio *models.AEPSBiometricDataModel, data *models.AEPSMerchantDetailsResponseModel, lat, lon string) (*models.UpdateAEPSMerchantResponseModel, error) {
	var res models.UpdateAEPSMerchantResponseModel
	biometricData := map[string]any{
		"encryptedAadhaar": bio.EncryptedAadhaar,
		"pidData":          bio.PIDData,
		"pidDataType":      bio.PIDDataType,
		"dc":               bio.DeviceCode,
		"dpId":             bio.DeviceProviderID,
		"rdsId":            bio.RegisteredDevicesServiceID,
		"rdsVer":           bio.RegisteredDeviceServiceVersion,
		"mi":               bio.ModelIdentifier,
		"mc":               bio.ModelCertificationCode,
		"ci":               bio.ModelCertificateExpiryDate,
		"sessionKey":       bio.SessionKey,
		"hmac":             bio.Hmac,
		"srno":             bio.SerialNumber,
		"sysid":            bio.SystemIdentifier,
		"ts":               bio.BiometricTimestamp,
		"nmPoints":         bio.NmPoints,
		"fCount":           bio.NumberOfFingerprintsCaptured,
		"fType":            bio.FingerType,
		"iCount":           bio.NumberOfIrisScanCaptured,
		"iType":            bio.IrisType,
		"pCount":           bio.NumberOfPhotosCaptured,
		"pType":            bio.PhotoType,
		"qScore":           bio.QualityScore,
		"errCode":          bio.ErrorCode,
		"errInfo":          bio.ErrorInfo,
	}
	if err := utils.PostRequest2(
		utils.PayntricAPI+utils.AEPSBiometricKYC,
		"token",
		utils.PayntricAPIToken,
		"username",
		utils.PayntricUsername,
		map[string]any{
			"subMerchantId": data.SubMerchantID,
			"referenceKey":  data.ReferenceKey,
			"latitude":      lat,
			"longitude":     lon,
			"externalRef":   uuid.NewString(),
			"captureType":   "finger",
			"biometricData": biometricData,
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
