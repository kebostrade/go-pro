# Exercise: Model Versioning & MLOps

## Problem 1: Model Registry Implementation

Implement a simple model registry:

```python
import json
import os
from datetime import datetime
from typing import Optional, List, Dict

class ModelRegistry:
    def __init__(self, registry_path: str = "./model_registry"):
        self.registry_path = registry_path
        self.models_file = os.path.join(registry_path, "models.json")
        os.makedirs(registry_path, exist_ok=True)
        self.models = self._load()
    
    def _load(self) -> dict:
        """Load registry from disk"""
        # Your implementation
        pass
    
    def _save(self):
        """Save registry to disk"""
        # Your implementation
        pass
    
    def register(
        self,
        name: str,
        version: str,
        model_path: str,
        metrics: dict,
        params: dict
    ) -> dict:
        """Register a new model"""
        # Your implementation
        pass
    
    def get(self, name: str, version: str = None) -> Optional[dict]:
        """Get model by name and version"""
        # Your implementation
        pass
    
    def list_models(self) -> List[dict]:
        """List all registered models"""
        # Your implementation
        pass
    
    def list_versions(self, name: str) -> List[dict]:
        """List all versions of a model"""
        # Your implementation
        pass
    
    def update_stage(self, name: str, version: str, stage: str):
        """Update model stage"""
        # Your implementation
        pass

# Test
registry = ModelRegistry()

# Register models
registry.register(
    name="sentiment-classifier",
    version="1.0.0",
    model_path="models/sentiment-v1.pt",
    metrics={"accuracy": 0.92, "f1": 0.91},
    params={"model_type": "bert", "max_length": 512}
)

registry.register(
    name="sentiment-classifier",
    version="1.1.0",
    model_path="models/sentiment-v1.1.pt",
    metrics={"accuracy": 0.95, "f1": 0.94},
    params={"model_type": "bert", "max_length": 512, "lora_rank": 16}
)

# List and get
print(registry.list_versions("sentiment-classifier"))
```

---

## Problem 2: Experiment Tracker

Implement an experiment tracking system:

```python
import time
import uuid
from typing import Dict, List, Optional

class ExperimentTracker:
    def __init__(self, storage_path: str = "./experiments"):
        self.storage_path = storage_path
        self.experiments = []
        os.makedirs(storage_path, exist_ok=True)
    
    def start_run(
        self,
        experiment_name: str,
        run_name: str,
        params: Dict
    ) -> str:
        """Start a new experiment run"""
        # Your implementation
        pass
    
    def log_metric(self, run_id: str, metric_name: str, value: float):
        """Log a metric"""
        # Your implementation
        pass
    
    def log_metrics(self, run_id: str, metrics: Dict):
        """Log multiple metrics"""
        # Your implementation
        pass
    
    def end_run(self, run_id: str, status: str = "completed"):
        """End an experiment run"""
        # Your implementation
        pass
    
    def get_run(self, run_id: str) -> Optional[dict]:
        """Get run details"""
        # Your implementation
        pass
    
    def search_runs(
        self,
        experiment_name: str,
        filter_metrics: Dict = None
    ) -> List[dict]:
        """Search runs by experiment and metrics"""
        # Your implementation
        pass
    
    def compare_runs(self, run_ids: List[str], metric: str) -> dict:
        """Compare multiple runs by metric"""
        # Your implementation
        pass

# Test
tracker = ExperimentTracker()

run_id = tracker.start_run(
    experiment_name="llm-finetuning",
    run_name="lora-rank-16",
    params={"lora_rank": 16, "learning_rate": 2e-4, "epochs": 3}
)

for epoch in range(3):
    tracker.log_metric(run_id, f"epoch_{epoch}_loss", 1.0 - (epoch * 0.3))
    tracker.log_metric(run_id, f"epoch_{epoch}_accuracy", 0.7 + (epoch * 0.1))

tracker.end_run(run_id)
```

---

## Problem 3: CI/CD Pipeline Config

Create a CI/CD configuration for ML:

```yaml
# Create .github/workflows/ml-training.yml
# Include:
# 1. Trigger on data/training code changes
# 2. Setup Python and GPU
# 3. Install dependencies
# 4. Validate data
# 5. Train model
# 6. Evaluate model
# 7. Register model if metrics improve
# 8. Run integration tests

# Your implementation below:

name: ML Training Pipeline

on:
  push:
    branches: [main]
    paths:
      - 'training/**'
      - 'data/**'
      - 'models/**'

jobs:
  validate:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      # Add validation steps

  train:
    needs: validate
    runs-on: gpu-runner
    steps:
      # Add training steps

  test:
    needs: train
    runs-on: ubuntu-latest
    steps:
      # Add testing steps

  deploy:
    needs: test
    if: github.ref == 'refs/heads/main'
    runs-on: ubuntu-latest
    steps:
      # Add deployment steps
```

---

## Problem 4: Data Validation

Implement data validation checks:

```python
import pandas as pd
from typing import List, Dict

class DataValidator:
    def __init__(self):
        self.validation_results = []
    
    def validate_schema(
        self,
        df: pd.DataFrame,
        required_columns: List[str],
        column_types: Dict[str, str]
    ) -> Dict:
        """Validate data schema"""
        # Your implementation
        pass
    
    def validate_quality(self, df: pd.DataFrame) -> Dict:
        """Validate data quality"""
        # Your implementation
        pass
    
    def validate_distribution(
        self,
        train_df: pd.DataFrame,
        test_df: pd.DataFrame,
        column: str,
        threshold: float = 0.1
    ) -> Dict:
        """Validate train/test distribution"""
        # Your implementation
        pass
    
    def run_all(self, train_df: pd.DataFrame, test_df: pd.DataFrame) -> Dict:
        """Run all validations"""
        # Your implementation
        pass

# Test
train_df = pd.DataFrame({
    "prompt": ["Hello", "World", "Test"],
    "completion": ["Hi", "Earth", "Demo"],
    "category": ["greeting", "fact", "demo"]
})

test_df = pd.DataFrame({
    "prompt": ["Hi", "Test"],
    "completion": ["Hello", "Testing"],
    "category": ["greeting", "demo"]
})

validator = DataValidator()
results = validator.run_all(train_df, test_df)
print(results)
```

---

## Problem 5: Model Lifecycle Manager

Implement model lifecycle management:

```python
class LifecycleManager:
    def __init__(self, registry):
        self.registry = registry
        self.stages = ["development", "staging", "production", "archived"]
    
    def deploy_to_stage(
        self,
        model_name: str,
        version: str,
        target_stage: str
    ) -> Dict:
        """Deploy model to a stage"""
        # Your implementation
        pass
    
    def promote(self, model_name: str, version: str) -> Dict:
        """Promote model to next stage"""
        # Your implementation
        pass
    
    def rollback(self, model_name: str) -> Dict:
        """Rollback to previous production model"""
        # Your implementation
        pass
    
    def archive(self, model_name: str, version: str):
        """Archive a model version"""
        # Your implementation
        pass
    
    def get_production_model(self, model_name: str) -> Dict:
        """Get current production model"""
        # Your implementation
        pass

# Test with registry from Problem 1
manager = LifecycleManager(registry)

# Deploy to stages
manager.deploy_to_stage("sentiment-classifier", "1.0.0", "production")
manager.deploy_to_stage("sentiment-classifier", "1.1.0", "staging")

# Check production
prod = manager.get_production_model("sentiment-classifier")
print(f"Production: {prod}")
```

---

## Problem 6: A/B Test Manager

Implement A/B testing for models:

```python
import hashlib
from typing import List, Dict
import random

class ABTestManager:
    def __init__(self):
        self.experiments = {}
    
    def create_experiment(
        self,
        name: str,
        model_a_config: Dict,
        model_b_config: Dict,
        traffic_split: float = 0.5
    ):
        """Create A/B test experiment"""
        # Your implementation
        pass
    
    def get_model_for_request(
        self,
        experiment_name: str,
        user_id: str
    ) -> str:
        """Get model (a or b) for request"""
        # Your implementation
        pass
    
    def record_outcome(
        self,
        experiment_name: str,
        model_variant: str,
        user_id: str,
        outcome: Dict
    ):
        """Record experiment outcome"""
        # Your implementation
        pass
    
    def get_results(self, experiment_name: str) -> Dict:
        """Get experiment results"""
        # Your implementation
        pass

# Test
ab_test = ABTestManager()

ab_test.create_experiment(
    name="model-comparison",
    model_a_config={"name": "model-v1", "type": "base"},
    model_b_config={"name": "model-v2", "type": "fine-tuned"},
    traffic_split=0.5
)

# Simulate traffic
for user_id in [f"user_{i}" for i in range(100)]:
    model = ab_test.get_model_for_request("model-comparison", user_id)
    
    # Record outcome (simulated)
    outcome = {"success": random.random() > 0.1, "latency_ms": random.randint(50, 200)}
    ab_test.record_outcome("model-comparison", model, user_id, outcome)

results = ab_test.get_results("model-comparison")
print(results)
```

---

## Problem 7: Experiment Comparison

Create an experiment comparison dashboard:

```python
def create_comparison_report(experiment_runs: List[dict]) -> str:
    """
    Create markdown comparison report from experiment runs
    """
    # Your implementation
    pass

# Sample data
runs = [
    {
        "run_id": "run_1",
        "params": {"lora_rank": 8, "lr": 2e-4},
        "metrics": {"val_loss": 0.8, "accuracy": 0.88}
    },
    {
        "run_id": "run_2", 
        "params": {"lora_rank": 16, "lr": 2e-4},
        "metrics": {"val_loss": 0.6, "accuracy": 0.92}
    },
    {
        "run_id": "run_3",
        "params": {"lora_rank": 32, "lr": 1e-4},
        "metrics": {"val_loss": 0.5, "accuracy": 0.94}
    }
]

report = create_comparison_report(runs)
print(report)
```

---

## Submission

Complete all implementations and be prepared to discuss:
- How model registries enable reproducibility
- CI/CD best practices for ML
- A/B testing methodology for model comparison
- Model lifecycle management strategies
