"use client";

import { useState } from "react";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Progress } from "@/components/ui/progress";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import ProgressTracker from "./progress-tracker";
import {
  BookOpen,
  Code2,
  Trophy,
  Play,
  ArrowRight,
  Clock,
  CheckCircle,
  Star,
  Calendar,
  TrendingUp,
  Target,
  Zap,
  Award,
  Flame,
  Sparkles
} from "lucide-react";
import Link from "next/link";

interface RecentActivity {
  id: string;
  type: "lesson" | "exercise" | "achievement";
  title: string;
  description: string;
  timestamp: string;
  icon: any;
  color: string;
}

interface LearningDashboardProps {
  userId?: string;
}

export function LearningDashboard({ userId = "demo-user" }: LearningDashboardProps) {
  const [activeTab, setActiveTab] = useState("overview");

  const recentActivities: RecentActivity[] = [
    {
      id: "1",
      type: "lesson",
      title: "Completed Lesson 2",
      description: "Variables, Constants, and Functions",
      timestamp: "2 hours ago",
      icon: CheckCircle,
      color: "text-green-500",
    },
    {
      id: "2",
      type: "exercise",
      title: "Solved Exercise",
      description: "Function Implementation Challenge",
      timestamp: "3 hours ago",
      icon: Code2,
      color: "text-blue-500",
    },
    {
      id: "3",
      type: "achievement",
      title: "Earned Achievement",
      description: "Code Warrior - Write 100 lines of Go code",
      timestamp: "1 day ago",
      icon: Trophy,
      color: "text-yellow-500",
    },
    {
      id: "4",
      type: "lesson",
      title: "Started Lesson 3",
      description: "Control Structures and Loops",
      timestamp: "2 days ago",
      icon: Play,
      color: "text-primary",
    },
  ];

  const upcomingLessons = [
    {
      id: 4,
      title: "Arrays, Slices, and Maps",
      description: "Data structures, manipulation, memory considerations",
      duration: "5-6 hours",
      difficulty: "Beginner",
      locked: false,
    },
    {
      id: 5,
      title: "Pointers and Memory Management",
      description: "Pointer basics, memory allocation, garbage collection",
      duration: "4-5 hours",
      difficulty: "Beginner",
      locked: true,
    },
    {
      id: 6,
      title: "Structs and Methods",
      description: "Struct definition, methods, receivers, method sets",
      duration: "5-6 hours",
      difficulty: "Intermediate",
      locked: true,
    },
  ];

  const weeklyGoal = {
    target: 5,
    completed: 3,
    description: "Complete 5 lessons this week",
  };

  return (
    <div className="min-h-screen bg-gradient-to-br from-background via-background to-accent/5">
      <div className="container-responsive padding-responsive-y">
        {/* Enhanced Header with Stats Preview */}
        <div className="margin-responsive mb-8">
          <div className="flex flex-col lg:flex-row lg:items-center lg:justify-between gap-6">
            <div className="space-y-3">
              <h1 className="text-responsive-heading font-bold tracking-tight bg-gradient-to-r from-primary via-primary/80 to-primary/60 bg-clip-text text-transparent animate-in fade-in slide-in-from-bottom-4 duration-700">
                Learning Dashboard
              </h1>
              <p className="text-responsive-body text-muted-foreground max-w-2xl animate-in fade-in slide-in-from-bottom-5 duration-700 delay-100">
                Track your progress and continue your Go programming journey
              </p>
            </div>

            {/* Quick Stats Preview */}
            <div className="flex items-center gap-4 animate-in fade-in slide-in-from-right-4 duration-700 delay-200">
              <div className="flex items-center gap-2 px-4 py-2 rounded-full bg-gradient-to-r from-orange-500/10 to-orange-500/5 border border-orange-500/20">
                <Flame className="h-4 w-4 text-orange-500" />
                <span className="text-sm font-bold text-orange-600 dark:text-orange-400">5 Day Streak</span>
              </div>
              <div className="flex items-center gap-2 px-4 py-2 rounded-full bg-gradient-to-r from-yellow-500/10 to-yellow-500/5 border border-yellow-500/20">
                <Star className="h-4 w-4 text-yellow-500" />
                <span className="text-sm font-bold text-yellow-600 dark:text-yellow-400">350 XP</span>
              </div>
            </div>
          </div>
        </div>

      <Tabs value={activeTab} onValueChange={setActiveTab} className="space-y-6">
        <TabsList className="grid w-full grid-cols-3 lg:w-[400px] p-1 bg-muted/50 backdrop-blur-sm">
          <TabsTrigger value="overview" className="data-[state=active]:bg-background data-[state=active]:shadow-sm">
            Overview
          </TabsTrigger>
          <TabsTrigger value="progress" className="data-[state=active]:bg-background data-[state=active]:shadow-sm">
            Progress
          </TabsTrigger>
          <TabsTrigger value="activity" className="data-[state=active]:bg-background data-[state=active]:shadow-sm">
            Activity
          </TabsTrigger>
        </TabsList>

        {/* Overview Tab */}
        <TabsContent value="overview" className="space-y-6 animate-in fade-in duration-500">
          {/* Enhanced Quick Stats with Gradients */}
          <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
            <Card className="group relative overflow-hidden border-blue-200/50 dark:border-blue-800/50 hover:shadow-lg hover:shadow-blue-500/10 transition-all duration-300 hover:-translate-y-1">
              <div className="absolute inset-0 bg-gradient-to-br from-blue-500/5 via-transparent to-transparent" />
              <CardContent className="relative p-6">
                <div className="flex items-center justify-between mb-3">
                  <div className="p-3 rounded-xl bg-gradient-to-br from-blue-500/20 to-blue-500/10 group-hover:from-blue-500/30 group-hover:to-blue-500/20 transition-all duration-300">
                    <BookOpen className="h-6 w-6 text-blue-600 dark:text-blue-400" />
                  </div>
                  <Sparkles className="h-4 w-4 text-blue-400 opacity-0 group-hover:opacity-100 transition-opacity duration-300" />
                </div>
                <div className="text-3xl font-bold mb-1 bg-gradient-to-br from-blue-600 to-blue-500 bg-clip-text text-transparent">3</div>
                <div className="text-sm font-medium text-muted-foreground">Lessons Completed</div>
              </CardContent>
            </Card>

            <Card className="group relative overflow-hidden border-green-200/50 dark:border-green-800/50 hover:shadow-lg hover:shadow-green-500/10 transition-all duration-300 hover:-translate-y-1">
              <div className="absolute inset-0 bg-gradient-to-br from-green-500/5 via-transparent to-transparent" />
              <CardContent className="relative p-6">
                <div className="flex items-center justify-between mb-3">
                  <div className="p-3 rounded-xl bg-gradient-to-br from-green-500/20 to-green-500/10 group-hover:from-green-500/30 group-hover:to-green-500/20 transition-all duration-300">
                    <Code2 className="h-6 w-6 text-green-600 dark:text-green-400" />
                  </div>
                  <Sparkles className="h-4 w-4 text-green-400 opacity-0 group-hover:opacity-100 transition-opacity duration-300" />
                </div>
                <div className="text-3xl font-bold mb-1 bg-gradient-to-br from-green-600 to-green-500 bg-clip-text text-transparent">18</div>
                <div className="text-sm font-medium text-muted-foreground">Exercises Solved</div>
              </CardContent>
            </Card>

            <Card className="group relative overflow-hidden border-yellow-200/50 dark:border-yellow-800/50 hover:shadow-lg hover:shadow-yellow-500/10 transition-all duration-300 hover:-translate-y-1">
              <div className="absolute inset-0 bg-gradient-to-br from-yellow-500/5 via-transparent to-transparent" />
              <CardContent className="relative p-6">
                <div className="flex items-center justify-between mb-3">
                  <div className="p-3 rounded-xl bg-gradient-to-br from-yellow-500/20 to-yellow-500/10 group-hover:from-yellow-500/30 group-hover:to-yellow-500/20 transition-all duration-300">
                    <Star className="h-6 w-6 text-yellow-600 dark:text-yellow-400" />
                  </div>
                  <Sparkles className="h-4 w-4 text-yellow-400 opacity-0 group-hover:opacity-100 transition-opacity duration-300" />
                </div>
                <div className="text-3xl font-bold mb-1 bg-gradient-to-br from-yellow-600 to-yellow-500 bg-clip-text text-transparent">350</div>
                <div className="text-sm font-medium text-muted-foreground">XP Earned</div>
              </CardContent>
            </Card>

            <Card className="group relative overflow-hidden border-orange-200/50 dark:border-orange-800/50 hover:shadow-lg hover:shadow-orange-500/10 transition-all duration-300 hover:-translate-y-1">
              <div className="absolute inset-0 bg-gradient-to-br from-orange-500/5 via-transparent to-transparent" />
              <CardContent className="relative p-6">
                <div className="flex items-center justify-between mb-3">
                  <div className="p-3 rounded-xl bg-gradient-to-br from-orange-500/20 to-orange-500/10 group-hover:from-orange-500/30 group-hover:to-orange-500/20 transition-all duration-300">
                    <Flame className="h-6 w-6 text-orange-600 dark:text-orange-400" />
                  </div>
                  <Sparkles className="h-4 w-4 text-orange-400 opacity-0 group-hover:opacity-100 transition-opacity duration-300" />
                </div>
                <div className="text-3xl font-bold mb-1 bg-gradient-to-br from-orange-600 to-orange-500 bg-clip-text text-transparent">5</div>
                <div className="text-sm font-medium text-muted-foreground">Day Streak</div>
              </CardContent>
            </Card>
          </div>

          {/* Enhanced Weekly Goal */}
          <Card className="relative overflow-hidden border-primary/20 bg-gradient-to-br from-primary/5 via-background to-background">
            <div className="absolute top-0 right-0 w-64 h-64 bg-gradient-to-br from-primary/10 to-transparent rounded-full blur-3xl -z-10" />
            <CardHeader>
              <div className="flex items-center justify-between">
                <div>
                  <CardTitle className="flex items-center text-xl">
                    <div className="p-2 rounded-lg bg-primary/10 mr-3">
                      <Target className="h-5 w-5 text-primary" />
                    </div>
                    Weekly Goal
                  </CardTitle>
                  <CardDescription className="mt-2">{weeklyGoal.description}</CardDescription>
                </div>
                <div className="text-right">
                  <div className="text-3xl font-bold text-primary">
                    {Math.round((weeklyGoal.completed / weeklyGoal.target) * 100)}%
                  </div>
                  <div className="text-xs text-muted-foreground mt-1">Complete</div>
                </div>
              </div>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="relative">
                <Progress
                  value={(weeklyGoal.completed / weeklyGoal.target) * 100}
                  className="h-3"
                />
                <div className="flex items-center justify-between mt-3">
                  <span className="text-sm font-medium text-muted-foreground">
                    {weeklyGoal.completed} of {weeklyGoal.target} lessons
                  </span>
                  <span className="text-sm font-semibold text-primary">
                    {weeklyGoal.target - weeklyGoal.completed} to go
                  </span>
                </div>
              </div>
              {weeklyGoal.completed >= weeklyGoal.target * 0.6 && (
                <div className="flex items-center gap-2 p-3 rounded-lg bg-green-500/10 border border-green-500/20">
                  <Award className="h-4 w-4 text-green-600 dark:text-green-400" />
                  <span className="text-sm font-medium text-green-700 dark:text-green-300">
                    Great progress! You're on track to meet your goal! 🎉
                  </span>
                </div>
              )}
            </CardContent>
          </Card>

          {/* Enhanced Continue Learning */}
          <Card className="border-primary/10">
            <CardHeader>
              <div className="flex items-center justify-between">
                <div>
                  <CardTitle className="flex items-center text-xl">
                    <div className="p-2 rounded-lg bg-gradient-to-br from-primary/20 to-primary/10 mr-3">
                      <Play className="h-5 w-5 text-primary" />
                    </div>
                    Continue Learning
                  </CardTitle>
                  <CardDescription className="mt-2">Pick up where you left off</CardDescription>
                </div>
              </div>
            </CardHeader>
            <CardContent>
              <div className="space-y-4">
                {upcomingLessons.slice(0, 2).map((lesson, index) => (
                  <div
                    key={lesson.id}
                    className="group relative overflow-hidden rounded-2xl bg-gradient-to-br from-muted/50 via-muted/30 to-background border border-border/50 hover:border-primary/40 hover:shadow-xl hover:shadow-primary/5 transition-all duration-500 hover:-translate-y-1"
                  >
                    <div className="absolute inset-0 bg-gradient-to-br from-primary/5 via-transparent to-transparent opacity-0 group-hover:opacity-100 transition-opacity duration-500" />
                    <div className="relative flex items-center justify-between p-5">
                      <div className="flex items-start space-x-4 flex-1">
                        <div className="relative">
                          <div className="p-4 rounded-2xl bg-gradient-to-br from-primary/20 via-primary/15 to-primary/10 group-hover:from-primary/30 group-hover:via-primary/25 group-hover:to-primary/20 transition-all duration-500 group-hover:scale-110">
                            <BookOpen className="h-6 w-6 text-primary" />
                          </div>
                          {index === 0 && (
                            <div className="absolute -top-1 -right-1 w-3 h-3 bg-green-500 rounded-full border-2 border-background animate-pulse" />
                          )}
                        </div>
                        <div className="flex-1 space-y-3">
                          <div className="flex items-center gap-2 flex-wrap">
                            <span className="text-xs font-bold text-primary/80 tracking-wider uppercase px-2 py-1 rounded-md bg-primary/10">
                              Lesson {lesson.id}
                            </span>
                            <div className="h-1 w-1 rounded-full bg-primary/30" />
                            <Badge
                              variant="outline"
                              className="text-xs font-medium border-primary/30 text-primary"
                            >
                              {lesson.difficulty}
                            </Badge>
                            {index === 0 && (
                              <Badge className="text-xs bg-green-500/10 text-green-700 dark:text-green-400 border-green-500/30">
                                In Progress
                              </Badge>
                            )}
                          </div>
                          <h3 className="font-bold text-lg leading-tight group-hover:text-primary transition-colors duration-300">
                            {lesson.title}
                          </h3>
                          <p className="text-sm text-muted-foreground leading-relaxed line-clamp-2">
                            {lesson.description}
                          </p>
                          <div className="flex items-center gap-4 pt-1">
                            <div className="flex items-center gap-2 text-xs text-muted-foreground">
                              <Clock className="h-4 w-4" />
                              <span className="font-medium">{lesson.duration}</span>
                            </div>
                            {index === 0 && (
                              <div className="flex items-center gap-2">
                                <Progress value={35} className="w-20 h-1.5" />
                                <span className="text-xs font-medium text-primary">35%</span>
                              </div>
                            )}
                          </div>
                        </div>
                      </div>
                      <Link href={lesson.locked ? "#" : `/learn/${lesson.id}`} className="ml-4">
                        <Button
                          size="lg"
                          disabled={lesson.locked}
                          className="group-hover:shadow-lg group-hover:shadow-primary/20 transition-all duration-300 group-hover:scale-105"
                        >
                          {lesson.locked ? "Locked" : index === 0 ? "Continue" : "Start"}
                          {!lesson.locked && <ArrowRight className="ml-2 h-4 w-4 group-hover:translate-x-1 transition-transform" />}
                        </Button>
                      </Link>
                    </div>
                  </div>
                ))}
              </div>
              <div className="mt-6 text-center">
                <Link href="/curriculum">
                  <Button variant="outline" size="lg" className="group">
                    View Full Curriculum
                    <ArrowRight className="ml-2 h-4 w-4 group-hover:translate-x-1 transition-transform" />
                  </Button>
                </Link>
              </div>
            </CardContent>
          </Card>
        </TabsContent>

        {/* Progress Tab */}
        <TabsContent value="progress">
          <ProgressTracker userId={userId} />
        </TabsContent>

        {/* Enhanced Activity Tab */}
        <TabsContent value="activity" className="space-y-6 animate-in fade-in duration-500">
          <Card className="border-primary/10">
            <CardHeader>
              <CardTitle className="flex items-center text-xl">
                <div className="p-2 rounded-lg bg-gradient-to-br from-primary/20 to-primary/10 mr-3">
                  <TrendingUp className="h-5 w-5 text-primary" />
                </div>
                Recent Activity
              </CardTitle>
              <CardDescription className="mt-2">Your learning activity over the past week</CardDescription>
            </CardHeader>
            <CardContent>
              <div className="relative space-y-4">
                {/* Timeline line */}
                <div className="absolute left-[29px] top-4 bottom-4 w-0.5 bg-gradient-to-b from-primary/50 via-primary/30 to-transparent" />

                {recentActivities.map((activity, index) => (
                  <div
                    key={activity.id}
                    className="relative flex items-start space-x-4 p-4 rounded-xl bg-gradient-to-br from-muted/50 to-background border border-border/50 hover:border-primary/30 hover:shadow-lg transition-all duration-300 hover:-translate-y-0.5 group"
                    style={{ animationDelay: `${index * 100}ms` }}
                  >
                    <div className="relative z-10">
                      <div className={`p-3 rounded-xl bg-gradient-to-br ${
                        activity.type === 'lesson' ? 'from-blue-500/20 to-blue-500/10' :
                        activity.type === 'exercise' ? 'from-green-500/20 to-green-500/10' :
                        'from-yellow-500/20 to-yellow-500/10'
                      } group-hover:scale-110 transition-transform duration-300`}>
                        <activity.icon className={`h-5 w-5 ${activity.color}`} />
                      </div>
                    </div>
                    <div className="flex-1 min-w-0">
                      <div className="flex items-start justify-between gap-2 mb-1">
                        <div className="font-semibold text-base group-hover:text-primary transition-colors">
                          {activity.title}
                        </div>
                        <Badge variant="outline" className="text-xs shrink-0">
                          {activity.type}
                        </Badge>
                      </div>
                      <div className="text-sm text-muted-foreground mb-2 leading-relaxed">
                        {activity.description}
                      </div>
                      <div className="flex items-center gap-2 text-xs text-muted-foreground">
                        <Calendar className="h-3.5 w-3.5" />
                        <span className="font-medium">{activity.timestamp}</span>
                      </div>
                    </div>
                  </div>
                ))}
              </div>

              {/* View More Button */}
              <div className="mt-6 text-center">
                <Button variant="outline" className="group">
                  View All Activity
                  <ArrowRight className="ml-2 h-4 w-4 group-hover:translate-x-1 transition-transform" />
                </Button>
              </div>
            </CardContent>
          </Card>
        </TabsContent>
      </Tabs>
      </div>
    </div>
  );
}

export default LearningDashboard;
