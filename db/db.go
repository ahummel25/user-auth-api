package db

import (
	"context"
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
func GetDBConnection() (*mongo.Client, context.Context, context.CancelFunc, error) {
	var (
		client   *mongo.Client
		err      error
		mongoURI = os.Getenv("MONGODB_URI")
	)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

	if client, err = mongo.Connect(ctx, options.Client().ApplyURI(mongoURI)); err != nil {
		return nil, ctx, cancel, err
	}

	return client, ctx, cancel, nil
}
