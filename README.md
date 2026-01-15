# HTTP Server with Database Connection and REST API

This is a REST API server implementation with database connection capabilities.

## Project Structure

```
.
├── api/
│   └── openapi/
│       └── api.yaml         # OpenAPI specification
├── cmd/
│   ├── app/
│   │   └── main.go          # Main entry point
│   └── cli/
│       └── main.go          # CLI entry point
├── internal/
│   ├── handlers/            # HTTP handlers
│   ├── models/              # Data models
│   │   ├── organization.go  # Organization data model
│   │   └── user.go          # User data model
│   ├── repository/          # Database connection and repositories
│   │   ├── db.go            # Database connection class
│   │   ├── interface.go     # Repository interface
│   │   ├── organization_repository.go  # Organization repository implementation
│   │   └── repository.go    # User repository implementation
│   └── server/              # HTTP server implementation
└── go.mod                   # Go modules file
```

## Database Connection Class

The database connection class is located in `internal/repository/db.go`. It provides:

- Connection to PostgreSQL database using github.com/jackc/pgx/v5
- Connection validation
- Methods for getting and closing database connection

### Usage Example

```go
// Create new database connection
db, err := repository.NewDB("host=localhost port=5432 user=postgres password=postgres dbname=testdb sslmode=disable")
if err != nil {
    log.Fatal("Failed to connect to database:", err)
}
defer db.Close()

// Get database connection
connection := db.GetConnection()

// The server manages the database connection lifecycle
// and passes it to handlers through the base handler
```

## Repository Pattern Implementation

The project implements repository pattern in `internal/repository/repository.go`:

- `UserRepository` interface defines methods for user operations (defined in `interface.go`)
- `userRepository` struct implements the interface
- Provides methods for CRUD operations on users
- Includes database initialization

## Handler Integration

HTTP handlers in `internal/handlers/` now use the repository pattern through a base handler:

- Base handler handles all user operations (create, update, get, delete)
- Handlers receive the server instance to access database connection
- The actual repository interaction is delegated to the base handler
- All handlers are registered via the server in `internal/server/server.go`

## API Endpoints


The server now supports the following endpoints:

`GET /` - Root endpoint

### Users

- `GET /api/users` - Get list of users (requires limit and offset query parameters)
- `GET /api/user` - Get user by ID (requires id query parameter)
- `POST /api/user` - Create a new user
- `PUT /api/user` - Update an existing user
- `DELETE /api/user` - Delete a user

### Organizations

- `GET /api/organizations` - Get list of organizations (requires limit and offset query parameters)
- `GET /api/organization` - Get organization by ID (requires id query parameter)
- `POST /api/organization` - Create a new organization
- `PUT /api/organization` - Update an existing organization
- `DELETE /api/organization` - Delete an organization
## Getting Started

1. Make sure you have Go installed
2. Install dependencies: `go mod tidy`
3. Run the server: `go run cmd/app/main.go`
4. Run the cli: `go run cmd/cli/main.go`