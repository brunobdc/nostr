package service

import (
	"encoding/json"
	"log"

	"github.com/brunobdc/nostr/relay/infra"
)

type CloseMessageService struct{}

func MakeCloseMessageService() *CloseMessageService {
	return &CloseMessageService{}
}

func (service CloseMessageService) Handle(ws *infra.Websocket, msgArray []json.RawMessage) {
	var subscriptionID string
	err := json.Unmarshal(msgArray[1], &subscriptionID)
	if err != nil {
		log.Println(err)
		return
	}
	ws.Subscriptions.CloseSubscription(subscriptionID)
}
