# 🪐 Backend Engineer Challenge: Rockets 🚀

---

## SOLUTION IMPLEMENTED

This is my solution to the Rockets challenge. I've tried to balance simplicity with production-ready patterns, keeping in mind the 6-hour time constraint while building something that could actually scale with the right tweaks.

### How to Run

```bash
# Build the server
make build

# Start it
./bin/rockets

# In another terminal, test with the provided program
./rockets launch "http://localhost:8088/messages" --message-delay=500ms --concurrency-level=1
```

The server runs on port 8088 by default. You can change it with `PORT=9000 ./bin/rockets`.

### What I Built

The solution is a REST API that ingests rocket telemetry messages asynchronously and maintains the current state of all rockets in memory. The architecture separates concerns into three main pieces:

1. **HTTP handlers** - Receive messages and validate input, return rocket data on query endpoints
2. **Message processing service** - Consumes messages from a pub/sub channel, handles ordering and deduplication, updates rocket state. Designed to be swapped out for a real message queue later.
3. **Repository layer** - Abstracts storage behind an interface, currently in-memory but designed to swap in PostgreSQL or similar

**API Endpoints:**
- `POST /messages` - Accepts rocket messages (returns 202 immediately, processes async)
- `GET /rockets` - Lists all rockets with optional sorting (`?sort=type|speed|mission|status`)
- `GET /rockets/:id` - Gets a specific rocket by channel UUID
- `GET /health` - Health check (thought useful to have for monitoring)

### Design Decisions and Trade-offs

**Async Pub/Sub with Channels**

I went with a channel-based pub/sub pattern where HTTP handlers publish messages to a buffered channel and a background goroutine processes them. This means the HTTP response returns immediately without waiting for processing, which is better for throughput.

The trade-off here is that if the channel fills up (buffer size is 1000), messages get dropped. For this assignment that seemed reasonable, but in production you'd want proper monitoring on dropped messages and probably a real message queue like Redis Streams or RabbitMQ to be able to scale. The pub/sub interface is designed to make that swap easy - it uses callbacks instead of exposing channels, so a Redis or Kafka implementation would fit right in.

**In-Memory Storage**

The repository uses an in-memory map protected by a RWMutex. This is fast and simple for the assignment scope. The obvious downside is that all data is lost on restart and you can't run multiple instances.

For production, you'd replace this with a proper database. I structured it as a repository interface so swapping to PostgreSQL would just mean implementing the same interface with database calls. The mutex logic would move to database transactions.

**Out-of-Order and Duplicate Handling**

Each rocket tracks the last processed message number. When a new message arrives, if its number is less than or equal to the last processed one, it gets ignored. This handles both out-of-order delivery and duplicates with a simple comparison.

The trade-off is that this only works because each rocket has its own independent message stream (channel ID). If messages could arrive across different channels for the same rocket, you'd need more complex event sourcing or a proper event store.

**Message Type Consolidation**

I merged `RocketSpeedIncreased` and `RocketSpeedDecreased` into a single `RocketSpeedChangedMessage` type since they're structurally identical (just a "by" field). The handler checks the message type to determine whether to add or subtract the speed delta. This reduced code duplication in the handlers.

**Generic Helper Functions**

To avoid repeating the same marshal/unmarshal logic in every message handler, I used a generic `parseMessage[T]` function. This cuts down on boilerplate and makes it easier to maintain. The marshal/unmarshal step is necessary because Gin unmarshals JSON into `map[string]interface{}` and we need typed structs.

**Validation**

Input validation happens at the HTTP layer before publishing. I validate that UUIDs are valid UUIDs, message types are recognized, required fields are present, etc. This prevents bad data from even entering the system. The validator functions are separated from the handlers following single responsibility principle.

**No Tests**

I didn't write automated tests given the time constraint and the instructions saying it's okay to submit what you have. In a real project, I'd have unit tests for the message handlers, integration tests with testcontainers, and probably some property-based tests for the ordering logic.

### What I'd Change for Production

**Persistence:** Replace in-memory repository with PostgreSQL or similar. Add proper migrations.

**Message Queue:** Swap channel pub/sub for Redis Streams or RabbitMQ. Add dead-letter queues for failed messages.

**Observability:** Add structured logging (probably zap or slog), Datadog metrics for message processing latency and error rates, distributed tracing with Datadog as well.

**Resilience:** Add retry logic with exponential backoff for transient errors, circuit breakers for external dependencies, rate limiting on the HTTP endpoints.

**Configuration:** Use proper config management (Viper or similar) instead of just environment variables.

**Deployment:** Containerize with Docker, add Kubernetes manifests, proper health checks.

**Monitoring:** Track channel buffer utilization, alert if it's consistently full, monitor duplicate message rate to detect issues.

**Testing:** Comprehensive unit and integration test suite, load testing to verify the async pattern holds up under pressure.

### Architecture Overview

The flow is: HTTP request comes in, handler validates it and publishes to channel, background goroutine picks it up, processes it according to the message type, updates the repository. Query endpoints read directly from the repository.

Services are cleanly separated - `MessageService` owns the async processing, `RocketService` owns the query logic. Neither knows about the other. Both depend on the repository interface.

This separation means you could theoretically run the message processor and the query API as separate processes if needed for scaling, though that wasn't a requirement here.

**Project Structure:**
```
cmd/server/          - Application entry point
internal/
  handler/          - HTTP request handlers (separated by concern)
  service/          - Business logic (message and rocket services)
  repository/       - Storage abstraction (interface + in-memory impl)
  pubsub/          - Pub/sub abstraction (interface + channel impl)
  models/          - Data structures
  api/             - Router setup
docs/              - Swagger documentation (auto-generated)
```

The project follows standard Go layout conventions. Dependencies are explicit and injected through constructors. Concurrency is visible - you can see the goroutine starts in main.go.

---

# 🪐 Backend Engineer Challenge: Rockets 🚀 (Original Challenge Description Below)

## Introduction 👋
Thank you for taking Lunar's code challenge for backend engineers! 

In the ZIP-file you have received, you will find a `README.md` (dah! of course) and folders 
containing executables for various operating systems and architectures.

> **Important:** If you cannot find an executable that works for you please reach out to us as soon as possible, 
> so we can get you one that works.

We hope you will enjoy this challenge - good luck.

## The Challenge 🧑‍💻
In this challenge you are going to build a service (or multiple) which consumes messages 
from a number of entities – i.e. a set of _rockets_ – and make the state of these 
available through a REST API. We imagine this API to be used by something like a dashboard.

As a minimum we expect endpoints which can:
1. Return the current state of a given rocket (type, speed, mission, etc.)
1. Return a list of all the rockets in the system; preferably with some kind of sorting.

The service should also expose an endpoint where the test program can post the messages to (see this [section](#running-the-test-program))

We are writing all our services in [Go](https://go.dev/) but there are no constrains on the programming language that you choose for 
solving the challenge. 
We prefer that you implement a great solution in a language that you feel comfortable in rather than trying to write 
in Go and implement a mediocre solution.

### The messages ✉️
Each rocket will be dispatching various messages (encoded as JSON) about its state changes through individual radio _channels_.
The channel is unique for each rocket and can therefore be treated as the ID of the rocket.

Apart from the channel each message also contains a _message number_ which expresses the order of the message within a channel, 
a _message time_ indicating when the message was sent and a _message type_ describing the event that occurred.

**Important:** Messages will arrive **out of order** and there is an **at-least-once guarantee** on messages 
meaning that you might receive the same message more than once.

Here is an example of a `RocketLaunch` message:

```json
{
    "metadata": {
        "channel": "193270a9-c9cf-404a-8f83-838e71d9ae67",
        "messageNumber": 1,    
        "messageTime": "2022-02-02T19:39:05.86337+01:00",                                          
        "messageType": "RocketLaunched"                             
    },
    "message": {                                                    
        "type": "Falcon-9",
        "launchSpeed": 500,
        "mission": "ARTEMIS"  
    }
}
```

The possible message types are:

#### `RocketLaunched`
Sent out once: when a rocket is launched for the first time.
```json
{
    "type": "Falcon-9",
    "launchSpeed": 500,
    "mission": "ARTEMIS"  
}
```

#### `RocketSpeedIncreased`
Continuously sent out: when the speed of a rocket is increased by a certain amount.
```json
{
    "by": 3000
}
```

#### `RocketSpeedDecreased`
Continuously sent out: when the speed of a rocket is decreased by a certain amount.
```json
{
    "by": 2500
}
```

#### `RocketExploded`
Sent out once: if a rocket explodes due to an accident/malfunction.
```json
{
    "reason": "PRESSURE_VESSEL_FAILURE"
}
```

#### `RocketMissionChanged`
Continuously sent out: when the mission for a rocket is changed.
```json
{
    "newMission":"SHUTTLE_MIR"
}
```

### Running the test program 💽
In the ZIP-file locate the executable that works for your system and run the following:

```bash
./rockets launch "http://localhost:8088/messages" --message-delay=500ms --concurrency-level=1
```

This launches the program which starts posting (request method: `POST`) messages to the URL provided with a delay of 500ms between each message.

To see all commands run `./rockets help` and for help on the `launch` command run `./rockets launch --help`.

> We are going to run the program against your solution with the default values.

### Your solution and our assessment 📝
Before submitting your solution please make sure that you have included all the necessary files/information for 
running and assessing your solution. You can either submit a ZIP-file or provide a link to an online version control provider like GitHub, GitLab or Bitbucket.

Any design of a software system as a solution to a given problem will be affected by the choices made between various 
trade-offs. When submitting your solution, we will be really excited if you have explicitly described the design choices
you made and which trade-offs they entail. If you consciously chose a certain design, but are well aware that a different
(but maybe more complex) solution exist not having a certain trade-off, then include this in your documentation. 

When reviewing your solution we are going to look at things such as:
- The documentation provided, i.e. is it clear how to run your service(s) and, perhaps, what considerations/shortcuts have you made.
- The overall design of your solution, e.g. how easy is the code to understand, can the service(s) scale and how maintainable your code is.
- The measures you have taken to verify that your code works, e.g. automated tests.

We do not expect you to spend more than **6 hours** on this challenge. 
If you do not succeed in completing everything, then submit what you have, so we have something to look at - that is much better than nothing! ☺️
