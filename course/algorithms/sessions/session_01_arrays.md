# Session 1: Arrays & Hash Maps

**Date**: _Fill in_
**Duration**: 60 min
**Focus**: Array fundamentals and hash map patterns

## Goals

- [ ] Complete 3 array problems
- [ ] Understand hash map pattern for lookups
- [ ] Track time for each problem

## Problems

### Problem 1: Two Sum
- **Difficulty**: Easy
- **Pattern**: Hash Map
- **LeetCode**: #1
- **Status**: ⬜ Not started
- **Time**: _fill in_
- **Notes**:
  - Approach: Use hash map to store seen values
  - Complexity: Time O(n), Space O(n)
  - Key insight: Trade space for time - O(n²) → O(n)

### Problem 2: Contains Duplicate
- **Difficulty**: Easy
- **Pattern**: Hash Set
- **LeetCode**: #217
- **Status**: ⬜ Not started
- **Time**: _fill in_
- **Notes**:
  - Approach: Use set to track seen elements
  - Complexity: Time O(n), Space O(n)
  - Key insight: Set is perfect for "seen before" checks

### Problem 3: Best Time to Buy and Sell Stock
- **Difficulty**: Easy
- **Pattern**: Single Pass
- **LeetCode**: #121
- **Status**: ⬜ Not started
- **Time**: _fill in_
- **Notes**:
  - Approach: Track min price, calculate max profit at each step
  - Complexity: Time O(n), Space O(1)
  - Key insight: Only need to track min so far, not all prices

## Key Learnings

1. Hash maps reduce O(n²) to O(n) for lookup problems
2. Sets are ideal for "duplicate" or "seen before" checks
3. Sometimes O(1) space is possible with clever tracking

## Struggles

1. _Fill in after completing_
2. _Fill in after completing_

## Patterns Identified

| Pattern | Problems | Confidence |
|---------|----------|------------|
| Hash Map Lookup | Two Sum | 1-5 |
| Hash Set | Contains Duplicate | 1-5 |
| Single Pass Tracking | Buy/Sell Stock | 1-5 |

## Next Session Plan

- Focus on: Two Pointers pattern
- Problems to retry: _if any_
- Concepts to review: _if any_

---

## Code Solutions

### Two Sum
```go
func twoSum(nums []int, target int) []int {
    seen := make(map[int]int)

    for i, num := range nums {
        complement := target - num
        if j, exists := seen[complement]; exists {
            return []int{j, i}
        }
        seen[num] = i
    }

    return []int{}
}
```

### Contains Duplicate
```go
func containsDuplicate(nums []int) bool {
    seen := make(map[int]bool)

    for _, num := range nums {
        if seen[num] {
            return true
        }
        seen[num] = true
    }

    return false
}
```

### Best Time to Buy and Sell Stock
```go
func maxProfit(prices []int) int {
    minPrice := prices[0]
    maxProfit := 0

    for _, price := range prices {
        if price < minPrice {
            minPrice = price
        } else {
            profit := price - minPrice
            if profit > maxProfit {
                maxProfit = profit
            }
        }
    }

    return maxProfit
}
```
