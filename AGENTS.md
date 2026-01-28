# AGENTS.md

This document provides guidance for agentic coding agents operating in this repository.

## Build/Lint/Test Commands

### Build Commands
```bash
# Build the main application
go build -o bin/app cmd/app/main.go

# Build the CLI application
go build -o bin/cli cmd/cli/main.go

# Generate protobuf files
make grpc-gen
```

### Lint Commands
```bash
# Install golangci-lint if not already installed
curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.54.2

# Run linter
golangci-lint run

# Run linter with specific configuration
golangci-lint run --config .golangci.yml
```

### Test Commands
```bash
# Run all tests
go test ./...

# Run tests with verbose output
go test -v ./...

# Run a specific test file
go test -v ./internal/handlers/login_test.go

# Run a specific test function
go test -v ./internal/handlers -run TestBaseHandler_Login

# Run tests with coverage
go test -cover ./...

# Run tests with race detector
go test -race ./...
```

## Code Style Guidelines

### Imports
- Group imports in the following order:
  1. Standard library packages
  2. External packages
  3. Internal packages
- Sort imports alphabetically within each group
- Use explicit package aliases for clarity when needed
- Example:
  ```go
  import (
      "fmt"
      "net/http"
      
      "github.com/gin-gonic/gin"
      "github.com/stretchr/testify/assert"
      
      "github.com/LiFeAiR/crud-ai/internal/models"
      "github.com/LiFeAiR/crud-ai/internal/repository"
  )
  ```

### Formatting
- Follow Go standard formatting (gofmt)
- Use 4-space indentation
- Limit lines to 120 characters
- Place opening braces on the same line as the declaration
- Use blank lines to separate logical sections of code

### Types
- Use meaningful type names that clearly indicate their purpose
- Prefer interfaces over concrete types when possible
- Define custom error types for specific error cases
- Use pointer receivers for large structs or when mutation is needed

### Naming Conventions
- Use camelCase for variables and functions
- Use PascalCase for exported types and methods
- Use lowercase for package names
- Use descriptive names that clearly indicate the purpose (avoid abbreviations)
- Use `ctx` for context parameters
- Use `err` for error parameters
- Use `ID` for identifiers
- Use `URL` for URLs

### Error Handling
- Always check and handle errors appropriately
- Return errors early in functions
- Use wrapped errors with fmt.Errorf for context
- Implement custom error types for domain-specific errors
- Don't ignore errors
- Use sentinel errors for well-known error cases

### Logging
- Use structured logging where possible
- Log at appropriate levels (debug, info, warn, error)
- Include contextual information in logs
- Avoid logging sensitive information

### Testing
- Write unit tests for all functions
- Use table-driven tests for multiple test cases
- Mock external dependencies using interfaces
- Test error cases and edge cases
- Use testify assertions for cleaner test code
- Test both happy path and failure scenarios
- Include integration tests for API endpoints
- Run tests regularly and maintain high coverage

### Documentation
- Document exported functions, types, and methods with comments
- Use godoc format for documentation
- Add inline comments for complex logic
- Include examples in documentation where appropriate

### Security
- Never expose sensitive information in logs
- Validate all inputs
- Use parameterized queries to prevent SQL injection
- Implement proper authentication and authorization
- Handle passwords securely (hash with salt)
- Use HTTPS in production environments

## Cursor/Copilot Rules

No specific Cursor or Copilot rules found in the repository.