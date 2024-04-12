package main

import (
	"context"

	"github.com/brunobdc/nostr/relay/src/infra"
	"github.com/brunobdc/nostr/relay/src/infra/server"
	"github.com/brunobdc/nostr/relay/src/subscription"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	infra.InitializeDB(context.Background())
	go subscription.Listen()
	server.StartNewRelay()
}
