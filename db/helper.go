package db

import (
	"context"

	"github.com/ahummel25/user-auth-api/service/user"
)

// SetupDBContext sets up the database context by getting the required collections and adding them to the context
func SetupDBContext(ctx context.Context) (context.Context, error) {
	collectionsToGet := []CollectionName{usersCollection}
	collections, err := globalDBManager.getCollections(ctx, collectionsToGet)
	if err != nil {
		return nil, err
	}
	return user.NewContext(ctx, user.GetUsersCollectionKey(), collections[usersCollection]), nil
}
