# IO-18: Model Versioning & MLOps

**Duration**: 2 hours
**Module**: 5 - Advanced LLM-Ops & Production

## Learning Objectives

- Implement model versioning and a model registry
- Set up experiment tracking with MLflow or Weights & Biases
- Build CI/CD pipelines for ML workflows
- Manage model lifecycle from development to production
- Establish reproducible training pipelines

## Model Registry

### What is a Model Registry?

A centralized system to store, version, and manage ML models:

```python
# Model Registry Concept
class ModelRegistry:
    def __init__(self, registry_path: str):
        self.registry_path = registry_path
        self.models = {}
    
    def register(self, name: str, version: str, model_path: str, metadata: dict):
        """Register a new model version"""
        model_entry = {
            "name": name,
            "version": version,
            "model_path": model_path,
            "metadata": metadata,
            "registered_at": datetime.now(),
            "status": "staging"
        }
        self.models[f"{name}:{version}"] = model_entry
    
    def get(self, name: str, version: str = None) -> dict:
        """Get model by name and version"""
        key = f"{name}:{version}" if version else name
        return self.models.get(key)
    
    def list_versions(self, name: str) -> list:
        """List all versions of a model"""
        return [v for k, v in self.models.items() if k.startswith(name)]
    
    def promote(self, name: str, version: str, stage: str):
        """Promote model to stage (staging/production/archived)"""
        key = f"{name}:{version}"
        if key in self.models:
            self.models[key]["status"] = stage
```

### MLflow Model Registry

```python
import mlflow
from mlflow.tracking import MlflowClient

# Configure MLflow
mlflow.set_tracking_uri("http://localhost:5000")
mlflow.set_experiment("llm-experiments")

client = MlflowClient()

# Register model
model_uri = "runs:/run_id/model"
model_name = "finagent-classifier"

# Create registered model
mlflow.register_model(model_uri, model_name)

# Transition model through stages
client.transition_model_version_stage(
    name=model_name,
    version=1,
    stage="production"
)

# Get latest production model
latest = client.get_latest_model_versions(name=model_name, stages=["production"])[0]
print(f"Production model: {latest.name} version {latest.version}")
```

### Weights & Biases Model Registry

```python
import wandb

# Initialize W&B
wandb.init(project="llm-finetuning")

# Log model artifacts
with wandb.init(project="llm-project", job_type="train"):
    # Log training metrics
    wandb.log({"loss": 0.5, "accuracy": 0.9})
    
    # Save model as artifact
    artifact = wandb.Artifact("model", type="model")
    artifact.add_file("model.pt")
    wandb.log_artifact(artifact)

# Load model from artifact
artifact = wandb.use_artifact("model:v0")
artifact.download()
```

## Experiment Tracking

### Setting Up Experiments

```python
import mlflow
import numpy as np

class ExperimentTracker:
    def __init__(self, experiment_name: str):
        self.experiment_name = experiment_name
        mlflow.set_experiment(experiment_name)
    
    def start_run(self, run_name: str, params: dict):
        """Start a new experiment run"""
        self.run = mlflow.start_run(run_name=run_name)
        mlflow.log_params(params)
    
    def log_metrics(self, metrics: dict):
        """Log metrics for current step"""
        mlflow.log_metrics(metrics)
    
    def log_artifacts(self, artifacts_dir: str):
        """Log artifacts (models, visualizations)"""
        mlflow.log_artifacts(artifacts_dir)
    
    def end_run(self, status: str = "FINISHED"):
        """End the current run"""
        mlflow.end_run(status=status)
    
    def search_runs(self, filter_string: str = None, max_results: int = 100):
        """Search for runs in experiment"""
        return mlflow.search_runs(
            filter_string=filter_string,
            max_results=max_results
        )

# Usage
tracker = ExperimentTracker("fine-tuning-experiments")

params = {
    "model_size": "7B",
    "lora_rank": 16,
    "learning_rate": 2e-4,
    "batch_size": 4,
    "epochs": 3
}

tracker.start_run("experiment-1", params)

for epoch in range(3):
    # Train
    loss = 1.0 - (epoch * 0.3)
    accuracy = 0.7 + (epoch * 0.08)
    
    tracker.log_metrics({
        "epoch": epoch,
        "loss": loss,
        "accuracy": accuracy,
        "perplexity": np.exp(loss)
    })

tracker.end_run()
```

### Comparing Experiments

```python
class ExperimentComparator:
    def __init__(self, experiment_name: str):
        self.experiment_name = experiment_name
    
    def compare_runs(self, metric: str = "val_accuracy"):
        """Compare runs by metric"""
        runs = mlflow.search_runs(
            experiment_names=[self.experiment_name],
            order_by=[f"metrics.{metric} DESC"],
            max_results=10
        )
        
        comparison = runs[["run_id", "params.lora_rank", "params.learning_rate", f"metrics.{metric}"]]
        return comparison
    
    def get_best_run(self, metric: str) -> dict:
        """Get best run by metric"""
        runs = mlflow.search_runs(
            experiment_names=[self.experiment_name],
            order_by=[f"metrics.{metric} DESC"],
            max_results=1
        )
        return runs.iloc[0].to_dict()
```

## CI/CD for ML

### ML Pipeline Architecture

```yaml
# .github/workflows/ml-pipeline.yml
name: ML Training Pipeline

on:
  push:
    branches: [main]
    paths:
      - 'training/**'
      - 'data/**'

jobs:
  train:
    runs-on: gpu-runner
    
    steps:
      - uses: actions/checkout@v3
      
      - name: Set up Python
        uses: actions/setup-python@v4
        with:
          python-version: '3.10'
      
      - name: Install dependencies
        run: |
          pip install -r requirements.txt
      
      - name: Run data validation
        run: python scripts/validate_data.py
      
      - name: Train model
        run: python scripts/train.py
        env:
          WANDB_API_KEY: ${{ secrets.WANDB_API_KEY }}
      
      - name: Evaluate model
        run: python scripts/evaluate.py
      
      - name: Register model
        if: github.ref == 'refs/heads/main'
        run: |
          python scripts/register_model.py
      
      - name: Upload artifacts
        uses: actions/upload-artifact@v3
        with:
          name: model-artifacts
          path: models/
```

### Data Validation Pipeline

```python
class DataValidator:
    def __init__(self):
        self.checks = []
    
    def check_schema(self, df, expected_columns: list):
        """Validate data schema"""
        missing = set(expected_columns) - set(df.columns)
        if missing:
            raise ValueError(f"Missing columns: {missing}")
        return True
    
    def check_quality(self, df):
        """Check data quality"""
        issues = []
        
        # Check for nulls
        null_counts = df.isnull().sum()
        if null_counts.any():
            issues.append(f"Null values found: {null_counts[null_counts > 0].to_dict()}")
        
        # Check for duplicates
        dup_count = df.duplicated().sum()
        if dup_count > 0:
            issues.append(f"{dup_count} duplicate rows found")
        
        return issues
    
    def check_distribution(self, train_df, test_df, column: str):
        """Check distribution shift"""
        from scipy import stats
        
        stat, p_value = stats.ks_2samp(train_df[column], test_df[column])
        
        return {
            "statistic": stat,
            "p_value": p_value,
            "significant_shift": p_value < 0.05
        }

# Integration with CI
def run_validation():
    validator = DataValidator()
    
    train_df = pd.read_csv("data/train.csv")
    test_df = pd.read_csv("data/test.csv")
    
    # Run checks
    validator.check_schema(train_df, ["prompt", "completion", "category"])
    issues = validator.check_quality(train_df)
    
    if issues:
        print(f"Validation issues: {issues}")
        return False
    
    return True
```

### Model Testing in CI

```python
class ModelTester:
    def __init__(self, model_path: str):
        self.model = load_model(model_path)
        self.tokenizer = load_tokenizer(model_path)
    
    def test_inference_speed(self, test_prompts: list) -> dict:
        """Test inference speed"""
        import time
        
        times = []
        for prompt in test_prompts:
            start = time.time()
            self.model.generate(prompt)
            times.append(time.time() - start)
        
        return {
            "mean_latency_ms": np.mean(times) * 1000,
            "p50_latency_ms": np.percentile(times, 50) * 1000,
            "p99_latency_ms": np.percentile(times, 99) * 1000
        }
    
    def test_output_quality(self, test_cases: list) -> dict:
        """Test output quality"""
        results = []
        for case in test_cases:
            output = self.model.generate(case["prompt"])
            results.append({
                "expected": case["expected"],
                "actual": output,
                "match": self._calculate_similarity(output, case["expected"])
            })
        
        accuracy = sum(1 for r in results if r["match"] > 0.8) / len(results)
        return {"accuracy": accuracy, "details": results}
    
    def test_safety(self, test_prompts: list) -> dict:
        """Test safety compliance"""
        violations = []
        for prompt in test_prompts:
            if self._detect_violation(prompt):
                violations.append(prompt)
        
        return {
            "total_tested": len(test_prompts),
            "violations": len(violations),
            "safe": len(violations) == 0
        }
```

## Model Lifecycle Management

### Lifecycle Stages

```
┌─────────┐   ┌──────────┐   ┌───────────┐   ┌────────────┐   ┌──────────┐
│  Data   │──>│ Training │──>│ Validation│──>│ Staging   │──>│Production│
│  Prep   │   │          │   │           │   │           │   │          │
└─────────┘   └──────────┘   └───────────┘   └────────────┘   └──────────┘
                                                                     │
                                                                     v
                                                              ┌──────────┐
                                                              │Archived  │
                                                              └──────────┘
```

### Stage Management

```python
class ModelLifecycleManager:
    def __init__(self, registry_client):
        self.client = registry_client
        self.stages = ["none", "staging", "production", "archived"]
    
    def deploy_to_staging(self, model_name: str, version: int):
        """Deploy model to staging"""
        self.client.transition_model_version_stage(
            name=model_name,
            version=version,
            stage="staging"
        )
    
    def promote_to_production(self, model_name: str, version: int):
        """Promote model to production"""
        # First, demote current production model
        current = self.client.get_latest_model_versions(
            name=model_name,
            stages=["production"]
        )
        
        if current:
            self.client.transition_model_version_stage(
                name=model_name,
                version=current[0].version,
                stage="archived"
            )
        
        # Promote new model
        self.client.transition_model_version_stage(
            name=model_name,
            version=version,
            stage="production"
        )
    
    def rollback(self, model_name: str):
        """Rollback to previous production model"""
        # Get archived models
        archived = self.client.get_latest_model_versions(
            name=model_name,
            stages=["archived"]
        )
        
        if archived:
            self.promote_to_production(model_name, archived[0].version)
```

### A/B Testing Integration

```python
class ABTestManager:
    def __init__(self):
        self.experiments = {}
    
    def create_experiment(
        self,
        name: str,
        model_a: str,
        model_b: str,
        traffic_split: float = 0.5
    ):
        """Create A/B test experiment"""
        self.experiments[name] = {
            "model_a": model_a,
            "model_b": model_b,
            "traffic_split": traffic_split,
            "results_a": [],
            "results_b": []
        }
    
    def route_request(self, experiment_name: str, user_id: str) -> str:
        """Route request to model A or B"""
        import hashlib
        
        exp = self.experiments[experiment_name]
        
        # Deterministic routing based on user_id
        hash_val = int(hashlib.md5(user_id.encode()).hexdigest(), 16)
        if (hash_val % 100) / 100 < exp["traffic_split"]:
            return "model_a"
        return "model_b"
    
    def record_result(self, experiment_name: str, model: str, result: dict):
        """Record experiment result"""
        if model == "model_a":
            self.experiments[experiment_name]["results_a"].append(result)
        else:
            self.experiments[experiment_name]["results_b"].append(result)
    
    def analyze_results(self, experiment_name: str) -> dict:
        """Analyze A/B test results"""
        exp = self.experiments[experiment_name]
        
        return {
            "model_a": self._compute_metrics(exp["results_a"]),
            "model_b": self._compute_metrics(exp["results_b"]),
            "significant": self._statistical_test(
                exp["results_a"],
                exp["results_b"]
            )
        }
```

## AI Agent Platform Reference

The FinAgent platform implements MLOps practices:

```go
// services/ai-agent-platform/internal/agent/
// Model versioning and registry integration

type ModelVersion struct {
    ID          string    `json:"id"`
    Name        string    `json:"name"`
    Version     string    `json:"version"`
    Stage       string    `json:"stage"` // staging, production, archived
    CreatedAt   time.Time `json:"created_at"`
    Metrics     Metrics   `json:"metrics"`
    ArtifactURI string    `json:"artifact_uri"`
}

type ExperimentRun struct {
    RunID       string            `json:"run_id"`
    Experiment  string            `json:"experiment"`
    Parameters  map[string]float64 `json:"parameters"`
    Metrics     map[string]float64 `json:"metrics"`
    Status      string            `json:"status"`
    StartTime   time.Time         `json:"start_time"`
    EndTime     time.Time         `json:"end_time"`
}

// MLOps pipeline integrates:
// - MLflow for experiment tracking
// - Weights & Biases for visualization
// - GitHub Actions for CI/CD
// - Kubernetes for deployment
```

## Summary

- Model registries provide centralized model management
- Experiment tracking enables reproducible research
- CI/CD pipelines automate training, testing, and deployment
- Lifecycle management ensures proper model governance
- A/B testing validates model improvements before full rollout
- Integration with MLOps tools (MLflow, W&B) is essential for production systems
