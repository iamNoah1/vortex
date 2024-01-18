package main

import (
	"fmt"
	"os"
)

type Event struct {
	Sequence  uint64
	EventType EventType
	Key       string
	Value     *string
}

type EventType uint64

const (
	EventDelete EventType = iota
	EventPut
)

type TransactionLogger interface {
	Put(key string, value string) error
	Delete(key string) error
}

type FileTransactionLogger struct {
	fileName string
}

func NewFileTransactionLogger(filename string) (TransactionLogger, error) {
	return &FileTransactionLogger{fileName: filename}, nil
}

func (f *FileTransactionLogger) Put(key string, value string) error {
	return AppendToFile(f.fileName, fmt.Sprintf("%d\t%d\t%s\t%s", EventPut, key, value))
}

func (f *FileTransactionLogger) Delete(key string) error {
	return nil
}

func AppendToFile(fileName string, text string) error {
	file, err := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	_, err = file.WriteString("\n" + text)
	if err != nil {
		return fmt.Errorf("failed to append to file: %w", err)
	}

	return nil
}
