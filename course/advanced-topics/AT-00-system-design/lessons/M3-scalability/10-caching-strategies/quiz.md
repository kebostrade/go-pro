# Quiz: Caching Strategies

## Question 1
In the cache-aside pattern, what happens during a write operation?

A) Data is written to cache only
B) Data is written to database only
C) Data is written to database, then cache is deleted
D) Data is written to cache and database simultaneously

## Question 2
What is cache stampede?

A) When the cache is too large
B) When many requests hit the database simultaneously after cache expires
C) When the cache is corrupted
D) When cache and database are out of sync

## Question 3
Which caching strategy provides the strongest consistency guarantee?

A) Cache-aside
B) Write-through
C) Write-behind
D) Refresh-ahead

## Question 4
What is the main advantage of write-behind caching?

A) Strong consistency
B) Fast write performance
C) Simpler implementation
D) Better cache hit rate

## Question 5
In a multi-level cache, which level is typically checked first?

A) Redis
B) Database
C) In-memory (local)
D) Disk

## Question 6
What is TTL in caching?

A) Time To Live - how long an item stays in cache
B) Total Transfer Load
C) Throughput to Latency ratio
D) Type, Length, Location

## Question 7
Which pattern is best for read-heavy workloads?

A) Write-through
B) Write-behind
C) Cache-aside
D) Refresh-ahead

## Question 8
What happens when you update data in the database but forget to invalidate the cache?

A) Nothing happens
B) The cache will automatically update
C) Stale data may be served from cache
D) The application will crash

## Question 9
What is the purpose of singleflight?

A) To make all requests faster
B) To prevent cache stampede by deduplicating concurrent requests
C) To compress cache data
D) To encrypt cache entries

## Question 10
Which is NOT a common cache eviction policy?

A) LRU (Least Recently Used)
B) LFU (Least Frequently Used)
C) FIFO (First In First Out)
D) Random
