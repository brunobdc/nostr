package relay

import (
	"encoding/json"
	"log"

	"github.com/brunobdc/nostr/relay/helpers"
	"github.com/brunobdc/nostr/relay/model"
	"github.com/brunobdc/nostr/relay/security"
)

type Handler struct {
	repository EventsRepository
}

func NewHandler(repository EventsRepository) *Handler {
	return &Handler{repository: repository}
}

func (handler Handler) HandleClose(ctx RelayContext) {
	var subscriptionID string
	err := json.Unmarshal(ctx.MsgArray[1], &subscriptionID)
	if err != nil {
		log.Println(err)
		return
	}
	ctx.Subscription.CloseSubscription(subscriptionID)
}

func (handler Handler) HandleReq(ctx RelayContext) {
	var subscriptionID string
	err := json.Unmarshal(ctx.MsgArray[1], &subscriptionID)
	if err != nil {
		log.Println(err)
		return
	}
	filters := make([]*model.Filters, len(ctx.MsgArray)-2)
	for i, v := range ctx.MsgArray[2:] {
		err := json.Unmarshal(v, &filters[i])
		if err != nil {
			log.Println(err)
			return
		}
	}

	err = handler.repository.FindWithFilters(
		ctx.Ctx,
		filters,
		func(event *model.Event) error {
			response, err := helpers.MakeEventResponse(subscriptionID, *event)
			if err != nil {
				return err
			}
			ctx.SendMessage(response)
			return nil
		},
	)
	if err != nil {
		log.Println(err)
		return
	}

	response, err := helpers.MakeEoseResponse(subscriptionID)
	if err != nil {
		log.Println(err)
	}
	ctx.SendMessage(response)

	ctx.Subscription.AddSubscription(subscriptionID, filters)
}

func (handler Handler) HandleEvent(ctx RelayContext) {
	event := model.NewEvent()
	if err := json.Unmarshal(ctx.MsgArray[1], event); err != nil {
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
		err := handler.repository.Save(ctx.Ctx, *event)
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
		err := handler.repository.SaveLatest(ctx.Ctx, *event)
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
		err := handler.repository.SaveParemeterizedLatest(ctx.Ctx, *event)
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

	ctx.eventChannel <- *event

	response, err := helpers.MakeOkResponse(event.ID, false, "")
	if err != nil {
		log.Println(err)
	} else {
		ctx.SendMessage(response)
	}
}
