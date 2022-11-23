package main

import (
	"testing"
	"time"
	"purchase-tracker-service/dao"
	"purchase-tracker-service/domain"
)

func TestListAllAccounts(t *testing.T) {
	var testedObject = dao.NewLocalPayerStore()
	var actualListing = testedObject.ListAllAccounts()
	expectNoAccounts(t, actualListing)

	var firstAccount = &domain.PayerAccount{
		"account-1",
		"Brand 1",
		time.Now(),
	}
	var secondAccount = &domain.PayerAccount{
		"account-2",
		"Brand 2",
		time.Now(),
	}
	testedObject.AddAccount(firstAccount)
	testedObject.AddAccount(secondAccount)
	actualListing = testedObject.ListAllAccounts()
	expectSpecificAccounts(t, actualListing, 2)
	expectAccountById(t, actualListing, firstAccount.Id)
	expectAccountById(t, actualListing, secondAccount.Id)
}

func TestGetWithId(t *testing.T) {
	var testedObject = dao.NewLocalPayerStore()
	var testAccount = &domain.PayerAccount{
		"account-1",
		"Foobar",
		time.Now(),
	}

	var actualAccount = testedObject.GetWithId(testAccount.Id)
	if actualAccount != nil {
		t.Fatalf("Did not expect Account to be found by id %s", testAccount.Id)
	}
	testedObject.AddAccount(testAccount)
	actualAccount = testedObject.GetWithId(testAccount.Id)
	if actualAccount == nil {
		t.Fatalf("Expected Account to be found by id %s", testAccount.Id)
	}
	if actualAccount.Id != testAccount.Id {
		t.Fatalf("Expected Account to have the id %s but was %s", testAccount.Id, actualAccount.Id)
	}
	if actualAccount.Name != testAccount.Name {
		t.Fatalf("Expected Account to have the name %s but was %s", testAccount.Name, actualAccount.Name)
	}
	if actualAccount.CreationTimestamp != testAccount.CreationTimestamp {
		t.Fatalf("Expected Account to have the creationTimestamp %s but was %s", testAccount.CreationTimestamp, actualAccount.CreationTimestamp)
	}
}

func TestGetWithName(t *testing.T) {
	var testedObject = dao.NewLocalPayerStore()
	var testAccount = &domain.PayerAccount{
		"account-1",
		"Foobar",
		time.Now(),
	}

	var actualAccount = testedObject.GetWithName(testAccount.Name)
	if actualAccount != nil {
		t.Fatalf("Did not expect Account to be found by id %s", testAccount.Id)
	}
	testedObject.AddAccount(testAccount)
	actualAccount = testedObject.GetWithName(testAccount.Name)
	if actualAccount == nil {
		t.Fatalf("Expected Account to be found by id %s", testAccount.Id)
	}
	if actualAccount.Id != testAccount.Id {
		t.Fatalf("Expected Account to have the id %s but was %s", testAccount.Id, actualAccount.Id)
	}
	if actualAccount.Name != testAccount.Name {
		t.Fatalf("Expected Account to have the name %s but was %s", testAccount.Name, actualAccount.Name)
	}
	if actualAccount.CreationTimestamp != testAccount.CreationTimestamp {
		t.Fatalf("Expected Account to have the creationTimestamp %s but was %s", testAccount.CreationTimestamp, actualAccount.CreationTimestamp)
	}
}

func expectNoAccounts(t *testing.T, actualListing []*domain.PayerAccount) {
	if len(actualListing) > 0 {
		t.Fatalf("Expected an empty list of accounts but the list has %d elements", len(actualListing))
	}
}

func expectSpecificAccounts(t *testing.T, actualListing []*domain.PayerAccount, expectedNumber int) {
	if len(actualListing) != expectedNumber {
		t.Fatalf("Expected an list of %d accounts but the list has %d elements", expectedNumber, len(actualListing))
	}
}

func expectAccountById(t *testing.T, actualListing []*domain.PayerAccount, expectedId string) *domain.PayerAccount {
	for _, actualAccount := range actualListing {
		if actualAccount.Id == expectedId {
			return actualAccount
		}
	}
	t.Fatalf("Expected to find Account with the %s but none was found.", expectedId)
	return nil
}
