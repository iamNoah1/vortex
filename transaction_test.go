package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
	"os/exec"
	"testing"
)

func TestMain(t *testing.T) {
	cmd := exec.Command("go", "run", ".")

	err := os.Setenv("TRANSACTION_LOG_FILE", "transaction_test.log")
	if err != nil {
		t.Fatalf("Failed to set environment variable: %v", err)
	}

	err = cmd.Start()
	if err != nil {
		t.Fatalf("cmd.Start() failed with '%s'\n", err)
	}

	data := map[string]interface{}{
		"value": "myValue",
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		t.Fatalf("Failed to marshal JSON: %v", err)
	}

	req, err := http.NewRequest("PUT", "http://localhost:8080/v1/kv/key1", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatalf("Failed to create PUT request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Failed to send PUT request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("Expected status OK; got %v", resp.StatusCode)
	}

	err = cmd.Process.Kill()
	if err != nil {
		t.Fatalf("cmd.Process.Kill() failed with '%s'\n", err)
	}
}
