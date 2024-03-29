package service

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

// NewContext returns a new Context with a Mongo collection added
func NewContext(ctx context.Context, collectionCtxKey any, collection *mongo.Collection) context.Context {
	return context.WithValue(ctx, collectionCtxKey, collection)
}

// FromContext returns the Mongo collection that was stored in the context, or nil if none was stored
func FromContext(ctx context.Context, collectionCtxKey any) *mongo.Collection {
	if s, ok := ctx.Value(collectionCtxKey).(*mongo.Collection); ok {
		return s
	}
	return nil
}
