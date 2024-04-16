package transaction

import (
	"fmt"

	"github.com/iamNoah1/vortex/store"
	"github.com/iamNoah1/vortex/utils"
)

type FileTransactionLogger struct {
	fileName string
	sequence uint64
}

func NewFileTransactionLogger(filename string) (TransactionLogger, error) {
	return &FileTransactionLogger{fileName: filename, sequence: 0}, nil
}

func (f *FileTransactionLogger) Put(key string, value string) error {
	f.sequence++
	return utils.AppendToFile(f.fileName, fmt.Sprintf("%d,%s,%s,%s", f.sequence, EventPut.String(), key, value))
}

func (f *FileTransactionLogger) Delete(key string) error {
	f.sequence++
	return utils.AppendToFile(f.fileName, fmt.Sprintf("%d,%s,%s,%s", f.sequence, EventDelete.String(), key, ""))
}

func (f *FileTransactionLogger) ReplayEvents() error {
	lines, err := utils.ReadFromFile(f.fileName)

	if err != nil {
		return fmt.Errorf("failed to read from file: %w", err)
	}

	for _, line := range lines {
		event, err := NewEventFromString(line)
		if err != nil {
			return fmt.Errorf("failed to parse line: %w", err)
		}

		if event.EventType == EventPut {
			err = store.Put(event.Key, *event.Value)
		}

		if err != nil {
			return fmt.Errorf("failed to replay event: %w", err)
		}
	}

	return nil
}
