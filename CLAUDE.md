# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Go SDK for OnlyOffice Document Server integration. Provides document editing configuration, format conversion, JWT authentication, version history management, and callback processing.

## Build & Test Commands

```bash
# Run all tests
go test ./...

# Run a specific test
go test -run TestBuildEditorConfig ./...

# Run tests with verbose output
go test -v ./...

# Build the example
go build -o example ./examples

# Run the example (starts callback server on :8080)
go run ./examples/main.go

# Check for issues
go vet ./...
```

## Architecture

**Core Client Pattern**: All functionality flows through `Client` struct initialized with `Config`. The client handles JWT signing/verification and HTTP operations.

```
client.go      → Client struct, config, JWT operations, BuildEditorConfig()
conversion.go  → ConvertDocument(), format utilities (CanConvert, GetInternalExtension)
callback.go    → ParseCallback(), ValidateCallback(), GetDownloadURL()
history.go     → CreateHistory(), GetHistory(), CountVersion() - file-based version storage
models/        → Data structures for OnlyOffice API (Config, Callback, Document, etc.)
```

**Key Flows**:
- Editor config: `NewClient()` → `BuildEditorConfig()` returns `models.Config` with JWT token
- Document conversion: `ConvertDocument()` POSTs to `/ConvertService.ashx`
- Callback handling: `ParseCallback()` → `ValidateCallback()` → `GetDownloadURL()`
- History: Stored in `{storagePath}/.history/{key}/changes.json`

**JWT Integration**: When `Config.JWTEnabled=true`, tokens are automatically created for editor configs and conversion requests, and validated on callbacks via Authorization header.

## Callback Status Codes

- `2` - Document ready/edited (valid)
- `6` - Must save (valid)
- `7` - Document corrupted (error)
