package resolvers_test

import (
	"context"
	"errors"
	"os"
	"testing"
	"time"

	"github.com/99designs/gqlgen/client"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/ahummel25/user-auth-api/graphql/directives"
	"github.com/ahummel25/user-auth-api/graphql/generated"
	"github.com/ahummel25/user-auth-api/graphql/model"
	"github.com/ahummel25/user-auth-api/graphql/resolvers"
	"github.com/ahummel25/user-auth-api/graphql/resolvers/mutations"
	userMutation "github.com/ahummel25/user-auth-api/graphql/resolvers/mutations/user"
	"github.com/ahummel25/user-auth-api/graphql/resolvers/query"
	userQuery "github.com/ahummel25/user-auth-api/graphql/resolvers/query/user"
	userMocks "github.com/ahummel25/user-auth-api/service/user/mocks"
	"github.com/ahummel25/user-auth-api/testutils"
)

var (
	ctxMatcher = mock.MatchedBy(func(ctx context.Context) bool {
		// This matcher will accept any context
		return true
	})
	errInvalidPassword   = errors.New("invalid password")
	errNoUserFound       = errors.New("no user found")
	errUserAlreadyExists = errors.New("user name or email already exists")
	mockUserID           = "dfb8fe7f-56e4-47dc-b5bc-f6f0f524402b"
	mockFirstName        = "Test"
	mockLastName         = "User"
	mockEmail            = "mock_email@gmail.com"
	mockInvalidEmail     = "mock_email@gmail"
	mockPassword         = "mockPassword123"
	mockUserName         = "mock_username"
	mockRole             = model.RoleUser
)

var (
	loginQuery = `query Login($usernameOrEmail: String!, $password: String!) { 
	  auth: login(params: {usernameOrEmail: $usernameOrEmail, password: $password}) {
		user {
		  id
		  firstName
		  lastName
		  email
		  userName
		  role
		  lastLoginDate
		}
	  }
	}`

	createUser = `mutation CreateUser($newUserInput: NewUserInput!) {
	  createUser(user: $newUserInput) {
		user {
		  id
		  email
		  firstName
		  lastName
		  userName
		  role
		}
	  }
	}`

	deleteUser = `mutation DeleteUser($userID: ID!) {
		deleteUser(userID: $userID)
	}`
)

func setup(t *testing.T) (*client.Client, *userMocks.MockAPI) {
	mockUserService := userMocks.NewMockAPI(t)
	mutationResolvers := mutations.MutationResolvers{
		Resolver: userMutation.Resolver{UserService: mockUserService},
	}
	queryResolvers := query.QueryResolvers{
		Resolver: userQuery.Resolver{UserService: mockUserService},
	}
	mockResolvers := resolvers.Services{
		MutationResolvers: mutationResolvers,
		QueryResolvers:    queryResolvers,
	}
	cfg := generated.Config{Resolvers: &mockResolvers}
	cfg.Directives.Binding = directives.Binding
	cfg.Directives.HasRole = directives.HasRole
	c := client.New(handler.NewDefaultServer(generated.NewExecutableSchema(cfg)))
	return c, mockUserService
}

func assertUserEqual(t *testing.T, expected, actual model.User) {
	assert.Equal(t, expected.ID, actual.ID)
	assert.Equal(t, expected.FirstName, actual.FirstName)
	assert.Equal(t, expected.LastName, actual.LastName)
	assert.Equal(t, expected.Email, actual.Email)
	assert.Equal(t, expected.UserName, actual.UserName)
	assert.Equal(t, expected.Role, actual.Role)
	if expected.LastLoginDate != nil {
		assert.NotNil(t, actual.LastLoginDate)
		assert.Equal(t, expected.LastLoginDate, actual.LastLoginDate)
	} else {
		assert.Nil(t, actual.LastLoginDate)
	}
}

func createMockUser() *model.User {
	now := testutils.CurrentTime.Now()
	return &model.User{
		ID:            mockUserID,
		FirstName:     mockFirstName,
		LastName:      mockLastName,
		Email:         mockEmail,
		UserName:      mockUserName,
		Role:          mockRole,
		LastLoginDate: &now,
	}
}

func createMockUserWithoutLastLoginDate() *model.User {
	return &model.User{
		ID:        mockUserID,
		FirstName: mockFirstName,
		LastName:  mockLastName,
		Email:     mockEmail,
		UserName:  mockUserName,
		Role:      mockRole,
	}
}

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

func Test_Login(t *testing.T) {
	tests := []struct {
		name                  string
		setupMock             func(*userMocks.MockAPI)
		expectedError         string
		expectedLastLoginDate *time.Time
	}{
		{
			name: "Success",
			setupMock: func(mockService *userMocks.MockAPI) {
				mockService.On("Login", ctxMatcher, mockUserName, mockPassword).
					Return(&model.UserObject{User: createMockUser()}, nil)
			},
			expectedLastLoginDate: testutils.TimePtr(testutils.CurrentTime.Now()),
		},
		{
			name: "Success with update failure",
			setupMock: func(mockService *userMocks.MockAPI) {
				user := createMockUser()
				oldTime := testutils.CurrentTime.Now().Add(-24 * time.Hour)
				user.LastLoginDate = &oldTime
				mockService.On("Login", ctxMatcher, mockUserName, mockPassword).
					Return(&model.UserObject{User: user}, nil)
			},
			expectedLastLoginDate: testutils.TimePtr(testutils.CurrentTime.Now().Add(-24 * time.Hour)),
		},
		{
			name: "User not found",
			setupMock: func(mockService *userMocks.MockAPI) {
				mockService.On("Login", ctxMatcher, mockUserName, mockPassword).
					Return(nil, errNoUserFound)
			},
			expectedError: errNoUserFound.Error(),
		},
		{
			name: "Invalid password",
			setupMock: func(mockService *userMocks.MockAPI) {
				mockService.On("Login", ctxMatcher, mockUserName, mockPassword).
					Return(nil, errInvalidPassword)
			},
			expectedError: errInvalidPassword.Error(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, mockUserService := setup(t)
			tt.setupMock(mockUserService)

			var response struct {
				Auth struct {
					User struct {
						ID            string
						FirstName     string
						LastName      string
						Email         string
						UserName      string
						Role          model.Role
						LastLoginDate *string
					}
				}
			}
			err := c.Post(loginQuery, &response,
				client.Var("usernameOrEmail", mockUserName),
				client.Var("password", mockPassword),
			)

			mockUserService.AssertExpectations(t)
			if tt.expectedError != "" {
				require.Error(t, err)
				require.EqualError(t, err, `[{"message":"`+tt.expectedError+`","path":["auth"]}]`)
				assert.Empty(t, response)
			} else {
				require.NoError(t, err)
				expectedUser := createMockUser()
				if tt.expectedLastLoginDate != nil {
					expectedUser.LastLoginDate = tt.expectedLastLoginDate
				}
				var lastLoginTime *time.Time
				if response.Auth.User.LastLoginDate != nil {
					parsedTime, err := time.Parse(time.RFC3339, *response.Auth.User.LastLoginDate)
					require.NoError(t, err)
					lastLoginTime = &parsedTime
				}
				actualUser := model.User{
					ID:            response.Auth.User.ID,
					FirstName:     response.Auth.User.FirstName,
					LastName:      response.Auth.User.LastName,
					Email:         response.Auth.User.Email,
					UserName:      response.Auth.User.UserName,
					Role:          response.Auth.User.Role,
					LastLoginDate: lastLoginTime,
				}
				assertUserEqual(t, *expectedUser, actualUser)
			}
		})
	}
}

func Test_CreateUser(t *testing.T) {
	tests := []struct {
		name              string
		input             model.NewUserInput
		setupMock         func(*userMocks.MockAPI, model.NewUserInput)
		expectedError     string
		expectedErrorPath string
	}{
		{
			name: "Success",
			input: model.NewUserInput{
				Email: mockEmail, FirstName: mockFirstName, LastName: mockLastName,
				UserName: mockUserName, Password: mockPassword,
			},
			setupMock: func(mockService *userMocks.MockAPI, input model.NewUserInput) {
				mockService.On("CreateUser", ctxMatcher, input).
					Return(&model.UserObject{User: createMockUserWithoutLastLoginDate()}, nil)
			},
		},
		{
			name: "User already exists",
			input: model.NewUserInput{
				Email: mockEmail, FirstName: mockFirstName, LastName: mockLastName,
				UserName: mockUserName, Password: mockPassword,
			},
			setupMock: func(mockService *userMocks.MockAPI, input model.NewUserInput) {
				mockService.On("CreateUser", ctxMatcher, input).
					Return(nil, errUserAlreadyExists)
			},
			expectedError:     errUserAlreadyExists.Error(),
			expectedErrorPath: `["createUser"]`,
		},
		{
			name: "Invalid password",
			input: model.NewUserInput{
				Email: mockEmail, FirstName: mockFirstName, LastName: mockLastName,
				UserName: mockUserName, Password: "123ABC",
			},
			expectedError:     "password must be at least 8 characters in length",
			expectedErrorPath: `["createUser","user","password"]`,
		},
		{
			name: "Invalid email",
			input: model.NewUserInput{
				Email: mockInvalidEmail, FirstName: mockFirstName, LastName: mockLastName,
				UserName: mockUserName, Password: mockPassword,
			},
			expectedError:     "email must be a valid email address",
			expectedErrorPath: `["createUser","user","email"]`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, mockUserService := setup(t)
			if tt.setupMock != nil {
				tt.setupMock(mockUserService, tt.input)
			}

			var response struct{ CreateUser struct{ model.User } }
			err := c.Post(createUser, &response, client.Var("newUserInput", tt.input))

			if tt.expectedError != "" {
				require.Error(t, err)
				require.EqualError(t, err, `[{"message":"`+tt.expectedError+`","path":`+tt.expectedErrorPath+`}]`)
				assert.Empty(t, response)
			} else {
				require.NoError(t, err)
				assertUserEqual(t, *createMockUserWithoutLastLoginDate(), response.CreateUser.User)
			}
			mockUserService.AssertExpectations(t)
		})
	}
}

func Test_DeleteUser(t *testing.T) {
	tests := []struct {
		name              string
		setupMock         func(*userMocks.MockAPI)
		expectedError     string
		expectedErrorPath string
	}{
		{
			name: "Success",
			setupMock: func(mockService *userMocks.MockAPI) {
				mockService.On("DeleteUser", ctxMatcher, mockUserID).Return(true, nil)
			},
		},
		{
			name: "User not found",
			setupMock: func(mockService *userMocks.MockAPI) {
				mockService.On("DeleteUser", ctxMatcher, mockUserID).Return(false, errNoUserFound)
			},
			expectedError:     errNoUserFound.Error(),
			expectedErrorPath: `["deleteUser"]`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, mockUserService := setup(t)
			tt.setupMock(mockUserService)

			var response struct{ DeleteUser bool }
			err := c.Post(deleteUser, &response, client.Var("userID", mockUserID))

			if tt.expectedError != "" {
				require.Error(t, err)
				require.EqualError(t, err, `[{"message":"`+tt.expectedError+`","path":`+tt.expectedErrorPath+`}]`)
				assert.False(t, response.DeleteUser)
			} else {
				require.NoError(t, err)
				assert.True(t, response.DeleteUser)
			}
			mockUserService.AssertExpectations(t)
		})
	}
}
