package infra

import (
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
	Subscriptions *model.Subscriptions
}

var upgrader = websocket.Upgrader{
	WriteBufferSize: 1024,
	ReadBufferSize:  1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func UpgradeToWebsocket(w http.ResponseWriter, r *http.Request) (*Websocket, error) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return nil, err
	}

	ws := &Websocket{
		conn:          conn,
		Subscriptions: model.NewSubscriptions(),
		ticker:        time.NewTicker(30 * time.Second),
		stopTicker:    make(chan struct{}),
	}

	go ws.handlePingPong()

	return ws, nil
}

func (ws *Websocket) WriteTextMessage(msg []byte) error {
	ws.mutex.Lock()
	defer ws.mutex.Unlock()
	return ws.conn.WriteMessage(websocket.TextMessage, msg)
}

func (ws *Websocket) ReadMessage() ([]byte, error) {
	_, msg, err := ws.conn.ReadMessage()
	return msg, err
}

func (ws *Websocket) Close() error {
	ws.ticker.Stop()
	ws.stopTicker <- struct{}{}
	return ws.conn.Close()
}

func (ws *Websocket) handlePingPong() {
	ws.conn.SetReadDeadline(time.Now().Add(time.Minute + (30 * time.Second)))
	ws.conn.SetPongHandler(func(string) error {
		return ws.conn.SetReadDeadline(time.Now().Add(time.Minute))
	})

	defer func() {
		ws.ticker.Stop()
		ws.conn.Close()
		for range ws.stopTicker {
		}
	}()
	for {
		select {
		case <-ws.ticker.C:
			err := ws.conn.WriteMessage(websocket.PingMessage, nil)
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
