package dao

import (
	"sort"
	"purchase-tracker-service/domain"
)

type RewardTransactionByTimestamp []*domain.RewardTransaction

func (t RewardTransactionByTimestamp) Len() int {
	return len(t)
}

func (t RewardTransactionByTimestamp) Less(i int, j int) bool {
	return (*t[i].TransactionTimestamp).Before(*t[j].TransactionTimestamp)
}

func (t RewardTransactionByTimestamp) Swap(i int, j int) {
	t[i], t[j] = t[j], t[i]
}

type TransactionsDao interface {
	// Add a new Transaction to the system.
	AddTransaction(transaction *domain.RewardTransaction)
	// Return all Transactions sorted by the Transaction Timestamp
	GetTransactionLog() []*domain.RewardTransaction
}

type LocalTransactionsStore struct {
	cache []*domain.RewardTransaction
}

func NewLocalTransactionsStore(payers *PayerAccountsDao) *LocalTransactionsStore {
	return &LocalTransactionsStore{}
}

func (store *LocalTransactionsStore) AddTransaction(transaction *domain.RewardTransaction) {
	store.cache = append(store.cache, transaction)
}

func (store *LocalTransactionsStore) GetTransactionLog() []*domain.RewardTransaction {
	// this definitely is inefficient in an actual business to store a transaction log by a field
	// persistence systems like a relational DB can hash-sort items by the timestamp, but for
	// this exercise we need only achieve the requested functionality.
	sort.Sort(RewardTransactionByTimestamp(store.cache))
	return store.cache
}

