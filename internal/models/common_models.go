package models

type DropdownItem struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type APIResponseModel struct {
	Error                 int    `json:"error"`
	Message               string `json:"msg"`
	Status                int    `json:"status"`
	OrderID               string `json:"orderid"`
	OperatorTransactionID string `json:"optransid"`
	PartnerRequestID      string `json:"partnerreqid"`
}

type PayntricAPIResponseModel struct {
	Status     string     `json:"status"`
	StatusCode string     `json:"statuscode"`
	Message    string     `json:"message"`
	Data       PayoutData `json:"data"`
}

type PayoutData struct {
	UTR             string  `json:"utr"`
	OpRefID         string  `json:"opRefId"`
	APITxnID        string  `json:"apiTxnId"`
	RequestID       string  `json:"requestId"`
	Amount          float64 `json:"amount"`
	BeneficiaryName string  `json:"beneficiaryName"`
	IFSCCode        string  `json:"ifscCode"`
	AccountNumber   string  `json:"accountNumber"`
	BankName        string  `json:"bankName"`
	ChargeAmount    float64 `json:"chargeAmount"`
	GSTAmount       float64 `json:"gstAmount"`
	TotalDeduction  float64 `json:"totalDeduction"`
	ClosingBalance  float64 `json:"closingBalance"`
}

type RechargeKitWalletBalanceResponseModel struct {
	Error    int     `json:"error"`
	Message  string  `json:"msg"`
	Balance  float64 `json:"wallet_amount"`
	Balance2 float32 `json:"payout_wallet_amount"`
}

type PayntricWalletBalanceResponseModel struct {
	WalletBalance float64 `json:"wallet_balance"`
	MerchantId    string  `json:"merchant_id"`
	Status        string  `json:"status"`
}
