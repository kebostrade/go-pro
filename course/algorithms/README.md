# Algorithms Course

Systematic practice and tracking for algorithms and data structures mastery.

## Structure

```
algorithms/
├── README.md                    # This file
├── ALGORITHMS_TRACKER.md        # Progress tracking template
├── PROBLEMS_TEMPLATE.go         # Go solution template
├── sessions/                    # Session notes
│   ├── session_01.md
│   ├── session_02.md
│   └── ...
├── problems/                    # Problem solutions by category
│   ├── arrays/
│   ├── linkedlists/
│   ├── trees/
│   ├── graphs/
│   ├── dp/
│   └── ...
└── notes/                       # Pattern notes
    ├── sliding_window.md
    ├── two_pointers.md
    └── ...
```

## Quick Start

### 1. Set Up Your Tracker

Copy the tracker template:
```bash
cp ALGORITHMS_TRACKER.md MY_PROGRESS.md
```

Update your progress as you solve problems.

### 2. Create Your First Solution

Use the template:
```bash
cp PROBLEMS_TEMPLATE.go problems/arrays/two_sum.go
```

### 3. Track Your Sessions

Create a session note:
```bash
cp sessions/session_template.md sessions/session_01.md
```

## Workflow

### Daily Practice Routine

1. **Pick a problem** from the tracker
2. **Time yourself** solving it
3. **Record your solution** in the appropriate category folder
4. **Update the tracker** with:
   - Status (⬜ → ✅)
   - Date completed
   - Time taken
   - Number of attempts
5. **Schedule review** using spaced repetition

### Problem-Solving Process

```
1. Understand (5 min)
   - Read problem carefully
   - Understand examples
   - Identify constraints

2. Plan (5-10 min)
   - Identify pattern
   - Write pseudocode
   - Consider edge cases

3. Implement (15-30 min)
   - Write clean code
   - Use meaningful names
   - Add comments

4. Test (5 min)
   - Run given examples
   - Test edge cases
   - Check constraints

5. Optimize (5-10 min)
   - Analyze complexity
   - Look for improvements
   - Consider alternatives
```

## Problem Categories

### Beginner (Start Here)
1. **Arrays & Strings** - Foundation for most problems
2. **Linked Lists** - Pointer manipulation
3. **Stacks & Queues** - LIFO/FIFO patterns

### Intermediate
4. **Hash Tables** - O(1) lookups
5. **Trees & BST** - Recursive thinking
6. **Heaps** - Priority-based processing

### Advanced
7. **Graphs** - Complex relationships
8. **Dynamic Programming** - Optimization
9. **Backtracking** - Exhaustive search

## Pattern Recognition

### Key Patterns to Master

| Pattern | Use Case | Example Problems |
|---------|----------|------------------|
| Sliding Window | Contiguous subarrays | Max sum subarray, Longest substring |
| Two Pointers | Sorted arrays, pairs | Two sum, Container with water |
| Fast/Slow Pointers | Cycle detection | Linked list cycle |
| Merge Intervals | Overlapping ranges | Merge meetings |
| BFS | Shortest path | Level order, Word ladder |
| DFS | Exhaustive search | Path finding, Backtracking |
| Binary Search | Sorted data | Search in rotated array |
| Top K Elements | K largest/smallest | Top K frequent |
| Trie | Prefix matching | Word search, Autocomplete |

### Pattern Decision Tree

```
Is the array sorted?
├─ Yes → Binary Search / Two Pointers
└─ No → Can we use a hash map?
    ├─ Yes → Hash Table pattern
    └─ No → Is it about contiguous elements?
        ├─ Yes → Sliding Window
        └─ No → Is it about subsets/permutations?
            ├─ Yes → Backtracking
            └─ No → Graph/tree traversal?
                ├─ Yes → BFS/DFS
                └─ No → Optimization problem?
                    ├─ Yes → Dynamic Programming
                    └─ No → Greedy / Other
```

## Review System

### Spaced Repetition Intervals

| Difficulty | Intervals |
|------------|-----------|
| Easy | 1d → 3d → 7d → 14d → 30d → 60d |
| Medium | 1d → 2d → 5d → 10d → 21d → 45d |
| Hard | 1d → 2d → 4d → 7d → 14d → 30d |

### How to Review

1. Check `ALGORITHMS_TRACKER.md` for problems due today
2. Attempt to solve without looking at solution
3. If stuck after 15 min, review your previous solution
4. Update review date based on difficulty
5. Mark as "need more practice" if you struggled

## Metrics to Track

### Per Problem
- [ ] Time to solve
- [ ] Number of attempts
- [ ] Space complexity
- [ ] Time complexity
- [ ] Alternative approaches

### Overall
- [ ] Problems solved per category
- [ ] Average time per difficulty
- [ ] Current streak
- [ ] Review accuracy

## Resources

### Recommended Order
1. **NeetCode 150** - Structured roadmap
2. **Blind 75** - Essential problems
3. **Grind 75** - Time-based study plan
4. **LeetCode Top Interview 150** - Company prep

### Learning Path

```
Week 1-2:  Arrays & Strings (25 problems)
Week 3:    Linked Lists (15 problems)
Week 4:    Stacks & Queues (12 problems)
Week 5-6:  Trees & BST (20 problems)
Week 7:    Heaps (10 problems)
Week 8:    Hash Tables (12 problems)
Week 9-11: Graphs (25 problems)
Week 12-15: Dynamic Programming (30 problems)
Week 16-17: Backtracking (15 problems)
Week 18:   Sorting & Searching (18 problems)
```

## Tips for Success

### DO
- ✅ Solve problems daily (consistency > intensity)
- ✅ Focus on understanding patterns, not memorizing
- ✅ Review problems you've solved
- ✅ Time yourself to build speed
- ✅ Write clean, readable code
- ✅ Explain your solution out loud

### DON'T
- ❌ Look at solutions too quickly
- ❌ Skip easy problems
- ❌ Ignore time/space complexity
- ❌ Forget to review
- ❌ Give up after one attempt
- ❌ Memorize without understanding

## Go-Specific Tips

### Common Packages
```go
import (
    "sort"      // Sorting slices
    "container/heap"  // Heap operations
    "container/list"  // Doubly linked list
)
```

### Useful Snippets

```go
// Sort a slice
sort.Ints(nums)
sort.Slice(items, func(i, j int) bool {
    return items[i].Val < items[j].Val
})

// Create a map
counts := make(map[int]int)

// Copy a slice
copySlice := make([]int, len(original))
copy(copySlice, original)

// Check if key exists
if val, exists := m[key]; exists {
    // use val
}
```

---

*Happy coding! Remember: consistent practice beats cramming.*
