package models

import "time"

type UPIATMCreateQRRequestModel struct {
	Amount float64 `json:"amount"`
	Mobile string  `json:"mobile"`
}

type UPIATMCreateQRAPIRequestModel struct {
	RequestData *UPIATMCreateQRRequestModel
	RetailerID  string
	RequestID   string
	Latitude    string
	Longitude   string
	OutletID    string
}

type UPIATMCreateQRAPIResponseModel struct {
	Status              string  `json:"status"`
	StatusCode          string  `json:"status_code"`
	Message             string  `json:"message"`
	TransactionID       string  `json:"transactionId"`
	RequestID           string  `json:"requestId"`
	QRStatus            string  `json:"qrStatus"`
	IpayID              string  `json:"ipayId"`
	QRString            string  `json:"qrString"`
	QRMobile            string  `json:"qrMobile"`
	ExpiryDate          string  `json:"expiryDt"`
	QRCreatedDate       string  `json:"qrCreatedDt"`
	DisplayExpirySecond int     `json:"displayExpirySec"`
	Amount              float64 `json:"amount"`
	PayableValue        float64 `json:"payableValue"`
	TransactionValue    float64 `json:"transactionValue"`
	CommissionAmount    float64 `json:"commissionAmount"`
	TDSAmount           float64 `json:"tdsAmount"`
	NETAmount           float64 `json:"netAmount"`
	SettlementStatus    string  `json:"settlementStatus"`
}

type UPIATMCheckQRTransactionStatusAPIRequestModel struct {
	TransactionID string `json:"transactionId"`
	IpayID        string `json:"ipayId"`
	RequestID     string `json:"requestId"`
	OutletID      string `json:"outletId"`
}

type UPIATMCheckQRTransactionStatusAPIResponseModel struct {
	Status           string  `json:"status"`
	StatusCode       string  `json:"statusCode"`
	Message          string  `json:"message"`
	TransactionID    string  `json:"transactionId"`
	RequestID        string  `json:"requestId"`
	QRStatus         string  `json:"qrStatus"`
	IpayID           string  `json:"ipayId"`
	Amount           float64 `json:"amount"`
	SettlementStatus string  `json:"settlementStatus"`
}

type UPIATMTransactionResponseModel struct {
	RetailerID       string    `json:"retailer_id"`
	RetailerName     string    `json:"retailer_name"`
	TransactionID    string    `json:"transactionId"`
	RequestID        string    `json:"requestId"`
	QRStatus         string    `json:"qrStatus"`
	IpayID           string    `json:"ipayId"`
	Amount           float64   `json:"amount"`
	PayableValue     float64   `json:"payableValue"`
	TransactionValue float64   `json:"transactionValue"`
	CommissionAmount float64   `json:"commissionAmount"`
	TDSAmount        float64   `json:"tdsAmount"`
	NETAmount        float64   `json:"netAmount"`
	SettlementStatus string    `json:"settlementStatus"`
	CreatedAT        time.Time `json:"created_at"`
	UpdatedAT        time.Time `json:"updated_at"`
}
