package middleware

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/shekshuev/gophertalk-backend/internal/utils"
)

func RequestAuth(secret string) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tokenString, err := utils.GetRawAccessToken(r)
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			_, err = utils.GetToken(tokenString, secret)
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			h.ServeHTTP(w, r)
		})
	}
}

func RequestAuthSameID(secret string) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tokenString, err := utils.GetRawAccessToken(r)
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			claims, err := utils.GetToken(tokenString, secret)
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			strId := chi.URLParam(r, "id")
			_, err = strconv.Atoi(strId)
			if err != nil {
				// if not number - ignore, handler will return 404
				h.ServeHTTP(w, r)
				return
			}
			if strId != claims.Subject {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			h.ServeHTTP(w, r)
		})
	}
}
