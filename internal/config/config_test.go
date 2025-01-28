package config

import (
	"flag"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetConfig_WithEnv(t *testing.T) {
	serverAddress := "localhost:3000"
	databaseDSN := "host=test port=5432 user=test password=test dbname=test sslmode=disable"
	accessTokenSecret := "test"
	refreshTokenSecret := "test"
	accessTokenExpires := time.Hour
	refreshTokenExpires := time.Hour * 24
	os.Setenv("SERVER_ADDRESS", serverAddress)
	os.Setenv("DATABASE_DSN", databaseDSN)
	os.Setenv("ACCESS_TOKEN_SECRET", accessTokenSecret)
	os.Setenv("REFRESH_TOKEN_SECRET", refreshTokenSecret)
	os.Setenv("ACCESS_TOKEN_EXPIRES", "1h")
	os.Setenv("REFRESH_TOKEN_EXPIRES", "24h")
	defer os.Unsetenv("SERVER_ADDRESS")
	defer os.Unsetenv("DATABASE_DSN")
	defer os.Unsetenv("ACCESS_TOKEN_SECRET")
	defer os.Unsetenv("REFRESH_TOKEN_SECRET")
	defer os.Unsetenv("ACCESS_TOKEN_EXPIRES")
	defer os.Unsetenv("REFRESH_TOKEN_EXPIRES")
	cfg := GetConfig()
	assert.Equal(t, cfg.ServerAddress, serverAddress)
	assert.Equal(t, cfg.DatabaseDSN, databaseDSN)
	assert.Equal(t, cfg.AccessTokenSecret, accessTokenSecret)
	assert.Equal(t, cfg.RefreshTokenSecret, refreshTokenSecret)
	assert.Equal(t, cfg.AccessTokenExpires, accessTokenExpires)
	assert.Equal(t, cfg.RefreshTokenExpires, refreshTokenExpires)
}

func TestGetConfig_WithoutEnv(t *testing.T) {
	os.Unsetenv("SERVER_ADDRESS")
	os.Unsetenv("DATABASE_DSN")
	os.Unsetenv("ACCESS_TOKEN_SECRET")
	os.Unsetenv("REFRESH_TOKEN_SECRET")
	os.Unsetenv("ACCESS_TOKEN_EXPIRES")
	os.Unsetenv("REFRESH_TOKEN_EXPIRES")
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	os.Args = []string{"cmd"}
	cfg := GetConfig()
	assert.Equal(t, cfg.ServerAddress, "")
	assert.Equal(t, cfg.DatabaseDSN, "")
	assert.Equal(t, cfg.AccessTokenExpires, time.Duration(0))
	assert.Equal(t, cfg.RefreshTokenExpires, time.Duration(0))
	assert.Equal(t, cfg.AccessTokenSecret, "")
	assert.Equal(t, cfg.RefreshTokenSecret, "")
}
