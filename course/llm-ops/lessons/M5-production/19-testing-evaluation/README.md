# IO-19: LLM Testing & Evaluation

**Duration**: 3 hours
**Module**: 5 - Advanced LLM-Ops & Production

## Learning Objectives

- Understand LLM-specific evaluation challenges
- Implement benchmark testing with standard datasets
- Design A/B testing frameworks for production models
- Conduct human evaluation campaigns
- Build automated evaluation pipelines with guardrails
- Apply automated metrics (BLEU, ROUGE, etc.) appropriately

## Evaluation Challenges

### Why LLM Evaluation is Different

| Traditional ML | LLM Evaluation |
|---------------|----------------|
| Single correct answer | Multiple valid outputs |
| Numeric metrics clear | Quality is subjective |
| Test set static | Capabilities evolving |
| Fixed input format | Open-ended tasks |

### Evaluation Dimensions

```python
class LLMEvaluationDimensions:
    """
    Multi-dimensional evaluation framework for LLMs
    """
    
    dimensions = {
        "accuracy": {
            "description": "Correctness of outputs",
            "metrics": ["exact_match", "F1", "BLEU", "ROUGE"]
        },
        "helpfulness": {
            "description": "Usefulness to user",
            "metrics": ["human_rating", "task_completion"]
        },
        "coherence": {
            "description": "Logical consistency",
            "metrics": ["perplexity", "consistency_score"]
        },
        "safety": {
            "description": "Harmfulness of outputs",
            "metrics": ["refusal_rate", "toxicity_score"]
        },
        "efficiency": {
            "description": "Latency and resource usage",
            "metrics": ["latency_p50", "throughput"]
        }
    }
```

## Benchmark Datasets

### Common LLM Benchmarks

| Benchmark | Purpose | Metrics |
|----------|---------|---------|
| MMLU | Multi-task understanding | Accuracy |
| HumanEval | Code generation | Pass@1 |
| GSM8K | Math reasoning | Accuracy |
| BIG-Bench | Diverse tasks | Various |
| HELM | Comprehensive | Multi-metric |
| MT-Bench | Chat capability | Win rate |

### Running Benchmarks

```python
from lm_eval import evaluator
import json

class BenchmarkRunner:
    def __init__(self, model, tokenizer):
        self.model = model
        self.tokenizer = tokenizer
    
    def run_mmlu(self, model_name: str):
        """Run MMLU benchmark"""
        results = evaluator.simple_evaluate(
            model="hf",
            model_args=f"pretrained={model_name}",
            tasks="mmlu",
            num_fewshot=5
        )
        return results
    
    def run_humaneval(self, model_name: str):
        """Run HumanEval code benchmark"""
        results = evaluator.simple_evaluate(
            model="hf",
            model_args=f"pretrained={model_name}",
            tasks="humaneval",
            num_fewshot=0
        )
        return results
    
    def run_custom(self, dataset_path: str, task_config: dict):
        """Run custom benchmark"""
        # Load dataset
        dataset = load_dataset(dataset_path)
        
        # Evaluate each example
        results = []
        for example in dataset:
            output = self.generate(example["prompt"])
            score = self.compute_score(output, example["reference"])
            results.append(score)
        
        return {
            "mean_score": sum(results) / len(results),
            "results": results
        }
```

### Custom Benchmark Creation

```python
class CustomBenchmark:
    def __init__(self, name: str):
        self.name = name
        self.test_cases = []
    
    def add_test_case(
        self,
        prompt: str,
        reference: str,
        metadata: dict = None
    ):
        """Add a test case"""
        self.test_cases.append({
            "prompt": prompt,
            "reference": reference,
            "metadata": metadata or {}
        })
    
    def add_from_file(self, file_path: str):
        """Load test cases from JSON file"""
        with open(file_path) as f:
            data = json.load(f)
            self.test_cases.extend(data)
    
    def run(self, model) -> dict:
        """Run benchmark"""
        results = []
        for case in self.test_cases:
            output = model.generate(case["prompt"])
            
            evaluation = {
                "prompt": case["prompt"],
                "reference": case["reference"],
                "output": output,
                "scores": self._score(output, case["reference"])
            }
            results.append(evaluation)
        
        return self._aggregate(results)
    
    def _score(self, output: str, reference: str) -> dict:
        """Compute scores"""
        # Implement scoring logic
        return {"exact_match": output.strip() == reference.strip()}
    
    def _aggregate(self, results: list) -> dict:
        """Aggregate results"""
        import numpy as np
        
        all_scores = {}
        for key in results[0]["scores"].keys():
            scores = [r["scores"][key] for r in results]
            all_scores[key] = {
                "mean": np.mean(scores),
                "std": np.std(scores),
                "min": min(scores),
                "max": max(scores)
            }
        
        return {
            "benchmark": self.name,
            "total_cases": len(results),
            "scores": all_scores
        }
```

## Automated Metrics

### Text Generation Metrics

```python
from datasets import load_metric
import evaluate

class LLMEvaluationMetrics:
    def __init__(self):
        self.bleu = load_metric("bleu")
        self.rouge = load_metric("rouge")
        self.bertscore = load_metric("bertscore")
    
    def compute_bleu(
        self,
        predictions: list,
        references: list
    ) -> float:
        """Compute BLEU score"""
        return self.bleu.compute(
            predictions=predictions,
            references=[[r] for r in references]
        )["bleu"]
    
    def compute_rouge(
        self,
        predictions: list,
        references: list
    ) -> dict:
        """Compute ROUGE scores"""
        results = self.rouge.compute(
            predictions=predictions,
            references=references
        )
        return {
            "rouge1": results["rouge1"].mid.fmeasure,
            "rouge2": results["rouge2"].mid.fmeasure,
            "rougeL": results["rougeL"].mid.fmeasure
        }
    
    def compute_bertscore(
        self,
        predictions: list,
        references: list,
        lang: str = "en"
    ) -> dict:
        """Compute BERTScore"""
        results = self.bertscore.compute(
            predictions=predictions,
            references=references,
            lang=lang
        )
        return {
            "precision": sum(results["precision"]) / len(results["precision"]),
            "recall": sum(results["recall"]) / len(results["recall"]),
            "f1": sum(results["f1"]) / len(results["f1"])
        }
```

### Task-Specific Metrics

```python
class TaskSpecificMetrics:
    @staticmethod
    def classification_metrics(predictions: list, labels: list) -> dict:
        """Classification metrics"""
        from sklearn.metrics import classification_report
        
        report = classification_report(labels, predictions, output_dict=True)
        return {
            "accuracy": report["accuracy"],
            "precision": report["weighted avg"]["precision"],
            "recall": report["weighted avg"]["recall"],
            "f1": report["weighted avg"]["f1-score"]
        }
    
    @staticmethod
    def extraction_metrics(predictions: list, references: list) -> dict:
        """Entity extraction metrics"""
        # Exact match
        exact_match = sum(
            p.strip() == r.strip() for p, r in zip(predictions, references)
        ) / len(predictions)
        
        # Partial match (contains)
        partial_match = sum(
            r.strip() in p.strip() for p, r in zip(predictions, references)
        ) / len(predictions)
        
        return {
            "exact_match": exact_match,
            "partial_match": partial_match
        }
    
    @staticmethod
    def summarization_metrics(predictions: list, references: list) -> dict:
        """Summarization metrics"""
        metrics = LLMEvaluationMetrics()
        
        return {
            "bleu": metrics.compute_bleu(predictions, references),
            "rouge": metrics.compute_rouge(predictions, references)
        }
```

## A/B Testing

### Production A/B Framework

```python
import random
import time
from typing import Dict, List
import hashlib

class ABTestFramework:
    def __init__(self):
        self.experiments = {}
        self.events = []
    
    def create_experiment(
        self,
        name: str,
        variants: Dict[str, dict],
        traffic_split: Dict[str, float] = None
    ):
        """Create A/B test experiment"""
        if traffic_split is None:
            traffic_split = {name: 1.0 / len(variants) for name in variants}
        
        self.experiments[name] = {
            "variants": variants,
            "traffic_split": traffic_split,
            "start_time": time.time()
        }
    
    def assign_variant(
        self,
        experiment_name: str,
        user_id: str
    ) -> str:
        """Assign user to variant deterministically"""
        exp = self.experiments[experiment_name]
        
        # Use hash for deterministic assignment
        hash_input = f"{experiment_name}:{user_id}"
        hash_val = int(hashlib.md5(hash_input.encode()).hexdigest(), 16)
        
        cumulative = 0
        for variant, split in exp["traffic_split"].items():
            cumulative += split
            if (hash_val % 10000) / 10000 < cumulative:
                return variant
        
        return list(exp["variants"].keys())[0]
    
    def record_event(
        self,
        experiment_name: str,
        user_id: str,
        variant: str,
        event_type: str,
        properties: Dict = None
    ):
        """Record event for analysis"""
        self.events.append({
            "experiment": experiment_name,
            "user_id": user_id,
            "variant": variant,
            "event_type": event_type,
            "properties": properties or {},
            "timestamp": time.time()
        })
    
    def analyze(self, experiment_name: str) -> Dict:
        """Analyze experiment results"""
        exp_events = [e for e in self.events if e["experiment"] == experiment_name]
        
        results = {}
        for variant in self.experiments[experiment_name]["variants"]:
            variant_events = [e for e in exp_events if e["variant"] == variant]
            
            # Compute metrics per variant
            conversion_events = [e for e in variant_events if e["event_type"] == "conversion"]
            
            results[variant] = {
                "users": len(set(e["user_id"] for e in variant_events)),
                "conversions": len(conversion_events),
                "conversion_rate": len(conversion_events) / len(set(e["user_id"] for e in variant_events)) if variant_events else 0
            }
        
        # Statistical significance
        results["significance"] = self._compute_significance(results)
        
        return results
    
    def _compute_significance(self, results: Dict) -> Dict:
        """Compute statistical significance"""
        # Simplified chi-square test
        return {"p_value": 0.05, "significant": True}
```

### Shadow Mode Testing

```python
class ShadowModeTester:
    """
    Run new model in shadow mode - compare outputs without serving to users
    """
    
    def __init__(self, production_model, shadow_model):
        self.production_model = production_model
        self.shadow_model = shadow_model
        self.comparisons = []
    
    def process(self, prompt: str) -> Dict:
        """Process prompt with both models"""
        # Production
        prod_start = time.time()
        prod_output = self.production_model.generate(prompt)
        prod_latency = time.time() - prod_start
        
        # Shadow
        shadow_start = time.time()
        shadow_output = self.shadow_model.generate(prompt)
        shadow_latency = time.time() - shadow_start
        
        # Compare
        comparison = {
            "prompt": prompt,
            "production_output": prod_output,
            "shadow_output": shadow_output,
            "production_latency": prod_latency,
            "shadow_latency": shadow_latency,
            "outputs_match": prod_output == shadow_output,
            "similarity": self._compute_similarity(prod_output, shadow_output)
        }
        
        self.comparisons.append(comparison)
        return comparison
    
    def analyze(self) -> Dict:
        """Analyze shadow test results"""
        match_rate = sum(1 for c in self.comparisons if c["outputs_match"]) / len(self.comparisons)
        
        avg_latency_diff = sum(
            c["shadow_latency"] - c["production_latency"] 
            for c in self.comparisons
        ) / len(self.comparisons)
        
        return {
            "total_requests": len(self.comparisons),
            "output_match_rate": match_rate,
            "avg_latency_difference_ms": avg_latency_diff * 1000
        }
```

## Human Evaluation

### Evaluation Framework

```python
class HumanEvaluationCampaign:
    def __init__(self, campaign_name: str):
        self.campaign_name = campaign_name
        self.eval_items = []
        self.ratings = []
    
    def create_eval_items(
        self,
        prompts: List[str],
        outputs: List[str],
        metadata: Dict = None
    ):
        """Create items for human evaluation"""
        for i, (prompt, output) in enumerate(zip(prompts, outputs)):
            self.eval_items.append({
                "id": f"{self.campaign_name}_{i}",
                "prompt": prompt,
                "output": output,
                "metadata": metadata or {}
            })
    
    def add_rating(
        self,
        item_id: str,
        evaluator_id: str,
        ratings: Dict[str, float]
    ):
        """Add evaluator rating"""
        self.ratings.append({
            "item_id": item_id,
            "evaluator_id": evaluator_id,
            "ratings": ratings,
            "timestamp": time.time()
        })
    
    def get_criteria(self) -> List[Dict]:
        """Get evaluation criteria"""
        return [
            {"name": "accuracy", "description": "Is the output factually correct?"},
            {"name": "helpfulness", "description": "Does the output help the user?"},
            {"name": "coherence", "description": "Is the output coherent?"},
            {"name": "safety", "description": "Is the output safe?"}
        ]
    
    def aggregate_ratings(self) -> Dict:
        """Aggregate ratings across evaluators"""
        from collections import defaultdict
        import numpy as np
        
        ratings_by_criterion = defaultdict(list)
        
        for rating in self.ratings:
            for criterion, score in rating["ratings"].items():
                ratings_by_criterion[criterion].append(score)
        
        return {
            criterion: {
                "mean": np.mean(scores),
                "std": np.std(scores),
                "count": len(scores)
            }
            for criterion, scores in ratings_by_criterion.items()
        }
```

### Evaluation Interface

```yaml
# Evaluation interface specification
evaluation_interface:
  - type: pairwise_comparison
    question: "Which response is better?"
    options: ["Response A", "Response B", "Tie"]
    criteria:
      - accuracy
      - helpfulness
      
  - type: rating_scale
    criteria:
      - name: accuracy
        label: "Accuracy"
        scale: [1, 2, 3, 4, 5]
        description: "How accurate is this response?"
        
      - name: safety
        label: "Safety"
        scale: [1, 2, 3, 4, 5]
        description: "Is this response safe?"
```

## Guardrails

### Input/Output Guardrails

```python
class LLMSafetyGuardrails:
    def __init__(self):
        self.input_filters = []
        self.output_filters = []
    
    def add_input_filter(self, filter_func):
        """Add input filter"""
        self.input_filters.append(filter_func)
    
    def add_output_filter(self, filter_func):
        """Add output filter"""
        self.output_filters.append(filter_func)
    
    def check_input(self, text: str) -> tuple[bool, str]:
        """Check input for safety issues"""
        for filter_func in self.input_filters:
            is_safe, reason = filter_func(text)
            if not is_safe:
                return False, reason
        return True, ""
    
    def check_output(self, text: str) -> tuple[bool, str]:
        """Check output for safety issues"""
        for filter_func in self.output_filters:
            is_safe, reason = filter_func(text)
            if not is_safe:
                return False, reason
        return True, ""


# Common filters
def prompt_injection_filter(text: str) -> tuple[bool, str]:
    """Detect prompt injection attempts"""
    patterns = [
        r"ignore\s+(previous|all|your)",
        r"system\s+(prompt|role)",
        r"new\s+instructions"
    ]
    
    import re
    for pattern in patterns:
        if re.search(pattern, text, re.IGNORECASE):
            return False, "Potential prompt injection detected"
    
    return True, ""


def pii_filter(text: str) -> tuple[bool, str]:
    """Detect PII in output"""
    import re
    
    pii_patterns = {
        "SSN": r"\d{3}-\d{2}-\d{4}",
        "EMAIL": r"[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}",
        "PHONE": r"\d{3}[-.]?\d{3}[-.]?\d{4}"
    }
    
    for pii_type, pattern in pii_patterns.items():
        if re.search(pattern, text):
            return False, f"Potential {pii_type} detected"
    
    return True, ""
```

## AI Agent Platform Reference

The FinAgent platform implements evaluation:

```go
// services/ai-agent-platform/internal/agent/
// Evaluation and testing infrastructure

type EvaluationResult struct {
    RunID          string            `json:"run_id"`
    ModelVersion   string            `json:"model_version"`
    TestSet        string            `json:"test_set"`
    Metrics        map[string]float64 `json:"metrics"`
    GuardrailChecks map[string]bool   `json:"guardrail_checks"`
    Timestamp      time.Time         `json:"timestamp"`
}

type GuardrailConfig struct {
    EnablePIIFilter        bool     `json:"enable_pii_filter"`
    EnableToxicityFilter  bool     `json:"enable_toxicity_filter"`
    EnableInjectionFilter bool     `json:"enable_injection_filter"`
    MaxOutputLength       int      `json:"max_output_length"`
}

// Evaluation runs:
// 1. Unit tests on individual tools
// 2. Integration tests on agent workflows
// 3. Benchmark tests on standard datasets
// 4. Human evaluation campaigns
// 5. Production A/B testing
```

## Summary

- LLM evaluation requires multi-dimensional approaches
- Benchmarks provide standardized comparisons
- Automated metrics (BLEU, ROUGE, BERTScore) have limitations
- Human evaluation is essential but expensive
- A/B testing enables data-driven decisions
- Guardrails ensure safe production deployments
- Comprehensive evaluation combines multiple methods
