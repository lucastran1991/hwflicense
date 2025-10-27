package crypto

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"math/big"
	"os"
)

// GenerateKeyPair generates a new ECDSA P-256 key pair
func GenerateKeyPair() (*ecdsa.PrivateKey, *ecdsa.PublicKey, error) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, nil, err
	}
	return privateKey, &privateKey.PublicKey, nil
}

// SignData signs arbitrary data with a private key
func SignData(data []byte, privateKey *ecdsa.PrivateKey) (string, error) {
	hash := sha256.Sum256(data)
	r, s, err := ecdsa.Sign(rand.Reader, privateKey, hash[:])
	if err != nil {
		return "", err
	}

	// Encode signature as base64
	sig := r.Bytes()
	sig = append(sig, s.Bytes()...)
	return base64.StdEncoding.EncodeToString(sig), nil
}

// VerifySignature verifies a signature against data and a public key
func VerifySignature(data []byte, signature string, publicKey *ecdsa.PublicKey) (bool, error) {
	hash := sha256.Sum256(data)
	sigBytes, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		return false, err
	}

	if len(sigBytes) != 64 {
		return false, fmt.Errorf("invalid signature length")
	}

	r := new(big.Int).SetBytes(sigBytes[:32])
	s := new(big.Int).SetBytes(sigBytes[32:])

	return ecdsa.Verify(publicKey, hash[:], r, s), nil
}

// SignJSON signs a JSON object
func SignJSON(obj interface{}, privateKey *ecdsa.PrivateKey) (string, error) {
	jsonData, err := json.Marshal(obj)
	if err != nil {
		return "", err
	}
	return SignData(jsonData, privateKey)
}

// VerifyJSON verifies a JSON object's signature
func VerifyJSON(obj interface{}, signature string, publicKey *ecdsa.PublicKey) (bool, error) {
	jsonData, err := json.Marshal(obj)
	if err != nil {
		return false, err
	}
	return VerifySignature(jsonData, signature, publicKey)
}

// PublicKeyToPEM encodes a public key to PEM format
func PublicKeyToPEM(publicKey *ecdsa.PublicKey) (string, error) {
	pubASN1, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return "", err
	}

	pubBlock := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: pubASN1,
	}

	pubPEM := pem.EncodeToMemory(pubBlock)
	return string(pubPEM), nil
}

// PEMToPublicKey decodes a PEM-encoded public key
func PEMToPublicKey(pemData string) (*ecdsa.PublicKey, error) {
	block, _ := pem.Decode([]byte(pemData))
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM block")
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	ecdsaPub, ok := pub.(*ecdsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("not an ECDSA public key")
	}

	return ecdsaPub, nil
}

// SaveKeyPair saves a key pair to files
func SaveKeyPair(privateKey *ecdsa.PrivateKey, publicKey *ecdsa.PublicKey, dir, prefix string) error {
	// Save private key (PEM format)
	privateKeyBytes, err := x509.MarshalECPrivateKey(privateKey)
	if err != nil {
		return err
	}

	privateBlock := &pem.Block{
		Type:  "EC PRIVATE KEY",
		Bytes: privateKeyBytes,
	}

	privatePEM := pem.EncodeToMemory(privateBlock)
	privateFile, err := os.Create(fmt.Sprintf("%s/%s_private.pem", dir, prefix))
	if err != nil {
		return err
	}
	defer privateFile.Close()

	if _, err := privateFile.Write(privatePEM); err != nil {
		return err
	}

	// Save public key
	publicPEM, err := PublicKeyToPEM(publicKey)
	if err != nil {
		return err
	}

	publicFile, err := os.Create(fmt.Sprintf("%s/%s_public.pem", dir, prefix))
	if err != nil {
		return err
	}
	defer publicFile.Close()

	if _, err := publicFile.Write([]byte(publicPEM)); err != nil {
		return err
	}

	return nil
}

// LoadPrivateKey loads a private key from a PEM file
func LoadPrivateKey(filename string) (*ecdsa.PrivateKey, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(data)
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM block")
	}

	privateKey, err := x509.ParseECPrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return privateKey, nil
}

// LoadPublicKey loads a public key from a PEM file
func LoadPublicKey(filename string) (*ecdsa.PublicKey, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	return PEMToPublicKey(string(data))
}
