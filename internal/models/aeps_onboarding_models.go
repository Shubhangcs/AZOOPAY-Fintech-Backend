package models

import (
	"time"
)

type CreateAEPSApplicationRequestModel struct {
	RetailerID      string `json:"retailer_id"`
	Latitude        string `json:"latitude"`
	Longitude       string `json:"longitude"`
	RetailerRemarks string `json:"retailer_remarks"`
}

type ChangeAEPSApplicationStatusRequestModel struct {
	RetailerID            string `json:"retailer_id"`
	AEPSApplicationStatus string `json:"aeps_application_status"`
	AdminRemarks          string `json:"admin_remarks"`
	RetailerRemarks       string `json:"retailer_remarks"`
}

type AEPSApplicationResponseModel struct {
	AEPSApplicationID     string `json:"aeps_application_id"`
	RetailerID            string `json:"retailer_id"`
	AEPSApplicationStatus string `json:"aeps_application_status"`
	RetailerRemarks       string `json:"retailer_remarks"`
	AdminRemarks          string `json:"admin_remarks"`
	Latitude              string `json:"latitude"`
	Longitude             string `json:"longitude"`
	RetailerDetails       struct {
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

type CreateAEPSMerchantResponseModel struct {
	Status            string `json:"status"`
	StatusCode        string `json:"statusCode"`
	Message           string `json:"message"`
	SubMerchantID     string `json:"subMerchantId"`
	ParentMerchantID  string `json:"parentMerchantId"`
	OutletID          string `json:"outletId"`
	MinKYCStatus      string `json:"minKycStatus"`
	EKYCStatus        string `json:"eKycStatus"`
	MobileChangeState string `json:"mobileChangeState"`
	IPayUUID          string `json:"ipayUuid"`
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

type UpdateAEPSMerchantResponseModel struct {
	Status       string `json:"status"`
	StatusCode   string `json:"statusCode"`
	Message      string `json:"message"`
	EKYCStatus   string `json:"ekycStatus"`
	EKYCAction   string `json:"ekycAction"`
	ReferenceKey string `json:"referenceKey,omitempty"`
	Data         struct {
		Status                  string `json:"status"`
		PidOptionWadh           string `json:"pidOptionWadh"`
		IsFaceAuthAvailable     bool   `json:"isFaceAuthAvailable,omitempty"`
		IsBiometricKycManditory bool   `json:"isBiometricKycManditory,omitempty"`
		BankName                string `json:"bankName"`
	} `json:"data"`
}

type AEPSMerchantDetailsResponseModel struct {
	AEPSMerchantID          string `json:"aeps_merchant_id"`
	RetailerID              string `json:"retailer_id"`
	SubMerchantID           string `json:"sub_merchant_id"`
	ParentMerchantID        string `json:"parent_merchant_id"`
	OutletID                string `json:"outlet_id"`
	MinKYCStatus            string `json:"min_kyc_status"`
	EKYCStatus              string `json:"ekyc_status"`
	EKYCAction              string `json:"ekyc_action"`
	ReferenceKey            string `json:"reference_key"`
	Status                  string `json:"status"`
	IsFaceAuthAvailable     bool   `json:"is_face_auth_available"`
	IsBiometricKycManditory bool   `json:"is_biometric_kyc_manditory"`
	BankName                string `json:"bank_name"`
	IsMerchantBlocked       bool   `json:"is_merchant_blocked"`
	PidOptionWadh           string `json:"pid_option_wadh"`
}
