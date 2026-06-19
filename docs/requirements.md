# Requirements

## Functional Requirements

### Session Management

FR-001: User can create a session.

FR-002: System generates unique session ID.

FR-003: User can join existing session.

FR-004: Multiple users can join one session.

### Terminal

FR-005: System creates a PTY-backed shell.

FR-006: Host can send terminal input.

FR-007: Shell output is broadcast to all connected clients.

FR-008: All connected clients must receive terminal output in the same order it is produced by the PTY.

### WebSocket

FR-009: Client communicates through WebSocket.

FR-010: System handles disconnects gracefully.

### Permissions

FR-011: Host role exists.

FR-012: Viewer role exists.

## Non Functional Requirements

NFR-001: Average latency < 200ms.

NFR-002: Session creation < 1 second.

NFR-003: Support at least 10 concurrent users per session.

NFR-004: Clear logging.

NFR-005: Modular architecture.

NFR-006: AI-agent-friendly codebase.

## Constraints

* Backend language: Go
* Transport: WebSocket
* Terminal: PTY
* Database: None for MVP

## Session Lifecycle

FR-013: A session is created by a host.

FR-014: Each session has exactly one host.

FR-015: A session is terminated when the host disconnects.

FR-016: New viewers cannot join a terminated session.

## Assumptions

- Sessions are ephemeral.
- Session data is stored in memory during MVP.
- Session persistence is out of scope.
- Authentication is out of scope.
- Docker isolation is out of scope.