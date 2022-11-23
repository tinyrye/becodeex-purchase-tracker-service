package main

import (
	"log"
	"strings"
	"testing"
	"time"
	"purchase-tracker-service/domain"
	"purchase-tracker-service/service"
)

func TestGetPointsProgressForPayer(t *testing.T) {
	var testedObject = service.NewLocalTransactionService()
	var testAccount = &domain.PayerAccount{
		"account-1",
		"Foobar",
		time.Now(),
	}

	var actualProgress, getError = testedObject.GetPointsProgressForPayer(testAccount.Id)
	if actualProgress != nil {
		t.Fatalf("Did not expect Account to be found by id %s", testAccount.Id)
	} else if getError == nil {
		t.Fatalf("Expected there to be an error on get")
	} else if !strings.HasPrefix(getError.Error(), "Payer was not found") {
		t.Fatalf("Expected the error to be a not found error")
	}
	testedObject.AddPayer(testAccount.Id, testAccount.Name)
	actualProgress, getError = testedObject.GetPointsProgressForPayer(testAccount.Id)
	if actualProgress == nil {
		t.Fatalf("Expected Account to be found by id %s", testAccount.Id)
	} else if getError != nil {
		t.Fatalf("Expected there to not be an error on get")
	}
	if actualProgress.Payer.Id != testAccount.Id {
		t.Fatalf("Expected Account to have the id %s but was %s", testAccount.Id, actualProgress.Payer.Id)
	}
	if actualProgress.Payer.Name != testAccount.Name {
		t.Fatalf("Expected Account to have the name %s but was %s", testAccount.Name, actualProgress.Payer.Name)
	}
	if actualProgress.Points != 0 {
		t.Fatalf("Expected Payer Points to be zero but was %d", actualProgress.Points)
	}
}

func TestReceiveNewPurchase(t *testing.T) {
	var testedObject = service.NewLocalTransactionService()
	var testAccount = &domain.PayerAccount{
		"account-1",
		"Foobar",
		time.Now(),
	}

	testedObject.AddPayer(testAccount.Id, testAccount.Name)

	var testTransaction = &domain.RewardTransaction {
		testAccount.Id,
		555,
		time.Time{},
	}

	testedObject.ReceiveNewPurchase(testTransaction)

	var actualProgress, getError = testedObject.GetPointsProgressForPayer(testAccount.Id)
	if actualProgress == nil {
		t.Fatalf("Expected Account to be found by id %s", testAccount.Id)
	} else if getError != nil {
		t.Fatalf("Expected there to not be an error on get")
	}
	if actualProgress.Points != testTransaction.Points {
		t.Fatalf("Expected Payer Points to be %d but was %d", testTransaction.Points, actualProgress.Points)
	}
}

func TestReceiveNewPurchase_PayerDoesNotExist(t *testing.T) {
	var testedObject = service.NewLocalTransactionService()
	var testAccount = &domain.PayerAccount{
		"account-1",
		"Foobar",
		time.Now(),
	}

	testedObject.AddPayer(testAccount.Id, testAccount.Name)

	var testTransaction = &domain.RewardTransaction {
		"account-2",
		555,
		time.Time{},
	}

	var actualProgress, txRecvError = testedObject.ReceiveNewPurchase(testTransaction)
	if actualProgress != nil {
		t.Fatal("Expected Progress to be nil")
	} else if txRecvError == nil {
		t.Fatalf("Expected there to be an error on get")
	} else if !strings.HasPrefix(txRecvError.Error(), "Payer was not found") {
		t.Fatalf("Expected the error to be a not found error")
	}
}

func TestSpendPoints(t *testing.T) {
	log.Printf("- - - - - - - - - - - - - TestSpendPoints")
	var testedObject = service.NewLocalTransactionService()
	var firstAccount = &domain.PayerAccount{
		"account-1",
		"Foobar",
		time.Time{},
	}
	var secondAccount = &domain.PayerAccount{
		"account-2",
		"Foobar 2",
		time.Time{},
	}

	testedObject.AddPayer(firstAccount.Id, firstAccount.Name)
	testedObject.AddPayer(secondAccount.Id, secondAccount.Name)

	// Transactions
	testedObject.ReceiveNewPurchase(&domain.RewardTransaction {
		firstAccount.Id,
		500,
		time.Time{},
	})
	testedObject.ReceiveNewPurchase(&domain.RewardTransaction {
		secondAccount.Id,
		349,
		time.Time{},
	})
	testedObject.ReceiveNewPurchase(&domain.RewardTransaction {
		secondAccount.Id,
		251,
		time.Time{},
	})
	// the goal is that when we want to spend 850 points, after the second TX,
	// the allocation is forced to go through all 4 transactions to spend later points.
	var actualAllocations = testedObject.SpendPoints(850)
	if len(actualAllocations) != 2 {
		t.Fatalf("Expected 2 alloations since there are 2 Payers set up in the system but %d were returned.", len(actualAllocations))
	}
	expectAllocationForPayer(t, actualAllocations, "account-1", -500)
	expectAllocationForPayer(t, actualAllocations, "account-2", -350)
}

func expectAllocationForPayer(t *testing.T, actualAllocations []*domain.RewardsSpendAllocation, expectedPayerId string, expectedPointsAllocated int) {
	for _, actualAllocation := range actualAllocations {
		if actualAllocation.Payer.Id == expectedPayerId {
			if actualAllocation.Points != expectedPointsAllocated {
				t.Fatalf("Expected %d points to be allocated from Payer account %s but %d was allocated instead", expectedPointsAllocated, expectedPayerId, actualAllocation.Points)
			}
			return
		}
	}
	t.Fatalf("Expected Payer %s to be among those listed with points allocations but it was not", expectedPayerId)
}
