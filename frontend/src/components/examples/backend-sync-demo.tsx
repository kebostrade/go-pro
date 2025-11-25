'use client';

import { useAuth } from '@/contexts/auth-context';
import { api } from '@/lib/api';
import { useEffect, useState } from 'react';

/**
 * Demo component showing backend integration with Firebase Auth
 *
 * This example demonstrates:
 * 1. Firebase user vs Backend user
 * 2. Token-based API calls
 * 3. Sync status and error handling
 * 4. Progress updates to backend
 */
export default function BackendSyncDemo() {
  const { user, userProfile, backendUser, syncWithBackend, loading } = useAuth();
  const [syncing, setSyncing] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [token, setToken] = useState<string | null>(null);

  // Get current ID token
  useEffect(() => {
    if (user) {
      user.getIdToken().then(setToken).catch(console.error);
    }
  }, [user]);

  const handleManualSync = async () => {
    setSyncing(true);
    setError(null);
    try {
      await syncWithBackend();
      console.log('✅ Manual sync successful');
    } catch (err: any) {
      setError(err.message);
      console.error('❌ Manual sync failed:', err);
    } finally {
      setSyncing(false);
    }
  };

  const testAuthenticatedAPI = async () => {
    if (!backendUser) {
      alert('Backend user not loaded. Sign in first.');
      return;
    }

    try {
      // Test authenticated API call
      const progress = await api.getUserProgress(backendUser.id, 1, 5);
      console.log('Progress from backend:', progress);
      alert(`Loaded ${progress.progress.length} progress records`);
    } catch (err: any) {
      console.error('API call failed:', err);
      alert(`API Error: ${err.message}`);
    }
  };

  const testProgressUpdate = async () => {
    if (!backendUser) {
      alert('Backend user not loaded. Sign in first.');
      return;
    }

    try {
      // Test progress update
      const result = await api.updateProgress(backendUser.id, 'lesson-1', {
        completed: true,
        score: 95,
      });
      console.log('Progress update result:', result);
      alert('Progress updated successfully!');
    } catch (err: any) {
      console.error('Progress update failed:', err);
      alert(`Update Error: ${err.message}`);
    }
  };

  if (loading) {
    return <div className="p-4">Loading...</div>;
  }

  if (!user) {
    return (
      <div className="p-4 border rounded-lg">
        <h2 className="text-xl font-bold mb-2">Backend Sync Demo</h2>
        <p className="text-gray-600">Please sign in to test backend integration</p>
      </div>
    );
  }

  const isSynced = backendUser?.firebase_uid === userProfile?.uid;

  return (
    <div className="p-6 border rounded-lg max-w-2xl">
      <h2 className="text-2xl font-bold mb-4">Backend Integration Status</h2>

      {/* Sync Status */}
      <div className="mb-6 p-4 bg-gray-50 rounded">
        <div className="flex items-center justify-between mb-2">
          <span className="font-semibold">Sync Status:</span>
          <span className={`px-2 py-1 rounded text-sm ${
            isSynced ? 'bg-green-100 text-green-800' : 'bg-yellow-100 text-yellow-800'
          }`}>
            {isSynced ? '✅ Synced' : '⚠️ Not Synced'}
          </span>
        </div>
        {error && (
          <div className="mt-2 p-2 bg-red-50 text-red-700 rounded text-sm">
            {error}
          </div>
        )}
      </div>

      {/* Firebase User Info */}
      <div className="mb-4 p-4 border rounded">
        <h3 className="font-semibold mb-2">🔥 Firebase User</h3>
        <div className="text-sm space-y-1">
          <div><strong>UID:</strong> {user.uid}</div>
          <div><strong>Email:</strong> {user.email}</div>
          <div><strong>Display Name:</strong> {user.displayName || 'Not set'}</div>
        </div>
      </div>

      {/* Backend User Info */}
      <div className="mb-4 p-4 border rounded">
        <h3 className="font-semibold mb-2">🖥️ Backend User</h3>
        {backendUser ? (
          <div className="text-sm space-y-1">
            <div><strong>ID:</strong> {backendUser.id}</div>
            <div><strong>Firebase UID:</strong> {backendUser.firebase_uid}</div>
            <div><strong>Email:</strong> {backendUser.email}</div>
            <div><strong>Role:</strong> {backendUser.role}</div>
            <div><strong>Created:</strong> {new Date(backendUser.created_at).toLocaleString()}</div>
          </div>
        ) : (
          <div className="text-sm text-gray-500">Not loaded</div>
        )}
      </div>

      {/* Token Info */}
      <div className="mb-4 p-4 border rounded">
        <h3 className="font-semibold mb-2">🔐 ID Token</h3>
        {token ? (
          <div className="text-xs font-mono bg-gray-100 p-2 rounded overflow-x-auto">
            {token.substring(0, 50)}...
          </div>
        ) : (
          <div className="text-sm text-gray-500">No token</div>
        )}
      </div>

      {/* Actions */}
      <div className="space-y-2">
        <button
          onClick={handleManualSync}
          disabled={syncing}
          className="w-full px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700 disabled:bg-gray-400"
        >
          {syncing ? 'Syncing...' : '🔄 Manual Sync with Backend'}
        </button>

        <button
          onClick={testAuthenticatedAPI}
          disabled={!backendUser}
          className="w-full px-4 py-2 bg-green-600 text-white rounded hover:bg-green-700 disabled:bg-gray-400"
        >
          📊 Test Get Progress API
        </button>

        <button
          onClick={testProgressUpdate}
          disabled={!backendUser}
          className="w-full px-4 py-2 bg-purple-600 text-white rounded hover:bg-purple-700 disabled:bg-gray-400"
        >
          ✏️ Test Update Progress API
        </button>
      </div>

      {/* Debug Info */}
      <details className="mt-4 text-xs">
        <summary className="cursor-pointer font-semibold">🐛 Debug Info</summary>
        <pre className="mt-2 p-2 bg-gray-100 rounded overflow-x-auto">
          {JSON.stringify({
            firebaseUser: {
              uid: user.uid,
              email: user.email,
              emailVerified: user.emailVerified,
            },
            backendUser,
            userProfile: userProfile ? {
              uid: userProfile.uid,
              email: userProfile.email,
              role: userProfile.role,
              progress: userProfile.progress,
            } : null,
          }, null, 2)}
        </pre>
      </details>
    </div>
  );
}
