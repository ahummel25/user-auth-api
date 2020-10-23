package user

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"

	"github.com/src/user-auth-api/graphql/model"
	"github.com/src/user-auth-api/service"
)

var (
	errInvalidPassword   = errors.New("invalid password")
	errNoUserFound       = errors.New("user not found")
	errUserAlreadyExists = errors.New("user name or email already exists")
)

// AuthenticateUser authenticates the user.
func (u userSvc) AuthenticateUser(ctx context.Context, usernameOrEmail string, password string) (*model.UserObject, error) {
	var user userDB
	usersCollection := service.FromContext(ctx, &usersCollectionCtxKey{})
	filter := bson.M{
		"$or": []bson.M{
			{"email": usernameOrEmail},
			{"user_name": usernameOrEmail},
		}}
	var err error
	if err = usersCollection.FindOne(ctx, filter).Decode(&user); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errNoUserFound
		}
		return nil, err
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		log.Printf("Error comparing user password on user login: %v\n", err)
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return nil, errInvalidPassword
		}
		return nil, err
	}

	loggedInUser := &model.User{
		ID:        user.UserID,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		UserName:  user.UserName,
		Role:      user.Role,
	}
	userObject := &model.UserObject{User: loggedInUser}
	return userObject, nil
}

// CreateUser creates a new user.
func (u userSvc) CreateUser(ctx context.Context, params model.NewUserInput) (*model.UserObject, error) {
	var (
		err       error
		hash      []byte
		userCount int64
	)

	usersCollection := service.FromContext(ctx, &usersCollectionCtxKey{})
	// Verify if the user name or email already exists
	filter := bson.M{
		"$or": []bson.M{
			{"email": params.Email},
			{"user_name": params.UserName},
		}}
	if userCount, err = usersCollection.CountDocuments(ctx, filter); err != nil {
		return nil, err
	} else if userCount > 0 {
		return nil, errUserAlreadyExists
	}

	// Generate a hash from the password to store in the DB
	if hash, err = bcrypt.GenerateFromPassword([]byte(params.Password), bcrypt.DefaultCost); err != nil {
		log.Printf("Error while encrypting user password on user creation: %v\n", err)
		return nil, err
	}
	newUserID := uuid.New().String()
	// Default to user role unless one is provided
	role := model.RoleUser
	if params.Role != nil {
		role = *params.Role
	}
	now := time.Now()
	newUserInput := bson.D{
		{
			Key: "user_id", Value: newUserID,
		},
		{
			Key: "email", Value: params.Email,
		},
		{
			Key: "user_name", Value: params.UserName,
		},
		{
			Key: "password", Value: string(hash),
		},
		{
			Key: "first_name", Value: params.FirstName,
		},
		{
			Key: "last_name", Value: params.LastName,
		},
		{
			Key: "role", Value: role,
		},
		{
			Key: "creation_date", Value: now,
		},
		{
			Key: "last_update_date", Value: now,
		},
	}

	if _, err = usersCollection.InsertOne(ctx, newUserInput); err != nil {
		return nil, err
	}

	newUser := &model.User{
		ID:        newUserID,
		Email:     params.Email,
		FirstName: params.FirstName,
		LastName:  params.LastName,
		UserName:  params.UserName,
		Role:      role,
	}
	user := &model.UserObject{User: newUser}
	return user, nil
}

// DeleteUser deletes an existing user.
func (u userSvc) DeleteUser(ctx context.Context, userID string) (bool, error) {
	var (
		deleteResult *mongo.DeleteResult
		err          error
	)
	usersCollection := service.FromContext(ctx, &usersCollectionCtxKey{})
	filter := bson.M{"user_id": userID}
	if deleteResult, err = usersCollection.DeleteOne(ctx, filter); err != nil {
		return false, err
	}
	if deleteResult.DeletedCount == 0 {
		return false, errNoUserFound
	}

	return true, nil
}
