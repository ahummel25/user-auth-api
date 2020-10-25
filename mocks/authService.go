package mocks

import (
	"errors"

	"github.com/src/user-auth-api/graph/model"
	"github.com/stretchr/testify/mock"
)

// MockedAuthService mocks the authentication services.
type MockedAuthService struct {
	mock.Mock
	ErrorInvalidPassword bool
}

var errInvalidPassword = errors.New("invalid password")

// AuthenticateUser mocks the user authentication function.
func (s *MockedAuthService) AuthenticateUser(username string, password string) (*model.User, error) {
	args := s.Called(username, password)

	if s.ErrorInvalidPassword {
		return nil, errInvalidPassword
	}

	return args.Get(0).(*model.User), nil
}
