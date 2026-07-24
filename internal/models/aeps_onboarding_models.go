package models

import (
	"errors"
	"time"
)

type ApplyForAEPSRequestModel struct {
	AEPSApplicationID     string    `json:"aeps_application_id,omitempty"`
	RetailerID            string    `json:"retailer_id,omitempty"`
	RetailerName          string    `json:"retailer_name,omitempty"`
	AEPSApplicationStatus string    `json:"aeps_application_status,omitempty"`
	Remarks               string    `json:"remarks,omitempty"`
	Latitude              string    `json:"latitude,omitempty"`
	Longitude             string    `json:"longitude,omitempty"`
	CreatedAT             time.Time `json:"created_at"`
	UpdatedAT             time.Time `json:"updated_at"`
}

func (st *ApplyForAEPSRequestModel) Validate() error {
	if st.RetailerID == "" {
		return errors.New("retailer id missing")
	}

	if st.RetailerName == "" {
		return errors.New("retailer name is missing")
	}

	if st.Latitude == "" {
		return errors.New("latitude is missing")
	}

	if st.Longitude == "" {
		return errors.New("longitude is missing")
	}

	return nil
}

type AEPSApplicationResponseModel struct {
	AEPSApplicationID     string    `json:"aeps_application_id,omitempty"`
	RetailerID            string    `json:"retailer_id,omitempty"`
	RetailerName          string    `json:"retailer_name,omitempty"`
	RetailerPAN           string    `json:"retailer_pan,omitempty"`
	RetailerAadhaar       string    `json:"retailer_aadhaar,omitempty"`
	RetailerPhone         string    `json:"retailer_phone,omitempty"`
	RetailerEmail         string    `json:"retailer_email,omitempty"`
	RetailerDateOfBirth   time.Time `json:"retailer_date_of_birth,omitempty"`
	RetailerGender        string    `json:"retailer_gender,omitempty"`
	RetailerAddress       string    `json:"retailer_address,omitempty"`
	RetailerCity          string    `json:"retailer_city,omitempty"`
	RetailerPincode       string    `json:"retailer_pincode,omitempty"`
	AEPSApplicationStatus string    `json:"aeps_application_status,omitempty"`
	Remarks               string    `json:"remarks,omitempty"`
	Latitude              string    `json:"latitude,omitempty"`
	Longitude             string    `json:"longitude,omitempty"`
	CreatedAT             time.Time `json:"created_at,omitempty"`
	UpdatedAT             time.Time `json:"updated_at,omitempty"`
}

type AEPSOnboardingSubMerchantSignupRequestModel struct {
	Mobile      string `json:"mobile"`
	Name        string `json:"name"`
	Gender      string `json:"gender"`
	Pan         string `json:"pan"`
	Email       string `json:"email"`
	Aadhaar     string `json:"aadhaar"`
	DateOfBirth string `json:"dateOfBirth"`
	Latitude    string `json:"latitude"`
	Longitude   string `json:"longitude"`
	Address     struct {
		Full    string `json:"full"`
		City    string `json:"city"`
		Pincode string `json:"pincode"`
	} `json:"address"`
}

type AEPSOnboardingSubMerchantSignupSuccessResponseModel struct {
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
	Timestamp         string `json:"timestamp"`
	MerchantData      struct {
		OutletID    string  `json:"outletId"`
		Name        string  `json:"name"`
		DateOfBirth string  `json:"dateOfBirth"`
		Gender      string  `json:"gender"`
		Pincode     string  `json:"pincode"`
		State       string  `json:"state"`
		City        string  `json:"city"`
		Address     string  `json:"address"`
		ProfilePic  *string `json:"profilePic"`
	} `json:"data"`
}

type AEPSOnboardingeKycCheckRequestModel struct {
	SubMerchantID      string `json:"subMerchantId"`
	ServiceProviderKey string `json:"spKey"`
	SpecialCode        string `json:"gw,omitempty"`
}

type AEPSOnboardingeKycCheckResponseModel struct {
	Status        string `json:"status"`
	StatusCode    string `json:"statusCode"`
	Message       string `json:"message"`
	SubMerchantID string `json:"subMerchantId"`
	OutletID      string `json:"outletId"`
	MinKycStatus  string `json:"minKycStatus"`
	EKYCStatus    string `json:"ekycStatus"`
	EKYCAction    string `json:"ekycAction"`
	ReferenceKey  string `json:"referenceKey,omitempty"`
	Data          struct {
		Status                  string `json:"status"`
		IsFaceAuthAvailable     bool   `json:"isFaceAuthAvailable,omitempty"`
		IsBiometricKycManditory bool   `json:"isBiometricKycManditory,omitempty"`
		BankName                string `json:"bankName"`
	} `json:"data"`
}

type AEPSOnboardingBiometricKYCRequestModel struct {
	SubMerchantID     string                 `json:"subMerchantId"`
	ReferenceKey      string                 `json:"referenceKey"`
	Latitude          string                 `json:"latitude"`
	Longitude         string                 `json:"longitude"`
	ExternalReference string                 `json:"externalRef"`
	CaptureType       string                 `json:"captureType"`
	BiometricData     AEPSBiometricDataModel `json:"biometricData"`
}

type AEPSOnboardingBiometricKYCResponseModel struct {
	Status        string `json:"status"`
	StatusCode    string `json:"statusCode"`
	Message       string `json:"message"`
	SubMerchantID string `json:"subMerchantId"`
	OutletID      string `json:"outletId"`
	MinKYCStatus  string `json:"minKycStatus"`
	EKYCStatus    string `json:"ekycStatus"`
	EKYCAction    string `json:"ekycAction"`
}

type AEPSOnboardingErrorResponse struct {
	Status     string `json:"status"`
	StatusCode string `json:"statusCode"`
	Message    string `json:"message"`
	Errors     struct {
		Mobile string `json:"mobile,omitempty"`
	} `json:"errors"`
}
