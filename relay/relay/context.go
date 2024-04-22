package relay

import (
	"context"
	"encoding/json"

	"github.com/brunobdc/nostr/relay/model"
	"github.com/fasthttp/websocket"
)

type RelayContext struct {
	ws           *Websocket
	MsgArray     []json.RawMessage
	Ctx          context.Context
	eventChannel chan model.Event
}

func (ctx *RelayContext) SendMessage(msg []byte) {
	ctx.ws.WriteMessage(websocket.TextMessage, msg)
}

func (ctx *RelayContext) AddSubscription(subscriptionId string, filters []*model.Filters) {
	ctx.ws.subscriptions.AddSubscription(subscriptionId, filters)
}

func (ctx *RelayContext) CloseSubscription(subscriptionID string) {
	delete(ctx.ws.subscriptions, subscriptionID)
}
