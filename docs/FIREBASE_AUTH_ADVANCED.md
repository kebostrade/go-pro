# Firebase Authentication - Advanced Features

This document covers the advanced authentication features implemented for the Go Pro learning platform, including Multi-Factor Authentication, Session Management, Admin Dashboard, and enhanced security features.

## Table of Contents

1. [Overview](#overview)
2. [Multi-Factor Authentication (MFA)](#multi-factor-authentication)
3. [Phone Authentication](#phone-authentication)
4. [Session Management](#session-management)
5. [Login History & Activity Tracking](#login-history--activity-tracking)
6. [Security Dashboard](#security-dashboard)
7. [Admin User Management](#admin-user-management)
8. [Custom Claims & Advanced RBAC](#custom-claims--advanced-rbac)
9. [Implementation Guide](#implementation-guide)
10. [API Reference](#api-reference)

## Overview

The advanced authentication system extends the basic auth features with:

- **Multi-Factor Authentication (MFA)** - SMS-based two-factor authentication
- **Phone Authentication** - Sign in with phone number
- **Session Management** - Track and revoke active sessions
- **Login History** - Audit trail of all login attempts
- **Security Dashboard** - User-facing security score and recommendations
- **Admin Dashboard** - Comprehensive user management interface
- **Enhanced Security** - Account lockout, failed login tracking, activity logging
- **Account Linking** - Connect multiple authentication providers

## Multi-Factor Authentication

### Overview

MFA adds an additional layer of security by requiring users to verify their identity using a second factor (SMS code) in addition to their password.

### Features

- **SMS-based verification** via Firebase Phone Auth
- **Enrollment management** - Users can enable/disable MFA
- **Multiple factors** - Support for multiple phone numbers
- **Backup codes** (future enhancement)

### Setup

#### 1. Enable Phone Authentication in Firebase

```bash
# Go to Firebase Console
https://console.firebase.google.com/project/go-pro-platform/authentication/providers

# Enable Phone sign-in method
```

#### 2. Configure reCAPTCHA

Phone auth requires reCAPTCHA verification:

```tsx
import { RecaptchaVerifier } from 'firebase/auth';
import { auth } from '@/lib/firebase';

// Create reCAPTCHA verifier (invisible or visible)
const recaptchaVerifier = new RecaptchaVerifier(auth, 'recaptcha-container', {
  size: 'invisible', // or 'normal' for visible widget
  callback: (response: any) => {
    // reCAPTCHA solved
  },
});
```

#### 3. Enroll MFA

```tsx
import { useAuth } from '@/contexts/auth-context-advanced';

function MFAEnrollment() {
  const { enrollMFA } = useAuth();
  const [phoneNumber, setPhoneNumber] = useState('');

  const handleEnrollMFA = async () => {
    const verifier = new RecaptchaVerifier(auth, 'recaptcha-container', {
      size: 'invisible',
    });

    try {
      await enrollMFA(phoneNumber, verifier);
      // Verify code in next step
    } catch (error) {
      console.error('MFA enrollment failed:', error);
    }
  };

  return (
    <div>
      <input
        type="tel"
        value={phoneNumber}
        onChange={(e) => setPhoneNumber(e.target.value)}
        placeholder="+1234567890"
      />
      <button onClick={handleEnrollMFA}>Enable 2FA</button>
      <div id="recaptcha-container"></div>
    </div>
  );
}
```

#### 4. Verify MFA During Sign-In

When MFA is enabled, the sign-in flow requires an additional verification step:

```tsx
const handleSignIn = async (email: string, password: string) => {
  try {
    const result = await signIn(email, password);
    // User signed in successfully
  } catch (error: any) {
    if (error.code === 'auth/multi-factor-auth-required') {
      // MFA required - prompt for verification code
      const resolver = error.resolver;
      // Handle MFA verification flow
    }
  }
};
```

### Best Practices

1. **Always offer backup methods** - Don't lock users out if they lose their phone
2. **Test thoroughly** - Phone auth can be tricky in development
3. **Use invisible reCAPTCHA** - Better UX than visible widget
4. **Store backup codes** - Generate backup codes for account recovery

## Phone Authentication

### Direct Phone Sign-In

Allow users to sign in using only their phone number:

```tsx
import { useAuth } from '@/contexts/auth-context-advanced';
import { RecaptchaVerifier } from 'firebase/auth';
import { auth } from '@/lib/firebase';

function PhoneSignIn() {
  const { signInWithPhone, verifyPhoneCode } = useAuth();
  const [phoneNumber, setPhoneNumber] = useState('');
  const [verificationCode, setVerificationCode] = useState('');
  const [verificationId, setVerificationId] = useState('');

  const handleSendCode = async () => {
    const verifier = new RecaptchaVerifier(auth, 'recaptcha-container', {
      size: 'normal',
    });

    try {
      const confirmationResult = await signInWithPhone(phoneNumber, verifier);
      setVerificationId(confirmationResult.verificationId);
    } catch (error) {
      console.error('Failed to send verification code:', error);
    }
  };

  const handleVerifyCode = async () => {
    try {
      await verifyPhoneCode(verificationId, verificationCode);
      // User signed in
    } catch (error) {
      console.error('Invalid verification code:', error);
    }
  };

  return (
    <div>
      {!verificationId ? (
        <>
          <input
            type="tel"
            value={phoneNumber}
            onChange={(e) => setPhoneNumber(e.target.value)}
            placeholder="+1234567890"
          />
          <button onClick={handleSendCode}>Send Code</button>
          <div id="recaptcha-container"></div>
        </>
      ) : (
        <>
          <input
            type="text"
            value={verificationCode}
            onChange={(e) => setVerificationCode(e.target.value)}
            placeholder="123456"
          />
          <button onClick={handleVerifyCode}>Verify</button>
        </>
      )}
    </div>
  );
}
```

### Link Phone Number to Existing Account

```tsx
const { linkPhoneNumber } = useAuth();

const handleLinkPhone = async (phoneNumber: string) => {
  const verifier = new RecaptchaVerifier(auth, 'recaptcha-container', {
    size: 'invisible',
  });

  try {
    const verificationId = await linkPhoneNumber(phoneNumber, verifier);
    // Prompt user for verification code
  } catch (error) {
    console.error('Failed to link phone number:', error);
  }
};
```

## Session Management

### Overview

Track all active sessions and allow users to revoke access from unknown devices.

### Features

- **Active session tracking** - Device info, IP address, location
- **Session revocation** - Sign out from specific devices
- **Revoke all sessions** - Emergency sign-out from all devices
- **Session persistence** - Remember trusted devices

### Implementation

#### Display Active Sessions

```tsx
import { SecuritySettings } from '@/components/auth/security-settings';

function AccountPage() {
  return <SecuritySettings />;
}
```

The SecuritySettings component automatically displays:
- Current device session
- All active sessions
- Device information (browser, OS, platform)
- Last activity timestamp
- Option to revoke individual sessions

#### Revoke Session

```tsx
const { revokeSession } = useAuth();

const handleRevokeSession = async (sessionId: string) => {
  try {
    await revokeSession(sessionId);
    alert('Session revoked successfully');
  } catch (error) {
    console.error('Failed to revoke session:', error);
  }
};
```

#### Revoke All Sessions

```tsx
const { revokeAllSessions } = useAuth();

const handleSignOutEverywhere = async () => {
  if (confirm('Sign out from all devices?')) {
    await revokeAllSessions();
  }
};
```

## Login History & Activity Tracking

### Overview

Comprehensive audit trail of all authentication events and user activity.

### Features

- **Login attempts** - Successful and failed
- **Authentication method** - Email, Google, GitHub, Phone
- **Device information** - Browser, OS, user agent
- **Timestamp tracking** - When events occurred
- **Activity logging** - Custom events and actions

### Viewing Login History

The Security Settings UI includes a Login History tab that displays:

```tsx
import { useAuth } from '@/contexts/auth-context-advanced';

function LoginHistory() {
  const { getLoginHistory } = useAuth();
  const [history, setHistory] = useState([]);

  useEffect(() => {
    loadHistory();
  }, []);

  const loadHistory = async () => {
    const history = await getLoginHistory(20); // Last 20 logins
    setHistory(history);
  };

  return (
    <div>
      {history.map((entry) => (
        <div key={entry.timestamp}>
          <p>{entry.success ? 'Success' : 'Failed'}</p>
          <p>{entry.method} - {new Date(entry.timestamp).toLocaleString()}</p>
          <p>{entry.userAgent}</p>
        </div>
      ))}
    </div>
  );
}
```

### Logging Custom Activity

```tsx
const { logActivity } = useAuth();

// Log when user completes an important action
await logActivity('completed_advanced_course', {
  courseId: 'golang-advanced',
  score: 95,
  timeSpent: 7200,
});

// Log security-related events
await logActivity('password_changed', {
  method: 'user_initiated',
});
```

## Security Dashboard

### Overview

User-facing security score and personalized recommendations to improve account security.

### Features

- **Security score** (0-100) based on:
  - Email verification
  - MFA enabled
  - Number of linked providers
  - Phone number added
  - Password strength (future)
- **Recommendations** - Specific actions to improve security
- **Visual indicators** - Color-coded score (red/yellow/green)
- **Progress tracking** - See security improve over time

### Using the Security Dashboard

```tsx
import { SecuritySettings } from '@/components/auth/security-settings';

export default function SecurityPage() {
  return (
    <div className="container mx-auto py-8">
      <SecuritySettings />
    </div>
  );
}
```

### Programmatic Security Check

```tsx
const { checkAccountSecurity } = useAuth();

const checkSecurity = async () => {
  const { score, recommendations } = await checkAccountSecurity();

  console.log(`Security Score: ${score}/100`);
  recommendations.forEach((rec) => {
    console.log(`- ${rec}`);
  });
};
```

## Admin User Management

### Overview

Comprehensive dashboard for administrators to manage all users, roles, and permissions.

### Features

- **User listing** - View all users with filtering and search
- **Role management** - Change user roles (Student/Instructor/Admin)
- **Account actions** - Disable/enable user accounts
- **User details** - View complete user profile and activity
- **Statistics** - User counts by role, verification status
- **Search & filter** - Find users by name, email, or role

### Using the Admin Dashboard

```tsx
import { ProtectedRoute } from '@/components/auth/protected-route';
import { UserManagementDashboard } from '@/components/admin/user-management-dashboard';

export default function AdminUsersPage() {
  return (
    <ProtectedRoute requiredRole="admin">
      <UserManagementDashboard />
    </ProtectedRoute>
  );
}
```

### Admin API Functions

```tsx
const {
  adminGetAllUsers,
  adminUpdateUserRole,
  adminDisableUser,
  adminDeleteUser,
} = useAuth();

// Get all users
const users = await adminGetAllUsers();

// Change user role
await adminUpdateUserRole('user-id-123', 'instructor');

// Disable user account
await adminDisableUser('user-id-123');

// Delete user (requires backend implementation)
// await adminDeleteUser('user-id-123');
```

### Admin Dashboard Features

#### User Statistics

- Total users count
- Users by role (Students, Instructors, Admins)
- Email verification rate
- MFA adoption rate

#### Search & Filter

```tsx
// Search by name or email
const filteredUsers = users.filter((user) =>
  user.displayName?.includes(searchTerm) || user.email?.includes(searchTerm)
);

// Filter by role
const instructors = users.filter((user) => user.role === 'instructor');
```

#### Bulk Actions

- Export user data (CSV)
- Send bulk notifications
- Mass role updates (careful!)

## Custom Claims & Advanced RBAC

### Overview

Custom claims provide fine-grained permissions beyond simple roles.

### Use Cases

- **Feature flags** - Enable beta features for specific users
- **Tenant isolation** - Multi-tenant applications
- **Permission sets** - Complex permission hierarchies
- **Time-limited access** - Temporary elevated permissions

### Implementation (Backend Required)

Custom claims must be set using Firebase Admin SDK on the backend:

```typescript
// Backend (Cloud Functions or Admin API)
import * as admin from 'firebase-admin';

// Set custom claims
await admin.auth().setCustomUserClaims(userId, {
  permissions: ['create_course', 'grade_assignments'],
  features: ['beta_ai_assistant'],
  organizationId: 'org-123',
  subscriptionTier: 'premium',
});
```

### Accessing Custom Claims in Frontend

```tsx
const { user, userProfile } = useAuth();

// Custom claims available in user token
const tokenResult = await user?.getIdTokenResult();
const customClaims = tokenResult?.claims;

if (customClaims?.permissions?.includes('create_course')) {
  // Show course creation UI
}

if (customClaims?.features?.includes('beta_ai_assistant')) {
  // Enable AI assistant
}
```

### Security Rules with Custom Claims

Update Firestore security rules to check custom claims:

```
rules_version = '2';
service cloud.firestore {
  match /databases/{database}/documents {
    function hasPermission(permission) {
      return request.auth.token.permissions != null &&
             permission in request.auth.token.permissions;
    }

    match /courses/{courseId} {
      allow create: if hasPermission('create_course');
    }
  }
}
```

## Implementation Guide

### Step 1: Replace Auth Context

Update your app to use the advanced auth context:

```tsx
// app/layout.tsx
import { AuthProvider } from '@/contexts/auth-context-advanced';

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

### Step 2: Add Security Settings Page

```tsx
// app/settings/security/page.tsx
import { ProtectedRoute } from '@/components/auth/protected-route';
import { SecuritySettings } from '@/components/auth/security-settings';

export default function SecuritySettingsPage() {
  return (
    <ProtectedRoute>
      <SecuritySettings />
    </ProtectedRoute>
  );
}
```

### Step 3: Add Admin Dashboard

```tsx
// app/admin/users/page.tsx
import { ProtectedRoute } from '@/components/auth/protected-route';
import { UserManagementDashboard } from '@/components/admin/user-management-dashboard';

export default function AdminUsersPage() {
  return (
    <ProtectedRoute requiredRole="admin">
      <UserManagementDashboard />
    </ProtectedRoute>
  );
}
```

### Step 4: Update Firestore Security Rules

The existing rules in `firestore.rules` already support the advanced features. Just deploy them:

```bash
firebase deploy --only firestore:rules
```

### Step 5: Add reCAPTCHA Site Key (for Phone Auth)

```tsx
// .env.local
NEXT_PUBLIC_RECAPTCHA_SITE_KEY=your-recaptcha-site-key
```

Get your reCAPTCHA site key from:
https://www.google.com/recaptcha/admin/create

## API Reference

### Advanced Auth Context

```typescript
interface AuthContextType {
  // Basic auth (same as before)
  user: User | null;
  userProfile: UserProfile | null;
  loading: boolean;
  error: string | null;
  signUp: (email: string, password: string, displayName?: string) => Promise<UserCredential>;
  signIn: (email: string, password: string) => Promise<UserCredential>;
  signInWithGoogle: () => Promise<UserCredential>;
  signInWithGithub: () => Promise<UserCredential>;
  signOut: () => Promise<void>;

  // Phone authentication
  signInWithPhone: (phoneNumber: string, verifier: ApplicationVerifier) => Promise<any>;
  verifyPhoneCode: (verificationId: string, code: string) => Promise<UserCredential>;
  linkPhoneNumber: (phoneNumber: string, verifier: ApplicationVerifier) => Promise<any>;

  // Multi-Factor Authentication
  enrollMFA: (phoneNumber: string, verifier: ApplicationVerifier) => Promise<void>;
  unenrollMFA: (factorUid: string) => Promise<void>;
  verifyMFACode: (verificationId: string, code: string) => Promise<UserCredential>;
  getMFAInfo: () => any[];

  // Account linking
  linkGoogleAccount: () => Promise<void>;
  linkGithubAccount: () => Promise<void>;
  unlinkProvider: (providerId: string) => Promise<void>;
  getLinkedProviders: () => string[];

  // Session management
  sessions: SessionInfo[];
  getCurrentSession: () => SessionInfo | null;
  getAllSessions: () => Promise<SessionInfo[]>;
  revokeSession: (sessionId: string) => Promise<void>;
  revokeAllSessions: () => Promise<void>;

  // Security & Activity
  loginHistory: LoginHistory[];
  getLoginHistory: (limit?: number) => Promise<LoginHistory[]>;
  logActivity: (activity: string, metadata?: any) => Promise<void>;
  checkAccountSecurity: () => Promise<{ score: number; recommendations: string[] }>;

  // Admin functions (admin role only)
  adminGetAllUsers: () => Promise<UserProfile[]>;
  adminUpdateUserRole: (userId: string, role: 'student' | 'instructor' | 'admin') => Promise<void>;
  adminDisableUser: (userId: string) => Promise<void>;
  adminDeleteUser: (userId: string) => Promise<void>;
}
```

### Extended User Profile

```typescript
interface UserProfile {
  uid: string;
  email: string | null;
  displayName: string | null;
  photoURL: string | null;
  emailVerified: boolean;
  phoneNumber: string | null;
  createdAt: Date;
  lastLoginAt: Date;
  role: 'student' | 'instructor' | 'admin';
  customClaims?: Record<string, any>;
  security?: SecuritySettings;
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
    twoFactorEnabled?: boolean;
  };
  metadata?: {
    creationTime: string;
    lastSignInTime: string;
    lastRefreshTime?: string;
  };
}

interface SecuritySettings {
  mfaEnabled: boolean;
  phoneNumberVerified: boolean;
  linkedProviders: string[];
  lastPasswordChange?: Date;
  accountLockout: boolean;
  failedLoginAttempts: number;
}

interface LoginHistory {
  timestamp: Date;
  ipAddress?: string;
  userAgent: string;
  location?: string;
  success: boolean;
  method: 'email' | 'google' | 'github' | 'phone';
}

interface SessionInfo {
  sessionId: string;
  createdAt: Date;
  lastActivity: Date;
  deviceInfo: {
    userAgent: string;
    platform: string;
    browser: string;
  };
  ipAddress?: string;
  location?: string;
}
```

## Security Considerations

### Best Practices

1. **Always use HTTPS** in production
2. **Enable MFA for admins** - Require 2FA for all admin accounts
3. **Rotate credentials** - Encourage password changes periodically
4. **Monitor login activity** - Alert users of suspicious logins
5. **Rate limiting** - Implement on backend to prevent brute force
6. **Audit logging** - Log all administrative actions
7. **Session timeout** - Auto-logout after inactivity
8. **IP whitelisting** - For admin access (optional)

### Common Vulnerabilities to Avoid

❌ **Storing sensitive data in Firestore without encryption**
✅ Use Firebase Security Rules and encrypt sensitive fields

❌ **Trusting client-side role checks**
✅ Always verify roles in Security Rules and backend

❌ **Not validating phone numbers**
✅ Use proper phone number validation (libphonenumber)

❌ **Exposing admin endpoints**
✅ Use Firebase callable functions with auth checks

❌ **Not rate limiting auth attempts**
✅ Implement exponential backoff and account lockout

## Troubleshooting

### Phone Auth Issues

**"reCAPTCHA verification failed"**
- Ensure reCAPTCHA is properly configured
- Check that domain is authorized in reCAPTCHA settings
- Try visible reCAPTCHA instead of invisible

**"Phone number already in use"**
- Phone numbers can only be linked to one account
- User must unlink from other account first

**"Code expired"**
- Verification codes expire after 5 minutes
- Prompt user to request a new code

### MFA Issues

**"Multi-factor auth required but not enrolled"**
- User enabled MFA but didn't complete enrollment
- Guide user through enrollment process again

**"Cannot enroll more than X factors"**
- Firebase limits number of enrolled factors
- Remove old factors before adding new ones

### Session Management Issues

**"Sessions not persisting"**
- Check Firestore Security Rules for sessions collection
- Ensure session data is being written correctly

**"Cannot revoke current session"**
- Current session revocation logs user out
- Implement confirmation dialog

## Next Steps

1. **Implement Backend Admin API** - For user deletion and bulk operations
2. **Add WebAuthn/Passkeys** - Hardware key support
3. **Biometric Authentication** - Face ID, Touch ID
4. **Risk-based Authentication** - Adaptive security based on context
5. **SSO Integration** - SAML, OIDC for enterprise
6. **Audit Logging** - Comprehensive audit trail
7. **IP Geolocation** - Detect suspicious locations
8. **Device Fingerprinting** - Enhanced fraud detection

## Support

- Firebase Console: https://console.firebase.google.com/project/go-pro-platform
- Firebase Auth Docs: https://firebase.google.com/docs/auth
- Firebase Admin SDK: https://firebase.google.com/docs/admin/setup
- Report Issues: See project documentation
