# Solution Documentation

## Overview

This solution implements a REST API service that receives rocket telemetry messages and maintains the current state of all rockets in the system. The service is built with Go using the Gin web framework.

## Architecture & Design Decisions

### 1. **Async Pub/Sub Architecture with Channels**
**Decision**: Use Go channels for asynchronous message processing with separate goroutines.

**Architecture Flow**:
```
HTTP Request → Handler → Publisher (channel) → Subscriber (goroutine) → Repository
```

**Rationale**:
- Decouples HTTP handling from message processing
- Non-blocking HTTP responses (true async)
- Idiomatic Go pattern using channels
- Easy to test and reason about
- Naturally handles backpressure with buffered channels

**Components**:
- **Publisher**: HTTP handler publishes messages to channel (non-blocking)
- **Subscriber**: Background goroutine processes messages from channel
- **Repository**: Interface-based storage with in-memory implementation
- **Mutex Protection**: Repository uses RWMutex for thread-safe access

**Trade-offs**:
- Slightly more complex than synchronous processing
- Messages in channel buffer are lost on crash
- Single-instance limitation (in-memory storage)

**Benefits**:
- HTTP endpoints return immediately (better performance)
- Message processing doesn't block HTTP threads
- Can easily swap repository implementation
- Clear separation of concerns

**Production considerations**:
- Replace channel with actual message queue (Kafka, RabbitMQ, SQS)
- Add persistent storage (PostgreSQL, MongoDB)
- Implement message persistence before acknowledgment
- Add retry logic and dead letter queues

### 2. **Repository Pattern with Interface**
**Decision**: Use repository interface with in-memory implementation.

**Rationale**:
- Separates data access from business logic
- Easy to test with mock repositories
- Simple to add new storage implementations
- Mutex protection isolated in repository layer

**Trade-offs**:
- Data is lost on service restart
- Limited to single-instance deployment
- No persistence across crashes

**Production considerations**:
- Implement PostgreSQL/MongoDB repository
- Add connection pooling
- Implement proper transaction handling
- Add caching layer (Redis)

### 2. **Message Ordering**
**Decision**: Use message number as the ordering mechanism, ignore messages with lower message numbers than already processed.

**Rationale**:
- Simple comparison-based approach
- Handles out-of-order delivery effectively
- Prevents processing older state over newer state

**Trade-offs**:
- Assumes message numbers are strictly sequential
- Gap in message numbers could be problematic
- No buffering of future messages if gaps exist

**Production considerations**:
- Implement a buffer/queue for out-of-order messages
- Add timeout mechanism to handle missing messages
- Consider using message timestamps as additional ordering criteria
- Implement dead letter queue for problematic messages

### 3. **Duplicate Message Handling**
**Decision**: Track last processed message number per channel and skip duplicates.

**Rationale**:
- Ensures idempotent message processing
- Simple comparison to detect duplicates
- Minimal memory overhead

**Production considerations**:
- Could use message IDs for more robust deduplication
- Consider bloom filters for memory-efficient duplicate detection at scale
- May need to handle partial duplicates (same number, different content)

### 4. **Concurrency Model**
**Decision**: Single mutex protecting the entire rocket state map.

**Rationale**:
- Simple to implement and reason about
- Prevents race conditions
- Adequate for expected load

**Trade-offs**:
- Coarse-grained locking limits concurrency
- All operations block on the same lock
- May become bottleneck under heavy load

**Production considerations**:
- Use per-rocket locks (sharded locking) for better concurrency
- Consider lock-free data structures
- Implement read-write locks (RWMutex) if reads dominate
- Use connection pooling and rate limiting

### 5. **No Persistence Layer**
**Decision**: Pure in-memory storage without any persistence.

**Trade-offs**:
- Fast and simple
- No durability guarantees
- Cannot recover from crashes

**Production considerations**:
- Add write-ahead log (WAL) for durability
- Implement snapshot mechanism for quick recovery
- Use event sourcing to replay messages
- Add database layer with proper transactions

### 6. **REST API Design**
**Decision**: Simple REST endpoints with JSON responses.

**Endpoints**:
- `POST /messages` - Receive rocket messages
- `GET /rockets` - List all rockets (with sorting)
- `GET /rockets/:id` - Get specific rocket state
- `GET /health` - Health check endpoint
- `GET /swagger/*` - Swagger documentation

**Rationale**:
- Follows REST conventions
- Easy to understand and use
- Standard HTTP status codes

**Production considerations**:
- Add pagination for list endpoint
- Implement filtering and advanced sorting
- Add rate limiting and authentication
- Consider GraphQL for more flexible queries
- Add versioning (e.g., /v1/rockets)

### 7. **Error Handling**
**Decision**: Basic error handling with HTTP status codes and JSON error responses.

**Production considerations**:
- Implement structured logging (e.g., with context)
- Add error tracking (Sentry, Rollbar)
- Create custom error types with error codes
- Add request tracing (OpenTelemetry)
- Implement circuit breakers for downstream services

### 8. **No Authentication/Authorization**
**Decision**: Open endpoints without authentication.

**Rationale**:
- Not specified in requirements
- Simplifies implementation for challenge
- Focus on core functionality

**Production considerations**:
- Add JWT or OAuth2 authentication
- Implement role-based access control (RBAC)
- Add API key validation for test program
- Implement rate limiting per client
- Add HTTPS/TLS support

### 9. **Sorting Implementation**
**Decision**: Support sorting by type, speed, and mission via query parameter.

**Rationale**:
- Simple and intuitive API
- Handles common use cases
- Easy to extend

**Production considerations**:
- Add support for multiple sort fields
- Implement ascending/descending order
- Add more sophisticated filtering
- Consider pre-computed indexes for common queries

## Project Structure

```
.
├── cmd/
│   └── server/
│       └── main.go              # Application entry point & wiring
├── internal/
│   ├── api/
│   │   ├── handler.go           # HTTP handlers (publishes to channel)
│   │   └── router.go            # Route definitions
│   ├── models/
│   │   └── rocket.go            # Data models
│   ├── pubsub/
│   │   └── pubsub.go            # Pub/Sub interfaces & channel implementation
│   ├── repository/
│   │   └── repository.go        # Repository interface & in-memory impl
│   └── service/
│       └── rocket_service.go    # Business logic & async processor
├── docs/                        # Swagger generated docs
├── Makefile                     # Build automation
├── go.mod                       # Go dependencies
├── README.md                    # Challenge description
└── SOLUTION.md                  # This file
```

### Architecture Layers

1. **HTTP Layer** (`internal/api`): Receives requests, validates input, publishes to channel
2. **Pub/Sub Layer** (`internal/pubsub`): Channel-based async messaging
3. **Service Layer** (`internal/service`): Subscribes to messages, processes them asynchronously
4. **Repository Layer** (`internal/repository`): Thread-safe data storage with mutex
5. **Models** (`internal/models`): Shared data structures

## Running the Application

### Prerequisites
- Go 1.21 or later
- Make (optional, but recommended)

### Setup
```bash
# Install development tools (swag, golangci-lint)
make install-tools

# Download dependencies
make deps
```

### Build and Run
```bash
# Build the application (includes swagger generation)
make build

# Run the binary
./bin/rockets

# Or run directly with Go
make run
```

### Development
```bash
# Run in development mode with swagger generation
make dev
```

### Other Commands
```bash
# Generate swagger documentation
make swagger

# Run linter
make lint

# Clean build artifacts
make clean
```

## Testing with the Rockets Program

Once the server is running (default port: 8088), launch the test program:

```bash
./rockets launch "http://localhost:8088/messages" --message-delay=500ms --concurrency-level=1
```

## API Documentation

Once the server is running, access the Swagger UI at:
```
http://localhost:8088/swagger/index.html
```

## Configuration

The server uses environment variables for configuration:
- `PORT`: Server port (default: 8088)
- `GIN_MODE`: Gin mode - debug, release, or test (default: debug)

Example:
```bash
PORT=9000 GIN_MODE=release ./bin/rockets
```

## Future Improvements for Production

### High Priority
1. **Persistent Storage**: Implement database layer for durability
2. **Message Queue**: Use Kafka/RabbitMQ for reliable message delivery
3. **Monitoring**: Add Prometheus metrics and Grafana dashboards
4. **Logging**: Structured logging with correlation IDs
5. **Authentication**: Secure API endpoints
6. **Containerization**: Docker and Kubernetes deployment

### Medium Priority
7. **Testing**: Comprehensive unit and integration tests
8. **CI/CD**: Automated testing and deployment pipeline
9. **Rate Limiting**: Protect against abuse
10. **Caching**: Add Redis for frequently accessed data
11. **Configuration Management**: Use proper config system (Viper)
12. **Graceful Shutdown**: Handle shutdown signals properly

### Low Priority
13. **WebSocket Support**: Real-time updates for dashboards
14. **GraphQL API**: More flexible querying
15. **Message Replay**: Ability to replay messages for recovery
16. **Audit Log**: Track all state changes
17. **Multi-region**: Support for geo-distributed deployments

## Time Spent (Approximate)

This solution was designed to be completable within 6 hours:
- Project setup and structure: 30 minutes
- Core message handling logic: 2 hours
- REST API implementation: 1.5 hours
- Swagger documentation: 30 minutes
- Makefile and tooling: 30 minutes
- Documentation: 1 hour

Total: ~6 hours

## Conclusion

This solution provides a solid foundation for the rocket telemetry system while acknowledging its limitations and providing clear paths for future improvements. The design prioritizes simplicity and clarity over premature optimization, which is appropriate for a time-constrained challenge while demonstrating understanding of production requirements.
module github.com/ahernandez9/rockets

go 1.21

require (
	github.com/gin-gonic/gin v1.9.1
	github.com/google/uuid v1.4.0
	github.com/swaggo/files v1.0.1
	github.com/swaggo/gin-swagger v1.6.0
	github.com/swaggo/swag v1.16.2
)

