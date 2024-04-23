package main

import (
	"log"
	"net/http"
	"os"

	"github.com/brunobdc/nostr/relay/application/controller"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Llongfile)

	http.Handle("GET /", controller.MakeWebscoketController())

	log.Fatal(
		http.ListenAndServe(
			":"+os.Getenv("SERVER_PORT"),
			nil,
		),
	)
}
