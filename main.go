package main

import (
	"encoding/json"
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

func putKeyValuePairHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]

	var kv KeyValuePair
	err := json.NewDecoder(r.Body).Decode(&kv)
	defer r.Body.Close()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = Put(key, kv.Value)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func getKeyValuePairHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]

	kv, err := GetKeyValuePair(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(kv)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func deleteKeyValuePairHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]

	err := Delete(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Key deleted successfully"))
}
