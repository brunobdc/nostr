package subscription

import (
	"github.com/brunobdc/nostr/relay/src/infra"
	"github.com/brunobdc/nostr/relay/src/model"
	"github.com/fasthttp/websocket"
)

type WebsocketSubscriptions struct {
	ws            *infra.Websocket
	subscriptions map[string][]*model.Filters
}

func newWebsocketSubscriptions(ws *infra.Websocket) *WebsocketSubscriptions {
	return &WebsocketSubscriptions{ws: ws, subscriptions: make(map[string][]*model.Filters)}
}

func (wsSub *WebsocketSubscriptions) AddSubscription(subscriptionID string, filters []*model.Filters) {
	wsSub.subscriptions[subscriptionID] = filters
}

func (wsSub *WebsocketSubscriptions) CloseSubscription(subscriptionID string) {
	delete(wsSub.subscriptions, subscriptionID)
}

func (wsSub *WebsocketSubscriptions) SendResponse(response []byte) {
	wsSub.ws.WriteMessage(websocket.TextMessage, response)
}

var websocketSubscriptions = make([]*WebsocketSubscriptions, 0)

func AddWebsocket(ws *infra.Websocket) *WebsocketSubscriptions {
	wsSubs := newWebsocketSubscriptions(ws)
	websocketSubscriptions = append(websocketSubscriptions, wsSubs)
	return wsSubs
}
