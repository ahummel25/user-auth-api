package user

import (
	"context"
	"errors"
	"log"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"

	"github.com/src/user-auth-api/graphql/model"
	"github.com/src/user-auth-api/service"
)

var (
	errInvalidPassword       = errors.New("invalid password")
	errNoUserFound           = errors.New("user not found")
	errUserNameAlreadyExists = errors.New("user name already exists")
)

// AuthenticateUser authenticates the user.
func (u *User) AuthenticateUser(ctx context.Context, username string, password string) (*model.UserObject, error) {
	var err error
	userDB := &userDB{}
	usersCollection := service.FromContext(ctx, &usersCollectionCtxKey)
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

	usersCollection := service.FromContext(ctx, &usersCollectionCtxKey)
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
	usersCollection := service.FromContext(ctx, &usersCollectionCtxKey)
	filter := bson.M{"user_id": params.UserID}
	if deleteResult, err = usersCollection.DeleteOne(ctx, filter); err != nil {
		return false, err
	}
	if deleteResult.DeletedCount == 0 {
		return false, errNoUserFound
	}

	return true, nil
}