'use client';

import { useEffect, useState } from 'react';
import {
  LineChart,
  Line,
  BarChart,
  Bar,
  PieChart,
  Pie,
  Cell,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  Legend,
  ResponsiveContainer,
} from 'recharts';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Download, TrendingUp, AlertTriangle, BookOpen } from 'lucide-react';

// Types for analytics data
interface StudentAnalytics {
  learningVelocity: Array<{ week: string; lessonsCompleted: number }>;
  errorPatterns: Array<{ error: string; count: number }>;
  timeDistribution: Array<{ topic: string; hours: number }>;
  peerComparison: {
    userPercentile: number;
    classAverage: number;
    topPercentile: number;
  };
  weakAreas: Array<{
    topic: string;
    score: number;
    status: 'weak' | 'needs-review' | 'good';
    recommendations: string[];
  }>;
}

const COLORS = ['#ef4444', '#f59e0b', '#10b981', '#3b82f6', '#8b5cf6', '#ec4899'];

export default function StudentAnalyticsPage() {
  const [analytics, setAnalytics] = useState<StudentAnalytics | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    fetchAnalytics();
  }, []);

  const fetchAnalytics = async () => {
    try {
      setLoading(true);
      // TODO: Replace with actual API call
      // const response = await api.getStudentAnalytics();
      // setAnalytics(response.data);

      // Mock data for development
      setAnalytics({
        learningVelocity: [
          { week: 'Week 1', lessonsCompleted: 2 },
          { week: 'Week 2', lessonsCompleted: 3 },
          { week: 'Week 3', lessonsCompleted: 4 },
          { week: 'Week 4', lessonsCompleted: 3 },
          { week: 'Week 5', lessonsCompleted: 5 },
          { week: 'Week 6', lessonsCompleted: 4 },
          { week: 'Week 7', lessonsCompleted: 6 },
          { week: 'Week 8', lessonsCompleted: 5 },
          { week: 'Week 9', lessonsCompleted: 4 },
          { week: 'Week 10', lessonsCompleted: 5 },
          { week: 'Week 11', lessonsCompleted: 6 },
          { week: 'Week 12', lessonsCompleted: 7 },
        ],
        errorPatterns: [
          { error: 'nil pointer dereference', count: 15 },
          { error: 'race condition', count: 12 },
          { error: 'type assertion error', count: 10 },
          { error: 'goroutine leak', count: 8 },
          { error: 'channel deadlock', count: 7 },
          { error: 'slice bounds error', count: 6 },
          { error: 'interface nil', count: 5 },
          { error: 'defer error', count: 4 },
          { error: 'context timeout', count: 3 },
          { error: 'panic in goroutine', count: 2 },
        ],
        timeDistribution: [
          { topic: 'Concurrency', hours: 12 },
          { topic: 'HTTP Servers', hours: 8 },
          { topic: 'Databases', hours: 10 },
          { topic: 'Testing', hours: 6 },
          { topic: 'Error Handling', hours: 7 },
        ],
        peerComparison: {
          userPercentile: 72,
          classAverage: 65,
          topPercentile: 95,
        },
        weakAreas: [
          {
            topic: 'Goroutine Management',
            score: 58,
            status: 'weak',
            recommendations: [
              'Review Concurrency in Go lesson',
              'Complete goroutine pooling exercise',
              'Practice with worker patterns',
            ],
          },
          {
            topic: 'Error Handling',
            score: 68,
            status: 'needs-review',
            recommendations: [
              'Review error wrapping patterns',
              'Complete custom errors exercise',
            ],
          },
          {
            topic: 'Channel Patterns',
            score: 75,
            status: 'needs-review',
            recommendations: [
              'Review buffered vs unbuffered channels',
              'Complete select statement exercise',
            ],
          },
          {
            topic: 'HTTP Handlers',
            score: 85,
            status: 'good',
            recommendations: [],
          },
          {
            topic: 'Database Design',
            score: 82,
            status: 'good',
            recommendations: [],
          },
        ],
      });
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
              <AlertTriangle className="h-5 w-5 text-red-500" />
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
          <h1 className="text-3xl font-bold tracking-tight">Your Analytics</h1>
          <p className="text-gray-600 dark:text-gray-400 mt-1">
            Track your learning progress and identify areas for improvement
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

      {/* Learning Velocity Chart */}
      <Card className="mb-8">
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <TrendingUp className="h-5 w-5" />
            Learning Velocity
          </CardTitle>
          <CardDescription>Lessons completed per week over the last 12 weeks</CardDescription>
        </CardHeader>
        <CardContent>
          <ResponsiveContainer width="100%" height={300}>
            <LineChart data={analytics.learningVelocity}>
              <CartesianGrid strokeDasharray="3 3" />
              <XAxis dataKey="week" />
              <YAxis />
              <Tooltip />
              <Legend />
              <Line
                type="monotone"
                dataKey="lessonsCompleted"
                stroke="#3b82f6"
                strokeWidth={2}
                name="Lessons Completed"
              />
            </LineChart>
          </ResponsiveContainer>
        </CardContent>
      </Card>

      {/* Error Patterns and Time Distribution */}
      <div className="grid grid-cols-1 md:grid-cols-2 gap-8 mb-8">
        <Card>
          <CardHeader>
            <CardTitle>Common Error Patterns</CardTitle>
            <CardDescription>Top 10 mistakes to review</CardDescription>
          </CardHeader>
          <CardContent>
            <ResponsiveContainer width="100%" height={300}>
              <BarChart data={analytics.errorPatterns} layout="vertical">
                <CartesianGrid strokeDasharray="3 3" />
                <XAxis type="number" />
                <YAxis dataKey="error" type="category" width={150} />
                <Tooltip />
                <Bar dataKey="count" fill="#ef4444" />
              </BarChart>
            </ResponsiveContainer>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>Time Distribution</CardTitle>
            <CardDescription>Hours spent by topic</CardDescription>
          </CardHeader>
          <CardContent>
            <ResponsiveContainer width="100%" height={300}>
              <PieChart>
                <Pie
                  data={analytics.timeDistribution}
                  cx="50%"
                  cy="50%"
                  labelLine={false}
                  label={({ topic, hours }) => `${topic}: ${hours}h`}
                  outerRadius={80}
                  fill="#8884d8"
                  dataKey="hours"
                >
                  {analytics.timeDistribution.map((entry, index) => (
                    <Cell key={`cell-${index}`} fill={COLORS[index % COLORS.length]} />
                  ))}
                </Pie>
                <Tooltip />
              </PieChart>
            </ResponsiveContainer>
          </CardContent>
        </Card>
      </div>

      {/* Peer Comparison */}
      <Card className="mb-8">
        <CardHeader>
          <CardTitle>Peer Comparison</CardTitle>
          <CardDescription>Your performance relative to the class</CardDescription>
        </CardHeader>
        <CardContent>
          <div className="space-y-4">
            <div>
              <div className="flex justify-between mb-2">
                <span className="text-sm font-medium">Your Percentile</span>
                <span className="text-sm font-bold text-blue-600">{analytics.peerComparison.userPercentile}%</span>
              </div>
              <div className="w-full bg-gray-200 dark:bg-gray-700 rounded-full h-2">
                <div
                  className="bg-blue-600 h-2 rounded-full"
                  style={{ width: `${analytics.peerComparison.userPercentile}%` }}
                ></div>
              </div>
            </div>
            <div>
              <div className="flex justify-between mb-2">
                <span className="text-sm font-medium">Class Average</span>
                <span className="text-sm font-bold text-gray-600">{analytics.peerComparison.classAverage}%</span>
              </div>
              <div className="w-full bg-gray-200 dark:bg-gray-700 rounded-full h-2">
                <div
                  className="bg-gray-500 h-2 rounded-full"
                  style={{ width: `${analytics.peerComparison.classAverage}%` }}
                ></div>
              </div>
            </div>
            <div>
              <div className="flex justify-between mb-2">
                <span className="text-sm font-medium">Top 10% Average</span>
                <span className="text-sm font-bold text-green-600">{analytics.peerComparison.topPercentile}%</span>
              </div>
              <div className="w-full bg-gray-200 dark:bg-gray-700 rounded-full h-2">
                <div
                  className="bg-green-600 h-2 rounded-full"
                  style={{ width: `${analytics.peerComparison.topPercentile}%` }}
                ></div>
              </div>
            </div>
          </div>
        </CardContent>
      </Card>

      {/* Weak Areas & Recommendations */}
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <BookOpen className="h-5 w-5" />
            Areas for Improvement
          </CardTitle>
          <CardDescription>Topics to review and study recommendations</CardDescription>
        </CardHeader>
        <CardContent>
          <div className="space-y-4">
            {analytics.weakAreas.map((area) => (
              <div key={area.topic} className="border rounded-lg p-4">
                <div className="flex items-center justify-between mb-2">
                  <h3 className="font-semibold">{area.topic}</h3>
                  <span
                    className={`px-3 py-1 rounded-full text-sm font-medium ${
                      area.status === 'weak'
                        ? 'bg-red-100 text-red-800 dark:bg-red-900 dark:text-red-200'
                        : area.status === 'needs-review'
                        ? 'bg-yellow-100 text-yellow-800 dark:bg-yellow-900 dark:text-yellow-200'
                        : 'bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-200'
                    }`}
                  >
                    {area.score}%
                  </span>
                </div>
                {area.recommendations.length > 0 && (
                  <div className="mt-3">
                    <p className="text-sm font-medium mb-2">Recommended actions:</p>
                    <ul className="list-disc list-inside text-sm text-gray-600 dark:text-gray-400 space-y-1">
                      {area.recommendations.map((rec, idx) => (
                        <li key={idx}>{rec}</li>
                      ))}
                    </ul>
                  </div>
                )}
              </div>
            ))}
          </div>
        </CardContent>
      </Card>
    </div>
  );
}
