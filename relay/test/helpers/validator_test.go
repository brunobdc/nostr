package helpers_test

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"strconv"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/brunobdc/nostr/relay/helpers"
	"github.com/brunobdc/nostr/relay/model"
	security_test "github.com/brunobdc/nostr/relay/test/security"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestValidateEvent(T *testing.T) {
	T.Parallel()

	T.Run("Should return false if hex encoded ID is not [32]byte", func(t *testing.T) {
		t.Parallel()

		var bytesID [30]byte
		gofakeit.Slice(&bytesID)
		id := make([]byte, hex.EncodedLen(len(bytesID)))
		hex.Encode(id, bytesID[:])

		event := model.Event{
			ID: string(id),
		}

		expectedFalse, msg, err := helpers.ValidateEvent(event, &security_test.MockSignature{})

		assert.Nil(t, err)
		assert.False(t, expectedFalse)
		assert.Equal(t, msg, "invalid: id bytes size is differente from 32")
	})

	T.Run("Should return false if hex encoded Public Key is not [32]byte", func(t *testing.T) {
		t.Parallel()

		var bytesID [32]byte
		gofakeit.Slice(&bytesID)
		id := make([]byte, hex.EncodedLen(len(bytesID)))
		hex.Encode(id, bytesID[:])

		var bytesPubKey [30]byte
		gofakeit.Slice(&bytesPubKey)
		pubKey := make([]byte, hex.EncodedLen(len(bytesPubKey)))
		hex.Encode(pubKey, bytesPubKey[:])

		event := model.Event{
			ID:        string(id),
			PublicKey: string(pubKey),
		}

		expectedFalse, msg, err := helpers.ValidateEvent(event, &security_test.MockSignature{})

		assert.Nil(t, err)
		assert.False(t, expectedFalse)
		assert.Equal(t, msg, "invalid: pubKey bytes size is differente from 32")
	})

	T.Run("Should return false if hex encoded Signature is not [64]byte", func(t *testing.T) {
		t.Parallel()

		var bytesID [32]byte
		gofakeit.Slice(&bytesID)
		id := make([]byte, hex.EncodedLen(len(bytesID)))
		hex.Encode(id, bytesID[:])

		var bytesPubKey [32]byte
		gofakeit.Slice(&bytesPubKey)
		pubKey := make([]byte, hex.EncodedLen(len(bytesPubKey)))
		hex.Encode(pubKey, bytesPubKey[:])

		var bytesSig [52]byte
		gofakeit.Slice(bytesSig)
		sig := make([]byte, hex.EncodedLen(len(bytesSig)))
		hex.Encode(sig, bytesSig[:])

		event := model.Event{
			ID:        string(id),
			PublicKey: string(pubKey),
			Signature: string(sig),
		}

		expectedFalse, msg, err := helpers.ValidateEvent(event, &security_test.MockSignature{})

		assert.Nil(t, err)
		assert.False(t, expectedFalse)
		assert.Equal(t, msg, "invalid: signature bytes size is differente from 64")
	})

	T.Run("Should return false if id hash doesn't match", func(t *testing.T) {
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

		event := model.Event{
			ID:        string(id),
			PublicKey: string(pubKey),
			Signature: string(sig),
			Tags:      make(model.Tags),
			CreatedAt: uint64(gofakeit.Date().Unix()),
			Kind:      model.EventKind(gofakeit.Uint16()),
			Content:   gofakeit.LoremIpsumSentence(int(gofakeit.Uint16())),
		}

		expectedFalse, msg, err := helpers.ValidateEvent(event, &security_test.MockSignature{})

		assert.Nil(t, err)
		assert.False(t, expectedFalse)
		assert.Equal(t, msg, "invalid: id hash doesn't match the data")
	})

	T.Run("Should return false if id hash doesn't match", func(t *testing.T) {
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

		event := model.Event{
			ID:        string(id),
			PublicKey: string(pubKey),
			Signature: string(sig),
			Tags:      make(model.Tags),
			CreatedAt: uint64(gofakeit.Date().Unix()),
			Kind:      model.EventKind(gofakeit.Uint16()),
			Content:   gofakeit.LoremIpsumSentence(int(gofakeit.Uint16())),
		}

		expectedFalse, msg, err := helpers.ValidateEvent(event, &security_test.MockSignature{})

		assert.Nil(t, err)
		assert.False(t, expectedFalse)
		assert.Equal(t, msg, "invalid: id hash doesn't match the data")
	})

	T.Run("Should return false if signature verify returns false", func(t *testing.T) {
		t.Parallel()

		var bytesPubKey [32]byte
		gofakeit.Slice(&bytesPubKey)
		pubKey := make([]byte, hex.EncodedLen(len(bytesPubKey)))
		hex.Encode(pubKey, bytesPubKey[:])

		var bytesSig [64]byte
		gofakeit.Slice(bytesSig)
		sig := make([]byte, hex.EncodedLen(len(bytesSig)))
		hex.Encode(sig, bytesSig[:])

		event := model.Event{
			PublicKey: string(pubKey),
			Signature: string(sig),
			Tags:      make(model.Tags),
			CreatedAt: uint64(gofakeit.Date().Unix()),
			Kind:      model.EventKind(gofakeit.Uint16()),
			Content:   gofakeit.LoremIpsumSentence(int(gofakeit.Uint16())),
		}

		var buffer bytes.Buffer
		buffer.WriteString("[0,\"")
		buffer.WriteString(event.PublicKey)
		buffer.WriteString("\",")
		buffer.WriteString(strconv.FormatUint(event.CreatedAt, 10))
		buffer.WriteString(",")
		buffer.WriteString(strconv.FormatUint(uint64(event.Kind), 10))
		buffer.WriteString(",")
		jsonTags, _ := event.Tags.MarshalJSON()
		buffer.Write(jsonTags)
		buffer.WriteString(",\"")
		buffer.WriteString(event.Content)
		buffer.WriteString("\"]")

		hash := sha256.Sum256(buffer.Bytes())
		event.ID = hex.EncodeToString(hash[:])

		mockSignature := &security_test.MockSignature{}
		mockSignature.On("VerifySignature", mock.Anything, mock.Anything, mock.Anything).Return(false, nil)

		expectedFalse, msg, err := helpers.ValidateEvent(event, mockSignature)

		assert.Nil(t, err)
		assert.False(t, expectedFalse)
		assert.Equal(t, msg, "invalid: signature doesn't match")
	})

	T.Run("Should return false if id hash doesn't match", func(t *testing.T) {
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

		event := model.Event{
			ID:        string(id),
			PublicKey: string(pubKey),
			Signature: string(sig),
			Tags:      make(model.Tags),
			CreatedAt: uint64(gofakeit.Date().Unix()),
			Kind:      model.EventKind(gofakeit.Uint16()),
			Content:   gofakeit.LoremIpsumSentence(int(gofakeit.Uint16())),
		}

		expectedFalse, msg, err := helpers.ValidateEvent(event, &security_test.MockSignature{})

		assert.Nil(t, err)
		assert.False(t, expectedFalse)
		assert.Equal(t, msg, "invalid: id hash doesn't match the data")
	})

	T.Run("Should return error if signature verify returns an error", func(t *testing.T) {
		t.Parallel()

		var bytesPubKey [32]byte
		gofakeit.Slice(&bytesPubKey)
		pubKey := make([]byte, hex.EncodedLen(len(bytesPubKey)))
		hex.Encode(pubKey, bytesPubKey[:])

		var bytesSig [64]byte
		gofakeit.Slice(bytesSig)
		sig := make([]byte, hex.EncodedLen(len(bytesSig)))
		hex.Encode(sig, bytesSig[:])

		event := model.Event{
			PublicKey: string(pubKey),
			Signature: string(sig),
			Tags:      make(model.Tags),
			CreatedAt: uint64(gofakeit.Date().Unix()),
			Kind:      model.EventKind(gofakeit.Uint16()),
			Content:   gofakeit.LoremIpsumSentence(int(gofakeit.Uint16())),
		}

		var buffer bytes.Buffer
		buffer.WriteString("[0,\"")
		buffer.WriteString(event.PublicKey)
		buffer.WriteString("\",")
		buffer.WriteString(strconv.FormatUint(event.CreatedAt, 10))
		buffer.WriteString(",")
		buffer.WriteString(strconv.FormatUint(uint64(event.Kind), 10))
		buffer.WriteString(",")
		jsonTags, _ := event.Tags.MarshalJSON()
		buffer.Write(jsonTags)
		buffer.WriteString(",\"")
		buffer.WriteString(event.Content)
		buffer.WriteString("\"]")

		hash := sha256.Sum256(buffer.Bytes())
		event.ID = hex.EncodeToString(hash[:])

		returnError := gofakeit.Error()
		mockSignature := &security_test.MockSignature{}
		mockSignature.On("VerifySignature", mock.Anything, mock.Anything, mock.Anything).Return(false, returnError)

		expectedFalse, msg, err := helpers.ValidateEvent(event, mockSignature)

		assert.Equal(t, returnError, err)
		assert.False(t, expectedFalse)
		assert.Equal(t, msg, "")
	})

	T.Run("Should return false if signature verify returns false", func(t *testing.T) {
		t.Parallel()

		var bytesPubKey [32]byte
		gofakeit.Slice(&bytesPubKey)
		pubKey := make([]byte, hex.EncodedLen(len(bytesPubKey)))
		hex.Encode(pubKey, bytesPubKey[:])

		var bytesSig [64]byte
		gofakeit.Slice(bytesSig)
		sig := make([]byte, hex.EncodedLen(len(bytesSig)))
		hex.Encode(sig, bytesSig[:])

		event := model.Event{
			PublicKey: string(pubKey),
			Signature: string(sig),
			Tags:      make(model.Tags),
			CreatedAt: uint64(gofakeit.Date().Unix()),
			Kind:      model.EventKind(gofakeit.Uint16()),
			Content:   gofakeit.LoremIpsumSentence(int(gofakeit.Uint16())),
		}

		var buffer bytes.Buffer
		buffer.WriteString("[0,\"")
		buffer.WriteString(event.PublicKey)
		buffer.WriteString("\",")
		buffer.WriteString(strconv.FormatUint(event.CreatedAt, 10))
		buffer.WriteString(",")
		buffer.WriteString(strconv.FormatUint(uint64(event.Kind), 10))
		buffer.WriteString(",")
		jsonTags, _ := event.Tags.MarshalJSON()
		buffer.Write(jsonTags)
		buffer.WriteString(",\"")
		buffer.WriteString(event.Content)
		buffer.WriteString("\"]")

		hash := sha256.Sum256(buffer.Bytes())
		event.ID = hex.EncodeToString(hash[:])

		mockSignature := &security_test.MockSignature{}
		mockSignature.On("VerifySignature", mock.Anything, mock.Anything, mock.Anything).Return(false, nil)

		expectedFalse, msg, err := helpers.ValidateEvent(event, mockSignature)

		assert.Nil(t, err)
		assert.False(t, expectedFalse)
		assert.Equal(t, msg, "invalid: signature doesn't match")
	})
}
