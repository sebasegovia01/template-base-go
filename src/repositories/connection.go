package repositories

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// mongodb connection
func Connect(uri string, db_name string) *mongo.Database {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal(err)
	}

	// Ping conn
	if err := client.Ping(ctx, nil); err != nil {
		log.Fatal(err)
	}

	log.Println("Succesfully connection to db")

	return client.Database(db_name)
}
