package utils

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
)

func TestCreateToken(t *testing.T) {
	secret := "testsecret"
	userId := "12345"
	exp := time.Hour

	token, err := CreateToken(secret, userId, exp)
	assert.Nil(t, err, "Error should be nil")
	assert.NotEmpty(t, token, "Token should not be empty")

	parsedToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	assert.Nil(t, err, "Error should be nil when parsing token")
	assert.NotNil(t, parsedToken, "Parsed token should not be nil")

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	assert.True(t, ok, "Claims should be of type jwt.MapClaims")
	assert.Equal(t, "GopherTalk", claims["iss"], "Issuer should match")
	assert.Equal(t, userId, claims["sub"], "Subject should match")
}

func TestGetToken(t *testing.T) {
	secret := "testsecret"
	userId := "12345"
	exp := time.Second * 1

	token, err := CreateToken(secret, userId, exp)
	assert.Nil(t, err, "Error should be nil")
	assert.NotEmpty(t, token, "Token should not be empty")

	claims, err := GetToken(token, secret)
	assert.Nil(t, err, "Error should be nil when validating token")
	assert.NotNil(t, claims, "Claims should not be nil")
	assert.Equal(t, userId, claims.Subject, "Subject should match")

	time.Sleep(2 * exp) // Simulate token expiration

	claims, err = GetToken(token, secret)
	assert.NotNil(t, err, "Error should not be nil for expired token")
	assert.Equal(t, ErrTokenExpired, err, "Error should indicate token expiration")
	assert.Nil(t, claims, "Claims should be nil for expired token")
}
