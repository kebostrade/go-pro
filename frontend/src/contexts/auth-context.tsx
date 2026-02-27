'use client';

import React, { createContext, useContext, useEffect, useState } from 'react';
import {
  User,
  UserCredential,
  createUserWithEmailAndPassword,
  signInWithEmailAndPassword,
  signInWithPopup,
  GoogleAuthProvider,
  GithubAuthProvider,
  RecaptchaVerifier,
  signInWithPhoneNumber,
  ConfirmationResult,
  signOut as firebaseSignOut,
  onAuthStateChanged,
  updateProfile,
  sendPasswordResetEmail,
  sendEmailVerification,
  updateEmail,
  updatePassword,
  reauthenticateWithCredential,
  EmailAuthProvider,
  setPersistence,
  browserLocalPersistence,
  browserSessionPersistence,
  inMemoryPersistence,
} from 'firebase/auth';
import { doc, setDoc, getDoc, updateDoc, serverTimestamp } from 'firebase/firestore';
import { getAuthInstance, getDbInstance } from '@/lib/firebase';
import { api, BackendUser, setTokenRefreshCallback } from '@/lib/api';
import {
  getAuthErrorMessage,
  authRateLimiter,
  SessionPersistence,
  PhoneValidator,
  PasswordValidator,
  EmailValidator,
  AuthAnalytics,
} from '@/lib/auth-utils';

export interface UserProfile {
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

interface AuthContextType {
  user: User | null;
  userProfile: UserProfile | null;
  backendUser: BackendUser | null;
  loading: boolean;
  error: string | null;

  // Authentication methods
  signUp: (email: string, password: string, displayName?: string) => Promise<UserCredential>;
  signIn: (email: string, password: string) => Promise<UserCredential>;
  signInWithGoogle: () => Promise<UserCredential>;
  signInWithGithub: () => Promise<UserCredential>;
  signOut: () => Promise<void>;

  // Phone authentication
  setupRecaptcha: (containerId: string) => RecaptchaVerifier;
  signInWithPhone: (phoneNumber: string, recaptchaVerifier: RecaptchaVerifier) => Promise<ConfirmationResult>;
  verifyPhoneCode: (confirmationResult: ConfirmationResult, code: string) => Promise<UserCredential>;

  // Profile management
  updateUserProfile: (data: Partial<UserProfile>) => Promise<void>;
  updateUserEmail: (newEmail: string, password: string) => Promise<void>;
  updateUserPassword: (currentPassword: string, newPassword: string) => Promise<void>;
  sendPasswordReset: (email: string) => Promise<void>;
  sendVerificationEmail: () => Promise<void>;

  // Session management
  setSessionPersistence: (persistence: SessionPersistence) => Promise<void>;

  // Progress tracking (backend-synced)
  updateProgress: (lessonId: string, completed: boolean) => Promise<void>;
  enrollInCourse: (courseId: string) => Promise<void>;

  // Backend sync
  syncWithBackend: () => Promise<void>;

  // Utilities
  validateEmail: (email: string) => boolean;
  validatePassword: (password: string) => { valid: boolean; message?: string };
  getPasswordStrength: (password: string) => { strength: 'weak' | 'medium' | 'strong'; percentage: number; color: string; label: string };
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
  const [backendUser, setBackendUser] = useState<BackendUser | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  // Backend sync function
  const syncWithBackend = async () => {
    if (!user) return;

    try {
      const idToken = await user.getIdToken(true); // Force refresh
      const response = await api.verifyToken(idToken);
      setBackendUser(response.user);
      console.log('✅ Synced with backend:', response.user);
    } catch (err: any) {
      // Silently skip if backend not configured
      if (err?.type !== 'backend_not_configured') {
        console.error('❌ Backend sync failed:', err);
        setError(err.message || 'Failed to sync with backend');
      }
    }
  };

  // Set up token refresh callback
  useEffect(() => {
    setTokenRefreshCallback(async () => {
      if (user) {
        await user.getIdToken(true); // Force token refresh
      }
    });
  }, [user]);

  // Create or update user profile in Firestore
  const createUserProfile = async (user: User, additionalData?: Partial<UserProfile>) => {
    const db = getDbInstance();
    const userRef = doc(db, 'users', user.uid);
    const userSnap = await getDoc(userRef);

    if (!userSnap.exists()) {
      // Create new user profile
      const profile: UserProfile = {
        uid: user.uid,
        email: user.email,
        displayName: user.displayName,
        photoURL: user.photoURL,
        emailVerified: user.emailVerified,
        createdAt: new Date(),
        lastLoginAt: new Date(),
        role: 'student',
        progress: {
          completedLessons: [],
          xp: 0,
          level: 1,
        },
        preferences: {
          theme: 'system',
          notifications: true,
          language: 'en',
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
      // Update last login
      await updateDoc(userRef, {
        lastLoginAt: serverTimestamp(),
      });

      return userSnap.data() as UserProfile;
    }
  };

  // Load user profile from Firestore and sync with backend
  const loadUserProfile = async (user: User) => {
    // Try Firestore profile (non-blocking)
    try {
      const profile = await createUserProfile(user);
      setUserProfile(profile);
    } catch (firestoreErr) {
      console.warn('⚠️ Firestore profile failed (permissions?):', firestoreErr);
      // Create a basic profile from Firebase Auth data
      setUserProfile({
        uid: user.uid,
        email: user.email,
        displayName: user.displayName,
        photoURL: user.photoURL,
        emailVerified: user.emailVerified,
        createdAt: new Date(),
        lastLoginAt: new Date(),
        role: 'student',
        progress: { completedLessons: [], xp: 0, level: 1 },
        preferences: { theme: 'system', notifications: true, language: 'en' },
      });
    }

    // Try backend sync (non-blocking) - skip silently if backend not configured
    try {
      const idToken = await user.getIdToken();
      const response = await api.verifyToken(idToken);
      setBackendUser(response.user);
      console.log('✅ Backend sync successful');
    } catch (backendErr: any) {
      // Silently skip if backend not configured (production without backend)
      if (backendErr?.type !== 'backend_not_configured') {
        console.warn('⚠️ Backend sync failed:', backendErr);
      }
      // User still authenticated, just no backend sync
    }
  };

  // Sign up with email and password
  const signUp = async (email: string, password: string, displayName?: string): Promise<UserCredential> => {
    try {
      setError(null);

      // Rate limiting check
      const rateLimitCheck = authRateLimiter.check(email);
      if (!rateLimitCheck.allowed) {
        const error = new Error(`Too many signup attempts. Please try again in ${rateLimitCheck.retryAfter} seconds.`);
        setError(error.message);
        AuthAnalytics.trackError(error.message, 'signup_rate_limit');
        throw error;
      }

      // Validate email
      if (!EmailValidator.isValid(email)) {
        const error = new Error('Please enter a valid email address.');
        setError(error.message);
        throw error;
      }

      // Validate password
      const passwordValidation = PasswordValidator.isValid(password);
      if (!passwordValidation.valid) {
        setError(passwordValidation.message || 'Invalid password');
        throw new Error(passwordValidation.message || 'Invalid password');
      }

      const result = await createUserWithEmailAndPassword(getAuthInstance(), email, password);

      if (displayName && result.user) {
        await updateProfile(result.user, { displayName });
      }

      await sendEmailVerification(result.user);
      await createUserProfile(result.user, { displayName });

      // Reset rate limiter on success
      authRateLimiter.reset(email);
      AuthAnalytics.trackSignUp('email');

      return result;
    } catch (err: any) {
      const errorMessage = getAuthErrorMessage(err);
      setError(errorMessage);
      AuthAnalytics.trackError(errorMessage, 'signup');
      throw new Error(errorMessage);
    }
  };

  // Sign in with email and password
  const signIn = async (email: string, password: string): Promise<UserCredential> => {
    try {
      setError(null);

      // Rate limiting check
      const rateLimitCheck = authRateLimiter.check(email);
      if (!rateLimitCheck.allowed) {
        const error = new Error(`Too many signin attempts. Please try again in ${rateLimitCheck.retryAfter} seconds.`);
        setError(error.message);
        AuthAnalytics.trackError(error.message, 'signin_rate_limit');
        throw error;
      }

      const result = await signInWithEmailAndPassword(getAuthInstance(), email, password);

      // Reset rate limiter on success
      authRateLimiter.reset(email);
      AuthAnalytics.trackSignIn('email');

      return result;
    } catch (err: any) {
      const errorMessage = getAuthErrorMessage(err);
      setError(errorMessage);
      AuthAnalytics.trackError(errorMessage, 'signin');
      throw new Error(errorMessage);
    }
  };

  // Sign in with Google
  const signInWithGoogle = async (): Promise<UserCredential> => {
    try {
      setError(null);
      const provider = new GoogleAuthProvider();
      provider.addScope('profile');
      provider.addScope('email');
      const result = await signInWithPopup(getAuthInstance(), provider);
      // Profile creation is non-blocking - don't fail login if Firestore write fails
      try {
        await createUserProfile(result.user);
      } catch (profileErr) {
        console.warn('⚠️ Profile creation failed (Firestore permissions?):', profileErr);
        // Continue - user is still authenticated
      }
      AuthAnalytics.trackSignIn('google');
      return result;
    } catch (err: any) {
      const errorMessage = getAuthErrorMessage(err);
      setError(errorMessage);
      AuthAnalytics.trackError(errorMessage, 'google_signin');
      throw new Error(errorMessage);
    }
  };

  // Sign in with GitHub
  const signInWithGithub = async (): Promise<UserCredential> => {
    try {
      setError(null);
      const provider = new GithubAuthProvider();
      provider.addScope('read:user');
      provider.addScope('user:email');
      const result = await signInWithPopup(getAuthInstance(), provider);
      // Profile creation is non-blocking - don't fail login if Firestore write fails
      try {
        await createUserProfile(result.user);
      } catch (profileErr) {
        console.warn('⚠️ Profile creation failed (Firestore permissions?):', profileErr);
        // Continue - user is still authenticated
      }
      AuthAnalytics.trackSignIn('github');
      return result;
    } catch (err: any) {
      const errorMessage = getAuthErrorMessage(err);
      setError(errorMessage);
      AuthAnalytics.trackError(errorMessage, 'github_signin');
      throw new Error(errorMessage);
    }
  };

  // Setup reCAPTCHA for phone auth
  const setupRecaptcha = (containerId: string): RecaptchaVerifier => {
    const recaptchaVerifier = new RecaptchaVerifier(getAuthInstance(), containerId, {
      size: 'invisible',
      callback: () => {
        console.log('reCAPTCHA solved');
      },
      'expired-callback': () => {
        console.log('reCAPTCHA expired');
      },
    });
    return recaptchaVerifier;
  };

  // Sign in with phone number
  const signInWithPhone = async (
    phoneNumber: string,
    recaptchaVerifier: RecaptchaVerifier
  ): Promise<ConfirmationResult> => {
    try {
      setError(null);

      // Validate phone number
      if (!PhoneValidator.isValid(phoneNumber)) {
        const error = new Error('Please enter a valid phone number in E.164 format (+1234567890).');
        setError(error.message);
        throw error;
      }

      // Rate limiting check
      const rateLimitCheck = authRateLimiter.check(phoneNumber);
      if (!rateLimitCheck.allowed) {
        const error = new Error(`Too many SMS requests. Please try again in ${rateLimitCheck.retryAfter} seconds.`);
        setError(error.message);
        AuthAnalytics.trackError(error.message, 'phone_rate_limit');
        throw error;
      }

      const confirmationResult = await signInWithPhoneNumber(getAuthInstance(), phoneNumber, recaptchaVerifier);
      return confirmationResult;
    } catch (err: any) {
      const errorMessage = getAuthErrorMessage(err);
      setError(errorMessage);
      AuthAnalytics.trackError(errorMessage, 'phone_signin');
      throw new Error(errorMessage);
    }
  };

  // Verify phone code
  const verifyPhoneCode = async (
    confirmationResult: ConfirmationResult,
    code: string
  ): Promise<UserCredential> => {
    try {
      setError(null);
      const result = await confirmationResult.confirm(code);
      await createUserProfile(result.user);
      AuthAnalytics.trackSignIn('phone');
      return result;
    } catch (err: any) {
      const errorMessage = getAuthErrorMessage(err);
      setError(errorMessage);
      AuthAnalytics.trackError(errorMessage, 'phone_verify');
      throw new Error(errorMessage);
    }
  };

  // Set session persistence
  const setSessionPersistence = async (persistence: SessionPersistence): Promise<void> => {
    try {
      setError(null);
      const persistenceMap = {
        [SessionPersistence.LOCAL]: browserLocalPersistence,
        [SessionPersistence.SESSION]: browserSessionPersistence,
        [SessionPersistence.NONE]: inMemoryPersistence,
      };
      await setPersistence(getAuthInstance(), persistenceMap[persistence]);
    } catch (err: any) {
      const errorMessage = getAuthErrorMessage(err);
      setError(errorMessage);
      throw new Error(errorMessage);
    }
  };

  // Sign out
  const signOut = async (): Promise<void> => {
    try {
      setError(null);
      await firebaseSignOut(getAuthInstance());
      setUserProfile(null);
      setBackendUser(null);
    } catch (err: any) {
      setError(err.message);
      throw err;
    }
  };

  // Update user profile
  const updateUserProfile = async (data: Partial<UserProfile>): Promise<void> => {
    if (!user) throw new Error('No user logged in');

    try {
      setError(null);
      const db = getDbInstance();
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

      setUserProfile((prev) => prev ? { ...prev, ...data } : null);
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
    } catch (err: any) {
      setError(err.message);
      throw err;
    }
  };

  // Send password reset email
  const sendPasswordReset = async (email: string): Promise<void> => {
    try {
      setError(null);
      await sendPasswordResetEmail(getAuthInstance(), email);
    } catch (err: any) {
      setError(err.message);
      throw err;
    }
  };

  // Send email verification
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

  // Update learning progress (backend-first, Firestore as cache)
  const updateProgress = async (lessonId: string, completed: boolean): Promise<void> => {
    if (!user || !userProfile || !backendUser) throw new Error('No user logged in');

    try {
      setError(null);

      // Calculate progress updates
      const completedLessons = userProfile.progress?.completedLessons || [];
      const updatedLessons = completed
        ? [...new Set([...completedLessons, lessonId])]
        : completedLessons.filter((id) => id !== lessonId);

      const xpGain = completed ? 10 : -10;
      const newXp = Math.max(0, (userProfile.progress?.xp || 0) + xpGain);
      const newLevel = Math.floor(newXp / 100) + 1;

      // Update backend first (source of truth)
      try {
        await api.updateProgress(backendUser.id, lessonId, {
          completed,
          score: newXp,
        });
        console.log('✅ Progress synced to backend');
      } catch (backendErr) {
        console.warn('⚠️ Backend progress update failed, using Firestore cache:', backendErr);
      }

      // Update Firestore cache for offline support
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

  // Auth state listener
  useEffect(() => {
    const unsubscribe = onAuthStateChanged(getAuthInstance(), async (user) => {
      setUser(user);
      if (user) {
        await loadUserProfile(user);
      } else {
        setUserProfile(null);
        setBackendUser(null);
      }
      setLoading(false);
    });

    return unsubscribe;
  }, []);

  // Periodic backend sync (every 5 minutes)
  useEffect(() => {
    if (!user) return;

    const syncInterval = setInterval(async () => {
      try {
        await syncWithBackend();
      } catch (err) {
        console.warn('Periodic sync failed:', err);
      }
    }, 5 * 60 * 1000); // 5 minutes

    return () => clearInterval(syncInterval);
  }, [user]);

  // Utility methods
  const validateEmail = (email: string): boolean => {
    return EmailValidator.isValid(email);
  };

  const validatePassword = (password: string): { valid: boolean; message?: string } => {
    return PasswordValidator.isValid(password);
  };

  const getPasswordStrength = (password: string) => {
    return PasswordValidator.getStrengthIndicator(password);
  };

  const value: AuthContextType = {
    user,
    userProfile,
    backendUser,
    loading,
    error,
    signUp,
    signIn,
    signInWithGoogle,
    signInWithGithub,
    setupRecaptcha,
    signInWithPhone,
    verifyPhoneCode,
    signOut,
    updateUserProfile,
    updateUserEmail,
    updateUserPassword,
    sendPasswordReset,
    sendVerificationEmail,
    setSessionPersistence,
    updateProgress,
    enrollInCourse,
    syncWithBackend,
    validateEmail,
    validatePassword,
    getPasswordStrength,
  };

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
}
