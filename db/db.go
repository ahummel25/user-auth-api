package db

import (
	"context"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// GetDBConnection will return a mongo client connection.
func GetDBConnection(ctx context.Context) (*mongo.Client, error) {
	var (
		client   *mongo.Client
		err      error
		mongoURI = os.Getenv("MONGODB_URI")
	)
	if client, err = mongo.Connect(ctx, options.Client().ApplyURI(mongoURI)); err != nil {
		return nil, err
	}
	return client, nil
}
