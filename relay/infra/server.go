package infra

import (
	"sync"

	"github.com/brunobdc/nostr/relay/model"
)

type server struct {
	Websockets    *ServerWebsockets
	EventsChannel chan model.Event
}

type ServerWebsockets struct {
	mutex      *sync.Mutex
	websockets map[*Websocket]struct{}
}

var singletonServer *server = nil

func Server() *server {
	if singletonServer == nil {
		singletonServer = &server{
			Websockets:    &ServerWebsockets{websockets: make(map[*Websocket]struct{})},
			EventsChannel: make(chan model.Event),
		}
	}
	return singletonServer
}

func (svWs *ServerWebsockets) NewWebsocket(ws *Websocket) {
	svWs.mutex.Lock()
	defer svWs.mutex.Unlock()
	svWs.websockets[ws] = struct{}{}
}

func (svWs *ServerWebsockets) RemoveWebsocket(ws *Websocket) {
	svWs.mutex.Lock()
	defer svWs.mutex.Unlock()
	delete(svWs.websockets, ws)
}
