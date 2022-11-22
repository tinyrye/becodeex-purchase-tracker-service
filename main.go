package main

import (
	"context"
	"encoding/json"
	"flag"
	"log"
	"net/http"

	"purchase-tracker-service/dao"
	"purchase-tracker-service/domain"
	"purchase-tracker-service/service"
	httptransport "github.com/go-kit/kit/transport/http"
)

func NewLocalTransactionService(payerStore *dao.LocalTransactionsStore, transactionsStore *dao.LocalTransactionsStore) *service.LocalTransactionService {
	return &service.LocalTransactionService{
		
	}
}

func main() {
	var flagSet = flag.NewFlagSet("http-server", flag.ExitOnError)
	var (
		serverHttpAddress = flagSet.String("http-address", ":8999", "The host:port address to bind to a server socket and listen for requests.")
	)
	var purchaseHandler = httptransport.NewServer(AddPurchaseTransaction, decodeSpendRequest, encodeResponse)
	http.Handle("/purchases", purchaseHandler)
	log.Fatal(http.ListenAndServe(*serverHttpAddress, nil))
}

func AddPurchaseTransaction(requestContext context.Context, request interface{}) (interface{}, error) {
	log.Printf("Hello world")
	transactionService := requestContext.Value("transactionService").(service.LocalTransactionService)
	transaction := request.(domain.RewardTransaction)
	return transactionService.ReceiveNewPurchase(&transaction)
}

func DecodePurchaseTransactionRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request domain.RewardTransaction
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func decodeSpendRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request domain.RewardTransaction
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}
