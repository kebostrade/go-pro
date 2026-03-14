# Exercise: Caching Strategies

## Problem 1: Cache Strategy Selection

For each scenario, recommend the best caching strategy:

| Scenario | Recommended Strategy | Rationale |
|----------|---------------------|-----------|
| User profile data (frequently read, rarely updated) | | |
| Shopping cart (frequently updated) | | |
| Product catalog (read-heavy, updates at night) | | |
| Session data (short-lived, user-specific) | | |
| Real-time stock prices (constantly changing) | | |

---

## Problem 2: Implementing Cache-Aside

Complete the cache-aside implementation:

```go
func (s *UserService) GetUser(ctx context.Context, id string) (*User, error) {
    key := "user:" + id
    
    // TODO: Step 1 - Check cache first
    
    
    // TODO: Step 2 - If cache miss, fetch from DB
    
    
    // TODO: Step 3 - Store in cache and return
    
    
    return user, nil
}
```

---

## Problem 3: Multi-Level Cache Design

Design a 3-level caching system:

```
┌─────────────────────────────────────┐
│              App                     │
└──────────────┬──────────────────────┘
               │
    ┌──────────▼──────────┐
    │    L1: Local        │
    │    (In-memory)      │
    └──────────┬──────────┘
               │
    ┌──────────▼──────────┐
    │    L2: Distributed  │
    │    (Redis)          │
    └──────────┬──────────┘
               │
    ┌──────────▼──────────┐
    │    L3: Database     │
    └─────────────────────┘
```

Write the Get function:

```go
func (m *Cache) Get(ctx context.Context, key string) (interface{}, error) {
    // Implement 3-level cache lookup
    
    
    return nil, nil
}
```

---

## Problem 4: Handling Cache Stampede

Implement protection against cache stampede using singleflight:

```go
var sf = singleflight.Group{}

func (s *UserService) GetUserSafe(ctx context.Context, id string) (*User, error) {
    key := "user:" + id
    
    // Use singleflight to prevent stampede
    
    
    return nil, nil
}
```

---

## Problem 5: Cache Invalidation

You have a user profile that is cached. When should you invalidate the cache?

| Event | Invalidate? | Strategy |
|-------|-------------|----------|
| User updates profile | ? | |
| User changes password | ? | |
| Admin deletes user | ? | |
| User logs out | ? | |
| TTL expires | ? | |

---

## Problem 6: Calculating Cache Size

Calculate the cache size needed:

**Requirements:**
- 10,000 requests per second
- 80% hit rate
- Average item size: 5KB
- TTL: 1 hour

**Questions:**
1. How many items per second need to be cached?
2. What's the cache size needed for 1 hour?
3. What's the cache size needed for 24 hours?
