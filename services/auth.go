package services

import (
	"errors"

	"github.com/src/user-auth-api/graph/model"
)

// AuthService contains signatures for any auth functions.
type AuthService interface {
	AuthenticateUser(email string, password string) (*model.User, error)
}

type authService struct{}

var (
	errInvalidPassword  = errors.New("invalid password")
	errUserDoesNotExist = errors.New("user does not exist")
)

var mockUserDB = map[string]string{
	"ahummel25@gmail.com": "Welcome123",
}

// NewAuthService returns a pointer to a new auth service.
func NewAuthService() *authService {
	return &authService{}
}

// AuthenticateUser authenticates the user.
func (a *authService) AuthenticateUser(email string, password string) (*model.User, error) {
	if mockUserDB[email] == "" {
		return nil, errUserDoesNotExist
	}

	dbUserPassword := mockUserDB[email]

	if dbUserPassword != password {
		return nil, errInvalidPassword
	}

	userID := "1"
	firstName := "Andrew"

	user := &model.User{
		UserID:    userID,
		FirstName: firstName,
	}

	return user, nil
}
