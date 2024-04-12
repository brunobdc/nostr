package model

import (
	"encoding/json"
	"slices"
)

type Filters struct {
	IDs     []string
	Authors []string
	Kinds   []EventKind
	Tags    Tags
	Since   uint64
	Until   uint64
	Limit   uint
}

func (filters Filters) Match(event Event) bool {
	if len(filters.IDs) > 0 && !slices.Contains(filters.IDs, event.ID) {
		return false
	}
	if len(filters.Authors) > 0 && !slices.Contains(filters.Authors, event.PublicKey) {
		return false
	}
	if len(filters.Kinds) > 0 && !slices.Contains(filters.Kinds, event.Kind) {
		return false
	}
	for tag, values := range filters.Tags {
		eventTagValues, ok := event.Tags[tag]
		if !ok || !slices.ContainsFunc(values, func(value string) bool { return slices.Contains(eventTagValues, value) }) {
			return false
		}
	}
	if filters.Since > event.CreatedAt || filters.Until < event.CreatedAt {
		return false
	}
	return true
}

func (filter *Filters) UnmarshalJSON(b []byte) error {
	var data map[string]json.RawMessage
	if err := json.Unmarshal(b, &data); err != nil {
		return err
	}

	for key, value := range data {
		switch key {
		case "ids":
			var ids []string
			if err := json.Unmarshal(value, &ids); err != nil {
				return err
			}
			filter.IDs = ids
		case "authors":
			var authors []string
			if err := json.Unmarshal(value, &authors); err != nil {
				return err
			}
			filter.Authors = authors
		case "kinds":
			var kinds []EventKind
			if err := json.Unmarshal(value, &kinds); err != nil {
				return err
			}
			filter.Kinds = kinds
		case "since":
			var since uint64
			if err := json.Unmarshal(value, &since); err != nil {
				return err
			}
			filter.Since = since
		case "until":
			var until uint64
			if err := json.Unmarshal(value, &until); err != nil {
				return err
			}
			filter.Until = until
		case "limit":
			var limit uint
			if err := json.Unmarshal(value, &limit); err != nil {
				return err
			}
			filter.Limit = limit
		default:
			if len(key) == 2 && string(key[0]) == "#" {
				var tag_values []string
				if err := json.Unmarshal(value, &tag_values); err != nil {
					return err
				}
				filter.Tags[string(key[1])] = tag_values
			}
		}
	}

	return nil
}
