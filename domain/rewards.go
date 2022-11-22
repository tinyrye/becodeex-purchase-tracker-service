package domain

type RewardsSpendAllocation struct {
	Payer *PayerAccount `json:"payer"`
	Points int `json:"points"`
}

type RewardsAccumulateProgress struct {
	// this field could be the Id or Name of the Payer since
	Payer *PayerAccount `json:"payer"`
	Points int
}
