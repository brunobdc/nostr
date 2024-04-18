package helpers

import (
	"bytes"
	"strconv"

	"github.com/brunobdc/nostr/relay/model"
)

func MakeEventResponse(subscriptionID string, event model.Event) ([]byte, error) {
	var buffer bytes.Buffer
	buffer.WriteString("[\"EVENT\",\"")
	buffer.WriteString(subscriptionID)
	buffer.WriteString("\",")
	eventJson, err := event.MarshalJSON()
	if err != nil {
		return nil, err
	}
	buffer.Write(eventJson)
	buffer.WriteString("]")
	return buffer.Bytes(), nil
}

func MakeOkResponse(eventID string, ok bool, msg string) ([]byte, error) {
	var buffer bytes.Buffer
	buffer.WriteString("[\"OK\",\"")
	buffer.WriteString(eventID)
	buffer.WriteString("\",")
	buffer.WriteString(strconv.FormatBool(ok))
	buffer.WriteString(",\"")
	buffer.WriteString(msg)
	buffer.WriteString("\"]")
	return buffer.Bytes(), nil
}

func MakeEoseResponse(subscriptionID string) ([]byte, error) {
	var buffer bytes.Buffer
	buffer.WriteString("[\"EOSE\",\"")
	buffer.WriteString(subscriptionID)
	buffer.WriteString("\"]")
	return buffer.Bytes(), nil
}
