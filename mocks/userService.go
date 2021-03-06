package mocks

import (
	"errors"

	"github.com/src/user-auth-api/graphql/model"
	"github.com/stretchr/testify/mock"
)

// MockedUserService mocks the authentication services.
type MockedUserService struct {
	mock.Mock
	ErrorInvalidPassword   bool
	ErrorNoUserFound       bool
	ErrorUserAlreadyExists bool
}

var (
	errInvalidPassword       = errors.New("invalid password")
	errNoUserFound           = errors.New("no user found!")
	errUserNameAlreadyExists = errors.New("user name already exists")
)

// AuthenticateUser mocks the user authentication function.
func (s *MockedUserService) AuthenticateUser(username string, password string) (*model.UserObject, error) {
	args := s.Called(username, password)

	if s.ErrorInvalidPassword {
		return nil, errInvalidPassword
	}

	if s.ErrorNoUserFound {
		return nil, errNoUserFound
	}

	return args.Get(0).(*model.UserObject), nil
}

// CreateUser mocks the create user authentication.
func (s *MockedUserService) CreateUser(params model.CreateUserInput) (*model.UserObject, error) {
	args := s.Called(params)

	if s.ErrorUserAlreadyExists {
		return nil, errUserNameAlreadyExists
	}

	return args.Get(0).(*model.UserObject), nil
}

// DeleteUser mocks the delete user function.
func (s *MockedUserService) DeleteUser(params model.DeleteUserInput) (string, error) {
	_ = s.Called(params)

	if s.ErrorNoUserFound {
		return "", errNoUserFound
	}

	return params.UserName + " successfully deleted", nil
}
