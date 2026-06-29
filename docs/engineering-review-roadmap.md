# Engineering Review Roadmap — Realtime Terminal Collaboration

## Project Overview (for context)

The project is a real-time collaborative terminal platform in Go. It follows a **modular monolith** with 4 layers: Domain → Application → Infrastructure → Transport. TASK-001 through TASK-014 are implemented (full backend). Frontend (TASK-015–018) and Hardening (TASK-019–021) are not started.

**Implemented files (15 Go files):**

| Layer | Files |
|---|---|
| `cmd/` | `main.go` |
| `domain/` | `session.go`, `client.go`, `role.go` |
| `application/` | `session_store.go`, `pty_provider.go`, `session_manager.go`, `connection_manager.go` |
| `infrastructure/` | `memory_store.go`, `pty_manager.go`, `pty_instance.go` |
| `transport/` | `websocket.go`, `session_handler.go`, `pty_streamer.go` |
| Tests | `session_manager_test.go` |

---

## Session 1: The Vision & Architecture Contract

**Goal:** Understand WHY this project exists, what problem it solves, and the rules the architecture enforces — before reading a single line of code.

**Files to read:**
- `README.md`
- `docs/vision.md`
- `docs/requirements.md`
- `docs/architecture.md`
- `docs/roadmap.md`
- `docs/tasks.md`

**Concepts to learn:**
- The problem: remote terminal collaboration that is lightweight, self-hostable, and extensible
- The success criteria: multiple users, low latency, stable sessions, AI-agent-readable codebase
- The 5 architectural principles (especially: "depend on interfaces, not implementations" and "business logic must not depend on transport or storage")
- The 4-layer structure: Transport → Application → Domain → Infrastructure
- The MVP scope boundaries — what is intentionally EXCLUDED (auth, Docker, persistence, recording)
- The task dependency graph: TASK-001 → ... → TASK-014

**Key question to hold:** Why was "AI-agent-friendly codebase" listed as a non-functional requirement? What design choices serve that goal?

**Expected outcome:** You can draw the architecture diagram from memory and explain every layer's responsibility. You know what the MVP is and isn't.

**Estimated duration:** 30 minutes

---

## Session 2: Domain Modeling — The Entities

**Goal:** Understand the core domain objects and why they were designed with minimal fields.

**Files to read:**
- `internal/domain/session.go`
- `internal/domain/client.go`
- `internal/domain/role.go`
- Re-read `docs/architecture.md` § "Domain Objects"

**Concepts to learn:**
- `Session` has 3 fields: ID, Clients, Host — why no PTY field here?
- `Client` has 2 fields: ID, Role — why no WebSocket connection stored here?
- `Role` as an `int` enum via `iota` — why not a string? Why not a separate file per role?
- The decision to keep domain models **pure** — no infrastructure imports, no business logic methods
- How these entities map to the requirements (FR-011, FR-012, FR-013, FR-014)

**Key question to hold:** What happens when you need to add a third role (e.g., Moderator)? Does this design accommodate it?

**Expected outcome:** You can explain why each domain model is intentionally minimal and what responsibilities live elsewhere.

**Estimated duration:** 25 minutes

---

## Session 3: The Dependency Inversion Pattern

**Goal:** Understand how interfaces create boundaries between layers, and why this matters for testing and future changes.

**Files to read:**
- `internal/application/session_store.go`
- `internal/application/pty_provider.go`
- `internal/infrastructure/memory_store.go`
- `internal/infrastructure/pty_manager.go`
- `internal/infrastructure/pty_instance.go`

**Concepts to learn:**
- **Interface location rule:** interfaces live in the `application` package (the consumer), not in `infrastructure` (the provider) — this is Go's "accept interfaces, return structs" idiom
- `SessionStore` interface: Create/Get/Delete — why no Update? Why no List?
- `PTYProvider` interface: Spawn/Stop — why is `PTYInstance` also an interface?
- `var _ application.SessionStore = (*MemoryStore)(nil)` — compile-time interface satisfaction check
- How `MemoryStore` uses a simple `map[string]*domain.Session` with no mutex
- How `PTYManager.Spawn()` spawns `bash` with `TERM=xterm` and starts a goroutine

**Key question to hold:** The `MemoryStore` has no concurrency protection (no mutex). Is this a bug, or is the assumption that `SessionManager` handles all synchronization?

**Expected outcome:** You can explain the dependency flow: `transport → application ← infrastructure` and why interfaces sit in the application layer.

**Estimated duration:** 40 minutes

---

## Session 4: SessionManager — The Orchestration Layer

**Goal:** Understand the core business logic that coordinates session lifecycle, PTY management, and storage.

**Files to read:**
- `internal/application/session_manager.go`
- `internal/application/session_manager_test.go`

**Concepts to learn:**
- `SessionManager` holds 3 things: `store` (SessionStore), `ptyManager` (PTYProvider), `ptys` (map of active PTYs)
- `CreateSession` flow: generate UUID → create domain object → store it → spawn PTY → register PTY → **rollback on failure** (delete session if PTY spawn fails)
- `DeleteSession` flow: stop PTY → remove from PTY map → delete from store — note the cleanup ordering
- `sync.RWMutex` protecting the `ptys` map — read lock for reads, write lock for writes
- `GetAllSessionIDs()` — used by PTYStreamer to discover sessions
- The **mock-based testing** approach: `mockStore`, `mockPTY`, `mockPTYProvider`
- Why tests verify behavior (not implementation) — e.g., checking `session.Host.ID == "host-1"` not checking internal map state

**Key question to hold:** `SessionManager` manages both sessions AND PTYs. Is this a single-responsibility violation, or a pragmatic choice for an MVP?

**Expected outcome:** You can trace the complete lifecycle of a session from creation to deletion, including failure paths.

**Estimated duration:** 45 minutes

---

## Session 5: The HTTP API — Exposing Business Logic

**Goal:** Understand how the session creation HTTP endpoint bridges transport to application.

**Files to read:**
- `internal/transport/session_handler.go`
- Re-read `internal/application/session_manager.go` (CreateSession method)

**Concepts to learn:**
- `SessionHandler` depends only on `*application.SessionManager` — it knows nothing about storage or PTY
- HTTP method check (POST only) — basic request validation
- `sm.CreateSession("host")` — note the hardcoded host ID `"host"` is a placeholder
- JSON response with `http.StatusCreated` (201)
- Error handling: 500 for internal errors, 405 for wrong method
- The separation: HTTP concerns (method, content-type, status code) stay in transport, business logic stays in application

**Key question to hold:** The host ID is hardcoded as `"host"`. When a user later connects via WebSocket, they might also be assigned as host. How does this interact?

**Expected outcome:** You can explain the full request lifecycle: HTTP POST → handler → SessionManager → MemoryStore → PTYManager → response.

**Estimated duration:** 25 minutes

---

## Session 6: WebSocket Gateway — The Real-Time Entry Point

**Goal:** Understand how WebSocket connections are established, how roles are assigned, and how the read loop works.

**Files to read:**
- `internal/transport/websocket.go`
- Re-read `internal/application/connection_manager.go`

**Concepts to learn:**
- `upgrader.CheckOrigin` returns `true` always — this is a development-only decision
- The WebSocket endpoint: `GET /ws?sessionID=<id>`
- Connection flow: validate session exists → upgrade HTTP → generate client ID → register connection → assign role → enter read loop
- **Role assignment logic:** if `session.Host == nil`, the connecting user becomes host; otherwise viewer — this means the FIRST WebSocket connection to a session becomes the host
- `readLoop`: reads messages from client → checks role → if host, forwards to PTY → if viewer, rejects with log
- Cleanup: `defer` handles unregister + close on disconnect
- The `ConnectionManager` stores `*websocket.Conn` keyed by client ID — note: it broadcasts to ALL connections globally, not per-session

**Key question to hold:** The `ConnectionManager.Broadcast()` sends to ALL connected clients, not just clients in the same session. Is this a bug or an acceptable MVP shortcut?

**Expected outcome:** You can trace the full lifecycle of a WebSocket connection from upgrade to disconnect.

**Estimated duration:** 40 minutes

---

## Session 7: The Composition Root — How Everything Connects

**Goal:** Understand how `main.go` wires all dependencies together and starts the server.

**Files to read:**
- `cmd/server/main.go`
- Re-read all interface definitions for reference

**Concepts to learn:**
- **Composition root pattern:** `main.go` is the ONLY place that knows about concrete implementations
- Wiring order: Store → PTYManager → SessionManager(store, ptyManager) → ConnectionManager() → PTYStreamer(cm) → Server(sm, cm)
- `go ptyStreamer.WatchSessions(sm)` — starts background goroutine for PTY streaming
- Two HTTP routes: `/ws` (WebSocket) and `/session` (HTTP API)
- `http.ListenAndServe(":8080", nil)` — uses default mux
- Why `main.go` imports all three layers (application, infrastructure, transport)

**Key question to hold:** Could you swap `MemoryStore` for a database-backed store without changing any file other than `main.go`? (Answer: yes — that's the point of the interface pattern)

**Expected outcome:** You can draw the dependency graph of all components and explain why `main.go` is the only coupling point.

**Estimated duration:** 25 minutes

---

## Session 8: PTY Abstraction — Interfaces for Terminal Operations

**Goal:** Understand how the PTY subsystem is abstracted and why it uses two separate interfaces.

**Files to read:**
- `internal/application/pty_provider.go`
- `internal/infrastructure/pty_instance.go`
- `internal/infrastructure/pty_manager.go`

**Concepts to learn:**
- `PTYProvider` interface (Spawn/Stop) — factory pattern for creating terminal instances
- `PTYInstance` interface (Write/Output/Close) — represents a running terminal
- Why two interfaces: the provider creates instances, instances are used independently
- `PTYInstance` concrete type: holds `exec.Cmd`, `*os.File` (the PTY), output channel, done channel
- `Write()` forwards data to the PTY file descriptor
- `Output()` returns a read-only channel — consumers receive raw terminal bytes
- `Close()` kills the process and closes the PTY file
- `readLoop()`: reads from PTY in 4096-byte chunks, copies data, sends to output channel with **drop-on-full** behavior (non-blocking send with default case)

**Key question to hold:** The output channel has a buffer of 256. What happens when output arrives faster than clients consume it? (Answer: data is dropped silently)

**Expected outcome:** You can explain how a real bash process is spawned, how output flows through channels, and what happens when the buffer fills.

**Estimated duration:** 40 minutes

---

## Session 9: The Streaming Pipeline — PTY to WebSocket

**Goal:** Understand how terminal output from a PTY reaches all connected WebSocket clients.

**Files to read:**
- `internal/transport/pty_streamer.go`
- Re-read `internal/transport/websocket.go` (for the broadcast target)
- Re-read `internal/infrastructure/pty_manager.go` (readLoop)

**Concepts to learn:**
- `PTYStreamer` discovers sessions by polling `GetAllSessionIDs()` — this is a polling-based design, not event-driven
- `WatchSessions()` creates an unbuffered channel `ticker` that is never sent to — the goroutine blocks on `range ticker` and never executes. **This is a bug/dead code.** Only the initial `checkSessions` call runs, then `select{}` blocks forever.
- `tryStream()`: deduplicates via `streams` map (one stream per session), starts a goroutine per session
- `stream()`: reads from PTY output channel → broadcasts to ALL connections
- **Critical design observation:** `cm.Broadcast()` sends to ALL connections globally, not scoped to a session. If session A's PTY produces output, viewers of session B also receive it. This is a significant architectural gap.

**Key question to hold:** The ticker channel is never written to. The `WatchSessions` polling loop is dead code. How should this be redesigned?

**Expected outcome:** You can trace the complete data path: PTY read → channel → stream goroutine → Broadcast → all WebSocket connections.

**Estimated duration:** 45 minutes

---

## Session 10: Input Handling — The Host Controls Everything

**Goal:** Understand the permission model and how host input reaches the terminal.

**Files to read:**
- `internal/transport/websocket.go` (readLoop, lines 61–93)
- Re-read `internal/domain/role.go`

**Concepts to learn:**
- Input is handled inside `readLoop` — same goroutine that detects disconnects
- Role check: only `RoleHost` can write to PTY; viewers' messages are silently dropped (logged)
- `pty.Write(message)` — raw bytes are forwarded directly to the PTY
- No message framing or protocol — the WebSocket message body IS the terminal input
- Disconnect detection: `websocket.IsUnexpectedCloseError` filters expected close codes
- The host identity problem: `CreateSession("host")` creates a host with ID `"host"`, but when the first WebSocket client connects, `session.Host` is set to that client's UUID — the original `"host"` client never connects

**Key question to hold:** Is the `CreateSession("host")` host assignment meaningful? When a real client connects and becomes host, the original host ID is orphaned.

**Expected outcome:** You can explain the full input path: browser keystroke → WebSocket message → readLoop → role check → PTY.Write → bash processes it.

**Estimated duration:** 30 minutes

---

## Session 11: Testing Strategy & Mock Patterns

**Goal:** Understand how the project tests business logic without infrastructure dependencies.

**Files to read:**
- `internal/application/session_manager_test.go`
- Re-read `internal/application/session_store.go` and `pty_provider.go` (the interfaces being mocked)

**Concepts to learn:**
- **Test doubles:** `mockStore` (in-memory map, no concurrency), `mockPTY` (channel-based, no real process), `mockPTYProvider` (returns mockPTY)
- The mocks implement the exact same interfaces as real implementations — this is possible because interfaces are in `application`, not `infrastructure`
- Test coverage: Create, Get, GetNotFound, Delete, DeleteNotFound — 5 test cases
- Tests verify **external behavior** (session has ID, host matches, errors on missing) — not internal state
- No table-driven tests used — each test is a standalone function
- No tests for PTY streaming, WebSocket handling, or connection management

**Key question to hold:** What would it take to test the WebSocket layer? What's blocking it?

**Expected outcome:** You can explain the mock pattern, what's tested, what's NOT tested, and how the interface design enables testing.

**Estimated duration:** 35 minutes

---

## Session 12: Architectural Gaps, Trade-offs & Concurrency

**Goal:** Identify what's missing, what's trade-offed, and what could break under load.

**Files to read:**
- ALL `.go` files (comparative pass)
- Re-read `docs/architecture.md` § "Architectural Principles"

**Concepts to learn:**

| Gap | Details |
|---|---|
| **No session-scoped broadcast** | `ConnectionManager.Broadcast()` sends to ALL connections globally — cross-session data leakage |
| **No concurrency in MemoryStore** | `map[string]*Session` with no mutex — relies on `SessionManager`'s mutex for safety, but `session_handler.go` calls `CreateSession` without explicit synchronization beyond what SessionManager provides |
| **Dead polling code** | `PTYStreamer.WatchSessions` ticker channel is never written to — the `for range ticker` goroutine never executes |
| **Hardcoded host ID** | `CreateSession("host")` creates a phantom host; real host is assigned on first WebSocket connect |
| **No graceful shutdown** | `ListenAndServe` runs forever; no signal handling; PTY processes not cleaned up on exit |
| **No session removal from PTY map** | `DeleteSession` removes from store but the PTY map cleanup depends on the caller |
| **PTY output drops silently** | Buffer of 256 with non-blocking send — output can be lost without any client notification |
| **No message framing** | WebSocket messages are raw bytes — no JSON envelope for typed messages (join, input, output, error) |
| **No viewer list** | `Session.Clients` slice exists but is never populated — no tracking of who's in a session |
| **CheckOrigin always true** | Security concern for production |

**Expected outcome:** You have a complete mental list of every architectural trade-off and can prioritize them by severity.

**Estimated duration:** 50 minutes

---

## Session 13: Redesign Exercise — Build It From Scratch

**Goal:** Mentally redesign the entire system, applying what you've learned, to solidify understanding.

**Exercise — answer these questions:**

1. **If you were to rebuild this, what would you change in the domain layer?** (e.g., should Session track its own PTY? Should Client have a connection reference?)

2. **How would you fix the broadcast scoping?** (per-session connection registry vs. global registry with session filtering)

3. **How would you replace the polling-based PTYStreamer with an event-driven design?** (session creation callback? channel-based notification?)

4. **What message protocol would you design for WebSocket communication?** (typed messages: `{type: "input", data: "..."}`, `{type: "output", data: "..."}`, etc.)

5. **How would you add authentication without breaking the architecture?** (middleware in transport layer? session-level token?)

6. **How would you add session persistence?** (swap MemoryStore for Postgres-backed store — what interfaces change?)

7. **How would you add Docker isolation?** (what new interface would PTYProvider need? How does container lifecycle map to session lifecycle?)

**Expected outcome:** You can explain every design decision in the project, justify each trade-off, and articulate how you would evolve the architecture for production.

**Estimated duration:** 45 minutes

---

## Summary

| # | Session | Duration | Task Coverage |
|---|---|---|---|
| 1 | Vision & Architecture | 30 min | All (docs only) |
| 2 | Domain Modeling | 25 min | TASK-002 |
| 3 | Dependency Inversion | 40 min | TASK-003, TASK-004, TASK-011, TASK-012 |
| 4 | SessionManager | 45 min | TASK-005, TASK-006 |
| 5 | HTTP API | 25 min | TASK-012.5 |
| 6 | WebSocket Gateway | 40 min | TASK-007, TASK-008, TASK-009 |
| 7 | Composition Root | 25 min | TASK-001 (wiring) |
| 8 | PTY Abstraction | 40 min | TASK-011, TASK-012 |
| 9 | Streaming Pipeline | 45 min | TASK-010, TASK-013 |
| 10 | Input Handling | 30 min | TASK-014 |
| 11 | Testing Patterns | 35 min | TASK-006 |
| 12 | Architecture Gaps | 50 min | Cross-cutting |
| 13 | Redesign Exercise | 45 min | All (synthesis) |

**Total: ~8.5 hours across 13 sessions**

After completing all sessions, you will have:
- Traced every data path (create session → join → type command → see output)
- Understood every interface and why it exists
- Identified every architectural trade-off
- Be able to explain or redesign any component independently
