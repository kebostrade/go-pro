# Exercise: Fine-Tuning Pipelines

## Problem 1: LoRA Configuration

Implement a LoRA configuration optimizer:

```python
from peft import LoraConfig, TaskType
from typing import Optional

class LoRAOptimizer:
    def __init__(self, model_size: str):
        self.model_size = model_size
        self.configs = {
            "7B": {"r": 16, "alpha": 32, "dropout": 0.05},
            "13B": {"r": 32, "alpha": 64, "dropout": 0.05},
            "70B": {"r": 64, "alpha": 128, "dropout": 0.1},
        }
    
    def get_config(self, custom_r: Optional[int] = None) -> LoRAConfig:
        """
        Get LoRA configuration for model size
        Customize rank if specified
        """
        # Your implementation
        pass
    
    def estimate_memory(self, config: LoRAConfig, batch_size: int) -> dict:
        """
        Estimate GPU memory requirements
        Returns: dict with memory estimates in GB
        """
        # Your implementation
        pass

# Test
optimizer = LoRAOptimizer("7B")
config = optimizer.get_config()
print(f"Rank: {config.r}, Alpha: {config.lora_alpha}")

memory = optimizer.estimate_memory(config, batch_size=4)
print(f"Estimated memory: {memory}")
```

---

## Problem 2: Dataset Preprocessor

Implement a dataset preprocessing pipeline:

```python
class DatasetPreprocessor:
    def __init__(self, tokenizer, format_type: str = "chatml"):
        self.tokenizer = tokenizer
        self.format_type = format_type
    
    def load_and_preprocess(self, data_path: str, max_length: int = 2048) -> dict:
        """
        Load JSON data and preprocess for training
        """
        # Your implementation
        pass
    
    def format_example(self, prompt: str, completion: str) -> str:
        """
        Format prompt/completion based on format_type
        """
        # Your implementation
        pass
    
    def tokenize(self, text: str, max_length: int) -> dict:
        """
        Tokenize text with proper truncation
        """
        # Your implementation
        pass
    
    def split_train_val(self, dataset, val_ratio: float = 0.1):
        """
        Split dataset into train/validation
        """
        # Your implementation
        pass

# Test with sample data
from transformers import AutoTokenizer
tokenizer = AutoTokenizer.from_pretrained("gpt2")

preprocessor = DatasetPreprocessor(tokenizer, "chatml")
sample_data = [
    {"prompt": "What is Python?", "completion": "Python is a programming language."},
    {"prompt": "Define machine learning", "completion": "ML is a subset of AI."}
]

# Process data
```

---

## Problem 3: Training Monitor

Implement a training monitor with health checks:

```python
import time
from typing import Dict, List

class TrainingMonitor:
    def __init__(self, window_size: int = 10):
        self.window_size = window_size
        self.history = []
        self.start_time = None
    
    def start(self):
        """Start monitoring"""
        self.start_time = time.time()
    
    def log_step(self, step: int, loss: float, lr: float, gpu_memory: float):
        """Log training step metrics"""
        # Your implementation
        pass
    
    def check_health(self) -> Dict:
        """
        Check training health
        Returns: dict with health status and warnings
        """
        # Your implementation
        pass
    
    def get_summary(self) -> Dict:
        """
        Get training summary
        """
        # Your implementation
        pass

# Test
monitor = TrainingMonitor(window_size=5)
monitor.start()

for step in range(20):
    loss = 2.0 - (step * 0.08) + (hash(str(step)) % 100) * 0.001
    monitor.log_step(step, loss, 2e-4, 20.5)

health = monitor.check_health()
print(f"Health: {health}")
```

---

## Problem 4: Evaluation Pipeline

Implement evaluation metrics computation:

```python
class FineTuneEvaluator:
    def __init__(self, model, tokenizer):
        self.model = model
        self.tokenizer = tokenizer
    
    def compute_exact_match(self, predictions: List[str], references: List[str]) -> float:
        """Compute exact match accuracy"""
        # Your implementation
        pass
    
    def compute_bleu(self, predictions: List[str], references: List[str]) -> float:
        """Compute BLEU score"""
        # Your implementation
        pass
    
    def compute_rouge(self, predictions: List[str], references: List[str]) -> Dict:
        """Compute ROUGE scores"""
        # Your implementation
        pass
    
    def compute_perplexity(self, texts: List[str]) -> float:
        """Compute perplexity on text"""
        # Your implementation
        pass
    
    def evaluate(self, test_data: List[dict]) -> Dict:
        """
        Run full evaluation
        test_data: [{"prompt": ..., "reference": ...}]
        """
        # Your implementation
        pass

# Test with sample data
```

---

## Problem 5: Checkpoint Manager

Implement a checkpoint management system:

```python
import os
import shutil
from typing import List, Optional

class CheckpointManager:
    def __init__(self, checkpoint_dir: str, keep_last_n: int = 3):
        self.checkpoint_dir = checkpoint_dir
        self.keep_last_n = keep_last_n
        self.checkpoints = []
    
    def save_checkpoint(self, model, optimizer, epoch: int, metrics: dict):
        """
        Save model checkpoint with metadata
        """
        # Your implementation
        pass
    
    def load_checkpoint(self, checkpoint_path: str):
        """
        Load checkpoint and return model, optimizer, state
        """
        # Your implementation
        pass
    
    def list_checkpoints(self) -> List[dict]:
        """
        List all available checkpoints with metadata
        """
        # Your implementation
        pass
    
    def get_best_checkpoint(self, metric: str = "loss") -> Optional[str]:
        """
        Get path to best checkpoint by metric
        """
        # Your implementation
        pass
    
    def cleanup_old(self):
        """
        Remove old checkpoints keeping only last N
        """
        # Your implementation
        pass

# Test
manager = CheckpointManager("./checkpoints", keep_last_n=2)
```

---

## Problem 6: Training Hyperparameter Search

Design a hyperparameter search:

```python
from itertools import product

class HyperparameterSearcher:
    def __init__(self, param_grid: dict):
        self.param_grid = param_grid
        self.results = []
    
    def generate_configs(self) -> List[dict]:
        """
        Generate all parameter combinations
        """
        # Your implementation
        pass
    
    def run_trial(self, config: dict, train_data, eval_data) -> dict:
        """
        Run single trial with config
        Returns: metrics
        """
        # Your implementation
        pass
    
    def search(self, train_data, eval_data, max_trials: int = 10) -> dict:
        """
        Run hyperparameter search
        """
        # Your implementation
        pass
    
    def get_best_config(self) -> dict:
        """
        Get best configuration found
        """
        # Your implementation
        pass

# Define search space
param_grid = {
    "learning_rate": [1e-4, 2e-4, 3e-4],
    "lora_r": [8, 16, 32],
    "lora_alpha": [16, 32, 64],
    "batch_size": [2, 4, 8]
}

searcher = HyperparameterSearcher(param_grid)
```

---

## Problem 7: Fine-Tuning Decision Matrix

Create a decision framework:

```python
class FineTuneDecisionMaker:
    def should_finetune(
        self,
        current_performance: float,
        target_performance: float,
        dataset_size: int,
        budget_hours: float,
        has_domain_data: bool
    ) -> dict:
        """
        Decide whether to fine-tune
        Returns: dict with decision and reasoning
        """
        # Your implementation
        pass
    
    def recommend_approach(
        self,
        model_size: str,
        gpu_memory_gb: int,
        dataset_size: int
    ) -> dict:
        """
        Recommend fine-tuning approach
        Returns: dict with approach and config
        """
        # Your implementation
        pass
    
    def estimate_cost(
        self,
        approach: str,
        model_size: str,
        dataset_size: int,
        epochs: int
    ) -> dict:
        """
        Estimate cost and time
        """
        # Your implementation
        pass

# Test scenarios
maker = FineTuneDecisionMaker()

# Scenario 1: Startup with limited resources
result1 = maker.should_finetune(
    current_performance=0.7,
    target_performance=0.85,
    dataset_size=5000,
    budget_hours=20,
    has_domain_data=True
)

# Scenario 2: Enterprise with large budget
result2 = maker.should_finetune(
    current_performance=0.6,
    target_performance=0.95,
    dataset_size=50000,
    budget_hours=200,
    has_domain_data=True
)

print(f"Scenario 1: {result1}")
print(f"Scenario 2: {result2}")
```

---

## Submission

Save your implementations and be prepared to discuss:
- How LoRA rank affects model quality vs. memory
- Dataset preprocessing best practices
- Training monitoring strategies
- Evaluation methodology selection
