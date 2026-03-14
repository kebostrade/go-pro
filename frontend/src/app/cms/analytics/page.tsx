'use client';

import { useEffect, useState } from 'react';
import { useRouter } from 'next/navigation';
import {
  BarChart,
  Bar,
  LineChart,
  Line,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  Legend,
  ResponsiveContainer,
  Cell,
} from 'recharts';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Badge } from '@/components/ui/badge';
import {
  Download,
  Users,
  TrendingUp,
  AlertCircle,
  Clock,
  Target,
} from 'lucide-react';

// Types for instructor analytics
interface InstructorAnalytics {
  classOverview: {
    totalStudents: number;
    activeThisWeek: number;
    averageCompletionRate: number;
  };
  engagementHeatmap: Array<{ day: string; hour: number; activity: number }>;
  completionFunnel: Array<{ stage: string; count: number }>;
  timeToCompletion: Array<{ days: string; count: number }>;
  strugglingStudents: Array<{
    id: string;
    name: string;
    score: number;
    classAverage: number;
    inactiveDays: number;
  }>;
  contentPerformance: Array<{
    lessonId: string;
    title: string;
    passRate: number;
    avgTimeSpent: number;
    engagement: number;
  }>;
  gradeDistribution: Array<{ range: string; count: number }>;
}

export default function InstructorAnalyticsPage() {
  const router = useRouter();
  const [analytics, setAnalytics] = useState<InstructorAnalytics | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [searchTerm, setSearchTerm] = useState('');

  useEffect(() => {
    fetchAnalytics();
  }, []);

  const fetchAnalytics = async () => {
    try {
      setLoading(true);
      // TODO: Replace with actual API call
      // const response = await api.getClassAnalytics();
      // setAnalytics(response.data);

      // Mock data for development
      setAnalytics({
        classOverview: {
          totalStudents: 150,
          activeThisWeek: 87,
          averageCompletionRate: 68,
        },
        engagementHeatmap: [
          { day: 'Mon', hour: 9, activity: 15 },
          { day: 'Mon', hour: 10, activity: 25 },
          { day: 'Mon', hour: 11, activity: 20 },
          { day: 'Tue', hour: 14, activity: 30 },
          { day: 'Tue', hour: 15, activity: 35 },
          { day: 'Wed', hour: 10, activity: 22 },
          { day: 'Wed', hour: 16, activity: 18 },
          { day: 'Thu', hour: 14, activity: 28 },
          { day: 'Thu', hour: 15, activity: 32 },
          { day: 'Fri', hour: 10, activity: 20 },
          { day: 'Fri', hour: 14, activity: 25 },
          { day: 'Sat', hour: 11, activity: 10 },
          { day: 'Sun', hour: 14, activity: 8 },
        ],
        completionFunnel: [
          { stage: 'Started', count: 150 },
          { stage: 'Submitted', count: 95 },
          { stage: 'Passed', count: 72 },
        ],
        timeToCompletion: [
          { days: '1-3', count: 25 },
          { days: '4-7', count: 40 },
          { days: '8-14', count: 30 },
          { days: '15-21', count: 15 },
          { days: '21+', count: 10 },
        ],
        strugglingStudents: [
          { id: '1', name: 'Alice Johnson', score: 52, classAverage: 68, inactiveDays: 8 },
          { id: '2', name: 'Bob Smith', score: 48, classAverage: 68, inactiveDays: 12 },
          { id: '3', name: 'Carol White', score: 55, classAverage: 68, inactiveDays: 7 },
          { id: '4', name: 'David Brown', score: 50, classAverage: 68, inactiveDays: 10 },
          { id: '5', name: 'Eve Davis', score: 58, classAverage: 68, inactiveDays: 6 },
        ],
        contentPerformance: [
          { lessonId: '1', title: 'Introduction to Go', passRate: 92, avgTimeSpent: 45, engagement: 95 },
          { lessonId: '2', title: 'Variables and Types', passRate: 88, avgTimeSpent: 60, engagement: 90 },
          { lessonId: '3', title: 'Control Structures', passRate: 75, avgTimeSpent: 75, engagement: 82 },
          { lessonId: '4', title: 'Functions', passRate: 68, avgTimeSpent: 90, engagement: 78 },
          { lessonId: '5', title: 'Concurrency Basics', passRate: 52, avgTimeSpent: 120, engagement: 70 },
        ],
        gradeDistribution: [
          { range: '90-100%', count: 20 },
          { range: '80-89%', count: 28 },
          { range: '70-79%', count: 24 },
          { range: '60-69%', count: 18 },
          { range: 'Below 60%', count: 10 },
        ],
      });
    } catch (err) {
      setError('Failed to load analytics');
      console.error(err);
    } finally {
      setLoading(false);
    }
  };

  const handleStudentClick = (studentId: string) => {
    router.push(`/cms/analytics/students/${studentId}`);
  };

  const handleExport = async (format: 'json' | 'csv') => {
    try {
      // TODO: Implement export functionality
      console.log(`Exporting analytics as ${format}`);
    } catch (err) {
      console.error('Export failed:', err);
    }
  };

  const filteredStudents = analytics?.strugglingStudents.filter((student) =>
    student.name.toLowerCase().includes(searchTerm.toLowerCase())
  );

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
              <AlertCircle className="h-5 w-5 text-red-500" />
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
          <h1 className="text-3xl font-bold tracking-tight">Class Analytics</h1>
          <p className="text-gray-600 dark:text-gray-400 mt-1">
            Monitor student progress and identify areas for intervention
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

      {/* Class Overview */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-6 mb-8">
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Total Students</CardTitle>
            <Users className="h-4 w-4 text-gray-600" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{analytics.classOverview.totalStudents}</div>
            <p className="text-xs text-gray-600 dark:text-gray-400 mt-1">
              {analytics.classOverview.activeThisWeek} active this week
            </p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Avg Completion Rate</CardTitle>
            <Target className="h-4 w-4 text-gray-600" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{analytics.classOverview.averageCompletionRate}%</div>
            <p className="text-xs text-gray-600 dark:text-gray-400 mt-1">
              Across all assessments
            </p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Struggling Students</CardTitle>
            <AlertCircle className="h-4 w-4 text-red-500" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-red-600">{analytics.strugglingStudents.length}</div>
            <p className="text-xs text-gray-600 dark:text-gray-400 mt-1">
              Require intervention
            </p>
          </CardContent>
        </Card>
      </div>

      {/* Engagement Heatmap and Completion Funnel */}
      <div className="grid grid-cols-1 md:grid-cols-2 gap-8 mb-8">
        <Card>
          <CardHeader>
            <CardTitle>Engagement Heatmap</CardTitle>
            <CardDescription>Activity by day and time of day</CardDescription>
          </CardHeader>
          <CardContent>
            <ResponsiveContainer width="100%" height={300}>
              <BarChart data={analytics.engagementHeatmap}>
                <CartesianGrid strokeDasharray="3 3" />
                <XAxis dataKey="day" />
                <YAxis label={{ value: 'Hour', angle: -90, position: 'insideLeft' }} />
                <Tooltip />
                <Bar dataKey="hour" fill="#3b82f6" name="Hour of Day" />
              </BarChart>
            </ResponsiveContainer>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>Completion Funnel</CardTitle>
            <CardDescription>Student progression through assessments</CardDescription>
          </CardHeader>
          <CardContent>
            <ResponsiveContainer width="100%" height={300}>
              <BarChart data={analytics.completionFunnel} layout="vertical">
                <CartesianGrid strokeDasharray="3 3" />
                <XAxis type="number" />
                <YAxis dataKey="stage" type="category" width={100} />
                <Tooltip />
                <Bar dataKey="count" fill="#10b981" />
              </BarChart>
            </ResponsiveContainer>
          </CardContent>
        </Card>
      </div>

      {/* Time to Completion and Grade Distribution */}
      <div className="grid grid-cols-1 md:grid-cols-2 gap-8 mb-8">
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <Clock className="h-5 w-5" />
              Time to Completion
            </CardTitle>
            <CardDescription>Distribution of completion time</CardDescription>
          </CardHeader>
          <CardContent>
            <ResponsiveContainer width="100%" height={300}>
              <BarChart data={analytics.timeToCompletion}>
                <CartesianGrid strokeDasharray="3 3" />
                <XAxis dataKey="days" />
                <YAxis />
                <Tooltip />
                <Bar dataKey="count" fill="#8b5cf6" />
              </BarChart>
            </ResponsiveContainer>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>Grade Distribution</CardTitle>
            <CardDescription>Assessment score ranges</CardDescription>
          </CardHeader>
          <CardContent>
            <ResponsiveContainer width="100%" height={300}>
              <BarChart data={analytics.gradeDistribution}>
                <CartesianGrid strokeDasharray="3 3" />
                <XAxis dataKey="range" />
                <YAxis />
                <Tooltip />
                <Bar dataKey="count" fill="#f59e0b">
                  {analytics.gradeDistribution.map((entry, index) => (
                    <Cell
                      key={`cell-${index}`}
                      fill={
                        index === 0
                          ? '#10b981'
                          : index === 1
                          ? '#22c55e'
                          : index === 2
                          ? '#84cc16'
                          : index === 3
                          ? '#eab308'
                          : '#ef4444'
                      }
                    />
                  ))}
                </Bar>
              </BarChart>
            </ResponsiveContainer>
          </CardContent>
        </Card>
      </div>

      {/* Struggling Students */}
      <Card className="mb-8">
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <AlertCircle className="h-5 w-5 text-red-500" />
            Struggling Students
          </CardTitle>
          <CardDescription>
            Students needing intervention (score below average - 1SD or inactive 7+ days)
          </CardDescription>
        </CardHeader>
        <CardContent>
          <div className="mb-4">
            <Input
              placeholder="Search students..."
              value={searchTerm}
              onChange={(e) => setSearchTerm(e.target.value)}
              className="max-w-sm"
            />
          </div>
          <div className="space-y-4">
            {filteredStudents?.map((student) => (
              <div
                key={student.id}
                className="flex items-center justify-between p-4 border rounded-lg hover:bg-gray-50 dark:hover:bg-gray-800 cursor-pointer transition-colors"
                onClick={() => handleStudentClick(student.id)}
              >
                <div>
                  <h3 className="font-semibold">{student.name}</h3>
                  <div className="flex items-center gap-4 mt-1 text-sm text-gray-600 dark:text-gray-400">
                    <span>Score: {student.score}%</span>
                    <span className="text-gray-400">|</span>
                    <span>Class Avg: {student.classAverage}%</span>
                    <span className="text-gray-400">|</span>
                    <span
                      className={student.inactiveDays >= 7 ? 'text-red-600 font-medium' : ''}
                    >
                      Inactive: {student.inactiveDays} days
                    </span>
                  </div>
                </div>
                <Badge variant={student.inactiveDays >= 7 ? 'destructive' : 'secondary'}>
                  {student.inactiveDays >= 7 ? 'Critical' : 'At Risk'}
                </Badge>
              </div>
            ))}
          </div>
        </CardContent>
      </Card>

      {/* Content Performance */}
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <TrendingUp className="h-5 w-5" />
            Content Performance
          </CardTitle>
          <CardDescription>Pass rates and engagement by lesson</CardDescription>
        </CardHeader>
        <CardContent>
          <div className="space-y-4">
            {analytics.contentPerformance.map((lesson) => (
              <div key={lesson.lessonId} className="border rounded-lg p-4">
                <div className="flex items-center justify-between mb-2">
                  <h3 className="font-semibold">{lesson.title}</h3>
                  <div className="flex items-center gap-2">
                    <Badge
                      variant={
                        lesson.passRate >= 80
                          ? 'default'
                          : lesson.passRate >= 60
                          ? 'secondary'
                          : 'destructive'
                      }
                    >
                      {lesson.passRate}% Pass Rate
                    </Badge>
                  </div>
                </div>
                <div className="grid grid-cols-2 gap-4 text-sm">
                  <div>
                    <span className="text-gray-600 dark:text-gray-400">Avg Time: </span>
                    <span className="font-medium">{lesson.avgTimeSpent} min</span>
                  </div>
                  <div>
                    <span className="text-gray-600 dark:text-gray-400">Engagement: </span>
                    <span className="font-medium">{lesson.engagement}%</span>
                  </div>
                </div>
              </div>
            ))}
          </div>
        </CardContent>
      </Card>
    </div>
  );
}
