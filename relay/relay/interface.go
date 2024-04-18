package relay

import (
	"context"

	"github.com/brunobdc/nostr/relay/model"
)

type EventsRepository interface {
	Save(ctx context.Context, event model.Event) error
	SaveLatest(ctx context.Context, event model.Event) error
	SaveParemeterizedLatest(ctx context.Context, event model.Event) error
	FindWithFilters(ctx context.Context, filters []*model.Filters) (func(ctx context.Context) (bool, *model.Event), error)
}

type MessageHandler interface {
	HandleEvent(ctx RelayContext)
	HandleReq(ctx RelayContext)
	HandleClose(ctx RelayContext)
}
