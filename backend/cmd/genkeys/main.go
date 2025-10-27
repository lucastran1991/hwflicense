package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"taskmaster-license/pkg/crypto"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run genkeys/main.go <key-name>")
		fmt.Println("Example: go run genkeys/main.go root")
		os.Exit(1)
	}

	keyName := os.Args[1]
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
	if keyName == "root" {
		publicKeyFile := filepath.Join(keysDir, "root_public.pem")
		data, err := os.ReadFile(publicKeyFile)
		if err == nil {
			fmt.Printf("\nAdd this to your .env file as ROOT_PUBLIC_KEY:\n")
			fmt.Printf("%s", string(data))
		}
	}
}
