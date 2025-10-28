'use client';

import { useState, useEffect } from 'react';
import { useParams, useRouter } from 'next/navigation';
import { apiClient } from '@/lib/api-client';
import {
  Box,
  Heading,
  Button,
  Card,
  CardHeader,
  CardBody,
  Grid,
  Badge,
  Spinner,
  Text,
  Alert,
  AlertIcon,
  Code,
  Flex,
  Divider,
  IconButton,
} from '@chakra-ui/react';
import { ArrowBackIcon, DownloadIcon, DeleteIcon, CheckIcon, WarningIcon } from '@chakra-ui/icons';
import { useToast } from '@chakra-ui/react';

export default function SiteDetailPage() {
  const params = useParams();
  const router = useRouter();
  const siteId = params.id as string;
  const toast = useToast();
  
  const [site, setSite] = useState<any>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [licenseData, setLicenseData] = useState<any>(null);
  const [fingerprint, setFingerprint] = useState<any>(null);

  useEffect(() => {
    loadSiteDetails();
  }, [siteId]);

  const loadSiteDetails = async () => {
    try {
      const response = await apiClient.getSite(siteId);
      const siteData = response.data.license;
      setSite(siteData);
      
      if (siteData.license_data) {
        try {
          const parsed = JSON.parse(siteData.license_data);
          setLicenseData(parsed);
        } catch (e) {
          console.error('Failed to parse license_data:', e);
        }
      }
      
      if (siteData.fingerprint) {
        try {
          const parsed = JSON.parse(siteData.fingerprint);
          setFingerprint(parsed);
        } catch (e) {
          setFingerprint({});
        }
      }
    } catch (err: any) {
      setError(err.response?.data?.error || 'Failed to load site details');
    } finally {
      setLoading(false);
    }
  };

  const handleDownload = async () => {
    try {
      const response = await apiClient.getSite(siteId);
      const licenseDataStr = response.data.license?.license_data;
      
      if (!licenseDataStr) {
        toast({
          title: 'Error',
          description: 'No license data available',
          status: 'error',
          duration: 3000,
        });
        return;
      }

      const blob = new Blob([licenseDataStr], { type: 'application/json' });
      const url = window.URL.createObjectURL(blob);
      const a = document.createElement('a');
      a.href = url;
      a.download = `site_${siteId}.lic`;
      document.body.appendChild(a);
      a.click();
      document.body.removeChild(a);
      window.URL.revokeObjectURL(url);
      toast({
        title: 'Success',
        description: 'License file downloaded',
        status: 'success',
        duration: 3000,
      });
    } catch (err) {
      toast({
        title: 'Error',
        description: 'Failed to download license file',
        status: 'error',
        duration: 3000,
      });
    }
  };

  const handleRevoke = async () => {
    if (!confirm(`Are you sure you want to revoke site ${siteId}?`)) return;
    
    try {
      await apiClient.deleteSite(siteId);
      router.push('/dashboard/sites');
      toast({
        title: 'Success',
        description: 'Site revoked successfully',
        status: 'success',
        duration: 3000,
      });
    } catch (err: any) {
      toast({
        title: 'Error',
        description: err.response?.data?.error || 'Failed to revoke site',
        status: 'error',
        duration: 3000,
      });
    }
  };

  if (loading) {
    return (
      <Box display="flex" justifyContent="center" alignItems="center" minH="400px">
        <Spinner size="xl" />
      </Box>
    );
  }

  if (error) {
    return (
      <Box maxW="7xl" mx="auto" py={8} px={4}>
        <Alert status="error">
          <AlertIcon />
          {error}
        </Alert>
      </Box>
    );
  }

  if (!site) {
    return (
      <Box maxW="7xl" mx="auto" py={8} px={4}>
        <Alert status="info">
          <AlertIcon />
          Site not found
        </Alert>
      </Box>
    );
  }

  return (
    <Box maxW="7xl" mx="auto" py={6} px={4}>
      <Flex justify="space-between" align="center" mb={6}>
        <Box>
          <Button
            variant="link"
            colorScheme="blue"
            leftIcon={<ArrowBackIcon />}
            onClick={() => router.push('/dashboard/sites')}
            mb={2}
          >
            Back to Sites
          </Button>
          <Heading size="xl">Site Details</Heading>
          <Text color="gray.600" mt={1}>{siteId}</Text>
        </Box>
        <Box>
          <Button
            leftIcon={<DownloadIcon />}
            colorScheme="blue"
            onClick={handleDownload}
            mr={2}
          >
            Download License
          </Button>
          {site.status !== 'revoked' && (
            <Button
              leftIcon={<DeleteIcon />}
              colorScheme="red"
              onClick={handleRevoke}
            >
              Revoke Site
            </Button>
          )}
        </Box>
      </Flex>

      <Grid templateColumns={{ base: '1fr', md: 'repeat(2, 1fr)' }} gap={6} mb={6}>
        <Card>
          <CardHeader>
            <Heading size="md">Basic Information</Heading>
          </CardHeader>
          <CardBody>
            <Box mb={4}>
              <Text fontSize="sm" color="gray.500" mb={1}>Site ID</Text>
              <Text fontWeight="medium">{site.site_id}</Text>
            </Box>
            <Divider mb={4} />
            <Box mb={4}>
              <Text fontSize="sm" color="gray.500" mb={1}>Status</Text>
              <Badge
                colorScheme={
                  site.status === 'active' ? 'green' :
                  site.status === 'revoked' ? 'red' : 'gray'
                }
              >
                {site.status}
              </Badge>
            </Box>
            <Divider mb={4} />
            <Box mb={4}>
              <Text fontSize="sm" color="gray.500" mb={1}>Issued At</Text>
              <Text>{new Date(site.issued_at).toLocaleString()}</Text>
            </Box>
            <Divider mb={4} />
            <Box>
              <Text fontSize="sm" color="gray.500" mb={1}>Last Seen</Text>
              <Text>{site.last_seen ? new Date(site.last_seen).toLocaleString() : 'Never'}</Text>
            </Box>
          </CardBody>
        </Card>

        <Card>
          <CardHeader>
            <Heading size="md">Fingerprint</Heading>
          </CardHeader>
          <CardBody>
            {fingerprint && Object.keys(fingerprint).length > 0 ? (
              Object.entries(fingerprint).map(([key, value], index) => (
                <Box key={key} mb={index < Object.keys(fingerprint).length - 1 ? 4 : 0}>
                  <Text fontSize="sm" color="gray.500" mb={1} textTransform="capitalize">
                    {key.replace('_', ' ')}
                  </Text>
                  <Text>{String(value)}</Text>
                  {index < Object.keys(fingerprint).length - 1 && <Divider mt={4} />}
                </Box>
              ))
            ) : (
              <Text fontSize="sm" color="gray.500">No fingerprint information available</Text>
            )}
          </CardBody>
        </Card>
      </Grid>

      <Card mb={6}>
        <CardHeader>
          <Heading size="md">License Data</Heading>
        </CardHeader>
        <CardBody>
          {licenseData ? (
            <Code p={4} borderRadius="md" display="block" whiteSpace="pre-wrap" overflowX="auto">
              {JSON.stringify(licenseData, null, 2)}
            </Code>
          ) : (
            <Text fontSize="sm" color="gray.500">No license data available</Text>
          )}
        </CardBody>
      </Card>

      <Card>
        <CardHeader>
          <Heading size="md">Signature</Heading>
        </CardHeader>
        <CardBody>
          <Code p={4} borderRadius="md" display="block" wordBreak="break-all">
            {site.signature}
          </Code>
          {site.signature ? (
            <Flex align="center" mt={2} color="green.600">
              <CheckIcon mr={2} />
              <Text fontSize="sm">Valid ECDSA signature</Text>
            </Flex>
          ) : (
            <Flex align="center" mt={2} color="yellow.600">
              <WarningIcon mr={2} />
              <Text fontSize="sm">Placeholder signature</Text>
            </Flex>
          )}
        </CardBody>
      </Card>
    </Box>
  );
}
