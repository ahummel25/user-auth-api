package config

import (
	"context"
	"encoding/json"
	"os"
	"strings"

	"github.com/aws/aws-secretsmanager-caching-go/secretcache"
)

// Config represents the key value map in the AWS secret
type config struct {
	UserName string `json:"DB_USER_NAME"`
	Password string `json:"DB_PASSWORD"`
	Cluster  string `json:"DB_CLUSTER_NAME"`
	Domain   string `json:"DB_DOMAIN"`
}

// configCtxKey is the context key for the Config value stored in the context
type configCtxKey struct{}

var (
	cfg            *config
	secretCache, _ = secretcache.New()
)

// new retrieves cached configuration secrets
func new() (config, error) {
	// Return the already populated config if we have it, no need to refetch here
	if cfg != nil {
		return *cfg, nil
	}
	secretString, err := secretCache.GetSecretString(os.Getenv("SECRET_NAME"))
	if err != nil {
		return config{}, err
	}
	if err = json.NewDecoder(strings.NewReader(secretString)).Decode(&cfg); err != nil {
		return config{}, err
	}
	return *cfg, nil
}

// NewContext returns a new context containing the config
func NewContext(ctx context.Context, c config) context.Context {
	return context.WithValue(ctx, configCtxKey{}, c)
}

// FromContext returns the config that was stored in the context, or a new one if none was stored
func FromContext(ctx context.Context) (config, error) {
	if i, ok := ctx.Value(configCtxKey{}).(config); ok {
		return i, nil
	}
	return new()
}
