package transistor

import (
	"fmt"
	"reflect"
	"runtime"
	"time"

	json "github.com/bww/go-json"
	log "github.com/codeamp/logger"
	"github.com/satori/go.uuid"
)

type Event struct {
	ID           uuid.UUID   `json:"id"`
	ParentID     uuid.UUID   `json:"parentId"`
	Name         string      `json:"name"`
	Payload      interface{} `json:"payload"`
	PayloadModel string      `json:"payloadModel"`
	Error        error       `json:"error"`
	CreatedAt    time.Time   `json:"createdAt"`
	Caller       Caller      `json:"caller"`
}

type Caller struct {
	File       string `json:"file"`
	LineNumber int    `json:"line_number"`
}

func name(payload interface{}) string {
	s := reflect.ValueOf(payload)

	if s.Kind() != reflect.Struct {
		return reflect.TypeOf(payload).String()
	}

	f := s.FieldByName("Action")
	if f.IsValid() {
		model := reflect.TypeOf(payload).String()
		action := f.String()
		if action != "" {
			return fmt.Sprintf("%v:%v", model, action)
		}
	}

	return reflect.TypeOf(payload).String()
}

func NewEvent(payload interface{}, err error) Event {
	event := Event{
		ID:           uuid.NewV4(),
		Name:         name(payload),
		Payload:      payload,
		PayloadModel: reflect.TypeOf(payload).String(),
		Error:        err,
		CreatedAt:    time.Now(),
	}

	// for debugging purposes
	_, file, no, ok := runtime.Caller(1)
	if ok {
		event.Caller = Caller{
			File:       file,
			LineNumber: no,
		}
	}

	return event
}

func (e *Event) NewEvent(payload interface{}, err error) Event {
	event := Event{
		ID:           uuid.NewV4(),
		ParentID:     e.ID,
		Name:         name(payload),
		Payload:      payload,
		PayloadModel: reflect.TypeOf(payload).String(),
		Error:        err,
		CreatedAt:    time.Now(),
	}

	// for debugging purposes
	_, file, no, ok := runtime.Caller(1)
	if ok {
		event.Caller = Caller{
			File:       file,
			LineNumber: no,
		}
	}

	return event
}

func (e *Event) Dump() {
	event, _ := json.MarshalRole("dummy", e)
	log.Info(string(event))
}
