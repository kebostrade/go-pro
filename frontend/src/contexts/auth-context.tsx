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
  signOut as firebaseSignOut,
  onAuthStateChanged,
  updateProfile,
  sendPasswordResetEmail,
  sendEmailVerification,
  updateEmail,
  updatePassword,
  reauthenticateWithCredential,
  EmailAuthProvider,
} from 'firebase/auth';
import { doc, setDoc, getDoc, updateDoc, serverTimestamp } from 'firebase/firestore';
import { auth, db } from '@/lib/firebase';

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
  loading: boolean;
  error: string | null;

  // Authentication methods
  signUp: (email: string, password: string, displayName?: string) => Promise<UserCredential>;
  signIn: (email: string, password: string) => Promise<UserCredential>;
  signInWithGoogle: () => Promise<UserCredential>;
  signInWithGithub: () => Promise<UserCredential>;
  signOut: () => Promise<void>;

  // Profile management
  updateUserProfile: (data: Partial<UserProfile>) => Promise<void>;
  updateUserEmail: (newEmail: string, password: string) => Promise<void>;
  updateUserPassword: (currentPassword: string, newPassword: string) => Promise<void>;
  sendPasswordReset: (email: string) => Promise<void>;
  sendVerificationEmail: () => Promise<void>;

  // Progress tracking
  updateProgress: (lessonId: string, completed: boolean) => Promise<void>;
  enrollInCourse: (courseId: string) => Promise<void>;
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

  // Create or update user profile in Firestore
  const createUserProfile = async (user: User, additionalData?: Partial<UserProfile>) => {
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

  // Load user profile from Firestore
  const loadUserProfile = async (user: User) => {
    try {
      const profile = await createUserProfile(user);
      setUserProfile(profile);
    } catch (err) {
      console.error('Error loading user profile:', err);
      setError('Failed to load user profile');
    }
  };

  // Sign up with email and password
  const signUp = async (email: string, password: string, displayName?: string): Promise<UserCredential> => {
    try {
      setError(null);
      const result = await createUserWithEmailAndPassword(auth, email, password);

      if (displayName && result.user) {
        await updateProfile(result.user, { displayName });
      }

      await sendEmailVerification(result.user);
      await createUserProfile(result.user, { displayName });

      return result;
    } catch (err: any) {
      setError(err.message);
      throw err;
    }
  };

  // Sign in with email and password
  const signIn = async (email: string, password: string): Promise<UserCredential> => {
    try {
      setError(null);
      const result = await signInWithEmailAndPassword(auth, email, password);
      return result;
    } catch (err: any) {
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
      return result;
    } catch (err: any) {
      setError(err.message);
      throw err;
    }
  };

  // Sign out
  const signOut = async (): Promise<void> => {
    try {
      setError(null);
      await firebaseSignOut(auth);
      setUserProfile(null);
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
      await sendPasswordResetEmail(auth, email);
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

  // Update learning progress
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
  };

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
}
