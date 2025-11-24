"use client";

import { useEffect, useState, useMemo, memo } from "react";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Progress } from "@/components/ui/progress";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import {
  BookOpen,
  Code2,
  Trophy,
  Play,
  ArrowRight,
  Clock,
  CheckCircle,
  Lock,
  Zap,
  Globe,
  TrendingUp,
  Award,
  Calendar,
  AlertCircle,
} from "lucide-react";
import Link from "next/link";
import { api, type Curriculum, type CurriculumPhase, type CurriculumLesson, type Project } from "@/lib/api";

// Icon mapping for phases
const iconMap: Record<string, any> = {
  zap: Zap,
  globe: Globe,
  trending: TrendingUp,
  award: Award,
  code: Code2,
};

// Difficulty badge variant mapping
const getDifficultyVariant = (difficulty: string) => {
  const lower = difficulty.toLowerCase();
  if (lower === 'beginner') return 'secondary';
  if (lower === 'intermediate') return 'default';
  if (lower === 'advanced') return 'destructive';
  return 'outline';
};

// Loading skeleton component - Memoized for better perceived performance
// Performance optimization: Prevents unnecessary re-renders during loading state
const LoadingSkeleton = memo(function LoadingSkeleton() {
  return (
    <div className="space-y-8 animate-in fade-in duration-700">
      {[1, 2, 3, 4].map(i => (
        <div key={i} className="animate-pulse">
          <div className="h-8 bg-gray-200 dark:bg-gray-700 rounded w-1/4 mb-4"></div>
          <div className="space-y-3">
            {[1, 2, 3].map(j => (
              <div key={j} className="h-24 bg-gray-100 dark:bg-gray-800 rounded"></div>
            ))}
          </div>
        </div>
      ))}
    </div>
  );
});

// Error message component
function ErrorMessage({ error, onRetry }: { error: string; onRetry: () => void }) {
  return (
    <div className="text-center py-12 animate-in fade-in slide-in-bottom duration-500">
      <div className="inline-flex items-center justify-center w-16 h-16 rounded-full bg-red-100 dark:bg-red-900/20 mb-4">
        <AlertCircle className="h-8 w-8 text-red-600 dark:text-red-400" />
      </div>
      <h3 className="text-lg font-semibold mb-2">Failed to Load Curriculum</h3>
      <p className="text-muted-foreground mb-4 max-w-md mx-auto">{error}</p>
      <Button onClick={onRetry} variant="outline">
        Try Again
      </Button>
    </div>
  );
}

// Empty state component - Memoized for performance
const EmptyState = memo(function EmptyState() {
  return (
    <div className="text-center py-12 animate-in fade-in duration-500">
      <div className="inline-flex items-center justify-center w-16 h-16 rounded-full bg-muted mb-4">
        <BookOpen className="h-8 w-8 text-muted-foreground" />
      </div>
      <h3 className="text-lg font-semibold mb-2">No Curriculum Available</h3>
      <p className="text-muted-foreground">The curriculum is being prepared. Check back soon!</p>
    </div>
  );
});

// Lesson card component - Memoized to prevent re-rendering unchanged lessons
// Performance optimization: Only re-renders when lesson data changes
const LessonCard = memo(function LessonCard({ lesson }: { lesson: CurriculumLesson }) {
  return (
    <Card className={`lesson-card ${lesson.locked ? 'opacity-60' : ''}`}>
      <CardHeader className="pb-3">
        <div className="flex items-start justify-between">
          <div className="flex-1">
            <div className="flex items-center space-x-2 mb-2">
              <Badge variant="outline" className="text-xs">
                Lesson {lesson.id}
              </Badge>
              <Badge variant={getDifficultyVariant(lesson.difficulty)} className="text-xs">
                {lesson.difficulty.charAt(0).toUpperCase() + lesson.difficulty.slice(1)}
              </Badge>
            </div>
            <CardTitle className="text-lg leading-tight">{lesson.title}</CardTitle>
            <CardDescription className="text-sm mt-1">
              {lesson.description}
            </CardDescription>
          </div>
          <div className="ml-3">
            {lesson.completed ? (
              <CheckCircle className="h-5 w-5 text-green-500" />
            ) : lesson.locked ? (
              <Lock className="h-5 w-5 text-muted-foreground" />
            ) : (
              <Play className="h-5 w-5 text-primary" />
            )}
          </div>
        </div>
      </CardHeader>
      <CardContent className="pt-0">
        <div className="flex items-center justify-between text-sm text-muted-foreground mb-3">
          <div className="flex items-center space-x-4">
            <div className="flex items-center space-x-1">
              <Clock className="h-3 w-3" />
              <span>{lesson.duration}</span>
            </div>
            <div className="flex items-center space-x-1">
              <Code2 className="h-3 w-3" />
              <span>{lesson.exercises} exercises</span>
            </div>
          </div>
        </div>
        <Link href={lesson.locked ? "#" : `/learn/${lesson.id}`}>
          <Button
            className="w-full"
            variant={lesson.completed ? "outline" : "default"}
            disabled={lesson.locked}
          >
            {lesson.completed ? (
              <>
                <CheckCircle className="mr-2 h-4 w-4" />
                Review Lesson
              </>
            ) : lesson.locked ? (
              <>
                <Lock className="mr-2 h-4 w-4" />
                Locked
              </>
            ) : (
              <>
                <Play className="mr-2 h-4 w-4" />
                Start Lesson
                <ArrowRight className="ml-2 h-4 w-4" />
              </>
            )}
          </Button>
        </Link>
      </CardContent>
    </Card>
  );
});

export default function CurriculumPage() {
  const [curriculum, setCurriculum] = useState<Curriculum | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [activePhase, setActivePhase] = useState<string>("");

  const fetchCurriculum = async () => {
    try {
      setLoading(true);
      setError(null);
      const data = await api.getCurriculum();
      setCurriculum(data);
      // Set first phase as active by default
      if (data.phases.length > 0) {
        setActivePhase(data.phases[0].id);
      }
    } catch (err) {
      setError(err instanceof Error ? err.message : 'An unexpected error occurred');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchCurriculum();
  }, []);

  if (loading) {
    return (
      <div className="min-h-screen animated-gradient">
        <div className="container-responsive padding-responsive-y">
          <LoadingSkeleton />
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="min-h-screen animated-gradient">
        <div className="container-responsive padding-responsive-y">
          <ErrorMessage error={error} onRetry={fetchCurriculum} />
        </div>
      </div>
    );
  }

  if (!curriculum || curriculum.phases.length === 0) {
    return (
      <div className="min-h-screen animated-gradient">
        <div className="container-responsive padding-responsive-y">
          <EmptyState />
        </div>
      </div>
    );
  }

  // Performance optimization: useMemo for expensive calculations
  // Prevents recalculation on every render unless curriculum changes
  const stats = useMemo(() => {
    if (!curriculum) return null;

    // Calculate overall progress
    const overallProgress = Math.round(
      curriculum.phases.reduce((acc, phase) => acc + phase.progress, 0) / curriculum.phases.length
    );

    // Count total lessons and exercises
    const totalLessons = curriculum.phases.reduce((acc, phase) => acc + phase.lessons.length, 0);
    const totalExercises = curriculum.phases.reduce(
      (acc, phase) => acc + phase.lessons.reduce((sum, lesson) => sum + lesson.exercises, 0),
      0
    );

    // Calculate total weeks
    const totalWeeks = curriculum.phases.reduce((acc, phase) => {
      const match = phase.weeks.match(/\d+/g);
      if (match) {
        const weeks = match.map(Number);
        return Math.max(acc, Math.max(...weeks));
      }
      return acc;
    }, 0);

    // Calculate total XP
    const totalXP = curriculum.projects.reduce((acc, project) => acc + project.points, 0);

    return { overallProgress, totalLessons, totalExercises, totalWeeks, totalXP };
  }, [curriculum]);

  // Extract for easier use
  const { overallProgress, totalLessons, totalExercises, totalWeeks, totalXP } = stats || {};

  return (
    <div className="min-h-screen animated-gradient">
      {/* Decorative elements */}
      <div className="absolute inset-0 overflow-hidden pointer-events-none">
        <div className="absolute -top-40 -right-40 w-96 h-96 bg-gradient-to-br from-primary/20 to-blue-500/20 rounded-full blur-3xl float-animation" />
        <div className="absolute -bottom-40 -left-40 w-96 h-96 bg-gradient-to-tr from-cyan-500/20 to-primary/20 rounded-full blur-3xl float-animation" style={{ animationDelay: '1s' }} />
        <div className="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 w-[600px] h-[600px] bg-gradient-to-r from-primary/10 via-transparent to-blue-500/10 rounded-full blur-3xl float-animation" style={{ animationDelay: '2s' }} />
      </div>

      <div className="container-responsive padding-responsive-y relative z-10">
        {/* Header */}
        <div className="margin-responsive">
          <div className="glass-card p-6 lg:p-8 rounded-2xl mb-8 border-2 animate-in fade-in slide-in-bottom duration-700">
            <div className="flex flex-col lg:flex-row lg:items-center lg:justify-between mb-6">
              <div className="mb-4 lg:mb-0">
                <Badge variant="outline" className="mb-3">
                  📚 Complete Curriculum
                </Badge>
                <h1 className="text-3xl lg:text-5xl font-bold tracking-tight mb-3 bg-gradient-to-r from-primary via-primary to-primary/70 bg-clip-text text-transparent">
                  {curriculum.title}
                </h1>
                <p className="text-base lg:text-xl text-muted-foreground max-w-2xl leading-relaxed">
                  {curriculum.description}
                </p>
              </div>
              <div className="text-left lg:text-right">
                <div className="inline-flex flex-col items-center lg:items-end p-4 rounded-xl bg-primary/10 border border-primary/20">
                  <div className="text-4xl lg:text-5xl font-bold text-primary mb-1">{overallProgress}%</div>
                  <div className="text-sm lg:text-base text-muted-foreground">Overall Progress</div>
                </div>
              </div>
            </div>

            <div className="grid grid-cols-2 md:grid-cols-4 gap-4 mb-6">
              <Card className="glass-card border-2 hover:border-primary/50 transition-all duration-300 group">
                <CardContent className="p-4 lg:p-6 text-center">
                  <div className="p-3 rounded-xl bg-blue-500/10 w-fit mx-auto mb-3 group-hover:scale-110 transition-transform">
                    <BookOpen className="h-6 w-6 lg:h-7 lg:w-7 text-blue-500" />
                  </div>
                  <div className="text-2xl lg:text-3xl font-bold mb-1">{totalLessons}</div>
                  <div className="text-xs lg:text-sm text-muted-foreground">Lessons</div>
                </CardContent>
              </Card>
              <Card className="glass-card border-2 hover:border-primary/50 transition-all duration-300 group">
                <CardContent className="p-4 lg:p-6 text-center">
                  <div className="p-3 rounded-xl bg-yellow-500/10 w-fit mx-auto mb-3 group-hover:scale-110 transition-transform">
                    <Trophy className="h-6 w-6 lg:h-7 lg:w-7 text-yellow-500" />
                  </div>
                  <div className="text-2xl lg:text-3xl font-bold mb-1">{curriculum.projects.length}</div>
                  <div className="text-xs lg:text-sm text-muted-foreground">Projects</div>
                </CardContent>
              </Card>
              <Card className="glass-card border-2 hover:border-primary/50 transition-all duration-300 group">
                <CardContent className="p-4 lg:p-6 text-center">
                  <div className="p-3 rounded-xl bg-green-500/10 w-fit mx-auto mb-3 group-hover:scale-110 transition-transform">
                    <Clock className="h-6 w-6 lg:h-7 lg:w-7 text-green-500" />
                  </div>
                  <div className="text-2xl lg:text-3xl font-bold mb-1">{totalWeeks}</div>
                  <div className="text-xs lg:text-sm text-muted-foreground">Weeks</div>
                </CardContent>
              </Card>
              <Card className="glass-card border-2 hover:border-primary/50 transition-all duration-300 group">
                <CardContent className="p-4 lg:p-6 text-center">
                  <div className="p-3 rounded-xl bg-purple-500/10 w-fit mx-auto mb-3 group-hover:scale-110 transition-transform">
                    <Award className="h-6 w-6 lg:h-7 lg:w-7 text-purple-500" />
                  </div>
                  <div className="text-2xl lg:text-3xl font-bold mb-1">{totalXP}</div>
                  <div className="text-xs lg:text-sm text-muted-foreground">Total XP</div>
                </CardContent>
              </Card>
            </div>

            <div className="space-y-2">
              <div className="flex items-center justify-between text-sm">
                <span className="text-muted-foreground">Progress</span>
                <span className="font-semibold text-primary">{overallProgress}% Complete</span>
              </div>
              <Progress value={overallProgress} className="h-3 shadow-sm" />
            </div>
          </div>
        </div>

        {/* Learning Path */}
        <Tabs value={activePhase} onValueChange={setActivePhase} className="space-y-6 animate-in fade-in duration-1000">
          <TabsList className="grid w-full grid-cols-2 lg:grid-cols-4 bg-card/50 backdrop-blur-sm border border-border/50 shadow-lg p-1">
            {curriculum.phases.map((phase) => {
              const PhaseIcon = iconMap[phase.icon.toLowerCase()] || Code2;
              return (
                <TabsTrigger
                  key={phase.id}
                  value={phase.id}
                  className="data-[state=active]:bg-primary data-[state=active]:text-primary-foreground"
                >
                  <PhaseIcon className="mr-2 h-4 w-4" />
                  <span className="hidden sm:inline">{phase.title}</span>
                  <span className="sm:hidden">{phase.title.slice(0, 4)}</span>
                </TabsTrigger>
              );
            })}
            <TabsTrigger
              value="projects"
              className="data-[state=active]:bg-primary data-[state=active]:text-primary-foreground"
            >
              <Award className="mr-2 h-4 w-4" />
              Projects
            </TabsTrigger>
          </TabsList>

          {/* Phase Content */}
          {curriculum.phases.map((phase) => {
            const PhaseIcon = iconMap[phase.icon.toLowerCase()] || Code2;
            return (
              <TabsContent key={phase.id} value={phase.id} className="space-y-6">
                <Card>
                  <CardHeader>
                    <div className="flex items-center justify-between">
                      <div className="flex items-center space-x-3">
                        <div className="p-2 rounded-lg bg-primary/10">
                          <PhaseIcon className={`h-6 w-6 ${phase.color}`} />
                        </div>
                        <div>
                          <CardTitle className="text-2xl">{phase.title}</CardTitle>
                          <CardDescription className="text-base">
                            {phase.description} • {phase.weeks}
                          </CardDescription>
                        </div>
                      </div>
                      <div className="text-right">
                        <div className="text-xl font-bold text-primary">{phase.progress}%</div>
                        <div className="text-sm text-muted-foreground">Complete</div>
                      </div>
                    </div>
                  </CardHeader>
                  <CardContent>
                    <Progress value={phase.progress} className="mb-4" />
                    <div className="grid grid-cols-1 lg:grid-cols-2 xl:grid-cols-3 gap-responsive">
                      {/* Performance optimization: Use memoized LessonCard component */}
                      {phase.lessons.map((lesson) => (
                        <LessonCard key={lesson.id} lesson={lesson} />
                      ))}
                    </div>
                  </CardContent>
                </Card>
              </TabsContent>
            );
          })}

          {/* Projects Tab */}
          <TabsContent value="projects" className="space-y-6">
            <Card>
              <CardHeader>
                <div className="flex items-center space-x-3">
                  <div className="p-2 rounded-lg bg-primary/10">
                    <Trophy className="h-6 w-6 text-yellow-500" />
                  </div>
                  <div>
                    <CardTitle className="text-2xl">Real Projects</CardTitle>
                    <CardDescription className="text-base">
                      Apply your skills to build production-ready applications
                    </CardDescription>
                  </div>
                </div>
              </CardHeader>
              <CardContent>
                <div className="grid grid-cols-1 lg:grid-cols-2 xl:grid-cols-3 gap-responsive">
                  {curriculum.projects.map((project) => (
                    <Card key={project.id} className={`lesson-card ${project.locked ? 'opacity-60' : ''}`}>
                      <CardHeader>
                        <div className="flex items-start justify-between">
                          <div className="flex-1">
                            <div className="flex items-center space-x-2 mb-2">
                              <Badge variant="outline" className="text-xs">
                                Project
                              </Badge>
                              <Badge variant={getDifficultyVariant(project.difficulty)} className="text-xs">
                                {project.difficulty.charAt(0).toUpperCase() + project.difficulty.slice(1)}
                              </Badge>
                              <Badge variant="secondary" className="text-xs">
                                {project.points} XP
                              </Badge>
                            </div>
                            <CardTitle className="text-lg">{project.title}</CardTitle>
                            <CardDescription className="text-sm mt-1">
                              {project.description}
                            </CardDescription>
                          </div>
                          <div className="ml-3">
                            {project.completed ? (
                              <CheckCircle className="h-5 w-5 text-green-500" />
                            ) : project.locked ? (
                              <Lock className="h-5 w-5 text-muted-foreground" />
                            ) : (
                              <Play className="h-5 w-5 text-primary" />
                            )}
                          </div>
                        </div>
                      </CardHeader>
                      <CardContent>
                        <div className="space-y-3">
                          <div className="flex items-center space-x-1 text-sm text-muted-foreground">
                            <Calendar className="h-3 w-3" />
                            <span>{project.duration}</span>
                          </div>
                          <div className="flex flex-wrap gap-1">
                            {project.skills.map((skill, index) => (
                              <Badge key={index} variant="outline" className="text-xs">
                                {skill}
                              </Badge>
                            ))}
                          </div>
                          <Link href={project.locked ? "#" : `/projects/${project.id}`}>
                            <Button
                              className="w-full"
                              variant={project.completed ? "outline" : "default"}
                              disabled={project.locked}
                            >
                              {project.completed ? (
                                <>
                                  <CheckCircle className="mr-2 h-4 w-4" />
                                  View Project
                                </>
                              ) : project.locked ? (
                                <>
                                  <Lock className="mr-2 h-4 w-4" />
                                  Complete Prerequisites
                                </>
                              ) : (
                                <>
                                  <Code2 className="mr-2 h-4 w-4" />
                                  Start Project
                                  <ArrowRight className="ml-2 h-4 w-4" />
                                </>
                              )}
                            </Button>
                          </Link>
                        </div>
                      </CardContent>
                    </Card>
                  ))}
                </div>
              </CardContent>
            </Card>
          </TabsContent>
        </Tabs>

        {/* CTA Section */}
        <Card className="mt-8 bg-gradient-to-r from-primary/5 to-primary/10 border-primary/20">
          <CardContent className="p-8 text-center">
            <h3 className="text-2xl font-bold mb-4">Ready to Start Your Go Journey?</h3>
            <p className="text-muted-foreground mb-6 max-w-2xl mx-auto">
              Begin with the foundations and work your way up to building production-ready microservices.
              Each lesson builds upon the previous one, ensuring you develop a solid understanding of Go.
            </p>
            <div className="flex flex-col sm:flex-row gap-4 justify-center">
              {curriculum.phases.length > 0 && curriculum.phases[0].lessons.length > 0 && (
                <Link href={`/learn/${curriculum.phases[0].lessons[0].id}`}>
                  <Button size="lg" className="go-gradient text-white">
                    <Play className="mr-2 h-5 w-5" />
                    Start First Lesson
                    <ArrowRight className="ml-2 h-5 w-5" />
                  </Button>
                </Link>
              )}
              <Link href="/learn">
                <Button size="lg" variant="outline">
                  <BookOpen className="mr-2 h-5 w-5" />
                  Continue Learning
                </Button>
              </Link>
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  );
}
