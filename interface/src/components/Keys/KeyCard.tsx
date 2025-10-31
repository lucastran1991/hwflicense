'use client';

import {
  Box,
  Heading,
  Text,
  HStack,
  VStack,
  Badge,
  IconButton,
  useColorModeValue,
  Tooltip,
} from '@chakra-ui/react';
import { FaDownload, FaSync, FaTrash } from 'react-icons/fa';
import type { KeyInfo } from '@/lib/types/keys';

interface KeyCardProps {
  keyData: KeyInfo;
  onRefresh?: (keyId: string) => void;
  onRevoke?: (keyId: string) => void;
  onDownload?: (keyId: string) => void;
}

export function KeyCard({ keyData, onRefresh, onRevoke, onDownload }: KeyCardProps) {
  const bgColor = useColorModeValue('white', 'gray.800');
  const borderColor = useColorModeValue('gray.200', 'gray.700');

  const isExpired = keyData.expired || new Date(keyData.expires_at) < new Date();
  const statusColor = isExpired
    ? 'red'
    : keyData.revoked || keyData.status === 'revoked'
    ? 'orange'
    : 'green';

  return (
    <Box
      bg={bgColor}
      borderColor={borderColor}
      borderWidth="1px"
      borderRadius="md"
      boxShadow="sm"
      overflow="hidden"
    >
      <Box p={4} borderBottomWidth="1px" borderBottomColor={borderColor}>
        <HStack justify="space-between" align="start">
          <VStack align="start" spacing={1}>
            <Heading size="sm">{keyData.key_id.substring(0, 8)}...</Heading>
            <HStack spacing={2}>
              <Badge colorScheme={keyData.key_type === 'symmetric' ? 'blue' : 'purple'}>
                {keyData.key_type}
              </Badge>
              <Badge colorScheme={statusColor}>
                {keyData.revoked ? 'revoked' : isExpired ? 'expired' : keyData.status || 'active'}
              </Badge>
            </HStack>
          </VStack>
        </HStack>
      </Box>
      <Box p={4}>
        <VStack align="stretch" spacing={2}>
          <Box>
            <Text fontSize="xs" color="gray.500">
              Created
            </Text>
            <Text fontSize="sm">
              {new Date(keyData.created_at).toLocaleString()}
            </Text>
          </Box>
          <Box>
            <Text fontSize="xs" color="gray.500">
              Expires
            </Text>
            <Text fontSize="sm">
              {new Date(keyData.expires_at).toLocaleString()}
            </Text>
          </Box>
          {keyData.public_key && (
            <Box>
              <Text fontSize="xs" color="gray.500">
                Public Key
              </Text>
              <Text fontSize="xs" fontFamily="mono" noOfLines={1}>
                {keyData.public_key.substring(0, 32)}...
              </Text>
            </Box>
          )}
          <HStack spacing={2} mt={2}>
            {onDownload && (
              <Tooltip label="Download key">
                <IconButton
                  aria-label="Download key"
                  size="sm"
                  variant="outline"
                  colorScheme="blue"
                  icon={<FaDownload />}
                  onClick={() => onDownload(keyData.key_id)}
                />
              </Tooltip>
            )}
            {onRefresh && (
              <Tooltip label="Refresh expiry">
                <IconButton
                  aria-label="Refresh key"
                  size="sm"
                  variant="outline"
                  icon={<FaSync />}
                  onClick={() => onRefresh(keyData.key_id)}
                />
              </Tooltip>
            )}
            {onRevoke && (
              <Tooltip label="Revoke key">
                <IconButton
                  aria-label="Revoke key"
                  size="sm"
                  variant="outline"
                  colorScheme="red"
                  icon={<FaTrash />}
                  onClick={() => onRevoke(keyData.key_id)}
                />
              </Tooltip>
            )}
          </HStack>
        </VStack>
      </Box>
    </Box>
  );
}

