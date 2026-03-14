# 🚀 Quick Reference - New Go Exercises

## One-Line Commands

### Run Standalone Examples
```bash
# Concurrency
go run examples/prime_numbers.go
go run examples/producer_consumer.go
go run examples/context_timeout.go
go run examples/rate_limiter.go

# Data Structures
go run examples/queue.go
go run examples/linked_list.go
go run examples/cache.go

# Algorithms
go run examples/binary_search.go
go run examples/merge_sort.go

# Utilities
go run examples/word_counter.go
go run examples/json_parser.go
```

### Run File I/O Examples
```bash
cd examples/12.\ File\ IO/01_read_file && go run main.go
cd examples/12.\ File\ IO/02_write_file && go run main.go
cd examples/12.\ File\ IO/03_append_file && go run main.go
cd examples/12.\ File\ IO/04_read_line_by_line && go run main.go
cd examples/12.\ File\ IO/05_file_info && go run main.go
cd examples/12.\ File\ IO/06_directory_operations && go run main.go
```

### Run Tests
```bash
cd examples/13.\ Testing/01_basic_test && go test -v
cd examples/13.\ Testing/02_table_driven_tests && go test -v
cd examples/13.\ Testing/03_benchmarks && go test -bench=. -benchmem
```

### Run AWS Lambda Examples
```bash
cd 14-serverless-lambda
sam local invoke LambdaHandler --event events/example.json
sam local start-api
curl http://localhost:3000/hello
```

### Run GraphQL Examples
```bash
cd 15-graphql-gqlgen
go run server.go
# Open http://localhost:8080 for GraphQL Playground
```

### Run Practice Exercises
```bash
# Basics
go run exercises/01_basics/fizzbuzz.go
go run exercises/01_basics/reverse_string.go

# Intermediate
go run exercises/02_intermediate/url_shortener.go

# Advanced
go run exercises/03_advanced/web_crawler.go
```

### Run Solutions
```bash
go run exercises/01_basics/fizzbuzz_solution.go
go run exercises/01_basics/reverse_string_solution.go
go run exercises/02_intermediate/url_shortener_solution.go
go run exercises/03_advanced/web_crawler_solution.go
```

## File Locations

```
basic/
├── examples/
│   ├── prime_numbers.go
│   ├── word_counter.go
│   ├── json_parser.go
│   ├── rate_limiter.go
│   ├── cache.go
│   ├── queue.go
│   ├── linked_list.go
│   ├── binary_search.go
│   ├── producer_consumer.go
│   ├── context_timeout.go
│   ├── merge_sort.go
│   ├── 12. File IO/
│   │   ├── 01_read_file/
│   │   ├── 02_write_file/
│   │   ├── 03_append_file/
│   │   ├── 04_read_line_by_line/
│   │   ├── 05_file_info/
│   │   └── 06_directory_operations/
│   └── 13. Testing/
│       ├── 01_basic_test/
│       ├── 02_table_driven_tests/
│       └── 03_benchmarks/
└── exercises/
    ├── 01_basics/
    │   ├── fizzbuzz.go
    │   ├── fizzbuzz_solution.go
    │   ├── reverse_string.go
    │   └── reverse_string_solution.go
    ├── 02_intermediate/
    │   ├── url_shortener.go
    │   └── url_shortener_solution.go
    └── 03_advanced/
        ├── web_crawler.go
        └── web_crawler_solution.go
```

## Concepts Map

| Want to Learn | Run This |
|---------------|----------|
| Goroutines & Channels | `producer_consumer.go`, `prime_numbers.go` |
| Context Package | `context_timeout.go` |
| Rate Limiting | `rate_limiter.go` |
| Caching | `cache.go` |
| Data Structures | `queue.go`, `linked_list.go` |
| Algorithms | `binary_search.go`, `merge_sort.go` |
| String Processing | `word_counter.go`, `reverse_string.go` |
| JSON | `json_parser.go` |
| File I/O | All in `12. File IO/` |
| Testing | All in `13. Testing/` |
| Web Services | `url_shortener.go`, `web_crawler.go` |

## Difficulty Levels

### ⭐ Easy
- fizzbuzz.go
- reverse_string.go
- word_counter.go
- queue.go
- File I/O exercises

### ⭐⭐ Medium
- json_parser.go
- binary_search.go
- linked_list.go
- url_shortener.go

### ⭐⭐⭐ Hard
- prime_numbers.go
- rate_limiter.go
- cache.go
- context_timeout.go

### ⭐⭐⭐⭐ Expert
- producer_consumer.go
- merge_sort.go
- web_crawler.go

## Common Commands

```bash
# Run a file
go run filename.go

# Run tests
go test -v

# Run benchmarks
go test -bench=. -benchmem

# Run with race detector
go run -race filename.go

# Format code
go fmt filename.go

# Check for issues
go vet filename.go
```

## Tips

1. **Start Simple**: Begin with ⭐ exercises
2. **Read Comments**: Each file has detailed explanations
3. **Experiment**: Modify the code and see what happens
4. **Test**: Try different inputs
5. **Compare**: Check solutions after attempting

## Documentation

- **Full Guide**: `NEW_EXERCISES.md`
- **Summary**: `EXERCISES_SUMMARY.md`
- **Exercise Guide**: `exercises/README.md`
- **Main README**: `README.md`


## Specialized Topics

### AWS Lambda Serverless (14-serverless-lambda/)
- Lambda handler patterns
- API Gateway integration
- DynamoDB CRUD operations
- S3 event processing
- AWS SAM deployment

### GraphQL with gqlgen (15-graphql-gqlgen/)
- Schema design and definition
- Query and mutation resolvers
- Authentication and authorization
- Data loading patterns
- Real-time subscriptions

## Advanced Topics Map

| Want to Learn | Location |
|---------------|----------|
| Serverless APIs | `14-serverless-lambda/` |
| GraphQL APIs | `15-graphql-gqlgen/` |
| Cloud Deployment | AWS SAM template.yaml |
| API Design | GraphQL schema.graphqls |
| Event Processing | S3 events example |

## Testing Specialized Topics

### AWS Lambda
```bash
cd 14-serverless-lambda
# Install dependencies
go mod tidy

# Test locally
sam local invoke LambdaHandler --event events/example.json

# Start API Gateway
sam local start-api

# Deploy to AWS
sam deploy --guided
```

### GraphQL
```bash
cd 15-graphql-gqlgen
# Install dependencies
go mod tidy

# Run server
go run server.go

# Access Playground
open http://localhost:8080

# Example queries
curl -X POST http://localhost:8080/query \
  -H "Content-Type: application/json" \
  -d '{"query": "{ users { id username } }"}'
```
