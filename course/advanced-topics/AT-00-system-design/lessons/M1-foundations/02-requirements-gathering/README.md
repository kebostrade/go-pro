# SD-02: Requirements Gathering & Analysis

Learn how to gather, analyze, and document system requirements effectively.

## Overview

Before designing any system, you must understand what the system needs to do. This lesson covers techniques for gathering and analyzing requirements.

## Learning Objectives

- Conduct effective requirements interviews
- Create user stories and use cases
- Identify functional and non-functional requirements
- Prioritize features

## Types of Requirements

### Functional Requirements

What the system must do:

- User authentication and authorization
- CRUD operations for entities
- Search and filtering capabilities
- Notifications and alerts

### Non-Functional Requirements

How the system should perform:

- **Performance**: Response time, throughput
- **Scalability**: Users, data volume
- **Availability**: Uptime, SLA
- **Security**: Authentication, encryption
- **Reliability**: Error handling, recovery

## Requirements Gathering Techniques

### 1. Stakeholder Interviews

```
Questions to ask:
- What problem are you trying to solve?
- Who are the end users?
- What are the current pain points?
- What success looks like?
- What constraints exist?
```

### 2. User Stories

Format: As a [user], I want [feature], so that [benefit]

```
Examples:
- As a user, I want to reset my password, so that I can recover access
- As an admin, I want to view analytics, so that I can make data-driven decisions
- As a developer, I want API documentation, so that I can integrate easily
```

### 3. Use Cases

```
Use Case: User Login
Actor: Registered User
Preconditions: User has valid credentials
Main Flow:
  1. User enters email/password
  2. System validates credentials
  3. System creates session
  4. User is redirected to dashboard
Postconditions: User is authenticated
```

## Requirement Prioritization

### MoSCoW Method

| Priority | Description |
|----------|-------------|
| Must Have | Critical for release |
| Should Have | Important but not critical |
| Could Have | Nice to have |
| Won't Have | Not in scope |

### RICE Scoring

```
RICE Score = (Reach × Impact × Confidence) / Effort

- Reach: How many users per quarter?
- Impact: How much does it help users?
- Confidence: How sure are we?
- Effort: Person-weeks needed
```

## Requirements Document Template

```markdown
# Requirement: User Authentication

## Type
Functional Requirement

## Description
Users must be able to authenticate using email/password

## Acceptance Criteria
- [ ] User can register with email
- [ ] User can login with credentials
- [ ] Invalid credentials show error
- [ ] Session expires after 24 hours

## Priority
Must Have

## Dependencies
- Email service for verification
- Database for user storage
```

## Examples

See `examples/` directory for:
- `user_stories.md` - Example user stories
- `use_cases.go` - Use case implementation
- `requirements_template.md` - Document template

## Exercises

See `exercises/problems.md` for hands-on practice.

## Quiz

Test your knowledge with `quiz.md`.

## Summary

- Gather requirements from multiple sources
- Document both functional and non-functional requirements
- Prioritize based on business value
- Validate requirements with stakeholders

## Next Steps

Continue to [SD-03: Capacity Planning](03-capacity-planning/README.md)
