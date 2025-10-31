# 🐘🔴 PostgreSQL & Redis with Go - Tutorial 20

A comprehensive tutorial on working with PostgreSQL and Redis in Go, covering database operations, caching patterns, and real-world integration examples.

## 📋 Table of Contents

- [Overview](#overview)
- [Prerequisites](#prerequisites)
- [Quick Start](#quick-start)
- [Project Structure](#project-structure)
- [Examples](#examples)
- [Patterns](#patterns)
- [Best Practices](#best-practices)

## 🎯 Overview

This tutorial demonstrates:

### PostgreSQL
- **pgx/v5**: High-performance PostgreSQL driver
- **GORM**: Full-featured ORM
- Connection pooling
- Transactions
- Batch operations
- Prepared statements

### Redis
- Basic data structures (String, Hash, List, Set, Sorted Set)
- Pub/Sub messaging
- Caching patterns
- Distributed locks
- Rate limiting
- Session management

### Combined Patterns
- Cache-Aside (Lazy Loading)
- Write-Through Cache
- Cache Invalidation
- Session Management with PostgreSQL + Redis

## 📦 Prerequisites

- Go 1.22 or higher
- Docker and Docker Compose
- PostgreSQL 16
- Redis 7

## 🚀 Quick Start

### 1. Clone and Setup

```bash
cd basic/projects/postgres-redis-go
cp .env.example .env
```

### 2. Start Databases

```bash
make docker-up
```

This starts:
- PostgreSQL on port 5432
- Redis on port 6379
- Adminer (PostgreSQL UI) on port 8080
- Redis Commander on port 8081

### 3. Install Dependencies

```bash
make deps
```

### 4. Run Examples

```bash
# PostgreSQL examples
make run-pgx          # pgx driver examples
make run-gorm         # GORM ORM examples

# Redis examples
make run-redis-basic  # Basic Redis operations
make run-redis-pubsub # Pub/Sub messaging
make run-redis-patterns # Common patterns

# Combined examples
make run-cache        # Cache-Aside pattern
make run-session      # Session management

# Run all examples
make run-all
```

## 📁 Project Structure

```
postgres-redis-go/
├── postgres/
│   ├── pgx/              # pgx driver examples
│   │   └── main.go       # Connection, CRUD, transactions, batching
│   ├── gorm/             # GORM ORM examples
│   │   └── main.go       # Models, associations, migrations, hooks
│   └── migrations/       # Database migrations
│       └── init.sql      # Initial schema
├── redis/
│   ├── basic/            # Basic Redis operations
│   │   └── main.go       # Strings, hashes, lists, sets, sorted sets
│   ├── pubsub/           # Pub/Sub messaging
│   │   └── main.go       # Channels, patterns, broadcasting
│   └── patterns/         # Common Redis patterns
│       └── main.go       # Locks, rate limiting, leaderboards
├── combined/
│   ├── cache/            # Caching patterns
│   │   └── main.go       # Cache-aside, write-through, invalidation
│   └── session/          # Session management
│       └── main.go       # User sessions with PostgreSQL + Redis
├── docker-compose.yml    # Docker services
├── Makefile              # Build and run commands
└── README.md             # This file
```

## 💡 Examples

### PostgreSQL with pgx

```go
// Connection pool
config, _ := pgxpool.ParseConfig(connString)
config.MaxConns = 25
config.MinConns = 5
pool, _ := pgxpool.NewWithConfig(ctx, config)

// Query
var user User
pool.QueryRow(ctx, "SELECT * FROM users WHERE id = $1", id).
    Scan(&user.ID, &user.Username, &user.Email)

// Transaction
tx, _ := pool.Begin(ctx)
tx.Exec(ctx, "INSERT INTO users ...")
tx.Commit(ctx)

// Batch operations
batch := &pgx.Batch{}
batch.Queue("INSERT INTO users ...")
results := pool.SendBatch(ctx, batch)
```

### PostgreSQL with GORM

```go
// Auto migration
db.AutoMigrate(&User{}, &Profile{}, &Post{})

// Create
user := User{Username: "john", Email: "john@example.com"}
db.Create(&user)

// Query
db.Where("username = ?", "john").First(&user)

// Update
db.Model(&user).Update("email", "new@example.com")

// Associations
db.Preload("Profile").Preload("Posts").Find(&users)

// Transaction
db.Transaction(func(tx *gorm.DB) error {
    tx.Create(&user)
    tx.Create(&profile)
    return nil
})
```

### Redis Basics

```go
// String operations
rdb.Set(ctx, "key", "value", 0)
val, _ := rdb.Get(ctx, "key").Result()

// Hash operations
rdb.HSet(ctx, "user:1", "name", "John")
rdb.HGetAll(ctx, "user:1")

// List operations
rdb.LPush(ctx, "tasks", "task1", "task2")
tasks, _ := rdb.LRange(ctx, "tasks", 0, -1).Result()

// Set operations
rdb.SAdd(ctx, "tags", "go", "redis")
members, _ := rdb.SMembers(ctx, "tags").Result()

// Sorted set (leaderboard)
rdb.ZAdd(ctx, "scores", redis.Z{Score: 100, Member: "player1"})
top, _ := rdb.ZRevRangeWithScores(ctx, "scores", 0, 9).Result()
```

### Redis Pub/Sub

```go
// Subscribe
pubsub := rdb.Subscribe(ctx, "notifications")
ch := pubsub.Channel()

go func() {
    for msg := range ch {
        fmt.Printf("Received: %s\n", msg.Payload)
    }
}()

// Publish
rdb.Publish(ctx, "notifications", "Hello, World!")

// Pattern subscription
pubsub := rdb.PSubscribe(ctx, "user:*:events")
```

## 🎯 Patterns

### 1. Cache-Aside Pattern

```go
func GetUser(id int) (*User, error) {
    // 1. Check cache
    cached, err := rdb.Get(ctx, fmt.Sprintf("user:%d", id)).Result()
    if err == nil {
        var user User
        json.Unmarshal([]byte(cached), &user)
        return &user, nil
    }

    // 2. Query database
    var user User
    db.QueryRow(ctx, "SELECT * FROM users WHERE id = $1", id).
        Scan(&user.ID, &user.Username, &user.Email)

    // 3. Cache result
    data, _ := json.Marshal(user)
    rdb.Set(ctx, fmt.Sprintf("user:%d", id), data, 5*time.Minute)

    return &user, nil
}
```

### 2. Distributed Lock

```go
// Acquire lock
acquired, _ := rdb.SetNX(ctx, "lock:resource", "unique-id", 10*time.Second).Result()

if acquired {
    // Critical section
    defer func() {
        // Release lock
        script := `
            if redis.call("get", KEYS[1]) == ARGV[1] then
                return redis.call("del", KEYS[1])
            end
        `
        rdb.Eval(ctx, script, []string{"lock:resource"}, "unique-id")
    }()
}
```

### 3. Rate Limiting

```go
func CheckRateLimit(userID string, limit int, window time.Duration) bool {
    key := fmt.Sprintf("rate:%s", userID)
    count, _ := rdb.Incr(ctx, key).Result()
    
    if count == 1 {
        rdb.Expire(ctx, key, window)
    }
    
    return count <= int64(limit)
}
```

### 4. Session Management

```go
// Create session
sessionData := map[string]interface{}{
    "user_id": "123",
    "username": "john",
}
rdb.HSet(ctx, "session:abc123", sessionData)
rdb.Expire(ctx, "session:abc123", 30*time.Minute)

// Get session
session, _ := rdb.HGetAll(ctx, "session:abc123").Result()

// Refresh session
rdb.Expire(ctx, "session:abc123", 30*time.Minute)

// Delete session
rdb.Del(ctx, "session:abc123")
```

## 🏆 Best Practices

### PostgreSQL

1. **Use Connection Pooling**
   ```go
   config.MaxConns = 25
   config.MinConns = 5
   config.MaxConnLifetime = time.Hour
   ```

2. **Use Prepared Statements**
   - Prevents SQL injection
   - Improves performance

3. **Use Transactions for Multiple Operations**
   ```go
   tx, _ := pool.Begin(ctx)
   defer tx.Rollback(ctx)
   // ... operations
   tx.Commit(ctx)
   ```

4. **Use Batch Operations for Bulk Inserts**
   - Significantly faster than individual inserts

5. **Always Use Context**
   - Enables timeouts and cancellation

### Redis

1. **Set Expiration on Keys**
   ```go
   rdb.Set(ctx, key, value, 5*time.Minute)
   ```

2. **Use Pipelining for Multiple Commands**
   ```go
   pipe := rdb.Pipeline()
   pipe.Set(ctx, "key1", "val1", 0)
   pipe.Set(ctx, "key2", "val2", 0)
   pipe.Exec(ctx)
   ```

3. **Use Appropriate Data Structures**
   - Strings: Simple values
   - Hashes: Objects
   - Lists: Queues
   - Sets: Unique items
   - Sorted Sets: Leaderboards

4. **Monitor Memory Usage**
   - Redis is in-memory
   - Set maxmemory policy

5. **Use Redis for What It's Good At**
   - Caching
   - Session storage
   - Real-time analytics
   - Message queuing

### Combined Patterns

1. **Cache Invalidation Strategy**
   - Invalidate on write
   - Set appropriate TTL
   - Use cache keys consistently

2. **Error Handling**
   - Cache failures shouldn't break app
   - Fall back to database

3. **Monitoring**
   - Track cache hit rate
   - Monitor database query times
   - Alert on connection pool exhaustion

## 🔧 Configuration

### Environment Variables

```env
# PostgreSQL
POSTGRES_HOST=localhost
POSTGRES_PORT=5432
POSTGRES_USER=postgres
POSTGRES_PASSWORD=postgres
POSTGRES_DB=go_tutorial

# Redis
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0
```

### Connection Pool Settings

```go
// PostgreSQL
config.MaxConns = 25              // Maximum connections
config.MinConns = 5               // Minimum idle connections
config.MaxConnLifetime = 1*time.Hour
config.MaxConnIdleTime = 30*time.Minute

// Redis
redis.Options{
    PoolSize: 10,                 // Maximum connections
    MinIdleConns: 5,              // Minimum idle connections
    MaxRetries: 3,                // Retry failed commands
}
```

## 📊 Performance Tips

1. **Use Connection Pooling** - Reuse connections
2. **Batch Operations** - Reduce round trips
3. **Pipelining** - Send multiple commands at once
4. **Appropriate Indexes** - Speed up queries
5. **Cache Frequently Accessed Data** - Reduce database load
6. **Monitor and Profile** - Identify bottlenecks

## 🧪 Testing

```bash
# Run tests
make test

# Run with coverage
go test -cover ./...

# Run with race detection
go test -race ./...
```

## 📚 Resources

- [pgx Documentation](https://github.com/jackc/pgx)
- [GORM Documentation](https://gorm.io/)
- [go-redis Documentation](https://redis.uptrace.dev/)
- [PostgreSQL Documentation](https://www.postgresql.org/docs/)
- [Redis Documentation](https://redis.io/docs/)

## 🎓 Learning Outcomes

After completing this tutorial, you will understand:

- ✅ PostgreSQL connection pooling and configuration
- ✅ CRUD operations with pgx and GORM
- ✅ Database transactions and batch operations
- ✅ Redis data structures and operations
- ✅ Pub/Sub messaging patterns
- ✅ Caching strategies (cache-aside, write-through)
- ✅ Distributed locks and rate limiting
- ✅ Session management
- ✅ Performance optimization techniques
- ✅ Production-ready patterns

## 📝 License

MIT License - feel free to use this tutorial for learning!

---

**Built with ❤️ for the Go learning community**

