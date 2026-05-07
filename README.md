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
├── go.sum
└── README.md
```

## Getting Started

### Prerequisites

- Go 1.23+
- Docker & Docker Compose
- PostgreSQL (if running locally without Docker)

### Running with Docker (Recommended)

1. Clone the repository:
   ```bash
   git clone https://github.com/digitalmaxing/devices-api.git
   cd devices-api
   ```

2. Start the services:
   ```bash
   docker-compose up --build
   ```

3. The API will be available at `http://localhost:8080`

4. Access Swagger UI at `http://localhost:8080/swagger/index.html`

### Environment Variables

See `.env.example` (create from it or use docker-compose defaults).

## API Endpoints

See Swagger documentation for full details and examples.

### Devices

- `POST /devices` - Create a new device
- `GET /devices` - List all devices (supports ?brand=...&state=...)
- `GET /devices/:id` - Get device by ID
- `PATCH /devices/:id` - Partially update device
- `DELETE /devices/:id` - Delete device

## Domain Validations

- Creation time is immutable
- Name and Brand cannot be changed if device state is "in-use"
- Devices in "in-use" state cannot be deleted

## Testing

```bash
go test ./... -v
```

## Future Improvements

- Add authentication & authorization (JWT)
- Rate limiting and request logging middleware
- Health check endpoints (/health, /ready)
- Structured logging with zap or zerolog
- CI/CD pipeline with GitHub Actions
- Database migrations with golang-migrate
- Metrics with Prometheus
- Better error handling with custom error types

## Notes / Known Limitations

- Uses GORM AutoMigrate for simplicity (production would use proper migrations)
- In-memory DB option not used; always Postgres
- Partial updates use map for flexibility

## License

MIT