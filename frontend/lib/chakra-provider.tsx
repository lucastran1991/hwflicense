'use client';

import { ChakraProvider, extendTheme } from '@chakra-ui/react';

const theme = extendTheme({
  config: {
    initialColorMode: 'dark',
    useSystemColorMode: false,
  },
  styles: {
    global: {
      body: {
        bg: 'linear-gradient(135deg, #1a0b2e 0%, #2d1b4e 50%, #4a2c6d 100%)',
        color: 'white',
      },
    },
  },
  colors: {
    brand: {
      50: '#fce7f3',
      100: '#fbcfe8',
      200: '#f9a8d4',
      300: '#f472b6',
      400: '#ec4899',
      500: '#d946ef',
      600: '#c026d3',
      700: '#a21caf',
      800: '#86198f',
      900: '#701a75',
    },
    purple: {
      50: '#faf5ff',
      100: '#f3e8ff',
      200: '#e9d5ff',
      300: '#d8b4fe',
      400: '#c084fc',
      500: '#a855f7',
      600: '#9333ea',
      700: '#7e22ce',
      800: '#6b21a8',
      900: '#581c87',
    },
  },
  components: {
    Card: {
      baseStyle: {
        container: {
          bg: 'rgba(74, 44, 109, 0.6)',
          backdropFilter: 'blur(10px)',
          borderColor: 'rgba(217, 70, 239, 0.3)',
        },
      },
    },
    Button: {
      variants: {
        solid: {
          bg: 'linear-gradient(135deg, #d946ef 0%, #c026d3 100%)',
          _hover: {
            bg: 'linear-gradient(135deg, #e91e63 0%, #d946ef 100%)',
          },
        },
      },
    },
    Badge: {
      baseStyle: {
        bg: 'rgba(217, 70, 239, 0.2)',
        color: '#f9a8d4',
        borderColor: '#d946ef',
      },
    },
  },
});

export function ChakraUIProvider({ children }: { children: React.ReactNode }) {
  return <ChakraProvider theme={theme}>{children}</ChakraProvider>;
}

