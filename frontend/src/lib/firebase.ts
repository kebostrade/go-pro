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

let app: FirebaseApp | null = null;
let authInstance: Auth | null = null;
let dbInstance: Firestore | null = null;
let storageInstance: FirebaseStorage | null = null;
let initError: Error | null = null;

function initializeFirebaseApp(): FirebaseApp | null {
  if (initError) {
    throw initError;
  }
  
  if (app) {
    return app;
  }

  const existingApps = getApps();
  if (existingApps.length > 0) {
    app = existingApps[0];
    return app;
  }

  const hasRequiredConfig = firebaseConfig.apiKey && firebaseConfig.projectId && firebaseConfig.appId;
  if (!hasRequiredConfig) {
    console.warn('⚠️ Firebase configuration missing - Firebase features disabled');
    return null;
  }

  try {
    app = initializeApp(firebaseConfig);
    console.log('✅ Firebase initialized successfully');
    return app;
  } catch (error) {
    initError = error as Error;
    console.error('❌ Firebase initialization failed:', error);
    throw error;
  }
}

// Getter for Firebase app - initializes on first access
export function getFirebaseApp(): FirebaseApp {
  const initializedApp = initializeFirebaseApp();
  if (!initializedApp) {
    throw new Error('Firebase not initialized - please add Firebase config to .env.local');
  }
  return initializedApp;
}

// Check if Firebase is available
export function isFirebaseReady(): boolean {
  try {
    return initializeFirebaseApp() !== null;
  } catch {
    return false;
  }
}

// Export auth - initialize on first access
export function getAuthInstance(): Auth {
  if (!authInstance || initError) {
    const firebaseApp = initializeFirebaseApp();
    if (!firebaseApp) {
      throw new Error('Firebase auth not available - check configuration');
    }
    authInstance = getAuth(firebaseApp);
  }
  return authInstance;
}

// Export db - initialize on first access
export function getDbInstance(): Firestore {
  if (!dbInstance || initError) {
    const firebaseApp = initializeFirebaseApp();
    if (!firebaseApp) {
      throw new Error('Firebase db not available - check configuration');
    }
    dbInstance = getFirestore(firebaseApp);
  }
  return dbInstance;
}

// Export storage - initialize on first access
export function getStorageInstance(): FirebaseStorage {
  if (!storageInstance || initError) {
    const firebaseApp = initializeFirebaseApp();
    if (!firebaseApp) {
      throw new Error('Firebase storage not available - check configuration');
    }
    storageInstance = getStorage(firebaseApp);
  }
  return storageInstance;
}

// For backward compatibility, provide getters
export const auth = {
  get currentUser() {
    try {
      return getAuthInstance().currentUser;
    } catch {
      return null;
    }
  }
} as Auth;

export const db: Firestore = new Proxy({} as Firestore, {
  get(_target, prop) {
    try {
      const instance = getDbInstance();
      return (instance as any)[prop];
    } catch {
      return undefined;
    }
  }
}) as any;

export const storage: FirebaseStorage = new Proxy({} as FirebaseStorage, {
  get(_target, prop) {
    try {
      const instance = getStorageInstance();
      return (instance as any)[prop];
    } catch {
      return undefined;
    }
  }
}) as any;

// Initialize Analytics (only in browser) - lazy
export const analytics = typeof window !== 'undefined' && firebaseConfig.apiKey
  ? isSupported().then(yes => yes ? getAnalytics(getFirebaseApp()) : null).catch(() => null)
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
