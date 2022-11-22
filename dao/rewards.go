package dao

import (
	"log"
	"purchase-tracker-service/domain"
)

type RewardsDao interface {
	AddTransaction(transaction *domain.RewardTransaction)
	GetPointsForPayer(payerId string) *int
}

type LocalRewardsStore struct {
	cache map[string]int
}

func NewLocalRewardsStore() *LocalRewardsStore {
	return &LocalRewardsStore{make(map[string]int)}
}

func (store *LocalRewardsStore) AddTransaction(transaction *domain.RewardTransaction) {
	log.Printf("Payer transaction %s : %d", transaction.Payer, transaction.Points)
	var currentProgress, exists = store.cache[transaction.Payer]
	if !exists {
		currentProgress = 0
	}
	store.cache[transaction.Payer] = currentProgress + transaction.Points
}

func (store *LocalRewardsStore) GetPointsForPayer(payerId string) *int {
	if currentProgress, tracked := store.cache[payerId]; tracked {
		return &currentProgress
	} else {
		return nil
	}
}
