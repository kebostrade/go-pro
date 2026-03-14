# Two Pointers Pattern

## When to Use

- **Sorted arrays** (or can be sorted)
- **Finding pairs** with specific sum/property
- **Palindromes**
- **Removing duplicates** from sorted array
- **Merging** sorted arrays

## Types

### 1. Opposite Direction (Start/End)
Pointers move toward each other

```go
func oppositeDirection(arr []int, target int) []int {
    left, right := 0, len(arr)-1

    for left < right {
        sum := arr[left] + arr[right]

        if sum == target {
            return []int{left, right}
        } else if sum < target {
            left++  // Need larger sum
        } else {
            right-- // Need smaller sum
        }
    }

    return []int{-1, -1}
}
```

### 2. Same Direction (Fast/Slow)
Both move forward at different speeds

```go
func sameDirection(arr []int) int {
    slow, fast := 0, 0

    for fast < len(arr) {
        // Process slow pointer
        if arr[fast] != arr[slow] {
            slow++
            arr[slow] = arr[fast]
        }
        fast++
    }

    return slow + 1  // Length of unique elements
}
```

### 3. Two Arrays
One pointer per array

```go
func mergeSorted(a, b []int) []int {
    result := make([]int, 0, len(a)+len(b))
    i, j := 0, 0

    for i < len(a) && j < len(b) {
        if a[i] <= b[j] {
            result = append(result, a[i])
            i++
        }} else {
            result = append(result, b[j])
            j++
        }
    }

    // Add remaining elements
    result = append(result, a[i:]...)
    result = append(result, b[j:]...)

    return result
}
```

## Common Problems

### Opposite Direction

| Problem | Description | LeetCode |
|---------|-------------|----------|
| Two Sum II | Find pair with target sum | #167 |
| 3Sum | Triplets summing to zero | #15 |
| Container With Most Water | Max area between lines | #11 |
| Trapping Rain Water | Water trapped between bars | #42 |
| Valid Palindrome | Check if palindrome | #125 |

### Same Direction

| Problem | Description | LeetCode |
|---------|-------------|----------|
| Remove Duplicates | In-place deduplication | #26 |
| Remove Element | Remove all occurrences | #27 |
| Move Zeroes | Move zeros to end | #283 |

### Two Arrays

| Problem | Description | LeetCode |
|---------|-------------|----------|
| Merge Sorted Arrays | Merge in-place | #88 |
| Intersection of Arrays | Common elements | #349 |
| Backspace String Compare | Compare after edits | #844 |

## Template Variations

### 3Sum (Nested Two Pointers)
```go
func threeSum(nums []int) [][]int {
    sort.Ints(nums)
    result := [][]int{}

    for i := 0; i < len(nums)-2; i++ {
        // Skip duplicates
        if i > 0 && nums[i] == nums[i-1] {
            continue
        }

        // Two pointers for remaining two numbers
        left, right := i+1, len(nums)-1
        for left < right {
            sum := nums[i] + nums[left] + nums[right]
            if sum == 0 {
                result = append(result, []int{nums[i], nums[left], nums[right]})
                left++
                right--
                // Skip duplicates
                for left < right && nums[left] == nums[left-1] {
                    left++
                }
                for left < right && nums[right] == nums[right+1] {
                    right--
                }
            } else if sum < 0 {
                left++
            } else {
                right--
            }
        }
    }

    return result
}
```

### Palindrome Check
```go
func isPalindrome(s string) bool {
    left, right := 0, len(s)-1

    for left < right {
        // Skip non-alphanumeric
        for left < right && !isAlphaNum(s[left]) {
            left++
        }
        for left < right && !isAlphaNum(s[right]) {
            right--
        }

        if toLower(s[left]) != toLower(s[right]) {
            return false
        }
        left++
        right--
    }

    return true
}
```

## Complexity

- **Time**: O(n) for single pass, O(n log n) if sorting needed
- **Space**: O(1) in-place, O(n) if creating result

## Key Insights

1. **Sorting enables two pointers**: Many problems become trivial after sorting
2. **Opposite for pairs**: Finding pairs that satisfy condition
3. **Same for in-place modification**: Slow tracks position, fast explores
4. **Skip duplicates**: When sorted, skip same values to avoid duplicates

## Common Mistakes

1. **Not sorting first**: Two sum requires sorted array
2. **Index out of bounds**: Check `left < right` not `left <= right`
3. **Missing duplicates**: Skip same values in sorted arrays
4. **Wrong movement**: Move correct pointer based on condition

## Practice Progress

| Problem | Status | Date | Notes |
|---------|--------|------|-------|
| Two Sum II | ⬜ | | |
| 3Sum | ⬜ | | |
| Container With Water | ⬜ | | |
| Trapping Rain | ⬜ | | |
| Valid Palindrome | ⬜ | | |
| Remove Duplicates | ⬜ | | |
| Move Zeroes | ⬜ | | |
| Merge Sorted | ⬜ | | |
