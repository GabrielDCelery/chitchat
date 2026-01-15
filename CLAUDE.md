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

---

### Session 3 - 2026-01-15

#### What We Built:

- Completed `Client` struct implementation in `internal/server/client.go`:
  - Implemented `writePump()` with proper error handling and graceful shutdown
  - Used `NextWriter()` for streaming WebSocket writes
  - Configured ping/pong mechanism for connection keepalive
  - Implemented client-specific logger with baked-in context fields
  - Ready to implement `readPump()` following guidance
- Fixed critical bugs:
  - Interface pointer error: Changed `encoder *protocol.Encoder` to `encoder protocol.Encoder`
  - Added missing `return` statements after write errors in `writePump`
  - Properly handled WebSocket encoding via `NextWriter` or buffer pattern

#### Key Decisions:

**Client Implementation:**

- Used `NextWriter()` for streaming writes (memory efficient)
- Set write deadlines before each write operation (prevents hanging)
- Implemented graceful shutdown with WebSocket close frames
- Logged close message errors but continued cleanup (best effort)
- Client-specific logger created via `logger.With()` for cleaner code

**Ping/Pong Mechanism:**

- Server sends PING every 54 seconds (via `writePump`)
- Client automatically responds with PONG (WebSocket library handles this)
- Server receives PONG and resets read deadline (via `pongHandler`)
- If no PONG in 60 seconds, connection is considered dead
- 6-second response window is generous for network latency

**Read vs Write Deadlines:**

- Write deadline: Set before EACH write (discrete operations)
- Read deadline: Set once, reset on activity (rolling timeout)
- Both use absolute time, not duration per operation

**WebSocket Encoding:**

- WebSocket connections don't implement `io.Writer` directly
- Two solutions: `NextWriter()` (streaming) or buffer-based
- User chose `NextWriter()` for better memory efficiency
- Both approaches bridge the gap between WebSocket frames and `io.Writer`

**Logging Pattern:**

- Store client-specific logger with fields baked in
- Use `logger.With()` to create contextualized logger in `NewClient()`
- Eliminates repetitive field logging (clientID, username)
- Can still add per-log fields when needed
- Follows zap's idiomatic design philosophy

#### Go Patterns Learned:

1. **Never use pointers to interfaces**: `encoder protocol.Encoder` not `*protocol.Encoder`
2. **Interfaces are already references**: Contain type info pointer and data pointer
3. **Write deadlines prevent hanging**: Must set before each write to avoid resource exhaustion
4. **Read deadlines are rolling**: Reset on any activity, not per-message
5. **Graceful vs abrupt close**:
   - `WriteMessage(CloseMessage)`: Sends close frame (polite)
   - `conn.Close()`: Closes TCP socket (cleanup)
   - Best practice: Use both (send frame, then defer closes socket)
6. **Close message errors are OK**: Log but don't fail, connection might already be dead
7. **WebSocket frame-based messaging**: Use `NextWriter`/`ReadMessage` or buffer pattern
8. **Zap logger composition**: Use `logger.With()` for context-specific loggers
9. **Error handling strategy**:
   - `return` for connection errors (writePump)
   - `continue` for message errors (readPump - more forgiving)
10. **Security principle**: Never trust client-provided identity (override `Sender` field)

#### Architecture Insights:

**Terminology Clarification:**

- The `Client` struct runs on the SERVER side
- It represents the server's view of a connected user
- Better names could be: `Connection`, `ClientConnection`, or `Session`
- But `Client` is conventional in WebSocket examples

**Connection Lifecycle:**

```
1. WebSocket upgrade creates Client struct
2. Two goroutines start: readPump() and writePump()
3. readPump: WebSocket → Server (receive messages)
4. writePump: Server → WebSocket (send messages)
5. On disconnect: both goroutines exit, defer cleanup runs
```

**Ping/Pong Timeline:**

```
Time 0s:    Set read deadline = now + 60s
Time 54s:   Send PING
Time 54.1s: Receive PONG, reset deadline = now + 60s (114.1s total)
Time 108s:  Send PING
Time 108.1s: Receive PONG, reset deadline = now + 60s (168.1s total)
(Pattern continues...)
```

**Rolling Deadline Pattern:**

- Any activity resets the timer (pongs, chat messages, etc.)
- Not "did this specific ping succeed?" but "is connection alive?"
- More efficient than tracking individual pings
- Standard idle timeout pattern in network protocols

#### Implementation Status:

**Completed:**
- ✅ `Client` struct with proper fields and logger
- ✅ `writePump()` with streaming writes via `NextWriter()`
- ✅ Ping mechanism for keepalive
- ✅ Graceful shutdown with close frames
- ✅ Error handling and logging
- ✅ Fixed interface pointer bug in `Server`

**In Progress:**
- 🔄 `readPump()` (guidance provided, ready to implement)

**Next:**
- Implement `readPump()` with message decoding and validation
- Implement `Room` struct for managing chat rooms
- Complete `Server` struct with hub logic
- Create WebSocket upgrade handler
- Write tests for server components

#### Detailed Conversation Topics:

See `ai-chat/` directory for full conversation details:
- `24-session-3-ping-pong-deep-dive.md`: Ping/pong mechanism explained
- `25-websocket-encoder-and-logging.md`: Interface errors, WebSocket encoding, logger patterns
- `26-writepump-review-and-readpump-guide.md`: writePump review, deadlines, close messages, readPump guide

#### Next Steps:

- Implement `readPump()` following the provided guide
- Implement `Room` struct for chat room management
- Complete `Server` struct with hub logic and channels
- Create WebSocket upgrade handler
- Write tests for Client, Room, and Server components
