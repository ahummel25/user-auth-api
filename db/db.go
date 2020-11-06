package db

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// IDBService contains signatures for any DB functions.
type IDBService interface {
	GetDBConnection() (*mongo.Client, context.Context)
}

type dbService struct{}

// GetDBConnection will return a mongo client connection.
func GetDBConnection() (*mongo.Client, context.Context, context.CancelFunc) {
	var (
		client   *mongo.Client
		err      error
		mongoURI = os.Getenv("MONGODB_URI")
	)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	if client, err = mongo.Connect(ctx, options.Client().ApplyURI(mongoURI)); err != nil {
		log.Fatal(err)
	}

	return client, ctx, cancel
}
