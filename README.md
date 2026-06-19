# Realtime Terminal Collaboration

A real-time collaborative terminal platform built with Go.

## Overview

Realtime Terminal Collaboration is a system that allows multiple users to share and collaborate on a terminal session in real time.

The goal of this project is to build a clean, scalable monolithic architecture while practicing an AI-assisted development workflow.

## Project Goals

- Build a real-time terminal sharing system
- Learn WebSocket-based communication
- Manage terminal sessions using PTY
- Create an AI-agent-friendly project structure
- Practice professional software development workflow

## MVP Features

- Create terminal sessions
- Join existing sessions
- Real-time terminal output synchronization
- Host-controlled terminal input
- Multiple users connected to one session

## Non-MVP Features

The following features are intentionally excluded from the first version:

- Authentication
- Persistent sessions
- Docker isolation
- Terminal recording
- Multi-user input control
- Production deployment

## Tech Stack

Backend:
- Go

Communication:
- WebSocket

Terminal:
- PTY

Storage:
- SessionStore Interface
- MemoryStore (MVP)

## Architecture

High-level architecture:

Client
|
WebSocket
|
Transport Layer
|
Application Layer
|
Session Manager
|
Infrastructure Layer
|
PTY + Session Storage

The project follows a modular monolithic architecture.

## Project Structure

.
├── cmd/
├── internal/
│ ├── domain/
│ ├── application/
│ ├── infrastructure/
│ └── transport/
│
├── docs/
│ ├── vision.md
│ ├── requirements.md
│ ├── architecture.md
│ └── roadmap.md
│
└── README.md


## Documentation

Project documentation:

- [Vision](docs/vision.md)
- [Requirements](docs/requirements.md)
- [Architecture](docs/architecture.md)
- [Roadmap](docs/roadmap.md)

## Development Status

🚧 Currently in MVP development.

## License

MIT