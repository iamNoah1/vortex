package main

import "errors"

var store = make(map[string]string)

var ErrorNoSuchKey = errors.New("no such key")

func GetValue(key string) (string, error) {
	value, ok := store[key]

	if !ok {
		return "", ErrorNoSuchKey
	}

	return value, nil
}

func GetKeyValuePair(key string) (KeyValuePair, error) {
	value, ok := store[key]

	if !ok {
		return KeyValuePair{}, ErrorNoSuchKey
	}

	return KeyValuePair{key, value}, nil
}

func Put(key, value string) error {
	store[key] = value
	return nil
}

func Delete(key string) error {
	delete(store, key)
	return nil
}
