package db

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-secretsmanager-caching-go/secretcache"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var secretCache, _ = secretcache.New()

const (
	authDB = "auth"
)

var (
	connection *mongo.Client
	collection *mongo.Collection
)

type dbConfig struct {
	UserName string `json:"DB_USER_NAME"`
	Password string `json:"DB_PASSWORD"`
	Cluster  string `json:"DB_CLUSTER_NAME"`
}

// getDBConnection returns a mongo client connection.
func getDBConnection(ctx context.Context) (*mongo.Client, error) {
	var (
		client *mongo.Client
		err    error
	)
	cfg, err := getDBConfig()
	if err != nil {
		return nil, err
	}
	dsn := fmt.Sprintf("mongodb+srv://%v:%v@%v.%v/?%v", cfg.UserName, cfg.Password, cfg.Cluster, "yohvj.mongodb.net", "retryWrites=true&w=majority")
	if client, err = mongo.Connect(ctx, options.Client().ApplyURI(dsn)); err != nil {
		return nil, err
	}
	return client, nil
}

// getDBConfig retrieves cached DB configuation secrets
func getDBConfig() (dbConfig, error) {
	var cfg dbConfig
	secretString, err := secretCache.GetSecretString(os.Getenv("SECRET_NAME"))
	if err != nil {
		return dbConfig{}, err
	}
	if err = json.NewDecoder(strings.NewReader(secretString)).Decode(&cfg); err != nil {
		return dbConfig{}, err
	}
	return cfg, nil
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
	db := connection.Database(authDB)
	collection = db.Collection(collectionName)
	return collection, nil
}
