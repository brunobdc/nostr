package command

import (
	"context"
	"encoding/json"

	"github.com/brunobdc/nostr/relay/src/infra"
	"github.com/brunobdc/nostr/relay/src/model"
	"github.com/brunobdc/nostr/relay/src/repository"
)

var WebsocketSubscriptions = make(map[*infra.Websocket]map[string][]*model.Filters)
var EventsChannel = make(chan *model.Event)

func handleReq(ws *infra.Websocket, data []json.RawMessage, repo repository.EventsRepository) {
	var subscription_id string
	json.Unmarshal(data[0], &subscription_id)
	filters := make([]*model.Filters, len(data)-1)
	for i, v := range data[1:] {
		json.Unmarshal(v, &filters[i])
	}

	cursor := repo.FindWithFilters(filters)
	for cursor.Next(context.TODO()) {
		var event *model.Event
		if err := cursor.Decode(event); err != nil {
			panic(err)
		}
		eventJson, err := event.MarshalJSON()
		if err != nil {
			panic(err)
		}
		ws.WriteJson([]any{"EVENT", subscription_id, string(eventJson)})
	}
	cursor.Close(context.TODO())

	if subs, ok := WebsocketSubscriptions[ws]; ok {
		subs[subscription_id] = filters
	} else {
		subs := make(map[string][]*model.Filters)
		subs[subscription_id] = filters
		WebsocketSubscriptions[ws] = subs
	}
}

func SubscriptionListener() {
	for event := range EventsChannel {
		for ws, subs := range WebsocketSubscriptions {
			for subId, filters := range subs {
				for _, filter := range filters {
					if filter.Match(*event) {
						eventJson, _ := event.MarshalJSON()
						ws.WriteJson([]any{"EVENT", subId, string(eventJson)})
						continue
					}
				}
			}
		}
	}
}
