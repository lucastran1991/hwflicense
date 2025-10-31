// Type definitions for Keys API matching Go backend

export type KeyType = 'symmetric' | 'asymmetric';

export type KeyStatus = 'active' | 'revoked';

export interface RegisterKeyRequest {
  key_type: KeyType;
  expires_in_seconds?: number;
  key_material?: string; // Base64 encoded, optional
}

export interface RegisterKeyResponse {
  key_id: string;
  key_type: KeyType;
  public_key?: string; // Base64 encoded, only for asymmetric keys
  expires_at: string; // ISO 8601 timestamp
  created_at: string; // ISO 8601 timestamp
}

export interface ValidateKeyRequest {
  key_id: string;
  key_material?: string; // Required for symmetric keys
  message?: string; // Required for asymmetric keys
  signature?: string; // Required for asymmetric keys, base64 encoded
}

export interface ValidateKeyResponse {
  valid: boolean;
  expired: boolean;
  revoked: boolean;
}

export interface RefreshKeyRequest {
  expires_in_seconds: number;
}

export interface RefreshKeyResponse {
  key_id: string;
  new_expires_at: string; // ISO 8601 timestamp
}

export interface RemoveKeyResponse {
  success: boolean;
  key_id: string;
}

export interface Key {
  key_id: string;
  key_type: KeyType;
  public_key?: string;
  expires_at: string;
  created_at: string;
  status?: KeyStatus;
}

export interface KeyInfo {
  key_id: string;
  key_type: KeyType;
  public_key?: string; // Base64 encoded, only for asymmetric keys
  expires_at: string; // ISO 8601 timestamp
  created_at: string; // ISO 8601 timestamp
  status: string;
  version: number;
  expired: boolean;
  revoked: boolean;
}

export interface ListKeysResponse {
  keys: KeyInfo[];
}

export interface DownloadKeyResponse {
  key_id: string;
  key_type: KeyType;
  public_key?: string; // Base64 encoded, only for asymmetric keys
  private_key?: string; // Base64 encoded decrypted key material (for asymmetric keys)
  symmetric_key?: string; // Base64 encoded decrypted key (for symmetric keys)
  created_at: string; // ISO 8601 timestamp
  expires_at: string; // ISO 8601 timestamp
  status: string;
  version: number;
  warning?: string; // Security warning
}

