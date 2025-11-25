# Backend Integration Guide

## Overview

Frontend now integrates Firebase ID tokens with the Go backend for user synchronization.

## Architecture

```
Firebase Auth (Identity) → Frontend → Go Backend (Source of Truth)
                              ↓
                         Firestore (Cache for offline)
```

## Changes Made

### 1. API Client (`src/lib/api.ts`)

**Token Management:**
- `getIdToken()`: Gets fresh Firebase ID token
- `getAuthHeaders()`: Adds `Authorization: Bearer <token>` header
- `setTokenRefreshCallback()`: Registers token refresh handler

**Error Handling:**
- 401 responses: Auto-retry with refreshed token
- 403 responses: Show permission denied errors
- Network errors: Graceful fallback to Firestore cache

**New Endpoints:**
```typescript
api.verifyToken(idToken)           // POST /api/v1/auth/verify
api.getCurrentUser()               // GET /api/v1/auth/me
api.updateBackendProfile(data)     // PUT /api/v1/auth/profile
```

**Updated Endpoints (now authenticated):**
- `api.getProgress(userId)` - Requires auth
- `api.updateProgress(userId, lessonId, data)` - Requires auth
- `api.submitExercise(exerciseId, code)` - Requires auth
- `api.completeLesson(lessonId)` - Requires auth
- `api.getUserProgress(userId)` - Requires auth
- `api.getProgressStats(userId)` - Requires auth
- `api.updateLessonProgress(userId, lessonId, status, score)` - Requires auth

### 2. Auth Context (`src/contexts/auth-context.tsx`)

**New State:**
```typescript
backendUser: BackendUser | null  // Backend user data
```

**New Methods:**
```typescript
syncWithBackend()  // Manual backend sync
```

**Backend Sync Flow:**
1. Firebase login → Get ID token
2. Call `/api/v1/auth/verify` with token
3. Backend validates token and returns/creates user
4. Store `backendUser` in context

**Progress Sync Strategy:**
- **Backend first**: Update backend as source of truth
- **Firestore cache**: Fallback for offline support
- **Periodic sync**: Every 5 minutes
- **On login**: Immediate sync

**Error Handling:**
- Backend sync failures: Keep Firestore profile for offline
- Token refresh on 401: Auto-retry with fresh token
- Graceful degradation: App works even if backend is down

## Usage

### Sign In Flow

```typescript
const { signIn, backendUser } = useAuth();

// Sign in with Firebase
await signIn(email, password);

// backendUser is automatically synced
console.log(backendUser.id);  // Backend user ID
```

### Progress Updates

```typescript
const { updateProgress, backendUser } = useAuth();

// Updates backend first, then Firestore cache
await updateProgress(lessonId, completed);
```

### Manual Sync

```typescript
const { syncWithBackend } = useAuth();

// Force sync with backend
await syncWithBackend();
```

### API Calls

```typescript
// Authenticated requests automatically include token
const progress = await api.getUserProgress(userId);

// Token refresh handled automatically on 401
const stats = await api.getProgressStats(userId);
```

## Backend Requirements

Backend must implement these endpoints:

### POST /api/v1/auth/verify
**Request:**
```json
{
  "id_token": "firebase-id-token"
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "user": {
      "id": "backend-user-id",
      "firebase_uid": "firebase-uid",
      "email": "user@example.com",
      "display_name": "John Doe",
      "photo_url": "https://...",
      "role": "student",
      "created_at": "2025-01-01T00:00:00Z",
      "updated_at": "2025-01-01T00:00:00Z"
    },
    "token": "jwt-token-if-needed"
  }
}
```

### GET /api/v1/auth/me
**Headers:**
```
Authorization: Bearer <firebase-id-token>
```

**Response:**
```json
{
  "success": true,
  "data": {
    "id": "backend-user-id",
    "firebase_uid": "firebase-uid",
    "email": "user@example.com",
    // ... other user fields
  }
}
```

### PUT /api/v1/auth/profile
**Headers:**
```
Authorization: Bearer <firebase-id-token>
```

**Request:**
```json
{
  "display_name": "New Name",
  "photo_url": "https://...",
  "preferences": {}
}
```

## Token Lifecycle

1. **Login**: Firebase generates ID token
2. **First Request**: Frontend sends token to `/auth/verify`
3. **Subsequent Requests**: Token auto-added to all authenticated requests
4. **Token Refresh**: Auto-refresh when expired (handled by Firebase SDK)
5. **401 Response**: Auto-retry with refreshed token
6. **Logout**: Clear both Firebase and backend state

## Offline Support

- **Firestore as cache**: User profile and progress cached locally
- **Background sync**: When online, sync to backend every 5 minutes
- **Progressive enhancement**: Works offline, syncs when online
- **Conflict resolution**: Backend is source of truth

## Security

- ID tokens are validated by backend using Firebase Admin SDK
- Tokens expire after 1 hour (Firebase default)
- Auto-refresh keeps users logged in
- 401 errors trigger re-authentication
- 403 errors indicate insufficient permissions

## Testing

### Check Token Flow
```typescript
// In browser console
const user = auth.currentUser;
const token = await user.getIdToken();
console.log('Token:', token);

// Verify backend
const response = await fetch('http://localhost:8080/api/v1/auth/verify', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({ id_token: token })
});
console.log('Backend response:', await response.json());
```

### Check Sync Status
```typescript
const { backendUser, userProfile } = useAuth();

console.log('Firebase UID:', userProfile?.uid);
console.log('Backend ID:', backendUser?.id);
console.log('Synced:', backendUser?.firebase_uid === userProfile?.uid);
```

## Troubleshooting

### Backend Sync Fails
- Check backend is running: `http://localhost:8080/api/v1/health`
- Check CORS settings in backend
- Check Firebase ID token is valid
- Check console for detailed error messages

### 401 Errors
- Token might be expired, should auto-retry
- Check Firebase user is signed in: `auth.currentUser`
- Check token refresh callback is registered

### Progress Not Syncing
- Check backend user is loaded: `backendUser !== null`
- Check network requests in DevTools
- Check backend logs for errors
- Firestore cache should still work offline

## Environment Variables

Add to `.env.local`:
```bash
NEXT_PUBLIC_API_URL=http://localhost:8080
NEXT_PUBLIC_FIREBASE_API_KEY=your-api-key
NEXT_PUBLIC_FIREBASE_AUTH_DOMAIN=your-auth-domain
NEXT_PUBLIC_FIREBASE_PROJECT_ID=your-project-id
# ... other Firebase config
```

## Migration Notes

### Breaking Changes
- All authenticated API methods now require auth
- `authToken` parameter removed (handled automatically)
- Progress API now requires backend user

### Compatibility
- Firestore still works as cache
- Offline mode still functional
- Existing user profiles preserved
