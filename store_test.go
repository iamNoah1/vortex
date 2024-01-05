package main

import (
	"sync"
	"testing"
)

// TestThreadSafety tests the thread safety of the map operations
func TestThreadSafety(t *testing.T) {
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
		// ... Add more key-value pairs as needed
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
