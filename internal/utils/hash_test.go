package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHashPassword(t *testing.T) {
	testCases := []struct {
		name        string
		password    string
		expectEmpty bool
	}{
		{
			name:        "Valid password",
			password:    "mySecureP@ssword",
			expectEmpty: false,
		},
		{
			name:        "Empty password",
			password:    "",
			expectEmpty: false,
		},
		{
			name:        "Long password over bcrypt limit of 72 bytes",
			password:    "testtdsdddddddddddddddddddddddddddddddddddddddddddddddsdddddddddddddddddd",
			expectEmpty: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			hash := HashPassword(tc.password)
			if tc.expectEmpty {
				assert.Empty(t, hash, "Hash should be empty")
			} else {
				assert.NotEmpty(t, hash, "Hash should not be empty")
			}
		})
	}
}

func TestVerifyPassword(t *testing.T) {
	testCases := []struct {
		name        string
		password    string
		hash        string
		expectValid bool
	}{
		{
			name:        "Valid password match",
			password:    "mySecureP@ssword",
			hash:        HashPassword("mySecureP@ssword"),
			expectValid: true,
		},
		{
			name:        "Invalid password",
			password:    "wrongPassword",
			hash:        HashPassword("mySecureP@ssword"),
			expectValid: false,
		},
		{
			name:        "Empty password",
			password:    "",
			hash:        HashPassword(""),
			expectValid: true,
		},
		{
			name:        "Empty hash",
			password:    "testtdsdddddddddddddddddddddddddddddddddddddddddddddddsdddddddddddddddddd",
			hash:        "",
			expectValid: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			valid := VerifyPassword(tc.password, tc.hash)
			assert.Equal(t, tc.expectValid, valid, "Password verification mismatch")
		})
	}
}
