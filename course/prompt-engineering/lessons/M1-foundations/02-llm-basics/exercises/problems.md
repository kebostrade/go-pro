# Exercise: How LLMs Work

## Problem 1: Tokenization

The sentence: "AI is amazing!" would be tokenized differently by different models.

### Your Task:
Think about how each word might be broken into tokens:
- "AI" - likely 1-2 tokens
- "is" - likely 1 token
- "amazing!" - likely 2-3 tokens

Why might "amazing!" be more tokens than "amazing"?

---

## Problem 2: Context Window

A model has a context window of 4,000 tokens.

### Your Task:
Calculate if each prompt + expected output would fit:
- Prompt: 100 tokens, Expected output: 3,800 tokens
- Prompt: 2,000 tokens, Expected output: 2,500 tokens
- Prompt: 3,999 tokens, Expected output: 1 token

---

## Problem 3: Attention Mechanism

Explain in your own words why the attention mechanism allows LLMs to understand context in the sentence:

"The dog chased the cat because it was hungry."

What does "it" refer to? How does attention help the model know?

---

## Problem 4: Temperature Setting

For each scenario, choose an appropriate temperature setting (0.0-1.0):

| Scenario | Recommended Temperature | Why |
|----------|-------------------------|-----|
| Code generation | | |
| Creative story writing | | |
| Factual question answering | | |
| Mathematical problem solving | | |

---

## Problem 5: Model Comparison

Research and compare two modern LLMs (e.g., Claude, GPT-4, Gemini):

| Feature | Model A | Model B |
|---------|---------|---------|
| Context window | | |
| Training cutoff | | |
| Strengths | | |
| Best use cases | | |
