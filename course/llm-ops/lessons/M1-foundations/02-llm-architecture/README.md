# IO-02: LLM Architecture Fundamentals

**Duration**: 3 hours
**Module**: 1 - Foundations

## Learning Objectives

- Understand the Transformer architecture and its components
- Explain how attention mechanisms work
- Learn about tokenization and model sizing
- Connect architecture concepts to infrastructure decisions

## The Transformer Architecture

The Transformer, introduced in the paper "Attention Is All You Need" (2017), is the foundation of all modern LLMs. It uses self-attention to process sequences in parallel rather than sequentially.

```
┌─────────────────────────────────────────────────────────────────┐
│                     Transformer Architecture                     │
├─────────────────────────────────────────────────────────────────┤
│  Input: "Analyze this transaction for fraud"                    │
├─────────────────────────────────────────────────────────────────┤
│  ┌─────────────┐    ┌─────────────┐    ┌─────────────┐        │
│  │   Input     │───▶│   Encoder   │───▶│   Output    │        │
│  │  Embedding  │    │    Stack    │    │   Projection│        │
│  └─────────────┘    └─────────────┘    └─────────────┘        │
│                            │                                    │
│                     ┌──────┴──────┐                             │
│                     │ Self-Attention│                           │
│                     │   + FFN      │                             │
│                     └─────────────┘                             │
│                            │                                    │
│                     (repeated N times)                          │
└─────────────────────────────────────────────────────────────────┘
```

### Key Components

1. **Token Embeddings**: Convert words to numerical vectors
2. **Positional Encoding**: Add position information to embeddings
3. **Self-Attention**: Model relationships between all tokens
4. **Feed-Forward Networks**: Process attention outputs
5. **Layer Normalization**: Stabilize training

## Understanding Attention

Attention allows the model to "focus" on relevant parts of the input when generating each output token.

### How Self-Attention Works

```
Input: "the fraud detection system flagged the transaction"

For the word "flagged":
┌─────────────────────────────────────────────────────────────────┐
│                    Self-Attention Weights                        │
├─────────────────────────────────────────────────────────────────┤
│  the        fraud      detection    system    flagged   the    │
│  0.05  ──▶ 0.45  ──▶   0.25    ──▶  0.10  ──▶ [self]  0.05   │
│         ▲                                              │       │
│         └──────────────────────────────────────────────┘       │
│                                                                  │
│  "flagged" attends most to "fraud" and "detection"             │
└─────────────────────────────────────────────────────────────────┘
```

### Multi-Head Attention

Instead of one attention mechanism, Transformers use multiple "heads" that learn different types of relationships:

- **Syntactic**: Subject-verb relationships
- **Semantic**: Meaning connections
- **Temporal**: Time-based dependencies
- **Entity**: Referencing the same entity

## Tokenization

LLMs don't process raw text - they work with tokens. Tokenization converts text into integer IDs that the model can process.

### Tokenization Approaches

| Method | Description | Example |
|--------|-------------|---------|
| **Word-based** | Split by spaces | ["Analyze", "this", "transaction"] |
| **BPE** | Byte Pair Encoding | ["An", "alyze", "this", "transaction"] |
| **WordPiece** | Google's approach | Similar to BPE, used in BERT |
| **SentencePiece** | Unsupervised tokenizer | Used in many LLMs |

### Token Limits and Context Window

| Model | Context Window | Approx. Tokens |
|-------|---------------|----------------|
| GPT-4 | 128K | ~100K words |
| Claude 3 | 200K | ~150K words |
| Llama 3 | 128K | ~100K words |
| Mistral 7B | 32K | ~24K words |

**Infrastructure Implication**: Longer context windows require more VRAM and increase latency.

## Model Sizes and Parameters

LLMs are defined by their parameter count - the weights learned during training.

### Common Model Sizes

| Model | Parameters | VRAM (FP16) | Use Case |
|-------|------------|-------------|----------|
| Phi-3 | 3.8B | ~8GB | Mobile/Edge |
| Llama 3 8B | 8B | ~16GB | Single GPU |
| Mistral 7B | 7B | ~14GB | Single GPU |
| Llama 3 70B | 70B | ~140GB (2xA100) | Multi-GPU |
| GPT-4 | ~1.7T (experts) | Cloud only | API only |

### Quantization

Reduce model size and memory requirements by using lower precision:

| Precision | Bits | Memory Reduction | Quality Loss |
|-----------|------|------------------|--------------|
| FP32 | 32 | 1x (baseline) | None |
| FP16/BF16 | 16 | 2x | Minimal |
| INT8 | 8 | 4x | Low |
| INT4 | 4 | 8x | Moderate |

**Infrastructure Impact**: Quantization allows running larger models on limited GPUs.

## Architecture Patterns for Production

### Causal vs. Encoder-Only vs. Encoder-Decoder

| Architecture | Use Case | Examples |
|--------------|----------|----------|
| **Causal (Decoder-only)** | Text generation | GPT, Llama, Claude |
| **Encoder-only** | Classification, embedding | BERT, RoBERTa |
| **Encoder-Decoder** | Seq2seq, translation | T5, BART |

### Go Implementation Pattern

The AI Agent Platform demonstrates LLM abstraction in Go:

```go
// From services/ai-agent-platform/internal/llm/
type Provider interface {
    Generate(ctx context.Context, req Request) (*Response, error)
    Stream(ctx context.Context, req Request) (<-chan *Response, error)
}

// Supports multiple providers: OpenAI, Anthropic, Ollama
// Allows infrastructure to swap providers based on cost/performance
```

## Key Terminology

| Term | Definition |
|------|------------|
| **Transformer** | Neural network architecture using self-attention |
| **Self-Attention** | Mechanism to weigh relationships between tokens |
| **Token** | Basic unit of text processed by LLMs (~4 chars) |
| **Context Window** | Maximum tokens the model can process |
| **Parameters** | Learned weights in the neural network |
| **Quantization** | Reducing precision to save memory |
| **VRAM** | Video RAM - memory on GPU for model loading |
| **Inference** | Generating output from a trained model |
| **Feed-Forward Network (FFN)** | Dense layers in transformer block |
| **Layer Normalization** | Stabilization technique in transformers |

## Understanding Latency Factors

```
Total Latency = Tokenization + Model Forward + Sampling + De-tokenization

Where:
├── Tokenization: ~1-5ms (CPU)
├── Model Forward: varies by model size and GPU
│   ├── 7B model on A100: ~50-100ms
│   ├── 70B model on 2xA100: ~200-400ms
│   └── GPT-4 API: ~500-2000ms
├── Sampling: ~1-10ms
└── De-tokenization: ~1-5ms
```

## Exercise

### Exercise 2.1: Attention Analysis

For the sentence: "The AI agent used the fraud detection tool to analyze suspicious transactions"

Draw or describe which words the word "analyzed" would attend to most strongly and why.

### Exercise 2.2: Memory Calculation

Calculate the VRAM required to run inference for:

1. Llama 3 8B in FP16
2. Llama 3 70B in INT8
3. GPT-4 (estimate, considering it's API-only)

Show your calculation methodology.

### Exercise 2.3: Context Window Impact

If you're building a RAG system with documents of 10,000 tokens each:
- How many documents can fit in a 128K context window?
- What infrastructure considerations change with larger context windows?

### Exercise 2.4: Architecture Comparison

Compare the following for a customer service chatbot:

| Factor | GPT-4 API | Self-hosted Llama 3 70B |
|--------|-----------|------------------------|
| Latency | | |
| Cost per 1K requests | | |
| Data privacy | | |
| Customization | | |
| Infrastructure complexity | | |

## Key Takeaways

- ✅ Transformers use self-attention to process sequences in parallel
- ✅ Tokenization converts text to model-processable integers
- ✅ Model size directly impacts VRAM requirements and latency
- ✅ Quantization enables larger models on limited hardware
- ✅ Architecture choice affects infrastructure decisions

## Next Steps

→ [IO-03: GPU Compute Fundamentals](../03-gpu-compute/README.md)

## Additional Resources

- [Attention Is All You Need (Paper)](https://arxiv.org/abs/1706.03762)
- [Transformer Explainer](https://transformer-explainer.github.io)
- [LLM Parameter Counting](https://blog.eleuther.ai/transformer-math/)
- [Quantization Techniques](https://arxiv.org/abs/2304.04367)
