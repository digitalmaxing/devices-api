# Devices API

A RESTful API for managing device resources, built with Go.

## Features

- Full CRUD operations for devices
- Filter devices by brand and state
- Domain validations enforced (e.g., cannot update name/brand for in-use devices, cannot delete in-use devices)
- Persistent storage with PostgreSQL
- API documentation with Swagger
- Containerized with Docker
- Comprehensive test coverage

## Tech Stack

- **Go 1.23+**
- **Gin** - Web framework
- **GORM** - ORM for PostgreSQL
- **PostgreSQL** - Database
- **Docker & Docker Compose** - Containerization
- **Swagger** - API documentation

## Project Structure

```
devices-api/
├── cmd/
│   ├── api/
│       └── main.go              # Application entrypoint
├── internal/
│   ├── handlers/                # HTTP handlers
│   │   └── device_handler.go
│   ├── models/                  # Domain models
│   │   └── device.go
│   ├── repository/              # Data access layer
│   │   └── device_repository.go
│   ├── service/                 # Business logic layer
│   │   └── device_service.go
│   └── config/                  # Configuration
│       └── config.go
├── Dockerfile
├── docker-compose.yml
├── go.mod
├── go.sum (generated on go mod tidy)
├── .env.example
└── README.md
```

## Getting Started

### Prerequisites

- Go 1.23+
- Docker & Docker Compose (recommended)
- PostgreSQL (optional for local non-Docker dev)
- `swag` CLI for regenerating Swagger docs (optional): `go install github.com/swaggo/swag/cmd/swag@latest`

### Running with Docker (Recommended - Zero Setup)

1. Clone the repository:
   ```bash
   git clone https://github.com/digitalmaxing/devices-api.git
   cd devices-api
   cp .env.example .env   # optional, defaults work
   ```

2. Start everything (builds image, starts Postgres + API with healthchecks):
   ```bash
   docker-compose up --build -d
   ```

3. API available at http://localhost:8080
   - Swagger UI: http://localhost:8080/swagger/index.html (interactive docs & testing)
   - Health: http://localhost:8080/health

4. Stop: `docker-compose down`

### Running Locally (Go + Postgres required)

1. Start Postgres (or use Docker for DB only):
   ```bash
   docker run -d --name pg -e POSTGRES_PASSWORD=postgres -e POSTGRES_DB=devices -p 5432:5432 postgres:16
   ```

2. Set env and run:
   ```bash
   export DB_HOST=localhost DB_USER=postgres DB_PASSWORD=postgres DB_NAME=devices
   go run ./cmd/api
   ```

### Generate / Update Swagger Docs

After code changes:
```bash
swag init --parseDependency --parseInternal -g cmd/api/main.go
```
Then rebuild/restart to see updated docs at /swagger.

### Testing

```bash
go test ./... -v -cover
```
Current coverage focuses on service layer validations and core flows (expandable with more handler/integration tests).

### Example API Usage (via curl or Swagger)

Create device:
```bash
curl -X POST http://localhost:8080/devices -H "Content-Type: application/json" -d '{"name":"Pixel 8","brand":"Google","state":"available"}'
```

List by brand:
```bash
curl "http://localhost:8080/devices?brand=Apple"
```

## Domain Validations

- Creation time is immutable
- Name and Brand cannot be changed if device state is "in-use"
- Devices in "in-use" state cannot be deleted

## Future Improvements

- Add authentication & authorization (JWT)
- Rate limiting and request logging middleware
- Structured logging with zap or zerolog
- CI/CD pipeline with GitHub Actions
- Database migrations with golang-migrate (instead of AutoMigrate)
- Metrics with Prometheus
- Better error handling with custom error types and HTTP status mapping
- Integration tests with testcontainers

## Notes / Known Limitations (per tips)

- Uses GORM AutoMigrate for simplicity (production would use proper migrations)
- No in-memory DB fallback (always requires Postgres per requirements)
- Partial updates implemented with map[string]interface{} for flexibility
- Error handling is basic string-based (could improve with sentinel errors)
- No auth/rate limiting yet (production readiness partial)
- go.sum will be generated on first `go mod tidy` after clone

## License

MIT