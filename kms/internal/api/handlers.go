package api

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/atprof/license-server/kms/internal/crypto"
	"github.com/atprof/license-server/kms/internal/licenses"
	"github.com/atprof/license-server/kms/internal/storage"
	"github.com/atprof/license-server/kms/pkg/errors"
)

// Handler holds dependencies for API handlers
type Handler struct {
	store     *storage.BoltStore
	masterKey []byte
}

// NewHandler creates a new API handler instance
func NewHandler(store *storage.BoltStore, masterKey []byte) *Handler {
	return &Handler{
		store:     store,
		masterKey: masterKey,
	}
}

// RegisterKeyRequest represents a request to register a key
type RegisterKeyRequest struct {
	KeyType          string `json:"key_type" binding:"required,oneof=symmetric asymmetric"`
	ExpiresInSeconds int64  `json:"expires_in_seconds"` // Optional, default 1 year
	KeyMaterial      string `json:"key_material,omitempty"` // Optional base64 encoded key for external keys
}

// RegisterKeyResponse represents a response from registering a key
type RegisterKeyResponse struct {
	KeyID     string    `json:"key_id"`
	KeyType   string    `json:"key_type"`
	PublicKey string    `json:"public_key,omitempty"` // Base64 encoded, only for asymmetric
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}

// RegisterKey handles POST /keys - Register or generate a key
func (h *Handler) RegisterKey(c *gin.Context) {
	var req RegisterKeyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var key *storage.Key
	var err error
	var publicKeyBase64 string

	now := time.Now().UTC()
	expiresIn := req.ExpiresInSeconds
	if expiresIn == 0 {
		expiresIn = 365 * 24 * 60 * 60 // Default: 1 year
	}
	expiresAt := now.Add(time.Duration(expiresIn) * time.Second)

	keyID := uuid.New().String()
	version := 1

	if req.KeyType == "symmetric" {
		var keyMaterial []byte
		if req.KeyMaterial != "" {
			keyMaterial, err = base64.StdEncoding.DecodeString(req.KeyMaterial)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid key_material: must be base64 encoded"})
				return
			}
			if len(keyMaterial) != crypto.SymmetricKeySize {
				c.JSON(http.StatusBadRequest, gin.H{"error": "symmetric key must be 32 bytes (256 bits)"})
				return
			}
		} else {
			keyMaterial, err = crypto.GenerateSymmetricKey()
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate key"})
				return
			}
		}

		// Encrypt the key material
		encryptedKey, err := crypto.EncryptKey(h.masterKey, keyMaterial)
		if err != nil {
			// Zero out key material before returning error
			for i := range keyMaterial {
				keyMaterial[i] = 0
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to encrypt key"})
			return
		}

		// Zero out plaintext key material
		for i := range keyMaterial {
			keyMaterial[i] = 0
		}

		key = &storage.Key{
			ID:                 keyID,
			KeyType:            storage.KeyTypeSymmetric,
			PublicKey:          nil,
			EncryptedPrivateKey: encryptedKey,
			ExpiresAt:          expiresAt,
			CreatedAt:          now,
			Status:             storage.KeyStatusActive,
			Version:            version,
		}

	} else if req.KeyType == "asymmetric" {
		var publicKey, privateKey []byte
		if req.KeyMaterial != "" {
			// For asymmetric, key_material would be the private key
			privateKey, err = base64.StdEncoding.DecodeString(req.KeyMaterial)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid key_material: must be base64 encoded"})
				return
			}
			// Extract public key from private key (Ed25519)
			if len(privateKey) != crypto.Ed25519PrivateKeySize {
				c.JSON(http.StatusBadRequest, gin.H{"error": "asymmetric private key must be 64 bytes (Ed25519)"})
				return
			}
			publicKey = privateKey[32:] // Last 32 bytes are the public key for Ed25519
		} else {
			publicKey, privateKey, err = crypto.GenerateAsymmetricKeyPair()
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate key pair"})
				return
			}
		}

		// Encrypt the private key
		encryptedPrivateKey, err := crypto.EncryptKey(h.masterKey, privateKey)
		if err != nil {
			// Zero out keys before returning error
			for i := range privateKey {
				privateKey[i] = 0
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to encrypt private key"})
			return
		}

		publicKeyBase64 = base64.StdEncoding.EncodeToString(publicKey)

		// Zero out plaintext private key
		for i := range privateKey {
			privateKey[i] = 0
		}

		key = &storage.Key{
			ID:                 keyID,
			KeyType:            storage.KeyTypeAsymmetric,
			PublicKey:          publicKey,
			EncryptedPrivateKey: encryptedPrivateKey,
			ExpiresAt:          expiresAt,
			CreatedAt:          now,
			Status:             storage.KeyStatusActive,
			Version:            version,
		}
	}

	// Store the key
	if err := h.store.StoreKey(key); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to store key"})
		return
	}

	resp := RegisterKeyResponse{
		KeyID:     key.ID,
		KeyType:   string(key.KeyType),
		ExpiresAt: key.ExpiresAt,
		CreatedAt: key.CreatedAt,
	}

	if key.KeyType == storage.KeyTypeAsymmetric {
		resp.PublicKey = publicKeyBase64
	}

	c.JSON(http.StatusOK, resp)
}

// ValidateKeyRequest represents a request to validate a key
type ValidateKeyRequest struct {
	KeyID      string `json:"key_id" binding:"required"`
	KeyMaterial string `json:"key_material,omitempty"`       // For symmetric keys
	Message     string `json:"message,omitempty"`             // For asymmetric signature validation
	Signature   string `json:"signature,omitempty"`           // Base64 encoded signature
}

// ValidateKeyResponse represents a response from validating a key
type ValidateKeyResponse struct {
	Valid   bool `json:"valid"`
	Expired bool `json:"expired"`
	Revoked bool `json:"revoked"`
}

// ValidateKey handles POST /keys/validate - Validate a key or signature
func (h *Handler) ValidateKey(c *gin.Context) {
	var req ValidateKeyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get the key from storage
	key, err := h.store.GetKey(req.KeyID)
	if err != nil {
		if err == errors.ErrKeyNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "key not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve key"})
		return
	}

	resp := ValidateKeyResponse{
		Expired: key.IsExpired(),
		Revoked: key.IsRevoked(),
		Valid:   false,
	}

	// Check if key is expired or revoked
	if resp.Expired || resp.Revoked {
		c.JSON(http.StatusOK, resp)
		return
	}

	// Validate based on key type
	if key.KeyType == storage.KeyTypeSymmetric {
		if req.KeyMaterial == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "key_material is required for symmetric keys"})
			return
		}

		providedKey, err := base64.StdEncoding.DecodeString(req.KeyMaterial)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid key_material: must be base64 encoded"})
			return
		}
		defer func() {
			for i := range providedKey {
				providedKey[i] = 0
			}
		}()

		valid, err := crypto.ValidateSymmetricKey(h.masterKey, key.EncryptedPrivateKey, providedKey)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to validate key"})
			return
		}

		resp.Valid = valid

	} else if key.KeyType == storage.KeyTypeAsymmetric {
		if req.Message == "" || req.Signature == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "message and signature are required for asymmetric keys"})
			return
		}

		signature, err := base64.StdEncoding.DecodeString(req.Signature)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid signature: must be base64 encoded"})
			return
		}

		valid, err := crypto.ValidateSignature(key.PublicKey, []byte(req.Message), signature)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to validate signature"})
			return
		}

		resp.Valid = valid
	}

	c.JSON(http.StatusOK, resp)
}

// RefreshKeyRequest represents a request to refresh a key's expiry
type RefreshKeyRequest struct {
	ExpiresInSeconds int64 `json:"expires_in_seconds" binding:"required"`
}

// RefreshKeyResponse represents a response from refreshing a key
type RefreshKeyResponse struct {
	KeyID       string    `json:"key_id"`
	NewExpiresAt time.Time `json:"new_expires_at"`
}

// RefreshKey handles POST /keys/:id/refresh - Refresh key expiry
func (h *Handler) RefreshKey(c *gin.Context) {
	keyID := c.Param("id")
	if keyID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "key_id is required"})
		return
	}

	var req RefreshKeyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get the key
	key, err := h.store.GetKey(keyID)
	if err != nil {
		if err == errors.ErrKeyNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "key not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve key"})
		return
	}

	// Check if key is revoked
	if key.IsRevoked() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cannot refresh revoked key"})
		return
	}

	// Calculate new expiry
	newExpiresAt := time.Now().UTC().Add(time.Duration(req.ExpiresInSeconds) * time.Second)

	// Update expiry
	if err := h.store.UpdateKeyExpiry(keyID, newExpiresAt); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update expiry"})
		return
	}

	c.JSON(http.StatusOK, RefreshKeyResponse{
		KeyID:       keyID,
		NewExpiresAt: newExpiresAt,
	})
}

// RemoveKeyResponse represents a response from removing a key
type RemoveKeyResponse struct {
	Success bool   `json:"success"`
	KeyID   string `json:"key_id"`
}

// RemoveKey handles DELETE /keys/:id - Revoke a key
func (h *Handler) RemoveKey(c *gin.Context) {
	keyID := c.Param("id")
	if keyID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "key_id is required"})
		return
	}

	// Check if key exists
	_, err := h.store.GetKey(keyID)
	if err != nil {
		if err == errors.ErrKeyNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "key not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve key"})
		return
	}

	// Revoke the key
	if err := h.store.RevokeKey(keyID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to revoke key"})
		return
	}

	c.JSON(http.StatusOK, RemoveKeyResponse{
		Success: true,
		KeyID:   keyID,
	})
}

// ListKeysResponse represents a response from listing keys
type ListKeysResponse struct {
	Keys []KeyInfo `json:"keys"`
}

// KeyInfo represents key information without private key material
type KeyInfo struct {
	KeyID     string    `json:"key_id"`
	KeyType   string    `json:"key_type"`
	PublicKey string    `json:"public_key,omitempty"` // Base64 encoded, only for asymmetric keys
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
	Status    string    `json:"status"`
	Version   int       `json:"version"`
	Expired   bool      `json:"expired"`
	Revoked   bool      `json:"revoked"`
}

// ListKeys handles GET /keys - List all keys
func (h *Handler) ListKeys(c *gin.Context) {
	// Get all keys from storage (without private key material)
	keys, err := h.store.ListKeys()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list keys"})
		return
	}

	// Convert storage.Key to KeyInfo
	keyInfos := make([]KeyInfo, 0, len(keys))
	for _, key := range keys {
		keyInfo := KeyInfo{
			KeyID:     key.ID,
			KeyType:   string(key.KeyType),
			ExpiresAt: key.ExpiresAt,
			CreatedAt: key.CreatedAt,
			Status:    string(key.Status),
			Version:   key.Version,
			Expired:   key.IsExpired(),
			Revoked:   key.IsRevoked(),
		}

		// Include public key for asymmetric keys
		if key.KeyType == storage.KeyTypeAsymmetric && key.PublicKey != nil {
			keyInfo.PublicKey = base64.StdEncoding.EncodeToString(key.PublicKey)
		}

		keyInfos = append(keyInfos, keyInfo)
	}

	c.JSON(http.StatusOK, ListKeysResponse{
		Keys: keyInfos,
	})
}

// DownloadKeyResponse represents a downloadable key file structure
type DownloadKeyResponse struct {
	KeyID       string `json:"key_id"`
	KeyType     string `json:"key_type"`
	PublicKey   string `json:"public_key,omitempty"`   // Base64 encoded, only for asymmetric keys
	PrivateKey  string `json:"private_key,omitempty"` // Base64 encoded decrypted key material
	SymmetricKey string `json:"symmetric_key,omitempty"` // Base64 encoded decrypted key (for symmetric keys)
	CreatedAt   string `json:"created_at"`             // ISO 8601 timestamp
	ExpiresAt   string `json:"expires_at"`             // ISO 8601 timestamp
	Status      string `json:"status"`
	Version     int    `json:"version"`
	Warning     string `json:"warning,omitempty"`       // Security warning
}

// DownloadKey handles GET /keys/:id/download - Download key material
func (h *Handler) DownloadKey(c *gin.Context) {
	keyID := c.Param("id")
	if keyID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "key_id is required"})
		return
	}

	// Get the key from storage (with private key material)
	key, err := h.store.GetKey(keyID)
	if err != nil {
		if err == errors.ErrKeyNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "key not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve key"})
		return
	}

	// Build response structure
	response := DownloadKeyResponse{
		KeyID:     key.ID,
		KeyType:   string(key.KeyType),
		CreatedAt: key.CreatedAt.Format(time.RFC3339),
		ExpiresAt: key.ExpiresAt.Format(time.RFC3339),
		Status:    string(key.Status),
		Version:   key.Version,
		Warning:   "⚠️ SECURITY WARNING: This file contains sensitive key material. Keep it secure and never share publicly!",
	}

	// Include public key for asymmetric keys
	if key.KeyType == storage.KeyTypeAsymmetric && key.PublicKey != nil {
		response.PublicKey = base64.StdEncoding.EncodeToString(key.PublicKey)
	}

	// Decrypt private key material
	var decryptedKey []byte
	if key.EncryptedPrivateKey != nil && len(key.EncryptedPrivateKey) > 0 {
		decryptedKey, err = crypto.DecryptKey(h.masterKey, key.EncryptedPrivateKey)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to decrypt key material"})
			return
		}
		defer func() {
			// Zero out decrypted key after use
			for i := range decryptedKey {
				decryptedKey[i] = 0
			}
		}()
	}

	// Set decrypted key material based on key type
	if key.KeyType == storage.KeyTypeSymmetric {
		// For symmetric keys, return as symmetric_key field
		response.SymmetricKey = base64.StdEncoding.EncodeToString(decryptedKey)
	} else if key.KeyType == storage.KeyTypeAsymmetric {
		// For asymmetric keys, return as private_key field
		response.PrivateKey = base64.StdEncoding.EncodeToString(decryptedKey)
	}

	// Set headers for file download
	c.Header("Content-Type", "application/json")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"key_%s.json\"", key.ID))

	// Return JSON response (will be downloaded as file)
	c.JSON(http.StatusOK, response)
}

// GenerateLicense handles POST /licenses/generate - Generate a license file
func (h *Handler) GenerateLicense(c *gin.Context) {
	var req licenses.GenerateLicenseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Retrieve key from storage
	key, err := h.store.GetKey(req.KeyID)
	if err != nil {
		if err == errors.ErrKeyNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "key not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve key"})
		return
	}

	// Verify key is valid (not expired/revoked)
	if !key.IsValid() {
		if key.IsExpired() {
			c.JSON(http.StatusBadRequest, gin.H{"error": "cannot generate license for expired key"})
			return
		}
		if key.IsRevoked() {
			c.JSON(http.StatusBadRequest, gin.H{"error": "cannot generate license for revoked key"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": "key is not valid"})
		return
	}

	// Generate license file
	_, licenseBytes, err := licenses.GenerateLicense(key, req.LicenseType, req.Metadata, h.masterKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Determine filename based on license type
	filename := req.LicenseType + ".lic"
	if filename == ".lic" {
		filename = "license.lic"
	}

	// Encode license file content to base64
	licenseBase64 := base64.StdEncoding.EncodeToString(licenseBytes)

	resp := licenses.GenerateLicenseResponse{
		LicenseFile: licenseBase64,
		Filename:    filename,
	}

	// Parse the generated license to get license ID
	var licenseFile licenses.LicenseFile
	if err := json.Unmarshal(licenseBytes, &licenseFile); err == nil {
		resp.LicenseID = licenseFile.LicenseID
	}

	c.JSON(http.StatusOK, resp)
}

// ValidateLicense handles POST /licenses/validate - Validate a license file
func (h *Handler) ValidateLicense(c *gin.Context) {
	var fileContent []byte
	var err error

	// Support two input methods: multipart file upload or JSON body
	contentType := c.GetHeader("Content-Type")
	
	if contentType == "application/json" || contentType == "" {
		// JSON body: expect license_content field
		var req licenses.ValidateLicenseRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if req.LicenseContent == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "license_content is required"})
			return
		}

		// Decode base64 content
		fileContent, err = base64.StdEncoding.DecodeString(req.LicenseContent)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid license_content: must be base64 encoded"})
			return
		}
	} else {
		// Multipart form data: expect file field
		file, err := c.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "file is required in multipart form data"})
			return
		}

		// Open uploaded file
		src, err := file.Open()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "failed to open uploaded file"})
			return
		}
		defer src.Close()

		// Read file content
		fileContent, err = io.ReadAll(src)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "failed to read file content"})
			return
		}
	}

	// Validate license file
	result, err := licenses.ValidateLicense(fileContent, h.store, h.masterKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Convert ValidationResult to ValidateLicenseResponse
	resp := licenses.ValidateLicenseResponse{
		Valid:       result.Valid,
		LicenseID:   result.LicenseID,
		LicenseType: result.LicenseType,
		KeyID:       result.KeyID,
		ExpiresAt:   result.ExpiresAt,
		Expired:     result.Expired,
		Revoked:     result.Revoked,
		Metadata:    result.Metadata,
		Error:       result.Error,
	}

	if !resp.Valid {
		c.JSON(http.StatusOK, resp)
		return
	}

	c.JSON(http.StatusOK, resp)
}

