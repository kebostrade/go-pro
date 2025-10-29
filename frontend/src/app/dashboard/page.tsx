"use client";

import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import CourseProgress from "@/components/dashboard/course-progress";
import LearningAnalytics from "@/components/dashboard/learning-analytics";
import LessonCard from "@/components/dashboard/lesson-card";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { 
  BookOpen, 
  TrendingUp, 
  Trophy, 
  Settings,
  Bell
} from "lucide-react";

export default function Dashboard() {
  // Sample data - in a real app, this would come from an API
  const courseData = {
    courseTitle: "Go Fundamentals",
    courseDescription: "Master the basics of Go programming language with hands-on exercises and real-world projects.",
    totalLessons: 12,
    completedLessons: 8,
    estimatedTime: "2 weeks",
    difficulty: "Beginner" as const,
    lessons: [
      {
        id: "1",
        title: "Introduction to Go",
        duration: "15 min",
        completed: true,
        locked: false,
        type: "lesson" as const,
      },
      {
        id: "2", 
        title: "Variables and Data Types",
        duration: "25 min",
        completed: true,
        locked: false,
        type: "lesson" as const,
      },
      {
        id: "3",
        title: "Control Structures",
        duration: "30 min", 
        completed: true,
        locked: false,
        type: "exercise" as const,
      },
      {
        id: "4",
        title: "Functions and Methods",
        duration: "35 min",
        completed: false,
        locked: false,
        type: "lesson" as const,
      },
      {
        id: "5",
        title: "Build a CLI Tool",
        duration: "45 min",
        completed: false,
        locked: true,
        type: "project" as const,
      },
    ]
  };

  const analyticsData = {
    totalHours: 24,
    streak: 7,
    completedCourses: 2,
    skillLevel: "Intermediate",
    weeklyProgress: [
      { day: "Mon", hours: 2 },
      { day: "Tue", hours: 1.5 },
      { day: "Wed", hours: 3 },
      { day: "Thu", hours: 2.5 },
      { day: "Fri", hours: 1 },
      { day: "Sat", hours: 4 },
      { day: "Sun", hours: 2 },
    ],
    skillBreakdown: [
      { skill: "Go Syntax", level: 8, maxLevel: 10 },
      { skill: "Concurrency", level: 5, maxLevel: 10 },
      { skill: "Web Development", level: 6, maxLevel: 10 },
      { skill: "Testing", level: 4, maxLevel: 10 },
      { skill: "Microservices", level: 2, maxLevel: 10 },
    ],
    achievements: [
      {
        id: "1",
        title: "First Steps",
        description: "Complete your first Go lesson",
        icon: "first-lesson",
        earned: true,
        earnedDate: "2 weeks ago",
      },
      {
        id: "2",
        title: "Week Warrior",
        description: "Maintain a 7-day learning streak",
        icon: "streak-7",
        earned: true,
        earnedDate: "1 week ago",
      },
      {
        id: "3",
        title: "Project Master",
        description: "Complete your first Go project",
        icon: "first-project",
        earned: false,
      },
    ]
  };

  const featuredLessons = [
    {
      id: "advanced-1",
      title: "Goroutines and Channels",
      description: "Master Go's concurrency model with practical examples and exercises.",
      duration: "45 min",
      difficulty: "Advanced" as const,
      type: "lesson" as const,
      completed: false,
      locked: false,
      progress: 0,
      rating: 4.8,
      enrolledCount: 1250,
      tags: ["concurrency", "goroutines", "channels"],
    },
    {
      id: "project-1",
      title: "Build a REST API",
      description: "Create a complete REST API with authentication, database integration, and testing.",
      duration: "2 hours",
      difficulty: "Intermediate" as const,
      type: "project" as const,
      completed: false,
      locked: false,
      progress: 35,
      rating: 4.9,
      enrolledCount: 890,
      tags: ["api", "http", "database", "testing"],
    },
    {
      id: "exercise-1",
      title: "Algorithm Challenges",
      description: "Solve common programming problems using Go's unique features and idioms.",
      duration: "30 min",
      difficulty: "Intermediate" as const,
      type: "exercise" as const,
      completed: true,
      locked: false,
      rating: 4.6,
      enrolledCount: 2100,
      tags: ["algorithms", "problem-solving", "practice"],
    },
  ];

  return (
    <div className="min-h-screen bg-gradient-to-br from-background via-background to-accent/5">
      <div className="container-responsive padding-responsive-y">
        {/* Header */}
        <div className="flex flex-col lg:flex-row lg:items-center lg:justify-between margin-responsive">
          <div className="mb-4 lg:mb-0">
            <h1 className="text-responsive-heading font-bold tracking-tight mb-3 bg-gradient-to-r from-primary to-primary/70 bg-clip-text text-transparent">
              Dashboard
            </h1>
            <p className="text-responsive-body text-muted-foreground max-w-2xl">
              Welcome back! Continue your Go learning journey.
            </p>
          </div>
          <div className="flex items-center space-x-2">
            <Button variant="outline" size="icon">
              <Bell className="h-4 w-4" />
            </Button>
            <Button variant="outline" size="icon">
              <Settings className="h-4 w-4" />
            </Button>
          </div>
        </div>

      {/* Main Content */}
      <Tabs defaultValue="overview" className="space-y-6">
        <TabsList className="grid w-full grid-cols-4 lg:w-[400px]">
          <TabsTrigger value="overview">Overview</TabsTrigger>
          <TabsTrigger value="courses">Courses</TabsTrigger>
          <TabsTrigger value="analytics">Analytics</TabsTrigger>
          <TabsTrigger value="explore">Explore</TabsTrigger>
        </TabsList>

        <TabsContent value="overview" className="space-y-6">
          <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
            <div className="lg:col-span-2">
              <CourseProgress {...courseData} />
            </div>
            <div className="space-y-6">
              <div className="grid grid-cols-2 gap-4">
                <div className="text-center p-4 rounded-lg bg-primary/5 border border-primary/20">
                  <TrendingUp className="h-8 w-8 text-primary mx-auto mb-2" />
                  <div className="text-2xl font-bold">7</div>
                  <div className="text-sm text-muted-foreground">Day Streak</div>
                </div>
                <div className="text-center p-4 rounded-lg bg-green-50 border border-green-200 dark:bg-green-950 dark:border-green-800">
                  <Trophy className="h-8 w-8 text-green-500 mx-auto mb-2" />
                  <div className="text-2xl font-bold">2</div>
                  <div className="text-sm text-muted-foreground">Completed</div>
                </div>
              </div>
              
              <div className="space-y-3">
                <h3 className="font-semibold">Quick Actions</h3>
                <Button className="w-full go-gradient text-white">
                  <BookOpen className="mr-2 h-4 w-4" />
                  Continue Learning
                </Button>
                <Button variant="outline" className="w-full">
                  <Trophy className="mr-2 h-4 w-4" />
                  View Achievements
                </Button>
              </div>
            </div>
          </div>
        </TabsContent>

        <TabsContent value="courses">
          <CourseProgress {...courseData} />
        </TabsContent>

        <TabsContent value="analytics">
          <LearningAnalytics data={analyticsData} />
        </TabsContent>

        <TabsContent value="explore" className="space-y-6">
          <div>
            <h2 className="text-2xl font-bold mb-4">Explore New Content</h2>
            <p className="text-muted-foreground mb-6">
              Discover new lessons, projects, and challenges to advance your Go skills.
            </p>
          </div>
          
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            {featuredLessons.map((lesson) => (
              <LessonCard
                key={lesson.id}
                {...lesson}
                onStart={() => console.log(`Starting lesson: ${lesson.title}`)}
                onContinue={() => console.log(`Continuing lesson: ${lesson.title}`)}
                onReview={() => console.log(`Reviewing lesson: ${lesson.title}`)}
              />
            ))}
          </div>
        </TabsContent>
      </Tabs>
      </div>
    </div>
  );
}
