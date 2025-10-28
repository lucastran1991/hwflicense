module.exports = {
  apps: [
    {
      name: 'license-server',
      script: './wrapper-license-server.sh',
      instances: 1,
      exec_mode: 'fork',
      autorestart: true,
      watch: false,
      max_memory_restart: '200M',
      error_file: './logs/system.log',
      out_file: './logs/system.log',
      merge_logs: true,
      log_date_format: 'YYYY-MM-DD HH:mm:ss',
      env: {
        LOG_PREFIX: '[LS]'
      }
    },
    {
      name: 'hub-backend',
      script: './wrapper-backend.sh',
      instances: 1,
      exec_mode: 'fork',
      autorestart: true,
      watch: false,
      max_memory_restart: '500M',
      error_file: './logs/system.log',
      out_file: './logs/system.log',
      merge_logs: true,
      log_date_format: 'YYYY-MM-DD HH:mm:ss',
      env: {
        LOG_PREFIX: '[BE]'
      }
    },
    {
      name: 'frontend',
      script: './wrapper-frontend.sh',
      instances: 1,
      exec_mode: 'fork',
      autorestart: true,
      watch: false,
      max_memory_restart: '300M',
      error_file: './logs/system.log',
      out_file: './logs/system.log',
      merge_logs: true,
      log_date_format: 'YYYY-MM-DD HH:mm:ss',
      env: {
        LOG_PREFIX: '[FE]'
      }
    }
  ]
};

