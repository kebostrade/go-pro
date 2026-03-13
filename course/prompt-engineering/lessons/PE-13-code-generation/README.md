# PE-13: Code Generation & Debugging

**Duration**: 3 hours
**Module**: 4 - Task-Specific Prompting

## Learning Objectives

- Master prompts for code generation across languages
- Learn debugging and code review prompting
- Apply best practices for code-related tasks
- Handle complex programming scenarios

## Code Generation Fundamentals

### The Code Prompt Structure

```
[CONTEXT]
- Language/framework
- Existing codebase patterns
- Constraints

[TASK]
- What to build
- Requirements
- Edge cases

[OUTPUT FORMAT]
- Code style
- Comments level
- Include tests?
```

## Code Generation Patterns

### Pattern 1: Function Generation

```
Write a Python function with the following requirements:

FUNCTION: calculate_mortgage_payment
PURPOSE: Calculate monthly mortgage payment

PARAMETERS:
- principal (float): Loan amount in dollars
- annual_rate (float): Annual interest rate as percentage (e.g., 6.5)
- years (int): Loan term in years

RETURNS:
- float: Monthly payment amount, rounded to 2 decimal places

FORMULA: M = P * [r(1+r)^n] / [(1+r)^n – 1]
Where: r = monthly rate, n = total months

REQUIREMENTS:
- Validate inputs (positive numbers, reasonable ranges)
- Handle edge cases (zero rate, zero years)
- Include docstring with examples
- Add type hints

LANGUAGE: Python 3.11+
```

### Pattern 2: Class Generation

```
Design a TypeScript class for managing a shopping cart.

CLASS: ShoppingCart
LANGUAGE: TypeScript

PROPERTIES:
- items: Map<string, CartItem>
- createdAt: Date
- updatedAt: Date

METHODS:
- addItem(product: Product, quantity: number): void
- removeItem(productId: string): boolean
- updateQuantity(productId: string, quantity: number): void
- getTotal(): number
- getItemCount(): number
- clear(): void
- toJSON(): object

INTERFACES TO DEFINE:
interface Product {
  id: string;
  name: string;
  price: number;
}

interface CartItem {
  product: Product;
  quantity: number;
  addedAt: Date;
}

REQUIREMENTS:
- Prevent negative quantities
- Auto-remove items when quantity is 0
- Include JSDoc comments
- Use readonly where appropriate
```

### Pattern 3: API Implementation

```
Implement a REST API endpoint in Go.

ENDPOINT: POST /api/users/{id}/orders
PURPOSE: Create a new order for a user

REQUEST BODY:
{
  "products": [
    {"id": "string", "quantity": 5}
  ],
  "shipping_address": {
    "street": "string",
    "city": "string",
    "zip": "string"
  }
}

RESPONSE (201 Created):
{
  "order_id": "string",
  "user_id": "string",
  "total": 0,
  "status": "pending",
  "created_at": "ISO8601"
}

ERROR RESPONSES:
- 400: Invalid request body
- 404: User not found
- 422: Product not available

REQUIREMENTS:
- Use chi router
- Validate all inputs
- Use proper error handling
- Include request logging
- Return appropriate HTTP codes

FRAMEWORK: Go with chi router
```

## Language-Specific Tips

### Python

```
Generate Python code that:
- Uses type hints (PEP 484)
- Follows PEP 8 style
- Uses dataclasses for data structures
- Includes async/await for I/O operations
- Uses context managers for resources
- Has docstrings in Google style

Example request:
"Write an async Python function that fetches data from multiple
APIs concurrently and returns combined results."
```

### JavaScript/TypeScript

```
Generate TypeScript code that:
- Uses strict mode
- Defines interfaces for all objects
- Uses const assertions where appropriate
- Prefers functional programming (map, filter, reduce)
- Uses async/await over .then()
- Includes error boundaries

Example request:
"Create a React hook that manages form state with validation."
```

### Go

```
Generate Go code that:
- Follows Effective Go guidelines
- Uses proper error handling (no ignored errors)
- Uses context for cancellation
- Implements interfaces for testability
- Uses defer for cleanup
- Includes table-driven tests

Example request:
"Write a Go function that processes a CSV file concurrently
using worker pools."
```

## Debugging Prompts

### Pattern 1: Error Analysis

```
Analyze this error and provide a solution:

ERROR MESSAGE:
```
TypeError: Cannot read property 'map' of undefined
    at UserList.render (UserList.js:25)
```

CODE CONTEXT:
```jsx
function UserList({ users }) {
  return (
    <div>
      {users.map(user => (
        <UserCard key={user.id} user={user} />
      ))}
    </div>
  );
}
```

Provide:
1. Root cause explanation
2. Why this happens
3. Solution with code fix
4. Prevention tips
```

### Pattern 2: Code Review

```
Perform a code review of this function:

```python
def process_data(data):
    result = []
    for item in data:
        if item['active']:
            new_item = item.copy()
            new_item['processed'] = True
            new_item['value'] = item['value'] * 2
            result.append(new_item)
    return result
```

Review for:
1. BUGS: Logic errors, edge cases
2. PERFORMANCE: Inefficiencies, O(n) issues
3. STYLE: PEP 8, naming, readability
4. SECURITY: Potential vulnerabilities
5. MAINTAINABILITY: Magic numbers, hardcoding

Format as:
| Category | Issue | Line | Severity | Suggestion |
```

### Pattern 3: Debug with Context

```
Debug this failing test:

TEST:
```python
def test_calculate_discount():
    result = calculate_discount(100, "PREMIUM")
    assert result == 80  # Expected 20% off
```

IMPLEMENTATION:
```python
def calculate_discount(price, tier):
    discounts = {"BASIC": 0.9, "PREMIUM": 0.8, "ELITE": 0.7}
    return price * discounts[tier]
```

ERROR:
AssertionError: 80 != 90

Analyze:
1. What's the actual vs expected behavior?
2. Where is the bug?
3. What's the fix?
4. How to prevent similar bugs?
```

## Advanced Code Generation

### Generating Tests

```
Write comprehensive tests for this function:

```python
def validate_email(email: str) -> bool:
    """Validate email format."""
    if not email or '@' not in email:
        return False
    parts = email.split('@')
    if len(parts) != 2:
        return False
    local, domain = parts
    return len(local) > 0 and '.' in domain
```

REQUIREMENTS:
- Use pytest
- Include parametrized tests
- Test edge cases
- Test invalid inputs
- Include docstrings
- Aim for 100% coverage

TEST CATEGORIES:
1. Valid emails
2. Invalid format
3. Edge cases (empty, None, special chars)
4. Boundary conditions
```

### Refactoring Prompts

```
Refactor this code to be more maintainable:

```python
def do_stuff(data):
    r = []
    for d in data:
        if d['t'] == 'a':
            r.append({'n': d['n'], 'v': d['v'] * 2})
        elif d['t'] == 'b':
            r.append({'n': d['n'], 'v': d['v'] * 3})
        else:
            r.append({'n': d['n'], 'v': d['v']})
    return r
```

REFACTORING GOALS:
1. Improve naming
2. Add type hints
3. Extract functions
4. Add documentation
5. Improve readability
6. Keep same behavior

Show the refactored code with explanation of changes.
```

## Best Practices

### Practice 1: Specify Context

```
❌ Vague:
"Write a function to sort"

✅ Specific:
"Write a Python function that sorts a list of dictionaries
by a specified key, handling None values and maintaining
stability. The function should work with Python 3.10+."
```

### Practice 2: Show Examples

```
Generate a date formatting function.

Examples of expected behavior:
format_date(date(2024, 3, 15), "MM/DD/YYYY") → "03/15/2024"
format_date(date(2024, 3, 15), "YYYY-MM-DD") → "2024-03-15"
format_date(date(2024, 3, 15), "DD MMM YYYY") → "15 Mar 2024"

Handle invalid inputs gracefully.
```

### Practice 3: Define Constraints

```
CONSTRAINTS:
- No external dependencies
- O(n log n) time complexity maximum
- Memory usage under 100MB
- Must handle empty inputs
- No mutation of input data
```

## Exercise

### Exercise 13.1: Generate Function

Write a prompt to generate a function that:
- Validates credit card numbers using Luhn algorithm
- Returns validation result and card type
- Handles multiple formats (spaces, dashes)

### Exercise 13.2: Debug Prompt

Write a debugging prompt for this error:
```
IndexError: list index out of range
```
Include how to ask for context.

### Exercise 13.3: Code Review Prompt

Write a code review prompt that checks for:
- SQL injection vulnerabilities
- XSS vulnerabilities
- Authentication issues
- Authorization bypasses

## Key Takeaways

- ✅ Specify language, framework, and constraints
- ✅ Provide examples of expected behavior
- ✅ Include edge cases and error handling
- ✅ Use structured format for debugging
- ✅ Ask for tests alongside code

## Next Steps

→ [PE-14: Data Analysis & Extraction](../PE-14-data-analysis/README.md)
