'use client';

import { useState, useEffect } from 'react';
import { useAuth, LoginHistory, SessionInfo } from '@/contexts/auth-context-advanced';

export function SecuritySettings() {
  const {
    user,
    userProfile,
    loginHistory,
    sessions,
    getMFAInfo,
    getLinkedProviders,
    checkAccountSecurity,
    getLoginHistory,
    getAllSessions,
    revokeSession,
    updateUserPassword,
  } = useAuth();

  const [securityScore, setSecurityScore] = useState(0);
  const [recommendations, setRecommendations] = useState<string[]>([]);
  const [activeTab, setActiveTab] = useState<'overview' | 'mfa' | 'sessions' | 'history' | 'password'>(
    'overview'
  );
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  // Password change states
  const [currentPassword, setCurrentPassword] = useState('');
  const [newPassword, setNewPassword] = useState('');
  const [confirmPassword, setConfirmPassword] = useState('');

  useEffect(() => {
    loadSecurityData();
  }, [user]);

  const loadSecurityData = async () => {
    if (!user) return;

    setLoading(true);
    try {
      const { score, recommendations: recs } = await checkAccountSecurity();
      setSecurityScore(score);
      setRecommendations(recs);

      await getLoginHistory(20);
      await getAllSessions();
    } catch (err: any) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  const handlePasswordChange = async (e: React.FormEvent) => {
    e.preventDefault();
    setError(null);

    if (newPassword !== confirmPassword) {
      setError('Passwords do not match');
      return;
    }

    if (newPassword.length < 8) {
      setError('Password must be at least 8 characters');
      return;
    }

    try {
      await updateUserPassword(currentPassword, newPassword);
      setCurrentPassword('');
      setNewPassword('');
      setConfirmPassword('');
      alert('Password updated successfully!');
    } catch (err: any) {
      setError(err.message);
    }
  };

  const handleRevokeSession = async (sessionId: string) => {
    if (!confirm('Are you sure you want to revoke this session?')) return;

    try {
      await revokeSession(sessionId);
      await getAllSessions();
    } catch (err: any) {
      setError(err.message);
    }
  };

  if (!user || !userProfile) {
    return (
      <div className="p-6 bg-yellow-50 border border-yellow-200 rounded-lg">
        <p className="text-yellow-800">Please sign in to manage your security settings.</p>
      </div>
    );
  }

  return (
    <div className="max-w-6xl mx-auto space-y-6">
      {/* Header */}
      <div>
        <h1 className="text-3xl font-bold text-gray-900 dark:text-white">Security Settings</h1>
        <p className="text-gray-600 dark:text-gray-400 mt-1">
          Manage your account security and privacy
        </p>
      </div>

      {error && (
        <div className="p-4 bg-red-50 border border-red-200 text-red-700 rounded-lg">{error}</div>
      )}

      {/* Security Score Card */}
      <div className="bg-white dark:bg-gray-800 rounded-lg shadow p-6">
        <div className="flex items-center justify-between mb-4">
          <h2 className="text-xl font-semibold text-gray-900 dark:text-white">Security Score</h2>
          <div className="text-right">
            <div className="text-3xl font-bold text-blue-600 dark:text-blue-400">{securityScore}%</div>
            <div className="text-sm text-gray-500 dark:text-gray-400">
              {securityScore >= 80 ? 'Excellent' : securityScore >= 60 ? 'Good' : 'Needs Improvement'}
            </div>
          </div>
        </div>

        {/* Score Bar */}
        <div className="w-full bg-gray-200 dark:bg-gray-700 rounded-full h-3 mb-4">
          <div
            className={`h-3 rounded-full transition-all duration-500 ${
              securityScore >= 80
                ? 'bg-green-600'
                : securityScore >= 60
                ? 'bg-yellow-600'
                : 'bg-red-600'
            }`}
            style={{ width: `${securityScore}%` }}
          ></div>
        </div>

        {/* Recommendations */}
        {recommendations.length > 0 && (
          <div className="mt-4">
            <h3 className="text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
              Recommendations to improve your security:
            </h3>
            <ul className="space-y-2">
              {recommendations.map((rec, index) => (
                <li key={index} className="flex items-start text-sm text-gray-600 dark:text-gray-400">
                  <span className="text-yellow-500 mr-2">⚠</span>
                  {rec}
                </li>
              ))}
            </ul>
          </div>
        )}
      </div>

      {/* Tabs */}
      <div className="border-b border-gray-200 dark:border-gray-700">
        <nav className="flex space-x-8">
          {[
            { id: 'overview', label: 'Overview' },
            { id: 'mfa', label: 'Two-Factor Auth' },
            { id: 'sessions', label: 'Active Sessions' },
            { id: 'history', label: 'Login History' },
            { id: 'password', label: 'Password' },
          ].map((tab) => (
            <button
              key={tab.id}
              onClick={() => setActiveTab(tab.id as any)}
              className={`py-2 px-1 border-b-2 font-medium text-sm ${
                activeTab === tab.id
                  ? 'border-blue-500 text-blue-600 dark:text-blue-400'
                  : 'border-transparent text-gray-500 hover:text-gray-700 dark:text-gray-400 dark:hover:text-gray-300'
              }`}
            >
              {tab.label}
            </button>
          ))}
        </nav>
      </div>

      {/* Tab Content */}
      <div className="bg-white dark:bg-gray-800 rounded-lg shadow">
        {activeTab === 'overview' && <OverviewTab userProfile={userProfile} />}
        {activeTab === 'mfa' && <MFATab getMFAInfo={getMFAInfo} />}
        {activeTab === 'sessions' && (
          <SessionsTab sessions={sessions} onRevoke={handleRevokeSession} />
        )}
        {activeTab === 'history' && <LoginHistoryTab history={loginHistory} />}
        {activeTab === 'password' && (
          <PasswordTab
            currentPassword={currentPassword}
            setCurrentPassword={setCurrentPassword}
            newPassword={newPassword}
            setNewPassword={setNewPassword}
            confirmPassword={confirmPassword}
            setConfirmPassword={setConfirmPassword}
            onSubmit={handlePasswordChange}
          />
        )}
      </div>
    </div>
  );
}

function OverviewTab({ userProfile }: { userProfile: any }) {
  return (
    <div className="p-6 space-y-6">
      <div>
        <h3 className="text-lg font-semibold text-gray-900 dark:text-white mb-4">
          Account Security Overview
        </h3>

        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          <SecurityItem
            icon="✓"
            label="Email Verification"
            value={userProfile.emailVerified ? 'Verified' : 'Not Verified'}
            status={userProfile.emailVerified ? 'good' : 'warning'}
          />
          <SecurityItem
            icon="🔒"
            label="Two-Factor Authentication"
            value={userProfile.security?.mfaEnabled ? 'Enabled' : 'Disabled'}
            status={userProfile.security?.mfaEnabled ? 'good' : 'warning'}
          />
          <SecurityItem
            icon="📱"
            label="Phone Number"
            value={userProfile.phoneNumber || 'Not Added'}
            status={userProfile.phoneNumber ? 'good' : 'neutral'}
          />
          <SecurityItem
            icon="🔗"
            label="Linked Accounts"
            value={`${userProfile.security?.linkedProviders?.length || 0} provider(s)`}
            status="neutral"
          />
        </div>
      </div>

      <div>
        <h3 className="text-lg font-semibold text-gray-900 dark:text-white mb-4">
          Recent Activity
        </h3>
        <div className="space-y-3">
          <ActivityItem
            label="Last Login"
            value={new Date(userProfile.lastLoginAt).toLocaleString()}
          />
          <ActivityItem
            label="Last Password Change"
            value={
              userProfile.security?.lastPasswordChange
                ? new Date(userProfile.security.lastPasswordChange).toLocaleString()
                : 'Never'
            }
          />
          <ActivityItem
            label="Account Created"
            value={new Date(userProfile.createdAt).toLocaleString()}
          />
        </div>
      </div>
    </div>
  );
}

function MFATab({ getMFAInfo }: { getMFAInfo: () => any[] }) {
  const enrolledFactors = getMFAInfo();

  return (
    <div className="p-6 space-y-6">
      <div>
        <h3 className="text-lg font-semibold text-gray-900 dark:text-white mb-2">
          Two-Factor Authentication
        </h3>
        <p className="text-gray-600 dark:text-gray-400 mb-4">
          Add an extra layer of security to your account by requiring a second form of verification.
        </p>

        {enrolledFactors.length === 0 ? (
          <div className="bg-yellow-50 dark:bg-yellow-900 border border-yellow-200 dark:border-yellow-700 rounded-lg p-4">
            <p className="text-yellow-800 dark:text-yellow-200">
              Two-factor authentication is not enabled. We strongly recommend enabling it to protect
              your account.
            </p>
            <button className="mt-3 bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded-md text-sm font-medium">
              Enable Two-Factor Authentication
            </button>
          </div>
        ) : (
          <div className="space-y-4">
            {enrolledFactors.map((factor: any) => (
              <div
                key={factor.uid}
                className="border border-gray-200 dark:border-gray-700 rounded-lg p-4"
              >
                <div className="flex items-center justify-between">
                  <div>
                    <p className="font-medium text-gray-900 dark:text-white">
                      {factor.displayName || 'SMS Verification'}
                    </p>
                    <p className="text-sm text-gray-500 dark:text-gray-400">
                      {factor.phoneNumber || 'Phone number'}
                    </p>
                  </div>
                  <button className="text-red-600 hover:text-red-700 text-sm font-medium">
                    Remove
                  </button>
                </div>
              </div>
            ))}
            <button className="bg-gray-100 dark:bg-gray-700 hover:bg-gray-200 dark:hover:bg-gray-600 text-gray-900 dark:text-white px-4 py-2 rounded-md text-sm font-medium">
              + Add Another Method
            </button>
          </div>
        )}
      </div>
    </div>
  );
}

function SessionsTab({
  sessions,
  onRevoke,
}: {
  sessions: SessionInfo[];
  onRevoke: (sessionId: string) => void;
}) {
  return (
    <div className="p-6">
      <h3 className="text-lg font-semibold text-gray-900 dark:text-white mb-4">Active Sessions</h3>
      <p className="text-gray-600 dark:text-gray-400 mb-6">
        Manage where you're signed in. If you see a session you don't recognize, revoke it
        immediately.
      </p>

      {sessions.length === 0 ? (
        <p className="text-gray-500 dark:text-gray-400">No active sessions found.</p>
      ) : (
        <div className="space-y-4">
          {sessions.map((session) => (
            <div
              key={session.sessionId}
              className="border border-gray-200 dark:border-gray-700 rounded-lg p-4"
            >
              <div className="flex items-start justify-between">
                <div className="flex-1">
                  <div className="flex items-center gap-2">
                    <p className="font-medium text-gray-900 dark:text-white">
                      {session.deviceInfo.browser} on {session.deviceInfo.platform}
                    </p>
                    <span className="px-2 py-1 bg-green-100 dark:bg-green-900 text-green-800 dark:text-green-200 text-xs rounded">
                      Current
                    </span>
                  </div>
                  <p className="text-sm text-gray-500 dark:text-gray-400 mt-1">
                    {session.ipAddress || 'Unknown location'}
                  </p>
                  <p className="text-xs text-gray-400 dark:text-gray-500 mt-1">
                    Last active: {new Date(session.lastActivity).toLocaleString()}
                  </p>
                </div>
                <button
                  onClick={() => onRevoke(session.sessionId)}
                  className="text-red-600 hover:text-red-700 text-sm font-medium"
                >
                  Revoke
                </button>
              </div>
            </div>
          ))}
        </div>
      )}

      <button className="mt-6 text-red-600 hover:text-red-700 text-sm font-medium">
        Revoke All Other Sessions
      </button>
    </div>
  );
}

function LoginHistoryTab({ history }: { history: LoginHistory[] }) {
  return (
    <div className="p-6">
      <h3 className="text-lg font-semibold text-gray-900 dark:text-white mb-4">Login History</h3>
      <p className="text-gray-600 dark:text-gray-400 mb-6">
        Review your recent login activity to spot any unauthorized access.
      </p>

      {history.length === 0 ? (
        <p className="text-gray-500 dark:text-gray-400">No login history available.</p>
      ) : (
        <div className="space-y-3">
          {history.map((entry, index) => (
            <div
              key={index}
              className="flex items-center justify-between border-b border-gray-200 dark:border-gray-700 pb-3"
            >
              <div>
                <div className="flex items-center gap-2">
                  <span
                    className={`w-2 h-2 rounded-full ${
                      entry.success ? 'bg-green-500' : 'bg-red-500'
                    }`}
                  ></span>
                  <p className="text-sm font-medium text-gray-900 dark:text-white">
                    {entry.success ? 'Successful login' : 'Failed login attempt'}
                  </p>
                </div>
                <p className="text-xs text-gray-500 dark:text-gray-400 mt-1">
                  {new Date(entry.timestamp).toLocaleString()} • {entry.method} • {entry.userAgent}
                </p>
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}

function PasswordTab({
  currentPassword,
  setCurrentPassword,
  newPassword,
  setNewPassword,
  confirmPassword,
  setConfirmPassword,
  onSubmit,
}: {
  currentPassword: string;
  setCurrentPassword: (v: string) => void;
  newPassword: string;
  setNewPassword: (v: string) => void;
  confirmPassword: string;
  setConfirmPassword: (v: string) => void;
  onSubmit: (e: React.FormEvent) => void;
}) {
  return (
    <div className="p-6">
      <h3 className="text-lg font-semibold text-gray-900 dark:text-white mb-4">Change Password</h3>
      <p className="text-gray-600 dark:text-gray-400 mb-6">
        Choose a strong password that you don't use elsewhere.
      </p>

      <form onSubmit={onSubmit} className="max-w-md space-y-4">
        <div>
          <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
            Current Password
          </label>
          <input
            type="password"
            value={currentPassword}
            onChange={(e) => setCurrentPassword(e.target.value)}
            required
            className="w-full px-4 py-2 border border-gray-300 dark:border-gray-600 rounded-md focus:ring-2 focus:ring-blue-500 focus:border-transparent dark:bg-gray-700 dark:text-white"
          />
        </div>

        <div>
          <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
            New Password
          </label>
          <input
            type="password"
            value={newPassword}
            onChange={(e) => setNewPassword(e.target.value)}
            required
            className="w-full px-4 py-2 border border-gray-300 dark:border-gray-600 rounded-md focus:ring-2 focus:ring-blue-500 focus:border-transparent dark:bg-gray-700 dark:text-white"
          />
          <p className="text-xs text-gray-500 dark:text-gray-400 mt-1">
            Must be at least 8 characters with uppercase, lowercase, and number
          </p>
        </div>

        <div>
          <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
            Confirm New Password
          </label>
          <input
            type="password"
            value={confirmPassword}
            onChange={(e) => setConfirmPassword(e.target.value)}
            required
            className="w-full px-4 py-2 border border-gray-300 dark:border-gray-600 rounded-md focus:ring-2 focus:ring-blue-500 focus:border-transparent dark:bg-gray-700 dark:text-white"
          />
        </div>

        <button
          type="submit"
          className="w-full bg-blue-600 hover:bg-blue-700 text-white font-medium py-2 px-4 rounded-md transition-colors"
        >
          Update Password
        </button>
      </form>
    </div>
  );
}

function SecurityItem({
  icon,
  label,
  value,
  status,
}: {
  icon: string;
  label: string;
  value: string;
  status: 'good' | 'warning' | 'neutral';
}) {
  const statusColors = {
    good: 'bg-green-50 dark:bg-green-900 border-green-200 dark:border-green-700',
    warning: 'bg-yellow-50 dark:bg-yellow-900 border-yellow-200 dark:border-yellow-700',
    neutral: 'bg-gray-50 dark:bg-gray-700 border-gray-200 dark:border-gray-600',
  };

  return (
    <div className={`border rounded-lg p-4 ${statusColors[status]}`}>
      <div className="flex items-center gap-3">
        <span className="text-2xl">{icon}</span>
        <div>
          <p className="text-sm text-gray-600 dark:text-gray-400">{label}</p>
          <p className="font-medium text-gray-900 dark:text-white">{value}</p>
        </div>
      </div>
    </div>
  );
}

function ActivityItem({ label, value }: { label: string; value: string }) {
  return (
    <div className="flex justify-between items-center">
      <span className="text-sm text-gray-600 dark:text-gray-400">{label}</span>
      <span className="text-sm font-medium text-gray-900 dark:text-white">{value}</span>
    </div>
  );
}
