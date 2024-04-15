package handler

import (
	"encoding/json"
	"log"

	"github.com/brunobdc/nostr/relay/src/helpers"
	"github.com/brunobdc/nostr/relay/src/model"
	"github.com/brunobdc/nostr/relay/src/repository"
	"github.com/brunobdc/nostr/relay/src/security"
	"github.com/brunobdc/nostr/relay/src/subscription"
	websocketcontext "github.com/brunobdc/nostr/relay/src/websocketContext"
)

func MessageHandler(ctx websocketcontext.WebsocketContext, msg []byte, websocketSubscriptions *subscription.WebsocketSubscriptions) {
	var msg_arr []json.RawMessage
	json.Unmarshal(msg, &msg_arr)

	var msg_typ string
	json.Unmarshal(msg_arr[0], &msg_typ)

	ctx.MsgArray = msg_arr[1:]
	switch msg_typ {
	case "EVENT":
		handleEvent(ctx, repository.NewMongoEventsRepository())
	case "REQ":
		handleReq(ctx, repository.NewMongoEventsRepository(), websocketSubscriptions)
	case "CLOSE":
		handleClose(ctx, websocketSubscriptions)
	}
}

func handleClose(ctx websocketcontext.WebsocketContext, websocketSubscriptions *subscription.WebsocketSubscriptions) {
	var subscriptionID string
	err := json.Unmarshal(ctx.MsgArray[0], &subscriptionID)
	if err != nil {
		log.Println(err)
		return
	}
	websocketSubscriptions.CloseSubscription(subscriptionID)
}

func handleReq(ctx websocketcontext.WebsocketContext, eventRepository repository.EventsRepository, websocketSubscriptions *subscription.WebsocketSubscriptions) {
	var subscriptionID string
	err := json.Unmarshal(ctx.MsgArray[0], &subscriptionID)
	if err != nil {
		log.Println(err)
		return
	}
	filters := make([]*model.Filters, len(ctx.MsgArray)-1)
	for i, v := range ctx.MsgArray[1:] {
		err := json.Unmarshal(v, &filters[i])
		if err != nil {
			log.Println(err)
			return
		}
	}

	iterate, err := eventRepository.FindWithFilters(ctx.Ctx, filters)
	if err != nil {
		log.Println(err)
		return
	}
	for ok, event := iterate(ctx.Ctx); ok; {
		response, err := helpers.MakeEventResponse(subscriptionID, *event)
		if err != nil {
			log.Println(err)
			continue
		}
		ctx.SendMessage(response)
	}

	response, err := helpers.MakeEoseResponse(subscriptionID)
	if err != nil {
		log.Println(err)
	}
	ctx.SendMessage(response)

	websocketSubscriptions.AddSubscription(subscriptionID, filters)
}

func handleEvent(ctx websocketcontext.WebsocketContext, eventRepository repository.EventsRepository) {
	var event model.Event
	if err := json.Unmarshal(ctx.MsgArray[0], &event); err != nil {
		log.Println(err)
		return
	}
	valid, msg, err := helpers.ValidateEvent(event, security.SchnorrSignature{})
	if err != nil {
		log.Println(err)
		response, err := helpers.MakeOkResponse(event.ID, false, "error: couldn't validate de event")
		if err != nil {
			log.Println(err)
		} else {
			ctx.SendMessage(response)
		}
		return
	}
	if !valid {
		response, err := helpers.MakeOkResponse(event.ID, false, msg)
		if err != nil {
			log.Println(err)
		} else {
			ctx.SendMessage(response)
		}
		return
	}

	if event.Kind == 1 || (event.Kind >= 1_000 && event.Kind < 10_000) {
		err := eventRepository.Save(ctx.Ctx, event)
		if err != nil {
			log.Println(err)
			response, err := helpers.MakeOkResponse(event.ID, false, "error: couldn't save the event in database")
			if err != nil {
				log.Println(err)
			} else {
				ctx.SendMessage(response)
			}
			return
		}
	} else if event.Kind == 0 || event.Kind == 3 || (event.Kind >= 10_000 && event.Kind < 20_000) {
		err := eventRepository.SaveLatest(ctx.Ctx, event)
		if err != nil {
			log.Println(err)
			response, err := helpers.MakeOkResponse(event.ID, false, "error: couldn't save the event in database")
			if err != nil {
				log.Println(err)
			} else {
				ctx.SendMessage(response)
			}
			return
		}
	} else if event.Kind >= 30_000 && event.Kind < 40_000 {
		err := eventRepository.SaveParemeterizedLatest(ctx.Ctx, event)
		if err != nil {
			log.Println(err)
			response, err := helpers.MakeOkResponse(event.ID, false, "error: couldn't save the event in database")
			if err != nil {
				log.Println(err)
			} else {
				ctx.SendMessage(response)
			}
			return
		}
	}

	subscription.NewEvent(&event)

	response, err := helpers.MakeOkResponse(event.ID, false, "")
	if err != nil {
		log.Println(err)
	} else {
		ctx.SendMessage(response)
	}
}
