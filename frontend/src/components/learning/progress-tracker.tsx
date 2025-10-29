"use client";

import { useState, useEffect } from "react";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Progress } from "@/components/ui/progress";
import {
  CheckCircle,
  Clock,
  Trophy,
  Target,
  Star,
  TrendingUp,
  BookOpen,
  Code2,
  Award,
  Zap
} from "lucide-react";

interface ProgressStats {
  totalLessons: number;
  completedLessons: number;
  totalExercises: number;
  completedExercises: number;
  totalProjects: number;
  completedProjects: number;
  totalXP: number;
  currentStreak: number;
  achievements: number;
}

interface ProgressTrackerProps {
  userId?: string;
  className?: string;
}

export function ProgressTracker({ userId = "demo-user", className = "" }: ProgressTrackerProps) {
  const [stats, setStats] = useState<ProgressStats>({
    totalLessons: 15,
    completedLessons: 3,
    totalExercises: 120,
    completedExercises: 18,
    totalProjects: 4,
    completedProjects: 0,
    totalXP: 350,
    currentStreak: 5,
    achievements: 4,
  });

  const [loading, setLoading] = useState(false);

  const progressPercentage = Math.round((stats.completedLessons / stats.totalLessons) * 100);
  const exerciseProgress = Math.round((stats.completedExercises / stats.totalExercises) * 100);

  const progressItems = [
    {
      icon: BookOpen,
      label: "Lessons",
      current: stats.completedLessons,
      total: stats.totalLessons,
      color: "text-blue-500",
      bgColor: "bg-blue-50 dark:bg-blue-950",
      borderColor: "border-blue-200 dark:border-blue-800",
    },
    {
      icon: Code2,
      label: "Exercises",
      current: stats.completedExercises,
      total: stats.totalExercises,
      color: "text-green-500",
      bgColor: "bg-green-50 dark:bg-green-950",
      borderColor: "border-green-200 dark:border-green-800",
    },
    {
      icon: Trophy,
      label: "Projects",
      current: stats.completedProjects,
      total: stats.totalProjects,
      color: "text-yellow-500",
      bgColor: "bg-yellow-50 dark:bg-yellow-950",
      borderColor: "border-yellow-200 dark:border-yellow-800",
    },
    {
      icon: Award,
      label: "Achievements",
      current: stats.achievements,
      total: 20,
      color: "text-purple-500",
      bgColor: "bg-purple-50 dark:bg-purple-950",
      borderColor: "border-purple-200 dark:border-purple-800",
    },
  ];

  return (
    <div className={`space-y-6 ${className} animate-in fade-in duration-500`}>
      {/* Enhanced Overall Progress */}
      <Card className="relative overflow-hidden border-primary/20 bg-gradient-to-br from-primary/5 via-background to-background">
        <div className="absolute top-0 right-0 w-96 h-96 bg-gradient-to-br from-primary/10 to-transparent rounded-full blur-3xl -z-10" />
        <CardHeader>
          <div className="flex flex-col lg:flex-row lg:items-center lg:justify-between gap-4">
            <div>
              <CardTitle className="flex items-center text-xl mb-2">
                <div className="p-2 rounded-lg bg-primary/10 mr-3">
                  <TrendingUp className="h-5 w-5 text-primary" />
                </div>
                Learning Progress
              </CardTitle>
              <CardDescription>Your journey through GO-PRO curriculum</CardDescription>
            </div>
            <div className="text-center lg:text-right">
              <div className="text-4xl font-bold bg-gradient-to-br from-primary to-primary/70 bg-clip-text text-transparent">
                {progressPercentage}%
              </div>
              <div className="text-sm text-muted-foreground mt-1">Complete</div>
            </div>
          </div>
        </CardHeader>
        <CardContent className="space-y-4">
          <div className="relative">
            <Progress value={progressPercentage} className="h-3" />
          </div>
          <div className="flex items-center justify-between text-sm">
            <span className="font-medium text-muted-foreground">
              {stats.completedLessons} of {stats.totalLessons} lessons completed
            </span>
            <span className="font-semibold text-primary">
              {stats.totalLessons - stats.completedLessons} remaining
            </span>
          </div>
          {progressPercentage >= 20 && (
            <div className="flex items-center gap-2 p-3 rounded-lg bg-primary/10 border border-primary/20">
              <Trophy className="h-4 w-4 text-primary" />
              <span className="text-sm font-medium text-primary">
                You're making great progress! Keep it up! 🚀
              </span>
            </div>
          )}
        </CardContent>
      </Card>

      {/* Enhanced Progress Grid */}
      <div className="grid grid-cols-2 lg:grid-cols-4 gap-4">
        {progressItems.map((item, index) => (
          <Card
            key={index}
            className={`group relative overflow-hidden ${item.bgColor} ${item.borderColor} hover:shadow-lg transition-all duration-300 hover:-translate-y-1`}
          >
            <div className="absolute inset-0 bg-gradient-to-br from-white/50 to-transparent dark:from-white/5 opacity-0 group-hover:opacity-100 transition-opacity duration-300" />
            <CardContent className="relative p-5 text-center">
              <div className="flex justify-center mb-3">
                <div className="p-3 rounded-xl bg-gradient-to-br from-white/50 to-white/20 dark:from-black/20 dark:to-black/10 group-hover:scale-110 transition-transform duration-300">
                  <item.icon className={`h-6 w-6 ${item.color}`} />
                </div>
              </div>
              <div className="text-3xl font-bold mb-1">{item.current}</div>
              <div className="text-xs font-medium text-muted-foreground mb-3">
                of {item.total} {item.label}
              </div>
              <div className="space-y-2">
                <Progress
                  value={(item.current / item.total) * 100}
                  className="h-2"
                />
                <div className="text-xs font-semibold" style={{ color: item.color.replace('text-', '') }}>
                  {Math.round((item.current / item.total) * 100)}%
                </div>
              </div>
            </CardContent>
          </Card>
        ))}
      </div>

      {/* Enhanced Stats Cards */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
        <Card className="group relative overflow-hidden bg-gradient-to-br from-yellow-500/10 via-yellow-500/5 to-background border-yellow-500/20 hover:shadow-lg hover:shadow-yellow-500/10 transition-all duration-300 hover:-translate-y-1">
          <div className="absolute inset-0 bg-gradient-to-br from-yellow-500/10 to-transparent opacity-0 group-hover:opacity-100 transition-opacity duration-300" />
          <CardContent className="relative p-6 text-center">
            <div className="flex justify-center mb-3">
              <div className="p-3 rounded-xl bg-gradient-to-br from-yellow-500/20 to-yellow-500/10 group-hover:scale-110 transition-transform duration-300">
                <Star className="h-6 w-6 text-yellow-600 dark:text-yellow-400" />
              </div>
            </div>
            <div className="text-3xl font-bold mb-1 bg-gradient-to-br from-yellow-600 to-yellow-500 bg-clip-text text-transparent">
              {stats.totalXP}
            </div>
            <div className="text-sm font-medium text-muted-foreground">Total XP</div>
            <div className="mt-3 text-xs text-yellow-600 dark:text-yellow-400 font-medium">
              +50 XP this week
            </div>
          </CardContent>
        </Card>

        <Card className="group relative overflow-hidden bg-gradient-to-br from-orange-500/10 via-orange-500/5 to-background border-orange-500/20 hover:shadow-lg hover:shadow-orange-500/10 transition-all duration-300 hover:-translate-y-1">
          <div className="absolute inset-0 bg-gradient-to-br from-orange-500/10 to-transparent opacity-0 group-hover:opacity-100 transition-opacity duration-300" />
          <CardContent className="relative p-6 text-center">
            <div className="flex justify-center mb-3">
              <div className="p-3 rounded-xl bg-gradient-to-br from-orange-500/20 to-orange-500/10 group-hover:scale-110 transition-transform duration-300">
                <Zap className="h-6 w-6 text-orange-600 dark:text-orange-400" />
              </div>
            </div>
            <div className="text-3xl font-bold mb-1 bg-gradient-to-br from-orange-600 to-orange-500 bg-clip-text text-transparent">
              {stats.currentStreak}
            </div>
            <div className="text-sm font-medium text-muted-foreground">Day Streak</div>
            <div className="mt-3 text-xs text-orange-600 dark:text-orange-400 font-medium">
              🔥 Keep it going!
            </div>
          </CardContent>
        </Card>

        <Card className="group relative overflow-hidden bg-gradient-to-br from-green-500/10 via-green-500/5 to-background border-green-500/20 hover:shadow-lg hover:shadow-green-500/10 transition-all duration-300 hover:-translate-y-1">
          <div className="absolute inset-0 bg-gradient-to-br from-green-500/10 to-transparent opacity-0 group-hover:opacity-100 transition-opacity duration-300" />
          <CardContent className="relative p-6 text-center">
            <div className="flex justify-center mb-3">
              <div className="p-3 rounded-xl bg-gradient-to-br from-green-500/20 to-green-500/10 group-hover:scale-110 transition-transform duration-300">
                <Target className="h-6 w-6 text-green-600 dark:text-green-400" />
              </div>
            </div>
            <div className="text-3xl font-bold mb-1 bg-gradient-to-br from-green-600 to-green-500 bg-clip-text text-transparent">
              {exerciseProgress}%
            </div>
            <div className="text-sm font-medium text-muted-foreground">Exercise Progress</div>
            <div className="mt-3 text-xs text-green-600 dark:text-green-400 font-medium">
              {stats.completedExercises} completed
            </div>
          </CardContent>
        </Card>
      </div>

      {/* Enhanced Next Steps */}
      <Card className="border-primary/10">
        <CardHeader>
          <CardTitle className="flex items-center text-xl">
            <div className="p-2 rounded-lg bg-gradient-to-br from-primary/20 to-primary/10 mr-3">
              <Clock className="h-5 w-5 text-primary" />
            </div>
            Next Steps
          </CardTitle>
          <CardDescription className="mt-2">Continue your learning journey</CardDescription>
        </CardHeader>
        <CardContent>
          <div className="space-y-3">
            <div className="group flex items-center justify-between p-4 rounded-xl bg-gradient-to-br from-primary/5 via-primary/3 to-background border border-primary/20 hover:border-primary/40 hover:shadow-lg transition-all duration-300 hover:-translate-y-0.5">
              <div className="flex items-center space-x-4 flex-1">
                <div className="relative">
                  <div className="p-3 rounded-xl bg-gradient-to-br from-primary/20 to-primary/10 group-hover:scale-110 transition-transform duration-300">
                    <BookOpen className="h-5 w-5 text-primary" />
                  </div>
                  <div className="absolute -top-1 -right-1 w-3 h-3 bg-green-500 rounded-full border-2 border-background animate-pulse" />
                </div>
                <div className="flex-1">
                  <div className="font-semibold text-base mb-1 group-hover:text-primary transition-colors">
                    Continue Lesson 4
                  </div>
                  <div className="text-sm text-muted-foreground">Arrays, Slices, and Maps</div>
                  <div className="flex items-center gap-2 mt-2">
                    <Progress value={35} className="w-24 h-1.5" />
                    <span className="text-xs font-medium text-primary">35%</span>
                  </div>
                </div>
              </div>
              <Button size="sm" className="group-hover:shadow-md transition-all duration-300">
                Continue
              </Button>
            </div>

            <div className="group flex items-center justify-between p-4 rounded-xl bg-gradient-to-br from-green-500/5 via-green-500/3 to-background border border-green-500/20 hover:border-green-500/40 hover:shadow-lg transition-all duration-300 hover:-translate-y-0.5">
              <div className="flex items-center space-x-4 flex-1">
                <div className="p-3 rounded-xl bg-gradient-to-br from-green-500/20 to-green-500/10 group-hover:scale-110 transition-transform duration-300">
                  <Code2 className="h-5 w-5 text-green-600 dark:text-green-400" />
                </div>
                <div className="flex-1">
                  <div className="font-semibold text-base mb-1 group-hover:text-green-600 dark:group-hover:text-green-400 transition-colors">
                    Practice Exercise
                  </div>
                  <div className="text-sm text-muted-foreground">Variable Declaration Challenge</div>
                  <Badge variant="outline" className="mt-2 text-xs border-green-500/30 text-green-600 dark:text-green-400">
                    Recommended
                  </Badge>
                </div>
              </div>
              <Button size="sm" variant="outline" className="group-hover:shadow-md transition-all duration-300">
                Practice
              </Button>
            </div>

            <div className="group flex items-center justify-between p-4 rounded-xl bg-muted/30 border border-border/50 opacity-60">
              <div className="flex items-center space-x-4 flex-1">
                <div className="p-3 rounded-xl bg-gradient-to-br from-yellow-500/20 to-yellow-500/10">
                  <Trophy className="h-5 w-5 text-yellow-600 dark:text-yellow-400" />
                </div>
                <div className="flex-1">
                  <div className="font-semibold text-base mb-1">First Project</div>
                  <div className="text-sm text-muted-foreground">CLI Task Manager</div>
                  <div className="text-xs text-muted-foreground mt-2">
                    Complete 2 more lessons to unlock
                  </div>
                </div>
              </div>
              <Badge variant="outline" className="text-xs">
                Locked
              </Badge>
            </div>
          </div>
        </CardContent>
      </Card>
    </div>
  );
}

export default ProgressTracker;
