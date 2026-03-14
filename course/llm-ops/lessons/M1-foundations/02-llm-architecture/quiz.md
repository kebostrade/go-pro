# Quiz: LLM Architecture Fundamentals

## Question 1

What is the key innovation of the Transformer architecture?

A) Recurrent neural networks
B) Convolutional layers
C) Self-attention mechanism
D) Dropout regularization

## Question 2

How does self-attention help LLMs process text?

A) It processes words sequentially one at a time
B) It allows the model to weigh relationships between all words simultaneously
C) It reduces the number of parameters needed
D) It eliminates the need for embeddings

## Question 3

What is a "token" in the context of LLMs?

A) A complete sentence
B) A word or subword unit that the model processes
C) A type of neural network layer
D) A model parameter

## Question 4

What does "context window" refer to?

A) The size of the computer screen for displaying outputs
B) The maximum number of tokens the model can process at once
C) The amount of training data used
D) The number of parameters in the model

## Question 5

If a model has 70 billion parameters and uses FP16 precision, approximately how much VRAM is needed just to load the model weights?

A) 70 GB
B) 140 GB
C) 35 GB
D) 280 GB

## Question 6

What is quantization in the context of LLMs?

A) Converting text to numbers
B) Reducing the precision of model weights to save memory
C) Adding more parameters to the model
D) Training the model on more data

## Question 7

Which architecture type is GPT (Generative Pre-trained Transformer)?

A) Encoder-only
B) Encoder-decoder
C) Decoder-only
D) Hybrid

## Question 8

What is the main advantage of INT4 quantization over FP16?

A) Higher accuracy
B) 4x reduction in memory usage
C) Faster training
D) Better tokenization

## Question 9

Why do larger context windows require more VRAM?

A) Because more tokens need to be stored in the KV cache
B) Because the model itself becomes larger
C) Because quantization becomes less effective
D) Because tokenization is slower

## Question 10

What is the primary function of positional encoding in Transformers?

A) To add noise to prevent overfitting
B) To indicate the position of each token in the sequence
C) To compress the input
D) To encode the meaning of words

## Question 11

In the AI Agent Platform, the LLM provider abstraction allows:

A) Training different models simultaneously
B) Swapping between OpenAI, Anthropic, and Ollama without changing application code
C) Reducing the number of parameters
D) Eliminating the need for GPUs

## Question 12

What happens during the "de-tokenization" step of LLM inference?

A) The input text is converted to tokens
B) The model's numerical output is converted back to human-readable text
C) The model weights are loaded into memory
D) The attention weights are calculated

## Question 13

Which component of the transformer architecture processes the attention outputs through dense layers?

A) Self-attention
B) Layer normalization
C) Feed-forward network
D) Positional encoding

## Question 14

A model with 8 billion parameters running on a single GPU would typically use which precision for efficient inference?

A) FP32
B) FP16 or BF16
C) INT2
D) Binary (1 bit)
