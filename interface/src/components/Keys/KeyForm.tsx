'use client';

import { useState } from 'react';
import {
  Box,
  Button,
  FormControl,
  FormLabel,
  Input,
  Select,
  VStack,
  NumberInput,
  NumberInputField,
  useToast,
} from '@chakra-ui/react';
import type { RegisterKeyRequest } from '@/lib/types/keys';
import { registerKey } from '@/lib/api/keys';

interface KeyFormProps {
  onSuccess?: () => void;
}

export function KeyForm({ onSuccess }: KeyFormProps) {
  const [keyType, setKeyType] = useState<'symmetric' | 'asymmetric'>('symmetric');
  const [expiresInSeconds, setExpiresInSeconds] = useState<string>('31536000'); // 1 year
  const [keyMaterial, setKeyMaterial] = useState<string>('');
  const [isLoading, setIsLoading] = useState(false);
  const toast = useToast();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setIsLoading(true);

    try {
      const request: RegisterKeyRequest = {
        key_type: keyType,
        expires_in_seconds: parseInt(expiresInSeconds, 10),
      };

      if (keyMaterial.trim()) {
        request.key_material = keyMaterial.trim();
      }

      await registerKey(request);
      
      toast({
        title: 'Key created successfully',
        status: 'success',
        duration: 3000,
        isClosable: true,
      });

      // Reset form
      setKeyType('symmetric');
      setExpiresInSeconds('31536000');
      setKeyMaterial('');
      
      if (onSuccess) {
        onSuccess();
      }
    } catch (error: any) {
      toast({
        title: 'Error creating key',
        description: error.message || 'An error occurred',
        status: 'error',
        duration: 5000,
        isClosable: true,
      });
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <Box as="form" onSubmit={handleSubmit}>
      <VStack spacing={4} align="stretch">
        <FormControl isRequired>
          <FormLabel>Key Type</FormLabel>
          <Select
            value={keyType}
            onChange={(e) => setKeyType(e.target.value as 'symmetric' | 'asymmetric')}
          >
            <option value="symmetric">Symmetric (AES-256)</option>
            <option value="asymmetric">Asymmetric (Ed25519)</option>
          </Select>
        </FormControl>

        <FormControl>
          <FormLabel>Expires In (seconds)</FormLabel>
          <NumberInput
            value={expiresInSeconds}
            onChange={(valueString) => setExpiresInSeconds(valueString)}
            min={1}
          >
            <NumberInputField />
          </NumberInput>
          <Box fontSize="xs" color="gray.500" mt={1}>
            Default: 31536000 (1 year)
          </Box>
        </FormControl>

        <FormControl>
          <FormLabel>
            Key Material (Base64, optional - will be generated if empty)
          </FormLabel>
          <Input
            value={keyMaterial}
            onChange={(e) => setKeyMaterial(e.target.value)}
            placeholder="Leave empty to auto-generate"
            fontFamily="mono"
            fontSize="sm"
          />
        </FormControl>

        <Button
          type="submit"
          colorScheme="blue"
          isLoading={isLoading}
          loadingText="Creating..."
        >
          Create Key
        </Button>
      </VStack>
    </Box>
  );
}

