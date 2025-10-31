'use client';

import {
  Box,
  Container,
  Heading,
  VStack,
  Tabs,
  TabList,
  TabPanels,
  Tab,
  TabPanel,
  useColorModeValue,
} from '@chakra-ui/react';
import { MainLayout } from '@/components/Layout/MainLayout';
import { KeyForm } from '@/components/Keys/KeyForm';
import { KeyList } from '@/components/Keys/KeyList';
import { KeyUpload } from '@/components/Keys/KeyUpload';

export default function KeysPage() {
  const bgColor = useColorModeValue('white', 'gray.800');

  return (
    <MainLayout>
      <Container maxW="container.xl">
          <Heading mb={6}>Keys Management</Heading>

          <Tabs>
            <TabList>
              <Tab>Create Key</Tab>
              <Tab>Upload Key</Tab>
              <Tab>Key List</Tab>
            </TabList>

            <TabPanels>
              <TabPanel>
                <Box bg={bgColor} p={6} borderRadius="md" boxShadow="sm" maxW="600px">
                  <KeyForm />
                </Box>
              </TabPanel>

              <TabPanel>
                <Box bg={bgColor} p={6} borderRadius="md" boxShadow="sm" maxW="800px">
                  <KeyUpload />
                </Box>
              </TabPanel>

              <TabPanel>
                <Box bg={bgColor} p={6} borderRadius="md" boxShadow="sm">
                  <KeyList />
                </Box>
              </TabPanel>
            </TabPanels>
          </Tabs>
        </Container>
    </MainLayout>
  );
}

