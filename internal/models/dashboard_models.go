package models

type ServiceStat struct {
	TotalAmount       float64 `json:"total_amount"`
	TotalTransactions int64   `json:"total_transactions"`
}

type RetailerDashboard struct {
	Payout          ServiceStat `json:"payout"`
	MobileRecharge  ServiceStat `json:"mobile_recharge"`
	DTHRecharge     ServiceStat `json:"dth_recharge"`
	ElectricityBill ServiceStat `json:"electricity_bill"`
	FundRequests    ServiceStat `json:"fund_requests"`
	WalletCredited  float64     `json:"wallet_credited"`
	WalletDebited   float64     `json:"wallet_debited"`
}
