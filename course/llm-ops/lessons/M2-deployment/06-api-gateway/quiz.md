# Quiz: API Gateway Patterns

## Question 1

What is the primary purpose of an API gateway in LLM applications?

A) To train the LLM model
B) To manage authentication, rate limiting, and request routing
C) To store embeddings
D) To preprocess training data

## Question 2

Which rate limiting algorithm is best for smooth, constant-rate traffic?

A) Token Bucket
B) Leaky Bucket
C) Fixed Window
D) Sliding Log

## Question 3

What header is commonly used to pass JWT tokens?

A) X-API-Key
B) Authorization: Bearer <token>
C) Content-Type
D) X-Auth-Token

## Question 4

In the token bucket algorithm, what happens when the bucket is empty?

A) Requests are queued until tokens are available
B) Requests are rejected immediately
C) New tokens are created instantly
D) The bucket is refilled to maximum capacity

## Question 5

What is the advantage of using Redis for caching LLM responses?

A) Lower cost than in-memory caching
B) Distributed caching across multiple instances
C) Better compression
D) Automatic model optimization

## Question 6

Which HTTP status code should be returned when rate limit is exceeded?

A) 401 Unauthorized
B) 403 Forbidden
C) 429 Too Many Requests
D) 503 Service Unavailable

## Question 7

What is the purpose of the X-RateLimit-Remaining header?

A) To authenticate the user
B) To show how many requests the user can still make
C) To specify the maximum tokens allowed
D) To identify the API version

## Question 8

In load balancing, what is a health check?

A) A test to verify user credentials
B) A periodic check to determine if a backend is operational
C) A rate limiting mechanism
D) A caching strategy

## Question 9

What type of content should NOT be cached in an LLM API?

A) Repeated identical prompts
B) Static system prompts
C) Responses containing user-specific data
D) Embeddings queries

## Question 10

Which component in the API gateway handles the actual forwarding to backend LLM servers?

A) Rate limiter
B) Authentication
C) Proxy/Load Balancer
D) Cache

## Question 11

What is the benefit of weighted load balancing?

A) Simpler implementation
B) Direct traffic based on backend capacity
C) Even distribution of requests
D) Automatic failover

## Question 12

In the context of API gateways, what does "tier" typically refer to?

A) The layer of the API (REST, GraphQL)
B) A service level (free, pro, enterprise)
C) The backend model version
D) The caching layer
