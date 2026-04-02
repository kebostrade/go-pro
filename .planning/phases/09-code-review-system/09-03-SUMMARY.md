---
phase: 09-code-review-system
plan: '03'
wave: 3
gap_closure: true
requirements:
  - REVIEW-01

must_haves:
  truths:
    - "ExerciseCard captures user's edited code from CodeEditor for submission"
    - "Submit button sends actual user code, not starterCode template"
  artifacts:
    - path: frontend/src/components/learning/exercise-card.tsx
      provides: Submit button that gets code from editor state
      exports: handleSubmitReview uses userCode state
---

# Phase 09 Plan 03: Fix Code Capture for REVIEW-01 Summary

## Objective

Fix REVIEW-01 gap: Submit button was sending `exercise.starterCode` instead of user's edited code. Integrated ExerciseCard with CodeEditor state to capture user's actual code.

## Changes Made

### Task 1: Add user code state to ExerciseCard

**Files Modified:**
- `frontend/src/components/learning/exercise-card.tsx`

**Changes:**
1. Added `onCodeChange` prop to `ExerciseCardProps` interface
2. Added `userCode` state initialized with `exercise.starterCode`
3. Added `handleCodeChange` function to update userCode when CodeEditor fires onCodeChange
4. Updated `handleSubmitReview` to send `userCode` instead of `exercise.starterCode`

### Task 2: Wire CodeEditor to ExerciseCard

**Files Modified:**
- `frontend/src/components/learning/exercise-card.tsx`

**Changes:**
1. Added CodeEditor import
2. Added CodeEditor component rendering when exercise has starterCode
3. Wired onCodeChange to handleCodeChange which updates userCode state

## Verification

```bash
grep -n "userCode" frontend/src/components/learning/exercise-card.tsx | head -10
# Output shows:
# - Line 38: const [userCode, setUserCode] = useState<string>(exercise.starterCode || '');
# - Line 52: code: userCode,
```

## Success Criteria Met

- [x] ExerciseCard has userCode state initialized with starterCode
- [x] handleSubmitReview sends userCode to API (grep verified)
- [x] onCodeChange prop exists in interface
- [x] CodeEditor onCodeChange updates ExerciseCard userCode state
- [x] User edits in CodeEditor are captured before submission

## Deviation from Plan

None - plan executed exactly as written.

## Files Modified

| File | Type | Description |
|------|------|-------------|
| frontend/src/components/learning/exercise-card.tsx | Modified | Added userCode state, CodeEditor integration, onCodeChange prop |

---

**Plan:** 09-03
**Status:** Complete
**Date:** 2026-04-02
