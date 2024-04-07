package main

import (
	"context"

	"github.com/brunobdc/nostr/relay/src/command"
	"github.com/brunobdc/nostr/relay/src/infra"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	infra.InitializeDB(context.Background())
	go command.SubscriptionListener()
}
