# Tasks

This document contains all implementation tasks for the MVP.

Rules:

* Every task must be completed independently.
* Every task must satisfy its acceptance criteria.
* Do not implement features outside the documented scope.
* Follow architecture.md and requirements.md.

---

# Phase 1 - Foundation

Goal:
Create the project foundation and core domain structure.

---

## TASK-001 - Initialize Project

Description:
Initialize Go project and repository structure.

Acceptance Criteria:

* go.mod exists
* Project compiles successfully
* Folder structure follows architecture.md

Dependencies:
None

---

## TASK-002 - Create Domain Models

Description:
Create core domain entities.

Acceptance Criteria:

* Session entity exists
* Client entity exists
* Role definitions exist
* Domain models contain no infrastructure code

Dependencies:
TASK-001

---

## TASK-003 - Create SessionStore Interface

Description:
Define storage abstraction for session persistence.

Acceptance Criteria:

* SessionStore interface exists
* Create operation defined
* Get operation defined
* Delete operation defined

Dependencies:
TASK-002

---

## TASK-004 - Implement MemoryStore

Description:
Implement in-memory storage for sessions.

Acceptance Criteria:

* MemoryStore implements SessionStore
* Sessions stored in memory
* Create works
* Get works
* Delete works

Dependencies:
TASK-003

---

## TASK-005 - Create Session Manager

Description:
Create application service responsible for session lifecycle.

Acceptance Criteria:

* CreateSession implemented
* GetSession implemented
* DeleteSession implemented
* Uses SessionStore interface only

Dependencies:
TASK-004

---

## TASK-006 - Unit Test Session Manager

Description:
Verify session lifecycle operations.

Acceptance Criteria:

* Create session test
* Get session test
* Delete session test

Dependencies:
TASK-005

---

# Phase 2 - Realtime Communication

Goal:
Enable realtime communication through WebSocket.

---

## TASK-007 - Create WebSocket Endpoint

Description:
Expose WebSocket endpoint for client connections.

Acceptance Criteria:

* Endpoint created
* Upgrade request works
* Connection established

Dependencies:
TASK-005

---

## TASK-008 - Create Connection Manager

Description:
Manage active WebSocket connections.

Acceptance Criteria:

* Register client
* Unregister client
* Track active connections

Dependencies:
TASK-007

---

## TASK-009 - Implement Session Join Flow

Description:
Allow users to join existing sessions.

Acceptance Criteria:

* Join request handled
* Session validated
* Client registered

Dependencies:
TASK-008

---

## TASK-010 - Implement Broadcast System

Description:
Broadcast messages to all session participants.

Acceptance Criteria:

* Broadcast function exists
* Multiple clients receive same message
* Failed connections handled safely

Dependencies:
TASK-009

---

# Phase 3 - Terminal Integration

Goal:
Attach a real terminal to each session.

---

## TASK-011 - Integrate PTY Library

Description:
Add PTY support using creack/pty.

Acceptance Criteria:

* PTY created successfully
* Shell process starts
* Terminal output readable

Dependencies:
TASK-005

---

## TASK-012 - Create PTY Manager

Description:
Manage PTY lifecycle.

Acceptance Criteria:

* Start PTY
* Stop PTY
* Read output

Dependencies:
TASK-011

---

## TASK-013 - Stream PTY Output

Description:
Forward terminal output to session clients.

Acceptance Criteria:

* Output captured
* Output forwarded
* Multiple clients receive output

Dependencies:
TASK-012
TASK-010

---

## TASK-014 - Handle Terminal Input

Description:
Send host input to terminal.

Acceptance Criteria:

* Host input accepted
* Viewer input rejected
* Commands executed successfully

Dependencies:
TASK-013

---

# Phase 4 - Frontend

Goal:
Provide a usable browser interface.

---

## TASK-015 - Create Frontend Skeleton

Description:
Initialize frontend application.

Acceptance Criteria:

* Frontend starts successfully
* Build process works

Dependencies:
TASK-014

---

## TASK-016 - Create Session Page

Description:
Display session interface.

Acceptance Criteria:

* Session page exists
* Session ID visible

Dependencies:
TASK-015

---

## TASK-017 - Integrate Terminal UI

Description:
Render terminal output in browser.

Acceptance Criteria:

* Output displayed
* Updates in realtime

Dependencies:
TASK-016

---

## TASK-018 - Implement Session Join UI

Description:
Allow users to join sessions.

Acceptance Criteria:

* Join form exists
* Session connection works

Dependencies:
TASK-017

---

# Phase 5 - Hardening

Goal:
Improve reliability and maintainability.

---

## TASK-019 - Add Structured Logging

Acceptance Criteria:

* Session events logged
* Connection events logged
* Errors logged

Dependencies:
TASK-014

---

## TASK-020 - Add Error Handling

Acceptance Criteria:

* Errors propagated correctly
* Unexpected failures handled safely

Dependencies:
TASK-019

---

## TASK-021 - Add Integration Tests

Acceptance Criteria:

* Session creation tested
* Join flow tested
* Terminal flow tested

Dependencies:
TASK-020

---

# MVP Completion Criteria

The MVP is considered complete when:

* Session creation works
* Session joining works
* WebSocket communication works
* Shared terminal works
* Host input works
* Viewer mode works
* Tests pass
