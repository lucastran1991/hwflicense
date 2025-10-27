# TaskMaster License Backend

License management backend for the Hub component of the TaskMaster-AI system.

## Setup

### Prerequisites
- Go 1.21 or higher
- SQLite3

### Installation

1. Install dependencies:
```bash
cd backend
go mod download
```

2. Generate root key pair:
```bash
go run cmd/genkeys/main.go root
```

3. (Optional) Copy environment variables template:
```bash
cp .env.example .env
```

### Running the Server

```bash
go run cmd/server/main.go
```

The server will:
- Create the database at `data/taskmaster_license.db`
- Run migrations automatically
- Start on port 8080 (configurable via `API_PORT` environment variable)

## Project Structure

```
backend/
├── cmd/
│   ├── server/         # Main Hub server
│   └── genkeys/        # Key generation utility
├── internal/
│   ├── api/            # HTTP handlers
│   ├── config/         # Configuration
│   ├── database/       # Database connection
│   ├── models/         # Data models
│   ├── repository/     # Database operations
│   └── service/        # Business logic
├── migrations/         # SQL migration files
├── pkg/
│   └── crypto/         # Cryptographic operations
└── keys/               # Generated keys (gitignored)
```

## Configuration

Environment variables (set in `.env` or as system environment):

- `DB_PATH`: Database file path (default: `data/taskmaster_license.db`)
- `JWT_SECRET`: Secret for JWT signing (default: `taskmaster-secret-key-change-in-production`)
- `API_PORT`: Server port (default: `8080`)
- `API_ENV`: Environment mode (default: `development`)
- `ROOT_PUBLIC_KEY`: Root public key for CML validation
- `ASTACK_MOCK_PORT`: Mock A-Stack server port (default: `8081`)

## Development Status

Currently in development. Core infrastructure is being implemented.
