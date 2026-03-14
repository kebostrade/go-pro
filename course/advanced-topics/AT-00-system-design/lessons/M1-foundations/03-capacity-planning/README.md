# SD-03: Capacity Planning & Estimation

Learn how to estimate system capacity and plan for growth.

## Overview

Capacity planning involves predicting system resource needs based on user traffic and data volume. This lesson covers estimation techniques and calculations.

## Learning Objectives

- Estimate traffic and storage needs
- Calculate resource requirements
- Plan for growth and peak loads
- Perform back-of-the-envelope calculations

## Key Metrics

### Traffic Estimates

```
Requests per second (RPS) = Users × Requests per user per day / 86,400

Example:
- 1 million daily users
- Each user makes 100 requests/day
- RPS = 1,000,000 × 100 / 86,400 ≈ 1,157 RPS
```

### Storage Estimates

```
Storage = Users × Data per user × Growth rate × Retention period

Example:
- 1 million users
- Each user stores 1MB
- 12 months retention
- Storage = 1M × 1MB × 12 = 12TB
```

### Bandwidth Estimates

```
Bandwidth = RPS × Average request size × 2 (in/out)

Example:
- 1,000 RPS
- 10KB average request
- Bandwidth = 1,000 × 10KB × 2 = 20MB/s = 160Mbps
```

## Capacity Calculation Examples

### Web Server Capacity

```
Required Servers = (Peak RPS × Average Response Time) / Concurrent Connections

Example:
- Peak RPS: 10,000
- Average response time: 100ms
- Max connections per server: 10,000
- Servers = (10,000 × 0.1) / 10,000 = 1 server
```

### Database Capacity

```
Connections Needed = Users × Connections per user

Example:
- 1,000 concurrent users
- Each user needs 2 connections
- Pool size: 2,000 connections
- With 25 connections per server: 80 servers needed
```

### Cache Capacity

```
Cache Size = (Requests × Hit Rate × Item Size) / Compression

Example:
- 10,000 RPS
- 90% hit rate
- 1KB average item
- Cache = 10,000 × 0.9 × 1KB = 9MB/s
- Keep 1 hour: 9MB × 3,600 = ~32GB
```

## Scaling Factors

| Component | Vertical Limit | Horizontal Support |
|-----------|----------------|-------------------|
| Web Server | ~10K RPS | Yes |
| Database | Varies | Yes (read replicas) |
| Cache | Limited by memory | Yes (Redis Cluster) |
| Storage | Limited by disk | Yes |

## Back-of-the-Envelope Calculations

Use these rules of thumb:

```
1. 80/20 Rule: 80% traffic to 20% of data
2. 3x Rule: Plan for 3x current capacity
3. 10x Rule: Design for 10x for viral products
4. 95th Percentile: Design for p95, not average
```

## Example: URL Shortener

```
Requirements:
- 100M URLs
- 1B clicks/month
- 10 years retention

Calculations:
- Storage: 100M × 1KB = 100GB
- Clicks: 1B / month = ~400 RPS
- Cache: 400 × 0.5 × 1KB × 3600 = ~7GB
- Read replicas: 1 primary + 2 replicas
```

## Tools for Planning

- **Google Cloud SQL**: Use calculator
- **AWS S3**: Use pricing calculator
- **New Relic**: Performance monitoring
- **Prometheus**: Metrics collection

## Examples

See `examples/` directory for:
- `capacity_calc.go` - Calculation examples
- `scaling_estimates.md` - Detailed estimates

## Exercises

See `exercises/problems.md` for hands-on practice.

## Quiz

Test your knowledge with `quiz.md`.

## Summary

- Start with traffic and storage estimates
- Use back-of-the-envelope calculations
- Plan for 3x current capacity
- Design for peak loads, not averages

## Next Steps

Continue to [SD-04: High-Level Design](04-high-level-design/README.md)
