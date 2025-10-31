import { Providers } from './providers';
import './globals.css';

export const metadata = {
  title: 'KMS - Key Management Service',
  description: 'Key Management Service Interface',
};

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="en">
      <body>
        <Providers>{children}</Providers>
      </body>
    </html>
  );
}
