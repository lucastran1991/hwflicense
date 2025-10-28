'use client';

import { useState } from 'react';
import { useAuth } from '@/lib/auth-context';
import {
  Box,
  VStack,
  Heading,
  Text,
  FormControl,
  FormLabel,
  Input,
  Button,
  Alert,
  AlertIcon,
  Container,
} from '@chakra-ui/react';

export default function LoginPage() {
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);
  const { login } = useAuth();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    setLoading(true);

    try {
      await login(username, password);
    } catch (err) {
      setError('Invalid credentials');
    } finally {
      setLoading(false);
    }
  };

  return (
    <Box minH="100vh" display="flex" alignItems="center" justifyContent="center">
      <Container maxW="md">
        <Box bg="rgba(74, 44, 109, 0.8)" backdropFilter="blur(20px)" p={8} borderRadius="lg" borderColor="rgba(217, 70, 239, 0.3)" borderWidth="1px">
          <VStack spacing={6} align="stretch">
            <Box>
              <Heading textAlign="center" size="xl" color="white">
                TaskMaster License Hub
              </Heading>
              <Text textAlign="center" color="whiteAlpha.800" mt={2}>
                Sign in to manage your licenses
              </Text>
            </Box>

            <form onSubmit={handleSubmit}>
              <VStack spacing={4}>
                {error && (
                  <Alert status="error" rounded="md">
                    <AlertIcon />
                    {error}
                  </Alert>
                )}

                <FormControl isRequired>
                  <FormLabel>Username</FormLabel>
                  <Input
                    id="username"
                    type="text"
                    value={username}
                    onChange={(e) => setUsername(e.target.value)}
                    placeholder="admin"
                  />
                </FormControl>

                <FormControl isRequired>
                  <FormLabel>Password</FormLabel>
                  <Input
                    id="password"
                    type="password"
                    value={password}
                    onChange={(e) => setPassword(e.target.value)}
                    placeholder="admin123"
                  />
                </FormControl>

                <Button
                  type="submit"
                  colorScheme="blue"
                  width="full"
                  isLoading={loading}
                  loadingText="Signing in..."
                >
                  Sign in
                </Button>

                <Text fontSize="sm" color="whiteAlpha.600" textAlign="center">
                  Default credentials: admin / admin123
                </Text>
              </VStack>
            </form>
          </VStack>
        </Box>
      </Container>
    </Box>
  );
}

