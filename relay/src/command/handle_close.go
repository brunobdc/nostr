package command

import (
	"encoding/json"

	"github.com/brunobdc/nostr/relay/src/infra"
)

func handleClose(ws *infra.Websocket, data []json.RawMessage) {
	var subscription_id string
	json.Unmarshal(data[0], &subscription_id)
	delete(WebsocketSubscriptions[ws], subscription_id)
}
