package config

import (
	"context"
	"encoding/json"
	"os"
	"strings"

	"github.com/aws/aws-secretsmanager-caching-go/secretcache"
)

var (
	configCtxKey   string // &configCtxKey is the context key for the Supplier
	secretCache, _ = secretcache.New()
)

// Config represents the key value map in the AWS secret
type Config struct {
	UserName string `json:"DB_USER_NAME"`
	Password string `json:"DB_PASSWORD"`
	Cluster  string `json:"DB_CLUSTER_NAME"`
	Domain   string `json:"DB_DOMAIN"`
}

// new retrieves cached configuration secrets
func new() (Config, error) {
	var cfg Config
	secretString, err := secretCache.GetSecretString(os.Getenv("SECRET_NAME"))
	if err != nil {
		return Config{}, err
	}
	if err = json.NewDecoder(strings.NewReader(secretString)).Decode(&cfg); err != nil {
		return Config{}, err
	}
	return cfg, nil
}

// NewContext returns a new context containing the config
func NewContext(ctx context.Context, c Config) context.Context {
	return context.WithValue(ctx, &configCtxKey, c)
}

// FromContext returns the interface that was stored in the context, or a new one if none was stored
func FromContext(ctx context.Context) (Config, error) {
	if i, ok := ctx.Value(&configCtxKey).(Config); ok {
		return i, nil
	}
	return new()
}
