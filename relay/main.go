package main

import (
	"log"

	"github.com/brunobdc/nostr/relay/relay"
	"github.com/brunobdc/nostr/relay/relay/repository"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Llongfile)
	repo := repository.NewMongoEventsRepository()
	handler := relay.NewHandler(repo)
	relay.Start(handler)
}
