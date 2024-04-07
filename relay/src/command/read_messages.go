package command

import (
	"encoding/json"

	"github.com/brunobdc/nostr/relay/src/infra"
	"github.com/brunobdc/nostr/relay/src/repository"
	"github.com/fasthttp/websocket"
)

func ReadMessages(ws *infra.Websocket) {
	for {
		typ, msg, err := ws.ReadMessage()
		if err != nil {
			return
		}
		if typ == websocket.PingMessage {
			ws.WriteMessage(websocket.PongMessage, nil)
			continue
		}

		var msg_arr []json.RawMessage
		json.Unmarshal(msg, &msg_arr)

		var msg_typ string
		json.Unmarshal(msg_arr[0], &msg_typ)

		switch msg_typ {
		case "EVENT":
			go handleEvent(ws, msg_arr[1:], repository.NewMongoEventsRepository())
		case "REQ":
			go handleReq(ws, msg_arr[1:], repository.NewMongoEventsRepository())
		case "CLOSE":
			go handleClose(ws, msg_arr[1:])
		}
	}
}
