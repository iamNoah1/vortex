package main

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/iamNoah1/vortex/store"
	"github.com/iamNoah1/vortex/transaction"
)

func putKeyValuePairHandler(w http.ResponseWriter, r *http.Request, logger transaction.TransactionLogger) {
	vars := mux.Vars(r)
	key := vars["key"]

	var kv store.KeyValuePair
	err := json.NewDecoder(r.Body).Decode(&kv)
	defer r.Body.Close()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = store.Put(key, kv.Value)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = logger.Put(key, kv.Value)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func getKeyValuePairHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]

	kv, err := store.GetKeyValuePair(key)
	if err != nil {
		if errors.Is(err, store.ErrorNoSuchKey) {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(kv)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func getAllKeyValuePairsHandler(w http.ResponseWriter, r *http.Request) {
	kvs, err := store.GetAllKeyValuePairs()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(kvs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func deleteKeyValuePairHandler(w http.ResponseWriter, r *http.Request, logger transaction.TransactionLogger) {
	vars := mux.Vars(r)
	key := vars["key"]

	err := store.Delete(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = logger.Delete(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Key deleted successfully"))
}
