"use client";

import { useState, useEffect } from "react";
import { useParams } from "next/navigation";
import { api } from "@/lib/api";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Progress } from "@/components/ui/progress";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import CodeEditor from "@/components/learning/code-editor";
import ExerciseSubmission from "@/components/learning/exercise-submission";
import LessonNotes from "@/components/learning/lesson-notes";
import LessonProgress from "@/components/learning/lesson-progress";
import MarkdownRenderer from "@/components/learning/markdown-renderer";
import InteractiveExample from "@/components/learning/interactive-example";
import TableOfContents from "@/components/learning/table-of-contents";
import RelatedLessons from "@/components/learning/related-lessons";
import {
  BookOpen,
  Code2,
  Trophy,
  ArrowRight,
  ArrowLeft,
  Clock,
  Target,
  CheckCircle,
  Home,
  List,
  ChevronRight,
  Lightbulb,
  FileText,
  Bookmark,
  BookmarkCheck,
  Star,
  Award,
  TrendingUp,
  BarChart3,
  Flame,
  Menu,
  X,
  StickyNote
} from "lucide-react";
import Link from "next/link";
import "../../../styles/lesson-animations.css";
import { performanceMonitor, defaultPerformanceBudget, checkPerformanceBudget, checkAccessibility } from "@/lib/performance";

interface LessonData {
  id: number;
  title: string;
  description: string;
  duration: string;
  difficulty: "Beginner" | "Intermediate" | "Advanced";
  phase: string;
  objectives: string[];
  theory: string;
  codeExample: string;
  solution: string;
  exercises: Exercise[];
  nextLessonId?: number;
  prevLessonId?: number;
}

interface Exercise {
  id: string;
  title: string;
  description: string;
  requirements: string[];
  initialCode: string;
  solution: string;
}

export default function LessonPage() {
  const params = useParams();
  const [activeTab, setActiveTab] = useState("lesson");
  const [lessonData, setLessonData] = useState<LessonData | null>(null);
  const [loading, setLoading] = useState(true);

  // Enhanced state for new features
  const [isBookmarked, setIsBookmarked] = useState(false);
  const [lessonProgress, setLessonProgress] = useState(0);
  const [showSidebar, setShowSidebar] = useState(false);
  const [completedObjectives, setCompletedObjectives] = useState<Set<number>>(new Set());
  const [timeSpent, setTimeSpent] = useState(0);
  const [startTime] = useState(Date.now());

  const lessonId = parseInt(params.id as string);

  useEffect(() => {
    const loadLessonData = async () => {
      setLoading(true);

      try {
        const lesson = await api.getLessonDetail(lessonId);

        // Convert API response to component format
        const lessonData: LessonData = {
          id: lesson.id,
          title: lesson.title,
          description: lesson.description,
          duration: lesson.duration,
          difficulty: lesson.difficulty.charAt(0).toUpperCase() + lesson.difficulty.slice(1) as "Beginner" | "Intermediate" | "Advanced",
          phase: lesson.phase,
          objectives: lesson.objectives,
          theory: lesson.theory,
          codeExample: lesson.code_example,
          solution: lesson.solution,
          exercises: lesson.exercises.map(ex => ({
            id: ex.id,
            title: ex.title,
            description: ex.description,
            requirements: ex.requirements,
            initialCode: ex.initial_code,
            solution: ex.solution,
          })),
          nextLessonId: lesson.next_lesson_id,
          prevLessonId: lesson.prev_lesson_id,
        };

        setLessonData(lessonData);
      } catch (error) {
        console.error('Failed to load lesson:', error);
        setLessonData(null);
      } finally {
        setLoading(false);
      }
    };

    loadLessonData();
  }, [lessonId]);

  // Time tracking effect
  useEffect(() => {
    const interval = setInterval(() => {
      setTimeSpent(Math.floor((Date.now() - startTime) / 1000));
    }, 1000);

    return () => clearInterval(interval);
  }, [startTime]);

  // Performance monitoring
  useEffect(() => {
    // Measure initial load time
    performanceMonitor.measureLoadTime();

    // Check performance budget after component mounts
    setTimeout(() => {
      const metrics = performanceMonitor.getMetrics();
      const budgetCheck = checkPerformanceBudget(metrics, defaultPerformanceBudget);

      if (!budgetCheck.passed) {
        console.warn('Performance budget exceeded:', budgetCheck.violations);
      }

      // Run accessibility check in development
      if (process.env.NODE_ENV === 'development') {
        checkAccessibility();
      }
    }, 1000);
  }, []);

  // Debounced progress update for performance (currently unused but may be needed for future optimization)
  // const debouncedProgressUpdate = debounce((progress: number) => {
  //   setLessonProgress(progress);
  // }, 300);

  // Fallback mock data for development (currently unused but kept for future reference)
  /*
  const getMockLessonData = (): LessonData => {
    return {
      id: 1,
      title: "Go Syntax and Basic Types",
      description: "Learn the fundamental syntax of Go and work with basic data types including integers, floats, strings, and booleans.",
      duration: "3-4 hours",
      difficulty: "Beginner",
      phase: "Foundations",
      objectives: [
        "Set up a Go development environment",
        "Understand Go's basic syntax and program structure",
        "Work with primitive data types (int, float, string, bool)",
        "Declare and use constants",
        "Perform type conversions",
        "Use the iota identifier for enumerated constants"
      ],
          theory: `
# Go Program Structure

Every Go program starts with a package declaration, followed by imports, and then the program code:

\`\`\`go
package main

import "fmt"

func main() {
    fmt.Println("Hello, Go!")
}
\`\`\`

## Basic Types

Go has several built-in basic types:

### Numeric Types
- **Integers**: int, int8, int16, int32, int64
- **Unsigned integers**: uint, uint8, uint16, uint32, uint64
- **Floating point**: float32, float64
- **Complex numbers**: complex64, complex128

### Other Types
- **Boolean**: bool (true or false)
- **String**: string (UTF-8 encoded)
- **Byte**: byte (alias for uint8)
- **Rune**: rune (alias for int32, represents Unicode code points)

## Variable Declarations

\`\`\`go
// Explicit type declaration
var name string = "Go"
var age int = 10

// Type inference
var language = "Go"
var version = 1.21

// Short variable declaration (inside functions only)
message := "Hello, World!"
count := 42
\`\`\`

## Constants

\`\`\`go
const Pi = 3.14159
const Language = "Go"

// Enumerated constants with iota
const (
    Sunday = iota    // 0
    Monday           // 1
    Tuesday          // 2
    Wednesday        // 3
    Thursday         // 4
    Friday           // 5
    Saturday         // 6
)
\`\`\`
          `,
          codeExample: `package main

import "fmt"

func main() {
    // Basic variable declarations
    var name string = "Go Programming"
    var version float64 = 1.21
    var isAwesome bool = true
    
    // Short variable declaration
    year := 2024
    
    // Constants
    const MaxUsers = 1000
    
    // Type conversion
    var x int = 42
    var y float64 = float64(x)
    
    // Print values
    fmt.Printf("Language: %s\\n", name)
    fmt.Printf("Version: %.2f\\n", version)
    fmt.Printf("Year: %d\\n", year)
    fmt.Printf("Is Awesome: %t\\n", isAwesome)
    fmt.Printf("Max Users: %d\\n", MaxUsers)
    fmt.Printf("Converted: %.1f\\n", y)
}`,
          solution: `package main

import "fmt"

func main() {
    // Basic variable declarations
    var name string = "Go Programming"
    var version float64 = 1.21
    var isAwesome bool = true
    
    // Short variable declaration
    year := 2024
    
    // Constants
    const MaxUsers = 1000
    
    // Type conversion
    var x int = 42
    var y float64 = float64(x)
    
    // Print values
    fmt.Printf("Language: %s\\n", name)
    fmt.Printf("Version: %.2f\\n", version)
    fmt.Printf("Year: %d\\n", year)
    fmt.Printf("Is Awesome: %t\\n", isAwesome)
    fmt.Printf("Max Users: %d\\n", MaxUsers)
    fmt.Printf("Converted: %.1f\\n", y)
    
    // Additional examples
    
    // Multiple variable declaration
    var (
        firstName = "John"
        lastName  = "Doe"
        age       = 30
    )
    
    fmt.Printf("Full Name: %s %s, Age: %d\\n", firstName, lastName, age)
    
    // Enumerated constants
    const (
        Red = iota
        Green
        Blue
    )
    
    fmt.Printf("Colors: Red=%d, Green=%d, Blue=%d\\n", Red, Green, Blue)
}`,
          exercises: [
            {
              id: "basic-variables",
              title: "Variable Declaration Practice",
              description: "Practice declaring variables of different types and using type conversions.",
              requirements: [
                "Declare a string variable for your name",
                "Declare an integer variable for your age",
                "Declare a boolean variable for whether you like programming",
                "Use short variable declaration for the current year",
                "Convert an integer to float64 and print both values"
              ],
              initialCode: `package main

import "fmt"

func main() {
    // TODO: Declare your variables here
    
    // TODO: Print the values
    
}`,
              solution: `package main

import "fmt"

func main() {
    // Variable declarations
    var name string = "Alice"
    var age int = 25
    var likesProgramming bool = true
    currentYear := 2024
    
    // Type conversion
    var score int = 95
    var percentage float64 = float64(score)
    
    // Print values
    fmt.Printf("Name: %s\\n", name)
    fmt.Printf("Age: %d\\n", age)
    fmt.Printf("Likes Programming: %t\\n", likesProgramming)
    fmt.Printf("Current Year: %d\\n", currentYear)
    fmt.Printf("Score: %d, Percentage: %.1f%%\\n", score, percentage)
}`
            }
          ],
      nextLessonId: 2
    };
  };
  */

  if (loading) {
    return (
      <div className="min-h-screen bg-gradient-to-br from-background via-background to-accent/5">
        <div className="container-responsive padding-responsive-y">
          <div className="flex items-center justify-center min-h-[60vh]">
            <div className="text-center">
              <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-primary mx-auto mb-4"></div>
              <p className="text-responsive text-muted-foreground">Loading lesson...</p>
            </div>
          </div>
        </div>
      </div>
    );
  }

  if (!lessonData) {
    return (
      <div className="min-h-screen bg-gradient-to-br from-background via-background to-accent/5">
        <div className="container-responsive padding-responsive-y">
          <div className="text-center py-16">
            <h1 className="text-responsive-heading font-bold mb-4">Lesson Not Found</h1>
            <p className="text-responsive text-muted-foreground mb-6">The lesson you're looking for doesn't exist.</p>
            <Link href="/curriculum">
              <Button size="lg">
                <ArrowLeft className="mr-2 h-4 w-4" />
                Back to Curriculum
              </Button>
            </Link>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen animated-gradient">
      <div className="container-responsive padding-responsive-y">
        {/* Enhanced Breadcrumb Navigation */}
        <div className="mb-6 lg:mb-8 animate-in fade-in slide-in-right duration-500">
          {/* Top Navigation Bar */}
          <div className="flex items-center justify-between mb-4">
            <div className="flex items-center space-x-2 text-sm text-muted-foreground">
              <Link href="/" className="hover:text-primary transition-colors">
                <Home className="h-4 w-4" />
              </Link>
              <ChevronRight className="h-4 w-4" />
              <Link href="/curriculum" className="hover:text-primary transition-colors">
                Curriculum
              </Link>
              <ChevronRight className="h-4 w-4" />
              <span className="text-foreground font-medium">Lesson {lessonData.id}</span>
            </div>

            {/* Quick Actions */}
            <div className="flex items-center space-x-2">
              <Button
                variant="ghost"
                size="sm"
                onClick={() => setIsBookmarked(!isBookmarked)}
                className="hover:bg-primary/10"
              >
                {isBookmarked ? (
                  <BookmarkCheck className="h-4 w-4 text-primary" />
                ) : (
                  <Bookmark className="h-4 w-4" />
                )}
              </Button>
              <Button
                variant="ghost"
                size="sm"
                onClick={() => setShowSidebar(!showSidebar)}
                className="lg:hidden hover:bg-primary/10"
              >
                {showSidebar ? <X className="h-4 w-4" /> : <Menu className="h-4 w-4" />}
              </Button>
            </div>
          </div>

          {/* Progress Indicator */}
          <div className="bg-card/50 backdrop-blur-sm border border-border/50 rounded-xl p-4">
            <div className="flex items-center justify-between mb-2">
              <div className="flex items-center space-x-2">
                <TrendingUp className="h-4 w-4 text-primary" />
                <span className="text-sm font-medium">Lesson Progress</span>
              </div>
              <div className="flex items-center space-x-4 text-xs text-muted-foreground">
                <div className="flex items-center space-x-1">
                  <Clock className="h-3 w-3" />
                  <span>{Math.floor(timeSpent / 60)}m</span>
                </div>
                <div className="flex items-center space-x-1">
                  <Target className="h-3 w-3" />
                  <span>{completedObjectives.size}/{lessonData.objectives.length}</span>
                </div>
              </div>
            </div>
            <Progress value={lessonProgress} className="h-2" />
          </div>
        </div>

      {/* Enhanced Lesson Header */}
      <div className="mb-8 glass-card-strong p-6 lg:p-8 rounded-2xl animate-in fade-in slide-in-bottom duration-700 relative overflow-hidden hover-lift particle-bg">
        {/* Enhanced Background Gradients */}
        <div className="absolute top-0 right-0 w-96 h-96 bg-gradient-to-br from-primary/20 to-transparent rounded-full blur-3xl -z-10 float-animation" />
        <div className="absolute bottom-0 left-0 w-64 h-64 bg-gradient-to-tr from-secondary/15 to-transparent rounded-full blur-2xl -z-10" />
        <div className="absolute inset-0 shimmer opacity-5 -z-10" />

        <div className="flex items-start justify-between mb-6">
          <div className="flex-1">
            {/* Enhanced Badges */}
            <div className="flex items-center flex-wrap gap-2 mb-4">
              <Badge variant="outline" className="shadow-sm hover:shadow-md transition-all hover-lift stagger-item hover-glow">
                <BookOpen className="mr-1 h-3 w-3" />
                Lesson {lessonData.id}
              </Badge>
              <Badge variant="secondary" className="shadow-sm hover:shadow-md transition-all hover-lift stagger-item">
                {lessonData.phase}
              </Badge>
              <Badge
                variant={lessonData.difficulty === 'Beginner' ? 'success' :
                        lessonData.difficulty === 'Intermediate' ? 'warning' : 'info'}
                className="shadow-sm hover:shadow-md transition-all hover-lift stagger-item pulse-glow"
              >
                {lessonData.difficulty}
              </Badge>
              {isBookmarked && (
                <Badge variant="outline" className="border-primary text-primary shadow-sm bounce-in hover-glow">
                  <Star className="mr-1 h-3 w-3 fill-current" />
                  Bookmarked
                </Badge>
              )}
            </div>

            {/* Enhanced Title */}
            <h1 className="text-3xl lg:text-4xl font-bold tracking-tight mb-4 bg-gradient-to-r from-primary via-blue-500 to-purple-600 bg-clip-text text-transparent stagger-item">
              {lessonData.title}
            </h1>

            {/* Enhanced Description */}
            <p className="text-muted-foreground text-base lg:text-lg mb-6 leading-relaxed max-w-4xl stagger-item">
              {lessonData.description}
            </p>

            {/* Enhanced Stats Grid */}
            <div className="grid grid-cols-2 lg:grid-cols-4 gap-4">
              <div className="flex items-center space-x-2 bg-gradient-to-r from-primary/10 to-primary/5 px-4 py-3 rounded-xl border border-primary/20 hover:border-primary/30 transition-all hover-lift stagger-item hover-glow">
                <Clock className="h-5 w-5 text-primary float-animation" />
                <div>
                  <div className="text-sm font-medium">{lessonData.duration}</div>
                  <div className="text-xs text-muted-foreground">Duration</div>
                </div>
              </div>

              <div className="flex items-center space-x-2 bg-gradient-to-r from-blue-500/10 to-blue-500/5 px-4 py-3 rounded-xl border border-blue-500/20 hover:border-blue-500/30 transition-all hover-lift stagger-item hover-glow">
                <Target className="h-5 w-5 text-blue-500 float-animation" />
                <div>
                  <div className="text-sm font-medium">{completedObjectives.size}/{lessonData.objectives.length}</div>
                  <div className="text-xs text-muted-foreground">Objectives</div>
                </div>
              </div>

              <div className="flex items-center space-x-2 bg-gradient-to-r from-green-500/10 to-green-500/5 px-4 py-3 rounded-xl border border-green-500/20 hover:border-green-500/30 transition-colors">
                <Award className="h-5 w-5 text-green-500" />
                <div>
                  <div className="text-sm font-medium">{Math.round(lessonProgress)}%</div>
                  <div className="text-xs text-muted-foreground">Complete</div>
                </div>
              </div>

              <div className="flex items-center space-x-2 bg-gradient-to-r from-orange-500/10 to-orange-500/5 px-4 py-3 rounded-xl border border-orange-500/20 hover:border-orange-500/30 transition-colors">
                <Flame className="h-5 w-5 text-orange-500" />
                <div>
                  <div className="text-sm font-medium">{Math.floor(timeSpent / 60)}m</div>
                  <div className="text-xs text-muted-foreground">Time Spent</div>
                </div>
              </div>
            </div>
          </div>

          {/* Enhanced Action Buttons */}
          <div className="ml-6 hidden lg:flex flex-col space-y-2">
            <Link href="/curriculum">
              <Button variant="outline" size="sm" className="shadow-sm hover:shadow-md transition-all">
                <List className="mr-2 h-4 w-4" />
                All Lessons
              </Button>
            </Link>
            <Button
              variant="ghost"
              size="sm"
              className="shadow-sm hover:shadow-md transition-all"
              onClick={() => setIsBookmarked(!isBookmarked)}
            >
              {isBookmarked ? (
                <BookmarkCheck className="mr-2 h-4 w-4 text-primary" />
              ) : (
                <Bookmark className="mr-2 h-4 w-4" />
              )}
              {isBookmarked ? 'Bookmarked' : 'Bookmark'}
            </Button>
          </div>
        </div>
      </div>

      {/* Lesson Sidebar (Mobile) */}
      {showSidebar && (
        <div className="fixed inset-0 z-50 lg:hidden">
          <div className="absolute inset-0 bg-background/80 backdrop-blur-sm" onClick={() => setShowSidebar(false)} />
          <div className="absolute right-0 top-0 h-full w-80 bg-card border-l border-border shadow-2xl">
            <div className="p-6">
              <div className="flex items-center justify-between mb-6">
                <h3 className="text-lg font-semibold">Lesson Navigation</h3>
                <Button variant="ghost" size="sm" onClick={() => setShowSidebar(false)}>
                  <X className="h-4 w-4" />
                </Button>
              </div>

              {/* Quick Stats */}
              <div className="space-y-3 mb-6">
                <div className="flex items-center justify-between text-sm">
                  <span className="text-muted-foreground">Progress</span>
                  <span className="font-medium">{Math.round(lessonProgress)}%</span>
                </div>
                <Progress value={lessonProgress} className="h-2" />

                <div className="flex items-center justify-between text-sm">
                  <span className="text-muted-foreground">Objectives</span>
                  <span className="font-medium">{completedObjectives.size}/{lessonData.objectives.length}</span>
                </div>

                <div className="flex items-center justify-between text-sm">
                  <span className="text-muted-foreground">Time Spent</span>
                  <span className="font-medium">{Math.floor(timeSpent / 60)}m</span>
                </div>
              </div>

              {/* Tab Navigation */}
              <div className="space-y-2 mb-6">
                <h4 className="text-sm font-medium text-muted-foreground">Sections</h4>
                <div className="space-y-1">
                  <Button
                    variant={activeTab === "lesson" ? "default" : "ghost"}
                    size="sm"
                    className="w-full justify-start"
                    onClick={() => {
                      setActiveTab("lesson");
                      setShowSidebar(false);
                    }}
                  >
                    <Lightbulb className="mr-2 h-4 w-4" />
                    Theory
                  </Button>
                  <Button
                    variant={activeTab === "practice" ? "default" : "ghost"}
                    size="sm"
                    className="w-full justify-start"
                    onClick={() => {
                      setActiveTab("practice");
                      setShowSidebar(false);
                    }}
                  >
                    <Code2 className="mr-2 h-4 w-4" />
                    Practice
                  </Button>
                  <Button
                    variant={activeTab === "exercise" ? "default" : "ghost"}
                    size="sm"
                    className="w-full justify-start"
                    onClick={() => {
                      setActiveTab("exercise");
                      setShowSidebar(false);
                    }}
                  >
                    <Trophy className="mr-2 h-4 w-4" />
                    Exercise
                  </Button>
                </div>
              </div>

              {/* Lesson Objectives Checklist */}
              <div className="space-y-2">
                <h4 className="text-sm font-medium text-muted-foreground">Objectives</h4>
                <div className="space-y-2">
                  {lessonData.objectives.map((objective, index) => (
                    <div
                      key={index}
                      className="flex items-start space-x-2 p-2 rounded-lg hover:bg-muted/50 cursor-pointer transition-colors"
                      onClick={() => {
                        const newCompleted = new Set(completedObjectives);
                        if (completedObjectives.has(index)) {
                          newCompleted.delete(index);
                        } else {
                          newCompleted.add(index);
                        }
                        setCompletedObjectives(newCompleted);
                        setLessonProgress((newCompleted.size / lessonData.objectives.length) * 100);
                      }}
                    >
                      <div className={`p-1 rounded-full ${completedObjectives.has(index) ? 'bg-green-500/20' : 'bg-muted'}`}>
                        <CheckCircle className={`h-3 w-3 ${completedObjectives.has(index) ? 'text-green-500' : 'text-muted-foreground'}`} />
                      </div>
                      <span className={`text-xs leading-relaxed ${completedObjectives.has(index) ? 'line-through text-muted-foreground' : ''}`}>
                        {objective}
                      </span>
                    </div>
                  ))}
                </div>
              </div>
            </div>
          </div>
        </div>
      )}

      {/* Lesson Content */}
      <div className="flex gap-6">
        {/* Desktop Sidebar */}
        <div className="hidden lg:block w-80 flex-shrink-0">
          <div className="sticky top-6 space-y-6">
            {/* Progress Card */}
            <Card className="glass-card">
              <CardHeader className="pb-3">
                <CardTitle className="text-lg flex items-center">
                  <BarChart3 className="mr-2 h-5 w-5 text-primary" />
                  Your Progress
                </CardTitle>
              </CardHeader>
              <CardContent className="space-y-4">
                <div className="space-y-2">
                  <div className="flex items-center justify-between text-sm">
                    <span className="text-muted-foreground">Completion</span>
                    <span className="font-medium">{Math.round(lessonProgress)}%</span>
                  </div>
                  <Progress value={lessonProgress} className="h-2" />
                </div>

                <div className="grid grid-cols-2 gap-3 text-sm">
                  <div className="text-center p-2 bg-muted/50 rounded-lg">
                    <div className="font-medium">{completedObjectives.size}/{lessonData.objectives.length}</div>
                    <div className="text-xs text-muted-foreground">Objectives</div>
                  </div>
                  <div className="text-center p-2 bg-muted/50 rounded-lg">
                    <div className="font-medium">{Math.floor(timeSpent / 60)}m</div>
                    <div className="text-xs text-muted-foreground">Time</div>
                  </div>
                </div>
              </CardContent>
            </Card>

            {/* Objectives Checklist */}
            <Card className="glass-card">
              <CardHeader className="pb-3">
                <CardTitle className="text-lg flex items-center">
                  <Target className="mr-2 h-5 w-5 text-primary" />
                  Learning Objectives
                </CardTitle>
              </CardHeader>
              <CardContent>
                <div className="space-y-2">
                  {lessonData.objectives.map((objective, index) => (
                    <div
                      key={index}
                      className="flex items-start space-x-3 p-2 rounded-lg hover:bg-muted/50 cursor-pointer transition-colors group"
                      onClick={() => {
                        const newCompleted = new Set(completedObjectives);
                        if (completedObjectives.has(index)) {
                          newCompleted.delete(index);
                        } else {
                          newCompleted.add(index);
                        }
                        setCompletedObjectives(newCompleted);
                        setLessonProgress((newCompleted.size / lessonData.objectives.length) * 100);
                      }}
                    >
                      <div className={`p-1 rounded-full transition-colors ${completedObjectives.has(index) ? 'bg-green-500/20' : 'bg-muted group-hover:bg-primary/20'}`}>
                        <CheckCircle className={`h-4 w-4 transition-colors ${completedObjectives.has(index) ? 'text-green-500' : 'text-muted-foreground group-hover:text-primary'}`} />
                      </div>
                      <span className={`text-sm leading-relaxed transition-all ${completedObjectives.has(index) ? 'line-through text-muted-foreground' : 'group-hover:text-foreground'}`}>
                        {objective}
                      </span>
                    </div>
                  ))}
                </div>
              </CardContent>
            </Card>

            {/* Table of Contents */}
            <TableOfContents
              items={[
                {
                  id: "lesson",
                  title: "Theory & Concepts",
                  level: 1,
                  timeEstimate: "10 min"
                },
                {
                  id: "practice",
                  title: "Interactive Practice",
                  level: 1,
                  timeEstimate: "15 min"
                },
                {
                  id: "exercise",
                  title: "Exercises",
                  level: 1,
                  timeEstimate: "20 min"
                },
                {
                  id: "notes",
                  title: "Notes",
                  level: 1,
                  timeEstimate: "5 min"
                },
                {
                  id: "progress",
                  title: "Progress",
                  level: 1,
                  timeEstimate: "2 min"
                }
              ]}
              currentSection={activeTab}
              onSectionClick={(sectionId) => setActiveTab(sectionId)}
            />
          </div>
        </div>

        {/* Main Content */}
        <div className="flex-1 min-w-0">
          <Tabs value={activeTab} onValueChange={setActiveTab} className="space-y-6 animate-in fade-in duration-1000">
        <TabsList className="grid w-full grid-cols-5 lg:w-[650px] bg-card/50 backdrop-blur-sm border border-border/50 shadow-lg">
          <TabsTrigger value="lesson" className="data-[state=active]:bg-primary data-[state=active]:text-primary-foreground">
            <Lightbulb className="mr-2 h-4 w-4" />
            Theory
          </TabsTrigger>
          <TabsTrigger value="practice" className="data-[state=active]:bg-primary data-[state=active]:text-primary-foreground">
            <Code2 className="mr-2 h-4 w-4" />
            Practice
          </TabsTrigger>
          <TabsTrigger value="exercise" className="data-[state=active]:bg-primary data-[state=active]:text-primary-foreground">
            <Trophy className="mr-2 h-4 w-4" />
            Exercise
          </TabsTrigger>
          <TabsTrigger value="notes" className="data-[state=active]:bg-primary data-[state=active]:text-primary-foreground">
            <StickyNote className="mr-2 h-4 w-4" />
            Notes
          </TabsTrigger>
          <TabsTrigger value="progress" className="data-[state=active]:bg-primary data-[state=active]:text-primary-foreground">
            <BarChart3 className="mr-2 h-4 w-4" />
            Progress
          </TabsTrigger>
        </TabsList>

        {/* Theory Tab */}
        <TabsContent value="lesson" className="space-y-6">
          <Card className="glass-card border-2">
            <CardHeader className="bg-gradient-to-r from-primary/10 to-primary/5 border-b">
              <CardTitle className="flex items-center text-xl">
                <div className="p-2 rounded-lg bg-primary/10 mr-3">
                  <Lightbulb className="h-5 w-5 text-primary" />
                </div>
                Learning Objectives
              </CardTitle>
            </CardHeader>
            <CardContent className="pt-6">
              <ul className="space-y-3">
                {lessonData.objectives.map((objective) => (
                  <li key={objective} className="flex items-start space-x-3 group">
                    <div className="p-1 rounded-full bg-green-500/10 group-hover:bg-green-500/20 transition-colors">
                      <CheckCircle className="h-4 w-4 text-green-500 flex-shrink-0" />
                    </div>
                    <span className="text-sm leading-relaxed">{objective}</span>
                  </li>
                ))}
              </ul>
            </CardContent>
          </Card>

          <Card className="glass-card border-2">
            <CardHeader className="bg-gradient-to-r from-primary/10 to-primary/5 border-b">
              <CardTitle className="flex items-center text-xl">
                <div className="p-2 rounded-lg bg-primary/10 mr-3">
                  <FileText className="h-5 w-5 text-primary" />
                </div>
                Theory & Concepts
              </CardTitle>
            </CardHeader>
            <CardContent className="pt-6">
              <MarkdownRenderer
                content={lessonData.theory}
                enableCodeHighlight={true}
                enableInteractiveExamples={true}
              />
            </CardContent>
          </Card>
        </TabsContent>

        {/* Practice Tab */}
        <TabsContent value="practice" className="space-y-6">
          <InteractiveExample
            title="Interactive Code Example"
            description="Try modifying this Go code to see how the concepts work. Experiment with different values and see the results!"
            initialCode={lessonData.codeExample}
            language="go"
            inputs={[
              { name: "name", type: "string", defaultValue: "World", description: "Name to greet" },
              { name: "count", type: "number", defaultValue: 3, description: "Number of times to repeat" }
            ]}
            expectedOutput="Hello, World!\nHello, World!\nHello, World!"
            explanation="This example demonstrates basic Go syntax including variables, functions, and loops. Try changing the input values to see how the output changes."
            concepts={["Variables", "Functions", "Loops", "String formatting"]}
          />

          <div className="glass-card p-6 rounded-2xl border-2">
            <CodeEditor
              title="Advanced Practice"
              description="Practice with more complex examples and test your understanding!"
              initialCode={lessonData.codeExample}
              solution={lessonData.solution}
              language="go"
              onCodeChange={(code) => console.log("Code changed:", code)}
            />
          </div>
        </TabsContent>

        {/* Exercise Tab */}
        <TabsContent value="exercise" className="space-y-6">
          {lessonData.exercises.map((exercise) => (
            <div key={exercise.id} className="grid grid-cols-1 lg:grid-cols-2 gap-6">
              <CodeEditor
                title={exercise.title}
                description={exercise.description}
                initialCode={exercise.initialCode}
                solution={exercise.solution}
                language="go"
                onCodeChange={(code) => console.log("Exercise code:", code)}
              />
              
              <ExerciseSubmission
                exerciseId={exercise.id}
                title={exercise.title}
                description={exercise.description}
                requirements={exercise.requirements}
                code={exercise.initialCode}
                previousSubmissions={[]}
              />
            </div>
          ))}
        </TabsContent>

        {/* Notes Tab */}
        <TabsContent value="notes" className="space-y-6">
          <LessonNotes
            lessonId={lessonData.id}
            lessonTitle={lessonData.title}
            onNotesChange={(notes) => console.log("Notes updated:", notes)}
          />
        </TabsContent>

        {/* Progress Tab */}
        <TabsContent value="progress" className="space-y-6">
          <LessonProgress
            lessonId={lessonData.id}
            objectives={lessonData.objectives}
            totalExercises={lessonData.exercises.length}
            onProgressUpdate={(progress) => console.log("Progress updated:", progress)}
          />
        </TabsContent>
      </Tabs>

      {/* Related Lessons */}
      <RelatedLessons
        currentLessonId={lessonData.id}
        lessons={[
          {
            id: lessonData.id + 1,
            title: "Advanced Go Concepts",
            description: "Dive deeper into Go with advanced topics like interfaces, goroutines, and channels.",
            difficulty: "Intermediate",
            duration: "45 min",
            progress: 0,
            category: "Core Go",
            tags: ["interfaces", "goroutines", "channels"],
            nextInPath: true,
            rating: 4.8,
            enrolledCount: 1250
          },
          {
            id: lessonData.id - 1,
            title: "Go Fundamentals",
            description: "Learn the basics of Go programming language from scratch.",
            difficulty: "Beginner",
            duration: "30 min",
            progress: 100,
            completed: true,
            category: "Core Go",
            tags: ["basics", "syntax", "variables"],
            prerequisite: true,
            rating: 4.9,
            enrolledCount: 2100
          },
          {
            id: lessonData.id + 2,
            title: "Building Web APIs with Go",
            description: "Create robust web APIs using Go's standard library and popular frameworks.",
            difficulty: "Advanced",
            duration: "60 min",
            progress: 25,
            category: "Web Development",
            tags: ["web", "api", "http"],
            recommended: true,
            rating: 4.7,
            enrolledCount: 890
          }
        ]}
        title="Continue Your Learning Journey"
        maxItems={6}
        showCategories={true}
        className="mt-12"
      />

      {/* Navigation */}
      <div className="glass-card p-6 rounded-2xl mt-12 border-2 animate-in fade-in duration-1000">
        <div className="flex flex-col sm:flex-row items-center justify-between gap-4">
          <div className="w-full sm:w-auto">
            {lessonData.prevLessonId ? (
              <Link href={`/learn/${lessonData.prevLessonId}`} className="block w-full sm:w-auto">
                <Button variant="outline" size="lg" className="w-full sm:w-auto shadow-sm hover:shadow-md">
                  <ArrowLeft className="mr-2 h-4 w-4" />
                  Previous Lesson
                </Button>
              </Link>
            ) : (
              <div className="w-full sm:w-auto" />
            )}
          </div>

          <div className="flex flex-col sm:flex-row items-center gap-3 w-full sm:w-auto">
            <Link href="/curriculum" className="w-full sm:w-auto">
              <Button variant="outline" size="lg" className="w-full sm:w-auto shadow-sm hover:shadow-md">
                <List className="mr-2 h-4 w-4" />
                View All Lessons
              </Button>
            </Link>

            {lessonData.nextLessonId && (
              <Link href={`/learn/${lessonData.nextLessonId}`} className="w-full sm:w-auto">
                <Button size="lg" className="go-gradient text-white w-full sm:w-auto shadow-lg hover:shadow-2xl">
                  Next Lesson
                  <ArrowRight className="ml-2 h-4 w-4" />
                </Button>
              </Link>
            )}
          </div>
          </div>
        </div>
      </div>
        </div>
      </div>
    </div>
  );
}
