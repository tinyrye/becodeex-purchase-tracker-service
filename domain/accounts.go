package domain

import (
	"time"
)

type PayerAccount struct {
	Id string `json:"id"`
	Name string `json:"name"`
	CreationTimestamp time.Time `json:"creationTimestamp"`
}
