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
