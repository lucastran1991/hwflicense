// API functions for Keys endpoints
import { apiClient, handleApiError } from './client';
import type {
  RegisterKeyRequest,
  RegisterKeyResponse,
  ValidateKeyRequest,
  ValidateKeyResponse,
  RefreshKeyRequest,
  RefreshKeyResponse,
  RemoveKeyResponse,
  ListKeysResponse,
  DownloadKeyResponse,
} from '../types/keys';

/**
 * Register a new key or generate one automatically
 */
export async function registerKey(data: RegisterKeyRequest): Promise<RegisterKeyResponse> {
  try {
    const response = await apiClient.post<RegisterKeyResponse>('/keys', data);
    return response.data;
  } catch (error) {
    throw handleApiError(error);
  }
}

/**
 * Validate a symmetric key or verify an Ed25519 signature
 */
export async function validateKey(data: ValidateKeyRequest): Promise<ValidateKeyResponse> {
  try {
    const response = await apiClient.post<ValidateKeyResponse>('/keys/validate', data);
    return response.data;
  } catch (error) {
    throw handleApiError(error);
  }
}

/**
 * Refresh key expiry by extending the expiry time
 */
export async function refreshKey(keyId: string, data: RefreshKeyRequest): Promise<RefreshKeyResponse> {
  try {
    const response = await apiClient.post<RefreshKeyResponse>(`/keys/${keyId}/refresh`, data);
    return response.data;
  } catch (error) {
    throw handleApiError(error);
  }
}

/**
 * Revoke a key by setting its status to revoked
 */
export async function removeKey(keyId: string): Promise<RemoveKeyResponse> {
  try {
    const response = await apiClient.delete<RemoveKeyResponse>(`/keys/${keyId}`);
    return response.data;
  } catch (error) {
    throw handleApiError(error);
  }
}

/**
 * List all keys (without private key material)
 */
export async function listKeys(): Promise<ListKeysResponse> {
  try {
    const response = await apiClient.get<ListKeysResponse>('/keys');
    return response.data;
  } catch (error) {
    throw handleApiError(error);
  }
}

/**
 * Download key material (decrypted key data)
 */
export async function downloadKey(keyId: string): Promise<void> {
  try {
    const response = await apiClient.get<DownloadKeyResponse>(`/keys/${keyId}/download`);
    
    // Create a blob from the JSON response
    const jsonString = JSON.stringify(response.data, null, 2);
    const blob = new Blob([jsonString], {
      type: 'application/json',
    });
    
    // Create a URL for the blob
    const url = window.URL.createObjectURL(blob);
    
    // Create a temporary link element and trigger download
    const link = document.createElement('a');
    link.href = url;
    link.download = `key_${keyId}.json`;
    document.body.appendChild(link);
    link.click();
    
    // Clean up
    document.body.removeChild(link);
    window.URL.revokeObjectURL(url);
  } catch (error) {
    throw handleApiError(error);
  }
}

