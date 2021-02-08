package services

import (
	"context"
	"errors"
	"log"
	"strings"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"

	dbHelper "github.com/src/user-auth-api/db"
	"github.com/src/user-auth-api/graphql/model"
)

// UserService contains signatures for any auth functions.
type UserService interface {
	AuthenticateUser(email string, password string) (*model.UserObject, error)
	CreateUser(params model.CreateUserInput) (*model.UserObject, error)
	DeleteUser(params model.DeleteUserInput) (string, error)
}

type User struct{}

var (
	errInvalidPassword       = errors.New("invalid password")
	errNoUserFound           = errors.New("no user found!")
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

// NewUserService returns a pointer to a new auth service.
func NewUserService() *User {
	return &User{}
}

func (u *User) getUsersCollection() (context.Context, func(), *mongo.Collection, error) {
	var (
		cancel context.CancelFunc
		conn   *mongo.Client
		ctx    context.Context
		err    error
	)

	if conn, ctx, cancel, err = dbHelper.GetDBConnection(); err != nil {
		log.Printf("Error connecting to MongoDB: %v\n", err)

		return nil, nil, nil, errors.New("error connecting to DB")
	}

	cancelFunc := func() {
		if err = conn.Disconnect(ctx); err != nil {
			log.Printf("Error disconnecting from MongoDB: %v\n", err)
		}
		cancel()
	}

	db := conn.Database("auth")
	usersCollection := db.Collection("users")

	if _, err = usersCollection.Indexes().CreateMany(
		ctx,
		[]mongo.IndexModel{
			{
				Keys: bson.M{
					"user_id": 1,
				},
				Options: options.Index().SetUnique(true),
			}, {
				Keys: bson.M{
					"user_name": 1,
				},
				Options: options.Index().SetUnique(true),
			},
		},
	); err != nil {
		log.Printf("Error creating index on users collection: %v\n", err)

		return nil, nil, nil, errors.New("error connecting to DB")
	}

	return ctx, cancelFunc, usersCollection, nil
}

// AuthenticateUser authenticates the user.
func (u *User) AuthenticateUser(email string, password string) (*model.UserObject, error) {
	var (
		err    error
		userDB userDB
	)

	log.Printf("In AuthenticateUser")

	ctx, cancelFunc, usersCollection, err := u.getUsersCollection()

	log.Printf("Got users collection")

	defer cancelFunc()

	if err != nil {
		return nil, err
	}

	filter := bson.M{"email": email}

	if err = usersCollection.FindOne(ctx, filter).Decode(&userDB); err != nil {
		log.Printf("Error while finding user: %v\n", err)

		if strings.Contains(err.Error(), "no documents in result") {
			return nil, errNoUserFound
		}

		return nil, err
	}

	if err = bcrypt.CompareHashAndPassword([]byte(userDB.Password), []byte(password)); err != nil {
		log.Printf("Error comparing user password on user login: %v\n", err)

		if strings.Contains(err.Error(), "not the hash of the given password") {
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
func (u *User) CreateUser(params model.CreateUserInput) (*model.UserObject, error) {
	var (
		err       error
		hash      []byte
		userCount int64
	)

	ctx, cancelFunc, usersCollection, err := u.getUsersCollection()

	defer cancelFunc()

	if err != nil {
		return nil, err
	}

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
func (u *User) DeleteUser(params model.DeleteUserInput) (string, error) {
	var (
		deleteResult *mongo.DeleteResult
		err          error
	)

	ctx, cancelFunc, usersCollection, err := u.getUsersCollection()

	defer cancelFunc()

	if err != nil {
		return "", err
	}

	filter := bson.M{"user_id": params.UserID}

	if deleteResult, err = usersCollection.DeleteOne(ctx, filter); err != nil {
		return "", err
	}

	if deleteResult.DeletedCount == 0 {
		return "", errNoUserFound
	}

	return params.UserName + " successfully deleted", nil
}
