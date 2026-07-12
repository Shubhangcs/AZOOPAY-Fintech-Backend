package models

type StatsPeriod struct {
	TotalAmount       float64 `json:"total_amount"`
	TotalTransactions int64   `json:"total_transactions"`
}

type StatsResponse struct {
	Today   StatsPeriod `json:"today"`
	Weekly  StatsPeriod `json:"weekly"`
	Monthly StatsPeriod `json:"monthly"`
	Yearly  StatsPeriod `json:"yearly"`
}

type StatsModel struct {
	TotalRetailerWalletBalance          float64 `json:"total_retailer_wallet_balance"`
	TotalDistributorWalletBalance       float64 `json:"total_distributor_wallet_balance"`
	TotalMasterDistributorWalletBalance float64 `json:"total_master_distributor_wallet_balance"`
	TotalRetailerAdvanceCredit          float64 `json:"total_retailer_advance_credit"`
	TotalDistributorAdvanceCredit       float64 `json:"total_distributor_advance_credit"`
	TotalMasterDistributorAdvanceCredit float64 `json:"total_master_distributor_advance_credit"`
}
