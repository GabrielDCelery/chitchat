# Project: Distributed Chat System with TUI Client

## Project structure:

```sh
  chitchat/
  ├── cmd/
  │   ├── chitchat/        # TUI client
  │   └── chitchatd/       # Server daemon
  ├── internal/
  │   ├── server/
  │   ├── client/
  │   ├── protocol/
  │   └── tui/
  ├── pkg/                 # Public APIs if you want library usage
  ├── docker-compose.yml
  ├── mise.toml
  ├── CLAUDE.md
  └── README.md
```

## Architecture

Backend (Go):

- WebSocket server handling connections
- Room management and presence tracking
- Message persistence (PostgreSQL/SQLite)
- User authentication/sessions
- Horizontal scaling with Redis pub/sub

## TUI Client (Go):

- Beautiful terminal interface (bubbletea + lipgloss)
- Multiple views: room list, chat, user list, settings
- Real-time message streaming
- Notifications, typing indicators
- Vim-style keybindings

## Implementation Plan

### Week 1 - Core functionality:

- WebSocket server with basic room/message handling
- Database schema and persistence layer
- Simple TUI that connects and sends/receives messages
- User authentication (JWT tokens)
- Room creation/joining

### Week 2 - Real-time features:

- Presence tracking (who's online, typing indicators)
- Redis pub/sub for multi-server deployment
- Direct messages
- Message history with pagination
- File/image sharing (with TUI preview for images using sixel/kitty protocols)

### Week 3 - Polish:

- Beautiful TUI with multiple themes
- Search functionality (messages, users, rooms)
- Comprehensive tests (especially WebSocket connection handling)
- Docker compose with server + Redis + Postgres
- Admin commands (/kick, /mute, etc)
- Excellent documentation and demo GIFs

### Why This Showcases Skills

- Concurrency: Managing thousands of WebSocket connections
- Distributed systems: Multi-server scaling with pub/sub
- Protocol design: WebSocket message format, event handling
- TUI development: Shows UI/UX thinking even in terminal
- Real-time systems: Handling connection failures, reconnection logic
- Testing: WebSocket testing, concurrent access patterns

### Bonus Points

- E2E encryption option
- Message reactions/threads (Slack-like)
- Bot API for integrations
- CLI tool alongside TUI (scriptable interface)
- Performance metrics (messages/sec, connection limits)

---

## Development Log

### Session 1 - 2026-01-13

#### What We Built:

- Initialized Go module: `github.com/GabrielDCelery/chitchat`
- Created project directory structure (cmd/, internal/, pkg/)
- Implemented `internal/protocol/message.go`:
  - `MessageType` enum using `iota` (Chat, Join, Leave, Typing)
  - `Message` struct with Type, Sender, Room, Content, Timestamp, Metadata fields
  - Custom JSON marshaling/unmarshaling for MessageType
  - Educational comments following Go conventions
- Wrote comprehensive tests in `internal/protocol/message_test.go`:
  - Table-driven tests for MarshalJSON
  - Table-driven tests for UnmarshalJSON
  - Error handling tests for invalid types
  - All tests passing ✓

#### Key Decisions:

**Protocol Design:**

- Using JSON encoding initially (easy to debug, human-readable)
- Designed with extensibility in mind (Metadata field for future features)
- Architecture supports swapping to gRPC later via Encoder interface pattern
- Content field is optional (pointer type) for flexibility

**Future Considerations:**

- Collaborative markdown editing would require Operational Transformation or CRDTs
- Simple file sharing can be added via Metadata without protocol changes
- Migration path from JSON to gRPC preserved via interface-based design

**Documentation Strategy:**

- Manual conversation snapshots in `ai-chat/` directory
- Educational comments embedded in code
- Session summaries in this file

#### Go Patterns Learned:

1. Custom JSON marshaling via `MarshalJSON()` and `UnmarshalJSON()` methods
2. Interface satisfaction (`json.Marshaler` and `json.Unmarshaler`)
3. Table-driven tests for comprehensive coverage
4. `iota` for type-safe enums with explicit type declaration
5. Pointer fields for optional struct members
6. `omitempty` struct tag for optional JSON fields

#### Next Steps:

- Implement `Encoder` interface (encoding.go) for transport abstraction
- Create `JSONEncoder` implementation
- Test full `Message` struct marshaling/unmarshaling
- Begin server package implementation

---

### Session 2 - 2026-01-14

#### What We Built:

- Implemented `internal/protocol/encoding.go`:
  - `Encoder` interface for transport abstraction (JSON, protobuf, etc.)
  - Defines `Encode(io.Writer, *Message)` and `Decode(io.Reader, *Message)` methods
  - Enables swapping encoding formats without changing server/client code
- Implemented `JSONEncoder` in `encoding.go`:
  - Uses `json.NewEncoder(w).Encode()` for streaming (idiomatic Go)
  - Uses `json.NewDecoder(r).Decode()` for reading
  - Clean, single-line implementations
- Wrote comprehensive tests in `internal/protocol/encoding_test.go`:
  - Encode tests: success case, writer failure (using mock `failWriter`)
  - Decode tests: success case, invalid message type
  - Test helpers: `ptr()` for string pointers, `mustParseTime()` for timestamps
  - Used `bytes.Buffer` for in-memory testing
  - Handled `json.Encoder` newline behavior with `strings.TrimSpace()`
  - All tests passing ✓

#### Key Decisions:

**Encoder Abstraction:**

- Interface-based design allows protocol flexibility
- Works with any `io.Reader`/`io.Writer` (WebSocket, TCP, UDP)
- Future-proof: Can add protobuf, MessagePack, or custom encodings
- Enables P2P extension later (same encoder for WebSocket signaling and UDP data)

**Testing Strategy:**

- Used `bytes.Buffer` as both reader and writer for testing
- Created mock `failWriter` to test error handling
- Helper functions (`ptr()`, `mustParseTime()`) make tests readable
- Comprehensive coverage: happy paths and error cases

**Architecture Insights:**

- Discussed P2P networking: NAT traversal, hole punching, relay servers
- Current server-based design naturally extends to P2P signaling server
- Transport-agnostic protocol enables future hybrid P2P/relay architecture

#### Go Patterns Learned:

1. **Pointer vs Value Receivers**:
   - Value receiver: When only reading fields or struct is empty
   - Pointer receiver: When mutating fields or struct is large
   - For immutable config, value receiver is clearer
2. **io.Writer Contract**: Must return non-nil error if `n < len(p)` (prevents silent data loss)
3. **Testing with io Interfaces**: `bytes.Buffer` implements both `io.Reader` and `io.Writer`
4. **Mock Objects**: Creating test doubles (like `failWriter`) for error path testing
5. **Test Helpers**: Small functions (`ptr()`, `mustParseTime()`) improve test readability
6. **Forward Declarations**: Can reference types in same package defined later

#### Server Architecture Planning:

Started planning WebSocket server with three core components:

**1. Client** (per-connection state):

- Represents a connected user with WebSocket connection
- Two goroutines: `readPump()` (reads from WS) and `writePump()` (writes to WS)
- Buffered `send` channel for non-blocking message delivery
- Implements heartbeat/keepalive with ping/pong

**2. Room** (chat room management):

- Groups clients by room ID
- Thread-safe client map with `sync.RWMutex`
- Broadcasts messages to all clients in room

**3. Server** (hub pattern):

- Central coordinator running single `Run()` goroutine
- Manages rooms and clients via channels (register, unregister, broadcast)
- Non-blocking, concurrent-safe message routing
- Uses `protocol.Encoder` for serialization

**Key Patterns**:

- Hub pattern with channel-based communication
- Goroutine per client for independent operation
- Buffered channels prevent slow clients from blocking others
- Graceful cleanup on disconnect

#### Next Steps:

- Implement `Client` struct and methods (`readPump`, `writePump`)
- Implement `Room` struct for managing chat rooms
- Implement `Server` struct with hub logic
- Create WebSocket upgrade handler
- Write tests for server components
