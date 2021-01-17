package database

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// Client creates a mongo client
var Client *mongo.Client

// InitializeDB initialized DB
func InitializeDB() {
	var err error
	uri := "mongodb://localhost:27017"

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	Client, err = mongo.Connect(ctx, options.Client().ApplyURI(uri))

	if err != nil {
		panic(err)
	}

	// defer func() {
	// 	if err = Client.Disconnect(ctx); err != nil {
	// 		panic(err)
	// 	}
	// }()

	if err := Client.Ping(ctx, readpref.Primary()); err != nil {
		panic(err)
	}

	log.Println("Connected to zode-projectsDB and pinged")

}
