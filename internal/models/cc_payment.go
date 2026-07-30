package models

import "time"

type CreateCreditCardBeneficiaryRequestModel struct {
	RetailerID      string `json:"retailer_id"`
	RetailerName    string `json:"retailer_name"`
	BeneficiaryName string `json:"beneficiary_name"`
	AccountNumber   string `json:"account_number"`
	IFSCCode        string `json:"ifsc_code"`
	BankName        string `json:"bank_name"`
	PhoneNumber     string `json:"phone_number"`
	OperatorName    string `json:"operator_name"`
	OperatorCode    string `json:"operator_code"`
}

type UpdateCreditCardBeneficiaryRequestModel struct {
	BeneficiaryID   int64  `json:"beneficiary_id"`
	BeneficiaryName string `json:"beneficiary_name"`
	AccountNumber   string `json:"account_number"`
	IFSCCode        string `json:"ifsc_code"`
	BankName        string `json:"bank_name"`
	PhoneNumber     string `json:"phone_number"`
	OperatorName    string `json:"operator_name"`
	OperatorCode    string `json:"operator_code"`
}

type GetCreditCardBeneficiaryDetailsResponseModel struct {
	BeneficiaryID   int64     `json:"beneficiary_id"`
	RetailerID      string    `json:"retailer_id"`
	RetailerName    string    `json:"retailer_name"`
	BeneficiaryName string    `json:"beneficiary_name"`
	AccountNumber   string    `json:"account_number"`
	IFSCCode        string    `json:"ifsc_code"`
	BankName        string    `json:"bank_name"`
	PhoneNumber     string    `json:"phone_number"`
	OperatorName    string    `json:"operator_name"`
	OperatorCode    string    `json:"operator_code"`
	CreatedAT       time.Time `json:"created_at"`
	UpdatedAT       time.Time `json:"updated_at"`
}

type CreateCreditCardPaymentTransactionRequestModel struct {
	BeneDetails      GetCreditCardBeneficiaryDetailsResponseModel
	Amount           string `json:"amount"`
	PartnerRequestID string `json:"partner_request_id"`
}

type UpdateCreditCardPaymentTransactionRequestModel struct {
	TransactionID         int    `json:"transaction_id"`
	Status                int    `json:"status"`
	OrderID               string `json:"orderid"`
	OperatorTransactionID string `json:"optransid"`
}

type CardProviderDetailsModel struct {
	OperatorID           string `json:"operator_id"`
	OperatorName         string `json:"operator_name"`
	ServiceName          string `json:"service_name"`
	OperatorCategory     string `json:"operator_category"`
	OperatorCategoryName string `json:"operator_category_name"`
}
