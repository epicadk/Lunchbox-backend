package database

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//ConnectDB lol
func MongoClient(ctx context.Context) *mongo.Client {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(
		//Add mongodb URI here
		"",
	))
	if err != nil {
		log.Fatal(err)
	}
	return client
}
