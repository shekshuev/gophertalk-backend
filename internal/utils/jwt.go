package utils

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

type ContextKey string

const (
	AccessTokenCookieName  = "X-Access-Token"
	RefreshTokenCookieName = "X-Refresh-Token"
	ContextClaimsKey       = ContextKey("user-claims")
)

var ErrInvalidSignature = fmt.Errorf("token signature is invalid")
var ErrTokenExpired = fmt.Errorf("token is expired")
var ErrTokenInvalid = fmt.Errorf("token is invalid")

func GetRawAccessToken(req *http.Request) (string, error) {
	authHeader := req.Header.Get("Authorization")
	if len(authHeader) == 0 || !strings.HasPrefix(authHeader, "Bearer ") {
		return "", ErrTokenInvalid
	}
	return strings.TrimPrefix(authHeader, "Bearer "), nil
}

func GetRawRefreshToken(req *http.Request) (string, error) {
	cookie, err := req.Cookie(RefreshTokenCookieName)
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}

func CreateToken(secret, userId string, exp time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "GopherTalk",
		Subject:   userId,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(exp)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ID:        uuid.New().String(),
	})
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func GetToken(tokenString, secret string) (*jwt.RegisteredClaims, error) {
	claims := &jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		if validationErr, ok := err.(*jwt.ValidationError); ok {
			if validationErr.Errors&jwt.ValidationErrorExpired != 0 {
				return nil, ErrTokenExpired
			}
			if validationErr.Errors&jwt.ValidationErrorSignatureInvalid != 0 {
				return nil, ErrInvalidSignature
			}
		}
		return nil, ErrTokenInvalid
	}

	if !token.Valid {
		return nil, ErrInvalidSignature
	}
	return claims, nil
}

func GetClaimsFromContext(ctx context.Context) (jwt.RegisteredClaims, bool) {
	claims, ok := ctx.Value(ContextClaimsKey).(jwt.RegisteredClaims)
	return claims, ok
}

func PutClaimsToContext(ctx context.Context, claims jwt.RegisteredClaims) context.Context {
	return context.WithValue(ctx, ContextClaimsKey, claims)
}
