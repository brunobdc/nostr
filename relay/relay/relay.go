package relay

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/brunobdc/nostr/relay/model"
	"github.com/fasthttp/websocket"
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

		pingTicker := time.NewTicker(30 * time.Second)
		stopPingTicker := make(chan struct{})

		subscriptions := make(Subscriptions)
		relay.subscriptions[ws] = subscriptions

		ws.conn.SetReadDeadline(time.Now().Add(time.Minute))
		ws.conn.SetPongHandler(func(string) error {
			ws.conn.SetReadDeadline(time.Now().Add(time.Minute))
			return nil
		})

		go func() {
			defer func() {
				pingTicker.Stop()
				stopPingTicker <- struct{}{}
				close(stopPingTicker)
				delete(relay.subscriptions, ws)
				ws.conn.Close()
			}()
			for {
				typ, msg, err := ws.conn.ReadMessage()
				if err != nil {
					log.Println(err)
					return
				}

				if typ == websocket.TextMessage {
					var msg_arr []json.RawMessage
					err := json.Unmarshal(msg, &msg_arr)
					if err != nil {
						log.Println(err)
						return
					}

					if len(msg_arr) < 2 {
						log.Println(errors.New("received less than 2 data in message"))
						return
					}

					var msg_typ string
					err = json.Unmarshal(msg_arr[0], &msg_typ)
					if err != nil {
						log.Println(err)
						return
					}

					var relayContext = RelayContext{
						ws:           ws,
						MsgArray:     msg_arr,
						Ctx:          context.Background(),
						Subscription: subscriptions,
						eventChannel: relay.eventChannel,
					}

					switch msg_typ {
					case "EVENT":
						go relay.handler.HandleEvent(relayContext)
					case "REQ":
						go relay.handler.HandleReq(relayContext)
					case "CLOSE":
						go relay.handler.HandleClose(relayContext)
					}
				}
			}
		}()

		go func() {
			defer func() {
				pingTicker.Stop()
				ws.conn.Close()
				for range stopPingTicker {
				}
			}()
			for {
				select {
				case <-pingTicker.C:
					err := ws.WriteMessage(websocket.PingMessage, nil)
					if err != nil {
						log.Println(err)
						return
					}
					log.Println("pinging")
				case <-stopPingTicker:
					return
				}
			}
		}()
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
