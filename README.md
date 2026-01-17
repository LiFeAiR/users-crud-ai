# HTTP Server with Database Connection and REST API

This is a REST API server implementation with database connection capabilities.

## Project Structure

```
.
├── api/
│   └── openapi/
│       └── api.yaml         # OpenAPI specification
├── build/                   # Build and deployment configurations
│   ├── .env                 # Environment variables for Docker Compose
│   └── docker-compose.yml   # Docker Compose configuration for PostgreSQL database
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
│   │   └── user_repository.go    # User repository implementation
│   └── server/              # HTTP server implementation
└── go.mod                   # Go modules file
```

## Build and Deployment

The `build` directory contains configuration files for building and deploying the application with Docker Compose and PostgreSQL:

- `.env`: Environment variables for database and pgAdmin configuration
- `docker-compose.yml`: Docker Compose configuration that sets up PostgreSQL database service

These files allow you to easily spin up the entire application environment with a single command using Docker Compose.

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

Note: In production environments, database connection management should be handled more carefully
to prevent premature closure of connections. Currently, the connection is closed after
the Start method completes, which is not ideal for long-running servers.
```

## Repository Pattern Implementation

The project implements repository pattern in `internal/repository/`:

- `UserRepository` interface defines methods for user operations (defined in `interface.go`)
- `userRepository` struct implements the interface (defined in `user_repository.go`)
- `OrganizationRepository` interface defines methods for organization operations (defined in `interface.go`)
- `organizationRepository` struct implements the interface (defined in `organization_repository.go`)
- Both repositories use the same database connection from `db.go`
- Provides methods for CRUD operations on users and organizations
- Includes database initialization for both tables

## Handler Integration

HTTP handlers in `internal/handlers/` now use the repository pattern through a base handler:

- Base handler receives repository instances for users and organizations
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

## Фичи

1. Связь пользователя и организации
  - нельзя удалить организацию, если она связана с пользователем
  - нельзя присвоить пользоватлю id несуществующей организации
  - может быть пользователь без организации
  - пользователь может состоять в нескольких организациях
2. ✅ Сохранение паролей в безопасном виде
3. ✅ Авторизация пользователя по паролю (id/email+password -> jwt)
4. Пользователь может редактировать только свои данные
5. ✅ Права
6. ✅ Роли (одна роль может содержать несколько прав)
7. ✅ Пользователь может получить определенный набор прав (по одному или пачкой получив роль)
8. Права отображаются в jwt
9. Пользователь с правом "Администратор" может редактировать всех пользователей
10. ✅ Организация может получить определенный набор прав (по одному или пачкой получив роль)
11. Пользователь с правом "Администратор" внутри организации, может редактировать всех пользователей в ней (кроме пароля)
12. ✅ Тарифы
13. ✅ Тарифы могут предоставлять ~~права~~ или роли с правами
14. Тарифы связаны с организациями и пользователями
15. Можно иметь несколько тарифов
16. Тарифы предоставляется пользователю или организации на ограниченный срок от, до
17. Тарифы могут иметь атрибуты (например цифровые)
18. Администратор может удалять пользователей из организации
19. Нельзя добавить нового пользователя в организацию если превышен лимит (атрибут из тарифа)
20. Пользователь добавляется в организацию только после подтверждения Администратор организации
21. Администратор организации может комбинировать права организации в роли (сложно)
    - может быть ограничиться перераспределением ролей