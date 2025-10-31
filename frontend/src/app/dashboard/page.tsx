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
    <div className="min-h-screen">
      {/* Hero Section */}
      <section className="relative overflow-hidden animated-gradient">
        {/* Decorative elements */}
        <div className="absolute inset-0 overflow-hidden pointer-events-none">
          <div className="absolute -top-40 -right-40 w-96 h-96 bg-gradient-to-br from-primary/20 to-blue-500/20 rounded-full blur-3xl float-animation" />
          <div className="absolute -bottom-40 -left-40 w-96 h-96 bg-gradient-to-tr from-cyan-500/20 to-primary/20 rounded-full blur-3xl float-animation" style={{ animationDelay: '1s' }} />
        </div>

        <div className="container-responsive py-12 sm:py-16 lg:py-20 relative z-10">
          <div className="flex flex-col lg:flex-row lg:items-center lg:justify-between">
            <div className="mb-6 lg:mb-0">
              <Badge variant="secondary" className="mb-4 text-sm pulse-badge animate-in fade-in slide-in-bottom duration-500">
                👋 Welcome Back
              </Badge>
              <h1 className="text-3xl lg:text-4xl xl:text-5xl font-bold tracking-tight mb-3 animate-in fade-in slide-in-bottom duration-700">
                Your <span className="go-gradient-text">Learning Dashboard</span>
              </h1>
              <p className="text-base lg:text-xl text-muted-foreground max-w-2xl leading-relaxed animate-in fade-in slide-in-bottom duration-1000">
                Continue your Go learning journey and track your progress
              </p>
            </div>
            <div className="flex items-center space-x-2 animate-in fade-in duration-1000">
              <Button variant="outline" size="icon" className="glass-card hover:shadow-lg transition-all duration-300">
                <Bell className="h-4 w-4" />
              </Button>
              <Button variant="outline" size="icon" className="glass-card hover:shadow-lg transition-all duration-300">
                <Settings className="h-4 w-4" />
              </Button>
            </div>
          </div>
        </div>
      </section>

      {/* Main Content */}
      <div className="container-responsive padding-responsive-y relative overflow-hidden">
        {/* Background gradient */}
        <div className="absolute inset-0 bg-gradient-to-b from-background via-accent/5 to-background pointer-events-none" />

      <Tabs defaultValue="overview" className="space-y-6 relative z-10 animate-in fade-in duration-1000">
        <TabsList className="grid w-full grid-cols-4 lg:w-[400px] bg-card/50 backdrop-blur-sm border border-border/50 shadow-lg">
          <TabsTrigger value="overview" className="data-[state=active]:bg-primary data-[state=active]:text-primary-foreground">Overview</TabsTrigger>
          <TabsTrigger value="courses" className="data-[state=active]:bg-primary data-[state=active]:text-primary-foreground">Courses</TabsTrigger>
          <TabsTrigger value="analytics" className="data-[state=active]:bg-primary data-[state=active]:text-primary-foreground">Analytics</TabsTrigger>
          <TabsTrigger value="explore" className="data-[state=active]:bg-primary data-[state=active]:text-primary-foreground">Explore</TabsTrigger>
        </TabsList>

        <TabsContent value="overview" className="space-y-6">
          <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
            <div className="lg:col-span-2">
              <CourseProgress {...courseData} />
            </div>
            <div className="space-y-6">
              <div className="grid grid-cols-2 gap-4">
                <div className="text-center p-4 sm:p-6 rounded-xl glass-card hover:shadow-2xl hover:shadow-primary/20 transition-all duration-500 hover:-translate-y-2 relative overflow-hidden group">
                  <div className="absolute inset-0 bg-gradient-to-br from-primary/10 via-transparent to-blue-500/10 opacity-0 group-hover:opacity-100 transition-opacity duration-500" />
                  <div className="relative z-10">
                    <div className="p-3 rounded-xl bg-gradient-to-br from-primary/20 to-primary/10 w-fit mx-auto mb-2 group-hover:scale-110 transition-transform duration-300">
                      <TrendingUp className="h-6 w-6 sm:h-8 sm:w-8 text-primary group-hover:animate-pulse" />
                    </div>
                    <div className="text-2xl font-bold bg-gradient-to-br from-foreground via-primary/80 to-foreground/70 bg-clip-text text-transparent group-hover:scale-110 transition-transform duration-300">7</div>
                    <div className="text-sm text-muted-foreground mt-1 group-hover:text-foreground/80 transition-colors">Day Streak</div>
                  </div>
                </div>
                <div className="text-center p-4 sm:p-6 rounded-xl glass-card hover:shadow-2xl hover:shadow-green-500/20 transition-all duration-500 hover:-translate-y-2 relative overflow-hidden group">
                  <div className="absolute inset-0 bg-gradient-to-br from-green-500/10 via-transparent to-emerald-500/10 opacity-0 group-hover:opacity-100 transition-opacity duration-500" />
                  <div className="relative z-10">
                    <div className="p-3 rounded-xl bg-gradient-to-br from-green-500/20 to-green-500/10 w-fit mx-auto mb-2 group-hover:scale-110 transition-transform duration-300">
                      <Trophy className="h-6 w-6 sm:h-8 sm:w-8 text-green-500 group-hover:animate-pulse" />
                    </div>
                    <div className="text-2xl font-bold bg-gradient-to-br from-foreground via-green-500/80 to-foreground/70 bg-clip-text text-transparent group-hover:scale-110 transition-transform duration-300">2</div>
                    <div className="text-sm text-muted-foreground mt-1 group-hover:text-foreground/80 transition-colors">Completed</div>
                  </div>
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
