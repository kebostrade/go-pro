# GraphQL with gqlgen in Go

Learn to build GraphQL APIs using gqlgen, the most popular GraphQL library for Go. This topic covers schema design, resolvers, mutations, and best practices.

## Learning Objectives

- Understand GraphQL fundamentals
- Design GraphQL schemas
- Write query and mutation resolvers
- Implement authentication and authorization
- Handle errors properly
- Use data loaders for performance
- Deploy GraphQL APIs

## Prerequisites

- Go 1.23+ installed
- Basic understanding of GraphQL concepts
- Familiarity with REST APIs (helpful but not required)

## Setup

```bash
# Install gqlgen
go install github.com/99designs/gqlgen@latest

# Create new project
mkdir graphql-api && cd graphql-api
go mod init github.com/username/graphql-api

# Initialize gqlgen
gqlgen init

# Run the server
go run server.go
```

## Examples

### 1. GraphQL Schema
`schema.graphqls` - Complete GraphQL schema definition

### 2. Server Setup
`server.go` - GraphQL HTTP server with handlers

### 3. Resolvers
`resolvers.go` - Query and mutation implementations

### 4. Data Models
`models.go` - Generated models and custom types

### 5. Complete API
`examples/graphql_api.go` - Full GraphQL API with authentication

## Running Examples

### Local Development

```bash
# Generate models from schema
go run github.com/99designs/gqlgen generate

# Run the server
go run server.go

# Access GraphQL Playground
open http://localhost:8080
```

### Testing Queries

```bash
# Using curl
curl -X POST http://localhost:8080/query \
  -H "Content-Type: application/json" \
  -d '{"query": "{ users { id username email } }"}'

# Using GraphQL Playground
# Open http://localhost:8080 in your browser
```

## GraphQL Basics

### Query Example

```graphql
query GetUsers {
  users {
    id
    username
    email
    posts {
      id
      title
    }
  }
}
```

### Mutation Example

```graphql
mutation CreateUser {
  createUser(input: {
    username: "alice"
    email: "alice@example.com"
    role: USER
  }) {
    id
    username
    email
  }
}
```

### Subscription Example

```graphql
subscription PostPublished {
  postPublished {
    id
    title
    author {
      username
    }
  }
}
```

## Key Concepts

### 1. Schema Design

```graphql
type User {
  id: ID!
  username: String!
  email: String!
  posts: [Post!]!
}

type Query {
  user(id: ID!): User
  users: [User!]!
}

type Mutation {
  createUser(input: CreateUserInput!): User!
}
```

### 2. Resolver Pattern

```go
func (r *queryResolver) User(ctx context.Context, id string) (*model.User, error) {
    // Fetch user from database
    user, err := r.DB.UserByID(id)
    if err != nil {
        return nil, err
    }
    return user, nil
}
```

### 3. Data Loading

```go
// Use data loaders to avoid N+1 queries
loaders := dataloaders.NewLoaders(r.DB)
ctx = dataloaders.ContextWithLoaders(ctx, loaders)

// In resolver
users := loaders.UserByID.LoadAll(ids)
```

### 4. Context Middleware

```go
func (r *resolver) Query() generated.QueryResolver {
    return &queryResolver{r}
}

func (r *queryResolver) User(ctx context.Context, id string) (*model.User, error) {
    // Access user from context (authentication)
    user := auth.ForContext(ctx)
    ...
}
```

## Best Practices

### 1. Schema Design

- **Use specific types**: Prefer specific scalars over String
- **Nullability**: Only make fields nullable if appropriate
- **Pagination**: Use cursor-based pagination for large lists
- **Enums**: Use enums for fixed sets of values
- **Input types**: Group related inputs into input types

### 2. Resolver Patterns

- **Separate business logic**: Keep resolvers thin
- **Use data loaders**: Prevent N+1 query problems
- **Error handling**: Return proper GraphQL errors
- **Validation**: Validate inputs in resolvers
- **Authorization**: Check permissions in resolvers

### 3. Performance

- **Data loaders**: Batch database queries
- **Caching**: Cache frequently accessed data
- **Query complexity**: Limit query complexity
- **Rate limiting**: Implement rate limiting
- **Query depth**: Limit query depth

### 4. Security

- **Authentication**: Verify user identity
- **Authorization**: Check permissions
- **Input validation**: Validate all inputs
- **Rate limiting**: Prevent abuse
- **Query cost analysis**: Limit expensive queries

## Authentication

### JWT Middleware

```go
func authMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        token := r.Header.Get("Authorization")
        claims, err := validateToken(token)

        ctx := context.WithValue(r.Context(), "user", claims)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}
```

### In Resolvers

```go
func (r *mutationResolver) CreateUser(ctx context.Context, input model.CreateUserInput) (*model.User, error) {
    // Check authentication
    user := auth.ForContext(ctx)
    if user == nil {
        return nil, fmt.Errorf("unauthorized")
    }

    // Check authorization
    if user.Role != "ADMIN" {
        return nil, fmt.Errorf("forbidden")
    }

    // Create user
    ...
}
```

## Error Handling

### Custom Errors

```go
type UserError struct {
    Message string
    Code    string
}

func (e UserError) Error() string {
    return e.Message
}

// Return error in resolver
return nil, UserError{
    Message: "User not found",
    Code:    "NOT_FOUND",
}
```

### Error Extensions

```go
func (r *queryResolver) User(ctx context.Context, id string) (*model.User, error) {
    user, err := r.DB.UserByID(id)
    if err != nil {
        if errors.Is(err, ErrNotFound) {
            return nil, gqlgen.Errorf("user not found: %s", id)
        }
        return nil, fmt.Errorf("failed to get user: %w", err)
    }
    return user, nil
}
```

## Testing

### Query Tests

```go
func TestGetUser(t *testing.T) {
    res := client.Query(`
        query {
            user(id: "1") {
                id
                username
            }
        }
    `)

    assert.NoError(t, res.Err)
    assert.Equal(t, "1", res.Get("user.id").String())
}
```

### Mutation Tests

```go
func TestCreateUser(t *testing.T) {
    res := client.Mutation(`
        mutation {
            createUser(input: {
                username: "test"
                email: "test@example.com"
                role: USER
            }) {
                id
                username
            }
        }
    `)

    assert.NoError(t, res.Err)
    assert.Equal(t, "test", res.Get("createUser.username").String())
}
```

## Deployment

### Docker

```dockerfile
FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o server

FROM alpine:latest
COPY --from=builder /app/server .
EXPOSE 8080
CMD ["./server"]
```

### Kubernetes

```yaml
apiVersion: v1
kind: Service
metadata:
  name: graphql-api
spec:
  selector:
    app: graphql-api
  ports:
  - port: 80
    targetPort: 8080
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: graphql-api
spec:
  replicas: 3
  selector:
    matchLabels:
      app: graphql-api
  template:
    metadata:
      labels:
        app: graphql-api
    spec:
      containers:
      - name: graphql-api
        image: graphql-api:latest
        ports:
        - containerPort: 8080
```

## Troubleshooting

### N+1 Query Problem

Use data loaders to batch queries:
```go
users := loaders.UserByID.LoadAll(userIDs)
```

### Slow Queries

- Implement query complexity analysis
- Add caching layer
- Optimize database queries
- Use persisted queries

### Schema Generation Issues

```bash
# Clean generated files
rm -rf graph/

# Regenerate
go run github.com/99designs/gqlgen generate
```

## Resources

- [gqlgen Documentation](https://gqlgen.com/)
- [GraphQL Specification](https://spec.graphql.org/)
- [GraphQL Best Practices](https://graphql.best practices/)
- [gqlgen GitHub](https://github.com/99designs/gqlgen)

## Next Steps

1. Complete all examples
2. Build a complete GraphQL API
3. Implement authentication and authorization
4. Add file uploads
5. Implement real-time subscriptions
6. Add caching and rate limiting
7. Set up CI/CD pipeline
8. Deploy to production
