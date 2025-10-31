package config

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	// DefaultDBPath is the default path for the BoltDB database
	DefaultDBPath = "./kms.db"
	// DefaultPort is the default server port
	DefaultPort = ":8080"
	// MasterKeySize is the required size for the master encryption key (256 bits = 32 bytes)
	MasterKeySize = 32
	// DefaultConfigPath is the default path to settings JSON file
	DefaultConfigPath = "./config/setting.json"
)

// Settings represents the settings from JSON file
type Settings struct {
	KMSDBPath string `json:"kms_db_path"`
	KMSPort   string `json:"kms_port"`
}

// Config holds the application configuration
type Config struct {
	MasterKey []byte
	DBPath    string
	Port      string
}

// loadSettingsFromFile loads settings from JSON file if it exists
func loadSettingsFromFile(configPath string) (*Settings, error) {
	// Try to get config path from environment, otherwise use default
	if configPath == "" {
		configPath = os.Getenv("KMS_CONFIG_PATH")
		if configPath == "" {
			configPath = DefaultConfigPath
		}
	}

	// Check if file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, nil // File doesn't exist, return nil (not an error)
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var settings Settings
	if err := json.Unmarshal(data, &settings); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return &settings, nil
}

// loadMasterKeyFromFile loads master key from a secure file
func loadMasterKeyFromFile(filePath string) (string, error) {
	// Check file permissions (should be 600 or more restrictive)
	info, err := os.Stat(filePath)
	if err != nil {
		return "", err // File doesn't exist or can't access
	}
	
	// Check if file has restrictive permissions (only owner can read/write)
	mode := info.Mode()
	if mode.Perm()&077 != 0 {
		return "", fmt.Errorf("master key file %s has insecure permissions (should be 600 or more restrictive)", filePath)
	}
	
	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	
	// Remove whitespace and newlines
	masterKey := strings.TrimSpace(string(data))
	return masterKey, nil
}

// Load loads configuration from settings file and environment variables
// Environment variables take precedence over settings file
func Load() (*Config, error) {
	// Load settings from file (if exists)
	configPath := os.Getenv("KMS_CONFIG_PATH")
	if configPath == "" {
		// Try to find config file relative to current directory or executable
		execPath, err := os.Executable()
		if err == nil {
			execDir := filepath.Dir(execPath)
			configPath = filepath.Join(execDir, DefaultConfigPath)
			if _, err := os.Stat(configPath); os.IsNotExist(err) {
				// Try relative to working directory
				configPath = DefaultConfigPath
			}
		} else {
			configPath = DefaultConfigPath
		}
	}

	settings, err := loadSettingsFromFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load settings from file: %w", err)
	}

	// Start with defaults or file settings
	dbPath := DefaultDBPath
	port := DefaultPort

	if settings != nil {
		if settings.KMSDBPath != "" {
			dbPath = settings.KMSDBPath
		}
		if settings.KMSPort != "" {
			port = settings.KMSPort
		}
	}

	// Environment variables override file settings
	if envDBPath := os.Getenv("KMS_DB_PATH"); envDBPath != "" {
		dbPath = envDBPath
	}

	if envPort := os.Getenv("KMS_PORT"); envPort != "" {
		port = envPort
	}

	// Normalize port format (ensure it has colon prefix)
	if port[0] != ':' {
		port = ":" + port
	}

	// Load master key from environment or secure file
	masterKeyBase64 := os.Getenv("KMS_MASTER_KEY")
	
	// If not in environment, try to load from secure file
	if masterKeyBase64 == "" {
		// Try common secure file locations
		secretPaths := []string{
			"./secrets/master.key",
			"../secrets/master.key",
			"/etc/kms/master.key",
		}
		
		for _, secretPath := range secretPaths {
			if key, err := loadMasterKeyFromFile(secretPath); err == nil {
				masterKeyBase64 = key
				break
			}
		}
	}
	
	if masterKeyBase64 == "" {
		return nil, fmt.Errorf("master key not found: set KMS_MASTER_KEY environment variable or create ./secrets/master.key file")
	}

	masterKey, err := base64.StdEncoding.DecodeString(masterKeyBase64)
	if err != nil {
		return nil, fmt.Errorf("failed to decode master key: %w", err)
	}

	if len(masterKey) != MasterKeySize {
		return nil, fmt.Errorf("master key must be exactly %d bytes (got %d bytes), base64 encoded", MasterKeySize, len(masterKey))
	}

	return &Config{
		MasterKey: masterKey,
		DBPath:    dbPath,
		Port:      port,
	}, nil
}

