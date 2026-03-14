# Building GraphQL APIs with Go and gqlgen

Develop GraphQL APIs using Go and the gqlgen library.

## Learning Objectives

- Define GraphQL schemas
- Implement resolvers
- Handle mutations and queries
- Add subscriptions
- Implement authentication
- Optimize with DataLoader

## Theory

### Schema Definition

```graphql
# schema.graphql

type User {
    id: ID!
    name: String!
    email: String!
    posts: [Post!]!
    createdAt: Time!
}

type Post {
    id: ID!
    title: String!
    content: String!
    author: User!
    comments: [Comment!]!
    createdAt: Time!
}

type Comment {
    id: ID!
    content: String!
    author: User!
    post: Post!
    createdAt: Time!
}

type Query {
    user(id: ID!): User
    users(limit: Int = 10, offset: Int = 0): [User!]!
    post(id: ID!): Post
    posts(limit: Int = 10): [Post!]!
    me: User!
}

type Mutation {
    createUser(input: CreateUserInput!): User!
    updateUser(id: ID!, input: UpdateUserInput!): User!
    deleteUser(id: ID!): Boolean!
    
    createPost(input: CreatePostInput!): Post!
    updatePost(id: ID!, input: UpdatePostInput!): Post!
    deletePost(id: ID!): Boolean!
    
    createComment(input: CreateCommentInput!): Comment!
}

type Subscription {
    onPostCreated: Post!
    onCommentAdded(postId: ID!): Comment!
}

input CreateUserInput {
    name: String!
    email: String!
    password: String!
}

input UpdateUserInput {
    name: String
    email: String
}

input CreatePostInput {
    title: String!
    content: String!
}

input UpdatePostInput {
    title: String
    content: String
}

input CreateCommentInput {
    postId: ID!
    content: String!
}

scalar Time
```

### Resolver Implementation

```go
func NewResolver(db *sql.DB) graph.Config {
    return graph.Config{
        Resolvers: &Resolver{db: db},
    }
}

type Resolver struct {
    db *sql.DB
}

func (r *queryResolver) User(ctx context.Context, id string) (*model.User, error) {
    var user model.User
    err := r.db.QueryRowContext(ctx,
        "SELECT id, name, email, created_at FROM users WHERE id = $1",
        id,
    ).Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt)
    
    if err == sql.ErrNoRows {
        return nil, gqlerror.Errorf("user not found")
    }
    if err != nil {
        return nil, err
    }
    
    return &user, nil
}

func (r *queryResolver) Users(ctx context.Context, limit *int, offset *int) ([]*model.User, error) {
    rows, err := r.db.QueryContext(ctx,
        "SELECT id, name, email, created_at FROM users ORDER BY created_at DESC LIMIT $1 OFFSET $2",
        *limit, *offset,
    )
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var users []*model.User
    for rows.Next() {
        var user model.User
        if err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt); err != nil {
            return nil, err
        }
        users = append(users, &user)
    }

    return users, nil
}

func (r *mutationResolver) CreateUser(ctx context.Context, input model.CreateUserInput) (*model.User, error) {
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
    if err != nil {
        return nil, err
    }

    var user model.User
    err = r.db.QueryRowContext(ctx,
        "INSERT INTO users (name, email, password) VALUES ($1, $2, $3) RETURNING id, name, email, created_at",
        input.Name, input.Email, hashedPassword,
    ).Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt)

    if err != nil {
        return nil, fmt.Errorf("create user: %w", err)
    }

    return &user, nil
}
```

### DataLoader for N+1 Prevention

```go
type Loaders struct {
    UserByID    *dataloader.Loader
    PostsByUser *dataloader.Loader
}

func NewLoaders(db *sql.DB) *Loaders {
    return &Loaders{
        UserByID: dataloader.NewBatchedLoader(func(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
            ids := make([]string, len(keys))
            for i, key := range keys {
                ids[i] = key.String()
            }

            users, err := batchGetUsers(ctx, db, ids)
            if err != nil {
                return make([]*dataloader.Result, len(keys))
            }

            results := make([]*dataloader.Result, len(keys))
            userMap := make(map[string]*model.User)
            for _, u := range users {
                userMap[u.ID] = u
            }

            for i, id := range ids {
                if user, ok := userMap[id]; ok {
                    results[i] = &dataloader.Result{Data: user}
                } else {
                    results[i] = &dataloader.Result{Error: fmt.Errorf("user not found")}
                }
            }

            return results
        }),
    }
}

func (r *userResolver) Posts(ctx context.Context, obj *model.User) ([]*model.Post, error) {
    loaders := ctx.Value("loaders").(*Loaders)
    thunk := loaders.PostsByUser.Load(ctx, dataloader.StringKey(obj.ID))
    result, err := thunk()
    if err != nil {
        return nil, err
    }
    return result.Data.([]*model.Post), nil
}
```

### Subscriptions

```go
type subscriptionResolver struct {
    *Resolver
    postChan   chan *model.Post
    commentChs map[string][]chan *model.Comment
    mu         sync.Mutex
}

func (r *subscriptionResolver) OnPostCreated(ctx context.Context) (<-chan *model.Post, error) {
    ch := make(chan *model.Post, 1)
    
    go func() {
        defer close(ch)
        for {
            select {
            case <-ctx.Done():
                return
            case post := <-r.postChan:
                ch <- post
            }
        }
    }()
    
    return ch, nil
}

func (r *mutationResolver) CreatePost(ctx context.Context, input model.CreatePostInput) (*model.Post, error) {
    userID := middleware.GetUserID(ctx)

    var post model.Post
    err := r.db.QueryRowContext(ctx,
        "INSERT INTO posts (title, content, author_id) VALUES ($1, $2, $3) RETURNING id, title, content, author_id, created_at",
        input.Title, input.Content, userID,
    ).Scan(&post.ID, &post.Title, &post.Content, &post.AuthorID, &post.CreatedAt)

    if err != nil {
        return nil, err
    }

    r.postChan <- &post

    return &post, nil
}
```

### Authentication Middleware

```go
func AuthMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        authHeader := r.Header.Get("Authorization")
        if authHeader == "" {
            next.ServeHTTP(w, r)
            return
        }

        tokenString := strings.TrimPrefix(authHeader, "Bearer ")
        claims, err := validateJWT(tokenString)
        if err != nil {
            next.ServeHTTP(w, r)
            return
        }

        ctx := context.WithValue(r.Context(), "user_id", claims.UserID)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}

func (r *queryResolver) Me(ctx context.Context) (*model.User, error) {
    userID, ok := ctx.Value("user_id").(string)
    if !ok {
        return nil, gqlerror.Errorf("unauthenticated")
    }

    return r.User(ctx, userID)
}
```

### Server Setup

```go
func main() {
    db := setupDB()
    defer db.Close()

    resolver := NewResolver(db)

    srv := handler.NewDefaultServer(generated.NewExecutableSchema(resolver))
    srv.AddTransport(transport.Websocket{
        Upgrader: websocket.Upgrader{
            CheckOrigin: func(r *http.Request) bool { return true },
        },
    })
    srv.Use(extension.Introspection{})

    http.Handle("/", playground.Handler("GraphQL", "/query"))
    http.Handle("/query", AuthMiddleware(srv))

    log.Println("Server started on :8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}
```

## Security Considerations

```go
func validateInput(input string) error {
    if len(input) > 10000 {
        return errors.New("input too long")
    }
    if strings.Contains(input, "<script>") {
        return errors.New("invalid input")
    }
    return nil
}

func complexityLimit() graphql.Handler {
    return func(ctx context.Context, next graphql.OperationHandler) graphql.ResponseHandler {
        op := graphql.GetOperationContext(ctx).Operation
        if op == nil {
            return next(ctx)
        }

        complexity := calculateComplexity(op)
        if complexity > 1000 {
            return graphql.ErrorResponse(ctx, "query too complex")
        }

        return next(ctx)
    }
}
```

## Performance Tips

```go
func (r *queryResolver) Users(ctx context.Context, limit *int, offset *int) ([]*model.User, error) {
    if *limit > 100 {
        *limit = 100
    }

    cacheKey := fmt.Sprintf("users:%d:%d", *limit, *offset)
    if cached, ok := r.cache.Get(cacheKey); ok {
        return cached.([]*model.User), nil
    }

    users, err := r.fetchUsers(ctx, *limit, *offset)
    if err != nil {
        return nil, err
    }

    r.cache.Set(cacheKey, users, time.Minute)
    return users, nil
}
```

## Exercises

1. Build a blog API with users and posts
2. Implement real-time notifications
3. Add pagination and filtering
4. Create a file upload mutation

## Validation

```bash
cd exercises
go generate ./...
go test -v ./...
curl -X POST http://localhost:8080/query -d '{"query":"{ users { id name } }"}'
```

## Key Takeaways

- Use DataLoader to prevent N+1 queries
- Implement proper authentication
- Limit query complexity
- Cache expensive queries
- Handle subscriptions carefully

## Next Steps

**[AT-00: System Design](../AT-00-system-design/README.md)**

---

GraphQL: ask for what you need, get exactly that. 📊
