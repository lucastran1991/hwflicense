const fs = require('fs');
const path = require('path');

// Load environment.json configuration
function loadEnvironmentConfig() {
  const defaultConfig = {
    backend: { port: 8080, db_path: './kms.db', host: 'localhost' },
    frontend: { port: 3000, api_url: 'http://localhost:8080' },
    database: { path: './kms.db', type: 'boltdb' },
    server: { read_timeout: 15, write_timeout: 15, idle_timeout: 60 }
  };

  // Try multiple possible paths for environment.json
  const possiblePaths = [
    path.join(__dirname, '../config/environment.json'),
    path.join(__dirname, '../../config/environment.json'),
    path.join(__dirname, './config/environment.json'),
    path.join(process.cwd(), 'config/environment.json'),
    path.join(process.cwd(), '../config/environment.json'),
  ];

  for (const configPath of possiblePaths) {
    try {
      if (fs.existsSync(configPath)) {
        const data = fs.readFileSync(configPath, 'utf8');
        const config = JSON.parse(data);
        return config;
      }
    } catch (error) {
      // Continue to next path if this one fails
      continue;
    }
  }

  return defaultConfig;
}

// Load config from environment.json
const envConfig = loadEnvironmentConfig();

// Determine backend API URL and frontend port
// Priority: env var > environment.json frontend.api_url > constructed from backend.port
const backendPort = envConfig.backend?.port || 8080;
const backendHost = envConfig.backend?.host || 'localhost';
const frontendPort = envConfig.frontend?.port || 3000;

// Use api_url from environment.json if available, otherwise construct from backend config
let backendApiUrl = process.env.NEXT_PUBLIC_API_URL;
if (!backendApiUrl) {
  if (envConfig.frontend?.api_url) {
    backendApiUrl = envConfig.frontend.api_url;
  } else {
    // Construct from backend host and port
    backendApiUrl = `http://${backendHost}:${backendPort}`;
  }
}

/** @type {import('next').NextConfig} */
const nextConfig = {
  reactStrictMode: true,
  swcMinify: true,
  // output: 'standalone', // Disabled for local deployment, enable for Docker/container deployment
  // Configure for Node.js 16 compatibility
  compiler: {
    emotion: true,
  },
  // Skip type checking during build (type checking done separately)
  typescript: {
    // !! WARN !!
    // Dangerously allow production builds to successfully complete even if
    // your project has type errors.
    // !! WARN !!
    ignoreBuildErrors: true,
  },
  // Environment variables
  env: {
    NEXT_PUBLIC_API_URL: backendApiUrl,
    FRONTEND_PORT: frontendPort.toString(),
    BACKEND_PORT: backendPort.toString(),
  },
}

module.exports = nextConfig

