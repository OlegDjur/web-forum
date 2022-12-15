package controller

import (
	"context"
	"forum/internal/models"
	"net/http"
	"time"
)

type ctxKey int8

const (
	ctxKeyUser ctxKey = iota
)

func (h *Handler) authenticateUser(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			user models.User
			err  error
		)

		cookie, err := r.Cookie("sessionID")
		if err != nil {
			http.Redirect(w, r, "/sign-in", http.StatusNotFound)
			return
		}

		user, err = h.services.GetSessionToken(cookie.Value)

		if user.ExpiresAt.Before(time.Now()) {
			next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ctxKeyUser, models.User{})))
			return
		}

		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ctxKeyUser, user)))
	}
}
