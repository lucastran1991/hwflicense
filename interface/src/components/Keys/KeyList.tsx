'use client';

import { useState, useEffect } from 'react';
import {
  Box,
  SimpleGrid,
  Text,
  useToast,
  Spinner,
  Center,
  Alert,
  AlertIcon,
  Button,
  HStack,
} from '@chakra-ui/react';
import { KeyCard } from './KeyCard';
import type { KeyInfo } from '@/lib/types/keys';
import { listKeys, downloadKey, refreshKey, removeKey } from '@/lib/api/keys';

export function KeyList() {
  const [keys, setKeys] = useState<KeyInfo[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const toast = useToast();

  const fetchKeys = async () => {
    try {
      setIsLoading(true);
      setError(null);
      const response = await listKeys();
      setKeys(response.keys);
    } catch (err: any) {
      setError(err.message || 'Failed to fetch keys');
      toast({
        title: 'Error loading keys',
        description: err.message || 'An error occurred',
        status: 'error',
        duration: 5000,
        isClosable: true,
      });
    } finally {
      setIsLoading(false);
    }
  };

  useEffect(() => {
    fetchKeys();
  }, []);

  const handleDownload = async (keyId: string) => {
    try {
      await downloadKey(keyId);
      toast({
        title: 'Key downloaded',
        description: `Key ${keyId.substring(0, 8)}... has been downloaded`,
        status: 'success',
        duration: 3000,
        isClosable: true,
      });
    } catch (err: any) {
      toast({
        title: 'Error downloading key',
        description: err.message || 'An error occurred',
        status: 'error',
        duration: 5000,
        isClosable: true,
      });
    }
  };

  const handleRefresh = async (keyId: string) => {
    try {
      // Refresh with 1 year (31536000 seconds)
      await refreshKey(keyId, { expires_in_seconds: 31536000 });
      toast({
        title: 'Key refreshed',
        description: `Key ${keyId.substring(0, 8)}... expiry has been extended`,
        status: 'success',
        duration: 3000,
        isClosable: true,
      });
      // Reload keys
      await fetchKeys();
    } catch (err: any) {
      toast({
        title: 'Error refreshing key',
        description: err.message || 'An error occurred',
        status: 'error',
        duration: 5000,
        isClosable: true,
      });
    }
  };

  const handleRevoke = async (keyId: string) => {
    if (!confirm(`Are you sure you want to revoke key ${keyId.substring(0, 8)}...? This action cannot be undone.`)) {
      return;
    }

    try {
      await removeKey(keyId);
      toast({
        title: 'Key revoked',
        description: `Key ${keyId.substring(0, 8)}... has been revoked`,
        status: 'success',
        duration: 3000,
        isClosable: true,
      });
      // Reload keys
      await fetchKeys();
    } catch (err: any) {
      toast({
        title: 'Error revoking key',
        description: err.message || 'An error occurred',
        status: 'error',
        duration: 5000,
        isClosable: true,
      });
    }
  };

  if (isLoading) {
    return (
      <Center py={8}>
        <Spinner size="xl" />
      </Center>
    );
  }

  if (error && keys.length === 0) {
    return (
      <Box>
        <Alert status="error" mb={4}>
          <AlertIcon />
          {error}
        </Alert>
        <Button onClick={fetchKeys} colorScheme="blue">
          Retry
        </Button>
      </Box>
    );
  }

  return (
    <Box>
      <HStack justify="space-between" mb={4}>
        <Text fontSize="lg" fontWeight="semibold">
          Keys ({keys.length})
        </Text>
        <Button onClick={fetchKeys} size="sm" colorScheme="blue">
          Refresh
        </Button>
      </HStack>

      {keys.length === 0 ? (
        <Alert status="info">
          <AlertIcon />
          No keys found. Create a new key to get started.
        </Alert>
      ) : (
        <SimpleGrid columns={{ base: 1, md: 2, lg: 3 }} spacing={4}>
          {keys.map((key) => (
            <KeyCard
              key={key.key_id}
              keyData={key}
              onDownload={handleDownload}
              onRefresh={handleRefresh}
              onRevoke={handleRevoke}
            />
          ))}
        </SimpleGrid>
      )}
    </Box>
  );
}

