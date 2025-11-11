# Rockets API - Quick Start Guide

## Prerequisites

- Go 1.21 or later
- Make (optional but recommended)

## Installation & Setup

### 1. Install Required Tools

```bash
make install-tools
```

This will install:
- `swag` - for generating Swagger documentation
- `golangci-lint` - for code linting

### 2. Download Dependencies

```bash
make deps
```

Or simply:
```bash
go mod download
```

## Running the Application

### Option 1: Using Make (Recommended)

```bash
# Build and run
make build
./bin/rockets

# Or run directly in development mode
make dev
```

### Option 2: Using Go Directly

```bash
# Generate swagger docs first
go install github.com/swaggo/swag/cmd/swag@latest
swag init -g cmd/server/main.go -o docs

# Run the application
go run cmd/server/main.go
```

The server will start on port `8088` by default.

## Testing with the Rockets Test Program

Once your server is running, use the provided rockets executable to send test messages:

```bash
./rockets launch "http://localhost:8088/messages" --message-delay=500ms --concurrency-level=1
```

## API Endpoints

### Swagger Documentation
- **URL**: http://localhost:8088/swagger/index.html
- Interactive API documentation with try-it-out functionality

### Core Endpoints

#### 1. Health Check
```bash
curl http://localhost:8088/health
```

#### 2. Post Messages (Used by test program)
```bash
curl -X POST http://localhost:8088/messages \
  -H "Content-Type: application/json" \
  -d '{
    "metadata": {
      "channel": "test-channel-1",
      "messageNumber": 1,
      "messageTime": "2024-01-01T12:00:00Z",
      "messageType": "RocketLaunched"
    },
    "message": {
      "type": "Falcon-9",
      "launchSpeed": 500,
      "mission": "ARTEMIS"
    }
  }'
```

#### 3. Get All Rockets
```bash
# Get all rockets
curl http://localhost:8088/rockets

# Sort by type
curl "http://localhost:8088/rockets?sort=type"

# Sort by speed
curl "http://localhost:8088/rockets?sort=speed"

# Sort by mission
curl "http://localhost:8088/rockets?sort=mission"
```

#### 4. Get Specific Rocket
```bash
curl http://localhost:8088/rockets/{rocket-id}
```

## Configuration

The application can be configured using environment variables:

```bash
# Change port
PORT=9000 ./bin/rockets

# Set Gin mode (debug, release, test)
GIN_MODE=release ./bin/rockets

# Combine both
PORT=9000 GIN_MODE=release ./bin/rockets
```

## Development Commands

### Generate Swagger Documentation
```bash
make swagger
```

### Run Linter
```bash
make lint
```

### Build Binary
```bash
make build
# Output: bin/rockets
```

### Clean Build Artifacts
```bash
make clean
```

## Project Structure

```
.
├── cmd/
│   └── server/
│       └── main.go              # Application entry point
├── internal/
│   ├── api/
│   │   ├── handler.go           # HTTP request handlers
│   │   └── router.go            # Route configuration
│   ├── models/
│   │   └── rocket.go            # Data models and types
│   └── service/
│       └── rocket_service.go    # Business logic
├── docs/                        # Swagger documentation (generated)
├── Makefile                     # Build automation
├── go.mod                       # Go module definition
├── QUICKSTART.md               # This file
└── SOLUTION.md                 # Design decisions and trade-offs
```

## Troubleshooting

### Port Already in Use
If port 8088 is already in use:
```bash
PORT=9000 make run
```

### Swagger Not Loading
Make sure to generate the docs:
```bash
make swagger
```

### Dependencies Not Found
Download dependencies:
```bash
go mod download
go mod tidy
```

## Next Steps

1. Read [SOLUTION.md](SOLUTION.md) for detailed design decisions and trade-offs
2. Explore the API using Swagger UI at http://localhost:8088/swagger/index.html
3. Test with the rockets executable
4. Review the code structure in the `internal/` directory

## Support

For questions or issues, please refer to the [SOLUTION.md](SOLUTION.md) documentation.
# Swagger Documentation

This directory contains auto-generated Swagger documentation.

The documentation is generated from code annotations using the `swag` tool.

To regenerate the documentation, run:
```bash
make swagger
```

Access the Swagger UI at: `http://localhost:8088/swagger/index.html` when the server is running.

