# Prompt Engineering Full Course

A comprehensive, hands-on course on prompt engineering for AI systems, covering fundamentals to advanced techniques.

## Course Overview

**Duration**: 8-10 weeks (self-paced)
**Level**: Beginner to Advanced
**Prerequisites**: Basic programming knowledge, familiarity with AI/LLMs helpful

## Learning Objectives

By the end of this course, you will be able to:
- Understand how LLMs process and respond to prompts
- Apply core prompting techniques (zero-shot, few-shot, chain-of-thought)
- Design prompts for specific tasks (code generation, analysis, creative writing)
- Implement advanced patterns (ReAct, Self-Consistency, Tree of Thoughts)
- Evaluate and optimize prompt performance
- Build production-ready prompt pipelines
- Debug and iterate on prompts systematically

## Course Structure

### Module 1: Foundations (Week 1-2)
| Lesson | Topic | Duration |
|--------|-------|----------|
| [PE-01](lessons/PE-01-introduction/README.md) | Introduction to Prompt Engineering | 2h |
| [PE-02](lessons/PE-02-llm-basics/README.md) | How LLMs Work: Tokens, Context, Attention | 3h |
| [PE-03](lessons/PE-03-prompt-anatomy/README.md) | Anatomy of a Good Prompt | 2h |
| [PE-04](lessons/PE-04-zero-shot/README.md) | Zero-Shot Prompting | 2h |

### Module 2: Core Techniques (Week 3-4)
| Lesson | Topic | Duration |
|--------|-------|----------|
| [PE-05](lessons/PE-05-few-shot/README.md) | Few-Shot Learning & Examples | 3h |
| [PE-06](lessons/PE-06-chain-of-thought/README.md) | Chain-of-Thought Prompting | 3h |
| [PE-07](lessons/PE-07-structured-output/README.md) | Structured Output & JSON | 2h |
| [PE-08](lessons/PE-08-role-prompting/README.md) | Role Prompting & Personas | 2h |

### Module 3: Advanced Patterns (Week 5-6)
| Lesson | Topic | Duration |
|--------|-------|----------|
| [PE-09](lessons/PE-09-react/README.md) | ReAct: Reasoning + Acting | 3h |
| [PE-10](lessons/PE-10-self-consistency/README.md) | Self-Consistency & Ensembling | 2h |
| [PE-11](lessons/PE-11-tree-of-thoughts/README.md) | Tree of Thoughts | 3h |
| [PE-12](lessons/PE-12-prompt-chaining/README.md) | Prompt Chaining & Workflows | 3h |

### Module 4: Task-Specific Prompting (Week 7)
| Lesson | Topic | Duration |
|--------|-------|----------|
| [PE-13](lessons/PE-13-code-generation/README.md) | Code Generation & Debugging | 3h |
| [PE-14](lessons/PE-14-data-analysis/README.md) | Data Analysis & Extraction | 2h |
| [PE-15](lessons/PE-15-creative-writing/README.md) | Creative Writing & Content | 2h |
| [PE-16](lessons/PE-16-agentic-prompts/README.md) | Agentic Prompts & Tool Use | 3h |

### Module 5: Production & Optimization (Week 8-10)
| Lesson | Topic | Duration |
|--------|-------|----------|
| [PE-17](lessons/PE-17-evaluation/README.md) | Prompt Evaluation & Metrics | 3h |
| [PE-18](lessons/PE-18-optimization/README.md) | Optimization Techniques | 2h |
| [PE-19](lessons/PE-19-security/README.md) | Security: Injection, Leaking | 2h |
| [PE-20](lessons/PE-20-production/README.md) | Production Systems & Best Practices | 3h |

## Projects

| Project | Description | Difficulty |
|---------|-------------|------------|
| [P1](projects/P1-prompt-library/README.md) | Build a Prompt Library System | Intermediate |
| [P2](projects/P2-rag-assistant/README.md) | RAG-Powered Document Assistant | Advanced |
| [P3](projects/P3-code-reviewer/README.md) | AI Code Review Agent | Advanced |
| [P4](projects/P4-multi-agent-system/README.md) | Multi-Agent Orchestration | Expert |

## Quick Start

```bash
# Clone or navigate to the course
cd course/prompt-engineering

# Start with Module 1
cat lessons/PE-01-introduction/README.md
```

## Tools & Resources

### Required
- Access to an LLM (Claude, GPT-4, etc.)
- Text editor or IDE
- Python (for projects)

### Recommended
- [Anthropic API](https://console.anthropic.com)
- [Promptfoo](https://promptfoo.dev) - Prompt testing
- [LangSmith](https://smith.langchain.com) - Prompt tracing

## File Structure

```
prompt-engineering/
├── README.md                 # This file
├── lessons/                  # Lesson content
│   ├── PE-01-introduction/
│   ├── PE-02-llm-basics/
│   └── ...
├── exercises/                # Practice problems
│   ├── solutions/
│   └── prompts/
├── examples/                 # Code examples
├── projects/                 # Course projects
├── cheatsheet.md             # Quick reference
└── resources/                # Additional materials
```

## Certificate Requirements

Complete all modules + 2 projects to earn certificate:
- [ ] All 20 lessons reviewed
- [ ] 80% exercise completion
- [ ] 2 projects completed with documentation

---

*Last Updated: March 2026*
