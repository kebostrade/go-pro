# Archive - Old Examples

This directory contains old, unorganized example files that have been superseded by better implementations.

## Why Archived?

These files were loose in the root directory and:
- All marked with `//go:build ignore` (simple demo code)
- Duplicated functionality already in `pkg/` (library code)
- Better organized versions exist in `cmd/examples/`
- Cluttered the root directory structure

## Current Organization

The project now follows a clean structure:

```
fun/
├── cmd/examples/          # Organized, runnable examples
│   ├── basics/           # Language fundamentals
│   ├── algorithms/       # Algorithm demonstrations
│   ├── concurrency/      # Concurrency patterns
│   ├── datastructures/   # Data structure examples
│   └── cache/            # Caching examples
├── pkg/                  # Reusable library code
│   ├── algorithms/       # Algorithm implementations
│   ├── cache/            # Cache implementations
│   ├── concurrency/      # Concurrency utilities
│   ├── datastructures/   # Data structure implementations
│   └── utils/            # Common utilities
└── _archive/             # Old examples (this directory)
    └── old_examples/     # Superseded example files
```

## What to Use Instead

| Old File | Use Instead |
|----------|-------------|
| `binary_search.go` | `pkg/algorithms/search.go` + `cmd/examples/algorithms/search_demo.go` |
| `LIFO.go` | `pkg/datastructures/stack.go` + `cmd/examples/datastructures/stack_demo.go` |
| `linked_list.go` | `pkg/datastructures/linked_list.go` + demos |
| `queue.go` | `pkg/datastructures/queue.go` + `cmd/examples/datastructures/queue_demo.go` |
| `merge_sort.go` | `pkg/algorithms/sort.go` + `cmd/examples/algorithms/sort_demo.go` |
| `rate_limiter.go` | `pkg/concurrency/rate_limiter.go` + demos |
| `producer_consumer.go` | `pkg/concurrency/producer_consumer.go` + demos |
| `palindromes.go`, `prime_numbers.go`, `fib.go` | `pkg/algorithms/math.go` + `pkg/algorithms/strings.go` |
| `function.go`, `loop.go`, `pointer.go`, etc. | `cmd/examples/basics/*_demo.go` (comprehensive demos) |

## Preservation

These files are kept for historical reference but should not be used for new development. They may be removed in a future cleanup.

**Created:** 2025-10-09
**Reason:** Directory reorganization for better code structure
