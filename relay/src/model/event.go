package model

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"strconv"

	"github.com/brunobdc/nostr/relay/src/schnorr"
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

func (e Event) IsValid() (bool, string) {
	bytesID, err := hex.DecodeString(e.ID)
	if err != nil {
		return false, "error: couldn't decode id hexadecimal"
	}
	if len(bytesID) != 32 {
		return false, "invalid: id bytes size is differente from 32"
	}
	bytesPubKey, err := hex.DecodeString(e.PublicKey)
	if err != nil {
		return false, "error: couldn't decode pubKey hexadecimal"
	}
	if len(bytesPubKey) != 32 {
		return false, "invalid: pubKey bytes size is differente from 32"
	}
	bytesSignature, err := hex.DecodeString(e.Signature)
	if err != nil {
		return false, "error: couldn't decode signature hexadecimal"
	}
	if len(bytesSignature) != 64 {
		return false, "invalid: signature bytes size is differente from 64"
	}

	var buffer bytes.Buffer
	buffer.WriteString("[0,\"")
	buffer.WriteString(e.PublicKey)
	buffer.WriteString("\",")
	buffer.WriteString(strconv.FormatUint(e.CreatedAt, 10))
	buffer.WriteString(",")
	buffer.WriteString(strconv.FormatUint(uint64(e.Kind), 10))
	buffer.WriteString(",")
	jsonTags, err := e.Tags.MarshalJSON()
	if err != nil {
		return false, "error: couldn't marshal the tags"
	}
	buffer.Write(jsonTags)
	buffer.WriteString(",\"")
	buffer.WriteString(e.Content)
	buffer.WriteString("\"]")

	hash := sha256.Sum256(buffer.Bytes())
	if [32]byte(bytesID) != hash {
		return false, "invalid: id hash doesn't match the data"
	}

	valid := schnorr.Verify(bytesSignature, bytesID, bytesPubKey)
	if !valid {
		return false, "invalid: signature doesn't match"
	}
	return true, ""
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
	buffer.WriteString(",\"tags\":[")
	tags, err := e.Tags.MarshalJSON()
	if err != nil {
		return nil, err
	}
	buffer.Write(tags)
	// content
	buffer.WriteString("],\"content\":\"")
	buffer.WriteString(e.Content)
	// sig
	buffer.WriteString("\",\"sig\":\"")
	buffer.WriteString(e.Signature)
	buffer.WriteString("\"}")

	return buffer.Bytes(), nil
}
