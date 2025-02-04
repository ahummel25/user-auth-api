package db

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/v2/mongo"

	"github.com/ahummel25/user-auth-api/service/user"
)

// DBContextManager defines an interface for database context management
type DBContextManager interface {
	getCollections(ctx context.Context, collectionNames []CollectionName) (map[CollectionName]*mongo.Collection, error)
}

// SetupUserDBContext sets up the database context specifically for user operations
func SetupUserDBContext(ctx context.Context, dbManager DBContextManager) (context.Context, error) {
	if dbManager == nil {
		dbManager = globalDBManager
	}

	collectionsToGet := []CollectionName{usersCollection}
	collections, err := dbManager.getCollections(ctx, collectionsToGet)
	if err != nil {
		return nil, fmt.Errorf("failed to get collections: %w", err)
	}

	userCollection, exists := collections[usersCollection]
	if !exists {
		return nil, fmt.Errorf("users collection not found in retrieved collections")
	}

	return user.NewContext(ctx, user.GetUsersCollectionKey(), userCollection), nil
}

// SetupDBContext is maintained for backward compatibility
func SetupDBContext(ctx context.Context) (context.Context, error) {
	return SetupUserDBContext(ctx, globalDBManager)
}
