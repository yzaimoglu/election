package models

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/yzaimoglu/election/info/utilities"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Connection Parameters
const (
	connectTimeout           = 5
	connectionStringTemplate = "mongodb://%s:%s@%s"
)

// Get a Mongo instance (Client, Context, Cancel)
func GetMongoInstance() (*mongo.Client, context.Context, context.CancelFunc) {
	// MongoDB Credentials from .env
	username := utilities.GetEnv("BILGI_DB_USER", "admin")
	password := utilities.GetEnv("BILGI_DB_PASSWORD", "admin")
	hostname := utilities.GetEnv("BILGI_HOSTNAME", "localhost")

	// Connection URI for the database
	connectionURI := fmt.Sprintf(connectionStringTemplate, username, password, hostname)

	// Create the mongo client
	client, err := mongo.NewClient(options.Client().ApplyURI(connectionURI))
	if err != nil {
		log.Printf("Failed to create client: %v", err)
	}

	// Create the context and the cancel function
	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout*time.Second)

	// Connect to MongoDB and check for connection error
	err = client.Connect(ctx)
	if err != nil {
		log.Printf("Failed to connect to the database: %v", err)
	}

	// Force a connection to verify our connection string
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Printf("Failed to ping the database: %v", err)
	}

	// Return mongo instance
	return client, ctx, cancel
}
