package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

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
	DefaultMaxSites     int
	DefaultCMLValidity  string
	DefaultFeaturePacks []string
}

type JSONConfig struct {
	Mode               string   `json:"mode"`
	DBPath             string   `json:"db_path"`
	JWTSecret          string   `json:"jwt_secret"`
	APIPort            string   `json:"api_port"`
	APIEnv             string   `json:"api_env"`
	EncryptionPassword string   `json:"encryption_password"`
	AStackMockPort     string   `json:"astack_mock_port"`
	RootPublicKey      string   `json:"root_public_key"`
	DefaultMaxSites     int      `json:"default_max_sites"`
	DefaultCMLValidity string   `json:"default_cml_validity"`
	DefaultFeaturePacks []string `json:"default_feature_packs"`
}

var AppConfig *Config

func LoadConfig() error {
	// Try to load .env file, but don't fail if it doesn't exist
	_ = godotenv.Load()

	// Try to load from JSON config file first
	configPath := "../../config/backend.json"
	if config, err := loadFromJSON(configPath); err == nil {
		AppConfig = &Config{
			DBPath:             config.DBPath,
			JWTSecret:          config.JWTSecret,
			APIPort:            config.APIPort,
			APIEnv:             config.APIEnv,
			RootPublicKey:      config.RootPublicKey,
			AStackMockPort:     config.AStackMockPort,
			EncryptionPassword: config.EncryptionPassword,
			DefaultMaxSites:    config.DefaultMaxSites,
			DefaultCMLValidity: config.DefaultCMLValidity,
			DefaultFeaturePacks: config.DefaultFeaturePacks,
		}
		return nil
	}

	// Fall back to environment variables
	config := &Config{
		DBPath:              getEnv("DB_PATH", "data/taskmaster_license.db"),
		JWTSecret:           getEnv("JWT_SECRET", "taskmaster-secret-key-change-in-production"),
		APIPort:             getEnv("API_PORT", "8080"),
		APIEnv:              getEnv("API_ENV", "development"),
		RootPublicKey:       getEnv("ROOT_PUBLIC_KEY", ""),
		AStackMockPort:      getEnv("ASTACK_MOCK_PORT", "8081"),
		EncryptionPassword:  getEnv("ENCRYPTION_PASSWORD", "change-this-password-in-production-12345"),
		DefaultMaxSites:     100,
		DefaultCMLValidity:  "",
		DefaultFeaturePacks: []string{"basic", "standard"},
	}

	AppConfig = config
	return nil
}

func loadFromJSON(configPath string) (*JSONConfig, error) {
	// Get absolute path
	absPath, err := filepath.Abs(configPath)
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(absPath)
	if err != nil {
		return nil, err
	}

	var config JSONConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
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
