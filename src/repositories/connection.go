package repositories

import (
	"context"
	"log"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	clientInstance     *mongo.Client
	clientInstanceOnce sync.Once
	clientInstanceErr  error
)

func Connect(uri string, dbName string) (*mongo.Database, error) {
	clientInstanceOnce.Do(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		clientInstance, clientInstanceErr = mongo.Connect(ctx, options.Client().ApplyURI(uri))
		if clientInstanceErr != nil {
			log.Printf("Error connecting to MongoDB: %v\n", clientInstanceErr)
			return
		}

		pingCtx, pingCancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer pingCancel()
		if clientInstanceErr = clientInstance.Ping(pingCtx, nil); clientInstanceErr != nil {
			log.Printf("Error pinging MongoDB: %v\n", clientInstanceErr)
			_ = clientInstance.Disconnect(context.Background())
		} else {
			log.Println("Successfully connected and pinged MongoDB")
		}
	})

	if clientInstanceErr != nil {
		return nil, clientInstanceErr
	}

	return clientInstance.Database(dbName), nil
}
