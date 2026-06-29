# Engineering Review Guide

## Purpose

This workbook is designed to guide an engineer through a complete architectural understanding of the Realtime Terminal Collaboration project. It is not a code review for finding bugs — it is a structured path to understanding every design decision, every trade-off, and every architectural boundary in the codebase.

By the end of this review, you should be able to explain why every file exists, justify every interface, and redesign the system from scratch without opening the code.

## How to Use This Workbook

1. **Follow the sessions in order.** Each session builds on the previous one.
2. **Answer every question before moving on.** Write your answers down. The workbook rewards active thinking.
3. **Complete every exercise.** Drawing diagrams and tracing flows forces deeper understanding than passive reading.
4. **Use the checklist at the end of each session.** If any box is unchecked, you missed something — go back.
5. **Fill in Notes and Questions for Later.** These become your personal review notes.

## Expected Outcome

After completing all 13 sessions:

- You can draw the full architecture from memory
- You can explain every interface and why it lives where it does
- You can trace every data flow end-to-end
- You can identify every architectural gap and propose fixes
- You can redesign the system for production with informed trade-offs

## Recommended Pace

| Pace | Sessions per day | Total time |
|------|-----------------|------------|
| Intensive | 4–5 | 2–3 days |
| Standard | 2–3 | 4–5 days |
| Thorough | 1 | 2 weeks |

## Estimated Total Time

~8.5 hours across 13 sessions.

---

# Table of Contents

| # | Session | Focus | Duration |
|---|---------|-------|----------|
| 1 | [The Vision & Architecture Contract](#session-1--the-vision--architecture-contract) | Why this project exists | 30 min |
| 2 | [Domain Modeling — The Entities](#session-2--domain-modeling--the-entities) | Core domain objects | 25 min |
| 3 | [The Dependency Inversion Pattern](#session-3--the-dependency-inversion-pattern) | Interfaces and abstractions | 40 min |
| 4 | [SessionManager — The Orchestration Layer](#session-4--sessionmanager--the-orchestration-layer) | Core business logic | 45 min |
| 5 | [The HTTP API — Exposing Business Logic](#session-5--the-http-api--exposing-business-logic) | REST endpoint design | 25 min |
| 6 | [WebSocket Gateway — The Real-Time Entry Point](#session-6--websocket-gateway--the-real-time-entry-point) | Connection lifecycle | 40 min |
| 7 | [The Composition Root — How Everything Connects](#session-7--the-composition-root--how-everything-connects) | Dependency wiring | 25 min |
| 8 | [PTY Abstraction — Interfaces for Terminal Operations](#session-8--pty-abstraction--interfaces-for-terminal-operations) | Terminal interface design | 40 min |
| 9 | [The Streaming Pipeline — PTY to WebSocket](#session-9--the-streaming-pipeline--pty-to-websocket) | Output delivery path | 45 min |
| 10 | [Input Handling — The Host Controls Everything](#session-10--input-handling--the-host-controls-everything) | Permission model and input flow | 30 min |
| 11 | [Testing Strategy & Mock Patterns](#session-11--testing-strategy--mock-patterns) | How the project tests itself | 35 min |
| 12 | [Architectural Gaps, Trade-offs & Concurrency](#session-12--architectural-gaps-trade-offs--concurrency) | What's missing and why | 50 min |
| 13 | [Redesign Exercise — Build It From Scratch](#session-13--redesign-exercise--build-it-from-scratch) | Full synthesis | 45 min |

---

# Session 1 — The Vision & Architecture Contract

## Goal

Understand why this project exists, what problem it solves, and the rules the architecture enforces — before reading a single line of code.

---

## Files to Read

- `README.md`
- `docs/vision.md`
- `docs/requirements.md`
- `docs/architecture.md`
- `docs/roadmap.md`
- `docs/tasks.md`

---

## Before You Read

1. What makes a "good" architecture for an MVP? What are you optimizing for — speed, extensibility, or clarity?
2. Why would someone list "AI-agent-friendly codebase" as a non-functional requirement? What design choices serve that goal?
3. If you were building this from scratch with no docs, what decisions would you need to make first?
4. The architecture defines 5 principles. Which one is hardest to enforce in practice?

---

## Concepts to Learn

**The Problem:** Remote collaboration on terminal-based workflows is difficult. Existing solutions are too heavy, unfriendly to self-hosting, hard to extend, and hard to learn from.

**The Solution:** A lightweight, web-based collaborative terminal. Multiple users join the same session and observe terminal output in real time.

**MVP Scope (intentionally excluded):**
- Authentication
- Persistent sessions
- Docker isolation
- Terminal recording
- Multi-user input control
- Production deployment

**The 4-Layer Architecture:**

```
Transport Layer    →  WebSocket, HTTP
Application Layer  →  Session Manager, Connection Manager
Domain Layer       →  Session, Client, Role
Infrastructure Layer →  MemoryStore, PTY Manager
```

**The 5 Architectural Principles:**

1. Depend on interfaces, not implementations
2. Business logic must not depend on transport layer
3. Business logic must not depend on storage implementation
4. Every module should have a single responsibility
5. Realtime communication is performed through WebSocket only

**The Task Dependency Graph:**

```
TASK-001 (Project setup)
  └─→ TASK-002 (Domain models)
       └─→ TASK-003 (SessionStore interface)
            └─→ TASK-004 (MemoryStore)
                 └─→ TASK-005 (SessionManager)
                      └─→ TASK-006 (Tests)
                      └─→ TASK-007 (WebSocket endpoint)
                           └─→ TASK-008 (ConnectionManager)
                                └─→ TASK-009 (Join flow)
                                     └─→ TASK-010 (Broadcast)
                                          └─→ TASK-013 (Stream output)
                                               └─→ TASK-014 (Input handling)
                      └─→ TASK-011 (PTY library)
                           └─→ TASK-012 (PTY Manager)
                                └─→ TASK-012.5 (HTTP API)
                                └─→ TASK-013 (Stream output)
```

---

## Trace Exercise

Draw this diagram from memory after reading:

```
Browser → HTTP POST /session → SessionHandler → SessionManager → MemoryStore
                                                              → PTYManager
Browser → WebSocket /ws → Server → SessionManager → ConnectionManager
                                                   → PTYInstance
PTYInstance.readLoop → output channel → PTYStreamer → ConnectionManager.Broadcast → WebSocket clients
```

---

## Architecture Questions

1. Why is the Domain layer at the BOTTOM of the dependency diagram, not the top?
2. Why does `architecture.md` explicitly state that business logic must not depend on transport or storage? What breaks if it does?
3. The roadmap has 5 phases. Why is "Hardening" (logging, metrics, tests) the LAST phase instead of being integrated from the start?
4. What would change if you replaced Go with a different language? Which architectural decisions are language-independent?

---

## Hands-on Exercise

Draw the folder structure from memory. For each package, write one sentence describing its responsibility. Then compare against `docs/architecture.md`.

---

## Review Checklist

- [ ] I can explain the problem this project solves
- [ ] I can list the MVP features and the non-features
- [ ] I can draw the 4-layer architecture from memory
- [ ] I can state all 5 architectural principles
- [ ] I understand the task dependency graph
- [ ] I understand why "AI-agent-friendly" is an NFR

---

## Notes

---

## Questions for Later

---

## Estimated Duration

30 minutes

---

# Session 2 — Domain Modeling — The Entities

## Goal

Understand the core domain objects and why they were designed with minimal fields. The domain layer is the foundation — every design decision above it depends on getting this right.

---

## Files to Read

- `internal/domain/session.go`
- `internal/domain/client.go`
- `internal/domain/role.go`
- Re-read `docs/architecture.md` § "Domain Objects"

---

## Before You Read

1. What fields would you put on a `Session` struct? Think about what a session needs to track.
2. Why might you intentionally KEEP fields OFF a domain model?
3. What is the difference between a "domain model" and a "data transfer object"?
4. Should a domain model contain behavior (methods) or just data?

---

## Concepts to Learn

**Session struct — 3 fields only:**

```go
type Session struct {
    ID      string
    Clients []*Client
    Host    *Client
}
```

Why no PTY reference? Why no WebSocket connection? Why no state field (active/terminated)?

**Client struct — 2 fields only:**

```go
type Client struct {
    ID   string
    Role Role
}
```

Why no `*websocket.Conn`? Why no channel? Why no address or metadata?

**Role — int enum via iota:**

```go
type Role int
const (
    RoleHost   Role = iota
    RoleViewer
)
```

Why `int` and not `string`? Why not a separate type per role? Why put it in its own file?

**Domain purity:** These three files have zero imports from infrastructure, transport, or application packages. This is enforced by architectural principle #2.

---

## Trace Exercise

Map each requirement to its domain representation:

| Requirement | Domain Field | Where Enforced |
|-------------|-------------|----------------|
| FR-011: Host role exists | `RoleHost` | `role.go` |
| FR-012: Viewer role exists | `RoleViewer` | `role.go` |
| FR-013: Session created by host | `Session.Host` | `session.go` |
| FR-014: Each session has exactly one host | `Session.Host *Client` (pointer, not slice) | `session.go` |
| FR-004: Multiple users can join | `Session.Clients []*Client` (slice) | `session.go` |
| FR-009: Client communicates via WebSocket | NOT in domain | Transport layer |

Fill in the last row yourself: Why is the WebSocket connection deliberately excluded from the domain?

---

## Architecture Questions

1. `Session.Host` is a `*Client` (single pointer). `Session.Clients` is a `[]*Client` (slice). What does this structural choice communicate about the business rules?
2. If you needed to add a third role (Moderator), what changes? Is the `iota` enum extensible enough?
3. The domain models have no methods. Is this a limitation or a deliberate choice? What would a domain method look like for `Session`?
4. `Session.Clients` is a slice of pointers. Could this cause issues with concurrent access? Who is responsible for thread safety?

---

## Hands-on Exercise

Write a one-paragraph explanation of why `Session` does not have a PTY field. Write a second paragraph explaining why `Client` does not have a WebSocket connection field. These two explanations capture the essence of the domain design.

---

## Review Checklist

- [ ] I can list every field on Session, Client, and Role
- [ ] I understand why PTY is NOT in the domain model
- [ ] I understand why WebSocket connection is NOT in the domain model
- [ ] I can explain the Host vs Clients structural difference
- [ ] I know what layer enforces each business rule

---

## Notes

---

## Questions for Later

---

## Estimated Duration

25 minutes

---

# Session 3 — The Dependency Inversion Pattern

## Goal

Understand how interfaces create boundaries between layers, why interfaces live in the application layer (not infrastructure), and how this pattern enables testing and future changes.

---

## Files to Read

- `internal/application/session_store.go`
- `internal/application/pty_provider.go`
- `internal/infrastructure/memory_store.go`
- `internal/infrastructure/pty_manager.go`
- `internal/infrastructure/pty_instance.go`

---

## Before You Read

1. In Go, where should an interface be defined — by the provider or the consumer?
2. What problem do interfaces solve that concrete types don't?
3. Why would you want `application` to never import `infrastructure`?
4. What is a compile-time interface check and why does it exist?

---

## Concepts to Learn

**Interface location rule (Go idiom):**

Interfaces live in the `application` package — the **consumer** of the behavior. They do NOT live in `infrastructure` — the **provider**. This means:

```
application/session_store.go    →  defines SessionStore interface
infrastructure/memory_store.go  →  implements SessionStore interface
application/pty_provider.go     →  defines PTYProvider + PTYInstance interfaces
infrastructure/pty_manager.go   →  implements PTYProvider
infrastructure/pty_instance.go  →  implements PTYInstance
```

**Compile-time interface assertion:**

```go
var _ application.SessionStore = (*MemoryStore)(nil)
var _ application.PTYProvider  = (*PTYManager)(nil)
```

These lines cause a compile error if the struct ever stops satisfying the interface. They are documentation + safety net.

**SessionStore — 3 operations, no more:**

```go
type SessionStore interface {
    Create(session *domain.Session) error
    Get(id string) (*domain.Session, error)
    Delete(id string) error
}
```

Why no `Update`? Why no `List`? What does this tell you about the MVP scope?

**PTYProvider — factory pattern:**

```go
type PTYProvider interface {
    Spawn() (PTYInstance, error)
    Stop(inst PTYInstance) error
}

type PTYInstance interface {
    Write(data []byte) error
    Output() <-chan []byte
    Close() error
}
```

Two interfaces, not one. Why separate the factory (Spawn/Stop) from the instance (Write/Output/Close)?

**MemoryStore — zero concurrency protection:**

`MemoryStore` uses a plain `map[string]*domain.Session` with no mutex. The assumption is that `SessionManager` handles all synchronization. Is this safe?

---

## Trace Exercise

Trace the dependency arrows:

```
application/session_store.go  ← defines interface
    ↑ imports
infrastructure/memory_store.go ← implements interface
    ↑ imports
application (SessionManager) ← uses interface

Result: infrastructure imports application, application NEVER imports infrastructure.
```

Draw the same diagram for PTYProvider/PTYInstance.

---

## Architecture Questions

1. What would happen if `SessionStore` interface was defined in `infrastructure` instead of `application`? Which architectural principle would break?
2. `MemoryStore` has no mutex. Is this a bug or a valid design choice? Under what conditions would it become a bug?
3. `PTYInstance` is also an interface, not just `PTYProvider`. Why? What does this enable?
4. If you added a PostgreSQL-backed store, what files would change? What would NOT change?

---

## Hands-on Exercise

Mentally replace `MemoryStore` with a PostgreSQL-backed store:

1. What new file would you create?
2. What interface would it implement?
3. What changes in `main.go`?
4. What changes in `SessionManager`?

Answer all four before reading any code.

---

## Review Checklist

- [ ] I can explain why interfaces live in the application layer
- [ ] I understand the compile-time interface assertion pattern
- [ ] I can list every operation on SessionStore and PTYProvider
- [ ] I understand the dependency direction: infrastructure → application
- [ ] I know what would change if I swapped MemoryStore for a database

---

## Notes

---

## Questions for Later

---

## Estimated Duration

40 minutes

---

# Session 4 — SessionManager — The Orchestration Layer

## Goal

Understand the core business logic that coordinates session lifecycle, PTY management, and storage. This is the heart of the application — everything else delegates to it.

---

## Files to Read

- `internal/application/session_manager.go`
- `internal/application/session_manager_test.go`

---

## Before You Read

1. What operations does a session manager need? Think about the full lifecycle: create → use → destroy.
2. How do you handle failure during creation? What if step 2 of 3 fails?
3. How do you protect a shared map from concurrent access in Go?
4. What makes a good unit test for a service that depends on external systems?

---

## Concepts to Learn

**SessionManager holds 3 dependencies:**

```go
type SessionManager struct {
    store      SessionStore
    ptyManager PTYProvider
    ptys       map[string]PTYInstance
    mu         sync.RWMutex
}
```

Two interfaces (injected) + one internal map + one mutex. Why is `ptys` an internal map rather than stored in `SessionStore`?

**CreateSession — transaction-like behavior with manual rollback:**

```go
func (s *SessionManager) CreateSession(hostID string) (*domain.Session, error) {
    // 1. Create domain object
    // 2. Store it
    // 3. Spawn PTY
    // 4. Register PTY
    // If step 3 fails → delete session (rollback step 2)
}
```

This is not a database transaction — it's manual compensation. Is this reliable? What happens if step 4 fails after step 3 succeeds?

**DeleteSession — cleanup ordering matters:**

```go
func (s *SessionManager) DeleteSession(id string) error {
    // 1. Stop PTY (release OS resources)
    // 2. Remove from PTY map
    // 3. Delete from store
}
```

Why stop the PTY before deleting from the store? What happens if you reverse the order?

**sync.RWMutex — read/write lock pattern:**

- `RLock/RUnlock` for `GetPTY`, `GetAllSessionIDs` — allows concurrent reads
- `Lock/Unlock` for `CreateSession`, `DeleteSession` — exclusive writes

**Mock-based testing:**

The test file defines `mockStore`, `mockPTY`, `mockPTYProvider` — all implementing the same interfaces as real implementations. Tests verify behavior, not implementation details.

---

## Trace Exercise

Trace the CreateSession flow step by step:

```
CreateSession("host-1")
    │
    ├─→ domain.Session{ID: uuid, Host: &Client{ID: "host-1", Role: Host}}
    │
    ├─→ store.Create(session)        // persists to memory
    │       │
    │       └─→ map[uuid] = session
    │
    ├─→ ptyManager.Spawn()           // spawns bash process
    │       │
    │       ├─→ exec.Command("bash")
    │       ├─→ pty.Start(cmd)
    │       └─→ go instance.readLoop()
    │
    ├─→ s.ptys[session.ID] = pty    // register PTY
    │
    └─→ return session, nil
```

Now trace the failure path: what happens if `ptyManager.Spawn()` returns an error?

---

## Architecture Questions

1. `SessionManager` manages both sessions AND PTYs. Is this a single-responsibility violation? Or is "session lifecycle" broad enough to include PTY?
2. The rollback in `CreateSession` (delete session if PTY spawn fails) is not atomic. What happens if `store.Delete` also fails?
3. `GetAllSessionIDs()` returns a snapshot. What happens if a session is deleted between the snapshot and the caller using it?
4. Why is `ptys` stored in `SessionManager` rather than in `SessionStore`? What would change if it were in the store?

---

## Hands-on Exercise

Write the test for `CreateSession` without looking at the code. You should produce:

1. Create a mock store and mock PTY provider
2. Call `CreateSession`
3. Assert session has a non-empty ID
4. Assert session has a host with the correct ID
5. Assert no error

Then compare with the actual test.

---

## Review Checklist

- [ ] I can trace CreateSession from start to finish
- [ ] I can trace DeleteSession from start to finish
- [ ] I understand the manual rollback pattern
- [ ] I know why sync.RWMutex is used
- [ ] I can explain the mock-based testing approach
- [ ] I understand why ptys is an internal map

---

## Notes

---

## Questions for Later

---

## Estimated Duration

45 minutes

---

# Session 5 — The HTTP API — Exposing Business Logic

## Goal

Understand how the session creation HTTP endpoint bridges transport to application, and why HTTP concerns stay in the transport layer.

---

## Files to Read

- `internal/transport/session_handler.go`
- Re-read `internal/application/session_manager.go` (CreateSession method only)

---

## Before You Read

1. What HTTP methods does session creation need? Why not GET?
2. What status codes should a creation endpoint return?
3. What does the handler need to know about the system? What should it NOT know?

---

## Concepts to Learn

**Handler depends only on SessionManager:**

```go
type SessionHandler struct {
    sm *application.SessionManager
}
```

It knows nothing about MemoryStore, PTYManager, or WebSocket connections. This is the transport layer calling into the application layer — exactly as architected.

**Request validation stays in transport:**

- HTTP method check (POST only) → 405 Method Not Allowed
- No body parsing needed for session creation
- Future: authentication would be checked here

**Response construction stays in transport:**

- Set `Content-Type: application/json`
- Set `201 Created`
- Encode session as JSON

**Error propagation:**

- `SessionManager.CreateSession` returns `error`
- Handler maps it to `500 Internal Server Error`
- No error differentiation (auth failure vs. internal failure — same 500)

**The hardcoded host ID `"host"`:**

```go
session, err := h.sm.CreateSession("host")
```

This is a placeholder. When a real client connects via WebSocket, the host assignment is overwritten. This creates an orphaned host ID.

---

## Trace Exercise

```
POST /session
    │
    ├─→ SessionHandler.CreateSession()
    │       │
    │       ├─→ method != POST? → 405
    │       │
    │       ├─→ sm.CreateSession("host")
    │       │       │
    │       │       ├─→ domain.Session{ID: uuid, Host: {ID: "host", Role: Host}}
    │       │       ├─→ store.Create(session)
    │       │       ├─→ ptyManager.Spawn()
    │       │       └─→ return session, nil
    │       │
    │       ├─→ err != nil? → 500
    │       │
    │       ├─→ 201 Created
    │       └─→ JSON: {"ID":"uuid","Host":{...},"Clients":null}
```

---

## Architecture Questions

1. Why does the handler call `sm.CreateSession("host")` with a hardcoded ID? What would be the correct host ID at this point?
2. The handler returns `500` for all errors. Should it distinguish between "session limit reached" (429) vs "internal failure" (500)?
3. If you added input validation (e.g., session name), where would it go — handler or SessionManager?
4. What security concerns exist with this endpoint in production?

---

## Hands-on Exercise

Write the curl command to create a session and the expected JSON response. Then mentally add a GET /session/:id endpoint — what handler method would you write, and what SessionManager method would it call?

---

## Review Checklist

- [ ] I can explain the handler's responsibility boundary
- [ ] I understand why method check is in the handler
- [ ] I can trace the request-response cycle
- [ ] I understand the hardcoded host ID issue
- [ ] I know where input validation belongs

---

## Notes

---

## Questions for Later

---

## Estimated Duration

25 minutes

---

# Session 6 — WebSocket Gateway — The Real-Time Entry Point

## Goal

Understand how WebSocket connections are established, how roles are assigned, and how the read loop handles input. This is the most complex transport-layer component.

---

## Files to Read

- `internal/transport/websocket.go`
- Re-read `internal/application/connection_manager.go`

---

## Before You Read

1. What happens when a browser opens a WebSocket connection? What is the HTTP upgrade process?
2. How do you determine who is "host" when multiple clients connect?
3. What is a "read loop" and why does it need its own goroutine?
4. How do you detect when a client disconnects?

---

## Concepts to Learn

**WebSocket upgrade:**

```go
var upgrader = websocket.Upgrader{
    CheckOrigin: func(r *http.Request) bool { return true },
}
```

`CheckOrigin` returns `true` always — this is a development-only decision. In production, you'd validate the origin.

**Connection flow:**

1. Client connects to `/ws?sessionID=<id>`
2. Handler validates session exists
3. HTTP connection is upgraded to WebSocket
4. Client gets a UUID
5. Connection is registered in ConnectionManager
6. Role is assigned (host or viewer)
7. Handler enters `readLoop` — blocks until disconnect

**Role assignment logic:**

```go
role := domain.RoleViewer
if session.Host == nil {
    role = domain.RoleHost
    session.Host = &domain.Client{ID: clientID, Role: domain.RoleHost}
}
```

The FIRST WebSocket connection to a session becomes the host. All subsequent connections are viewers. This overwrites the `"host"` placeholder from session creation.

**readLoop — single goroutine, dual purpose:**

```
readLoop
    │
    ├─→ Read message from client
    │       │
    │       ├─→ Error? → cleanup and return
    │       │
    │       ├─→ role != Host? → log and continue (reject input)
    │       │
    │       └─→ role == Host? → write to PTY
    │
    └─→ (loop)
```

**Cleanup via defer:**

```go
defer func() {
    s.cm.Unregister(clientID)
    conn.Close()
}()
```

Guarantees cleanup on any exit path — normal disconnect or error.

**ConnectionManager.Broadcast — global, not session-scoped:**

`ConnectionManager` stores ALL connections across ALL sessions in a single map. `Broadcast()` sends to every connection. This means PTY output from session A reaches viewers of session B.

---

## Trace Exercise

```
Browser opens WebSocket to /ws?sessionID=abc-123
    │
    ├─→ Server.HandleWebSocket()
    │       │
    │       ├─→ sessionID = "abc-123"
    │       ├─→ sm.GetSession("abc-123") → session found
    │       ├─→ upgrader.Upgrade() → *websocket.Conn
    │       ├─→ clientID = uuid.New()
    │       ├─→ cm.Register(clientID, conn)
    │       ├─→ session.Host == nil? → YES → role = Host
    │       └─→ s.readLoop(conn, sessionID, clientID, Host)
    │
    readLoop (blocking):
    │
    ├─→ conn.ReadMessage() → message bytes
    │       │
    │       ├─→ role == Host → sm.GetPTY(sessionID) → pty.Write(message)
    │       │
    │       └─→ role == Viewer → log "attempted input, rejected"
    │
    (on disconnect):
    ├─→ cm.Unregister(clientID)
    └─→ conn.Close()
```

---

## Architecture Questions

1. `ConnectionManager` is in the `application` layer but stores `*websocket.Conn` — a transport-layer type. Is this a layer violation? The architecture decision says it was moved from transport to application intentionally. Justify this.
2. `Broadcast()` sends to ALL connections globally. What happens with 3 sessions, each with 5 viewers? Session A's PTY output reaches 14 viewers who don't care about it.
3. The `readLoop` does input forwarding AND disconnect detection in the same loop. Should these be separated?
4. `session.Host` is mutated directly in the WebSocket handler (`session.Host = &domain.Client{...}`). Is this safe with concurrent access?

---

## Hands-on Exercise

Draw the sequence diagram for the WebSocket lifecycle. Include:
- HTTP upgrade
- Role assignment
- Message flow (input → PTY)
- Disconnect cleanup

---

## Review Checklist

- [ ] I can trace the WebSocket upgrade flow
- [ ] I understand the role assignment logic
- [ ] I can explain what readLoop does step by step
- [ ] I understand why Broadcast is global (and why it's a problem)
- [ ] I know what happens on disconnect
- [ ] I understand the CheckOrigin security implication

---

## Notes

---

## Questions for Later

---

## Estimated Duration

40 minutes

---

# Session 7 — The Composition Root — How Everything Connects

## Goal

Understand how `main.go` wires all dependencies together and starts the server. This is the only place in the codebase where concrete implementations are coupled.

---

## Files to Read

- `cmd/server/main.go`
- Re-read all interface definitions in `application/` for reference

---

## Before You Read

1. What is a "composition root"? Why is it important?
2. How many concrete types does `main.go` instantiate?
3. Which packages does `main.go` import? Which does it NOT import?

---

## Concepts to Learn

**Composition root pattern:**

`main.go` is the ONLY file that knows about all concrete implementations. Every other file depends on interfaces. This means you can swap implementations by changing only `main.go`.

**Wiring order:**

```
1. store       = infrastructure.NewMemoryStore()
2. ptyManager  = infrastructure.NewPTYManager()
3. sm          = application.NewSessionManager(store, ptyManager)
4. cm          = application.NewConnectionManager()
5. ptyStreamer  = transport.NewPTYStreamer(cm)
6. go ptyStreamer.WatchSessions(sm)   // background goroutine
7. server      = transport.NewServer(sm, cm)
8. sessionHandler = transport.NewSessionHandler(sm)
```

**Two HTTP routes:**

```
/ws      → server.HandleWebSocket    (WebSocket)
/session → sessionHandler.CreateSession (HTTP API)
```

**The background goroutine:**

`go ptyStreamer.WatchSessions(sm)` starts a goroutine that watches for new sessions and streams their PTY output. This goroutine runs for the lifetime of the server.

**Import graph:**

```
cmd/server imports:
    application    (SessionManager, ConnectionManager)
    infrastructure (MemoryStore, PTYManager)
    transport      (Server, SessionHandler, PTYStreamer)

No other package imports all three layers.
```

---

## Trace Exercise

Draw the dependency graph. For each arrow, write which package imports which:

```
cmd/server ──→ transport ──→ application ──→ domain
    │                              ↑
    └──→ infrastructure ───────────┘
```

Verify: does `transport` import `infrastructure`? Does `application` import `infrastructure`? Does `domain` import anything?

---

## Architecture Questions

1. Could you swap `MemoryStore` for a database-backed store without changing any file other than `main.go`? Verify your answer by checking every import.
2. Why is `ptyStreamer.WatchSessions` started as a goroutine instead of being called synchronously?
3. What happens if `http.ListenAndServe` fails? Are PTY processes cleaned up?
4. If you added a `/health` endpoint, which files would change?

---

## Hands-on Exercise

Replace `MemoryStore` with a hypothetical `PostgresStore` in your head:

1. Create `internal/infrastructure/postgres_store.go`
2. It implements `application.SessionStore`
3. Change `main.go` line 1 from `NewMemoryStore()` to `NewPostgresStore(conn)`
4. Verify: no other file changes

Write down what you just did. This is the power of the composition root pattern.

---

## Review Checklist

- [ ] I can list every dependency wired in main.go
- [ ] I understand why main.go imports all three layers
- [ ] I can explain the background goroutine's purpose
- [ ] I know what would change if I swapped MemoryStore
- [ ] I understand the composition root pattern

---

## Notes

---

## Questions for Later

---

## Estimated Duration

25 minutes

---

# Session 8 — PTY Abstraction — Interfaces for Terminal Operations

## Goal

Understand how the PTY subsystem is abstracted into interfaces, why two separate interfaces exist, and how a real bash process is managed.

---

## Files to Read

- `internal/application/pty_provider.go`
- `internal/infrastructure/pty_instance.go`
- `internal/infrastructure/pty_manager.go`

---

## Before You Read

1. What is a PTY (pseudo-terminal)? Why not just use `os/exec` directly?
2. What operations do you need to manage a terminal process?
3. Why would you separate "creating a terminal" from "using a terminal"?
4. How do you read output from a process that produces data asynchronously?

---

## Concepts to Learn

**Two interfaces, two responsibilities:**

```go
// Factory — creates and destroys terminal instances
type PTYProvider interface {
    Spawn() (PTYInstance, error)
    Stop(inst PTYInstance) error
}

// Instance — a running terminal you can read/write
type PTYInstance interface {
    Write(data []byte) error
    Output() <-chan []byte
    Close() error
}
```

The factory pattern separates creation from usage. `SessionManager` uses `PTYProvider` to create. `PTYStreamer` and `Server` use `PTYInstance` to read/write.

**PTYManager — spawns bash with creack/pty:**

```go
func (pm *PTYManager) Spawn() (application.PTYInstance, error) {
    cmd := exec.Command("bash")
    cmd.Env = append(cmd.Env, "TERM=xterm")
    ptmx, err := pty.Start(cmd)
    // ...
}
```

The `creack/pty` library creates a pseudo-terminal pair and attaches it to the bash process. `TERM=xterm` tells bash it's running in a terminal.

**PTYInstance — holds OS resources:**

```go
type PTYInstance struct {
    cmd    *exec.Cmd
    pty    *os.File
    output chan []byte
    done   chan struct{}
}
```

Four fields: the process, the PTY file descriptor, an output channel, and a done signal.

**readLoop — goroutine that reads PTY output:**

```go
func (p *PTYInstance) readLoop() {
    defer close(p.done)
    defer close(p.output)
    buf := make([]byte, 4096)
    for {
        n, err := p.pty.Read(buf)
        if err != nil { return }
        data := make([]byte, n)
        copy(data, buf[:n])
        select {
        case p.output <- data:
        default:
            log.Printf("pty output buffer full, dropping data")
        }
    }
}
```

Key decisions:
- 4096-byte read buffer (matches typical PTY block size)
- Copy before sending (prevents buffer reuse issues)
- Non-blocking send with `select/default` — drops data if buffer is full
- Output channel has buffer of 256

**Close — kills process and closes file:**

```go
func (p *PTYInstance) Close() error {
    p.cmd.Process.Kill()
    return p.pty.Close()
}
```

No graceful shutdown (no SIGTERM → wait → SIGKILL). Just kills immediately.

---

## Trace Exercise

```
PTYManager.Spawn()
    │
    ├─→ exec.Command("bash")
    ├─→ cmd.Env = ["TERM=xterm"]
    ├─→ pty.Start(cmd)
    │       │
    │       ├─→ Creates PTY pair (/dev/ptyXX, /dev/ttyXX)
    │       ├─→ Forks bash process attached to ttyXX
    │       └─→ Returns ptmx (*os.File)
    │
    ├─→ PTYInstance{cmd, ptmx, output: chan[256], done: chan}
    ├─→ go instance.readLoop()
    │       │
    │       └─→ Loop: pty.Read(buf) → copy → output <- data
    │
    └─→ return instance, nil
```

---

## Architecture Questions

1. The output channel has a buffer of 256. Each message can be up to 4096 bytes. What's the maximum data in-flight before dropping? Is this enough for `cat /usr/share/dict/words`?
2. `Close()` calls `Process.Kill()` which sends SIGKILL. Should it try SIGTERM first?
3. `readLoop` closes `output` channel when PTY exits. What happens to consumers reading from that channel?
4. Why copy the buffer data before sending to the channel? What happens if you send `buf[:n]` directly?

---

## Hands-on Exercise

Calculate the maximum output buffering capacity:

- Channel buffer: 256 messages
- Max message size: 4096 bytes
- Total: 256 × 4096 = 1 MB

If bash outputs faster than WebSocket clients consume, when does data start dropping? What's the impact?

---

## Review Checklist

- [ ] I can explain why PTYProvider and PTYInstance are separate interfaces
- [ ] I understand how creack/pty spawns a bash process
- [ ] I can trace the readLoop goroutine
- [ ] I understand the non-blocking send pattern
- [ ] I know what happens when the output buffer fills
- [ ] I can explain the Close() behavior

---

## Notes

---

## Questions for Later

---

## Estimated Duration

40 minutes

---

# Session 9 — The Streaming Pipeline — PTY to WebSocket

## Goal

Understand the complete data path from PTY output to WebSocket clients, and identify the critical architectural issues in the current implementation.

---

## Files to Read

- `internal/transport/pty_streamer.go`
- Re-read `internal/transport/websocket.go` (broadcast target)
- Re-read `internal/infrastructure/pty_manager.go` (readLoop)

---

## Before You Read

1. How does the system discover that a new session has a PTY ready?
2. How does PTY output reach WebSocket clients?
3. What happens when multiple sessions are running simultaneously?
4. What is the difference between polling and event-driven discovery?

---

## Concepts to Learn

**PTYStreamer — bridges PTY output to broadcast:**

```go
type PTYStreamer struct {
    cm      *application.ConnectionManager
    streams map[string]bool
    mu      sync.Mutex
}
```

Tracks which sessions are already being streamed (deduplication).

**WatchSessions — polling-based discovery (buggy):**

```go
func (s *PTYStreamer) WatchSessions(sm *application.SessionManager) {
    ticker := make(chan struct{})     // unbuffered, never sent to
    go func() {
        for range ticker {            // NEVER executes after initial check
            s.checkSessions(sm)
        }
    }()
    s.checkSessions(sm)               // only this runs
    select {}                         // blocks forever
}
```

The ticker channel is never written to. The goroutine that was supposed to poll periodically never fires after the initial call. This is dead code.

**tryStream — exactly-once goroutine per session:**

```go
func (s *PTYStreamer) tryStream(sessionID string, sm *application.SessionManager) {
    s.mu.Lock()
    if s.streams[sessionID] { s.mu.Unlock(); return }  // already streaming
    pty, ok := sm.GetPTY(sessionID)
    if !ok { s.mu.Unlock(); return }                     // no PTY yet
    s.streams[sessionID] = true
    s.mu.Unlock()
    go s.stream(sessionID, pty.Output())
}
```

Uses a map + mutex to ensure exactly one goroutine per session.

**stream — reads channel, broadcasts to ALL connections:**

```go
func (s *PTYStreamer) stream(sessionID string, output <-chan []byte) {
    for data := range output {
        s.cm.Broadcast(data)       // broadcasts to ALL connections globally
    }
    s.mu.Lock()
    delete(s.streams, sessionID)
    s.mu.Unlock()
}
```

**Critical issue: `cm.Broadcast(data)` sends to EVERY connection in the system, not just connections in the same session.**

---

## Trace Exercise

Trace the complete output path:

```
bash produces output
    │
    ├─→ PTYInstance.readLoop() reads from pty.Read()
    │       │
    │       └─→ output <- data  (channel, buffer 256)
    │
    ├─→ PTYStreamer.stream() goroutine reads from output channel
    │       │
    │       └─→ cm.Broadcast(data)
    │               │
    │               ├─→ Connection 1 (session A) → WriteMessage → WebSocket
    │               ├─→ Connection 2 (session A) → WriteMessage → WebSocket
    │               ├─→ Connection 3 (session B) → WriteMessage → WebSocket  ← WRONG
    │               └─→ Connection 4 (session B) → WriteMessage → WebSocket  ← WRONG
    │
    └─→ All 4 clients receive session A's output
```

---

## Architecture Questions

1. The `WatchSessions` ticker is dead code. How would you fix this? What are the options?
2. `Broadcast` sends globally. How would you implement session-scoped broadcast? What data structure changes are needed?
3. When a PTY exits (output channel closes), the stream goroutine cleans up. But what happens to clients connected to that session? Are they notified?
4. If you had 100 sessions with 10 viewers each, how many goroutines would `PTYStreamer` create? Is this scalable?

---

## Hands-on Exercise

Design session-scoped broadcast:

1. What changes in `ConnectionManager`?
2. What new method is needed?
3. How does `PTYStreamer.stream()` call it differently?
4. Draw the new data flow.

---

## Review Checklist

- [ ] I can trace the complete PTY → WebSocket data path
- [ ] I understand the WatchSessions polling bug
- [ ] I understand the deduplication mechanism
- [ ] I can explain the global broadcast problem
- [ ] I know how many goroutines exist per session

---

## Notes

---

## Questions for Later

---

## Estimated Duration

45 minutes

---

# Session 10 — Input Handling — The Host Controls Everything

## Goal

Understand the permission model, how host input reaches the terminal, and the design trade-offs in the current implementation.

---

## Files to Read

- `internal/transport/websocket.go` (readLoop, lines 61–93)
- Re-read `internal/domain/role.go`

---

## Before You Read

1. How does the system distinguish between host and viewer input?
2. What happens when a viewer tries to type?
3. How are keystrokes transmitted — character by character or in batches?
4. What message format is used for terminal input?

---

## Concepts to Learn

**Role-based input control:**

```go
if role != domain.RoleHost {
    log.Printf("viewer %s attempted input, rejected", clientID)
    continue
}
```

Simple role check: only the host can send input. Viewers' messages are silently dropped (with a log entry).

**No message framing protocol:**

The WebSocket message body IS the terminal input. No JSON envelope, no message type field, no序列化. Raw bytes flow directly from browser to PTY.

**Input path:**

```
Browser keystroke → WebSocket message → readLoop → role check → pty.Write(message) → bash processes it
```

**Host identity problem:**

1. `CreateSession("host")` creates a session with `Host.ID = "host"`
2. First WebSocket client connects → `session.Host` is overwritten with that client's UUID
3. The original `"host"` ID is orphaned — it never connects

---

## Trace Exercise

```
Host types "ls\n" in browser
    │
    ├─→ Browser sends WebSocket message: [108, 115, 10]  (ASCII for "ls\n")
    │
    ├─→ readLoop receives message
    │       │
    │       ├─→ role == Host? → YES
    │       │
    │       ├─→ sm.GetPTY(sessionID) → PTYInstance
    │       │
    │       └─→ pty.Write([108, 115, 10])
    │               │
    │               └─→ p.pty.Write(data)  // writes to /dev/ptyXX
    │                       │
    │                       └─→ bash reads from stdin, executes "ls"
    │                               │
    │                               └─→ output → readLoop → stream → broadcast
```

Now trace what happens when a viewer tries to type:

```
Viewer types "rm -rf /\n"
    │
    ├─→ readLoop receives message
    │       │
    │       ├─→ role == Host? → NO
    │       │
    │       └─→ log "viewer X attempted input, rejected"
    │           (message is discarded, PTY never sees it)
```

---

## Architecture Questions

1. No message framing means browser and PTY share a raw byte stream. What happens if the browser sends binary data? What if it sends a partial UTF-8 sequence?
2. The host check is `role != domain.RoleHost` — a simple inequality. Should this be a method on `Role`? (e.g., `role.CanWrite()`)
3. `readLoop` handles input AND disconnect detection in the same goroutine. Should these be separated? What are the trade-offs?
4. What happens if the host disconnects? Does the session end? Does the PTY keep running?

---

## Hands-on Exercise

Design a message protocol. Instead of raw bytes, define a JSON message format:

```json
{
    "type": "input",
    "data": "ls -la\n"
}
```

What other message types would you need? What changes in the readLoop?

---

## Review Checklist

- [ ] I can explain the role-based input control
- [ ] I understand the raw byte input path
- [ ] I know the host identity problem
- [ ] I understand what happens when a viewer types
- [ ] I can design a basic message protocol

---

## Notes

---

## Questions for Later

---

## Estimated Duration

30 minutes

---

# Session 11 — Testing Strategy & Mock Patterns

## Goal

Understand how the project tests business logic without infrastructure dependencies, and identify what is NOT tested.

---

## Files to Read

- `internal/application/session_manager_test.go`
- Re-read `internal/application/session_store.go` and `pty_provider.go` (the interfaces being mocked)

---

## Before You Read

1. How do you test code that depends on a database or external process?
2. What is a mock? What is a stub? What is a fake?
3. Why define mocks in `_test.go` files instead of separate packages?
4. How many test cases should a simple CRUD service have?

---

## Concepts to Learn

**Three test doubles:**

| Mock | Replaces | Key Behavior |
|------|----------|-------------|
| `mockStore` | `MemoryStore` | In-memory map, returns errors for missing sessions |
| `mockPTY` | `PTYInstance` | Buffered channel, no real process |
| `mockPTYProvider` | `PTYManager` | Returns `mockPTY` from `Spawn()` |

**Why mocks work here:**

Interfaces are defined in `application`. Mocks implement the same interfaces. Tests import only `application` — no `infrastructure` needed. No import cycles.

**Test coverage:**

| Test | What it verifies |
|------|-----------------|
| `TestCreateSession` | Session has non-empty ID, host has correct ID |
| `TestGetSession` | Fetching a created session returns it |
| `TestGetSessionNotFound` | Missing session returns error |
| `TestDeleteSession` | Deleted session returns error on get |
| `TestDeleteSessionNotFound` | Deleting missing session returns error |

**What is NOT tested:**

- WebSocket connection handling
- PTY streaming
- Broadcast behavior
- ConnectionManager operations
- HTTP API endpoints
- PTYManager.Spawn()
- Concurrent access patterns

**Test style:**

- Standalone functions (not table-driven)
- Behavior-focused assertions (not implementation-focused)
- One assertion topic per test

---

## Trace Exercise

Trace how the mock works:

```
TestCreateSession
    │
    ├─→ newMockStore()           // in-memory map, no real storage
    ├─→ &mockPTYProvider{}       // returns mockPTY, no real process
    ├─→ NewSessionManager(mockStore, mockPTYProvider)
    │
    ├─→ manager.CreateSession("host-1")
    │       │
    │       ├─→ mockStore.Create(session)  // stores in map
    │       ├─→ mockPTYProvider.Spawn()    // returns mockPTY (no real bash)
    │       └─→ return session, nil
    │
    └─→ Assert session.ID != ""
    └─→ Assert session.Host.ID == "host-1"
```

---

## Architecture Questions

1. `mockStore.Create` always returns `nil` error. Should you test the failure case? How?
2. `mockPTY.Write` always returns `nil` error. Is this sufficient for testing input handling?
3. There are no tests for `ConnectionManager`, `PTYStreamer`, or WebSocket handlers. Why might these be harder to test?
4. If you wanted to test the WebSocket handler, what would you need? (Hint: look at `httptest` and `gorilla/websocket` test utilities)

---

## Hands-on Exercise

Write a test for `ConnectionManager.Broadcast()`:

1. Create a ConnectionManager
2. Register 3 mock WebSocket connections
3. Call Broadcast with a message
4. Verify all 3 connections received the message

What mock do you need for `*websocket.Conn`? Is this straightforward?

---

## Review Checklist

- [ ] I can list all three mocks and what they replace
- [ ] I understand why mocks are in `_test.go` files
- [ ] I know what IS tested and what is NOT tested
- [ ] I can explain the behavior-focused assertion approach
- [ ] I understand why WebSocket testing is harder

---

## Notes

---

## Questions for Later

---

## Estimated Duration

35 minutes

---

# Session 12 — Architectural Gaps, Trade-offs & Concurrency

## Goal

Identify every architectural weakness, trade-off, and potential concurrency issue in the codebase. This session turns understanding into critical analysis.

---

## Files to Read

- ALL `.go` files (comparative pass)
- Re-read `docs/architecture.md` § "Architectural Principles"

---

## Before You Read

1. What would break if two users created sessions at the exact same time?
2. What happens when a PTY produces output faster than clients can receive it?
3. How does the system handle a crashed PTY process?
4. What security vulnerabilities exist in the current implementation?

---

## Concepts to Learn

**Architectural gap inventory:**

| # | Gap | Severity | Details |
|---|-----|----------|---------|
| 1 | Global broadcast | Critical | `ConnectionManager.Broadcast()` sends to ALL connections — cross-session data leakage |
| 2 | WatchSessions dead code | High | `ticker` channel never fires — polling loop is inoperative |
| 3 | No graceful shutdown | High | No signal handling; PTY processes leak on server exit |
| 4 | No session-scoped connection tracking | High | `Session.Clients` slice exists but is never populated |
| 5 | Hardcoded host ID | Medium | `CreateSession("host")` creates phantom host; real host assigned on WebSocket connect |
| 6 | No message framing | Medium | Raw bytes instead of typed messages — no protocol for join/output/error |
| 7 | PTY output drops silently | Medium | Buffer of 256 with non-blocking send — data lost without notification |
| 8 | No MemoryStore concurrency | Low | Plain map with no mutex — relies on SessionManager's mutex for safety |
| 9 | CheckOrigin always true | Low | Security concern for production |
| 10 | No graceful PTY shutdown | Low | `Process.Kill()` instead of SIGTERM → wait → SIGKILL |

**Concurrency analysis:**

```
MemoryStore.sessions map
    │
    ├─→ Accessed by SessionManager.CreateSession (via store.Create)
    ├─→ Accessed by SessionManager.GetSession (via store.Get)
    ├─→ Accessed by SessionManager.DeleteSession (via store.Delete)
    │
    ├─→ Protected by SessionManager.mu (RWMutex)
    │
    └─→ BUT: store.Create/Get/Delete are called INSIDE SessionManager's lock
         → MemoryStore itself has no lock
         → Safe ONLY if all access goes through SessionManager
         → If anyone accesses MemoryStore directly → data race
```

**Goroutine count analysis:**

```
Per session:
    1 goroutine: PTYInstance.readLoop (reads PTY stdout)
    1 goroutine: PTYStreamer.stream (reads output channel, broadcasts)
    1 goroutine: Server.readLoop (reads WebSocket messages from host)

Per server:
    1 goroutine: PTYStreamer.WatchSessions (dead code — never fires)
    1 goroutine: http.ListenAndServe (main thread)

Total per session: 3 goroutines
Total for N sessions: 3N + 2 goroutines
```

---

## Trace Exercise

Trace the failure scenario: PTY process crashes mid-session.

```
bash process crashes (SIGKILL from OOM killer)
    │
    ├─→ pty.Read() returns error
    │
    ├─→ PTYInstance.readLoop:
    │       │
    │       ├─→ err != io.EOF? → log "pty read error"
    │       ├─→ close(p.done)
    │       └─→ close(p.output)
    │
    ├─→ PTYStreamer.stream:
    │       │
    │       ├─→ range output → channel closed, loop exits
    │       ├─→ log "stream ended for session X"
    │       └─→ delete(s.streams, sessionID)
    │
    ├─→ WebSocket clients:
    │       │
    │       └─→ NO NOTIFICATION → clients see frozen terminal
    │
    └─→ Session still exists in MemoryStore
        Host still connected via WebSocket
        But no PTY — input goes nowhere (GetPTY returns false → return)
```

---

## Architecture Questions

1. Which of the 10 gaps would you fix first for production? Why?
2. How would you implement session-scoped broadcast without breaking the existing interface?
3. The `Connection.Conn` field is `*websocket.Conn` (concrete type). The architecture decision says this was chosen for type safety. What's the alternative and why was it rejected?
4. If you added authentication, where would token validation happen? Which layer?

---

## Hands-on Exercise

For each of the 10 gaps, write a one-sentence fix:

1. Global broadcast →
2. WatchSessions dead code →
3. No graceful shutdown →
4. No session-scoped connections →
5. Hardcoded host ID →
6. No message framing →
7. PTY output drops →
8. No MemoryStore concurrency →
9. CheckOrigin always true →
10. No graceful PTY shutdown →

---

## Review Checklist

- [ ] I can list all 10 architectural gaps
- [ ] I can explain why each gap exists (MVP trade-off vs. bug)
- [ ] I understand the concurrency model
- [ ] I can count goroutines per session
- [ ] I know what happens when a PTY crashes
- [ ] I can prioritize fixes for production

---

## Notes

---

## Questions for Later

---

## Estimated Duration

50 minutes

---

# Session 13 — Redesign Exercise — Build It From Scratch

## Goal

Synthesize everything you've learned into a complete mental model. Without opening the code, you should be able to explain every design decision and propose improvements.

---

## Files to Read

None. Close the codebase. This session is pure synthesis.

---

## Before You Read

1. Can you draw the full architecture from memory?
2. Can you explain why each interface exists and where it lives?
3. Can you trace every data flow (create → join → type → see output)?
4. Can you identify every gap and propose a fix?

---

## Concepts to Learn

This session has no new concepts. It tests whether you've internalized the previous 12 sessions.

---

## Trace Exercise

Without opening any file, answer these questions:

**Architecture from memory:**

Draw the 4-layer diagram. Label every package. Draw every dependency arrow.

**Interface catalog:**

| Interface | Package | Methods | Implemented By |
|-----------|---------|---------|----------------|
| SessionStore | application | Create, Get, Delete | MemoryStore |
| PTYProvider | application | Spawn, Stop | PTYManager |
| PTYInstance | application | Write, Output, Close | PTYInstance (struct) |

**Goroutine inventory:**

List every goroutine in the system, what it does, and when it starts/stops.

**Mutex inventory:**

List every mutex, what it protects, and whether it's a read or write lock.

**Data flow traces:**

1. Session creation: POST /session → ?
2. WebSocket connection: /ws?sessionID=X → ?
3. Host types "ls": browser → ?
4. PTY output reaches viewers: bash → ?
5. Host disconnects: browser closes → ?

---

## Architecture Questions

**Deep redesign questions:**

1. How would you redesign broadcasting to be session-scoped? What data structure changes are needed in `ConnectionManager`?
2. How would you replace the polling-based `PTYStreamer` with an event-driven design? What notification mechanism would you use?
3. What message protocol would you design for WebSocket communication? Define at least 4 message types.
4. How would you add authentication without breaking the architecture? Where does token validation happen?
5. How would you add session persistence? What changes if you swap `MemoryStore` for Postgres?
6. How would you isolate PTYs in Docker containers? What new interface would `PTYProvider` need?
7. How would you scale horizontally? What state is per-server vs. shareable?
8. What would you change FIRST for production? Justify your priority.

---

## Hands-on Exercise

**The complete redesign:**

Write a 1-page design document for a production version of this system. Include:

1. Architecture diagram
2. All interfaces (new and modified)
3. Message protocol
4. Session-scoped broadcast design
5. Authentication approach
6. Graceful shutdown strategy
7. Monitoring and observability

---

## Review Checklist

- [ ] I can draw the full architecture from memory
- [ ] I can explain every interface
- [ ] I can explain every goroutine
- [ ] I can explain every mutex
- [ ] I can trace every data flow
- [ ] I can explain every layer
- [ ] I can explain every design decision
- [ ] I can explain every trade-off
- [ ] I can explain every architectural gap
- [ ] I can propose a production redesign

---

## Notes

---

## Questions for Later

---

## Estimated Duration

45 minutes

---

# Common Architectural Mistakes

All known weaknesses discovered during the review, organized by severity.

## Critical

### Global Broadcast Instead of Session-Scoped Broadcast

**What:** `ConnectionManager.Broadcast()` sends messages to ALL connected clients across ALL sessions.

**Why it's critical:** PTY output from session A is delivered to viewers of session B. In any multi-session deployment, this leaks data between unrelated sessions. It also wastes bandwidth — every client processes every message regardless of session membership.

**Root cause:** `ConnectionManager` was designed as a flat connection registry without session awareness. The `connections` map is keyed by client ID, not session ID.

**Impact:** Functional in single-session scenarios. Broken for any real multi-session use.

---

## High

### WatchSessions Polling Loop Is Dead Code

**What:** `PTYStreamer.WatchSessions` creates an unbuffered channel (`ticker`) that is never sent to. The goroutine iterating over `ticker` never executes after the initial `checkSessions` call.

**Why it's high:** New sessions created after server startup will never have their PTY output streamed. Only sessions that existed when `WatchSessions` was first called will work.

**Root cause:** Channel-based ticker was implemented instead of `time.NewTicker`. The intent was periodic polling, but the implementation never triggers.

---

### No Graceful Shutdown

**What:** When the server process is killed, PTY processes (bash instances) are not terminated. They become orphaned.

**Why it's high:** Orphaned processes consume system resources. On a server with many sessions, repeated restarts accumulate zombie bash processes.

**Root cause:** No signal handling (`SIGTERM`, `SIGINT`) in `main.go`. `ListenAndServe` blocks until fatal error.

---

### Session.Clients Is Never Populated

**What:** The `Session` struct has a `Clients []*Client` field, but no code ever adds clients to it. The field exists in the domain model but is unused.

**Why it's high:** The domain model promises a feature (tracking connected clients) that doesn't exist. Any future code relying on `Session.Clients` will get an empty slice.

---

## Medium

### Hardcoded Host ID Creates Orphaned Entity

**What:** `CreateSession("host")` creates a `Session.Host` with ID `"host"`. When the first WebSocket client connects, `session.Host` is overwritten with the client's UUID. The `"host"` entity is orphaned.

**Why it's medium:** Creates a misleading domain state. The session briefly has a host that doesn't correspond to any real connection.

---

### No Message Framing Protocol

**What:** WebSocket messages are raw bytes with no envelope. There's no way to distinguish "input data" from "join request" from "error notification."

**Why it's medium:** Limits extensibility. Adding new message types requires ad-hoc parsing. No way to send structured metadata alongside terminal data.

---

### PTY Output Drops Silently

**What:** When the PTY output channel buffer (256 messages) is full, new output is dropped with a log message. Clients never know data was lost.

**Why it's medium:** Terminal state can become inconsistent between what the host sees and what viewers see. No recovery mechanism exists.

---

## Low

### MemoryStore Has No Mutex

**What:** `MemoryStore` uses a plain `map[string]*domain.Session` with no synchronization primitive.

**Why it's low:** Currently safe because all access goes through `SessionManager`, which holds a mutex. However, this is a latent risk — any direct access to `MemoryStore` from a new goroutine would cause a data race.

---

### CheckOrigin Always Returns True

**What:** The WebSocket upgrader accepts connections from any origin.

**Why it's low:** Acceptable for local development. Would need to validate allowed origins in production to prevent cross-site WebSocket hijacking.

---

### No Graceful PTY Shutdown

**What:** `PTYInstance.Close()` calls `Process.Kill()` (SIGKILL) instead of trying SIGTERM first and waiting for graceful exit.

**Why it's low:** Bash doesn't typically need graceful shutdown, but child processes (editors, long-running commands) may not clean up properly.

---

# Final Review Challenge

Without opening the code, the engineer should be able to:

## Architecture

- [ ] Draw the 4-layer architecture diagram with all packages labeled
- [ ] Draw every dependency arrow between packages
- [ ] Explain why each layer exists and what it's responsible for

## Interfaces

- [ ] Name every interface, its package, its methods, and its implementor
- [ ] Explain why each interface lives where it does
- [ ] Explain what would break if any interface were moved

## Goroutines

- [ ] List every goroutine in the system
- [ ] Explain what each goroutine does
- [ ] Explain when each goroutine starts and stops
- [ ] Calculate total goroutine count for N sessions

## Mutexes

- [ ] List every mutex in the system
- [ ] Explain what each mutex protects
- [ ] Explain the read/write lock pattern

## Data Flows

- [ ] Trace session creation from HTTP request to stored session
- [ ] Trace WebSocket connection from browser to readLoop
- [ ] Trace host input from browser keystroke to PTY
- [ ] Trace PTY output from bash to all connected viewers
- [ ] Trace disconnect cleanup from browser close to resource release

## Layers

- [ ] Explain what belongs in transport vs. application vs. domain vs. infrastructure
- [ ] Explain the dependency direction between layers
- [ ] Give an example of a layer violation and why it's harmful

## Design Decisions

- [ ] Explain why domain models have no infrastructure imports
- [ ] Explain why interfaces live in the application layer
- [ ] Explain why ConnectionManager stores `*websocket.Conn` directly
- [ ] Explain why PTY and Session are not coupled in the domain
- [ ] Explain why tests use local mocks instead of real implementations

## Trade-offs

- [ ] Explain the MVP trade-offs (global broadcast, no auth, no persistence)
- [ ] Explain what you gain and what you lose with each trade-off
- [ ] Prioritize which trade-offs to address first for production

## Architectural Gaps

- [ ] List all 10 gaps from Session 12
- [ ] Explain the severity of each
- [ ] Propose a fix for each

---

# Production Redesign Challenge

Open-ended redesign questions. There are no single correct answers — the goal is to demonstrate deep understanding of the trade-offs involved.

## Broadcasting

How would you redesign broadcasting to be session-scoped? What data structure changes are needed? How do you handle session creation/deletion in the connection registry?

## Polling Elimination

How would you replace the polling-based `PTYStreamer` with an event-driven design? What notification mechanism would you use when a new session is created? How do you ensure exactly-once streaming?

## Authentication

How would you support authentication? Where does token validation happen? How does it interact with the session join flow? What changes in the WebSocket handshake?

## Persistence

How would you support session persistence? What changes if you swap `MemoryStore` for Postgres? What happens to in-flight PTY processes when the server restarts? How do you handle session recovery?

## PTY Isolation

How would you isolate PTYs in Docker containers? What new interface would `PTYProvider` need? How does container lifecycle map to session lifecycle? What about resource limits?

## Horizontal Scaling

How would you scale horizontally? What state is per-server vs. shareable? How do you handle PTY output reaching viewers connected to different servers? What about sticky sessions?

## Priority

If you could change ONE thing for production, what would it be and why? Justify your choice considering impact, effort, and risk.

---

*End of Engineering Review Workbook*
