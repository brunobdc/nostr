package controller

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/brunobdc/nostr/relay/application/service"
	"github.com/brunobdc/nostr/relay/infra"
	"github.com/brunobdc/nostr/relay/model"
)

type WebsocketController struct {
	ServerWebsockets *infra.ServerWebsockets
	EventChannerl    chan model.Event
	MessageServices  map[string]service.MessageService
}

func MakeWebscoketController() *WebsocketController {
	server := infra.Server()
	return &WebsocketController{
		ServerWebsockets: server.Websockets,
		EventChannerl:    server.EventsChannel,
		MessageServices:  service.MessageServices(),
	}
}

func (wsController WebsocketController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Upgrade") == "websocket" {
		ws, err := infra.UpgradeToWebsocket(w, r)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("500 - Couldn't upgrade to websocket connection!\n" + err.Error()))
			return
		}

		wsController.ServerWebsockets.NewWebsocket(ws)

		go wsController.readMessages(ws)
	}
}

func (wsController WebsocketController) readMessages(ws *infra.Websocket) {
	defer func() {
		ws.Close()
		wsController.ServerWebsockets.RemoveWebsocket(ws)
	}()
	for {
		msg, err := ws.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}

		var msgArr []json.RawMessage
		err = json.Unmarshal(msg, &msgArr)
		if err != nil {
			log.Println(err)
			return
		}

		if len(msgArr) < 2 {
			log.Println(errors.New("message array length is less than 2"))
			continue
		}

		var typ string
		err = json.Unmarshal(msgArr[0], &typ)
		if err != nil {
			log.Println(err)
			continue
		}

		if service, ok := wsController.MessageServices[typ]; ok {
			service.Handle(ws, msgArr)
		} else {
			log.Println(errors.New("invalid message type received: " + typ))
		}
	}
}
