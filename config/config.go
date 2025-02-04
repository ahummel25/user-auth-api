package config

import (
	"context"
	"encoding/json"
	"os"
	"strings"
	"sync"

	"github.com/aws/aws-secretsmanager-caching-go/v2/secretcache"
)

// Supplier provides both database credentials and system configuration parameters
type Supplier interface {
	GetConfig() (config, error)
}

// config represents the key value map in the AWS secret
type config struct {
	AppName    string `json:"APP_NAME"`
	IAMRoleARN string `json:"IAM_ROLE_ARN"`
	Cluster    string `json:"DB_CLUSTER_NAME"`
	Domain     string `json:"DB_DOMAIN"`
}

// configCtxKey is the context key for the Config value stored in the context
type configCtxKey struct{}

// envConfigSupplier exists as a function receiver to implement the GetConfig interface function
type envConfigSupplier struct {
	// It can literally be an empty struct since it doesn't need to store any state!
}

var (
	cfg            *config
	secretCache, _ = secretcache.New()
	once           sync.Once
)

// GetConfig retrieves cached configuration secrets
func (cs *envConfigSupplier) GetConfig() (config, error) {
	var err error
	once.Do(func() {
		// For local development, read directly from environment variables
		if os.Getenv("AWS_LAMBDA_FUNCTION_NAME") == "" {
			cfg = &config{
				AppName:    os.Getenv("APP_NAME"),
				IAMRoleARN: os.Getenv("IAM_ROLE_ARN"),
				Cluster:    os.Getenv("DB_CLUSTER_NAME"),
				Domain:     os.Getenv("DB_DOMAIN"),
			}
			return
		}

		// For production, use AWS Secrets Manager
		secretString, secretErr := secretCache.GetSecretString(os.Getenv("SECRET_NAME"))
		if secretErr != nil {
			err = secretErr
			return
		}
		cfg = &config{}
		if decodeErr := json.NewDecoder(strings.NewReader(secretString)).Decode(cfg); decodeErr != nil {
			err = decodeErr
			cfg = nil
		}
	})

	if err != nil {
		return config{}, err
	}
	return *cfg, nil
}

// NewContext returns a new context containing the config
func NewContext(ctx context.Context, s Supplier) context.Context {
	return context.WithValue(ctx, configCtxKey{}, s)
}

// FromContext returns the config that was stored in the context, or a new one if none was stored
func FromContext(ctx context.Context) (Supplier, error) {
	if s, ok := ctx.Value(configCtxKey{}).(Supplier); ok {
		return s, nil
	}
	return &envConfigSupplier{}, nil
}
