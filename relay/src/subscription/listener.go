package subscription

import (
	"log"

	"github.com/brunobdc/nostr/relay/src/helpers"
)

func Listen() {
	for event := range eventChannel {
		for _, wsSubs := range websocketSubscriptions {
			for subId, filters := range wsSubs.subscriptions {
				for _, filter := range filters {
					if filter.Match(*event) {
						response, err := helpers.MakeEventResponse(subId, *event)
						if err != nil {
							log.Println(err)
						} else {
							wsSubs.SendResponse(response)
						}
						break
					}
				}
			}
		}
	}
}
