# Mastering Coding Interviews: A Comprehensive Guide

## Table of Contents
1. [Mindset & Strategy](#mindset--strategy)
2. [Data Structures Deep Dive](#data-structures-deep-dive)
3. [Algorithm Patterns](#algorithm-patterns)
4. [Problem-Solving Framework](#problem-solving-framework)
5. [System Design](#system-design)
6. [Behavioral Interviews](#behavioral-interviews)
7. [Practice Roadmap](#practice-roadmap)
8. [Company-Specific Prep](#company-specific-prep)
9. [Common Pitfalls](#common-pitfalls)
10. [Quick Reference](#quick-reference)

---

## Mindset & Strategy

### The Interview Game

**What They're Actually Testing:**
- Problem decomposition ability
- Communication under pressure
- Code quality and cleanliness
- Edge case awareness
- Trade-off reasoning

**What They're NOT Testing:**
- Memorized solutions
- Perfect syntax on first try
- Knowing obscure algorithms
- Speed over clarity

### Winning Mindset

```
BEFORE: "I need to solve this perfectly in 10 minutes"
AFTER:  "I need to demonstrate clear thinking and communication"
```

**Key Principles:**
1. **Think Out Loud** - Silence is your enemy
2. **Clarify First** - Never assume the problem
3. **Start Simple** - Brute force is better than nothing
4. **Iterate** - Optimize after you have something working
5. **Test** - Walk through your code with examples

---

## Data Structures Deep Dive

### Tier 1: Must Master (80% of Problems)

#### Arrays & Strings
**When to Use:** Sequential data, index-based access needed

**Key Operations:**
- Two pointers (opposite ends)
- Sliding window (contiguous subarrays)
- Prefix sums (range queries)

**Common Patterns:**
```go
// Two Sum - Two Pointers
func twoSum(nums []int, target int) []int {
    left, right := 0, len(nums)-1
    for left < right {
        sum := nums[left] + nums[right]
        if sum == target {
            return []int{left, right}
        } else if sum < target {
            left++
        } else {
            right--
        }
    }
    return nil
}

// Sliding Window - Maximum Subarray Sum (size k)
func maxSumSubarray(nums []int, k int) int {
    windowSum := 0
    for i := 0; i < k; i++ {
        windowSum += nums[i]
    }
    maxSum := windowSum

    for i := k; i < len(nums); i++ {
        windowSum = windowSum - nums[i-k] + nums[i]
        if windowSum > maxSum {
            maxSum = windowSum
        }
    }
    return maxSum
}

// Prefix Sum - Range Sum Query
type PrefixSum struct {
    prefix []int
}

func NewPrefixSum(nums []int) *PrefixSum {
    prefix := make([]int, len(nums)+1)
    for i := 0; i < len(nums); i++ {
        prefix[i+1] = prefix[i] + nums[i]
    }
    return &PrefixSum{prefix: prefix}
}

func (p *PrefixSum) RangeSum(left, right int) int {
    return p.prefix[right+1] - p.prefix[left]
}
```

**Time Complexities:**
- Access: O(1)
- Search: O(n)
- Insert/Delete: O(n)

#### Hash Maps (Maps/Dicts)
**When to Use:** Fast lookups, counting, caching

**Key Operations:**
- Check existence: O(1) average
- Count frequencies
- Map values to indices

```go
// Frequency Counter
func frequencyCounter(s string) map[rune]int {
    freq := make(map[rune]int)
    for _, ch := range s {
        freq[ch]++
    }
    return freq
}

// First Non-Repeating Character
func firstUniqueChar(s string) int {
    freq := make(map[byte]int)
    for i := 0; i < len(s); i++ {
        freq[s[i]]++
    }
    for i := 0; i < len(s); i++ {
        if freq[s[i]] == 1 {
            return i
        }
    }
    return -1
}

// Group Anagrams
func groupAnagrams(strs []string) [][]string {
    groups := make(map[string][]string)
    for _, str := range strs {
        bytes := []byte(str)
        sort.Slice(bytes, func(i, j int) bool { return bytes[i] < bytes[j] })
        key := string(bytes)
        groups[key] = append(groups[key], str)
    }
    result := make([][]string, 0, len(groups))
    for _, group := range groups {
        result = append(result, group)
    }
    return result
}
```

#### Linked Lists
**When to Use:** Frequent insertions/deletions, unknown size

```go
type ListNode struct {
    Val  int
    Next *ListNode
}

// Reverse Linked List (Iterative)
func reverseList(head *ListNode) *ListNode {
    var prev *ListNode
    curr := head
    for curr != nil {
        next := curr.Next
        curr.Next = prev
        prev = curr
        curr = next
    }
    return prev
}

// Detect Cycle (Floyd's Algorithm)
func hasCycle(head *ListNode) bool {
    if head == nil || head.Next == nil {
        return false
    }
    slow, fast := head, head.Next
    for fast != nil && fast.Next != nil {
        if slow == fast {
            return true
        }
        slow = slow.Next
        fast = fast.Next.Next
    }
    return false
}

// Find Middle Node
func findMiddle(head *ListNode) *ListNode {
    slow, fast := head, head
    for fast != nil && fast.Next != nil {
        slow = slow.Next
        fast = fast.Next.Next
    }
    return slow
}

// Merge Two Sorted Lists
func mergeTwoLists(l1, l2 *ListNode) *ListNode {
    dummy := &ListNode{}
    curr := dummy
    for l1 != nil && l2 != nil {
        if l1.Val < l2.Val {
            curr.Next = l1
            l1 = l1.Next
        } else {
            curr.Next = l2
            l2 = l2.Next
        }
        curr = curr.Next
    }
    if l1 != nil {
        curr.Next = l1
    } else {
        curr.Next = l2
    }
    return dummy.Next
}
```

#### Trees (Binary Trees & BST)
**When to Use:** Hierarchical data, sorted structure needed

```go
type TreeNode struct {
    Val   int
    Left  *TreeNode
    Right *TreeNode
}

// In-order Traversal (Recursive)
func inorderTraversal(root *TreeNode) []int {
    result := []int{}
    var inorder func(node *TreeNode)
    inorder = func(node *TreeNode) {
        if node == nil {
            return
        }
        inorder(node.Left)
        result = append(result, node.Val)
        inorder(node.Right)
    }
    inorder(root)
    return result
}

// Level-order Traversal (BFS)
func levelOrder(root *TreeNode) [][]int {
    if root == nil {
        return [][]int{}
    }
    result := [][]int{}
    queue := []*TreeNode{root}

    for len(queue) > 0 {
        levelSize := len(queue)
        level := []int{}
        for i := 0; i < levelSize; i++ {
            node := queue[0]
            queue = queue[1:]
            level = append(level, node.Val)
            if node.Left != nil {
                queue = append(queue, node.Left)
            }
            if node.Right != nil {
                queue = append(queue, node.Right)
            }
        }
        result = append(result, level)
    }
    return result
}

// Maximum Depth
func maxDepth(root *TreeNode) int {
    if root == nil {
        return 0
    }
    return 1 + max(maxDepth(root.Left), maxDepth(root.Right))
}

// Validate BST
func isValidBST(root *TreeNode) bool {
    return validate(root, nil, nil)
}

func validate(node *TreeNode, min, max *int) bool {
    if node == nil {
        return true
    }
    if min != nil && node.Val <= *min {
        return false
    }
    if max != nil && node.Val >= *max {
        return false
    }
    return validate(node.Left, min, &node.Val) &&
           validate(node.Right, &node.Val, max)
}
```

#### Stacks & Queues
**When to Use:** LIFO (stack) or FIFO (queue) order needed

```go
// Valid Parentheses
func isValidParentheses(s string) bool {
    stack := []byte{}
    pairs := map[byte]byte{')': '(', '}': '{', ']': '['}

    for i := 0; i < len(s); i++ {
        ch := s[i]
        if ch == '(' || ch == '{' || ch == '[' {
            stack = append(stack, ch)
        } else {
            if len(stack) == 0 || stack[len(stack)-1] != pairs[ch] {
                return false
            }
            stack = stack[:len(stack)-1]
        }
    }
    return len(stack) == 0
}

// Next Greater Element
func nextGreaterElement(nums []int) []int {
    n := len(nums)
    result := make([]int, n)
    for i := range result {
        result[i] = -1
    }
    stack := []int{}

    for i := 0; i < n; i++ {
        for len(stack) > 0 && nums[i] > nums[stack[len(stack)-1]] {
            idx := stack[len(stack)-1]
            stack = stack[:len(stack)-1]
            result[idx] = nums[i]
        }
        stack = append(stack, i)
    }
    return result
}

// Min Stack
type MinStack struct {
    stack    []int
    minStack []int
}

func NewMinStack() MinStack {
    return MinStack{}
}

func (s *MinStack) Push(val int) {
    s.stack = append(s.stack, val)
    if len(s.minStack) == 0 || val <= s.minStack[len(s.minStack)-1] {
        s.minStack = append(s.minStack, val)
    }
}

func (s *MinStack) Pop() {
    if s.stack[len(s.stack)-1] == s.minStack[len(s.minStack)-1] {
        s.minStack = s.minStack[:len(s.minStack)-1]
    }
    s.stack = s.stack[:len(s.stack)-1]
}

func (s *MinStack) Top() int {
    return s.stack[len(s.stack)-1]
}

func (s *MinStack) GetMin() int {
    return s.minStack[len(s.minStack)-1]
}
```

#### Heaps (Priority Queues)
**When to Use:** Need min/max element frequently, top K problems

```go
import "container/heap"

// Min Heap implementation
type MinHeap []int

func (h MinHeap) Len() int           { return len(h) }
func (h MinHeap) Less(i, j int) bool { return h[i] < h[j] }
func (h MinHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *MinHeap) Push(x interface{}) {
    *h = append(*h, x.(int))
}

func (h *MinHeap) Pop() interface{} {
    old := *h
    n := len(old)
    x := old[n-1]
    *h = old[0 : n-1]
    return x
}

// Kth Largest Element
func findKthLargest(nums []int, k int) int {
    h := &MinHeap{}
    heap.Init(h)
    for _, num := range nums {
        heap.Push(h, num)
        if h.Len() > k {
            heap.Pop(h)
        }
    }
    return (*h)[0]
}
```

### Tier 2: Important (15% of Problems)

#### Graphs
**When to Use:** Network relationships, paths, connectivity

```go
// Graph Representation
type Graph struct {
    vertices int
    adjList  map[int][]int
}

// BFS - Shortest Path
func (g *Graph) BFS(start, end int) int {
    if start == end {
        return 0
    }
    visited := make(map[int]bool)
    queue := []struct{ node, dist int }{{start, 0}}
    visited[start] = true

    for len(queue) > 0 {
        curr := queue[0]
        queue = queue[1:]
        for _, neighbor := range g.adjList[curr.node] {
            if neighbor == end {
                return curr.dist + 1
            }
            if !visited[neighbor] {
                visited[neighbor] = true
                queue = append(queue, struct{ node, dist int }{neighbor, curr.dist + 1})
            }
        }
    }
    return -1
}

// Number of Islands (2D Grid DFS)
func numIslands(grid [][]byte) int {
    if len(grid) == 0 {
        return 0
    }
    rows, cols := len(grid), len(grid[0])
    count := 0

    var dfs func(int, int)
    dfs = func(r, c int) {
        if r < 0 || r >= rows || c < 0 || c >= cols || grid[r][c] == '0' {
            return
        }
        grid[r][c] = '0'
        dfs(r+1, c)
        dfs(r-1, c)
        dfs(r, c+1)
        dfs(r, c-1)
    }

    for r := 0; r < rows; r++ {
        for c := 0; c < cols; c++ {
            if grid[r][c] == '1' {
                count++
                dfs(r, c)
            }
        }
    }
    return count
}

// Union-Find (Disjoint Set)
type UnionFind struct {
    parent []int
    rank   []int
}

func NewUnionFind(n int) *UnionFind {
    parent := make([]int, n)
    rank := make([]int, n)
    for i := range parent {
        parent[i] = i
    }
    return &UnionFind{parent, rank}
}

func (uf *UnionFind) Find(x int) int {
    if uf.parent[x] != x {
        uf.parent[x] = uf.Find(uf.parent[x])
    }
    return uf.parent[x]
}

func (uf *UnionFind) Union(x, y int) {
    px, py := uf.Find(x), uf.Find(y)
    if px == py {
        return
    }
    if uf.rank[px] < uf.rank[py] {
        uf.parent[px] = py
    } else if uf.rank[px] > uf.rank[py] {
        uf.parent[py] = px
    } else {
        uf.parent[py] = px
        uf.rank[px]++
    }
}
```

#### Tries (Prefix Trees)
**When to Use:** String prefix matching, autocomplete

```go
type TrieNode struct {
    children map[byte]*TrieNode
    isEnd    bool
}

type Trie struct {
    root *TrieNode
}

func NewTrie() *Trie {
    return &Trie{root: &TrieNode{children: make(map[byte]*TrieNode)}}
}

func (t *Trie) Insert(word string) {
    node := t.root
    for i := 0; i < len(word); i++ {
        ch := word[i]
        if _, exists := node.children[ch]; !exists {
            node.children[ch] = &TrieNode{children: make(map[byte]*TrieNode)}
        }
        node = node.children[ch]
    }
    node.isEnd = true
}

func (t *Trie) Search(word string) bool {
    node := t.root
    for i := 0; i < len(word); i++ {
        ch := word[i]
        if _, exists := node.children[ch]; !exists {
            return false
        }
        node = node.children[ch]
    }
    return node.isEnd
}

func (t *Trie) StartsWith(prefix string) bool {
    node := t.root
    for i := 0; i < len(prefix); i++ {
        ch := prefix[i]
        if _, exists := node.children[ch]; !exists {
            return false
        }
        node = node.children[ch]
    }
    return true
}
```

---

## Algorithm Patterns

### Pattern 1: Two Pointers

**When:** Array/string, need to compare/move elements

```go
// Remove Duplicates (Same Direction)
func removeDuplicates(nums []int) int {
    if len(nums) == 0 {
        return 0
    }
    slow := 0
    for fast := 1; fast < len(nums); fast++ {
        if nums[fast] != nums[slow] {
            slow++
            nums[slow] = nums[fast]
        }
    }
    return slow + 1
}

// Container With Most Water (Opposite Direction)
func maxArea(height []int) int {
    left, right := 0, len(height)-1
    maxWater := 0
    for left < right {
        width := right - left
        h := min(height[left], height[right])
        maxWater = max(maxWater, width*h)
        if height[left] < height[right] {
            left++
        } else {
            right--
        }
    }
    return maxWater
}
```

### Pattern 2: Sliding Window

**When:** Contiguous subarray/substring problems

```go
// Longest Substring Without Repeating Characters
func lengthOfLongestSubstring(s string) int {
    charIndex := make(map[byte]int)
    maxLen := 0
    left := 0

    for right := 0; right < len(s); right++ {
        ch := s[right]
        if idx, exists := charIndex[ch]; exists && idx >= left {
            left = idx + 1
        }
        charIndex[ch] = right
        maxLen = max(maxLen, right-left+1)
    }
    return maxLen
}

// Minimum Window Substring
func minWindow(s string, t string) string {
    if len(s) < len(t) {
        return ""
    }

    need := make(map[byte]int)
    for i := 0; i < len(t); i++ {
        need[t[i]]++
    }

    have := make(map[byte]int)
    needCount := len(need)
    haveCount := 0
    left := 0
    minLen := len(s) + 1
    result := ""

    for right := 0; right < len(s); right++ {
        ch := s[right]
        have[ch]++
        if need[ch] > 0 && have[ch] == need[ch] {
            haveCount++
        }

        for haveCount == needCount {
            if right-left+1 < minLen {
                minLen = right - left + 1
                result = s[left : right+1]
            }
            leftCh := s[left]
            have[leftCh]--
            if need[leftCh] > 0 && have[leftCh] < need[leftCh] {
                haveCount--
            }
            left++
        }
    }
    return result
}
```

### Pattern 3: Binary Search

**When:** Sorted array, or can apply binary search on answer space

```go
// Classic Binary Search
func binarySearch(nums []int, target int) int {
    left, right := 0, len(nums)-1
    for left <= right {
        mid := left + (right-left)/2
        if nums[mid] == target {
            return mid
        } else if nums[mid] < target {
            left = mid + 1
        } else {
            right = mid - 1
        }
    }
    return -1
}

// Find First Position
func findFirst(nums []int, target int) int {
    left, right := 0, len(nums)-1
    result := -1
    for left <= right {
        mid := left + (right-left)/2
        if nums[mid] == target {
            result = mid
            right = mid - 1
        } else if nums[mid] < target {
            left = mid + 1
        } else {
            right = mid - 1
        }
    }
    return result
}

// Search in Rotated Sorted Array
func searchRotated(nums []int, target int) int {
    left, right := 0, len(nums)-1
    for left <= right {
        mid := left + (right-left)/2
        if nums[mid] == target {
            return mid
        }
        if nums[left] <= nums[mid] {
            if target >= nums[left] && target < nums[mid] {
                right = mid - 1
            } else {
                left = mid + 1
            }
        } else {
            if target > nums[mid] && target <= nums[right] {
                left = mid + 1
            } else {
                right = mid - 1
            }
        }
    }
    return -1
}
```

### Pattern 4: DFS/BFS

```go
// DFS - Recursive
func dfs(grid [][]byte, r, c int) {
    rows, cols := len(grid), len(grid[0])
    if r < 0 || r >= rows || c < 0 || c >= cols || grid[r][c] == '0' {
        return
    }
    grid[r][c] = '0'
    dfs(grid, r+1, c)
    dfs(grid, r-1, c)
    dfs(grid, r, c+1)
    dfs(grid, r, c-1)
}

// BFS - Iterative
func bfs(grid [][]byte, startR, startC int) {
    rows, cols := len(grid), len(grid[0])
    queue := [][]int{{startR, startC}}
    directions := [][]int{{1, 0}, {-1, 0}, {0, 1}, {0, -1}}

    for len(queue) > 0 {
        r, c := queue[0][0], queue[0][1]
        queue = queue[1:]
        for _, dir := range directions {
            nr, nc := r+dir[0], c+dir[1]
            if nr >= 0 && nr < rows && nc >= 0 && nc < cols && grid[nr][nc] == '1' {
                grid[nr][nc] = '0'
                queue = append(queue, []int{nr, nc})
            }
        }
    }
}
```

### Pattern 5: Dynamic Programming

**When:** Overlapping subproblems, optimal substructure

```go
// Fibonacci (Space Optimized)
func fib(n int) int {
    if n <= 1 {
        return n
    }
    prev, curr := 0, 1
    for i := 2; i <= n; i++ {
        prev, curr = curr, prev+curr
    }
    return curr
}

// Longest Common Subsequence
func longestCommonSubsequence(text1, text2 string) int {
    m, n := len(text1), len(text2)
    dp := make([][]int, m+1)
    for i := range dp {
        dp[i] = make([]int, n+1)
    }

    for i := 1; i <= m; i++ {
        for j := 1; j <= n; j++ {
            if text1[i-1] == text2[j-1] {
                dp[i][j] = dp[i-1][j-1] + 1
            } else {
                dp[i][j] = max(dp[i-1][j], dp[i][j-1])
            }
        }
    }
    return dp[m][n]
}

// Coin Change
func coinChange(coins []int, amount int) int {
    dp := make([]int, amount+1)
    for i := range dp {
        dp[i] = amount + 1
    }
    dp[0] = 0

    for i := 1; i <= amount; i++ {
        for _, coin := range coins {
            if coin <= i {
                dp[i] = min(dp[i], dp[i-coin]+1)
            }
        }
    }
    if dp[amount] > amount {
        return -1
    }
    return dp[amount]
}

// House Robber
func rob(nums []int) int {
    if len(nums) == 0 {
        return 0
    }
    if len(nums) == 1 {
        return nums[0]
    }
    prev2, prev1 := 0, nums[0]
    for i := 1; i < len(nums); i++ {
        curr := max(prev1, prev2+nums[i])
        prev2, prev1 = prev1, curr
    }
    return prev1
}
```

### Pattern 6: Backtracking

**When:** Generate all combinations/permutations

```go
// Subsets
func subsets(nums []int) [][]int {
    result := [][]int{}
    var backtrack func(start int, subset []int)
    backtrack = func(start int, subset []int) {
        temp := make([]int, len(subset))
        copy(temp, subset)
        result = append(result, temp)

        for i := start; i < len(nums); i++ {
            subset = append(subset, nums[i])
            backtrack(i+1, subset)
            subset = subset[:len(subset)-1]
        }
    }
    backtrack(0, []int{})
    return result
}

// Permutations
func permute(nums []int) [][]int {
    result := [][]int{}
    var backtrack func(used []bool, perm []int)
    backtrack = func(used []bool, perm []int) {
        if len(perm) == len(nums) {
            temp := make([]int, len(perm))
            copy(temp, perm)
            result = append(result, temp)
            return
        }
        for i := 0; i < len(nums); i++ {
            if used[i] {
                continue
            }
            used[i] = true
            perm = append(perm, nums[i])
            backtrack(used, perm)
            perm = perm[:len(perm)-1]
            used[i] = false
        }
    }
    backtrack(make([]bool, len(nums)), []int{})
    return result
}

// Combination Sum
func combinationSum(candidates []int, target int) [][]int {
    result := [][]int{}
    var backtrack func(start, remaining int, combo []int)
    backtrack = func(start, remaining int, combo []int) {
        if remaining == 0 {
            temp := make([]int, len(combo))
            copy(temp, combo)
            result = append(result, temp)
            return
        }
        if remaining < 0 {
            return
        }
        for i := start; i < len(candidates); i++ {
            combo = append(combo, candidates[i])
            backtrack(i, remaining-candidates[i], combo)
            combo = combo[:len(combo)-1]
        }
    }
    backtrack(0, target, []int{})
    return result
}
```

### Pattern 7: Greedy

**When:** Local optimal leads to global optimal

```go
// Jump Game
func canJump(nums []int) bool {
    maxReach := 0
    for i := 0; i < len(nums); i++ {
        if i > maxReach {
            return false
        }
        maxReach = max(maxReach, i+nums[i])
    }
    return true
}

// Merge Intervals
func merge(intervals [][]int) [][]int {
    if len(intervals) <= 1 {
        return intervals
    }
    sort.Slice(intervals, func(i, j int) bool {
        return intervals[i][0] < intervals[j][0]
    })
    result := [][]int{intervals[0]}
    for i := 1; i < len(intervals); i++ {
        last := result[len(result)-1]
        if intervals[i][0] <= last[1] {
            last[1] = max(last[1], intervals[i][1])
        } else {
            result = append(result, intervals[i])
        }
    }
    return result
}

// Meeting Rooms II
func minMeetingRooms(intervals [][]int) int {
    if len(intervals) == 0 {
        return 0
    }
    starts := make([]int, len(intervals))
    ends := make([]int, len(intervals))
    for i, interval := range intervals {
        starts[i] = interval[0]
        ends[i] = interval[1]
    }
    sort.Ints(starts)
    sort.Ints(ends)

    rooms := 0
    endPtr := 0
    for i := 0; i < len(starts); i++ {
        if starts[i] < ends[endPtr] {
            rooms++
        } else {
            endPtr++
        }
    }
    return rooms
}
```

---

## Problem-Solving Framework

### The REACTO Method

**R**epeat - Restate the problem
**E**xamples - Work through examples
**A**pproach - Describe your approach
**C**ode - Write the code
**T**est - Test with examples
**O**ptimize - Discuss optimizations

### Step-by-Step Process

```
1. UNDERSTAND (5 min)
   - Restate problem in own words
   - Ask clarifying questions
   - Identify inputs/outputs
   - Discuss constraints

2. EXAMPLES (2 min)
   - Simple case
   - Edge cases
   - Large input case

3. BRUTE FORCE (5 min)
   - Describe naive solution
   - State time/space complexity
   - Don't code yet

4. OPTIMIZE (5 min)
   - Identify bottlenecks
   - Consider patterns
   - Discuss trade-offs

5. CODE (15 min)
   - Write clean code
   - Use meaningful names
   - Add comments for clarity

6. TEST (5 min)
   - Walk through with examples
   - Check edge cases
   - Fix bugs if any

7. ANALYZE (3 min)
   - Time complexity
   - Space complexity
   - Possible improvements
```

### Input Size → Algorithm Guide

```
n ≤ 30      → Backtracking, recursion
n ≤ 100     → O(n²) acceptable
n ≤ 10⁵     → O(n log n) or O(n)
n ≤ 10⁷     → O(n), must be linear
```

### Communication Templates

**Understanding:**
"I want to make sure I understand. We're given [X] and we need to [Y]. Is that correct?"

**Clarifying:**
"Should I handle the case where [edge case]?"

**Approaching:**
"Let me start with a brute force approach. Then I'll optimize it."

**Coding:**
"I'll use a hash map to store [X]. This gives us O(1) lookups."

**Testing:**
"Let me trace through with input [X]. First iteration..."

---

## System Design

### Core Concepts

| Concept | Key Points |
|---------|------------|
| Scalability | Horizontal vs Vertical, Load Balancing |
| Caching | LRU, CDN, Redis, Cache invalidation |
| Databases | SQL vs NoSQL, Indexing, Sharding |
| Microservices | API Gateway, Service Mesh, REST/gRPC |
| CAP Theorem | Consistency, Availability, Partition Tolerance |

### Design Process

```
1. Requirements Clarification (2-3 min)
   - What are we building?
   - Who are the users?
   - What are the key features?
   - Scale expectations?

2. High-Level Design (5-10 min)
   - Core components
   - Data flow
   - API design

3. Deep Dive (10-15 min)
   - Bottlenecks identification
   - Scaling strategy
   - Trade-off discussion

4. Wrap-up (2-3 min)
   - Additional features
   - Monitoring
   - Error handling
```

### Common Design Problems

| Problem | Key Technologies |
|---------|------------------|
| URL Shortener | Hash functions, Database sharding |
| Twitter Timeline | Fan-out, Caching, Pagination |
| YouTube/Netflix | CDN, Video encoding, Streaming |
| Distributed Cache | Eviction policies, Consistency |
| Rate Limiter | Token bucket, Leaky bucket |
| Search Autocomplete | Trie, Caching, Ranking |

### Example: URL Shortener

**Requirements:**
- Shorten long URLs
- Redirect short to long
- Custom short URLs
- Analytics (optional)

**Capacity:**
- 100M URLs/month
- Read:Write ratio = 100:1
- 10B URLs in 10 years

**API:**
```
POST /api/v1/shorten
  Request: { "long_url": "https://..." }
  Response: { "short_url": "http://short.url/abc123" }

GET /{short_code}
  Response: 301 Redirect to long_url
```

**Architecture:**
```
Client → Load Balancer → API Servers → Cache (Redis)
                                  ↓
                            Database (NoSQL)
```

---

## Behavioral Interviews

### The STAR Method

**S**ituation - Set the context
**T**ask - What you needed to do
**A**ction - What you actually did
**R**esult - The outcome

### Common Questions

**Leadership:**
- "Tell me about a time you led a project"
- "How do you handle disagreements?"

**Technical:**
- "Describe a challenging bug you fixed"
- "Tell me about a system you designed"

**Growth:**
- "What's a mistake you learned from?"
- "How do you stay current?"

### Example Answers

**Challenge:**
```
S: "Our payment service was timing out under load"
T: "I needed to reduce latency by 50%"
A: "I profiled the code, found N+1 queries, implemented caching"
R: "Latency dropped 70%, saved $50K/month in infrastructure costs"
```

**Conflict:**
```
S: "Team disagreed on database choice"
T: "We needed to pick one and move forward"
A: "I created a comparison matrix, ran benchmarks, facilitated discussion"
R: "Team reached consensus, project delivered on time"
```

### Red Flags to Avoid

- Blaming others
- Taking all credit
- Vague answers
- No specific examples
- Negative attitude

---

## Practice Roadmap

### 8-Week Plan

**Week 1-2: Data Structures**
- Day 1-2: Arrays & Strings
- Day 3-4: Linked Lists
- Day 5-6: Stacks & Queues
- Day 7: Review

**Week 3-4: Algorithms**
- Day 1-2: Sorting & Searching
- Day 3-4: Trees & Graphs
- Day 5-6: Dynamic Programming
- Day 7: Review

**Week 5-6: Patterns**
- Day 1: Two Pointers
- Day 2: Sliding Window
- Day 3: Binary Search
- Day 4: DFS/BFS
- Day 5: Backtracking
- Day 6: Greedy
- Day 7: Review

**Week 7: System Design**
- Day 1-2: Core Concepts
- Day 3-4: Practice Problems
- Day 5-6: Case Studies
- Day 7: Mock Interview

**Week 8: Final Prep**
- Day 1-3: Mock Interviews
- Day 4-5: Behavioral Prep
- Day 6: Rest
- Day 7: Final Review

### Daily Routine

```
Morning (1 hour):
- 1 new problem
- Review previous solutions

Evening (1-2 hours):
- 2-3 problems
- Focus on weak areas
- Write clean solutions
```

### LeetCode Problem Selection

**Easy (Foundation):**
- Two Sum
- Valid Parentheses
- Merge Two Sorted Lists
- Maximum Subarray
- Best Time to Buy/Sell Stock

**Medium (Core):**
- Longest Substring Without Repeating
- 3Sum
- Validate BST
- Number of Islands
- Coin Change
- Course Schedule

**Hard (Advanced):**
- Merge K Sorted Lists
- Trapping Rain Water
- Word Ladder
- Serialize/Deserialize Binary Tree
- LRU Cache

---

## Company-Specific Prep

### FAANG+ Patterns

**Google:**
- Focus: Algorithms, graphs, dynamic programming
- Style: Collaborative, multiple solutions
- Tip: Think out loud, consider edge cases

**Meta (Facebook):**
- Focus: Speed, multiple problems
- Style: Behavioral heavy
- Tip: Practice speed, know behavioral stories

**Amazon:**
- Focus: System design, leadership principles
- Style: Behavioral driven
- Tip: Prepare STAR stories for LPs

**Apple:**
- Focus: Practical problems, domain knowledge
- Style: Team-specific
- Tip: Know your domain deeply

**Netflix:**
- Focus: System design, culture fit
- Style: Senior-level expectations
- Tip: Read culture memo

**Microsoft:**
- Focus: Balanced, collaborative
- Style: Growth mindset
- Tip: Show learning attitude

---

## Common Pitfalls

### Technical Mistakes

1. **Jumping to Code** - Always clarify first
2. **Silence** - Think out loud
3. **No Edge Cases** - Consider empty, single, duplicates, negatives
4. **Wrong Complexity** - Double-check your math
5. **Messy Code** - Use meaningful names, format properly

### Behavioral Mistakes

1. **Vague Stories** - Be specific with metrics
2. **No Result** - Always include outcome
3. **Negative Attitude** - Focus on learning

---

## Quick Reference

### Complexity Cheatsheet

| Data Structure | Access | Search | Insert | Delete |
|----------------|--------|--------|--------|--------|
| Array | O(1) | O(n) | O(n) | O(n) |
| Linked List | O(n) | O(n) | O(1) | O(1) |
| Hash Table | N/A | O(1) | O(1) | O(1) |
| BST | O(log n) | O(log n) | O(log n) | O(log n) |
| Heap | O(1) | O(n) | O(log n) | O(log n) |

### Sorting Algorithms

| Algorithm | Best | Average | Worst | Space |
|-----------|------|---------|-------|-------|
| Merge Sort | O(n log n) | O(n log n) | O(n log n) | O(n) |
| Quick Sort | O(n log n) | O(n log n) | O(n²) | O(log n) |
| Heap Sort | O(n log n) | O(n log n) | O(n log n) | O(1) |

### Pattern Recognition

| Keywords | Likely Pattern |
|----------|----------------|
| "contiguous subarray" | Sliding Window |
| "sorted array" | Binary Search |
| "all combinations" | Backtracking |
| "shortest path" | BFS |
| "optimal" | Dynamic Programming |
| "top K" | Heap |
| "prefix" | Trie |
| "cycle" | Fast/Slow Pointer or Union-Find |

---

## Resources

### Books
- "Cracking the Coding Interview" - Gayle Laakmann McDowell
- "Elements of Programming Interviews"
- "Designing Data-Intensive Applications" - Martin Kleppmann

### Platforms
- LeetCode (leetcode.com)
- HackerRank (hackerrank.com)
- Pramp (pramp.com) - Mock interviews

### Time Investment

**Minimum Viable Prep:**
- 100-150 problems
- 4-6 weeks
- 1-2 hours/day

**Strong Prep:**
- 200-300 problems
- 8-12 weeks
- 2-3 hours/day

---

## Final Tips

### Day Before Interview
1. Rest - Don't cram
2. Prepare stories
3. Test equipment
4. Sleep 8 hours

### Day of Interview
1. Relax - You're prepared
2. Focus - One problem at a time
3. Communicate - Think out loud
4. Be Honest - If stuck, say so

### Success Formula
```
Mastery = (Problems × Quality) + (Patterns × Understanding) + (Communication × Clarity)
```

**Remember:**
- Every expert was once a beginner
- Consistency beats intensity
- The interview is as much about communication as coding

Good luck! 🚀

---
*Last Updated: 2026*
