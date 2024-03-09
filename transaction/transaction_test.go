package transaction

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
	"os/exec"
	"testing"
)

func TestReplayEvents(t *testing.T) {
	cmd := startServer(t)

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

	stopServer(cmd)

	// Start the server again and check if the value was persisted
	cmd = startServer(t)

	req, err = http.NewRequest("GET", "http://localhost:8080/v1/kv/key1", nil)
	if err != nil {
		t.Fatalf("Failed to create GET request: %v", err)
	}

	resp, err = client.Do(req)
	if err != nil {
		t.Fatalf("Failed to send GET request: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status OK; got %v", resp.StatusCode)
	}

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)

	value := result["value"].(string)
	if value != "myValue" {
		t.Fatalf("Expected 'myValue'; got %v", value)
	}

	stopServer(cmd)

	t.Cleanup(func() {
		stopServer(cmd)
		os.Remove("transaction_test.log")
	})
}

func startServer(t *testing.T) *exec.Cmd {
	cmd := exec.Command("go", "run", ".")
	err := os.Setenv("TRANSACTION_LOG_FILE", "transaction_test.log")
	if err != nil {
		t.Fatalf("Failed to set environment variable: %v", err)
	}

	err = cmd.Start()
	if err != nil {
		t.Fatalf("cmd.Start() failed with '%s'\n", err)
	}

	return cmd
}

func stopServer(cmd *exec.Cmd) {
	err := cmd.Process.Kill()
	if err != nil {
		panic(err)
	}
}
