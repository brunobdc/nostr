package model_test

import (
	"encoding/json"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/brunobdc/nostr/relay/src/model"
	"github.com/stretchr/testify/assert"
)

func TestTagsUnmarshalJSON(T *testing.T) {
	T.Parallel()

	T.Run("Should return error if json.Unmarshal as [][]string fails", func(t *testing.T) {
		t.Parallel()

		var randomBytes []byte
		gofakeit.Slice(&randomBytes)

		var tags model.Tags
		err := tags.UnmarshalJSON(randomBytes)

		assert.Error(t, err)
	})

	T.Run("Should set the tags map correctly if input is a valid json string of [][]string", func(t *testing.T) {
		t.Parallel()

		var stringMatrix [][]string
		gofakeit.Slice(&stringMatrix)

		jsonString, err := json.Marshal(stringMatrix)
		assert.Nil(t, err)

		tags := make(model.Tags)
		err = tags.UnmarshalJSON(jsonString)

		assert.Nil(t, err)
		for _, strArray := range stringMatrix {
			if len(strArray) >= 2 {
				tagFound := false
				for tag, values := range tags {
					if tag == strArray[0] {
						tagFound = true
						assert.Equal(t, values, strArray[1:])
						break
					}
				}
				assert.True(t, tagFound)
			}
		}
	})
}

func TestTagsMarshalJSON(T *testing.T) {
	T.Parallel()

	T.Run("Should return a valid json string", func(t *testing.T) {
		t.Parallel()

		var stringMatrix [][]string
		gofakeit.Slice(&stringMatrix)

		tags := make(model.Tags)
		for _, values := range stringMatrix {
			tags[values[0]] = values[1:]
		}

		expectedBytes, err := json.Marshal(stringMatrix)
		assert.Nil(t, err)

		gotBytes, err := tags.MarshalJSON()
		assert.Nil(t, err)

		expected := make(model.Tags)
		got := make(model.Tags)
		err = expected.UnmarshalJSON(expectedBytes)
		assert.Nil(t, err)
		err = got.UnmarshalJSON(gotBytes)
		assert.Nil(t, err)

		assert.Equal(t, got, expected)
	})
}
