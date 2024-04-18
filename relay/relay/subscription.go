package relay

import (
	"log"

	"github.com/brunobdc/nostr/relay/helpers"
	"github.com/brunobdc/nostr/relay/model"
	"github.com/fasthttp/websocket"
)

type RelaySubscriptions struct {
	ws            *Websocket
	subscriptions map[string][]*model.Filters
}

func newRelaySubscriptions(ws *Websocket) *RelaySubscriptions {
	return &RelaySubscriptions{ws: ws, subscriptions: make(map[string][]*model.Filters)}
}

func (rsSub *RelaySubscriptions) AddSubscription(subscriptionID string, filters []*model.Filters) {
	rsSub.subscriptions[subscriptionID] = filters
}

func (rsSub *RelaySubscriptions) CloseSubscription(subscriptionID string) {
	delete(rsSub.subscriptions, subscriptionID)
}

func (rsSub *RelaySubscriptions) SendResponse(response []byte) {
	rsSub.ws.WriteMessage(websocket.TextMessage, response)
}

func AddWebsocket(ws *Websocket) *RelaySubscriptions {
	rsSub := newRelaySubscriptions(ws)
	subscriptions = append(subscriptions, rsSub)
	return rsSub
}

func SubscriptionListener() {
	for event := range eventChannel {
		for _, rsSub := range subscriptions {
			for subId, filters := range rsSub.subscriptions {
				for _, filter := range filters {
					if filter.Match(*event) {
						response, err := helpers.MakeEventResponse(subId, *event)
						if err != nil {
							log.Println(err)
						} else {
							rsSub.SendResponse(response)
						}
						break
					}
				}
			}
		}
	}
}
