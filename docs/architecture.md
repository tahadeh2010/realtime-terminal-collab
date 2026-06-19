# Architecture

## High Level Design

Client
↓
WebSocket Gateway
↓
Session Manager
↓
PTY Manager
↓
Shell Process

## Components

### WebSocket Gateway

Responsibilities:

* Connection handling
* Message routing
* Session join requests

### Session Manager

Responsibilities:

* Session lifecycle
* Client registration
* Broadcast events

### PTY Manager

Responsibilities:

* Spawn terminal
* Read output
* Forward output to Session Manager

### Shell Process

Responsibilities:

* Execute commands
* Produce terminal output

## Domain Objects

Session

* ID
* Clients
* Host
* PTY

Client

* ID
* Connection
* Role

Roles

* Host
* Viewer

## Future Extensions

* Docker isolation
* Recording
* Authentication
* Session persistence
* Multi-host mode

## Architectural Principles

1. Depend on interfaces, not implementations.
2. Business logic must not depend on transport layer.
3. Business logic must not depend on storage implementation.
4. Every module should have a single responsibility.
5. Realtime communication is performed through WebSocket only.

## Layers

Transport Layer
- WebSocket
- HTTP

Application Layer
- Session Service
- Permission Service

Domain Layer
- Session
- Client
- Role

Infrastructure Layer
- MemoryStore
- PTY

## Folder Structure

cmd/
    server/

internal/
    domain/
    application/
    infrastructure/
    transport/

docs/

## Responsibilities

Session Service
- Create session
- Delete session
- Join session

PTY Manager
- Create PTY
- Read PTY output

WebSocket Gateway
- Manage connections
- Route messages

## Deferred Features

The following features are intentionally excluded from MVP:

- Authentication
- Docker Isolation
- Session Persistence
- Recording
- Multi Cursor