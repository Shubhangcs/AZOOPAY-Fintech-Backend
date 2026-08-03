package models

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
