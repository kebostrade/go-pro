'use client';

import { useEffect, useState } from 'react';
import {
  LineChart,
  Line,
  BarChart,
  Bar,
  AreaChart,
  Area,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  Legend,
  ResponsiveContainer,
} from 'recharts';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import {
  Download,
  Users,
  Activity,
  HardDrive,
  Server,
  TrendingUp,
  Zap,
} from 'lucide-react';

// Types for admin analytics
interface AdminAnalytics {
  platformMetrics: {
    totalUsers: number;
    dailyActiveUsers: number;
    weeklyActiveUsers: number;
    monthlyActiveUsers: number;
    contentGrowth: number;
  };
  userGrowth: Array<{ month: string; signups: number; churn: number; retention: number }>;
  contentPopularity: Array<{
    lessonId: string;
    title: string;
    views: number;
    avgTimeSpent: number;
  }>;
  systemHealth: {
    apiResponseTimes: {
      p50: number;
      p95: number;
      p99: number;
    };
    errorRate: number;
    databaseQueryTime: number;
  };
  storageUsage: {
    videos: number;
    downloads: number;
    submissions: number;
    total: number;
    limit: number;
  };
  sandboxUsage: {
    successRate: number;
    avgExecutionTime: number;
    resourceUsage: number;
    totalExecutions: number;
  };
}

export default function AdminAnalyticsPage() {
  const [analytics, setAnalytics] = useState<AdminAnalytics | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [lastUpdated, setLastUpdated] = useState<Date>(new Date());

  useEffect(() => {
    fetchAnalytics();
    const interval = setInterval(fetchAnalytics, 60000); // Refresh every 60 seconds
    return () => clearInterval(interval);
  }, []);

  const fetchAnalytics = async () => {
    try {
      setLoading(true);
      // TODO: Replace with actual API call
      // const response = await api.getPlatformAnalytics();
      // setAnalytics(response.data);

      // Mock data for development
      setAnalytics({
        platformMetrics: {
          totalUsers: 5420,
          dailyActiveUsers: 1250,
          weeklyActiveUsers: 3200,
          monthlyActiveUsers: 4100,
          contentGrowth: 145,
        },
        userGrowth: [
          { month: 'Jan', signups: 320, churn: 15, retention: 85 },
          { month: 'Feb', signups: 450, churn: 22, retention: 87 },
          { month: 'Mar', signups: 520, churn: 18, retention: 86 },
          { month: 'Apr', signups: 610, churn: 25, retention: 88 },
          { month: 'May', signups: 580, churn: 20, retention: 89 },
          { month: 'Jun', signups: 720, churn: 28, retention: 90 },
          { month: 'Jul', signups: 850, churn: 30, retention: 91 },
          { month: 'Aug', signups: 780, churn: 25, retention: 90 },
          { month: 'Sep', signups: 920, churn: 32, retention: 92 },
          { month: 'Oct', signups: 1050, churn: 35, retention: 93 },
          { month: 'Nov', signups: 1120, churn: 38, retention: 94 },
          { month: 'Dec', signups: 1250, churn: 40, retention: 95 },
        ],
        contentPopularity: [
          { lessonId: '1', title: 'Introduction to Go', views: 3420, avgTimeSpent: 45 },
          { lessonId: '2', title: 'Variables and Types', views: 2890, avgTimeSpent: 60 },
          { lessonId: '3', title: 'Control Structures', views: 2560, avgTimeSpent: 75 },
          { lessonId: '4', title: 'Functions', views: 2240, avgTimeSpent: 90 },
          { lessonId: '5', title: 'Concurrency Basics', views: 1980, avgTimeSpent: 120 },
          { lessonId: '6', title: 'Error Handling', views: 890, avgTimeSpent: 55 },
          { lessonId: '7', title: 'Interfaces', views: 650, avgTimeSpent: 85 },
          { lessonId: '8', title: 'Goroutines', views: 420, avgTimeSpent: 110 },
        ],
        systemHealth: {
          apiResponseTimes: {
            p50: 120,
            p95: 380,
            p99: 650,
          },
          errorRate: 0.12,
          databaseQueryTime: 45,
        },
        storageUsage: {
          videos: 45.2,
          downloads: 12.8,
          submissions: 8.5,
          total: 66.5,
          limit: 1000,
        },
        sandboxUsage: {
          successRate: 94.8,
          avgExecutionTime: 3200,
          resourceUsage: 72,
          totalExecutions: 15420,
        },
      });
      setLastUpdated(new Date());
    } catch (err) {
      setError('Failed to load analytics');
      console.error(err);
    } finally {
      setLoading(false);
    }
  };

  const handleExport = async (format: 'json' | 'csv') => {
    try {
      // TODO: Implement export functionality
      console.log(`Exporting analytics as ${format}`);
    } catch (err) {
      console.error('Export failed:', err);
    }
  };

  if (loading) {
    return (
      <div className="flex items-center justify-center min-h-screen">
        <div className="text-center">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-gray-900 dark:border-gray-100 mx-auto mb-4"></div>
          <p className="text-gray-600 dark:text-gray-400">Loading analytics...</p>
        </div>
      </div>
    );
  }

  if (error || !analytics) {
    return (
      <div className="flex items-center justify-center min-h-screen">
        <Card className="max-w-md">
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              Error Loading Analytics
            </CardTitle>
          </CardHeader>
          <CardContent>
            <p className="text-gray-600 dark:text-gray-400 mb-4">{error || 'Unknown error'}</p>
            <Button onClick={fetchAnalytics}>Retry</Button>
          </CardContent>
        </Card>
      </div>
    );
  }

  return (
    <div className="container mx-auto px-4 py-8 max-w-7xl">
      <div className="flex items-center justify-between mb-8">
        <div>
          <h1 className="text-3xl font-bold tracking-tight">Platform Analytics</h1>
          <p className="text-gray-600 dark:text-gray-400 mt-1">
            Platform-wide metrics and system health
          </p>
          <p className="text-xs text-gray-500 mt-1">
            Last updated: {lastUpdated.toLocaleTimeString()}
          </p>
        </div>
        <div className="flex gap-2">
          <Button variant="outline" onClick={() => handleExport('json')}>
            <Download className="h-4 w-4 mr-2" />
            Export JSON
          </Button>
          <Button variant="outline" onClick={() => handleExport('csv')}>
            <Download className="h-4 w-4 mr-2" />
            Export CSV
          </Button>
        </div>
      </div>

      {/* Platform Metrics */}
      <div className="grid grid-cols-1 md:grid-cols-5 gap-6 mb-8">
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Total Users</CardTitle>
            <Users className="h-4 w-4 text-gray-600" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{analytics.platformMetrics.totalUsers.toLocaleString()}</div>
            <p className="text-xs text-gray-600 dark:text-gray-400 mt-1">
              +{analytics.platformMetrics.contentGrowth} this month
            </p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Daily Active</CardTitle>
            <Activity className="h-4 w-4 text-blue-600" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{analytics.platformMetrics.dailyActiveUsers.toLocaleString()}</div>
            <p className="text-xs text-gray-600 dark:text-gray-400 mt-1">
              {((analytics.platformMetrics.dailyActiveUsers / analytics.platformMetrics.totalUsers) * 100).toFixed(1)}% of total
            </p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Weekly Active</CardTitle>
            <Activity className="h-4 w-4 text-green-600" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{analytics.platformMetrics.weeklyActiveUsers.toLocaleString()}</div>
            <p className="text-xs text-gray-600 dark:text-gray-400 mt-1">
              {((analytics.platformMetrics.weeklyActiveUsers / analytics.platformMetrics.totalUsers) * 100).toFixed(1)}% of total
            </p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Monthly Active</CardTitle>
            <Activity className="h-4 w-4 text-purple-600" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{analytics.platformMetrics.monthlyActiveUsers.toLocaleString()}</div>
            <p className="text-xs text-gray-600 dark:text-gray-400 mt-1">
              {((analytics.platformMetrics.monthlyActiveUsers / analytics.platformMetrics.totalUsers) * 100).toFixed(1)}% of total
            </p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Content Items</CardTitle>
            <TrendingUp className="h-4 w-4 text-gray-600" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{analytics.platformMetrics.contentGrowth}</div>
            <p className="text-xs text-gray-600 dark:text-gray-400 mt-1">
              Lessons created
            </p>
          </CardContent>
        </Card>
      </div>

      {/* User Growth Chart */}
      <Card className="mb-8">
        <CardHeader>
          <CardTitle>User Growth Trends</CardTitle>
          <CardDescription>12-month signup, churn, and retention analysis</CardDescription>
        </CardHeader>
        <CardContent>
          <ResponsiveContainer width="100%" height={400}>
            <AreaChart data={analytics.userGrowth}>
              <CartesianGrid strokeDasharray="3 3" />
              <XAxis dataKey="month" />
              <YAxis />
              <Tooltip />
              <Legend />
              <Area type="monotone" dataKey="signups" stackId="1" stroke="#3b82f6" fill="#3b82f6" name="Signups" />
              <Area type="monotone" dataKey="churn" stackId="2" stroke="#ef4444" fill="#ef4444" name="Churn" />
              <Line type="monotone" dataKey="retention" stroke="#10b981" strokeWidth={3} name="Retention Rate %" />
            </AreaChart>
          </ResponsiveContainer>
        </CardContent>
      </Card>

      {/* Content Popularity */}
      <Card className="mb-8">
        <CardHeader>
          <CardTitle>Content Performance</CardTitle>
          <CardDescription>Most and least viewed lessons</CardDescription>
        </CardHeader>
        <CardContent>
          <div className="space-y-4">
            {analytics.contentPopularity.map((lesson) => (
              <div key={lesson.lessonId} className="flex items-center justify-between p-4 border rounded-lg">
                <div className="flex-1">
                  <h3 className="font-semibold">{lesson.title}</h3>
                  <p className="text-sm text-gray-600 dark:text-gray-400">
                    Avg time spent: {lesson.avgTimeSpent} min
                  </p>
                </div>
                <div className="text-right">
                  <div className="text-2xl font-bold">{lesson.views.toLocaleString()}</div>
                  <p className="text-xs text-gray-600 dark:text-gray-400">views</p>
                </div>
              </div>
            ))}
          </div>
        </CardContent>
      </Card>

      {/* System Health and Storage */}
      <div className="grid grid-cols-1 md:grid-cols-2 gap-8 mb-8">
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <Server className="h-5 w-5" />
              System Health
            </CardTitle>
            <CardDescription>API response times and error rates</CardDescription>
          </CardHeader>
          <CardContent>
            <div className="space-y-6">
              <div>
                <h4 className="font-semibold mb-3">API Response Times (ms)</h4>
                <div className="space-y-2">
                  <div className="flex justify-between items-center">
                    <span className="text-sm">p50</span>
                    <span className="font-mono font-bold text-green-600">{analytics.systemHealth.apiResponseTimes.p50}ms</span>
                  </div>
                  <div className="flex justify-between items-center">
                    <span className="text-sm">p95</span>
                    <span className="font-mono font-bold text-yellow-600">{analytics.systemHealth.apiResponseTimes.p95}ms</span>
                  </div>
                  <div className="flex justify-between items-center">
                    <span className="text-sm">p99</span>
                    <span className="font-mono font-bold text-orange-600">{analytics.systemHealth.apiResponseTimes.p99}ms</span>
                  </div>
                </div>
              </div>

              <div>
                <h4 className="font-semibold mb-3">Error Rate</h4>
                <div className="text-2xl font-bold text-red-600">{analytics.systemHealth.errorRate}%</div>
              </div>

              <div>
                <h4 className="font-semibold mb-3">Database Query Time</h4>
                <div className="text-2xl font-bold">{analytics.systemHealth.databaseQueryTime}ms</div>
                <p className="text-xs text-gray-600 dark:text-gray-400 mt-1">Average query latency</p>
              </div>
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <HardDrive className="h-5 w-5" />
              Storage Usage
            </CardTitle>
            <CardDescription>S3 storage by type</CardDescription>
          </CardHeader>
          <CardContent>
            <div className="space-y-4">
              <div>
                <div className="flex justify-between mb-2">
                  <span className="text-sm font-medium">Videos</span>
                  <span className="text-sm font-bold">{analytics.storageUsage.videos} GB</span>
                </div>
                <div className="w-full bg-gray-200 dark:bg-gray-700 rounded-full h-2">
                  <div
                    className="bg-blue-600 h-2 rounded-full"
                    style={{ width: `${(analytics.storageUsage.videos / analytics.storageUsage.limit) * 100}%` }}
                  ></div>
                </div>
              </div>

              <div>
                <div className="flex justify-between mb-2">
                  <span className="text-sm font-medium">Downloads</span>
                  <span className="text-sm font-bold">{analytics.storageUsage.downloads} GB</span>
                </div>
                <div className="w-full bg-gray-200 dark:bg-gray-700 rounded-full h-2">
                  <div
                    className="bg-green-600 h-2 rounded-full"
                    style={{ width: `${(analytics.storageUsage.downloads / analytics.storageUsage.limit) * 100}%` }}
                  ></div>
                </div>
              </div>

              <div>
                <div className="flex justify-between mb-2">
                  <span className="text-sm font-medium">Submissions</span>
                  <span className="text-sm font-bold">{analytics.storageUsage.submissions} GB</span>
                </div>
                <div className="w-full bg-gray-200 dark:bg-gray-700 rounded-full h-2">
                  <div
                    className="bg-purple-600 h-2 rounded-full"
                    style={{ width: `${(analytics.storageUsage.submissions / analytics.storageUsage.limit) * 100}%` }}
                  ></div>
                </div>
              </div>

              <div className="pt-4 border-t">
                <div className="flex justify-between mb-2">
                  <span className="text-sm font-semibold">Total Used</span>
                  <span className="text-sm font-bold">{analytics.storageUsage.total} GB</span>
                </div>
                <div className="w-full bg-gray-200 dark:bg-gray-700 rounded-full h-2">
                  <div
                    className="bg-orange-600 h-2 rounded-full"
                    style={{ width: `${(analytics.storageUsage.total / analytics.storageUsage.limit) * 100}%` }}
                  ></div>
                </div>
                <p className="text-xs text-gray-600 dark:text-gray-400 mt-2">
                  of {analytics.storageUsage.limit} GB limit
                </p>
              </div>
            </div>
          </CardContent>
        </Card>
      </div>

      {/* Sandbox Usage */}
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <Zap className="h-5 w-5" />
            Code Execution Sandbox
          </CardTitle>
          <CardDescription>Docker sandbox performance and resource usage</CardDescription>
        </CardHeader>
        <CardContent>
          <div className="grid grid-cols-1 md:grid-cols-4 gap-6">
            <div>
              <h4 className="text-sm font-medium mb-2">Success Rate</h4>
              <div className="text-2xl font-bold text-green-600">{analytics.sandboxUsage.successRate}%</div>
              <p className="text-xs text-gray-600 dark:text-gray-400 mt-1">
                {analytics.sandboxUsage.totalExecutions.toLocaleString()} total executions
              </p>
            </div>

            <div>
              <h4 className="text-sm font-medium mb-2">Avg Execution Time</h4>
              <div className="text-2xl font-bold">{(analytics.sandboxUsage.avgExecutionTime / 1000).toFixed(2)}s</div>
              <p className="text-xs text-gray-600 dark:text-gray-400 mt-1">
                95th percentile
              </p>
            </div>

            <div>
              <h4 className="text-sm font-medium mb-2">Resource Usage</h4>
              <div className="text-2xl font-bold">{analytics.sandboxUsage.resourceUsage}%</div>
              <p className="text-xs text-gray-600 dark:text-gray-400 mt-1">
                Avg memory/CPU utilization
              </p>
            </div>

            <div>
              <h4 className="text-sm font-medium mb-2">Total Executions</h4>
              <div className="text-2xl font-bold">{analytics.sandboxUsage.totalExecutions.toLocaleString()}</div>
              <p className="text-xs text-gray-600 dark:text-gray-400 mt-1">
                All time
              </p>
            </div>
          </div>
        </CardContent>
      </Card>
    </div>
  );
}
