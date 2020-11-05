package db

import (
	"context"
	"log"
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
		mongoURI = "mongodb+srv://user_auth_admin:rlippi7-yxyeEr@userauthmongoclusterdev.yohvj.mongodb.net/default?retryWrites=true&w=majority"
	)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	if client, err = mongo.Connect(ctx, options.Client().ApplyURI(mongoURI)); err != nil {
		log.Fatal(err)
	}

	return client, ctx, cancel
}
