"use client";

import { useState, useEffect } from "react";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Progress } from "@/components/ui/progress";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import {
  Code2,
  Trophy,
  Clock,
  Star,
  Users,
  GitBranch,
  Play,
  CheckCircle,
  Lock,
  ArrowRight,
  Terminal,
  Server,
  MessageSquare,
  Layers,
  Target,
  BookOpen,
  Zap,
  Award,
  TrendingUp,
  Calendar,
  Database,
  Globe,
  Cpu
} from "lucide-react";
import Link from "next/link";

interface Project {
  id: string;
  title: string;
  description: string;
  longDescription: string;
  difficulty: "Beginner" | "Intermediate" | "Advanced";
  estimatedTime: string;
  technologies: string[];
  prerequisites: string[];
  learningOutcomes: string[];
  chapters: Chapter[];
  completed: boolean;
  locked: boolean;
  progress: number;
  icon: any;
  category: string;
  githubRepo?: string;
  liveDemo?: string;
}

interface Chapter {
  id: string;
  title: string;
  description: string;
  estimatedTime: string;
  completed: boolean;
  locked: boolean;
}

export default function ProjectsPage() {
  const [activeTab, setActiveTab] = useState("all");
  const [selectedDifficulty, setSelectedDifficulty] = useState<string>("all");

  // Mock data for projects
  const projects: Project[] = [
    {
      id: "cli-task-manager",
      title: "CLI Task Manager",
      description: "Build a command-line task management application with file persistence",
      longDescription: "Create a full-featured command-line task manager that demonstrates Go's CLI capabilities, file I/O, JSON handling, and clean architecture patterns. This project covers fundamental Go concepts while building something practical.",
      difficulty: "Beginner",
      estimatedTime: "8-10 hours",
      technologies: ["Go", "CLI", "JSON", "File I/O", "Cobra"],
      prerequisites: ["Go Basics", "Functions", "Structs", "Error Handling"],
      learningOutcomes: [
        "Build command-line applications with Cobra",
        "Handle file I/O operations safely",
        "Work with JSON serialization/deserialization",
        "Implement clean architecture patterns",
        "Create comprehensive CLI interfaces",
        "Handle user input validation"
      ],
      chapters: [
        {
          id: "setup",
          title: "Project Setup & Architecture",
          description: "Set up the project structure and define the core architecture",
          estimatedTime: "1 hour",
          completed: true,
          locked: false
        },
        {
          id: "models",
          title: "Data Models & Storage",
          description: "Create task models and implement JSON file storage",
          estimatedTime: "2 hours",
          completed: true,
          locked: false
        },
        {
          id: "commands",
          title: "CLI Commands Implementation",
          description: "Build add, list, complete, and delete commands",
          estimatedTime: "3 hours",
          completed: false,
          locked: false
        },
        {
          id: "advanced",
          title: "Advanced Features",
          description: "Add filtering, sorting, and due date functionality",
          estimatedTime: "2 hours",
          completed: false,
          locked: true
        },
        {
          id: "testing",
          title: "Testing & Documentation",
          description: "Write comprehensive tests and documentation",
          estimatedTime: "2 hours",
          completed: false,
          locked: true
        }
      ],
      completed: false,
      locked: false,
      progress: 40,
      icon: Terminal,
      category: "CLI Applications",
      githubRepo: "https://github.com/go-pro/cli-task-manager",
    },
    {
      id: "rest-api-server",
      title: "REST API Server",
      description: "Build a scalable REST API with database integration and authentication",
      longDescription: "Develop a production-ready REST API server that demonstrates Go's web capabilities, database integration, authentication, middleware, and API design best practices. Perfect for backend development skills.",
      difficulty: "Intermediate",
      estimatedTime: "12-15 hours",
      technologies: ["Go", "HTTP", "PostgreSQL", "JWT", "Docker", "Gin"],
      prerequisites: ["HTTP Basics", "Database Concepts", "JSON", "Authentication"],
      learningOutcomes: [
        "Build RESTful APIs with proper HTTP methods",
        "Integrate with PostgreSQL database",
        "Implement JWT authentication",
        "Create middleware for logging and auth",
        "Handle database migrations",
        "Write API documentation",
        "Deploy with Docker"
      ],
      chapters: [
        {
          id: "setup",
          title: "API Setup & Routing",
          description: "Set up the HTTP server and define API routes",
          estimatedTime: "2 hours",
          completed: false,
          locked: false
        },
        {
          id: "database",
          title: "Database Integration",
          description: "Connect to PostgreSQL and implement data models",
          estimatedTime: "3 hours",
          completed: false,
          locked: false
        },
        {
          id: "auth",
          title: "Authentication System",
          description: "Implement JWT-based authentication",
          estimatedTime: "3 hours",
          completed: false,
          locked: true
        },
        {
          id: "endpoints",
          title: "CRUD Endpoints",
          description: "Build complete CRUD operations for resources",
          estimatedTime: "3 hours",
          completed: false,
          locked: true
        },
        {
          id: "deployment",
          title: "Testing & Deployment",
          description: "Write tests and deploy with Docker",
          estimatedTime: "3 hours",
          completed: false,
          locked: true
        }
      ],
      completed: false,
      locked: false,
      progress: 0,
      icon: Server,
      category: "Web APIs",
      githubRepo: "https://github.com/go-pro/rest-api-server",
      liveDemo: "https://api-demo.go-pro.dev"
    },
    {
      id: "realtime-chat",
      title: "Real-time Chat Server",
      description: "Create a WebSocket-based chat application with rooms and user management",
      longDescription: "Build a real-time chat application using WebSockets, demonstrating Go's concurrency features, real-time communication, and scalable architecture patterns for handling multiple concurrent connections.",
      difficulty: "Intermediate",
      estimatedTime: "10-12 hours",
      technologies: ["Go", "WebSockets", "Redis", "Goroutines", "Channels"],
      prerequisites: ["Concurrency", "HTTP", "JSON", "Basic Networking"],
      learningOutcomes: [
        "Implement WebSocket communication",
        "Handle concurrent connections with goroutines",
        "Use channels for message passing",
        "Integrate Redis for session management",
        "Build real-time features",
        "Handle connection lifecycle",
        "Implement chat rooms and user management"
      ],
      chapters: [
        {
          id: "websockets",
          title: "WebSocket Foundation",
          description: "Set up WebSocket server and handle connections",
          estimatedTime: "2 hours",
          completed: false,
          locked: false
        },
        {
          id: "concurrency",
          title: "Concurrent Message Handling",
          description: "Implement goroutines and channels for message routing",
          estimatedTime: "3 hours",
          completed: false,
          locked: true
        },
        {
          id: "rooms",
          title: "Chat Rooms & User Management",
          description: "Create chat rooms and user authentication",
          estimatedTime: "3 hours",
          completed: false,
          locked: true
        },
        {
          id: "features",
          title: "Advanced Chat Features",
          description: "Add private messages, file sharing, and notifications",
          estimatedTime: "3 hours",
          completed: false,
          locked: true
        },
        {
          id: "scaling",
          title: "Scaling & Performance",
          description: "Optimize for high concurrent connections",
          estimatedTime: "2 hours",
          completed: false,
          locked: true
        }
      ],
      completed: false,
      locked: true,
      progress: 0,
      icon: MessageSquare,
      category: "Real-time Applications",
      githubRepo: "https://github.com/go-pro/realtime-chat",
      liveDemo: "https://chat-demo.go-pro.dev"
    },
    {
      id: "microservices-system",
      title: "Microservices System",
      description: "Design and implement a complete microservices architecture with service discovery",
      longDescription: "Build a comprehensive microservices system that demonstrates advanced Go concepts including service discovery, API gateways, distributed tracing, and inter-service communication patterns.",
      difficulty: "Advanced",
      estimatedTime: "20-25 hours",
      technologies: ["Go", "gRPC", "Docker", "Kubernetes", "Consul", "Prometheus"],
      prerequisites: ["HTTP APIs", "Docker", "Database Design", "Distributed Systems"],
      learningOutcomes: [
        "Design microservices architecture",
        "Implement gRPC communication",
        "Set up service discovery with Consul",
        "Build API gateway patterns",
        "Implement distributed tracing",
        "Handle service resilience",
        "Deploy with Kubernetes",
        "Monitor with Prometheus"
      ],
      chapters: [
        {
          id: "architecture",
          title: "System Architecture Design",
          description: "Design the overall microservices architecture",
          estimatedTime: "3 hours",
          completed: false,
          locked: false
        },
        {
          id: "services",
          title: "Core Services Implementation",
          description: "Build user, product, and order services",
          estimatedTime: "6 hours",
          completed: false,
          locked: true
        },
        {
          id: "communication",
          title: "Service Communication",
          description: "Implement gRPC and HTTP communication",
          estimatedTime: "4 hours",
          completed: false,
          locked: true
        },
        {
          id: "gateway",
          title: "API Gateway & Discovery",
          description: "Set up API gateway and service discovery",
          estimatedTime: "4 hours",
          completed: false,
          locked: true
        },
        {
          id: "observability",
          title: "Monitoring & Deployment",
          description: "Add monitoring, logging, and deploy to Kubernetes",
          estimatedTime: "5 hours",
          completed: false,
          locked: true
        }
      ],
      completed: false,
      locked: true,
      progress: 0,
      icon: Layers,
      category: "Distributed Systems",
      githubRepo: "https://github.com/go-pro/microservices-system"
    }
  ];

  const difficulties = ["all", "Beginner", "Intermediate", "Advanced"];
  const categories = ["all", "CLI Applications", "Web APIs", "Real-time Applications", "Distributed Systems"];

  const filteredProjects = projects.filter(project => {
    const matchesDifficulty = selectedDifficulty === "all" || project.difficulty === selectedDifficulty;
    const matchesTab = activeTab === "all" || 
                      (activeTab === "available" && !project.locked) ||
                      (activeTab === "completed" && project.completed) ||
                      (activeTab === "in-progress" && project.progress > 0 && !project.completed);
    
    return matchesDifficulty && matchesTab;
  });

  const getDifficultyColor = (difficulty: string) => {
    switch (difficulty) {
      case "Beginner": return "text-green-600 bg-green-50 border-green-200";
      case "Intermediate": return "text-yellow-600 bg-yellow-50 border-yellow-200";
      case "Advanced": return "text-red-600 bg-red-50 border-red-200";
      default: return "text-gray-600 bg-gray-50 border-gray-200";
    }
  };

  const projectStats = {
    totalProjects: projects.length,
    completedProjects: projects.filter(p => p.completed).length,
    inProgressProjects: projects.filter(p => p.progress > 0 && !p.completed).length,
    totalHours: projects.reduce((sum, p) => sum + parseInt(p.estimatedTime.split('-')[0]), 0)
  };

  return (
    <div className="min-h-screen bg-gradient-to-br from-background via-background to-accent/5">
      <div className="container max-w-7xl mx-auto px-4 py-8 sm:px-6 sm:py-10 lg:px-8 lg:py-12">
        {/* Header */}
        <div className="mb-10 lg:mb-12">
          <div className="flex flex-col lg:flex-row lg:items-center lg:justify-between mb-6 lg:mb-8">
            <div className="mb-4 lg:mb-0">
              <h1 className="text-3xl lg:text-4xl xl:text-5xl font-bold tracking-tight mb-3 bg-gradient-to-r from-primary to-primary/70 bg-clip-text text-transparent">
                Real-World Projects
              </h1>
              <p className="text-muted-foreground text-lg lg:text-xl max-w-2xl">
                Build production-ready applications and master Go through hands-on projects
              </p>
            </div>
            <div className="text-left lg:text-right">
              <div className="text-3xl lg:text-4xl font-bold text-primary">{projectStats.completedProjects}/{projectStats.totalProjects}</div>
              <div className="text-sm lg:text-base text-muted-foreground">Projects Completed</div>
            </div>
          </div>

          {/* Project Stats */}
          <div className="grid grid-cols-2 md:grid-cols-4 gap-4 lg:gap-6 mb-8">
          <Card>
            <CardContent className="p-4 text-center">
              <Trophy className="h-6 w-6 text-yellow-500 mx-auto mb-2" />
              <div className="text-2xl font-bold">{projectStats.completedProjects}</div>
              <div className="text-sm text-muted-foreground">Completed</div>
            </CardContent>
          </Card>
          <Card>
            <CardContent className="p-4 text-center">
              <Zap className="h-6 w-6 text-blue-500 mx-auto mb-2" />
              <div className="text-2xl font-bold">{projectStats.inProgressProjects}</div>
              <div className="text-sm text-muted-foreground">In Progress</div>
            </CardContent>
          </Card>
          <Card>
            <CardContent className="p-4 text-center">
              <Clock className="h-6 w-6 text-purple-500 mx-auto mb-2" />
              <div className="text-2xl font-bold">{projectStats.totalHours}+</div>
              <div className="text-sm text-muted-foreground">Total Hours</div>
            </CardContent>
          </Card>
          <Card>
            <CardContent className="p-4 text-center">
              <Target className="h-6 w-6 text-green-500 mx-auto mb-2" />
              <div className="text-2xl font-bold">{projects.filter(p => !p.locked).length}</div>
              <div className="text-sm text-muted-foreground">Available</div>
            </CardContent>
          </Card>
        </div>

        <Progress 
          value={(projectStats.completedProjects / projectStats.totalProjects) * 100} 
          className="h-2" 
        />
      </div>

      {/* Main Content */}
      <Tabs value={activeTab} onValueChange={setActiveTab} className="space-y-6">
        <div className="flex flex-col md:flex-row md:items-center md:justify-between gap-4">
          <TabsList className="grid w-full grid-cols-4 md:w-[500px]">
            <TabsTrigger value="all">All Projects</TabsTrigger>
            <TabsTrigger value="available">Available</TabsTrigger>
            <TabsTrigger value="in-progress">In Progress</TabsTrigger>
            <TabsTrigger value="completed">Completed</TabsTrigger>
          </TabsList>

          <div className="flex gap-2">
            <select
              value={selectedDifficulty}
              onChange={(e) => setSelectedDifficulty(e.target.value)}
              className="px-3 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-primary"
            >
              {difficulties.map(difficulty => (
                <option key={difficulty} value={difficulty}>
                  {difficulty === "all" ? "All Levels" : difficulty}
                </option>
              ))}
            </select>
          </div>
        </div>

        {/* Projects Grid */}
        <TabsContent value={activeTab} className="space-y-8">
          <div className="grid grid-cols-1 lg:grid-cols-2 xl:grid-cols-3 gap-6 lg:gap-8">
            {filteredProjects.map((project) => (
              <Card key={project.id} className={`relative ${project.locked ? 'opacity-60' : ''}`}>
                {project.locked && (
                  <div className="absolute top-4 right-4 z-10">
                    <Lock className="h-5 w-5 text-muted-foreground" />
                  </div>
                )}
                <CardHeader>
                  <div className="flex items-start justify-between">
                    <div className="flex items-center space-x-3">
                      <div className="p-2 rounded-lg bg-primary/10">
                        <project.icon className="h-6 w-6 text-primary" />
                      </div>
                      <div>
                        <CardTitle className="text-xl mb-1">{project.title}</CardTitle>
                        <CardDescription className="text-sm">
                          {project.description}
                        </CardDescription>
                      </div>
                    </div>
                  </div>
                  <div className="flex items-center gap-2 mt-3">
                    <Badge className={getDifficultyColor(project.difficulty)}>
                      {project.difficulty}
                    </Badge>
                    <Badge variant="outline">{project.category}</Badge>
                    {project.completed && (
                      <Badge className="bg-green-100 text-green-800 border-green-200">
                        <CheckCircle className="mr-1 h-3 w-3" />
                        Completed
                      </Badge>
                    )}
                  </div>
                </CardHeader>
                <CardContent>
                  <div className="space-y-4">
                    <p className="text-sm text-muted-foreground line-clamp-3">
                      {project.longDescription}
                    </p>

                    <div className="flex items-center justify-between text-sm">
                      <div className="flex items-center space-x-4">
                        <div className="flex items-center space-x-1">
                          <Clock className="h-4 w-4 text-muted-foreground" />
                          <span>{project.estimatedTime}</span>
                        </div>
                        <div className="flex items-center space-x-1">
                          <BookOpen className="h-4 w-4 text-muted-foreground" />
                          <span>{project.chapters.length} chapters</span>
                        </div>
                      </div>
                    </div>

                    {/* Progress */}
                    {project.progress > 0 && (
                      <div className="space-y-2">
                        <div className="flex justify-between text-sm">
                          <span>Progress</span>
                          <span>{project.progress}%</span>
                        </div>
                        <Progress value={project.progress} className="h-2" />
                      </div>
                    )}

                    {/* Technologies */}
                    <div className="flex flex-wrap gap-1">
                      {project.technologies.slice(0, 4).map((tech) => (
                        <Badge key={tech} variant="secondary" className="text-xs">
                          {tech}
                        </Badge>
                      ))}
                      {project.technologies.length > 4 && (
                        <Badge variant="secondary" className="text-xs">
                          +{project.technologies.length - 4} more
                        </Badge>
                      )}
                    </div>

                    {/* Action Buttons */}
                    <div className="flex items-center justify-between pt-2">
                      <div className="flex items-center space-x-2">
                        {project.githubRepo && (
                          <Button variant="outline" size="sm" asChild>
                            <a href={project.githubRepo} target="_blank" rel="noopener noreferrer">
                              <GitBranch className="mr-1 h-3 w-3" />
                              Code
                            </a>
                          </Button>
                        )}
                        {project.liveDemo && (
                          <Button variant="outline" size="sm" asChild>
                            <a href={project.liveDemo} target="_blank" rel="noopener noreferrer">
                              <Globe className="mr-1 h-3 w-3" />
                              Demo
                            </a>
                          </Button>
                        )}
                      </div>

                      <Link href={`/projects/${project.id}`}>
                        <Button 
                          disabled={project.locked}
                          variant={project.completed ? "outline" : "default"}
                          className={!project.completed && !project.locked ? "go-gradient text-white" : ""}
                        >
                          {project.locked ? (
                            <>
                              <Lock className="mr-2 h-4 w-4" />
                              Locked
                            </>
                          ) : project.completed ? (
                            <>
                              <CheckCircle className="mr-2 h-4 w-4" />
                              Review
                            </>
                          ) : project.progress > 0 ? (
                            <>
                              <Play className="mr-2 h-4 w-4" />
                              Continue
                            </>
                          ) : (
                            <>
                              <Play className="mr-2 h-4 w-4" />
                              Start Project
                            </>
                          )}
                        </Button>
                      </Link>
                    </div>
                  </div>
                </CardContent>
              </Card>
            ))}
          </div>
        </TabsContent>
      </Tabs>
      </div>
    </div>
  );
}
