package model

import (
	"os"
	"time"
)

type Event struct {
	EventID   string    `json:"event_id" yaml:"event_id" xml:"event_id"`
	Type      string    `json:"type" yaml:"type" xml:"type"`
	Timestamp time.Time `json:"timestamp" yaml:"timestamp" xml:"timestamp"`
	Payload   Payload   `json:"payload" yaml:"payload" xml:"payload"`
}

type Payload struct {
	Amount        int      `json:"amount" yaml:"amount" xml:"amount"`
	Currency      string   `json:"currency" yaml:"currency" xml:"currency"`
	OriginalEvent string   `json:"original_event,omitempty" yaml:"original_event,omitempty" xml:"original_event,omitempty"`
	Details       *Details `json:"details,omitempty" yaml:"details,omitempty" xml:"details,omitempty"` // необязательное поле, а указатель на nil просто не существует, по этому записи даже нет (xml, yaml не любят пустые поля)
}

type Details struct {
	Method string `json:"method" yaml:"method" xml:"method"`
	Card   *Card  `json:"card,omitempty" yaml:"card,omitempty" xml:"card,omitempty"`
}

type Card struct {
	Last4  string `json:"last4" yaml:"last4" xml:"last4"`
	Issuer string `json:"issuer" yaml:"issuer" xml:"issuer"`
}

type FileInfo struct {
	Info *os.File
	Path string
	Data []byte
}

type Amm struct {
	Ammount float64
	Event
}

var Amms []Amm
