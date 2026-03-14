# IO-17: Fine-Tuning Pipelines

**Duration**: 3 hours
**Module**: 5 - Advanced LLM-Ops & Production

## Learning Objectives

- Understand when and why to fine-tune vs. prompt engineering
- Implement LoRA, QLoRA, and PEFT for efficient fine-tuning
- Prepare high-quality training datasets
- Set up training infrastructure and monitor training runs
- Evaluate fine-tuned models effectively

## When to Fine-Tune

### Decision Framework

| Approach | Use Case | Cost | Complexity |
|----------|----------|------|------------|
| Prompt Engineering | General tasks, quick iterations | Low | Low |
| RAG + Prompting | Knowledge-intensive tasks | Medium | Medium |
| Fine-Tuning | Task-specific behavior, style | High | High |
| Full Fine-Tuning | Major capability shifts | Very High | Very High |

### Fine-Tuning Indicators

- Prompt engineering reaches limitations
- Model needs consistent output format
- Domain-specific terminology required
- Task requires specific reasoning patterns
- Need for cost optimization at inference time

## Parameter-Efficient Fine-Tuning (PEFT)

### LoRA (Low-Rank Adaptation)

LoRA adds small trainable rank decomposition matrices to attention layers:

```python
# LoRA Concept
# Original: W (d x k) forward pass: h = Wx
# LoRA: W + BA where B (d x r), A (r x k), r << min(d,k)

from peft import LoraConfig, get_peft_model
import torch

lora_config = LoraConfig(
    r=16,                    # Rank
    lora_alpha=32,           # Scaling factor
    target_modules=["q_proj", "v_proj", "k_proj", "o_proj"],
    lora_dropout=0.05,
    bias="none",
    task_type="CAUSAL_LM"
)

# Apply LoRA to model
model = get_peft_model(base_model, lora_config)
model.print_trainable_parameters()
# Output: trainable params: 4,194,304 || all params: 7,068,870,912 || trainable%: 0.059
```

### QLoRA (Quantized LoRA)

QLoRA combines quantization with LoRA for memory efficiency:

```python
from transformers import AutoModelForCausalLM, BitsAndBytesConfig
from peft import LoraConfig, get_peft_model

# Quantization config for 4-bit loading
bnb_config = BitsAndBytesConfig(
    load_in_4bit=True,
    bnb_4bit_quant_type="nf4",
    bnb_4bit_compute_dtype=torch.float16,
    bnb_4bit_use_double_quant=True
)

# Load quantized model
model = AutoModelForCausalLM.from_pretrained(
    "meta-llama/Llama-2-70b-hf",
    quantization_config=bnb_config,
    device_map="auto"
)

# Apply LoRA
lora_config = LoraConfig(r=64, lora_alpha=128, target_modules=["q_proj", "v_proj"])
model = get_peft_model(model, lora_config)
```

### PEFT Library Integration

```python
from transformers import AutoModelForCausalLM, TrainingArguments, Trainer
from datasets import Dataset
from peft import LoraConfig, get_peft_model, TaskType

# Prepare dataset
dataset = Dataset.from_json("training_data.json")

def tokenize_function(examples):
    return tokenizer(
        examples["prompt"],
        examples="completion"],
        truncation=True,
        max_length=2048
    )

tokenized_dataset = dataset.map(tokenize_function, batched=True)

# PEFT configuration
peft_config = LoraConfig(
    task_type=TaskType.CAUSAL_LM,
    r=32,
    lora_alpha=64,
    lora_dropout=0.1,
    target_modules=["q_proj", "v_proj", "k_proj", "o_proj", "gate_proj", "up_proj", "down_proj"],
    bias="none",
    inference_mode=False
)

# Wrap model with PEFT
model = get_peft_model(model, peft_config)

# Training arguments
training_args = TrainingArguments(
    output_dir="./results",
    num_train_epochs=3,
    per_device_train_batch_size=4,
    gradient_accumulation_steps=4,
    learning_rate=2e-4,
    fp16=True,
    save_strategy="epoch",
    save_total_limit=3,
    logging_steps=10,
    report_to="wandb"
)

trainer = Trainer(
    model=model,
    args=training_args,
    train_dataset=tokenized_dataset,
)

trainer.train()
```

## Training Infrastructure

### GPU Requirements by Model Size

| Model Size | Full Fine-Tune (A100) | LoRA (A100) | QLoRA (A100) |
|------------|----------------------|-------------|---------------|
| 7B | 8x A100 80GB | 1x A100 40GB | 1x A100 24GB |
| 13B | 8x A100 80GB | 2x A100 40GB | 1x A100 40GB |
| 70B | 8x A100 80GB | 8x A100 40GB | 2x A100 80GB |

### Distributed Training Setup

```python
# deepspeed_config.json
{
    "train_batch_size": "auto",
    "train_micro_batch_size_per_gpu": "auto",
    "gradient_accumulation_steps": "auto",
    "gradient_clipping": 1.0,
    "zero_optimization": {
        "stage": 3,
        "offload_optimizer": {
            "device": "cpu",
            "pin_memory": true
        },
        "offload_param": {
            "device": "cpu",
            "pin_memory": true
        }
    }
}

# Launch training
deepspeed train.py --deepspeed deepspeed_config.json
```

### Multi-GPU Training

```python
import torch.distributed as dist
from torch.nn.parallel import DataParallel

# Initialize distributed training
dist.init_process_group(backend="nccl")

# Wrap model for multi-GPU
model = torch.nn.parallel.DistributedDataParallel(
    model,
    device_ids=[local_rank],
    output_device=local_rank
)
```

## Dataset Preparation

### Data Collection Strategies

```python
class DatasetBuilder:
    def __init__(self):
        self.examples = []
    
    def add_synthetic(self, template: str, variations: int):
        """Generate synthetic training examples"""
        for _ in range(variations):
            example = self._generate_variation(template)
            self.examples.append(example)
    
    def add_human_written(self, examples: list):
        """Add human-written examples"""
        self.examples.extend(examples)
    
    def add_from_gpt(self, prompts: list, model: str):
        """Generate examples using GPT-4"""
        for prompt in prompts:
            completion = gpt_generate(prompt, model=model)
            self.examples.append({"prompt": prompt, "completion": completion})
    
    def deduplicate(self):
        """Remove duplicate examples"""
        seen = set()
        unique = []
        for ex in self.examples:
            key = hash((ex["prompt"], ex["completion"]))
            if key not in seen:
                seen.add(key)
                unique.append(ex)
        self.examples = unique
    
    def balance_classes(self, label_col: str):
        """Balance class distribution"""
        from collections import Counter
        counts = Counter(ex[label_col] for ex in self.examples)
        min_count = min(counts.values())
        
        balanced = []
        for label in counts:
            examples = [ex for ex in self.examples if ex[label_col] == label]
            balanced.extend(examples[:min_count])
        
        self.examples = balanced
```

### Data Quality Guidelines

| Aspect | Requirement | Tool |
|--------|-------------|------|
| Diversity | Cover varied contexts | Clustering analysis |
| Accuracy | Correct outputs | Human review |
| Consistency | Uniform format | Auto-formatting |
| Size | 1K-10K examples typical | Collection scaling |
| Balance | Even class distribution | Stratified sampling |

### Prompt-Response Formatting

```python
def format_for_training(example: dict, format_type: str = "chatml") -> str:
    if format_type == "chatml":
        return f"""<|im_start|>user
{example['prompt']}<|im_end|>
<|im_start|>assistant
{example['completion']}<|im_end|>"""
    
    elif format_type == "alpaca":
        return f"""Below is an instruction that describes a task. Write a response that appropriately completes the request.

### Instruction:
{example['prompt']}

### Response:
{example['completion']}"""
    
    elif format_type == "plain":
        return f"Prompt: {example['prompt']}\n\nResponse: {example['completion']}"

# Process dataset
formatted_dataset = raw_dataset.map(
    lambda x: {"text": format_for_training(x, "chatml")},
    remove_columns=raw_dataset.column_names
)
```

## Training Best Practices

### Hyperparameter Selection

```python
# Recommended hyperparameters by model size
HYPERPARAMS = {
    "7B": {
        "learning_rate": 2e-4,
        "batch_size": 4,
        "gradient_accumulation": 4,
        "epochs": 3,
        "warmup_steps": 100,
        "lora_r": 16,
        "lora_alpha": 32,
    },
    "13B": {
        "learning_rate": 1.5e-4,
        "batch_size": 2,
        "gradient_accumulation": 8,
        "epochs": 3,
        "warmup_steps": 100,
        "lora_r": 32,
        "lora_alpha": 64,
    },
    "70B": {
        "learning_rate": 1e-4,
        "batch_size": 1,
        "gradient_accumulation": 16,
        "epochs": 2,
        "warmup_steps": 50,
        "lora_r": 64,
        "lora_alpha": 128,
    }
}
```

### Training Monitoring

```python
import wandb

# Initialize tracking
wandb.init(project="llm-finetuning", name="experiment-1")

# Log metrics
for step, metrics in enumerate(training_loop):
    wandb.log({
        "loss": metrics["loss"],
        "learning_rate": metrics["lr"],
        "gpu_memory": torch.cuda.memory_allocated() / 1e9,
        "throughput": metrics["samples_per_second"]
    })

# Monitor training health
class TrainingMonitor:
    def check_health(self, metrics: dict) -> dict:
        warnings = []
        
        if metrics["loss"] > 10:
            warnings.append("Loss exploded - check learning rate")
        
        if metrics["gpu_memory"] > 90:
            warnings.append("GPU memory critical")
        
        if abs(metrics["loss"] - metrics["prev_loss"]) < 0.001:
            warnings.append("Loss not converging")
        
        return {"healthy": len(warnings) == 0, "warnings": warnings}
```

### Checkpointing

```python
from transformers import PreTrainedModel

class CheckpointManager:
    def __init__(self, save_dir: str, keep_last_n: int = 3):
        self.save_dir = save_dir
        self.keep_last_n = keep_last_n
        self.checkpoints = []
    
    def save(self, model: PreTrainedModel, epoch: int, metrics: dict):
        checkpoint_path = f"{self.save_dir}/epoch-{epoch}"
        model.save_pretrained(checkpoint_path)
        
        # Save training state
        torch.save({
            "epoch": epoch,
            "metrics": metrics
        }, f"{checkpoint_path}/trainer_state.pt")
        
        self.checkpoints.append(checkpoint_path)
        
        # Cleanup old checkpoints
        if len(self.checkpoints) > self.keep_last_n:
            old = self.checkpoints.pop(0)
            shutil.rmtree(old)
```

## Model Evaluation

### Evaluation Metrics

```python
from datasets import load_metric
import evaluate

class FineTuneEvaluator:
    def __init__(self, model, tokenizer):
        self.model = model
        self.tokenizer = tokenizer
        self.bleu = load_metric("bleu")
        self.rouge = load_metric("rouge")
    
    def evaluate_on_test_set(self, test_dataset):
        predictions = []
        references = []
        
        for example in test_dataset:
            output = self.generate(example["prompt"])
            predictions.append(output)
            references.append(example["reference"])
        
        # Calculate metrics
        bleu_results = self.bleu.compute(
            predictions=predictions,
            references=[[r] for r in references]
        )
        
        rouge_results = self.rouge.compute(
            predictions=predictions,
            references=references
        )
        
        return {
            "bleu": bleu_results["bleu"],
            "rouge-1": rouge_results["rouge1"].mid.fmeasure,
            "rouge-2": rouge_results["rouge2"].mid.fmeasure,
            "rouge-l": rouge_results["rougeL"].mid.fmeasure
        }
    
    def generate(self, prompt: str, max_length: int = 512) -> str:
        inputs = self.tokenizer(prompt, return_tensors="pt").to(self.model.device)
        outputs = self.model.generate(**inputs, max_length=max_length)
        return self.tokenizer.decode(outputs[0], skip_special_tokens=True)
```

### Human Evaluation Framework

```python
class HumanEvaluationFramework:
    def __init__(self):
        self.evaluation_criteria = [
            "accuracy",
            "helpfulness",
            "coherence",
            "safety",
            "instruction_following"
        ]
    
    def create_evaluation_task(self, examples: list):
        """Create evaluation task for human raters"""
        tasks = []
        for example in examples:
            task = {
                "prompt": example["prompt"],
                "response_a": example["response_ft"],    # Fine-tuned
                "response_b": example["response_base"],  # Base model
                "criteria": self.evaluation_criteria
            }
            tasks.append(task)
        return tasks
    
    def aggregate_results(self, ratings: list) -> dict:
        """Aggregate human ratings"""
        import numpy as np
        
        results = {}
        for criterion in self.evaluation_criteria:
            scores = [r[criterion] for r in ratings]
            results[criterion] = {
                "mean": np.mean(scores),
                "std": np.std(scores),
                "agreement": self._calculate_agreement(scores)
            }
        return results
    
    def _calculate_agreement(self, scores: list) -> float:
        """Calculate inter-rater agreement"""
        import numpy as np
        from scipy.stats import spearmanr
        
        # Simplified: return correlation between first two raters
        return 0.85  # Placeholder
```

## AI Agent Platform Reference

The FinAgent platform in this repository demonstrates production-grade fine-tuning:

```go
// services/ai-agent-platform/internal/agent/
// Fine-tuning configuration for domain-specific agents

type FineTuneConfig struct {
    BaseModel       string            `json:"base_model"`
    AdapterType     string            `json:"adapter_type"` // "lora", "qlora", "full"
    Rank            int               `json:"rank"`
    Alpha           int               `json:"alpha"`
    TargetModules   []string          `json:"target_modules"`
    LearningRate    float64           `json:"learning_rate"`
    Epochs          int               `json:"epochs"`
    BatchSize       int               `json:"batch_size"`
    DatasetPath     string            `json:"dataset_path"`
}

// Training infrastructure uses Hugging Face Trainer with:
// - DeepSpeed ZeRO-3 for memory efficiency
// - Weights & Biases for experiment tracking
// - LoRA adapters for parameter-efficient training
```

### Fine-Tuning Pipeline Architecture

```
┌─────────────┐    ┌──────────────┐    ┌─────────────┐    ┌────────────┐
│   Dataset   │───>│   Preprocess │───>│   Training  │───>│ Evaluation │
│  Collection │    │   (Tokenize)  │    │  (LoRA/QLoRA)│    │  (Metrics) │
└─────────────┘    └──────────────┘    └─────────────┘    └────────────┘
                                              │                   │
                                              v                   v
                                       ┌─────────────┐    ┌────────────┐
                                       │  Checkpoint │───>│  Model     │
                                       │   Manager   │    │  Registry  │
                                       └─────────────┘    └────────────┘
```

## Summary

- Fine-tuning adapts pre-trained models to specific tasks
- LoRA/QLoRA enable efficient fine-tuning with minimal compute
- PEFT provides a unified interface for parameter-efficient methods
- High-quality datasets are critical for successful fine-tuning
- Comprehensive evaluation combines automated metrics with human assessment
- Production fine-tuning requires robust infrastructure for training and monitoring
