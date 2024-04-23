package infra

import (
	"context"

	"github.com/brunobdc/nostr/relay/infra/repository"
	"github.com/brunobdc/nostr/relay/model"
)

type EventsRepository interface {
	Save(ctx context.Context, event model.Event) error
	SaveLatest(ctx context.Context, event model.Event) error
	SaveParemeterizedLatest(ctx context.Context, event model.Event) error
	FindWithFilters(ctx context.Context, filters []*model.Filters, foreachCb func(event *model.Event) error) error
}

func MakeEvenstRepository() EventsRepository {
	return repository.MakeMongoEventsRepository()
}
