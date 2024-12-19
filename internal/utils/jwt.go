package utils

import (
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

const (
	AccessTokenCookieName  = "X-Access-Token"
	RefreshTokenCookieName = "X-Refresh-Token"
)

var ErrInvalidSignature = fmt.Errorf("token signature is invalid")
var ErrTokenExpired = fmt.Errorf("token is expired")
var ErrTokenInvalid = fmt.Errorf("token is invalid")

func GetAccessTokenCookie(req *http.Request) (string, error) {
	return getCookieByName(req, AccessTokenCookieName)
}

func GetRefreshTokenCookie(req *http.Request) (string, error) {
	return getCookieByName(req, RefreshTokenCookieName)
}

func getCookieByName(req *http.Request, name string) (string, error) {
	cookie, err := req.Cookie(name)
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
