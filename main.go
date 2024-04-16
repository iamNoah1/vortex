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

func createTransactionLogger() (transaction.TransactionLogger, error) {
	transactionLoggerType := os.Getenv("TRANSACTION_LOGGER_TYPE")

	if transactionLoggerType == "" {
		transactionLoggerType = "file"
	}

	switch transactionLoggerType {
	case "file":
		transactionFilePath := os.Getenv("TRANSACTION_LOG_FILE")
		if transactionFilePath == "" {
			transactionFilePath = "./transactions.log"
		}
		transaction.NewFileTransactionLogger(transactionFilePath)
	case "psql":
		host := os.Getenv("DB_HOST")
		dbName := os.Getenv("DB_NAME")
		user := os.Getenv("DB_USER")
		password := os.Getenv("DB_PASSWORD")

		param := transaction.PostgresDBParams{
			DbName:   dbName,
			Host:     host,
			User:     user,
			Password: password,
		}

		//TODO make tablename configurable
		return transaction.NewPsqlTransactionLogger(param, "transactions")
	}
	//TODO
	return nil, nil
}
