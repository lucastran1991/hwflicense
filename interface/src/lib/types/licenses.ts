// Type definitions for Licenses API matching Go backend

export interface GenerateLicenseRequest {
  key_id: string;
  license_type: string;
  metadata?: Record<string, string>;
}

export interface GenerateLicenseResponse {
  license_file: string; // Base64 encoded license file content
  filename: string; // Suggested filename (e.g., "enterprise.lic")
  license_id: string;
}

export interface ValidateLicenseRequest {
  license_content?: string; // Base64 encoded license file (for JSON body)
  // Note: Also supports multipart file upload (handled separately)
}

export interface ValidateLicenseResponse {
  valid: boolean;
  license_id?: string;
  license_type?: string;
  key_id?: string;
  expires_at?: string; // ISO 8601 timestamp
  expired: boolean;
  revoked: boolean;
  metadata?: Record<string, string>;
  error?: string; // Error message if validation failed
}

export interface LicenseFile {
  license_id: string;
  license_type: string;
  key_id: string;
  key_type: string;
  public_key?: string; // Base64 encoded, only for asymmetric keys
  issued_at: string; // ISO 8601 timestamp
  expires_at: string; // ISO 8601 timestamp
  metadata?: Record<string, string>;
  signature: string; // Base64 encoded HMAC-SHA256 signature
}

