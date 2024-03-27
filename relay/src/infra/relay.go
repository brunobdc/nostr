package infra

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"nhooyr.io/websocket"
)

type Relay struct {
	serveMux http.ServeMux
}

func (relay *Relay) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Upgrade") == "websocket" {
		websocket_handler(w, r)
	} else {
		relay.serveMux.ServeHTTP(w, r)
	}
}

func websocket_handler(w http.ResponseWriter, r *http.Request) {
	conn, err := websocket.Accept(w, r, nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 - Couldn't make websocket connection\n" + err.Error()))
		return
	}
	defer conn.CloseNow()

	ctx := context.WithValue(r.Context(), "id", uuid.New())

	for {
		_, data, err := conn.Read(ctx)
		if err != nil {
			return
		}

		var msg_array []any
		if err = json.Unmarshal(data, &msg_array); err != nil {
			return
		}

		switch msg_array[0] {
		case "EVENT":
			// TODO
		case "REQ":
			// TODO
		case "CLOSE":
			// TODO
		}
	}
}
