package middleware

import (
	"net/http"

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
