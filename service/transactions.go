package service

import (
	"fmt"
	"log"
	"math"
	"time"
	"purchase-tracker-service/dao"
	"purchase-tracker-service/domain"
)

type TransactionService interface {
	AddPayer(id string, name string) error
	// Get the Current Points Balance/Progress for all known Payers.
	GetAllPointsProgressesForPayers() []*domain.RewardsAccumulateProgress
	// Get the Current Points Balance/Progress for a single Payer.
	GetPointsProgressForPayer(payerId string) (*domain.RewardsAccumulateProgress, error)
	// When a Purchaser makes a new Purchase, this will accumulate Points under a Payer.
	ReceiveNewPurchase(transaction *domain.RewardTransaction) (*domain.RewardsAccumulateProgress, error)
	// Spend Points using internal allocation logic gather values from Partners' balances.
	SpendPoints(numberOfPoints int) []*domain.RewardsSpendAllocation
}

type LocalTransactionService struct {
	payerStore *dao.LocalPayerStore
	transactionsStore *dao.LocalTransactionsStore
	rewardsStore *dao.LocalRewardsStore
}

func NewLocalTransactionService() *LocalTransactionService {
	return &LocalTransactionService{
		dao.NewLocalPayerStore(),
		dao.NewLocalTransactionsStore(),
		dao.NewLocalRewardsStore(),
	}
}

func (s *LocalTransactionService) AddPayer(id string, name string) error {
	return s.payerStore.AddAccount(&domain.PayerAccount{
		id,
		name,
		time.Now(),
	})
}

func (s *LocalTransactionService) GetPointsProgressForPayer(payerId string) (*domain.RewardsAccumulateProgress, error) {
	var payer = s.payerStore.GetWithId(payerId)
	if payer == nil {
		return nil, PayerNotFoundError{payerId}
	}
	return s.getPointsProgressWithPayer(payer), nil
}

func (s *LocalTransactionService) getPointsProgressWithPayer(payer *domain.PayerAccount) *domain.RewardsAccumulateProgress {
	return &domain.RewardsAccumulateProgress{payer, s.getPointsForPayer(payer.Id)}
}

func (s *LocalTransactionService) getPointsForPayer(payerId string) int {
	var pointsForPayer = s.rewardsStore.GetPointsForPayer(payerId)
	if pointsForPayer == nil {
		return 0
	} else {
		return *pointsForPayer
	}
}

func (s *LocalTransactionService) GetAllPointsProgressesForPayers() []*domain.RewardsAccumulateProgress {
	var allPayers = s.payerStore.ListAllAccounts()
	var allProgresses []*domain.RewardsAccumulateProgress
	for _, payer := range allPayers {
		allProgresses = append(allProgresses, s.getPointsProgressWithPayer(payer))
	}
	return allProgresses
}

func (s *LocalTransactionService) getAllPointsForPayers() map[string]int {
	var allPayers = s.payerStore.ListAllAccounts()
	var pointsByPayer map[string]int = make(map[string]int)
	for _, payer := range allPayers {
		pointsByPayer[payer.Id] = s.getPointsForPayer(payer.Id)
	}
	return pointsByPayer
}

func (s *LocalTransactionService) ReceiveNewPurchase(transaction *domain.RewardTransaction) (*domain.RewardsAccumulateProgress, error) {
	var payer = s.payerStore.GetWithId(transaction.Payer)
	if payer == nil {
		return nil, PayerNotFoundError{transaction.Payer}
	}
	transaction.TransactionTimestamp = time.Now()
	log.Printf("Adding Transaction %s", transaction)
	s.addTransaction(transaction)
	return s.getPointsProgressWithPayer(payer), nil
}

func (s *LocalTransactionService) creditPayer(payerId string, pointsToCredit int) {
	log.Printf("Payer %s being credited %d", payerId, pointsToCredit)
	s.addTransaction(&domain.RewardTransaction{
		payerId,
		-pointsToCredit,
		time.Now(),
	})
}

func (s *LocalTransactionService) addTransaction(transaction *domain.RewardTransaction) {
	s.transactionsStore.AddTransaction(transaction)
	s.rewardsStore.AddTransaction(transaction)
}

func (s *LocalTransactionService) SpendPoints(numberOfPoints int) []*domain.RewardsSpendAllocation {
	var txLog = s.transactionsStore.GetTransactionLog()
	var totalSpendBalance int = numberOfPoints
	var currentBalances = s.getAllPointsForPayers()
	var payerSpendAllocation map[string]int = make(map[string]int)
	var payerSpendBalance map[string]int = make(map[string]int)
	for payer, balance := range currentBalances {
		payerSpendBalance[payer] = balance
		payerSpendAllocation[payer] = 0
	}
	for _, tx := range txLog {
		currentPayerBalance := payerSpendBalance[tx.Payer]
		currentPayerAllocation := payerSpendAllocation[tx.Payer]
		if tx.Points != 0 {
			// we are applying points to a rewards to removing points from the reward
			// depending on whether the transaction is a negative or positive.
			amountToApply := int(math.Min(math.Min(float64(currentPayerBalance), float64(tx.Points)), float64(totalSpendBalance)))
			payerSpendBalance[tx.Payer] = currentPayerBalance - amountToApply
			payerSpendAllocation[tx.Payer] = currentPayerAllocation + amountToApply
			totalSpendBalance = totalSpendBalance - amountToApply
			log.Printf("Apply %d points from Payer %s and resulting in a points allocation balance of %d", amountToApply, tx.Payer, totalSpendBalance)
		}
		if totalSpendBalance == 0 {
			break
		}
	}
	// Now, we have to credit these payer accounts the amount of Points being spent here
	s.creditPayerAccountsViaAllocation(payerSpendAllocation)
	return s.buildRewardAllocations(payerSpendAllocation)
}

func (s *LocalTransactionService) creditPayerAccountsViaAllocation(spendAllocationByPayerId map[string]int) {
	for payerId, pointsSpent := range spendAllocationByPayerId {
		s.creditPayer(payerId, pointsSpent)
	}
}

func (s *LocalTransactionService) buildRewardAllocations(spendAllocationByPayerId map[string]int) []*domain.RewardsSpendAllocation {
	var payerAllocations []*domain.RewardsSpendAllocation
	for payerId, pointsSpent := range spendAllocationByPayerId {
		payerAllocations = append(payerAllocations, &domain.RewardsSpendAllocation{
			s.payerStore.GetWithId(payerId),
			-pointsSpent,
		})
	}
	return payerAllocations
}

type PayerNotFoundError struct {
	PayerId string
}

func (e PayerNotFoundError) Error() string {
	return fmt.Sprintf("Payer was not found in the system '%s'", e.PayerId)
}
