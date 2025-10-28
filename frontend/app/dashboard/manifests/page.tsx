'use client';

import { useState, useEffect } from 'react';
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
  Badge,
  Spinner,
  Code,
  VStack,
  Text,
  IconButton,
  useToast,
} from '@chakra-ui/react';
import { ViewIcon, DownloadIcon, ExternalLinkIcon, CloseIcon } from '@chakra-ui/icons';

export default function ManifestsPage() {
  const [manifests, setManifests] = useState<any[]>([]);
  const [loading, setLoading] = useState(true);
  const [period, setPeriod] = useState('');
  const [selectedManifest, setSelectedManifest] = useState<any>(null);
  const { isOpen: isDetailOpen, onOpen: onDetailOpen, onClose: onDetailClose } = useDisclosure();
  const { isOpen: isGenerateOpen, onOpen: onGenerateOpen, onClose: onGenerateClose } = useDisclosure();
  const toast = useToast();

  useEffect(() => {
    loadManifests();
  }, []);

  const loadManifests = async () => {
    setLoading(true);
    try {
      const response = await apiClient.listManifests();
      setManifests(response.data.manifests || []);
    } catch (error) {
      console.error('Failed to load manifests:', error);
      toast({
        title: 'Error',
        description: 'Failed to load manifests',
        status: 'error',
        duration: 3000,
      });
    } finally {
      setLoading(false);
    }
  };

  const handleGenerate = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      await apiClient.generateManifest(period);
      setPeriod('');
      onGenerateClose();
      loadManifests();
      toast({
        title: 'Success',
        description: 'Manifest generated successfully',
        status: 'success',
        duration: 3000,
      });
    } catch (error: any) {
      toast({
        title: 'Error',
        description: error.response?.data?.error || 'Failed to generate manifest',
        status: 'error',
        duration: 3000,
      });
    }
  };

  const handleSendToAStack = async (manifestId: string) => {
    const astackEndpoint = prompt('Enter A-Stack endpoint:', 'http://localhost:8081/api/manifests/receive');
    if (!astackEndpoint) return;

    try {
      await apiClient.sendManifest(manifestId, astackEndpoint);
      loadManifests();
      toast({
        title: 'Success',
        description: 'Manifest successfully sent to A-Stack!',
        status: 'success',
        duration: 3000,
      });
    } catch (error: any) {
      toast({
        title: 'Error',
        description: error.response?.data?.error || 'Failed to send manifest to A-Stack',
        status: 'error',
        duration: 3000,
      });
    }
  };

  const handleViewDetails = async (manifestId: string) => {
    try {
      const response = await apiClient.getManifest(manifestId);
      setSelectedManifest(response.data.manifest);
      onDetailOpen();
    } catch (error) {
      toast({
        title: 'Error',
        description: 'Failed to load manifest details',
        status: 'error',
        duration: 3000,
      });
    }
  };

  const handleDownload = async (manifestId: string) => {
    try {
      const response = await apiClient.downloadManifest(manifestId);
      const blob = new Blob([response.data], { type: 'application/json' });
      const url = window.URL.createObjectURL(blob);
      const a = document.createElement('a');
      a.href = url;
      a.download = `manifest_${manifestId}.json`;
      document.body.appendChild(a);
      a.click();
      document.body.removeChild(a);
      window.URL.revokeObjectURL(url);
      toast({
        title: 'Success',
        description: 'Manifest downloaded',
        status: 'success',
        duration: 3000,
      });
    } catch (error) {
      toast({
        title: 'Error',
        description: 'Failed to download manifest',
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
    <>
      <Modal isOpen={isGenerateOpen} onClose={onGenerateClose} size="md">
        <ModalOverlay />
        <ModalContent>
          <form onSubmit={handleGenerate}>
            <ModalHeader>Generate Usage Manifest</ModalHeader>
            <ModalCloseButton />
            <ModalBody>
              <FormControl isRequired>
                <FormLabel>Period (YYYY-MM format)</FormLabel>
                <Input
                  type="text"
                  value={period}
                  onChange={(e) => setPeriod(e.target.value)}
                  placeholder="2024-01"
                  pattern="\d{4}-\d{2}"
                />
              </FormControl>
            </ModalBody>
            <ModalFooter>
              <Button type="submit" colorScheme="blue" mr={3}>
                Generate
              </Button>
              <Button variant="ghost" onClick={onGenerateClose}>
                Cancel
              </Button>
            </ModalFooter>
          </form>
        </ModalContent>
      </Modal>

      <Modal isOpen={isDetailOpen} onClose={onDetailClose} size="xl">
        <ModalOverlay />
        <ModalContent maxH="80vh">
          <ModalHeader>Manifest Details</ModalHeader>
          <ModalCloseButton />
          <ModalBody>
            {selectedManifest && (
              <>
                <Box bg="gray.50" p={4} borderRadius="md" overflow="auto" mb={4}>
                  <Code display="block" whiteSpace="pre-wrap" fontSize="sm">
                    {JSON.stringify(JSON.parse(selectedManifest.manifest_data), null, 2)}
                  </Code>
                </Box>
                <Box mb={4}>
                  <Text fontSize="sm" fontWeight="medium" color="gray.700" mb={1}>
                    Signature:
                  </Text>
                  <Code fontSize="xs" whiteSpace="pre-wrap" wordBreak="break-all">
                    {selectedManifest.signature}
                  </Code>
                </Box>
              </>
            )}
          </ModalBody>
          <ModalFooter>
            {selectedManifest && (
              <Button
                leftIcon={<DownloadIcon />}
                colorScheme="blue"
                onClick={() => handleDownload(selectedManifest.id)}
                mr={2}
              >
                Download
              </Button>
            )}
            <Button onClick={onDetailClose}>Close</Button>
          </ModalFooter>
        </ModalContent>
      </Modal>

      <Box maxW="7xl" mx="auto" py={6} px={4}>
        <Box display="flex" justifyContent="space-between" alignItems="center" mb={6}>
          <Heading size="xl">Usage Manifests</Heading>
          <Button colorScheme="blue" onClick={onGenerateOpen}>
            Generate Manifest
          </Button>
        </Box>

        <Box bg="white" shadow="md" rounded="lg" overflow="hidden">
          <Table variant="simple">
            <Thead bg="gray.50">
              <Tr>
                <Th>Period</Th>
                <Th>Created At</Th>
                <Th>Status</Th>
                <Th>Actions</Th>
              </Tr>
            </Thead>
            <Tbody>
              {manifests.length === 0 ? (
                <Tr>
                  <Td colSpan={4} textAlign="center" color="gray.500">
                    No manifests generated yet
                  </Td>
                </Tr>
              ) : (
                manifests.map((manifest) => (
                  <Tr key={manifest.id}>
                    <Td fontWeight="medium">{manifest.period}</Td>
                    <Td>{new Date(manifest.created_at).toLocaleString()}</Td>
                    <Td>
                      <Badge colorScheme={manifest.sent_to_astack ? 'green' : 'yellow'}>
                        {manifest.sent_to_astack ? 'Sent' : 'Pending'}
                      </Badge>
                    </Td>
                    <Td>
                      <Button
                        size="sm"
                        leftIcon={<ViewIcon />}
                        onClick={() => handleViewDetails(manifest.id)}
                        mr={2}
                        colorScheme="blue"
                        variant="outline"
                      >
                        View
                      </Button>
                      <Button
                        size="sm"
                        leftIcon={<DownloadIcon />}
                        onClick={() => handleDownload(manifest.id)}
                        mr={2}
                        colorScheme="green"
                        variant="outline"
                      >
                        Download
                      </Button>
                      {!manifest.sent_to_astack && (
                        <Button
                          size="sm"
                          leftIcon={<ExternalLinkIcon />}
                          onClick={() => handleSendToAStack(manifest.id)}
                          colorScheme="orange"
                          variant="outline"
                        >
                          Send to A-Stack
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
    </>
  );
}

