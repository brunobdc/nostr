package model

import (
	"bytes"
	"encoding/json"
)

type Tags map[string][]string

func (tags *Tags) UnmarshalJSON(data []byte) error {
	var tagsMatrix [][]string
	if err := json.Unmarshal(data, &tagsMatrix); err != nil {
		return err
	}
	for _, tagsSlice := range tagsMatrix {
		if len(tagsSlice) >= 2 {
			(*tags)[tagsSlice[0]] = tagsSlice[1:]
		}
	}
	return nil
}

func (tags *Tags) MarshalJSON() ([]byte, error) {
	var buffer bytes.Buffer
	buffer.WriteString("[")
	isFirst := true
	for tag, values := range *tags {
		if isFirst {
			buffer.WriteString("[\"")
			isFirst = false
		} else {
			buffer.WriteString(",[\"")
		}
		buffer.WriteString(tag)
		buffer.WriteString("\"")
		for _, value := range values {
			buffer.WriteString(",\"")
			buffer.WriteString(value)
			buffer.WriteString("\"")
		}
		buffer.WriteString("]")
	}
	buffer.WriteString("]")

	return buffer.Bytes(), nil
}
