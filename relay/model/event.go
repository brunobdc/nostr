package model

import (
	"bytes"
	"strconv"
)

type EventKind uint16

const (
	EventKindMetaData   = 0
	EventKindTextNote   = 1
	EventKindFollowList = 3
)

type Event struct {
	ID        string    `json:"id" bson:"_id"`
	PublicKey string    `json:"pubKey"`
	CreatedAt uint64    `json:"created_at"`
	Kind      EventKind `json:"kind"`
	Tags      Tags      `json:"tags"`
	Content   string    `json:"content"`
	Signature string    `json:"sig"`
}

func NewEvent() *Event {
	return &Event{Tags: make(Tags)}
}

func (e *Event) MarshalJSON() ([]byte, error) {
	var buffer bytes.Buffer
	// id
	buffer.WriteString("{\"id\":\"")
	buffer.WriteString(e.ID)
	// pubKey
	buffer.WriteString("\",\"pubKey\":\"")
	buffer.WriteString(e.PublicKey)
	// created_at
	buffer.WriteString("\",\"created_at\":")
	buffer.WriteString(strconv.FormatUint(e.CreatedAt, 10))
	// kind
	buffer.WriteString(",\"kind\":")
	buffer.WriteString(strconv.FormatUint(uint64(e.Kind), 10))
	// tags
	buffer.WriteString(",\"tags\":")
	tags, err := e.Tags.MarshalJSON()
	if err != nil {
		return nil, err
	}
	buffer.Write(tags)
	// content
	buffer.WriteString(",\"content\":\"")
	buffer.WriteString(e.Content)
	// sig
	buffer.WriteString("\",\"sig\":\"")
	buffer.WriteString(e.Signature)
	buffer.WriteString("\"}")

	return buffer.Bytes(), nil
}
