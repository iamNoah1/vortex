package main

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gorilla/mux"
)

func putKeyValuePairHandler(w http.ResponseWriter, r *http.Request, logger TransactionLogger) {
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

	kv, err := GetKeyValuePair(key)
	if err != nil {
		if errors.Is(err, ErrorNoSuchKey) {
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

func deleteKeyValuePairHandler(w http.ResponseWriter, r *http.Request, logger TransactionLogger) {
	vars := mux.Vars(r)
	key := vars["key"]

	err := Delete(key)
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
