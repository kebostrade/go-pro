# SD-10: Caching Strategies

Learn various caching strategies and when to apply each for optimal performance.

## Overview

Caching is one of the most effective ways to improve system performance. This lesson covers multiple caching strategies and their trade-offs.

## Learning Objectives

- Understand cache-aside, write-through, and other patterns
- Implement multi-level caching
- Handle cache stampede
- Choose appropriate cache invalidation strategies

## Caching Strategies

### 1. Cache-Aside (Lazy Loading)

The most common pattern:

```
┌─────────┐     ┌─────────┐     ┌─────────┐
│  App    │────▶│  Cache  │────▶│   DB    │
└─────────┘     └─────────┘     └─────────┘
```

**Read:**
1. Check cache first
2. If miss, load from database
3. Store in cache
4. Return data

**Write:**
1. Write to database
2. Delete from cache (not update)

```go
func (s *UserService) GetUser(ctx context.Context, id string) (*User, error) {
    // Try cache first
    key := "user:" + id
    cached, err := s.cache.Get(ctx, key)
    if err == nil {
        var user User
        json.Unmarshal([]byte(cached), &user)
        return &user, nil
    }
    
    // Cache miss - load from DB
    user, err := s.repo.GetByID(ctx, id)
    if err != nil {
        return nil, err
    }
    
    // Store in cache
    data, _ := json.Marshal(user)
    s.cache.Set(ctx, key, string(data), time.Hour)
    
    return user, nil
}
```

### 2. Write-Through

Write to cache and database simultaneously:

```
┌─────────┐     ┌─────────┐     ┌─────────┐
│  App    │────▶│  Cache  │────▶│   DB    │
└─────────┘     └─────────┘     └─────────┘
```

```go
func (s *UserService) CreateUser(ctx context.Context, user *User) error {
    // Write to cache first
    key := "user:" + user.ID
    data, _ := json.Marshal(user)
    s.cache.Set(ctx, key, string(data), 0) // No expiry
    
    // Then write to DB
    return s.repo.Create(ctx, user)
}
```

**Pros:** Always fresh data in cache
**Cons:** Slower writes

### 3. Write-Behind (Write-Back)

Write to cache, async write to DB:

```go
func (s *UserService) CreateUser(ctx context.Context, user *User) error {
    // Write to cache only
    key := "user:" + user.ID
    data, _ := json.Marshal(user)
    s.cache.Set(ctx, key, string(data), time.Hour)
    
    // Async write to DB (use queue)
    s.writeQueue <- user
    
    return nil
}
```

**Pros:** Fast writes
**Cons:** Risk of data loss

### 4. Refresh-Ahead

Proactively refresh cache before expiration:

```go
func (s *UserService) GetUser(ctx context.Context, id string) (*User, error) {
    key := "user:" + id
    
    user, err := s.getFromCache(ctx, key)
    if err != nil {
        return s.loadUser(ctx, id)
    }
    
    // If close to expiration, refresh in background
    if s.cache.TTL(key) < time.Minute {
        go s.refreshUserCache(ctx, id)
    }
    
    return user, nil
}
```

## Multi-Level Caching

```
┌─────────────────────────────────────────┐
│              Application                 │
└─────────────────┬───────────────────────┘
                  │
        ┌─────────▼─────────┐
        │   L1: Local Cache  │  (In-memory, fast)
        │   e.g., BigCache   │
        └─────────┬───────────┘
                  │
        ┌─────────▼─────────┐
        │   L2: Redis       │  (Distributed)
        └─────────┬─────────┘
                  │
        ┌─────────▼─────────┐
        │   Database        │  (Slowest)
        └───────────────────┘
```

```go
type MultiLevelCache struct {
    local  sync.Map   // L1: In-memory
    redis *redis.Client  // L2: Redis
    db    Database
}

func (m *MultiLevelCache) Get(ctx context.Context, key string) (interface{}, error) {
    // L1: Check local cache
    if val, ok := m.local.Load(key); ok {
        return val, nil
    }
    
    // L2: Check Redis
    if val, err := m.redis.Get(ctx, key).Result(); err == nil {
        m.local.Store(key, val)  // Populate L1
        return val, nil
    }
    
    // L3: Query DB
    val, err := m.db.Get(ctx, key)
    if err != nil {
        return nil, err
    }
    
    // Populate caches
    m.local.Store(key, val)
    m.redis.Set(ctx, key, val, time.Hour)
    
    return val, nil
}
```

## Cache Invalidation Strategies

| Strategy | Description | Use Case |
|----------|-------------|----------|
| TTL | Time-based expiration | General purpose |
| LRU | Evict least recently used | Limited cache size |
| LFU | Evict least frequently used | Popular items |
| Write-through | Update on write | Data consistency |
|主动失效 | Manual invalidation | Known events |

## Handling Cache Stampede

When cache expires and many requests hit DB simultaneously:

```go
func (s *UserService) GetUserWithLock(ctx context.Context, id string) (*User, error) {
    key := "user:" + id
    
    // Singleflight ensures only one request hits DB
    result, err := s.sf.Do(key, func() (interface{}, error) {
        return s.getFromCacheOrDB(ctx, key, id)
    })
    
    return result.(*User), err
}
```

## Examples

See `examples/` directory for:
- `cache_strategies.go` - All strategies implemented
- `multi_level_cache.go` - L1/L2 caching
- `cache_stampede.go` - Stampede prevention

## Exercises

See `exercises/problems.md` for hands-on practice.

## Quiz

Test your knowledge with `quiz.md`.

## Summary

- Cache-aside is most common for reads
- Write-through ensures consistency
- Multi-level caching optimizes performance
- Handle stampede with singleflight

## Next Steps

Continue to [SD-11: Load Balancing](11-load-balancing/README.md)
