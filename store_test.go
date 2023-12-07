package main

import (
	"math/rand"
	"sync"
	"testing"
	"time"
)

// TestThreadSafety tests the thread safety of the map operations
func TestThreadSafety(t *testing.T) {
	var wg sync.WaitGroup
	keyValuePairs := map[string]string{
		"key1": "value1",
		"key2": "value2",
		// ... Add more key-value pairs as needed
	}

	for i := 0; i < 100; i++ { // Repeat multiple times to increase chance of overlap
		for key, value := range keyValuePairs {
			wg.Add(1)
			go func(k, v string) {
				defer wg.Done()
				time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
				err := Put(k, v)
				if err != nil {
					t.Errorf("Failed to put key-value pair: %s, %s", k, v)
				}
			}(key, value)

			wg.Add(1)
			go func(k, expected string) {
				defer wg.Done()
				time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
				kvPair, err := GetKeyValuePair(k)
				if err != nil || kvPair.Value != expected {
					t.Errorf("Unexpected result. Expected: %s, Got: %s", k, kvPair.Value)
				}
			}(key, value)
		}
	}
	wg.Wait() // Wait for all operations to complete

	// Additional consistency checks can be performed here
}
