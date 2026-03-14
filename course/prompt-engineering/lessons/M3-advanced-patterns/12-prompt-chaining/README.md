# PE-12: Prompt Chaining & Workflows

**Duration**: 3 hours
**Module**: 3 - Advanced Patterns

## Learning Objectives

- Design multi-step prompt pipelines
- Implement sequential and parallel workflows
- Handle state and context between prompts
- Build robust, fault-tolerant chains

## What is Prompt Chaining?

Prompt chaining connects multiple prompts in sequence, where each prompt's output becomes the next prompt's input.

```
┌─────────────────────────────────────────────────────────────┐
│                    PROMPT CHAIN                             │
│                                                             │
│  Input → [Prompt 1] → Output 1 → [Prompt 2] → Output 2 →   │
│          → [Prompt 3] → Output 3 → Final Result             │
│                                                             │
│  Example:                                                   │
│  Article → [Summarize] → Summary → [Extract Topics] →       │
│            Topics → [Generate Tweet] → Tweet                 │
└─────────────────────────────────────────────────────────────┘
```

## Why Chain Prompts?

1. **Complexity management**: Break big tasks into smaller ones
2. **Quality improvement**: Each step can be optimized
3. **Reusability**: Use same steps in different workflows
4. **Debugging**: Identify which step fails
5. **Cost control**: Use cheaper models for simple steps

## Chain Patterns

### Pattern 1: Sequential Chain

Simple linear pipeline:

```python
def sequential_chain(article):
    # Step 1: Summarize
    summary = prompt(f"Summarize: {article}")

    # Step 2: Extract key points
    key_points = prompt(f"Extract 3 key points from: {summary}")

    # Step 3: Generate social posts
    posts = prompt(f"Write social media posts about: {key_points}")

    return posts
```

### Pattern 2: Parallel Chain

Process multiple aspects simultaneously:

```python
async def parallel_chain(article):
    # Run these in parallel
    summary = await prompt_async(f"Summarize: {article}")
    sentiment = await prompt_async(f"Analyze sentiment: {article}")
    topics = await prompt_async(f"Extract topics: {article}")

    # Combine results
    return {
        "summary": summary,
        "sentiment": sentiment,
        "topics": topics
    }
```

### Pattern 3: Branching Chain

Conditional paths based on content:

```python
def branching_chain(article):
    # Classify first
    category = prompt(f"Classify this article: {article}")

    if category == "technical":
        return technical_pipeline(article)
    elif category == "news":
        return news_pipeline(article)
    else:
        return general_pipeline(article)
```

### Pattern 4: Map-Reduce Chain

Process multiple items then combine:

```python
def map_reduce_chain(articles):
    # Map: Process each article
    summaries = [prompt(f"Summarize: {a}") for a in articles]

    # Reduce: Combine summaries
    combined = prompt(f"Synthesize these summaries: {summaries}")

    return combined
```

## Chain Design Principles

### Principle 1: Clear Interfaces

Define what each step expects and produces:

```python
@dataclass
class StepInput:
    text: str
    metadata: dict

@dataclass
class StepOutput:
    result: str
    confidence: float

def summarize_step(input: StepInput) -> StepOutput:
    # Clear input/output contract
    ...
```

### Principle 2: Single Responsibility

Each step does ONE thing well:

```
❌ Bad: "Summarize and extract entities and translate to Spanish"
✅ Good:
  Step 1: Summarize
  Step 2: Extract entities
  Step 3: Translate to Spanish
```

### Principle 3: State Management

Pass context through the chain:

```python
class ChainState:
    def __init__(self, input_data):
        self.data = input_data
        self.history = []
        self.metadata = {}

    def add_result(self, step_name, result):
        self.history.append({"step": step_name, "result": result})
        self.data = result

    def get_context(self, n_steps=3):
        """Get context from last n steps."""
        return self.history[-n_steps:]
```

## Complete Chain Example

### Document Processing Pipeline

```python
import anthropic
from typing import Dict, List
from dataclasses import dataclass

@dataclass
class DocumentPipeline:
    """Process documents through multiple stages."""

    client: anthropic.Anthropic
    model: str = "claude-sonnet-4-6-20250514"

    def run(self, document: str) -> Dict:
        """Execute the full pipeline."""
        state = {
            "original": document,
            "steps": []
        }

        # Step 1: Clean and normalize
        cleaned = self._clean(document)
        state["steps"].append({"name": "clean", "output": cleaned[:100]})

        # Step 2: Extract structure
        structure = self._extract_structure(cleaned)
        state["steps"].append({"name": "structure", "output": structure})

        # Step 3: Summarize
        summary = self._summarize(cleaned, structure)
        state["steps"].append({"name": "summary", "output": summary})

        # Step 4: Extract entities
        entities = self._extract_entities(cleaned)
        state["steps"].append({"name": "entities", "output": entities})

        # Step 5: Generate output
        final = self._format_output(summary, entities, structure)
        state["result"] = final

        return state

    def _clean(self, text: str) -> str:
        """Clean and normalize text."""
        prompt = f"""Clean this text by:
- Removing artifacts and formatting issues
- Normalizing whitespace
- Fixing obvious typos
- Preserve all meaningful content

Text: {text}

Cleaned text:"""

        response = self.client.messages.create(
            model=self.model,
            max_tokens=2000,
            messages=[{"role": "user", "content": prompt}]
        )
        return response.content[0].text

    def _extract_structure(self, text: str) -> Dict:
        """Extract document structure."""
        prompt = f"""Analyze the structure of this document and return JSON:
{{
  "title": "main title or topic",
  "sections": ["list of main sections"],
  "type": "article|report|email|other"
}}

Document: {text[:2000]}

JSON:"""

        response = self.client.messages.create(
            model=self.model,
            max_tokens=500,
            messages=[{"role": "user", "content": prompt}]
        )
        # Parse JSON...
        return response.content[0].text

    def _summarize(self, text: str, structure: Dict) -> str:
        """Generate summary."""
        prompt = f"""Summarize this document in 3 sentences.

Structure: {structure}

Document: {text[:3000]}

Summary:"""

        response = self.client.messages.create(
            model=self.model,
            max_tokens=200,
            messages=[{"role": "user", "content": prompt}]
        )
        return response.content[0].text

    def _extract_entities(self, text: str) -> List[Dict]:
        """Extract named entities."""
        prompt = f"""Extract entities as JSON array:
[{{"type": "person|org|location|date", "name": "value"}}]

Text: {text[:2000]}

Entities:"""

        response = self.client.messages.create(
            model=self.model,
            max_tokens=500,
            messages=[{"role": "user", "content": prompt}]
        )
        return response.content[0].text

    def _format_output(self, summary, entities, structure) -> str:
        """Format final output."""
        return f"""
# Document Analysis

## Summary
{summary}

## Structure
Type: {structure.get('type', 'Unknown')}
Sections: {', '.join(structure.get('sections', []))}

## Key Entities
{entities}
"""
```

## Error Handling in Chains

### Retry Logic

```python
from tenacity import retry, stop_after_attempt, wait_exponential

@retry(stop=stop_after_attempt(3), wait=wait_exponential(multiplier=1, min=2, max=10))
def robust_prompt(prompt_text, client):
    """Prompt with automatic retry."""
    try:
        response = client.messages.create(...)
        return response.content[0].text
    except anthropic.RateLimitError:
        raise  # Let retry handle it
    except anthropic.APIError as e:
        raise  # Let retry handle it
```

### Fallback Strategies

```python
def prompt_with_fallback(prompt_text, primary_model, fallback_model):
    """Try primary model, fall back to secondary."""
    try:
        return run_prompt(prompt_text, primary_model)
    except Exception as e:
        print(f"Primary failed: {e}, trying fallback")
        return run_prompt(prompt_text, fallback_model)
```

### Validation Gates

```python
def validated_step(input_data, validator_fn):
    """Run step with validation."""
    result = run_prompt(input_data)

    if not validator_fn(result):
        raise ValidationError(f"Invalid output: {result[:100]}")

    return result

def validate_summary(summary):
    """Check summary meets requirements."""
    return (
        len(summary) > 50 and
        len(summary) < 500 and
        summary.endswith('.')
    )
```

## Chain Orchestration

### Using LangChain

```python
from langchain.chains import SequentialChain, LLMChain
from langchain_anthropic import ChatAnthropic

llm = ChatAnthropic(model="claude-sonnet-4-6-20250514")

# Define chains
summarize_chain = LLMChain(
    llm=llm,
    prompt=PromptTemplate(
        template="Summarize: {text}",
        input_variables=["text"]
    ),
    output_key="summary"
)

extract_chain = LLMChain(
    llm=llm,
    prompt=PromptTemplate(
        template="Extract key points from: {summary}",
        input_variables=["summary"]
    ),
    output_key="key_points"
)

# Combine
full_chain = SequentialChain(
    chains=[summarize_chain, extract_chain],
    input_variables=["text"],
    output_variables=["summary", "key_points"]
)

result = full_chain({"text": long_article})
```

## Monitoring & Debugging

### Logging

```python
import logging

logging.basicConfig(level=logging.INFO)
logger = logging.getLogger("prompt_chain")

def logged_step(step_name, prompt, input_data):
    logger.info(f"Starting step: {step_name}")
    logger.debug(f"Input: {str(input_data)[:200]}")

    start_time = time.time()
    result = run_prompt(prompt.format(**input_data))
    elapsed = time.time() - start_time

    logger.info(f"Completed {step_name} in {elapsed:.2f}s")
    logger.debug(f"Output: {str(result)[:200]}")

    return result
```

### Debugging Hooks

```python
class DebugChain:
    def __init__(self, steps, debug=False):
        self.steps = steps
        self.debug = debug
        self.history = []

    def run(self, input_data):
        current = input_data

        for step in self.steps:
            result = step(current)

            if self.debug:
                self.history.append({
                    "step": step.__name__,
                    "input": current[:200],
                    "output": result[:200]
                })
                print(f"[{step.__name__}]")
                print(f"  In: {current[:100]}...")
                print(f"  Out: {result[:100]}...")

            current = result

        return current
```

## Exercise

### Exercise 12.1: Design a Chain

Design a prompt chain for:
- Input: Raw customer feedback
- Output: Structured report with sentiment, topics, and action items

Show each step with inputs/outputs.

### Exercise 12.2: Add Error Handling

Add error handling to this chain:
```python
def analyze_sentiment(text):
    return prompt(f"Sentiment of: {text}")

def extract_topics(text):
    return prompt(f"Topics in: {text}")

def generate_response(sentiment, topics):
    return prompt(f"Write response for {sentiment} about {topics}")
```

### Exercise 12.3: Parallel Chain

Design a parallel chain that:
- Processes a document
- Runs 3 analyses simultaneously
- Combines results

## Key Takeaways

- ✅ Chains break complex tasks into manageable steps
- ✅ Use sequential, parallel, branching, or map-reduce patterns
- ✅ Each step should have single responsibility
- ✅ Implement retry, fallback, and validation
- ✅ Log and monitor for debugging

## Next Steps

→ [PE-13: Code Generation & Debugging](../PE-13-code-generation/README.md)
