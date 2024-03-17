package kamu

import (
	"bytes"
	"encoding/json"
	"log"
)

type requestCommand string

const (
	commandPost  requestCommand = "POST"
	commandStart requestCommand = "START"
)

type clientTimezone string

const (
	clientTimezoneAsiaBangkok clientTimezone = "Asia/Bangkok"
)

type requestType string

const (
	requestTypeActionLink requestType = "action_link"
	requestTypeText       requestType = "text"
)

type request struct {
	Command          requestCommand `json:"command"`
	FilterValues     []string       `json:"filter_values,omitempty"`
	Language         string         `json:"language,omitempty"`
	DisableHumanChat bool           `json:"disable_human_chat,omitempty"`

	Type           requestType    `json:"type,omitempty"`
	ID             string         `json:"id,omitempty"`
	ConversationID string         `json:"conversation_id,omitempty"`
	ClientTimezone clientTimezone `json:"client_timezone,omitempty"`

	Value string `json:"value,omitempty"`
}

func (r *request) NewReader() *bytes.Reader {
	if marshalled, err := json.Marshal(r); err != nil {
		log.Fatalf("impossible to marshall request: %s", err)
	} else {
		return bytes.NewReader(marshalled)
	}

	return nil
}

func (r *request) Read(p []byte) (n int, err error) {
	return r.NewReader().Read(p)
}
