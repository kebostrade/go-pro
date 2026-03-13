# PE-05: Few-Shot Learning

**Duration**: 3 hours
**Module**: 2 - Core Techniques

## Learning Objectives

- Understand when and why to use few-shot prompting
- Learn to design effective examples
- Master example selection and ordering strategies
- Handle edge cases in few-shot scenarios

## What is Few-Shot Learning?

Few-shot prompting provides examples to guide the model's behavior.

```
┌─────────────────────────────────────────────────────────────┐
│                      FEW-SHOT                               │
│                                                             │
│  Example 1:                                                 │
│    Input: "happy" → Output: "sad"                           │
│  Example 2:                                                 │
│    Input: "big" → Output: "small"                           │
│  Example 3:                                                 │
│    Input: "fast" → Output: "slow"                           │
│                                                             │
│  Now you:                                                   │
│    Input: "hot" → Output: ?                                 │
│                                                             │
│  Model infers the pattern (antonyms) from examples          │
└─────────────────────────────────────────────────────────────┘
```

## Zero-Shot vs Few-Shot vs Many-Shot

| Technique | Examples | When to Use |
|-----------|----------|-------------|
| Zero-Shot | 0 | Simple, common tasks |
| One-Shot | 1 | Format guidance needed |
| Few-Shot | 2-5 | Complex patterns, specific formats |
| Many-Shot | 6+ | Very complex, nuanced tasks |

## When to Use Few-Shot

### ✅ Use Few-Shot When:

1. **Output format is specific**
   ```
   Example: Specific JSON structure
   ```

2. **Task has subtle rules**
   ```
   Example: "Rewrite formally but keep it friendly"
   ```

3. **Domain conventions matter**
   ```
   Example: Medical coding, legal citations
   ```

4. **Zero-shot is inconsistent**
   ```
   Example: Classification with ambiguous categories
   ```

### ❌ Skip Few-Shot When:

1. Task is simple and common (translation, summarization)
2. Examples would consume too many tokens
3. Zero-shot already works well
4. You don't have good examples

## Few-Shot Prompt Structure

### Basic Structure

```
[Optional: Task Description]

Example 1:
Input: [example input 1]
Output: [example output 1]

Example 2:
Input: [example input 2]
Output: [example output 2]

Example 3:
Input: [example input 3]
Output: [example output 3]

Now process this:
Input: [actual input]
Output:
```

### Example: Sentiment Classification

```
Classify the sentiment of customer reviews as positive, negative, or neutral.

Review: "This product exceeded my expectations! Best purchase ever."
Sentiment: positive

Review: "The item arrived damaged and customer service was unhelpful."
Sentiment: negative

Review: "It's okay. Does what it says, nothing special."
Sentiment: neutral

Review: "Shipping was fast but the product quality is mediocre at best."
Sentiment:
```

## Designing Effective Examples

### Principle 1: Diversity

Include examples that cover different cases:

```
❌ Poor diversity (all similar):
Input: "I love this!" → positive
Input: "This is great!" → positive
Input: "Amazing product!" → positive

✅ Good diversity (different patterns):
Input: "I love this!" → positive
Input: "Not bad, actually quite good" → positive
Input: "After initial issues, they fixed everything perfectly" → positive
```

### Principle 2: Clarity

Make the pattern obvious:

```
❌ Unclear pattern:
Input: "Review the code" → "LGTM, minor nits"
Input: "Check this" → "Looks good"

✅ Clear pattern:
Input: "Review this PR for a login feature"
Output: "## Code Review\n- Authentication logic: ✓ Correct\n- Error handling: ⚠️ Add rate limiting\n- Tests: ✓ 95% coverage"

Input: "Review this PR for a payment feature"
Output: "## Code Review\n- Payment logic: ✓ PCI compliant\n- Error handling: ⚠️ Handle timeout cases\n- Tests: ✓ 88% coverage"
```

### Principle 3: Consistency

Use consistent formatting across all examples:

```
✅ Consistent:
Input: "Convert: 100 USD to EUR"
Output: "100 USD = 92.50 EUR (rate: 0.925)"

Input: "Convert: 50 GBP to USD"
Output: "50 GBP = 63.00 USD (rate: 1.26)"

Input: "Convert: 1000 JPY to USD"
Output: "1000 JPY = 6.70 USD (rate: 0.0067)"
```

### Principle 4: Edge Cases

Include challenging examples:

```
Input: "This product is terrible but the support team was nice."
Sentiment: mixed

Input: "" (empty review)
Sentiment: unknown

Input: "Best thing ever!!! 🎉🎉🎉 love it!!!"
Sentiment: positive
```

## Example Selection Strategies

### Strategy 1: Random Sampling

Good for diverse, representative examples.

### Strategy 2: Cluster-Based Selection

Pick examples from different clusters of your data.

```
Cluster 1: Short reviews
Cluster 2: Long detailed reviews
Cluster 3: Reviews with mixed sentiment
→ Pick 1-2 from each cluster
```

### Strategy 3: Difficulty-Based

Start simple, increase complexity:

```
Example 1: Simple case
Example 2: Medium case
Example 3: Complex/edge case
```

### Strategy 4: Similarity-Based

Pick examples most similar to input:

```
For input: "The camera quality is decent but battery life is poor"
Use examples about: product features, mixed reviews, comparisons
```

## Example Ordering

### Order Matters!

```
❌ Random order:
[Complex] → [Simple] → [Complex] → [Simple]

✅ Logical progression:
[Simple] → [Medium] → [Complex] → [Input]

✅ Most relevant last:
[General] → [Related] → [Very Similar] → [Input]
```

### Recency Bias

Models are influenced more by recent examples:

```
If your input is technical:
Place technical examples last.

If your input is casual:
Place casual examples last.
```

## Few-Shot Patterns

### Pattern 1: Input-Output Pairs

```
Input: "The cat sat."
Output: "Sat the cat." (reversed words)

Input: "Hello world."
Output: "World hello."
```

### Pattern 2: Question-Answer Pairs

```
Q: What is 2+2?
A: 4

Q: What is 10-3?
A: 7
```

### Pattern 3: Before-After Pairs

```
Before: "u r going to luv this!!!"
After: "You are going to love this!"

Before: "idk if this works tbh"
After: "I'm not sure if this works, to be honest."
```

### Pattern 4: Context-Response Pairs

```
Context: User is frustrated with slow app
Response: I understand how frustrating slow performance can be.
Let me help troubleshoot this issue.

Context: User wants to upgrade their plan
Response: Great choice! I'd be happy to help you explore our upgrade options.
```

## Advanced Few-Shot Techniques

### Technique 1: Chain-of-Thought Examples

```
Q: Roger has 5 tennis balls. He buys 2 cans of 3 balls each.
   How many balls does he have now?
A: Roger started with 5 balls.
   2 cans of 3 balls = 6 new balls.
   5 + 6 = 11 balls total.
   The answer is 11.

Q: The cafeteria had 23 apples. They used 8 for lunch and
   bought 5 more. How many apples remain?
A:
```

### Technique 2: Self-Generated Examples

Let the model generate examples first:

```
Generate 3 examples of converting informal text to formal text.
Then convert this: "gonna grab some food brb"
```

### Technique 3: Negative Examples

Show what NOT to do:

```
✓ Good: "The weather is nice." → Formal: "The weather conditions are pleasant."
✗ Bad: "The weather is nice." → Formal: "It would appear that meteorological conditions..."

Now formalize: "The meeting went well."
```

## Few-Shot with Different Formats

### JSON Output

```
Extract person information as JSON:

Text: "John Smith, 32, engineer at Google"
JSON: {"name": "John Smith", "age": 32, "occupation": "engineer", "company": "Google"}

Text: "Dr. Sarah Chen, 45, works at Stanford Medical"
JSON: {"name": "Sarah Chen", "age": 45, "occupation": "doctor", "company": "Stanford Medical"}

Text: "Meet Alex, 28-year-old freelance designer from Brooklyn"
JSON:
```

### Code Generation

```
# Convert description to Python function

Description: "A function that adds two numbers"
Code:
def add(a, b):
    return a + b

Description: "A function that checks if a number is even"
Code:
def is_even(n):
    return n % 2 == 0

Description: "A function that finds the maximum in a list"
Code:
```

## Common Pitfalls

### Pitfall 1: Overfitting to Examples

```
Problem: Model copies example patterns too rigidly

Solution: Use diverse examples with varying structures
```

### Pitfall 2: Too Many Examples

```
Problem: 10+ examples can confuse or use too many tokens

Solution: 3-5 well-chosen examples usually suffice
```

### Pitfall 3: Conflicting Examples

```
Problem: Examples send mixed signals

Solution: Ensure all examples follow the same pattern
```

## Exercise

### Exercise 5.1: Design Few-Shot Examples

Create 3 few-shot examples for:
- Task: Converting active voice to passive voice
- Include one edge case

### Exercise 5.2: Example Selection

Given these potential examples, select the best 3 for a "professional email writer" prompt:

1. "yo dude can u help" → "Could you please assist me?"
2. "need the report asap" → "I would appreciate it if you could provide the report at your earliest convenience."
3. "thx for ur help" → "Thank you for your assistance."
4. "hey" → "Hello, I hope this message finds you well."
5. "sup" → "Good morning."

### Exercise 5.3: Chain-of-Thought Few-Shot

Create a few-shot prompt that teaches the model to solve word problems step-by-step.

## Key Takeaways

- ✅ Few-shot improves performance on specific, complex tasks
- ✅ Examples should be diverse, clear, consistent, and include edge cases
- ✅ Order matters—place most relevant examples last
- ✅ 3-5 examples usually optimal
- ✅ Chain-of-thought examples improve reasoning

## Next Steps

→ [PE-06: Chain-of-Thought Prompting](../PE-06-chain-of-thought/README.md)
