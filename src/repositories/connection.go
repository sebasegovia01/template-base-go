package repositories

import (
	"context"
	"log"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Declaración de variables globales para almacenar la instancia del cliente de MongoDB y controlar la inicialización única.
var (
	clientInstance     *mongo.Client
	clientInstanceOnce sync.Once
	clientInstanceErr  error
)

func Connect(uri string, dbName string) (*mongo.Database, error) {
	// Inicializar la instancia del cliente de MongoDB solo una vez.
	clientInstanceOnce.Do(func() {
		// Se ajusta el contexto para el proceso de conexión a MongoDB.
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second) // Ajuste de 8 segundos para el contexto inicial.
		defer cancel()

		// Usar Connect para crear y conectar el cliente en un solo paso.
		clientInstance, clientInstanceErr = mongo.Connect(ctx, options.Client().ApplyURI(uri))
		if clientInstanceErr != nil {
			log.Printf("Error connecting to MongoDB: %v\n", clientInstanceErr)
			return
		}

		// Verificar la conexión con Ping usando un contexto separado.
		pingCtx, pingCancel := context.WithTimeout(context.Background(), 60*time.Second) // Ajuste de 5 segundos para el ping.
		defer pingCancel()
		if clientInstanceErr = clientInstance.Ping(pingCtx, nil); clientInstanceErr != nil {
			log.Printf("Error pinging MongoDB: %v\n", clientInstanceErr)
			_ = clientInstance.Disconnect(context.Background()) // Desconectar sin contexto de timeout.
		} else {
			log.Println("Successfully connected and pinged MongoDB")
		}
	})

	if clientInstanceErr != nil {
		return nil, clientInstanceErr
	}

	return clientInstance.Database(dbName), nil
}
