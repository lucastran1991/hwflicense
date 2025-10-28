import "./globals.css";
import { AuthProvider } from "@/lib/auth-context";
import { ChakraUIProvider } from "@/lib/chakra-provider";

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="en">
      <body>
        <ChakraUIProvider>
          <AuthProvider>{children}</AuthProvider>
        </ChakraUIProvider>
      </body>
    </html>
  );
}
