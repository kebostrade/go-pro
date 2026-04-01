# Phase 4: Specialized Domains — Context

**Phase:** 4
**Name:** Specialized Domains
**Topics:** ML/Gorgonia, Blockchain/Ethereum, IoT/MQTT, System Design
**Status:** Pending Research & Planning

---

## Decisions (Locked)

### Topic 12: Machine Learning with Gorgonia
- **Template:** `basic/projects/ml-gorgonia/`
- **Focus:** Tensor operations, model inference, neural network execution
- **Library:** gorgonia.org/gorgonia (active fork) or go-ml/gorgonia
- **Deliverable:** ML inference pipeline with pre-trained model loading

### Topic 13: Blockchain with Ethereum
- **Template:** `basic/projects/blockchain/`
- **Focus:** Smart contracts, wallet integration, web3 interactions
- **Library:** go-ethereum (ethereum/go-ethereum) for full node; web3go.xyz/web3 for lighter weight
- **Deliverable:** Ethereum wallet and smart contract interaction demo

### Topic 14: IoT with MQTT
- **Template:** `basic/projects/iot-mqtt/`
- **Focus:** Device telemetry, command/control, broker patterns
- **Library:** github.com/eclipse/paho.mqtt.golang
- **Broker:** eclipse-mosquitto (Docker)
- **Deliverable:** IoT sensor network with MQTT backend

### Topic 15: System Design with Golang
- **Template:** `basic/projects/system-design/`
- **Focus:** Architecture patterns, trade-offs, case studies
- **Deliverable:** System design case study implementations

### Standard Stack Decisions
| Component | Decision | Rationale |
|-----------|----------|-----------|
| Go Version | 1.23+ | Standardized across all modules |
| ML Framework | gorgonia | Native Go tensor operations |
| Blockchain | go-ethereum | Official Ethereum Foundation Go client |
| MQTT Client | eclipse/paho.mqtt.golang | Industry standard Go MQTT |
| MQTT Broker | mosquitto | Lightweight, production-grade broker |

---

## the agent's Discretion

The following are **NOT locked** — the planner researches and recommends:

1. **Gorgonia vs alternative ML**: Should we use gorgonia, or focus on ONNX runtime for Go?
2. **Ethereum library choice**: go-ethereum (full featured) vs web3go (lighter)?
3. **Smart contract language**: Solidity only, or also Vyper examples?
4. **MQTT QoS strategy**: Which QoS levels to emphasize for learning?
5. **System Design format**: Architecture diagrams only, or code-implemented case studies?

---

## Deferred Ideas (Out of Scope for Phase 4)

- ~~Kubernetes operator pattern (Phase 3)~~
- ~~NATS JetStream (Phase 3)~~
- ~~AWS Lambda (Phase 3)~~
- ~~GPT/broader LLM integration (separate track)~~
- ~~Custom blockchain (not Ethereum) — too complex~~
- ~~Real-time OS (RTOS) — out of scope~~
- ~~Hardware/circuit design — out of scope~~

---

## Dependencies

- **Requires:** Phase 3 cloud infrastructure patterns (especially NATS for event-driven)
- **Leverages:** existing advanced-topics/ content (IoT already exists)
- **Optional:** Phase 2 WebSocket patterns for real-time IoT dashboards

---

## Phase 4 Task Breakdown

```
Phase 4: Specialized Domains
├── Task 12: Template - ML with Gorgonia (tensor ops, model inference)
├── Task 13: Template - Blockchain with Ethereum (smart contracts, web3)
├── Task 14: Template - IoT with MQTT (mosquitto, sensors)
└── Task 15: Template - System Design (architecture patterns, case studies)
```

---

## Quality Gates (per template)

- [ ] `go build ./...` passes
- [ ] `go test ./...` passes with >80% coverage
- [ ] `golangci-lint run` passes
- [ ] `docker build` succeeds
- [ ] `docker-compose up` runs without errors (for IoT)
- [ ] CI pipeline green on GitHub Actions

---

## Notes

### Phase 4 Scale
This is the **largest phase** with 4 independent topics. Each topic:
- Is self-contained (no shared state between topics)
- Can be implemented in parallel by different agents
- Has its own template structure

### ML in Go Ecosystem
Gorgonia development has slowed. The community has fragmented:
- Original: goutm/gorgonia (archived)
- Active fork: gorgonia/gorgonia
- Alternative: go-ml (various packages)

Recommendation: Use gorgonia/gorgonia with fallback to ONNX runtime for production ML.

### Blockchain Complexity
Ethereum development requires:
- Test network (Sepolia or local Ganache)
- Wallet/accounts
- Smart contract ABIs
- Gas estimation

Keep the learning template focused on **read operations** (querying contracts) first, then writes.

### IoT Already Exists
The `advanced-topics/13-iot-mqtt/` directory has excellent content:
- README with architecture
- Examples directory (publisher, subscriber, device, processor, command/control)
- Broker config

The template should **extend** this, not replace it.
