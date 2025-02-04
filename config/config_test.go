package config

import (
	"context"
	"encoding/json"
	"os"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
		os.Setenv(k, v)
	}
	secretCache = originalSecretCacheInstance

	os.Exit(code)
}

func setupTest() {
	// Reset env vars before each test
	for k := range originalEnvVars {
		os.Unsetenv(k)
	}
	// Reset singleton for each test
	cfg = nil
	once = sync.Once{}
}

func TestGetConfig_Development(t *testing.T) {
	setupTest()

	// Set test env vars
	os.Setenv("APP_NAME", "test-app")
	os.Setenv("IAM_ROLE_ARN", "test-role")
	os.Setenv("DB_CLUSTER_NAME", "test-cluster")
	os.Setenv("DB_DOMAIN", "test-domain")

	supplier := &envConfigSupplier{}
	config, err := supplier.GetConfig()

	require.NoError(t, err)
	assert.True(t, config.IsDev)
	assert.Equal(t, "test-app", config.AppName)
	assert.Equal(t, "test-role", config.IAMRoleARN)
	assert.Equal(t, "test-cluster", config.Cluster)
	assert.Equal(t, "test-domain", config.Domain)
}

func TestGetConfig_Production(t *testing.T) {
	setupTest()

	// Set test env vars for production mode
	os.Setenv("AWS_LAMBDA_FUNCTION_NAME", "test-lambda")
	os.Setenv("SECRET_NAME", "test-secret")

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

	require.NoError(t, err)
	assert.False(t, config.IsDev)
	assert.Equal(t, "prod-app", config.AppName)
	assert.Equal(t, "prod-role", config.IAMRoleARN)
	assert.Equal(t, "prod-cluster", config.Cluster)
	assert.Equal(t, "prod-domain", config.Domain)
}

func TestGetConfig_ProductionError(t *testing.T) {
	setupTest()

	// Set production mode
	os.Setenv("AWS_LAMBDA_FUNCTION_NAME", "test-lambda")

	// Set up mock secret cache with error
	secretCache = &mockSecretCache{err: assert.AnError}

	supplier := &envConfigSupplier{}
	_, err := supplier.GetConfig()

	assert.Error(t, err)
}

func TestGetConfig_DecodeError(t *testing.T) {
	setupTest()

	// Set test env vars for production mode
	os.Setenv("AWS_LAMBDA_FUNCTION_NAME", "test-lambda")
	os.Setenv("SECRET_NAME", "test-secret")

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

	assert.Error(t, err)
	var syntaxErr *json.SyntaxError
	assert.ErrorAs(t, err, &syntaxErr)
}

func TestContext(t *testing.T) {
	setupTest()

	ctx := context.Background()
	mockSupplier := &mockSupplier{
		config: config{
			AppName: "test-app",
			IsDev:   true,
		},
	}

	// Test NewContext
	ctxWithConfig := NewContext(ctx, mockSupplier)
	assert.NotNil(t, ctxWithConfig)

	// Test FromContext with config
	supplier, err := FromContext(ctxWithConfig)
	require.NoError(t, err)
	config, err := supplier.GetConfig()
	require.NoError(t, err)
	assert.Equal(t, "test-app", config.AppName)
	assert.True(t, config.IsDev)

	// Test FromContext without config (should return default envConfigSupplier)
	supplier, err = FromContext(context.Background())
	require.NoError(t, err)
	assert.IsType(t, &envConfigSupplier{}, supplier)
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

func (m *mockSecretCache) GetSecretStringWithStage(secretID string, stage string) (string, error) {
	return m.GetSecretString(secretID)
}

func (m *mockSecretCache) GetSecretBinary(secretID string) ([]byte, error) {
	if m.err != nil {
		return nil, m.err
	}
	return []byte(m.secretString), nil
}

func (m *mockSecretCache) GetSecretBinaryWithStage(secretID string, stage string) ([]byte, error) {
	return m.GetSecretBinary(secretID)
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
