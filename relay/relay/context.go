package relay

import (
	"context"
	"encoding/json"

	"github.com/fasthttp/websocket"
)

type RelayContext struct {
	ws           *Websocket
	MsgArray     []json.RawMessage
	Ctx          context.Context
	Subscription *RelaySubscriptions
}

func (ctx *RelayContext) SendMessage(msg []byte) {
	ctx.ws.WriteMessage(websocket.TextMessage, msg)
}
