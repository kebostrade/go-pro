# IO-20: Production Best Practices

**Duration**: 3 hours
**Module**: 5 - Advanced LLM-Ops & Production

## Learning Objectives

- Implement various deployment strategies (blue-green, canary, rolling)
- Design rollback mechanisms for quick recovery
- Build incident response procedures for LLM systems
- Manage SLAs and handle violations
- Implement health checks and monitoring
- Design for reliability and fault tolerance

## Deployment Strategies

### Blue-Green Deployment

```python
class BlueGreenDeployment:
    """
    Blue-Green deployment with instant rollback
    """
    
    def __init__(self):
        self.blue = {"version": "v1", "healthy": True}
        self.green = {"version": "v2", "healthy": True}
        self.active = "blue"
    
    def deploy(self, new_version: str) -> dict:
        """Deploy new version to inactive environment"""
        target = "green" if self.active == "blue" else "blue"
        
        # Deploy to inactive environment
        self._deploy_to(target, new_version)
        
        # Run smoke tests
        if self._smoke_test(target):
            # Switch traffic
            self.active = target
            return {"status": "success", "active": self.active}
        else:
            # Keep old version
            return {"status": "failed", "reason": "smoke_test_failed"}
    
    def rollback(self) -> dict:
        """Instant rollback to previous version"""
        target = "green" if self.active == "blue" else "blue"
        self.active = target
        return {"status": "rolled_back", "active": self.active}
    
    def _deploy_to(self, target: str, version: str):
        """Deploy version to target"""
        pass
    
    def _smoke_test(self, target: str) -> bool:
        """Run smoke tests"""
        return True
```

### Canary Deployment

```python
class CanaryDeployment:
    """
    Gradual rollout with traffic shifting
    """
    
    def __init__(self):
        self.versions = {}  # version -> traffic_percentage
        self.metrics = {}
    
    def add_version(self, version: str, initial_traffic: float = 0):
        """Add new version"""
        self.versions[version] = initial_traffic
    
    def increase_traffic(self, version: str, increment: float) -> dict:
        """Increase traffic to version"""
        self.versions[version] = min(100, self.versions.get(version, 0) + increment)
        
        # Decrease traffic from others
        remaining = 100 - self.versions[version]
        other_versions = [v for v in self.versions if v != version]
        
        for v in other_versions:
            self.versions[v] = remaining / len(other_versions)
        
        return {"status": "updated", "versions": self.versions}
    
    def rollback_version(self, version: str) -> dict:
        """Rollback specific version"""
        if version in self.versions:
            traffic = self.versions[version]
            del self.versions[version]
            
            # Redistribute traffic
            remaining = 100 - traffic
            for v in self.versions:
                self.versions[v] = remaining / len(self.versions)
        
        return {"status": "rolled_back", "versions": self.versions}
    
    def get_metrics(self, version: str) -> dict:
        """Get version-specific metrics"""
        return self.metrics.get(version, {})
    
    def record_request(self, version: str, latency: float, success: bool):
        """Record request for metrics"""
        if version not in self.metrics:
            self.metrics[version] = {"requests": [], "errors": 0}
        
        self.metrics[version]["requests"].append(latency)
        if not success:
            self.metrics[version]["errors"] += 1
```

### Rolling Deployment

```python
class RollingDeployment:
    """
    Gradual replacement of instances
    """
    
    def __init__(self, total_instances: int):
        self.total_instances = total_instances
        self.current_version = None
        self.instances = [{"id": i, "version": None} for i in range(total_instances)]
    
    def deploy(self, new_version: str, batch_size: int = 1) -> dict:
        """Deploy with rolling updates"""
        self.current_version = new_version
        deployed = 0
        
        # Deploy in batches
        for i in range(0, self.total_instances, batch_size):
            batch = self.instances[i:i + batch_size]
            
            # Update batch
            for instance in batch:
                instance["version"] = new_version
            
            deployed += len(batch)
            
            # Wait for health
            if not self._wait_for_health(batch):
                return {"status": "failed", "deployed": deployed}
        
        return {"status": "complete", "version": new_version}
    
    def _wait_for_health(self, batch: list, timeout: int = 60) -> bool:
        """Wait for instances to become healthy"""
        return True
    
    def rollback(self) -> dict:
        """Rollback to previous version"""
        # Implementation
        return {"status": "rolled_back"}
```

## Rollback Mechanisms

### Automated Rollback

```python
class RollbackManager:
    def __init__(self):
        self.checkpoints = []
        self.current_version = None
    
    def create_checkpoint(self, version: str, state: dict):
        """Create rollback checkpoint"""
        checkpoint = {
            "version": version,
            "state": state,
            "timestamp": datetime.now()
        }
        self.checkpoints.append(checkpoint)
    
    def rollback_to(self, checkpoint_index: int) -> dict:
        """Rollback to specific checkpoint"""
        if 0 <= checkpoint_index < len(self.checkpoints):
            checkpoint = self.checkpoints[checkpoint_index]
            self.current_version = checkpoint["version"]
            return {"status": "success", "version": checkpoint["version"]}
        return {"status": "failed", "reason": "invalid_checkpoint"}
    
    def auto_rollback(self, condition: str) -> dict:
        """Automatic rollback based on conditions"""
        conditions = {
            "high_error_rate": self._check_error_rate,
            "high_latency": self._check_latency,
            "health_check_failed": self._check_health
        }
        
        if condition in conditions and conditions[condition]():
            # Rollback to last stable version
            if len(self.checkpoints) > 1:
                return self.rollback_to(len(self.checkpoints) - 2)
        
        return {"status": "no_rollback"}
    
    def _check_error_rate(self) -> bool:
        """Check if error rate is too high"""
        return False
    
    def _check_latency(self) -> bool:
        """Check if latency is too high"""
        return False
    
    def _check_health(self) -> bool:
        """Check if health checks failed"""
        return False
```

### Database Rollback

```python
class ModelVersionControl:
    def __init__(self):
        self.versions = {}
        self.current = None
    
    def save_version(self, version: str, model_path: str, config: dict):
        """Save model version"""
        self.versions[version] = {
            "model_path": model_path,
            "config": config,
            "created_at": datetime.now()
        }
    
    def switch_version(self, version: str) -> bool:
        """Switch to version"""
        if version in self.versions:
            self.current = version
            return True
        return False
    
    def get_current(self) -> dict:
        """Get current version"""
        if self.current:
            return self.versions[self.current]
        return None
```

## Incident Response

### Incident Classification

```python
class IncidentClassifier:
    SEVERITY_LEVELS = {
        "P1": {"name": "Critical", "response_time": "15 min", "description": "Service down"},
        "P2": {"name": "High", "response_time": "1 hour", "description": "Major feature impacted"},
        "P3": {"name": "Medium", "response_time": "4 hours", "description": "Partial impact"},
        "P4": {"name": "Low", "response_time": "24 hours", "description": "Minor issue"}
    }
    
    def classify(self, symptoms: list) -> dict:
        """Classify incident based on symptoms"""
        # Check for critical symptoms
        critical = ["service_down", "data_loss", "security_breach"]
        
        for symptom in symptoms:
            if symptom in critical:
                return self.SEVERITY_LEVELS["P1"]
        
        return self.SEVERITY_LEVELS["P3"]
    
    def get_response_plan(self, severity: str) -> dict:
        """Get response plan for severity"""
        return {
            "severity": severity,
            **self.SEVERITY_LEVELS.get(severity, {})
        }
```

### Incident Response Workflow

```python
class IncidentResponse:
    def __init__(self):
        self.incidents = []
        self.escalation_contacts = {
            "P1": ["oncall-engineer", "engineering-lead", "vp-engineering"],
            "P2": ["oncall-engineer", "engineering-lead"],
            "P3": ["oncall-engineer"],
            "P4": ["team-channel"]
        }
    
    def create_incident(
        self,
        title: str,
        description: str,
        severity: str
    ) -> dict:
        """Create new incident"""
        incident = {
            "id": len(self.incidents) + 1,
            "title": title,
            "description": description,
            "severity": severity,
            "status": "open",
            "created_at": datetime.now(),
            "timeline": []
        }
        
        self.incidents.append(incident)
        
        # Notify escalations
        self._notify(severity, incident)
        
        return incident
    
    def update_status(self, incident_id: int, status: str):
        """Update incident status"""
        for inc in self.incidents:
            if inc["id"] == incident_id:
                inc["status"] = status
                inc["timeline"].append({
                    "event": f"Status changed to {status}",
                    "timestamp": datetime.now()
                })
    
    def add_update(self, incident_id: int, update: str):
        """Add update to incident"""
        for inc in self.incidents:
            if inc["id"] == incident_id:
                inc["timeline"].append({
                    "event": update,
                    "timestamp": datetime.now()
                })
    
    def resolve(self, incident_id: int, resolution: str):
        """Resolve incident"""
        self.update_status(incident_id, "resolved")
        for inc in self.incidents:
            if inc["id"] == incident_id:
                inc["resolution"] = resolution
    
    def _notify(self, severity: str, incident: dict):
        """Notify escalation contacts"""
        contacts = self.escalation_contacts.get(severity, [])
        # Implementation: send notifications
        print(f"Notifying {contacts} about incident: {incident['title']}")
```

### Post-Incident Review

```python
class PostIncidentReview:
    def __init__(self):
        self.reviews = []
    
    def create_review(self, incident: dict) -> dict:
        """Create post-incident review"""
        review = {
            "incident_id": incident["id"],
            "timeline": incident["timeline"],
            "root_cause": "",
            "impact": "",
            "action_items": [],
            "questions": []
        }
        
        self.reviews.append(review)
        return review
    
    def add_root_cause(self, review_id: int, cause: str):
        """Add root cause analysis"""
        for review in self.reviews:
            if review["incident_id"] == review_id:
                review["root_cause"] = cause
    
    def add_action_item(self, review_id: int, item: dict):
        """Add action item"""
        for review in self.reviews:
            if review["incident_id"] == review_id:
                review["action_items"].append(item)
    
    def generate_report(self, review_id: int) -> str:
        """Generate PIR report"""
        for review in self.reviews:
            if review["incident_id"] == review_id:
                return f"""
# Post-Incident Review

## Incident
ID: {review['incident_id']}

## Root Cause
{review['root_cause'] or 'TBD'}

## Action Items
{chr(10).join(f"- {item['description']}" for item in review['action_items'])}

## Timeline
{chr(10).join(f"- {t['event']}" for t in review['timeline'])}
"""
```

## SLA Management

### SLA Definition

```python
class SLA:
    def __init__(self):
        self.metrics = {}
    
    def define_sla(
        self,
        name: str,
        target: float,
        window: str,
        penalty: dict
    ):
        """Define SLA"""
        self.metrics[name] = {
            "target": target,
            "window": window,
            "penalty": penalty
        }
    
    def check_compliance(self, actual: dict) -> dict:
        """Check SLA compliance"""
        results = {}
        
        for metric, target in self.metrics.items():
            if metric in actual:
                actual_value = actual[metric]
                target_value = target["target"]
                
                results[metric] = {
                    "actual": actual_value,
                    "target": target_value,
                    "compliant": actual_value >= target_value if "latency" not in metric else actual_value <= target_value,
                    "penalty": self._calculate_penalty(target, actual_value) if actual_value != target_value else 0
                }
        
        return results
    
    def _calculate_penalty(self, target: dict, actual: float) -> float:
        """Calculate penalty for SLA violation"""
        # Simplified penalty calculation
        return target.get("penalty", {}).get("percentage", 0)
```

### Common LLM SLAs

| Metric | Target | Window | Penalty |
|--------|--------|--------|---------|
| Availability | 99.9% | Monthly | Credit |
| Latency p50 | < 500ms | Hourly | Credit |
| Latency p99 | < 2000ms | Hourly | Credit |
| Error rate | < 1% | Hourly | Credit |
| Throughput | > 100 rps | Hourly | Credit |

## Health Checks

### Health Check Implementation

```python
class LLMHealthCheck:
    def __init__(self, model):
        self.model = model
    
    def check(self) -> dict:
        """Run all health checks"""
        checks = {
            "model_loaded": self._check_model_loaded(),
            "gpu_available": self._check_gpu(),
            "memory_sufficient": self._check_memory(),
            "inference_working": self._check_inference(),
            "response_quality": self._check_quality()
        }
        
        overall = all(checks.values())
        
        return {
            "healthy": overall,
            "checks": checks,
            "timestamp": datetime.now()
        }
    
    def _check_model_loaded(self) -> bool:
        """Check if model is loaded"""
        return self.model is not None
    
    def _check_gpu(self) -> bool:
        """Check GPU availability"""
        import torch
        return torch.cuda.is_available()
    
    def _check_memory(self) -> bool:
        """Check memory availability"""
        import torch
        if torch.cuda.is_available():
            memory_allocated = torch.cuda.memory_allocated() / 1e9
            memory_total = torch.cuda.get_device_properties(0).total_memory / 1e9
            return memory_allocated / memory_total < 0.9
        return True
    
    def _check_inference(self) -> bool:
        """Check if inference works"""
        try:
            result = self.model.generate("test")
            return len(result) > 0
        except:
            return False
    
    def _check_quality(self) -> bool:
        """Check response quality"""
        # Simple sanity check
        try:
            result = self.model.generate("What is 2+2?")
            return "4" in result or "four" in result.lower()
        except:
            return False
```

## Reliability Patterns

### Circuit Breaker

```python
class CircuitBreaker:
    def __init__(self, failure_threshold: int = 5, timeout: int = 60):
        self.failure_threshold = failure_threshold
        self.timeout = timeout
        self.failures = 0
        self.last_failure_time = None
        self.state = "closed"  # closed, open, half-open
    
    def call(self, func, *args, **kwargs):
        """Call function through circuit breaker"""
        if self.state == "open":
            if time.time() - self.last_failure_time > self.timeout:
                self.state = "half-open"
            else:
                raise Exception("Circuit breaker open")
        
        try:
            result = func(*args, **kwargs)
            
            if self.state == "half-open":
                self.state = "closed"
                self.failures = 0
            
            return result
            
        except Exception as e:
            self.failures += 1
            self.last_failure_time = time.time()
            
            if self.failures >= self.failure_threshold:
                self.state = "open"
            
            raise
```

### Rate Limiting

```python
class TokenBucketRateLimiter:
    def __init__(self, rate: int, capacity: int):
        self.rate = rate
        self.capacity = capacity
        self.tokens = capacity
        self.last_refill = time.time()
    
    def allow_request(self) -> bool:
        """Check if request is allowed"""
        self._refill()
        
        if self.tokens >= 1:
            self.tokens -= 1
            return True
        return False
    
    def _refill(self):
        """Refill tokens"""
        now = time.time()
        elapsed = now - self.last_refill
        new_tokens = elapsed * self.rate
        
        self.tokens = min(self.capacity, self.tokens + new_tokens)
        self.last_refill = now
    
    def get_wait_time(self) -> float:
        """Get wait time for next token"""
        if self.tokens >= 1:
            return 0
        return (1 - self.tokens) / self.rate
```

## AI Agent Platform Reference

The FinAgent platform implements production practices:

```go
// services/ai-agent-platform/internal/agent/
// Production deployment configuration

type DeploymentConfig struct {
    Strategy      string `json:"strategy"` // canary, blue-green, rolling
    MaxSurge      int    `json:"max_surge"`
    MaxUnavailable int   `json:"max_unavailable"`
    HealthCheck   HealthCheckConfig `json:"health_check"`
    RateLimit     RateLimitConfig   `json:"rate_limit"`
    CircuitBreaker CircuitBreakerConfig `json:"circuit_breaker"`
}

type HealthCheckConfig struct {
    Enabled        bool   `json:"enabled"`
    IntervalSeconds int   `json:"interval_seconds"`
    TimeoutSeconds  int   `json:"timeout_seconds"`
    FailureThreshold int  `json:"failure_threshold"`
}

type IncidentResponseConfig struct {
    OnCallRotation string            `json:"on_call_rotation"`
    Escalations   []EscalationLevel `json:"escalations"`
    Runbooks      map[string]string `json:"runbooks"`
}

// Kubernetes deployment uses:
// - HorizontalPodAutoscaler for scaling
// - PodDisruptionBudget for availability
// - Readiness/Liveness probes for health
// - Service mesh for traffic management
```

## Summary

- Multiple deployment strategies enable safe rollouts
- Automated rollback prevents extended outages
- Incident response procedures minimize MTTR
- SLA management sets clear expectations
- Health checks ensure system reliability
- Circuit breakers and rate limiting prevent cascade failures
- Comprehensive monitoring is essential for production systems
