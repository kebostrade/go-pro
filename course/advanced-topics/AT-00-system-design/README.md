# System Design with Golang

Design scalable, maintainable systems using Go principles and patterns.

## Course Overview

**Duration**: 8-10 weeks (self-paced)
**Level**: Intermediate to Advanced
**Prerequisites**: Basic Go knowledge, understanding of computer science fundamentals

## Learning Objectives

- Apply SOLID principles in Go
- Design for scalability and maintainability
- Choose appropriate architectural patterns
- Implement clean architecture
- Evaluate trade-offs in system design

## Course Structure

### [Module 1: Foundations](lessons/M1-foundations/README.md) (Week 1-2)
| Lesson | Topic | Duration |
|--------|-------|----------|
| [SD-01](lessons/M1-foundations/01-design-principles/README.md) | System Design Principles | 2h |
| [SD-02](lessons/M1-foundations/02-requirements-gathering/README.md) | Requirements Gathering & Analysis | 2h |
| [SD-03](lessons/M1-foundations/03-capacity-planning/README.md) | Capacity Planning & Estimation | 2h |
| [SD-04](lessons/M1-foundations/04-high-level-design/README.md) | High-Level Design | 3h |

### [Module 2: Architecture Patterns](lessons/M2-architecture/README.md) (Week 3-4)
| Lesson | Topic | Duration |
|--------|-------|----------|
| [SD-05](lessons/M2-architecture/05-layered-architecture/README.md) | Layered Architecture | 2h |
| [SD-06](lessons/M2-architecture/06-clean-architecture/README.md) | Clean Architecture | 3h |
| [SD-07](lessons/M2-architecture/07-microservices/README.md) | Microservices Architecture | 3h |
| [SD-08](lessons/M2-architecture/08-event-driven/README.md) | Event-Driven Architecture | 2h |

### [Module 3: Scalability](lessons/M3-scalability/README.md) (Week 5-6)
| Lesson | Topic | Duration |
|--------|-------|----------|
| [SD-09](lessons/M3-scalability/09-horizontal-vertical-scaling/README.md) | Horizontal & Vertical Scaling | 2h |
| [SD-10](lessons/M3-scalability/10-caching-strategies/README.md) | Caching Strategies | 3h |
| [SD-11](lessons/M3-scalability/11-load-balancing/README.md) | Load Balancing | 2h |
| [SD-12](lessons/M3-scalability/12-database-scaling/README.md) | Database Scaling | 3h |

### [Module 4: Data Management](lessons/M4-data/README.md) (Week 7)
| Lesson | Topic | Duration |
|--------|-------|----------|
| [SD-13](lessons/M4-data/13-database-design/README.md) | Database Design Patterns | 3h |
| [SD-14](lessons/M4-data/14-cqrs-event-sourcing/README.md) | CQRS & Event Sourcing | 3h |
| [SD-15](lessons/M4-data/15-data-consistency/README.md) | Data Consistency Models | 2h |
| [SD-16](lessons/M4-data/16-cdn-content-delivery/README.md) | CDN & Content Delivery | 2h |

### [Module 5: Production Systems](lessons/M5-production/README.md) (Week 8-9)
| Lesson | Topic | Duration |
|--------|-------|----------|
| [SD-17](lessons/M5-production/17-monitoring-observability/README.md) | Monitoring & Observability | 3h |
| [SD-18](lessons/M5-production/18-fault-tolerance/README.md) | Fault Tolerance & Resilience | 3h |
| [SD-19](lessons/M5-production/19-security/README.md) | Security Best Practices | 2h |
| [SD-20](lessons/M5-production/20-ci-cd-deployment/README.md) | CI/CD & Deployment | 2h |

## File Structure

```
AT-00-system-design/
├── README.md
├── lessons/
│   ├── M1-foundations/
│   │   ├── 01-design-principles/
│   │   ├── 02-requirements-gathering/
│   │   ├── 03-capacity-planning/
│   │   ├── 04-high-level-design/
│   │   └── README.md
│   ├── M2-architecture/
│   │   ├── 05-layered-architecture/
│   │   ├── 06-clean-architecture/
│   │   ├── 07-microservices/
│   │   ├── 08-event-driven/
│   │   └── README.md
│   ├── M3-scalability/
│   │   ├── 09-horizontal-vertical-scaling/
│   │   ├── 10-caching-strategies/
│   │   ├── 11-load-balancing/
│   │   ├── 12-database-scaling/
│   │   └── README.md
│   ├── M4-data/
│   │   ├── 13-database-design/
│   │   ├── 14-cqrs-event-sourcing/
│   │   ├── 15-data-consistency/
│   │   ├── 16-cdn-content-delivery/
│   │   └── README.md
│   └── M5-production/
│       ├── 17-monitoring-observability/
│       ├── 18-fault-tolerance/
│       ├── 19-security/
│       ├── 20-ci-cd-deployment/
│       └── README.md
├── examples/
├── exercises/
└── projects/
```

## Quick Start

```bash
# Navigate to the course
cd course/advanced-topics/AT-00-system-design

# Start with Module 1
cat lessons/M1-foundations/README.md
cat lessons/M1-foundations/01-design-principles/README.md
```

## Key Takeaways

- Prefer composition over inheritance
- Design interfaces at boundaries
- Keep packages cohesive
- Optimize for readability first
- Use dependency injection

## Certificate Requirements

Complete all modules + 1 project to earn your certificate:
- [ ] All 20 lessons reviewed
- [ ] 80% exercise completion
- [ ] 1 project completed with documentation

## Next Steps

**[AT-01: RESTful APIs](../AT-01-rest-apis/README.md)**

---

*Last Updated: March 2026*
