package websocketcontext

import (
	"context"
	"encoding/json"

	"github.com/brunobdc/nostr/relay/src/infra"
	"github.com/fasthttp/websocket"
)

type WebsocketContext struct {
	ws       *infra.Websocket
	MsgArray []json.RawMessage
	Ctx      context.Context
}

func New(ws *infra.Websocket, ctx context.Context) *WebsocketContext {
	return &WebsocketContext{ws: ws}
}

func (wsCtx *WebsocketContext) SendMessage(msg []byte) {
	wsCtx.ws.WriteMessage(websocket.TextMessage, msg)
}
