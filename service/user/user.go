package user

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/src/user-auth-api/graphql/model"
)

// API contains signatures for any auth functions.
type API interface {
	AuthenticateUser(ctx context.Context, username string, password string) (*model.UserObject, error)
	CreateUser(ctx context.Context, params model.CreateUserInput) (*model.UserObject, error)
	DeleteUser(ctx context.Context, params model.DeleteUserInput) (bool, error)
}

type User struct{}

var usersCollectionCtxKey = "usersDB"

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

// NewContext returns a new Context that carries value s.
func NewContext(ctx context.Context, m *mongo.Collection) context.Context {
	return context.WithValue(ctx, &usersCollectionCtxKey, m)
}

// fromContext returns the *mongo.Collection that was stored in the context, or nil if none was stored.
func fromContext(ctx context.Context) *mongo.Collection {
	if s, ok := ctx.Value(&usersCollectionCtxKey).(*mongo.Collection); ok {
		return s
	}
	return nil
}
