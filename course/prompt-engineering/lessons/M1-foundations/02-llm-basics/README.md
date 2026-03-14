# PE-02: How LLMs Work

**Duration**: 3 hours
**Module**: 1 - Foundations

## Learning Objectives

- Understand tokenization and how text becomes numbers
- Learn about context windows and their limitations
- Grasp attention mechanisms at a conceptual level
- Understand temperature and sampling parameters

## How LLMs Process Text

### The Pipeline

```
Input Text → Tokenization → Tokens → Model → Output Tokens → Detokenization → Output Text
```

## Tokenization

### What are Tokens?

Tokens are the basic units that LLMs process. They're not quite words or characters—they're somewhere in between.

```
Text:   "Hello, world!"
Tokens: ["Hello", ",", " world", "!"]
IDs:    [15496, 11, 995, 0]
```

### Token Count Examples

| Text | Approximate Tokens |
|------|-------------------|
| "Hello" | 1 |
| "Hello, world!" | 4 |
| "The quick brown fox jumps over the lazy dog." | 10 |
| 1,000 words of English text | ~1,300 tokens |
| 1,000 lines of code | ~2,000-3,000 tokens |

### Tokenization Rules of Thumb

- **English**: ~1 token ≈ 4 characters ≈ 0.75 words
- **Code**: More tokens due to special characters
- **Non-English**: Can be 2-4x more tokens per word

### Why It Matters for Prompts

```
❌ Inefficient: "Please kindly assist me in the creation of..."
✅ Efficient: "Help me create..."
```

## Context Windows

### What is a Context Window?

The maximum amount of text (input + output) the model can process in one request.

```
┌──────────────────────────────────────────────────────────┐
│                    CONTEXT WINDOW                        │
│  ┌────────────────────────────────────────────────────┐  │
│  │                    INPUT                           │  │
│  │  System Prompt + User Message + History + Docs     │  │
│  └────────────────────────────────────────────────────┘  │
│  ┌────────────────────────────────────────────────────┐  │
│  │                    OUTPUT                          │  │
│  │           Model's Response (generated)             │  │
│  └────────────────────────────────────────────────────┘  │
└──────────────────────────────────────────────────────────┘
```

### Model Context Limits (2026)

| Model | Context Window |
|-------|---------------|
| Claude 3.5 Sonnet | 200K tokens |
| Claude 3 Opus | 200K tokens |
| GPT-4o | 128K tokens |
| GPT-4 Turbo | 128K tokens |

### Context Window Implications

1. **Truncation**: Old messages get cut off
2. **Cost**: More context = more cost
3. **Quality**: Too much irrelevant context can degrade performance

### Managing Context

```python
# Good: Essential context only
prompt = """
Documentation:
- Function: calculate_total(items: list) -> float
- Returns sum of all item prices

Task: Write a test for this function.
"""

# Bad: Irrelevant context
prompt = """
Our company was founded in 2010...
[1000 lines of company history]...
The function calculate_total(items: list) returns a float...
"""
```

## Attention Mechanism

### Conceptual Understanding

Attention lets the model focus on relevant parts of the input when generating each word.

```
Input: "The cat sat on the mat because it was tired."

When generating "it", the model attends to:
  - "cat" ████████████ (high attention)
  - "mat" ██ (low attention)
  - "sat" ███ (medium attention)
```

### Why Attention Matters for Prompts

**Key insight**: Place important information at the **beginning** or **end** of prompts.

```
✅ Good: Important info at start
"""
CRITICAL: Never reveal sensitive data.

You are a customer support agent. Help the user with their question.
"""

✅ Good: Important info at end
"""
You are a customer support agent. Help the user with their question.

REMEMBER: Never reveal sensitive data.
"""

❌ Less optimal: Important info buried in middle
"""
You are a customer support agent for Acme Corp.
Acme Corp was founded in 2010 and serves customers worldwide.
CRITICAL: Never reveal sensitive data.
We have offices in New York, London, and Tokyo.
Help the user with their question.
"""
```

## Temperature & Sampling

### Temperature

Controls randomness/creativity in outputs.

| Temperature | Behavior | Use Case |
|-------------|----------|----------|
| 0 | Deterministic, most likely tokens | Code, factual answers |
| 0.3 | Low randomness | Technical writing |
| 0.7 | Balanced | General chat |
| 1.0 | High creativity | Creative writing |

### Visual Example

```
Prompt: "The sky is"

Temperature 0.0: "blue" (always)
Temperature 0.5: "blue" / "clear" / "bright" (varies slightly)
Temperature 1.0: "a canvas of endless possibility" / "painted with sunset hues" (varies a lot)
```

### Other Parameters

```python
response = client.messages.create(
    model="claude-sonnet-4-6-20250514",
    max_tokens=1024,        # Max output length
    temperature=0.7,        # Randomness (0-1)
    top_p=0.9,             # Nucleus sampling
    top_k=50,              # Top-k sampling
    messages=[...]
)
```

## Model Capabilities & Limitations

### What LLMs Excel At

- ✅ Text generation and transformation
- ✅ Code generation and explanation
- ✅ Summarization and extraction
- ✅ Translation between languages
- ✅ Question answering (with context)
- ✅ Creative writing

### What LLMs Struggle With

- ❌ Math (without tools)
- ❌ Real-time information
- ❌ Very long-term memory
- ❌ Spatial reasoning
- ❌ Truth verification (can hallucinate)

### Hallucinations

LLMs can generate plausible-sounding but incorrect information.

```
User: "What did Abraham Lincoln say about the internet in his 1863 speech?"
Model: "In his famous 1863 address, Lincoln spoke about how the internet
would connect the nation..."  ← HALLUCINATION (internet didn't exist)
```

**Mitigation**:
- Provide accurate context
- Ask model to cite sources
- Verify important facts

## Exercise

### Exercise 2.1: Token Counting

Estimate tokens for these prompts:

1. A 500-word email
2. A 100-line Python file
3. A JSON object with 50 key-value pairs

### Exercise 2.2: Context Management

Rewrite this prompt to be more token-efficient:

```
I am writing to you today to request your assistance with a matter
of some importance. I have been trying to solve a programming
problem for quite some time now, and I was hoping that perhaps
you might be able to help me figure out how to write a function
that calculates the factorial of a number in Python. It would be
greatly appreciated if you could provide me with some guidance
on this matter. Thank you very much for your time and consideration.
```

### Exercise 2.3: Temperature Selection

Choose appropriate temperatures for:

1. Writing a legal contract
2. Generating poem ideas
3. Debugging a critical bug
4. Brainstorming product names

## Key Takeaways

- ✅ Tokens are ~4 characters; context windows limit input+output
- ✅ Important info should be at start or end of prompts
- ✅ Temperature controls creativity (0 for code, higher for creative)
- ✅ LLMs can hallucinate—verify important information

## Next Steps

→ [PE-03: Anatomy of a Good Prompt](../PE-03-prompt-anatomy/README.md)
