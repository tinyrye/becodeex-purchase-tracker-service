package domain

import (
	"time"
)

type RewardTransaction struct {
	// this field could be the Id of the Payer.
	//  - eg. "DANNON" will match a PayerAccount with an Id of 'DANNON".
	Payer string `json:"payer"`
	Points int `json:"points"`
	TransactionTimestamp time.Time
}

type PointsSpendTransaction struct {
	Points int `json:"points"`
}
