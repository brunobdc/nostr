package helpers

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"strconv"

	"github.com/brunobdc/nostr/relay/model"
	"github.com/brunobdc/nostr/relay/security"
)

func ValidateEvent(event model.Event, signature security.Signature) (bool, string, error) {
	bytesID, err := hex.DecodeString(event.ID)
	if err != nil {
		return false, "", err
	}
	if len(bytesID) != 32 {
		return false, "invalid: id bytes size is differente from 32", nil
	}
	bytesPubKey, err := hex.DecodeString(event.PublicKey)
	if err != nil {
		return false, "", err
	}
	if len(bytesPubKey) != 32 {
		return false, "invalid: pubKey bytes size is differente from 32", nil
	}
	bytesSignature, err := hex.DecodeString(event.Signature)
	if err != nil {
		return false, "", err
	}
	if len(bytesSignature) != 64 {
		return false, "invalid: signature bytes size is differente from 64", nil
	}

	var buffer bytes.Buffer
	buffer.WriteString("[0,\"")
	buffer.WriteString(event.PublicKey)
	buffer.WriteString("\",")
	buffer.WriteString(strconv.FormatUint(event.CreatedAt, 10))
	buffer.WriteString(",")
	buffer.WriteString(strconv.FormatUint(uint64(event.Kind), 10))
	buffer.WriteString(",")
	jsonTags, err := event.Tags.MarshalJSON()
	if err != nil {
		return false, "", err
	}
	buffer.Write(jsonTags)
	buffer.WriteString(",\"")
	buffer.WriteString(event.Content)
	buffer.WriteString("\"]")

	hash := sha256.Sum256(buffer.Bytes())
	if [32]byte(bytesID) != hash {
		return false, "invalid: id hash doesn't match the data", nil
	}

	valid, err := signature.VerifySignature([64]byte(bytesSignature), [32]byte(bytesID), [32]byte(bytesPubKey))
	if err != nil {
		return false, "", err
	}
	if !valid {
		return false, "invalid: signature doesn't match", nil
	}
	return true, "", nil
}
