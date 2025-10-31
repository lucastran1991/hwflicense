// API Client utility for communicating with Go KMS backend
import axios, { AxiosInstance, AxiosError } from 'axios';

const API_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';

// Create axios instance with default config
export const apiClient: AxiosInstance = axios.create({
  baseURL: API_URL,
  headers: {
    'Content-Type': 'application/json',
  },
  timeout: 30000, // 30 seconds
});

// Error handler
export function handleApiError(error: unknown): { message: string; status?: number } {
  if (axios.isAxiosError(error)) {
    const axiosError = error as AxiosError<{ error?: string; message?: string }>;
    return {
      message: axiosError.response?.data?.error || axiosError.response?.data?.message || axiosError.message || 'An error occurred',
      status: axiosError.response?.status,
    };
  }
  return {
    message: error instanceof Error ? error.message : 'An unknown error occurred',
  };
}

// Health check
export async function checkHealth(): Promise<{ status: string }> {
  const response = await apiClient.get<{ status: string }>('/health');
  return response.data;
}

