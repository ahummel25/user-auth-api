package db

import (
	"context"
	"fmt"
	"sync"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/ahummel25/user-auth-api/config"
)

// Constants for database and collection names
const (
	UsersDB         = "users"
	UsersCollection = "users"
)

// DBManager manages the database connection and collections
type DBManager struct {
	connection     *mongo.Client
	collection     *mongo.Collection
	collectionToDB map[string]string
	mu             sync.Mutex
}

// NewDBManager creates a new DBManager
func NewDBManager(collectionToDB map[string]string) *DBManager {
	return &DBManager{
		collectionToDB: collectionToDB,
	}
}

// Global instance of DBManager
var globalDBManager = NewDBManager(map[string]string{
	UsersCollection: UsersDB,
})

// getDBConnection returns a mongo client connection.
func (m *DBManager) getDBConnection(ctx context.Context) (*mongo.Client, error) {
	configSupplier, err := config.FromContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get config from context: %w", err)
	}

	cfg, err := configSupplier.GetConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to get config: %w", err)
	}

	dsn := fmt.Sprintf("mongodb+srv://%s:%s@%s.%s/?retryWrites=true&w=majority",
		cfg.UserName, cfg.Password, cfg.Cluster, cfg.Domain)
	clientOptions := options.Client().
		ApplyURI(dsn).
		SetServerAPIOptions(options.ServerAPI(options.ServerAPIVersion1))

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	// Ping the database to verify the connection
	if err = client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	return client, nil
}

// getCollection returns the mongo collection based on the collection name parameter
func getCollection(ctx context.Context, collectionName string) (*mongo.Collection, error) {
	globalDBManager.mu.Lock()
	defer globalDBManager.mu.Unlock()

	if globalDBManager.connection != nil && globalDBManager.collection != nil {
		return globalDBManager.collection, nil
	}

	dbName, ok := globalDBManager.collectionToDB[collectionName]
	if !ok {
		return nil, fmt.Errorf("invalid collection provided: %s", collectionName)
	}

	var err error
	if globalDBManager.connection == nil {
		globalDBManager.connection, err = globalDBManager.getDBConnection(ctx)
		if err != nil {
			return nil, fmt.Errorf("error connecting to DB: %w", err)
		}
	}

	globalDBManager.collection = globalDBManager.connection.Database(dbName).Collection(collectionName)
	return globalDBManager.collection, nil
}

// CloseConnection closes the database connection
func CloseConnection(ctx context.Context) error {
	globalDBManager.mu.Lock()
	defer globalDBManager.mu.Unlock()

	if globalDBManager.connection != nil {
		if err := globalDBManager.connection.Disconnect(ctx); err != nil {
			return fmt.Errorf("error disconnecting from DB: %w", err)
		}
		globalDBManager.connection = nil
		globalDBManager.collection = nil
	}
	return nil
}
