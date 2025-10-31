'use client';

import { useState, useEffect } from 'react';
import {
  Box,
  Container,
  Heading,
  Text,
  SimpleGrid,
  Stat,
  StatLabel,
  StatNumber,
  StatHelpText,
  useColorModeValue,
  Alert,
  AlertIcon,
  Spinner,
  Center,
} from '@chakra-ui/react';
import { MainLayout } from '@/components/Layout/MainLayout';
import { checkHealth } from '@/lib/api/client';

export default function Dashboard() {
  const [healthStatus, setHealthStatus] = useState<'ok' | 'error' | 'loading'>('loading');
  const bgColor = useColorModeValue('white', 'gray.800');

  useEffect(() => {
    async function checkBackendHealth() {
      try {
        await checkHealth();
        setHealthStatus('ok');
      } catch (error) {
        setHealthStatus('error');
      }
    }
    checkBackendHealth();
  }, []);

  return (
    <MainLayout>
      <Container maxW="container.xl">
          <Heading mb={6}>Dashboard</Heading>

          {healthStatus === 'loading' && (
            <Center py={8}>
              <Spinner size="xl" />
            </Center>
          )}

          {healthStatus === 'error' && (
            <Alert status="error" mb={6}>
              <AlertIcon />
              Cannot connect to backend. Please check if the KMS service is running.
            </Alert>
          )}

          {healthStatus === 'ok' && (
            <>
              <SimpleGrid columns={{ base: 1, md: 3 }} spacing={6} mb={6}>
                <Stat bg={bgColor} p={4} borderRadius="md" boxShadow="sm">
                  <StatLabel>Backend Status</StatLabel>
                  <StatNumber color="green.500">Online</StatNumber>
                  <StatHelpText>KMS Service is running</StatHelpText>
                </Stat>

                <Stat bg={bgColor} p={4} borderRadius="md" boxShadow="sm">
                  <StatLabel>API Endpoint</StatLabel>
                  <StatNumber fontSize="lg">
                    {process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080'}
                  </StatNumber>
                  <StatHelpText>Backend URL</StatHelpText>
                </Stat>

                <Stat bg={bgColor} p={4} borderRadius="md" boxShadow="sm">
                  <StatLabel>Version</StatLabel>
                  <StatNumber fontSize="lg">v0.1.0</StatNumber>
                  <StatHelpText>Interface version</StatHelpText>
                </Stat>
              </SimpleGrid>

              <Box bg={bgColor} p={6} borderRadius="md" boxShadow="sm">
                <Heading size="md" mb={4}>
                  Quick Start
                </Heading>
                <Box>
                  <Text mb={2}>
                    <strong>1. Create a Key:</strong> Go to the Keys page to generate symmetric or asymmetric keys.
                  </Text>
                  <Text mb={2}>
                    <strong>2. Generate License:</strong> Use the Licenses page to create license files from your keys.
                  </Text>
                  <Text mb={2}>
                    <strong>3. Validate License:</strong> Upload or paste a license file to validate its signature and status.
                  </Text>
                </Box>
              </Box>
            </>
          )}
        </Container>
    </MainLayout>
  );
}

