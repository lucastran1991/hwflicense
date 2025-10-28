'use client';

import { useState } from 'react';
import {
  Box,
  Button,
  Card,
  CardBody,
  Code,
  Flex,
  Heading,
  Text,
  VStack,
  HStack,
  Grid,
  GridItem,
  Badge,
  Alert,
  AlertIcon,
  Spinner,
  Textarea,
  UnorderedList,
  ListItem,
  Divider,
} from '@chakra-ui/react';
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

  const getMethodColor = (method: string) => {
    switch (method) {
      case 'GET': return 'blue';
      case 'POST': return 'green';
      case 'DELETE': return 'red';
      default: return 'gray';
    }
  };

  return (
    <Box minH="100vh" bg="gray.50" py={8} px={4}>
      <Box maxW="7xl" mx="auto">
        <VStack align="stretch" spacing={6} mb={8}>
          <Heading size="xl" mb={2}>API Testing & Documentation</Heading>
          <Text color="gray.600">Test all 18 backend API endpoints with live examples</Text>
          
          {/* Auth Token Display */}
          <Alert status={token ? 'success' : 'warning'} borderRadius="md">
            <AlertIcon />
            <Flex w="100%" justify="space-between" align="center">
              <Box>
                <Text fontWeight="semibold" mb={1}>Authentication</Text>
                {token ? (
                  <Flex align="center" gap={2}>
                    <Text fontSize="sm">Token:</Text>
                    <Code bg="green.100" px={2} py={1}>{token.substring(0, 20)}...</Code>
                  </Flex>
                ) : (
                  <Text fontSize="sm">No token. Login first to access protected endpoints.</Text>
                )}
              </Box>
              {!token && (
                <Button
                  size="sm"
                  colorScheme="blue"
                  onClick={() => {
                    const loginEndpoint = API_ENDPOINTS.find(e => e.path === '/api/auth/login');
                    if (loginEndpoint) setSelectedEndpoint(loginEndpoint);
                  }}
                >
                  Login Now
                </Button>
              )}
            </Flex>
          </Alert>
        </VStack>

        <Grid templateColumns={{ base: '1fr', lg: 'repeat(2, 1fr)' }} gap={6}>
          {/* Left: API List */}
          <VStack spacing={4} align="stretch">
            <Card>
              <CardBody>
                <Heading size="md" mb={4}>All API Endpoints (18)</Heading>
                <Box maxH="600px" overflowY="auto">
                  {API_ENDPOINTS.map((endpoint, index) => (
                    <Box
                      key={index}
                      p={3}
                      borderWidth="1px"
                      borderRadius="md"
                      cursor="pointer"
                      bg={selectedEndpoint?.path === endpoint.path ? 'indigo.50' : 'white'}
                      borderColor={selectedEndpoint?.path === endpoint.path ? 'indigo.500' : 'gray.200'}
                      onClick={() => {
                        setSelectedEndpoint(endpoint);
                        setRequestBody(endpoint.example ? JSON.stringify(endpoint.example, null, 2) : '');
                        setResponse(null);
                      }}
                      mb={2}
                      _hover={{ borderColor: 'gray.300' }}
                    >
                      <Flex align="center" gap={2} mb={1}>
                        <Badge colorScheme={getMethodColor(endpoint.method)}>{endpoint.method}</Badge>
                        <Text fontSize="xs" color="gray.600">{endpoint.auth ? '🔒 Auth' : '🌐 Public'}</Text>
                      </Flex>
                      <Code fontSize="sm">{endpoint.path}</Code>
                      <Text fontSize="xs" color="gray.600" mt={1}>{endpoint.description}</Text>
                    </Box>
                  ))}
                </Box>
              </CardBody>
            </Card>
          </VStack>

          {/* Right: Test Panel */}
          <VStack spacing={4} align="stretch">
            {selectedEndpoint && (
              <>
                {/* Endpoint Details */}
                <Card>
                  <CardBody>
                    <Flex justify="space-between" align="flex-start" mb={4}>
                      <Box>
                        <Flex align="center" gap={2} mb={2}>
                          <Badge colorScheme={getMethodColor(selectedEndpoint.method)}>
                            {selectedEndpoint.method}
                          </Badge>
                          <Code fontSize="sm">{selectedEndpoint.path}</Code>
                        </Flex>
                        <Text color="gray.700">{selectedEndpoint.description}</Text>
                      </Box>
                    </Flex>

                    {/* Parameters */}
                    {selectedEndpoint.params && selectedEndpoint.params.length > 0 && (
                      <Box mb={4}>
                        <Text fontWeight="semibold" fontSize="sm" mb={2}>Query Parameters:</Text>
                        <UnorderedList spacing={1} fontSize="xs">
                          {selectedEndpoint.params.map((param, idx) => (
                            <ListItem key={idx} color="gray.600">
                              <Code>{param.name}</Code> ({param.type}) {param.required ? '(required)' : '(optional)'} - {param.description}
                            </ListItem>
                          ))}
                        </UnorderedList>
                      </Box>
                    )}

                    {/* Body Parameters */}
                    {selectedEndpoint.body && selectedEndpoint.body.length > 0 && (
                      <Box mb={4}>
                        <Text fontWeight="semibold" fontSize="sm" mb={2}>Request Body:</Text>
                        <Button
                          size="xs"
                          colorScheme="indigo"
                          variant="ghost"
                          onClick={addExampleToEditor}
                          mb={2}
                        >
                          📋 Load Example
                        </Button>
                        <Textarea
                          value={requestBody}
                          onChange={(e) => setRequestBody(e.target.value)}
                          fontSize="xs"
                          fontFamily="mono"
                          rows={6}
                          placeholder="Enter JSON request body"
                        />
                      </Box>
                    )}

                    {/* Test Button */}
                    <Button
                      onClick={() => handleTest(selectedEndpoint)}
                      colorScheme="indigo"
                      isDisabled={loading || (selectedEndpoint.auth && !token)}
                      w="100%"
                    >
                      {loading ? 'Testing...' : '🚀 Test Endpoint'}
                    </Button>

                    {selectedEndpoint.auth && !token && (
                      <Alert status="warning" borderRadius="md">
                        <AlertIcon />
                        <Text fontSize="xs">
                          This endpoint requires authentication. Please login first.
                        </Text>
                      </Alert>
                    )}
                  </CardBody>
                </Card>

                {/* Response Display */}
                {response && (
                  <Card>
                    <CardBody>
                      <Flex justify="space-between" align="center" mb={4}>
                        <Heading size="md">Response</Heading>
                        <Badge colorScheme={response.status >= 200 && response.status < 300 ? 'green' : 'red'}>
                          Status: {response.status}
                        </Badge>
                      </Flex>
                      <Box
                        bg="gray.50"
                        p={4}
                        borderRadius="md"
                        overflow="auto"
                        maxH="400px"
                      >
                        <Code fontSize="xs" display="block" whiteSpace="pre-wrap">
                          {JSON.stringify(response.data || response.error, null, 2)}
                        </Code>
                      </Box>
                    </CardBody>
                  </Card>
                )}

                {error && (
                  <Alert status="error" borderRadius="md">
                    <AlertIcon />
                    <Box>
                      <Text fontWeight="semibold">Error:</Text>
                      <Text fontSize="xs">{error}</Text>
                    </Box>
                  </Alert>
                )}
              </>
            )}

            {!selectedEndpoint && (
              <Card>
                <CardBody>
                  <Text textAlign="center" color="gray.500">
                    Select an API endpoint from the list to start testing
                  </Text>
                </CardBody>
              </Card>
            )}
          </VStack>
        </Grid>

        {/* Quick Stats */}
        <Grid templateColumns={{ base: '1fr', md: 'repeat(4, 1fr)' }} gap={4} mt={8}>
          <Card>
            <CardBody>
              <Text fontSize="2xl" fontWeight="bold" color="indigo.600">{API_ENDPOINTS.length}</Text>
              <Text fontSize="sm" color="gray.600">Total Endpoints</Text>
            </CardBody>
          </Card>
          <Card>
            <CardBody>
              <Text fontSize="2xl" fontWeight="bold" color="green.600">
                {API_ENDPOINTS.filter(e => e.auth).length}
              </Text>
              <Text fontSize="sm" color="gray.600">Protected</Text>
            </CardBody>
          </Card>
          <Card>
            <CardBody>
              <Text fontSize="2xl" fontWeight="bold" color="blue.600">
                {API_ENDPOINTS.filter(e => !e.auth).length}
              </Text>
              <Text fontSize="sm" color="gray.600">Public</Text>
            </CardBody>
          </Card>
          <Card>
            <CardBody>
              <Text fontSize="2xl" fontWeight="bold" color="purple.600">
                {API_ENDPOINTS.filter(e => e.method === 'POST').length}
              </Text>
              <Text fontSize="sm" color="gray.600">POST Methods</Text>
            </CardBody>
          </Card>
        </Grid>
      </Box>
    </Box>
  );
}

