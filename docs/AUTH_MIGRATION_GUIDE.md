# Authentication System Migration Guide

This guide helps you migrate from the basic authentication system to the advanced authentication system with MFA, session management, and admin dashboard.

## Migration Options

You have two options:

### Option 1: Keep Both Systems (Recommended)

Keep the basic auth system and gradually migrate to advanced features:
- Less risky
- Test advanced features without breaking existing auth
- Migrate users gradually

### Option 2: Full Migration

Replace the basic auth context entirely:
- All features immediately available
- Single codebase to maintain
- Requires more testing

## Option 1: Gradual Migration (Recommended)

### Step 1: Add Advanced Context Alongside Basic

Keep both auth contexts available:

```tsx
// app/layout.tsx
import { AuthProvider as BasicAuthProvider } from '@/contexts/auth-context';
import { AuthProvider as AdvancedAuthProvider } from '@/contexts/auth-context-advanced';

export default function RootLayout({ children }) {
  return (
    <html lang="en">
      <body>
        {/* Use basic auth for now */}
        <BasicAuthProvider>
          {children}
        </BasicAuthProvider>
      </body>
    </html>
  );
}
```

### Step 2: Add Advanced Features to Specific Pages

Create new pages using advanced auth:

```tsx
// app/settings/security/page.tsx
import { AuthProvider as AdvancedAuthProvider } from '@/contexts/auth-context-advanced';
import { SecuritySettings } from '@/components/auth/security-settings';

export default function SecurityPage() {
  return (
    <AdvancedAuthProvider>
      <SecuritySettings />
    </AdvancedAuthProvider>
  );
}
```

### Step 3: Test Advanced Features

Test new features without affecting existing functionality:
- Security dashboard
- Session management
- Login history
- Admin dashboard (if you're an admin)

### Step 4: Migrate Page by Page

Once tested, migrate existing pages:

```tsx
// Before
import { useAuth } from '@/contexts/auth-context';

// After
import { useAuth } from '@/contexts/auth-context-advanced';
// API is backward compatible!
```

### Step 5: Complete Migration

Once all pages are migrated, switch the root layout:

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

## Option 2: Full Migration

### Step 1: Update Root Layout

```tsx
// app/layout.tsx
- import { AuthProvider } from '@/contexts/auth-context';
+ import { AuthProvider } from '@/contexts/auth-context-advanced';

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

### Step 2: Update All Imports

Find and replace across your codebase:

```bash
# Find all uses of basic auth context
grep -r "from '@/contexts/auth-context'" frontend/src

# Replace with advanced context
# Use your IDE's find-and-replace feature
```

### Step 3: Test All Auth Flows

Test every authentication flow:
- ✅ Sign up
- ✅ Sign in (email/password)
- ✅ Sign in with Google
- ✅ Sign in with GitHub
- ✅ Sign out
- ✅ Password reset
- ✅ Email verification
- ✅ Profile updates
- ✅ Protected routes

### Step 4: Remove Old Context

Once migration is complete and tested:

```bash
rm frontend/src/contexts/auth-context.tsx
```

## Backward Compatibility

The advanced auth context is **100% backward compatible** with the basic context. All these work identically:

```tsx
const {
  user,
  userProfile,
  loading,
  error,
  signUp,
  signIn,
  signInWithGoogle,
  signInWithGithub,
  signOut,
  updateUserProfile,
  updateUserEmail,
  updateUserPassword,
  sendPasswordReset,
  sendVerificationEmail,
  updateProgress,
  enrollInCourse,
} = useAuth();
```

New features are additive:

```tsx
const {
  // NEW: Phone authentication
  signInWithPhone,
  verifyPhoneCode,
  linkPhoneNumber,

  // NEW: Multi-Factor Authentication
  enrollMFA,
  unenrollMFA,
  verifyMFACode,
  getMFAInfo,

  // NEW: Session management
  sessions,
  getCurrentSession,
  getAllSessions,
  revokeSession,
  revokeAllSessions,

  // NEW: Security & Activity
  loginHistory,
  getLoginHistory,
  logActivity,
  checkAccountSecurity,

  // NEW: Account linking
  linkGoogleAccount,
  linkGithubAccount,
  unlinkProvider,
  getLinkedProviders,

  // NEW: Admin functions
  adminGetAllUsers,
  adminUpdateUserRole,
  adminDisableUser,
  adminDeleteUser,
} = useAuth();
```

## Data Migration

### User Profiles

The advanced system extends the existing user profile structure. Existing profiles automatically work with new features.

**No data migration required!** Existing data structure:

```typescript
// Basic profile (still supported)
{
  uid: string;
  email: string;
  displayName: string;
  role: 'student' | 'instructor' | 'admin';
  progress: { ... };
}

// Advanced profile (extended)
{
  uid: string;
  email: string;
  displayName: string;
  role: 'student' | 'instructor' | 'admin';
  progress: { ... },
  // NEW FIELDS (optional, added automatically)
  phoneNumber?: string;
  security?: {
    mfaEnabled: boolean;
    linkedProviders: string[];
    ...
  };
  metadata?: { ... };
}
```

### Security Rules

Update Firestore security rules to support new features:

```bash
# Deploy updated rules
firebase deploy --only firestore:rules
```

The updated rules in `firestore.rules` already support all features.

## New Features Setup

### 1. Phone Authentication

Enable in Firebase Console:

```
https://console.firebase.google.com/project/go-pro-platform/authentication/providers
→ Phone → Enable
```

### 2. reCAPTCHA (for Phone Auth)

Get site key from:
```
https://www.google.com/recaptcha/admin/create
```

Add to environment:
```bash
# .env.local
NEXT_PUBLIC_RECAPTCHA_SITE_KEY=your-site-key
```

### 3. Social Auth Providers

Enable Google and GitHub if not already enabled:

```
https://console.firebase.google.com/project/go-pro-platform/authentication/providers
```

## Testing Checklist

### Basic Auth Features (Must Work)
- [ ] Sign up with email/password
- [ ] Sign in with email/password
- [ ] Sign in with Google
- [ ] Sign in with GitHub
- [ ] Sign out
- [ ] Password reset email
- [ ] Email verification
- [ ] Update profile (name, photo)
- [ ] Update email
- [ ] Update password
- [ ] Protected routes
- [ ] Role-based access
- [ ] Progress tracking
- [ ] Course enrollment

### Advanced Features (New)
- [ ] View security dashboard
- [ ] Check security score
- [ ] View login history
- [ ] View active sessions
- [ ] Revoke session
- [ ] Change password from settings
- [ ] View linked providers
- [ ] Phone sign-in (if enabled)
- [ ] MFA enrollment (if enabled)
- [ ] Admin dashboard (if admin)
- [ ] User role management (if admin)

## Rollback Plan

If issues arise, quickly rollback:

### Step 1: Revert Root Layout

```tsx
// app/layout.tsx
- import { AuthProvider } from '@/contexts/auth-context-advanced';
+ import { AuthProvider } from '@/contexts/auth-context';
```

### Step 2: Revert Import Changes

```bash
# If you did find-replace, use git to revert
git checkout -- frontend/src/app
git checkout -- frontend/src/components
```

### Step 3: Keep Advanced Features (Optional)

You can keep advanced features in separate pages:

```tsx
// app/settings/security/page.tsx
import { AuthProvider as AdvancedAuthProvider } from '@/contexts/auth-context-advanced';

export default function SecurityPage() {
  return (
    <AdvancedAuthProvider>
      {/* Advanced features work here */}
    </AdvancedAuthProvider>
  );
}
```

## Common Issues & Solutions

### Issue: TypeScript Errors

**Problem**: TypeScript complains about missing properties

**Solution**: The advanced context has all the same properties as the basic context. Check:
1. Import path is correct
2. Context is wrapped properly
3. TypeScript cache is cleared: `rm -rf .next`

### Issue: User Profile Missing New Fields

**Problem**: `userProfile.security` is undefined

**Solution**: New fields are optional and added when user logs in. First login after migration will populate them.

### Issue: Sessions Not Showing

**Problem**: `sessions` array is empty

**Solution**: Sessions are tracked after migration. Old sessions won't appear. They'll populate on next login.

### Issue: Login History Empty

**Problem**: No login history displayed

**Solution**: History starts tracking after migration. Users need to log in again to generate history.

### Issue: Admin Dashboard Access Denied

**Problem**: Can't access admin dashboard despite being admin

**Solution**:
1. Check `userProfile.role === 'admin'`
2. Verify in Firestore Console: `users/{uid}/role`
3. If not admin, use Firebase Console to update role
4. Sign out and sign back in

## Performance Considerations

The advanced system adds minimal overhead:

- **Bundle size**: +15KB (gzipped)
- **Initial load**: No noticeable difference
- **Firestore reads**: +1 read per login (for login history)
- **Memory**: +~500KB (for session tracking)

Optimize if needed:
- Lazy load admin dashboard
- Paginate login history
- Limit session tracking
- Cache security score

## Best Practices

1. **Test in development first**
   - Use Firebase emulators
   - Test all auth flows
   - Check Firestore rules

2. **Communicate with users**
   - Announce new security features
   - Encourage MFA enrollment
   - Highlight security dashboard

3. **Monitor after migration**
   - Watch error rates
   - Check login success rates
   - Monitor Firestore usage

4. **Gradual feature rollout**
   - Start with security dashboard
   - Add MFA later
   - Enable admin features last

5. **Document for your team**
   - Update onboarding docs
   - Train support team
   - Create user guides

## Support

If you encounter issues:

1. **Check the documentation**
   - `FIREBASE_AUTH.md` - Basic auth
   - `FIREBASE_AUTH_ADVANCED.md` - Advanced features
   - This guide - Migration

2. **Review logs**
   - Browser console
   - Firebase Console
   - Firestore rules debugger

3. **Test in isolation**
   - Create a test page
   - Use advanced context only there
   - Debug without affecting main app

4. **Ask for help**
   - Firebase support
   - Project team
   - Community forums

## Timeline Recommendation

### Week 1: Preparation
- Review documentation
- Set up Firebase emulators
- Test advanced context in isolation

### Week 2: Testing
- Add advanced context to test pages
- Test all new features
- Get feedback from beta users

### Week 3: Migration
- Migrate page by page (Option 1)
- OR: Full migration (Option 2)
- Monitor for issues

### Week 4: Optimization
- Enable additional features (MFA, phone auth)
- Optimize performance
- Document learnings

## Success Metrics

Track these metrics to measure success:

- **Auth success rate**: Should remain 95%+
- **User complaints**: Should not increase
- **Page load time**: Should not increase >100ms
- **Error rates**: Should remain stable
- **Security score**: Average should be >70
- **MFA adoption**: Target 30% in 3 months

## Conclusion

The advanced authentication system provides enterprise-grade security while maintaining full backward compatibility. Whether you choose gradual migration or full migration, your existing auth flows will continue to work seamlessly.

Take your time, test thoroughly, and enjoy the enhanced security features!
