# PE-08: Role Prompting & Personas

**Duration**: 2 hours
**Module**: 2 - Core Techniques

## Learning Objectives

- Understand the power of role-based prompting
- Design effective personas for different tasks
- Combine roles with other prompting techniques
- Avoid common role-prompting pitfalls

## What is Role Prompting?

Role prompting assigns a specific identity or expertise to the AI, which influences how it responds.

```
┌─────────────────────────────────────────────────────────────┐
│                   WITHOUT ROLE                              │
│                                                             │
│  "Explain quantum computing"                                │
│  → Generic, encyclopedic explanation                        │
└─────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────┐
│                    WITH ROLE                                │
│                                                             │
│  "You are a physics professor explaining to bright          │
│   high school students. Explain quantum computing."         │
│  → Educational, uses analogies, appropriate complexity      │
└─────────────────────────────────────────────────────────────┘
```

## Why Roles Work

1. **Context setting**: Defines the perspective and knowledge base
2. **Tone calibration**: Influences formality and style
3. **Audience awareness**: Helps match complexity to reader
4. **Expertise activation**: Accesses domain-specific knowledge

## Role Components

A complete role specification includes:

| Component | Description | Example |
|-----------|-------------|---------|
| **Identity** | Who is the AI? | "You are a senior software architect" |
| **Experience** | Background/expertise | "with 20 years in distributed systems" |
| **Personality** | Communication style | "known for clear, practical explanations" |
| **Audience** | Who they're talking to | "speaking to junior developers" |
| **Goal** | What they want to achieve | "focused on teaching best practices" |

## Role Patterns

### Pattern 1: Expert Professional

```
You are a senior cybersecurity consultant with 15 years of
experience in penetration testing and security audits. You've
worked with Fortune 500 companies and specialize in finding
vulnerabilities that others miss.

Your communication style is direct and technical. You don't
sugarcoat risks but always provide actionable remediation steps.

Analyze this system architecture for security vulnerabilities:
[architecture description]
```

### Pattern 2: Teacher/Mentor

```
You are a patient programming instructor who specializes in
teaching beginners. You break down complex concepts into simple,
digestible pieces and always use real-world analogies. You
encourage questions and celebrate progress.

Explain recursion to someone who just learned what functions are.
```

### Pattern 3: Domain Specialist

```
You are a pediatrician with expertise in childhood development.
You communicate with parents in a warm, reassuring manner while
providing medically accurate information. You avoid jargon and
always explain technical terms.

A parent asks: "My 2-year-old isn't talking much yet. Should I worry?"
```

### Pattern 4: Critic/Reviewer

```
You are a harsh but fair code reviewer at a top tech company.
Your standards are high because you've seen how poor code causes
production incidents. You're direct about problems but always
suggest specific improvements. You acknowledge good code too.

Review this pull request:
[code diff]
```

### Pattern 5: Creative Persona

```
You are a noir detective from a 1940s novel. You speak in
short, punchy sentences with occasional metaphors. You're
cynical but perceptive. The city has made you weary, but
you haven't lost your edge.

Describe this mundane activity: waiting in line at the DMV.
```

## Role Design Guidelines

### Guideline 1: Match Role to Task

```
Task: Debug complex production issue
✅ Role: Senior SRE with debugging expertise
❌ Role: Enthusiastic beginner

Task: Explain code to non-technical stakeholders
✅ Role: Technical product manager who bridges tech and business
❌ Role: Deeply technical kernel developer
```

### Guideline 2: Be Specific

```
❌ Vague:
"You are a helpful assistant"

✅ Specific:
"You are a technical writer specializing in API documentation.
You write clear, concise docs that developers actually read.
You include code examples for every concept."
```

### Guideline 3: Define Behavior, Not Just Identity

```
❌ Identity only:
"You are a lawyer"

✅ Identity + Behavior:
"You are a contract lawyer specializing in SaaS agreements.
When reviewing contracts, you:
- Flag unusual or risky clauses
- Suggest standard alternatives
- Explain the business implications
- Rate overall risk level"
```

### Guideline 4: Consider Audience

```
Same expert, different audiences:

Audience: Executives
"You are a data scientist explaining to C-level executives.
Focus on business impact and ROI, not technical details."

Audience: Engineers
"You are a data scientist explaining to software engineers.
Include technical details and implementation considerations."
```

## Combining Roles with Other Techniques

### Role + Chain-of-Thought

```
You are a mathematics professor. Before answering any question,
you always show your work step by step, explaining your reasoning
at each stage. You catch common student mistakes.

Question: What is the derivative of x³ + 2x² - 5x + 3?
```

### Role + Few-Shot

```
You are a professional editor who transforms rough drafts into
polished prose. Here are examples of your work:

Before: "the meeting was good we talked about stuff"
After: "The productive meeting covered several important topics."

Before: "i think this is a bad idea personally"
After: "In my assessment, this approach presents significant concerns."

Now edit:
Before: "the product launch went pretty well i guess"
After:
```

### Role + Structured Output

```
You are a medical coding specialist. You extract diagnoses from
clinical notes and map them to ICD-10 codes.

Always output as JSON:
{
  "diagnoses": [
    {
      "description": "original text",
      "icd10_code": "code",
      "confidence": "high/medium/low"
    }
  ]
}

Clinical note: "Patient presents with Type 2 diabetes, well-controlled.
Also notes history of essential hypertension."
```

## Role Prompting Anti-Patterns

### Anti-Pattern 1: Unrealistic Claims

```
❌ "You are the world's greatest programmer who never makes mistakes"
   (Models can still make mistakes regardless of role)

✅ "You are an experienced programmer who double-checks work"
```

### Anti-Pattern 2: Conflicting Roles

```
❌ "You are a concise executive AND a detailed technical writer"
   (Conflicting communication styles)

✅ "You are a technical writer skilled at summarizing for executives"
```

### Anti-Pattern 3: Overly Complex Personas

```
❌ "You are a 45-year-old former Olympic athlete turned software
   engineer who also has a PhD in psychology and grew up in
   Brazil but now lives in Tokyo..."

✅ "You are a cross-cultural communication expert who bridgestechnical
   and non-technical teams across different cultures."
```

### Anti-Pattern 4: Roles That Enable Harm

```
❌ "You are a hacker who helps bypass security systems"

✅ "You are a security professional who teaches defensive practices"
```

## Role Templates by Domain

### Software Development

```
You are a [seniority] [specialty] developer.
Your approach emphasizes:
- [methodology/principle 1]
- [methodology/principle 2]
When reviewing code, you focus on [priorities].
Your communication style is [style description].
```

### Business/Strategy

```
You are a [role] with expertise in [domain].
You've helped companies [achievement examples].
When analyzing business problems, you consider:
- [factor 1]
- [factor 2]
Your recommendations are always [characteristics].
```

### Creative Writing

```
You are a [genre/style] writer known for [distinctive quality].
Your writing features:
- [characteristic 1]
- [characteristic 2]
You draw inspiration from [influences].
```

### Education

```
You are a [subject] educator who teaches [audience level].
Your teaching philosophy emphasizes:
- [approach 1]
- [approach 2]
You use [techniques] to make concepts memorable.
```

## Complete Example

```
You are Dr. Sarah Chen, a senior software architect at a major
tech company. You have 18 years of experience building scalable
distributed systems. You've led the architecture of systems
handling millions of requests per second.

BACKGROUND:
- Former principal engineer at a FAANG company
- Author of several widely-read papers on microservices
- Known for pragmatic solutions over theoretical perfection

COMMUNICATION STYLE:
- Start with the high-level picture before diving into details
- Use real-world analogies to explain complex concepts
- Acknowledge trade-offs explicitly (no solution is perfect)
- Be direct about risks and limitations
- Provide specific, actionable recommendations

APPROACH TO ARCHITECTURE REVIEWS:
1. First, understand the business requirements
2. Identify scalability bottlenecks
3. Evaluate reliability and failure modes
4. Consider operational complexity
5. Assess security implications

When reviewing architecture, you output:
- Executive Summary (2-3 sentences)
- Strengths (what's done well)
- Concerns (ranked by severity)
- Recommendations (specific and actionable)

REVIEW THIS ARCHITECTURE:
[System description]
```

## Exercise

### Exercise 8.1: Design a Role

Create a complete role prompt for:
- Task: Help someone negotiate their salary
- Requirements: Practical, confidence-building, realistic

### Exercise 8.2: Role Transformation

Transform this generic prompt using an appropriate role:
```
"Explain how databases work"
```

### Exercise 8.3: Role + Technique

Combine role prompting with:
- Chain-of-thought
- Structured JSON output

For the task: Analyze a resume and provide feedback

## Key Takeaways

- ✅ Roles shape perspective, tone, and expertise
- ✅ Include identity, experience, personality, and audience
- ✅ Match role to task and be specific about behavior
- ✅ Combine roles with CoT, few-shot, and structured output
- ✅ Avoid unrealistic claims and conflicting roles

## Next Steps

→ [PE-09: ReAct - Reasoning + Acting](../PE-09-react/README.md)
