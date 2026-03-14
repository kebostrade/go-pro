# PE-18: Prompt Optimization Techniques

**Duration**: 2 hours
**Module**: 5 - Production & Optimization

## Learning Objectives

- Optimize prompts for cost and performance
- Reduce token usage without sacrificing quality
- Improve latency and response times
- Apply systematic optimization methods

## Optimization Goals

| Goal | Why It Matters |
|------|----------------|
| Reduce tokens | Lower cost, faster processing |
| Improve accuracy | Better user experience |
| Reduce latency | Faster response times |
| Increase consistency | Reliable outputs |

## Token Optimization

### Remove Redundancy

```
❌ BEFORE (85 tokens):
"Please carefully read through the following text that I am providing
to you below, and then after you have finished reading it, please
provide a comprehensive summary of the main points."

✅ AFTER (25 tokens):
"Summarize the key points of this text:"
```

### Use Concise Language

```
❌ VERBOSE:
"In the event that you encounter an error or issue while processing
the user's request, you should respond by informing the user about
the problem and suggesting potential solutions."

✅ CONCISE:
"If an error occurs, inform the user and suggest solutions."
```

### Eliminate Redundant Instructions

```
❌ REDUNDANT:
"Be accurate. Be precise. Be correct in your answers. Make sure
your responses are factually correct and accurate."

✅ STREAMLINED:
"Provide accurate, factually correct answers."
```

### Use Examples Instead of Explanations

```
❌ LONG EXPLANATION (50+ tokens):
"Format your response as a bulleted list where each bullet point
starts with a dash followed by a space, and each item should be
on its own line. Make sure there is a blank line between sections..."

✅ EXAMPLE (15 tokens):
"Format:
- Item 1
- Item 2
- Item 3"
```

## Accuracy Optimization

### Add Explicit Constraints

```
❌ WITHOUT CONSTRAINTS:
"Answer questions about the product."

✅ WITH CONSTRAINTS:
"Answer questions about the product using ONLY information from
the provided documentation. If the answer is not in the docs,
respond: 'I don't have that information.'"

CONSTRAINT TYPES:
- Scope limits: "Only discuss X"
- Knowledge limits: "Use only provided context"
- Format limits: "Output only JSON"
- Length limits: "Maximum 100 words"
```

### Improve Few-Shot Examples

```
❌ POOR EXAMPLES (inconsistent):
Input: "happy" → positive
Input: "This is great!" → The sentiment is positive

✅ CONSISTENT EXAMPLES:
Input: "happy" → positive
Input: "This is great!" → positive
Input: "terrible" → negative
```

### Add Verification Steps

```
"Before responding:
1. Verify your answer addresses the question
2. Check for factual accuracy
3. Ensure output format is correct
4. Remove any unnecessary information"
```

## Latency Optimization

### Reduce Output Tokens

```
❌ SLOW (many output tokens):
"Write a detailed explanation of the concept, including background,
examples, and implications. Your response should be comprehensive..."

✅ FASTER (fewer output tokens):
"Explain in 2-3 sentences."
```

### Use Streaming

```python
# Non-streaming (wait for complete response)
response = client.messages.create(
    model="claude-sonnet-4-6-20250514",
    messages=[{"role": "user", "content": prompt}]
)
# User waits until complete...

# Streaming (show response as it generates)
with client.messages.stream(
    model="claude-sonnet-4-6-20250514",
    messages=[{"role": "user", "content": prompt}]
) as stream:
    for text in stream.text_stream:
        print(text, end="", flush=True)
# User sees response immediately
```

### Parallelize Independent Calls

```python
import asyncio

async def parallel_analysis(text):
    """Run multiple analyses in parallel."""

    # These run simultaneously
    sentiment_task = analyze_sentiment(text)
    topics_task = extract_topics(text)
    summary_task = summarize(text)

    # Wait for all to complete
    sentiment, topics, summary = await asyncio.gather(
        sentiment_task, topics_task, summary_task
    )

    return {
        "sentiment": sentiment,
        "topics": topics,
        "summary": summary
    }
```

## Systematic Optimization Process

### Step 1: Baseline Measurement

```python
def measure_baseline(prompt, test_cases):
    """Measure current performance before optimization."""
    results = {
        "total_tokens": 0,
        "total_cost": 0,
        "accuracy": 0,
        "avg_latency": 0
    }

    for case in test_cases:
        start = time.time()
        response = run_prompt(prompt, case["input"])
        latency = time.time() - start

        results["total_tokens"] += response.usage.total_tokens
        results["total_cost"] += calculate_cost(response.usage)
        results["avg_latency"] += latency

        if check_accuracy(response.content, case["expected"]):
            results["accuracy"] += 1

    results["accuracy"] /= len(test_cases)
    results["avg_latency"] /= len(test_cases)

    return results
```

### Step 2: Identify Bottlenecks

```python
def analyze_prompt(prompt, test_cases):
    """Identify optimization opportunities."""

    analysis = {
        "input_token_breakdown": {},
        "output_token_breakdown": {},
        "failure_patterns": [],
        "latency_distribution": []
    }

    # Token analysis
    prompt_parts = prompt.split("\n\n")
    for i, part in enumerate(prompt_parts):
        analysis["input_token_breakdown"][f"section_{i}"] = count_tokens(part)

    # Failure analysis
    for case in test_cases:
        response = run_prompt(prompt, case["input"])
        if not check_accuracy(response.content, case["expected"]):
            analysis["failure_patterns"].append({
                "input": case["input"],
                "expected": case["expected"],
                "actual": response.content,
                "failure_type": classify_failure(response.content, case["expected"])
            })

    return analysis
```

### Step 3: Apply Optimizations

```python
# Optimization techniques to apply

def optimize_prompt(original_prompt, analysis):
    """Apply optimizations based on analysis."""

    optimized = original_prompt

    # Remove redundant sections
    if analysis["input_token_breakdown"]["section_redundant"]:
        optimized = remove_redundancy(optimized)

    # Fix failure patterns
    for pattern in analysis["failure_patterns"]:
        if pattern["failure_type"] == "format_error":
            optimized = add_format_constraints(optimized)
        elif pattern["failure_type"] == "scope_error":
            optimized = add_scope_limits(optimized)

    # Reduce verbosity
    optimized = make_concise(optimized)

    return optimized
```

### Step 4: Validate Improvement

```python
def validate_optimization(original_prompt, optimized_prompt, test_cases):
    """Compare performance before and after optimization."""

    baseline = measure_baseline(original_prompt, test_cases)
    optimized = measure_baseline(optimized_prompt, test_cases)

    return {
        "token_reduction": baseline["total_tokens"] - optimized["total_tokens"],
        "token_reduction_pct": (1 - optimized["total_tokens"] / baseline["total_tokens"]) * 100,
        "cost_reduction": baseline["total_cost"] - optimized["total_cost"],
        "latency_improvement": baseline["avg_latency"] - optimized["avg_latency"],
        "accuracy_change": optimized["accuracy"] - baseline["accuracy"],
        "recommendation": "adopt" if optimized["accuracy"] >= baseline["accuracy"] * 0.95 else "reject"
    }
```

## Optimization Checklist

```
PROMPT OPTIMIZATION CHECKLIST

□ Remove redundant phrases
□ Replace verbose instructions with examples
□ Eliminate duplicate guidance
□ Consolidate similar instructions
□ Use shorter synonyms where possible
□ Remove unnecessary pleasantries
□ Check example consistency
□ Verify constraint clarity
□ Test with edge cases
□ Measure before/after metrics
□ Ensure accuracy maintained
□ Document changes made
```

## Common Optimization Patterns

### Pattern 1: Template Extraction

```
# Before: Repeated structure in prompt
"Format product 1 as: Name: X, Price: Y, Stock: Z"
"Format product 2 as: Name: X, Price: Y, Stock: Z"

# After: Template reference
"Format all products using template: {Name, Price, Stock}"
```

### Pattern 2: Reference Compression

```
# Before: Full documentation in prompt
[500 lines of API documentation]

# After: Summary + reference
"API Reference: See attached documentation.
Key endpoints: GET /users, POST /orders, PUT /items/{id}"
```

### Pattern 3: Conditional Instructions

```
# Before: Always include all instructions
"Translate to Spanish. Also translate to French. Also translate to German."

# After: Conditional
"Translate to: {{target_languages}}"
```

## Exercise

### Exercise 18.1: Token Reduction

Optimize this prompt to use fewer tokens while maintaining functionality:
```
"I would like you to please carefully analyze the following text
and provide me with a comprehensive and detailed summary that
captures all of the most important points and key information
contained within the text that I am about to share with you."
```

### Exercise 18.2: Accuracy Improvement

Improve this prompt to be more accurate:
```
"Answer questions about history."
```

### Exercise 18.3: Optimization Pipeline

Write a function that:
1. Takes a prompt and test cases
2. Measures baseline performance
3. Applies optimizations
4. Validates improvements
5. Returns the optimized prompt

## Key Takeaways

- ✅ Remove redundancy and verbosity
- ✅ Use examples instead of long explanations
- ✅ Add explicit constraints for accuracy
- ✅ Use streaming for perceived latency
- ✅ Measure before and after optimization

## Next Steps

→ [PE-19: Security - Injection & Leaking](../PE-19-security/README.md)
