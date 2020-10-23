package user

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"golang.org/x/crypto/bcrypt"

	"github.com/ahummel25/user-auth-api/graphql/model"
)

var (
	errInvalidPassword   = errors.New("invalid password")
	errNoUserFound       = errors.New("user not found")
	errUserAlreadyExists = errors.New("user name or email already exists")
)

// Helper function to get user collection from context
func getUserCollection(ctx context.Context) (UserCollection, error) {
	userCollection, err := FromContext(ctx)
	if err != nil {
		return nil, err
	}
	return userCollection, nil
}

// Helper function to find user by username or email
func findUserByUsernameOrEmail(ctx context.Context, userCollection UserCollection, usernameOrEmail string) (*userDB, error) {
	filter := bson.M{
		"$or": []bson.M{
			{"email": usernameOrEmail},
			{"user_name": usernameOrEmail},
		}}
	var user userDB
	if err := userCollection.FindOne(ctx, filter).Decode(&user); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errNoUserFound
		}
		return nil, err
	}
	return &user, nil
}

// Helper function to update last login date
func updateLastLoginDate(ctx context.Context, userCollection UserCollection, userID string, time time.Time) error {
	update := bson.M{"$set": bson.M{"last_login_date": time}}
	if _, err := userCollection.UpdateOne(ctx, bson.M{"user_id": userID}, update); err != nil {
		return err
	}
	return nil
}

// Login authenticates the user.
func (u *userSvc) Login(ctx context.Context, usernameOrEmail string, password string) (*model.UserObject, error) {
	userCollection, err := getUserCollection(ctx)
	if err != nil {
		return nil, err
	}

	user, err := findUserByUsernameOrEmail(ctx, userCollection, usernameOrEmail)
	if err != nil {
		return nil, err
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return nil, errInvalidPassword
		}
		return nil, err
	}

	// Update last_login_date
	now := time.Now().UTC()
	if err = updateLastLoginDate(ctx, userCollection, user.UserID, now); err != nil {
		// Set last_login_date to the previously fetched value and log the error, but don't fail the login process
		now = user.LastLoginDate.UTC()
		slog.Error("Failed to update last_login_date", "error", err)
	}

	loggedInUser := &model.User{
		ID:            user.UserID,
		Email:         user.Email,
		FirstName:     user.FirstName,
		LastName:      user.LastName,
		UserName:      user.UserName,
		Role:          user.Role,
		LastLoginDate: &now,
	}
	userObject := &model.UserObject{User: loggedInUser}
	return userObject, nil
}

// CreateUser creates a new user.
func (u *userSvc) CreateUser(ctx context.Context, params model.NewUserInput) (*model.UserObject, error) {
	userCollection, err := getUserCollection(ctx)
	if err != nil {
		return nil, err
	}

	// Verify if the user name or email already exists
	filter := bson.M{
		"$or": []bson.M{
			{"email": params.Email},
			{"user_name": params.UserName},
		}}
	userCount, err := userCollection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, err
	} else if userCount > 0 {
		return nil, errUserAlreadyExists
	}

	// Generate a hash from the password to store in the DB
	hash, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	newUserID := uuid.New().String()
	// Default to user role unless one is provided
	role := model.RoleUser
	if params.Role != nil {
		role = *params.Role
	}
	now := time.Now().UTC()
	newUserInput := bson.D{
		{Key: "user_id", Value: newUserID},
		{Key: "email", Value: params.Email},
		{Key: "user_name", Value: params.UserName},
		{Key: "password", Value: string(hash)},
		{Key: "first_name", Value: params.FirstName},
		{Key: "last_name", Value: params.LastName},
		{Key: "role", Value: role},
		{Key: "creation_date", Value: now},
		{Key: "last_update_date", Value: now},
		{Key: "last_login_date", Value: nil}, // Initialize last_login_date as nil
	}

	if _, err = userCollection.InsertOne(ctx, newUserInput); err != nil {
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
func (u *userSvc) DeleteUser(ctx context.Context, userID string) (bool, error) {
	userCollection, err := getUserCollection(ctx)
	if err != nil {
		return false, err
	}
	filter := bson.M{"user_id": userID}
	deleteResult, err := userCollection.DeleteOne(ctx, filter)
	if err != nil {
		return false, err
	}
	if deleteResult.DeletedCount == 0 {
		return false, errNoUserFound
	}

	return true, nil
}
