# Rockets Telemetry API

A REST API service built with Go that receives rocket telemetry messages and maintains the current state of all rockets in the system.

## Features

- ✅ **Async message processing** with Go channels (pub/sub pattern)
- ✅ **Non-blocking HTTP responses** (returns 202 immediately)
- ✅ Receives and processes rocket telemetry messages asynchronously
- ✅ Maintains current state for all rockets
- ✅ Handles out-of-order message delivery
- ✅ Prevents duplicate message processing
- ✅ **Interface-based repository** pattern (easy to swap storage)
- ✅ **Thread-safe** operations with mutex protection
- ✅ REST API with sorting capabilities
- ✅ Interactive Swagger documentation
- ✅ **Graceful shutdown** handling
- ✅ Built with Go best practices
- ✅ Production-ready Makefile
- ✅ Comprehensive error handling

## Quick Start

### Prerequisites

- Go 1.21 or later
- Make (optional, but recommended)

### Installation

```bash
# Clone or navigate to the project directory
cd /Users/albherna/Desktop/projects/rockets

# Download dependencies
go mod download
```

### Running the Server

```bash
# Option 1: Using Make (recommended)
make run

# Option 2: Direct execution
go run cmd/server/main.go

# Option 3: Build and run binary
make build
./bin/rockets
```

The server will start on **http://localhost:8088**

### Testing with the Rockets Program

Once your server is running, launch the provided test program:

```bash
./rockets launch "http://localhost:8088/messages" --message-delay=500ms --concurrency-level=1
```

## API Documentation

### Swagger UI

Access the interactive Swagger documentation at:
```
http://localhost:8088/swagger/index.html
```

### Endpoints

#### 1. Health Check
```bash
GET /health
```
Returns the health status of the service.

**Example:**
```bash
curl http://localhost:8088/health
```

#### 2. Post Messages
```bash
POST /messages
```
Receives rocket telemetry messages (used by the test program).

**Example:**
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

#### 3. List All Rockets
```bash
GET /rockets?sort={type|speed|mission|status}
```
Retrieves all rockets with optional sorting.

**Examples:**
```bash
# Get all rockets
curl http://localhost:8088/rockets

# Sort by type
curl "http://localhost:8088/rockets?sort=type"

# Sort by speed (descending)
curl "http://localhost:8088/rockets?sort=speed"

# Sort by mission
curl "http://localhost:8088/rockets?sort=mission"

# Sort by status
curl "http://localhost:8088/rockets?sort=status"
```

#### 4. Get Specific Rocket
```bash
GET /rockets/{id}
```
Retrieves the current state of a specific rocket by ID (channel).

**Example:**
```bash
curl http://localhost:8088/rockets/193270a9-c9cf-404a-8f83-838e71d9ae67
```

## Message Types

The system supports the following telemetry message types:

### RocketLaunched
Initializes a new rocket in the system.
```json
{
  "metadata": {
    "channel": "uuid",
    "messageNumber": 1,
    "messageTime": "2024-01-01T12:00:00Z",
    "messageType": "RocketLaunched"
  },
  "message": {
    "type": "Falcon-9",
    "launchSpeed": 500,
    "mission": "ARTEMIS"
  }
}
```

### RocketSpeedIncreased
Increases the rocket's speed.
```json
{
  "metadata": {...},
  "message": {
    "by": 3000
  }
}
```

### RocketSpeedDecreased
Decreases the rocket's speed.
```json
{
  "metadata": {...},
  "message": {
    "by": 2500
  }
}
```

### RocketExploded
Marks a rocket as exploded.
```json
{
  "metadata": {...},
  "message": {
    "reason": "PRESSURE_VESSEL_FAILURE"
  }
}
```

### RocketMissionChanged
Updates the rocket's mission.
```json
{
  "metadata": {...},
  "message": {
    "newMission": "SHUTTLE_MIR"
  }
}
```

## Development

### Makefile Commands

```bash
# Display help
make help

# Generate Swagger documentation
make swagger

# Run linter (requires golangci-lint)
make lint

# Build the application
make build

# Run the application
make run

# Run in development mode (with swagger generation)
make dev

# Clean build artifacts
make clean

# Download and tidy dependencies
make deps
```

### Project Structure

```
.
├── cmd/
│   └── server/
│       └── main.go              # Application entry point & wiring
├── internal/
│   ├── api/
│   │   ├── handler.go           # HTTP handlers (publishes to channel)
│   │   └── router.go            # Route configuration
│   ├── models/
│   │   └── rocket.go            # Data models and types
│   ├── pubsub/
│   │   └── pubsub.go            # Pub/Sub interfaces & channel implementation
│   ├── repository/
│   │   └── repository.go        # Repository interface & in-memory impl
│   └── service/
│       └── rocket_service.go    # Business logic & async message processor
├── docs/                        # Swagger documentation (generated)
├── bin/                         # Compiled binaries (generated)
├── Makefile                     # Build automation
├── go.mod                       # Go module definition
├── go.sum                       # Dependency checksums
├── README.md                    # Project overview
├── ARCHITECTURE.md              # Async architecture explained
├── QUICKSTART.md               # Quick start guide
├── SOLUTION.md                 # Design decisions and trade-offs
└── .golangci.yml               # Linter configuration
```

### Architecture Layers

1. **HTTP Layer** (`internal/api`): Receives requests, validates, publishes to channel
2. **Pub/Sub Layer** (`internal/pubsub`): Channel-based async messaging
3. **Service Layer** (`internal/service`): Subscribes & processes messages in background
4. **Repository Layer** (`internal/repository`): Thread-safe storage with mutex
5. **Models** (`internal/models`): Shared data structures

See [ARCHITECTURE.md](ARCHITECTURE.md) for detailed explanation.

## Configuration

The application can be configured using environment variables:

```bash
# Change the server port (default: 8088)
PORT=9000 ./bin/rockets

# Set Gin mode: debug, release, or test (default: debug)
GIN_MODE=release ./bin/rockets

# Combine multiple variables
PORT=9000 GIN_MODE=release ./bin/rockets
```

## Design Decisions & Trade-offs

See [SOLUTION.md](SOLUTION.md) for detailed documentation on:
- Architecture decisions
- Trade-offs and considerations
- Production readiness improvements
- Time estimation

## Key Features Explained

### Async Pub/Sub Architecture
The application uses Go channels for asynchronous message processing:
- **HTTP Handler** publishes messages to a buffered channel and returns 202 immediately
- **Background Goroutine** subscribes to the channel and processes messages
- **Repository** stores state with mutex protection
- **Non-blocking**: HTTP responses don't wait for message processing

See [ARCHITECTURE.md](ARCHITECTURE.md) for detailed explanation with diagrams.

### Out-of-Order Message Handling
The service tracks the last processed message number for each rocket and ignores messages with lower message numbers, ensuring state consistency. Messages are validated at processing time, not receipt time.

### Duplicate Prevention
Messages with the same or lower message numbers are automatically skipped during processing, providing idempotent message processing.

### Concurrency Safety
All state mutations are protected by a mutex in the repository layer, ensuring thread-safe operations in a concurrent environment.

### Interface-Based Design
- **Publisher/Subscriber interfaces**: Easy to swap channel implementation for Kafka/RabbitMQ
- **Repository interface**: Can replace in-memory storage with database
- **Testable**: Each component can be tested independently with mocks

### In-Memory Storage
For simplicity and speed, the current implementation uses in-memory storage with mutex protection. See [SOLUTION.md](SOLUTION.md) for production persistence recommendations.

## Testing

### Manual Testing

1. Start the server:
```bash
make run
```

2. In another terminal, run the rockets test program:
```bash
./rockets launch "http://localhost:8088/messages" --message-delay=500ms --concurrency-level=1
```

3. Check the rocket states:
```bash
# List all rockets
curl http://localhost:8088/rockets

# Get a specific rocket
curl http://localhost:8088/rockets/{rocket-id}
```

### Expected Behavior

- Messages are processed and stored correctly
- Out-of-order messages are handled properly
- Duplicate messages are ignored
- Rocket state reflects the latest processed message
- API returns correct data with proper sorting

## Troubleshooting

### Port Already in Use
```bash
PORT=9000 make run
```

### Swagger Documentation Not Loading
```bash
make swagger
```

### Dependencies Issues
```bash
go mod download
go mod tidy
```

### Build Errors
```bash
# Clean and rebuild
make clean
make build
```

## Production Considerations

This implementation is suitable for a 6-hour coding challenge. For production deployment, consider:

1. **Persistent Storage**: Add database layer (PostgreSQL, MongoDB)
2. **Message Queue**: Use Kafka/RabbitMQ for reliable messaging
3. **Monitoring**: Add Prometheus metrics and logging
4. **Authentication**: Implement API authentication
5. **Containerization**: Create Docker images and Kubernetes configs
6. **Testing**: Add comprehensive unit and integration tests
7. **CI/CD**: Implement automated testing and deployment

See [SOLUTION.md](SOLUTION.md) for detailed production recommendations.

## License

MIT License - See LICENSE file for details

## Author

Built as part of a backend engineering coding challenge.

## Support

For questions or issues:
1. Review the [SOLUTION.md](SOLUTION.md) documentation
2. Check the [QUICKSTART.md](QUICKSTART.md) guide
3. Explore the API via Swagger UI at http://localhost:8088/swagger/index.html

