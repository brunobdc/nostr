package relay

import (
	"log"

	"github.com/brunobdc/nostr/relay/helpers"
	"github.com/brunobdc/nostr/relay/model"
	"github.com/fasthttp/websocket"
)

type RelaySubscriptions map[*Websocket]struct{}

type Subscriptions map[string][]*model.Filters

func (subs Subscriptions) AddSubscription(subscriptionID string, filters []*model.Filters) {
	subs[subscriptionID] = filters
}

func (subs Subscriptions) CloseSubscription(subscriptionID string) {
	delete(subs, subscriptionID)
}

func SubscriptionListener(relaySubs RelaySubscriptions, eventChannel chan model.Event) {
	for event := range eventChannel {
		for ws := range relaySubs {
			for subId, filters := range ws.subscriptions {
				for _, filter := range filters {
					if filter.Match(event) {
						response, err := helpers.MakeEventResponse(subId, event)
						if err != nil {
							log.Println(err)
						} else {
							ws.WriteMessage(websocket.TextMessage, response)
						}
						break
					}
				}
			}
		}
	}
}
