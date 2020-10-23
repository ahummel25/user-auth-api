package db

import (
	"context"
	"fmt"
	"sync"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/ahummel25/user-auth-api/config"
)

// collectionToDBMap maps collection names to their respective database names
var collectionToDBMap = map[string]string{
	"users": "users",
}

// DBManager manages the database connection and collections
type DBManager struct {
	connection     *mongo.Client
	collections    map[string]*mongo.Collection
	collectionToDB map[string]string
	mu             sync.Mutex
}

// NewDBManager creates a new DBManager
func NewDBManager() *DBManager {
	return &DBManager{
		collectionToDB: collectionToDBMap,
		collections:    make(map[string]*mongo.Collection),
	}
}

// Global instance of DBManager
var globalDBManager = NewDBManager()

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

// getCollections fetches and stores multiple collections at once
func (m *DBManager) getCollections(ctx context.Context, collectionNames []string) (map[string]*mongo.Collection, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.connection == nil {
		var err error
		m.connection, err = m.getDBConnection(ctx)
		if err != nil {
			return nil, fmt.Errorf("error connecting to DB: %w", err)
		}
	}

	for _, collectionName := range collectionNames {
		if _, exists := m.collections[collectionName]; !exists {
			dbName, ok := m.collectionToDB[collectionName]
			if !ok {
				return nil, fmt.Errorf("invalid collection provided: %s", collectionName)
			}
			m.collections[collectionName] = m.connection.Database(dbName).Collection(collectionName)
		}
	}

	return m.collections, nil
}
