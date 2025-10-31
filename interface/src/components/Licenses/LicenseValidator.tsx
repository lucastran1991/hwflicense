'use client';

import { useState, useRef } from 'react';
import {
  Box,
  Button,
  FormControl,
  FormLabel,
  Textarea,
  VStack,
  useToast,
  Alert,
  AlertIcon,
  Code,
  Text,
  HStack,
} from '@chakra-ui/react';
import { validateLicense } from '@/lib/api/licenses';
import type { ValidateLicenseResponse } from '@/lib/types/licenses';

export function LicenseValidator() {
  const [licenseContent, setLicenseContent] = useState<string>('');
  const [file, setFile] = useState<File | null>(null);
  const [isLoading, setIsLoading] = useState(false);
  const [validationResult, setValidationResult] = useState<ValidateLicenseResponse | null>(null);
  const fileInputRef = useRef<HTMLInputElement>(null);
  const toast = useToast();

  const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const selectedFile = e.target.files?.[0];
    if (selectedFile) {
      setFile(selectedFile);
      const reader = new FileReader();
      reader.onload = (event) => {
        const content = event.target?.result as string;
        // Convert to base64
        const base64 = btoa(content);
        setLicenseContent(base64);
      };
      reader.readAsText(selectedFile);
    }
  };

  const handleValidateFile = async () => {
    if (!file) {
      toast({
        title: 'Please select a file',
        status: 'warning',
        duration: 3000,
        isClosable: true,
      });
      return;
    }

    setIsLoading(true);
    try {
      const formData = new FormData();
      formData.append('file', file);

      const result = await validateLicense(formData);
      setValidationResult(result);
    } catch (error: any) {
      toast({
        title: 'Error validating license',
        description: error.message || 'An error occurred',
        status: 'error',
        duration: 5000,
        isClosable: true,
      });
    } finally {
      setIsLoading(false);
    }
  };

  const handleValidateContent = async () => {
    if (!licenseContent.trim()) {
      toast({
        title: 'Please enter license content',
        status: 'warning',
        duration: 3000,
        isClosable: true,
      });
      return;
    }

    setIsLoading(true);
    try {
      // License content should be base64 encoded JSON
      const result = await validateLicense({
        license_content: licenseContent,
      });
      setValidationResult(result);
    } catch (error: any) {
      toast({
        title: 'Error validating license',
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
    <Box>
      <VStack spacing={4} align="stretch" mb={6}>
        <FormControl>
          <FormLabel>Upload License File</FormLabel>
          <HStack>
            <input
              ref={fileInputRef}
              type="file"
              accept=".lic"
              onChange={handleFileChange}
              style={{ display: 'none' }}
            />
            <Button
              onClick={() => fileInputRef.current?.click()}
              colorScheme="blue"
              variant="outline"
              size="sm"
            >
              Choose License File
            </Button>
            {file && (
              <Button
                onClick={handleValidateFile}
                isLoading={isLoading}
                colorScheme="blue"
              >
                Validate File
              </Button>
            )}
            {file && (
              <Text fontSize="sm" color="gray.600">
                {file.name}
              </Text>
            )}
          </HStack>
        </FormControl>

        <Box>
          <Text mb={2} fontSize="sm" fontWeight="semibold">
            OR
          </Text>
        </Box>

        <FormControl>
          <FormLabel>License Content (Base64 encoded JSON)</FormLabel>
          <Textarea
            value={licenseContent}
            onChange={(e) => setLicenseContent(e.target.value)}
            placeholder="Paste base64 encoded license content here"
            fontFamily="mono"
            fontSize="sm"
            rows={6}
          />
          <Button
            mt={2}
            onClick={handleValidateContent}
            isLoading={isLoading}
            colorScheme="blue"
          >
            Validate Content
          </Button>
        </FormControl>
      </VStack>

      {validationResult && (
        <Box mt={6}>
          <Alert
            status={validationResult.valid ? 'success' : 'error'}
            mb={4}
          >
            <AlertIcon />
            {validationResult.valid
              ? 'License is valid'
              : validationResult.error || 'License is invalid'}
          </Alert>

          {validationResult.valid && (
            <Box>
              <Text mb={2} fontWeight="semibold">
                License Details:
              </Text>
              <Code p={4} display="block" whiteSpace="pre-wrap">
                {JSON.stringify(
                  {
                    license_id: validationResult.license_id,
                    license_type: validationResult.license_type,
                    key_id: validationResult.key_id,
                    expires_at: validationResult.expires_at,
                    metadata: validationResult.metadata,
                  },
                  null,
                  2
                )}
              </Code>
            </Box>
          )}

          {(validationResult.expired || validationResult.revoked) && (
            <Alert
              status="warning"
              mt={4}
            >
              <AlertIcon />
              {validationResult.expired && 'License has expired. '}
              {validationResult.revoked && 'License key has been revoked.'}
            </Alert>
          )}
        </Box>
      )}
    </Box>
  );
}

