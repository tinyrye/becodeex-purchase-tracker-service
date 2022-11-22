package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"github.com/gorilla/mux"
	"purchase-tracker-service/domain"
	"purchase-tracker-service/service"
)

type Application struct {
	transactionService *service.LocalTransactionService
	context context.Context
}

func main() {
	var flagSet = flag.NewFlagSet("http-server", flag.ExitOnError)
	var (
		serverHttpAddress = flagSet.String("http-address", ":8999", "The host:port address to bind to a server socket and listen for requests.")
	)
	var transactionService = service.NewLocalTransactionService()
	var application = Application{
		transactionService,
		context.Background(),
	}
	var httpRouter = mux.NewRouter()
	transactionService.AddPayer("DANNON", "Dannon")
	transactionService.AddPayer("UNILEVER", "Unilever")
	transactionService.AddPayer("MILLER COORS", "Miller Coors")
	httpRouter.Handle("/payers/balances", application.HandleGetAllPayersBalances()).Methods("GET")
	httpRouter.Handle("/payers/${payerId}/balances", application.HandleGetPayerBalances()).Methods("GET")
	httpRouter.Handle("/purchases", application.HandleAddPurchaseTransaction()).Methods("POST")
	httpRouter.Handle("/rewards/spend", application.HandleNewPointsSpendTransaction()).Methods("POST")
	http.Handle("/", httpRouter)
	log.Printf("Listening with HTTP server on %s", *serverHttpAddress)
	log.Fatal(http.ListenAndServe(*serverHttpAddress, nil))
}

func (a *Application) HandleGetAllPayersBalances() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var result = a.GetAllPayersBalances()
		WriteServiceResponse(w, result, nil)
	})
}

func (a *Application) HandleGetPayerBalances() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var payerId = mux.Vars(r)["payerId"]
		var result, serviceError = a.GetPayerBalance(payerId)
		WriteServiceResponse(w, result, serviceError)
	})
}

func (a *Application) HandleAddPurchaseTransaction() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		transaction, requestDecodeErr := decodePurchaseTransactionRequest(a.context, r)
		if requestDecodeErr != nil {
			WriteDecodeErrorResponse(w, requestDecodeErr)
		} else {
			result, serviceError := a.AddPurchaseTransaction(transaction)
			WriteServiceResponse(w, result, serviceError)
		}
	})
}

func (a *Application) HandleNewPointsSpendTransaction() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		transaction, requestDecodeErr := decodePointsSpendTransactionRequest(a.context, r)
		if requestDecodeErr != nil {
			WriteDecodeErrorResponse(w, requestDecodeErr)
		} else {
			result := a.SpendPoints(transaction)
			WriteServiceResponse(w, result, nil)
		}
	})
}

func (a *Application) GetAllPayersBalances() []*domain.RewardsAccumulateProgress {
	return a.transactionService.GetAllPointsProgressesForPayers()
}

func (a *Application) GetPayerBalance(payerId string) (*domain.RewardsAccumulateProgress, error) {
	return a.transactionService.GetPointsProgressForPayer(payerId)
}

func (a *Application) AddPurchaseTransaction(transaction *domain.RewardTransaction) (*domain.RewardsAccumulateProgress, error) {
	return a.transactionService.ReceiveNewPurchase(transaction)
}

func (a *Application) SpendPoints(transaction *domain.PointsSpendTransaction) []*domain.RewardsAccumulateProgress {
	a.transactionService.SpendPoints(transaction.Points)
	return a.GetAllPayersBalances()
}

func decodePurchaseTransactionRequest(_ context.Context, r *http.Request) (*domain.RewardTransaction, error) {
	var request domain.RewardTransaction
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return &request, nil
}

func decodePointsSpendTransactionRequest(_ context.Context, r *http.Request) (*domain.PointsSpendTransaction, error) {
	var request domain.PointsSpendTransaction
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return &request, nil
}

func WriteDecodeErrorResponse(w http.ResponseWriter, requestDecodeErr error) {
	w.WriteHeader(422)
	json.NewEncoder(w).Encode(map[string]string {
		"status": "UNPROCESSIBLE ENITTY",
		"message": fmt.Sprintf("Request body is missing or invalid: %s", requestDecodeErr),
	})
}

func WriteServiceResponse(w http.ResponseWriter, result interface{}, error error) {
	if error != nil {
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(map[string]string {
			"status": "Internal Failure",
			"message": error.Error(),
		})
	} else {
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(result)
	}
}

func WriteNotMappedResponse(w http.ResponseWriter) {
	w.WriteHeader(404)
	json.NewEncoder(w).Encode(map[string]string {
		"status": "NOT FOUND",
		"message": "Endpoint is not mapped",
	})
}

func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}
