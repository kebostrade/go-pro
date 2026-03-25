"use client";

import { useState, useEffect } from "react";
import { useParams, useRouter } from "next/navigation";
import { api, LessonDetail, LessonExercise } from "@/lib/api";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Progress } from "@/components/ui/progress";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { Skeleton } from "@/components/ui/skeleton";
import { Alert, AlertDescription, AlertTitle } from "@/components/ui/alert";
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
  StickyNote,
  Eye,
  EyeOff,
  AlertCircle,
  RefreshCw,
  Play,
  Copy,
  Check
} from "lucide-react";
import Link from "next/link";

export default function LessonPage() {
  const params = useParams();
  const router = useRouter();
  const [activeTab, setActiveTab] = useState("lesson");
  const [lessonData, setLessonData] = useState<LessonDetail | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  // Enhanced state
  const [isBookmarked, setIsBookmarked] = useState(false);
  const [lessonProgress, setLessonProgress] = useState(0);
  const [showSidebar, setShowSidebar] = useState(false);
  const [completedObjectives, setCompletedObjectives] = useState<Set<number>>(new Set());
  const [timeSpent, setTimeSpent] = useState(0);
  const [startTime] = useState(Date.now());
  const [showSolution, setShowSolution] = useState(false);
  const [solutionCode, setSolutionCode] = useState("");
  const [showSolutionWarning, setShowSolutionWarning] = useState(true);
  const [isCompleting, setIsCompleting] = useState(false);
  const [completedAt, setCompletedAt] = useState<string | null>(null);
  const [copiedCode, setCopiedCode] = useState(false);

  const lessonId = parseInt((params?.id as string) || "1", 10);
  const isValidId = !isNaN(lessonId);

  // Load lesson data
  useEffect(() => {
    if (!isValidId) {
      setLoading(false);
      setError("Invalid lesson ID");
      return;
    }

    const loadLessonData = async () => {
      setLoading(true);
      setError(null);

      try {
        const lesson = await api.getLessonDetail(lessonId);
        setLessonData(lesson);
        setSolutionCode(lesson.solution);

        // Load saved progress from localStorage
        const savedBookmark = localStorage.getItem(`lesson-${lessonId}-bookmarked`);
        if (savedBookmark) {
          setIsBookmarked(JSON.parse(savedBookmark));
        }

        const savedCompletion = localStorage.getItem(`lesson-${lessonId}-completed`);
        if (savedCompletion) {
          setCompletedAt(savedCompletion);
        }
      } catch (err) {
        console.error("Failed to load lesson:", err);
        setError(err instanceof Error ? err.message : "Failed to load lesson");
      } finally {
        setLoading(false);
      }
    };

    loadLessonData();
  }, [lessonId, isValidId]);

  // Time tracking
  useEffect(() => {
    const interval = setInterval(() => {
      setTimeSpent(Math.floor((Date.now() - startTime) / 1000));
    }, 1000);

    return () => clearInterval(interval);
  }, [startTime]);

  // Calculate progress based on completed objectives
  useEffect(() => {
    if (lessonData) {
      const progress = (completedObjectives.size / lessonData.objectives.length) * 100;
      setLessonProgress(progress);
    }
  }, [completedObjectives, lessonData]);

  // Handle bookmark toggle
  const handleBookmarkToggle = () => {
    const newBookmarked = !isBookmarked;
    setIsBookmarked(newBookmarked);
    localStorage.setItem(`lesson-${lessonId}-bookmarked`, JSON.stringify(newBookmarked));
  };

  // Handle objective toggle
  const handleObjectiveToggle = (index: number) => {
    const newCompleted = new Set(completedObjectives);
    if (newCompleted.has(index)) {
      newCompleted.delete(index);
    } else {
      newCompleted.add(index);
    }
    setCompletedObjectives(newCompleted);
  };

  // Handle lesson completion
  const handleCompleteLesson = async () => {
    if (!lessonData) return;

    setIsCompleting(true);
    try {
      // In a real app, you'd call the API with auth token
      // await api.completeLesson(lessonId.toString(), authToken);

      const completionTime = new Date().toISOString();
      setCompletedAt(completionTime);
      localStorage.setItem(`lesson-${lessonId}-completed`, completionTime);

      // Mark all objectives as completed
      const allObjectives = new Set(lessonData.objectives.map((_, i) => i));
      setCompletedObjectives(allObjectives);

      // Show success message
      alert("Congratulations! You've completed this lesson! 🎉");
    } catch (err) {
      console.error("Failed to mark lesson as complete:", err);
      alert("Failed to mark lesson as complete. Please try again.");
    } finally {
      setIsCompleting(false);
    }
  };

  // Handle solution reveal
  const handleRevealSolution = () => {
    if (showSolutionWarning) {
      const confirmed = confirm(
        "⚠️ Warning: Viewing the solution before attempting the exercise yourself will reduce your learning effectiveness. Are you sure you want to continue?"
      );
      if (!confirmed) return;
      setShowSolutionWarning(false);
    }
    setShowSolution(true);
  };

  // Handle code copy
  const handleCopyCode = async (code: string) => {
    try {
      await navigator.clipboard.writeText(code);
      setCopiedCode(true);
      setTimeout(() => setCopiedCode(false), 2000);
    } catch (err) {
      console.error("Failed to copy code:", err);
    }
  };

  // Handle retry loading
  const handleRetry = () => {
    setError(null);
    setLoading(true);
    // Trigger reload by changing state
    router.refresh();
  };

  // Loading skeleton
  if (loading) {
    return (
      <div className="min-h-screen bg-gradient-to-br from-background via-background to-accent/5">
        <div className="container mx-auto px-4 py-8">
          {/* Header Skeleton */}
          <div className="mb-6">
            <Skeleton className="h-4 w-64 mb-4" />
            <Skeleton className="h-2 w-full mb-6" />
          </div>

          {/* Title Skeleton */}
          <Card className="mb-8 p-6">
            <Skeleton className="h-8 w-3/4 mb-4" />
            <Skeleton className="h-4 w-full mb-2" />
            <Skeleton className="h-4 w-5/6" />
            <div className="grid grid-cols-4 gap-4 mt-6">
              <Skeleton className="h-20" />
              <Skeleton className="h-20" />
              <Skeleton className="h-20" />
              <Skeleton className="h-20" />
            </div>
          </Card>

          {/* Content Skeleton */}
          <div className="flex gap-6">
            <div className="hidden lg:block w-80">
              <Skeleton className="h-96" />
            </div>
            <div className="flex-1">
              <Skeleton className="h-12 mb-6" />
              <Skeleton className="h-96" />
            </div>
          </div>
        </div>
      </div>
    );
  }

  // Error state
  if (error) {
    return (
      <div className="min-h-screen bg-gradient-to-br from-background via-background to-accent/5">
        <div className="container mx-auto px-4 py-8">
          <Alert variant="destructive" className="max-w-2xl mx-auto">
            <AlertCircle className="h-4 w-4" />
            <AlertTitle>Error Loading Lesson</AlertTitle>
            <AlertDescription className="mt-2">
              {error}
              <div className="mt-4 flex gap-2">
                <Button onClick={handleRetry} variant="outline" size="sm">
                  <RefreshCw className="mr-2 h-4 w-4" />
                  Retry
                </Button>
                <Link href="/learn">
                  <Button variant="outline" size="sm">
                    <ArrowLeft className="mr-2 h-4 w-4" />
                    Back to Lessons
                  </Button>
                </Link>
              </div>
            </AlertDescription>
          </Alert>
        </div>
      </div>
    );
  }

  // No lesson data
  if (!lessonData) {
    return (
      <div className="min-h-screen bg-gradient-to-br from-background via-background to-accent/5">
        <div className="container mx-auto px-4 py-8">
          <div className="text-center py-16">
            <h1 className="text-3xl font-bold mb-4">Lesson Not Found</h1>
            <p className="text-muted-foreground mb-6">
              The lesson you're looking for doesn't exist.
            </p>
            <Link href="/learn">
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

  const difficultyColor = {
    beginner: "success",
    intermediate: "warning",
    advanced: "info",
    expert: "destructive"
  } as const;

  const difficultyVariant = difficultyColor[lessonData.difficulty] || "outline";

  return (
    <div className="min-h-screen bg-gradient-to-br from-background via-background to-accent/5">
      <div className="container mx-auto px-4 py-8">
        {/* Breadcrumb Navigation */}
        <div className="mb-6">
          <div className="flex items-center justify-between mb-4">
            <div className="flex items-center space-x-2 text-sm text-muted-foreground">
              <Link href="/" className="hover:text-primary transition-colors">
                <Home className="h-4 w-4" />
              </Link>
              <ChevronRight className="h-4 w-4" />
              <Link href="/learn" className="hover:text-primary transition-colors">
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
                onClick={handleBookmarkToggle}
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
                  <span>
                    {completedObjectives.size}/{lessonData.objectives.length}
                  </span>
                </div>
              </div>
            </div>
            <Progress value={lessonProgress} className="h-2" />
          </div>
        </div>

        {/* Lesson Header */}
        <Card className="mb-8 p-6 lg:p-8 border-2 relative overflow-hidden">
          <div className="absolute top-0 right-0 w-96 h-96 bg-gradient-to-br from-primary/20 to-transparent rounded-full blur-3xl -z-10" />
          <div className="absolute bottom-0 left-0 w-64 h-64 bg-gradient-to-tr from-secondary/15 to-transparent rounded-full blur-2xl -z-10" />

          <div className="flex items-start justify-between mb-6">
            <div className="flex-1">
              {/* Badges */}
              <div className="flex items-center flex-wrap gap-2 mb-4">
                <Badge variant="outline">
                  <BookOpen className="mr-1 h-3 w-3" />
                  Lesson {lessonData.id}
                </Badge>
                <Badge variant="secondary">{lessonData.phase}</Badge>
                <Badge variant={difficultyVariant as any}>
                  {lessonData.difficulty.charAt(0).toUpperCase() + lessonData.difficulty.slice(1)}
                </Badge>
                {completedAt && (
                  <Badge variant="success">
                    <CheckCircle className="mr-1 h-3 w-3" />
                    Completed
                  </Badge>
                )}
                {isBookmarked && (
                  <Badge variant="outline" className="border-primary text-primary">
                    <Star className="mr-1 h-3 w-3 fill-current" />
                    Bookmarked
                  </Badge>
                )}
              </div>

              {/* Title */}
              <h1 className="text-3xl lg:text-4xl font-bold mb-4 bg-gradient-to-r from-primary to-blue-600 bg-clip-text text-transparent">
                {lessonData.title}
              </h1>

              {/* Description */}
              <p className="text-muted-foreground text-base lg:text-lg mb-6 leading-relaxed max-w-4xl">
                {lessonData.description}
              </p>

              {/* Stats Grid */}
              <div className="grid grid-cols-2 lg:grid-cols-4 gap-4">
                <div className="flex items-center space-x-2 bg-gradient-to-r from-primary/10 to-primary/5 px-4 py-3 rounded-xl border border-primary/20">
                  <Clock className="h-5 w-5 text-primary" />
                  <div>
                    <div className="text-sm font-medium">{lessonData.duration}</div>
                    <div className="text-xs text-muted-foreground">Duration</div>
                  </div>
                </div>

                <div className="flex items-center space-x-2 bg-gradient-to-r from-blue-500/10 to-blue-500/5 px-4 py-3 rounded-xl border border-blue-500/20">
                  <Target className="h-5 w-5 text-blue-500" />
                  <div>
                    <div className="text-sm font-medium">
                      {completedObjectives.size}/{lessonData.objectives.length}
                    </div>
                    <div className="text-xs text-muted-foreground">Objectives</div>
                  </div>
                </div>

                <div className="flex items-center space-x-2 bg-gradient-to-r from-green-500/10 to-green-500/5 px-4 py-3 rounded-xl border border-green-500/20">
                  <Award className="h-5 w-5 text-green-500" />
                  <div>
                    <div className="text-sm font-medium">{Math.round(lessonProgress)}%</div>
                    <div className="text-xs text-muted-foreground">Complete</div>
                  </div>
                </div>

                <div className="flex items-center space-x-2 bg-gradient-to-r from-orange-500/10 to-orange-500/5 px-4 py-3 rounded-xl border border-orange-500/20">
                  <Flame className="h-5 w-5 text-orange-500" />
                  <div>
                    <div className="text-sm font-medium">{Math.floor(timeSpent / 60)}m</div>
                    <div className="text-xs text-muted-foreground">Time Spent</div>
                  </div>
                </div>
              </div>
            </div>

            {/* Action Buttons */}
            <div className="ml-6 hidden lg:flex flex-col space-y-2">
              <Link href="/learn">
                <Button variant="outline" size="sm">
                  <List className="mr-2 h-4 w-4" />
                  All Lessons
                </Button>
              </Link>
              <Button variant="ghost" size="sm" onClick={handleBookmarkToggle}>
                {isBookmarked ? (
                  <BookmarkCheck className="mr-2 h-4 w-4 text-primary" />
                ) : (
                  <Bookmark className="mr-2 h-4 w-4" />
                )}
                {isBookmarked ? "Bookmarked" : "Bookmark"}
              </Button>
            </div>
          </div>
        </Card>

        {/* Mobile Sidebar */}
        {showSidebar && (
          <div className="fixed inset-0 z-50 lg:hidden">
            <div
              className="absolute inset-0 bg-background/80 backdrop-blur-sm"
              onClick={() => setShowSidebar(false)}
            />
            <div className="absolute right-0 top-0 h-full w-80 bg-card border-l border-border shadow-2xl overflow-y-auto">
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
                    <span className="font-medium">
                      {completedObjectives.size}/{lessonData.objectives.length}
                    </span>
                  </div>

                  <div className="flex items-center justify-between text-sm">
                    <span className="text-muted-foreground">Time Spent</span>
                    <span className="font-medium">{Math.floor(timeSpent / 60)}m</span>
                  </div>
                </div>

                {/* Objectives Checklist */}
                <div className="space-y-2">
                  <h4 className="text-sm font-medium text-muted-foreground">Objectives</h4>
                  <div className="space-y-2">
                    {lessonData.objectives.map((objective, index) => (
                      <div
                        key={index}
                        className="flex items-start space-x-2 p-2 rounded-lg hover:bg-muted/50 cursor-pointer transition-colors"
                        onClick={() => handleObjectiveToggle(index)}
                      >
                        <CheckCircle
                          className={`h-4 w-4 mt-0.5 ${
                            completedObjectives.has(index)
                              ? "text-green-500"
                              : "text-muted-foreground"
                          }`}
                        />
                        <span
                          className={`text-xs leading-relaxed ${
                            completedObjectives.has(index)
                              ? "line-through text-muted-foreground"
                              : ""
                          }`}
                        >
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
              <Card>
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
                      <div className="font-medium">
                        {completedObjectives.size}/{lessonData.objectives.length}
                      </div>
                      <div className="text-xs text-muted-foreground">Objectives</div>
                    </div>
                    <div className="text-center p-2 bg-muted/50 rounded-lg">
                      <div className="font-medium">{Math.floor(timeSpent / 60)}m</div>
                      <div className="text-xs text-muted-foreground">Time</div>
                    </div>
                  </div>

                  {/* Mark as Complete Button */}
                  <Button
                    onClick={handleCompleteLesson}
                    disabled={isCompleting || completedAt !== null}
                    className="w-full"
                    variant={completedAt ? "outline" : "default"}
                  >
                    {isCompleting ? (
                      <>
                        <RefreshCw className="mr-2 h-4 w-4 animate-spin" />
                        Completing...
                      </>
                    ) : completedAt ? (
                      <>
                        <CheckCircle className="mr-2 h-4 w-4" />
                        Completed
                      </>
                    ) : (
                      <>
                        <CheckCircle className="mr-2 h-4 w-4" />
                        Mark as Complete
                      </>
                    )}
                  </Button>

                  {completedAt && (
                    <p className="text-xs text-center text-muted-foreground">
                      Completed on {new Date(completedAt).toLocaleDateString()}
                    </p>
                  )}
                </CardContent>
              </Card>

              {/* Objectives Checklist */}
              <Card>
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
                        onClick={() => handleObjectiveToggle(index)}
                      >
                        <CheckCircle
                          className={`h-4 w-4 mt-0.5 transition-colors ${
                            completedObjectives.has(index)
                              ? "text-green-500"
                              : "text-muted-foreground group-hover:text-primary"
                          }`}
                        />
                        <span
                          className={`text-sm leading-relaxed transition-all ${
                            completedObjectives.has(index)
                              ? "line-through text-muted-foreground"
                              : "group-hover:text-foreground"
                          }`}
                        >
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
                    title: "Code Example",
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
                    id: "solution",
                    title: "Solution",
                    level: 1,
                    timeEstimate: "5 min"
                  }
                ]}
                currentSection={activeTab}
                onSectionClick={(sectionId) => setActiveTab(sectionId)}
              />
            </div>
          </div>

          {/* Main Content */}
          <div className="flex-1 min-w-0">
            <Tabs value={activeTab} onValueChange={setActiveTab} className="space-y-6">
              <TabsList className="grid w-full grid-cols-4 lg:w-[600px]">
                <TabsTrigger value="lesson">
                  <Lightbulb className="mr-2 h-4 w-4" />
                  Theory
                </TabsTrigger>
                <TabsTrigger value="practice">
                  <Code2 className="mr-2 h-4 w-4" />
                  Practice
                </TabsTrigger>
                <TabsTrigger value="exercise">
                  <Trophy className="mr-2 h-4 w-4" />
                  Exercise
                </TabsTrigger>
                <TabsTrigger value="solution">
                  <Eye className="mr-2 h-4 w-4" />
                  Solution
                </TabsTrigger>
              </TabsList>

              {/* Theory Tab */}
              <TabsContent value="lesson" className="space-y-6">
                {/* Objectives */}
                <Card className="border-2">
                  <CardHeader className="bg-gradient-to-r from-primary/10 to-primary/5 border-b">
                    <CardTitle className="flex items-center text-xl">
                      <div className="p-2 rounded-lg bg-primary/10 mr-3">
                        <Target className="h-5 w-5 text-primary" />
                      </div>
                      Learning Objectives
                    </CardTitle>
                  </CardHeader>
                  <CardContent className="pt-6">
                    <ul className="space-y-3">
                      {lessonData.objectives.map((objective, index) => (
                        <li
                          key={index}
                          className="flex items-start space-x-3 group cursor-pointer"
                          onClick={() => handleObjectiveToggle(index)}
                        >
                          <CheckCircle
                            className={`h-4 w-4 mt-0.5 flex-shrink-0 transition-colors ${
                              completedObjectives.has(index)
                                ? "text-green-500"
                                : "text-muted-foreground group-hover:text-primary"
                            }`}
                          />
                          <span
                            className={`text-sm leading-relaxed transition-all ${
                              completedObjectives.has(index)
                                ? "line-through text-muted-foreground"
                                : ""
                            }`}
                          >
                            {objective}
                          </span>
                        </li>
                      ))}
                    </ul>
                  </CardContent>
                </Card>

                {/* Theory Content */}
                <Card className="border-2">
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
                <Card className="border-2">
                  <CardHeader className="bg-gradient-to-r from-primary/10 to-primary/5 border-b">
                    <CardTitle className="flex items-center text-xl">
                      <div className="p-2 rounded-lg bg-primary/10 mr-3">
                        <Code2 className="h-5 w-5 text-primary" />
                      </div>
                      Interactive Code Example
                    </CardTitle>
                  </CardHeader>
                  <CardContent className="pt-6">
                    <CodeEditor
                      title="Practice Code"
                      description="Try modifying this code and run it to see the results!"
                      initialCode={lessonData.code_example}
                      solution={lessonData.solution}
                      language="go"
                      onCodeChange={(code) => console.log("Code changed:", code)}
                    />
                  </CardContent>
                </Card>
              </TabsContent>

              {/* Exercise Tab */}
              <TabsContent value="exercise" className="space-y-6">
                {lessonData.exercises.map((exercise: LessonExercise) => (
                  <Card key={exercise.id} className="border-2">
                    <CardHeader className="bg-gradient-to-r from-primary/10 to-primary/5 border-b">
                      <CardTitle className="flex items-center text-xl">
                        <div className="p-2 rounded-lg bg-primary/10 mr-3">
                          <Trophy className="h-5 w-5 text-primary" />
                        </div>
                        {exercise.title}
                      </CardTitle>
                    </CardHeader>
                    <CardContent className="pt-6">
                      <div className="space-y-6">
                        <div>
                          <h4 className="font-medium mb-2">Description</h4>
                          <p className="text-sm text-muted-foreground">{exercise.description}</p>
                        </div>

                        <div>
                          <h4 className="font-medium mb-2">Requirements</h4>
                          <ul className="space-y-2">
                            {exercise.requirements.map((req, idx) => (
                              <li key={idx} className="flex items-start space-x-2 text-sm">
                                <CheckCircle className="h-4 w-4 text-green-500 mt-0.5 flex-shrink-0" />
                                <span>{req}</span>
                              </li>
                            ))}
                          </ul>
                        </div>

                        <CodeEditor
                          title={`Exercise: ${exercise.title}`}
                          description="Complete the exercise based on the requirements above."
                          initialCode={exercise.initial_code}
                          solution={exercise.solution}
                          language="go"
                          onCodeChange={(code) => console.log("Exercise code:", code)}
                        />

                        <Link href={`/exercise/${exercise.id}`}>
                          <Button className="w-full">
                            <Play className="mr-2 h-4 w-4" />
                            Start Exercise
                          </Button>
                        </Link>
                      </div>
                    </CardContent>
                  </Card>
                ))}
              </TabsContent>

              {/* Solution Tab */}
              <TabsContent value="solution" className="space-y-6">
                <Card className="border-2">
                  <CardHeader className="bg-gradient-to-r from-yellow-500/10 to-yellow-500/5 border-b border-yellow-500/20">
                    <CardTitle className="flex items-center text-xl">
                      <div className="p-2 rounded-lg bg-yellow-500/10 mr-3">
                        <AlertCircle className="h-5 w-5 text-yellow-500" />
                      </div>
                      Solution Code
                    </CardTitle>
                  </CardHeader>
                  <CardContent className="pt-6 space-y-6">
                    {!showSolution ? (
                      <div className="text-center py-12">
                        <AlertCircle className="h-12 w-12 text-yellow-500 mx-auto mb-4" />
                        <h3 className="text-lg font-semibold mb-2">Solution Hidden</h3>
                        <p className="text-sm text-muted-foreground mb-6 max-w-md mx-auto">
                          Try solving the exercise yourself first. Viewing the solution too early can
                          reduce your learning effectiveness.
                        </p>
                        <Button onClick={handleRevealSolution} variant="outline">
                          <Eye className="mr-2 h-4 w-4" />
                          Reveal Solution
                        </Button>
                      </div>
                    ) : (
                      <div className="space-y-4">
                        <Alert variant="default" className="border-yellow-500/20 bg-yellow-500/5">
                          <Lightbulb className="h-4 w-4" />
                          <AlertTitle>Study the Solution</AlertTitle>
                          <AlertDescription>
                            Review this solution carefully and make sure you understand each part before
                            moving on.
                          </AlertDescription>
                        </Alert>

                        <Card>
                          <CardHeader>
                            <div className="flex items-center justify-between">
                              <CardTitle className="text-lg">Solution Code</CardTitle>
                              <Button
                                variant="outline"
                                size="sm"
                                onClick={() => handleCopyCode(solutionCode)}
                              >
                                {copiedCode ? (
                                  <Check className="h-4 w-4 text-green-500" />
                                ) : (
                                  <Copy className="h-4 w-4" />
                                )}
                              </Button>
                            </div>
                          </CardHeader>
                          <CardContent>
                            <pre className="p-4 rounded-lg bg-muted overflow-x-auto">
                              <code className="language-go text-sm">{solutionCode}</code>
                            </pre>
                          </CardContent>
                        </Card>

                        <Button
                          variant="ghost"
                          onClick={() => setShowSolution(false)}
                          className="w-full"
                        >
                          <EyeOff className="mr-2 h-4 w-4" />
                          Hide Solution
                        </Button>
                      </div>
                    )}
                  </CardContent>
                </Card>
              </TabsContent>
            </Tabs>

            {/* Navigation */}
            <Card className="p-6 mt-12 border-2">
              <div className="flex flex-col sm:flex-row items-center justify-between gap-4">
                <div className="w-full sm:w-auto">
                  {lessonData.prev_lesson_id ? (
                    <Link href={`/learn/${lessonData.prev_lesson_id}`} className="block w-full sm:w-auto">
                      <Button variant="outline" size="lg" className="w-full sm:w-auto">
                        <ArrowLeft className="mr-2 h-4 w-4" />
                        Previous Lesson
                      </Button>
                    </Link>
                  ) : (
                    <div className="w-full sm:w-auto" />
                  )}
                </div>

                <div className="flex flex-col sm:flex-row items-center gap-3 w-full sm:w-auto">
                  <Link href="/learn" className="w-full sm:w-auto">
                    <Button variant="outline" size="lg" className="w-full sm:w-auto">
                      <List className="mr-2 h-4 w-4" />
                      All Lessons
                    </Button>
                  </Link>

                  {lessonData.next_lesson_id && (
                    <Link href={`/learn/${lessonData.next_lesson_id}`} className="w-full sm:w-auto">
                      <Button size="lg" className="w-full sm:w-auto">
                        Next Lesson
                        <ArrowRight className="ml-2 h-4 w-4" />
                      </Button>
                    </Link>
                  )}
                </div>
              </div>
            </Card>
          </div>
        </div>
      </div>
    </div>
  );
}
