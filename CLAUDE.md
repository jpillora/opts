# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project

`github.com/jpillora/opts` — A Go library for building CLI applications from struct definitions using struct tags. Flags, positional arguments, subcommands, help text, env vars, JSON config, and shell completion are all derived from annotated Go structs.

## Commands

```bash
# Build
go build -v -o /dev/null

# Test
go test -v ./...

# Run a single test
go test -v -run TestName ./...
```

## Architecture

The library exposes two interfaces (`Opts` for configuration, `ParsedOpts` for results) via `opts.go`, with entry points `New(config)` and `Parse(config)`.

Internally, a tree of `node` structs represents the command hierarchy:

- **`node.go`** — Core `node` struct (command tree node, embeds `item`)
- **`node_build.go`** — Builder pattern methods implementing `Opts` (Name, Version, ConfigPath, UseEnv, etc.)
- **`node_parse.go`** — Main parsing engine: reflects over struct fields, builds `flag.FlagSet`, applies env/config defaults, validates, recurses into subcommands
- **`node_commands.go`** — Subcommand registration, `Run()` dispatch, `Selected()` traversal
- **`node_help.go`** — Help text generation via Go templates with customizable sections and padding
- **`node_complete.go`** — Shell completion via `github.com/posener/complete`

Field-level logic lives in:

- **`item.go`** — Represents a single struct field as a flag/arg/command. Handles type detection (primitives, `time.Duration`, `encoding.TextUnmarshaler`, `flag.Value`, custom `Setter`, and slices). Implements `flag.Value` interface.
- **`strings.go`** — `camel2dash` (field names → flag names), `camel2const` (env var names), text wrapping, singularization

## Key Design Decisions

- Two error types: `authorError` (programmer mistakes, panics) vs `exitError` (user input errors, prints and exits)
- Struct tag format: `` `opts:"key=value,key=value"` `` with keys: `name`, `help`, `mode` (flag/arg/embedded/cmd/cmdname), `short`, `group` (for flags/embedded/cmd), `env`, `min`, `max`, `-` (ignore)
- Field names are auto-converted from CamelCase to kebab-case for flags and CONSTANT_CASE for env vars
- Slices support repeated flags and positional args with min/max validation
- The `Setter` interface (`Set(string) error`) allows custom flag types
