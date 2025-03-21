package middleware

import (
	"context"
	"net/http"

	"version-backend/internal/db"
)

// DBKey is the context key for the database connection
type DBKey struct{}

// WithDB middleware injects the database connection into the request context
func WithDB(db *db.DB) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), DBKey{}, db)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
