package subscription

import "github.com/brunobdc/nostr/relay/src/model"

var eventChannel = make(chan *model.Event)

func NewEvent(event *model.Event) {
	eventChannel <- event
}
