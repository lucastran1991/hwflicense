'use client';

import { useState, useRef } from 'react';
import {
  Box,
  Button,
  VStack,
  Text,
  Alert,
  AlertIcon,
  useToast,
  Heading,
  Divider,
  HStack,
  Badge,
  Code,
  Textarea,
} from '@chakra-ui/react';
import type { DownloadKeyResponse } from '@/lib/types/keys';
import { registerKey } from '@/lib/api/keys';

export function KeyUpload() {
  const [selectedFile, setSelectedFile] = useState<File | null>(null);
  const [keyData, setKeyData] = useState<DownloadKeyResponse | null>(null);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [jsonText, setJsonText] = useState<string>('');
  const fileInputRef = useRef<HTMLInputElement>(null);
  const toast = useToast();

  const handleFileSelect = (event: React.ChangeEvent<HTMLInputElement>) => {
    const file = event.target.files?.[0];
    if (!file) return;

    if (!file.name.endsWith('.json')) {
      setError('Please select a JSON file');
      return;
    }

    setSelectedFile(file);
    setError(null);
    setKeyData(null);

    const reader = new FileReader();
    reader.onload = (e) => {
      try {
        const text = e.target?.result as string;
        setJsonText(text);
        const parsed = JSON.parse(text) as DownloadKeyResponse;
        
        // Validate structure
        if (!parsed.key_id || !parsed.key_type) {
          setError('Invalid key JSON file: missing required fields');
          return;
        }

        setKeyData(parsed);
      } catch (err: any) {
        setError(`Failed to parse JSON: ${err.message}`);
      }
    };
    reader.onerror = () => {
      setError('Failed to read file');
    };
    reader.readAsText(file);
  };

  const handleJsonTextChange = (event: React.ChangeEvent<HTMLTextAreaElement>) => {
    const text = event.target.value;
    setJsonText(text);
    setError(null);
    setKeyData(null);

    if (!text.trim()) {
      return;
    }

    try {
      const parsed = JSON.parse(text) as DownloadKeyResponse;
      
      // Validate structure
      if (!parsed.key_id || !parsed.key_type) {
        setError('Invalid key JSON: missing required fields');
        return;
      }

      setKeyData(parsed);
    } catch (err: any) {
      setError(`Invalid JSON: ${err.message}`);
    }
  };

  const handleRegister = async () => {
    if (!keyData) {
      setError('No key data to register');
      return;
    }

    setIsLoading(true);
    setError(null);

    try {
      // Extract key material based on key type
      let keyMaterial: string | undefined;
      if (keyData.key_type === 'symmetric' && keyData.symmetric_key) {
        keyMaterial = keyData.symmetric_key;
      } else if (keyData.key_type === 'asymmetric' && keyData.private_key) {
        keyMaterial = keyData.private_key;
      } else {
        setError('Key material not found in JSON file');
        setIsLoading(false);
        return;
      }

      // Calculate expires_in_seconds from expires_at if available
      // If expires_at is in the past, use default 1 year
      let expiresInSeconds: number | undefined;
      if (keyData.expires_at) {
        const expiresAt = new Date(keyData.expires_at);
        const now = new Date();
        const diffSeconds = Math.floor((expiresAt.getTime() - now.getTime()) / 1000);
        if (diffSeconds > 0) {
          expiresInSeconds = diffSeconds;
        } else {
          // If expired, default to 1 year from now
          expiresInSeconds = 365 * 24 * 60 * 60;
        }
      }

      // Register the key
      const response = await registerKey({
        key_type: keyData.key_type,
        key_material: keyMaterial,
        expires_in_seconds: expiresInSeconds,
      });

      toast({
        title: 'Key registered successfully',
        description: `New key ID: ${response.key_id}`,
        status: 'success',
        duration: 5000,
        isClosable: true,
      });

      // Reset form
      setSelectedFile(null);
      setKeyData(null);
      setJsonText('');
      if (fileInputRef.current) {
        fileInputRef.current.value = '';
      }
    } catch (err: any) {
      const errorMessage = err.message || 'Failed to register key';
      setError(errorMessage);
      toast({
        title: 'Error registering key',
        description: errorMessage,
        status: 'error',
        duration: 5000,
        isClosable: true,
      });
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <VStack align="stretch" spacing={6}>
      <Box>
        <Heading size="md" mb={2}>
          Upload Key from JSON
        </Heading>
        <Text fontSize="sm" color="gray.500">
          Upload a key JSON file (downloaded from the system) to register it in the system.
        </Text>
      </Box>

      <Divider />

      {/* File Upload */}
      <Box>
        <Text fontWeight="semibold" mb={2}>
          Option 1: Upload JSON File
        </Text>
        <input
          ref={fileInputRef}
          type="file"
          accept=".json"
          onChange={handleFileSelect}
          style={{ display: 'none' }}
        />
        <Button
          onClick={() => fileInputRef.current?.click()}
          colorScheme="blue"
          variant="outline"
          size="sm"
        >
          Choose JSON File
        </Button>
        {selectedFile && (
          <Text fontSize="sm" color="gray.600" mt={2}>
            Selected: {selectedFile.name}
          </Text>
        )}
      </Box>

      {/* Or JSON Text */}
      <Box>
        <Text fontWeight="semibold" mb={2}>
          Option 2: Paste JSON Text
        </Text>
        <Textarea
          value={jsonText}
          onChange={handleJsonTextChange}
          placeholder="Paste key JSON here..."
          rows={8}
          fontFamily="mono"
          fontSize="xs"
        />
      </Box>

      {/* Error Display */}
      {error && (
        <Alert status="error">
          <AlertIcon />
          {error}
        </Alert>
      )}

      {/* Key Info Display */}
      {keyData && !error && (
        <Box borderWidth="1px" borderRadius="md" p={4} bg="gray.50" _dark={{ bg: 'gray.700' }}>
          <Heading size="sm" mb={4}>
            Key Information
          </Heading>
          <VStack align="stretch" spacing={3}>
            <HStack>
              <Text fontWeight="semibold" minW="120px">
                Key ID:
              </Text>
              <Code>{keyData.key_id}</Code>
            </HStack>
            <HStack>
              <Text fontWeight="semibold" minW="120px">
                Key Type:
              </Text>
              <Badge colorScheme={keyData.key_type === 'symmetric' ? 'blue' : 'purple'}>
                {keyData.key_type}
              </Badge>
            </HStack>
            {keyData.public_key && (
              <HStack align="start">
                <Text fontWeight="semibold" minW="120px">
                  Public Key:
                </Text>
                <Code fontSize="xs" wordBreak="break-all">
                  {keyData.public_key.substring(0, 64)}...
                </Code>
              </HStack>
            )}
            <HStack>
              <Text fontWeight="semibold" minW="120px">
                Created:
              </Text>
              <Text>{new Date(keyData.created_at).toLocaleString()}</Text>
            </HStack>
            <HStack>
              <Text fontWeight="semibold" minW="120px">
                Expires:
              </Text>
              <Text>{new Date(keyData.expires_at).toLocaleString()}</Text>
            </HStack>
            <HStack>
              <Text fontWeight="semibold" minW="120px">
                Status:
              </Text>
              <Badge colorScheme={keyData.status === 'active' ? 'green' : 'red'}>
                {keyData.status}
              </Badge>
            </HStack>
            <HStack>
              <Text fontWeight="semibold" minW="120px">
                Version:
              </Text>
              <Text>{keyData.version}</Text>
            </HStack>
            {(keyData.symmetric_key || keyData.private_key) && (
              <Alert status="info" mt={2}>
                <AlertIcon />
                Key material found. Ready to register.
              </Alert>
            )}
          </VStack>
        </Box>
      )}

      {/* Register Button */}
      {keyData && !error && (
        <Button
          onClick={handleRegister}
          colorScheme="green"
          isLoading={isLoading}
          loadingText="Registering..."
        >
          Register Key in System
        </Button>
      )}
    </VStack>
  );
}

