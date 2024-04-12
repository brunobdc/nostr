package infra

import (
	"net/http"
	"sync"

	"github.com/fasthttp/websocket"
)

type Websocket struct {
	conn  *websocket.Conn
	mutex sync.Mutex
}

func UpgradeToWebsocket(w http.ResponseWriter, r *http.Request) (*Websocket, error) {
	upgrader := websocket.Upgrader{}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return nil, err
	}
	return &Websocket{conn: conn}, nil
}

func (ws *Websocket) WriteMessage(typ int, msg []byte) error {
	ws.mutex.Lock()
	defer ws.mutex.Unlock()
	return ws.conn.WriteMessage(typ, msg)
}

func (ws *Websocket) ReadMessages(messageHandler func([]byte)) error {
	for {
		typ, msg, err := ws.conn.ReadMessage()
		if err != nil {
			return err
		}
		if typ == websocket.PingMessage {
			ws.WriteMessage(websocket.PongMessage, nil)
			continue
		}

		go messageHandler(msg)
	}
}
