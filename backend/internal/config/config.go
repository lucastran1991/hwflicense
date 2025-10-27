package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBPath              string
	JWTSecret           string
	APIPort             string
	APIEnv              string
	RootPublicKey       string
	AStackMockPort      string
	EncryptionPassword  string
}

var AppConfig *Config

func LoadConfig() error {
	// Try to load .env file, but don't fail if it doesn't exist
	_ = godotenv.Load()

	config := &Config{
		DBPath:             getEnv("DB_PATH", "data/taskmaster_license.db"),
		JWTSecret:          getEnv("JWT_SECRET", "taskmaster-secret-key-change-in-production"),
		APIPort:            getEnv("API_PORT", "8080"),
		APIEnv:             getEnv("API_ENV", "development"),
		RootPublicKey:      getEnv("ROOT_PUBLIC_KEY", ""),
		AStackMockPort:     getEnv("ASTACK_MOCK_PORT", "8081"),
		EncryptionPassword: getEnv("ENCRYPTION_PASSWORD", "change-this-password-in-production-12345"),
	}

	AppConfig = config
	return nil
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func (c *Config) GetDatabaseConnectionString() string {
	return c.DBPath
}

func (c *Config) IsProduction() bool {
	return c.APIEnv == "production"
}

func (c *Config) ServerAddress() string {
	return fmt.Sprintf(":%s", c.APIPort)
}

func (c *Config) AStackAddress() string {
	return fmt.Sprintf(":%s", c.AStackMockPort)
}
