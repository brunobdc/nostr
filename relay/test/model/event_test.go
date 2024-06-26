package model_test

import (
	"encoding/hex"
	"encoding/json"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/brunobdc/nostr/relay/model"
	"github.com/stretchr/testify/assert"
)

func TestEventMarshalJSON(T *testing.T) {
	T.Parallel()

	T.Run("Should return a valid json string", func(t *testing.T) {
		t.Parallel()
		var bytesID [32]byte
		gofakeit.Slice(&bytesID)
		id := make([]byte, hex.EncodedLen(len(bytesID)))
		hex.Encode(id, bytesID[:])
		var bytesPubKey [32]byte
		gofakeit.Slice(&bytesPubKey)
		pubKey := make([]byte, hex.EncodedLen(len(bytesPubKey)))
		hex.Encode(pubKey, bytesPubKey[:])
		var bytesSig [64]byte
		gofakeit.Slice(bytesSig)
		sig := make([]byte, hex.EncodedLen(len(bytesSig)))
		hex.Encode(sig, bytesSig[:])

		var stringMatrix [][]string
		gofakeit.Slice(&stringMatrix)

		tags := make(model.Tags)
		for _, values := range stringMatrix {
			tags[values[0]] = values[1:]
		}

		event := model.Event{
			ID:        string(id),
			PublicKey: string(pubKey),
			CreatedAt: uint64(gofakeit.Date().Unix()),
			Kind:      model.EventKind(gofakeit.Uint16()),
			Tags:      tags,
			Content:   gofakeit.LoremIpsumSentence(gofakeit.IntN(100)),
			Signature: string(sig),
		}

		jsonString, err := event.MarshalJSON()
		assert.Nil(t, err)

		unmarshaledEvent := model.Event{Tags: make(model.Tags)}
		err = json.Unmarshal(jsonString, &unmarshaledEvent)
		assert.Nil(t, err)

		assert.Equal(t, event, unmarshaledEvent)
	})
}
