'use client';

import { useState, useEffect } from 'react';
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
  Table,
  Thead,
  Tbody,
  Tr,
  Th,
  Td,
  Link,
} from '@chakra-ui/react';
import { CheckIcon } from '@chakra-ui/icons';

export default function DashboardPage() {
  const [cml, setCml] = useState<any>(null);
  const [sites, setSites] = useState<any[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    loadData();
  }, []);

  const loadData = async () => {
    try {
      const [cmlData, sitesData] = await Promise.all([
        apiClient.getCML().catch(() => null),
        apiClient.listSites().catch(() => null),
      ]);
      
      if (cmlData?.data) {
        setCml(cmlData.data);
      }
      
      if (sitesData) {
        const sitesList = (sitesData as any).data?.sites || (sitesData as any).sites || [];
        setSites(sitesList);
      }
    } catch (error) {
      console.error('Failed to load data:', error);
    } finally {
      setLoading(false);
    }
  };

  if (loading) {
    return (
      <Box display="flex" justifyContent="center" alignItems="center" minH="400px">
        <Spinner size="xl" />
      </Box>
    );
  }

  return (
    <Box maxW="7xl" mx="auto" py={6} px={4}>
      <Heading size="xl" mb={6}>Dashboard</Heading>

      {cml ? (
        <Card mb={6}>
          <CardHeader>
            <Heading size="md">CML Status</Heading>
          </CardHeader>
          <CardBody>
            <Grid templateColumns="repeat(3, 1fr)" gap={4}>
              <Box>
                <Text fontSize="sm" color="gray.500" mb={1}>Organization ID</Text>
                <Text fontSize="lg" fontWeight="medium">{cml.cml?.org_id || 'N/A'}</Text>
              </Box>
              <Box>
                <Text fontSize="sm" color="gray.500" mb={1}>Max Sites</Text>
                <Text fontSize="lg" fontWeight="medium">{cml.cml?.max_sites || 'N/A'}</Text>
              </Box>
              <Box>
                <Text fontSize="sm" color="gray.500" mb={1}>Status</Text>
                <Badge colorScheme="green">Active</Badge>
              </Box>
            </Grid>
          </CardBody>
        </Card>
      ) : (
        <Alert status="warning" mb={6}>
          No CML configured. Please upload a CML to get started.
        </Alert>
      )}

      <Card>
        <CardHeader>
          <Box display="flex" justifyContent="space-between" alignItems="center">
            <Heading size="md">Sites Overview</Heading>
            <Box>
              <Button colorScheme="gray" mr={2} as="a" href="/dashboard/sites">
                Manage Sites
              </Button>
              <Button colorScheme="blue" as="a" href="/dashboard/manifests">
                View Manifests
              </Button>
            </Box>
          </Box>
        </CardHeader>
        <CardBody>
          {sites.length === 0 ? (
            <Text color="gray.500">No sites created yet.</Text>
          ) : (
            <Box overflowX="auto">
              <Table variant="simple">
                <Thead>
                  <Tr>
                    <Th>Site ID</Th>
                    <Th>Status</Th>
                    <Th>Issued At</Th>
                    <Th>Actions</Th>
                  </Tr>
                </Thead>
                <Tbody>
                  {sites.map((site) => (
                    <Tr key={site.site_id}>
                      <Td fontWeight="medium">{site.site_id}</Td>
                      <Td>
                        <Badge colorScheme={site.status === 'active' ? 'green' : 'red'}>
                          {site.status}
                        </Badge>
                      </Td>
                      <Td>{new Date(site.issued_at).toLocaleDateString()}</Td>
                      <Td>
                        <Link color="blue.600" href={`/dashboard/sites/${site.site_id}`}>
                          View
                        </Link>
                      </Td>
                    </Tr>
                  ))}
                </Tbody>
              </Table>
            </Box>
          )}
        </CardBody>
      </Card>
    </Box>
  );
}
