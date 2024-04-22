package relay

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/brunobdc/nostr/relay/model"
	"github.com/fasthttp/websocket"
)

type Websocket struct {
	conn          *websocket.Conn
	mutex         sync.Mutex
	ticker        *time.Ticker
	stopTicker    chan struct{}
	subscriptions Subscriptions
}

func UpgradeToWebsocket(w http.ResponseWriter, r *http.Request) (*Websocket, error) {
	upgrader := websocket.Upgrader{
		WriteBufferSize: 1024,
		ReadBufferSize:  1024,
		CheckOrigin:     func(r *http.Request) bool { return true },
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return nil, err
	}

	conn.SetReadDeadline(time.Now().Add(time.Minute))
	conn.SetPongHandler(func(string) error {
		return conn.SetReadDeadline(time.Now().Add(time.Minute))
	})

	return &Websocket{conn: conn, subscriptions: make(Subscriptions), ticker: time.NewTicker(30 * time.Second)}, nil
}

func (ws *Websocket) WriteMessage(typ int, msg []byte) error {
	ws.mutex.Lock()
	defer ws.mutex.Unlock()
	return ws.conn.WriteMessage(typ, msg)
}

func (ws *Websocket) Handle(relaySubs RelaySubscriptions, handler MessageHandler, eventChannel chan model.Event) {
	go ws.readMessagesLoop(relaySubs, handler, eventChannel)
	go ws.handlePingPong()
}

func (ws *Websocket) readMessagesLoop(relaySubs RelaySubscriptions, handler MessageHandler, eventChannel chan model.Event) {
	defer func() {
		ws.ticker.Stop()
		ws.stopTicker <- struct{}{}
		close(ws.stopTicker)
		delete(relaySubs, ws)
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
				eventChannel: eventChannel,
			}

			switch msg_typ {
			case "EVENT":
				go handler.HandleEvent(relayContext)
			case "REQ":
				go handler.HandleReq(relayContext)
			case "CLOSE":
				go handler.HandleClose(relayContext)
			}
		}
	}
}

func (ws *Websocket) handlePingPong() {
	defer func() {
		ws.ticker.Stop()
		ws.conn.Close()
		for range ws.stopTicker {
		}
	}()
	for {
		select {
		case <-ws.ticker.C:
			err := ws.WriteMessage(websocket.PingMessage, nil)
			if err != nil {
				log.Println(err)
				return
			}
			log.Println("pinging")
		case <-ws.stopTicker:
			return
		}
	}
}
