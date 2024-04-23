package service

import (
	"encoding/json"

	"github.com/brunobdc/nostr/relay/infra"
)

type MessageService interface {
	Handle(ws *infra.Websocket, msg []json.RawMessage)
}

func MessageServices() map[string]MessageService {
	return map[string]MessageService{
		"EVENT": MakeEventMessageService(infra.MakeEvenstRepository()),
		"REQ":   MakeReqMessageService(infra.MakeEvenstRepository()),
		"CLOSE": MakeCloseMessageService(),
	}
}
