# Quiz: Fine-Tuning Pipelines

## Question 1

What is the main advantage of LoRA over full fine-tuning?

A) Higher accuracy
B) Fewer trainable parameters
C) Faster inference
D) Larger model capacity

## Question 2

What does QLoRA combine?

A) Quantization and fine-tuning
B) LoRA and RAG
C) Prefix tuning and LoRA
D) 8-bit training with full fine-tuning

## Question 3

Which PEFT method adds rank decomposition matrices to attention layers?

A) Prefix Tuning
B) LoRA
C) Adapter
D) Prompt Tuning

## Question 4

What is a typical LoRA rank (r) value for a 7B model?

A) 1-4
B) 4-8
C) 8-32
D) 64-128

## Question 5

What is the recommended dataset size for fine-tuning?

A) 100-500 examples
B) 1K-10K examples
C) 100K-1M examples
D) More is always better

## Question 6

Which metric is NOT typically used for LLM fine-tuning evaluation?

A) BLEU
B) ROUGE
C) Perplexity
D) RMSE

## Question 7

What does the lora_alpha parameter control in LoRA?

A) Learning rate scaling
B) Dropout rate
C) Rank
D) Attention heads

## Question 8

Which target_modules are typically used for LoRA in causal language models?

A) Only q_proj
B) q_proj and v_proj
C) All linear layers
D) Embedding layers

## Question 9

What is the purpose of gradient accumulation?

A) Reduce memory usage
B) Increase batch size with limited GPU memory
C) Speed up training
D) Improve gradient accuracy

## Question 10

What does DeepSpeed ZeRO-3 provide?

A) Better optimizers
B) Memory-efficient distributed training
C) Faster tokenization
D) Automatic checkpointing

## Question 11

Which format is commonly used for instruction fine-tuning?

A) CSV only
B) JSON only
C) ChatML or Alpaca format
D) XML only

## Question 12

What is the recommended learning rate for LoRA fine-tuning?

A) 1e-5
B) 1e-3
C) 1e-4 to 3e-4
D) 1e-2

## Question 13

What should be monitored during training to detect convergence issues?

A) Only loss
B) Loss, learning rate, and GPU memory
C) Only GPU memory
D) Only throughput

## Question 14

What is catastrophic forgetting?

A) Model forgetting during inference
B) Model losing original capabilities when fine-tuned
C) Training instability
D) Overfitting to small dataset

## Question 15

What is the purpose of warmup steps in training?

A) Increase learning rate gradually
B) Stabilize early training
C) Reduce memory usage
D) Speed up convergence
