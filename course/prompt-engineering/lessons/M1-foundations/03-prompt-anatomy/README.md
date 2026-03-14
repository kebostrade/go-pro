# PE-03: Anatomy of a Good Prompt

**Duration**: 2 hours
**Module**: 1 - Foundations

## Learning Objectives

- Master the components of effective prompts
- Learn the CREATE framework for prompt structure
- Understand system vs user prompts
- Apply constraints and formatting requirements

## The Anatomy of a Prompt

### Core Components

```
┌─────────────────────────────────────────────────────────────┐
│                        PROMPT                                │
├─────────────────────────────────────────────────────────────┤
│  ┌─────────────────────────────────────────────────────┐    │
│  │  1. ROLE/PERSONA          (Who is the AI?)          │    │
│  │     "You are a senior software engineer..."          │    │
│  └─────────────────────────────────────────────────────┘    │
│  ┌─────────────────────────────────────────────────────┐    │
│  │  2. CONTEXT               (Background info)          │    │
│  │     "We're building a fintech app..."                │    │
│  └─────────────────────────────────────────────────────┘    │
│  ┌─────────────────────────────────────────────────────┐    │
│  │  3. TASK                  (What to do)               │    │
│  │     "Review this code for security issues..."        │    │
│  └─────────────────────────────────────────────────────┘    │
│  ┌─────────────────────────────────────────────────────┐    │
│  │  4. CONSTRAINTS           (Rules/limits)             │    │
│  │     "Focus on SQL injection and XSS..."              │    │
│  └─────────────────────────────────────────────────────┘    │
│  ┌─────────────────────────────────────────────────────┐    │
│  │  5. OUTPUT FORMAT         (How to respond)           │    │
│  │     "Return as JSON with keys: severity, issue..."   │    │
│  └─────────────────────────────────────────────────────┘    │
│  ┌─────────────────────────────────────────────────────┐    │
│  │  6. EXAMPLES              (Few-shot samples)         │    │
│  │     "Example: {severity: 'high', issue: '...'}"      │    │
│  └─────────────────────────────────────────────────────┘    │
└─────────────────────────────────────────────────────────────┘
```

## The CREATE Framework

A systematic approach to building prompts:

| Letter | Component | Description |
|--------|-----------|-------------|
| **C** | Context | Background information needed |
| **R** | Role | Who the AI should be |
| **E** | Explicit Task | Clear, specific instruction |
| **A** | Audience | Who the output is for |
| **T** | Tone/Format | Style and structure of output |
| **E** | Examples | Sample inputs/outputs |

### CREATE in Action

```
[C] Context: You're working on a Python web application that processes
    user payments. The app uses Flask and Stripe API.

[R] Role: Act as a security auditor specializing in payment systems.

[E] Explicit Task: Review the provided payment processing code and
    identify security vulnerabilities.

[A] Audience: The development team (intermediate skill level).

[T] Tone/Format: Professional, technical. Output as a markdown table
    with columns: Vulnerability, Severity, Line Number, Recommendation.

[E] Examples:
    Input: payment.process(amount, user_id)
    Output: | Vulnerability | Severity | Line | Recommendation |
            | No input validation | High | 15 | Validate amount > 0 |
```

## System vs User Prompts

### System Prompt

Defines the AI's behavior, personality, and constraints. Set once per conversation.

```python
system_prompt = """
You are a concise, accurate technical assistant.
- Always verify facts before stating them
- If uncertain, say "I'm not certain" rather than guessing
- Format code blocks with language specification
- Keep responses under 200 words unless more detail is requested
"""
```

### User Prompt

The actual task or question from the user.

```python
user_prompt = """
Write a Python function to validate email addresses.
Handle edge cases like subdomains and plus addressing.
"""
```

### API Usage

```python
response = client.messages.create(
    model="claude-sonnet-4-6-20250514",
    system=system_prompt,  # System prompt here
    messages=[
        {"role": "user", "content": user_prompt}  # User prompt here
    ]
)
```

## Constraints & Guardrails

### Types of Constraints

| Type | Example | Purpose |
|------|---------|---------|
| **Scope** | "Only discuss Python" | Limit topic |
| **Length** | "Max 100 words" | Control verbosity |
| **Format** | "Return JSON only" | Structure output |
| **Tone** | "Be professional" | Control style |
| **Safety** | "No code execution" | Security |

### Constraint Examples

```
SCOPE CONSTRAINT:
"Answer only questions about the provided documentation.
If the question is outside scope, respond: 'I can only help with documentation questions.'"

LENGTH CONSTRAINT:
"Your response must be exactly 3 bullet points, no more, no less."

FORMAT CONSTRAINT:
"Output must be valid JSON. No markdown formatting, no explanations outside JSON."

NEGATIVE CONSTRAINT:
"Do not use the words 'leverage', 'synergy', or 'robust'."
```

## Output Formatting

### Specifying Format

```
✅ Clear Format Specification:
"""
Return the analysis as a JSON object with this structure:
{
  "sentiment": "positive" | "negative" | "neutral",
  "confidence": <float between 0 and 1>,
  "keywords": ["array", "of", "keywords"],
  "summary": "one sentence summary"
}
"""
```

### Format Enforcement Techniques

```
TECHNIQUE 1: Explicit Schema
"""
Output format (strict):
- Title: [one line]
- Summary: [2-3 sentences]
- Points: [bullet list, max 5]
- Recommendation: [one sentence]
"""

TECHNIQUE 2: Delimiters
"""
Put your analysis between <analysis> and </analysis> tags.
Put your recommendation between <recommendation> and </recommendation> tags.
"""

TECHNIQUE 3: Example-Driven
"""
Format your response exactly like this example:
---
**Analysis**: The code has a race condition in line 42.
**Severity**: High
**Fix**: Add a mutex lock before accessing shared state.
---
"""
```

## Complete Example

### Task: Code Review Prompt

```
You are an expert code reviewer with 15 years of experience in Python
and web security. You specialize in identifying performance bottlenecks
and security vulnerabilities.

CONTEXT:
We're reviewing code for a financial services API that handles
thousands of requests per second. Performance and security are critical.

TASK:
Review the following Python function and provide:
1. Security issues (injection, auth, data validation)
2. Performance concerns (algorithmic complexity, database queries)
3. Code quality issues (readability, maintainability)

CODE TO REVIEW:
```python
def get_user_data(user_id, db):
    query = f"SELECT * FROM users WHERE id = {user_id}"
    result = db.execute(query)
    users = []
    for row in result:
        users.append(dict(row))
    return users
```

CONSTRAINTS:
- Focus on critical and high-severity issues only
- Ignore minor style preferences
- Assume the code runs in production

OUTPUT FORMAT:
Return a JSON array of issues:
[
  {
    "type": "security" | "performance" | "quality",
    "severity": "critical" | "high" | "medium" | "low",
    "line": <line_number>,
    "issue": "<description>",
    "fix": "<recommended_fix>"
  }
]

If no issues found, return: []
```

## Exercise

### Exercise 3.1: Deconstruct a Prompt

Analyze this prompt using the CREATE framework:

```
"Summarize this article about climate change for a middle school
science class. Make it engaging and use simple vocabulary.
Include 3 fun facts and format as bullet points."
```

### Exercise 3.2: Build a Prompt

Using CREATE, build a prompt for:
- Task: Generate API documentation
- Input: A Python function
- Output: Markdown documentation with parameters, returns, examples

### Exercise 3.3: Add Constraints

Add appropriate constraints to this prompt:

```
"Help me write a blog post about AI."
```

## Key Takeaways

- ✅ Good prompts have: Role, Context, Task, Constraints, Format, Examples
- ✅ CREATE framework ensures systematic prompt construction
- ✅ System prompts define behavior; user prompts contain tasks
- ✅ Constraints prevent unwanted outputs and ensure consistency

## Next Steps

→ [PE-04: Zero-Shot Prompting](../PE-04-zero-shot/README.md)
