package services

import (
	"context"
	"errors"
	"log"
	"strings"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"

	dbHelper "github.com/src/user-auth-api/db"
	"github.com/src/user-auth-api/graph/model"
)

// UserService contains signatures for any auth functions.
type UserService interface {
	AuthenticateUser(email string, password string) (*model.User, error)
	CreateUser(params model.CreateUserInput) (*model.User, error)
}

type userService struct{}

var (
	errInvalidPassword       = errors.New("invalid password")
	errNoUserFound           = errors.New("no user found by that email address")
	errUserNameAlreadyExists = errors.New("user name already exists")
)

type userDB struct {
	model.User
	Password string
}

// NewUserService returns a pointer to a new auth service.
func NewUserService() *userService {
	return &userService{}
}

func (u *userService) getUsersCollection() (context.Context, func(), *mongo.Collection) {
	var err error
	conn, ctx, cancel := dbHelper.GetDBConnection()

	cancelFunc := func() {
		if err = conn.Disconnect(ctx); err != nil {
			log.Fatal(err)
		}
		cancel()
	}

	db := conn.Database("default")
	return ctx, cancelFunc, db.Collection("users")
}

// AuthenticateUser authenticates the user.
func (u *userService) AuthenticateUser(email string, password string) (*model.User, error) {
	var (
		err    error
		userDB bson.M
	)

	ctx, cancelFunc, usersCollection := u.getUsersCollection()

	defer cancelFunc()

	filter := bson.M{"email": email}

	if err = usersCollection.FindOne(ctx, filter).Decode(&userDB); err != nil {
		log.Printf("Error while finding user: %v\n", err)

		if strings.Contains(err.Error(), "no documents in result") {
			return nil, errNoUserFound
		}

		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(userDB["password"].(string)), []byte(password)); err != nil {
		log.Printf("Error comparing user password on user login: %v\n", err)

		if strings.Contains(err.Error(), "not the hash of the given password") {
			return nil, errInvalidPassword
		}

		return nil, err
	}

	loggedInUser := &model.User{
		UserID:    userDB["user_id"].(string),
		Email:     userDB["email"].(string),
		FirstName: userDB["first_name"].(string),
		LastName:  userDB["last_name"].(string),
		UserName:  userDB["user_name"].(string),
	}

	return loggedInUser, nil
}

// CreateUser authenticates the user.
func (u *userService) CreateUser(params model.CreateUserInput) (*model.User, error) {
	var (
		err       error
		hash      []byte
		userCount int64
	)

	ctx, cancelFunc, usersCollection := u.getUsersCollection()

	defer cancelFunc()

	filter := bson.M{"user_name": params.UserName}

	if userCount, err = usersCollection.CountDocuments(ctx, filter); err != nil {
		return nil, err
	}

	if userCount > 0 {
		return nil, errUserNameAlreadyExists
	}

	if hash, err = bcrypt.GenerateFromPassword([]byte(params.Password), bcrypt.DefaultCost); err != nil {
		log.Printf("Error while encrypting user password on user creation: %v\n", err)
		return nil, err
	}

	newUserID := uuid.New().String()

	userInput := bson.D{
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

	if _, err = usersCollection.InsertOne(ctx, userInput); err != nil {
		return nil, err
	}

	user := &model.User{
		UserID:    newUserID,
		Email:     params.Email,
		FirstName: params.FirstName,
		LastName:  params.LastName,
		UserName:  params.UserName,
	}

	return user, nil
}
