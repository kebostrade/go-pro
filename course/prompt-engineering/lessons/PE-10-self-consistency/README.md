# PE-10: Self-Consistency & Ensembling

**Duration**: 2 hours
**Module**: 3 - Advanced Patterns

## Learning Objectives

- Understand self-consistency and why it improves accuracy
- Implement majority voting for reasoning tasks
- Apply ensemble techniques to prompt engineering
- Know when self-consistency helps (and when it doesn't)

## What is Self-Consistency?

Self-consistency generates multiple reasoning paths and selects the most common answer through majority voting.

```
┌─────────────────────────────────────────────────────────────┐
│                   SELF-CONSISTENCY                          │
│                                                             │
│   Question → Run 1 → Answer A                               │
│            → Run 2 → Answer A  ──► Majority: A ✓            │
│            → Run 3 → Answer B                               │
│            → Run 4 → Answer A                               │
│            → Run 5 → Answer A                               │
│                                                             │
│   4/5 said A → Final Answer: A                              │
└─────────────────────────────────────────────────────────────┘
```

## Why It Works

1. **Error independence**: Different reasoning paths make different errors
2. **Signal amplification**: Correct answers tend to cluster
3. **Noise reduction**: Random errors average out

## When to Use Self-Consistency

### ✅ Works Well For:

| Task Type | Example |
|-----------|---------|
| Math problems | "What is 23 × 47?" |
| Logic puzzles | "If A > B and B > C, is A > C?" |
| Multiple choice | "Which option is correct?" |
| Classification | "Spam or not spam?" |
| Factual QA | "What's the capital of Australia?" |

### ❌ Doesn't Help With:

| Task Type | Example |
|-----------|---------|
| Creative writing | "Write a poem" |
| Open-ended generation | "Summarize this text" |
| Subjective tasks | "Rate this movie" |
| Deterministic tasks | "Translate to Spanish" |

## Implementation

### Basic Self-Consistency

```python
import anthropic
from collections import Counter

client = anthropic.Anthropic()

def self_consistency(question, n_samples=5):
    """Generate multiple answers and return the most common one."""

    prompt = f"""
Solve this problem. Show your reasoning, then give a final answer
in the format "Answer: [your answer]"

Question: {question}
"""

    answers = []

    for i in range(n_samples):
        response = client.messages.create(
            model="claude-sonnet-4-6-20250514",
            max_tokens=1024,
            temperature=0.7,  # Some randomness for diversity
            messages=[{"role": "user", "content": prompt}]
        )

        output = response.content[0].text

        # Extract final answer
        if "Answer:" in output:
            answer = output.split("Answer:")[-1].strip().split("\n")[0]
            answers.append(answer)

    # Count and return most common
    counter = Counter(answers)
    most_common = counter.most_common(1)[0]

    return {
        "answer": most_common[0],
        "count": most_common[1],
        "confidence": most_common[1] / n_samples,
        "all_answers": answers
    }

# Example
result = self_consistency("What is 127 × 89?")
print(f"Answer: {result['answer']}")
print(f"Confidence: {result['confidence']:.0%}")
print(f"Vote distribution: {result['all_answers']}")
```

### With Chain-of-Thought

```python
def cot_self_consistency(question, n_samples=5):
    """Self-consistency with chain-of-thought reasoning."""

    prompt = f"""
Solve this step by step. Show all your work.

Question: {question}

Let's think step by step:
"""

    answers = []
    reasonings = []

    for _ in range(n_samples):
        response = client.messages.create(
            model="claude-sonnet-4-6-20250514",
            max_tokens=1024,
            temperature=0.7,
            messages=[{"role": "user", "content": prompt}]
        )

        output = response.content[0].text
        reasonings.append(output)

        # Extract numerical answer
        import re
        numbers = re.findall(r'\d+\.?\d*', output)
        if numbers:
            answers.append(numbers[-1])  # Usually last number is answer

    counter = Counter(answers)
    return counter.most_common(1)[0][0], reasonings
```

## Self-Consistency Patterns

### Pattern 1: Majority Vote

```
Run 1: Answer = 42
Run 2: Answer = 42
Run 3: Answer = 41
Run 4: Answer = 42
Run 5: Answer = 43

Result: 42 (3/5 votes)
```

### Pattern 2: Weighted Voting

Weight by confidence or reasoning quality:

```python
def weighted_vote(results):
    """Weight votes by model confidence."""
    votes = {}
    for answer, confidence in results:
        votes[answer] = votes.get(answer, 0) + confidence
    return max(votes, key=votes.get)
```

### Pattern 3: Agreement Detection

```python
def check_agreement(answers, threshold=0.8):
    """Check if answers agree above threshold."""
    counter = Counter(answers)
    top_answer, count = counter.most_common(1)[0]
    agreement = count / len(answers)

    if agreement >= threshold:
        return top_answer, True
    else:
        return None, False  # No clear consensus
```

## Advanced Ensemble Techniques

### Technique 1: Different Prompts

Use varied prompt phrasings for true ensemble:

```python
prompts = [
    "Solve: {question}",
    "Calculate the answer to: {question}",
    "Work out this problem: {question}",
    "Find the solution: {question}",
    "Determine the result: {question}"
]

answers = []
for prompt_template in prompts:
    response = run_prompt(prompt_template.format(question=question))
    answers.append(extract_answer(response))
```

### Technique 2: Different Temperatures

```python
temperatures = [0.0, 0.3, 0.5, 0.7, 1.0]

for temp in temperatures:
    response = client.messages.create(
        model="claude-sonnet-4-6-20250514",
        temperature=temp,
        ...
    )
```

### Technique 3: Different Models

```python
models = ["claude-sonnet-4-6-20250514", "claude-opus-4-6"]

for model in models:
    response = client.messages.create(model=model, ...)
```

## Sample Size Considerations

| Samples | Use Case |
|---------|----------|
| 3-5 | Quick verification, low stakes |
| 5-10 | Standard use, moderate confidence needed |
| 10-20 | High-stakes, need high confidence |
| 20+ | Critical decisions, research |

### Diminishing Returns

```
Accuracy improvement vs samples:
3 samples: +5%
5 samples: +8%
10 samples: +10%
20 samples: +11%
50 samples: +11.5%
```

## Practical Example

### Math Problem with Self-Consistency

```
Question: A store sells apples for $1.50 each. If you buy 8 apples
and pay with a $20 bill, how much change do you receive?

Running 5 times with temperature 0.7...

Run 1:
Cost = 8 × $1.50 = $12.00
Change = $20 - $12.00 = $8.00
Answer: $8.00

Run 2:
8 apples at $1.50 each
8 × 1.50 = 12
20 - 12 = 8
Answer: $8.00

Run 3:
Price per apple: $1.50
Number of apples: 8
Total cost: 1.50 × 8 = $12
Payment: $20
Change: 20 - 12 = $8
Answer: $8

Run 4:
$1.50 × 8 = $12.00
$20.00 - $12.00 = $8.00
Answer: $8.00

Run 5:
First, 1.50 times 8 is 12.
Then 20 minus 12 is 8.
Answer: $8

VOTES:
"$8.00" - 2 votes
"$8" - 2 votes
"$8.00" - 1 vote

After normalization: "8" wins with 5/5 votes
Final Answer: $8.00
```

## Handling Disagreement

When answers don't agree:

```python
def resolve_disagreement(answers, reasonings):
    """Handle cases where there's no clear consensus."""

    counter = Counter(answers)
    top = counter.most_common(2)

    # Check if it's a tie
    if len(top) > 1 and top[0][1] == top[1][1]:
        # It's a tie - need tiebreaker
        return {
            "status": "tie",
            "answers": top,
            "action": "manual_review",
            "reasonings": reasonings
        }

    # Check confidence threshold
    confidence = top[0][1] / len(answers)
    if confidence < 0.5:
        return {
            "status": "low_confidence",
            "best_answer": top[0][0],
            "confidence": confidence,
            "action": "verify_manually"
        }

    return {
        "status": "consensus",
        "answer": top[0][0],
        "confidence": confidence
    }
```

## Exercise

### Exercise 10.1: Implement Self-Consistency

Write Python code that:
1. Runs a prompt 5 times with temperature 0.7
2. Extracts numerical answers
3. Returns the majority vote with confidence

### Exercise 10.2: When to Use

For each task, decide if self-consistency would help:
1. "What is 15% of 847?"
2. "Write a haiku about programming"
3. "Classify this email as spam or not spam"
4. "Summarize this article"
5. "Solve: 3x + 7 = 22"

### Exercise 10.3: Design Ensemble

Design an ensemble approach using:
- Different prompts
- Different temperatures
- Different models

For the task: "Is this news article real or fake?"

## Key Takeaways

- ✅ Self-consistency = multiple runs + majority voting
- ✅ Works best for math, logic, classification tasks
- ✅ Use temperature > 0 for diversity
- ✅ 5-10 samples usually sufficient
- ✅ Consider different prompts/models for true ensemble

## Next Steps

→ [PE-11: Tree of Thoughts](../PE-11-tree-of-thoughts/README.md)
