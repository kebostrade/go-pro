"use client";

import { useState, useEffect } from "react";
import { useParams, useRouter } from "next/navigation";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Progress } from "@/components/ui/progress";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import CodeEditor from "@/components/learning/code-editor";
import {
  ArrowLeft,
  Clock,
  Star,
  Target,
  CheckCircle,
  Play,
  BookOpen,
  Code2,
  GitBranch,
  Globe,
  Home,
  ChevronRight,
  Terminal,
  FileText,
  Lightbulb,
  Trophy,
  Users,
  Download,
  ExternalLink,
  ArrowRight,
  Lock,
  MessageSquare,
  Share2
} from "lucide-react";
import Link from "next/link";

interface Chapter {
  id: string;
  title: string;
  description: string;
  estimatedTime: string;
  completed: boolean;
  locked: boolean;
  content: string;
  codeExample?: string;
  tasks: Task[];
}

interface Task {
  id: string;
  title: string;
  description: string;
  completed: boolean;
}

interface ProjectData {
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
  progress: number;
  icon: any;
  category: string;
  githubRepo?: string;
  liveDemo?: string;
}

export default function ProjectDetailPage() {
  const params = useParams();
  const router = useRouter();
  const [activeTab, setActiveTab] = useState("overview");
  const [projectData, setProjectData] = useState<ProjectData | null>(null);
  const [loading, setLoading] = useState(true);
  const [currentChapter, setCurrentChapter] = useState(0);

  const projectId = params.id as string;

  useEffect(() => {
    // Mock data loading based on project ID
    const mockProject: ProjectData = {
      id: projectId,
      title: "CLI Task Manager",
      description: "Build a command-line task management application with file persistence",
      longDescription: "Create a full-featured command-line task manager that demonstrates Go's CLI capabilities, file I/O, JSON handling, and clean architecture patterns. This project covers fundamental Go concepts while building something practical that you can use in your daily workflow.",
      difficulty: "Beginner",
      estimatedTime: "8-10 hours",
      technologies: ["Go", "CLI", "JSON", "File I/O", "Cobra", "Testing"],
      prerequisites: ["Go Basics", "Functions", "Structs", "Error Handling", "JSON"],
      learningOutcomes: [
        "Build command-line applications with Cobra framework",
        "Handle file I/O operations safely and efficiently",
        "Work with JSON serialization and deserialization",
        "Implement clean architecture patterns in Go",
        "Create comprehensive CLI interfaces with subcommands",
        "Handle user input validation and error cases",
        "Write unit tests for CLI applications",
        "Package and distribute Go applications"
      ],
      chapters: [
        {
          id: "setup",
          title: "Project Setup & Architecture",
          description: "Set up the project structure and define the core architecture",
          estimatedTime: "1 hour",
          completed: true,
          locked: false,
          content: `# Project Setup & Architecture

In this chapter, we'll set up our CLI task manager project and establish a clean architecture that will make our code maintainable and testable.

## Project Structure

We'll organize our project using Go's standard project layout:

\`\`\`
task-manager/
├── cmd/
│   └── root.go          # Root command and CLI setup
├── internal/
│   ├── models/          # Data models
│   ├── storage/         # File storage logic
│   └── commands/        # Command implementations
├── pkg/
│   └── utils/           # Utility functions
├── main.go              # Application entry point
├── go.mod
└── README.md
\`\`\`

## Architecture Principles

1. **Separation of Concerns**: Each package has a single responsibility
2. **Dependency Injection**: Dependencies are injected rather than hardcoded
3. **Interface-based Design**: Use interfaces for testability
4. **Error Handling**: Proper error handling throughout the application

## Getting Started

Let's initialize our Go module and set up the basic structure.`,
          codeExample: `package main

import (
    "fmt"
    "os"
    
    "github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
    Use:   "task",
    Short: "A simple CLI task manager",
    Long:  "A command-line task manager built with Go and Cobra",
}

func Execute() {
    if err := rootCmd.Execute(); err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
}

func main() {
    Execute()
}`,
          tasks: [
            {
              id: "init-module",
              title: "Initialize Go Module",
              description: "Create a new Go module for the project",
              completed: true
            },
            {
              id: "setup-cobra",
              title: "Set up Cobra CLI Framework",
              description: "Install and configure Cobra for command-line interface",
              completed: true
            },
            {
              id: "create-structure",
              title: "Create Project Structure",
              description: "Set up the directory structure and basic files",
              completed: false
            }
          ]
        },
        {
          id: "models",
          title: "Data Models & Storage",
          description: "Create task models and implement JSON file storage",
          estimatedTime: "2 hours",
          completed: true,
          locked: false,
          content: `# Data Models & Storage

In this chapter, we'll define our data models and implement persistent storage using JSON files.

## Task Model

Our task model will include all the essential fields for a task management system:

- ID: Unique identifier
- Title: Task description
- Completed: Status flag
- CreatedAt: Timestamp
- DueDate: Optional due date
- Priority: Task priority level

## Storage Interface

We'll create a storage interface that can be implemented by different storage backends (file, database, etc.).

## JSON File Storage

For this project, we'll implement file-based storage using JSON format for simplicity and portability.`,
          codeExample: `package models

import (
    "time"
)

type Priority int

const (
    Low Priority = iota
    Medium
    High
)

type Task struct {
    ID        string    \`json:"id"\`
    Title     string    \`json:"title"\`
    Completed bool      \`json:"completed"\`
    CreatedAt time.Time \`json:"created_at"\`
    DueDate   *time.Time \`json:"due_date,omitempty"\`
    Priority  Priority  \`json:"priority"\`
}

type TaskStorage interface {
    Save(tasks []Task) error
    Load() ([]Task, error)
    Exists() bool
}`,
          tasks: [
            {
              id: "define-models",
              title: "Define Task Model",
              description: "Create the Task struct with all necessary fields",
              completed: true
            },
            {
              id: "storage-interface",
              title: "Create Storage Interface",
              description: "Define the interface for task storage operations",
              completed: true
            },
            {
              id: "json-storage",
              title: "Implement JSON Storage",
              description: "Create JSON file storage implementation",
              completed: false
            }
          ]
        },
        {
          id: "commands",
          title: "CLI Commands Implementation",
          description: "Build add, list, complete, and delete commands",
          estimatedTime: "3 hours",
          completed: false,
          locked: false,
          content: `# CLI Commands Implementation

Now we'll implement the core commands for our task manager: add, list, complete, and delete.

## Command Structure

Each command will be implemented as a separate Cobra command with its own logic and validation.

## Commands to Implement

1. **add**: Add a new task
2. **list**: List all tasks with filtering options
3. **complete**: Mark a task as completed
4. **delete**: Remove a task
5. **edit**: Edit an existing task

## Input Validation

We'll implement proper input validation and error handling for each command.`,
          codeExample: `package commands

import (
    "fmt"
    "strconv"
    "time"
    
    "github.com/spf13/cobra"
    "your-module/internal/models"
    "your-module/internal/storage"
)

func NewAddCommand(storage models.TaskStorage) *cobra.Command {
    var priority string
    var dueDate string
    
    cmd := &cobra.Command{
        Use:   "add [task description]",
        Short: "Add a new task",
        Args:  cobra.MinimumNArgs(1),
        RunE: func(cmd *cobra.Command, args []string) error {
            title := strings.Join(args, " ")
            
            task := models.Task{
                ID:        generateID(),
                Title:     title,
                Completed: false,
                CreatedAt: time.Now(),
                Priority:  parsePriority(priority),
            }
            
            if dueDate != "" {
                due, err := time.Parse("2006-01-02", dueDate)
                if err != nil {
                    return fmt.Errorf("invalid due date format: %v", err)
                }
                task.DueDate = &due
            }
            
            return addTask(storage, task)
        },
    }
    
    cmd.Flags().StringVarP(&priority, "priority", "p", "medium", "Task priority (low, medium, high)")
    cmd.Flags().StringVarP(&dueDate, "due", "d", "", "Due date (YYYY-MM-DD)")
    
    return cmd
}`,
          tasks: [
            {
              id: "add-command",
              title: "Implement Add Command",
              description: "Create command to add new tasks",
              completed: false
            },
            {
              id: "list-command",
              title: "Implement List Command",
              description: "Create command to list tasks with filtering",
              completed: false
            },
            {
              id: "complete-command",
              title: "Implement Complete Command",
              description: "Create command to mark tasks as completed",
              completed: false
            },
            {
              id: "delete-command",
              title: "Implement Delete Command",
              description: "Create command to delete tasks",
              completed: false
            }
          ]
        },
        {
          id: "advanced",
          title: "Advanced Features",
          description: "Add filtering, sorting, and due date functionality",
          estimatedTime: "2 hours",
          completed: false,
          locked: true,
          content: `# Advanced Features

In this chapter, we'll add advanced features to make our task manager more powerful and user-friendly.

## Features to Implement

1. **Filtering**: Filter tasks by status, priority, or due date
2. **Sorting**: Sort tasks by different criteria
3. **Search**: Search tasks by title or description
4. **Statistics**: Show task statistics and productivity metrics
5. **Export**: Export tasks to different formats

## Enhanced User Experience

We'll also improve the user experience with better formatting, colors, and interactive features.`,
          tasks: [
            {
              id: "filtering",
              title: "Add Task Filtering",
              description: "Implement filtering by status, priority, and date",
              completed: false
            },
            {
              id: "sorting",
              title: "Add Task Sorting",
              description: "Implement sorting by different criteria",
              completed: false
            },
            {
              id: "search",
              title: "Add Search Functionality",
              description: "Implement task search by title",
              completed: false
            }
          ]
        },
        {
          id: "testing",
          title: "Testing & Documentation",
          description: "Write comprehensive tests and documentation",
          estimatedTime: "2 hours",
          completed: false,
          locked: true,
          content: `# Testing & Documentation

In this final chapter, we'll add comprehensive tests and documentation to ensure our application is reliable and maintainable.

## Testing Strategy

1. **Unit Tests**: Test individual functions and methods
2. **Integration Tests**: Test command interactions
3. **CLI Tests**: Test the complete CLI workflow

## Documentation

1. **Code Documentation**: Add comprehensive comments
2. **README**: Create detailed usage instructions
3. **Examples**: Provide usage examples

## Distribution

Learn how to build and distribute your CLI application.`,
          tasks: [
            {
              id: "unit-tests",
              title: "Write Unit Tests",
              description: "Create unit tests for core functionality",
              completed: false
            },
            {
              id: "integration-tests",
              title: "Write Integration Tests",
              description: "Create integration tests for CLI commands",
              completed: false
            },
            {
              id: "documentation",
              title: "Create Documentation",
              description: "Write comprehensive documentation and examples",
              completed: false
            }
          ]
        }
      ],
      completed: false,
      progress: 40,
      icon: Terminal,
      category: "CLI Applications",
      githubRepo: "https://github.com/go-pro/cli-task-manager"
    };

    setProjectData(mockProject);
    setLoading(false);
  }, [projectId]);

  const getDifficultyColor = (difficulty: string) => {
    switch (difficulty) {
      case "Beginner": return "text-green-600 bg-green-50 border-green-200";
      case "Intermediate": return "text-yellow-600 bg-yellow-50 border-yellow-200";
      case "Advanced": return "text-red-600 bg-red-50 border-red-200";
      default: return "text-gray-600 bg-gray-50 border-gray-200";
    }
  };

  if (loading || !projectData) {
    return (
      <div className="container max-w-screen-2xl px-4 py-8">
        <div className="flex items-center justify-center h-64">
          <div className="text-center">
            <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-primary mx-auto mb-4"></div>
            <p className="text-muted-foreground">Loading project...</p>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gradient-to-br from-background via-background to-accent/5">
      <div className="container max-w-7xl mx-auto px-4 py-8 sm:px-6 sm:py-10 lg:px-8 lg:py-12">
        {/* Breadcrumb Navigation */}
        <div className="flex items-center space-x-2 text-sm text-muted-foreground mb-6 lg:mb-8">
        <Link href="/" className="hover:text-primary">
          <Home className="h-4 w-4" />
        </Link>
        <ChevronRight className="h-4 w-4" />
        <Link href="/projects" className="hover:text-primary">
          Projects
        </Link>
        <ChevronRight className="h-4 w-4" />
        <span className="text-foreground">{projectData.title}</span>
      </div>

      {/* Project Header */}
      <div className="mb-10 lg:mb-12">
        <div className="flex flex-col lg:flex-row lg:items-start lg:justify-between mb-6 lg:mb-8">
          <div className="flex-1 mb-6 lg:mb-0">
            <div className="flex items-center space-x-4 mb-4">
              <div className="p-3 rounded-xl bg-primary/10">
                <projectData.icon className="h-10 w-10 lg:h-12 lg:w-12 text-primary" />
              </div>
              <div className="flex flex-wrap gap-2">
                <Badge className={getDifficultyColor(projectData.difficulty)}>
                  {projectData.difficulty}
                </Badge>
                <Badge variant="outline">{projectData.category}</Badge>
              </div>
            </div>
            <h1 className="text-3xl lg:text-4xl xl:text-5xl font-bold tracking-tight mb-4 bg-gradient-to-r from-primary to-primary/70 bg-clip-text text-transparent">
              {projectData.title}
            </h1>
            <p className="text-muted-foreground text-lg lg:text-xl mb-6 max-w-3xl">{projectData.description}</p>
            
            <div className="flex items-center space-x-6 text-sm">
              <div className="flex items-center space-x-1">
                <Clock className="h-4 w-4" />
                <span>{projectData.estimatedTime}</span>
              </div>
              <div className="flex items-center space-x-1">
                <BookOpen className="h-4 w-4" />
                <span>{projectData.chapters.length} chapters</span>
              </div>
              <div className="flex items-center space-x-1">
                <Target className="h-4 w-4" />
                <span>{projectData.learningOutcomes.length} outcomes</span>
              </div>
            </div>
          </div>
          <div className="ml-6 flex flex-col gap-2">
            <Link href="/projects">
              <Button variant="outline" size="sm">
                <ArrowLeft className="mr-2 h-4 w-4" />
                Back to Projects
              </Button>
            </Link>
            {projectData.githubRepo && (
              <Button variant="outline" size="sm" asChild>
                <a href={projectData.githubRepo} target="_blank" rel="noopener noreferrer">
                  <GitBranch className="mr-2 h-4 w-4" />
                  View Code
                </a>
              </Button>
            )}
          </div>
        </div>

        {/* Progress */}
        {projectData.progress > 0 && (
          <div className="mb-4">
            <div className="flex items-center justify-between text-sm mb-2">
              <span>Overall Progress</span>
              <span>{projectData.progress}%</span>
            </div>
            <Progress value={projectData.progress} className="h-2" />
          </div>
        )}
      </div>

      {/* Main Content */}
      <Tabs value={activeTab} onValueChange={setActiveTab} className="space-y-6">
        <TabsList className="grid w-full grid-cols-4 lg:w-[500px]">
          <TabsTrigger value="overview">Overview</TabsTrigger>
          <TabsTrigger value="chapters">Chapters</TabsTrigger>
          <TabsTrigger value="resources">Resources</TabsTrigger>
          <TabsTrigger value="community">Community</TabsTrigger>
        </TabsList>

        {/* Overview Tab */}
        <TabsContent value="overview" className="space-y-6">
          <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
            {/* Project Description */}
            <div className="lg:col-span-2 space-y-6">
              <Card>
                <CardHeader>
                  <CardTitle className="flex items-center">
                    <FileText className="mr-2 h-5 w-5" />
                    Project Description
                  </CardTitle>
                </CardHeader>
                <CardContent>
                  <p className="text-muted-foreground leading-relaxed">
                    {projectData.longDescription}
                  </p>
                </CardContent>
              </Card>

              <Card>
                <CardHeader>
                  <CardTitle className="flex items-center">
                    <Target className="mr-2 h-5 w-5" />
                    Learning Outcomes
                  </CardTitle>
                  <CardDescription>
                    What you'll learn by completing this project
                  </CardDescription>
                </CardHeader>
                <CardContent>
                  <ul className="space-y-2">
                    {projectData.learningOutcomes.map((outcome, index) => (
                      <li key={index} className="flex items-start space-x-2">
                        <CheckCircle className="h-4 w-4 text-green-500 mt-0.5 flex-shrink-0" />
                        <span className="text-sm">{outcome}</span>
                      </li>
                    ))}
                  </ul>
                </CardContent>
              </Card>
            </div>

            {/* Sidebar */}
            <div className="space-y-6">
              <Card>
                <CardHeader>
                  <CardTitle className="flex items-center">
                    <Code2 className="mr-2 h-5 w-5" />
                    Technologies
                  </CardTitle>
                </CardHeader>
                <CardContent>
                  <div className="flex flex-wrap gap-2">
                    {projectData.technologies.map((tech) => (
                      <Badge key={tech} variant="secondary">
                        {tech}
                      </Badge>
                    ))}
                  </div>
                </CardContent>
              </Card>

              <Card>
                <CardHeader>
                  <CardTitle className="flex items-center">
                    <Lightbulb className="mr-2 h-5 w-5" />
                    Prerequisites
                  </CardTitle>
                </CardHeader>
                <CardContent>
                  <ul className="space-y-2">
                    {projectData.prerequisites.map((prereq, index) => (
                      <li key={index} className="flex items-center space-x-2">
                        <div className="w-2 h-2 rounded-full bg-primary"></div>
                        <span className="text-sm">{prereq}</span>
                      </li>
                    ))}
                  </ul>
                </CardContent>
              </Card>

              <Card>
                <CardHeader>
                  <CardTitle>Quick Start</CardTitle>
                </CardHeader>
                <CardContent>
                  <div className="space-y-3">
                    <Button className="w-full go-gradient text-white">
                      <Play className="mr-2 h-4 w-4" />
                      Start Project
                    </Button>
                    {projectData.githubRepo && (
                      <Button variant="outline" className="w-full" asChild>
                        <a href={projectData.githubRepo} target="_blank" rel="noopener noreferrer">
                          <Download className="mr-2 h-4 w-4" />
                          Clone Repository
                        </a>
                      </Button>
                    )}
                  </div>
                </CardContent>
              </Card>
            </div>
          </div>
        </TabsContent>

        {/* Chapters Tab */}
        <TabsContent value="chapters" className="space-y-6">
          <div className="grid grid-cols-1 lg:grid-cols-4 gap-6">
            {/* Chapter List */}
            <div className="lg:col-span-1">
              <Card>
                <CardHeader>
                  <CardTitle>Chapters</CardTitle>
                  <CardDescription>
                    {projectData.chapters.filter(c => c.completed).length} of {projectData.chapters.length} completed
                  </CardDescription>
                </CardHeader>
                <CardContent>
                  <div className="space-y-2">
                    {projectData.chapters.map((chapter, index) => (
                      <button
                        key={chapter.id}
                        onClick={() => setCurrentChapter(index)}
                        disabled={chapter.locked}
                        className={`w-full text-left p-3 rounded-lg border transition-colors ${
                          currentChapter === index
                            ? 'border-primary bg-primary/5'
                            : chapter.locked
                              ? 'border-gray-200 bg-gray-50 opacity-60 cursor-not-allowed'
                              : 'border-gray-200 hover:border-primary/50 hover:bg-primary/5'
                        }`}
                      >
                        <div className="flex items-center justify-between mb-1">
                          <span className="font-medium text-sm">
                            Chapter {index + 1}
                          </span>
                          <div className="flex items-center space-x-1">
                            {chapter.completed && (
                              <CheckCircle className="h-4 w-4 text-green-500" />
                            )}
                            {chapter.locked && (
                              <Lock className="h-4 w-4 text-gray-400" />
                            )}
                          </div>
                        </div>
                        <div className="text-sm font-medium mb-1">{chapter.title}</div>
                        <div className="text-xs text-muted-foreground">
                          {chapter.estimatedTime}
                        </div>
                      </button>
                    ))}
                  </div>
                </CardContent>
              </Card>
            </div>

            {/* Chapter Content */}
            <div className="lg:col-span-3">
              {projectData.chapters[currentChapter] && (
                <div className="space-y-6">
                  <Card>
                    <CardHeader>
                      <div className="flex items-center justify-between">
                        <div>
                          <CardTitle className="flex items-center">
                            Chapter {currentChapter + 1}: {projectData.chapters[currentChapter].title}
                            {projectData.chapters[currentChapter].completed && (
                              <CheckCircle className="ml-2 h-5 w-5 text-green-500" />
                            )}
                          </CardTitle>
                          <CardDescription>
                            {projectData.chapters[currentChapter].description} • {projectData.chapters[currentChapter].estimatedTime}
                          </CardDescription>
                        </div>
                        <Badge variant={projectData.chapters[currentChapter].completed ? "default" : "secondary"}>
                          {projectData.chapters[currentChapter].completed ? "Completed" : "In Progress"}
                        </Badge>
                      </div>
                    </CardHeader>
                    <CardContent>
                      <div className="prose dark:prose-invert max-w-none">
                        <div dangerouslySetInnerHTML={{
                          __html: projectData.chapters[currentChapter].content
                            .replace(/\n/g, '<br>')
                            .replace(/```([\s\S]*?)```/g, '<pre><code>$1</code></pre>')
                            .replace(/`([^`]+)`/g, '<code>$1</code>')
                            .replace(/\*\*(.*?)\*\*/g, '<strong>$1</strong>')
                            .replace(/## (.*?)(<br>|$)/g, '<h2>$1</h2>')
                            .replace(/# (.*?)(<br>|$)/g, '<h1>$1</h1>')
                        }} />
                      </div>
                    </CardContent>
                  </Card>

                  {/* Code Example */}
                  {projectData.chapters[currentChapter].codeExample && (
                    <CodeEditor
                      title="Code Example"
                      description="Study this code example and try running it"
                      initialCode={projectData.chapters[currentChapter].codeExample}
                      language="go"
                      readOnly={false}
                    />
                  )}

                  {/* Chapter Tasks */}
                  <Card>
                    <CardHeader>
                      <CardTitle className="flex items-center">
                        <Trophy className="mr-2 h-5 w-5" />
                        Chapter Tasks
                      </CardTitle>
                      <CardDescription>
                        Complete these tasks to finish the chapter
                      </CardDescription>
                    </CardHeader>
                    <CardContent>
                      <div className="space-y-3">
                        {projectData.chapters[currentChapter].tasks.map((task) => (
                          <div key={task.id} className={`p-3 rounded-lg border ${
                            task.completed ? 'bg-green-50 border-green-200' : 'bg-gray-50 border-gray-200'
                          }`}>
                            <div className="flex items-start space-x-3">
                              <div className="mt-1">
                                {task.completed ? (
                                  <CheckCircle className="h-4 w-4 text-green-600" />
                                ) : (
                                  <div className="w-4 h-4 rounded-full border-2 border-gray-300"></div>
                                )}
                              </div>
                              <div className="flex-1">
                                <div className={`font-medium ${task.completed ? 'text-green-800' : 'text-foreground'}`}>
                                  {task.title}
                                </div>
                                <div className={`text-sm ${task.completed ? 'text-green-600' : 'text-muted-foreground'}`}>
                                  {task.description}
                                </div>
                              </div>
                            </div>
                          </div>
                        ))}
                      </div>
                    </CardContent>
                  </Card>

                  {/* Chapter Navigation */}
                  <div className="flex items-center justify-between">
                    <Button
                      variant="outline"
                      onClick={() => setCurrentChapter(Math.max(0, currentChapter - 1))}
                      disabled={currentChapter === 0}
                    >
                      <ArrowLeft className="mr-2 h-4 w-4" />
                      Previous Chapter
                    </Button>

                    <div className="flex items-center space-x-2">
                      {!projectData.chapters[currentChapter].completed && (
                        <Button className="go-gradient text-white">
                          <CheckCircle className="mr-2 h-4 w-4" />
                          Mark Complete
                        </Button>
                      )}

                      <Button
                        variant="outline"
                        onClick={() => setCurrentChapter(Math.min(projectData.chapters.length - 1, currentChapter + 1))}
                        disabled={currentChapter === projectData.chapters.length - 1 || projectData.chapters[currentChapter + 1]?.locked}
                      >
                        Next Chapter
                        <ArrowRight className="ml-2 h-4 w-4" />
                      </Button>
                    </div>
                  </div>
                </div>
              )}
            </div>
          </div>
        </TabsContent>

        {/* Resources Tab */}
        <TabsContent value="resources" className="space-y-6">
          <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center">
                  <GitBranch className="mr-2 h-5 w-5" />
                  Code Repository
                </CardTitle>
                <CardDescription>
                  Access the complete source code and examples
                </CardDescription>
              </CardHeader>
              <CardContent>
                <div className="space-y-4">
                  <p className="text-sm text-muted-foreground">
                    The complete source code for this project is available on GitHub. You can clone it, study the implementation, and use it as a reference.
                  </p>
                  <div className="flex gap-2">
                    {projectData.githubRepo && (
                      <Button asChild>
                        <a href={projectData.githubRepo} target="_blank" rel="noopener noreferrer">
                          <GitBranch className="mr-2 h-4 w-4" />
                          View on GitHub
                        </a>
                      </Button>
                    )}
                    <Button variant="outline">
                      <Download className="mr-2 h-4 w-4" />
                      Download ZIP
                    </Button>
                  </div>
                </div>
              </CardContent>
            </Card>

            <Card>
              <CardHeader>
                <CardTitle className="flex items-center">
                  <BookOpen className="mr-2 h-5 w-5" />
                  Documentation
                </CardTitle>
                <CardDescription>
                  Additional resources and documentation
                </CardDescription>
              </CardHeader>
              <CardContent>
                <div className="space-y-3">
                  <a href="#" className="flex items-center justify-between p-3 rounded-lg border hover:bg-muted/50 transition-colors">
                    <div>
                      <div className="font-medium">API Reference</div>
                      <div className="text-sm text-muted-foreground">Complete API documentation</div>
                    </div>
                    <ExternalLink className="h-4 w-4 text-muted-foreground" />
                  </a>
                  <a href="#" className="flex items-center justify-between p-3 rounded-lg border hover:bg-muted/50 transition-colors">
                    <div>
                      <div className="font-medium">Go CLI Best Practices</div>
                      <div className="text-sm text-muted-foreground">Guidelines for building CLI apps</div>
                    </div>
                    <ExternalLink className="h-4 w-4 text-muted-foreground" />
                  </a>
                  <a href="#" className="flex items-center justify-between p-3 rounded-lg border hover:bg-muted/50 transition-colors">
                    <div>
                      <div className="font-medium">Cobra Framework Guide</div>
                      <div className="text-sm text-muted-foreground">Official Cobra documentation</div>
                    </div>
                    <ExternalLink className="h-4 w-4 text-muted-foreground" />
                  </a>
                </div>
              </CardContent>
            </Card>

            <Card>
              <CardHeader>
                <CardTitle className="flex items-center">
                  <Terminal className="mr-2 h-5 w-5" />
                  Quick Commands
                </CardTitle>
                <CardDescription>
                  Essential commands for this project
                </CardDescription>
              </CardHeader>
              <CardContent>
                <div className="space-y-3">
                  <div className="p-3 bg-muted rounded-lg">
                    <div className="text-sm font-mono">go mod init task-manager</div>
                    <div className="text-xs text-muted-foreground mt-1">Initialize Go module</div>
                  </div>
                  <div className="p-3 bg-muted rounded-lg">
                    <div className="text-sm font-mono">go get github.com/spf13/cobra</div>
                    <div className="text-xs text-muted-foreground mt-1">Install Cobra CLI framework</div>
                  </div>
                  <div className="p-3 bg-muted rounded-lg">
                    <div className="text-sm font-mono">go build -o task</div>
                    <div className="text-xs text-muted-foreground mt-1">Build the application</div>
                  </div>
                  <div className="p-3 bg-muted rounded-lg">
                    <div className="text-sm font-mono">go test ./...</div>
                    <div className="text-xs text-muted-foreground mt-1">Run all tests</div>
                  </div>
                </div>
              </CardContent>
            </Card>

            <Card>
              <CardHeader>
                <CardTitle className="flex items-center">
                  <Lightbulb className="mr-2 h-5 w-5" />
                  Tips & Tricks
                </CardTitle>
                <CardDescription>
                  Helpful tips for completing this project
                </CardDescription>
              </CardHeader>
              <CardContent>
                <div className="space-y-3 text-sm">
                  <div className="p-3 bg-blue-50 border border-blue-200 rounded-lg">
                    <div className="font-medium text-blue-800">💡 Pro Tip</div>
                    <div className="text-blue-700 mt-1">Use interfaces to make your code testable and maintainable</div>
                  </div>
                  <div className="p-3 bg-green-50 border border-green-200 rounded-lg">
                    <div className="font-medium text-green-800">✅ Best Practice</div>
                    <div className="text-green-700 mt-1">Always validate user input and handle errors gracefully</div>
                  </div>
                  <div className="p-3 bg-yellow-50 border border-yellow-200 rounded-lg">
                    <div className="font-medium text-yellow-800">⚠️ Common Pitfall</div>
                    <div className="text-yellow-700 mt-1">Don't forget to handle file permissions when reading/writing files</div>
                  </div>
                </div>
              </CardContent>
            </Card>
          </div>
        </TabsContent>

        {/* Community Tab */}
        <TabsContent value="community" className="space-y-6">
          <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
            <div className="lg:col-span-2 space-y-6">
              <Card>
                <CardHeader>
                  <CardTitle className="flex items-center">
                    <Users className="mr-2 h-5 w-5" />
                    Project Discussions
                  </CardTitle>
                  <CardDescription>
                    Join the conversation about this project
                  </CardDescription>
                </CardHeader>
                <CardContent>
                  <div className="space-y-4">
                    {[
                      {
                        title: "Best practices for CLI error handling",
                        author: "Sarah Johnson",
                        replies: 12,
                        time: "2 hours ago"
                      },
                      {
                        title: "How to add color output to CLI commands?",
                        author: "Mike Chen",
                        replies: 8,
                        time: "5 hours ago"
                      },
                      {
                        title: "Testing strategies for CLI applications",
                        author: "Alex Rodriguez",
                        replies: 15,
                        time: "1 day ago"
                      }
                    ].map((discussion, index) => (
                      <div key={index} className="p-4 border rounded-lg hover:bg-muted/50 transition-colors cursor-pointer">
                        <div className="flex items-start justify-between">
                          <div className="flex-1">
                            <h4 className="font-medium mb-1">{discussion.title}</h4>
                            <div className="flex items-center space-x-4 text-sm text-muted-foreground">
                              <span>by {discussion.author}</span>
                              <span>{discussion.replies} replies</span>
                              <span>{discussion.time}</span>
                            </div>
                          </div>
                          <ArrowRight className="h-4 w-4 text-muted-foreground" />
                        </div>
                      </div>
                    ))}
                  </div>
                  <div className="pt-4 border-t">
                    <Link href="/community">
                      <Button variant="outline" className="w-full">
                        <MessageSquare className="mr-2 h-4 w-4" />
                        View All Discussions
                      </Button>
                    </Link>
                  </div>
                </CardContent>
              </Card>
            </div>

            <div className="space-y-6">
              <Card>
                <CardHeader>
                  <CardTitle className="flex items-center">
                    <Trophy className="mr-2 h-5 w-5" />
                    Project Stats
                  </CardTitle>
                </CardHeader>
                <CardContent>
                  <div className="space-y-4">
                    <div className="flex justify-between">
                      <span className="text-sm">Completed by</span>
                      <span className="font-bold">1,247 learners</span>
                    </div>
                    <div className="flex justify-between">
                      <span className="text-sm">Average rating</span>
                      <div className="flex items-center space-x-1">
                        <Star className="h-4 w-4 fill-yellow-400 text-yellow-400" />
                        <span className="font-bold">4.8</span>
                      </div>
                    </div>
                    <div className="flex justify-between">
                      <span className="text-sm">Success rate</span>
                      <span className="font-bold">89%</span>
                    </div>
                  </div>
                </CardContent>
              </Card>

              <Card>
                <CardHeader>
                  <CardTitle>Share Your Progress</CardTitle>
                </CardHeader>
                <CardContent>
                  <div className="space-y-3">
                    <Button className="w-full" variant="outline">
                      <Share2 className="mr-2 h-4 w-4" />
                      Share on Social
                    </Button>
                    <Button className="w-full" variant="outline">
                      <Users className="mr-2 h-4 w-4" />
                      Find Study Partners
                    </Button>
                  </div>
                </CardContent>
              </Card>
            </div>
          </div>
        </TabsContent>
      </Tabs>
      </div>
    </div>
  );
}
