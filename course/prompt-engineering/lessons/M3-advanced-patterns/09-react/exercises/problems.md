# Exercise: ReAct Prompting

## Problem 1: Understanding ReAct

ReAct combines reasoning and acting. Fill in the diagram:

```
Thought: I need to find the current temperature in London.
Action: search
Action Input: "current temperature London London"
Observation: The current temperature in London is 18°C.

Thought: 
Action: 
Action Input: 
Observation: 
...
```

### Complete this ReAct sequence to answer:
"What's the weather like in London and should I bring an umbrella?"

---

## Problem 2: Designing a ReAct Prompt

Create a ReAct prompt for a research assistant that can:
1. Search the web
2. Read files
3. Calculate

### Template:
```
Question: [User's question]

You have access to the following tools:
- search: Search the web for information
- read_file: Read content from a file
- calculate: Perform mathematical calculations

Use the following format:
Thought: [your reasoning]
Action: [tool name]
Action Input: [tool input]
Observation: [tool output]
...
Thought: I have enough information to answer the question.
Final Answer: [your answer]
```

### Your Task:
Write a complete ReAct prompt to answer:
"What is the population of Tokyo multiplied by 2?"

---

## Problem 3: ReAct vs Chain-of-Thought

Compare how ReAct and CoT would solve this problem:

**Problem:** "Who was the first person to win both an Oscar and a Nobel Prize?"

### Chain-of-Thought approach:
```
Let me think about this step by step...
```

### ReAct approach:
```
Thought: I need to search for information about people who won both Oscar and Nobel.
Action: search
Action Input: "first person to win Oscar and Nobel Prize"
```

### Your Analysis:
| Aspect | CoT | ReAct |
|--------|-----|-------|
| Can access external info? | | |
| Shows reasoning? | | |
| Better for factual questions? | | |

---

## Problem 4: Implementing ReAct

Design a simple ReAct loop for a question-answering system:

```python
# Simplified ReAct loop
def react_loop(question, tools):
    # Initialize
    thought = "Let me answer this question."
    
    # Maximum 5 iterations
    for i in range(5):
        # YOUR CODE HERE: Decide action based on thought
        action = None  
        action_input = None
        
        # Execute action
        observation = execute_tool(action, action_input)
        
        # Update thought based on observation
        thought = f"After getting '{observation}', I think..."
        
        # Check if we can answer
        if can_answer(thought):
            return extract_answer(thought)
    
    return "Could not find answer"
```

### Your Task:
Write the decision logic for choosing actions:

---

## Problem 5: ReAct Limitations

For each scenario, identify potential issues with ReAct:

| Scenario | Potential Issue |
|----------|-----------------|
| Tool gives wrong information | |
| Infinite loop of actions | |
| Tools are unavailable | |
| Action produces unexpected output | |
