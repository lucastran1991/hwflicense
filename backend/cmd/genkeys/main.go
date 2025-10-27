package main

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"

	"taskmaster-license/internal/config"
	"taskmaster-license/internal/database"
	"taskmaster-license/internal/models"
	"taskmaster-license/internal/repository"
	"taskmaster-license/pkg/crypto"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "root":
		generateRootKeys()
	case "org":
		if len(os.Args) < 4 {
			fmt.Println("Usage: go run genkeys/main.go org <org-id> <dev|prod>")
			os.Exit(1)
		}
		orgID := os.Args[2]
		keyType := os.Args[3]
		if keyType != "dev" && keyType != "prod" {
			log.Fatal("key_type must be 'dev' or 'prod'")
		}
		generateOrgKey(orgID, keyType)
	default:
		fmt.Printf("Unknown command: %s\n", command)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("Usage:")
	fmt.Println("  go run genkeys/main.go root                    - Generate root keys (for file storage)")
	fmt.Println("  go run genkeys/main.go org <org-id> <dev|prod> - Generate org keys (for database storage)")
	fmt.Println("\nExamples:")
	fmt.Println("  go run genkeys/main.go root")
	fmt.Println("  go run genkeys/main.go org org_12345 dev")
}

func generateRootKeys() {
	keyName := "root"
	keysDir := "keys"

	// Create keys directory if it doesn't exist
	if err := os.MkdirAll(keysDir, 0755); err != nil {
		log.Fatalf("Failed to create keys directory: %v", err)
	}

	// Generate key pair
	fmt.Printf("Generating %s key pair...\n", keyName)
	privateKey, publicKey, err := crypto.GenerateKeyPair()
	if err != nil {
		log.Fatalf("Failed to generate key pair: %v", err)
	}

	// Save keys
	if err := crypto.SaveKeyPair(privateKey, publicKey, keysDir, keyName); err != nil {
		log.Fatalf("Failed to save key pair: %v", err)
	}

	fmt.Printf("Keys saved to %s/\n", keysDir)
	fmt.Printf("  - %s_private.pem\n", keyName)
	fmt.Printf("  - %s_public.pem\n", keyName)

	// Display public key
	publicPEM, err := crypto.PublicKeyToPEM(publicKey)
	if err != nil {
		log.Fatalf("Failed to encode public key: %v", err)
	}

	fmt.Printf("\nPublic Key:\n%s\n", publicPEM)

	// Check if these are root keys and should be added to config
	publicKeyFile := filepath.Join(keysDir, "root_public.pem")
	data, err := os.ReadFile(publicKeyFile)
	if err == nil {
		fmt.Printf("\nAdd this to your .env file as ROOT_PUBLIC_KEY:\n")
		fmt.Printf("%s", string(data))
	}
}

func generateOrgKey(orgID, keyType string) {
	// Load configuration
	if err := config.LoadConfig(); err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize database
	dbPath := config.AppConfig.GetDatabaseConnectionString()
	db, err := database.NewDB(dbPath)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Initialize repository
	repo := repository.NewRepository(db)

	fmt.Printf("Generating org key for org_id=%s, key_type=%s...\n", orgID, keyType)

	// Check if key already exists
	existingKey, err := repo.GetOrgKey(orgID, keyType)
	if err == nil && existingKey != nil {
		log.Fatalf("Org key already exists for org_id=%s, key_type=%s (ID: %s)", orgID, keyType, existingKey.ID)
	}

	// Generate ECDSA P-256 key pair
	privateKey, publicKey, err := crypto.GenerateKeyPair()
	if err != nil {
		log.Fatalf("Failed to generate key pair: %v", err)
	}

	// Encode private key to PEM
	privateKeyBytes, err := x509.MarshalECPrivateKey(privateKey)
	if err != nil {
		log.Fatalf("Failed to marshal private key: %v", err)
	}

	privateBlock := &pem.Block{
		Type:  "EC PRIVATE KEY",
		Bytes: privateKeyBytes,
	}
	privateKeyPEM := pem.EncodeToMemory(privateBlock)

	// Encrypt private key
	encryptedPrivateKey, err := crypto.EncryptPrivateKey(privateKeyPEM, config.AppConfig.EncryptionPassword)
	if err != nil {
		log.Fatalf("Failed to encrypt private key: %v", err)
	}

	// Encode public key to PEM
	publicKeyPEM, err := crypto.PublicKeyToPEM(publicKey)
	if err != nil {
		log.Fatalf("Failed to encode public key: %v", err)
	}

	// Create org key model
	orgKey := &models.OrgKey{
		ID:                 uuid.New().String(),
		OrgID:              orgID,
		KeyType:            keyType,
		PrivateKeyEncrypted: encryptedPrivateKey,
		PublicKey:          publicKeyPEM,
		CreatedAt:          time.Now(),
	}

	// Store in database
	if err := repo.CreateOrgKey(orgKey); err != nil {
		log.Fatalf("Failed to store org key: %v", err)
	}

	fmt.Printf("✓ Org key created successfully!\n")
	fmt.Printf("  - ID: %s\n", orgKey.ID)
	fmt.Printf("  - Org ID: %s\n", orgID)
	fmt.Printf("  - Key Type: %s\n", keyType)
	fmt.Printf("\nPublic Key:\n%s\n", publicKeyPEM)
	fmt.Printf("\n⚠️  Private key is encrypted and stored in database.\n")
	fmt.Printf("⚠️  Keep ENCRYPTION_PASSWORD environment variable secure!\n")
}
