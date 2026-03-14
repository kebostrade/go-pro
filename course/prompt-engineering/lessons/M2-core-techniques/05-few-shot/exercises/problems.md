# Exercise: Few-Shot Learning

## Problem 1: Creating Effective Examples

Create 3-shot examples for sentiment analysis:

### Task: Classify the sentiment as POSITIVE or NEGATIVE

**Examples to include:**
```
Input: I loved this movie! It was amazing!
Output: POSITIVE

Input: This product is terrible. Worst purchase ever.
Output: NEGATIVE

Input: The service was okay, nothing special.
Output: NEUTRAL
```

### Your Task:
Now test this prompt with a new input:
```
Input: I'm really disappointed with the quality.
Output: ?
```

---

## Problem 2: Example Selection

You need a prompt to classify customer support tickets.

### Step 1: Select 3 examples that cover:
- A simple billing question
- A technical bug report
- A feature request

### Step 2: Write the complete prompt with your examples:

```
<!-- Write your prompt here -->
```

---

## Problem 3: Formatting Examples

Re-format the following examples for a math tutor:

**Raw data:**
- Q: What is 2+2? A: 4
- Q: What is 5+3? A: 8
- Q: What is 10+7? A: 17

### Choose a format and write the prompt:
- Option A: Q: ... A: ...
- Option B: Input: ... Output: ...
- Option C: JSON

---

## Problem 4: Counterfactual Examples

Explain why including both positive AND negative examples matters:

```
Positive Example: "Great product!" → POSITIVE
Negative Example: "Bad service." → NEGATIVE
```

Why not just include POSITIVE examples?

---

## Problem 5: Optimal Number of Examples

For each scenario, recommend 1-shot, 3-shot, or 5-shot:

| Scenario | Recommendation | Rationale |
|----------|---------------|-----------|
| Simple yes/no classification | | |
| Nuanced emotion detection | | |
| Identifying rare edge cases | | |
