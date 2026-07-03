# AGENTS.md

## Overview

HttpForge is a lightweight, config-driven HTTP API testing toolkit written in Go. It allows developers to build, execute, and validate HTTP requests from YAML/JSON configuration files.

## Architecture

### Core Modules

1. **models/** - Data structures (Request, Response, Assertion, etc.)
2. **parser/** - YAML/JSON config file loading
3. **template/** - Variable substitution engine
4. **engine/** - HTTP request execution
5. **assert/** - Response validation engine
6. **output/** - Response formatting (text, JSON, CSV)
7. **export/** - Code generation (curl, Python, JavaScript)
8. **cmd/** - CLI commands (Cobra)

### Key Design Decisions

- **Cobra CLI** - Industry standard for Go CLI applications
- **Template variables** - `${var}` syntax for dynamic values
- **Pluggable assertions** - Extensible validation system
- **Multiple output formats** - Text for humans, JSON for scripts

## Development

### Building

```bash
go build -o httpforge .
```

### Testing

```bash
# Run all tests
go test ./...

# Run with verbose output
go test -v ./...

# Run specific package tests
go test ./parser/...
go test ./assert/...
```

### Adding New Assertions

1. Add field handler in `assert/assert.go`
2. Add operation handler in `evaluateOp()`
3. Add tests in `assert/assert_test.go`

### Adding New Export Formats

1. Add format constant in `models/models.go`
2. Add export function in `export/export.go`
3. Update `Export()` function
4. Add tests in `export/export_test.go`

## Dependencies

- `github.com/spf13/cobra` - CLI framework
- `github.com/fatih/color` - Terminal colors
- `gopkg.in/yaml.v3` - YAML parsing
- `github.com/tidwall/gjson` - JSON path queries

## Testing Guidelines

- Test all public functions
- Test edge cases (empty input, invalid input)
- Test error paths
- Use table-driven tests where appropriate
- Aim for >80% coverage
