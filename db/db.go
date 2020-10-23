package db

import (
	"context"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/src/user-auth-api/config"
)

// DBs
const (
	USERS_DB = "users"
)

// Collections
const (
	USERS_COLLECTION = "users"
)

var (
	connection     *mongo.Client
	collection     *mongo.Collection
	collectionToDB = map[string]string{
		USERS_COLLECTION: USERS_DB,
	}
)

// getDBConnection returns a mongo client connection.
func getDBConnection(ctx context.Context) (*mongo.Client, error) {
	serverAPIOptions := options.ServerAPI(options.ServerAPIVersion1)
	configSupplier, err := config.FromContext(ctx)
	if err != nil {
		return nil, err
	}
	cfg, err := configSupplier.GetConfig()
	if err != nil {
		return nil, err
	}
	dsn := fmt.Sprintf("mongodb+srv://%v:%v@%v.%v/?%v", cfg.UserName, cfg.Password, cfg.Cluster, cfg.Domain, "retryWrites=true&w=majority")
	clientOptions := options.Client().ApplyURI(dsn).SetServerAPIOptions(serverAPIOptions)
	return mongo.Connect(ctx, clientOptions)
}

// GetCollection returns the mongo collection based on the collection name parameter
func GetCollection(ctx context.Context, collectionName string) (*mongo.Collection, error) {
	if dbName, ok := collectionToDB[collectionName]; !ok {
		return nil, errors.New("invalid collection provided")
	} else {
		if connection != nil {
			if collection != nil {
				return collection, nil
			}
			collection = connection.Database(dbName).Collection(collectionName)
			return collection, nil
		}
		var err error
		if connection, err = getDBConnection(ctx); err != nil {
			return nil, errors.New("error connecting to DB")
		}
		collection = connection.Database(dbName).Collection(collectionName)
		return collection, nil
	}
}
