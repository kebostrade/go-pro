# Monaco Editor Setup Guide

Complete setup and usage guide for the Monaco Code Editor component.

## Installation Complete ✅

The Monaco editor has been successfully installed and configured:

```bash
# Dependencies installed
@monaco-editor/react: ^4.6.0
monaco-editor: ^0.45.0
```

## Component Location

```
frontend/src/components/learning/
├── monaco-code-editor.tsx        # Main editor component
├── editor-error-boundary.tsx     # Error boundary wrapper
└── README.md                     # Detailed documentation
```

## Quick Start

### 1. Basic Usage

```tsx
import MonacoCodeEditor from '@/components/learning/monaco-code-editor';

function MyExercisePage() {
  const initialCode = `package main

import "fmt"

func main() {
    fmt.Println("Hello, World!")
}`;

  const handleSubmit = async (code: string) => {
    const response = await fetch('/api/exercises/submit', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ code }),
    });
    return response.json();
  };

  return (
    <MonacoCodeEditor
      initialCode={initialCode}
      exerciseId="my-exercise"
      onSubmit={handleSubmit}
    />
  );
}
```

### 2. With Error Boundary (Recommended)

```tsx
import MonacoCodeEditor from '@/components/learning/monaco-code-editor';
import EditorErrorBoundary from '@/components/learning/editor-error-boundary';

function MyExercisePage() {
  return (
    <EditorErrorBoundary>
      <MonacoCodeEditor
        initialCode={code}
        exerciseId="exercise-1"
        onSubmit={handleSubmit}
      />
    </EditorErrorBoundary>
  );
}
```

## Complete Example

See the full working example at:
```
frontend/src/app/exercises/[id]/page.tsx
```

This example demonstrates:
- ✅ Full exercise page layout
- ✅ Instructions and hints sidebar
- ✅ Monaco editor integration
- ✅ Test results display
- ✅ Loading states
- ✅ Error handling

## Features Implemented

### Editor Features
- ✅ **Go Syntax Highlighting**: Full language support
- ✅ **IntelliSense**: Auto-completion and suggestions
- ✅ **Error Detection**: Real-time syntax checking
- ✅ **Line Numbers**: Clear code navigation
- ✅ **Minimap**: Visual code overview
- ✅ **Code Folding**: Collapse/expand code blocks

### User Experience
- ✅ **Theme Toggle**: Light/dark themes
- ✅ **Font Size Control**: 12px - 20px
- ✅ **Fullscreen Mode**: Immersive coding
- ✅ **Reset Code**: Restore initial state
- ✅ **Copy to Clipboard**: Quick code sharing
- ✅ **Auto-save**: localStorage persistence
- ✅ **Keyboard Shortcuts**:
  - `Ctrl+Enter`: Run code
  - `Ctrl+Shift+Enter`: Submit code

### Code Execution
- ✅ **Run Code**: Test without submission
- ✅ **Submit Code**: Official grading
- ✅ **Loading States**: Visual feedback
- ✅ **Test Results**: Comprehensive feedback

### Test Results Display
- ✅ **Pass/Fail Status**: Clear indicators
- ✅ **Expected vs Actual**: Side-by-side comparison
- ✅ **Error Messages**: Detailed debugging info
- ✅ **Execution Time**: Performance metrics
- ✅ **Score Display**: Percentage-based grading

## API Integration

### Expected Backend Response Format

```typescript
interface ExerciseResult {
  success: boolean;              // API call succeeded
  passed: boolean;               // All tests passed
  score: number;                 // 0-100
  results: TestResult[];         // Individual test results
  execution_time_ms: number;     // Total time
  message: string;               // Result message
}

interface TestResult {
  name: string;                  // Test case name
  passed: boolean;               // Pass/fail
  expected: string;              // Expected output
  actual: string;                // Actual output
  error?: string;                // Error if failed
  execution_time_ms?: number;    // Test time
}
```

### Example API Endpoint

```typescript
// pages/api/exercises/[id]/submit.ts
export default async function handler(req, res) {
  const { id } = req.query;
  const { code } = req.body;

  // Run tests against code
  const results = await runTests(code);

  res.json({
    success: true,
    passed: results.allPassed,
    score: results.score,
    results: results.testResults,
    execution_time_ms: results.executionTime,
    message: results.allPassed
      ? 'All tests passed! Great work!'
      : 'Some tests failed. Keep trying!',
  });
}
```

## Configuration Options

### Component Props

```typescript
interface MonacoCodeEditorProps {
  initialCode: string;           // Required: Starting code
  exerciseId: string;            // Required: Unique ID
  language?: 'go' | 'javascript' | 'python';  // Default: 'go'
  height?: string;               // Default: '500px'
  readOnly?: boolean;            // Default: false
  onSubmit?: (code: string) => Promise<ExerciseResult>;
  testCases?: TestCase[];
}
```

### Editor Options

Default Monaco configuration:
```typescript
{
  minimap: { enabled: true },
  fontSize: 14,
  lineNumbers: 'on',
  tabSize: 4,                    // Go convention
  wordWrap: 'on',
  formatOnPaste: true,
  formatOnType: true,
  folding: true,
}
```

## Keyboard Shortcuts

| Shortcut | Action |
|----------|--------|
| `Ctrl+Enter` | Run code |
| `Ctrl+Shift+Enter` | Submit code |
| `Ctrl+Z` | Undo |
| `Ctrl+Y` | Redo |
| `Ctrl+F` | Find |
| `Ctrl+H` | Find and replace |
| `Alt+Up/Down` | Move line up/down |
| `Ctrl+/` | Toggle comment |

## localStorage Behavior

The editor automatically saves code to localStorage:
- **Key Format**: `exercise-{exerciseId}-code`
- **Auto-save**: 1-second debounce after typing stops
- **Auto-load**: Loads saved code on mount
- **Auto-clear**: Clears on successful submission or reset

### Managing localStorage

```typescript
// Manually clear localStorage for an exercise
localStorage.removeItem('exercise-my-exercise-code');

// Clear all exercise data
Object.keys(localStorage)
  .filter(key => key.startsWith('exercise-'))
  .forEach(key => localStorage.removeItem(key));
```

## Theming

Two built-in themes:
- **vs-dark** (default): Dark theme
- **light**: Light theme

Users can toggle via the theme selector in the toolbar.

## Browser Support

- ✅ Chrome 90+
- ✅ Firefox 88+
- ✅ Safari 14+
- ✅ Edge 90+

## Performance

### Bundle Size
- Editor component: ~40KB (gzipped)
- Monaco editor: ~2MB (lazy-loaded, browser-cached)

### Optimization
- Code splitting for Monaco assets
- Debounced auto-save (1s)
- Conditional rendering of results
- Automatic layout adjustment

## Troubleshooting

### Monaco Not Loading

**Symptom**: Loading spinner indefinitely

**Solution**: Check Next.js configuration for Monaco assets

```javascript
// next.config.js
module.exports = {
  webpack: (config) => {
    config.module.rules.push({
      test: /\.ttf$/,
      type: 'asset/resource',
    });
    return config;
  },
};
```

### Code Not Saving

**Symptom**: Code lost on refresh

**Solution**: Check localStorage availability

```typescript
if (typeof window !== 'undefined' && window.localStorage) {
  // localStorage available
} else {
  // Use session storage or cookies as fallback
}
```

### Theme Not Applying

**Symptom**: Theme doesn't change

**Solution**: Monaco themes are applied on mount. Some require editor restart.

## Testing

### Development Server

```bash
cd frontend
bun run dev
# Visit http://localhost:3000/exercises/1
```

### Production Build

```bash
bun run build
bun start
```

### Type Checking

```bash
bun run type-check
```

## Next Steps

1. **Connect to Backend API**
   - Implement actual code execution endpoint
   - Add user authentication
   - Track exercise progress

2. **Add More Languages**
   - JavaScript/TypeScript support
   - Python support
   - Configure language-specific settings

3. **Enhanced Features**
   - Code snippets library
   - Real-time collaboration
   - AI-powered hints
   - Code review comments

4. **Analytics**
   - Track completion rates
   - Monitor execution times
   - Identify common errors

## Resources

- Component Docs: `frontend/src/components/learning/README.md`
- Example Page: `frontend/src/app/exercises/[id]/page.tsx`
- Monaco Docs: https://microsoft.github.io/monaco-editor/
- React Monaco: https://github.com/suren-atoyan/monaco-react

## Support

For issues or questions:
1. Check the detailed README in `frontend/src/components/learning/README.md`
2. Review the example implementation
3. Check browser console for errors
4. Verify API response format matches expected structure
