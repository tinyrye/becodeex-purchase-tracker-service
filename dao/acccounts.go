package dao

import (
	"fmt"
	"log"
	"regexp"
	"strings"
	"purchase-tracker-service/domain"
)

const (
	whitespaceTokenizerPattern = "\\s"
)

type PayerAccountsDao interface {
	AddAccount(*domain.PayerAccount)
	ListAllAccounts() []*domain.PayerAccount
	GetWithId(id string) *domain.PayerAccount
	GetWithName(name string) *domain.PayerAccount
	SearchWithNameQuery(query string) []*domain.PayerAccount
}

type LocalPayerStore struct {
	cacheById map[string]*domain.PayerAccount
	cacheByName map[string]*domain.PayerAccount
	cacheByTokens map[string][]*domain.PayerAccount
}

func NewLocalPayerStore() *LocalPayerStore {
	return &LocalPayerStore{
		make(map[string]*domain.PayerAccount),
		make(map[string]*domain.PayerAccount),
		make(map[string][]*domain.PayerAccount),
	}
}

type AccountExistsError struct {
	PayerId string
}

func (e AccountExistsError) Error() string {
	return fmt.Sprintf("Payer Account already registered: %s", e.PayerId)
}

func (s *LocalPayerStore) AddAccount(payer *domain.PayerAccount) error {
	s.cacheById[payer.Id] = payer
	s.cacheByName[payer.Name] = payer
	nameTokens := tokenizeSearchableTerm(payer.Name)
	for _, nameToken := range nameTokens {
		cacheOfToken, tokenEntryExists := s.cacheByTokens[nameToken]
		if !tokenEntryExists {
			cacheOfToken = make([]*domain.PayerAccount, 0)
			s.cacheByTokens[nameToken] = cacheOfToken
		}
		cacheOfToken = append(cacheOfToken, payer)
	}
	log.Printf("Added account %s with tokens %s", payer.Id, nameTokens)
	return nil
}

func (s *LocalPayerStore) ListAllAccounts() []*domain.PayerAccount {
	var allPayers []*domain.PayerAccount
	for _, payer := range s.cacheById {
		allPayers = append(allPayers, payer)
	}
	return allPayers
}

func (s *LocalPayerStore) GetWithId(id string) *domain.PayerAccount {
	return s.cacheById[id]
}

func (s *LocalPayerStore) GetWithName(name string) *domain.PayerAccount {
	return s.cacheByName[name]
}

func (s *LocalPayerStore) SearchWithNameQuery(query string) []*domain.PayerAccount {
	var queryTokens = tokenizeSearchableTerm(query)
	for _, queryToken := range queryTokens {
		cacheItems, cacheHit := s.cacheByTokens[queryToken]
		if cacheHit {
			return cacheItems
		}
	}
	return make([]*domain.PayerAccount, 0)
}

// derive a set of searchable tokens found within a name.
func tokenizeSearchableTerm(name string) []string {
	var searchableTokens []string = make([]string, 0)
	var whitespaceTokenizer, whitespaceTokenizerErr = regexp.Compile("\\s")
	if whitespaceTokenizerErr != nil {
		log.Fatalf("REGEX error %s", whitespaceTokenizerErr)
	}
	var tokensAroundWhitespace = whitespaceTokenizer.FindAll([]byte(name), -1)
	for _, tokenAroundWhitespace := range tokensAroundWhitespace {
		searchableTokens = append(searchableTokens, strings.ToLower(string(tokenAroundWhitespace)))
	}
	for lengthSubSet := 1; lengthSubSet < len(name) - 1; lengthSubSet++ {
		namePrefix := name[0:lengthSubSet]
		searchableTokens = append(searchableTokens, strings.ToLower(namePrefix))
	}
	return searchableTokens
}
