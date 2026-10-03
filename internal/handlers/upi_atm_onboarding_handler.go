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

type UPIATMOnboardingHandler struct {
	logger                *slog.Logger
	UPIATMOnboardingStore store.UPIATMOnboardingStore
}

func NewUPIATMOnboardingHandler(logger *slog.Logger, UPIATMOnboardingStore store.UPIATMOnboardingStore) *UPIATMOnboardingHandler {
	return &UPIATMOnboardingHandler{
		logger,
		UPIATMOnboardingStore,
	}
}

func (uh *UPIATMOnboardingHandler) HandleCreateUPIATMApplication(w http.ResponseWriter, r *http.Request) {
	var req models.CreateUPIATMApplicationRequestModel
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.BadRequest(w, uh.logger, "create upi atm application", err)
		return
	}

	if err := uh.UPIATMOnboardingStore.CreateUPIATMApplication(&req); err != nil {
		utils.ServerError(w, uh.logger, "create upi atm application", err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"message": "upi atm application created successfully"})
}

func (uh *UPIATMOnboardingHandler) HandleChangeUPIATMApplicationStatus(w http.ResponseWriter, r *http.Request) {
	var req models.ChangeUPIATMApplicationStatusRequestModel
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.BadRequest(w, uh.logger, "change upi atm application status", err)
		return
	}

	if err := uh.UPIATMOnboardingStore.ChangeUPIATMApplicationStatus(&req); err != nil {
		utils.ServerError(w, uh.logger, "change upi atm application status", err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"message": "upi atm application status updated successfully"})
}

func (uh *UPIATMOnboardingHandler) HandleSignupUPIATMMerchant(w http.ResponseWriter, r *http.Request) {
	retailerId, err := utils.ReadParamID(r)
	if err != nil {
		utils.BadRequest(w, uh.logger, "signup upi atm merchant", err)
		return
	}

	details, err := uh.UPIATMOnboardingStore.GetUPIATMApplicationByRetailerID(retailerId)
	if err != nil {
		utils.ServerError(w, uh.logger, "signup upi atm merchant", err)
		return
	}

	apiRes, err := upiAtmMerchantSignup(details)
	if err != nil {
		utils.BadRequest(w, uh.logger, "signup upi atm merchant", err)
		return
	}

	if err := uh.UPIATMOnboardingStore.CreateUPIATMMerchant(retailerId, apiRes); err != nil {
		utils.ServerError(w, uh.logger, "signup upi atm merchant", err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"message": "upi atm merchant signup successfull", "api_response": apiRes})
}

func (uh *UPIATMOnboardingHandler) HandleCheckEKYCStatus(w http.ResponseWriter, r *http.Request) {
	retailerId, err := utils.ReadParamID(r)
	if err != nil {
		utils.BadRequest(w, uh.logger, "check upi atm ekyc status", err)
		return
	}

	details, err := uh.UPIATMOnboardingStore.GetUPIATMMerchantDetailsByRetailerID(retailerId)
	if err != nil {
		utils.ServerError(w, uh.logger, "check upi atm ekyc status", err)
		return
	}

	// gw: AC (Airtel Payments Bank), NA (NSDL Bank), JA (Jio Payments Bank)
	gw := r.URL.Query().Get("gw")
	if gw == "" {
		gw = "AC"
	}

	res, err := upiAtmCheckMerchantEKYC(details.SubMerchantID, gw)
	if err != nil {
		utils.BadRequest(w, uh.logger, "check upi atm ekyc status", err)
		return
	}

	if err := uh.UPIATMOnboardingStore.UpdateUPIATMMerchant(retailerId, res); err != nil {
		utils.ServerError(w, uh.logger, "check upi atm ekyc status", err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"message": "ekyc check successfull", "api_response": res})
}

func (uh *UPIATMOnboardingHandler) HandleBiometricKYC(w http.ResponseWriter, r *http.Request) {
	retailerId, err := utils.ReadParamID(r)
	if err != nil {
		utils.BadRequest(w, uh.logger, "upi atm biometric kyc", err)
		return
	}

	var req models.UPIATMBiometricKYCRequestModel
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.BadRequest(w, uh.logger, "upi atm biometric kyc", err)
		return
	}

	details, err := uh.UPIATMOnboardingStore.GetUPIATMMerchantDetailsByRetailerID(retailerId)
	if err != nil {
		utils.ServerError(w, uh.logger, "upi atm biometric kyc", err)
		return
	}

	if details.ReferenceKey == "" {
		utils.BadRequest(w, uh.logger, "upi atm biometric kyc", errors.New("reference key not found, check ekyc status first"))
		return
	}

	res, err := upiAtmBiometricKYC(&req, details)
	if err != nil {
		utils.BadRequest(w, uh.logger, "upi atm biometric kyc", err)
		return
	}

	if err := uh.UPIATMOnboardingStore.UpdateUPIATMMerchant(retailerId, res); err != nil {
		utils.ServerError(w, uh.logger, "upi atm biometric kyc", err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"message": "biometric kyc submitted successfully", "api_response": res})
}

func (uh *UPIATMOnboardingHandler) HandleGetAllUPIATMApplications(w http.ResponseWriter, r *http.Request) {
	res, err := uh.UPIATMOnboardingStore.GetAllUPIATMApplications()
	if err != nil {
		utils.ServerError(w, uh.logger, "get all upi atm applications", err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"message": "all upi atm applications fetched successfully", "applications": res})
}

func (uh *UPIATMOnboardingHandler) HandleGetUPIATMApplicationByRetailerID(w http.ResponseWriter, r *http.Request) {
	retailerId, err := utils.ReadParamID(r)
	if err != nil {
		utils.BadRequest(w, uh.logger, "get upi atm application", err)
		return
	}

	res, err := uh.UPIATMOnboardingStore.GetUPIATMApplicationByRetailerID(retailerId)
	if err != nil {
		utils.ServerError(w, uh.logger, "get upi atm application", err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"message": "upi atm application fetched successfully", "application": res})
}

func (uh *UPIATMOnboardingHandler) HandleGetUPIATMMerchantDetails(w http.ResponseWriter, r *http.Request) {
	retailerId, err := utils.ReadParamID(r)
	if err != nil {
		utils.BadRequest(w, uh.logger, "get upi atm merchant details", err)
		return
	}

	res, err := uh.UPIATMOnboardingStore.GetUPIATMMerchantDetailsByRetailerID(retailerId)
	if err != nil {
		utils.ServerError(w, uh.logger, "get upi atm merchant details", err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"message": "upi atm merchant details fetched successfully", "merchant": res})
}

func (uh *UPIATMOnboardingHandler) HandleGetAllUPIATMMerchants(w http.ResponseWriter, r *http.Request) {
	res, err := uh.UPIATMOnboardingStore.GetAllUPIATMMerchants()
	if err != nil {
		utils.ServerError(w, uh.logger, "get all upi atm merchants", err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"message": "upi atm merchants fetched successfully", "merchants": res})
}

func (uh *UPIATMOnboardingHandler) HandleChangeUPIATMMerchantBlockStatus(w http.ResponseWriter, r *http.Request) {
	retailerId, err := utils.ReadParamID(r)
	if err != nil {
		utils.BadRequest(w, uh.logger, "change upi atm merchant block status", err)
		return
	}

	var req struct {
		IsMerchantBlocked bool `json:"is_merchant_blocked"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.BadRequest(w, uh.logger, "change upi atm merchant block status", err)
		return
	}

	if err := uh.UPIATMOnboardingStore.ChangeUPIATMMerchantBlockStatus(retailerId, req.IsMerchantBlocked); err != nil {
		utils.ServerError(w, uh.logger, "change upi atm merchant block status", err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"message": "upi atm merchant block status updated successfully"})
}

func isUPIATMAPIFailure(status string) bool {
	return status == "FAILED" || status == "FAILURE" || status == "Failure"
}

func upiAtmMerchantSignup(data *models.UPIATMApplicationResponseModel) (*models.CreateUPIATMMerchantResponseModel, error) {
	var res models.CreateUPIATMMerchantResponseModel

	// Payntric accepts M / F / O
	var gender string
	switch data.RetailerDetails.RetailerGender {
	case "MALE", "M":
		gender = "M"
	case "FEMALE", "F":
		gender = "F"
	default:
		gender = "O"
	}

	if err := utils.PostRequest2(
		utils.PayntricAPI+utils.UPIATMSubMerchantSignup,
		"username",
		utils.PayntricUsername,
		"token",
		utils.PayntricAPIToken,
		map[string]any{
			"mobile":      data.RetailerDetails.RetailerPhone,
			"name":        data.RetailerDetails.RetailerName,
			"gender":      gender,
			"pan":         data.RetailerDetails.RetailerPanNumber,
			"email":       data.RetailerDetails.RetailerEmail,
			"aadhaar":     data.RetailerDetails.RetailerAadhaarNumber,
			"dateOfBirth": data.RetailerDetails.RetailerDateOfBirth.Format("2006-01-02"),
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
		return nil, err
	}

	if isUPIATMAPIFailure(res.Status) || res.SubMerchantID == "" {
		if res.Message == "" {
			return nil, errors.New("upi atm merchant signup failed")
		}
		return nil, errors.New(res.Message)
	}

	return &res, nil
}

func upiAtmCheckMerchantEKYC(subMerchantId, gw string) (*models.UpdateUPIATMMerchantResponseModel, error) {
	var res models.UpdateUPIATMMerchantResponseModel
	if err := utils.PostRequest2(
		utils.PayntricAPI+utils.UPIATMMerchantEKYCStatusCheck,
		"username",
		utils.PayntricUsername,
		"token",
		utils.PayntricAPIToken,
		map[string]any{
			"subMerchantId": subMerchantId,
			"spKey":         "WAP",
			"gw":            gw,
		},
		&res,
	); err != nil {
		return nil, err
	}

	if isUPIATMAPIFailure(res.Status) {
		return nil, errors.New(res.Message)
	}

	return &res, nil
}

func upiAtmBiometricKYC(req *models.UPIATMBiometricKYCRequestModel, data *models.UPIATMMerchantDetailsResponseModel) (*models.UpdateUPIATMMerchantResponseModel, error) {
	var res models.UpdateUPIATMMerchantResponseModel

	captureType := req.CaptureType
	if captureType == "" {
		captureType = "FMR"
	}

	bio := req.BiometricData
	biometricData := map[string]any{
		"encryptedAadhaar": bio.EncryptedAadhaar,
		"dc":               bio.DeviceCode,
		"dpId":             bio.DeviceProviderID,
		"mc":               bio.ModelCertificationCode,
		"mi":               bio.ModelIdentifier,
		"rdsId":            bio.RegisteredDevicesServiceID,
		"rdsVer":           bio.RegisteredDeviceServiceVersion,
		"ci":               bio.ModelCertificateExpiryDate,
		"sessionKey":       bio.SessionKey,
		"pidData":          bio.PIDData,
		"pidDataType":      bio.PIDDataType,
		"hmac":             bio.Hmac,
		"errCode":          bio.ErrorCode,
		"errInfo":          bio.ErrorInfo,
		"fCount":           bio.NumberOfFingerprintsCaptured,
		"fType":            bio.FingerType,
		"iCount":           bio.NumberOfIrisScanCaptured,
		"iType":            bio.IrisType,
		"pCount":           bio.NumberOfPhotosCaptured,
		"pType":            bio.PhotoType,
		"qScore":           bio.QualityScore,
		"nmPoints":         bio.NmPoints,
		"srno":             bio.SerialNumber,
		"sysid":            bio.SystemIdentifier,
		"ts":               bio.BiometricTimestamp,
	}

	if err := utils.PostRequest2(
		utils.PayntricAPI+utils.UPIATMBiometricKYC,
		"username",
		utils.PayntricUsername,
		"token",
		utils.PayntricAPIToken,
		map[string]any{
			"subMerchantId": data.SubMerchantID,
			"referenceKey":  data.ReferenceKey,
			"latitude":      req.Latitude,
			"longitude":     req.Longitude,
			"externalRef":   uuid.NewString(),
			"captureType":   captureType,
			"biometricData": biometricData,
		},
		&res,
	); err != nil {
		return nil, err
	}

	if isUPIATMAPIFailure(res.Status) {
		return nil, errors.New(res.Message)
	}

	return &res, nil
}
