package store

import (
	"sync"
	"testing"
)

func TestReadWriteNotRacing(t *testing.T) {
	var wg sync.WaitGroup
	keyValuePairs := map[string]string{
		"key1":                  "value1",
		"key2":                  "value2",
		"user123":               "John Doe",
		"email123":              "johndoe@example.com",
		"productID45":           "Laptop",
		"order456":              "Completed",
		"sessionID9a2b":         "Active",
		"config:timeout":        "30s",
		"config:maxConnections": "100",
		"status:server1":        "Online",
		"lastLogin:user123":     "2023-03-15",
		"itemCount:cart456":     "5",
		"price:productID45":     "$999",
		"location:office":       "Building 3, Floor 2",
		"mode:system":           "Auto",
	}

	// Pre-populate the map
	for key, value := range keyValuePairs {
		err := Put(key, value)
		if err != nil {
			t.Fatalf("Failed to pre-populate map: %s", err)
		}
	}

	// Perform concurrent reads and writes
	for i := 0; i < 100; i++ {
		for key, value := range keyValuePairs {
			wg.Add(1)
			go func(k, v string) {
				defer wg.Done()
				err := Put(k, v)
				if err != nil {
					t.Errorf("Failed to put key-value pair: %s, %s", k, v)
				}
			}(key, value)

			wg.Add(1)
			go func(k, expected string) {
				defer wg.Done()
				kvPair, err := GetKeyValuePair(k)
				if err != nil || kvPair.Value != expected {
					t.Errorf("Unexpected result for key: %s", k)
				}
			}(key, value)
		}
	}
	wg.Wait()
}

func TestWriteDeleteNotRacing(t *testing.T) {
	var wg sync.WaitGroup
	keyValuePairs := map[string]string{
		"email123":              "johndoe@example.com",
		"productID45":           "Laptop",
		"order456":              "Completed",
		"sessionID9a2b":         "Active",
		"config:timeout":        "30s",
		"config:maxConnections": "100",
		"status:server1":        "Online",
		"lastLogin:user123":     "2023-03-15",
		"itemCount:cart456":     "5",
		"price:productID45":     "$999",
		"location:office":       "Building 3, Floor 2",
		"mode:system":           "Auto",
	}

	for key, value := range keyValuePairs {
		err := Put(key, value)
		if err != nil {
			t.Fatalf("Failed to pre-populate map: %s", err)
		}
	}

	// Perform concurrent writes and deletes
	for i := 0; i < 100; i++ {
		for key, value := range keyValuePairs {
			wg.Add(1)
			go func(k, v string) {
				defer wg.Done()
				err := Put(k, v)
				if err != nil {
					t.Errorf("Failed to put key-value pair: %s, %s", k, v)
				}
			}(key, value)

			wg.Add(1)
			go func(k string) {
				defer wg.Done()

				err := Delete(k)
				if err != nil {
					t.Errorf("Error deleting: %s", k)
				}
			}(key)
		}
	}
	wg.Wait()
}

func TestGetAllKeyValuePairs(t *testing.T) {
	Clear()
	keyValuePairs := map[string]string{
		"email123":       "noah@github.com",
		"productID45":    "Laptop",
		"order456":       "Completed",
		"sessionID9a2b":  "Active",
		"config:timeout": "30s",
	}

	for key, value := range keyValuePairs {
		err := Put(key, value)
		if err != nil {
			t.Fatalf("Failed to pre-populate map: %s", err)
		}
	}

	kvs, err := GetAllKeyValuePairs()
	if err != nil {
		t.Fatalf("Failed to get all key-value pairs: %s", err)
	}

	if len(kvs) != len(keyValuePairs) {
		t.Fatalf("Expected %d key-value pairs; got %d", len(keyValuePairs), len(kvs))
	}

	for _, kv := range kvs {
		if keyValuePairs[kv.Key] != kv.Value {
			t.Errorf("Unexpected value for key: %s", kv.Key)
		}
	}
}
