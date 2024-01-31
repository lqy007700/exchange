package internal

type Event struct {
	Type EventType `json:"type"`
	Data any       `json:"data"`
}

type EventType string

const (
	Create EventType = "create"
	Close  EventType = "close"
)
