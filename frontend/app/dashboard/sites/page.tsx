'use client';

import { useState, useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { apiClient } from '@/lib/api-client';
import {
  Box,
  Heading,
  Button,
  Table,
  Thead,
  Tbody,
  Tr,
  Th,
  Td,
  Modal,
  ModalOverlay,
  ModalContent,
  ModalHeader,
  ModalFooter,
  ModalBody,
  ModalCloseButton,
  useDisclosure,
  FormControl,
  FormLabel,
  Input,
  Select,
  Grid,
  Badge,
  useToast,
  Spinner,
} from '@chakra-ui/react';
import { CheckIcon, DownloadIcon, ViewIcon, DeleteIcon } from '@chakra-ui/icons';

export default function SitesPage() {
  const [sites, setSites] = useState<any[]>([]);
  const [loading, setLoading] = useState(true);
  const [newSiteId, setNewSiteId] = useState('');
  const [keyType, setKeyType] = useState<'production' | 'dev'>('production');
  const [fingerprint, setFingerprint] = useState({ address: '', dns_suffix: '', deployment_tag: '' });
  const router = useRouter();
  const { isOpen, onOpen, onClose } = useDisclosure();
  const toast = useToast();

  useEffect(() => {
    loadSites();
  }, []);

  const loadSites = async () => {
    setLoading(true);
    try {
      const response = await apiClient.listSites();
      setSites(response.data.sites || []);
    } catch (error) {
      toast({
        title: 'Error',
        description: 'Failed to load sites',
        status: 'error',
        duration: 3000,
      });
    } finally {
      setLoading(false);
    }
  };

  const handleCreate = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      const fingerprintData: any = {};
      if (fingerprint.address) fingerprintData.address = fingerprint.address;
      if (fingerprint.dns_suffix) fingerprintData.dns_suffix = fingerprint.dns_suffix;
      if (fingerprint.deployment_tag) fingerprintData.deployment_tag = fingerprint.deployment_tag;
      
      await apiClient.createSite(newSiteId, Object.keys(fingerprintData).length > 0 ? fingerprintData : undefined);
      onClose();
      setNewSiteId('');
      setFingerprint({ address: '', dns_suffix: '', deployment_tag: '' });
      loadSites();
      toast({
        title: 'Success',
        description: 'Site created successfully',
        status: 'success',
        duration: 3000,
      });
    } catch (error: any) {
      toast({
        title: 'Error',
        description: error.response?.data?.error || 'Failed to create site',
        status: 'error',
        duration: 3000,
      });
    }
  };

  const handleDownload = async (siteId: string) => {
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

  const handleDelete = async (siteId: string) => {
    try {
      await apiClient.deleteSite(siteId);
      loadSites();
      toast({
        title: 'Success',
        description: 'Site revoked successfully',
        status: 'success',
        duration: 3000,
      });
    } catch (error) {
      toast({
        title: 'Error',
        description: 'Failed to revoke site',
        status: 'error',
        duration: 3000,
      });
    }
  };

  if (loading) {
    return (
      <Box display="flex" justifyContent="center" alignItems="center" minH="200px">
        <Spinner size="xl" />
      </Box>
    );
  }

  return (
    <Box maxW="7xl" mx="auto" py={6} px={4}>
      <Box display="flex" justifyContent="space-between" alignItems="center" mb={6}>
        <Heading size="xl">Site Licenses</Heading>
        <Button colorScheme="blue" leftIcon={<CheckIcon />} onClick={onOpen}>
          Create Site License
        </Button>
      </Box>

      <Modal isOpen={isOpen} onClose={onClose} size="xl">
        <ModalOverlay />
        <ModalContent>
          <form onSubmit={handleCreate}>
            <ModalHeader>Create New Site License</ModalHeader>
            <ModalCloseButton />
            <ModalBody>
              <FormControl mb={4} isRequired>
                <FormLabel>Site ID</FormLabel>
                <Input
                  value={newSiteId}
                  onChange={(e) => setNewSiteId(e.target.value)}
                  placeholder="Enter site ID"
                />
              </FormControl>

              <FormControl mb={4}>
                <FormLabel>Fingerprint (Optional)</FormLabel>
                <Grid templateColumns="repeat(3, 1fr)" gap={4}>
                  <Box>
                    <FormLabel fontSize="xs" color="gray.500">Address</FormLabel>
                    <Input
                      value={fingerprint.address}
                      onChange={(e) => setFingerprint({...fingerprint, address: e.target.value})}
                      placeholder="192.168.1.1"
                      size="sm"
                    />
                  </Box>
                  <Box>
                    <FormLabel fontSize="xs" color="gray.500">DNS Suffix</FormLabel>
                    <Input
                      value={fingerprint.dns_suffix}
                      onChange={(e) => setFingerprint({...fingerprint, dns_suffix: e.target.value})}
                      placeholder="hwf.local"
                      size="sm"
                    />
                  </Box>
                  <Box>
                    <FormLabel fontSize="xs" color="gray.500">Deployment Tag</FormLabel>
                    <Input
                      value={fingerprint.deployment_tag}
                      onChange={(e) => setFingerprint({...fingerprint, deployment_tag: e.target.value})}
                      placeholder="production"
                      size="sm"
                    />
                  </Box>
                </Grid>
              </FormControl>

              <FormControl>
                <FormLabel>Key Type</FormLabel>
                <Select value={keyType} onChange={(e) => setKeyType(e.target.value as 'production' | 'dev')}>
                  <option value="production">Production</option>
                  <option value="dev">Development</option>
                </Select>
              </FormControl>
            </ModalBody>

            <ModalFooter>
              <Button variant="ghost" mr={3} onClick={onClose}>
                Cancel
              </Button>
              <Button type="submit" colorScheme="blue">
                Create
              </Button>
            </ModalFooter>
          </form>
        </ModalContent>
      </Modal>

      <Box bg="white" shadow="md" rounded="lg" overflow="hidden">
        <Table variant="simple">
          <Thead bg="gray.50">
            <Tr>
              <Th>Site ID</Th>
              <Th>Status</Th>
              <Th>Issued At</Th>
              <Th>Last Seen</Th>
              <Th>Actions</Th>
            </Tr>
          </Thead>
          <Tbody>
            {sites.length === 0 ? (
              <Tr>
                <Td colSpan={5} textAlign="center" color="gray.500">
                  No sites found
                </Td>
              </Tr>
            ) : (
              sites.map((site) => (
                <Tr key={site.site_id}>
                  <Td fontWeight="medium">{site.site_id}</Td>
                  <Td>
                    <Badge
                      colorScheme={
                        site.status === 'active' ? 'green' :
                        site.status === 'revoked' ? 'red' : 'gray'
                      }
                    >
                      {site.status}
                    </Badge>
                  </Td>
                  <Td>{new Date(site.issued_at).toLocaleString()}</Td>
                  <Td>{site.last_seen ? new Date(site.last_seen).toLocaleString() : 'Never'}</Td>
                  <Td>
                    <Button
                      size="sm"
                      leftIcon={<ViewIcon />}
                      onClick={() => router.push(`/dashboard/sites/${site.site_id}`)}
                      mr={2}
                      colorScheme="blue"
                      variant="outline"
                    >
                      View
                    </Button>
                    <Button
                      size="sm"
                      leftIcon={<DownloadIcon />}
                      onClick={() => handleDownload(site.site_id)}
                      mr={2}
                      colorScheme="green"
                      variant="outline"
                    >
                      Download
                    </Button>
                    {site.status !== 'revoked' && (
                      <Button
                        size="sm"
                        leftIcon={<DeleteIcon />}
                        onClick={() => handleDelete(site.site_id)}
                        colorScheme="red"
                        variant="outline"
                      >
                        Revoke
                      </Button>
                    )}
                  </Td>
                </Tr>
              ))
            )}
          </Tbody>
        </Table>
      </Box>
    </Box>
  );
}
