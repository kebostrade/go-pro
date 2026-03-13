# PE-04: Zero-Shot Prompting

**Duration**: 2 hours
**Module**: 1 - Foundations

## Learning Objectives

- Understand zero-shot prompting and when to use it
- Learn techniques to improve zero-shot performance
- Apply instruction tuning principles
- Handle edge cases and failures

## What is Zero-Shot Prompting?

Zero-shot prompting asks the model to perform a task without providing examples.

```
┌─────────────────────────────────────────────────────────┐
│                   ZERO-SHOT                             │
│                                                         │
│  Prompt: "Translate to French: Hello, how are you?"    │
│  Output: "Bonjour, comment allez-vous?"                │
│                                                         │
│  No examples provided → Model uses pre-trained knowledge │
└─────────────────────────────────────────────────────────┘
```

### Comparison

| Technique | Examples Provided | Use Case |
|-----------|------------------|----------|
| Zero-Shot | 0 | Simple, well-defined tasks |
| One-Shot | 1 | When you need format guidance |
| Few-Shot | 2-5+ | Complex tasks, specific formats |

## When Zero-Shot Works Well

### ✅ Good for Zero-Shot

1. **Common tasks** the model has seen in training
   - Translation
   - Summarization
   - Basic Q&A
   - Sentiment analysis

2. **Clear, unambiguous instructions**
   - "Convert to uppercase"
   - "Count the words"
   - "Extract all email addresses"

3. **Standard formats**
   - JSON generation
   - Markdown formatting
   - Code syntax

### ❌ Poor for Zero-Shot

1. **Novel or unique formats**
2. **Domain-specific conventions**
3. **Subtle distinctions**
4. **Multi-step reasoning** (without chain-of-thought)

## Zero-Shot Techniques

### 1. Direct Instruction

```
Simple and direct:
"Summarize this text in 3 bullet points."

"Classify this email as 'spam' or 'not spam'."

"Fix the grammar in this sentence."
```

### 2. Role-Based Zero-Shot

```
You are a professional editor. Rewrite this paragraph to be more concise
while maintaining the original meaning.

[paragraph]
```

### 3. Constraint-Based Zero-Shot

```
Generate a password that:
- Is exactly 12 characters
- Contains at least 2 uppercase letters
- Contains at least 2 numbers
- Contains at least 1 special character
- Does not contain the word "password"
```

### 4. Format-Specified Zero-Shot

```
Extract the key information from this product description as JSON:
{
  "name": "",
  "price": "",
  "features": [],
  "rating": null
}

Product: "The UltraWidget Pro is a revolutionary gadget priced at $49.99.
It features Bluetooth 5.0, water resistance, and a 4.5-star rating from
over 2,000 reviews."
```

## Improving Zero-Shot Performance

### Technique 1: Be Explicit About Output

```
❌ Vague:
"Analyze this text."

✅ Explicit:
"Analyze this text and identify:
1. Main topic
2. Tone (formal/informal)
3. Target audience
4. Key arguments"
```

### Technique 2: Specify What NOT to Do

```
❌ Without negatives:
"Summarize this article."

✅ With negatives:
"Summarize this article. Do NOT include:
- Personal opinions
- Information not in the original text
- More than 3 sentences"
```

### Technique 3: Add Confidence Indicators

```
Answer the following question. If you're uncertain about any part
of your answer, indicate this with [UNCERTAIN] before that part.

Question: What is the population of Tokyo in 2024?
```

### Technique 4: Request Structured Reasoning

```
Before answering, think through this step by step:
1. What information do I need?
2. What does the provided text say?
3. What can I conclude?

Then provide your answer.
```

## Zero-Shot Prompt Patterns

### Pattern 1: Classification

```
Classify the following customer feedback into one of these categories:
- Product Issue
- Service Complaint
- Feature Request
- General Inquiry

Feedback: "I've been waiting 2 weeks for my order and nobody responds to my emails."

Category:
```

### Pattern 2: Extraction

```
Extract all dates and their associated events from this text:

"On March 15th, 2024, we launched version 2.0. The beta testing began
January 5th and concluded February 28th. Version 3.0 is scheduled for
release on December 1st, 2024."

Output format:
- Date: [date] | Event: [event]
```

### Pattern 3: Transformation

```
Convert this SQL query to a pandas DataFrame operation:

SQL: SELECT name, age FROM users WHERE age > 18 ORDER BY age DESC

Pandas:
```

### Pattern 4: Validation

```
Review this email for professionalism issues:
- Check for informal language
- Identify potential tone problems
- Flag any missing elements (greeting, closing, etc.)

Email:
"hey can u send me the files? need them asap. thanks"

Issues found:
```

## Handling Zero-Shot Failures

### Failure Mode 1: Format Inconsistency

```
Problem: Model doesn't follow requested format consistently.

Solution: Add format enforcement:
"IMPORTANT: Your response must be ONLY valid JSON. No other text.
No explanations. No markdown. Just the JSON object."
```

### Failure Mode 2: Task Misunderstanding

```
Problem: Model misinterprets what you're asking.

Solution: Clarify with definitions:
"By 'summarize', I mean create a concise version that captures
only the essential points, omitting details and examples."
```

### Failure Mode 3: Edge Case Failure

```
Problem: Model fails on unusual inputs.

Solution: Add explicit handling instructions:
"If the input is empty or invalid, respond with 'INVALID_INPUT'
and do not attempt to process it."
```

### Failure Mode 4: Hallucination

```
Problem: Model makes up information.

Solution: Constrain knowledge sources:
"Answer using ONLY the information provided below. If the answer
is not in the provided text, respond with 'Information not provided.'

[Provided text]"
```

## Complete Zero-Shot Example

```
TASK: Product Description Generator

You are a professional copywriter for an e-commerce platform.

Generate a product description for the following item:
- Product: Wireless Earbuds
- Brand: SoundMax
- Features: 24-hour battery, noise cancellation, waterproof IPX7
- Price: $79.99
- Target: Fitness enthusiasts

REQUIREMENTS:
- 150-200 words
- Professional but energetic tone
- Include a call-to-action
- Structure: Hook → Features → Benefits → CTA
- Do NOT make claims not supported by the provided features

OUTPUT:
```

## Exercise

### Exercise 4.1: Improve Zero-Shot Prompts

Improve these zero-shot prompts:

1. "Make this better"
2. "Is this good?"
3. "Extract the important stuff"
4. "Fix the errors"

### Exercise 4.2: Zero-Shot Classification

Write a zero-shot prompt to classify customer support tickets into:
- Technical Issue
- Billing Question
- Feature Request
- General Feedback

### Exercise 4.3: Zero-Shot with Constraints

Write a zero-shot prompt that:
- Generates a weekly meal plan
- Accommodates dietary restrictions (user will specify)
- Outputs as a markdown table
- Includes calorie estimates

## Key Takeaways

- ✅ Zero-shot works for common, well-defined tasks
- ✅ Be explicit about output format and constraints
- ✅ Specify what NOT to do, not just what to do
- ✅ Add structured reasoning for complex tasks
- ✅ Handle failures with explicit instructions

## Next Steps

→ [PE-05: Few-Shot Learning](../PE-05-few-shot/README.md)
