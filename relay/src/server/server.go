package server

import (
	"log"
	"net/http"
	"os"

	"github.com/brunobdc/nostr/relay/src/handler"
)

type Relay struct {
	serveMux *http.ServeMux
}

func (relay *Relay) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Upgrade") == "websocket" {
		handler.WebsocketHandler(w, r)
	} else {
		relay.serveMux.ServeHTTP(w, r)
	}
}

func StartNewRelay() {
	log.Fatal(http.ListenAndServe(":"+os.Getenv("SERVER_PORT"), &Relay{serveMux: http.DefaultServeMux}))
}
