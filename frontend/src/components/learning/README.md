# Monaco Code Editor Component

A feature-rich code editor component built with Monaco Editor for Go programming exercises.

## Features

### Core Functionality
- ✅ **Syntax Highlighting**: Go language support with full IntelliSense
- ✅ **Auto-completion**: Context-aware code suggestions
- ✅ **Error Detection**: Real-time syntax and semantic error highlighting
- ✅ **Line Numbers**: Clear code navigation with line numbering
- ✅ **Minimap**: Visual overview of code structure

### User Experience
- ✅ **Theme Toggle**: Switch between light and dark themes
- ✅ **Font Size Control**: Adjustable font size (12px - 20px)
- ✅ **Fullscreen Mode**: Immersive coding experience
- ✅ **Reset Code**: Restore initial exercise code
- ✅ **Copy to Clipboard**: Quick code copying
- ✅ **Auto-save**: Automatic localStorage persistence
- ✅ **Keyboard Shortcuts**: Fast code execution
  - `Ctrl+Enter`: Run code
  - `Ctrl+Shift+Enter`: Submit code

### Code Execution
- ✅ **Run Code**: Test locally without submission
- ✅ **Submit Code**: Official grading submission
- ✅ **Loading States**: Visual feedback during execution
- ✅ **Test Results**: Comprehensive test case feedback
- ✅ **Score Display**: Real-time performance metrics

### Test Results Display
- ✅ **Pass/Fail Indicators**: Clear visual status
- ✅ **Expected vs Actual**: Side-by-side comparison
- ✅ **Error Messages**: Detailed error information
- ✅ **Execution Time**: Performance metrics per test
- ✅ **Overall Score**: Percentage-based grading

## Installation

```bash
npm install @monaco-editor/react monaco-editor
```

## Usage

### Basic Example

```tsx
import MonacoCodeEditor from '@/components/learning/monaco-code-editor';

function ExercisePage() {
  const initialCode = `package main

import "fmt"

func main() {
    fmt.Println("Hello, World!")
}`;

  const handleSubmit = async (code: string) => {
    const response = await fetch('/api/exercises/submit', {
      method: 'POST',
      body: JSON.stringify({ code }),
    });
    return response.json();
  };

  return (
    <MonacoCodeEditor
      initialCode={initialCode}
      exerciseId="exercise-1"
      onSubmit={handleSubmit}
    />
  );
}
```

### With Error Boundary

```tsx
import MonacoCodeEditor from '@/components/learning/monaco-code-editor';
import EditorErrorBoundary from '@/components/learning/editor-error-boundary';

function ExercisePage() {
  return (
    <EditorErrorBoundary>
      <MonacoCodeEditor
        initialCode={initialCode}
        exerciseId="exercise-1"
        onSubmit={handleSubmit}
      />
    </EditorErrorBoundary>
  );
}
```

### Advanced Configuration

```tsx
<MonacoCodeEditor
  initialCode={exercise.initial_code}
  exerciseId={exercise.id}
  language="go"
  height="600px"
  readOnly={false}
  onSubmit={handleSubmit}
  testCases={[
    {
      name: "Test Case 1",
      input: "5",
      expected: "120",
    },
    {
      name: "Test Case 2",
      input: "10",
      expected: "3628800",
    },
  ]}
/>
```

## API Reference

### Props

#### `MonacoCodeEditorProps`

| Prop | Type | Default | Description |
|------|------|---------|-------------|
| `initialCode` | `string` | **required** | Initial code to display in editor |
| `exerciseId` | `string` | **required** | Unique identifier for exercise (used for localStorage) |
| `language` | `'go' \| 'javascript' \| 'python'` | `'go'` | Programming language for syntax highlighting |
| `height` | `string` | `'500px'` | Editor height (CSS value) |
| `readOnly` | `boolean` | `false` | Whether editor is read-only |
| `onSubmit` | `(code: string) => Promise<ExerciseResult>` | `undefined` | Async function to handle code submission |
| `testCases` | `TestCase[]` | `undefined` | Array of test cases for display |

### Types

#### `TestCase`

```typescript
interface TestCase {
  name: string;        // Test case name
  input?: string;      // Optional input data
  expected: string;    // Expected output
}
```

#### `ExerciseResult`

```typescript
interface ExerciseResult {
  success: boolean;              // Whether API call succeeded
  passed: boolean;               // Whether all tests passed
  score: number;                 // Score percentage (0-100)
  results: TestResult[];         // Individual test results
  execution_time_ms: number;     // Total execution time
  message: string;               // Result message
}
```

#### `TestResult`

```typescript
interface TestResult {
  name: string;                  // Test name
  passed: boolean;               // Pass/fail status
  expected: string;              // Expected output
  actual: string;                // Actual output
  error?: string;                // Error message if failed
  execution_time_ms?: number;    // Test execution time
}
```

## Local Storage

The component automatically saves code to localStorage with the key pattern:
```
exercise-{exerciseId}-code
```

This allows users to resume work across browser sessions. The saved code is cleared when:
- User clicks "Reset" button
- Exercise is successfully submitted with passing grade

## Keyboard Shortcuts

| Shortcut | Action |
|----------|--------|
| `Ctrl+Enter` (Mac: `Cmd+Enter`) | Run code |
| `Ctrl+Shift+Enter` (Mac: `Cmd+Shift+Enter`) | Submit code |
| `Ctrl+Z` (Mac: `Cmd+Z`) | Undo |
| `Ctrl+Y` (Mac: `Cmd+Y`) | Redo |
| `Ctrl+F` (Mac: `Cmd+F`) | Find |
| `Ctrl+H` (Mac: `Cmd+H`) | Find and replace |

## Editor Options

Default Monaco editor configuration:

```typescript
{
  minimap: { enabled: true },           // Show minimap
  fontSize: 14,                         // Font size in pixels
  lineNumbers: 'on',                    // Show line numbers
  scrollBeyondLastLine: false,          // No extra scroll space
  automaticLayout: true,                // Auto-resize with container
  tabSize: 4,                           // Tab width for Go
  wordWrap: 'on',                       // Wrap long lines
  formatOnPaste: true,                  // Format code on paste
  formatOnType: true,                   // Format as you type
  scrollbar: {
    vertical: 'visible',
    horizontal: 'visible',
  },
  suggestOnTriggerCharacters: true,     // Auto-complete triggers
  quickSuggestions: true,               // Enable suggestions
  folding: true,                        // Code folding
  renderWhitespace: 'selection',        // Show spaces in selection
}
```

## Theming

Two built-in themes:
- `vs-dark` (default): Dark theme for reduced eye strain
- `light`: Light theme for bright environments

Users can toggle between themes via the theme selector in the toolbar.

## Error Handling

### Editor Error Boundary

Wrap the editor in `EditorErrorBoundary` to gracefully handle Monaco editor crashes:

```tsx
<EditorErrorBoundary>
  <MonacoCodeEditor {...props} />
</EditorErrorBoundary>
```

### API Error Handling

The component handles submission errors gracefully:

```typescript
try {
  const result = await onSubmit(code);
  setResults(result);
} catch (error) {
  setResults({
    success: false,
    passed: false,
    score: 0,
    results: [],
    execution_time_ms: 0,
    message: error.message || 'Failed to submit code',
  });
}
```

## Accessibility

The component follows WCAG 2.1 guidelines:

- ✅ All buttons have `aria-label` attributes
- ✅ Keyboard navigation fully supported
- ✅ Focus indicators visible on all interactive elements
- ✅ Color contrast meets AA standards
- ✅ Screen reader compatible

## Performance

### Optimization Features

- **Code Splitting**: Monaco editor loaded asynchronously
- **Debounced Auto-save**: 1-second delay before saving to localStorage
- **Conditional Rendering**: Results only render when available
- **Memoization**: Editor options memoized to prevent unnecessary re-renders

### Bundle Size

- `@monaco-editor/react`: ~40KB (gzipped)
- `monaco-editor`: ~2MB (loaded on demand, cached by browser)

## Browser Support

- Chrome 90+
- Firefox 88+
- Safari 14+
- Edge 90+

## Examples

### Read-Only Viewer

```tsx
<MonacoCodeEditor
  initialCode={solutionCode}
  exerciseId="solution-view"
  readOnly={true}
  height="400px"
/>
```

### Custom Submit Handler

```tsx
const handleSubmit = async (code: string) => {
  // Custom validation
  if (!code.includes('package main')) {
    throw new Error('Code must include package main');
  }

  // Submit to API
  const response = await fetch('/api/submit', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ exerciseId, code, userId }),
  });

  if (!response.ok) {
    throw new Error('Submission failed');
  }

  return response.json();
};
```

### Multiple Languages

```tsx
const [language, setLanguage] = useState<'go' | 'javascript' | 'python'>('go');

<select value={language} onChange={(e) => setLanguage(e.target.value)}>
  <option value="go">Go</option>
  <option value="javascript">JavaScript</option>
  <option value="python">Python</option>
</select>

<MonacoCodeEditor
  initialCode={code}
  exerciseId="multi-lang"
  language={language}
/>
```

## Troubleshooting

### Editor Not Loading

**Problem**: Editor shows loading spinner indefinitely

**Solution**: Check browser console for errors. Monaco may be blocked by CSP or CORS policies.

```typescript
// Add to next.config.js
module.exports = {
  async headers() {
    return [
      {
        source: '/(.*)',
        headers: [
          {
            key: 'Content-Security-Policy',
            value: "script-src 'self' 'unsafe-eval' 'unsafe-inline'",
          },
        ],
      },
    ];
  },
};
```

### Code Not Saving

**Problem**: Code changes lost on page refresh

**Solution**: Check localStorage quota and permissions:

```typescript
// Check if localStorage is available
if (typeof window !== 'undefined' && window.localStorage) {
  localStorage.setItem(storageKey, code);
}
```

### Theme Not Applying

**Problem**: Theme changes don't take effect

**Solution**: Monaco requires page refresh for some theme changes. Use built-in themes only.

## Contributing

To add new features or fix bugs:

1. Modify component in `frontend/src/components/learning/monaco-code-editor.tsx`
2. Update types as needed
3. Test with example page at `frontend/src/app/exercises/[id]/page.tsx`
4. Update this README with new features

## License

MIT
