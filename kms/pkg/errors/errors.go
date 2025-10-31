package errors

import "fmt"

// Error definitions for the key management service

var (
	// ErrKeyNotFound indicates the requested key was not found
	ErrKeyNotFound = fmt.Errorf("key not found")
	
	// ErrKeyExpired indicates the key has expired
	ErrKeyExpired = fmt.Errorf("key expired")
	
	// ErrKeyRevoked indicates the key has been revoked
	ErrKeyRevoked = fmt.Errorf("key revoked")
	
	// ErrInvalidKeyMaterial indicates the provided key material is invalid
	ErrInvalidKeyMaterial = fmt.Errorf("invalid key material")
	
	// ErrInvalidSignature indicates the signature verification failed
	ErrInvalidSignature = fmt.Errorf("invalid signature")
	
	// ErrEncryptionFailed indicates encryption operation failed
	ErrEncryptionFailed = fmt.Errorf("encryption failed")
	
	// ErrDecryptionFailed indicates decryption operation failed
	ErrDecryptionFailed = fmt.Errorf("decryption failed")
	
	// ErrInvalidLicenseFile indicates the license file is invalid or malformed
	ErrInvalidLicenseFile = fmt.Errorf("invalid license file")
	
	// ErrLicenseSignatureInvalid indicates the license signature verification failed
	ErrLicenseSignatureInvalid = fmt.Errorf("invalid license signature")
	
	// ErrLicenseExpired indicates the license has expired
	ErrLicenseExpired = fmt.Errorf("license expired")
	
	// ErrLicenseRevoked indicates the license has been revoked
	ErrLicenseRevoked = fmt.Errorf("license revoked")
)

