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

type AEPSDailyLoginRequestModel struct {
	Type              string                 `json:"type"`
	ExternalReference string                 `json:"externalRef"`
	OutletID          string                 `json:"outletId"`
	Latitude          string                 `json:"latitude"`
	Longitude         string                 `json:"longitude"`
	CaptureType       string                 `json:"captureType"`
	BiometricData     AEPSBiometricDataModel `json:"biometricData"`
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
	Data                struct {
		AdditionalProp1 map[any]any `json:"additionalProp1"`
		AdditionalProp2 map[any]any `json:"additionalProp2"`
		AdditionalProp3 map[any]any `json:"additionalProp3"`
	} `json:"data"`
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
	Data                struct {
		AdditionalProp1 map[any]any `json:"additionalProp1"`
		AdditionalProp2 map[any]any `json:"additionalProp2"`
		AdditionalProp3 map[any]any `json:"additionalProp3"`
	} `json:"data"`
}

type AEPSMiniStatementRequestModel struct {
	RequestID                string                 `json:"requestId"`
	OutletID                 string                 `json:"outletId"`
	BankIdentificationNumber string                 `json:"bankin"`
	Mobile                   string                 `json:"mobile"`
	Amount                   string                 `json:"amount"`
	Latitude                 string                 `json:"latitude"`
	Longitude                string                 `json:"longitude"`
	CustomerName             string                 `json:"customerName"`
	CaptureType              string                 `json:"captureType"`
	BiometricData            AEPSBiometricDataModel `json:"biometricData"`
}

type AEPSMiniStatementResponseModel struct {
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

type AEPSCashWithdrawalRequestModel struct {
	RequestID                string                 `json:"requestId"`
	OutletID                 string                 `json:"outletId"`
	BankIdentificationNumber string                 `json:"bankin"`
	Mobile                   string                 `json:"mobile"`
	Amount                   string                 `json:"amount"`
	Latitude                 string                 `json:"latitude"`
	Longitude                string                 `json:"longitude"`
	CustomerName             string                 `json:"customerName"`
	CaptureType              string                 `json:"captureType"`
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
	OutletID                 string                 `json:"outletId"`
	BankIdentificationNumber string                 `json:"bankin"`
	Mobile                   string                 `json:"mobile"`
	Amount                   string                 `json:"amount"`
	Latitude                 string                 `json:"latitude"`
	Longitude                string                 `json:"longitude"`
	CustomerName             string                 `json:"customerName"`
	CaptureType              string                 `json:"captureType"`
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

type AEPSBalanceEnquiryRequestModel struct {
	RequestID                string                 `json:"requestId"`
	OutletID                 string                 `json:"outletId"`
	BankIdentificationNumber string                 `json:"bankin"`
	Mobile                   string                 `json:"mobile"`
	Amount                   string                 `json:"amount"`
	Latitude                 string                 `json:"latitude"`
	Longitude                string                 `json:"longitude"`
	CustomerName             string                 `json:"customerName"`
	CaptureType              string                 `json:"captureType"`
	Aadhaar                  string                 `json:"aadhaar"`
	BiometricData            AEPSBiometricDataModel `json:"biometricData"`
}

type AEPSBalanceEnquiryResponseModel struct {
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
	BankID                   string `json:"bankId"`
	BankName                 string `json:"name"`
	BankIdentificationNumber string `json:"iin"`
	AEPSEnabled              bool   `json:"aepsEnabled"`
	AadhaarPayEnabled        bool   `json:"aadhaarpayEnabled"`
	AEPSFailureRate          string `json:"aepsFailureRate"`
	AadhaarPayFailureRate    string `json:"aadhaarpayFailureRate"`
}
