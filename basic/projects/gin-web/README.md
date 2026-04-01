# Web Application Template with Go and Gin

A production-ready web application project template using Go and the Gin web framework.

## Features

- **Gin v1.12**: High-performance HTTP web framework
- **Middleware Stack**: RequestID, CORS, Error handling, Recovery, Logger
- **HTML Templates**: Go's html/template package for server-side rendering
- **Static Files**: CSS and JavaScript asset serving
- **Docker Support**: Multi-stage build for small images
- **CI/CD**: GitHub Actions with test, lint, and docker build

## Project Structure

```
basic/projects/gin-web/
├── cmd/server/main.go           # Application entry point
├── internal/
│   ├── handler/home.go        # HTTP handlers
│   └── middleware/middleware.go # Custom middleware
├── internal/views/             # HTML templates
├── static/                    # CSS and JS assets
│   ├── css/style.css
│   └── js/app.js
├── Dockerfile
├── docker-compose.yml
├── .github/workflows/ci.yml
├── go.mod
└── Makefile
```

## Prerequisites

- Go 1.23 or later
- Docker (optional)
- Make

## Quick Start

### Local Development

```bash
# Run the server
cd basic/projects/gin-web
go run ./cmd/server

# Or use make
make run
```

The server starts on http://localhost:8080

### Docker

```bash
# Build and run with docker-compose
docker-compose up --build

# Or build and run manually
docker build -t gin-web .
docker run -p 8080:8080 gin-web
```

## Endpoints

| Method | Endpoint           | Description              |
|--------|-------------------|--------------------------|
| GET    | /                 | Home page                |
| GET    | /about            | About page               |
| GET    | /api/v1/health    | Health check (JSON)      |
| GET    | /static/*         | Static files             |

## Example Requests

```bash
# Home page
curl http://localhost:8080/

# About page
curl http://localhost:8080/about

# Health check
curl http://localhost:8080/api/v1/health
# Response: {"status":"ok","version":"1.0.0"}

# Static CSS
curl http://localhost:8080/static/css/style.css

# Static JS
curl http://localhost:8080/static/js/app.js
```

## Middleware

The application includes the following middleware:

1. **Recovery** - Catches panics and returns 500 error
2. **Logger** - Logs HTTP requests
3. **RequestID** - Generates unique ID for each request (X-Request-ID header)
4. **CORS** - Handles Cross-Origin Resource Sharing
5. **ErrorHandler** - Catches panics and returns JSON error response

## Testing

```bash
# Run all tests
go test -v ./...

# Run with coverage
go test -v -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## Make Commands

```bash
make run       # Run the server
make build     # Build the binary
make test      # Run tests
make lint      # Run golangci-lint
make docker    # Build Docker image
make docker-up # Run with docker-compose
make clean     # Clean up
```

## License

MIT
