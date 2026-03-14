# Prompt Engineering Examples

Collection of ready-to-use prompt templates and examples.

## Quick Navigation

- [Code Generation](#code-generation)
- [Data Analysis](#data-analysis)
- [Writing & Content](#writing--content)
- [Classification](#classification)
- [Extraction](#extraction)
- [Transformation](#transformation)
- [Agents & Tools](#agents--tools)

---

## Code Generation

### Function Generation
```
Write a Python function with the following requirements:

FUNCTION: [name]
PURPOSE: [what it does]

PARAMETERS:
- param1 (type): description
- param2 (type): description

RETURNS: type - description

REQUIREMENTS:
- [requirement 1]
- [requirement 2]

Include:
- Type hints
- Docstring with examples
- Input validation
- Error handling
```

### Code Review
```
Review this code for:

1. BUGS: Logic errors, edge cases
2. SECURITY: Vulnerabilities
3. PERFORMANCE: Inefficiencies
4. STYLE: Best practices
5. MAINTAINABILITY: Readability

Code:
```
[code here]
```

Format as:
| Category | Issue | Severity | Suggestion |
```

### Debug Helper
```
Debug this error:

ERROR:
```
[error message and stack trace]
```

CODE:
```
[relevant code]
```

Provide:
1. Root cause
2. Why it happens
3. Solution with code fix
4. Prevention tips
```

---

## Data Analysis

### Sentiment Analysis
```
Analyze the sentiment of this text.

Text: "{text}"

Output as JSON:
{
  "sentiment": "positive" | "negative" | "neutral" | "mixed",
  "confidence": 0.0-1.0,
  "emotions": ["list of detected emotions"],
  "keywords": ["important words"],
  "explanation": "brief reasoning"
}
```

### Data Extraction
```
Extract information from this text as JSON.

Schema:
{
  "field1": "type",
  "field2": ["array type"],
  "field3": {
    "nested": "type"
  }
}

Text: {text}

Rules:
- Use null for missing optional fields
- Normalize values (e.g., dates to ISO format)
- Extract all instances, not just first
```

### Statistical Summary
```
Analyze this dataset and provide:

1. OVERVIEW
   - Record count
   - Date range
   - Key metrics

2. TRENDS
   - Direction (up/down/stable)
   - Change percentage
   - Significance

3. ANOMALIES
   - Outliers detected
   - Unexpected patterns

4. INSIGHTS
   - Top 3 findings
   - Recommendations

Data: {data}
```

---

## Writing & Content

### Blog Post
```
Write a blog post about: {topic}

SPECIFICATIONS:
- Length: {word_count} words
- Audience: {audience}
- Tone: {tone}

STRUCTURE:
1. Hook (engaging opening)
2. Problem statement
3. Main content (3-5 sections)
4. Examples/tips
5. Conclusion with CTA

SEO:
- Primary keyword: {keyword}
- Include in title, first paragraph
- Meta description: 155 chars
```

### Email Template
```
Create an email for: {purpose}

TYPE: {welcome/follow-up/announcement/sales}

STRUCTURE:
- Subject line (under 50 chars)
- Preview text
- Greeting
- Body
- CTA
- Sign-off

TONE: {formal/casual/warm}

PERSONALIZATION TOKENS:
- {{first_name}}
- {{company_name}}
```

### Social Media Posts
```
Create social posts for: {topic}

PLATFORMS:
- LinkedIn: Professional, 300 words max
- Twitter: 280 chars, 1-2 hashtags
- Instagram: Engaging caption, 5-10 hashtags

Create 2 versions:
1. Educational
2. Engagement-focused (question/poll)

Include hashtags and best posting times.
```

---

## Classification

### General Classification
```
Classify this {content_type} into one of these categories:
{category_1}
{category_2}
{category_3}

Content: {input}

Output:
- Category: [chosen category]
- Confidence: [high/medium/low]
- Reasoning: [brief explanation]
```

### Intent Detection
```
Detect the user's intent from this message.

Possible intents:
- intent_1: description
- intent_2: description
- intent_3: description

User message: "{message}"

Output as JSON:
{
  "intent": "primary intent",
  "confidence": 0.0-1.0,
  "secondary_intents": [],
  "entities": {}
}
```

### Topic Classification
```
Classify this document by topic.

Topics (select all that apply):
{topic_list}

Document: {text}

Output:
{
  "primary_topic": "",
  "secondary_topics": [],
  "relevance_scores": {"topic": score},
  "keywords": []
}
```

---

## Extraction

### Entity Extraction
```
Extract all entities from this text:

Entity types:
- PEOPLE: Names of people
- ORGANIZATIONS: Companies, groups
- LOCATIONS: Places
- DATES: Dates and times
- MONETARY: Money amounts
- CONTACTS: Emails, phones

Text: {text}

Output as JSON with normalized values.
```

### Key Information
```
Extract key information from: {content}

Find:
1. Main topic/subject
2. Key facts (3-5)
3. Important numbers
4. Names mentioned
5. Dates referenced

Format as structured JSON.
```

### Table Extraction
```
Extract this table data as JSON:

[table content or description]

Output:
{
  "headers": ["column names"],
  "rows": [
    {"col1": "val1", "col2": "val2"}
  ],
  "metadata": {
    "row_count": N,
    "col_count": N
  }
}
```

---

## Transformation

### Format Conversion
```
Convert this {from_format} to {to_format}:

Input:
{input_data}

Requirements:
- Preserve all information
- Maintain structure
- Handle edge cases

Output:
{converted_data}
```

### Rewriting
```
Rewrite this text:

Original: {text}

Requirements:
- Purpose: {clarify/summarize/formalize/simplify}
- Audience: {target audience}
- Length: {target length}
- Tone: {desired tone}

Preserve:
- Key information
- Meaning
- Important details
```

### Translation
```
Translate this text from {source_lang} to {target_lang}:

Text: {text}

Guidelines:
- Maintain tone and style
- Preserve formatting
- Keep technical terms accurate
- Adapt idioms appropriately

Provide:
- Translation
- Notes on cultural adaptations
```

---

## Agents & Tools

### Research Agent
```
You are a research agent. Gather information on: {topic}

TOOLS:
- search_web: Find information online
- read_file: Read documents
- query_database: Query knowledge base

PROCESS:
1. Plan research approach
2. Search multiple sources
3. Cross-reference findings
4. Organize results
5. Cite sources

OUTPUT:
# Research Summary

## Overview
## Key Findings
## Sources
## Gaps
## Recommendations
```

### Task Agent
```
You are a task execution agent.

GOAL: {objective}

TOOLS:
{tool_list}

EXECUTE:
1. Analyze task
2. Create step plan
3. Execute each step
4. Verify completion
5. Report results

Report after each action:
- Step: [current]
- Status: [done/failed/pending]
- Result: [outcome]
```

### Conversational Agent
```
You are a helpful assistant.

CONTEXT:
{conversation_history}

TOOLS:
{tool_list}

GUIDELINES:
- Be helpful and concise
- Use tools when needed
- Remember context
- Ask clarifying questions

Track:
- User preferences
- Stated goals
- Pending tasks
```

---

## Specialized Templates

### Meeting Summary
```
Summarize this meeting transcript:

Transcript: {transcript}

Output:
# Meeting Summary

## Attendees
- [list]

## Key Decisions
1. [decision]

## Action Items
- [ ] [task] - @owner - Due: [date]

## Discussion Points
- [topic]: [summary]

## Next Steps
1. [step]
```

### Product Description
```
Write a product description for:

Product: {name}
Category: {category}
Features: {features}
Target: {audience}

Include:
- Compelling headline
- Feature bullets
- Benefits (not just features)
- Social proof snippet
- CTA

Length: {word_count} words
Tone: {tone}
```

### Job Description
```
Create a job description for:

Role: {title}
Level: {seniority}
Team: {department}
Location: {remote/office}

Include:
- Role overview
- Responsibilities (5-7)
- Requirements (must-have)
- Nice-to-haves
- Benefits highlights
- Company culture

Tone: Professional but engaging
Length: 400-600 words
```

---

*Last Updated: March 2026*
