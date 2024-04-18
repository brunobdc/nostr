package relay

import "github.com/brunobdc/nostr/relay/model"

var eventChannel = make(chan *model.Event)
