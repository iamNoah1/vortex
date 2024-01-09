package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

func TestPutKeyValuePairHandler(t *testing.T) {
	req, err := http.NewRequest("PUT", "/v1/kv/test_key", bytes.NewBuffer([]byte(`{"value":"test_value"}`)))
	if err != nil {
		t.Fatal(err)
	}

	responseRecorder := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/v1/kv/{key}", putKeyValuePairHandler).Methods("PUT")

	router.ServeHTTP(responseRecorder, req)

	if status := responseRecorder.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
	}
}

func TestGetKeyValuePairHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/v1/kv/test_key", nil)
	if err != nil {
		t.Fatal(err)
	}

	responseRecorder := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/v1/kv/{key}", getKeyValuePairHandler).Methods("GET")

	router.ServeHTTP(responseRecorder, req)

	if status := responseRecorder.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}

func TestGetKeyValuePairHandlerKeyNotFound(t *testing.T) {
	req, err := http.NewRequest("GET", "/v1/kv/non_existing_key", nil)
	if err != nil {
		t.Fatal(err)
	}

	responseRecorder := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/v1/kv/{key}", getKeyValuePairHandler).Methods("GET")

	router.ServeHTTP(responseRecorder, req)

	if status := responseRecorder.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNotFound)
	}
}

func TestDeleteKeyValuePairHandler(t *testing.T) {
	req, err := http.NewRequest("DELETE", "/v1/kv/test_key", nil)
	if err != nil {
		t.Fatal(err)
	}

	responseRecorder := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/v1/kv/{key}", deleteKeyValuePairHandler).Methods("DELETE")

	router.ServeHTTP(responseRecorder, req)

	if status := responseRecorder.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Add more assertions based on your application's logic
}
