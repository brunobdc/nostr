package service

import (
	"context"
	"encoding/json"
	"log"

	"github.com/brunobdc/nostr/relay/helpers"
	"github.com/brunobdc/nostr/relay/infra"
	"github.com/brunobdc/nostr/relay/model"
	"github.com/brunobdc/nostr/relay/security"
)

type EventMessageService struct {
	Repository   infra.EventsRepository
	EventChannel chan model.Event
}

func MakeEventMessageService(repository infra.EventsRepository) *EventMessageService {
	return &EventMessageService{
		Repository:   repository,
		EventChannel: infra.Server().EventsChannel,
	}
}

func (service EventMessageService) Handle(ws *infra.Websocket, msgArray []json.RawMessage) {
	event := model.NewEvent()
	if err := json.Unmarshal(msgArray[1], event); err != nil {
		log.Println(err)
		return
	}
	valid, msg, err := helpers.ValidateEvent(*event, security.SchnorrSignature{})
	if err != nil {
		log.Println(err)
		response, err := helpers.MakeOkResponse(event.ID, false, "error: couldn't validate de event")
		if err != nil {
			log.Println(err)
		} else {
			ws.WriteTextMessage(response)
		}
		return
	}
	if !valid {
		response, err := helpers.MakeOkResponse(event.ID, false, msg)
		if err != nil {
			log.Println(err)
		} else {
			ws.WriteTextMessage(response)
		}
		return
	}

	if event.Kind == 1 || (event.Kind >= 1_000 && event.Kind < 10_000) {
		err := service.Repository.Save(context.Background(), *event)
		if err != nil {
			log.Println(err)
			response, err := helpers.MakeOkResponse(event.ID, false, "error: couldn't save the event in database")
			if err != nil {
				log.Println(err)
			} else {
				ws.WriteTextMessage(response)
			}
			return
		}
	} else if event.Kind == 0 || event.Kind == 3 || (event.Kind >= 10_000 && event.Kind < 20_000) {
		err := service.Repository.SaveLatest(context.Background(), *event)
		if err != nil {
			log.Println(err)
			response, err := helpers.MakeOkResponse(event.ID, false, "error: couldn't save the event in database")
			if err != nil {
				log.Println(err)
			} else {
				ws.WriteTextMessage(response)
			}
			return
		}
	} else if event.Kind >= 30_000 && event.Kind < 40_000 {
		err := service.Repository.SaveParemeterizedLatest(context.Background(), *event)
		if err != nil {
			log.Println(err)
			response, err := helpers.MakeOkResponse(event.ID, false, "error: couldn't save the event in database")
			if err != nil {
				log.Println(err)
			} else {
				ws.WriteTextMessage(response)
			}
			return
		}
	}

	service.EventChannel <- *event

	response, err := helpers.MakeOkResponse(event.ID, true, "")
	if err != nil {
		log.Println(err)
	} else {
		ws.WriteTextMessage(response)
	}
}
