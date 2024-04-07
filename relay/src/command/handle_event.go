package command

import (
	"encoding/json"

	"github.com/brunobdc/nostr/relay/src/infra"
	"github.com/brunobdc/nostr/relay/src/model"
	"github.com/brunobdc/nostr/relay/src/repository"
)

func handleEvent(ws *infra.Websocket, data []json.RawMessage, repo repository.EventsRepository) {
	var event model.Event
	if err := json.Unmarshal(data[0], &event); err != nil {
		return
	}
	if valid, msg := event.IsValid(); !valid {
		ws.WriteJson([]any{"OK", event.ID, false, msg})
	}

	if event.Kind == 1 || (event.Kind >= 1_000 && event.Kind < 10_000) {
		repo.Save(event)
	} else if event.Kind == 0 || event.Kind == 3 || (event.Kind >= 10_000 && event.Kind < 20_000) {
		repo.SaveLatest(event)
	} else if event.Kind >= 30_000 && event.Kind < 40_000 {
		repo.SaveParemeterizedLatest(event)
	}

	ws.WriteJson([]any{"OK", event.ID, true, ""})

	EventsChannel <- &event
}
