# Tutorial 2: Prompt Engineering - The Art of Talking to AI

```
┌──────────────────────────────────────────────────────────────────────────┐
│ 🟡 INTERMEDIATE                                ⏱️  25 minutes             │
├──────────────────────────────────────────────────────────────────────────┤
│                                                                          │
│  🎯 GOAL: Master prompt design for better AI responses                  │
│                                                                          │
│  📚 WHAT YOU'LL LEARN:                                                   │
│     ✓ Prompt engineering fundamentals                                  │
│     ✓ System, user, and assistant roles                               │
│     ✓ Few-shot learning techniques                                     │
│     ✓ Chain-of-thought prompting                                       │
│     ✓ Prompt templates and patterns                                    │
│                                                                          │
│  🛠️ PROJECT: Prompt Template System                                     │
│                                                                          │
└──────────────────────────────────────────────────────────────────────────┘
```

---

## 🎨 What is Prompt Engineering?

**Prompt Engineering** is the practice of designing inputs (prompts) to get the best possible outputs from LLMs.

```
┌─────────────────────────────────────────────────────────────────┐
│  Bad Prompt:                                                     │
│    "Write code"                                                  │
│                                                                  │
│  Good Prompt:                                                    │
│    "Write a Go function that validates email addresses using    │
│     regex. Include error handling and unit tests."              │
│                                                                  │
│  Result: 10x better output with clear instructions!             │
└─────────────────────────────────────────────────────────────────┘
```

---

## 🏗️ Prompt Structure

### The Three Roles

```go
messages := []openai.ChatCompletionMessage{
    // 1. SYSTEM: Sets behavior and context
    {
        Role: "system",
        Content: "You are an expert Go programmer with 10 years of experience.",
    },
    
    // 2. USER: Your request/question
    {
        Role: "user",
        Content: "How do I implement a binary search tree in Go?",
    },
    
    // 3. ASSISTANT: AI's response (or examples)
    {
        Role: "assistant",
        Content: "Here's a binary search tree implementation...",
    },
}
```

### System Message Best Practices

```
┌─────────────────────────────────────────────────────────────────┐
│  ✅ Good System Messages:                                        │
│                                                                  │
│  "You are a helpful coding assistant that writes clean,         │
│   well-documented Go code following best practices."            │
│                                                                  │
│  "You are a technical writer. Explain complex topics in         │
│   simple terms with examples. Use analogies when helpful."      │
│                                                                  │
│  "You are a code reviewer. Analyze code for bugs, security      │
│   issues, and performance problems. Be constructive."           │
│                                                                  │
│  ❌ Avoid:                                                       │
│  "You are smart" (too vague)                                    │
│  "Do what I say" (not helpful)                                  │
└─────────────────────────────────────────────────────────────────┘
```

---

## 📝 Prompt Patterns

### Pattern 1: Clear Instructions

**Bad:**
```
"Make this better"
```

**Good:**
```
"Refactor this Go function to:
1. Add error handling
2. Improve variable names
3. Add comments
4. Follow Go conventions"
```

### Pattern 2: Provide Context

**Bad:**
```
"Fix this code"
```

**Good:**
```
"This Go function should validate user input but it's not working.
The input should be a valid email address.
Current error: regex always returns false.

Code:
func validateEmail(email string) bool {
    // ... code here
}
"
```

### Pattern 3: Specify Format

**Bad:**
```
"Explain goroutines"
```

**Good:**
```
"Explain goroutines in Go using this format:
1. Definition (2 sentences)
2. Simple code example
3. Common use case
4. One gotcha to avoid

Target audience: Beginner Go developers"
```

### Pattern 4: Use Examples (Few-Shot Learning)

```go
systemMsg := `You are a code commenter. Add clear comments to code.

Example 1:
Input: func add(a, b int) int { return a + b }
Output: 
// add returns the sum of two integers
func add(a, b int) int {
    return a + b
}

Example 2:
Input: func max(nums []int) int { m := nums[0]; for _, n := range nums { if n > m { m = n } }; return m }
Output:
// max returns the largest number in a slice
func max(nums []int) int {
    m := nums[0]
    for _, n := range nums {
        if n > m {
            m = n
        }
    }
    return m
}

Now comment the user's code following this pattern.`
```

---

## 🧠 Chain-of-Thought Prompting

Asking the AI to "think step by step" improves reasoning.

### Without Chain-of-Thought

```
User: "Is this Go code thread-safe?"

AI: "No, it's not thread-safe."
```

### With Chain-of-Thought

```
User: "Is this Go code thread-safe? Think step by step:
1. Identify shared state
2. Check for synchronization
3. Look for race conditions
4. Provide verdict with explanation"

AI: "Let me analyze step by step:

1. Shared state: The 'counter' variable is accessed by multiple goroutines
2. Synchronization: No mutex or channel is used
3. Race conditions: Yes, counter++ is not atomic
4. Verdict: NOT thread-safe. Multiple goroutines can read and write 
   'counter' simultaneously, causing race conditions.

Fix: Add sync.Mutex or use atomic operations."
```

---

## 🎯 Project: Prompt Template System

Let's build a reusable prompt template system.

### Step 1: Create Template Structure

```go
package main

import (
	"bytes"
	"fmt"
	"text/template"
)

// PromptTemplate represents a reusable prompt
type PromptTemplate struct {
	Name        string
	System      string
	UserTemplate string
	Examples    []Example
}

// Example for few-shot learning
type Example struct {
	Input  string
	Output string
}

// Execute fills the template with data
func (pt *PromptTemplate) Execute(data map[string]interface{}) (string, error) {
	tmpl, err := template.New("prompt").Parse(pt.UserTemplate)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}

// GetMessages returns OpenAI-compatible messages
func (pt *PromptTemplate) GetMessages(userPrompt string) []map[string]string {
	messages := []map[string]string{
		{"role": "system", "content": pt.System},
	}

	// Add examples (few-shot)
	for _, ex := range pt.Examples {
		messages = append(messages,
			map[string]string{"role": "user", "content": ex.Input},
			map[string]string{"role": "assistant", "content": ex.Output},
		)
	}

	// Add actual user prompt
	messages = append(messages,
		map[string]string{"role": "user", "content": userPrompt},
	)

	return messages
}
```

### Step 2: Define Templates

```go
// Code Review Template
var CodeReviewTemplate = PromptTemplate{
	Name: "code-review",
	System: `You are an expert Go code reviewer. Analyze code for:
- Bugs and errors
- Security vulnerabilities
- Performance issues
- Best practices violations
- Code style and readability

Provide constructive feedback with specific suggestions.`,
	UserTemplate: `Review this {{.Language}} code:

{{.Code}}

Focus areas: {{.FocusAreas}}`,
	Examples: []Example{
		{
			Input: "Review this Go code:\n\nfunc divide(a, b int) int {\n    return a / b\n}",
			Output: `Issues found:
1. ❌ No error handling for division by zero
2. ⚠️  Should return error as second value

Suggested fix:
func divide(a, b int) (int, error) {
    if b == 0 {
        return 0, errors.New("division by zero")
    }
    return a / b, nil
}`,
		},
	},
}

// Code Explanation Template
var CodeExplainTemplate = PromptTemplate{
	Name: "code-explain",
	System: `You are a programming teacher. Explain code clearly and simply.
Use analogies when helpful. Break down complex concepts.`,
	UserTemplate: `Explain this {{.Language}} code to a {{.Level}} developer:

{{.Code}}`,
}

// Bug Fix Template
var BugFixTemplate = PromptTemplate{
	Name: "bug-fix",
	System: `You are a debugging expert. Find and fix bugs in code.
Explain what was wrong and why your fix works.`,
	UserTemplate: `This code has a bug:

{{.Code}}

Expected behavior: {{.Expected}}
Actual behavior: {{.Actual}}

Find and fix the bug.`,
}
```

### Step 3: Use Templates

```go
func main() {
	// Example 1: Code Review
	data := map[string]interface{}{
		"Language":   "Go",
		"Code":       "func process(data []int) { for i := 0; i < len(data); i++ { data[i] = data[i] * 2 } }",
		"FocusAreas": "performance, style",
	}

	prompt, _ := CodeReviewTemplate.Execute(data)
	messages := CodeReviewTemplate.GetMessages(prompt)

	fmt.Println("=== Code Review Prompt ===")
	for _, msg := range messages {
		fmt.Printf("[%s]: %s\n\n", msg["role"], msg["content"])
	}

	// Example 2: Code Explanation
	explainData := map[string]interface{}{
		"Language": "Go",
		"Level":    "beginner",
		"Code":     "ch := make(chan int, 10)",
	}

	explainPrompt, _ := CodeExplainTemplate.Execute(explainData)
	fmt.Println("=== Explanation Prompt ===")
	fmt.Println(explainPrompt)
}
```

---

## 🎨 Advanced Techniques

### 1. Role Prompting

```go
systemMsg := `You are a senior software architect at Google.
You have 15 years of experience building distributed systems.
You value simplicity, scalability, and maintainability.`
```

### 2. Constraint Setting

```go
userMsg := `Design a caching system with these constraints:
- Must handle 10,000 requests/second
- Maximum 100ms latency
- Budget: $500/month
- Team size: 2 developers
- Timeline: 2 weeks

Provide a realistic solution.`
```

### 3. Output Formatting

```go
userMsg := `Analyze this code and respond in JSON format:
{
  "bugs": ["list of bugs"],
  "security_issues": ["list of security issues"],
  "suggestions": ["list of improvements"],
  "severity": "low|medium|high"
}`
```

### 4. Iterative Refinement

```go
// First prompt
"Write a function to validate emails"

// Refine
"Add support for international domains"

// Refine more
"Add unit tests covering edge cases"

// Final refinement
"Add documentation and examples"
```

---

## 📊 Prompt Quality Checklist

```
┌──────────────────────────────────────────────────────────────────┐
│  Before sending a prompt, check:                                 │
│                                                                   │
│  ✅ Clear objective: What do you want?                           │
│  ✅ Sufficient context: Background information                   │
│  ✅ Specific constraints: Limitations, requirements              │
│  ✅ Desired format: How should output look?                      │
│  ✅ Examples provided: Show what you want (few-shot)             │
│  ✅ Role defined: Who should the AI be?                          │
│  ✅ Tone specified: Formal, casual, technical?                   │
└──────────────────────────────────────────────────────────────────┘
```

---

## 🎯 Challenges

### Challenge 1: Create Templates
Build templates for:
- Unit test generation
- API documentation
- Code refactoring
- Performance optimization

### Challenge 2: Template Library
Create a library of reusable templates with:
- Template registry
- Template validation
- Variable substitution
- Example management

### Challenge 3: A/B Testing
Test different prompt variations:
- Measure response quality
- Compare token usage
- Track success rate

---

## ✅ What You Learned

```
┌──────────────────────────────────────────────────────────────────┐
│  ✓ Prompt engineering fundamentals                               │
│  ✓ System, user, assistant roles                                 │
│  ✓ Few-shot learning with examples                               │
│  ✓ Chain-of-thought prompting                                    │
│  ✓ Prompt templates and patterns                                 │
│  ✓ Advanced techniques (role, constraints, formatting)           │
│  ✓ Quality checklist for prompts                                 │
└──────────────────────────────────────────────────────────────────┘
```

---

## 🚀 Next Steps

**Practice:**
1. Build the prompt template system
2. Create 5 custom templates
3. Test with different models
4. Measure quality improvements

**Next Tutorial:**
[Tutorial 3: Embeddings & Vectors →](03_EMBEDDINGS_VECTORS.md)

Learn how to convert text to vectors for semantic search!

---

## 📚 Resources

- [OpenAI Prompt Engineering Guide](https://platform.openai.com/docs/guides/prompt-engineering)
- [Prompt Engineering Guide by DAIR.AI](https://www.promptingguide.ai/)
- [Learn Prompting](https://learnprompting.org/)

---

**💡 Pro Tip**: Keep a library of your best prompts. Good prompts are reusable and save time!

