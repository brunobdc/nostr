package model

import "github.com/google/uuid"

type Subscription struct {
	ID      uuid.UUID `bson:"_id"`
	SubID   string
	Filters []map[string]string
}
