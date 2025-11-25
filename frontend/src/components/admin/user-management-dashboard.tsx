'use client';

import { useState, useEffect } from 'react';
import { useAuth, UserProfile } from '@/contexts/auth-context-advanced';

export function UserManagementDashboard() {
  const { userProfile, adminGetAllUsers, adminUpdateUserRole, adminDisableUser } = useAuth();
  const [users, setUsers] = useState<UserProfile[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [searchTerm, setSearchTerm] = useState('');
  const [filterRole, setFilterRole] = useState<'all' | 'student' | 'instructor' | 'admin'>('all');
  const [selectedUser, setSelectedUser] = useState<UserProfile | null>(null);

  useEffect(() => {
    if (userProfile?.role === 'admin') {
      loadUsers();
    }
  }, [userProfile]);

  const loadUsers = async () => {
    setLoading(true);
    setError(null);
    try {
      const allUsers = await adminGetAllUsers();
      setUsers(allUsers);
    } catch (err: any) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  const handleRoleChange = async (userId: string, newRole: 'student' | 'instructor' | 'admin') => {
    try {
      await adminUpdateUserRole(userId, newRole);
      await loadUsers();
    } catch (err: any) {
      setError(err.message);
    }
  };

  const handleDisableUser = async (userId: string) => {
    if (!confirm('Are you sure you want to disable this user account?')) return;

    try {
      await adminDisableUser(userId);
      await loadUsers();
    } catch (err: any) {
      setError(err.message);
    }
  };

  const filteredUsers = users.filter((user) => {
    const matchesSearch =
      user.displayName?.toLowerCase().includes(searchTerm.toLowerCase()) ||
      user.email?.toLowerCase().includes(searchTerm.toLowerCase());

    const matchesRole = filterRole === 'all' || user.role === filterRole;

    return matchesSearch && matchesRole;
  });

  if (userProfile?.role !== 'admin') {
    return (
      <div className="p-6 bg-red-50 border border-red-200 rounded-lg">
        <h2 className="text-xl font-bold text-red-800">Access Denied</h2>
        <p className="text-red-600 mt-2">You do not have admin access to view this page.</p>
      </div>
    );
  }

  if (loading) {
    return (
      <div className="flex items-center justify-center p-12">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold text-gray-900 dark:text-white">User Management</h1>
          <p className="text-gray-600 dark:text-gray-400 mt-1">
            Manage user accounts, roles, and permissions
          </p>
        </div>
        <div className="bg-blue-50 dark:bg-blue-900 px-4 py-2 rounded-lg">
          <p className="text-sm text-blue-600 dark:text-blue-300">Total Users</p>
          <p className="text-2xl font-bold text-blue-900 dark:text-blue-100">{users.length}</p>
        </div>
      </div>

      {error && (
        <div className="p-4 bg-red-50 border border-red-200 text-red-700 rounded-lg">{error}</div>
      )}

      {/* Stats Cards */}
      <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
        <StatCard title="Students" count={users.filter((u) => u.role === 'student').length} color="blue" />
        <StatCard
          title="Instructors"
          count={users.filter((u) => u.role === 'instructor').length}
          color="green"
        />
        <StatCard title="Admins" count={users.filter((u) => u.role === 'admin').length} color="purple" />
        <StatCard
          title="Verified"
          count={users.filter((u) => u.emailVerified).length}
          color="yellow"
        />
      </div>

      {/* Filters */}
      <div className="bg-white dark:bg-gray-800 p-4 rounded-lg shadow">
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div>
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
              Search Users
            </label>
            <input
              type="text"
              value={searchTerm}
              onChange={(e) => setSearchTerm(e.target.value)}
              placeholder="Search by name or email..."
              className="w-full px-4 py-2 border border-gray-300 dark:border-gray-600 rounded-md focus:ring-2 focus:ring-blue-500 focus:border-transparent dark:bg-gray-700 dark:text-white"
            />
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
              Filter by Role
            </label>
            <select
              value={filterRole}
              onChange={(e) => setFilterRole(e.target.value as any)}
              className="w-full px-4 py-2 border border-gray-300 dark:border-gray-600 rounded-md focus:ring-2 focus:ring-blue-500 focus:border-transparent dark:bg-gray-700 dark:text-white"
            >
              <option value="all">All Roles</option>
              <option value="student">Students</option>
              <option value="instructor">Instructors</option>
              <option value="admin">Admins</option>
            </select>
          </div>
        </div>
      </div>

      {/* Users Table */}
      <div className="bg-white dark:bg-gray-800 rounded-lg shadow overflow-hidden">
        <div className="overflow-x-auto">
          <table className="w-full">
            <thead className="bg-gray-50 dark:bg-gray-700">
              <tr>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">
                  User
                </th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">
                  Role
                </th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">
                  Status
                </th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">
                  Progress
                </th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">
                  Joined
                </th>
                <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">
                  Actions
                </th>
              </tr>
            </thead>
            <tbody className="divide-y divide-gray-200 dark:divide-gray-700">
              {filteredUsers.map((user) => (
                <tr key={user.uid} className="hover:bg-gray-50 dark:hover:bg-gray-700">
                  <td className="px-6 py-4 whitespace-nowrap">
                    <div className="flex items-center">
                      {user.photoURL ? (
                        <img
                          src={user.photoURL}
                          alt={user.displayName || 'User'}
                          className="w-10 h-10 rounded-full"
                        />
                      ) : (
                        <div className="w-10 h-10 rounded-full bg-blue-100 dark:bg-blue-900 flex items-center justify-center">
                          <span className="text-blue-600 dark:text-blue-300 font-medium">
                            {user.displayName?.[0] || user.email?.[0] || '?'}
                          </span>
                        </div>
                      )}
                      <div className="ml-4">
                        <div className="text-sm font-medium text-gray-900 dark:text-white">
                          {user.displayName || 'Unnamed User'}
                        </div>
                        <div className="text-sm text-gray-500 dark:text-gray-400">{user.email}</div>
                      </div>
                    </div>
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap">
                    <select
                      value={user.role}
                      onChange={(e) => handleRoleChange(user.uid, e.target.value as any)}
                      className="text-sm border border-gray-300 dark:border-gray-600 rounded px-2 py-1 dark:bg-gray-700 dark:text-white"
                      disabled={user.uid === userProfile?.uid}
                    >
                      <option value="student">Student</option>
                      <option value="instructor">Instructor</option>
                      <option value="admin">Admin</option>
                    </select>
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap">
                    <div className="flex flex-col gap-1">
                      {user.emailVerified ? (
                        <span className="inline-flex items-center px-2 py-1 rounded-full text-xs font-medium bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-200">
                          ✓ Verified
                        </span>
                      ) : (
                        <span className="inline-flex items-center px-2 py-1 rounded-full text-xs font-medium bg-yellow-100 text-yellow-800 dark:bg-yellow-900 dark:text-yellow-200">
                          ⚠ Unverified
                        </span>
                      )}
                      {user.security?.mfaEnabled && (
                        <span className="inline-flex items-center px-2 py-1 rounded-full text-xs font-medium bg-blue-100 text-blue-800 dark:bg-blue-900 dark:text-blue-200">
                          🔒 MFA
                        </span>
                      )}
                    </div>
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900 dark:text-gray-100">
                    <div>Level {user.progress?.level || 1}</div>
                    <div className="text-xs text-gray-500 dark:text-gray-400">
                      {user.progress?.xp || 0} XP
                    </div>
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500 dark:text-gray-400">
                    {new Date(user.createdAt).toLocaleDateString()}
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
                    <button
                      onClick={() => setSelectedUser(user)}
                      className="text-blue-600 hover:text-blue-900 dark:text-blue-400 dark:hover:text-blue-200 mr-3"
                    >
                      View
                    </button>
                    {user.uid !== userProfile?.uid && (
                      <button
                        onClick={() => handleDisableUser(user.uid)}
                        className="text-red-600 hover:text-red-900 dark:text-red-400 dark:hover:text-red-200"
                        disabled={user.security?.accountLockout}
                      >
                        {user.security?.accountLockout ? 'Disabled' : 'Disable'}
                      </button>
                    )}
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>

        {filteredUsers.length === 0 && (
          <div className="text-center py-12 text-gray-500 dark:text-gray-400">
            No users found matching your criteria
          </div>
        )}
      </div>

      {/* User Detail Modal */}
      {selectedUser && (
        <UserDetailModal user={selectedUser} onClose={() => setSelectedUser(null)} />
      )}
    </div>
  );
}

function StatCard({ title, count, color }: { title: string; count: number; color: string }) {
  const colors = {
    blue: 'bg-blue-50 dark:bg-blue-900 text-blue-600 dark:text-blue-300',
    green: 'bg-green-50 dark:bg-green-900 text-green-600 dark:text-green-300',
    purple: 'bg-purple-50 dark:bg-purple-900 text-purple-600 dark:text-purple-300',
    yellow: 'bg-yellow-50 dark:bg-yellow-900 text-yellow-600 dark:text-yellow-300',
  };

  return (
    <div className={`p-4 rounded-lg ${colors[color as keyof typeof colors]}`}>
      <p className="text-sm font-medium">{title}</p>
      <p className="text-3xl font-bold mt-1">{count}</p>
    </div>
  );
}

function UserDetailModal({ user, onClose }: { user: UserProfile; onClose: () => void }) {
  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4 z-50">
      <div className="bg-white dark:bg-gray-800 rounded-lg shadow-xl max-w-2xl w-full max-h-[90vh] overflow-y-auto">
        <div className="p-6">
          <div className="flex items-center justify-between mb-6">
            <h2 className="text-2xl font-bold text-gray-900 dark:text-white">User Details</h2>
            <button
              onClick={onClose}
              className="text-gray-500 hover:text-gray-700 dark:text-gray-400 dark:hover:text-gray-200"
            >
              <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth={2}
                  d="M6 18L18 6M6 6l12 12"
                />
              </svg>
            </button>
          </div>

          <div className="space-y-6">
            {/* Profile Info */}
            <div>
              <h3 className="text-lg font-semibold text-gray-900 dark:text-white mb-3">
                Profile Information
              </h3>
              <div className="grid grid-cols-2 gap-4">
                <InfoField label="Display Name" value={user.displayName || 'N/A'} />
                <InfoField label="Email" value={user.email || 'N/A'} />
                <InfoField label="Phone" value={user.phoneNumber || 'N/A'} />
                <InfoField label="Role" value={user.role} />
                <InfoField
                  label="Email Verified"
                  value={user.emailVerified ? 'Yes' : 'No'}
                />
                <InfoField
                  label="MFA Enabled"
                  value={user.security?.mfaEnabled ? 'Yes' : 'No'}
                />
              </div>
            </div>

            {/* Progress */}
            <div>
              <h3 className="text-lg font-semibold text-gray-900 dark:text-white mb-3">
                Learning Progress
              </h3>
              <div className="grid grid-cols-3 gap-4">
                <InfoField label="Level" value={user.progress?.level || 1} />
                <InfoField label="XP" value={user.progress?.xp || 0} />
                <InfoField
                  label="Completed Lessons"
                  value={user.progress?.completedLessons?.length || 0}
                />
              </div>
            </div>

            {/* Security */}
            <div>
              <h3 className="text-lg font-semibold text-gray-900 dark:text-white mb-3">
                Security Information
              </h3>
              <div className="grid grid-cols-2 gap-4">
                <InfoField
                  label="Linked Providers"
                  value={user.security?.linkedProviders?.join(', ') || 'None'}
                />
                <InfoField
                  label="Failed Login Attempts"
                  value={user.security?.failedLoginAttempts || 0}
                />
                <InfoField
                  label="Account Status"
                  value={user.security?.accountLockout ? 'Locked' : 'Active'}
                />
              </div>
            </div>

            {/* Metadata */}
            <div>
              <h3 className="text-lg font-semibold text-gray-900 dark:text-white mb-3">Metadata</h3>
              <div className="grid grid-cols-2 gap-4">
                <InfoField
                  label="Created At"
                  value={new Date(user.createdAt).toLocaleString()}
                />
                <InfoField
                  label="Last Login"
                  value={new Date(user.lastLoginAt).toLocaleString()}
                />
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}

function InfoField({ label, value }: { label: string; value: string | number }) {
  return (
    <div>
      <p className="text-sm text-gray-500 dark:text-gray-400">{label}</p>
      <p className="text-sm font-medium text-gray-900 dark:text-white">{value}</p>
    </div>
  );
}
