package user

import (
	"context"
	"errors"
	"log"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"

	dbHelper "github.com/src/user-auth-api/db"
	"github.com/src/user-auth-api/graphql/model"
)

// API contains signatures for any auth functions.
type API interface {
	AuthenticateUser(ctx context.Context, username string, password string) (*model.UserObject, error)
	CreateUser(ctx context.Context, params model.CreateUserInput) (*model.UserObject, error)
	DeleteUser(ctx context.Context, params model.DeleteUserInput) (bool, error)
}

type User struct{}

// DB names
const (
	authDB = "auth"
)

var (
	conn                  *mongo.Client
	usersCollection       *mongo.Collection
	UsersCollectionCtxKey = "usersDB"
)

var (
	errInvalidPassword       = errors.New("invalid password")
	errNoUserFound           = errors.New("user not found")
	errUserNameAlreadyExists = errors.New("user name already exists")
)

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
	return context.WithValue(ctx, &UsersCollectionCtxKey, m)
}

// fromContext returns the *mongo.Collection that was stored in the context, or nil if none was stored.
func fromContext(ctx context.Context) *mongo.Collection {
	if s, ok := ctx.Value(&UsersCollectionCtxKey).(*mongo.Collection); ok {
		return s
	}
	return nil
}

func GetUsersCollection(ctx context.Context) (*mongo.Collection, error) {
	var err error
	if conn != nil && usersCollection != nil {
		log.Println("Connection and collection are already active!")
		return usersCollection, nil
	}
	if conn, err = dbHelper.GetDBConnection(ctx); err != nil {
		log.Printf("Error connecting to MongoDB: %v\n", err)
		return nil, errors.New("error connecting to DB")
	}
	db := conn.Database(authDB)
	usersCollection = db.Collection("users")
	return usersCollection, nil
}

// AuthenticateUser authenticates the user.
func (u *User) AuthenticateUser(ctx context.Context, username string, password string) (*model.UserObject, error) {
	var err error
	userDB := &userDB{}
	usersCollection := fromContext(ctx)
	filter := bson.M{"user_name": username}
	if err = usersCollection.FindOne(ctx, filter).Decode(userDB); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errNoUserFound
		}
		return nil, err
	}

	if err = bcrypt.CompareHashAndPassword([]byte(userDB.Password), []byte(password)); err != nil {
		log.Printf("Error comparing user password on user login: %v\n", err)
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return nil, errInvalidPassword
		}
		return nil, err
	}

	loggedInUser := &model.User{
		UserID:    userDB.UserID,
		Email:     userDB.Email,
		FirstName: userDB.FirstName,
		LastName:  userDB.LastName,
		UserName:  userDB.UserName,
	}
	user := &model.UserObject{User: loggedInUser}
	return user, nil
}

// CreateUser creates a new user.
func (u *User) CreateUser(ctx context.Context, params model.CreateUserInput) (*model.UserObject, error) {
	var (
		err       error
		hash      []byte
		userCount int64
	)

	usersCollection := fromContext(ctx)
	// Verify if the user name already exists
	filter := bson.M{"user_name": params.UserName}
	if userCount, err = usersCollection.CountDocuments(ctx, filter); err != nil {
		return nil, err
	}
	if userCount > 0 {
		return nil, errUserNameAlreadyExists
	}

	// Generate a hash from the password to store in the DB
	if hash, err = bcrypt.GenerateFromPassword([]byte(params.Password), bcrypt.DefaultCost); err != nil {
		log.Printf("Error while encrypting user password on user creation: %v\n", err)
		return nil, err
	}
	newUserID := uuid.New().String()
	newUserInput := bson.D{
		{
			Key: "user_id", Value: newUserID,
		},
		{
			Key: "email", Value: params.Email,
		},
		{
			Key: "first_name", Value: params.FirstName,
		},
		{
			Key: "last_name", Value: params.LastName,
		},
		{
			Key: "user_name", Value: params.UserName,
		},
		{
			Key: "password", Value: string(hash),
		},
	}

	if _, err = usersCollection.InsertOne(ctx, newUserInput); err != nil {
		return nil, err
	}

	newUser := &model.User{
		UserID:    newUserID,
		Email:     params.Email,
		FirstName: params.FirstName,
		LastName:  params.LastName,
		UserName:  params.UserName,
	}
	user := &model.UserObject{User: newUser}
	return user, nil
}

// DeleteUser deletes an existing user.
func (u *User) DeleteUser(ctx context.Context, params model.DeleteUserInput) (bool, error) {
	var (
		deleteResult *mongo.DeleteResult
		err          error
	)
	usersCollection := fromContext(ctx)
	filter := bson.M{"user_id": params.UserID}
	if deleteResult, err = usersCollection.DeleteOne(ctx, filter); err != nil {
		return false, err
	}
	if deleteResult.DeletedCount == 0 {
		return false, errNoUserFound
	}

	return true, nil
}
