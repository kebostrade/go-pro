# Tutorial 5: AI Agents - ReAct Pattern

```
┌──────────────────────────────────────────────────────────────────────────┐
│ 🔴 ADVANCED                                    ⏱️  60 minutes             │
├──────────────────────────────────────────────────────────────────────────┤
│                                                                          │
│  🎯 PROJECT: Coding Assistant Agent                                     │
│                                                                          │
│  📚 WHAT YOU'LL LEARN:                                                   │
│     ✓ AI agent architecture                                             │
│     ✓ ReAct pattern (Reasoning + Acting)                                │
│     ✓ Tool calling and function execution                               │
│     ✓ Multi-step reasoning                                              │
│     ✓ Agent loops and state management                                  │
│                                                                          │
│  🛠️ TECH STACK: OpenAI Function Calling, Go                            │
│                                                                          │
└──────────────────────────────────────────────────────────────────────────┘
```

---

## 🤖 Understanding AI Agents

### What is an AI Agent?

An **AI Agent** is an autonomous system that can:
1. **Perceive** its environment
2. **Reason** about what to do
3. **Act** using tools
4. **Learn** from results

```
┌─────────────────────────────────────────────────────────────────┐
│  Traditional LLM:                                                │
│    Input → LLM → Output (one-shot)                              │
│                                                                  │
│  AI Agent:                                                       │
│    Task → [Think → Act → Observe] → ... → Result               │
│           └─────── Loop until done ──────┘                      │
│                                                                  │
│  Agent can:                                                      │
│    ✓ Use tools (search, calculator, APIs)                       │
│    ✓ Make multi-step plans                                      │
│    ✓ Correct mistakes                                           │
│    ✓ Adapt to new information                                   │
└─────────────────────────────────────────────────────────────────┘
```

---

## 🧠 The ReAct Pattern

**ReAct** = **Rea**soning + **Act**ing

```
┌──────────────────────────────────────────────────────────────────┐
│  Step 1: THOUGHT                                                 │
│    "I need to find the current weather in San Francisco"        │
│                                                                   │
│  Step 2: ACTION                                                  │
│    Tool: get_weather("San Francisco")                           │
│                                                                   │
│  Step 3: OBSERVATION                                             │
│    Result: "72°F, Sunny"                                         │
│                                                                   │
│  Step 4: THOUGHT                                                 │
│    "Now I have the weather, I can answer the user"              │
│                                                                   │
│  Step 5: FINAL ANSWER                                            │
│    "The weather in San Francisco is 72°F and sunny."            │
└──────────────────────────────────────────────────────────────────┘
```

---

## 🚀 Project: Coding Assistant Agent

### Step 1: Define Tools

```go
package main

import (
    "context"
    "encoding/json"
    "fmt"
    "os/exec"
    "strings"
)

type Tool struct {
    Name        string
    Description string
    Parameters  map[string]interface{}
    Execute     func(ctx context.Context, args map[string]interface{}) (string, error)
}

// Tool: Run Go code
var RunCodeTool = Tool{
    Name:        "run_go_code",
    Description: "Execute Go code and return the output",
    Parameters: map[string]interface{}{
        "type": "object",
        "properties": map[string]interface{}{
            "code": map[string]interface{}{
                "type":        "string",
                "description": "The Go code to execute",
            },
        },
        "required": []string{"code"},
    },
    Execute: func(ctx context.Context, args map[string]interface{}) (string, error) {
        code := args["code"].(string)
        
        // Write code to temp file
        tmpFile := "/tmp/main.go"
        if err := os.WriteFile(tmpFile, []byte(code), 0644); err != nil {
            return "", err
        }
        
        // Run the code
        cmd := exec.CommandContext(ctx, "go", "run", tmpFile)
        output, err := cmd.CombinedOutput()
        
        if err != nil {
            return fmt.Sprintf("Error: %s\nOutput: %s", err, output), nil
        }
        
        return string(output), nil
    },
}

// Tool: Search documentation
var SearchDocsTool = Tool{
    Name:        "search_docs",
    Description: "Search Go documentation for a topic",
    Parameters: map[string]interface{}{
        "type": "object",
        "properties": map[string]interface{}{
            "query": map[string]interface{}{
                "type":        "string",
                "description": "The search query",
            },
        },
        "required": []string{"query"},
    },
    Execute: func(ctx context.Context, args map[string]interface{}) (string, error) {
        query := args["query"].(string)
        
        // Simulate doc search (in real app, use actual search)
        docs := map[string]string{
            "goroutines": "Goroutines are lightweight threads managed by Go runtime...",
            "channels":   "Channels are typed conduits for communication...",
            "interfaces": "Interfaces define behavior through method sets...",
        }
        
        for key, doc := range docs {
            if strings.Contains(strings.ToLower(query), key) {
                return doc, nil
            }
        }
        
        return "No documentation found for: " + query, nil
    },
}

// Tool: Analyze code
var AnalyzeCodeTool = Tool{
    Name:        "analyze_code",
    Description: "Analyze Go code for issues and improvements",
    Parameters: map[string]interface{}{
        "type": "object",
        "properties": map[string]interface{}{
            "code": map[string]interface{}{
                "type":        "string",
                "description": "The Go code to analyze",
            },
        },
        "required": []string{"code"},
    },
    Execute: func(ctx context.Context, args map[string]interface{}) (string, error) {
        code := args["code"].(string)
        
        issues := []string{}
        
        // Simple static analysis
        if !strings.Contains(code, "package") {
            issues = append(issues, "Missing package declaration")
        }
        if strings.Contains(code, "panic(") {
            issues = append(issues, "Consider using error handling instead of panic")
        }
        if !strings.Contains(code, "error") && strings.Contains(code, "func") {
            issues = append(issues, "Consider adding error handling")
        }
        
        if len(issues) == 0 {
            return "Code looks good! No major issues found.", nil
        }
        
        return "Issues found:\n- " + strings.Join(issues, "\n- "), nil
    },
}
```

### Step 2: Build Agent Core

```go
type Agent struct {
    client   *openai.Client
    tools    map[string]Tool
    messages []openai.ChatCompletionMessage
    maxSteps int
}

func NewAgent(apiKey string) *Agent {
    agent := &Agent{
        client:   openai.NewClient(apiKey),
        tools:    make(map[string]Tool),
        maxSteps: 10,
        messages: []openai.ChatCompletionMessage{
            {
                Role: openai.ChatMessageRoleSystem,
                Content: `You are a helpful coding assistant. You can:
1. Run Go code
2. Search documentation
3. Analyze code for issues

Think step by step and use tools when needed.`,
            },
        },
    }
    
    // Register tools
    agent.RegisterTool(RunCodeTool)
    agent.RegisterTool(SearchDocsTool)
    agent.RegisterTool(AnalyzeCodeTool)
    
    return agent
}

func (a *Agent) RegisterTool(tool Tool) {
    a.tools[tool.Name] = tool
}

func (a *Agent) getToolDefinitions() []openai.FunctionDefinition {
    defs := make([]openai.FunctionDefinition, 0, len(a.tools))
    
    for _, tool := range a.tools {
        defs = append(defs, openai.FunctionDefinition{
            Name:        tool.Name,
            Description: tool.Description,
            Parameters:  tool.Parameters,
        })
    }
    
    return defs
}
```

### Step 3: Implement ReAct Loop

```go
func (a *Agent) Run(ctx context.Context, task string) (string, error) {
    // Add user task
    a.messages = append(a.messages, openai.ChatCompletionMessage{
        Role:    openai.ChatMessageRoleUser,
        Content: task,
    })
    
    for step := 0; step < a.maxSteps; step++ {
        fmt.Printf("\n--- Step %d ---\n", step+1)
        
        // 1. THINK: Ask LLM what to do next
        resp, err := a.client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
            Model:     openai.GPT4,
            Messages:  a.messages,
            Functions: a.getToolDefinitions(),
            Temperature: 0.1,
        })
        
        if err != nil {
            return "", fmt.Errorf("LLM error: %w", err)
        }
        
        message := resp.Choices[0].Message
        
        // 2. Check if agent wants to use a tool
        if message.FunctionCall != nil {
            // ACT: Execute the tool
            fmt.Printf("Thought: Using tool '%s'\n", message.FunctionCall.Name)
            
            result, err := a.executeTool(ctx, message.FunctionCall)
            if err != nil {
                return "", fmt.Errorf("tool execution error: %w", err)
            }
            
            // OBSERVE: Add tool result to conversation
            fmt.Printf("Observation: %s\n", result)
            
            a.messages = append(a.messages,
                message,
                openai.ChatCompletionMessage{
                    Role:    openai.ChatMessageRoleFunction,
                    Name:    message.FunctionCall.Name,
                    Content: result,
                },
            )
            
            continue
        }
        
        // 3. Agent has final answer
        fmt.Printf("Final Answer: %s\n", message.Content)
        return message.Content, nil
    }
    
    return "", fmt.Errorf("max steps reached without answer")
}

func (a *Agent) executeTool(ctx context.Context, call *openai.FunctionCall) (string, error) {
    tool, ok := a.tools[call.Name]
    if !ok {
        return "", fmt.Errorf("unknown tool: %s", call.Name)
    }
    
    // Parse arguments
    var args map[string]interface{}
    if err := json.Unmarshal([]byte(call.Arguments), &args); err != nil {
        return "", fmt.Errorf("failed to parse arguments: %w", err)
    }
    
    // Execute tool
    result, err := tool.Execute(ctx, args)
    if err != nil {
        return "", fmt.Errorf("tool execution failed: %w", err)
    }
    
    return result, nil
}
```

### Step 4: Use the Agent

```go
func main() {
    ctx := context.Background()
    agent := NewAgent(os.Getenv("OPENAI_API_KEY"))
    
    tasks := []string{
        "Write a Go function that calculates fibonacci numbers and run it",
        "Explain how goroutines work",
        "Analyze this code: func divide(a, b int) int { return a / b }",
    }
    
    for i, task := range tasks {
        fmt.Printf("\n\n========== Task %d ==========\n", i+1)
        fmt.Printf("User: %s\n", task)
        
        answer, err := agent.Run(ctx, task)
        if err != nil {
            fmt.Printf("Error: %v\n", err)
            continue
        }
        
        fmt.Printf("\nAgent: %s\n", answer)
    }
}
```

**📤 Expected Output:**
```
========== Task 1 ==========
User: Write a Go function that calculates fibonacci numbers and run it

--- Step 1 ---
Thought: Using tool 'run_go_code'
Observation: 0 1 1 2 3 5 8 13 21 34

--- Step 2 ---
Final Answer: I've created and executed a fibonacci function. The output shows 
the first 10 fibonacci numbers: 0, 1, 1, 2, 3, 5, 8, 13, 21, 34.

Agent: I've created and executed a fibonacci function...


========== Task 2 ==========
User: Explain how goroutines work

--- Step 1 ---
Thought: Using tool 'search_docs'
Observation: Goroutines are lightweight threads managed by Go runtime...

--- Step 2 ---
Final Answer: Goroutines are lightweight threads managed by the Go runtime. 
They allow concurrent execution with minimal overhead...

Agent: Goroutines are lightweight threads...


========== Task 3 ==========
User: Analyze this code: func divide(a, b int) int { return a / b }

--- Step 1 ---
Thought: Using tool 'analyze_code'
Observation: Issues found:
- Missing package declaration
- Consider adding error handling

--- Step 2 ---
Final Answer: The code has two issues:
1. Missing package declaration
2. No error handling for division by zero

Suggested fix:
func divide(a, b int) (int, error) {
    if b == 0 {
        return 0, errors.New("division by zero")
    }
    return a / b, nil
}

Agent: The code has two issues...
```

---

## 🎨 Advanced Agent Patterns

### 1. Planning Agent

```go
func (a *Agent) Plan(ctx context.Context, task string) ([]string, error) {
    prompt := fmt.Sprintf(`Break down this task into steps:
"%s"

Return a numbered list of steps.`, task)
    
    resp, err := a.client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
        Model: openai.GPT4,
        Messages: []openai.ChatCompletionMessage{
            {Role: openai.ChatMessageRoleUser, Content: prompt},
        },
    })
    
    if err != nil {
        return nil, err
    }
    
    steps := strings.Split(resp.Choices[0].Message.Content, "\n")
    return steps, nil
}
```

### 2. Self-Correction

```go
func (a *Agent) Verify(ctx context.Context, answer string) (bool, string, error) {
    prompt := fmt.Sprintf(`Verify this answer is correct:
"%s"

If incorrect, explain why and provide the correct answer.`, answer)
    
    // Use LLM to verify
    // Return: isCorrect, correction, error
}
```

### 3. Memory and State

```go
type AgentMemory struct {
    ShortTerm []Message  // Recent conversation
    LongTerm  []Fact     // Learned facts
    WorkingMemory map[string]interface{} // Current task state
}

func (a *Agent) Remember(key string, value interface{}) {
    a.memory.WorkingMemory[key] = value
}

func (a *Agent) Recall(key string) (interface{}, bool) {
    val, ok := a.memory.WorkingMemory[key]
    return val, ok
}
```

---

## 🎯 Challenges

### Challenge 1: Multi-Tool Agent
Build an agent that can use 5+ different tools.

### Challenge 2: Conversational Agent
Add conversation history and context management.

### Challenge 3: Error Recovery
Implement retry logic and error recovery.

### Challenge 4: Agent Evaluation
Build a system to evaluate agent performance.

---

## ✅ What You Learned

```
┌──────────────────────────────────────────────────────────────────┐
│  ✓ AI agent architecture                                         │
│  ✓ ReAct pattern (Reasoning + Acting)                            │
│  ✓ Tool definition and registration                              │
│  ✓ Function calling with OpenAI                                  │
│  ✓ Multi-step reasoning loops                                    │
│  ✓ State management                                              │
│  ✓ Planning and verification                                     │
│  ✓ Error handling in agents                                      │
└──────────────────────────────────────────────────────────────────┘
```

---

## 🚀 Next Steps

**Practice:**
1. Build the coding assistant
2. Add more tools
3. Test with complex tasks
4. Implement planning

**Next Tutorial:**
[Tutorial 6: Production AI Systems →](06_PRODUCTION_AI.md)

Learn how to deploy AI agents to production!

---

**💡 Pro Tip**: Start with simple tools and gradually add complexity. Test each tool independently first!

