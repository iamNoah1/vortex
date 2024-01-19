package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type KeyValuePair struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

var logger TransactionLogger

func main() {
	logger, err := NewFileTransactionLogger("transaction.log")
	if err != nil {
		log.Fatal("failed to create event logger: ", err)
	}
	err = logger.ReplayEvents()
	if err != nil {
		log.Fatal("failed to replay events: ", err)
	}

	r := mux.NewRouter()
	r.HandleFunc("/v1/kv/{key}", func(w http.ResponseWriter, r *http.Request) {
		putKeyValuePairHandler(w, r, logger)
	}).Methods("PUT")

	r.HandleFunc("/v1/kv/{key}", func(w http.ResponseWriter, r *http.Request) {
		deleteKeyValuePairHandler(w, r, logger)
	}).Methods("DELETE")

	r.HandleFunc("/v1/kv/{key}", getKeyValuePairHandler).Methods("GET")

	log.Fatal(http.ListenAndServe(":8080", r))
}
