# Quiz: Prompt Caching Strategies

## Question 1

What does KV cache store in transformer models?

A) The query matrices
B) The key and value matrices from attention
C) The model weights
D) The embeddings

## Question 2

What is the main advantage of semantic caching over exact-match caching?

A) It's faster to implement
B) It caches responses for similar but not identical prompts
C) It uses less memory
D) It doesn't require Redis

## Question 3

What similarity threshold is typical for semantic caching?

A) 0.50
B) 0.70
C) 0.85
D) 0.99

## Question 4

Which Redis data structure is best for storing cached embeddings?

A) String
B) List
C) Hash
D) Sorted Set

## Question 5

What is the primary benefit of KV cache in LLM inference?

A) Reduces API costs
B) Enables longer context windows
C) Reduces per-token generation latency
D) Improves response quality

## Question 6

In the AI Agent Platform, what cache configuration option enables semantic matching?

A) EnableSemantic
B) SemanticMatch
C) EnableSimilarity
D) SemanticCache

## Question 7

What is a good target for cache hit rate?

A) > 10%
B) > 25%
C) > 40%
D) > 90%

## Question 8

Which caching strategy is best for a FAQ bot with identical questions?

A) Semantic caching
B) Exact match caching
C) No caching needed
D) KV cache only

## Question 9

What does TTL stand for in caching?

A) Token Transfer Time
B) Time To Live
C) Total Token Length
D) Throughput To Latency

## Question 10

What is PagedAttention in vLLM?

A) A pagination API
B) Memory-efficient KV cache management
C) A caching strategy
D) A load balancing technique

## Question 11

Which metric measures cost reduction from caching?

A) Hit rate
B) Latency improvement
C) Cost savings percentage
D) Memory usage

## Question 12

Why should system prompts be stable for effective caching?

A) They add tokens but not value
B) Any change invalidates the entire cache
C) They don't matter for caching
D) System prompts shouldn't be cached

## Question 13

What is the main trade-off with semantic caching?

A) Memory vs. speed
B) Accuracy vs. similarity threshold
C) Latency vs. cost savings
D) Implementation complexity vs. savings

## Question 14

Which Redis command is used to set an expiration on a key?

A) Redis.EXPIRE
B) SETEX
C) EXPIRE
D) TTL

## Question 15

What is the expected cost savings from semantic caching for typical workloads?

A) 5-15%
B) 20-40%
C) 50-80%
D) 90-100%
