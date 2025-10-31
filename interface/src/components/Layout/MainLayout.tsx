'use client';

import { Box, useColorModeValue } from '@chakra-ui/react';
import { Header } from './Header';
import { Sidebar } from './Sidebar';

interface MainLayoutProps {
  children: React.ReactNode;
}

export function MainLayout({ children }: MainLayoutProps) {
  const bgColor = useColorModeValue('white', 'gray.800');

  return (
    <>
      <Header />
      <Sidebar />
      <Box ml="240px" minH="calc(100vh - 73px)" bg={bgColor} p={6}>
        {children}
      </Box>
    </>
  );
}

