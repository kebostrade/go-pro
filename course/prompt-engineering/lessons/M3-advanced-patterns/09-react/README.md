# PE-09: ReAct - Reasoning + Acting

**Duration**: 3 hours
**Module**: 3 - Advanced Patterns

## Learning Objectives

- Understand the ReAct pattern and its components
- Implement ReAct for tool-use scenarios
- Design effective tool definitions
- Build agents that reason before acting

## What is ReAct?

ReAct (Reasoning + Acting) is a pattern where the AI alternates between thinking about a problem and taking actions to gather information or make changes.

```
┌─────────────────────────────────────────────────────────────┐
│                    REACT CYCLE                              │
│                                                             │
│   ┌──────────┐     ┌──────────┐     ┌──────────┐           │
│   │  THOUGHT │ ──► │  ACTION  │ ──► │OBSERVATION│           │
│   │          │     │          │     │          │           │
│   └──────────┘     └──────────┘     └──────────┘           │
│        ▲                                  │                 │
│        └──────────────────────────────────┘                 │
│                   (repeat until done)                       │
│                                                             │
│   THOUGHT: "I need to find the current stock price"         │
│   ACTION: search("AAPL stock price today")                  │
│   OBSERVATION: "AAPL is trading at $178.52"                 │
│   THOUGHT: "Now I can answer the user's question"           │
└─────────────────────────────────────────────────────────────┘
```

## ReAct Components

### 1. Thought
The model's internal reasoning about what to do next.

### 2. Action
A specific tool or function to call.

### 3. Observation
The result of the action, fed back to the model.

### 4. Answer
The final response to the user (when done).

## Basic ReAct Structure

```
You have access to the following tools:
- search(query): Search the web for information
- calculator(expression): Evaluate a mathematical expression
- weather(location): Get current weather for a location

Use the following format:

Thought: [your reasoning about what to do]
Action: [tool name]
Action Input: [input for the tool]
Observation: [tool result - this will be provided]
... (repeat Thought/Action/Observation as needed)
Thought: [final reasoning]
Final Answer: [your answer to the user]

Question: What's the weather in Paris, and is it warmer than London?
```

## Complete ReAct Example

```
You are a helpful assistant with access to these tools:

TOOLS:
- search(query): Search the web
- calculator(expression): Do math calculations
- get_stock_price(ticker): Get current stock price

FORMAT:
Thought: [reasoning]
Action: [tool_name]
Action Input: [input]
Observation: [result]
... (repeat)
Final Answer: [response]

QUESTION: If I bought 10 shares of Apple stock at $150 each,
what's my profit if I sell at today's price?

---

Thought: I need to find Apple's current stock price first.
Action: get_stock_price
Action Input: AAPL
Observation: $178.52

Thought: Now I need to calculate the profit.
- Purchase price: 10 shares × $150 = $1,500
- Current value: 10 shares × $178.52 = $1,785.20
- Profit: $1,785.20 - $1,500 = $285.20

Let me verify this calculation.
Action: calculator
Action Input: (10 * 178.52) - (10 * 150)
Observation: 285.2

Thought: I have confirmed the calculation. The profit is $285.20.
Final Answer: Your profit would be $285.20. You bought 10 shares
at $150 each (total $1,500), and at today's price of $178.52,
your holdings are worth $1,785.20, giving you a profit of $285.20.
```

## Designing Tool Definitions

### Tool Definition Structure

```
TOOL: [name]
DESCRIPTION: [what it does]
PARAMETERS:
  - [param1] (type): [description]
  - [param2] (type): [description]
RETURNS: [what it returns]
EXAMPLE USAGE: [example]
```

### Example Tool Definitions

```
TOOLS:

TOOL: search
DESCRIPTION: Search the web for current information
PARAMETERS:
  - query (string): The search query
RETURNS: Top 3 search results with titles and snippets
EXAMPLE: search("current population of Tokyo")

TOOL: get_weather
DESCRIPTION: Get current weather conditions
PARAMETERS:
  - location (string): City name or coordinates
  - units (string, optional): "celsius" or "fahrenheit"
RETURNS: Temperature, conditions, humidity
EXAMPLE: get_weather("New York", "fahrenheit")

TOOL: query_database
DESCRIPTION: Execute a SQL query on the company database
PARAMETERS:
  - query (string): SQL SELECT query (read-only)
RETURNS: Query results as JSON
EXAMPLE: query_database("SELECT * FROM users LIMIT 10")

TOOL: send_email
DESCRIPTION: Send an email to a recipient
PARAMETERS:
  - to (string): Recipient email address
  - subject (string): Email subject line
  - body (string): Email content
RETURNS: Confirmation of send
EXAMPLE: send_email("user@example.com", "Hello", "Message body")
```

## ReAct Best Practices

### Practice 1: Clear Thought Process

```
❌ Vague thought:
Thought: I'll look it up.
Action: search
Action Input: stuff

✅ Clear thought:
Thought: The user wants to know the population of Tokyo.
I need to search for current population data since my
training data may be outdated.
Action: search
Action Input: "Tokyo population 2024"
```

### Practice 2: Verify Important Results

```
Thought: I found that the population is 14 million.
Let me verify this with another source.
Action: search
Action Input: "Tokyo metropolitan area population official statistics"
```

### Practice 3: Know When to Stop

```
Thought: I have gathered all the information needed to answer
the user's question. I have:
- Current stock price: $178.52
- Purchase price: $150
- Number of shares: 10
I can now calculate and provide the final answer.
Final Answer: ...
```

### Practice 4: Handle Tool Failures

```
Observation: Error: Could not fetch stock price. API timeout.

Thought: The stock API failed. Let me try searching for
the current price instead.
Action: search
Action Input: "AAPL stock price today"
```

## ReAct Variations

### Variation 1: Simplified ReAct

```
Think: What should I do?
Act: [tool](input)
Observe: [result]
[repeat]

Question: What's 25% of 847?

Think: I need to calculate 0.25 × 847
Act: calculator(0.25 * 847)
Observe: 211.75
Answer: 25% of 847 is 211.75
```

### Variation 2: Plan-and-Execute

```
First, plan your approach:
1. [step 1]
2. [step 2]
3. [step 3]

Then execute each step:

Step 1: [action]
Result: [observation]

Step 2: [action]
Result: [observation]

Final Answer: [response]
```

### Variation 3: Tool Selection First

```
Available tools: [list]

Select which tools you'll need for this task:
- Tool 1: [reason]
- Tool 2: [reason]

Now proceed with ReAct:
Thought: ...
```

## ReAct in Code

### Python Implementation

```python
import anthropic
import json
import re

client = anthropic.Anthropic()

# Define tools with safe implementations
def calculator(expression):
    """Safely evaluate basic math expressions."""
    # Only allow numbers, operators, parentheses, and spaces
    if not re.match(r'^[\d\s+\-*/().]+$', expression):
        return "Error: Invalid expression"
    try:
        # Use a simple expression parser instead of eval
        # In production, use a proper math library
        result = float(eval(expression))  # Note: In production, use a safer alternative
        return str(result)
    except Exception as e:
        return f"Error: {e}"

tools = {
    "search": lambda q: f"Search results for: {q}",
    "calculator": calculator,
    "weather": lambda loc: f"Weather in {loc}: Sunny, 72°F"
}

def run_react(question, max_iterations=5):
    prompt = f"""
You are a ReAct agent with these tools:
- search(query): Search the web
- calculator(expression): Do math
- weather(location): Get weather

Format:
Thought: [reasoning]
Action: [tool]
Action Input: [input]
Observation: [result]
Final Answer: [answer when done]

Question: {question}
"""

    conversation = [{"role": "user", "content": prompt}]

    for _ in range(max_iterations):
        response = client.messages.create(
            model="claude-sonnet-4-6-20250514",
            max_tokens=1024,
            messages=conversation
        )

        output = response.content[0].text
        conversation.append({"role": "assistant", "content": output})

        # Check if done
        if "Final Answer:" in output:
            return output.split("Final Answer:")[-1].strip()

        # Parse and execute action
        if "Action:" in output and "Action Input:" in output:
            action = output.split("Action:")[1].split("\n")[0].strip()
            action_input = output.split("Action Input:")[1].split("\n")[0].strip()

            # Execute tool
            if action in tools:
                result = tools[action](action_input)
                observation = f"Observation: {result}"
                conversation.append({"role": "user", "content": observation})

    return "Max iterations reached"

# Run
answer = run_react("What's the weather in Tokyo?")
print(answer)
```

## ReAct vs Other Patterns

| Pattern | Best For | Limitations |
|---------|----------|-------------|
| Zero-Shot | Simple, single-step tasks | No tool use |
| CoT | Multi-step reasoning | No external info |
| ReAct | Tasks requiring tools/external data | More complex setup |

## When to Use ReAct

### ✅ Use ReAct When:

- Task requires external information
- Multiple tools/actions needed
- Verification across sources required
- Complex multi-step workflows

### ❌ Skip ReAct When:

- Task is simple and self-contained
- No tools or external data needed
- Speed is critical
- Single-step answer suffices

## Exercise

### Exercise 9.1: Design Tools

Design 3 tools for a "travel planning assistant" with:
- Tool names and descriptions
- Parameters and return types
- Example usage

### Exercise 9.2: Write ReAct Trace

Write a complete ReAct trace for:
"Find the population of France and calculate what percentage
it is of the world population."

### Exercise 9.3: Handle Failure

Write a ReAct trace that handles this failure:
```
Observation: Error: API rate limit exceeded. Try again later.
```

## Key Takeaways

- ✅ ReAct = Reasoning + Acting in a loop
- ✅ Components: Thought → Action → Observation → (repeat)
- ✅ Design clear, well-documented tools
- ✅ Model should know when to stop
- ✅ Handle tool failures gracefully

## Next Steps

→ [PE-10: Self-Consistency & Ensembling](../PE-10-self-consistency/README.md)
