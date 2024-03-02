package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
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

func (e *Event) ToString() string {
	return fmt.Sprintf("%d,%s,%s,%s", e.Sequence, e.EventType.String(), e.Key, *e.Value)
}

func NewEventFromString(eventStr string) (*Event, error) {
	parts := strings.Split(eventStr, ",")
	if len(parts) != 4 {
		return nil, fmt.Errorf("invalid event string")
	}

	id, err := strconv.Atoi(parts[0])
	if err != nil {
		return nil, fmt.Errorf("invalid ID: %v", err)
	}

	var eventType EventType
	switch parts[1] {
	case "Put":
		eventType = EventPut
	case "Delete":
		eventType = EventDelete
	default:
		return nil, fmt.Errorf("invalid event type")
	}

	return &Event{
		Sequence:  uint64(id),
		EventType: eventType,
		Key:       parts[2],
		Value:     &parts[3],
	}, nil
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
	lines, err := ReadFromFile(f.fileName)

	if err != nil {
		return fmt.Errorf("failed to read from file: %w", err)
	}

	for _, line := range lines {
		event, err := NewEventFromString(line)
		if err != nil {
			return fmt.Errorf("failed to parse line: %w", err)
		}

		if event.EventType == EventPut {
			err = Put(event.Key, *event.Value)
		}

		if err != nil {
			return fmt.Errorf("failed to replay event: %w", err)
		}
	}

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

func ReadFromFile(fileName string) ([]string, error) {
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_RDONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return lines, nil
}
