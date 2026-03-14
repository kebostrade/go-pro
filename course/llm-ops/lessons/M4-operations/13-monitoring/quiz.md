# Quiz: Monitoring AI Systems

## Question 1

What are the three pillars of observability?

A) Metrics, Logs, Traces
B) CPU, Memory, Network
C) Input, Processing, Output
D) Latency, Cost, Quality

## Question 2

Which Prometheus metric type is best for measuring request duration?

A) Counter
B) Gauge
C) Histogram
D) Summary

## Question 3

What does P99 latency represent?

A) Average response time
B) 99% of requests are faster than this
C) The slowest request
D) 1% of requests are slower than this

## Question 4

What is the primary purpose of distributed tracing?

A) To track costs across services
B) To understand request flow across multiple services
C) To monitor CPU usage
D) To store application logs

## Question 5

Which tool is specifically designed for LLM observability?

A) Prometheus
B) Grafana
C) LangSmith
D) Jaeger

## Question 6

What does OpenTelemetry provide?

A) A time-series database
B) A visualization platform
C) Vendor-agnostic instrumentation SDKs
D) A log aggregation system

## Question 7

In the AI Agent Platform, what metric tracks token usage?

A) LLMRequestsTotal
B) LLMTokensUsed
C) LLMLatency
D) LLMCacheHits

## Question 8

What is a histogram bucket in Prometheus?

A) A time range for aggregation
B) A range of values for measuring distribution
C) A type of metric storage
D) A label for grouping

## Question 9

Why is structured logging important for LLM applications?

A) It reduces storage costs
B) It enables detailed debugging of complex request flows
C) It improves performance
D) It replaces metrics

## Question 10

What does the `rate()` function in PromQL calculate?

A) Total count per second
B) Average per second over a time range
C) Maximum value
D) Minimum value

## Question 11

What is cache hit rate?

A) Number of cache entries / total entries
B) Cache retrievals / total requests
C) Cache size / memory
D) Cache misses / total requests

## Question 12

Which alert condition would indicate a severe problem?

A) Cache hit rate < 50%
B) Error rate > 5%
C) P99 latency > 10 seconds
D) Cost > $100/day

## Question 13

What is the purpose of the `for` field in Prometheus alerts?

A) How long to wait before firing
B) How often to check
C) How long to keep alerts
D) How quickly to escalate

## Question 14

In Grafana, what is a datasource?

A) A type of visualization
B) A source of metrics (e.g., Prometheus)
C) A dashboard template
D) An alerting channel

## Question 15

What information should you include in a trace span?

A) Only the duration
B) Model name, token counts, operation type
C) Only the error message
D) Only the request ID
