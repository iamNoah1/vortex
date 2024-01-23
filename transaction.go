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

func (s EventType) String() string {
	return [...]string{"Delete", "Put"}[s]
}

type TransactionLogger interface {
	Put(key string, value string) error
	Delete(key string) error
	ReplayEvents() error
}

type FileTransactionLogger struct {
	fileName string
	sequence uint64
}

func NewFileTransactionLogger(filename string) (TransactionLogger, error) {
	return &FileTransactionLogger{fileName: filename, sequence: 0}, nil
}

func (f *FileTransactionLogger) Put(key string, value string) error {
	f.sequence++
	return AppendToFile(f.fileName, fmt.Sprintf("%d,%s,%s,%s", f.sequence, EventPut.String(), key, value))
}

func (f *FileTransactionLogger) Delete(key string) error {
	f.sequence++
	return AppendToFile(f.fileName, fmt.Sprintf("%d,%s,%s,%s", f.sequence, EventDelete.String(), key, ""))
}

func (f *FileTransactionLogger) ReplayEvents() error {
	return nil
}

func AppendToFile(fileName string, text string) error {
	file, err := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	_, err = file.WriteString(text + "\n")
	if err != nil {
		return fmt.Errorf("failed to append to file: %w", err)
	}

	return nil
}
