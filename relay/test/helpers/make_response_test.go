package helpers_test

import (
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/brunobdc/nostr/relay/src/helpers"
	"github.com/stretchr/testify/assert"
)

func TestMakeEoseResponse(T *testing.T) {
	T.Parallel()
	T.Run("Should return a valid json string of EOSE response", func(t *testing.T) {
		t.Parallel()
		subscriptionID := gofakeit.UUID()

		result, err := helpers.MakeEoseResponse(subscriptionID)
		assert.Nil(t, err)

		expected := []byte("[\"EOSE\",\"" + subscriptionID + "\"]")
		assert.Equal(t, result, expected)
	})
}
