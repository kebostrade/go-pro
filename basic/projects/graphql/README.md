# GraphQL API Template

A production-ready GraphQL API template using Go, gqlgen, chi router, JWT auth, and WebSocket subscriptions.

## Features

- **Schema-first GraphQL** with gqlgen v0.17+
- **chi v5 router** for HTTP handling
- **JWT authentication** middleware
- **WebSocket subscriptions** for real-time updates
- **Data loaders** for N+1 query prevention
- **In-memory database** for quick prototyping
- **Docker** support for containerized deployment
- **CI/CD** with GitHub Actions

## Project Structure

```
basic/projects/graphql/
├── cmd/server/main.go          # Application entry point
├── internal/
│   ├── graph/
│   │   ├── schema.graphqls    # GraphQL schema
│   │   ├── resolver.go         # Resolver implementations
│   │   └── types.go            # Type definitions
│   ├── loader/                 # Data loaders
│   └── pubsub/                 # PubSub for subscriptions
├── pkg/
│   ├── auth/jwt.go             # JWT utilities
│   └── models/db.go            # Database models
├── .github/workflows/ci.yml    # CI pipeline
├── Dockerfile                   # Container build
├── docker-compose.yml           # Local development
├── gqlgen.yml                  # gqlgen configuration
├── Makefile                    # Build commands
└── README.md                   # This file
```

## Prerequisites

- Go 1.23+
- Docker (optional)
- Make (optional)

## Quick Start

### Local Development

```bash
# Clone and navigate
cd basic/projects/graphql

# Download dependencies
go mod download

# Run the server
go run ./cmd/server

# Or use make
make run
```

The server starts at `http://localhost:8080`

### GraphQL Playground

Open `http://localhost:8080` in your browser to access the GraphQL Playground.

### Example Queries

**Get all users:**
```graphql
query {
  users {
    id
    username
    email
    role
  }
}
```

**Get a single user:**
```graphql
query {
  user(id: "1") {
    id
    username
    email
  }
}
```

**Create a user:**
```graphql
mutation {
  createUser(input: {
    username: "newuser"
    email: "new@example.com"
    password: "password123"
    role: USER
  }) {
    id
    username
  }
}
```

**Get posts with filtering:**
```graphql
query {
  posts(authorId: "2", published: true) {
    id
    title
    content
    author {
      username
    }
  }
}
```

### Authentication

Some operations require JWT authentication. Include the token in the Authorization header:

```bash
curl -X POST http://localhost:8080/query \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <your-jwt-token>" \
  -d '{"query":"mutation { createPost(input: {title: \"New Post\", content: \"Content\"}) { id } }"}'
```

### Health Check

```bash
curl http://localhost:8080/health
# Returns: {"status":"healthy"}
```

## Docker

### Build Image

```bash
docker build -t graphql-api .
```

### Run with Docker Compose

```bash
docker-compose up -d
```

### Run Container

```bash
docker run -p 8080:8080 graphql-api
```

## Make Commands

```bash
make run          # Run the application
make build        # Build the binary
make test         # Run tests with race detection
make test-coverage # Generate coverage report
make vet          # Run go vet
make lint         # Run golangci-lint
make clean        # Clean build artifacts
make docker-build # Build Docker image
make docker-run   # Run with Docker Compose
make docker-stop  # Stop Docker Compose
```

## Environment Variables

| Variable    | Default | Description                           |
|-------------|---------|---------------------------------------|
| PORT        | 8080    | Server port                           |
| ENV         | development | Environment (development/production) |
| JWT_SECRET  | (dev default) | JWT signing secret               |

## GraphQL Schema

The API provides:

**Types:**
- `User` - User account with role-based access
- `Post` - Blog post with author and comments
- `Comment` - Comment on a post
- `Stats` - System statistics

**Queries:**
- `user(id)` - Get user by ID
- `users` - List users with filtering
- `me` - Get current authenticated user
- `post(id)` - Get post by ID
- `posts` - List posts with filtering
- `comments(postId)` - List comments for a post
- `stats` - Get system statistics

**Mutations:**
- `createUser` - Create a new user
- `updateUser` - Update user details
- `deleteUser` - Delete a user
- `createPost` - Create a new post
- `updatePost` - Update a post
- `deletePost` - Delete a post
- `publishPost` / `unpublishPost` - Toggle publication status
- `createComment` / `deleteComment` - Manage comments

**Subscriptions:**
- `userCreated` - New user created
- `postCreated` - New post created
- `commentAdded(postId)` - New comment on a post

## Testing

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## CI/CD

The project includes a GitHub Actions workflow that:

1. Runs tests with race detection
2. Builds the application
3. Builds the Docker image
4. Runs linter

## License

MIT
