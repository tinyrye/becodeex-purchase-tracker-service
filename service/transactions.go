package service

import (
	"fmt"
	"math"
	"purchase-tracker-service/dao"
	"purchase-tracker-service/domain"
)

type TransactionService interface {
	// When a Purchaser makes a new Purchase, this will accumulate Points under a Payer.
	ReceiveNewPurchase(transaction *domain.RewardTransaction) (*domain.RewardsAccumulateProgress, error)
	// Get the Current Points Balance/Progress for a single Payer.
	GetPointsProgressForPayer(payerId string) *domain.RewardsAccumulateProgress
	// Get the Current Points Balance/Progress for all known Payers.
	GetAllPointsProgressesForPayers() []*domain.RewardsAccumulateProgress
	// Spend Points using internal allocation logic gather values from Partners' balances.
	SpendPoints(numberOfPoints int) []*domain.RewardsAccumulateProgress
}

type LocalTransactionService struct {
	payerStore *dao.LocalPayerStore
	transactionsStore *dao.LocalTransactionsStore
	rewardsStore *dao.LocalRewardsStore
}

func NewLocalTransactionService() *LocalTransactionService {
	return &LocalTransactionService{
		&dao.LocalPayerStore{},
		&dao.LocalTransactionsStore{},
		&dao.LocalRewardsStore{},
	}
}

func (s *LocalTransactionService) ReceiveNewPurchase(transaction *domain.RewardTransaction) (*domain.RewardsAccumulateProgress, error) {
	var payer = s.payerStore.GetWithId(transaction.Payer)
	if payer == nil {
		return nil, PayerNotFoundError{transaction.Payer}
	}
	s.transactionsStore.AddTransaction(transaction)
	s.rewardsStore.AddTransaction(transaction)
	return s.getPointsProgressWithPayer(payer), nil
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
	var pointsByPayer map[string]int
	for _, payer := range allPayers {
		pointsByPayer[payer.Id] = s.getPointsForPayer(payer.Id)
	}
	return pointsByPayer
}

func (s *LocalTransactionService) SpendPoints(numberOfPoints int) []*domain.RewardsSpendAllocation {
	var txLog = s.transactionsStore.GetTransactionLog()
	var totalSpendBalance int = numberOfPoints
	var currentBalances = s.getAllPointsForPayers()
	var payerSpendAllocation map[string]int
	var payerSpendBalance map[string]int
	for payer, balance := range currentBalances {
		payerSpendBalance[payer] = balance
		payerSpendAllocation[payer] = 0
	}
	for _, tx := range txLog {
		currentPayerBalance := payerSpendBalance[tx.Payer]
		currentPayerAllocation := payerSpendAllocation[tx.Payer]
		if tx.Points > 0 {
			// we are taking from the Payer
			amountToTake := int(math.Min(float64(currentPayerBalance), float64(tx.Points)))
			payerSpendBalance[tx.Payer] = currentPayerBalance - amountToTake
			payerSpendAllocation[tx.Payer] = currentPayerAllocation + amountToTake
			totalSpendBalance = totalSpendBalance - amountToTake
		} else if tx.Points < 0 {
			// we are taking from the Payer
			amountToGive := int(math.Max(float64(currentPayerBalance), float64(tx.Points)))
			payerSpendBalance[tx.Payer] = currentPayerBalance + amountToGive
			payerSpendAllocation[tx.Payer] = currentPayerAllocation - amountToGive
			totalSpendBalance = totalSpendBalance + amountToGive
		}
		if totalSpendBalance == 0 {
			break
		}
	}
	return s.buildRewardAllocations(payerSpendAllocation)
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
