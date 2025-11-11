# 🪐 Backend Engineer Challenge: Rockets 🚀

---

## SOLUTION IMPLEMENTED

This is my solution to the Rockets challenge. I've tried to balance simplicity with production-ready patterns, keeping in mind the 6-hour time constraint while building something that could actually scale with the right tweaks.

### How to Run

**Prerequisites:**
- Go 1.21 or later ([download here](https://go.dev/dl/))
- Make (optional but recommended, usually pre-installed on macOS/Linux)
- The rockets test program provided in the challenge

**Minimum versions:**
- Go: 1.21+ (uses generics and other modern features)
- Make: Any recent version

**First-time setup (if you don't have Go):**

1. Install Go from https://go.dev/dl/
2. Verify installation:
   ```bash
   go version
   # Should output: go version go1.21.x or later
   ```

3. Clone or extract this project:
   ```bash
   cd /path/to/rockets
   ```

4. Download dependencies:
   ```bash
   go mod download
   ```

**Running the service:**

```bash
# Option 1: Using Make (recommended)
make build
./bin/rockets

# Option 2: Without Make
go build -o bin/rockets cmd/server/main.go
./bin/rockets

# Option 3: Run directly without building
go run cmd/server/main.go
```

**Testing with the rockets program:**

In another terminal, run the test program provided in the challenge:
```bash
./rockets launch "http://localhost:8088/messages" --message-delay=500ms --concurrency-level=1
```

**Configuration:**

The server runs on port 8088 by default. You can change it with an environment variable:
```bash
PORT=9000 ./bin/rockets
```

**Verify it's working:**
```bash
# Check health endpoint
curl http://localhost:8088/health
```

**Troubleshooting:**

If you get "command not found: make":
- Use the direct Go commands instead (Option 2 above)
- Or install Make: `brew install make` (macOS) or `apt-get install build-essential` (Linux)

If you get module errors:
```bash
go mod tidy
go mod download
```

### What I Built

The solution is a REST API that ingests rocket telemetry messages asynchronously and maintains the current state of all rockets in memory. The architecture separates concerns into three main pieces:

1. **HTTP handlers** - Receive messages and validate input, return rocket data on query endpoints
2. **Message processing service** - Consumes messages from a pub/sub channel, handles ordering and deduplication, updates rocket state. Designed to be swapped out for a real message queue later.
3. **Repository layer** - Abstracts storage behind an interface, currently in-memory but designed to swap in PostgreSQL or similar

**API Endpoints:**
- `POST /messages` - Accepts rocket messages
- `GET /rockets` - Lists all rockets with optional sorting (`?sort=type|speed|mission|status`)
- `GET /rockets/:id` - Gets a specific rocket by channel UUID
- `GET /health` - Health check (thought useful to have for monitoring)

### Design Decisions and Trade-offs

**Event-Driven with Async Processing**

I went with an event-driven approach where HTTP handlers publish messages to a channel and a background worker processes them. The API returns immediately without waiting for storage, which keeps it fast and responsive.

Why this pattern:
- Decouples message intake from processing (can scale them separately)
- Natural fit for telemetry streams (continuous events updating state)
- Shows how you'd build this in production with proper message queues

The trade-off is messages can be lost if the service crashes between accepting and processing them. The channel buffer is also limited to 1000 messages - if processing slows down and it fills up, new messages get rejected.

**In-Memory Storage**

Everything lives in memory with a mutex-protected map. Fast and simple for a demo, but obviously all data disappears on restart. Also means you can't run multiple instances.

For production you'd swap in PostgreSQL. I structured it as a repository interface so that's a straightforward change - just implement the same interface with database calls instead of map operations.

**Handling Out-of-Order Messages**

Each rocket tracks its last processed message number. When a message arrives, if its number is ≤ the last one we saw, we ignore it. This handles both out-of-order delivery and duplicates in one shot.

Works well because each rocket has its own independent stream (channel ID). If messages for the same rocket could arrive on different channels, you'd need something more complex.

**Testing**

Included a test suite for the `GetRocket` handler to show the testing approach - table-driven tests with gomock for mocking. Didn't write full test coverage given the time constraint, but the pattern would be the same for other handlers.

### What's Missing for Production

The big ones:
- **Database**: Swap in PostgreSQL instead of in-memory storage
- **Real queue**: Use Redis Streams or RabbitMQ instead of Go channels (with persistence and horizontal scaling)
- **Observability**: Structured logging, metrics, distributed tracing
- **Tests**: Full test coverage, integration tests, load testing

This was built for a 6-hour scope - enough to show the architecture and patterns, but you'd need these pieces before going live.

### Architecture Overview

The flow is: HTTP request comes in, handler validates it and publishes to channel, background goroutine picks it up, processes it according to the message type, updates the repository. Query endpoints read directly from the repository.

Services are cleanly separated - `MessageService` owns the async processing, `RocketService` owns the query logic. Neither knows about the other. Both depend on the repository interface.

This separation means you could theoretically run the message processor and the query API as separate processes if needed for scaling, though that wasn't a requirement here.

The project follows standard Go layout conventions. Dependencies are explicit and injected through constructors. Concurrency is visible - you can see the goroutine starts in main.go.

### Time Spent

Roughly 6 hours total, spent mostly on getting the async pattern working cleanly, building the validation layer, and making sure the interfaces are swappable for production. The Swagger setup and test suite took a bit longer than expected but felt worth it for demonstration purposes.

---