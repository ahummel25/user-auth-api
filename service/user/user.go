package user

import (
	"context"
	"errors"
	"time"

	"github.com/ahummel25/user-auth-api/db/mongo"
	"github.com/ahummel25/user-auth-api/graphql/model"
)

// API is the interface that wraps the methods for user operations.
type API interface {
	Login(ctx context.Context, usernameOrEmail string, password string) (*model.UserObject, error)
	CreateUser(ctx context.Context, params model.NewUserInput) (*model.UserObject, error)
	DeleteUser(ctx context.Context, userID string) (bool, error)
}

// UserCollection is an interface that wraps the database.Collection interface
type UserCollection interface {
	mongo.Collection
}

// userSvc serves as the function receiver for implementation of the API interface
type userSvc struct{}

// usersCollectionCtxKey represents the context key of the users Mongo collection
type usersCollectionCtxKey struct{}

type userDB struct {
	UserID        string     `bson:"user_id"`
	Email         string     `bson:"email"`
	FirstName     string     `bson:"first_name"`
	LastName      string     `bson:"last_name"`
	UserName      string     `bson:"user_name"`
	Role          model.Role `bson:"role"`
	Password      string     `bson:"password"`
	LastLoginDate *time.Time `bson:"last_login_date"`
}

// New returns a pointer to a new auth service.
func New() *userSvc {
	return &userSvc{}
}

func NewContext(ctx context.Context, collectionCtxKey any, collection UserCollection) context.Context {
	return context.WithValue(ctx, collectionCtxKey, collection)
}

// FromContext returns the UserCollection from the context, or an error if not found
func FromContext(ctx context.Context) (UserCollection, error) {
	if c, ok := ctx.Value(GetUsersCollectionKey()).(UserCollection); ok {
		return c, nil
	}
	return nil, errors.New("user collection not found in context")
}

// GetUsersCollectionKey is a wrapper function around the usersCollectionCtxKey returning a pointer to that value
func GetUsersCollectionKey() *usersCollectionCtxKey {
	return &usersCollectionCtxKey{}
}
