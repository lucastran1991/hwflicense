package tests

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/atprof/license-server/kms/internal/config"
)

// TestLoadConfigFromFile tests loading configuration from setting.json
func TestLoadConfigFromFile(t *testing.T) {
	// Create a temporary config file
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "setting.json")

	// Create test config
	testConfig := map[string]string{
		"kms_db_path": "./test_db.db",
		"kms_port":    "9090",
	}

	data, err := json.Marshal(testConfig)
	if err != nil {
		t.Fatalf("Failed to marshal test config: %v", err)
	}

	if err := os.WriteFile(configPath, data, 0644); err != nil {
		t.Fatalf("Failed to write test config: %v", err)
	}

	// Set environment variables
	// Generate a valid 32-byte master key (base64 encoded)
	masterKey := "dGVzdG1hc3RlcmtleWZvcnRlc3RpbmcwMTIzNDU2Nzg=" // 32 bytes base64
	os.Setenv("KMS_MASTER_KEY", masterKey)
	os.Setenv("KMS_CONFIG_PATH", configPath)
	defer os.Unsetenv("KMS_MASTER_KEY")
	defer os.Unsetenv("KMS_CONFIG_PATH")

	// Load config
	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// Verify values from file
	if cfg.DBPath != "./test_db.db" {
		t.Errorf("Expected DBPath './test_db.db', got '%s'", cfg.DBPath)
	}

	if cfg.Port != ":9090" {
		t.Errorf("Expected Port ':9090', got '%s'", cfg.Port)
	}
}

// TestEnvVarsOverrideFile tests that environment variables override file settings
func TestEnvVarsOverrideFile(t *testing.T) {
	// Create a temporary config file
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "setting.json")

	// Create test config
	testConfig := map[string]string{
		"kms_db_path": "./test_db.db",
		"kms_port":    "9090",
	}

	data, err := json.Marshal(testConfig)
	if err != nil {
		t.Fatalf("Failed to marshal test config: %v", err)
	}

	if err := os.WriteFile(configPath, data, 0644); err != nil {
		t.Fatalf("Failed to write test config: %v", err)
	}

	// Set environment variables (should override file)
	// Generate a valid 32-byte master key (base64 encoded)
	masterKey := "dGVzdG1hc3RlcmtleWZvcnRlc3RpbmcwMTIzNDU2Nzg=" // 32 bytes base64
	os.Setenv("KMS_MASTER_KEY", masterKey)
	os.Setenv("KMS_CONFIG_PATH", configPath)
	os.Setenv("KMS_DB_PATH", "./env_db.db")
	os.Setenv("KMS_PORT", "8080")
	defer func() {
		os.Unsetenv("KMS_MASTER_KEY")
		os.Unsetenv("KMS_CONFIG_PATH")
		os.Unsetenv("KMS_DB_PATH")
		os.Unsetenv("KMS_PORT")
	}()

	// Load config
	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// Verify environment variables override file settings
	if cfg.DBPath != "./env_db.db" {
		t.Errorf("Expected DBPath './env_db.db' (from env), got '%s'", cfg.DBPath)
	}

	if cfg.Port != ":8080" {
		t.Errorf("Expected Port ':8080' (from env), got '%s'", cfg.Port)
	}
}

