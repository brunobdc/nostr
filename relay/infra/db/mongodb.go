package db

import (
	"context"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var mongoClient *mongo.Client = nil

func MongoRelayDB() *mongo.Database {
	if mongoClient == nil {
		client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(os.Getenv("MONGODB_URI")))
		if err != nil {
			log.Fatal(err)
		}
		mongoClient = client
	}

	return mongoClient.Database("relay")
}
