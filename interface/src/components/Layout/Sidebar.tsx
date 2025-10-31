'use client';

import {
  Box,
  VStack,
  Link as ChakraLink,
  useColorModeValue,
} from '@chakra-ui/react';
import Link from 'next/link';
import { usePathname } from 'next/navigation';

interface NavItem {
  label: string;
  href: string;
  icon?: string;
}

const navItems: NavItem[] = [
  { label: 'Dashboard', href: '/' },
  { label: 'Keys', href: '/keys' },
  { label: 'Licenses', href: '/licenses' },
];

export function Sidebar() {
  const pathname = usePathname();
  const bgColor = useColorModeValue('white', 'gray.800');
  const borderColor = useColorModeValue('gray.200', 'gray.700');
  const activeBg = useColorModeValue('blue.50', 'blue.900');
  const activeColor = useColorModeValue('blue.600', 'blue.300');
  const hoverBg = useColorModeValue('gray.50', 'gray.700');

  return (
    <Box
      as="aside"
      w="240px"
      h="calc(100vh - 73px)"
      bg={bgColor}
      borderRight="1px"
      borderColor={borderColor}
      position="fixed"
      left={0}
      top="73px"
      overflowY="auto"
    >
      <VStack spacing={1} align="stretch" p={4}>
        {navItems.map((item) => {
          const isActive = pathname === item.href;
          return (
            <ChakraLink
              key={item.href}
              as={Link}
              href={item.href}
              px={4}
              py={3}
              borderRadius="md"
              bg={isActive ? activeBg : 'transparent'}
              color={isActive ? activeColor : 'inherit'}
              fontWeight={isActive ? 'semibold' : 'normal'}
              _hover={{
                bg: isActive ? activeBg : hoverBg,
                textDecoration: 'none',
              }}
            >
              {item.label}
            </ChakraLink>
          );
        })}
      </VStack>
    </Box>
  );
}

