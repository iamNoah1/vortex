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

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/v1/kv/{key}", putKeyValuePairHandler).Methods("PUT")
	r.HandleFunc("/v1/kv/{key}", getKeyValuePairHandler).Methods("GET")
	r.HandleFunc("/v1/kv/{key}", deleteKeyValuePairHandler).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":8080", r))
}
