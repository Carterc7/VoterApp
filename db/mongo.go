package db

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client

func Connect(uri string) *mongo.Client {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal("Mongo Connect Error:", err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal("Mongo Ping Error:", err)
	}

	Client = client
	log.Println("Connected to MongoDB!")
	return client
}

func GetPollsCollection() *mongo.Collection {
	return Client.Database("voterapp").Collection("polls")
}

func GetUsersCollection() *mongo.Collection {
	return Client.Database("voterapp").Collection("users")
}
