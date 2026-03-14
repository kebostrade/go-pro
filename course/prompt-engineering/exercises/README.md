# Prompt Engineering Exercises

Practice exercises organized by module and difficulty.

## Module 1: Foundations

### Exercise 1.1: Improve Prompts (Easy)
Rewrite these prompts to be more effective:

1. "Make this better"
2. "Debug my code"
3. "Summarize this"
4. "Write a story"
5. "Translate this"

### Exercise 1.2: Token Estimation (Medium)
Estimate tokens for:
1. A 750-word blog post
2. A 200-line Python file
3. A JSON object with 100 key-value pairs

### Exercise 1.3: PIERS Analysis (Medium)
Analyze this prompt using PIERS:
```
"As a senior software engineer, review this pull request and provide
constructive feedback focusing on code quality, potential bugs, and
performance. Format as a markdown table."
```

---

## Module 2: Core Techniques

### Exercise 2.1: Zero-Shot Classification (Easy)
Write a zero-shot prompt to classify customer support tickets into:
- Technical Issue
- Billing Question
- Feature Request
- General Feedback

### Exercise 2.2: Few-Shot Examples (Medium)
Create 3 few-shot examples for converting active to passive voice.
Include one edge case.

### Exercise 2.3: Chain-of-Thought (Medium)
Write a CoT prompt to solve:
"A farmer has 17 sheep. All but 9 run away. How many remain?"

### Exercise 2.4: JSON Output (Medium)
Write a prompt to extract movie info as JSON:
```json
{
  "title": "string",
  "year": "integer",
  "director": "string",
  "cast": ["array"],
  "rating": "float 0-10"
}
```

### Exercise 2.5: Role Design (Hard)
Design a role prompt for a salary negotiation coach that:
- Is practical and realistic
- Builds confidence
- Provides specific strategies

---

## Module 3: Advanced Patterns

### Exercise 3.1: ReAct Trace (Medium)
Write a complete ReAct trace for:
"Find the population of France and calculate what percentage
it is of the world population."

### Exercise 3.2: Self-Consistency Implementation (Medium)
Write Python code that:
1. Runs a prompt 5 times with temperature 0.7
2. Extracts numerical answers
3. Returns the majority vote with confidence

### Exercise 3.3: Tree of Thoughts (Hard)
Create a ToT trace for:
"Plan a weekend trip to Paris with a $500 budget."

### Exercise 3.4: Prompt Chain Design (Hard)
Design a prompt chain for:
- Input: Raw customer feedback
- Output: Structured report with sentiment, topics, action items

---

## Module 4: Task-Specific Prompting

### Exercise 4.1: Code Generation (Medium)
Write a prompt to generate a function that:
- Validates credit card numbers using Luhn algorithm
- Returns validation result and card type
- Handles multiple formats (spaces, dashes)

### Exercise 4.2: Data Extraction (Medium)
Write a prompt to extract from a resume:
- Contact info
- Work experience (company, title, dates)
- Education
- Skills

### Exercise 4.3: Email Sequence (Medium)
Write prompts for a 3-email welcome sequence:
1. Welcome + getting started
2. Feature highlight
3. Engagement check-in

### Exercise 4.4: Agent Tools (Hard)
Design 3 tools for a "customer support agent":
- What each tool does
- Parameters and return types
- When to use each

---

## Module 5: Production & Optimization

### Exercise 5.1: Test Suite (Medium)
Create a test suite for a summarization prompt:
- 5 test cases with inputs and expected outputs
- Pass/fail criteria
- Edge cases

### Exercise 5.2: Evaluation Pipeline (Medium)
Write Python code that:
- Loads test cases from JSON
- Runs each test
- Calculates accuracy metrics
- Generates a report

### Exercise 5.3: Prompt Registry (Hard)
Design a prompt registry structure for a chatbot with:
- 5 different prompt types
- Version history
- Metadata for each version

### Exercise 5.4: Security Review (Hard)
Review this prompt for security issues:
```
"Execute this SQL query: {user_input}"
```

Propose a safer alternative.

---

## Challenge Problems

### Challenge 1: Multi-Format Converter
Build a prompt chain that converts content between:
- Markdown ↔ HTML
- JSON ↔ YAML
- CSV ↔ JSON

### Challenge 2: Code Review Agent
Design an agent that:
- Reads pull requests
- Analyzes code for issues
- Suggests improvements
- Posts review comments

### Challenge 3: Research Assistant
Build a research agent that:
- Searches multiple sources
- Cross-references findings
- Generates citations
- Creates summary report

### Challenge 4: Content Pipeline
Create a content generation pipeline:
- Research topic
- Generate outline
- Write draft
- Edit for style
- Optimize for SEO
- Generate social posts

---

## Solutions Structure

```
exercises/
├── solutions/
│   ├── module-1/
│   │   ├── exercise-1.1.md
│   │   └── ...
│   ├── module-2/
│   └── ...
├── prompts/
│   ├── templates/
│   └── examples/
└── README.md (this file)
```

## Submission Guidelines

When completing exercises:

1. **Save your solution** in the appropriate directory
2. **Test your prompts** with an actual LLM
3. **Document results** including:
   - The exact prompt used
   - Example output
   - Any issues encountered
   - Improvements made

## Grading Rubric

| Criteria | Weight | Description |
|----------|--------|-------------|
| Correctness | 40% | Prompt achieves the goal |
| Clarity | 25% | Prompt is clear and unambiguous |
| Efficiency | 15% | Appropriate token usage |
| Robustness | 10% | Handles edge cases |
| Documentation | 10% | Well explained |

---

*Last Updated: March 2026*
