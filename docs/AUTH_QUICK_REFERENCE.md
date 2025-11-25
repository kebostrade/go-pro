# Firebase Authentication - Quick Reference

One-page reference for common authentication operations.

## 🚀 Quick Start

```tsx
import { useAuth } from '@/contexts/auth-context-advanced';

function MyComponent() {
  const { user, userProfile, signIn, signOut } = useAuth();

  if (!user) return <SignInButton />;

  return (
    <div>
      <p>Welcome, {userProfile?.displayName}!</p>
      <button onClick={signOut}>Sign Out</button>
    </div>
  );
}
```

## 🔐 Authentication Operations

### Sign Up
```tsx
const { signUp } = useAuth();

await signUp('user@example.com', 'password123', 'John Doe');
```

### Sign In
```tsx
const { signIn, signInWithGoogle, signInWithGithub } = useAuth();

// Email/Password
await signIn('user@example.com', 'password123');

// Google
await signInWithGoogle();

// GitHub
await signInWithGithub();
```

### Sign Out
```tsx
const { signOut } = useAuth();

await signOut();
```

## 📱 Phone Authentication

### Sign In with Phone
```tsx
import { RecaptchaVerifier } from 'firebase/auth';
import { auth } from '@/lib/firebase';

const verifier = new RecaptchaVerifier(auth, 'recaptcha-container', {
  size: 'invisible',
});

const confirmationResult = await signInWithPhone('+1234567890', verifier);
await verifyPhoneCode(confirmationResult.verificationId, '123456');
```

### Link Phone to Account
```tsx
const verificationId = await linkPhoneNumber('+1234567890', verifier);
// Verify with code
```

## 🔒 Multi-Factor Authentication

### Enroll MFA
```tsx
const verifier = new RecaptchaVerifier(auth, 'recaptcha-container');
await enrollMFA('+1234567890', verifier);
// Verify with code
```

### Check MFA Status
```tsx
const factors = getMFAInfo();
const isMFAEnabled = factors.length > 0;
```

### Unenroll MFA
```tsx
const factors = getMFAInfo();
await unenrollMFA(factors[0].uid);
```

## 👤 Profile Management

### Update Profile
```tsx
await updateUserProfile({
  displayName: 'New Name',
  photoURL: 'https://example.com/photo.jpg',
});
```

### Change Password
```tsx
await updateUserPassword('currentPassword', 'newPassword123');
```

### Change Email
```tsx
await updateUserEmail('newemail@example.com', 'currentPassword');
```

### Send Password Reset
```tsx
await sendPasswordReset('user@example.com');
```

### Send Verification Email
```tsx
await sendVerificationEmail();
```

## 🔗 Account Linking

### Link Provider
```tsx
await linkGoogleAccount();
await linkGithubAccount();
```

### Unlink Provider
```tsx
await unlinkProvider('google.com');
await unlinkProvider('github.com');
```

### Check Linked Providers
```tsx
const providers = getLinkedProviders();
// ['password', 'google.com', 'github.com']
```

## 💻 Session Management

### View Active Sessions
```tsx
const { sessions } = useAuth();

sessions.map(session => (
  <div>
    <p>{session.deviceInfo.browser} on {session.deviceInfo.platform}</p>
    <p>Last active: {session.lastActivity}</p>
  </div>
));
```

### Revoke Session
```tsx
await revokeSession(sessionId);
```

### Revoke All Sessions
```tsx
await revokeAllSessions(); // Signs out everywhere
```

## 📊 Security & Activity

### Check Security Score
```tsx
const { score, recommendations } = await checkAccountSecurity();

console.log(`Score: ${score}/100`);
recommendations.forEach(rec => console.log(`- ${rec}`));
```

### View Login History
```tsx
const history = await getLoginHistory(20); // Last 20 logins

history.map(entry => (
  <div>
    <p>{entry.success ? '✓' : '✗'} {entry.method}</p>
    <p>{new Date(entry.timestamp).toLocaleString()}</p>
  </div>
));
```

### Log Custom Activity
```tsx
await logActivity('completed_course', {
  courseId: 'golang-101',
  score: 95,
});
```

## 🎓 Progress Tracking

### Update Progress
```tsx
await updateProgress('lesson-123', true); // Mark complete
await updateProgress('lesson-123', false); // Mark incomplete
```

### Enroll in Course
```tsx
await enrollInCourse('course-456');
```

### Check Progress
```tsx
const { userProfile } = useAuth();

const xp = userProfile?.progress?.xp || 0;
const level = userProfile?.progress?.level || 1;
const completed = userProfile?.progress?.completedLessons || [];
```

## 🛡️ Protected Routes

### Basic Protection
```tsx
import { ProtectedRoute } from '@/components/auth/protected-route';

export default function DashboardPage() {
  return (
    <ProtectedRoute>
      <div>Protected content</div>
    </ProtectedRoute>
  );
}
```

### Require Email Verification
```tsx
<ProtectedRoute requireEmailVerification>
  <div>Verified users only</div>
</ProtectedRoute>
```

### Require Specific Role
```tsx
<ProtectedRoute requiredRole="instructor">
  <div>Instructors only</div>
</ProtectedRoute>

<ProtectedRoute requiredRole="admin">
  <div>Admins only</div>
</ProtectedRoute>
```

### Custom Fallback
```tsx
<ProtectedRoute fallbackPath="/auth/signin">
  <div>Protected content</div>
</ProtectedRoute>
```

## 👨‍💼 Admin Operations

### Get All Users
```tsx
const users = await adminGetAllUsers();
```

### Change User Role
```tsx
await adminUpdateUserRole('user-id', 'instructor');
await adminUpdateUserRole('user-id', 'admin');
```

### Disable User
```tsx
await adminDisableUser('user-id');
```

### Admin Dashboard
```tsx
import { UserManagementDashboard } from '@/components/admin/user-management-dashboard';

<ProtectedRoute requiredRole="admin">
  <UserManagementDashboard />
</ProtectedRoute>
```

## 🎨 UI Components

### Sign In Form
```tsx
import { SignInForm } from '@/components/auth/sign-in-form';

export default function SignInPage() {
  return (
    <div className="min-h-screen flex items-center justify-center">
      <SignInForm />
    </div>
  );
}
```

### Sign Up Form
```tsx
import { SignUpForm } from '@/components/auth/sign-up-form';

export default function SignUpPage() {
  return (
    <div className="min-h-screen flex items-center justify-center">
      <SignUpForm />
    </div>
  );
}
```

### Security Settings
```tsx
import { SecuritySettings } from '@/components/auth/security-settings';

export default function SecurityPage() {
  return (
    <ProtectedRoute>
      <SecuritySettings />
    </ProtectedRoute>
  );
}
```

## ⚙️ Configuration

### Environment Variables
```bash
# .env.local
NEXT_PUBLIC_FIREBASE_API_KEY=your-api-key
NEXT_PUBLIC_FIREBASE_AUTH_DOMAIN=your-project.firebaseapp.com
NEXT_PUBLIC_FIREBASE_PROJECT_ID=your-project
NEXT_PUBLIC_FIREBASE_STORAGE_BUCKET=your-project.firebasestorage.app
NEXT_PUBLIC_FIREBASE_MESSAGING_SENDER_ID=123456789
NEXT_PUBLIC_FIREBASE_APP_ID=1:123:web:abc123

# Optional
NEXT_PUBLIC_USE_FIREBASE_EMULATORS=false
NEXT_PUBLIC_RECAPTCHA_SITE_KEY=your-recaptcha-key
```

### Firebase Emulators
```bash
# Start emulators
firebase emulators:start

# Set environment variable
export NEXT_PUBLIC_USE_FIREBASE_EMULATORS=true
```

## 🔍 User State Checks

### Check Authentication
```tsx
const { user, loading } = useAuth();

if (loading) return <LoadingSpinner />;
if (!user) return <SignInPrompt />;
return <AuthenticatedContent />;
```

### Check Role
```tsx
const { userProfile } = useAuth();

const isStudent = userProfile?.role === 'student';
const isInstructor = userProfile?.role === 'instructor';
const isAdmin = userProfile?.role === 'admin';
```

### Check Verification
```tsx
const { user } = useAuth();

if (!user?.emailVerified) {
  return <EmailVerificationPrompt />;
}
```

### Check Security Features
```tsx
const { userProfile } = useAuth();

const hasMFA = userProfile?.security?.mfaEnabled;
const hasPhone = !!userProfile?.phoneNumber;
const linkedProviders = userProfile?.security?.linkedProviders?.length || 0;
```

## 🎯 Common Patterns

### Sign In Flow
```tsx
const handleSignIn = async (email: string, password: string) => {
  try {
    await signIn(email, password);
    router.push('/dashboard');
  } catch (error: any) {
    if (error.code === 'auth/multi-factor-auth-required') {
      // Handle MFA verification
      setShowMFAPrompt(true);
    } else {
      setError('Invalid credentials');
    }
  }
};
```

### Profile Update Flow
```tsx
const handleUpdateProfile = async (data: Partial<UserProfile>) => {
  try {
    await updateUserProfile(data);
    toast.success('Profile updated!');
  } catch (error) {
    toast.error('Failed to update profile');
  }
};
```

### Role-Based Rendering
```tsx
const { userProfile } = useAuth();

return (
  <div>
    {/* Everyone sees this */}
    <PublicContent />

    {/* Students and above */}
    {userProfile?.role && <StudentContent />}

    {/* Instructors and above */}
    {['instructor', 'admin'].includes(userProfile?.role) && (
      <InstructorContent />
    )}

    {/* Admins only */}
    {userProfile?.role === 'admin' && <AdminContent />}
  </div>
);
```

## 🚨 Error Handling

### Common Error Codes
```tsx
const getErrorMessage = (code: string) => {
  switch (code) {
    case 'auth/invalid-email':
      return 'Invalid email address';
    case 'auth/user-not-found':
      return 'No account with this email';
    case 'auth/wrong-password':
      return 'Incorrect password';
    case 'auth/too-many-requests':
      return 'Too many attempts. Try again later';
    case 'auth/email-already-in-use':
      return 'Email already registered';
    case 'auth/weak-password':
      return 'Password too weak';
    case 'auth/multi-factor-auth-required':
      return 'MFA verification required';
    default:
      return 'An error occurred';
  }
};
```

## 📦 Type Definitions

### UserProfile
```typescript
interface UserProfile {
  uid: string;
  email: string | null;
  displayName: string | null;
  photoURL: string | null;
  emailVerified: boolean;
  phoneNumber: string | null;
  role: 'student' | 'instructor' | 'admin';
  security?: SecuritySettings;
  progress?: {
    completedLessons: string[];
    currentCourse?: string;
    xp: number;
    level: number;
  };
}
```

### SecuritySettings
```typescript
interface SecuritySettings {
  mfaEnabled: boolean;
  phoneNumberVerified: boolean;
  linkedProviders: string[];
  lastPasswordChange?: Date;
  accountLockout: boolean;
  failedLoginAttempts: number;
}
```

## 🔗 Useful Links

- **Firebase Console**: https://console.firebase.google.com/project/go-pro-platform
- **Authentication Providers**: /authentication/providers
- **Firestore Database**: /firestore/data
- **Usage & Billing**: /usage

## 💡 Tips & Best Practices

1. **Always check loading state** before accessing user
2. **Use ProtectedRoute** for auth-required pages
3. **Enable MFA** for admin and instructor accounts
4. **Log important activities** for audit trails
5. **Check security score** regularly
6. **Monitor failed login attempts** for security
7. **Test with emulators** before production
8. **Keep Firebase SDK updated** for security patches
9. **Use TypeScript** for type safety
10. **Read error messages** - they're helpful!

## 📚 Documentation

- `FIREBASE_AUTH.md` - Full authentication guide
- `FIREBASE_AUTH_ADVANCED.md` - Advanced features
- `AUTH_MIGRATION_GUIDE.md` - Migration instructions
- `AUTHENTICATION_UPGRADE_SUMMARY.md` - Complete overview

---

**Quick help**: Check browser console for detailed error messages
