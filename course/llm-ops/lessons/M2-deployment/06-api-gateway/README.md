# IO-06: API Gateway Patterns

**Duration**: 2 hours
**Module**: 2 - Deployment & Serving

## Learning Objectives

- Implement API gateway patterns for LLM traffic management
- Design rate limiting strategies for LLM workloads
- Build authentication and authorization for LLM APIs
- Implement caching strategies for cost optimization

## What is an API Gateway?

An API gateway is the entry point for all client requests to backend services. For LLM applications, it handles:

- Authentication and authorization
- Rate limiting and quota management
- Request routing and load balancing
- Response caching
- Monitoring and logging

```
┌─────────────────────────────────────────────────────────────────┐
│                       Client Requests                            │
│                    (Web, Mobile, API Clients)                    │
├─────────────────────────────────────────────────────────────────┐
│                      API Gateway                                 │
├─────────────────────────────────────────────────────────────────┤
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐           │
│  │  AuthN/Z    │  │Rate Limiting│  │   Caching   │           │
│  │  (JWT, API  │  │ (Token,     │  │ (Redis,     │           │
│  │   Keys)     │  │  User)      │  │  Memory)    │           │
│  └─────────────┘  └─────────────┘  └─────────────┘           │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐           │
│  │  Routing    │  │   Load      │  │  Metrics &  │           │
│  │  (Path,     │  │   Balancing │  │  Logging    │           │
│  │   Header)   │  │  (Round-    │  │  (Prometheus│           │
│  │             │  │   Robin)    │  │   , Tracing)│           │
│  └─────────────┘  └─────────────┘  └─────────────┘           │
├─────────────────────────────────────────────────────────────────┤
│                    LLM Inference Servers                         │
│                 (vLLM, TGI, OpenAI API)                        │
└─────────────────────────────────────────────────────────────────┘
```

## 1. Rate Limiting

LLM APIs require sophisticated rate limiting due to:

- High cost per request (token-based pricing)
- Variable response times
- GPU resource constraints

### Rate Limiting Strategies

| Strategy | Description | Use Case |
|----------|-------------|----------|
| **Token Bucket** | Fixed tokens per user, replenished over time | General API protection |
| **Leaky Bucket** | Fixed rate, queue excess requests | Smooth traffic |
| **Fixed Window** | Fixed requests per time window | Simple quotas |
| **Sliding Window** | Rolling time window for accuracy | Precise limits |

### Go Implementation

```go
package ratelimit

import (
	"sync"
	"time"
)

type TokenBucket struct {
	tokens     int
	maxTokens  int
	refillRate time.Duration
	mu         sync.Mutex
	lastRefill time.Time
}

func NewTokenBucket(maxTokens int, refillRate time.Duration) *TokenBucket {
	return &TokenBucket{
		tokens:     maxTokens,
		maxTokens:  maxTokens,
		refillRate: refillRate,
		lastRefill: time.Now(),
	}
}

func (tb *TokenBucket) Allow() bool {
	tb.mu.Lock()
	defer tb.mu.Unlock()
	
	tb.refill()
	
	if tb.tokens > 0 {
		tb.tokens--
		return true
	}
	return false
}

func (tb *TokenBucket) refill() {
	now := time.Now()
	elapsed := now.Sub(tb.lastRefill)
	
	refills := int(elapsed / tb.refillRate)
	if refills > 0 {
		tb.tokens = min(tb.maxTokens, tb.tokens+refills)
		tb.lastRefill = now
	}
}
```

### Rate Limiting with User Tiers

```go
type TierConfig struct {
	RequestsPerMinute int
	TokensPerMinute   int
	MaxConcurrent    int
}

var TierLimits = map[string]TierConfig{
	"free": {
		RequestsPerMinute: 10,
		TokensPerMinute:   10000,
		MaxConcurrent:    1,
	},
	"pro": {
		RequestsPerMinute: 100,
		TokensPerMinute:   100000,
		MaxConcurrent:    5,
	},
	"enterprise": {
		RequestsPerMinute: 1000,
		TokensPerMinute:   1000000,
		MaxConcurrent:    50,
	},
}
```

## 2. Authentication

### JWT Authentication

```go
package auth

import (
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type JWTClaims struct {
	UserID string `json:"user_id"`
	Tier   string `json:"tier"`
	jwt.RegisteredClaims
}

func ValidateJWT(tokenString, secret string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, 
		func(t *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		})
	
	if err != nil {
		return nil, err
	}
	
	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}
	
	return nil, jwt.ErrSignatureInvalid
}

func GenerateJWT(userID, tier, secret string) (string, error) {
	claims := JWTClaims{
		UserID: userID,
		Tier:   tier,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}
```

### API Key Authentication

```go
type APIKeyStore struct {
	keys map[string]APIKey
	mu   sync.RWMutex
}

type APIKey struct {
	Key       string
	UserID    string
	Tier      string
	CreatedAt time.Time
	ExpiresAt *time.Time
	Revoked   bool
}

func (s *APIKeyStore) Validate(key string) (*APIKey, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	apiKey, ok := s.keys[key]
	if !ok {
		return nil, ErrInvalidKey
	}
	
	if apiKey.Revoked {
		return nil, ErrRevokedKey
	}
	
	if apiKey.ExpiresAt != nil && time.Now().After(*apiKey.ExpiresAt) {
		return nil, ErrExpiredKey
	}
	
	return &apiKey, nil
}
```

## 3. Request Routing

### Multi-Backend Routing

```go
package proxy

import (
	"net/http"
	"net/url"
)

type BackendPool struct {
	backends []*url.URL
	current  int
}

func NewBackendPool(backends []string) (*BackendPool, error) {
	urls := make([]*url.URL, len(backends))
	for i, b := range backends {
		u, err := url.Parse(b)
		if err != nil {
			return nil, err
		}
		urls[i] = u
	}
	
	return &BackendPool{backends: urls}, nil
}

func (p *BackendPool) Next() *url.URL {
	backend := p.backends[p.current]
	p.current = (p.current + 1) % len(p.backends)
	return backend
}

func (p *BackendPool) Proxy(w http.ResponseWriter, r *http.Request) {
	backend := p.Next()
	
	r.URL.Scheme = backend.Scheme
	r.URL.Host = backend.Host
	r.Host = backend.Host
	
	http.DefaultTransport.RoundTrip(r)
}
```

### Header-Based Routing

```go
func (p *Proxy) routeRequest(r *http.Request) *url.URL {
	// Route based on model header
	if model := r.Header.Get("X-Model"); model != "" {
		if backend, ok := p.modelBackends[model]; ok {
			return backend
		}
	}
	
	// Route based on user tier
	claims := getClaims(r.Context())
	if claims.Tier == "enterprise" {
		return p.premiumBackend
	}
	
	return p.defaultBackend
}
```

## 4. Caching Strategies

### Response Caching for LLM

```go
package cache

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"
	
	"github.com/redis/go-redis/v9"
)

type LLMCache struct {
	client *redis.Client
	ttl    time.Duration
}

type CacheKey struct {
	Model      string   `json:"model"`
	Messages   []string `json:"messages"`
	Temperature float64 `json:"temperature"`
	MaxTokens  int      `json:"max_tokens"`
}

func (c *LLMCache) Get(ctx context.Context, prompt string) (string, bool) {
	key := c.makeKey(prompt)
	
	result, err := c.client.Get(ctx, key).Result()
	if err != nil {
		return "", false
	}
	
	return result, true
}

func (c *LLMCache) Set(ctx context.Context, prompt, response string) {
	key := c.makeKey(prompt)
	c.client.Set(ctx, key, response, c.ttl)
}

func (c *LLMCache) makeKey(prompt string) string {
	hash := sha256.Sum256([]byte(prompt))
	return fmt.Sprintf("llm:response:%s", hex.EncodeToString(hash[:]))
}
```

### Prompt Caching

```go
// Cache system prompts that don't change often
type PromptCache struct {
	systemPrompts map[string]string
	mu            sync.RWMutex
}

func (pc *PromptCache) GetSystemPrompt(template string) string {
	pc.mu.RLock()
	defer pc.mu.RUnlock()
	return pc.systemPrompts[template]
}

func (pc *PromptCache) SetSystemPrompt(template, prompt string) {
	pc.mu.Lock()
	defer pc.mu.Unlock()
	pc.systemPrompts[template] = prompt
}
```

## 5. Load Balancing

### Weighted Load Balancing

```go
type WeightedBackend struct {
	URL     *url.URL
	Weight  int
	Healthy bool
}

type WeightedPool struct {
	backends []WeightedBackend
	rng      *rand.Rand
	mu       sync.Mutex
}

func (p *WeightedPool) Next() *url.URL {
	p.mu.Lock()
	defer p.mu.Unlock()
	
	var totalWeight int
	for _, b := range p.backends {
		if b.Healthy {
			totalWeight += b.Weight
		}
	}
	
	// Weighted random selection
	r := p.rng.Intn(totalWeight)
	current := 0
	
	for _, b := range p.backends {
		if !b.Healthy {
			continue
		}
		current += b.Weight
		if r <= current {
			return b.URL
		}
	}
	
	return p.backends[0].URL
}
```

## Real-World Example: API Gateway

The [API Gateway in this repository](services/api-gateway/) demonstrates production patterns:

```go
// From services/api-gateway/internal/proxy/proxy.go
func NewProxy(cfg *config.Config) *Proxy {
    return &Proxy{
        client:      &http.Client{Timeout: cfg.Timeout},
        backends:    cfg.Backends,
        rateLimiter: ratelimit.New(cfg.RateLimit),
        auth:        auth.NewJWT(cfg.JWTSecret),
    }
}

func (p *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    // 1. Authenticate
    if !p.authenticate(r) {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }

    // 2. Rate limit
    if !p.rateLimiter.Allow(r) {
        http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
        return
    }

    // 3. Proxy to backend
    p.proxyRequest(w, r)
}
```

## Key Takeaways

- ✅ Rate limiting protects against abuse and manages costs
- ✅ JWT and API key auth provide secure access control
- ✅ Caching reduces latency and costs for repeated prompts
- ✅ Load balancing ensures high availability and performance

## Next Steps

→ [IO-07: Containerization with Docker](../07-containerization/README.md)

## Additional Resources

- [Kong Gateway](https://konghq.com/)
- [API Gateway Patterns](https://microservices.io/patterns/apigateway.html)
- [Rate Limiting Algorithms](https://en.wikipedia.org/wiki/Rate_limiting)
- [JWT Documentation](https://jwt.io/introduction/)
