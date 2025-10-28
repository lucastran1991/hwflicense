const fs = require('fs');
const path = require('path');

// Load frontend.json config
let config = {};
try {
  const configPath = path.join(__dirname, '../config/frontend.json');
  const configData = fs.readFileSync(configPath, 'utf8');
  config = JSON.parse(configData);
} catch (error) {
  console.log('No config file found, using defaults');
}

module.exports = {
  env: {
    NEXT_PUBLIC_API_URL: config.api_url || process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api',
  },
};
