"use client";

import { useState } from "react";
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
  Calendar
} from "lucide-react";
import Link from "next/link";

interface Lesson {
  id: number;
  title: string;
  description: string;
  duration: string;
  exercises: number;
  difficulty: "Beginner" | "Intermediate" | "Advanced" | "Expert";
  completed: boolean;
  locked: boolean;
}

interface Phase {
  id: string;
  title: string;
  description: string;
  weeks: string;
  icon: any;
  color: string;
  lessons: Lesson[];
  progress: number;
}

interface Project {
  id: string;
  title: string;
  description: string;
  duration: string;
  difficulty: "Intermediate" | "Advanced" | "Expert";
  skills: string[];
  points: number;
  completed: boolean;
  locked: boolean;
}

export default function CurriculumPage() {
  const [activePhase, setActivePhase] = useState("foundations");

  // Curriculum data based on syllabus.md
  const phases: Phase[] = [
    {
      id: "foundations",
      title: "Foundations",
      description: "Master Go basics and core concepts",
      weeks: "Weeks 1-2",
      icon: Zap,
      color: "text-blue-500",
      progress: 80,
      lessons: [
        {
          id: 1,
          title: "Go Syntax and Basic Types",
          description: "Go installation, basic syntax, primitive types, constants and iota",
          duration: "3-4 hours",
          exercises: 5,
          difficulty: "Beginner",
          completed: true,
          locked: false,
        },
        {
          id: 2,
          title: "Variables, Constants, and Functions",
          description: "Variable declarations, scope, function definitions, multiple return values",
          duration: "4-5 hours",
          exercises: 6,
          difficulty: "Beginner",
          completed: true,
          locked: false,
        },
        {
          id: 3,
          title: "Control Structures and Loops",
          description: "if/else, switch statements, for loops, defer statements",
          duration: "3-4 hours",
          exercises: 7,
          difficulty: "Beginner",
          completed: true,
          locked: false,
        },
        {
          id: 4,
          title: "Arrays, Slices, and Maps",
          description: "Data structures, manipulation, memory considerations",
          duration: "5-6 hours",
          exercises: 8,
          difficulty: "Beginner",
          completed: true,
          locked: false,
        },
        {
          id: 5,
          title: "Pointers and Memory Management",
          description: "Pointer basics, memory allocation, garbage collection",
          duration: "4-5 hours",
          exercises: 8,
          difficulty: "Beginner",
          completed: true,
          locked: false,
        },
      ],
    },
    {
      id: "intermediate",
      title: "Intermediate",
      description: "Object-oriented concepts and concurrency",
      weeks: "Weeks 3-5",
      icon: Globe,
      color: "text-green-500",
      progress: 60,
      lessons: [
        {
          id: 6,
          title: "Structs and Methods",
          description: "Struct definition, methods, receivers, method sets",
          duration: "5-6 hours",
          exercises: 8,
          difficulty: "Intermediate",
          completed: true,
          locked: false,
        },
        {
          id: 7,
          title: "Interfaces and Polymorphism",
          description: "Interface definition, type assertions, composition",
          duration: "6-7 hours",
          exercises: 9,
          difficulty: "Intermediate",
          completed: false,
          locked: false,
        },
        {
          id: 8,
          title: "Error Handling Patterns",
          description: "Custom errors, wrapping, panic/recover, best practices",
          duration: "4-5 hours",
          exercises: 7,
          difficulty: "Intermediate",
          completed: false,
          locked: false,
        },
        {
          id: 9,
          title: "Goroutines and Channels",
          description: "Concurrency, channel operations, select statements",
          duration: "7-8 hours",
          exercises: 10,
          difficulty: "Intermediate",
          completed: false,
          locked: false,
        },
        {
          id: 10,
          title: "Packages and Modules",
          description: "Package organization, Go modules, dependency management",
          duration: "5-6 hours",
          exercises: 6,
          difficulty: "Intermediate",
          completed: false,
          locked: false,
        },
      ],
    },
    {
      id: "advanced",
      title: "Advanced",
      description: "Production-ready development skills",
      weeks: "Weeks 6-8",
      icon: TrendingUp,
      color: "text-purple-500",
      progress: 40,
      lessons: [
        {
          id: 11,
          title: "Advanced Concurrency Patterns",
          description: "Worker pools, pipelines, context package, sync primitives",
          duration: "8-9 hours",
          exercises: 12,
          difficulty: "Advanced",
          completed: true,
          locked: false,
        },
        {
          id: 12,
          title: "Testing and Benchmarking",
          description: "Unit testing, table-driven tests, benchmarking, profiling",
          duration: "6-7 hours",
          exercises: 8,
          difficulty: "Advanced",
          completed: false,
          locked: false,
        },
        {
          id: 13,
          title: "HTTP Servers and REST APIs",
          description: "HTTP server basics, routing, middleware, authentication",
          duration: "8-9 hours",
          exercises: 10,
          difficulty: "Advanced",
          completed: true,
          locked: false,
        },
        {
          id: 14,
          title: "Database Integration",
          description: "Database/sql package, connection pooling, transactions",
          duration: "7-8 hours",
          exercises: 9,
          difficulty: "Advanced",
          completed: true,
          locked: false,
        },
        {
          id: 15,
          title: "Microservices Architecture",
          description: "Design principles, service communication, monitoring",
          duration: "9-10 hours",
          exercises: 11,
          difficulty: "Advanced",
          completed: false,
          locked: false,
        },
      ],
    },
    {
      id: "expert",
      title: "Expert",
      description: "Advanced patterns and production systems",
      weeks: "Weeks 9-12",
      icon: Award,
      color: "text-orange-500",
      progress: 0,
      lessons: [
        {
          id: 16,
          title: "Performance Optimization and Profiling",
          description: "Memory optimization, CPU profiling, benchmarking techniques",
          duration: "8-10 hours",
          exercises: 10,
          difficulty: "Expert",
          completed: false,
          locked: false,
        },
        {
          id: 17,
          title: "Security Best Practices",
          description: "Authentication, encryption, vulnerability prevention",
          duration: "7-9 hours",
          exercises: 9,
          difficulty: "Expert",
          completed: false,
          locked: false,
        },
        {
          id: 18,
          title: "Deployment and DevOps",
          description: "Docker, CI/CD, cloud deployment, monitoring",
          duration: "9-11 hours",
          exercises: 12,
          difficulty: "Expert",
          completed: false,
          locked: false,
        },
        {
          id: 19,
          title: "Advanced Design Patterns",
          description: "Functional programming, generics, architectural patterns",
          duration: "8-10 hours",
          exercises: 11,
          difficulty: "Expert",
          completed: false,
          locked: false,
        },
        {
          id: 20,
          title: "Building Production Systems",
          description: "Complete system design, observability, scalability",
          duration: "10-12 hours",
          exercises: 15,
          difficulty: "Expert",
          completed: false,
          locked: false,
        },
      ],
    },
  ];

  const projects: Project[] = [
    {
      id: "cli-task-manager",
      title: "CLI Task Manager",
      description: "Command-line application with file persistence",
      duration: "1 week",
      difficulty: "Intermediate",
      skills: ["File I/O", "JSON handling", "CLI design"],
      points: 100,
      completed: false,
      locked: false,
    },
    {
      id: "rest-api-server",
      title: "REST API with Database",
      description: "Full REST API with PostgreSQL integration",
      duration: "1.5 weeks",
      difficulty: "Advanced",
      skills: ["HTTP servers", "Database integration", "Testing"],
      points: 150,
      completed: false,
      locked: false,
    },
    {
      id: "realtime-chat",
      title: "Real-time Chat Server",
      description: "WebSocket-based chat with concurrent users",
      duration: "1.5 weeks",
      difficulty: "Advanced",
      skills: ["Concurrency", "WebSockets", "Real-time communication"],
      points: 200,
      completed: false,
      locked: false,
    },
    {
      id: "microservices-system",
      title: "Microservices System",
      description: "Multi-service system with API gateway",
      duration: "2 weeks",
      difficulty: "Advanced",
      skills: ["Microservices", "Service mesh", "Monitoring"],
      points: 250,
      completed: false,
      locked: false,
    },
    {
      id: "distributed-cache",
      title: "Distributed Cache System",
      description: "High-performance distributed cache similar to Redis",
      duration: "3 weeks",
      difficulty: "Expert",
      skills: ["Distributed Systems", "Networking", "Performance"],
      points: 300,
      completed: false,
      locked: false,
    },
    {
      id: "event-driven-system",
      title: "Event-Driven Architecture",
      description: "Complete event-driven system with CQRS and event sourcing",
      duration: "3 weeks",
      difficulty: "Expert",
      skills: ["Event Sourcing", "CQRS", "Message Queues"],
      points: 350,
      completed: false,
      locked: false,
    },
    {
      id: "observability-platform",
      title: "Monitoring & Observability Platform",
      description: "Comprehensive monitoring platform for distributed systems",
      duration: "3 weeks",
      difficulty: "Expert",
      skills: ["Observability", "Time Series", "Real-time Analytics"],
      points: 400,
      completed: false,
      locked: false,
    },
  ];

  const currentPhase = phases.find(p => p.id === activePhase) || phases[0];
  const overallProgress = Math.round(phases.reduce((acc, phase) => acc + phase.progress, 0) / phases.length);

  return (
    <div className="min-h-screen animated-gradient">
      <div className="container-responsive padding-responsive-y">
        {/* Header */}
        <div className="margin-responsive">
          <div className="glass-card p-6 lg:p-8 rounded-2xl mb-8 border-2 animate-in fade-in slide-in-bottom duration-700">
            <div className="flex flex-col lg:flex-row lg:items-center lg:justify-between mb-6">
              <div className="mb-4 lg:mb-0">
                <Badge variant="outline" className="mb-3">
                  📚 Complete Curriculum
                </Badge>
                <h1 className="text-3xl lg:text-5xl font-bold tracking-tight mb-3 bg-gradient-to-r from-primary via-primary to-primary/70 bg-clip-text text-transparent">
                  GO-PRO Curriculum
                </h1>
                <p className="text-base lg:text-xl text-muted-foreground max-w-2xl leading-relaxed">
                  Complete Go Programming Mastery - From Basics to Microservices
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
                  <div className="text-2xl lg:text-3xl font-bold mb-1">20</div>
                  <div className="text-xs lg:text-sm text-muted-foreground">Lessons</div>
                </CardContent>
              </Card>
              <Card className="glass-card border-2 hover:border-primary/50 transition-all duration-300 group">
                <CardContent className="p-4 lg:p-6 text-center">
                  <div className="p-3 rounded-xl bg-yellow-500/10 w-fit mx-auto mb-3 group-hover:scale-110 transition-transform">
                    <Trophy className="h-6 w-6 lg:h-7 lg:w-7 text-yellow-500" />
                  </div>
                  <div className="text-2xl lg:text-3xl font-bold mb-1">7</div>
                  <div className="text-xs lg:text-sm text-muted-foreground">Projects</div>
                </CardContent>
              </Card>
              <Card className="glass-card border-2 hover:border-primary/50 transition-all duration-300 group">
                <CardContent className="p-4 lg:p-6 text-center">
                  <div className="p-3 rounded-xl bg-green-500/10 w-fit mx-auto mb-3 group-hover:scale-110 transition-transform">
                    <Clock className="h-6 w-6 lg:h-7 lg:w-7 text-green-500" />
                  </div>
                  <div className="text-2xl lg:text-3xl font-bold mb-1">16</div>
                  <div className="text-xs lg:text-sm text-muted-foreground">Weeks</div>
                </CardContent>
              </Card>
              <Card className="glass-card border-2 hover:border-primary/50 transition-all duration-300 group">
                <CardContent className="p-4 lg:p-6 text-center">
                  <div className="p-3 rounded-xl bg-purple-500/10 w-fit mx-auto mb-3 group-hover:scale-110 transition-transform">
                    <Award className="h-6 w-6 lg:h-7 lg:w-7 text-purple-500" />
                  </div>
                  <div className="text-2xl lg:text-3xl font-bold mb-1">700</div>
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
          <TabsTrigger value="foundations" className="data-[state=active]:bg-primary data-[state=active]:text-primary-foreground">
            <Zap className="mr-2 h-4 w-4" />
            <span className="hidden sm:inline">Foundations</span>
            <span className="sm:hidden">Basic</span>
          </TabsTrigger>
          <TabsTrigger value="intermediate" className="data-[state=active]:bg-primary data-[state=active]:text-primary-foreground">
            <Code2 className="mr-2 h-4 w-4" />
            <span className="hidden sm:inline">Intermediate</span>
            <span className="sm:hidden">Inter</span>
          </TabsTrigger>
          <TabsTrigger value="advanced" className="data-[state=active]:bg-primary data-[state=active]:text-primary-foreground">
            <Trophy className="mr-2 h-4 w-4" />
            <span className="hidden sm:inline">Advanced</span>
            <span className="sm:hidden">Adv</span>
          </TabsTrigger>
          <TabsTrigger value="projects" className="data-[state=active]:bg-primary data-[state=active]:text-primary-foreground">
            <Award className="mr-2 h-4 w-4" />
            Projects
          </TabsTrigger>
        </TabsList>

        {/* Phase Content */}
        {phases.map((phase) => (
          <TabsContent key={phase.id} value={phase.id} className="space-y-6">
            <Card>
              <CardHeader>
                <div className="flex items-center justify-between">
                  <div className="flex items-center space-x-3">
                    <div className={`p-2 rounded-lg bg-primary/10`}>
                      <phase.icon className={`h-6 w-6 ${phase.color}`} />
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
                  {phase.lessons.map((lesson) => (
                    <Card key={lesson.id} className={`lesson-card ${lesson.locked ? 'opacity-60' : ''}`}>
                      <CardHeader className="pb-3">
                        <div className="flex items-start justify-between">
                          <div className="flex-1">
                            <div className="flex items-center space-x-2 mb-2">
                              <Badge variant="outline" className="text-xs">
                                Lesson {lesson.id}
                              </Badge>
                              <Badge
                                variant={lesson.difficulty === 'Beginner' ? 'secondary' :
                                        lesson.difficulty === 'Intermediate' ? 'default' :
                                        lesson.difficulty === 'Advanced' ? 'destructive' : 'outline'}
                                className="text-xs"
                              >
                                {lesson.difficulty}
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
                  ))}
                </div>
              </CardContent>
            </Card>
          </TabsContent>
        ))}

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
                    Apply your skills to build production-ready applications • Weeks 9-12
                  </CardDescription>
                </div>
              </div>
            </CardHeader>
            <CardContent>
              <div className="grid grid-cols-1 lg:grid-cols-2 xl:grid-cols-3 gap-responsive">
                {projects.map((project) => (
                  <Card key={project.id} className={`lesson-card ${project.locked ? 'opacity-60' : ''}`}>
                    <CardHeader>
                      <div className="flex items-start justify-between">
                        <div className="flex-1">
                          <div className="flex items-center space-x-2 mb-2">
                            <Badge variant="outline" className="text-xs">
                              Project
                            </Badge>
                            <Badge
                              variant={project.difficulty === 'Intermediate' ? 'default' :
                                      project.difficulty === 'Advanced' ? 'destructive' : 'outline'}
                              className="text-xs"
                            >
                              {project.difficulty}
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
            <Link href="/learn/1">
              <Button size="lg" className="go-gradient text-white">
                <Play className="mr-2 h-5 w-5" />
                Start First Lesson
                <ArrowRight className="ml-2 h-5 w-5" />
              </Button>
            </Link>
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
