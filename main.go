package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

type KeyValuePair struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func main() {
	transactionFilePath := os.Getenv("TRANSACTION_LOG_FILE")
	if transactionFilePath == "" {
		transactionFilePath = "./transactions.log"
	}

	transactions, err := NewFileTransactionLogger(transactionFilePath)
	if err != nil {
		log.Fatal("failed to create event logger: ", err)
	}
	err = transactions.ReplayEvents()
	if err != nil {
		log.Fatal("failed to replay events: ", err)
	}

	r := mux.NewRouter()
	r.HandleFunc("/v1/kv/{key}", func(w http.ResponseWriter, r *http.Request) {
		putKeyValuePairHandler(w, r, transactions)
	}).Methods("PUT")

	r.HandleFunc("/v1/kv/{key}", func(w http.ResponseWriter, r *http.Request) {
		deleteKeyValuePairHandler(w, r, transactions)
	}).Methods("DELETE")

	r.HandleFunc("/v1/kv/{key}", getKeyValuePairHandler).Methods("GET")

	log.Fatal(http.ListenAndServe(":8080", r))
}
