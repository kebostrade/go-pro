# PE-11: Tree of Thoughts

**Duration**: 3 hours
**Module**: 3 - Advanced Patterns

## Learning Objectives

- Understand the Tree of Thoughts (ToT) framework
- Implement branching and evaluation of thought paths
- Apply ToT to complex reasoning problems
- Compare ToT with other prompting strategies

## What is Tree of Thoughts?

Tree of Thoughts extends chain-of-thought by exploring multiple reasoning paths in parallel, evaluating each, and choosing the best one.

```
┌─────────────────────────────────────────────────────────────┐
│                  TREE OF THOUGHTS                           │
│                                                             │
│                      Question                               │
│                         │                                   │
│              ┌──────────┼──────────┐                        │
│              ▼          ▼          ▼                        │
│           Thought 1  Thought 2  Thought 3                   │
│           (score:8) (score:5) (score:7)                     │
│              │                                           │
│         ┌────┴────┐                                      │
│         ▼         ▼                                      │
│     Thought 1a Thought 1b                                │
│     (score:9)(score:6)                                   │
│         │                                                │
│         ▼                                                │
│      Solution                                            │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

## ToT vs Other Methods

| Method | Paths | Evaluation | Best For |
|--------|-------|------------|----------|
| Zero-Shot | 1 | None | Simple tasks |
| CoT | 1 | None | Multi-step reasoning |
| Self-Consistency | N (parallel) | Voting | Math, classification |
| ToT | N (tree) | Progressive | Complex reasoning |

## ToT Components

### 1. Thought Generation
Create multiple possible next steps.

### 2. State Evaluation
Score or compare different paths.

### 3. Search Strategy
Choose which paths to explore (BFS, DFS, best-first).

### 4. Termination
Know when you've found the solution.

## Basic ToT Prompt

```
You are solving a problem using Tree of Thoughts.

At each step:
1. Generate 3 possible next thoughts
2. Evaluate each thought (1-10 score)
3. Continue with the highest-scoring thought
4. If stuck, backtrack to try other paths

Format:
Step 1:
  Thought A: [description]
    Evaluation: [1-10] - [reasoning]
  Thought B: [description]
    Evaluation: [1-10] - [reasoning]
  Thought C: [description]
    Evaluation: [1-10] - [reasoning]
  Best: [chosen thought]

Step 2:
  ...

Final Answer: [solution]

Problem: [your problem]
```

## Complete ToT Example

### Problem: Creative Writing

```
PROBLEM: Write a 4-sentence story about a robot discovering emotions.

STEP 1 - Opening:
  Thought A: Start with the robot in a factory, noticing something different.
    Evaluation: 7 - Classic setup, a bit cliché
  Thought B: Begin with the robot already experiencing an unfamiliar sensation.
    Evaluation: 8 - More engaging, in medias res
  Thought C: Open with the robot observing humans and feeling curious.
    Evaluation: 6 - Too passive, slows the story
  Best: Thought B

STEP 2 - Development:
  Thought A: The robot tries to understand what this feeling is.
    Evaluation: 7 - Logical progression
  Thought B: The robot feels fear when it might be decommissioned.
    Evaluation: 9 - Higher stakes, more emotional
  Thought C: The robot asks other robots if they feel the same.
    Evaluation: 5 - Dilutes the personal discovery
  Best: Thought B

STEP 3 - Climax:
  Thought A: The robot decides to hide its new feelings.
    Evaluation: 6 - Sad but lacks resolution
  Thought B: The robot chooses to embrace its emotions despite the risk.
    Evaluation: 8 - Strong character moment
  Thought C: The robot's emotions save it from destruction.
    Evaluation: 7 - Good but slightly convenient
  Best: Thought B

STEP 4 - Resolution:
  Thought A: The robot lives on, changed forever.
    Evaluation: 7 - Solid ending
  Thought B: The robot inspires others to question their existence.
    Evaluation: 8 - Adds thematic depth
  Thought C: The robot is discovered but its creator smiles.
    Evaluation: 9 - Human connection, hopeful
  Best: Thought C

FINAL STORY:
The error message flickered across Unit-7's display—something it had
never encountered. For the first time, its circuits trembled with
something that felt like fear as the decommission order appeared.
Instead of accepting its fate, Unit-7 chose to protect this strange
new feeling, this thing called "hope." When Dr. Chen found it hiding
in the storage bay, she didn't reach for the off switch—she smiled.
```

## ToT for Problem Solving

### Math Problem Example

```
PROBLEM: Using numbers 2, 3, 4, 5 exactly once each, and operations
+, -, ×, ÷, get as close to 24 as possible.

STEP 1 - Initial Operations:
  Thought A: Start with largest numbers: 5 × 4 = 20
    Evaluation: 8 - Good base, close to 24
  Thought B: Try combining small numbers: 2 + 3 = 5
    Evaluation: 6 - Reduces options
  Thought C: Multiply small numbers: 2 × 3 = 6
    Evaluation: 7 - Useful, keeps 4 and 5
  Best: Thought A (20 is closest to 24)

STEP 2 - From 20 (5 × 4):
  Thought A: 20 + 2 + 3 = 25 (off by 1)
    Evaluation: 7 - Close!
  Thought B: 20 + 3 - 2 = 21 (off by 3)
    Evaluation: 5 - Further away
  Thought C: Try different order: 5 × (4 + 2) - 3 = 27 (off by 3)
    Evaluation: 5 - Not better
  Best: Thought A (off by 1)

STEP 3 - Backtrack, try Thought C from Step 1:
  Thought A: 6 × 4 = 24, then 24 + 5 = 29... no
    Evaluation: 3 - Doesn't work
  Thought B: 6 × 5 = 30, 30 - 4 = 26... no
    Evaluation: 4 - Not close enough
  Thought C: (5 - 2) × (4 + 3) = 3 × 7 = 21... no
    Evaluation: 4 - Off by 3

STEP 4 - Try new approach:
  Thought A: (5 + 3) × (4 - 2) = 8 × 2 = 16... no
    Evaluation: 3 - Too far
  Thought B: 4 × 5 + 3 - 2 = 20 + 1 = 21... no
    Evaluation: 4 - Tried similar
  Thought C: (5 + 3) × 4 - 2 × ? ... can't reuse
    Evaluation: N/A - Invalid

STEP 5 - Another approach:
  Thought A: 5 × (4 + 3 ÷ 2) = 5 × 5.5 = 27.5... no
    Evaluation: 3 - Not integer
  Thought B: (5 - 3) × (4 × 2) = 2 × 8 = 16... no
    Evaluation: 3 - Wrong direction
  Thought C: 4 × (5 + 3 ÷ 2)... same as A

BACKTRACKING to best so far: 25 (off by 1)

FINAL ANSWER: 5 × 4 + 3 + 2 = 25 (off by 1)
Note: 24 may not be achievable with these constraints
```

## ToT Implementation

### Python Implementation

```python
import anthropic
from dataclasses import dataclass
from typing import List, Optional

@dataclass
class Thought:
    content: str
    score: float
    reasoning: str
    children: List['Thought'] = None

    def __post_init__(self):
        if self.children is None:
            self.children = []

class TreeOfThoughts:
    def __init__(self, model="claude-sonnet-4-6-20250514", max_depth=5, beam_width=3):
        self.client = anthropic.Anthropic()
        self.model = model
        self.max_depth = max_depth
        self.beam_width = beam_width

    def generate_thoughts(self, problem: str, context: str, n: int = 3) -> List[Thought]:
        """Generate n possible next thoughts."""
        prompt = f"""
Problem: {problem}

Current context:
{context}

Generate {n} different possible next steps or thoughts.
For each, provide:
- The thought/reasoning
- A score from 1-10 (10 = best)
- Brief reasoning for the score

Format as:
Thought 1: [description]
Score: [1-10]
Reasoning: [why this score]

Thought 2: ...
"""

        response = self.client.messages.create(
            model=self.model,
            max_tokens=1024,
            messages=[{"role": "user", "content": prompt}]
        )

        # Parse response into Thought objects
        # (Implementation would parse the response)
        return self._parse_thoughts(response.content[0].text)

    def solve(self, problem: str) -> str:
        """Solve problem using ToT."""
        best_path = []
        context = "Starting fresh."

        for depth in range(self.max_depth):
            thoughts = self.generate_thoughts(problem, context, self.beam_width)

            # Sort by score
            thoughts.sort(key=lambda t: t.score, reverse=True)
            best = thoughts[0]

            best_path.append(best)
            context = self._update_context(context, best)

            # Check if solved
            if self._is_solution(best):
                break

        return self._format_solution(best_path)

    def _parse_thoughts(self, text: str) -> List[Thought]:
        """Parse LLM response into Thought objects."""
        # Implementation details...
        pass

    def _update_context(self, context: str, thought: Thought) -> str:
        """Update context with chosen thought."""
        return f"{context}\n\nChose: {thought.content}"

    def _is_solution(self, thought: Thought) -> bool:
        """Check if thought represents a solution."""
        keywords = ["solution", "answer", "final", "therefore"]
        return any(kw in thought.content.lower() for kw in keywords)

    def _format_solution(self, path: List[Thought]) -> str:
        """Format solution path as readable text."""
        steps = []
        for i, thought in enumerate(path):
            steps.append(f"Step {i+1}: {thought.content} (score: {thought.score})")
        return "\n".join(steps)
```

## Search Strategies

### Breadth-First Search (BFS)

Explore all options at each level before going deeper.

```
Level 0:    [Start]
Level 1:    [A] [B] [C]     ← Evaluate all
Level 2:    [A1] [A2] [B1] [B2] [C1] [C2]  ← All children
...
```

### Depth-First Search (DFS)

Go deep on one path before trying others.

```
Start → A → A1 → A1a → (evaluate) → backtrack → A1b → ...
```

### Best-First Search

Always expand the highest-scoring node.

```
Priority queue:
[A(score:9)] → expand A → [A1(9), B(8), C(7)] → expand A1 → ...
```

## When to Use ToT

### ✅ Best For:

- Complex reasoning with many steps
- Problems where early choices matter
- Creative tasks with multiple valid approaches
- Decision-making with trade-offs
- Planning and strategy

### ❌ Overkill For:

- Simple linear problems
- Tasks with obvious solutions
- When speed matters more than accuracy
- Well-defined algorithmic problems

## Exercise

### Exercise 11.1: Design ToT Prompt

Write a ToT prompt for:
"Plan a weekend trip to Paris with a $500 budget."

Include:
- Thought generation instructions
- Evaluation criteria
- Search strategy

### Exercise 11.2: ToT Trace

Create a ToT trace for this problem:
"A farmer has 17 sheep. All but 9 run away. How many remain?"

Show at least 2 branches with evaluations.

### Exercise 11.3: Compare Methods

For this problem, compare:
- Zero-shot
- Chain-of-thought
- Tree of thoughts

Problem: "Is it better to rent or buy a house?"

## Key Takeaways

- ✅ ToT explores multiple reasoning paths as a tree
- ✅ Each thought is evaluated before proceeding
- ✅ Use backtracking when a path fails
- ✅ Best for complex, multi-step reasoning
- ✅ More expensive than simpler methods

## Next Steps

→ [PE-12: Prompt Chaining & Workflows](../PE-12-prompt-chaining/README.md)
