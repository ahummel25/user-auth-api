package mocks

import (
	"errors"

	"github.com/src/user-auth-api/graph/model"
	"github.com/stretchr/testify/mock"
)

// MockedUserService mocks the authentication services.
type MockedUserService struct {
	mock.Mock
	ErrorInvalidPassword bool
}

var errInvalidPassword = errors.New("invalid password")

// AuthenticateUser mocks the user authentication function.
func (s *MockedUserService) AuthenticateUser(username string, password string) (*model.User, error) {
	args := s.Called(username, password)

	if s.ErrorInvalidPassword {
		return nil, errInvalidPassword
	}

	return args.Get(0).(*model.User), nil
}

// CreateUser mocks the user authentication function.
func (s *MockedUserService) CreateUser(params model.CreateUserInput) (*model.User, error) {
	args := s.Called(params.UserName, params.Password)

	if s.ErrorInvalidPassword {
		return nil, errInvalidPassword
	}

	return args.Get(0).(*model.User), nil
}
