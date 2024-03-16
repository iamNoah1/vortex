package main

import (
	"net/http"
	"os"

	"github.com/rs/zerolog/log"

	"github.com/gorilla/mux"
	"github.com/iamNoah1/vortex/transaction"
)

func main() {
	transactionFilePath := os.Getenv("TRANSACTION_LOG_FILE")
	if transactionFilePath == "" {
		transactionFilePath = "./transactions.log"
	}

	transactions, err := transaction.NewFileTransactionLogger(transactionFilePath)
	if err != nil {
		log.Fatal().Msgf("failed to create event logger: %s", err)
	}
	err = transactions.ReplayEvents()
	if err != nil {
		log.Fatal().Msgf("failed to replay events: %s", err)
	}

	router := mux.NewRouter()
	router.HandleFunc("/v1/kv/{key}", func(w http.ResponseWriter, r *http.Request) {
		putKeyValuePairHandler(w, r, transactions)
	}).Methods("PUT")

	router.HandleFunc("/v1/kv/{key}", func(w http.ResponseWriter, r *http.Request) {
		deleteKeyValuePairHandler(w, r, transactions)
	}).Methods("DELETE")

	router.HandleFunc("/v1/kv/{key}", getKeyValuePairHandler).Methods("GET")

	router.HandleFunc("/v1/kv", getAllKeyValuePairsHandler).Methods("GET")

	log.Info().Msg("server started on port 8080")
	err = http.ListenAndServe(":8080", router)
	if err != nil {
		log.Fatal().Msgf("failed to start server: %s", err)
	}
}
