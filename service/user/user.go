package user

import (
	"context"

	"github.com/src/user-auth-api/graphql/model"
	"github.com/src/user-auth-api/service"
)

//go:generate mockery --name API
// Install mockery from https://github.com/vektra/mockery

// API contains signatures for any auth functions.
type API interface {
	AuthenticateUser(ctx context.Context, username string, password string) (*model.UserObject, error)
	CreateUser(ctx context.Context, params model.CreateUserInput) (*model.UserObject, error)
	DeleteUser(ctx context.Context, params model.DeleteUserInput) (bool, error)
}

type User struct{}

var usersCollectionCtxKey service.CollectionCtxKey = "usersDB"

type userDB struct {
	UserID    string `bson:"user_id"`
	Email     string `bson:"email"`
	FirstName string `bson:"first_name"`
	LastName  string `bson:"last_name"`
	UserName  string `bson:"user_name"`
	Password  string `bson:"password"`
}

// New returns a pointer to a new auth service.
func New() *User {
	return &User{}
}

// GetUsersCollectionKey is a wrapper function around the usersCollectionCtxKey returning a pointer to that value
func GetUsersCollectionKey() *service.CollectionCtxKey {
	return &usersCollectionCtxKey
}
