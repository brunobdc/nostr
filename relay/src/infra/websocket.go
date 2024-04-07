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

func (ws *Websocket) WriteJson(data any) error {
	ws.mutex.Lock()
	defer ws.mutex.Unlock()
	return ws.conn.WriteJSON(data)
}

func (ws *Websocket) WriteMessage(typ int, message []byte) error {
	ws.mutex.Lock()
	defer ws.mutex.Unlock()
	return ws.conn.WriteMessage(typ, message)
}

func (ws *Websocket) ReadMessage() (int, []byte, error) {
	return ws.conn.ReadMessage()
}
