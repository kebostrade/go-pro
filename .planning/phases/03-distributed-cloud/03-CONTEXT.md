# Phase 3: Distributed & Cloud — Context

**Phase:** 3  
**Name:** Distributed & Cloud  
**Topics:** Kubernetes, NATS Events, AWS Lambda  
**Status:** Pending Research & Planning

---

## Decisions (Locked)

### Topic 9: Kubernetes Cloud-Native
- **Template:** `basic/projects/kubernetes/`
- **Focus:** K8s manifests, Helm charts, operator patterns
- **Libraries:** client-go, controller-runtime, operator-sdk
- **Deliverable:** Production-grade K8s deployment with operators

### Topic 10: NATS Event-Driven
- **Template:** `basic/projects/nats-events/`
- **Focus:** JetStream, publish/subscribe, queue groups
- **Library:** nats-io/nats.go v1.49+ (already in course/)
- **Deliverable:** Event-driven architecture with NATS JetStream

### Topic 11: AWS Lambda Serverless
- **Template:** `basic/projects/serverless/`
- **Focus:** SAM, serverless.yaml, Lambda URLs
- **Library:** aws/aws-lambda-go
- **Deliverable:** Serverless Go functions with SAM deployment

### Standard Stack Decisions
| Component | Decision | Rationale |
|-----------|----------|-----------|
| Go Version | 1.23+ | Standardized across all modules |
| K8s Client | client-go + controller-runtime | Official Kubernetes Go client |
| NATS Client | nats-io/nats.go | Official NATS client, JetStream support |
| Lambda Runtime | aws-lambda-go | Official AWS Lambda Go SDK |
| SAM CLI | AWS SAM CLI | AWS-provided serverless deployment |

---

## the agent's Discretion

The following are **NOT locked** — the planner researches and recommends:

1. **Operator SDK vs Raw client-go**: Should we use operator-sdk scaffolding or raw client-go for the Kubernetes operator template?
2. **Helm vs Kustomize**: Which Kubernetes templating approach to teach?
3. **NATS JetStream vs Core NATS**: Should we focus on JetStream persistence or core pub/sub?
4. **SAM vs Terraform for Lambda**: Which infrastructure-as-code approach?
5. **Lambda URLs vs API Gateway**: Should we use Lambda function URLs (simpler) or API Gateway (more features)?
6. **Lambda handler patterns**: Best practices for Go Lambda handler signatures

---

## Deferred Ideas (Out of Scope for Phase 3)

- ~~Kubernetes operator pattern (already in scope)~~
- ~~NATS JetStream (already in scope)~~
- ~~AWS Lambda (already in scope)~~
- Multi-region deployment patterns (Phase 4+)
- Service mesh (Istio/Linkerd) — Phase 4+ topic
- Cross-cloud Kubernetes (EKS/GKE/AKS) — Phase 4+ topic
- Lambda@Edge patterns — Phase 4+ topic
- NATS clustered/deployment — Phase 4+ topic

---

## Dependencies

- **Requires:** Phase 2 Microservices template (Docker Compose DNS, service discovery)
- **Leverages:** course/ module NATS dependencies (already have nats-io/nats.go v1.49.0)

---

## Phase 3 Task Breakdown

```
Phase 3: Distributed & Cloud
├── Task 9:  Template - Kubernetes (K8s manifests, Helm, operators)
├── Task 10: Template - NATS (JetStream, publish/subscribe)
└── Task 11: Template - AWS Lambda (SAM, serverless.yaml)
```

---

## Quality Gates (per template)

- [ ] `go build ./...` passes
- [ ] `go test ./...` passes with >80% coverage
- [ ] `golangci-lint run` passes
- [ ] `docker build` succeeds
- [ ] `docker-compose up` runs without errors
- [ ] K8s manifests validate (`kubectl apply --dry-run`)
- [ ] SAM template validates (`sam validate`)
- [ ] CI pipeline green on GitHub Actions
