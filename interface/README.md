# KMS Interface - Next.js 13 Frontend

A Next.js 13 frontend application for the Key Management Service (KMS) backend, built with ChakraUI and TypeScript. Designed for deployment on AWS EC2 with Node.js 16 compatibility.

## Tech Stack

- **Next.js**: 13.5.6 (compatible with Node.js 16)
- **React**: 18.2.0
- **ChakraUI**: 1.8.8 (v1.x compatible with Node.js 16)
- **TypeScript**: 5.x
- **Axios**: For API communication

## Prerequisites

- Node.js 16.x or higher
- npm 8.x or higher
- Access to the KMS backend service (default: http://localhost:8080)

## Installation

1. Install dependencies:

```bash
npm install
```

2. Create environment file:

```bash
cp .env.local.example .env.local
```

3. Edit `.env.local` and set your backend API URL:

```env
NEXT_PUBLIC_API_URL=http://localhost:8080
```

## Development

Start the development server:

```bash
npm run dev
```

The application will be available at [http://localhost:3000](http://localhost:3000).

## Building for Production

Build the application:

```bash
npm run build
```

Start the production server:

```bash
npm start
```

## Project Structure

```
core/interface/
├── src/
│   ├── app/              # Next.js App Router pages
│   │   ├── layout.tsx    # Root layout
│   │   ├── page.tsx      # Dashboard page
│   │   ├── keys/         # Keys management page
│   │   └── licenses/     # Licenses management page
│   ├── components/        # React components
│   │   ├── Keys/         # Key management components
│   │   ├── Licenses/     # License management components
│   │   └── Layout/       # Layout components (Header, Sidebar)
│   ├── lib/
│   │   ├── api/          # API client functions
│   │   ├── types/        # TypeScript type definitions
│   │   └── theme.ts      # ChakraUI theme configuration
│   └── hooks/            # Custom React hooks
├── public/               # Static assets
├── package.json
├── tsconfig.json
├── next.config.js
└── README.md
```

## Features

### Key Management
- Create symmetric (AES-256) or asymmetric (Ed25519) keys
- Set custom expiration times
- Provide custom key material or auto-generate

### License Management
- Generate license files from keys
- Add custom metadata to licenses
- Validate license files (file upload or base64 content)
- Download generated license files

### Dashboard
- Backend health check status
- API endpoint information
- Quick start guide

## API Endpoints

The frontend communicates with the KMS backend API:

- `GET /health` - Health check
- `POST /keys` - Register/create a key
- `POST /keys/validate` - Validate a key
- `POST /keys/:id/refresh` - Refresh key expiry
- `DELETE /keys/:id` - Revoke a key
- `POST /licenses/generate` - Generate a license file
- `POST /licenses/validate` - Validate a license file

## Environment Variables

### `NEXT_PUBLIC_API_URL`

The base URL of the KMS backend API. Default: `http://localhost:8080`

Example:
```env
NEXT_PUBLIC_API_URL=http://localhost:8080
```

For production:
```env
NEXT_PUBLIC_API_URL=https://api.example.com
```

## Deployment to AWS EC2

### Prerequisites

- AWS EC2 instance with Node.js 16 installed
- KMS backend service running (can be on the same instance or separate)
- Network access between frontend and backend

### Steps

1. **Build the application locally:**

```bash
npm run build
```

2. **Transfer files to EC2:**

```bash
# Copy the entire interface directory to EC2
scp -r core/interface/ user@ec2-instance:/path/to/destination/
```

Or use a deployment tool like:
- AWS CodeDeploy
- GitHub Actions
- CI/CD pipeline

3. **On EC2 instance:**

```bash
# Navigate to the interface directory
cd /path/to/interface

# Install dependencies (if not already done)
npm install --production

# Set environment variables
export NEXT_PUBLIC_API_URL=http://localhost:8080

# Build (if not built locally)
npm run build

# Start the production server
npm start
```

4. **Run as a service (optional):**

Create a systemd service file `/etc/systemd/system/kms-interface.service`:

```ini
[Unit]
Description=KMS Interface Next.js Application
After=network.target

[Service]
Type=simple
User=your-user
WorkingDirectory=/path/to/interface
Environment=NEXT_PUBLIC_API_URL=http://localhost:8080
ExecStart=/usr/bin/npm start
Restart=always

[Install]
WantedBy=multi-user.target
```

Enable and start the service:

```bash
sudo systemctl enable kms-interface
sudo systemctl start kms-interface
sudo systemctl status kms-interface
```

### Using PM2 (Alternative)

```bash
# Install PM2
npm install -g pm2

# Start the application
pm2 start npm --name "kms-interface" -- start

# Save PM2 configuration
pm2 save

# Setup PM2 to start on boot
pm2 startup
```

### Nginx Configuration (Optional)

If you want to use Nginx as a reverse proxy:

```nginx
server {
    listen 80;
    server_name your-domain.com;

    location / {
        proxy_pass http://localhost:3000;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_cache_bypass $http_upgrade;
    }
}
```

## CORS Configuration

If the frontend and backend are on different origins, you may need to configure CORS in the Go backend. The backend should allow requests from the frontend's origin.

## Troubleshooting

### Backend Connection Issues

- Check if the KMS backend is running
- Verify `NEXT_PUBLIC_API_URL` is correct
- Check network connectivity between frontend and backend
- Review browser console for CORS errors

### Build Issues

- Ensure Node.js version is 16.x or higher
- Clear `.next` directory and rebuild: `rm -rf .next && npm run build`
- Check for TypeScript errors: `npm run lint`

### Runtime Errors

- Check browser console for errors
- Verify environment variables are set correctly
- Ensure all dependencies are installed: `npm install`

## Development Notes

- The application uses Next.js 13 App Router (not Pages Router)
- All components that use client-side features must have `'use client'` directive
- ChakraUI v1.8.8 is used for Node.js 16 compatibility (v2.x requires Node.js 20+)
- The backend doesn't have a key list endpoint, so keys need to be tracked separately or added to the backend

## License

[Add your license here]

