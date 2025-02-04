package config

import (
	"context"
	"encoding/json"
	"os"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

var (
	// Store original env vars
	originalEnvVars = map[string]string{
		"AWS_LAMBDA_FUNCTION_NAME": "",
		"APP_NAME":                 "",
		"IAM_ROLE_ARN":             "",
		"DB_CLUSTER_NAME":          "",
		"DB_DOMAIN":                "",
		"SECRET_NAME":              "",
	}
	originalSecretCacheInstance SecretCache
)

type ConfigTestSuite struct {
	suite.Suite
}

func TestMain(m *testing.M) {
	// Save all original env vars
	for k := range originalEnvVars {
		originalEnvVars[k] = os.Getenv(k)
	}
	// Save original secret cache
	originalSecretCacheInstance = secretCache

	// Run tests
	code := m.Run()

	// Cleanup
	for k, v := range originalEnvVars {
		_ = os.Setenv(k, v)
	}
	secretCache = originalSecretCacheInstance

	os.Exit(code)
}

func (suite *ConfigTestSuite) SetupTest() {
	// Reset env vars before each test
	for k := range originalEnvVars {
		_ = os.Unsetenv(k)
	}
	// Reset singleton for each test
	cfg = nil
	once = sync.Once{}
}

func (suite *ConfigTestSuite) TestGetConfig_Development() {
	// Set test env vars
	_ = os.Setenv("APP_NAME", "test-app")
	_ = os.Setenv("IAM_ROLE_ARN", "test-role")
	_ = os.Setenv("DB_CLUSTER_NAME", "test-cluster")
	_ = os.Setenv("DB_DOMAIN", "test-domain")

	supplier := &envConfigSupplier{}
	config, err := supplier.GetConfig()

	suite.Require().NoError(err)
	suite.Assert().True(config.IsDev)
	suite.Assert().Equal("test-app", config.AppName)
	suite.Assert().Equal("test-role", config.IAMRoleARN)
	suite.Assert().Equal("test-cluster", config.Cluster)
	suite.Assert().Equal("test-domain", config.Domain)
}

func (suite *ConfigTestSuite) TestGetConfig_Production() {
	// Set test env vars for production mode
	_ = os.Setenv("AWS_LAMBDA_FUNCTION_NAME", "test-lambda")
	_ = os.Setenv("SECRET_NAME", "test-secret")

	// Set up mock secret cache
	mockSecretString := `{
		"APP_NAME": "prod-app",
		"IAM_ROLE_ARN": "prod-role",
		"DB_CLUSTER_NAME": "prod-cluster",
		"DB_DOMAIN": "prod-domain"
	}`
	secretCache = &mockSecretCache{secretString: mockSecretString}

	supplier := &envConfigSupplier{}
	config, err := supplier.GetConfig()

	suite.Require().NoError(err)
	suite.Assert().False(config.IsDev)
	suite.Assert().Equal("prod-app", config.AppName)
	suite.Assert().Equal("prod-role", config.IAMRoleARN)
	suite.Assert().Equal("prod-cluster", config.Cluster)
	suite.Assert().Equal("prod-domain", config.Domain)
}

func (suite *ConfigTestSuite) TestGetConfig_ProductionError() {
	// Set production mode
	_ = os.Setenv("AWS_LAMBDA_FUNCTION_NAME", "test-lambda")

	// Set up mock secret cache with error
	secretCache = &mockSecretCache{err: assert.AnError}

	supplier := &envConfigSupplier{}
	_, err := supplier.GetConfig()

	suite.Assert().Error(err)
	suite.Assert().ErrorAs(err, &assert.AnError)
}

func (suite *ConfigTestSuite) TestGetConfig_DecodeError() {
	// Set test env vars for production mode
	_ = os.Setenv("AWS_LAMBDA_FUNCTION_NAME", "test-lambda")
	_ = os.Setenv("SECRET_NAME", "test-secret")

	// Set up mock secret cache with invalid JSON
	invalidJSON := `{
		"APP_NAME": "prod-app",
		"IAM_ROLE_ARN": "prod-role", INVALID JSON HERE
		"DB_CLUSTER_NAME": "prod-cluster",
		"DB_DOMAIN": "prod-domain"
	}`
	secretCache = &mockSecretCache{secretString: invalidJSON}

	supplier := &envConfigSupplier{}
	_, err := supplier.GetConfig()

	suite.Assert().Error(err)
	var syntaxErr *json.SyntaxError
	suite.Assert().ErrorAs(err, &syntaxErr)
}

func (suite *ConfigTestSuite) TestContext() {
	ctx := context.Background()
	mockSupplier := &mockSupplier{
		config: config{
			AppName: "test-app",
			IsDev:   true,
		},
	}

	// Test NewContext
	ctxWithConfig := NewContext(ctx, mockSupplier)
	suite.Assert().NotNil(ctxWithConfig)

	// Test FromContext with config
	supplier, err := FromContext(ctxWithConfig)
	suite.Require().NoError(err)
	config, err := supplier.GetConfig()
	suite.Require().NoError(err)
	suite.Assert().Equal("test-app", config.AppName)
	suite.Assert().True(config.IsDev)

	// Test FromContext without config (should return default envConfigSupplier)
	supplier, err = FromContext(context.Background())
	suite.Require().NoError(err)
	suite.Assert().IsType(&envConfigSupplier{}, supplier)
}

// Mock implementations for testing
type mockSecretCache struct {
	secretString string
	err          error
}

func (m *mockSecretCache) GetSecretString(secretID string) (string, error) {
	if m.err != nil {
		return "", m.err
	}
	return m.secretString, nil
}

type mockSupplier struct {
	config config
	err    error
}

func (m *mockSupplier) GetConfig() (config, error) {
	if m.err != nil {
		return config{}, m.err
	}
	return m.config, nil
}

func TestConfigSuite(t *testing.T) {
	suite.Run(t, new(ConfigTestSuite))
}
