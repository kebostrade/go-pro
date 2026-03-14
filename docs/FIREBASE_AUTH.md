# Firebase Authentication Setup Guide

This document explains the comprehensive Firebase Authentication system implemented for the Go Pro learning platform.

## Overview

The authentication system includes:
- **Email/Password Authentication**
- **Google OAuth**
- **GitHub OAuth**
- **User Profile Management**
- **Role-Based Access Control (RBAC)**
- **Progress Tracking**
- **Secure Firestore Rules**

## Architecture

### Components

```
frontend/src/
├── lib/
│   └── firebase.ts                    # Firebase initialization
├── contexts/
│   └── auth-context.tsx               # Authentication context & hooks
├── components/
│   └── auth/
│       ├── sign-in-form.tsx          # Sign-in UI
│       ├── sign-up-form.tsx          # Sign-up UI
│       └── protected-route.tsx       # Route protection wrapper
```

## Setup Instructions

### 1. Environment Configuration

Copy the environment template:
```bash
cd frontend
cp .env.local.example .env.local
```

The configuration is already set up with your Firebase project credentials:
- **Project ID**: `go-pro-platform`
- **App ID**: `1:434643680939:web:0ce27b5f6cda53789781ee`

### 2. Enable Authentication Methods

Go to [Firebase Console](https://console.firebase.google.com/project/go-pro-platform/authentication/providers):

1. **Email/Password**:
   - Already enabled by default
   - No additional configuration needed

2. **Google OAuth**:
   - Go to Authentication → Sign-in method
   - Click "Google" → Enable
   - Add authorized domains if needed

3. **GitHub OAuth**:
   - Go to Authentication → Sign-in method
   - Click "GitHub" → Enable
   - Create a GitHub OAuth App:
     - Go to GitHub Settings → Developer settings → OAuth Apps
     - Create new OAuth App
     - **Homepage URL**: `https://go-pro-platform.firebaseapp.com`
     - **Callback URL**: `https://go-pro-platform.firebaseapp.com/__/auth/handler`
   - Copy Client ID and Client Secret to Firebase

### 3. Deploy Security Rules

Deploy the Firestore security rules:
```bash
cd /home/dima/Desktop/FUN/go-pro
firebase deploy --only firestore:rules
```

### 4. Integrate in Your App

#### Wrap your app with AuthProvider

In your root layout (`app/layout.tsx`):

```tsx
import { AuthProvider } from '@/contexts/auth-context';

export default function RootLayout({ children }) {
  return (
    <html lang="en">
      <body>
        <AuthProvider>
          {children}
        </AuthProvider>
      </body>
    </html>
  );
}
```

#### Use the auth hook

```tsx
import { useAuth } from '@/contexts/auth-context';

function MyComponent() {
  const { user, userProfile, signOut } = useAuth();

  if (!user) {
    return <div>Please sign in</div>;
  }

  return (
    <div>
      <h1>Welcome, {userProfile?.displayName}!</h1>
      <p>Email: {user.email}</p>
      <p>Role: {userProfile?.role}</p>
      <button onClick={signOut}>Sign Out</button>
    </div>
  );
}
```

#### Protect routes

```tsx
import { ProtectedRoute } from '@/components/auth/protected-route';

export default function DashboardPage() {
  return (
    <ProtectedRoute>
      <div>
        <h1>Protected Dashboard</h1>
        {/* Your dashboard content */}
      </div>
    </ProtectedRoute>
  );
}

// Require email verification
export default function VerifiedPage() {
  return (
    <ProtectedRoute requireEmailVerification>
      <div>Email verified users only</div>
    </ProtectedRoute>
  );
}

// Require specific role
export default function InstructorPage() {
  return (
    <ProtectedRoute requiredRole="instructor">
      <div>Instructors only</div>
    </ProtectedRoute>
  );
}
```

## User Roles

The system implements three roles with hierarchical permissions:

1. **Student** (Level 1)
   - Default role for new users
   - Access to learning content
   - Can track their own progress
   - Can submit code exercises

2. **Instructor** (Level 2)
   - All student permissions
   - Create and edit courses/lessons/tutorials
   - View all student submissions
   - Access analytics

3. **Admin** (Level 3)
   - All instructor permissions
   - Manage user accounts
   - Change user roles
   - Delete any content
   - Full system access

### Role Hierarchy

```
Admin ≥ Instructor ≥ Student
```

When checking roles, higher roles automatically have access to lower-level features.

## Firestore Security Rules

The implemented rules provide:

### User Profiles (`/users/{userId}`)
- ✅ Read: Any authenticated user (for leaderboards)
- ✅ Create: Users can create their own profile (role defaults to 'student')
- ✅ Update: Users can update their own profile (except role/uid)
- ✅ Update: Admins can update any profile including roles
- ✅ Delete: Admins only

### Courses/Lessons/Tutorials
- ✅ Read: Public access (anyone can view content)
- ✅ Write: Instructors and admins only

### User Progress (`/progress/{userId}`)
- ✅ Read: Users see their own progress; admins see all
- ✅ Write: Users can update their own progress
- ❌ Delete: No one can delete progress (audit trail)

### Submissions (`/submissions/{submissionId}`)
- ✅ Read: Users see their own; instructors see all
- ✅ Create/Update: Users for their own submissions
- ❌ Delete: No deletions (audit trail)

### Comments (`/comments/{commentId}`)
- ✅ Read: All authenticated users
- ✅ Create: Authenticated users with verified email
- ✅ Update/Delete: Users for their own comments
- ✅ Delete: Admins can delete any (moderation)

## API Reference

### Authentication Methods

```tsx
const {
  // State
  user,              // Firebase User object
  userProfile,       // Extended user profile from Firestore
  loading,           // Loading state
  error,             // Error message

  // Auth methods
  signUp,            // (email, password, displayName?) => Promise<UserCredential>
  signIn,            // (email, password) => Promise<UserCredential>
  signInWithGoogle,  // () => Promise<UserCredential>
  signInWithGithub,  // () => Promise<UserCredential>
  signOut,           // () => Promise<void>

  // Profile management
  updateUserProfile,     // (data: Partial<UserProfile>) => Promise<void>
  updateUserEmail,       // (newEmail, password) => Promise<void>
  updateUserPassword,    // (currentPassword, newPassword) => Promise<void>
  sendPasswordReset,     // (email) => Promise<void>
  sendVerificationEmail, // () => Promise<void>

  // Progress tracking
  updateProgress,    // (lessonId, completed) => Promise<void>
  enrollInCourse,    // (courseId) => Promise<void>
} = useAuth();
```

### UserProfile Type

```typescript
interface UserProfile {
  uid: string;
  email: string | null;
  displayName: string | null;
  photoURL: string | null;
  emailVerified: boolean;
  createdAt: Date;
  lastLoginAt: Date;
  role: 'student' | 'instructor' | 'admin';
  progress?: {
    completedLessons: string[];
    currentCourse?: string;
    xp: number;
    level: number;
  };
  preferences?: {
    theme: 'light' | 'dark' | 'system';
    notifications: boolean;
    language: string;
  };
}
```

## Usage Examples

### Sign Up Flow

```tsx
import { SignUpForm } from '@/components/auth/sign-up-form';

export default function SignUpPage() {
  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-50">
      <SignUpForm />
    </div>
  );
}
```

### Sign In Flow

```tsx
import { SignInForm } from '@/components/auth/sign-in-form';

export default function SignInPage() {
  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-50">
      <SignInForm />
    </div>
  );
}
```

### Track Learning Progress

```tsx
function LessonPage({ lessonId }) {
  const { updateProgress } = useAuth();

  const handleComplete = async () => {
    await updateProgress(lessonId, true);
    alert('Progress saved!');
  };

  return (
    <div>
      <h1>Lesson Content</h1>
      <button onClick={handleComplete}>Mark as Complete</button>
    </div>
  );
}
```

### Display User Profile

```tsx
function ProfilePage() {
  const { userProfile, updateUserProfile } = useAuth();

  const handleUpdateTheme = async (theme: 'light' | 'dark' | 'system') => {
    await updateUserProfile({
      preferences: {
        ...userProfile?.preferences,
        theme,
      },
    });
  };

  return (
    <div>
      <h1>{userProfile?.displayName}</h1>
      <p>XP: {userProfile?.progress?.xp}</p>
      <p>Level: {userProfile?.progress?.level}</p>
      <p>Completed Lessons: {userProfile?.progress?.completedLessons.length}</p>
    </div>
  );
}
```

## Deployment

### Deploy to Firebase Hosting

1. Build your Next.js app:
   ```bash
   cd frontend
   bun run build
   ```

2. Deploy to Firebase:
   ```bash
   firebase deploy
   ```

3. Your app will be live at:
   - **Hosting**: https://go-pro-platform.web.app
   - **Auth Domain**: https://go-pro-platform.firebaseapp.com

## Security Best Practices

1. **Email Verification**: Users receive verification emails on signup
2. **Password Requirements**: Minimum 8 characters with uppercase, lowercase, and numbers
3. **Role Protection**: Users cannot elevate their own roles
4. **Audit Trail**: Submissions and progress cannot be deleted
5. **Token Validation**: All Firestore rules validate authentication tokens
6. **Rate Limiting**: Firebase automatically implements rate limiting

## Testing

### Test User Accounts

Create test accounts with different roles:

```tsx
// Student account (default)
await signUp('student@test.com', 'Test1234', 'Test Student');

// For instructor/admin roles, manually update Firestore:
// Go to Firestore Console → users → {userId} → Edit role field
```

### Test Security Rules

```bash
# Install Firebase Emulators
firebase init emulators

# Start emulators
firebase emulators:start

# Set environment variable
export NEXT_PUBLIC_USE_FIREBASE_EMULATORS=true

# Run your app
bun run dev
```

## Troubleshooting

### Common Issues

**"Firebase: Error (auth/operation-not-allowed)"**
- Go to Firebase Console → Authentication → Sign-in method
- Enable the authentication provider you're trying to use

**"Missing or insufficient permissions"**
- Deploy Firestore rules: `firebase deploy --only firestore:rules`
- Check user's role in Firestore Console

**"Email not verified"**
- User needs to click verification link in email
- Resend with `sendVerificationEmail()`

**Social auth popup closes immediately**
- Check OAuth app configuration (callback URLs)
- Ensure domain is authorized in Firebase Console

## Next Steps

1. **Customize UI**: Style the auth forms to match your design system
2. **Add More Features**:
   - Multi-factor authentication (MFA)
   - Phone authentication
   - Custom claims for fine-grained permissions
   - Session management
3. **Analytics**: Track authentication events with Firebase Analytics
4. **Monitoring**: Set up alerts for failed login attempts

## Support

For issues or questions:
- Firebase Console: https://console.firebase.google.com/project/go-pro-platform
- Firebase Documentation: https://firebase.google.com/docs/auth
- Project Documentation: See other docs in `/docs` directory
