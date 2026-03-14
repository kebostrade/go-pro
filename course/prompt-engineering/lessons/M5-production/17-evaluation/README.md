# PE-17: Prompt Evaluation & Metrics

**Duration**: 3 hours
**Module**: 5 - Production & Optimization

## Learning Objectives

- Define evaluation criteria for prompts
- Implement automated testing for prompts
- Measure prompt performance with metrics
- Build evaluation datasets and benchmarks

## Why Evaluate Prompts?

| Reason | Impact |
|--------|--------|
| Quality Assurance | Catch issues before production |
| Regression Prevention | Ensure updates don't break things |
| Optimization | Compare prompt versions |
| Cost Control | Identify inefficient prompts |
| User Satisfaction | Measure actual usefulness |

## Evaluation Dimensions

### 1. Accuracy

Does the prompt produce correct outputs?

```
EVALUATION CRITERIA:
- Factual correctness
- Logical consistency
- Mathematical accuracy
- Code correctness (does it run?)

MEASUREMENT:
- Exact match: Output == expected
- F1 score: Precision and recall
- Human evaluation: Expert review
```

### 2. Relevance

Does the output address the actual request?

```
EVALUATION CRITERIA:
- Answers the question asked
- Stays on topic
- Addresses all parts of request
- Doesn't add irrelevant information

MEASUREMENT:
- Relevance score (1-5)
- Coverage (% of sub-questions answered)
- Topic consistency check
```

### 3. Consistency

Does the prompt produce similar outputs for similar inputs?

```
EVALUATION CRITERIA:
- Deterministic (same input → same output at temp=0)
- Coherent style
- Consistent format
- Stable across runs

MEASUREMENT:
- Variance across N runs
- Format compliance rate
- Style consistency score
```

### 4. Efficiency

How much does it cost to run?

```
METRICS:
- Token count (input + output)
- Latency (response time)
- Cost per request
- Success rate (no errors)

TARGETS:
- Input tokens: < 2000
- Output tokens: Task-dependent
- Latency: < 5 seconds
- Success rate: > 99%
```

## Evaluation Methods

### Method 1: Unit Testing

```python
import pytest
import anthropic

client = anthropic.Anthropic()

def run_prompt(prompt, input_text):
    response = client.messages.create(
        model="claude-sonnet-4-6-20250514",
        max_tokens=1024,
        messages=[{"role": "user", "content": prompt.format(input=input_text)}]
    )
    return response.content[0].text

# Test cases
class TestSentimentPrompt:
    PROMPT = """Classify sentiment as positive, negative, or neutral.
    Output only one word.

    Text: {input}
    Sentiment:"""

    @pytest.mark.parametrize("text,expected", [
        ("I love this product!", "positive"),
        ("This is terrible.", "negative"),
        ("It's okay.", "neutral"),
        ("Best day ever!", "positive"),
        ("Worst experience of my life.", "negative"),
    ])
    def test_sentiment_classification(self, text, expected):
        result = run_prompt(self.PROMPT, text)
        assert result.strip().lower() == expected

    def test_handles_empty_input(self):
        result = run_prompt(self.PROMPT, "")
        # Should not crash, should return something
        assert result is not None
```

### Method 2: Evaluation Dataset

```python
# eval_dataset.json
{
    "evaluations": [
        {
            "id": "sentiment_001",
            "input": "I love this product!",
            "expected_output": "positive",
            "category": "sentiment",
            "difficulty": "easy"
        },
        {
            "id": "sentiment_002",
            "input": "The product is good but shipping was slow.",
            "expected_output": "mixed",
            "category": "sentiment",
            "difficulty": "medium"
        }
    ]
}

# eval_runner.py
import json

def run_evaluation(prompt, dataset_path):
    with open(dataset_path) as f:
        dataset = json.load(f)

    results = []
    for eval_case in dataset["evaluations"]:
        output = run_prompt(prompt, eval_case["input"])

        # Check if output matches expected
        passed = eval_case["expected_output"].lower() in output.lower()

        results.append({
            "id": eval_case["id"],
            "passed": passed,
            "expected": eval_case["expected_output"],
            "actual": output,
            "category": eval_case["category"]
        })

    # Calculate metrics
    total = len(results)
    passed = sum(1 for r in results if r["passed"])

    return {
        "accuracy": passed / total,
        "total": total,
        "passed": passed,
        "results": results
    }
```

### Method 3: LLM-as-Judge

```python
def llm_evaluate(prompt_output, criteria, reference=None):
    """Use another LLM to evaluate the output."""

    judge_prompt = f"""
Evaluate this AI output based on the criteria.

CRITERIA: {criteria}

{"REFERENCE ANSWER: " + reference if reference else ""}

OUTPUT TO EVALUATE:
{prompt_output}

Rate each criterion from 1-5:
- Accuracy: [1-5]
- Relevance: [1-5]
- Completeness: [1-5]
- Clarity: [1-5]

Provide brief justification for each score.
Overall score: [average]
"""

    response = client.messages.create(
        model="claude-sonnet-4-6-20250514",
        max_tokens=500,
        messages=[{"role": "user", "content": judge_prompt}]
    )

    return response.content[0].text
```

### Method 4: Human Evaluation

```python
# human_eval_interface.py
def collect_human_feedback(output_id, output, criteria):
    """Template for human evaluation interface."""
    return {
        "output_id": output_id,
        "output": output,
        "questions": [
            {
                "id": "accuracy",
                "question": "Is the output factually accurate?",
                "type": "rating",
                "scale": "1-5"
            },
            {
                "id": "helpfulness",
                "question": "How helpful is this response?",
                "type": "rating",
                "scale": "1-5"
            },
            {
                "id": "issues",
                "question": "Are there any issues with this output?",
                "type": "multiselect",
                "options": [
                    "Factual error",
                    "Missing information",
                    "Too verbose",
                    "Off-topic",
                    "Inappropriate tone"
                ]
            },
            {
                "id": "comments",
                "question": "Additional comments",
                "type": "text"
            }
        ]
    }
```

## Key Metrics

### Classification Metrics

```python
def calculate_metrics(true_labels, predicted_labels):
    """Calculate classification metrics."""
    from sklearn.metrics import accuracy_score, precision_score, recall_score, f1_score

    return {
        "accuracy": accuracy_score(true_labels, predicted_labels),
        "precision": precision_score(true_labels, predicted_labels, average='weighted'),
        "recall": recall_score(true_labels, predicted_labels, average='weighted'),
        "f1": f1_score(true_labels, predicted_labels, average='weighted')
    }
```

### Generation Metrics

```python
def calculate_generation_metrics(reference, generated):
    """Calculate text generation metrics."""

    metrics = {}

    # Exact match
    metrics["exact_match"] = reference.strip() == generated.strip()

    # BLEU score (for translation/summarization)
    from nltk.translate.bleu_score import sentence_bleu
    metrics["bleu"] = sentence_bleu([reference.split()], generated.split())

    # ROUGE score (for summarization)
    from rouge import Rouge
    rouge = Rouge()
    scores = rouge.get_scores(generated, reference)[0]
    metrics["rouge_1"] = scores["rouge-1"]["f"]
    metrics["rouge_2"] = scores["rouge-2"]["f"]
    metrics["rouge_l"] = scores["rouge-l"]["f"]

    # Token overlap
    ref_tokens = set(reference.lower().split())
    gen_tokens = set(generated.lower().split())
    metrics["token_overlap"] = len(ref_tokens & gen_tokens) / len(ref_tokens)

    return metrics
```

### Cost Metrics

```python
def calculate_cost_metrics(prompt_tokens, completion_tokens, model):
    """Calculate cost metrics for a prompt run."""

    # Pricing per 1K tokens (example prices)
    pricing = {
        "claude-sonnet-4-6-20250514": {"input": 0.003, "output": 0.015},
        "claude-opus-4-6": {"input": 0.015, "output": 0.075}
    }

    model_pricing = pricing.get(model, {"input": 0, "output": 0})

    input_cost = (prompt_tokens / 1000) * model_pricing["input"]
    output_cost = (completion_tokens / 1000) * model_pricing["output"]

    return {
        "prompt_tokens": prompt_tokens,
        "completion_tokens": completion_tokens,
        "total_tokens": prompt_tokens + completion_tokens,
        "input_cost": input_cost,
        "output_cost": output_cost,
        "total_cost": input_cost + output_cost
    }
```

## Evaluation Dashboard

```python
class PromptEvaluator:
    """Comprehensive prompt evaluation system."""

    def __init__(self, prompt, test_cases):
        self.prompt = prompt
        self.test_cases = test_cases
        self.results = []

    def run_evaluation(self):
        """Run full evaluation suite."""
        for case in self.test_cases:
            output = self.run_prompt(case["input"])
            result = {
                "id": case["id"],
                "output": output,
                "passed": self.check_output(output, case.get("expected")),
                "latency": self.measure_latency(case["input"]),
                "tokens": self.count_tokens(case["input"], output)
            }
            self.results.append(result)

        return self.generate_report()

    def generate_report(self):
        """Generate evaluation report."""
        total = len(self.results)
        passed = sum(1 for r in self.results if r["passed"])

        return {
            "summary": {
                "total_cases": total,
                "passed": passed,
                "failed": total - passed,
                "accuracy": passed / total if total > 0 else 0
            },
            "performance": {
                "avg_latency": sum(r["latency"] for r in self.results) / total,
                "avg_tokens": sum(r["tokens"]["total"] for r in self.results) / total
            },
            "details": self.results
        }
```

## Exercise

### Exercise 17.1: Create Test Suite

Create a test suite for a "summarization" prompt:
- 5 test cases with inputs and expected outputs
- Pass/fail criteria
- Edge cases to test

### Exercise 17.2: Design Evaluation Criteria

Define evaluation criteria for:
- A code generation prompt
- A translation prompt
- A creative writing prompt

### Exercise 17.3: Build Evaluation Pipeline

Write Python code that:
- Loads test cases from JSON
- Runs each test
- Calculates metrics
- Generates a report

## Key Takeaways

- ✅ Evaluate accuracy, relevance, consistency, and efficiency
- ✅ Use unit tests for deterministic prompts
- ✅ Build evaluation datasets for comprehensive testing
- ✅ Consider LLM-as-judge for subjective tasks
- ✅ Track metrics over time to catch regressions

## Next Steps

→ [PE-18: Prompt Optimization Techniques](../PE-18-optimization/README.md)
