import { initializeApp, getApps, FirebaseApp } from 'firebase/app';
import { getAuth, connectAuthEmulator, Auth } from 'firebase/auth';
import { getFirestore, connectFirestoreEmulator, Firestore } from 'firebase/firestore';
import { getStorage, connectStorageEmulator, FirebaseStorage } from 'firebase/storage';
import { getAnalytics, isSupported } from 'firebase/analytics';

const firebaseConfig = {
  apiKey: process.env.NEXT_PUBLIC_FIREBASE_API_KEY || "",
  authDomain: process.env.NEXT_PUBLIC_FIREBASE_AUTH_DOMAIN || "",
  projectId: process.env.NEXT_PUBLIC_FIREBASE_PROJECT_ID || "",
  storageBucket: process.env.NEXT_PUBLIC_FIREBASE_STORAGE_BUCKET || "",
  messagingSenderId: process.env.NEXT_PUBLIC_FIREBASE_MESSAGING_SENDER_ID || "",
  appId: process.env.NEXT_PUBLIC_FIREBASE_APP_ID || ""
};

// Initialize Firebase with lazy validation
function initializeFirebaseApp(): FirebaseApp {
  // Check if already initialized
  const existingApps = getApps();
  if (existingApps.length > 0) {
    return existingApps[0];
  }

  // Validate config only when actually initializing
  const hasRequiredConfig = firebaseConfig.apiKey && firebaseConfig.projectId && firebaseConfig.appId;
  if (!hasRequiredConfig) {
    console.error('❌ Firebase configuration missing:', {
      hasApiKey: !!firebaseConfig.apiKey,
      hasProjectId: !!firebaseConfig.projectId,
      hasAppId: !!firebaseConfig.appId,
      env: process.env.NODE_ENV
    });
    throw new Error('Firebase not initialized - missing critical configuration (apiKey, projectId, or appId)');
  }

  try {
    const app = initializeApp(firebaseConfig);
    console.log('✅ Firebase initialized successfully');
    return app;
  } catch (error) {
    console.error('❌ Firebase initialization failed:', error);
    throw error;
  }
}

// Lazy initialize Firebase
let app: FirebaseApp | null = null;
let authInstance: Auth | null = null;
let dbInstance: Firestore | null = null;
let storageInstance: FirebaseStorage | null = null;

// Getter for Firebase app - initializes on first access
function getFirebaseApp(): FirebaseApp {
  if (!app) {
    app = initializeFirebaseApp();
  }
  return app;
}

// Export auth - initialize on first access
export function getAuthInstance(): Auth {
  if (!authInstance) {
    authInstance = getAuth(getFirebaseApp());
  }
  return authInstance;
}

// Export db - initialize on first access
export function getDbInstance(): Firestore {
  if (!dbInstance) {
    dbInstance = getFirestore(getFirebaseApp());
  }
  return dbInstance;
}

// Export storage - initialize on first access
export function getStorageInstance(): FirebaseStorage {
  if (!storageInstance) {
    storageInstance = getStorage(getFirebaseApp());
  }
  return storageInstance;
}

// For backward compatibility, provide getters
export const auth = {
  get currentUser() {
    return getAuthInstance().currentUser;
  }
} as Auth;

export const db: Firestore = new Proxy({} as Firestore, {
  get(_target, prop) {
    const instance = getDbInstance();
    return (instance as any)[prop];
  }
}) as any;

export const storage: FirebaseStorage = new Proxy({} as FirebaseStorage, {
  get(_target, prop) {
    const instance = getStorageInstance();
    return (instance as any)[prop];
  }
}) as any;

// Initialize Analytics (only in browser) - lazy
export const analytics = typeof window !== 'undefined'
  ? isSupported().then(yes => yes ? getAnalytics(getFirebaseApp()) : null)
  : Promise.resolve(null);

// Connect to emulators in development - deferred
if (typeof window !== 'undefined') {
  // Wait for first access to initialize emulators
  const initEmulators = () => {
    if (process.env.NODE_ENV === 'development') {
      const useEmulators = process.env.NEXT_PUBLIC_USE_FIREBASE_EMULATORS === 'true';

      if (useEmulators && authInstance && dbInstance && storageInstance) {
        try {
          connectAuthEmulator(authInstance, 'http://localhost:9099', { disableWarnings: true });
          connectFirestoreEmulator(dbInstance, 'localhost', 8080);
          connectStorageEmulator(storageInstance, 'localhost', 9199);
          console.log('✅ Firebase emulators connected');
        } catch (err) {
          console.warn('⚠️ Emulator connection skipped (already connected)');
        }
      }
    }
  };

  // Initialize emulators after a short delay to allow lazy init
  setTimeout(initEmulators, 100);
}

export default getFirebaseApp;
