package handler

import (
	"net/http"

	"github.com/brunobdc/nostr/relay/src/command"
	"github.com/brunobdc/nostr/relay/src/infra"
)

func WebsocketHandler(w http.ResponseWriter, r *http.Request) {
	ws, err := infra.UpgradeToWebsocket(w, r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 - Couldn't upgrade to websocket connection!\n" + err.Error()))
		return
	}

	go command.ReadMessages(ws)
}
