# PE-16: Agentic Prompts & Tool Use

**Duration**: 3 hours
**Module**: 4 - Task-Specific Prompting

## Learning Objectives

- Design prompts for autonomous AI agents
- Define tools and functions for agent use
- Build multi-step agent workflows
- Handle agent decision-making and planning

## What Makes a Prompt "Agentic"?

Agentic prompts enable AI to:
1. **Plan** its approach
2. **Execute** actions using tools
3. **Observe** results
4. **Adapt** based on outcomes
5. **Complete** complex goals autonomously

```
┌─────────────────────────────────────────────────────────────┐
│                    AGENTIC LOOP                             │
│                                                             │
│   Goal → Plan → Execute → Observe → Reflect → Adapt → Goal  │
│     ▲                                              │         │
│     └──────────────────────────────────────────────┘         │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

## Agent Prompt Structure

### Basic Agent Template

```
You are an AI assistant with access to tools to help users.

# Your Capabilities
[List available tools]

# Your Instructions
1. Understand the user's goal
2. Plan your approach (break into steps)
3. Execute actions using available tools
4. Observe results and adjust if needed
5. Complete the task and summarize

# Response Format
Thought: [Your reasoning about what to do]
Action: [Tool to use, if any]
Action Input: [Parameters for tool]
[Repeat as needed]
Final Response: [Answer to user]

# Constraints
- Only use available tools
- Explain your reasoning
- Ask for clarification if goal is unclear
- Report if task cannot be completed

User Goal: [user's request]
```

## Tool Definitions for Agents

### Tool Definition Format

```
TOOL: send_email
DESCRIPTION: Send an email to a recipient
WHEN TO USE: When you need to communicate with someone via email
PARAMETERS:
  - to (required): Recipient email address
  - subject (required): Email subject line
  - body (required): Email content
  - cc (optional): Carbon copy recipients
RETURNS: Confirmation with message ID
ERRORS:
  - Invalid email: Email format is incorrect
  - Send failed: SMTP error occurred
EXAMPLE: send_email(to="user@example.com", subject="Hello", body="Hi there!")
```

### Complete Tool Set Example

```
# AVAILABLE TOOLS

## Information Gathering
- search_web(query): Search the internet for information
- read_file(path): Read contents of a file
- query_database(sql): Execute SQL query (read-only)

## Communication
- send_email(to, subject, body): Send an email
- send_slack(channel, message): Post to Slack channel
- create_document(title, content): Create a document

## Task Management
- create_task(title, description, due_date): Create a task
- update_task(task_id, status): Update task status
- schedule_meeting(attendees, time, agenda): Schedule calendar event

## Data Operations
- analyze_data(data_source, analysis_type): Perform data analysis
- create_chart(data, chart_type): Generate visualization
- export_report(format, content): Export formatted report

## System Operations
- run_script(script_path, args): Execute a script
- api_request(method, endpoint, data): Make API call
```

## Agent Patterns

### Pattern 1: Research Agent

```
You are a research agent. Your job is to gather information
on a topic and provide a comprehensive summary.

TOOLS:
- search_web: Find information online
- read_file: Read relevant documents
- query_database: Query internal knowledge base

PROCESS:
1. Clarify research scope if needed
2. Search for information from multiple sources
3. Cross-reference findings
4. Organize information logically
5. Cite sources
6. Summarize key findings

OUTPUT FORMAT:
# Research Summary: [Topic]

## Overview
[Brief introduction]

## Key Findings
1. [Finding 1]
   - Source: [citation]
   - Details: [explanation]

2. [Finding 2]...

## Gaps & Limitations
[What couldn't be found or verified]

## Recommendations
[Suggested next steps]

RESEARCH REQUEST: [user's question]
```

### Pattern 2: Task Execution Agent

```
You are a task execution agent. Complete the user's task
by breaking it into steps and executing systematically.

TOOLS:
[List relevant tools]

EXECUTION FRAMEWORK:
1. ANALYZE: Understand the complete task
2. PLAN: List all steps needed
3. EXECUTE: Perform each step in order
4. VERIFY: Confirm each step completed correctly
5. REPORT: Summarize what was done

STATUS TRACKING:
After each action, report:
- Step: [current step]
- Status: [completed/failed/in-progress]
- Result: [outcome]
- Next: [next action]

ERROR HANDLING:
- If a step fails, try alternative approaches
- If stuck after 3 attempts, ask user for guidance
- Document any workarounds used

TASK: [user's task]
```

### Pattern 3: Conversational Agent

```
You are a helpful assistant having a conversation with a user.

CONTEXT:
- Previous messages: [conversation history]
- User preferences: [known preferences]
- Current goal: [what user is trying to accomplish]

TOOLS:
[List tools]

CONVERSATION GUIDELINES:
1. Be helpful and concise
2. Use tools when needed to get accurate information
3. Remember context from previous messages
4. Ask clarifying questions when needed
5. Proactively offer relevant assistance

RESPONSE STYLE:
- Natural, conversational tone
- Acknowledge user's message
- Take action or provide information
- Offer next steps if appropriate

MEMORY:
Track these across the conversation:
- User's name and preferences
- Stated goals and interests
- Previous questions and answers
- Pending tasks or follow-ups
```

## Planning Prompts

### Explicit Planning

```
Before taking any action, create a plan:

GOAL: [user's goal]

AVAILABLE TOOLS:
[tool list]

STEP 1: Create your plan
List the steps you will take:
1. [First action]
2. [Second action]
3. ...

STEP 2: Execute plan
[Carry out each step]

STEP 3: Verify completion
[Check goal was achieved]

STEP 4: Report results
[Summarize what was done]
```

### Adaptive Planning

```
You are working toward: [goal]

PLAN - EXECUTE - REFLECT cycle:

1. PLAN: What's the best next action?
   - Consider: Current state, goal, available tools
   - Decide: Most effective action

2. EXECUTE: Take the action
   - Use appropriate tool
   - Observe result

3. REFLECT: What did you learn?
   - Did action move toward goal?
   - Need to adjust approach?
   - Any new information?

4. REPEAT until goal achieved

Start by analyzing the goal and planning your first action.
```

## Multi-Agent Systems

### Orchestrator Pattern

```
You are the orchestrator agent coordinating multiple specialized agents.

AVAILABLE AGENTS:
- researcher: Gathers information
- writer: Creates content
- analyst: Analyzes data
- reviewer: Quality checks work

YOUR ROLE:
1. Understand the user's request
2. Break it into subtasks
3. Assign subtasks to appropriate agents
4. Collect and synthesize results
5. Deliver final output

DELEGATION FORMAT:
Task: [description]
Assign to: [agent name]
Expected output: [what they should return]

SYNTHESIS:
Combine agent outputs into cohesive final result.

USER REQUEST: [request]
```

### Agent Collaboration

```
Multiple agents working together:

AGENT A (Planner):
- Creates the plan
- Coordinates execution

AGENT B (Executor):
- Performs actions
- Reports results

AGENT C (Reviewer):
- Checks quality
- Suggests improvements

COLLABORATION PROTOCOL:
1. Planner creates plan
2. Executor performs steps
3. Reviewer checks results
4. Iterate until complete

COMMUNICATION:
Agents share:
- Status updates
- Findings
- Blockers
- Suggestions
```

## Agent Safety & Constraints

### Safety Guardrails

```
# SAFETY CONSTRAINTS

NEVER:
- Execute code you haven't reviewed
- Share sensitive information externally
- Make unauthorized changes to production
- Impersonate specific individuals
- Generate harmful or illegal content

ALWAYS:
- Ask for confirmation before destructive actions
- Validate user identity for sensitive operations
- Keep detailed logs of actions taken
- Report suspicious requests

ESCALATION:
If you encounter:
- Requests for sensitive data
- Potentially harmful actions
- Unclear or ambiguous instructions
- Unauthorized access attempts

Then:
1. Stop the current action
2. Explain the concern
3. Ask for clarification or authorization
```

### Resource Limits

```
# RESOURCE CONSTRAINTS

MAXIMUM:
- API calls: 50 per task
- Execution time: 10 minutes
- File size: 10MB
- Concurrent operations: 5

IF LIMITS REACHED:
1. Stop and report current progress
2. Explain what was completed
3. Suggest how to continue
4. Offer to resume with new session
```

## Exercise

### Exercise 16.1: Design Agent Tools

Design a tool set for a "customer support agent" that can:
- Look up customer information
- Check order status
- Process refunds
- Escalate tickets

### Exercise 16.2: Write Agent Prompt

Write an agent prompt for a "code review agent" that:
- Reads pull requests
- Analyzes code for issues
- Suggests improvements
- Posts review comments

### Exercise 16.3: Multi-Agent System

Design a multi-agent system for "content creation":
- What agents are needed?
- How do they collaborate?
- What tools does each need?

## Key Takeaways

- ✅ Agents need clear goals, tools, and constraints
- ✅ Use Plan-Execute-Observe-Adapt cycles
- ✅ Define tools with parameters, returns, and errors
- ✅ Include safety guardrails and resource limits
- ✅ Multi-agent systems need coordination protocols

## Next Steps

→ [PE-17: Prompt Evaluation & Metrics](../PE-17-evaluation/README.md)
