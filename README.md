# HttpForge

A lightweight, config-driven HTTP API testing toolkit for developers and DevOps engineers.

## Features

- **YAML/JSON Request Configs** - Define requests in human-readable files
- **Template Variables** - Use `${variable}` syntax for dynamic values
- **Response Assertions** - Validate responses with built-in assertion engine
- **Code Export** - Generate curl, Python, or JavaScript code from requests
- **Environment Management** - Separate configs for dev/staging/prod
- **Multiple Output Formats** - Text, JSON, and CSV output
- **CLI-First** - Fast, single-binary, no GUI required

## Quick Start

### Install

```bash
go install github.com/EdgarOrtegaRamirez/httpforge@latest
```

### Run a Request

```bash
# Simple GET request
httpforge run https://httpbin.org/get

# From a YAML file
httpforge run request.yaml

# With verbose output
httpforge run request.yaml -v
```

### Create a Request Config

```yaml
# request.yaml
request:
  name: Get User
  method: GET
  url: https://api.example.com/users/123
  headers:
    Authorization: Bearer ${API_TOKEN}
    Accept: application/json
  assertions:
    - field: status_code
      op: eq
      value: "200"
    - field: json_path
      op: eq
      value: "user.name"
```

## Commands

| Command | Description |
|---------|-------------|
| `httpforge run <file|url>` | Execute an HTTP request |
| `httpforge export <file>` | Export request as code |
| `httpforge info <file>` | Display request configuration |
| `httpforge validate <file>` | Validate request file format |
| `httpforge version` | Show version information |

## Request File Format

### YAML Example

```yaml
request:
  name: Create User
  method: POST
  url: https://api.example.com/users
  headers:
    Content-Type: application/json
    Authorization: Bearer ${API_TOKEN}
  body: |
    {
      "name": "John Doe",
      "email": "john@example.com"
    }
  body_type: json
  query_params:
    send_welcome: "true"
  assertions:
    - field: status_code
      op: eq
      value: "201"
    - field: json_path
      op: exists
      value: "user.id"
    - field: json_path
      op: type
      value: "number"
```

### JSON Example

```json
{
  "request": {
    "name": "Get Users",
    "method": "GET",
    "url": "https://api.example.com/users",
    "headers": {
      "Accept": "application/json"
    },
    "query_params": {
      "limit": "10",
      "offset": "0"
    }
  }
}
```

## Template Variables

HttpForge supports variable substitution using `${variable_name}` syntax:

```yaml
request:
  url: https://${HOST}:${PORT}/api/${RESOURCE}
  headers:
    Authorization: Bearer ${API_TOKEN}
```

Variables are resolved from (in order of priority):
1. Command-line `--var` flags
2. Environment file (`--env`)
3. `HTTPFORGE_*` environment variables
4. System environment variables

## Assertions

Validate responses with built-in assertions:

```yaml
assertions:
  # Status code
  - field: status_code
    op: eq
    value: "200"

  # JSON path
  - field: json_path
    op: eq
    value: "user.name"

  # Response header
  - field: header
    op: contains
    value: "Content-Type"

  # Response body length
  - field: body.length
    op: gt
    value: "0"

  # Response timing
  - field: timing
    op: lt
    value: "5s"
```

### Supported Operations

| Operation | Description |
|-----------|-------------|
| `eq`, `==` | Equals |
| `ne`, `!=` | Not equals |
| `gt`, `>` | Greater than |
| `gte`, `>=` | Greater than or equal |
| `lt`, `<` | Less than |
| `lte`, `<=` | Less than or equal |
| `contains` | Contains substring |
| `starts_with` | Starts with |
| `ends_with` | Ends with |
| `exists` | Value exists |
| `empty` | Value is empty |
| `type` | Value type check |

## Code Export

Export requests as executable code:

```bash
# Export as curl
httpforge export request.yaml -f curl

# Export as Python
httpforge export request.yaml -f python

# Export as JavaScript
httpforge export request.yaml -f javascript
```

## Environment Files

Create separate environment configs for different stages:

```yaml
# .env.dev.yaml
name: development
variables:
  HOST: localhost
  PORT: "8080"
  API_TOKEN: dev-token-123
```

```bash
# Run with dev environment
httpforge run request.yaml --env .env.dev.yaml

# Override variables
httpforge run request.yaml --var "API_TOKEN=custom-token"
```

## Output Formats

```bash
# Text output (default)
httpforge run request.yaml -o text

# JSON output
httpforge run request.yaml -o json

# CSV output
httpforge run request.yaml -o csv
```

## Use Cases

- **API Testing** - Test REST APIs during development
- **CI/CD** - Validate API endpoints in pipelines
- **Documentation** - Generate code examples from request configs
- **Load Testing** - Create request sequences for testing
- **Debugging** - Inspect API responses with verbose output

## Architecture

```
httpforge/
├── models/          # Data structures
├── parser/          # YAML/JSON config loading
├── template/        # Variable substitution
├── engine/          # HTTP request execution
├── assert/          # Response validation
├── output/          # Response formatting
├── export/          # Code generation
├── cmd/             # CLI commands
└── tests/           # Test data
```

## Security

- Environment variables for secrets (never commit tokens)
- Template variables prevent hardcoded credentials
- Input validation on all request parameters
- No remote code execution
- Local-only operation (no telemetry)

## License

MIT License
