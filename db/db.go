package db

import (
	"context"
	"errors"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/src/user-auth-api/config"
)

const (
	usersDB = "users"
)

var (
	connection *mongo.Client
	collection *mongo.Collection
)

// getDBConnection returns a mongo client connection.
func getDBConnection(ctx context.Context) (*mongo.Client, error) {
	serverAPIOptions := options.ServerAPI(options.ServerAPIVersion1)
	cfg, err := config.FromContext(ctx)
	if err != nil {
		return nil, err
	}
	dsn := fmt.Sprintf("mongodb+srv://%v:%v@%v.%v/?%v", cfg.UserName, cfg.Password, cfg.Cluster, cfg.Domain, "retryWrites=true&w=majority")
	clientOptions := options.Client().ApplyURI(dsn).SetServerAPIOptions(serverAPIOptions)
	return mongo.Connect(ctx, clientOptions)
}

// GetCollection returns the collection name based on the parameter
func GetCollection(ctx context.Context, collectionName string) (*mongo.Collection, error) {
	var err error
	if connection != nil && collection != nil {
		log.Println("Connection and collection are already active!")
		return collection, nil
	}
	if connection, err = getDBConnection(ctx); err != nil {
		log.Printf("Error connecting to MongoDB: %v\n", err)
		return nil, errors.New("error connecting to DB")
	}
	db := connection.Database(usersDB)
	collection = db.Collection(collectionName)
	return collection, nil
}
