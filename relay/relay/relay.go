package relay

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"

	"github.com/fasthttp/websocket"
)

type Relay struct {
	serveMux *http.ServeMux
	handler  MessageHandler
}

func (relay *Relay) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Upgrade") == "websocket" {
		ws, err := UpgradeToWebsocket(w, r)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("500 - Couldn't upgrade to websocket connection!\n" + err.Error()))
			return
		}

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
					Subscription: newRelaySubscriptions(ws),
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
	} else {
		relay.serveMux.ServeHTTP(w, r)
	}
}

func Start(handler MessageHandler) {
	log.Fatal(
		http.ListenAndServe(
			":"+os.Getenv("SERVER_PORT"),
			&Relay{
				serveMux: http.DefaultServeMux,
				handler:  handler,
			},
		),
	)
}
