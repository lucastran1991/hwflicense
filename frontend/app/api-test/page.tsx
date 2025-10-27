'use client';

import { useState } from 'react';
import { apiClient } from '@/lib/api-client';

interface APIEndpoint {
  method: string;
  path: string;
  description: string;
  auth: boolean;
  params?: { name: string; type: string; required: boolean; description: string }[];
  body?: { name: string; type: string; required: boolean; description: string }[];
  example?: any;
}

const API_ENDPOINTS: APIEndpoint[] = [
  // Health Check
  {
    method: 'GET',
    path: '/api/health',
    description: 'Health check endpoint',
    auth: false,
  },

  // Authentication
  {
    method: 'POST',
    path: '/api/auth/login',
    description: 'Login and get JWT token',
    auth: false,
    body: [
      { name: 'username', type: 'string', required: true, description: 'Username (default: admin)' },
      { name: 'password', type: 'string', required: true, description: 'Password (default: admin123)' },
    ],
    example: { username: 'admin', password: 'admin123' },
  },

  // CML Management
  {
    method: 'POST',
    path: '/api/cml/upload',
    description: 'Upload Customer Master License',
    auth: true,
    body: [
      { name: 'cml_data', type: 'string', required: true, description: 'CML JSON data' },
      { name: 'signature', type: 'string', required: true, description: 'ECDSA signature' },
      { name: 'public_key', type: 'string', required: true, description: 'Public key for verification' },
    ],
  },
  {
    method: 'GET',
    path: '/api/cml',
    description: 'Get CML information',
    auth: true,
  },
  {
    method: 'POST',
    path: '/api/cml/refresh',
    description: 'Refresh CML with new keys',
    auth: true,
  },

  // Site License Management
  {
    method: 'POST',
    path: '/api/sites/create',
    description: 'Create a new site license',
    auth: true,
    body: [
      { name: 'site_id', type: 'string', required: true, description: 'Unique site identifier' },
      { name: 'fingerprint', type: 'object', required: false, description: 'Fingerprint (address, dns_suffix, deployment_tag)' },
    ],
    example: { site_id: 'site_001', fingerprint: { address: '192.168.1.1' } },
  },
  {
    method: 'GET',
    path: '/api/sites',
    description: 'List all site licenses',
    auth: true,
    params: [
      { name: 'org_id', type: 'string', required: false, description: 'Filter by organization ID' },
      { name: 'status', type: 'string', required: false, description: 'Filter by status (active/revoked)' },
      { name: 'limit', type: 'number', required: false, description: 'Limit results (default: 50)' },
      { name: 'offset', type: 'number', required: false, description: 'Offset for pagination' },
    ],
  },
  {
    method: 'GET',
    path: '/api/sites/:site_id',
    description: 'Get specific site license',
    auth: true,
  },
  {
    method: 'DELETE',
    path: '/api/sites/:site_id',
    description: 'Revoke a site license',
    auth: true,
  },
  {
    method: 'POST',
    path: '/api/sites/:site_id/heartbeat',
    description: 'Update site heartbeat (last_seen)',
    auth: true,
  },

  // License Validation
  {
    method: 'POST',
    path: '/api/validate',
    description: 'Validate a license (public endpoint)',
    auth: false,
    body: [
      { name: 'license', type: 'object', required: true, description: 'License data' },
      { name: 'fingerprint', type: 'object', required: false, description: 'Fingerprint for matching' },
    ],
  },

  // Usage Ledger
  {
    method: 'GET',
    path: '/api/ledger',
    description: 'Get usage ledger entries',
    auth: true,
    params: [
      { name: 'org_id', type: 'string', required: false, description: 'Filter by organization ID' },
      { name: 'limit', type: 'number', required: false, description: 'Limit results' },
      { name: 'offset', type: 'number', required: false, description: 'Offset for pagination' },
    ],
  },

  // Manifest Management
  {
    method: 'POST',
    path: '/api/manifests/generate',
    description: 'Generate usage manifest',
    auth: true,
    body: [
      { name: 'period', type: 'string', required: true, description: 'Period in YYYY-MM format' },
    ],
    example: { period: '2024-01' },
  },
  {
    method: 'GET',
    path: '/api/manifests',
    description: 'List manifests',
    auth: true,
    params: [
      { name: 'period', type: 'string', required: false, description: 'Filter by period' },
    ],
  },
  {
    method: 'GET',
    path: '/api/manifests/:manifest_id',
    description: 'Get specific manifest',
    auth: true,
  },
  {
    method: 'GET',
    path: '/api/manifests/:manifest_id/download',
    description: 'Download manifest as JSON',
    auth: true,
  },
  {
    method: 'POST',
    path: '/api/manifests/send',
    description: 'Send manifest to A-Stack',
    auth: true,
    body: [
      { name: 'manifest_id', type: 'string', required: true, description: 'Manifest ID' },
      { name: 'astack_endpoint', type: 'string', required: true, description: 'A-Stack endpoint URL' },
    ],
    example: { manifest_id: 'xxx', astack_endpoint: 'http://localhost:8081/api/manifests/receive' },
  },
];

export default function APITestPage() {
  const [selectedEndpoint, setSelectedEndpoint] = useState<APIEndpoint | null>(null);
  const [requestBody, setRequestBody] = useState<string>('');
  const [response, setResponse] = useState<any>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [token, setToken] = useState<string>(apiClient.getToken() || '');
  const [testSiteId, setTestSiteId] = useState<string>('site_test_001');

  const getResponseColor = () => {
    if (!response) return '';
    const status = response.status;
    if (status >= 200 && status < 300) return 'text-green-600';
    if (status >= 400) return 'text-red-600';
    return 'text-yellow-600';
  };

  const handleTest = async (endpoint: APIEndpoint) => {
    setLoading(true);
    setError(null);
    setResponse(null);

    try {
      let result;
      const method = endpoint.method.toLowerCase();

      switch (endpoint.path) {
        case '/api/health':
          result = await apiClient.client.get('/health');
          break;
        case '/api/auth/login':
          const loginData = requestBody ? JSON.parse(requestBody) : { username: 'admin', password: 'admin123' };
          result = await apiClient.login(loginData.username, loginData.password);
          if (result.token) setToken(result.token);
          break;
        case '/api/cml':
          result = await apiClient.getCML();
          break;
        case '/api/cml/upload':
          const uploadData = JSON.parse(requestBody);
          result = await apiClient.uploadCML(uploadData.cml_data, uploadData.signature, uploadData.public_key);
          break;
        case '/api/sites/create':
          const createData = requestBody ? JSON.parse(requestBody) : { site_id: testSiteId };
          result = await apiClient.createSite(createData.site_id, createData.fingerprint);
          break;
        case '/api/sites':
          result = await apiClient.listSites();
          break;
        case '/api/sites/:site_id':
          const siteId = prompt('Enter site ID:') || testSiteId;
          result = await apiClient.getSite(siteId);
          break;
        case '/api/sites/:site_id/heartbeat':
          const heartbeatSiteId = prompt('Enter site ID:') || testSiteId;
          result = await apiClient.heartbeat(heartbeatSiteId);
          break;
        case '/api/manifests/generate':
          const manifestData = requestBody ? JSON.parse(requestBody) : { period: '2024-01' };
          result = await apiClient.generateManifest(manifestData.period);
          break;
        case '/api/manifests':
          result = await apiClient.listManifests();
          break;
        case '/api/ledger':
          result = await apiClient.getLedger('default', 100, 0);
          break;
        case '/api/validate':
          const validateData = JSON.parse(requestBody || '{}');
          result = await apiClient.validateLicense(validateData.license, validateData.fingerprint);
          break;
        case '/api/manifests/send':
          const sendData = JSON.parse(requestBody);
          result = await apiClient.sendManifest(sendData.manifest_id, sendData.astack_endpoint);
          break;
        default:
          throw new Error(`Endpoint ${endpoint.path} not implemented in test client`);
      }

      setResponse({ status: 200, data: result.data });
    } catch (err: any) {
      setResponse({
        status: err.response?.status || 500,
        error: err.response?.data || err.message,
      });
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  const addExampleToEditor = () => {
    if (selectedEndpoint?.example) {
      setRequestBody(JSON.stringify(selectedEndpoint.example, null, 2));
    }
  };

  return (
    <div className="min-h-screen bg-gray-50">
      <div className="max-w-7xl mx-auto py-8 px-4 sm:px-6 lg:px-8">
        <div className="mb-8">
          <h1 className="text-3xl font-bold text-gray-900 mb-2">API Testing & Documentation</h1>
          <p className="text-gray-600">Test all 18 backend API endpoints with live examples</p>
          
          {/* Auth Token Display */}
          <div className="mt-4 p-4 bg-blue-50 border border-blue-200 rounded-lg">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm font-medium text-blue-900">Authentication</p>
                <p className="text-xs text-blue-700 mt-1">
                  {token ? (
                    <>
                      ✅ Token: <code className="bg-blue-100 px-2 py-1 rounded">{token.substring(0, 20)}...</code>
                    </>
                  ) : (
                    '❌ No token. Login first to access protected endpoints.'
                  )}
                </p>
              </div>
              {!token && (
                <button
                  onClick={() => {
                    const loginEndpoint = API_ENDPOINTS.find(e => e.path === '/api/auth/login');
                    if (loginEndpoint) setSelectedEndpoint(loginEndpoint);
                  }}
                  className="text-sm bg-blue-600 text-white px-3 py-1 rounded hover:bg-blue-700"
                >
                  Login Now
                </button>
              )}
            </div>
          </div>
        </div>

        <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
          {/* Left: API List */}
          <div className="space-y-4">
            <div className="bg-white shadow rounded-lg p-4">
              <h2 className="text-xl font-semibold mb-4">All API Endpoints (18)</h2>
              <div className="space-y-2 max-h-[600px] overflow-y-auto">
                {API_ENDPOINTS.map((endpoint, index) => (
                  <div
                    key={index}
                    onClick={() => {
                      setSelectedEndpoint(endpoint);
                      setRequestBody(endpoint.example ? JSON.stringify(endpoint.example, null, 2) : '');
                      setResponse(null);
                    }}
                    className={`p-3 border rounded-lg cursor-pointer transition-colors ${
                      selectedEndpoint?.path === endpoint.path
                        ? 'bg-indigo-50 border-indigo-500'
                        : 'border-gray-200 hover:border-gray-300'
                    }`}
                  >
                    <div className="flex items-center gap-2 mb-1">
                      <span className={`text-xs font-semibold px-2 py-1 rounded ${
                        endpoint.method === 'GET' ? 'bg-blue-100 text-blue-800' :
                        endpoint.method === 'POST' ? 'bg-green-100 text-green-800' :
                        endpoint.method === 'DELETE' ? 'bg-red-100 text-red-800' :
                        'bg-gray-100 text-gray-800'
                      }`}>
                        {endpoint.method}
                      </span>
                      <span className="text-xs font-medium text-gray-600">{endpoint.auth ? '🔒 Auth' : '🌐 Public'}</span>
                    </div>
                    <p className="font-mono text-sm">{endpoint.path}</p>
                    <p className="text-xs text-gray-600 mt-1">{endpoint.description}</p>
                  </div>
                ))}
              </div>
            </div>
          </div>

          {/* Right: Test Panel */}
          <div className="space-y-4">
            {selectedEndpoint && (
              <>
                {/* Endpoint Details */}
                <div className="bg-white shadow rounded-lg p-6">
                  <div className="flex items-start justify-between mb-4">
                    <div>
                      <div className="flex items-center gap-2 mb-2">
                        <span className={`text-xs font-semibold px-2 py-1 rounded ${
                          selectedEndpoint.method === 'GET' ? 'bg-blue-100 text-blue-800' :
                          selectedEndpoint.method === 'POST' ? 'bg-green-100 text-green-800' :
                          selectedEndpoint.method === 'DELETE' ? 'bg-red-100 text-red-800' :
                          'bg-gray-100 text-gray-800'
                        }`}>
                          {selectedEndpoint.method}
                        </span>
                        <code className="text-sm font-mono">{selectedEndpoint.path}</code>
                      </div>
                      <p className="text-gray-700">{selectedEndpoint.description}</p>
                    </div>
                  </div>

                  {/* Parameters */}
                  {selectedEndpoint.params && selectedEndpoint.params.length > 0 && (
                    <div className="mb-4">
                      <p className="text-sm font-semibold mb-2">Query Parameters:</p>
                      <ul className="space-y-1 text-xs">
                        {selectedEndpoint.params.map((param, idx) => (
                          <li key={idx} className="text-gray-600">
                            <code>{param.name}</code> ({param.type}) {param.required ? '(required)' : '(optional)'} - {param.description}
                          </li>
                        ))}
                      </ul>
                    </div>
                  )}

                  {/* Body Parameters */}
                  {selectedEndpoint.body && selectedEndpoint.body.length > 0 && (
                    <div className="mb-4">
                      <p className="text-sm font-semibold mb-2">Request Body:</p>
                      <div className="mb-2">
                        <button
                          onClick={addExampleToEditor}
                          className="text-xs text-indigo-600 hover:text-indigo-800"
                        >
                          📋 Load Example
                        </button>
                      </div>
                      <textarea
                        value={requestBody}
                        onChange={(e) => setRequestBody(e.target.value)}
                        className="w-full font-mono text-xs border border-gray-300 rounded p-2"
                        rows={6}
                        placeholder="Enter JSON request body"
                      />
                    </div>
                  )}

                  {/* Test Button */}
                  <button
                    onClick={() => handleTest(selectedEndpoint)}
                    disabled={loading || (selectedEndpoint.auth && !token)}
                    className="w-full bg-indigo-600 text-white px-4 py-2 rounded-md hover:bg-indigo-700 disabled:opacity-50 disabled:cursor-not-allowed"
                  >
                    {loading ? 'Testing...' : '🚀 Test Endpoint'}
                  </button>

                  {selectedEndpoint.auth && !token && (
                    <p className="text-xs text-red-600 mt-2 text-center">
                      ⚠️ This endpoint requires authentication. Please login first.
                    </p>
                  )}
                </div>

                {/* Response Display */}
                {response && (
                  <div className="bg-white shadow rounded-lg p-6">
                    <div className="flex items-center justify-between mb-4">
                      <h3 className="text-lg font-semibold">Response</h3>
                      <span className={`text-sm font-semibold ${getResponseColor()}`}>
                        Status: {response.status}
                      </span>
                    </div>
                    <pre className="bg-gray-50 p-4 rounded text-xs overflow-auto max-h-[400px]">
                      {JSON.stringify(response.data || response.error, null, 2)}
                    </pre>
                  </div>
                )}

                {error && (
                  <div className="bg-red-50 border border-red-200 rounded-lg p-4">
                    <p className="text-sm text-red-800 font-semibold">Error:</p>
                    <p className="text-xs text-red-700 mt-1">{error}</p>
                  </div>
                )}
              </>
            )}

            {!selectedEndpoint && (
              <div className="bg-white shadow rounded-lg p-8 text-center text-gray-500">
                <p>Select an API endpoint from the list to start testing</p>
              </div>
            )}
          </div>
        </div>

        {/* Quick Stats */}
        <div className="mt-8 grid grid-cols-1 md:grid-cols-4 gap-4">
          <div className="bg-white shadow rounded-lg p-4">
            <p className="text-2xl font-bold text-indigo-600">{API_ENDPOINTS.length}</p>
            <p className="text-sm text-gray-600">Total Endpoints</p>
          </div>
          <div className="bg-white shadow rounded-lg p-4">
            <p className="text-2xl font-bold text-green-600">
              {API_ENDPOINTS.filter(e => e.auth).length}
            </p>
            <p className="text-sm text-gray-600">Protected</p>
          </div>
          <div className="bg-white shadow rounded-lg p-4">
            <p className="text-2xl font-bold text-blue-600">
              {API_ENDPOINTS.filter(e => !e.auth).length}
            </p>
            <p className="text-sm text-gray-600">Public</p>
          </div>
          <div className="bg-white shadow rounded-lg p-4">
            <p className="text-2xl font-bold text-purple-600">
              {API_ENDPOINTS.filter(e => e.method === 'POST').length}
            </p>
            <p className="text-sm text-gray-600">POST Methods</p>
          </div>
        </div>
      </div>
    </div>
  );
}

