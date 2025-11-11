# 🚀 Rockets API - Quick Reference Card

## 📋 At a Glance
- **Language**: Go 1.21+
- **Framework**: Gin
- **Port**: 8088 (configurable)
- **Status**: ✅ Complete & Tested

## ⚡ Quick Commands

```bash
# Build
make build

# Run
./bin/rockets

# Test
./rockets launch "http://localhost:8088/messages" --message-delay=500ms --concurrency-level=1

# Query
curl http://localhost:8088/rockets
```

## 🔗 API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/messages` | Post telemetry message |
| GET | `/rockets` | List all rockets |
| GET | `/rockets/:id` | Get specific rocket |
| GET | `/health` | Health check |
| GET | `/swagger/*` | Swagger UI |

## 📊 Sorting Options

```bash
curl "http://localhost:8088/rockets?sort=type"     # By rocket type
curl "http://localhost:8088/rockets?sort=speed"    # By speed (desc)
curl "http://localhost:8088/rockets?sort=mission"  # By mission
curl "http://localhost:8088/rockets?sort=status"   # By status
```

## 📝 Message Types

1. **RocketLaunched** - Initialize rocket
2. **RocketSpeedIncreased** - Increase speed
3. **RocketSpeedDecreased** - Decrease speed
4. **RocketExploded** - Mark as exploded
5. **RocketMissionChanged** - Update mission

## 📚 Documentation

- **README.md** - Challenge + Solution overview
- **PROJECT_README.md** - Complete API docs
- **SOLUTION.md** - Design & trade-offs
- **QUICKSTART.md** - Setup guide
- **CHECKLIST.md** - Completion checklist
- **Swagger** - http://localhost:8088/swagger/index.html

## 🛠️ Make Commands

```bash
make help      # Show all commands
make build     # Build application
make run       # Run application
make swagger   # Generate docs
make clean     # Clean artifacts
make deps      # Update dependencies
make lint      # Run linter
```

## 🔧 Configuration

```bash
PORT=9000 ./bin/rockets              # Custom port
GIN_MODE=release ./bin/rockets       # Release mode
```

## ✨ Key Features

- ✅ Out-of-order message handling
- ✅ Duplicate message prevention
- ✅ Thread-safe operations
- ✅ RESTful API design
- ✅ Swagger documentation
- ✅ Sortable results

## 📁 Project Structure

```
cmd/server/main.go           # Entry point
internal/
  ├── api/                   # HTTP layer
  ├── models/                # Data models
  └── service/               # Business logic
docs/                        # Swagger (generated)
```

## 🎯 Example Workflow

```bash
# 1. Start server
make build && ./bin/rockets

# 2. Post a message
curl -X POST http://localhost:8088/messages \
  -H "Content-Type: application/json" \
  -d '{
    "metadata": {
      "channel": "rocket-1",
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

# 3. Verify
curl http://localhost:8088/rockets

# 4. Get specific rocket
curl http://localhost:8088/rockets/rocket-1
```

## 🐛 Troubleshooting

| Issue | Solution |
|-------|----------|
| Port in use | `PORT=9000 ./bin/rockets` |
| Build fails | `go mod tidy && make build` |
| Swagger 404 | `make swagger` |

## 📖 Learn More

1. Start with **QUICKSTART.md**
2. Read **PROJECT_README.md** for API details
3. Review **SOLUTION.md** for design decisions
4. Explore **Swagger UI** for interactive docs

---

**Status: Ready for Review** ✅

