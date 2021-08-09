package client

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Attemps to connect to mongoDB database instance.
// If process takes longer than 5 seconds, cancel is called
func NewClient() (*mongo.Client, error) {
	uri := os.Getenv("DB_URI")

	opts := options.Client().ApplyURI(uri)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		return nil, err
	}

	if err = client.Ping(ctx, nil); err != nil {
		return nil, err
	}

	log.Println("New client connection to mongoDB!")

	return client, nil

}
