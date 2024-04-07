package repository

import (
	"context"

	"github.com/brunobdc/nostr/relay/src/model"
)

type EventsRepository interface {
	Save(event model.Event)
	SaveLatest(event model.Event)
	SaveParemeterizedLatest(event model.Event)
	FindWithFilters(filters []*model.Filters) EventCursor
}

type EventCursor interface {
	Close(ctx context.Context)
	Next(ctx context.Context) bool
	Decode(event *model.Event) error
}
