# PE-06: Chain-of-Thought Prompting

**Duration**: 3 hours
**Module**: 2 - Core Techniques

## Learning Objectives

- Understand chain-of-thought (CoT) and why it works
- Master zero-shot CoT and few-shot CoT techniques
- Apply CoT to math, logic, and reasoning tasks
- Learn when CoT helps and when it doesn't

## What is Chain-of-Thought?

Chain-of-thought prompting encourages the model to show its reasoning step-by-step before giving a final answer.

```
┌─────────────────────────────────────────────────────────────┐
│                  WITHOUT CHAIN-OF-THOUGHT                   │
│                                                             │
│  Q: If I have 5 apples and eat 2, then buy 3 more, how     │
│     many do I have?                                         │
│  A: 6                                                       │
│  ❌ (Model might guess without reasoning)                   │
└─────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────┐
│                   WITH CHAIN-OF-THOUGHT                     │
│                                                             │
│  Q: If I have 5 apples and eat 2, then buy 3 more, how     │
│     many do I have?                                         │
│  A: Starting with 5 apples.                                 │
│     Eating 2 means: 5 - 2 = 3 apples left.                  │
│     Buying 3 more means: 3 + 3 = 6 apples.                  │
│     Final answer: 6 apples.                                 │
│  ✅ (Model shows work, catches errors)                      │
└─────────────────────────────────────────────────────────────┘
```

## Why Chain-of-Thought Works

1. **More computation** = better answers (more tokens to think)
2. **Decomposition** breaks complex problems into simple steps
3. **Error detection** happens at each step
4. **Transparency** lets humans verify reasoning

## Zero-Shot CoT

The simplest form: just add "think step by step"

### The Magic Phrase

```
[Your question]

Think step by step.
```

### Examples

```
❌ Without CoT:
Q: A bat and ball cost $1.10 total. The bat costs $1 more than
   the ball. How much does the ball cost?
A: $0.10

✅ With Zero-Shot CoT:
Q: A bat and ball cost $1.10 total. The bat costs $1 more than
   the ball. How much does the ball cost?

Let's think step by step.
A: Let's denote the ball's cost as x.
   Then the bat's cost is x + $1.
   Together: x + (x + $1) = $1.10
   2x + $1 = $1.10
   2x = $0.10
   x = $0.05

   The ball costs $0.05.
```

### Variations of Zero-Shot CoT

```
"Let's think step by step."

"Think through this carefully, one step at a time."

"Break this down into steps and solve."

"Walk through your reasoning step by step."
```

## Few-Shot CoT

Provide examples that demonstrate step-by-step reasoning.

### Structure

```
Example 1:
Q: [question]
A: [step 1]
   [step 2]
   [step 3]
   Therefore, [answer]

Example 2:
Q: [question]
A: [step 1]
   [step 2]
   Therefore, [answer]

Q: [your question]
A:
```

### Example: Math Word Problems

```
Q: The school cafeteria had 50 sandwiches. They sold 32 at lunch
   and made 15 more for the afternoon. How many sandwiches do they have?
A: Starting sandwiches: 50
   Sold at lunch: 50 - 32 = 18 remaining
   Made more: 18 + 15 = 33
   Therefore, they have 33 sandwiches.

Q: A library has 1,200 books. They lend out 145 on Monday and
   receive 32 new books on Tuesday. How many books now?
A: Starting books: 1,200
   Lent out Monday: 1,200 - 145 = 1,055
   Received Tuesday: 1,055 + 32 = 1,087
   Therefore, the library has 1,087 books.

Q: A store has 240 customers on Saturday. Sunday had 35 more
   customers than Saturday. How many customers total for the weekend?
A:
```

## When to Use CoT

### ✅ CoT Helps With:

| Task Type | Example |
|-----------|---------|
| Math word problems | "If 3 people share $45..." |
| Multi-step logic | "If A then B, and B then C..." |
| Common sense reasoning | "Can you fit an elephant in a car?" |
| Symbolic manipulation | "If X = 5 and Y = X + 3..." |
| Complex decision making | "Should I buy or rent?" |

### ❌ CoT Doesn't Help With:

| Task Type | Example |
|-----------|---------|
| Simple recall | "What's the capital of France?" |
| Single-step tasks | "Translate to Spanish" |
| Pattern matching | "Classify sentiment" |
| Creative generation | "Write a poem" |

## CoT Best Practices

### Practice 1: Encourage Explicit Calculations

```
✅ Good:
"The discount is 20% of $150.
 20% = 0.20
 0.20 × $150 = $30
 So the discount is $30."

❌ Bad:
"The discount is 20% so that's $30."
```

### Practice 2: State Assumptions

```
"I need to find the arrival time.

Assumptions:
- The train departs on time
- There are no delays
- The journey time is accurate

Step 1: Departure is 2:30 PM
Step 2: Journey is 3 hours 45 minutes
Step 3: 2:30 PM + 3h 45m = 6:15 PM

The train arrives at 6:15 PM."
```

### Practice 3: Verify at Each Step

```
Step 1: Calculate total items
  Items: 5 apples + 3 oranges = 8 items
  ✓ This looks correct

Step 2: Calculate cost per item
  Total cost: $24
  Cost per item: $24 ÷ 8 = $3
  ✓ 8 × $3 = $24, correct

Step 3: Apply discount
  Discount: 10% of $24 = $2.40
  Final price: $24 - $2.40 = $21.60
  ✓ Verified: 0.9 × $24 = $21.60
```

### Practice 4: Use Structured Format

```
Let me solve this systematically:

📋 Given:
- Original price: $200
- Discount: 15%
- Tax rate: 8%

📝 Steps:
1. Calculate discount: $200 × 0.15 = $30
2. Price after discount: $200 - $30 = $170
3. Calculate tax: $170 × 0.08 = $13.60
4. Final price: $170 + $13.60 = $183.60

✅ Answer: $183.60
```

## Advanced CoT Techniques

### Technique 1: Self-Consistency

Generate multiple reasoning paths and take majority vote:

```
Run 1:
Step 1: x = 5
Step 2: y = x + 3 = 8
Answer: 8

Run 2:
Step 1: x = 5
Step 2: y = 5 + 3
Answer: 8

Run 3:
Step 1: x = 5
Step 2: y = 5 + 3 = 9 (error)
Answer: 9

Majority vote: 8 ✓
```

### Technique 2: Least-to-Most Prompting

Break complex problems into sub-problems:

```
Q: What is the sum of all even numbers from 2 to 100?

First, let's identify what we need to solve:
1. What are the even numbers from 2 to 100?
2. What is the sum of these numbers?

Sub-problem 1: Even numbers from 2 to 100
Answer: 2, 4, 6, 8, ..., 98, 100 (50 numbers total)

Sub-problem 2: Sum of arithmetic sequence
- First term: 2
- Last term: 100
- Number of terms: 50
- Sum = (50/2) × (2 + 100) = 25 × 102 = 2,550

Answer: 2,550
```

### Technique 3: Faithful CoT

Require the model to justify each step:

```
Step 1: [claim]
Justification: [why this claim is true]

Step 2: [claim]
Justification: [why this follows from step 1]

...
```

## CoT for Code

### Debugging with CoT

```
Find the bug in this code:

def calculate_average(numbers):
    total = 0
    for n in numbers:
        total += n
    return total / len(numbers)

Let me trace through this step by step:

Step 1: Identify the function's purpose
- Calculate average of a list of numbers

Step 2: Trace the logic
- Initialize total to 0
- Loop through numbers, add each to total
- Divide total by count

Step 3: Identify potential issues
- What if numbers is empty? len(numbers) = 0
- Division by zero would occur

Step 4: Confirm the bug
- If numbers = [], then total = 0 and len = 0
- 0 / 0 raises ZeroDivisionError

Bug found: No check for empty list.
Fix: Add a check before division.
```

## Common CoT Mistakes

### Mistake 1: Skipping Steps

```
❌ Too brief:
"5 + 3 = 8, then 8 × 2 = 16. Answer: 16"

✅ Complete:
"First, I need to add 5 and 3.
 5 + 3 = 8
Next, I multiply the result by 2.
 8 × 2 = 16
Therefore, the answer is 16."
```

### Mistake 2: Not Checking Work

```
❌ No verification:
"Therefore, the answer is 16."

✅ With verification:
"Therefore, the answer is 16.
Let me verify: (5 + 3) × 2 = 8 × 2 = 16 ✓"
```

### Mistake 3: Using CoT for Simple Tasks

```
❌ Overkill:
Q: What is 2 + 2?
A: Let's think step by step.
   First, we have 2.
   Then we add 2 more.
   2 + 2 = 4.
   The answer is 4.

✅ Appropriate:
Q: What is 2 + 2?
A: 4
```

## Exercise

### Exercise 6.1: Zero-Shot CoT

Apply zero-shot CoT to improve this prompt:
```
"If a shirt costs $25 and is on sale for 20% off,
 what's the final price?"
```

### Exercise 6.2: Few-Shot CoT Examples

Create 2 few-shot CoT examples for teaching the model to solve:
- Percentage increase/decrease problems

### Exercise 6.3: CoT for Logic

Write a CoT prompt to solve this logic puzzle:
```
"Three friends (Alice, Bob, Carol) are sitting in a row.
- Alice is not at either end
- Bob is to the left of Carol
Who is sitting in the middle?"
```

## Key Takeaways

- ✅ CoT improves reasoning by showing work step-by-step
- ✅ Zero-shot CoT: just add "think step by step"
- ✅ Few-shot CoT: provide examples with reasoning
- ✅ Best for math, logic, multi-step problems
- ✅ State assumptions, verify at each step

## Next Steps

→ [PE-07: Structured Output & JSON](../PE-07-structured-output/README.md)
