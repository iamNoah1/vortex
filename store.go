package main

import (
	"errors"
	"sync"
)

var store = struct {
	sync.RWMutex
	m map[string]string
}{m: make(map[string]string)}

var ErrorNoSuchKey = errors.New("no such key")

func GetKeyValuePair(key string) (KeyValuePair, error) {
	store.RLock()
	value, ok := store.m[key]
	store.RUnlock()

	if !ok {
		return KeyValuePair{}, ErrorNoSuchKey
	}

	return KeyValuePair{key, value}, nil
}

func Put(key, value string) error {
	store.Lock()
	store.m[key] = value
	store.Unlock()

	return nil
}

func Delete(key string) error {
	store.Lock()
	delete(store.m, key)
	store.Unlock()

	return nil
}
