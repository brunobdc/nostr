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
	Subscription Subscriptions
	eventChannel chan model.Event
}

func (ctx *RelayContext) SendMessage(msg []byte) {
	ctx.ws.WriteMessage(websocket.TextMessage, msg)
}
