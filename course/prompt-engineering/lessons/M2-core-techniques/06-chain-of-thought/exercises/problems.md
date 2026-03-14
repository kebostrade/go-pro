# Exercise: Chain-of-Thought Prompting

## Problem 1: Adding Reasoning

Take this simple math problem and add chain-of-thought reasoning:

### Basic Prompt:
```
Q: If a store has 50 apples and sells 12 apples on Monday, then 18 apples on Tuesday, how many apples remain?
A:
```

### Chain-of-Thought Version:
```
Q: If a store has 50 apples and sells 12 apples on Monday, then 18 apples on Tuesday, how many apples remain?
A: Let's think step by step.
```

### Your Task:
Complete the response with explicit steps:

Step 1: Start with 50 apples
Step 2: Subtract 12 apples sold on Monday → 50 - 12 = 38
Step 3: Subtract 18 apples sold on Tuesday → 38 - 18 = 20
Step 4: Final answer: 20 apples remain

---

## Problem 2: CoT for Logic Puzzles

Write a chain-of-thought prompt for this logic problem:

**Problem:**
"Three friends - Alice, Bob, and Charlie - have different favorite colors: red, blue, and green. Alice doesn't like green. Bob's favorite color is blue. What is each person's favorite color?"

### Your Prompt:
```
Let's think step by step.
```

---

## Problem 3: When NOT to Use CoT

For each scenario, indicate if CoT is helpful or not:

| Scenario | CoT Helpful? | Why or Why Not |
|----------|-------------|----------------|
| "What is 2+2?" | | |
| "Explain quantum entanglement" | | |
| "List all US presidents" | | |
| "If all roses are flowers and some flowers fade quickly, do all roses fade quickly?" | | |
| "Translate 'hello' to Spanish" | | |

---

## Problem 4: Implicit vs Explicit CoT

Compare these two approaches:

### Implicit CoT:
```
Solve: 15 × 4 + 23
Answer:
```

### Explicit CoT:
```
Solve: 15 × 4 + 23
Let's solve this step by step:
Step 1: 15 × 4 = 60
Step 2: 60 + 23 = 83
Final answer: 83
```

### Questions:
1. Which produces more reliable results?
2. Which uses more tokens?
3. When would you prefer each?

---

## Problem 5: CoT for Code Debugging

Write a CoT prompt to debug this code:

```python
def find_max(numbers):
    max_num = 0
    for num in numbers:
        if num > max_num:
            max_num = num
    return max_num
```

**Problem:** It fails for negative numbers.
