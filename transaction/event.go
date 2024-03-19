package transaction

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/rs/zerolog/log"
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
	log.Debug().Msgf("Incoming Event: %s", eventStr)

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
