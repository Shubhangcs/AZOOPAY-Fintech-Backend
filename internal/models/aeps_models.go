package models

import (
	"database/sql"
	"time"
)

type AEPSBiometricDataModel struct {
	EncryptedAadhaar               string `json:"encryptedAadhaar"`
	PIDData                        string `json:"pidData"`
	PIDDataType                    string `json:"pidDataType"`
	DeviceCode                     string `json:"dc"`
	DeviceProviderID               string `json:"dpId"`
	RegisteredDevicesServiceID     string `json:"rdsId"`
	RegisteredDeviceServiceVersion string `json:"rdsVer"`
	ModelIdentifier                string `json:"mi"`
	ModelCertificationCode         string `json:"mc"`
	ModelCertificateExpiryDate     string `json:"ci"`
	SessionKey                     string `json:"sessionKey"`
	Hmac                           string `json:"hmac"`
	SerialNumber                   string `json:"srno"`
	SystemIdentifier               string `json:"sysid"`
	BiometricTimestamp             string `json:"ts"`
	NmPoints                       string `json:"nmPoints"`
	NumberOfFingerprintsCaptured   string `json:"fCount"`
	FingerType                     string `json:"fType"`
	NumberOfIrisScanCaptured       string `json:"iCount"`
	IrisType                       string `json:"iType"`
	NumberOfPhotosCaptured         string `json:"pCount"`
	PhotoType                      string `json:"pType"`
	QualityScore                   string `json:"qScore"`
	ErrorCode                      string `json:"errCode"`
	ErrorInfo                      string `json:"errInfo"`
}

type AEPSDetailsModel struct {
	AadhaarNumber string `json:"aadhaar_number"`
	OutletID      string `json:"outlet_id"`
	Latitude      string `json:"latitude"`
	Longitude     string `json:"longitude"`
	ReferenceKey  string `json:"reference_key"`
}

type AEPSMiniStatementRequestModel struct {
	RetailerID    string                 `json:"retailer_id"`
	RequestID     string                 `json:"request_id"`
	BankID        string                 `json:"bankiin"`
	MobileNumber  string                 `json:"mobile"`
	Latitude      string                 `json:"latitude"`
	Longitude     string                 `json:"longitude"`
	CustomerName  string                 `json:"customer_name"`
	AadhaarNumber string                 `json:"aadhaar_number"`
	BiometricData AEPSBiometricDataModel `json:"biometric_data"`
}

type AEPSMiniStatementResponseModel struct {
	Status              string  `json:"status"`
	StatusCode          string  `json:"statusCode"`
	Message             string  `json:"message"`
	TransactionID       string  `json:"transactionId"`
	RequestID           string  `json:"requestId"`
	OutletID            string  `json:"outletId"`
	Operation           string  `json:"operation"`
	TransactionStatus   string  `json:"txnStatus"`
	IPayID              string  `json:"ipayId"`
	OrderID             string  `json:"orderId"`
	Amount              float64 `json:"amount"`
	BankName            string  `json:"bankName"`
	AccountNumber       string  `json:"accountNumber"`
	TransactionMode     string  `json:"transactionMode"`
	BankAccountBalance  float64 `json:"bankAccountBalance"`
	IsOnusTransaction   bool    `json:"isOnusTxn"`
	ExternalReferenceID string  `json:"externalRef"`
	MiniStatement       []any   `json:"miniStatemenet"`
	CommisionAmount     float64 `json:"commisionAmount"`
	TDSAmount           float64 `json:"tdsAmount"`
	GSTAmount           float64 `json:"gstAmount"`
	NETAmount           float64 `json:"netAmount"`
	BalanceAfter        float64 `json:"balanceAfter"`
	LedgerID            float64 `json:"ledgerId"`
	SettlementStatus    string  `json:"settlementStatus"`
	SettlementMode      string  `json:"settlementMode"`
}

type AEPSBalanceEnquiryRequestModel struct {
	RetailerID    string                 `json:"retailer_id"`
	RequestID     string                 `json:"request_id"`
	BankID        string                 `json:"bankiin"`
	MobileNumber  string                 `json:"mobile"`
	Latitude      string                 `json:"latitude"`
	Longitude     string                 `json:"longitude"`
	CustomerName  string                 `json:"customer_name"`
	AadhaarNumber string                 `json:"aadhaar_number"`
	BiometricData AEPSBiometricDataModel `json:"biometric_data"`
}

type AEPSDailyLoginRequestModel struct {
	Latitude      string                 `json:"latitude"`
	Longitude     string                 `json:"longitude"`
	BiometricData AEPSBiometricDataModel `json:"biometricData"`
}

type AEPSDailyLoginResponseModel struct {
	Status              string `json:"status"`
	StatusCode          string `json:"statusCode"`
	ActionCode          string `json:"actionCode"`
	Message             string `json:"message"`
	IPayID              string `json:"ipayId"`
	IPayUUID            string `json:"ipayUuid"`
	ExternalReferenceID string `json:"externalRef"`
	OutletID            string `json:"outletId"`
	ReferenceKey        string `json:"referenceKey"`
	// Data                struct {
	// 	AdditionalProp1 map[any]any `json:"additionalProp1"`
	// 	AdditionalProp2 map[any]any `json:"additionalProp2"`
	// 	AdditionalProp3 map[any]any `json:"additionalProp3"`
	// } `json:"data"`
}

type AEPSOutletLoginResponseModel struct {
	Status              string `json:"status"`
	StatusCode          string `json:"statusCode"`
	ActionCode          string `json:"actionCode"`
	Message             string `json:"message"`
	IPayID              string `json:"ipayId"`
	IPayUUID            string `json:"ipayUuid"`
	ExternalReferenceID string `json:"externalRef"`
	OutletID            string `json:"outletId"`
	// Data                struct {
	// 	AdditionalProp1 map[any]any `json:"additionalProp1"`
	// 	AdditionalProp2 map[any]any `json:"additionalProp2"`
	// 	AdditionalProp3 map[any]any `json:"additionalProp3"`
	// } `json:"data"`
}

type AEPSCashWithdrawalRequestModel struct {
	RequestID                string                 `json:"requestId"`
	BankIdentificationNumber string                 `json:"bankin"`
	Mobile                   string                 `json:"mobile"`
	Amount                   float64                `json:"amount"`
	Latitude                 string                 `json:"latitude"`
	Longitude                string                 `json:"longitude"`
	CustomerName             string                 `json:"customerName"`
	Aadhaar                  string                 `json:"aadhaar"`
	BiometricData            AEPSBiometricDataModel `json:"biometricData"`
}

type AEPSCashWithdrawalResponseModel struct {
	Status              string  `json:"status"`
	StatusCode          string  `json:"statusCode"`
	Message             string  `json:"message"`
	TransactionID       string  `json:"transactionId"`
	RequestID           string  `json:"requestId"`
	OutletID            string  `json:"outletId"`
	Operation           string  `json:"operation"`
	TransactionStatus   string  `json:"txnStatus"`
	IPayID              string  `json:"ipayId"`
	OrderID             string  `json:"orderId"`
	Amount              float64 `json:"amount"`
	BankName            string  `json:"bankName"`
	AccountNumber       string  `json:"accountNumber"`
	TransactionMode     string  `json:"transactionMode"`
	BankAccountBalance  float64 `json:"bankAccountBalance"`
	IsOnusTransaction   bool    `json:"isOnusTxn"`
	ExternalReferenceID string  `json:"externalRef"`
	MiniStatement       []any   `json:"miniStatemenet"`
	CommisionAmount     float64 `json:"commisionAmount"`
	TDSAmount           float64 `json:"tdsAmount"`
	GSTAmount           float64 `json:"gstAmount"`
	NETAmount           float64 `json:"netAmount"`
	BalanceAfter        float64 `json:"balanceAfter"`
	LedgerID            float64 `json:"ledgerId"`
	SettlementStatus    string  `json:"settlementStatus"`
	SettlementMode      string  `json:"settlementMode"`
}

type AEPSCashDepositRequestModel struct {
	RequestID                string                 `json:"requestId"`
	BankIdentificationNumber string                 `json:"bankin"`
	Mobile                   string                 `json:"mobile"`
	Amount                   string                 `json:"amount"`
	Latitude                 string                 `json:"latitude"`
	Longitude                string                 `json:"longitude"`
	CustomerName             string                 `json:"customerName"`
	Aadhaar                  string                 `json:"aadhaar"`
	BiometricData            AEPSBiometricDataModel `json:"biometricData"`
}

type AEPSCashDepositResponseModel struct {
	Status              string `json:"status"`
	StatusCode          string `json:"statusCode"`
	Message             string `json:"message"`
	TransactionID       string `json:"transactionId"`
	RequestID           string `json:"requestId"`
	OutletID            string `json:"outletId"`
	Operation           string `json:"operation"`
	TransactionStatus   string `json:"txnStatus"`
	IPayID              string `json:"ipayId"`
	OrderID             string `json:"orderId"`
	Amount              string `json:"amount"`
	BankName            string `json:"bankName"`
	AccountNumber       string `json:"accountNumber"`
	TransactionMode     string `json:"transactionMode"`
	BankAccountBalance  string `json:"bankAccountBalance"`
	IsOnusTransaction   string `json:"isOnusTxn"`
	ExternalReferenceID string `json:"externalRef"`
	MiniStatement       []any  `json:"miniStatemenet"`
	CommisionAmount     string `json:"commisionAmount"`
	TDSAmount           string `json:"tdsAmount"`
	GSTAmount           string `json:"gstAmount"`
	NETAmount           string `json:"netAmount"`
	BalanceAfter        string `json:"balanceAfter"`
	LedgerID            string `json:"ledgerId"`
	SettlementStatus    string `json:"settlementStatus"`
	SettlementMode      string `json:"settlementMode"`
}

type AEPSBalanceEnquiryResponseModel struct {
	Status              string  `json:"status"`
	StatusCode          string  `json:"statusCode"`
	Message             string  `json:"message"`
	TransactionID       string  `json:"transactionId"`
	RequestID           string  `json:"requestId"`
	OutletID            string  `json:"outletId"`
	Operation           string  `json:"operation"`
	TransactionStatus   string  `json:"txnStatus"`
	IPayID              string  `json:"ipayId"`
	OrderID             string  `json:"orderId"`
	Amount              float64 `json:"amount"`
	BankName            string  `json:"bankName"`
	AccountNumber       string  `json:"accountNumber"`
	TransactionMode     string  `json:"transactionMode"`
	BankAccountBalance  float64 `json:"bankAccountBalance"`
	IsOnusTransaction   bool    `json:"isOnusTxn"`
	ExternalReferenceID string  `json:"externalRef"`
	MiniStatement       []any   `json:"miniStatemenet"`
	CommisionAmount     float64 `json:"commisionAmount"`
	TDSAmount           float64 `json:"tdsAmount"`
	GSTAmount           float64 `json:"gstAmount"`
	NETAmount           float64 `json:"netAmount"`
	BalanceAfter        float64 `json:"balanceAfter"`
	LedgerID            float64 `json:"ledgerId"`
	SettlementStatus    string  `json:"settlementStatus"`
	SettlementMode      string  `json:"settlementMode"`
}

type AEPSTransactionOTPRequestModel struct {
	ExternalReference        string `json:"externalRef"`
	OutletID                 string `json:"outletId"`
	BankIdentificationNumber string `json:"bankin"`
	Mobile                   string `json:"mobile"`
	Amount                   string `json:"amount"`
	Latitude                 string `json:"latitude"`
	Longitude                string `json:"longitude"`
	Aadhaar                  string `json:"aadhaar"`
}

type AEPSTransactionOTPResponseModel struct {
	Status              string `json:"status"`
	StatusCode          string `json:"statusCode"`
	ActionCode          string `json:"actionCode"`
	Message             string `json:"message"`
	IPayID              string `json:"ipayId"`
	IPayUUID            string `json:"ipayUuid"`
	ExternalReferenceID string `json:"externalRef"`
	OutletID            string `json:"outletId"`
	Data                struct {
		AdditionalProp1 map[any]any `json:"additionalProp1"`
		AdditionalProp2 map[any]any `json:"additionalProp2"`
		AdditionalProp3 map[any]any `json:"additionalProp3"`
	} `json:"data"`
}

type AEPSBankModel struct {
	BankID                   int    `json:"bankId"`
	BankName                 string `json:"name"`
	BankIdentificationNumber string `json:"iin"`
	AEPSEnabled              bool   `json:"aepsEnabled"`
	AadhaarPayEnabled        bool   `json:"aadhaarpayEnabled"`
	AEPSFailureRate          string `json:"aepsFailureRate"`
	AadhaarPayFailureRate    string `json:"aadhaarpayFailureRate"`
}

type AepsTransactionResponse struct {
	AepsTransactionID  int64           `json:"aeps_transaction_id" db:"aeps_transaction_id"`
	ReferenceID        string          `json:"reference_id" db:"reference_id"`
	TransactionID      string          `json:"transaction_id" db:"transaction_id"`
	OrderID            string          `json:"order_id" db:"order_id"`
	CustomerName       string          `json:"customer_name" db:"customer_name"`
	CustomerPhone      string          `json:"customer_phone" db:"customer_phone"`
	CustomerAadhaar    string          `json:"customer_aadhaar" db:"customer_aadhaar"`
	Amount             float64         `json:"amount" db:"amount"`
	MdCommission       float64         `json:"md_commission" db:"md_commision"`
	DisCommission      float64         `json:"dis_commission" db:"dis_commision"`
	RetailerCommission float64         `json:"retailer_commission" db:"retailer_commision"`
	TransactionStatus  string          `json:"transaction_status" db:"transaction_status"`
	CreatedAt          time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt          time.Time       `json:"updated_at" db:"updated_at"`
	RetailerName       string          `json:"retailer_name" db:"retailer_name"`
	BeforeBalance      sql.NullFloat64 `json:"before_balance" db:"before_balance"`
	AfterBalance       sql.NullFloat64 `json:"after_balance" db:"after_balance"`
	Reason             sql.NullString  `json:"reason" db:"transaction_reason"`
	Remarks            sql.NullString  `json:"remarks" db:"remarks"`
}

type AepsTdsDeductionResponse struct {
	TdsID                      int64     `json:"tds_id" db:"tds_id"`
	AepsTransactionID          int64     `json:"aeps_transaction_id" db:"aeps_transaction_id"`
	UserID                     string    `json:"user_id" db:"user_id"`
	UserType                   string    `json:"user_type" db:"user_type"`
	CreatedAt                  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt                  time.Time `json:"updated_at" db:"updated_at"`
	CustomerName               string    `json:"customer_name" db:"customer_name"`
	UserName                   string    `json:"user_name" db:"user_name"`
	TdsAmount                  float64   `json:"tds_amount" db:"tds_amount"`
	PanNumber                  string    `json:"pan_number"`
	RetailerCommision          float64   `json:"retailer_commision"`
	MasterDistributorCommision float64   `json:"master_distributor_commision"`
	DistributorCommision       float64   `json:"distributor_commision"`
}

type AEPSGetOTPRequestModel struct {
	ExternalReference        string `json:"externalRef"`
	OutletID                 string `json:"outletId"`
	BankIdentificationNumber string `json:"bankin"`
	Mobile                   string `json:"mobile"`
	Amount                   string `json:"amount"`
	Latitude                 string `json:"latitude"`
	Longitude                string `json:"longitude"`
	Aadhaar                  string `json:"aadhaar"`
	ReferenceKey             string `json:"referenceKey"`
}
