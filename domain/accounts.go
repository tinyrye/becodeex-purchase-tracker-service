package domain

import (
	"time"
)

type PayerAccount struct {
	Id string
	Name string
	Balance int
	CreationTimestamp *time.Time
}
