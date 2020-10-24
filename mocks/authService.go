package mocks

import (
	"github.com/src/user-auth-api/graph/model"
	"github.com/stretchr/testify/mock"
)

// MockedAuthService mocks the authentication services.
type MockedAuthService struct {
	mock.Mock
}

// AuthenticateUser mocks the user authentication function.
func (s *MockedAuthService) AuthenticateUser(username string, password string) (*model.User, error) {
	args := s.Called(username, password)
	return args.Get(0).(*model.User), nil
}
