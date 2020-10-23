package main

import (
	"net/http"

	"github.com/src/user-auth-api/config"
	"github.com/src/user-auth-api/db"
	"github.com/src/user-auth-api/service"
	"github.com/src/user-auth-api/service/user"
)

// injectDBCollection is a middleware that injects the Users collection into the request context
func injectDBCollection(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		cfg, err := config.FromContext(ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		ctx = config.NewContext(ctx, cfg)
		usersCollection, err := db.GetCollection(ctx, db.USERS_COLLECTION)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		reqWithCtx := r.WithContext(service.NewContext(ctx, user.GetUsersCollectionKey(), usersCollection))
		next.ServeHTTP(w, reqWithCtx)
	})
}
