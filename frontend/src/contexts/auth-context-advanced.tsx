'use client';

import React, { createContext, useContext, useEffect, useState } from 'react';
import {
  User,
  UserCredential,
  createUserWithEmailAndPassword,
  signInWithEmailAndPassword,
  signInWithPopup,
  signInWithPhoneNumber,
  GoogleAuthProvider,
  GithubAuthProvider,
  signOut as firebaseSignOut,
  onAuthStateChanged,
  updateProfile,
  sendPasswordResetEmail,
  sendEmailVerification,
  updateEmail,
  updatePassword,
  reauthenticateWithCredential,
  EmailAuthProvider,
  multiFactor,
  PhoneAuthProvider,
  PhoneMultiFactorGenerator,
  RecaptchaVerifier,
  ApplicationVerifier,
  linkWithCredential,
  unlink,
  fetchSignInMethodsForEmail,
} from 'firebase/auth';
import {
  doc,
  setDoc,
  getDoc,
  updateDoc,
  serverTimestamp,
  collection,
  addDoc,
  query,
  where,
  orderBy,
  limit as firestoreLimit,
  getDocs,
  Timestamp,
} from 'firebase/firestore';
import { auth, db } from '@/lib/firebase';

export interface LoginHistory {
  timestamp: Date;
  ipAddress?: string;
  userAgent: string;
  location?: string;
  success: boolean;
  method: 'email' | 'google' | 'github' | 'phone';
}

export interface SecuritySettings {
  mfaEnabled: boolean;
  phoneNumberVerified: boolean;
  linkedProviders: string[];
  lastPasswordChange?: Date;
  accountLockout: boolean;
  failedLoginAttempts: number;
}

export interface SessionInfo {
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

export interface UserProfile {
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

interface AuthContextType {
  user: User | null;
  userProfile: UserProfile | null;
  loading: boolean;
  error: string | null;
  sessions: SessionInfo[];
  loginHistory: LoginHistory[];

  // Basic authentication
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

  // Profile management
  updateUserProfile: (data: Partial<UserProfile>) => Promise<void>;
  updateUserEmail: (newEmail: string, password: string) => Promise<void>;
  updateUserPassword: (currentPassword: string, newPassword: string) => Promise<void>;
  sendPasswordReset: (email: string) => Promise<void>;
  sendVerificationEmail: () => Promise<void>;

  // Account linking
  linkGoogleAccount: () => Promise<void>;
  linkGithubAccount: () => Promise<void>;
  unlinkProvider: (providerId: string) => Promise<void>;
  getLinkedProviders: () => string[];

  // Session management
  getCurrentSession: () => SessionInfo | null;
  getAllSessions: () => Promise<SessionInfo[]>;
  revokeSession: (sessionId: string) => Promise<void>;
  revokeAllSessions: () => Promise<void>;

  // Security & Activity
  getLoginHistory: (limit?: number) => Promise<LoginHistory[]>;
  logActivity: (activity: string, metadata?: any) => Promise<void>;
  checkAccountSecurity: () => Promise<{ score: number; recommendations: string[] }>;

  // Progress tracking
  updateProgress: (lessonId: string, completed: boolean) => Promise<void>;
  enrollInCourse: (courseId: string) => Promise<void>;

  // Admin functions (admin role only)
  adminGetAllUsers: () => Promise<UserProfile[]>;
  adminUpdateUserRole: (userId: string, role: 'student' | 'instructor' | 'admin') => Promise<void>;
  adminDisableUser: (userId: string) => Promise<void>;
  adminDeleteUser: (userId: string) => Promise<void>;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export function useAuth() {
  const context = useContext(AuthContext);
  if (context === undefined) {
    throw new Error('useAuth must be used within an AuthProvider');
  }
  return context;
}

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [user, setUser] = useState<User | null>(null);
  const [userProfile, setUserProfile] = useState<UserProfile | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [sessions, setSessions] = useState<SessionInfo[]>([]);
  const [loginHistory, setLoginHistory] = useState<LoginHistory[]>([]);

  // Device info helper
  const getDeviceInfo = () => {
    const ua = navigator.userAgent;
    return {
      userAgent: ua,
      platform: navigator.platform,
      browser: getBrowserInfo(ua),
    };
  };

  const getBrowserInfo = (ua: string): string => {
    if (ua.includes('Firefox')) return 'Firefox';
    if (ua.includes('Chrome')) return 'Chrome';
    if (ua.includes('Safari')) return 'Safari';
    if (ua.includes('Edge')) return 'Edge';
    return 'Unknown';
  };

  // Create or update user profile
  const createUserProfile = async (user: User, additionalData?: Partial<UserProfile>) => {
    const userRef = doc(db, 'users', user.uid);
    const userSnap = await getDoc(userRef);

    const securitySettings: SecuritySettings = {
      mfaEnabled: multiFactor(user).enrolledFactors.length > 0,
      phoneNumberVerified: !!user.phoneNumber,
      linkedProviders: user.providerData.map((p) => p.providerId),
      accountLockout: false,
      failedLoginAttempts: 0,
    };

    if (!userSnap.exists()) {
      const profile: UserProfile = {
        uid: user.uid,
        email: user.email,
        displayName: user.displayName,
        photoURL: user.photoURL,
        emailVerified: user.emailVerified,
        phoneNumber: user.phoneNumber,
        createdAt: new Date(),
        lastLoginAt: new Date(),
        role: 'student',
        security: securitySettings,
        progress: {
          completedLessons: [],
          xp: 0,
          level: 1,
        },
        preferences: {
          theme: 'system',
          notifications: true,
          language: 'en',
          twoFactorEnabled: false,
        },
        metadata: {
          creationTime: user.metadata.creationTime || '',
          lastSignInTime: user.metadata.lastSignInTime || '',
        },
        ...additionalData,
      };

      await setDoc(userRef, {
        ...profile,
        createdAt: serverTimestamp(),
        lastLoginAt: serverTimestamp(),
      });

      return profile;
    } else {
      await updateDoc(userRef, {
        lastLoginAt: serverTimestamp(),
        'security.mfaEnabled': securitySettings.mfaEnabled,
        'security.linkedProviders': securitySettings.linkedProviders,
      });

      return userSnap.data() as UserProfile;
    }
  };

  // Log login activity
  const logLoginActivity = async (
    userId: string,
    success: boolean,
    method: 'email' | 'google' | 'github' | 'phone'
  ) => {
    try {
      const activityRef = collection(db, 'users', userId, 'loginHistory');
      await addDoc(activityRef, {
        timestamp: serverTimestamp(),
        userAgent: navigator.userAgent,
        success,
        method,
        deviceInfo: getDeviceInfo(),
      });
    } catch (err) {
      console.error('Failed to log login activity:', err);
    }
  };

  // Load user profile
  const loadUserProfile = async (user: User) => {
    try {
      const profile = await createUserProfile(user);
      setUserProfile(profile);
      await logLoginActivity(user.uid, true, 'email');
      await loadLoginHistory(user.uid);
    } catch (err) {
      console.error('Error loading user profile:', err);
      setError('Failed to load user profile');
    }
  };

  // Load login history
  const loadLoginHistory = async (userId: string, historyLimit: number = 10) => {
    try {
      const historyRef = collection(db, 'users', userId, 'loginHistory');
      const q = query(historyRef, orderBy('timestamp', 'desc'), firestoreLimit(historyLimit));
      const snapshot = await getDocs(q);

      const history = snapshot.docs.map((doc) => {
        const data = doc.data();
        return {
          ...data,
          timestamp: data.timestamp?.toDate() || new Date(),
        } as LoginHistory;
      });

      setLoginHistory(history);
    } catch (err) {
      console.error('Failed to load login history:', err);
    }
  };

  // Sign up
  const signUp = async (email: string, password: string, displayName?: string): Promise<UserCredential> => {
    try {
      setError(null);
      const result = await createUserWithEmailAndPassword(auth, email, password);

      if (displayName && result.user) {
        await updateProfile(result.user, { displayName });
      }

      await sendEmailVerification(result.user);
      await createUserProfile(result.user, { displayName });
      await logLoginActivity(result.user.uid, true, 'email');

      return result;
    } catch (err: any) {
      setError(err.message);
      throw err;
    }
  };

  // Sign in
  const signIn = async (email: string, password: string): Promise<UserCredential> => {
    try {
      setError(null);
      const result = await signInWithEmailAndPassword(auth, email, password);
      await logLoginActivity(result.user.uid, true, 'email');
      return result;
    } catch (err: any) {
      if (user) {
        await logLoginActivity(user.uid, false, 'email');
      }
      setError(err.message);
      throw err;
    }
  };

  // Sign in with Google
  const signInWithGoogle = async (): Promise<UserCredential> => {
    try {
      setError(null);
      const provider = new GoogleAuthProvider();
      provider.addScope('profile');
      provider.addScope('email');
      const result = await signInWithPopup(auth, provider);
      await createUserProfile(result.user);
      await logLoginActivity(result.user.uid, true, 'google');
      return result;
    } catch (err: any) {
      setError(err.message);
      throw err;
    }
  };

  // Sign in with GitHub
  const signInWithGithub = async (): Promise<UserCredential> => {
    try {
      setError(null);
      const provider = new GithubAuthProvider();
      provider.addScope('read:user');
      provider.addScope('user:email');
      const result = await signInWithPopup(auth, provider);
      await createUserProfile(result.user);
      await logLoginActivity(result.user.uid, true, 'github');
      return result;
    } catch (err: any) {
      setError(err.message);
      throw err;
    }
  };

  // Phone authentication
  const signInWithPhone = async (phoneNumber: string, verifier: ApplicationVerifier) => {
    try {
      setError(null);
      const confirmationResult = await signInWithPhoneNumber(auth, phoneNumber, verifier);
      return confirmationResult;
    } catch (err: any) {
      setError(err.message);
      throw err;
    }
  };

  const verifyPhoneCode = async (verificationId: string, code: string): Promise<UserCredential> => {
    try {
      const credential = PhoneAuthProvider.credential(verificationId, code);
      // This would need to be completed with signInWithCredential
      throw new Error('Implementation needed based on your flow');
    } catch (err: any) {
      setError(err.message);
      throw err;
    }
  };

  const linkPhoneNumber = async (phoneNumber: string, verifier: ApplicationVerifier) => {
    if (!user) throw new Error('No user logged in');
    try {
      const provider = new PhoneAuthProvider(auth);
      const verificationId = await provider.verifyPhoneNumber(phoneNumber, verifier);
      return verificationId;
    } catch (err: any) {
      setError(err.message);
      throw err;
    }
  };

  // Multi-Factor Authentication
  const enrollMFA = async (phoneNumber: string, verifier: ApplicationVerifier): Promise<void> => {
    if (!user) throw new Error('No user logged in');
    try {
      const session = await multiFactor(user).getSession();
      const phoneAuthProvider = new PhoneAuthProvider(auth);
      const verificationId = await phoneAuthProvider.verifyPhoneNumber(
        { phoneNumber, session },
        verifier
      );
      // Store verificationId for the next step
      return;
    } catch (err: any) {
      setError(err.message);
      throw err;
    }
  };

  const unenrollMFA = async (factorUid: string): Promise<void> => {
    if (!user) throw new Error('No user logged in');
    try {
      const enrolledFactors = multiFactor(user).enrolledFactors;
      const factor = enrolledFactors.find((f) => f.uid === factorUid);
      if (factor) {
        await multiFactor(user).unenroll(factor);
        await updateUserProfile({ 'preferences.twoFactorEnabled': false } as any);
      }
    } catch (err: any) {
      setError(err.message);
      throw err;
    }
  };

  const verifyMFACode = async (verificationId: string, code: string): Promise<UserCredential> => {
    try {
      const cred = PhoneAuthProvider.credential(verificationId, code);
      const multiFactorAssertion = PhoneMultiFactorGenerator.assertion(cred);
      // Complete MFA sign-in
      throw new Error('Implementation needed based on your MFA flow');
    } catch (err: any) {
      setError(err.message);
      throw err;
    }
  };

  const getMFAInfo = () => {
    if (!user) return [];
    return multiFactor(user).enrolledFactors;
  };

  // Account linking
  const linkGoogleAccount = async (): Promise<void> => {
    if (!user) throw new Error('No user logged in');
    try {
      const provider = new GoogleAuthProvider();
      const result = await linkWithCredential(user, provider.credential(null as any, null as any));
      await updateUserProfile({});
    } catch (err: any) {
      setError(err.message);
      throw err;
    }
  };

  const linkGithubAccount = async (): Promise<void> => {
    if (!user) throw new Error('No user logged in');
    try {
      const provider = new GithubAuthProvider();
      // Implementation would use linkWithPopup
      throw new Error('Use linkWithPopup in actual implementation');
    } catch (err: any) {
      setError(err.message);
      throw err;
    }
  };

  const unlinkProvider = async (providerId: string): Promise<void> => {
    if (!user) throw new Error('No user logged in');
    try {
      await unlink(user, providerId);
      await updateUserProfile({});
    } catch (err: any) {
      setError(err.message);
      throw err;
    }
  };

  const getLinkedProviders = (): string[] => {
    if (!user) return [];
    return user.providerData.map((p) => p.providerId);
  };

  // Session management
  const getCurrentSession = (): SessionInfo | null => {
    if (!user) return null;
    return {
      sessionId: user.uid + '-' + Date.now(),
      createdAt: new Date(),
      lastActivity: new Date(),
      deviceInfo: getDeviceInfo(),
    };
  };

  const getAllSessions = async (): Promise<SessionInfo[]> => {
    // Implementation would query Firestore for active sessions
    return sessions;
  };

  const revokeSession = async (sessionId: string): Promise<void> => {
    // Implementation would remove session from Firestore
    setSessions((prev) => prev.filter((s) => s.sessionId !== sessionId));
  };

  const revokeAllSessions = async (): Promise<void> => {
    if (!user) throw new Error('No user logged in');
    // This would sign out from all devices
    await firebaseSignOut(auth);
    setSessions([]);
  };

  // Security & Activity
  const getLoginHistory = async (historyLimit: number = 20): Promise<LoginHistory[]> => {
    if (!user) return [];
    await loadLoginHistory(user.uid, historyLimit);
    return loginHistory;
  };

  const logActivity = async (activity: string, metadata?: any): Promise<void> => {
    if (!user) throw new Error('No user logged in');
    try {
      const activityRef = collection(db, 'users', user.uid, 'activity');
      await addDoc(activityRef, {
        activity,
        metadata,
        timestamp: serverTimestamp(),
        deviceInfo: getDeviceInfo(),
      });
    } catch (err: any) {
      console.error('Failed to log activity:', err);
    }
  };

  const checkAccountSecurity = async (): Promise<{ score: number; recommendations: string[] }> => {
    if (!user || !userProfile) {
      return { score: 0, recommendations: ['Please sign in'] };
    }

    let score = 50;
    const recommendations: string[] = [];

    if (user.emailVerified) {
      score += 15;
    } else {
      recommendations.push('Verify your email address');
    }

    if (userProfile.security?.mfaEnabled) {
      score += 20;
    } else {
      recommendations.push('Enable two-factor authentication');
    }

    if (userProfile.security?.linkedProviders && userProfile.security.linkedProviders.length > 1) {
      score += 10;
    } else {
      recommendations.push('Link additional authentication methods');
    }

    if (user.phoneNumber) {
      score += 5;
    } else {
      recommendations.push('Add a phone number for account recovery');
    }

    return { score, recommendations };
  };

  // Sign out
  const signOut = async (): Promise<void> => {
    try {
      setError(null);
      await firebaseSignOut(auth);
      setUserProfile(null);
      setSessions([]);
      setLoginHistory([]);
    } catch (err: any) {
      setError(err.message);
      throw err;
    }
  };

  // Update profile
  const updateUserProfile = async (data: Partial<UserProfile>): Promise<void> => {
    if (!user) throw new Error('No user logged in');

    try {
      setError(null);
      const userRef = doc(db, 'users', user.uid);
      await updateDoc(userRef, {
        ...data,
        updatedAt: serverTimestamp(),
      });

      if (data.displayName || data.photoURL) {
        await updateProfile(user, {
          displayName: data.displayName || user.displayName,
          photoURL: data.photoURL || user.photoURL,
        });
      }

      setUserProfile((prev) => (prev ? { ...prev, ...data } : null));
    } catch (err: any) {
      setError(err.message);
      throw err;
    }
  };

  // Update email
  const updateUserEmail = async (newEmail: string, password: string): Promise<void> => {
    if (!user || !user.email) throw new Error('No user logged in');

    try {
      setError(null);
      const credential = EmailAuthProvider.credential(user.email, password);
      await reauthenticateWithCredential(user, credential);
      await updateEmail(user, newEmail);
      await updateUserProfile({ email: newEmail });
    } catch (err: any) {
      setError(err.message);
      throw err;
    }
  };

  // Update password
  const updateUserPassword = async (currentPassword: string, newPassword: string): Promise<void> => {
    if (!user || !user.email) throw new Error('No user logged in');

    try {
      setError(null);
      const credential = EmailAuthProvider.credential(user.email, currentPassword);
      await reauthenticateWithCredential(user, credential);
      await updatePassword(user, newPassword);

      await updateUserProfile({
        security: {
          ...userProfile?.security,
          lastPasswordChange: new Date(),
        } as SecuritySettings,
      });
    } catch (err: any) {
      setError(err.message);
      throw err;
    }
  };

  // Send password reset
  const sendPasswordReset = async (email: string): Promise<void> => {
    try {
      setError(null);
      await sendPasswordResetEmail(auth, email);
    } catch (err: any) {
      setError(err.message);
      throw err;
    }
  };

  // Send verification email
  const sendVerificationEmail = async (): Promise<void> => {
    if (!user) throw new Error('No user logged in');

    try {
      setError(null);
      await sendEmailVerification(user);
    } catch (err: any) {
      setError(err.message);
      throw err;
    }
  };

  // Progress tracking
  const updateProgress = async (lessonId: string, completed: boolean): Promise<void> => {
    if (!user || !userProfile) throw new Error('No user logged in');

    try {
      setError(null);
      const completedLessons = userProfile.progress?.completedLessons || [];
      const updatedLessons = completed
        ? [...new Set([...completedLessons, lessonId])]
        : completedLessons.filter((id) => id !== lessonId);

      const xpGain = completed ? 10 : -10;
      const newXp = Math.max(0, (userProfile.progress?.xp || 0) + xpGain);
      const newLevel = Math.floor(newXp / 100) + 1;

      await updateUserProfile({
        progress: {
          completedLessons: updatedLessons,
          currentCourse: userProfile.progress?.currentCourse,
          xp: newXp,
          level: newLevel,
        },
      });
    } catch (err: any) {
      setError(err.message);
      throw err;
    }
  };

  // Enroll in course
  const enrollInCourse = async (courseId: string): Promise<void> => {
    if (!user || !userProfile) throw new Error('No user logged in');

    try {
      setError(null);
      await updateUserProfile({
        progress: {
          ...userProfile.progress,
          completedLessons: userProfile.progress?.completedLessons || [],
          currentCourse: courseId,
          xp: userProfile.progress?.xp || 0,
          level: userProfile.progress?.level || 1,
        },
      });
    } catch (err: any) {
      setError(err.message);
      throw err;
    }
  };

  // Admin functions
  const adminGetAllUsers = async (): Promise<UserProfile[]> => {
    if (!userProfile || userProfile.role !== 'admin') {
      throw new Error('Admin access required');
    }

    try {
      const usersRef = collection(db, 'users');
      const snapshot = await getDocs(usersRef);
      return snapshot.docs.map((doc) => doc.data() as UserProfile);
    } catch (err: any) {
      setError(err.message);
      throw err;
    }
  };

  const adminUpdateUserRole = async (
    userId: string,
    role: 'student' | 'instructor' | 'admin'
  ): Promise<void> => {
    if (!userProfile || userProfile.role !== 'admin') {
      throw new Error('Admin access required');
    }

    try {
      const userRef = doc(db, 'users', userId);
      await updateDoc(userRef, { role, updatedAt: serverTimestamp() });
    } catch (err: any) {
      setError(err.message);
      throw err;
    }
  };

  const adminDisableUser = async (userId: string): Promise<void> => {
    if (!userProfile || userProfile.role !== 'admin') {
      throw new Error('Admin access required');
    }

    try {
      const userRef = doc(db, 'users', userId);
      await updateDoc(userRef, {
        'security.accountLockout': true,
        updatedAt: serverTimestamp(),
      });
    } catch (err: any) {
      setError(err.message);
      throw err;
    }
  };

  const adminDeleteUser = async (userId: string): Promise<void> => {
    if (!userProfile || userProfile.role !== 'admin') {
      throw new Error('Admin access required');
    }

    // Note: Actual user deletion requires Firebase Admin SDK on backend
    throw new Error('User deletion must be performed via backend admin endpoint');
  };

  // Auth state listener
  useEffect(() => {
    const unsubscribe = onAuthStateChanged(auth, async (user) => {
      setUser(user);
      if (user) {
        await loadUserProfile(user);
      } else {
        setUserProfile(null);
      }
      setLoading(false);
    });

    return unsubscribe;
  }, []);

  const value: AuthContextType = {
    user,
    userProfile,
    loading,
    error,
    sessions,
    loginHistory,
    signUp,
    signIn,
    signInWithGoogle,
    signInWithGithub,
    signInWithPhone,
    verifyPhoneCode,
    linkPhoneNumber,
    enrollMFA,
    unenrollMFA,
    verifyMFACode,
    getMFAInfo,
    signOut,
    updateUserProfile,
    updateUserEmail,
    updateUserPassword,
    sendPasswordReset,
    sendVerificationEmail,
    linkGoogleAccount,
    linkGithubAccount,
    unlinkProvider,
    getLinkedProviders,
    getCurrentSession,
    getAllSessions,
    revokeSession,
    revokeAllSessions,
    getLoginHistory,
    logActivity,
    checkAccountSecurity,
    updateProgress,
    enrollInCourse,
    adminGetAllUsers,
    adminUpdateUserRole,
    adminDisableUser,
    adminDeleteUser,
  };

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
}
