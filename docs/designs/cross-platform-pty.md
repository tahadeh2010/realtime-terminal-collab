# Cross-Platform PTY Support

**Status:** Proposed  
**Date:** 2026-07-01  
**Scope:** Infrastructure layer only — application and domain layers remain unchanged  
**Type:** Architecture Decision Document (ADD)

---

## 1. Problem Statement

The PTY implementation relies on `github.com/creack/pty`, which uses POSIX `forkpty`/`openpty` syscalls. This makes the project impossible to compile or run on Windows. The project must support both Linux and Windows while preserving the existing Clean Architecture.

**Affected files:**

| File | Issue |
|------|-------|
| `internal/infrastructure/pty_manager.go` | Imports `creack/pty`, calls `pty.Start()`, hardcodes `bash` |
| `internal/infrastructure/pty_instance.go` | Uses `*os.File` for PTY file descriptor |
| `go.mod` | `creack/pty` dependency fails on Windows |

**Unaffected files (no changes needed):**

| File | Reason |
|------|--------|
| `internal/application/pty_provider.go` | Pure interface definitions — no OS imports |
| `internal/application/session_manager.go` | Consumes `PTYProvider` interface only |
| `internal/domain/*` | No OS awareness |
| `internal/transport/*` | Consumes `PTYInstance.Output()` channel only |
| `cmd/server/main.go` | Requires only a platform-appropriate factory call |

---

## 2. Current Architecture Analysis

### Dependency Flow

```
┌─────────────────────────────────────────────────────┐
│  cmd/server/main.go                                 │
│  ┌───────────────────────────────────────────────┐  │
│  │  infrastructure.NewPTYManager() → PTYProvider │  │
│  └───────────────────┬───────────────────────────┘  │
└──────────────────────┼──────────────────────────────┘
                       │ implements
┌──────────────────────▼──────────────────────────────┐
│  application/session_manager.go                     │
│  ┌───────────────────────────────────────────────┐  │
│  │  SessionManager.store   → SessionStore (iface) │  │
│  │  SessionManager.ptyManager → PTYProvider (iface)│  │
│  └───────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────┘
                       │ depends on
┌──────────────────────▼──────────────────────────────┐
│  application/pty_provider.go                        │
│  ┌───────────────────────────────────────────────┐  │
│  │  PTYProvider interface                        │  │
│  │    Spawn() → (PTYInstance, error)             │  │
│  │    Stop(PTYInstance) → error                  │  │
│  │                                               │  │
│  │  PTYInstance interface                        │  │
│  │    Write([]byte) → error                      │  │
│  │    Output() → <-chan []byte                   │  │
│  │    Close() → error                            │  │
│  └───────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────┘
                       │ implemented by
┌──────────────────────▼──────────────────────────────┐
│  infrastructure/pty_manager.go  ← Unix-only (HERE) │
│  infrastructure/pty_instance.go ← Unix-only (HERE) │
│  ┌───────────────────────────────────────────────┐  │
│  │  imports creack/pty                           │  │
│  │  uses exec.Command("bash")                    │  │
│  │  uses *os.File for PTY fd                     │  │
│  └───────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────┘
```

### Key Finding

The application layer already provides a correct dependency inversion boundary via `PTYProvider` and `PTYInstance` interfaces. The platform-specific code is entirely confined to 2 files in `internal/infrastructure/`. The fix is surgical — only the infrastructure layer needs platform-specific implementations.

### Separation of Concerns Issue

The current `PTYManager.Spawn()` mixes two distinct responsibilities:

1. **Shell discovery** — determining which shell executable to launch (`bash`)
2. **PTY lifecycle** — creating the pseudo-terminal, starting the process, reading output

These are independent concerns. Shell discovery is a platform-specific policy decision that should not be coupled to PTY mechanics. Extracting it into a dedicated component makes both responsibilities independently testable and replaceable.

---

## 3. Design Options Considered

### Option A: Go Build Tags

Use `//go:build` constraints to compile different source files per OS. Same package, different files selected at build time.

| Aspect | Assessment |
|--------|-----------|
| Idiomatic Go | Yes — used by the standard library |
| Runtime overhead | None — compile-time selection |
| Interface preservation | Perfect — same interfaces, different implementations |
| Maintenance burden | Two files to maintain, sharing the same contract |
| Build complexity | Minimal — `go build` handles tag selection automatically |

### Option B: Runtime Detection

Check `runtime.GOOS` at initialization and instantiate the correct implementation.

| Aspect | Assessment |
|--------|-----------|
| Idiomatic Go | No — build tags exist for this purpose |
| Runtime overhead | Minor but unnecessary |
| Compile-time safety | No — both implementations must compile for all targets |
| Practicality | Fails because Unix C headers are unavailable when cross-compiling to Windows |

### Option C: Separate Platform Packages

Create `infrastructure/pty_unix/` and `infrastructure/pty_windows/` packages.

| Aspect | Assessment |
|--------|-----------|
| Separation | Maximum |
| Boilerplate | High — requires factory/registration pattern in main.go |
| Over-engineering | Same solution as build tags with extra indirection |

### Option D: Third-Party Cross-Platform Library

Use a single library that wraps both platforms under one API.

| Aspect | Assessment |
|--------|-----------|
| Maturity | No mature cross-platform PTY library exists for Go |
| Control | Less control over ConPTY tuning and error handling |
| Dependency risk | Adds external dependency for a small API surface |

### Decision: Option A — Go Build Tags

Build tags are the simplest, most idiomatic, and most maintainable approach. The architecture already supports it perfectly through the existing interface boundaries.

---

## 4. Recommended Design

### Principle

The `PTYProvider` and `PTYInstance` interfaces in `internal/application/pty_provider.go` remain the sole contract. They are **not modified**. All platform-specific logic lives in `internal/infrastructure/` behind build tags.

### Target File Layout

```
internal/
  application/
    pty_provider.go              # UNCHANGED — interface definitions
  infrastructure/
    shell_discovery.go           # NEW — ShellFinder interface + platform detection
    shell_unix.go                # NEW — Unix shell discovery (build tag: !windows)
    shell_windows.go             # NEW — Windows shell discovery (build tag: windows)
    pty_manager.go               # RENAMED from current → Unix PTYManager (build tag: !windows)
    pty_instance.go              # RENAMED from current → Unix PTYInstance (build tag: !windows)
    pty_windows.go               # NEW — Windows PTYManager + PTYInstance (build tag: windows)
    memory_store.go              # UNCHANGED
```

**Rationale for renaming instead of deleting:** Renaming files (via `git mv`) preserves Git history. When reviewing `git log --follow` for `pty_manager.go`, the full history of the Unix implementation remains traceable. Deleting and recreating files severs that lineage.

---

## 5. Windows PTY Implementation — Engineering Decision

### Decision

Use `github.com/UserExistsError/conpty` as the Windows ConPTY wrapper.

### Background

Windows provides the Pseudo Console (ConPTY) API since Windows 10 version 1809 (build 17763, released October 2018). ConPTY is the official Windows mechanism for hosting terminal applications — it is the same API used by Windows Terminal, VS Code's integrated terminal, and Go's own `x/term` package.

### Alternatives Evaluated

| Option | Description | Assessment |
|--------|-------------|------------|
| `github.com/UserExistsError/conpty` | Thin ConPTY wrapper. 43 stars, MIT license, ~360 lines. Uses `golang.org/x/sys/windows` for syscalls. | **Selected.** Mature enough for production use. Minimal surface area. Clean Go API. |
| Direct syscall wrapper | Implement ConPTY calls manually using `golang.org/x/sys/windows`. | Rejected. ConPTY setup involves complex attribute list manipulation (`InitializeProcThreadAttributeList`, `UpdateProcThreadAttribute`). The `conpty` library already handles this correctly. Reimplementing risks subtle bugs for no architectural benefit. |
| `github.com/nicholasgasior/go-winpty` | Older winpty-based wrapper using `winpty.dll`. | Rejected. winpty is deprecated in favor of ConPTY. Adds a native DLL dependency. |
| No Windows support | Document Windows as unsupported. | Rejected. Windows support is a stated requirement. |

### Library Characteristics

- **API surface:** `conpty.Start(commandLine, options...)` returns a `*ConPty` with `Read()`, `Write()`, `Close()`, `Wait()`, `Resize()`.
- **Availability check:** `conpty.IsConPtyAvailable()` checks at runtime whether ConPTY syscalls exist on the current OS.
- **Error handling:** Returns `conpty.ErrConPtyUnsupported` on pre-1809 Windows, allowing graceful degradation.
- **Dependencies:** Only `golang.org/x/sys/windows` (Go team-maintained).
- **Build tag:** `//go:build windows` — the library itself is Windows-only.

### Long-Term Maintenance

- The ConPTY API is stable and maintained by Microsoft as part of the Windows platform.
- The `conpty` library is a thin wrapper (~360 lines) over 5 kernel32.dll functions. The maintenance surface is small.
- If the library becomes unmaintained, the wrapper is simple enough to fork or replace — it depends on no other third-party code.

### Compatibility Requirements

| Requirement | Value |
|-------------|-------|
| Minimum Windows version | Windows 10 version 1809 (build 17763) |
| Recommended Windows version | Windows 10 2004+ or Windows 11 |
| Fallback on older Windows | Return error with clear message at session creation time |

---

## 6. Shell Discovery — Separated Responsibility

### Problem

The current `PTYManager.Spawn()` hardcodes `exec.Command("bash")`. This is wrong on Windows (where `bash` may not exist) and fragile on Unix (where `zsh` or `fish` might be preferred).

Shell discovery is a distinct concern from PTY lifecycle management. It involves:
- Reading platform-specific environment variables
- Checking executable availability on `$PATH`
- Applying a fallback chain

This logic should not live inside `PTYManager` because:
1. **Testability** — Shell discovery can be unit-tested independently of PTY creation.
2. **Replaceability** — A future Docker or SSH provider might need different shell resolution logic.
3. **Clarity** — PTYManager's responsibility is clear: create a PTY and attach a process. How the process path is resolved is a separate question.

### Design

```
┌─────────────────────────────────────┐
│  ShellFinder interface (app layer)  │
│    FindShell() (string, error)      │
└──────────────┬──────────────────────┘
               │ implemented by
    ┌──────────┴──────────┐
    │                     │
┌───▼──────────┐  ┌──────▼──────────┐
│ shell_unix.go │  │ shell_windows.go │
│ (build tag:   │  │ (build tag:      │
│  !windows)    │  │  windows)        │
└───────────────┘  └─────────────────┘
```

### Platform-Specific Logic

**Unix (`shell_unix.go`):**

```
1. Read $SHELL environment variable
2. If set and executable → use it
3. Fallback chain: bash → zsh → sh
4. Verify each candidate exists on $PATH using exec.LookPath()
5. Return error if no shell found
```

**Windows (`shell_windows.go`):**

```
1. Read %COMSPEC% environment variable (typically C:\Windows\system32\cmd.exe)
2. If set and executable → use it
3. Fallback: powershell.exe (search $PATH)
4. Fallback: cmd.exe (absolute path from %SystemRoot%)
5. Return error if no shell found
```

### Interface Location

The `ShellFinder` interface belongs in `internal/application/` alongside `PTYProvider`, since it is a dependency that PTYManager consumes. It is injected into PTYManager at construction time.

```
// internal/application/shell_finder.go (NEW — interface only, no OS imports)

type ShellFinder interface {
    FindShell() (string, error)
}
```

### Benefits

- Shell discovery logic is independently testable with mock shell environments.
- Different PTY providers (Docker, SSH) can supply their own `ShellFinder` implementations.
- PTYManager remains focused on PTY mechanics only.

---

## 7. Build Tag Strategy

### Option Evaluated: `!windows` vs `linux || darwin`

| Tag | Platforms Covered | Extensibility | Maintenance |
|-----|-------------------|---------------|-------------|
| `//go:build !windows` | Linux, macOS, FreeBSD, OpenBSD, NetBSD, all non-Windows | Automatic — new Unix-like OSes compile without changes | Single file for all Unix-like platforms |
| `//go:build linux \|\| darwin` | Linux and macOS only | Requires manual addition for BSDs or other Unixes | More explicit, but incomplete |
| `//go:build linux` | Linux only | Requires separate files for macOS, BSDs | Maximum control, maximum files |

### Recommendation: `//go:build !windows`

Using `!windows` for the Unix file is preferable because:

1. **Future-proof.** If the project is ever built for FreeBSD, OpenBSD, or Plan 9, the Unix implementation compiles automatically. `creack/pty` supports all these platforms.
2. **Simplicity.** Two files total (Unix + Windows) covers all current and likely future platforms.
3. **Go standard library precedent.** The Go standard library uses `!windows` and `!plan9` patterns extensively for platform families.
4. **Override mechanism.** If a specific non-Windows platform needs special handling (e.g., macOS terminal quirks), a file with a more specific tag (e.g., `//go:build darwin`) overrides the catch-all. Go resolves from most-specific to least-specific.

If explicit platform enumeration is later preferred for safety, the tag can be narrowed to `//go:build linux || darwin || freebsd || openbsd || netbsd` without any interface or architecture changes.

---

## 8. Constructor Naming — `NewPTYManager` vs `NewPTYProvider`

### Option A: Keep `NewPTYManager()`

```go
ptyManager := infrastructure.NewPTYManager()
sm := application.NewSessionManager(store, ptyManager)
```

| Aspect | Assessment |
|--------|-----------|
| Clarity | Describes what it creates (a PTY manager) |
| Consistency | Matches existing codebase naming (`NewMemoryStore`, `NewSessionManager`) |
| Interface alignment | Returns concrete type, but assigned to interface at call site |

### Option B: Rename to `NewPTYProvider()`

```go
ptyProvider := infrastructure.NewPTYProvider()
sm := application.NewSessionManager(store, ptyProvider)
```

| Aspect | Assessment |
|--------|-----------|
| Clarity | Expresses that it produces the interface contract |
| Misleading | Implies it returns the interface type, but Go constructors return concrete types |
| Inconsistency | Breaks the existing pattern where infrastructure constructors return their struct |

### Recommendation: Keep `NewPTYManager()`

The constructor should remain `NewPTYManager()`. Reasons:

1. **Consistency with codebase conventions.** `NewMemoryStore()` returns `*MemoryStore`, not `SessionStore`. `NewPTYManager()` should return `*PTYManager`, not `PTYProvider`. This is standard Go practice.
2. **No architectural benefit.** Renaming the constructor adds churn without changing behavior. The caller already assigns to the interface type at the call site.
3. **Discoverability.** Developers searching for `PTYManager` in the codebase find the constructor. Searching for `PTYProvider` would not find it.

---

## 9. Data Flow

Both platforms produce the same data flow. The transport and application layers see no difference.

### Unix

```
WebSocket input → transport.readLoop → PTYInstance.Write()
  → os.File.Write() → PTY fd → Shell process

Shell output → PTY fd → os.File.Read() → readLoop goroutine
  → output channel → PTYStreamer → Broadcast to clients
```

### Windows

```
WebSocket input → transport.readLoop → PTYInstance.Write()
  → ConPty.Write() → cmdIn pipe → Shell process

Shell output → cmdOut pipe → ConPty.Read() → readLoop goroutine
  → output channel → PTYStreamer → Broadcast to clients
```

### Interface Contract (shared)

```
PTYInstance.Output() → <-chan []byte    ← both platforms produce this
PTYInstance.Write([]byte) → error      ← both platforms consume this
PTYInstance.Close() → error            ← both platforms support this
```

---

## 10. Dependency Changes

### `go.mod` additions

| Dependency | Platform | Purpose |
|-----------|----------|---------|
| `github.com/UserExistsError/conpty` | Windows | ConPTY wrapper |
| `golang.org/x/sys` | Windows | Transitive dependency of `conpty` |

No new dependencies for Unix. `creack/pty` remains.

### Module management caveat

`go mod tidy` on Windows will attempt to remove `creack/pty` from `go.mod` because it only appears in files excluded by the `!windows` build tag. To handle this:

- CI always runs `go mod tidy` on a Unix runner
- `go.mod` and `go.sum` are committed from Unix CI
- Windows CI runs `go build` and `go test` only (not `go mod tidy`)

---

## 11. Impact on Existing Code

| Layer | Files | Impact |
|-------|-------|--------|
| Domain | `session.go`, `client.go`, `role.go` | **None** |
| Application | `pty_provider.go` | **None** — interfaces unchanged |
| Application | `shell_finder.go` | **Created** — new `ShellFinder` interface |
| Application | `session_manager.go` | **None** — uses interfaces |
| Application | `session_manager_test.go` | **None** — mock PTYProvider already works |
| Application | `connection_manager.go`, `session_store.go` | **None** |
| Transport | `websocket.go`, `session_handler.go`, `pty_streamer.go` | **None** — uses interfaces |
| Infrastructure | `memory_store.go` | **None** |
| Infrastructure | `pty_manager.go` | **Renamed** — becomes Unix-only (build tag) |
| Infrastructure | `pty_instance.go` | **Renamed** — becomes Unix-only (build tag) |
| Infrastructure | `shell_unix.go` | **Created** — Unix shell discovery |
| Infrastructure | `shell_windows.go` | **Created** — Windows shell discovery |
| Infrastructure | `pty_windows.go` | **Created** — Windows PTY implementation |
| Entry point | `cmd/server/main.go` | **Minor** — instantiate platform-appropriate `ShellFinder` |
| Module | `go.mod` | **Minor** — add `conpty` and `golang.org/x/sys` |

**Summary:** 0 interfaces changed, 0 domain changes, 0 transport changes. 1 new application-layer interface (`ShellFinder`). 2 infrastructure files renamed. 3 infrastructure files created. 1 application file created.

---

## 12. Migration Plan

### Phase 1: Extract Shell Discovery

1. Create `internal/application/shell_finder.go` — define `ShellFinder` interface
2. Create `internal/infrastructure/shell_unix.go` with `//go:build !windows`
   - Implement `UnixShellFinder` using `$SHELL` + fallback chain
3. Create `internal/infrastructure/shell_windows.go` with `//go:build windows`
   - Implement `WindowsShellFinder` using `%COMSPEC%` + fallback chain

### Phase 2: Split PTY Files by Platform

4. Rename `internal/infrastructure/pty_manager.go` → `pty_unix.go` (via `git mv`)
   - Add `//go:build !windows` tag
   - Inject `ShellFinder` into `PTYManager` constructor
   - Replace hardcoded `exec.Command("bash")` with `ShellFinder.FindShell()`
5. Rename `internal/infrastructure/pty_instance.go` → keep as `pty_instance.go`
   - Add `//go:build !windows` tag
   - No structural changes — struct fields are already Unix-appropriate

### Phase 3: Implement Windows PTY

6. Create `internal/infrastructure/pty_windows.go` with `//go:build windows`
   - Implement `PTYManager` using `github.com/UserExistsError/conpty`
   - Implement `PTYInstance` wrapping `conpty.ConPty`
   - Inject `ShellFinder` into constructor

### Phase 4: Update Wiring

7. Update `cmd/server/main.go`
   - Instantiate platform-appropriate `ShellFinder`
   - Pass it to `NewPTYManager(shellFinder)`
8. Update `internal/application/session_manager.go`
   - `SessionManager` now holds `ShellFinder` for injection into PTYManager (or PTYManager receives it at construction)

### Phase 5: Update Dependencies

9. Add `github.com/UserExistsError/conpty` and `golang.org/x/sys` to `go.mod`
10. Run `go mod tidy` on Unix to update `go.sum`

### Phase 6: Verify

11. `go build ./...` on Linux — must succeed
12. `go build ./...` on Windows — must succeed
13. `go test ./...` on both platforms
14. Manual test: create session, verify terminal I/O on both platforms

### Phase 7: Update Documentation

15. Update `docs/architecture.md` — document cross-platform support and ShellFinder
16. Update `README.md` — mention Windows support and requirements

---

## 13. Testing Strategy

### Unit Tests (no changes needed)

The existing `session_manager_test.go` uses `mockPTYProvider` and `mockPTY`. These tests already prove the application layer is platform-independent. They pass on all platforms without modification.

### Shell Discovery Tests

Create `internal/infrastructure/shell_unix_test.go` and `shell_windows_test.go`:

- Verify fallback chain when `$SHELL` is unset
- Verify preferred shell is selected when `$SHELL` is set
- Verify error when no shell is found

### Platform-Specific Integration Tests

Create `internal/infrastructure/pty_integration_test.go` (no build tags — runs on all platforms):

- Test that `NewPTYManager()` returns a valid provider
- Test that `Spawn()` creates a working PTY
- Test that writing input produces output
- Test that `Close()` cleans up without leaking

These tests use the `PTYProvider` interface and call the platform-specific implementation automatically.

### CI Matrix

| Platform | Build | Unit Tests | Integration Tests |
|----------|-------|------------|-------------------|
| Linux (`ubuntu-latest`) | `go build ./...` | `go test ./...` | Full PTY integration |
| Windows (`windows-latest`) | `go build ./...` | `go test ./...` | Full PTY integration |
| macOS (`macos-latest`) | `go build ./...` | `go test ./...` | Full PTY integration |

### Manual Verification Checklist

- [ ] Create session on Windows — terminal appears and accepts input
- [ ] Create session on Linux — terminal appears and accepts input
- [ ] Commands execute and output streams to connected clients
- [ ] Host disconnect terminates session on both platforms
- [ ] Multiple viewers join and all receive synchronized output
- [ ] Shell detection falls back correctly when preferred shell is unavailable
- [ ] `go mod tidy` on Unix does not break Windows build
- [ ] Windows build fails gracefully on pre-1809 Windows (clear error message)

---

## 14. Risks and Mitigations

| Risk | Severity | Mitigation |
|------|----------|------------|
| ConPTY requires Windows 10 1809+ | Medium | `conpty.IsConPtyAvailable()` checked at spawn time. Return clear error. |
| `go mod tidy` on Windows removes Unix deps | Medium | CI runs `go mod tidy` on Linux only. Document constraint. |
| `conpty` library becomes unmaintained | Low | Library is ~360 lines over 5 kernel32.dll calls. Simple to fork or inline. |
| Shell availability differs per platform | Low | Dedicated `ShellFinder` with per-platform fallback chains. |
| Test coverage gap on Windows CI | Medium | Windows runner in CI matrix. |
| Process cleanup semantics differ | Low | Both platforms support `cmd.Process.Kill()`. ConPTY close handles cleanup. |

---

## 15. Backward Compatibility

- **API compatibility:** 100%. No HTTP endpoints, WebSocket protocols, or domain models change.
- **Binary compatibility:** Unix binary behavior is identical. Windows binary is new.
- **Configuration compatibility:** No configuration files affected.
- **Dependency compatibility:** `creack/pty` remains for Unix. No existing dependency removed.
- **Test compatibility:** All existing tests pass without modification.

---

## 16. Future Extensibility

### PTY Provider Architecture

The `PTYProvider` interface enables multiple implementations without modifying application or domain code:

```
                        Application Layer
                              │
                              ▼
                      ┌───────────────┐
                      │  PTYProvider  │
                      │  interface    │
                      └───────┬───────┘
                              │
              ┌───────────────┼───────────────────┐
              │               │                   │
      ┌───────▼──────┐ ┌─────▼────────┐ ┌────────▼────────┐
      │  UnixPTY     │ │ WindowsPTY   │ │  Future Providers│
      │  Provider    │ │ Provider     │ │                  │
      │  (creack/pty)│ │ (conpty)     │ │  ├── Docker      │
      └──────────────┘ └──────────────┘ │  ├── SSH         │
                                        │  └── Remote Exec │
                                        └─────────────────┘
```

Each provider implements `PTYProvider` and `PTYInstance`. The application layer consumes only the interface. Adding a new provider (e.g., Docker-based isolation for security) requires:
1. A new file in `internal/infrastructure/`
2. A constructor that returns the concrete type
3. No changes to `session_manager.go`, `pty_streamer.go`, or any domain code

### Terminal Recording

The `PTYInstance.Output()` channel is the natural injection point. A decorator pattern wraps any platform-specific instance:

```
type RecordingPTYInstance struct {
    inner    application.PTYInstance
    recorder io.Writer
}
```

This wraps `Write()`, `Output()`, and `Close()` — recording input and output transparently.

### Platform-Specific Overrides

If a platform later needs special handling beyond what `!windows` provides (e.g., macOS-specific terminal settings), a more specific build tag file overrides the catch-all:

```
//go:build darwin

// Darwin-specific overrides
```

Go resolves build tags from most-specific to least-specific, so `darwin` overrides `!windows`.

---

## 17. Architecture Review Outcome

### Architectural Strengths

1. **Existing dependency inversion is correct.** The `PTYProvider`/`PTYInstance` interfaces in the application layer already provide the exact abstraction needed. No interface changes are required.
2. **Clean layer separation.** Domain, application, and transport layers have zero OS awareness. Platform code is fully confined to infrastructure.
3. **Existing tests validate the abstraction.** `session_manager_test.go` uses mock PTY types, proving the application layer is already platform-independent.
4. **Build tags are the correct tool.** The problem is compile-time, not runtime. Build tags select the right implementation without runtime overhead or conditional logic.

### Remaining Risks

1. **ConPTY version floor.** Windows 10 1809 is the minimum. Older Windows versions (7, 8, 8.1, pre-1809 10) are unsupported. Mitigated by `conpty.IsConPtyAvailable()` and clear error messages.
2. **`go mod tidy` asymmetry.** Running `go mod tidy` on Windows would remove Unix-only dependencies from `go.mod`. Mitigated by CI policy (only run `go mod tidy` on Unix).
3. **Third-party dependency.** `github.com/UserExistsError/conpty` is a small, focused library. Risk is low but non-zero. The library's simplicity (~360 lines) makes forking feasible if needed.

### Open Questions

1. **ShellFinder injection scope.** Should `ShellFinder` be injected into `PTYManager` at construction, or should it be a package-level function? Construction injection is more testable but adds a parameter to `NewPTYManager()`. Package-level functions are simpler but harder to mock.
2. **ConPTY pipe semantics.** The `conpty.ConPty` type exposes `Read()` and `Write()` directly. The readLoop goroutine pattern from the Unix implementation needs adaptation to ensure the output channel semantics match. This is an implementation detail, not an architectural concern.
3. **TERM environment variable.** The current Unix implementation sets `TERM=xterm`. Windows ConPTY does not use `TERM`. The `ShellFinder` or `PTYManager` should handle this environment variable per-platform.

### Final Recommendation

**Proceed with the design as documented.** The architecture is sound. The changes are confined to the infrastructure layer with zero impact on business logic. The existing interface boundaries make this a clean, low-risk cross-platform extension. The `ShellFinder` extraction improves separation of concerns beyond the original scope. The migration strategy preserves Git history through file renaming rather than deletion.
