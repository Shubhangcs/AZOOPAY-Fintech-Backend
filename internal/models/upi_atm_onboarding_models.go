package models

import (
	"time"
)

type CreateUPIATMApplicationRequestModel struct {
	RetailerID      string `json:"retailer_id"`
	Latitude        string `json:"latitude"`
	Longitude       string `json:"longitude"`
	RetailerRemarks string `json:"retailer_remarks"`
}

type ChangeUPIATMApplicationStatusRequestModel struct {
	RetailerID              string `json:"retailer_id"`
	UPIATMApplicationStatus string `json:"upi_atm_application_status"`
	AdminRemarks            string `json:"admin_remarks"`
	RetailerRemarks         string `json:"retailer_remarks"`
}

type UPIATMApplicationResponseModel struct {
	UPIATMApplicationID     string `json:"upi_atm_application_id"`
	RetailerID              string `json:"retailer_id"`
	UPIATMApplicationStatus string `json:"upi_atm_application_status"`
	RetailerRemarks         string `json:"retailer_remarks"`
	AdminRemarks            string `json:"admin_remarks"`
	Latitude                string `json:"latitude"`
	Longitude               string `json:"longitude"`
	RetailerDetails         struct {
		RetailerName          string    `json:"retailer_name"`
		RetailerPhone         string    `json:"retailer_phone"`
		RetailerEmail         string    `json:"retailer_email"`
		RetailerAadhaarNumber string    `json:"retailer_aadhaar_number"`
		RetailerPanNumber     string    `json:"retailer_pan_number"`
		RetailerFullAddress   string    `json:"retailer_full_address"`
		RetailerCity          string    `json:"retailer_city"`
		RetailerPincode       string    `json:"retailer_pincode"`
		RetailerDateOfBirth   time.Time `json:"retailer_date_of_birth"`
		RetailerGender        string    `json:"retailer_gender"`
	} `json:"retailer_details"`
	CreatedAT time.Time `json:"created_at"`
	UpdatedAT time.Time `json:"updated_at"`
}

type CreateUPIATMMerchantResponseModel struct {
	Status            string            `json:"status"`
	StatusCode        string            `json:"statusCode"`
	Message           string            `json:"message"`
	SubMerchantID     string            `json:"subMerchantId"`
	ParentMerchantID  string            `json:"parentMerchantId"`
	OutletID          string            `json:"outletId"`
	MinKYCStatus      string            `json:"minKycStatus"`
	EKYCStatus        string            `json:"ekycStatus"`
	MobileChangeState string            `json:"mobileChangeState"`
	IPayUUID          string            `json:"ipayUuid"`
	Timestamp         string            `json:"timestamp"`
	Errors            map[string]string `json:"errors,omitempty"`
	MerchantData      struct {
		Name        string `json:"name"`
		DateOfBirth string `json:"dateOfBirth"`
		Gender      string `json:"gender"`
		Pincode     string `json:"pincode"`
		State       string `json:"state"`
		City        string `json:"city"`
		Address     string `json:"address"`
	} `json:"data"`
}

type UpdateUPIATMMerchantResponseModel struct {
	Status        string `json:"status"`
	StatusCode    string `json:"statusCode"`
	Message       string `json:"message"`
	SubMerchantID string `json:"subMerchantId"`
	OutletID      string `json:"outletId"`
	MinKYCStatus  string `json:"minKycStatus"`
	EKYCStatus    string `json:"ekycStatus"`
	EKYCAction    string `json:"ekycAction"`
	ReferenceKey  string `json:"referenceKey,omitempty"`
	Data          struct {
		Status                  string `json:"status"`
		IsFaceAuthAvailable     *bool  `json:"isFaceAuthAvailable,omitempty"`
		IsBiometricKycMandatory *bool  `json:"isBiometricKycMandatory,omitempty"`
		BankName                string `json:"bankName"`
	} `json:"data"`
}

type UPIATMMerchantDetailsResponseModel struct {
	UPIATMMerchantID        string    `json:"upi_atm_merchant_id"`
	RetailerID              string    `json:"retailer_id"`
	SubMerchantID           string    `json:"sub_merchant_id"`
	ParentMerchantID        string    `json:"parent_merchant_id"`
	OutletID                string    `json:"outlet_id"`
	MinKYCStatus            string    `json:"min_kyc_status"`
	EKYCStatus              string    `json:"ekyc_status"`
	EKYCAction              string    `json:"ekyc_action"`
	ReferenceKey            string    `json:"reference_key"`
	Status                  string    `json:"status"`
	IsFaceAuthAvailable     bool      `json:"is_face_auth_available"`
	IsBiometricKycMandatory bool      `json:"is_biometric_kyc_mandatory"`
	BankName                string    `json:"bank_name"`
	MobileChangeState       string    `json:"mobile_change_state"`
	IPayUUID                string    `json:"ipay_uuid"`
	IsMerchantBlocked       bool      `json:"is_merchant_blocked"`
	CreatedAT               time.Time `json:"created_at"`
	UpdatedAT               time.Time `json:"updated_at"`
}

type UPIATMBiometricKYCRequestModel struct {
	Latitude      string                 `json:"latitude"`
	Longitude     string                 `json:"longitude"`
	CaptureType   string                 `json:"capture_type"`
	BiometricData AEPSBiometricDataModel `json:"biometric_data"`
}
