# Developer Contribution & Onboarding Guide

Welcome to the **dcd-be** backend service repository! This guide will help you get set up locally, understand the project architecture, and follow the standard workflows for adding features and collaborating.

---

## API Documentation & Playground
We use **Apidog** for interactive API documentation, schema references, and client-side mocking:
* **Interactive API Docs**: `https://74upocjbme.apidog.io/`

---

## 1. Quick Start Local Setup

### Prerequisites
* Go 1.26.2+
* Docker & Docker Compose
* [Task](https://taskfile.dev) (Task runner tool: `go install github.com/go-task/task/v3/cmd/task@latest`)
* [Air](https://github.com/air-verse/air) (Live reloading tool: `go install github.com/air-verse/air@latest`)

### Steps
1. **Clone and Configure**
   Copy the example environment configuration:
   ```bash
   cp .env.example .env
   ```
   *By default, the `.env.example` points to `localhost:5432` for running the database inside Docker.*

2. **Start Local Docker Services**
   Spin up the local PostgreSQL database service in the background:
   ```bash
   task docker-run
   ```

3. **Run Database Migrations & Seeds**
   Run the database schema migrations and seed the initial service package listings:
   ```bash
   task migrate
   ```

4. **Start the API Server (with Live Reload)**
   Start the application locally with hot-reloading:
   ```bash
   task watch
   ```
   The API will be available at `http://localhost:8080`.

---

## 2. Project Architecture & Directory Layout

The project follows a standard Go directory structure designed for layered architecture, clean interfaces, and clean separation of concerns.

```text
├── cmd/
│   ├── api/          # Application entrypoint (Server setup, DB connection pool, structured logging setup)
│   └── migrate/      # Database migrations and seeding runner command
├── internal/
│   ├── config/       # Environment parsing (.env file loading, Config struct schema)
│   ├── database/     # DB initialization, connection pooling settings, GORM model definitions, and migrations
│   ├── server/       # HTTP server, routing, custom middlewares (Request ID, custom Logger), and route handlers
│   └── service/      # Core business logic services and DTOs (interfaces and implementations)
├── README.md         # General deployment and starting guide
└── Taskfile.yml      # Shortcuts for building, running, and testing tasks
```

### Dependency Injection Flow
Dependencies flow inward. Handlers depend on Services, and Services depend on the Database interface:
`cmd/api/main.go` ➔ `server.NewServer(cfg, db)` ➔ `service.NewService(db)` ➔ `database.New(cfg)`

---

## 3. Common Developer Workflows

### A. Adding or Modifying a Database Schema
Our database schema is updated using a programmatic, versioned migration runner located in [internal/database/migrations.go](file:///e:/Web%20projects/dcd-be/internal/database/migrations.go).

1. **Define the GORM Model**:
   Update or add your struct inside [internal/database/models.go](file:///e:/Web%20projects/dcd-be/internal/database/models.go).
2. **Add a Migration Step**:
   Open [internal/database/migrations.go](file:///e:/Web%20projects/dcd-be/internal/database/migrations.go), write a migration function, and append it to the `migrationList` with a unique version key:
   ```go
   {
       id: "202608280002_add_field_to_reservation",
       run: func(db *gorm.DB) error {
           return db.Migrator().AddColumn(&Reservation{}, "Status")
       },
   },
   ```
3. **Execute**:
   Run `task migrate` locally to test your migration. Applied migrations are automatically recorded in the `migrations` table to ensure they only run once.

---

### B. Adding a New API Endpoint
1. **Define the Handler**:
   Create or edit the handler function inside `internal/server/`. Bind JSON inputs into request DTOs using Gin bindings:
   ```go
   type CreateItemDTO struct {
       Name string `json:"name" binding:"required"`
   }
   ```
2. **Implement Business Logic**:
   Write the interface and service implementation inside `internal/service/` that interacts with the database.
3. **Wire up Dependency Injection**:
   Configure the handlers inside `internal/server/server.go` to inject the service packages, and register the route inside `internal/server/routes.go`.
4. **Error Handling**:
   Use `NewAppError(httpStatus, message, err)` to propagate errors up to the centralized error middleware.

---

## 4. Observability & Logging

### Structured Logging (`log/slog`)
We use Go's standard structured logging library:
* **Local Development (`APP_ENV=local`)**: Logs are printed as readable text-based key-values.
* **Production/Cloud (`APP_ENV=production`)**: Logs are output in raw JSON format for monitoring integrations (such as Datadog or CloudWatch).

### Correlation Tracking
Every request gets assigned a unique correlation ID via `X-Request-ID`.
* The ID is injected into headers, HTTP context, and all custom log output for request tracing.
* Use `slog.InfoContext(c.Request.Context(), ...)` to ensure the ID is printed alongside contextual attributes.

---

## 5. Testing Guidelines

Tests are split into two suites:

### Unit & Handler Tests
Test individual router handlers, configurations, and services using mock providers:
```bash
task test
```

### Integration Tests
Verify database connections, complex GORM queries, transactions, and repository layers against an actual database using **Testcontainers**:
```bash
task itest
```
*Note: Make sure Docker is running on your machine prior to starting integration tests.*
