# Project dcd-be

DCD Backend service built with Go, Gin, GORM, PostgreSQL, and Docker Compose. It features live-reloading (via Air).

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes.

### Prerequisites

- Go (1.21+)
- Docker and Docker Compose
- Air (optional, for live reload: `go install github.com/air-verse/air@latest`)

### Configuration

Before running the application, copy the example environment configuration:
```bash
cp .env.example .env
```
Ensure the values match your local development environment settings.

## Running the Project with Task

This project uses [Task](https://taskfile.dev) as a task runner. Install it on your machine (e.g., `go install github.com/go-task/task/v3/cmd/task@latest` or `brew install go-task`).

To list all available tasks, run:
```bash
task --list
```

### Build & Test

Run the default task (lists all available tasks):
```bash
task
```

### Build the application binary:
```bash
task build
```

Run the Go test suite:
```bash
task test
```

### Development & Docker Services

Spin up the local docker container services (PostgreSQL) in the background:
```bash
task docker-run
```

Stop the docker containers:
```bash
task docker-down
```

Run the application locally:
```bash
task run
```

Live reload the application (uses Air):
```bash
task watch
```

Run database integration tests:
```bash
task itest
```

Clean up build artifacts:
```bash
task clean
```
