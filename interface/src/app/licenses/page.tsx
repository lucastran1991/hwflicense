'use client';

import {
  Box,
  Container,
  Heading,
  Tabs,
  TabList,
  TabPanels,
  Tab,
  TabPanel,
  useColorModeValue,
} from '@chakra-ui/react';
import { MainLayout } from '@/components/Layout/MainLayout';
import { LicenseGenerator } from '@/components/Licenses/LicenseGenerator';
import { LicenseValidator } from '@/components/Licenses/LicenseValidator';

export default function LicensesPage() {
  const bgColor = useColorModeValue('white', 'gray.800');

  return (
    <MainLayout>
      <Container maxW="container.xl">
          <Heading mb={6}>License Management</Heading>

          <Tabs>
            <TabList>
              <Tab>Generate License</Tab>
              <Tab>Validate License</Tab>
            </TabList>

            <TabPanels>
              <TabPanel>
                <Box bg={bgColor} p={6} borderRadius="md" boxShadow="sm" maxW="800px">
                  <LicenseGenerator />
                </Box>
              </TabPanel>

              <TabPanel>
                <Box bg={bgColor} p={6} borderRadius="md" boxShadow="sm" maxW="800px">
                  <LicenseValidator />
                </Box>
              </TabPanel>
            </TabPanels>
          </Tabs>
        </Container>
    </MainLayout>
  );
}

