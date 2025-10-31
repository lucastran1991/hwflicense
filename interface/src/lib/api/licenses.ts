// API functions for Licenses endpoints
import { apiClient, handleApiError } from './client';
import type {
  GenerateLicenseRequest,
  GenerateLicenseResponse,
  ValidateLicenseRequest,
  ValidateLicenseResponse,
} from '../types/licenses';

/**
 * Generate a license file (.lic) containing key information and metadata
 */
export async function generateLicense(data: GenerateLicenseRequest): Promise<GenerateLicenseResponse> {
  try {
    const response = await apiClient.post<GenerateLicenseResponse>('/licenses/generate', data);
    return response.data;
  } catch (error) {
    throw handleApiError(error);
  }
}

/**
 * Validate a license file
 * Can accept either JSON with base64 content or FormData for file upload
 */
export async function validateLicense(
  data: ValidateLicenseRequest | FormData
): Promise<ValidateLicenseResponse> {
  try {
    let response;
    
    if (data instanceof FormData) {
      // Multipart file upload
      response = await apiClient.post<ValidateLicenseResponse>('/licenses/validate', data, {
        headers: {
          'Content-Type': 'multipart/form-data',
        },
      });
    } else {
      // JSON body with base64 content
      response = await apiClient.post<ValidateLicenseResponse>('/licenses/validate', data);
    }
    
    return response.data;
  } catch (error) {
    throw handleApiError(error);
  }
}

