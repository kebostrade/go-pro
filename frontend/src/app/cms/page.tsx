'use client';

import { useEffect, useState } from 'react';
import Link from 'next/link';
import {
  BookOpen,
  Users,
  FileEdit,
  GraduationCap,
  TrendingUp,
  Clock,
  AlertCircle,
  ArrowRight,
} from 'lucide-react';

interface DashboardStats {
  totalLessons: number;
  totalStudents: number;
  pendingSubmissions: number;
  averageScore: number;
}

interface RecentActivity {
  id: string;
  type: 'lesson_published' | 'submission' | 'enrollment';
  message: string;
  timestamp: string;
}

export default function CMSDashboardPage() {
  const [stats, setStats] = useState<DashboardStats>({
    totalLessons: 0,
    totalStudents: 0,
    pendingSubmissions: 0,
    averageScore: 0,
  });
  const [recentActivity, setRecentActivity] = useState<RecentActivity[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    // TODO: Fetch real data from API
    // For now, using mock data
    setTimeout(() => {
      setStats({
        totalLessons: 24,
        totalStudents: 156,
        pendingSubmissions: 12,
        averageScore: 78.5,
      });
      setRecentActivity([
        {
          id: '1',
          type: 'lesson_published',
          message: 'Published "Introduction to Go Concurrency"',
          timestamp: '2 hours ago',
        },
        {
          id: '2',
          type: 'submission',
          message: '5 new project submissions received',
          timestamp: '4 hours ago',
        },
        {
          id: '3',
          type: 'enrollment',
          message: '3 students enrolled in Web Development path',
          timestamp: '6 hours ago',
        },
      ]);
      setLoading(false);
    }, 500);
  }, []);

  if (loading) {
    return (
      <div className="flex items-center justify-center h-64">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      {/* Header */}
      <div>
        <h1 className="text-3xl font-bold text-gray-900">Dashboard</h1>
        <p className="mt-2 text-sm text-gray-600">
          Welcome back! Here's what's happening with your courses.
        </p>
      </div>

      {/* Stats Grid */}
      <div className="grid grid-cols-1 gap-6 sm:grid-cols-2 lg:grid-cols-4">
        <div className="bg-white rounded-lg shadow p-6">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm font-medium text-gray-600">Total Lessons</p>
              <p className="mt-2 text-3xl font-bold text-gray-900">{stats.totalLessons}</p>
            </div>
            <div className="rounded-full bg-blue-100 p-3">
              <BookOpen className="h-6 w-6 text-blue-600" />
            </div>
          </div>
        </div>

        <div className="bg-white rounded-lg shadow p-6">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm font-medium text-gray-600">Total Students</p>
              <p className="mt-2 text-3xl font-bold text-gray-900">{stats.totalStudents}</p>
            </div>
            <div className="rounded-full bg-green-100 p-3">
              <Users className="h-6 w-6 text-green-600" />
            </div>
          </div>
        </div>

        <div className="bg-white rounded-lg shadow p-6">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm font-medium text-gray-600">Pending Reviews</p>
              <p className="mt-2 text-3xl font-bold text-gray-900">{stats.pendingSubmissions}</p>
            </div>
            <div className="rounded-full bg-yellow-100 p-3">
              <FileEdit className="h-6 w-6 text-yellow-600" />
            </div>
          </div>
        </div>

        <div className="bg-white rounded-lg shadow p-6">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm font-medium text-gray-600">Average Score</p>
              <p className="mt-2 text-3xl font-bold text-gray-900">{stats.averageScore}%</p>
            </div>
            <div className="rounded-full bg-purple-100 p-3">
              <GraduationCap className="h-6 w-6 text-purple-600" />
            </div>
          </div>
        </div>
      </div>

      {/* Quick Actions & Recent Activity */}
      <div className="grid grid-cols-1 gap-6 lg:grid-cols-2">
        {/* Quick Actions */}
        <div className="bg-white rounded-lg shadow">
          <div className="p-6 border-b border-gray-200">
            <h2 className="text-lg font-semibold text-gray-900">Quick Actions</h2>
          </div>
          <div className="p-6 space-y-3">
            <Link
              href="/cms/content/lessons/new"
              className="flex items-center justify-between p-4 rounded-lg border border-gray-200 hover:border-blue-500 hover:bg-blue-50 transition-colors group"
            >
              <div className="flex items-center space-x-3">
                <div className="rounded-full bg-blue-100 p-2 group-hover:bg-blue-200 transition-colors">
                  <FileEdit className="h-5 w-5 text-blue-600" />
                </div>
                <div>
                  <p className="font-medium text-gray-900">Create New Lesson</p>
                  <p className="text-sm text-gray-500">Add a new lesson to the curriculum</p>
                </div>
              </div>
              <ArrowRight className="h-5 w-5 text-gray-400 group-hover:text-blue-600 transition-colors" />
            </Link>

            <Link
              href="/cms/grading/submissions"
              className="flex items-center justify-between p-4 rounded-lg border border-gray-200 hover:border-green-500 hover:bg-green-50 transition-colors group"
            >
              <div className="flex items-center space-x-3">
                <div className="rounded-full bg-green-100 p-2 group-hover:bg-green-200 transition-colors">
                  <GraduationCap className="h-5 w-5 text-green-600" />
                </div>
                <div>
                  <p className="font-medium text-gray-900">Review Submissions</p>
                  <p className="text-sm text-gray-500">{stats.pendingSubmissions} pending reviews</p>
                </div>
              </div>
              <ArrowRight className="h-5 w-5 text-gray-400 group-hover:text-green-600 transition-colors" />
            </Link>

            <Link
              href="/cms/analytics"
              className="flex items-center justify-between p-4 rounded-lg border border-gray-200 hover:border-purple-500 hover:bg-purple-50 transition-colors group"
            >
              <div className="flex items-center space-x-3">
                <div className="rounded-full bg-purple-100 p-2 group-hover:bg-purple-200 transition-colors">
                  <TrendingUp className="h-5 w-5 text-purple-600" />
                </div>
                <div>
                  <p className="font-medium text-gray-900">View Analytics</p>
                  <p className="text-sm text-gray-500">Track student performance</p>
                </div>
              </div>
              <ArrowRight className="h-5 w-5 text-gray-400 group-hover:text-purple-600 transition-colors" />
            </Link>
          </div>
        </div>

        {/* Recent Activity */}
        <div className="bg-white rounded-lg shadow">
          <div className="p-6 border-b border-gray-200">
            <h2 className="text-lg font-semibold text-gray-900">Recent Activity</h2>
          </div>
          <div className="p-6">
            {recentActivity.length === 0 ? (
              <div className="text-center py-8">
                <Clock className="h-12 w-12 text-gray-300 mx-auto mb-3" />
                <p className="text-sm text-gray-500">No recent activity</p>
              </div>
            ) : (
              <div className="space-y-4">
                {recentActivity.map((activity) => (
                  <div key={activity.id} className="flex items-start space-x-3">
                    <div className="rounded-full bg-blue-100 p-2 mt-0.5">
                      {activity.type === 'lesson_published' && (
                        <FileEdit className="h-4 w-4 text-blue-600" />
                      )}
                      {activity.type === 'submission' && (
                        <GraduationCap className="h-4 w-4 text-blue-600" />
                      )}
                      {activity.type === 'enrollment' && (
                        <Users className="h-4 w-4 text-blue-600" />
                      )}
                    </div>
                    <div className="flex-1">
                      <p className="text-sm text-gray-900">{activity.message}</p>
                      <p className="text-xs text-gray-500 mt-1">{activity.timestamp}</p>
                    </div>
                  </div>
                ))}
              </div>
            )}
          </div>
        </div>
      </div>

      {/* Alerts */}
      {stats.pendingSubmissions > 10 && (
        <div className="bg-yellow-50 border border-yellow-200 rounded-lg p-4 flex items-start space-x-3">
          <AlertCircle className="h-6 w-6 text-yellow-600 flex-shrink-0 mt-0.5" />
          <div className="flex-1">
            <h3 className="text-sm font-medium text-yellow-800">
              High submission volume
            </h3>
            <p className="text-sm text-yellow-700 mt-1">
              You have {stats.pendingSubmissions} pending reviews. Consider setting aside time
              for grading or invite peer reviewers to help with the workload.
            </p>
          </div>
        </div>
      )}
    </div>
  );
}
