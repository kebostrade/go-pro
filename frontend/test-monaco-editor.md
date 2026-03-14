# Monaco Editor Integration Test Checklist

## Pre-Test Setup

```bash
cd /home/dima/Desktop/FUN/go-pro/frontend
bun install
bun run dev
```

## Test Checklist

### 1. Component Loading ✓
- [ ] Navigate to http://localhost:3000/exercises/1
- [ ] Editor loads without errors
- [ ] Initial code displays correctly
- [ ] Line numbers visible
- [ ] Minimap visible on right side

### 2. Theme Switching ✓
- [ ] Click theme dropdown
- [ ] Select "Light Theme"
- [ ] Editor background changes to light
- [ ] Select "Dark Theme"
- [ ] Editor background changes to dark

### 3. Font Size Control ✓
- [ ] Click font size dropdown
- [ ] Select 12px - text gets smaller
- [ ] Select 20px - text gets larger
- [ ] Select 14px - returns to normal

### 4. Code Editing ✓
- [ ] Type in editor
- [ ] Auto-completion popup appears
- [ ] Syntax highlighting works
- [ ] Line numbers update
- [ ] Minimap updates

### 5. Toolbar Actions ✓
- [ ] Click "Reset" - code reverts to initial
- [ ] Edit code
- [ ] Click "Copy" - shows "Copied!"
- [ ] Paste in text editor - code copied correctly
- [ ] Click "Fullscreen" - editor goes fullscreen
- [ ] Click "Exit" - returns to normal

### 6. Auto-save ✓
- [ ] Edit code
- [ ] Wait 2 seconds
- [ ] Refresh page
- [ ] Edited code persists (loaded from localStorage)

### 7. Code Execution ✓
- [ ] Click "Run Code" button
- [ ] Button shows "Running..." spinner
- [ ] After 1-2 seconds, results appear below
- [ ] Green header shows "All tests passed!"
- [ ] Score shows 100%
- [ ] Execution time displays (e.g., "3ms")
- [ ] Individual test cases show green checkmarks

### 8. Code Submission ✓
- [ ] Click "Submit" button
- [ ] Button shows "Submitting..." spinner
- [ ] After 1-2 seconds, results appear
- [ ] Results display same as "Run Code"
- [ ] On success, localStorage cleared

### 9. Test Results Display ✓
- [ ] Pass/fail indicators visible
- [ ] Test case names displayed
- [ ] Execution times per test shown
- [ ] Overall score displayed
- [ ] Execution time displayed

### 10. Keyboard Shortcuts ✓
- [ ] Press Ctrl+Enter (Cmd+Enter on Mac)
- [ ] Code runs (same as clicking "Run Code")
- [ ] Press Ctrl+Shift+Enter (Cmd+Shift+Enter on Mac)
- [ ] Code submits (same as clicking "Submit")

### 11. Error Boundary ✓
- [ ] Open browser console
- [ ] No React errors
- [ ] No Monaco errors
- [ ] If error occurs, error boundary catches it

### 12. Fullscreen Mode ✓
- [ ] Click "Fullscreen" button
- [ ] Editor fills entire screen
- [ ] Toolbar remains visible
- [ ] Action buttons remain visible
- [ ] Press Escape or click "Exit"
- [ ] Returns to normal view

### 13. Instructions Sidebar ✓
- [ ] Left sidebar shows exercise description
- [ ] Instructions list displays
- [ ] Click "Show" under Hints
- [ ] Hints display with 💡 icon
- [ ] Click "Hide" - hints disappear

### 14. Responsive Design ✓
- [ ] Resize browser to mobile width
- [ ] Editor remains functional
- [ ] Toolbar adapts to smaller screen
- [ ] Buttons remain accessible
- [ ] Text remains readable

### 15. localStorage Persistence ✓
- [ ] Edit code
- [ ] Close browser tab
- [ ] Reopen http://localhost:3000/exercises/1
- [ ] Code persists
- [ ] Click "Reset"
- [ ] Code reverts to initial
- [ ] localStorage cleared

## Browser Compatibility

### Chrome/Edge
- [ ] All tests pass
- [ ] No console errors
- [ ] Smooth performance

### Firefox
- [ ] All tests pass
- [ ] No console errors
- [ ] Smooth performance

### Safari (Mac)
- [ ] All tests pass
- [ ] No console errors
- [ ] Smooth performance

## Performance Tests

### Initial Load
- [ ] Editor loads in < 2 seconds
- [ ] No blocking JavaScript
- [ ] Smooth interaction immediately

### Code Execution
- [ ] "Run" completes in < 3 seconds (mock)
- [ ] "Submit" completes in < 3 seconds (mock)
- [ ] No UI freezing during execution

### Auto-save Performance
- [ ] Typing feels responsive
- [ ] No lag when typing quickly
- [ ] Auto-save doesn't block UI

## Accessibility Tests

### Keyboard Navigation
- [ ] Tab through all toolbar buttons
- [ ] Focus indicators visible
- [ ] Enter activates buttons
- [ ] Escape exits fullscreen

### Screen Reader
- [ ] aria-label attributes present
- [ ] Button purposes announced
- [ ] Editor content accessible
- [ ] Test results readable

### Color Contrast
- [ ] Text readable in light theme
- [ ] Text readable in dark theme
- [ ] Button text has good contrast
- [ ] Status indicators distinguishable

## Integration Tests

### API Integration (Mock)
- [ ] onSubmit function called with code
- [ ] Correct parameters passed
- [ ] Response handled correctly
- [ ] Errors handled gracefully

### localStorage
- [ ] Correct key format used
- [ ] Code saved on change
- [ ] Code loaded on mount
- [ ] Code cleared on reset/submit

## Error Scenarios

### Network Error
- [ ] Simulate API failure
- [ ] Error message displays
- [ ] User can retry
- [ ] No crash

### Invalid Code
- [ ] Submit invalid Go code
- [ ] Syntax errors highlighted
- [ ] Error messages clear
- [ ] Recovery possible

### Browser Storage Full
- [ ] Simulate localStorage quota exceeded
- [ ] Graceful degradation
- [ ] Editor remains functional
- [ ] User notified if needed

## Production Build Test

```bash
bun run build
bun start
```

- [ ] Production build succeeds
- [ ] No build warnings
- [ ] Editor works in production mode
- [ ] Assets loaded correctly
- [ ] Performance acceptable

## Test Results

**Test Date**: ___________
**Tester**: ___________
**Browser**: ___________
**OS**: ___________

**Overall Result**: ✅ Pass / ❌ Fail

**Issues Found**:
1. ___________________________________________
2. ___________________________________________
3. ___________________________________________

**Notes**:
___________________________________________
___________________________________________
___________________________________________
