'use client';

import { useEffect, useState } from 'react';
import { useAuth } from '@/contexts/auth-context';
import { api, Progress, ProgressStats, Curriculum } from '@/lib/api';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Progress as ProgressBar } from '@/components/ui/progress';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Skeleton } from '@/components/ui/skeleton';
import {
  BookOpen,
  TrendingUp,
  Trophy,
  Clock,
  Target,
  Flame,
  Award,
  PlayCircle,
  CheckCircle2,
  Circle,
  ArrowRight,
  Calendar,
  BarChart3,
  Zap,
} from 'lucide-react';

interface ActivityItem {
  id: string;
  lesson_id: string;
  lesson_title: string;
  status: 'not_started' | 'in_progress' | 'completed';
  score: number;
  timestamp: string;
}

export default function Dashboard() {
  // Handle SSR - auth is client-side only
  let authContext;
  try {
    authContext = useAuth();
  } catch {
    authContext = { user: null, userProfile: null };
  }
  const { user, userProfile } = authContext;
  const [loading, setLoading] = useState(true);
  const [progressData, setProgressData] = useState<Progress[]>([]);
  const [stats, setStats] = useState<ProgressStats | null>(null);
  const [curriculum, setCurriculum] = useState<Curriculum | null>(null);
  const [recentActivity, setRecentActivity] = useState<ActivityItem[]>([]);

  useEffect(() => {
    async function loadDashboardData() {
      if (!user) return;

      try {
        setLoading(true);

        // Load all data in parallel
        const [progressRes, statsRes, curriculumRes] = await Promise.all([
          api.getUserProgress(user.uid, 1, 20),
          api.getProgressStats(user.uid),
          api.getCurriculum(),
        ]);

        setProgressData(progressRes.progress);
        setStats(statsRes);
        setCurriculum(curriculumRes);

        // Generate recent activity from progress data
        const activities: ActivityItem[] = progressRes.progress
          .sort((a, b) => {
            const timeA = a.completed_at || a.started_at || '';
            const timeB = b.completed_at || b.started_at || '';
            return new Date(timeB).getTime() - new Date(timeA).getTime();
          })
          .slice(0, 5)
          .map((p) => ({
            id: p.id,
            lesson_id: p.lesson_id,
            lesson_title: getLessonTitle(curriculumRes, p.lesson_id),
            status: p.status,
            score: p.score,
            timestamp: p.completed_at || p.started_at || '',
          }));

        setRecentActivity(activities);
      } catch (error) {
        console.error('Failed to load dashboard data:', error);
      } finally {
        setLoading(false);
      }
    }

    loadDashboardData();
  }, [user]);

  // Helper to get lesson title from curriculum
  function getLessonTitle(curriculum: Curriculum, lessonId: string): string {
    for (const phase of curriculum.phases) {
      const lesson = phase.lessons.find((l) => l.id.toString() === lessonId);
      if (lesson) return lesson.title;
    }
    return 'Unknown Lesson';
  }

  // Calculate learning streak (consecutive days)
  const calculateStreak = (): number => {
    if (!progressData.length) return 0;

    const dates = progressData
      .map((p) => p.completed_at || p.started_at)
      .filter(Boolean)
      .map((d) => new Date(d!).toDateString())
      .sort((a, b) => new Date(b).getTime() - new Date(a).getTime());

    if (!dates.length) return 0;

    let streak = 1;
    const today = new Date().toDateString();
    const yesterday = new Date(Date.now() - 86400000).toDateString();

    if (dates[0] !== today && dates[0] !== yesterday) return 0;

    for (let i = 1; i < dates.length; i++) {
      const prevDate = new Date(dates[i - 1]);
      const currDate = new Date(dates[i]);
      const dayDiff = Math.floor((prevDate.getTime() - currDate.getTime()) / 86400000);

      if (dayDiff === 1) {
        streak++;
      } else {
        break;
      }
    }

    return streak;
  };

  // Get last lesson for "Continue Learning" CTA
  const getLastLesson = () => {
    const inProgressLessons = progressData.filter((p) => p.status === 'in_progress');
    return inProgressLessons.length > 0 ? inProgressLessons[0] : null;
  };

  // Calculate phase completion percentages
  const getPhaseStats = () => {
    if (!curriculum) return [];

    return curriculum.phases.map((phase) => {
      const totalLessons = phase.lessons.length;
      const completedLessons = phase.lessons.filter((l) =>
        progressData.some((p) => p.lesson_id === l.id.toString() && p.status === 'completed')
      ).length;

      return {
        id: phase.id,
        title: phase.title,
        color: phase.color,
        total: totalLessons,
        completed: completedLessons,
        percentage: totalLessons > 0 ? Math.round((completedLessons / totalLessons) * 100) : 0,
      };
    });
  };

  // Mock achievements based on progress
  const getAchievements = () => {
    const achievements = [];
    const completedCount = stats?.completed_lessons || 0;

    if (completedCount >= 1) {
      achievements.push({
        id: 'first-lesson',
        icon: '🎯',
        title: 'First Steps',
        description: 'Completed your first lesson',
      });
    }

    if (calculateStreak() >= 3) {
      achievements.push({
        id: 'streak-3',
        icon: '🔥',
        title: 'On Fire',
        description: '3-day learning streak',
      });
    }

    if (calculateStreak() >= 7) {
      achievements.push({
        id: 'streak-7',
        icon: '⚡',
        title: 'Week Warrior',
        description: '7-day learning streak',
      });
    }

    if (completedCount >= 5) {
      achievements.push({
        id: 'lesson-5',
        icon: '🏆',
        title: 'Knowledge Seeker',
        description: 'Completed 5 lessons',
      });
    }

    if (completedCount >= 10) {
      achievements.push({
        id: 'lesson-10',
        icon: '💎',
        title: 'Rising Star',
        description: 'Completed 10 lessons',
      });
    }

    if (stats && stats.average_score >= 80) {
      achievements.push({
        id: 'high-scorer',
        icon: '⭐',
        title: 'High Achiever',
        description: 'Average score above 80%',
      });
    }

    return achievements;
  };

  const streak = calculateStreak();
  const lastLesson = getLastLesson();
  const phaseStats = getPhaseStats();
  const achievements = getAchievements();

  if (loading) {
    return (
      <div className="min-h-screen bg-gradient-to-b from-background via-accent/5 to-background">
        <div className="container mx-auto px-4 py-8 max-w-7xl">
          <div className="space-y-8">
            <Skeleton className="h-32 w-full" />
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
              {[...Array(4)].map((_, i) => (
                <Skeleton key={i} className="h-32" />
              ))}
            </div>
            <Skeleton className="h-96 w-full" />
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gradient-to-b from-background via-accent/5 to-background">
      {/* Header */}
      <div className="border-b bg-card/50 backdrop-blur-sm">
        <div className="container mx-auto px-4 py-6 max-w-7xl">
          <div className="flex flex-col md:flex-row md:items-center md:justify-between gap-4">
            <div>
              <h1 className="text-3xl md:text-4xl font-bold tracking-tight">
                Welcome back, <span className="bg-gradient-to-r from-primary to-blue-600 bg-clip-text text-transparent">{userProfile?.displayName || 'Learner'}</span>
              </h1>
              <p className="text-muted-foreground mt-2">Track your progress and continue your Go journey</p>
            </div>
            {lastLesson && (
              <Button size="lg" className="bg-gradient-to-r from-primary to-blue-600 hover:opacity-90 transition-opacity">
                <PlayCircle className="mr-2 h-5 w-5" />
                Continue Learning
                <ArrowRight className="ml-2 h-4 w-4" />
              </Button>
            )}
          </div>
        </div>
      </div>

      <div className="container mx-auto px-4 py-8 max-w-7xl">
        <div className="space-y-8">
          {/* Stats Overview */}
          <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4 md:gap-6">
            {/* Total Lessons */}
            <Card className="border-none shadow-lg hover:shadow-xl transition-all duration-300 hover:-translate-y-1">
              <CardContent className="p-6">
                <div className="flex items-center justify-between">
                  <div className="flex-1">
                    <p className="text-sm font-medium text-muted-foreground">Total Lessons</p>
                    <p className="text-3xl font-bold mt-2">{stats?.total_lessons || 0}</p>
                    <p className="text-xs text-muted-foreground mt-1">Available to learn</p>
                  </div>
                  <div className="p-3 rounded-full bg-blue-100 dark:bg-blue-900/20">
                    <BookOpen className="h-6 w-6 text-blue-600 dark:text-blue-400" />
                  </div>
                </div>
              </CardContent>
            </Card>

            {/* Completed */}
            <Card className="border-none shadow-lg hover:shadow-xl transition-all duration-300 hover:-translate-y-1">
              <CardContent className="p-6">
                <div className="flex items-center justify-between">
                  <div className="flex-1">
                    <p className="text-sm font-medium text-muted-foreground">Completed</p>
                    <p className="text-3xl font-bold mt-2">{stats?.completed_lessons || 0}</p>
                    <p className="text-xs text-green-600 dark:text-green-400 mt-1">
                      {stats?.total_lessons && stats.total_lessons > 0
                        ? `${Math.round((stats.completed_lessons / stats.total_lessons) * 100)}% progress`
                        : '0% progress'}
                    </p>
                  </div>
                  <div className="p-3 rounded-full bg-green-100 dark:bg-green-900/20">
                    <CheckCircle2 className="h-6 w-6 text-green-600 dark:text-green-400" />
                  </div>
                </div>
              </CardContent>
            </Card>

            {/* Learning Streak */}
            <Card className="border-none shadow-lg hover:shadow-xl transition-all duration-300 hover:-translate-y-1">
              <CardContent className="p-6">
                <div className="flex items-center justify-between">
                  <div className="flex-1">
                    <p className="text-sm font-medium text-muted-foreground">Learning Streak</p>
                    <p className="text-3xl font-bold mt-2">{streak}</p>
                    <p className="text-xs text-orange-600 dark:text-orange-400 mt-1">
                      {streak > 0 ? 'Days in a row' : 'Start learning today'}
                    </p>
                  </div>
                  <div className="p-3 rounded-full bg-orange-100 dark:bg-orange-900/20">
                    <Flame className="h-6 w-6 text-orange-600 dark:text-orange-400" />
                  </div>
                </div>
              </CardContent>
            </Card>

            {/* Average Score */}
            <Card className="border-none shadow-lg hover:shadow-xl transition-all duration-300 hover:-translate-y-1">
              <CardContent className="p-6">
                <div className="flex items-center justify-between">
                  <div className="flex-1">
                    <p className="text-sm font-medium text-muted-foreground">Average Score</p>
                    <p className="text-3xl font-bold mt-2">{Math.round(stats?.average_score || 0)}%</p>
                    <p className="text-xs text-purple-600 dark:text-purple-400 mt-1">
                      {stats?.average_score && stats.average_score >= 80 ? 'Excellent!' : 'Keep improving'}
                    </p>
                  </div>
                  <div className="p-3 rounded-full bg-purple-100 dark:bg-purple-900/20">
                    <Target className="h-6 w-6 text-purple-600 dark:text-purple-400" />
                  </div>
                </div>
              </CardContent>
            </Card>
          </div>

          {/* Main Content Grid */}
          <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
            {/* Left Column - Progress & Activity */}
            <div className="lg:col-span-2 space-y-6">
              {/* Curriculum Progress */}
              <Card className="border-none shadow-lg">
                <CardHeader>
                  <div className="flex items-center justify-between">
                    <div>
                      <CardTitle className="flex items-center gap-2">
                        <BarChart3 className="h-5 w-5 text-primary" />
                        Curriculum Progress
                      </CardTitle>
                      <CardDescription className="mt-1">Track your progress across all learning phases</CardDescription>
                    </div>
                  </div>
                </CardHeader>
                <CardContent className="space-y-6">
                  {phaseStats.map((phase) => (
                    <div key={phase.id} className="space-y-2">
                      <div className="flex items-center justify-between text-sm">
                        <div className="flex items-center gap-2">
                          <div className={`w-3 h-3 rounded-full`} style={{ backgroundColor: phase.color }} />
                          <span className="font-medium">{phase.title}</span>
                        </div>
                        <span className="text-muted-foreground">
                          {phase.completed}/{phase.total} lessons
                        </span>
                      </div>
                      <ProgressBar value={phase.percentage} className="h-2" />
                      <p className="text-xs text-muted-foreground text-right">{phase.percentage}% complete</p>
                    </div>
                  ))}
                </CardContent>
              </Card>

              {/* Recent Activity */}
              <Card className="border-none shadow-lg">
                <CardHeader>
                  <CardTitle className="flex items-center gap-2">
                    <Clock className="h-5 w-5 text-primary" />
                    Recent Activity
                  </CardTitle>
                  <CardDescription>Your last 5 lesson interactions</CardDescription>
                </CardHeader>
                <CardContent>
                  <div className="space-y-4">
                    {recentActivity.length === 0 ? (
                      <div className="text-center py-8">
                        <Circle className="h-12 w-12 text-muted-foreground mx-auto mb-4 opacity-50" />
                        <p className="text-muted-foreground">No recent activity yet</p>
                        <p className="text-sm text-muted-foreground mt-1">Start learning to see your activity here</p>
                      </div>
                    ) : (
                      recentActivity.map((activity) => (
                        <div key={activity.id} className="flex items-center gap-4 p-3 rounded-lg hover:bg-accent/50 transition-colors">
                          <div className={`p-2 rounded-full ${
                            activity.status === 'completed' ? 'bg-green-100 dark:bg-green-900/20' :
                            activity.status === 'in_progress' ? 'bg-blue-100 dark:bg-blue-900/20' :
                            'bg-gray-100 dark:bg-gray-900/20'
                          }`}>
                            {activity.status === 'completed' ? (
                              <CheckCircle2 className="h-4 w-4 text-green-600 dark:text-green-400" />
                            ) : activity.status === 'in_progress' ? (
                              <PlayCircle className="h-4 w-4 text-blue-600 dark:text-blue-400" />
                            ) : (
                              <Circle className="h-4 w-4 text-gray-600 dark:text-gray-400" />
                            )}
                          </div>
                          <div className="flex-1 min-w-0">
                            <p className="font-medium text-sm truncate">{activity.lesson_title}</p>
                            <p className="text-xs text-muted-foreground">
                              {activity.status === 'completed' ? 'Completed' :
                               activity.status === 'in_progress' ? 'In Progress' :
                               'Not Started'}
                              {activity.score > 0 && ` • Score: ${activity.score}%`}
                            </p>
                          </div>
                          <div className="text-xs text-muted-foreground whitespace-nowrap">
                            {activity.timestamp ? new Date(activity.timestamp).toLocaleDateString() : 'N/A'}
                          </div>
                        </div>
                      ))
                    )}
                  </div>
                </CardContent>
              </Card>
            </div>

            {/* Right Column - Achievements & Stats */}
            <div className="space-y-6">
              {/* Continue Learning CTA */}
              {lastLesson && (
                <Card className="border-none shadow-lg bg-gradient-to-br from-primary/10 via-background to-blue-500/10">
                  <CardContent className="p-6">
                    <div className="space-y-4">
                      <div className="flex items-center gap-2">
                        <Zap className="h-5 w-5 text-primary" />
                        <h3 className="font-semibold">Continue Learning</h3>
                      </div>
                      <div>
                        <p className="text-sm text-muted-foreground mb-1">Last lesson</p>
                        <p className="font-medium">{getLessonTitle(curriculum!, lastLesson.lesson_id)}</p>
                      </div>
                      <Button className="w-full bg-gradient-to-r from-primary to-blue-600" onClick={() => window.location.href = `/learn/${lastLesson.lesson_id}`}>
                        <PlayCircle className="mr-2 h-4 w-4" />
                        Resume
                      </Button>
                    </div>
                  </CardContent>
                </Card>
              )}

              {/* Achievements */}
              <Card className="border-none shadow-lg">
                <CardHeader>
                  <CardTitle className="flex items-center gap-2">
                    <Trophy className="h-5 w-5 text-yellow-500" />
                    Achievements
                  </CardTitle>
                  <CardDescription>Earned badges</CardDescription>
                </CardHeader>
                <CardContent>
                  <div className="grid grid-cols-3 gap-3">
                    {achievements.length === 0 ? (
                      <div className="col-span-3 text-center py-6">
                        <Award className="h-12 w-12 text-muted-foreground mx-auto mb-2 opacity-50" />
                        <p className="text-sm text-muted-foreground">Complete lessons to earn badges</p>
                      </div>
                    ) : (
                      achievements.map((achievement) => (
                        <div key={achievement.id} className="group relative">
                          <div className="aspect-square rounded-lg bg-gradient-to-br from-yellow-100 to-orange-100 dark:from-yellow-900/20 dark:to-orange-900/20 flex items-center justify-center text-3xl hover:scale-110 transition-transform cursor-pointer">
                            {achievement.icon}
                          </div>
                          <div className="absolute inset-0 rounded-lg bg-black/80 opacity-0 group-hover:opacity-100 transition-opacity p-2 flex flex-col items-center justify-center text-center">
                            <p className="text-xs font-medium text-white">{achievement.title}</p>
                            <p className="text-xs text-gray-300 mt-1">{achievement.description}</p>
                          </div>
                        </div>
                      ))
                    )}
                  </div>
                </CardContent>
              </Card>

              {/* Quick Stats */}
              <Card className="border-none shadow-lg">
                <CardHeader>
                  <CardTitle className="flex items-center gap-2">
                    <TrendingUp className="h-5 w-5 text-primary" />
                    Quick Stats
                  </CardTitle>
                </CardHeader>
                <CardContent className="space-y-4">
                  <div className="flex items-center justify-between p-3 rounded-lg bg-accent/50">
                    <span className="text-sm text-muted-foreground">In Progress</span>
                    <span className="font-bold text-lg">{stats?.in_progress_lessons || 0}</span>
                  </div>
                  <div className="flex items-center justify-between p-3 rounded-lg bg-accent/50">
                    <span className="text-sm text-muted-foreground">Total Time</span>
                    <span className="font-bold text-lg">
                      {stats?.total_time_spent ? `${Math.round(stats.total_time_spent / 60)}h` : '0h'}
                    </span>
                  </div>
                  <div className="flex items-center justify-between p-3 rounded-lg bg-accent/50">
                    <span className="text-sm text-muted-foreground">XP Level</span>
                    <span className="font-bold text-lg">
                      {userProfile?.progress?.level || 1}
                    </span>
                  </div>
                </CardContent>
              </Card>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
