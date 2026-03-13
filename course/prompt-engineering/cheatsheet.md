# Prompt Engineering Cheatsheet

Quick reference for effective prompt engineering.

## The PIERS Framework

| Component | Description | Example |
|-----------|-------------|---------|
| **P**urpose | What do you want? | "Summarize this article" |
| **I**nformation | Context needed | "For a business audience" |
| **E**xamples | Show, don't tell | Few-shot examples |
| **R**ole | Who should AI be? | "You are a senior editor" |
| **S**tyle | Output format | "Return as JSON" |

## The CREATE Framework

| Letter | Component | Description |
|--------|-----------|-------------|
| **C**ontext | Background info | Domain, situation, constraints |
| **R**ole | Persona | Who the AI should be |
| **E**xplicit Task | Clear instruction | Specific action required |
| **A**udience | Target reader | Who will consume output |
| **T**one/Format | Style | Formal, casual, JSON, markdown |
| **E**xamples | Samples | Input/output pairs |

## Prompt Components

```
┌─────────────────────────────────────────────────┐
│ 1. ROLE/PERSONA    - "You are a [expert]..."    │
│ 2. CONTEXT         - Background information     │
│ 3. TASK            - What to do specifically    │
│ 4. CONSTRAINTS     - Rules and limitations      │
│ 5. OUTPUT FORMAT   - How to structure response  │
│ 6. EXAMPLES        - Few-shot samples           │
└─────────────────────────────────────────────────┘
```

## Core Techniques

### Zero-Shot
```
[Task description]
[Input]
```

### Few-Shot
```
Example 1:
Input: [example]
Output: [example]

Example 2:
Input: [example]
Output: [example]

Now:
Input: [actual input]
Output:
```

### Chain-of-Thought
```
[Task]
Let's think step by step.
```

### ReAct
```
Thought: [reasoning]
Action: [tool]
Action Input: [parameters]
Observation: [result]
[repeat]
Final Answer: [response]
```

## Output Format Templates

### JSON Output
```
Return as JSON with this structure:
{
  "field1": "type",
  "field2": ["array"],
  "field3": {
    "nested": "value"
  }
}

No markdown, no explanations, only valid JSON.
```

### Structured Text
```
Output format:
**Title**: [one line]
**Summary**: [2-3 sentences]
**Points**:
- [point 1]
- [point 2]
**Recommendation**: [one sentence]
```

## Temperature Guide

| Value | Use Case |
|-------|----------|
| 0.0 | Code, factual answers |
| 0.3 | Technical writing |
| 0.5 | Balanced responses |
| 0.7 | General conversation |
| 1.0 | Creative writing |

## Token Estimation

| Content Type | Approximate Ratio |
|--------------|-------------------|
| English text | 1 token ≈ 4 chars ≈ 0.75 words |
| Code | 1 token ≈ 3-4 chars |
| Non-English | 2-4x more tokens per word |

## Common Patterns

### Classification
```
Classify into: [category1, category2, category3]

Input: [text]
Category:
```

### Extraction
```
Extract [entities] as JSON:
{"field": "value"}

Text: [input]
```

### Transformation
```
Transform from [format A] to [format B]:

Input: [example]
Output: [example]

Now transform:
Input: [actual]
Output:
```

### Summarization
```
Summarize this [content type] for [audience].

Key requirements:
- [requirement 1]
- [requirement 2]

Length: [word count]
```

## Best Practices

### DO ✅
- Be specific and explicit
- Provide examples for complex tasks
- Define output format clearly
- Use consistent structure
- Include constraints
- State what NOT to do

### DON'T ❌
- Be vague ("make it better")
- Assume context is understood
- Request multiple formats
- Use ambiguous language
- Skip edge cases
- Over-constrain (too many rules)

## Debugging Prompts

### If output is too long:
```
Keep response under [N] words.
Be concise.
```

### If output is wrong format:
```
Output ONLY [format].
No explanations outside [format].
```

### If output is inconsistent:
```
Follow this EXACT format:
[show format]
```

### If task is misunderstood:
```
Your goal is to [specific action].
This means: [clarification].
NOT: [what it's not].
```

## Agentic Prompt Template

```
You are [agent description].

GOAL: [objective]

TOOLS:
- tool1: [description]
- tool2: [description]

FORMAT:
Thought: [reasoning]
Action: [tool]
Action Input: [input]
Observation: [result]
Final Answer: [when done]

CONSTRAINTS:
- [constraint 1]
- [constraint 2]

Begin!
```

## Evaluation Metrics

| Metric | Formula | Use When |
|--------|---------|----------|
| Accuracy | correct/total | Classification |
| F1 Score | 2*(P*R)/(P+R) | Imbalanced classes |
| BLEU | n-gram overlap | Translation |
| ROUGE | recall of n-grams | Summarization |
| Latency | time to response | Performance |
| Cost | tokens × rate | Budget tracking |

## Security Checklist

- [ ] Never include credentials in prompts
- [ ] Validate all inputs
- [ ] Sanitize outputs before display
- [ ] Rate limit requests
- [ ] Log all interactions
- [ ] Implement content filtering
- [ ] Test for injection attacks

## Version Control Best Practices

```
prompts/
├── registry.yaml      # All prompt definitions
├── sentiment/
│   ├── v1.yaml        # Deprecated
│   ├── v2.yaml        # Current stable
│   └── v3.yaml        # Canary testing
└── tests/
    └── sentiment/
        └── test_cases.json
```

## Quick Formula

```
Good Prompt = Role + Context + Task + Constraints + Format + Examples

Token Budget = Context Window - Input Tokens - Output Buffer
             = 200,000 - prompt_tokens - 1,000

Cost = (input_tokens × input_rate) + (output_tokens × output_rate)
```

---
*Last Updated: March 2026*
