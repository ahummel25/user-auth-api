package user

import (
	"context"
	"errors"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"golang.org/x/crypto/bcrypt"

	"github.com/ahummel25/user-auth-api/graphql/model"
	userMocks "github.com/ahummel25/user-auth-api/service/user/mocks"
	"github.com/ahummel25/user-auth-api/testutils"
)

func TestMain(m *testing.M) {
	// Set a fixed time for all tests
	fixedTime := time.Date(2024, 8, 30, 23, 41, 18, 0, time.UTC)
	testutils.SetFixedTime(fixedTime)

	// Run the tests
	code := m.Run()

	// Reset the time after tests
	testutils.ResetTime()

	os.Exit(code)
}

func TestLogin(t *testing.T) {
	t.Run("successful authentication", func(t *testing.T) {
		mockColl := userMocks.NewMockUserCollection(t)
		ctx := createContextWithMockCollection(mockColl)

		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
		user := userDB{
			UserID:    "test-id",
			Email:     "test@example.com",
			UserName:  "testuser",
			Password:  string(hashedPassword),
			FirstName: "Test",
			LastName:  "User",
			Role:      model.RoleUser,
		}

		expectedFilter := bson.M{
			"$or": []bson.M{
				{"email": "testuser"},
				{"user_name": "testuser"},
			},
		}

		// Create a mock SingleResult
		mockResult := mongo.NewSingleResultFromDocument(user, nil, nil)
		// Update the mock expectations
		mockColl.On("FindOne", ctx, expectedFilter).Return(mockResult)

		// Add expectation for UpdateOne
		updateFilter := bson.M{"user_id": user.UserID}
		mockColl.On("UpdateOne", ctx, updateFilter, mock.AnythingOfType("bson.M")).Return(&mongo.UpdateResult{}, nil)

		userSvc := &userSvc{}
		result, err := userSvc.Login(ctx, "testuser", "password")

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, user.UserID, result.User.ID)
		assert.Equal(t, user.Email, result.User.Email)
		assert.NotNil(t, result.User.LastLoginDate)
		mockColl.AssertExpectations(t)
	})

	t.Run("user not found", func(t *testing.T) {
		mockColl := userMocks.NewMockUserCollection(t)
		ctx := createContextWithMockCollection(mockColl)

		expectedFilter := bson.M{
			"$or": []bson.M{
				{"email": "nonexistent@example.com"},
				{"user_name": "nonexistent@example.com"},
			},
		}

		mockResult := mongo.NewSingleResultFromDocument(userDB{}, mongo.ErrNoDocuments, nil)
		mockColl.On("FindOne", ctx, expectedFilter).Return(mockResult)

		userSvc := &userSvc{}
		result, err := userSvc.Login(ctx, "nonexistent@example.com", "password")

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, errNoUserFound, err)
		mockColl.AssertExpectations(t)
	})

	t.Run("database error", func(t *testing.T) {
		mockColl := userMocks.NewMockUserCollection(t)
		ctx := createContextWithMockCollection(mockColl)

		expectedFilter := bson.M{
			"$or": []bson.M{
				{"email": "testuser"},
				{"user_name": "testuser"},
			},
		}

		// Simulate a database error on FindOne
		mockResult := mongo.NewSingleResultFromDocument(userDB{}, errors.New("database error"), nil)
		mockColl.On("FindOne", ctx, expectedFilter).Return(mockResult)

		userSvc := &userSvc{}
		result, err := userSvc.Login(ctx, "testuser", "password")

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "database error")
		mockColl.AssertExpectations(t)
	})

	t.Run("invalid password", func(t *testing.T) {
		mockColl := userMocks.NewMockUserCollection(t)
		ctx := createContextWithMockCollection(mockColl)

		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("correctpassword"), bcrypt.DefaultCost)
		user := userDB{
			UserID:   "test-id",
			Email:    "test@example.com",
			UserName: "testuser",
			Password: string(hashedPassword),
		}

		expectedFilter := bson.M{
			"$or": []bson.M{
				{"email": "testuser"},
				{"user_name": "testuser"},
			},
		}

		mockResult := mongo.NewSingleResultFromDocument(user, nil, nil)
		mockColl.On("FindOne", ctx, expectedFilter).Return(mockResult)

		userSvc := &userSvc{}
		result, err := userSvc.Login(ctx, "testuser", "wrongpassword")

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, errInvalidPassword, err)
		mockColl.AssertExpectations(t)
	})

	t.Run("update last login date fails", func(t *testing.T) {
		mockColl := userMocks.NewMockUserCollection(t)
		ctx := createContextWithMockCollection(mockColl)

		oldLoginDate := testutils.CurrentTime.Now().Add(-24 * time.Hour)
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
		user := userDB{
			UserID:        "test-id",
			Email:         "test@example.com",
			UserName:      "testuser",
			Password:      string(hashedPassword),
			FirstName:     "Test",
			LastName:      "User",
			Role:          model.RoleUser,
			LastLoginDate: &oldLoginDate,
		}

		expectedFilter := bson.M{
			"$or": []bson.M{
				{"email": "testuser"},
				{"user_name": "testuser"},
			},
		}

		mockResult := mongo.NewSingleResultFromDocument(user, nil, nil)
		mockColl.On("FindOne", ctx, expectedFilter).Return(mockResult)

		// Add expectation for UpdateOne to fail
		updateFilter := bson.M{"user_id": user.UserID}
		mockColl.On("UpdateOne", ctx, updateFilter, mock.AnythingOfType("bson.M")).Return(nil, errors.New("update failed"))

		userSvc := &userSvc{}
		result, err := userSvc.Login(ctx, "testuser", "password")

		assert.NoError(t, err) // Login should still succeed
		assert.NotNil(t, result)
		assert.Equal(t, user.UserID, result.User.ID)
		assert.Equal(t, user.Email, result.User.Email)
		assert.NotNil(t, result.User.LastLoginDate)
		assert.Equal(t, oldLoginDate, *result.User.LastLoginDate) // Check that the old login date is used
		mockColl.AssertExpectations(t)
	})
}

func TestCreateUser(t *testing.T) {
	docMatcher := mock.MatchedBy(func(doc interface{}) bool {
		bsonDoc, ok := doc.(bson.D)
		if !ok {
			return false
		}
		// Check if all required fields are present
		requiredFields := []string{
			"user_id",
			"email",
			"user_name",
			"password",
			"first_name",
			"last_name",
			"role",
			"creation_date",
			"last_update_date",
		}
		for _, field := range requiredFields {
			if !containsKey(bsonDoc, field) {
				return false
			}
		}
		return true
	})
	t.Run("successful user creation", func(t *testing.T) {
		mockColl := userMocks.NewMockUserCollection(t)
		ctx := createContextWithMockCollection(mockColl)

		newUser := model.NewUserInput{
			Email:     "test@example.com",
			UserName:  "testuser",
			Password:  "password",
			FirstName: "Test",
			LastName:  "User",
		}

		expectedCountFilter := bson.M{
			"$or": []bson.M{
				{"email": newUser.Email},
				{"user_name": newUser.UserName},
			},
		}
		mockColl.On("CountDocuments", ctx, expectedCountFilter).Return(int64(0), nil)

		mockColl.On("InsertOne", ctx, docMatcher).Return(&mongo.InsertOneResult{}, nil)

		userSvc := &userSvc{}
		result, err := userSvc.CreateUser(ctx, newUser)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, newUser.Email, result.User.Email)
		assert.Equal(t, newUser.UserName, result.User.UserName)
		mockColl.AssertExpectations(t)
	})

	t.Run("user already exists", func(t *testing.T) {
		mockColl := userMocks.NewMockUserCollection(t)
		ctx := createContextWithMockCollection(mockColl)

		newUser := model.NewUserInput{
			Email:    "existing@example.com",
			UserName: "existinguser",
			Password: "password",
		}

		expectedCountFilter := bson.M{
			"$or": []bson.M{
				{"email": newUser.Email},
				{"user_name": newUser.UserName},
			},
		}
		mockColl.On("CountDocuments", ctx, expectedCountFilter).Return(int64(1), nil)

		userSvc := &userSvc{}
		result, err := userSvc.CreateUser(ctx, newUser)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, errUserAlreadyExists, err)
		mockColl.AssertExpectations(t)
	})

	t.Run("database error", func(t *testing.T) {
		mockColl := userMocks.NewMockUserCollection(t)
		ctx := createContextWithMockCollection(mockColl)

		newUser := model.NewUserInput{
			Email:    "test@example.com",
			UserName: "testuser",
			Password: "password",
		}

		expectedCountFilter := bson.M{
			"$or": []bson.M{
				{"email": newUser.Email},
				{"user_name": newUser.UserName},
			},
		}
		mockColl.On("CountDocuments", ctx, expectedCountFilter).Return(int64(0), nil)

		mockColl.On("InsertOne", ctx, docMatcher).Return(nil, errors.New("database error"))

		userSvc := &userSvc{}
		result, err := userSvc.CreateUser(ctx, newUser)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "database error")
		mockColl.AssertExpectations(t)
	})
}

func TestDeleteUser(t *testing.T) {
	t.Run("successful user deletion", func(t *testing.T) {
		mockColl := userMocks.NewMockUserCollection(t)
		ctx := createContextWithMockCollection(mockColl)

		userID := "test-user-id"
		expectedFilter := bson.M{"user_id": userID}
		mockColl.On("DeleteOne", ctx, expectedFilter).Return(&mongo.DeleteResult{DeletedCount: 1}, nil)

		userSvc := &userSvc{}
		success, err := userSvc.DeleteUser(ctx, userID)

		assert.NoError(t, err)
		assert.True(t, success)
		mockColl.AssertExpectations(t)
	})

	t.Run("user not found", func(t *testing.T) {
		mockColl := userMocks.NewMockUserCollection(t)
		ctx := createContextWithMockCollection(mockColl)

		userID := "nonexistent-user-id"
		expectedFilter := bson.M{"user_id": userID}
		mockColl.On("DeleteOne", ctx, expectedFilter).Return(&mongo.DeleteResult{DeletedCount: 0}, nil)

		userSvc := &userSvc{}
		success, err := userSvc.DeleteUser(ctx, userID)

		assert.Error(t, err)
		assert.False(t, success)
		assert.Equal(t, errNoUserFound, err)
		mockColl.AssertExpectations(t)
	})

	t.Run("database error", func(t *testing.T) {
		mockColl := userMocks.NewMockUserCollection(t)
		ctx := createContextWithMockCollection(mockColl)

		userID := "test-user-id"
		expectedFilter := bson.M{"user_id": userID}
		mockColl.On("DeleteOne", ctx, expectedFilter).Return(nil, errors.New("database error"))

		userSvc := &userSvc{}
		success, err := userSvc.DeleteUser(ctx, userID)

		assert.Error(t, err)
		assert.False(t, success)
		assert.Contains(t, err.Error(), "database error")
		mockColl.AssertExpectations(t)
	})
}

// Helper function to create a context with mock collection
func createContextWithMockCollection(collection UserCollection) context.Context {
	ctx := context.Background()
	return NewContext(ctx, GetUsersCollectionKey(), collection)
}

// Helper function to check if a bson.D contains a key
func containsKey(doc bson.D, key string) bool {
	for _, elem := range doc {
		if elem.Key == key {
			return true
		}
	}
	return false
}
