'use client';

import {
  Box,
  Flex,
  Heading,
  HStack,
  Text,
  useColorModeValue,
} from '@chakra-ui/react';
import Link from 'next/link';

export function Header() {
  const bgColor = useColorModeValue('white', 'gray.800');
  const borderColor = useColorModeValue('gray.200', 'gray.700');

  return (
    <Box
      as="header"
      bg={bgColor}
      borderBottom="1px"
      borderColor={borderColor}
      px={6}
      py={4}
      position="sticky"
      top={0}
      zIndex={100}
      boxShadow="sm"
    >
      <Flex align="center" justify="space-between">
        <HStack spacing={4}>
          <Heading size="lg" as={Link} href="/">
            KMS
          </Heading>
          <Text fontSize="sm" color="gray.500">
            Key Management Service
          </Text>
        </HStack>
        <HStack spacing={4}>
          <Text fontSize="sm" color="gray.600">
            Interface v0.1.0
          </Text>
        </HStack>
      </Flex>
    </Box>
  );
}

