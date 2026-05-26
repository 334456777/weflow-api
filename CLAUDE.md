# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build & Install

```bash
make build          # go build with ldflags (version/commit/date)
make install        # build + sudo install to /usr/local/bin/weflow
make test           # go test ./...
make fmt            # gofmt -s -w .
make vet            # go vet ./...
make tidy           # go mod tidy
```

## Release Workflow

1. Commit changes
2. `git tag v0.1.x`
3. `make build && make install`
4. `git push origin main && git push origin v0.1.x`

Always include `make install` after building.

## Architecture

- `cmd/` — Cobra CLI commands; `root.go` has `newRenderer()` dispatcher
- `internal/client/` — HTTP client for WeFlow local API (127.0.0.1:5031)
- `internal/output/` — Renderers: table (default), JSON (`--json`), text (`--text`)
- `internal/config/` — Viper config; flag > env > file > defaults
- `internal/version/` — Injected via ldflags at build time

## Conventions

- Output is in Chinese (commit messages, CLI help, stderr progress)
- `internal/output/format.go` — `FormatContent()` handles WeChat XML→text (share cards, revoke, etc.)
- `internal/output/text.go` — System messages (sender contains `@chatroom`) skip sender display
- `cmd/messages.go` — First request returning empty auto-retries once (API warmup)
