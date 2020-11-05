package mocks

import (
	"errors"

	"github.com/src/user-auth-api/graph/model"
	"github.com/stretchr/testify/mock"
)

// MockedUserService mocks the authentication services.
type MockedUserService struct {
	mock.Mock
	ErrorInvalidPassword   bool
	ErrorUserAlreadyExists bool
}

var (
	errInvalidPassword       = errors.New("invalid password")
	errNoUserFound           = errors.New("no user found by that email address")
	errUserNameAlreadyExists = errors.New("user name already exists")
)

// AuthenticateUser mocks the user authentication function.
func (s *MockedUserService) AuthenticateUser(username string, password string) (*model.UserObject, error) {
	args := s.Called(username, password)

	if s.ErrorInvalidPassword {
		return nil, errInvalidPassword
	}

	return args.Get(0).(*model.UserObject), nil
}

// CreateUser mocks the user authentication function.
func (s *MockedUserService) CreateUser(params model.CreateUserInput) (*model.UserObject, error) {
	args := s.Called(params)

	if s.ErrorUserAlreadyExists {
		return nil, errUserNameAlreadyExists
	}

	return args.Get(0).(*model.UserObject), nil
}
