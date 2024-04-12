package server

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/brunobdc/nostr/relay/src/handler"
	"github.com/brunobdc/nostr/relay/src/infra"
	"github.com/brunobdc/nostr/relay/src/subscription"
	websocketcontext "github.com/brunobdc/nostr/relay/src/websocketContext"
)

type Relay struct {
	serveMux *http.ServeMux
}

func (relay *Relay) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Upgrade") == "websocket" {
		ws, err := infra.UpgradeToWebsocket(w, r)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("500 - Couldn't upgrade to websocket connection!\n" + err.Error()))
			return
		}
		websocketSubscriptions := subscription.AddWebsocket(ws)
		ws.ReadMessages(func(msg []byte) {
			handler.MessageHandler(*websocketcontext.New(ws, context.Background()), msg, websocketSubscriptions)
		})
	} else {
		relay.serveMux.ServeHTTP(w, r)
	}
}

func StartNewRelay() {
	log.Fatal(http.ListenAndServe(":"+os.Getenv("SERVER_PORT"), &Relay{serveMux: http.DefaultServeMux}))
}
