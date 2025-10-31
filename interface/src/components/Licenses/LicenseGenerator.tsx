'use client';

import { useState, useRef } from 'react';
import {
  Box,
  Button,
  FormControl,
  FormLabel,
  Input,
  VStack,
  Textarea,
  useToast,
  Code,
  Text,
  HStack,
  Badge,
  Alert,
  AlertIcon,
  Divider,
  Heading,
  useColorModeValue,
} from '@chakra-ui/react';
import { generateLicense } from '@/lib/api/licenses';
import type { GenerateLicenseRequest, LicenseFile } from '@/lib/types/licenses';
import type { DownloadKeyResponse } from '@/lib/types/keys';

export function LicenseGenerator() {
  const [keyId, setKeyId] = useState<string>('');
  const [licenseType, setLicenseType] = useState<string>('');
  const [metadata, setMetadata] = useState<string>(''); // JSON string
  const [isLoading, setIsLoading] = useState(false);
  const [licenseFile, setLicenseFile] = useState<string | null>(null);
  const [filename, setFilename] = useState<string | null>(null);
  const [uploadedKeyData, setUploadedKeyData] = useState<DownloadKeyResponse | null>(null);
  const [uploadedLicenseData, setUploadedLicenseData] = useState<LicenseFile | null>(null);
  const [uploadError, setUploadError] = useState<string | null>(null);
  const fileInputRef = useRef<HTMLInputElement>(null);
  const toast = useToast();
  
  const bgColor = useColorModeValue('white', 'gray.800');
  const borderColor = useColorModeValue('gray.200', 'gray.700');
  const sectionBgColor = useColorModeValue('gray.50', 'gray.700');

  const handleFileSelect = (event: React.ChangeEvent<HTMLInputElement>) => {
    const file = event.target.files?.[0];
    if (!file) return;

    if (!file.name.endsWith('.json')) {
      setUploadError('Please select a JSON file');
      return;
    }

    setUploadError(null);

    const reader = new FileReader();
    reader.onload = (e) => {
      try {
        const text = e.target?.result as string;
        const parsed = JSON.parse(text);
        
        // Check if it's a license file (has license_id and license_type)
        if (parsed.license_id && parsed.license_type) {
          const licenseData = parsed as LicenseFile;
          
          // Validate license structure
          if (!licenseData.key_id) {
            setUploadError('Invalid license JSON file: missing key_id');
            return;
          }

          setUploadedLicenseData(licenseData);
          setUploadedKeyData(null);
          setKeyId(licenseData.key_id); // Auto-fill key_id
          setLicenseType(licenseData.license_type); // Auto-fill license_type
          
          toast({
            title: 'License loaded successfully',
            description: `License Type: ${licenseData.license_type}, Key ID: ${licenseData.key_id.substring(0, 8)}...`,
            status: 'success',
            duration: 3000,
            isClosable: true,
          });
        } 
        // Check if it's a key file (has key_id and key_type)
        else if (parsed.key_id && parsed.key_type) {
          const keyData = parsed as DownloadKeyResponse;
          
          setUploadedKeyData(keyData);
          setUploadedLicenseData(null);
          setKeyId(keyData.key_id); // Auto-fill key_id
          // License type is not available in key JSON, so don't set it
          
          toast({
            title: 'Key loaded successfully',
            description: `Key ID: ${keyData.key_id.substring(0, 8)}...`,
            status: 'success',
            duration: 3000,
            isClosable: true,
          });
        } else {
          setUploadError('Invalid JSON file: must be either a key JSON (with key_id, key_type) or license JSON (with license_id, license_type)');
        }
      } catch (err: any) {
        setUploadError(`Failed to parse JSON: ${err.message}`);
      }
    };
    reader.onerror = () => {
      setUploadError('Failed to read file');
    };
    reader.readAsText(file);
  };

  const handleJsonTextChange = (event: React.ChangeEvent<HTMLTextAreaElement>) => {
    const text = event.target.value;
    setUploadError(null);

    if (!text.trim()) {
      setUploadedKeyData(null);
      setUploadedLicenseData(null);
      return;
    }

    try {
      const parsed = JSON.parse(text);
      
      // Check if it's a license file (has license_id and license_type)
      if (parsed.license_id && parsed.license_type) {
        const licenseData = parsed as LicenseFile;
        
        // Validate license structure
        if (!licenseData.key_id) {
          setUploadError('Invalid license JSON: missing key_id');
          setUploadedKeyData(null);
          setUploadedLicenseData(null);
          return;
        }

        setUploadedLicenseData(licenseData);
        setUploadedKeyData(null);
        setKeyId(licenseData.key_id); // Auto-fill key_id
        setLicenseType(licenseData.license_type); // Auto-fill license_type
        
        toast({
          title: 'License loaded successfully',
          description: `License Type: ${licenseData.license_type}, Key ID: ${licenseData.key_id.substring(0, 8)}...`,
          status: 'success',
          duration: 3000,
          isClosable: true,
        });
      } 
      // Check if it's a key file (has key_id and key_type)
      else if (parsed.key_id && parsed.key_type) {
        const keyData = parsed as DownloadKeyResponse;
        
        setUploadedKeyData(keyData);
        setUploadedLicenseData(null);
        setKeyId(keyData.key_id); // Auto-fill key_id
        // License type is not available in key JSON, so don't set it
        
        toast({
          title: 'Key loaded successfully',
          description: `Key ID: ${keyData.key_id.substring(0, 8)}...`,
          status: 'success',
          duration: 3000,
          isClosable: true,
        });
      } else {
        setUploadError('Invalid JSON: must be either a key JSON (with key_id, key_type) or license JSON (with license_id, license_type)');
        setUploadedKeyData(null);
        setUploadedLicenseData(null);
      }
    } catch (err: any) {
      setUploadError(`Invalid JSON: ${err.message}`);
      setUploadedKeyData(null);
      setUploadedLicenseData(null);
    }
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setIsLoading(true);

    try {
      let parsedMetadata: Record<string, string> | undefined;
      if (metadata.trim()) {
        try {
          parsedMetadata = JSON.parse(metadata);
        } catch (error) {
          toast({
            title: 'Invalid metadata JSON',
            status: 'error',
            duration: 3000,
            isClosable: true,
          });
          setIsLoading(false);
          return;
        }
      }

      const request: GenerateLicenseRequest = {
        key_id: keyId,
        license_type: licenseType,
        metadata: parsedMetadata,
      };

      const response = await generateLicense(request);
      setLicenseFile(response.license_file);
      setFilename(response.filename);

      toast({
        title: 'License generated successfully',
        status: 'success',
        duration: 3000,
        isClosable: true,
      });
    } catch (error: any) {
      toast({
        title: 'Error generating license',
        description: error.message || 'An error occurred',
        status: 'error',
        duration: 5000,
        isClosable: true,
      });
    } finally {
      setIsLoading(false);
    }
  };

  const handleDownload = () => {
    if (!licenseFile || !filename) return;

    const decoded = atob(licenseFile);
    const blob = new Blob([decoded], { type: 'application/json' });
    const url = URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = filename;
    document.body.appendChild(a);
    a.click();
    document.body.removeChild(a);
    URL.revokeObjectURL(url);
  };

  return (
    <Box>
      {/* Upload Key JSON Section */}
      <Box
        mb={6}
        p={4}
        bg={sectionBgColor}
        borderWidth="1px"
        borderRadius="md"
        borderColor={borderColor}
        boxShadow="sm"
      >
        <Heading size="sm" mb={3}>
          Upload Key/License JSON (Optional)
        </Heading>
        <Text fontSize="sm" color="gray.500" mb={4}>
          Upload a key JSON file to auto-fill Key ID, or a license JSON file to auto-fill both Key ID and License Type.
        </Text>

        <VStack spacing={3} align="stretch">
          {/* File Upload */}
          <Box>
            <Text fontWeight="semibold" fontSize="sm" mb={2}>
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
          </Box>

          {/* Or Paste JSON */}
          <Box>
            <Text fontWeight="semibold" fontSize="sm" mb={2}>
              Option 2: Paste JSON Text
            </Text>
            <Textarea
              onChange={handleJsonTextChange}
              placeholder="Paste key JSON here to auto-fill Key ID..."
              rows={4}
              fontFamily="mono"
              fontSize="xs"
            />
          </Box>

          {/* Error Display */}
          {uploadError && (
            <Alert status="error" size="sm">
              <AlertIcon />
              {uploadError}
            </Alert>
          )}

          {/* Uploaded License Info */}
          {uploadedLicenseData && !uploadError && (
            <Box p={3} bg="green.50" _dark={{ bg: 'green.900' }} borderRadius="md" borderWidth="1px" borderColor="green.200" _dark={{ borderColor: 'green.700' }}>
              <HStack spacing={2} mb={2}>
                <Text fontWeight="semibold" fontSize="sm">
                  License Loaded:
                </Text>
                <Badge colorScheme="green" fontSize="xs">
                  {uploadedLicenseData.license_type}
                </Badge>
              </HStack>
              <VStack align="start" spacing={1} fontSize="xs">
                <HStack>
                  <Text color="gray.600" _dark={{ color: 'gray.400' }}>
                    License ID:
                  </Text>
                  <Code fontSize="xs">{uploadedLicenseData.license_id}</Code>
                </HStack>
                <HStack>
                  <Text color="gray.600" _dark={{ color: 'gray.400' }}>
                    Key ID:
                  </Text>
                  <Code fontSize="xs">{uploadedLicenseData.key_id}</Code>
                </HStack>
              </VStack>
            </Box>
          )}

          {/* Uploaded Key Info */}
          {uploadedKeyData && !uploadedLicenseData && !uploadError && (
            <Box p={3} bg="blue.50" _dark={{ bg: 'blue.900' }} borderRadius="md" borderWidth="1px" borderColor="blue.200" _dark={{ borderColor: 'blue.700' }}>
              <HStack spacing={2} mb={2}>
                <Text fontWeight="semibold" fontSize="sm">
                  Key Loaded:
                </Text>
                <Code fontSize="xs">{uploadedKeyData.key_id}</Code>
                <Badge colorScheme={uploadedKeyData.key_type === 'symmetric' ? 'blue' : 'purple'} fontSize="xs">
                  {uploadedKeyData.key_type}
                </Badge>
              </HStack>
              {uploadedKeyData.public_key && (
                <Text fontSize="xs" color="gray.600" _dark={{ color: 'gray.400' }}>
                  Public Key: {uploadedKeyData.public_key.substring(0, 32)}...
                </Text>
              )}
            </Box>
          )}
        </VStack>
      </Box>

      <Divider mb={6} />

      <Box as="form" onSubmit={handleSubmit} mb={6}>
        <VStack spacing={4} align="stretch">
          <FormControl isRequired>
            <FormLabel>Key ID</FormLabel>
            <Input
              value={keyId}
              onChange={(e) => setKeyId(e.target.value)}
              placeholder="Enter key ID or upload JSON above"
            />
          </FormControl>

          <FormControl isRequired>
            <FormLabel>License Type</FormLabel>
            <Input
              value={licenseType}
              onChange={(e) => setLicenseType(e.target.value)}
              placeholder="e.g., enterprise, site, trial"
            />
          </FormControl>

          <FormControl>
            <FormLabel>Metadata (JSON, optional)</FormLabel>
            <Textarea
              value={metadata}
              onChange={(e) => setMetadata(e.target.value)}
              placeholder='{"customer_id": "CUST001", "site_name": "Main Office", "max_users": "100"}'
              fontFamily="mono"
              fontSize="sm"
              rows={4}
            />
          </FormControl>

          <Button
            type="submit"
            colorScheme="blue"
            isLoading={isLoading}
            loadingText="Generating..."
          >
            Generate License
          </Button>
        </VStack>
      </Box>

      {licenseFile && (
        <Box mt={6}>
          <Text mb={2} fontWeight="semibold">
            Generated License File:
          </Text>
          <Code
            p={4}
            display="block"
            whiteSpace="pre-wrap"
            maxH="300px"
            overflowY="auto"
            fontSize="xs"
          >
            {atob(licenseFile)}
          </Code>
          {filename && (
            <Button
              mt={4}
              colorScheme="green"
              onClick={handleDownload}
            >
              Download {filename}
            </Button>
          )}
        </Box>
      )}
    </Box>
  );
}

