package config

import (
	"context"
	"os"
	"sync"
)

// Supplier provides both database credentials and system configuration parameters
type Supplier interface {
	GetConfig() (config, error)
}

// config represents the configuration values from environment variables
type config struct {
	AppName    string
	IAMRoleARN string
	Cluster    string
	Domain     string
	IsDev      bool // Computed at runtime based on AWS_LAMBDA_FUNCTION_NAME
}

// configCtxKey is the context key for the Config value stored in the context
type configCtxKey struct{}

// envConfigSupplier exists as a function receiver to implement the GetConfig interface function
type envConfigSupplier struct {
	// It can literally be an empty struct since it doesn't need to store any state!
}

var (
	cfg  *config
	once sync.Once
)

// GetConfig retrieves configuration from environment variables
func (cs *envConfigSupplier) GetConfig() (config, error) {
	once.Do(func() {
		isDev := os.Getenv("AWS_LAMBDA_FUNCTION_NAME") == ""
		cfg = &config{
			AppName:    os.Getenv("APP_NAME"),
			IAMRoleARN: os.Getenv("IAM_ROLE_ARN"),
			Cluster:    os.Getenv("DB_CLUSTER_NAME"),
			Domain:     os.Getenv("DB_DOMAIN"),
			IsDev:      isDev,
		}
	})
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
