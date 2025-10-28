import type { NextConfig } from "next";
import * as fs from "fs";
import * as path from "path";

// Load frontend.json config
let config: any = {};
try {
  const configPath = path.join(process.cwd(), "../config/frontend.json");
  const configData = fs.readFileSync(configPath, "utf8");
  config = JSON.parse(configData);
} catch (error) {
  console.log("No config file found, using defaults");
}

const nextConfig: NextConfig = {
  env: {
    NEXT_PUBLIC_API_URL: config.api_url || process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080/api",
  },
};

export default nextConfig;
