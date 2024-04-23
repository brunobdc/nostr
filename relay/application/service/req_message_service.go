package service

import (
	"context"
	"encoding/json"
	"log"

	"github.com/brunobdc/nostr/relay/helpers"
	"github.com/brunobdc/nostr/relay/infra"
	"github.com/brunobdc/nostr/relay/model"
)

type ReqMessageService struct {
	Repository infra.EventsRepository
}

func MakeReqMessageService(repository infra.EventsRepository) *ReqMessageService {
	return &ReqMessageService{
		Repository: repository,
	}
}

func (service ReqMessageService) Handle(ws *infra.Websocket, msgArray []json.RawMessage) {
	var subscriptionID string
	err := json.Unmarshal(msgArray[1], &subscriptionID)
	if err != nil {
		log.Println(err)
		return
	}
	filters := make([]*model.Filters, len(msgArray)-2)
	for i, v := range msgArray[2:] {
		err := json.Unmarshal(v, &filters[i])
		if err != nil {
			log.Println(err)
			return
		}
	}

	err = service.Repository.FindWithFilters(
		context.Background(),
		filters,
		func(event *model.Event) error {
			response, err := helpers.MakeEventResponse(subscriptionID, *event)
			if err != nil {
				return err
			}
			ws.WriteTextMessage(response)
			return nil
		},
	)
	if err != nil {
		log.Println(err)
		return
	}

	response, err := helpers.MakeEoseResponse(subscriptionID)
	if err != nil {
		log.Println(err)
		return
	}
	ws.WriteTextMessage(response)

	ws.Subscriptions.AddSubscription(subscriptionID, filters)
}
