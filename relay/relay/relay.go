package relay

import (
	"log"
	"net/http"
	"os"

	"github.com/brunobdc/nostr/relay/model"
)

type Relay struct {
	serveMux      *http.ServeMux
	handler       MessageHandler
	subscriptions RelaySubscriptions
	eventChannel  chan model.Event
}

func (relay *Relay) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Upgrade") == "websocket" {
		ws, err := UpgradeToWebsocket(w, r)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("500 - Couldn't upgrade to websocket connection!\n" + err.Error()))
			return
		}
		relay.subscriptions[ws] = struct{}{}

		ws.Handle(relay.subscriptions, relay.handler, relay.eventChannel)
	} else {
		relay.serveMux.ServeHTTP(w, r)
	}
}

func Start(handler MessageHandler) {
	relaySubscriptions := make(RelaySubscriptions)
	eventChannel := make(chan model.Event)

	go SubscriptionListener(relaySubscriptions, eventChannel)

	log.Fatal(
		http.ListenAndServe(
			":"+os.Getenv("SERVER_PORT"),
			&Relay{
				serveMux:      http.DefaultServeMux,
				handler:       handler,
				subscriptions: relaySubscriptions,
				eventChannel:  eventChannel,
			},
		),
	)
}
