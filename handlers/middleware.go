package handlers

import (
	internalCtx "UserService/internal/context"
	"UserService/internal/models"
	"UserService/pkg/jwt"
	"context"
	"net/http"
	"strings"
)

type Middleware struct {
	repo IMiddleware
}

type IMiddleware interface {
	Get(ctx context.Context, userID int) (models.User, error)
}

func NewMiddleware(repo IMiddleware) *Middleware {
	return &Middleware{
		repo: repo,
	}
}

func (m *Middleware) AuthAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")

		if auth == "" {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		if !strings.HasPrefix(auth, "Bearer ") {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		token := strings.TrimPrefix(auth, "Bearer ")

		claims, err := jwt.ParseToken(token)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		user, err := m.repo.Get(r.Context(), claims.UserID)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		if user.Role != "admin" {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), internalCtx.UserIDKey, claims.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")

		if auth == "" {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		if !strings.HasPrefix(auth, "Bearer ") {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		token := strings.TrimPrefix(auth, "Bearer ")

		claims, err := jwt.ParseToken(token)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), internalCtx.RoleKey, claims.Role)
		ctx = context.WithValue(ctx, internalCtx.UserIDKey, claims.UserID)
		ctx = context.WithValue(ctx, internalCtx.EmailKey, claims.Email)

		next.ServeHTTP(w, r.WithContext(ctx))

	})
}
