package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoad_Defaults(t *testing.T) {
	// Clear env vars to test defaults
	os.Clearenv()

	cfg := Load()

	assert.Equal(t, "8080", cfg.ServerPort)
	assert.Equal(t, "localhost", cfg.DBHost)
	assert.Equal(t, "5432", cfg.DBPort)
	assert.Equal(t, "postgres", cfg.DBUser)
	assert.Equal(t, "postgres", cfg.DBPassword)
	assert.Equal(t, "devices", cfg.DBName)
	assert.Equal(t, "disable", cfg.DBSSLMode)
}

func TestLoad_WithEnvVars(t *testing.T) {
	os.Setenv("PORT", "9090")
	os.Setenv("DB_HOST", "db.example.com")
	os.Setenv("DB_PORT", "5433")
	os.Setenv("DB_USER", "admin")
	os.Setenv("DB_PASSWORD", "secret")
	os.Setenv("DB_NAME", "mydb")
	os.Setenv("DB_SSLMODE", "require")

	cfg := Load()

	assert.Equal(t, "9090", cfg.ServerPort)
	assert.Equal(t, "db.example.com", cfg.DBHost)
	assert.Equal(t, "5433", cfg.DBPort)
	assert.Equal(t, "admin", cfg.DBUser)
	assert.Equal(t, "secret", cfg.DBPassword)
	assert.Equal(t, "mydb", cfg.DBName)
	assert.Equal(t, "require", cfg.DBSSLMode)
}

func TestGetDBDSN(t *testing.T) {
	cfg := &Config{
		DBHost:     "localhost",
		DBPort:     "5432",
		DBUser:     "postgres",
		DBPassword: "postgres",
		DBName:     "devices",
		DBSSLMode:  "disable",
	}

	dsn := cfg.GetDBDSN()

	assert.Contains(t, dsn, "host=localhost")
	assert.Contains(t, dsn, "port=5432")
	assert.Contains(t, dsn, "user=postgres")
	assert.Contains(t, dsn, "password=postgres")
	assert.Contains(t, dsn, "dbname=devices")
	assert.Contains(t, dsn, "sslmode=disable")
}

func TestGetServerAddr(t *testing.T) {
	cfg := &Config{ServerPort: "8080"}
	assert.Equal(t, ":8080", cfg.GetServerAddr())
}