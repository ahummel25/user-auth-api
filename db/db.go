package db

import (
	"context"
	"fmt"
	"net/url"
	"sync"

	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials/stscreds"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	"github.com/ahummel25/user-auth-api/config"
)

type DBName string
type CollectionName string

var (
	// DB and collection names for users
	usersDB         DBName         = "users"
	usersCollection CollectionName = "users"
)

// collectionToDBMap maps collection names to their respective database names
var collectionToDBMap = map[CollectionName]DBName{
	usersCollection: usersDB,
}

// DBManager manages the database connection and collections
type DBManager struct {
	connection     *mongo.Client
	collections    map[CollectionName]*mongo.Collection
	collectionToDB map[CollectionName]DBName
	mu             sync.Mutex
}

// NewDBManager creates a new DBManager
func NewDBManager() *DBManager {
	return &DBManager{
		collectionToDB: collectionToDBMap,
		collections:    make(map[CollectionName]*mongo.Collection),
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

	// Load AWS SDK configuration
	awsCfg, err := awsConfig.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	// Create STS client
	stsClient := sts.NewFromConfig(awsCfg)

	// Assume IAM role to get temporary credentials
	provider := stscreds.NewAssumeRoleProvider(stsClient, cfg.IAMRoleARN)

	creds, err := provider.Retrieve(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to assume IAM role: %w", err)
	}

	// URL encode the credentials
	encodedAccessKeyID := url.QueryEscape(creds.AccessKeyID)
	encodedSecretAccessKey := url.QueryEscape(creds.SecretAccessKey)
	encodedSessionToken := url.QueryEscape(creds.SessionToken)

	// Construct MongoDB connection string with encoded IAM credentials
	dsn := fmt.Sprintf("mongodb+srv://%s:%s@%s.%s/?authSource=%%24external&authMechanism=MONGODB-AWS&retryWrites=true&w=majority&authMechanismProperties=AWS_SESSION_TOKEN:%s&appName=%s&readPreference=secondary&ssl=true&logLevel=1",
		encodedAccessKeyID,
		encodedSecretAccessKey,
		cfg.Cluster,
		cfg.Domain,
		encodedSessionToken,
		cfg.AppName)

	clientOptions := options.Client().
		ApplyURI(dsn).
		SetServerAPIOptions(options.ServerAPI(options.ServerAPIVersion1))

	client, err := mongo.Connect(clientOptions)
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
func (m *DBManager) getCollections(ctx context.Context, collectionNames []CollectionName) (map[CollectionName]*mongo.Collection, error) {
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
			m.collections[collectionName] = m.connection.Database(string(dbName)).Collection(string(collectionName))
		}
	}

	return m.collections, nil
}
