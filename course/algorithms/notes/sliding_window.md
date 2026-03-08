# Sliding Window Pattern

## When to Use

- **Contiguous subarrays or substrings**
- **"Longest", "Shortest", "Maximum", "Minimum" with constraints**
- **Fixed or variable window size**
- **Array or string input**

## Key Indicators

- Subarray/substring problems
- "At most K", "At least K" constraints
- Need to track elements in a window
- Optimization (max/min) over contiguous elements

## Template

### Fixed Window Size
```go
func fixedWindow(arr []int, k int) int {
    // Initialize window sum/count
    windowSum := 0
    for i := 0; i < k; i++ {
        windowSum += arr[i]
    }

    result := windowSum

    // Slide the window
    for i := k; i < len(arr); i++ {
        windowSum += arr[i] - arr[i-k]  // Add new, remove old
        result = max(result, windowSum)
    }

    return result
}
```

### Variable Window Size
```go
func variableWindow(arr []int, constraint int) int {
    left := 0
    result := 0
    windowState := 0  // Could be sum, count, map, etc.

    for right := 0; right < len(arr); right++ {
        // Expand: add arr[right] to window
        windowState += arr[right]

        // Shrink: while constraint violated
        for windowState > constraint {
            windowState -= arr[left]
            left++
        }

        // Update result (window is now valid)
        result = max(result, right-left+1)
    }

    return result
}
```

### With Frequency Map
```go
func windowWithMap(s string, k int) int {
    left := 0
    maxLen := 0
    maxCount := 0
    freq := make(map[byte]int)

    for right := 0; right < len(s); right++ {
        freq[s[right]]++
        maxCount = max(maxCount, freq[s[right]])

        // Window size - maxCount = number of chars to replace
        // If > k, shrink window
        for (right-left+1)-maxCount > k {
            freq[s[left]]--
            left++
        }

        maxLen = max(maxLen, right-left+1)
    }

    return maxLen
}
```

## Common Problems

| Problem | Type | LeetCode |
|---------|------|----------|
| Maximum Sum Subarray of Size K | Fixed | - |
| Longest Substring Without Repeating | Variable | #3 |
| Minimum Window Substring | Variable | #76 |
| Longest Repeating Character Replacement | Variable | #424 |
| Permutation in String | Fixed | #567 |
| Find All Anagrams | Fixed | #438 |
| Sliding Window Maximum | Fixed + Deque | #239 |
| Fruit Into Baskets | Variable | #904 |
| Longest Substring with At Most K Distinct | Variable | #340 |

## Complexity

- **Time**: O(n) - each element visited at most twice (left and right pointers)
- **Space**: O(k) where k is size of character set or O(1) for simple tracking

## Common Mistakes

1. **Off-by-one errors**: Remember window is `[left, right]` inclusive
2. **Wrong shrink condition**: Only shrink when constraint violated
3. **Forgetting to update result**: Must update inside the loop
4. **Incorrect window state**: Track the right state (sum, count, frequency)

## Practice Progress

| Problem | Status | Date | Notes |
|---------|--------|------|-------|
| Max Sum Subarray Size K | ⬜ | | |
| Longest Substring No Repeat | ⬜ | | |
| Min Window Substring | ⬜ | | |
| Longest Repeating Char | ⬜ | | |
| Permutation in String | ⬜ | | |
| Find All Anagrams | ⬜ | | |
