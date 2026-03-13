# PE-14: Data Analysis & Extraction

**Duration**: 2 hours
**Module**: 4 - Task-Specific Prompting

## Learning Objectives

- Extract structured data from unstructured text
- Perform text analysis and classification
- Transform and normalize data
- Generate reports and summaries

## Data Extraction

### Entity Extraction

```
Extract all entities from this text as JSON:

TEXT:
"Dr. Sarah Chen met with John Smith from Acme Corp on March 15,
2024 at their headquarters in San Francisco. They discussed the
$2.5 million contract for Q2 2024."

OUTPUT SCHEMA:
{
  "people": [
    {"name": "", "title": null, "organization": null}
  ],
  "organizations": [
    {"name": "", "location": null}
  ],
  "locations": ["list of places"],
  "dates": [
    {"text": "", "normalized": "YYYY-MM-DD"}
  ],
  "monetary_values": [
    {"amount": 0, "currency": "USD", "context": ""}
  ],
  "time_periods": [
    {"text": "", "start_date": null, "end_date": null}
  ]
}

RULES:
- Include all entities found
- Use null for missing optional fields
- Normalize dates to ISO format
- Extract monetary values as numbers
```

### Key-Value Extraction

```
Extract information from this product description:

DESCRIPTION:
"The UltraWidget Pro (SKU: UWP-2024-BLK) is our flagship product.
Priced at $149.99, it comes in black or silver. Dimensions:
8.5 x 5.2 x 1.1 inches. Weight: 12.3 oz. Battery life: 18 hours.
1-year warranty included."

OUTPUT:
{
  "product_name": "",
  "sku": "",
  "price": {
    "amount": 0,
    "currency": "USD"
  },
  "colors": [],
  "dimensions": {
    "length": 0,
    "width": 0,
    "height": 0,
    "unit": "inches"
  },
  "weight": {
    "value": 0,
    "unit": "oz"
  },
  "battery_life": {
    "value": 0,
    "unit": "hours"
  },
  "warranty": {
    "duration": 0,
    "unit": "years"
  }
}
```

### Table Extraction

```
Extract this table data as JSON array:

INVOICE TABLE:
| Item        | Qty | Unit Price | Total  |
|-------------|-----|------------|--------|
| Widget A    | 5   | $12.99     | $64.95 |
| Gadget B    | 2   | $45.00     | $90.00 |
| Service Fee | 1   | $25.00     | $25.00 |
|-------------|-----|------------|--------|
| Subtotal    |     |            | $179.95|
| Tax (8%)    |     |            | $14.40 |
| TOTAL       |     |            | $194.35|

OUTPUT:
{
  "items": [
    {
      "name": "",
      "quantity": 0,
      "unit_price": 0,
      "total": 0
    }
  ],
  "subtotal": 0,
  "tax": {
    "rate": 0,
    "amount": 0
  },
  "total": 0
}
```

## Text Classification

### Sentiment Analysis

```
Classify the sentiment of these customer reviews:

OUTPUT SCHEMA:
{
  "sentiment": "positive" | "negative" | "neutral" | "mixed",
  "confidence": 0.0-1.0,
  "aspects": [
    {
      "aspect": "product" | "service" | "shipping" | "price",
      "sentiment": "positive" | "negative" | "neutral",
      "reason": "brief explanation"
    }
  ],
  "summary": "one sentence summary"
}

REVIEWS:

1. "The product itself is great, but shipping took forever.
   Customer service was helpful when I called."

2. "Worst purchase ever. Complete waste of money."

3. "It's okay. Does what it says. Nothing special but nothing wrong either."
```

### Topic Classification

```
Classify this document into topics:

DOCUMENT:
[article about climate change policy]

TOPICS (select all that apply):
- Politics
- Environment
- Science
- Economics
- Health
- Technology
- International Relations

OUTPUT:
{
  "primary_topic": "",
  "secondary_topics": [],
  "relevance_scores": {
    "topic": 0.0-1.0
  },
  "keywords": ["top 5 keywords"],
  "summary": "2-3 sentence summary"
}
```

### Intent Classification

```
Classify the user intent:

USER MESSAGE: "I need to change my shipping address"

INTENTS:
- track_order: Check order status
- modify_order: Change order details
- return_item: Initiate return
- product_question: Ask about products
- account_issue: Account problems
- general_inquiry: Other questions

OUTPUT:
{
  "intent": "",
  "confidence": 0.0-1.0,
  "entities": {
    "order_id": null,
    "product": null,
    "address": null
  },
  "suggested_response": "brief response template"
}
```

## Data Transformation

### Format Conversion

```
Convert this data between formats:

INPUT (CSV):
name,age,city
John,30,New York
Sarah,25,Los Angeles
Mike,35,Chicago

OUTPUT (JSON):
{
  "people": [
    {"name": "John", "age": 30, "city": "New York"},
    {"name": "Sarah", "age": 25, "city": "Los Angeles"},
    {"name": "Mike", "age": 35, "city": "Chicago"}
  ],
  "metadata": {
    "count": 3,
    "fields": ["name", "age", "city"]
  }
}
```

### Data Normalization

```
Normalize these addresses to a consistent format:

ADDRESSES:
- "123 main st, new york, ny"
- "456 Elm Street, Suite 100, Los Angeles CA 90001"
- "789 Oak Ave., Chicago, Illinois"

OUTPUT FORMAT:
{
  "street": "number + street name",
  "unit": "apartment/suite number or null",
  "city": "city name",
  "state": "2-letter state code",
  "zip": "5-digit zip or null",
  "country": "USA"
}

RULES:
- Capitalize properly
- Expand abbreviations (St → Street)
- Standardize state codes
- Extract zip codes
```

### Unit Conversion

```
Convert measurements to standard units:

INPUT:
- "5 feet 10 inches"
- "180 lbs"
- "98.6 degrees Fahrenheit"
- "2.5 miles"

OUTPUT:
{
  "measurements": [
    {
      "original": "5 feet 10 inches",
      "normalized": {
        "value": 177.8,
        "unit": "cm"
      }
    },
    {
      "original": "180 lbs",
      "normalized": {
        "value": 81.65,
        "unit": "kg"
      }
    }
  ]
}
```

## Analysis & Reporting

### Statistical Summary

```
Analyze this dataset and provide a summary:

DATA: [sales figures, customer data, etc.]

OUTPUT:
{
  "overview": {
    "total_records": 0,
    "date_range": {"start": "", "end": ""},
    "key_metric": 0
  },
  "trends": [
    {
      "metric": "",
      "direction": "up" | "down" | "stable",
      "change_percent": 0,
      "significance": "high" | "medium" | "low"
    }
  ],
  "anomalies": [
    {
      "date": "",
      "metric": "",
      "value": 0,
      "expected": 0,
      "deviation": 0
    }
  ],
  "insights": [
    "actionable insight 1",
    "actionable insight 2"
  ]
}
```

### Comparative Analysis

```
Compare these two products:

PRODUCT A: [specifications]
PRODUCT B: [specifications]

OUTPUT:
{
  "comparison_matrix": {
    "feature": {
      "product_a": "value",
      "product_b": "value",
      "winner": "A" | "B" | "tie"
    }
  },
  "summary": {
    "product_a_strengths": [],
    "product_b_strengths": [],
    "product_a_weaknesses": [],
    "product_b_weaknesses": []
  },
  "recommendation": {
    "choose_a_if": ["condition 1", "condition 2"],
    "choose_b_if": ["condition 1", "condition 2"]
  }
}
```

## Best Practices

### Practice 1: Define Schema First

```
❌ Vague:
"Extract the important info"

✅ Specific:
"Extract as JSON with this schema:
{ name: string, email: string, phone: string }"
```

### Practice 2: Handle Missing Data

```
"If information is not available:
- Use null for optional fields
- Use 'N/A' for display fields
- Use 0 for numeric fields
- Note in a 'missing_fields' array"
```

### Practice 3: Validate Output

```
"Before responding, verify:
- JSON is valid
- All required fields present
- Data types are correct
- Values are in expected ranges"
```

## Exercise

### Exercise 14.1: Resume Extraction

Write a prompt to extract from a resume:
- Contact info
- Work experience (company, title, dates)
- Education
- Skills

### Exercise 14.2: Receipt Parsing

Write a prompt to extract from a receipt image:
- Store name and address
- Items purchased
- Prices and totals
- Date and time
- Payment method

### Exercise 14.3: Survey Analysis

Write a prompt to analyze survey responses:
- Calculate statistics
- Identify themes in open-ended responses
- Find correlations
- Generate summary report

## Key Takeaways

- ✅ Define explicit schemas for extracted data
- ✅ Handle missing/null data explicitly
- ✅ Normalize data to consistent formats
- ✅ Use structured output for programmatic use
- ✅ Include confidence scores when appropriate

## Next Steps

→ [PE-15: Creative Writing & Content](../PE-15-creative-writing/README.md)
