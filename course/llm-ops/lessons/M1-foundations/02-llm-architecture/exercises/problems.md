# Exercise: LLM Architecture Fundamentals

## Problem 1: Self-Attention Visualization

Given the sentence: "The fraud detection model identified suspicious activity and blocked the transaction"

For each of the following target words, identify the top 3 words they would attend to most strongly and explain why:

### a) "suspicious"
1. 
2. 
3. 

### b) "blocked"
1. 
2. 
3. 

---

## Problem 2: VRAM Calculation

Calculate the GPU memory (VRAM) required for inference with each configuration:

### Configuration A: Llama 3 8B
- Parameters: 8 billion
- Precision: FP16 (16 bits per parameter)

**Calculation:**
```
8B × 16 bits = ? bits = ? bytes = ? GB
```

### Configuration B: Mistral 7B with INT8 quantization
- Parameters: 7 billion
- Precision: INT8 (8 bits per parameter)

**Calculation:**

### Configuration C: Llama 3 70B with INT4 quantization
- Parameters: 70 billion  
- Precision: INT4 (4 bits per parameter)
- Plus 1.2x overhead for kv cache

**Calculation:**

---

## Problem 3: Model Selection Decision

You're building a real-time customer support chatbot for a fintech company.

Requirements:
- Response time: < 2 seconds
- Must run on company infrastructure (cannot use external APIs)
- Must handle financial domain terminology accurately
- Budget: Can purchase up to 2x NVIDIA A100 (80GB each)

Answer the following:

### a) Which model would you choose and why?

| Candidate | Pros | Cons | Verdict |
|-----------|------|------|---------|
| Phi-3 4B | | | |
| Mistral 7B | | | |
| Llama 3 8B | | | |
| Llama 3 70B | | | |

### b) What quantization level would you use?

### c) What's your estimated per-request latency?

---

## Problem 4: Context Window Analysis

A RAG system needs to process documents of varying lengths:

| Document Type | Avg Tokens | Max Tokens |
|---------------|-----------|-------------|
| Support tickets | 500 | 2,000 |
| Policy documents | 5,000 | 20,000 |
| Legal contracts | 15,000 | 50,000 |

### Questions:

a) Which documents fit entirely within a 32K context window?

b) Which require chunking with a 128K context window?

c) What's the trade-off between longer context windows and infrastructure cost?

---

## Problem 5: Architecture Pattern Matching

Match each use case to the most appropriate LLM architecture:

**Use Cases:**
1. Generating creative marketing copy
2. Classifying customer messages as urgent/not urgent
3. Translating English to Spanish
4. Building a chatbot that generates detailed responses

**Architectures:**
A) Encoder-only (BERT-style)
B) Decoder-only (GPT-style)  
C) Encoder-Decoder (T5-style)

| Use Case | Architecture | Reason |
|----------|--------------|--------|
| 1 | | |
| 2 | | |
| 3 | | |
| 4 | | |

---

## Problem 6: Latency Breakdown

For a self-hosted Llama 3 8B model running on an NVIDIA A100:

```
Total Latency = Tokenization + Model Forward + Sampling + De-tokenization
```

Estimate each component:
- Tokenization (CPU): 
- Model Forward (GPU): 
- Sampling: 
- De-tokenization: 

**Total Estimate:** 

If you process a 100-token input generating a 200-token response, what's the total expected latency?

---

## Submission

Complete all problems and submit your answers. Be prepared to discuss your reasoning in the next lesson.
