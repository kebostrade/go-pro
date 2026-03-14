# Exercise: Production Best Practices

## Problem 1: Canary Deployment

Implement canary deployment:

```python
import time

class CanaryDeployer:
    def __init__(self):
        self.versions = {}
        self.traffic分配 = {}
    
    def add_version(self, version: str, initial_traffic: float = 0):
        """Add new version"""
        # Your implementation
        pass
    
    def shift_traffic(self, version: str, target_percentage: float) -> dict:
        """Shift traffic to version"""
        # Your implementation
        pass
    
    def rollback(self, version: str) -> dict:
        """Rollback version"""
        # Your implementation
        pass
    
    def get_traffic_status(self) -> dict:
        """Get current traffic allocation"""
        # Your implementation
        pass

# Test
deployer = CanaryDeployer()

deployer.add_version("v1.0.0", 100)  # Initial version with all traffic
deployer.add_version("v2.0.0", 0)     # New version

# Gradually shift traffic
deployer.shift_traffic("v2.0.0", 10)   # 10% to new version
deployer.shift_traffic("v2.0.0", 25)   # 25% to new version
deployer.shift_traffic("v2.0.0", 50)   # 50% to new version

# Rollback if issues
deployer.rollback("v2.0.0")

print(deployer.get_traffic_status())
```

---

## Problem 2: Health Check System

Implement health checks:

```python
import random

class HealthChecker:
    def __init__(self):
        self.checks = []
    
    def register_check(self, name: str, check_func):
        """Register health check"""
        # Your implementation
        pass
    
    def run_checks(self) -> dict:
        """Run all health checks"""
        # Your implementation
        pass
    
    def get_status(self) -> str:
        """Get overall status"""
        # Your implementation
        pass

# Test with mock checks
def check_model_loaded():
    return random.choice([True, True, True, False])  # 75% healthy

def check_gpu_memory():
    return random.random() < 0.9  # 90% healthy

def check_inference():
    return random.choice([True, True, True, True, False])  # 80% healthy

checker = HealthChecker()
checker.register_check("model_loaded", check_model_loaded)
checker.register_check("gpu_memory", check_gpu_memory)
checker.register_check("inference", check_inference)

# Run multiple checks
for i in range(5):
    result = checker.run_checks()
    print(f"Run {i+1}: {checker.get_status()}")
```

---

## Problem 3: Circuit Breaker

Implement circuit breaker:

```python
class CircuitBreaker:
    def __init__(self, failure_threshold: int = 3, timeout: int = 10):
        self.failure_threshold = failure_threshold
        self.timeout = timeout
        self.state = "closed"  # closed, open, half-open
        self.failures = 0
        self.last_failure_time = None
    
    def call(self, func, *args, **kwargs):
        """Execute function through circuit breaker"""
        # Your implementation
        pass
    
    def reset(self):
        """Reset circuit breaker"""
        # Your implementation
        pass

# Test
def unreliable_function():
    if random.random() < 0.3:
        raise Exception("Random failure")
    return "Success"

cb = CircuitBreaker(failure_threshold=3)

# Try calling multiple times
results = []
for i in range(20):
    try:
        result = cb.call(unreliable_function)
        results.append("success")
    except Exception as e:
        results.append(f"failed: {str(e)[:20]}")

print(results)
```

---

## Problem 4: Incident Response

Implement incident management:

```python
class IncidentManager:
    def __init__(self):
        self.incidents = []
        self.escalation_levels = ["P4", "P3", "P2", "P1"]
    
    def create_incident(
        self,
        title: str,
        description: str,
        severity: str
    ) -> dict:
        """Create incident"""
        # Your implementation
        pass
    
    def add_update(self, incident_id: int, update: str):
        """Add update to incident"""
        # Your implementation
        pass
    
    def escalate(self, incident_id: int) -> dict:
        """Escalate incident"""
        # Your implementation
        pass
    
    def resolve(self, incident_id: int, resolution: str):
        """Resolve incident"""
        # Your implementation
        pass
    
    def get_incident(self, incident_id: int) -> dict:
        """Get incident details"""
        # Your implementation
        pass

# Test
manager = IncidentManager()

# Create incident
incident = manager.create_incident(
    "API returning 500 errors",
    "Users reporting API failures",
    "P3"
)
print(f"Created: {incident}")

# Add updates
manager.add_update(incident["id"], "Investigating load balancer")
manager.add_update(incident["id"], "Found high memory usage")

# Escalate
manager.escalate(incident["id"])

# Resolve
manager.resolve(incident["id"], "Restarted API pods, monitoring for 30 min")

print(manager.get_incident(incident["id"]))
```

---

## Problem 5: Rollback System

Implement rollback system:

```python
class RollbackSystem:
    def __init__(self):
        self.checkpoints = []
        self.max_checkpoints = 5
    
    def create_checkpoint(self, version: str, metadata: dict):
        """Create rollback checkpoint"""
        # Your implementation
        pass
    
    def rollback(self, checkpoint_id: int = None) -> dict:
        """Rollback to checkpoint"""
        # Your implementation
        pass
    
    def list_checkpoints(self) -> list:
        """List all checkpoints"""
        # Your implementation
        pass
    
    def auto_rollback_on_error(self, error_rate_threshold: float) -> bool:
        """Check if auto-rollback needed"""
        # Your implementation
        pass

# Test
rollback = RollbackSystem()

# Create checkpoints
rollback.create_checkpoint("v1.0", {"deployment_time": "2024-01-01"})
rollback.create_checkpoint("v1.1", {"deployment_time": "2024-01-15"})
rollback.create_checkpoint("v2.0", {"deployment_time": "2024-02-01"})
rollback.create_checkpoint("v2.1", {"deployment_time": "2024-02-15"})

print("Checkpoints:")
for cp in rollback.list_checkpoints():
    print(f"  {cp}")

# Rollback
result = rollback.rollback(2)  # Rollback to checkpoint 2
print(f"\nRollback: {result}")
```

---

## Problem 6: Rate Limiter

Implement token bucket rate limiter:

```python
import time

class TokenBucketRateLimiter:
    def __init__(self, rate: float, capacity: int):
        self.rate = rate  # tokens per second
        self.capacity = capacity
        self.tokens = capacity
        self.last_refill = time.time()
    
    def allow_request(self) -> tuple[bool, int]:
        """
        Check if request allowed
        Returns: (allowed, remaining_tokens)
        """
        # Your implementation
        pass
    
    def wait_time(self) -> float:
        """Get wait time until next token"""
        # Your implementation
        pass

# Test
limiter = TokenBucketRateLimiter(rate=5, capacity=10)

# Simulate requests
for i in range(15):
    allowed, remaining = limiter.allow_request()
    print(f"Request {i+1}: Allowed={allowed}, Remaining={remaining}")
    time.sleep(0.1)  # Small delay between requests
```

---

## Problem 7: SLA Monitor

Implement SLA monitoring:

```python
class SLAMonitor:
    def __init__(self):
        self.slas = {}
        self.measurements = {}
    
    def define_sla(
        self,
        name: str,
        target: float,
        unit: str,
        threshold_direction: str  # "lower" or "higher" is better
    ):
        """Define SLA"""
        # Your implementation
        pass
    
    def record(self, name: str, value: float):
        """Record measurement"""
        # Your implementation
        pass
    
    def check_compliance(self, time_window: str = "1h") -> dict:
        """Check SLA compliance"""
        # Your implementation
        pass
    
    def get_violations(self) -> list:
        """Get current violations"""
        # Your implementation
        pass

# Test
monitor = SLAMonitor()

# Define SLAs
monitor.define_sla("availability", 99.9, "%", "higher")
monitor.define_sla("latency_p99", 1000, "ms", "lower")
monitor.define_sla("error_rate", 1.0, "%", "lower")

# Record measurements
for i in range(100):
    monitor.record("availability", 99.5 + random.random() * 0.5)
    monitor.record("latency_p99", 800 + random.randint(-200, 400))
    monitor.record("error_rate", random.random() * 2)

violations = monitor.get_violations()
print(f"Violations: {violations}")
```

---

## Problem 8: Post-Incident Review

Design PIR template:

```python
def create_pir_template(incident_data: dict) -> str:
    """
    Create Post-Incident Review template
    """
    # Your implementation
    pass

# Sample incident
incident = {
    "id": 42,
    "title": "LLM API outage",
    "severity": "P1",
    "start_time": "2024-02-15T14:30:00Z",
    "end_time": "2024-02-15T15:45:00Z",
    "affected_users": 5000,
    "timeline": [
        {"time": "14:30", "event": "Alert: High error rate"},
        {"time": "14:35", "event": "On-call engaged"},
        {"time": "14:45", "event": "Root cause identified: GPU OOM"},
        {"time": "15:00", "event": "Rollback initiated"},
        {"time": "15:30", "event": "Service restored"},
        {"time": "15:45", "event": "Incident closed"}
    ]
}

pir = create_pir_template(incident)
print(pir)
```

---

## Submission

Complete all implementations and be prepared to discuss:
- Deployment strategy selection criteria
- Health check best practices
- Incident response procedures
- SLA management approaches
