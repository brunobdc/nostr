package main

import (
	"github.com/brunobdc/nostr/relay/relay"
	"github.com/brunobdc/nostr/relay/relay/repository"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	repo := repository.NewMongoEventsRepository()
	handler := relay.NewHandler(repo)
	go relay.SubscriptionListener()
	relay.Start(handler)
}
