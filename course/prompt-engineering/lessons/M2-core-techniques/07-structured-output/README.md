# PE-07: Structured Output & JSON

**Duration**: 2 hours
**Module**: 2 - Core Techniques

## Learning Objectives

- Generate reliable structured output from LLMs
- Master JSON, XML, and other formats
- Handle parsing and validation
- Implement output constraints effectively

## Why Structured Output Matters

Structured output enables:
- **Programmatic processing**: Parse and use output in code
- **Consistency**: Same format every time
- **Integration**: Connect to APIs, databases, UIs
- **Validation**: Check output meets requirements

```
┌─────────────────────────────────────────────────────────────┐
│                 UNSTRUCTURED OUTPUT                         │
│                                                             │
│  "The customer seemed pretty happy I think. They mentioned  │
│   liking the product but had some issues with shipping.     │
│   Overall probably positive sentiment."                     │
│                                                             │
│  ❌ Hard to parse programmatically                          │
└─────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────┐
│                  STRUCTURED OUTPUT                          │
│                                                             │
│  {                                                          │
│    "sentiment": "positive",                                 │
│    "confidence": 0.85,                                      │
│    "topics": ["product quality", "shipping"],               │
│    "summary": "Customer likes product, shipping issues"     │
│  }                                                          │
│                                                             │
│  ✅ Easy to parse, consistent, processable                  │
└─────────────────────────────────────────────────────────────┘
```

## JSON Output

### Basic JSON Structure

```
Extract information from this product review as JSON:

Review: "I bought the BlueWidget Pro last week. The quality is
excellent and shipping was fast. Price was $49.99. Would recommend!"

Output format:
{
  "product_name": "",
  "sentiment": "",
  "price": null,
  "aspects_mentioned": [],
  "would_recommend": false
}
```

### Schema Specification

Be explicit about the schema:

```
Analyze this text and return a JSON object with this EXACT structure:

{
  "entities": {
    "people": ["array of person names"],
    "organizations": ["array of company/organization names"],
    "locations": ["array of place names"]
  },
  "relationships": [
    {
      "person": "name",
      "organization": "company",
      "role": "their position"
    }
  ],
  "summary": "one sentence summary",
  "confidence_score": 0.0-1.0
}

Text: "Elon Musk, CEO of Tesla, announced plans to expand
production at the Austin, Texas facility."
```

### Handling Arrays and Nested Objects

```
Extract recipe information:

{
  "name": "recipe name",
  "prep_time_minutes": number,
  "cook_time_minutes": number,
  "servings": number,
  "ingredients": [
    {
      "item": "ingredient name",
      "amount": "quantity",
      "unit": "measurement unit"
    }
  ],
  "steps": [
    "step 1",
    "step 2"
  ],
  "tags": ["easy", "vegetarian", etc.]
}

Recipe: [paste recipe text here]
```

## Techniques for Reliable JSON

### Technique 1: Explicit Schema

```
Return ONLY valid JSON matching this schema:

{
  "type": "object",
  "required": ["name", "age"],
  "properties": {
    "name": {"type": "string"},
    "age": {"type": "integer", "minimum": 0},
    "email": {"type": "string", "format": "email"}
  }
}

No markdown, no explanations, only JSON.
```

### Technique 2: Example-Driven

```
Parse the contact information:

Example 1:
Input: "Call John at 555-1234 or email john@email.com"
Output: {"name": "John", "phone": "555-1234", "email": "john@email.com"}

Example 2:
Input: "Reach Sarah via sarah@work.org"
Output: {"name": "Sarah", "phone": null, "email": "sarah@work.org"}

Now parse:
Input: "Contact Mike Smith: 555-9999, mike.smith@company.net"
Output:
```

### Technique 3: Preamble Warning

```
You must respond with ONLY valid JSON. Do not include:
- Markdown code blocks
- Explanations before or after
- Any text outside the JSON object

Your entire response will be parsed as JSON.
```

### Technique 4: Post-Processing Instructions

```
After generating the JSON:
1. Verify all required fields are present
2. Ensure all strings are properly escaped
3. Confirm the JSON is valid before responding
4. If any field is uncertain, use null

Output ONLY the JSON object, nothing else.
```

## Handling Edge Cases

### Missing Data

```
If information is not available, use null:
{
  "name": "John",
  "phone": null,
  "email": null,
  "address": null
}
```

### Arrays with Unknown Length

```
"ingredients": ["list", "all", "ingredients", "mentioned"]

Include ALL items found, not just a subset.
```

### Special Characters

```
Escape special characters in strings:
- Quotes: \" inside strings
- Backslashes: \\
- Newlines: \n

Example: {"review": "He said \"great product!\""}
```

## Other Structured Formats

### XML Output

```
Extract contact information as XML:

<contact>
  <name>Full Name</name>
  <phone type="mobile">555-1234</phone>
  <email>email@example.com</email>
  <address>
    <street>Street Address</street>
    <city>City</city>
    <state>State</state>
    <zip>12345</zip>
  </address>
</contact>

Input: "Jane Doe lives at 123 Main St, Boston, MA 02101.
Contact at jane@email.com or 555-9999."
```

### Markdown Tables

```
Convert this data to a markdown table:

| Product | Price | Rating | In Stock |
|---------|-------|--------|----------|
| Widget  | $10   | 4.5    | Yes      |
| Gadget  | $25   | 4.0    | No       |

Data: [your data here]
```

### CSV Format

```
Convert to CSV format (no markdown, just raw CSV):

name,price,quantity
Widget,10.99,50
Gadget,25.00,30

Data: [your data here]
```

### YAML Output

```
Extract configuration as YAML:

server:
  host: localhost
  port: 8080
database:
  name: mydb
  user: admin
  pool_size: 10

Settings: [describe settings here]
```

## Validation Strategies

### Strategy 1: Field Validation

```
Output must satisfy these constraints:
- "name": non-empty string, max 100 chars
- "age": integer between 0 and 150
- "email": valid email format or null
- "phone": format "XXX-XXX-XXXX" or null

If a value doesn't match, use null.
```

### Strategy 2: Enum Constraints

```
The "status" field must be exactly one of:
- "pending"
- "approved"
- "rejected"
- "cancelled"

No other values are allowed.
```

### Strategy 3: Required Fields

```
Required fields (must not be null):
- id
- name
- created_at

Optional fields (can be null):
- description
- tags
- metadata
```

## Programmatic Parsing

### Python Example

```python
import json
import anthropic

client = anthropic.Anthropic()

prompt = """
Extract person information as JSON:
{"name": "", "age": null, "email": null}

Text: "John is 30 years old. Contact: john@email.com"

Output ONLY valid JSON, nothing else.
"""

response = client.messages.create(
    model="claude-sonnet-4-6-20250514",
    max_tokens=1024,
    messages=[{"role": "user", "content": prompt}]
)

# Parse the JSON response
try:
    data = json.loads(response.content[0].text)
    print(f"Name: {data['name']}")
    print(f"Age: {data['age']}")
except json.JSONDecodeError as e:
    print(f"Failed to parse JSON: {e}")
```

### With Schema Validation

```python
import jsonschema

schema = {
    "type": "object",
    "required": ["name", "age"],
    "properties": {
        "name": {"type": "string"},
        "age": {"type": "integer", "minimum": 0},
        "email": {"type": "string", "format": "email"}
    }
}

try:
    jsonschema.validate(data, schema)
    print("Valid!")
except jsonschema.ValidationError as e:
    print(f"Validation error: {e}")
```

## Handling Parsing Failures

### Retry Strategy

```python
def get_structured_output(prompt, max_retries=3):
    for attempt in range(max_retries):
        response = client.messages.create(...)
        text = response.content[0].text

        # Try to extract JSON from response
        try:
            # Handle markdown code blocks
            if "```json" in text:
                text = text.split("```json")[1].split("```")[0]
            elif "```" in text:
                text = text.split("```")[1].split("```")[0]

            return json.loads(text.strip())
        except json.JSONDecodeError:
            if attempt < max_retries - 1:
                # Ask model to fix
                fix_prompt = f"Fix this invalid JSON: {text}"
                continue
            raise
```

## Complete Example

```
TASK: Invoice Data Extraction

You are an invoice data extractor. Extract information from invoice
text and return it as JSON.

SCHEMA:
{
  "invoice_number": "string",
  "date": "YYYY-MM-DD format",
  "vendor": {
    "name": "string",
    "address": "string or null"
  },
  "customer": {
    "name": "string",
    "email": "string or null"
  },
  "items": [
    {
      "description": "string",
      "quantity": "integer",
      "unit_price": "float",
      "total": "float"
    }
  ],
  "subtotal": "float",
  "tax": "float or null",
  "total": "float"
}

RULES:
- All monetary values as numbers (no currency symbols)
- Use null for missing optional fields
- Calculate item totals if not stated
- Dates must be YYYY-MM-DD

OUTPUT: Only valid JSON, no markdown, no explanations.

INVOICE TEXT:
Invoice #INV-2024-1234
Date: March 15, 2024

From: Tech Supplies Inc.
      123 Business Ave, Suite 100

Bill To: ABC Company
         john@abccompany.com

Description          Qty    Unit Price    Total
Widget Pro           5      $29.99        $149.95
Gadget Plus          2      $49.99        $99.99
Service Fee          1      $25.00        $25.00

                              Subtotal:   $274.94
                              Tax (8%):   $21.99
                              TOTAL:      $296.93
```

## Exercise

### Exercise 7.1: JSON Schema

Write a prompt to extract movie information with this schema:
```json
{
  "title": "string",
  "year": "integer",
  "director": "string",
  "cast": ["array of names"],
  "genres": ["array of genres"],
  "rating": "float 0-10",
  "summary": "string max 100 chars"
}
```

### Exercise 7.2: Handle Edge Cases

Write instructions for handling:
- Missing required fields
- Invalid data types
- Arrays that might be empty

### Exercise 7.3: XML Output

Convert this JSON schema to an XML format specification.

## Key Takeaways

- ✅ Specify exact schema for reliable structured output
- ✅ Use "only JSON, nothing else" instructions
- ✅ Handle edge cases with null values
- ✅ Validate output programmatically
- ✅ Implement retry logic for parsing failures

## Next Steps

→ [PE-08: Role Prompting & Personas](../PE-08-role-prompting/README.md)
