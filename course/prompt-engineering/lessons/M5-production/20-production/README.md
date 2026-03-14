# PE-20: Production Systems & Best Practices

**Duration**: 3 hours
**Module**: 5 - Production & Optimization

## Learning Objectives

- Design production-ready prompt architectures
- Implement versioning and deployment strategies
- Build monitoring and observability systems
- Apply best practices for scale

## Production Architecture

### System Components

```
┌─────────────────────────────────────────────────────────────┐
│                 PRODUCTION ARCHITECTURE                     │
│                                                             │
│  ┌─────────┐    ┌─────────────┐    ┌─────────────────┐     │
│  │ Request │───►│ Prompt      │───►│ LLM Provider    │     │
│  │ Router  │    │ Manager     │    │ (Anthropic/etc) │     │
│  └─────────┘    └─────────────┘    └─────────────────┘     │
│       │              │                     │               │
│       │              ▼                     │               │
│       │        ┌───────────┐               │               │
│       │        │ Version   │               │               │
│       │        │ Control   │               │               │
│       │        └───────────┘               │               │
│       │              │                     │               │
│       ▼              ▼                     ▼               │
│  ┌─────────────────────────────────────────────────────┐   │
│  │              MONITORING & LOGGING                    │   │
│  └─────────────────────────────────────────────────────┘   │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

### Prompt Manager Pattern

```python
from dataclasses import dataclass
from typing import Dict, Optional
import yaml

@dataclass
class PromptConfig:
    name: str
    version: str
    template: str
    variables: list
    metadata: dict

class PromptManager:
    """Centralized prompt management."""

    def __init__(self, config_path: str):
        self.prompts: Dict[str, PromptConfig] = {}
        self._load_configs(config_path)

    def _load_configs(self, path: str):
        """Load prompt configurations from files."""
        with open(path) as f:
            configs = yaml.safe_load(f)
            for name, config in configs.items():
                self.prompts[name] = PromptConfig(**config)

    def get_prompt(self, name: str, version: str = "latest") -> PromptConfig:
        """Retrieve a specific prompt version."""
        key = f"{name}:{version}" if version != "latest" else name
        return self.prompts.get(key)

    def render(self, name: str, variables: dict, version: str = "latest") -> str:
        """Render a prompt with variables."""
        config = self.get_prompt(name, version)
        if not config:
            raise ValueError(f"Prompt {name} not found")

        # Validate required variables
        missing = set(config.variables) - set(variables.keys())
        if missing:
            raise ValueError(f"Missing variables: {missing}")

        # Render template
        return config.template.format(**variables)
```

## Version Control

### Prompt Registry

```yaml
# prompts/registry.yaml
sentiment_analysis:
  v1:
    template: |
      Classify sentiment: {text}
      Output: positive/negative/neutral
    variables: [text]
    metadata:
      created: 2024-01-15
      author: team-ml
      deprecated: false

  v2:
    template: |
      Analyze the sentiment of this text.
      Consider: tone, word choice, context.

      Text: {text}

      Output as JSON:
      {{"sentiment": "positive|negative|neutral", "confidence": 0.0-1.0}}
    variables: [text]
    metadata:
      created: 2024-02-20
      author: team-ml
      deprecated: false
      improvement: "Added confidence scoring and JSON output"

  v3:
    template: |
      [Current best version]
    variables: [text]
    metadata:
      created: 2024-03-10
      deprecated: false
```

### Version Migration

```python
class PromptMigrator:
    """Handle prompt version migrations."""

    def __init__(self, prompt_manager: PromptManager):
        self.manager = prompt_manager

    def migrate(self, prompt_name: str, from_version: str, to_version: str,
                input_data: dict) -> dict:
        """Migrate input data between prompt versions."""

        # Get version configs
        from_config = self.manager.get_prompt(prompt_name, from_version)
        to_config = self.manager.get_prompt(prompt_name, to_version)

        # Transform data if needed
        migrated_data = input_data.copy()

        # Handle variable renames
        if from_config.metadata.get("variable_mapping"):
            mapping = from_config.metadata["variable_mapping"]
            for old_key, new_key in mapping.items():
                if old_key in migrated_data:
                    migrated_data[new_key] = migrated_data.pop(old_key)

        return migrated_data
```

## Deployment Strategies

### Strategy 1: Blue-Green Deployment

```python
class BlueGreenDeployment:
    """Blue-green deployment for prompts."""

    def __init__(self, prompt_name: str):
        self.prompt_name = prompt_name
        self.blue_version = "v2"  # Current production
        self.green_version = "v3"  # New version

    def get_prompt(self, use_green: bool = False) -> PromptConfig:
        """Get prompt based on traffic allocation."""
        version = self.green_version if use_green else self.blue_version
        return self.manager.get_prompt(self.prompt_name, version)

    def switch_traffic(self, percentage: int):
        """Gradually shift traffic to green."""
        # Use feature flag or config to control traffic split
        pass

    def promote_green(self):
        """Promote green to blue (full production)."""
        self.blue_version = self.green_version
        self.log_deployment("promoted", self.blue_version)
```

### Strategy 2: Canary Release

```python
class CanaryRelease:
    """Canary release for prompt testing."""

    def __init__(self, prompt_name: str, stable_version: str, canary_version: str):
        self.prompt_name = prompt_name
        self.stable = stable_version
        self.canary = canary_version
        self.canary_percentage = 0

    def get_version_for_request(self, request_id: str) -> str:
        """Determine which version to use for this request."""
        import hashlib

        # Deterministic assignment based on request ID
        hash_value = int(hashlib.md5(request_id.encode()).hexdigest(), 16)
        bucket = (hash_value % 100) + 1  # 1-100

        if bucket <= self.canary_percentage:
            return self.canary
        return self.stable

    def increase_canary(self, percentage: int):
        """Increase canary traffic."""
        self.canary_percentage = min(100, percentage)
```

## Monitoring & Observability

### Metrics Collection

```python
from dataclasses import dataclass
from datetime import datetime
import statistics

@dataclass
class PromptMetrics:
    prompt_name: str
    prompt_version: str
    timestamp: datetime
    input_tokens: int
    output_tokens: int
    latency_ms: float
    success: bool
    error_message: str = None

class MetricsCollector:
    """Collect and aggregate prompt metrics."""

    def __init__(self):
        self.metrics: list[PromptMetrics] = []

    def record(self, metrics: PromptMetrics):
        """Record a metrics event."""
        self.metrics.append(metrics)

    def get_aggregate_stats(self, prompt_name: str, window_minutes: int = 60) -> dict:
        """Get aggregated statistics for a prompt."""
        cutoff = datetime.now() - timedelta(minutes=window_minutes)
        recent = [m for m in self.metrics
                  if m.prompt_name == prompt_name and m.timestamp > cutoff]

        if not recent:
            return {}

        return {
            "total_requests": len(recent),
            "success_rate": sum(1 for m in recent if m.success) / len(recent),
            "avg_latency_ms": statistics.mean(m.latency_ms for m in recent),
            "p95_latency_ms": statistics.quantiles([m.latency_ms for m in recent], n=20)[18],
            "avg_input_tokens": statistics.mean(m.input_tokens for m in recent),
            "avg_output_tokens": statistics.mean(m.output_tokens for m in recent),
            "total_cost": sum(calculate_cost(m) for m in recent)
        }
```

### Alerting

```python
class PromptMonitor:
    """Monitor prompts and trigger alerts."""

    def __init__(self, metrics_collector: MetricsCollector):
        self.collector = metrics_collector
        self.alerts = []

    def check_health(self, prompt_name: str) -> dict:
        """Check health of a prompt."""
        stats = self.collector.get_aggregate_stats(prompt_name)

        issues = []

        # Check success rate
        if stats.get("success_rate", 1) < 0.95:
            issues.append({
                "type": "low_success_rate",
                "value": stats["success_rate"],
                "threshold": 0.95
            })

        # Check latency
        if stats.get("p95_latency_ms", 0) > 5000:
            issues.append({
                "type": "high_latency",
                "value": stats["p95_latency_ms"],
                "threshold": 5000
            })

        # Check cost
        if stats.get("total_cost", 0) > 100:
            issues.append({
                "type": "high_cost",
                "value": stats["total_cost"],
                "threshold": 100
            })

        return {
            "healthy": len(issues) == 0,
            "issues": issues,
            "stats": stats
        }

    def alert(self, issue: dict):
        """Send alert for an issue."""
        # Send to Slack, PagerDuty, etc.
        self.alerts.append(issue)
        print(f"ALERT: {issue}")
```

## Best Practices Checklist

### Development

```
□ Version control all prompts
□ Document prompt purpose and variables
□ Include examples in prompt definitions
□ Write tests for prompt behavior
□ Review prompts before deployment
□ Use consistent naming conventions
```

### Testing

```
□ Unit test each prompt with expected outputs
□ Integration test prompt in full workflow
□ Load test for latency under traffic
□ A/B test new versions before rollout
□ Regression test after changes
□ Edge case testing (empty, long, special chars)
```

### Deployment

```
□ Use staged rollouts (canary/blue-green)
□ Monitor metrics after deployment
□ Have rollback plan ready
□ Document deployment history
□ Test in staging before production
□ Coordinate with dependent systems
```

### Operations

```
□ Monitor success rate and latency
□ Alert on anomalies
□ Track cost per prompt
□ Log requests for debugging
□ Regular performance reviews
□ Security audit periodically
```

## Production Template

```python
# production_prompt_system.py

class ProductionPromptSystem:
    """Complete production prompt system."""

    def __init__(self, config):
        self.manager = PromptManager(config.prompts_path)
        self.collector = MetricsCollector()
        self.monitor = PromptMonitor(self.collector)
        self.deployment = CanaryRelease(
            config.prompt_name,
            config.stable_version,
            config.canary_version
        )

    def process_request(self, request_id: str, prompt_name: str,
                        variables: dict) -> dict:
        """Process a request through the prompt system."""

        start_time = time.time()

        try:
            # Get appropriate version
            version = self.deployment.get_version_for_request(request_id)

            # Render prompt
            prompt = self.manager.render(prompt_name, variables, version)

            # Call LLM
            response = self.call_llm(prompt)

            # Record metrics
            metrics = PromptMetrics(
                prompt_name=prompt_name,
                prompt_version=version,
                timestamp=datetime.now(),
                input_tokens=response.usage.input_tokens,
                output_tokens=response.usage.output_tokens,
                latency_ms=(time.time() - start_time) * 1000,
                success=True
            )
            self.collector.record(metrics)

            return {
                "success": True,
                "output": response.content,
                "version": version
            }

        except Exception as e:
            # Record failure
            metrics = PromptMetrics(
                prompt_name=prompt_name,
                prompt_version=version,
                timestamp=datetime.now(),
                input_tokens=0,
                output_tokens=0,
                latency_ms=(time.time() - start_time) * 1000,
                success=False,
                error_message=str(e)
            )
            self.collector.record(metrics)

            return {
                "success": False,
                "error": str(e)
            }
```

## Exercise

### Exercise 20.1: Design Prompt Registry

Design a prompt registry structure for a chatbot that has:
- 5 different prompt types
- Version history
- Metadata for each version

### Exercise 20.2: Implement Monitoring

Write code to:
- Collect metrics for each prompt call
- Calculate success rate, latency, cost
- Trigger alerts when thresholds exceeded

### Exercise 20.3: Deployment Pipeline

Design a deployment pipeline for prompts:
- Development → Staging → Production
- Rollback capability
- A/B testing support

## Key Takeaways

- ✅ Use centralized prompt management
- ✅ Version control all prompts with metadata
- ✅ Deploy with staged rollouts (canary/blue-green)
- ✅ Monitor success rate, latency, and cost
- ✅ Have rollback plans ready
- ✅ Test thoroughly before production

## Course Complete! 🎉

Congratulations! You've completed the Prompt Engineering Full Course.

### What You've Learned

- **Module 1**: Foundations (tokens, context, prompt anatomy)
- **Module 2**: Core Techniques (zero-shot, few-shot, CoT)
- **Module 3**: Advanced Patterns (ReAct, ToT, chaining)
- **Module 4**: Task-Specific Prompting (code, data, creative)
- **Module 5**: Production Systems (evaluation, security, deployment)

### Next Steps

1. Complete the course projects
2. Practice with real-world scenarios
3. Build your own prompt library
4. Stay updated on new techniques
