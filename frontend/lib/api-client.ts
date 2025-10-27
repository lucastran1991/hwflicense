import axios, { AxiosInstance } from 'axios';

const API_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api';

class ApiClient {
  client: AxiosInstance;

  constructor() {
    this.client = axios.create({
      baseURL: API_URL,
      headers: {
        'Content-Type': 'application/json',
      },
    });

    // Request interceptor to add auth token
    this.client.interceptors.request.use((config) => {
      const token = this.getToken();
      if (token) {
        config.headers.Authorization = `Bearer ${token}`;
      }
      return config;
    });

    // Response interceptor for error handling
    this.client.interceptors.response.use(
      (response) => response,
      (error) => {
        if (error.response?.status === 401) {
          // Handle unauthorized
          this.removeToken();
          if (typeof window !== 'undefined') {
            window.location.href = '/login';
          }
        }
        return Promise.reject(error);
      }
    );
  }

  // Token management
  setToken(token: string) {
    if (typeof window !== 'undefined') {
      localStorage.setItem('token', token);
    }
  }

  getToken(): string | null {
    if (typeof window !== 'undefined') {
      return localStorage.getItem('token');
    }
    return null;
  }

  removeToken() {
    if (typeof window !== 'undefined') {
      localStorage.removeItem('token');
    }
  }

  // API methods
  async login(username: string, password: string) {
    const response = await this.client.post('/auth/login', { username, password });
    const token = response.data.token;
    this.setToken(token);
    return response.data;
  }

  logout() {
    this.removeToken();
  }

  // CML endpoints
  async uploadCML(cmlData: string, signature: string, publicKey: string) {
    return this.client.post('/cml/upload', { cml_data: cmlData, signature, public_key: publicKey });
  }

  async getCML(orgId?: string) {
    const params = orgId ? { org_id: orgId } : {};
    return this.client.get('/cml', { params });
  }

  async refreshCML(orgId: string, cmlData: string, signature: string) {
    return this.client.post('/cml/refresh', { cml_data: cmlData, signature }, { params: { org_id: orgId } });
  }

  // Site license endpoints
  async createSite(siteId: string, fingerprint?: Record<string, any>, orgId?: string) {
    const params = orgId ? { org_id: orgId } : {};
    return this.client.post('/sites/create', { site_id: siteId, fingerprint }, { params });
  }

  async listSites(orgId?: string, status?: string, limit: number = 50, offset: number = 0) {
    const params: any = { limit, offset };
    if (orgId) params.org_id = orgId;
    if (status) params.status = status;
    return this.client.get('/sites', { params });
  }

  async getSite(siteId: string) {
    return this.client.get(`/sites/${siteId}`);
  }

  async deleteSite(siteId: string) {
    return this.client.delete(`/sites/${siteId}`);
  }

  async heartbeat(siteId: string) {
    return this.client.post(`/sites/${siteId}/heartbeat`, { timestamp: new Date().toISOString() });
  }

  // Manifest endpoints
  async generateManifest(period: string) {
    return this.client.post('/manifests/generate', { period });
  }

  async listManifests(period?: string) {
    const params = period ? { period } : {};
    return this.client.get('/manifests', { params });
  }

  async getManifest(manifestId: string) {
    return this.client.get(`/manifests/${manifestId}`);
  }

  async downloadManifest(manifestId: string) {
    return this.client.get(`/manifests/${manifestId}/download`, { responseType: 'blob' });
  }

  async sendManifest(manifestId: string, astackEndpoint: string) {
    return this.client.post('/manifests/send', { manifest_id: manifestId, astack_endpoint: astackEndpoint });
  }

  // Ledger endpoints
  async getLedger(orgId: string, limit: number = 100, offset: number = 0) {
    return this.client.get('/ledger', { params: { org_id: orgId, limit, offset } });
  }

  // Validation endpoint (public)
  async validateLicense(license: Record<string, any>, fingerprint?: Record<string, any>) {
    return this.client.post('/validate', { license, fingerprint });
  }
}

export const apiClient = new ApiClient();

