package services

import (
	"errors"

	"github.com/src/user-auth-api/graph/model"
)

// AuthService contains signatures for any auth functions.
type AuthService interface {
	AuthenticateUser(username string, password string) (*model.User, error)
}

type authService struct{}

var errInvalidPassword = errors.New("invalid password")
var errUserDoesNotExist = errors.New("user does not exist")

var mockUserDB = map[string]string{
	"ahummel25": "Welcome123",
}

// NewAuthService returns a pointer to a new UserService.
func NewAuthService() *authService {
	return &authService{}
}

// AuthenticateUser resolves the user.
func (a *authService) AuthenticateUser(username string, password string) (*model.User, error) {
	if mockUserDB[username] == "" {
		return nil, errUserDoesNotExist
	}

	dbUserPassword := mockUserDB[username]

	if dbUserPassword != password {
		return nil, errInvalidPassword
	}

	id := "1"
	name := "Andrew"

	user := &model.User{
		ID:   id,
		Name: name,
	}

	return user, nil
}
