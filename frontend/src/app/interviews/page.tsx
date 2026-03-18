"use client";

import { useState } from "react";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { Progress } from "@/components/ui/progress";
import {
  Brain,
  Code2,
  Target,
  Clock,
  BookOpen,
  Zap,
  Trophy,
  CheckCircle,
  ChevronRight,
  Lightbulb,
  Users,
  TrendingUp,
  MessageSquare,
  Layers,
  Binary,
  GitBranch,
  Database,
  LayoutGrid,
  Search,
  Settings,
  Star,
  Play,
  ArrowRight,
  Scale,
  ArrowRightLeft,
  RefreshCw,
  Shield,
} from "lucide-react";
import Link from "next/link";

interface Topic {
  id: string;
  title: string;
  description: string;
  difficulty: "Beginner" | "Intermediate" | "Advanced";
  category: string;
  estimatedTime: string;
  completed: boolean;
  problems: number;
  key: string;
}

interface Pattern {
  id: string;
  name: string;
  description: string;
  problems: string[];
  icon: React.ReactNode;
  color: string;
}

export default function InterviewsPage() {
  const [activeTab, setActiveTab] = useState("roadmap");

  const topics: Topic[] = [
    {
      id: "arrays",
      title: "Arrays & Strings",
      description: "Two pointers, sliding window, prefix sums",
      difficulty: "Beginner",
      category: "Data Structures",
      estimatedTime: "4 hours",
      completed: true,
      problems: 25,
      key: "O(n) time, O(1) space patterns"
    },
    {
      id: "hashmaps",
      title: "Hash Maps & Sets",
      description: "Fast lookups, counting, caching",
      difficulty: "Beginner",
      category: "Data Structures",
      estimatedTime: "3 hours",
      completed: true,
      problems: 20,
      key: "O(1) average lookup time"
    },
    {
      id: "linkedlists",
      title: "Linked Lists",
      description: "Reversal, cycle detection, merge",
      difficulty: "Beginner",
      category: "Data Structures",
      estimatedTime: "3 hours",
      completed: false,
      problems: 15,
      key: "Fast/slow pointers pattern"
    },
    {
      id: "trees",
      title: "Trees & BST",
      description: "Traversals, validation, LCA",
      difficulty: "Intermediate",
      category: "Data Structures",
      estimatedTime: "5 hours",
      completed: false,
      problems: 30,
      key: "DFS/BFS recursion patterns"
    },
    {
      id: "graphs",
      title: "Graphs",
      description: "BFS, DFS, Union-Find, Topological Sort",
      difficulty: "Intermediate",
      category: "Data Structures",
      estimatedTime: "6 hours",
      completed: false,
      problems: 25,
      key: "Adjacency list, visited set"
    },
    {
      id: "dp",
      title: "Dynamic Programming",
      description: "1D/2D DP, knapsack, LCS",
      difficulty: "Advanced",
      category: "Algorithms",
      estimatedTime: "8 hours",
      completed: false,
      problems: 35,
      key: "State definition → recurrence"
    },
    {
      id: "backtracking",
      title: "Backtracking",
      description: "Subsets, permutations, combinations",
      difficulty: "Intermediate",
      category: "Algorithms",
      estimatedTime: "4 hours",
      completed: false,
      problems: 15,
      key: "Choose → explore → unchoose"
    },
    {
      id: "system",
      title: "System Design",
      description: "Scalability, caching, databases",
      difficulty: "Advanced",
      category: "Design",
      estimatedTime: "10 hours",
      completed: false,
      problems: 12,
      key: "Trade-offs, CAP theorem"
    },
  ];

  const patterns: Pattern[] = [
    {
      id: "two-pointers",
      name: "Two Pointers",
      description: "Use when array is sorted or need to compare elements",
      problems: ["Two Sum II", "Container With Most Water", "3Sum"],
      icon: <ArrowRight className="h-5 w-5" />,
      color: "from-blue-500 to-cyan-500"
    },
    {
      id: "sliding-window",
      name: "Sliding Window",
      description: "Contiguous subarray/substring problems",
      problems: ["Longest Substring", "Max Sum Subarray", "Min Window"],
      icon: <Layers className="h-5 w-5" />,
      color: "from-purple-500 to-pink-500"
    },
    {
      id: "binary-search",
      name: "Binary Search",
      description: "Sorted array or searchable answer space",
      problems: ["Search in Rotated Array", "Find Peak", "Koko Bananas"],
      icon: <Binary className="h-5 w-5" />,
      color: "from-green-500 to-emerald-500"
    },
    {
      id: "dfs-bfs",
      name: "DFS / BFS",
      description: "Tree/graph traversal, shortest path",
      problems: ["Number of Islands", "Level Order", "Clone Graph"],
      icon: <GitBranch className="h-5 w-5" />,
      color: "from-orange-500 to-red-500"
    },
    {
      id: "heap",
      name: "Heap / Priority Queue",
      description: "Top K elements, scheduling, merging",
      problems: ["Kth Largest", "Merge K Lists", "Task Scheduler"],
      icon: <Database className="h-5 w-5" />,
      color: "from-yellow-500 to-orange-500"
    },
    {
      id: "dp",
      name: "Dynamic Programming",
      description: "Overlapping subproblems, optimal substructure",
      problems: ["Coin Change", "LCS", "House Robber"],
      icon: <LayoutGrid className="h-5 w-5" />,
      color: "from-indigo-500 to-purple-500"
    },
  ];

  const complexityTable = [
    { structure: "Array", access: "O(1)", search: "O(n)", insert: "O(n)", delete: "O(n)" },
    { structure: "Linked List", access: "O(n)", search: "O(n)", insert: "O(1)", delete: "O(1)" },
    { structure: "Hash Table", access: "N/A", search: "O(1)", insert: "O(1)", delete: "O(1)" },
    { structure: "BST", access: "O(log n)", search: "O(log n)", insert: "O(log n)", delete: "O(log n)" },
    { structure: "Heap", access: "O(1)", search: "O(n)", insert: "O(log n)", delete: "O(log n)" },
  ];

  const userStats = {
    topicsCompleted: topics.filter(t => t.completed).length,
    totalTopics: topics.length,
    problemsSolved: 45,
    totalProblems: 177,
    studyHours: 12,
    currentStreak: 3,
  };

  const getDifficultyColor = (difficulty: string) => {
    switch (difficulty) {
      case "Beginner": return "bg-green-100 text-green-700 border-green-200";
      case "Intermediate": return "bg-yellow-100 text-yellow-700 border-yellow-200";
      case "Advanced": return "bg-red-100 text-red-700 border-red-200";
      default: return "bg-gray-100 text-gray-700 border-gray-200";
    }
  };

  return (
    <div className="min-h-screen">
      {/* Hero Section */}
      <section className="relative overflow-hidden animated-gradient">
        <div className="absolute inset-0 overflow-hidden pointer-events-none">
          <div className="absolute -top-40 -right-40 w-96 h-96 bg-gradient-to-br from-primary/20 to-blue-500/20 rounded-full blur-3xl float-animation" />
          <div className="absolute -bottom-40 -left-40 w-96 h-96 bg-gradient-to-tr from-cyan-500/20 to-primary/20 rounded-full blur-3xl float-animation" style={{ animationDelay: '1s' }} />
        </div>

        <div className="container max-w-7xl mx-auto px-4 py-12 sm:py-16 relative z-10">
          <div className="flex flex-col lg:flex-row lg:items-center lg:justify-between mb-8">
            <div>
              <Badge variant="secondary" className="mb-4 text-sm">
                <Brain className="mr-2 h-4 w-4" />
                Interview Preparation
              </Badge>
              <h1 className="text-3xl lg:text-4xl xl:text-5xl font-bold tracking-tight mb-3">
                <span className="go-gradient-text">Master Coding Interviews</span>
              </h1>
              <p className="text-base lg:text-xl text-muted-foreground max-w-2xl">
                Comprehensive guide to data structures, algorithms, and system design
              </p>
            </div>
            <div className="mt-6 lg:mt-0 flex gap-4">
              <div className="text-center p-4 rounded-xl bg-card/50 backdrop-blur border">
                <div className="text-2xl font-bold go-gradient-text">{userStats.problemsSolved}</div>
                <div className="text-sm text-muted-foreground">Solved</div>
              </div>
              <div className="text-center p-4 rounded-xl bg-card/50 backdrop-blur border">
                <div className="text-2xl font-bold go-gradient-text">{userStats.studyHours}h</div>
                <div className="text-sm text-muted-foreground">Study Time</div>
              </div>
            </div>
          </div>

          {/* Quick Stats */}
          <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
            <div className="p-4 rounded-xl bg-card/50 backdrop-blur border hover:shadow-lg transition-all">
              <div className="flex items-center gap-3">
                <div className="p-2 rounded-lg bg-blue-500/10">
                  <Target className="h-5 w-5 text-blue-500" />
                </div>
                <div>
                  <div className="text-2xl font-bold">{userStats.topicsCompleted}/{userStats.totalTopics}</div>
                  <div className="text-sm text-muted-foreground">Topics</div>
                </div>
              </div>
            </div>
            <div className="p-4 rounded-xl bg-card/50 backdrop-blur border hover:shadow-lg transition-all">
              <div className="flex items-center gap-3">
                <div className="p-2 rounded-lg bg-green-500/10">
                  <CheckCircle className="h-5 w-5 text-green-500" />
                </div>
                <div>
                  <div className="text-2xl font-bold">{userStats.problemsSolved}</div>
                  <div className="text-sm text-muted-foreground">Problems</div>
                </div>
              </div>
            </div>
            <div className="p-4 rounded-xl bg-card/50 backdrop-blur border hover:shadow-lg transition-all">
              <div className="flex items-center gap-3">
                <div className="p-2 rounded-lg bg-orange-500/10">
                  <Zap className="h-5 w-5 text-orange-500" />
                </div>
                <div>
                  <div className="text-2xl font-bold">{userStats.currentStreak}</div>
                  <div className="text-sm text-muted-foreground">Day Streak</div>
                </div>
              </div>
            </div>
            <div className="p-4 rounded-xl bg-card/50 backdrop-blur border hover:shadow-lg transition-all">
              <div className="flex items-center gap-3">
                <div className="p-2 rounded-lg bg-purple-500/10">
                  <Trophy className="h-5 w-5 text-purple-500" />
                </div>
                <div>
                  <div className="text-2xl font-bold">{Math.round((userStats.problemsSolved / userStats.totalProblems) * 100)}%</div>
                  <div className="text-sm text-muted-foreground">Progress</div>
                </div>
              </div>
            </div>
          </div>

          <Progress
            value={(userStats.topicsCompleted / userStats.totalTopics) * 100}
            className="h-2 mt-6"
          />
        </div>
      </section>

      {/* AI Mock Interviews Section */}
      <section className="container max-w-7xl mx-auto px-4 py-8">
        <div className="mb-8">
          <div className="flex items-center gap-3 mb-2">
            <div className="p-2 rounded-lg bg-gradient-to-br from-purple-500 to-blue-600 text-white">
              <Brain className="h-6 w-6" />
            </div>
            <div>
              <h2 className="text-2xl font-bold">AI Mock Interviews</h2>
              <p className="text-muted-foreground">Practice with AI-powered interview simulations</p>
            </div>
          </div>
        </div>

        <div className="grid grid-cols-1 md:grid-cols-3 gap-6 mb-8">
          {/* Coding Interview */}
          <Card className="relative overflow-hidden group hover:shadow-xl transition-all">
            <div className="absolute top-0 left-0 right-0 h-1 bg-gradient-to-r from-blue-500 to-cyan-500" />
            <CardHeader>
              <div className="flex items-center gap-3">
                <div className="p-3 rounded-xl bg-blue-500/10">
                  <Code2 className="h-8 w-8 text-blue-500" />
                </div>
                <div>
                  <CardTitle>Coding Interview</CardTitle>
                  <CardDescription>Algorithm & data structure problems</CardDescription>
                </div>
              </div>
            </CardHeader>
            <CardContent className="space-y-4">
              <p className="text-sm text-muted-foreground">
                Solve coding challenges with real-time AI feedback on your solutions.
              </p>
              <div className="flex flex-wrap gap-2 mb-4">
                <Badge variant="secondary">Arrays</Badge>
                <Badge variant="secondary">Trees</Badge>
                <Badge variant="secondary">DP</Badge>
                <Badge variant="secondary">Graphs</Badge>
              </div>
              <div className="flex gap-2">
                <Link href="/interviews/session?type=coding&difficulty=beginner" className="flex-1">
                  <Button variant="outline" size="sm" className="w-full">Beginner</Button>
                </Link>
                <Link href="/interviews/session?type=coding&difficulty=intermediate" className="flex-1">
                  <Button variant="outline" size="sm" className="w-full">Intermediate</Button>
                </Link>
                <Link href="/interviews/session?type=coding&difficulty=advanced" className="flex-1">
                  <Button variant="outline" size="sm" className="w-full">Advanced</Button>
                </Link>
              </div>
              <Link href="/interviews/practice" className="block">
                <Button className="w-full group-hover:bg-blue-600 transition-colors">
                  <Play className="h-4 w-4 mr-2" />
                  Start Practice
                </Button>
              </Link>
            </CardContent>
          </Card>

          {/* Behavioral Interview */}
          <Card className="relative overflow-hidden group hover:shadow-xl transition-all">
            <div className="absolute top-0 left-0 right-0 h-1 bg-gradient-to-r from-purple-500 to-pink-500" />
            <CardHeader>
              <div className="flex items-center gap-3">
                <div className="p-3 rounded-xl bg-purple-500/10">
                  <MessageSquare className="h-8 w-8 text-purple-500" />
                </div>
                <div>
                  <CardTitle>Behavioral Interview</CardTitle>
                  <CardDescription>STAR method & soft skills</CardDescription>
                </div>
              </div>
            </CardHeader>
            <CardContent className="space-y-4">
              <p className="text-sm text-muted-foreground">
                Practice behavioral questions with AI analyzing your STAR responses.
              </p>
              <div className="flex flex-wrap gap-2 mb-4">
                <Badge variant="secondary">Leadership</Badge>
                <Badge variant="secondary">Conflict</Badge>
                <Badge variant="secondary">Growth</Badge>
                <Badge variant="secondary">Teamwork</Badge>
              </div>
              <div className="flex gap-2">
                <Link href="/interviews/session?type=behavioral&difficulty=beginner" className="flex-1">
                  <Button variant="outline" size="sm" className="w-full">Beginner</Button>
                </Link>
                <Link href="/interviews/session?type=behavioral&difficulty=intermediate" className="flex-1">
                  <Button variant="outline" size="sm" className="w-full">Intermediate</Button>
                </Link>
                <Link href="/interviews/session?type=behavioral&difficulty=advanced" className="flex-1">
                  <Button variant="outline" size="sm" className="w-full">Advanced</Button>
                </Link>
              </div>
              <Link href="/interviews/practice" className="block">
                <Button className="w-full group-hover:bg-purple-600 transition-colors">
                  <Play className="h-4 w-4 mr-2" />
                  Start Practice
                </Button>
              </Link>
            </CardContent>
          </Card>

          {/* System Design Interview */}
          <Card className="relative overflow-hidden group hover:shadow-xl transition-all">
            <div className="absolute top-0 left-0 right-0 h-1 bg-gradient-to-r from-orange-500 to-red-500" />
            <CardHeader>
              <div className="flex items-center gap-3">
                <div className="p-3 rounded-xl bg-orange-500/10">
                  <Layers className="h-8 w-8 text-orange-500" />
                </div>
                <div>
                  <CardTitle>System Design</CardTitle>
                  <CardDescription>Architecture & scalability</CardDescription>
                </div>
              </div>
            </CardHeader>
            <CardContent className="space-y-4">
              <p className="text-sm text-muted-foreground">
                Design scalable systems with AI feedback on trade-offs and architecture.
              </p>
              <div className="flex flex-wrap gap-2 mb-4">
                <Badge variant="secondary">Scalability</Badge>
                <Badge variant="secondary">Distributed</Badge>
                <Badge variant="secondary">Database</Badge>
                <Badge variant="secondary">Caching</Badge>
              </div>
              <div className="flex gap-2">
                <Link href="/interviews/session?type=system_design&difficulty=beginner" className="flex-1">
                  <Button variant="outline" size="sm" className="w-full">Beginner</Button>
                </Link>
                <Link href="/interviews/session?type=system_design&difficulty=intermediate" className="flex-1">
                  <Button variant="outline" size="sm" className="w-full">Intermediate</Button>
                </Link>
                <Link href="/interviews/session?type=system_design&difficulty=advanced" className="flex-1">
                  <Button variant="outline" size="sm" className="w-full">Advanced</Button>
                </Link>
              </div>
              <Link href="/interviews/practice" className="block">
                <Button className="w-full group-hover:bg-orange-600 transition-colors">
                  <Play className="h-4 w-4 mr-2" />
                  Start Practice
                </Button>
              </Link>
            </CardContent>
          </Card>
        </div>

        {/* Quick Stats & Recent Sessions */}
        <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
          <Card className="lg:col-span-2">
            <CardHeader>
              <CardTitle className="flex items-center gap-2">
                <Trophy className="h-5 w-5 text-yellow-500" />
                Your Progress
              </CardTitle>
            </CardHeader>
            <CardContent>
              <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
                <div className="text-center p-4 rounded-lg bg-muted/50">
                  <div className="text-3xl font-bold text-blue-600">12</div>
                  <div className="text-sm text-muted-foreground">Sessions</div>
                </div>
                <div className="text-center p-4 rounded-lg bg-muted/50">
                  <div className="text-3xl font-bold text-green-600">78%</div>
                  <div className="text-sm text-muted-foreground">Avg Score</div>
                </div>
                <div className="text-center p-4 rounded-lg bg-muted/50">
                  <div className="text-3xl font-bold text-purple-600">5</div>
                  <div className="text-sm text-muted-foreground">This Week</div>
                </div>
                <div className="text-center p-4 rounded-lg bg-muted/50">
                  <div className="text-3xl font-bold text-orange-600">4.2h</div>
                  <div className="text-sm text-muted-foreground">Practice Time</div>
                </div>
              </div>
            </CardContent>
          </Card>

          <Card>
            <CardHeader>
              <CardTitle className="flex items-center gap-2">
                <Star className="h-5 w-5 text-yellow-500" />
                Quick Actions
              </CardTitle>
            </CardHeader>
            <CardContent className="space-y-3">
              <Link href="/interviews/practice">
                <Button variant="outline" className="w-full justify-start">
                  <Play className="h-4 w-4 mr-2" />
                  Quick Practice
                </Button>
              </Link>
              <Link href="/interviews/feedback">
                <Button variant="outline" className="w-full justify-start">
                  <TrendingUp className="h-4 w-4 mr-2" />
                  View Feedback
                </Button>
              </Link>
              <Link href="#roadmap">
                <Button variant="outline" className="w-full justify-start">
                  <BookOpen className="h-4 w-4 mr-2" />
                  Study Roadmap
                </Button>
              </Link>
            </CardContent>
          </Card>
        </div>
      </section>

      {/* Main Content */}
      <div className="container max-w-7xl mx-auto px-4 py-8">
        <Tabs value={activeTab} onValueChange={setActiveTab} className="space-y-6">
          <TabsList className="grid w-full grid-cols-6 lg:w-[700px]">
            <TabsTrigger value="roadmap">Roadmap</TabsTrigger>
            <TabsTrigger value="system-design">System Design</TabsTrigger>
            <TabsTrigger value="patterns">Patterns</TabsTrigger>
            <TabsTrigger value="reference">Reference</TabsTrigger>
            <TabsTrigger value="behavioral">Behavioral</TabsTrigger>
            <TabsTrigger value="tips">Tips</TabsTrigger>
          </TabsList>

          {/* Roadmap Tab */}
          <TabsContent value="roadmap" className="space-y-6">
            <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
              {["Beginner", "Intermediate", "Advanced"].map((level) => (
                <div key={level} className="space-y-4">
                  <div className="flex items-center gap-2">
                    <Badge className={getDifficultyColor(level)}>{level}</Badge>
                    <span className="text-sm text-muted-foreground">
                      {topics.filter(t => t.difficulty === level).length} topics
                    </span>
                  </div>
                  {topics
                    .filter(t => t.difficulty === level)
                    .map((topic) => (
                      <Card key={topic.id} className={`cursor-pointer transition-all hover:shadow-lg ${topic.completed ? 'border-green-500/50' : ''}`}>
                        <CardHeader className="pb-2">
                          <div className="flex items-start justify-between">
                            <div className="flex items-center gap-2">
                              {topic.completed && <CheckCircle className="h-5 w-5 text-green-500" />}
                              <CardTitle className="text-lg">{topic.title}</CardTitle>
                            </div>
                            <Badge variant="outline">{topic.category}</Badge>
                          </div>
                          <CardDescription>{topic.description}</CardDescription>
                        </CardHeader>
                        <CardContent>
                          <div className="flex items-center justify-between text-sm text-muted-foreground mb-3">
                            <div className="flex items-center gap-4">
                              <span className="flex items-center gap-1">
                                <Clock className="h-4 w-4" />
                                {topic.estimatedTime}
                              </span>
                              <span className="flex items-center gap-1">
                                <Code2 className="h-4 w-4" />
                                {topic.problems} problems
                              </span>
                            </div>
                          </div>
                          <div className="p-2 rounded bg-muted/50 text-xs font-mono">
                            Key: {topic.key}
                          </div>
                        </CardContent>
                      </Card>
                    ))}
                </div>
              ))}
            </div>
          </TabsContent>

          {/* System Design Tab */}
          <TabsContent value="system-design" className="space-y-6">
            <div className="flex items-center justify-between mb-4">
              <div>
                <h3 className="text-xl font-semibold">System Design Fundamentals</h3>
                <p className="text-muted-foreground">Master distributed systems architecture concepts</p>
              </div>
              <Link href="/interviews/system-design">
                <Button>
                  View Full Guide
                  <ChevronRight className="h-4 w-4 ml-2" />
                </Button>
              </Link>
            </div>

            {/* System Design Categories */}
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
              {[
                {
                  id: 'fundamentals',
                  title: 'Scalability & Load Balancing',
                  description: 'Horizontal vs vertical scaling, load balancer algorithms',
                  icon: <Scale className="h-5 w-5" />,
                  color: 'from-blue-500 to-cyan-500',
                  topics: ['Scalability', 'Load Balancing', 'High Availability'],
                },
                {
                  id: 'data',
                  title: 'Data Layer',
                  description: 'Database sharding, replication, caching strategies',
                  icon: <Database className="h-5 w-5" />,
                  color: 'from-green-500 to-emerald-500',
                  topics: ['Database Sharding', 'Caching', 'CAP Theorem'],
                },
                {
                  id: 'communication',
                  title: 'Communication Patterns',
                  description: 'API design, message queues, microservices',
                  icon: <ArrowRightLeft className="h-5 w-5" />,
                  color: 'from-purple-500 to-pink-500',
                  topics: ['REST/GraphQL/gRPC', 'Message Queues', 'Microservices'],
                },
                {
                  id: 'patterns',
                  title: 'Design Patterns',
                  description: 'CDN, rate limiting, distributed ID generation',
                  icon: <Layers className="h-5 w-5" />,
                  color: 'from-orange-500 to-red-500',
                  topics: ['CDN', 'Rate Limiting', 'ID Generation'],
                },
                {
                  id: 'scalability',
                  title: 'Advanced Scalability',
                  description: 'Consistent hashing, event sourcing, CQRS',
                  icon: <RefreshCw className="h-5 w-5" />,
                  color: 'from-indigo-500 to-violet-500',
                  topics: ['Consistent Hashing', 'Event Sourcing', 'CQRS'],
                },
                {
                  id: 'practice',
                  title: 'Practice Problems',
                  description: 'Design URL shortener, chat system, and more',
                  icon: <Target className="h-5 w-5" />,
                  color: 'from-yellow-500 to-orange-500',
                  topics: ['URL Shortener', 'Chat System', 'Rate Limiter'],
                },
              ].map((category) => (
                <Card key={category.id} className="overflow-hidden hover:shadow-lg transition-all cursor-pointer group">
                  <div className={`h-2 bg-gradient-to-r ${category.color}`} />
                  <CardHeader className="pb-2">
                    <div className="flex items-center gap-3">
                      <div className={`p-2 rounded-lg bg-gradient-to-r ${category.color} text-white`}>
                        {category.icon}
                      </div>
                      <CardTitle className="text-lg">{category.title}</CardTitle>
                    </div>
                    <CardDescription>{category.description}</CardDescription>
                  </CardHeader>
                  <CardContent>
                    <div className="flex flex-wrap gap-2">
                      {category.topics.map((topic) => (
                        <Badge key={topic} variant="secondary" className="text-xs">
                          {topic}
                        </Badge>
                      ))}
                    </div>
                  </CardContent>
                </Card>
              ))}
            </div>

            {/* Quick Reference */}
            <Card className="border-0 shadow-lg bg-gradient-to-r from-indigo-50 to-purple-50">
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <Shield className="h-5 w-5 text-indigo-500" />
                  System Design Interview Framework
                </CardTitle>
              </CardHeader>
              <CardContent>
                <div className="grid grid-cols-1 md:grid-cols-4 gap-6">
                  {[
                    {
                      step: '1',
                      title: 'Requirements',
                      items: ['Functional', 'Non-functional', 'Capacity', 'Constraints'],
                    },
                    {
                      step: '2',
                      title: 'High-Level',
                      items: ['APIs', 'Data Model', 'Components', 'Data Flow'],
                    },
                    {
                      step: '3',
                      title: 'Deep Dive',
                      items: ['Schema', 'Scaling', 'Caching', 'Load Balancing'],
                    },
                    {
                      step: '4',
                      title: 'Trade-offs',
                      items: ['C vs A', 'Latency', 'Cost', 'Complexity'],
                    },
                  ].map((section) => (
                    <div key={section.step}>
                      <div className="flex items-center gap-2 mb-2">
                        <div className="w-6 h-6 rounded-full bg-primary text-primary-foreground flex items-center justify-center text-sm font-bold">
                          {section.step}
                        </div>
                        <span className="font-semibold">{section.title}</span>
                      </div>
                      <ul className="text-sm text-muted-foreground space-y-1">
                        {section.items.map((item) => (
                          <li key={item}>• {item}</li>
                        ))}
                      </ul>
                    </div>
                  ))}
                </div>
              </CardContent>
            </Card>

            {/* Start Practice CTA */}
            <div className="flex gap-4">
              <Link href="/interviews/system-design">
                <Button size="lg">
                  <BookOpen className="h-5 w-5 mr-2" />
                  Full System Design Guide
                </Button>
              </Link>
              <Link href="/interviews/session?type=system_design&difficulty=beginner">
                <Button size="lg" variant="outline">
                  <Play className="h-5 w-5 mr-2" />
                  Practice System Design
                </Button>
              </Link>
            </div>
          </TabsContent>

          {/* Patterns Tab */}
          <TabsContent value="patterns" className="space-y-6">
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
              {patterns.map((pattern) => (
                <Card key={pattern.id} className="overflow-hidden hover:shadow-lg transition-all">
                  <div className={`h-2 bg-gradient-to-r ${pattern.color}`} />
                  <CardHeader>
                    <div className="flex items-center gap-3">
                      <div className={`p-2 rounded-lg bg-gradient-to-r ${pattern.color} text-white`}>
                        {pattern.icon}
                      </div>
                      <CardTitle className="text-lg">{pattern.name}</CardTitle>
                    </div>
                    <CardDescription>{pattern.description}</CardDescription>
                  </CardHeader>
                  <CardContent>
                    <div className="space-y-2">
                      <p className="text-sm font-medium">Example Problems:</p>
                      <div className="flex flex-wrap gap-2">
                        {pattern.problems.map((problem) => (
                          <Badge key={problem} variant="secondary" className="text-xs">
                            {problem}
                          </Badge>
                        ))}
                      </div>
                    </div>
                  </CardContent>
                </Card>
              ))}
            </div>

            {/* Pattern Recognition Guide */}
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <Search className="h-5 w-5" />
                  Pattern Recognition Guide
                </CardTitle>
                <CardDescription>Match problem keywords to the right approach</CardDescription>
              </CardHeader>
              <CardContent>
                <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
                  {[
                    { keyword: "contiguous subarray", pattern: "Sliding Window" },
                    { keyword: "sorted array", pattern: "Binary Search" },
                    { keyword: "all combinations", pattern: "Backtracking" },
                    { keyword: "shortest path", pattern: "BFS" },
                    { keyword: "optimal", pattern: "Dynamic Programming" },
                    { keyword: "top K", pattern: "Heap" },
                    { keyword: "prefix", pattern: "Trie" },
                    { keyword: "cycle", pattern: "Fast/Slow Pointer" },
                    { keyword: "connected", pattern: "Union-Find" },
                  ].map((item) => (
                    <div key={item.keyword} className="flex items-center gap-2 p-2 rounded bg-muted/50">
                      <code className="text-xs bg-background px-2 py-1 rounded">"{item.keyword}"</code>
                      <ChevronRight className="h-4 w-4 text-muted-foreground" />
                      <span className="text-sm font-medium">{item.pattern}</span>
                    </div>
                  ))}
                </div>
              </CardContent>
            </Card>
          </TabsContent>

          {/* Reference Tab */}
          <TabsContent value="reference" className="space-y-6">
            {/* Complexity Table */}
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <Database className="h-5 w-5" />
                  Time Complexity Reference
                </CardTitle>
                <CardDescription>Quick lookup for data structure operations</CardDescription>
              </CardHeader>
              <CardContent>
                <div className="overflow-x-auto">
                  <table className="w-full text-sm">
                    <thead>
                      <tr className="border-b">
                        <th className="text-left p-3 font-medium">Data Structure</th>
                        <th className="text-center p-3 font-medium">Access</th>
                        <th className="text-center p-3 font-medium">Search</th>
                        <th className="text-center p-3 font-medium">Insert</th>
                        <th className="text-center p-3 font-medium">Delete</th>
                      </tr>
                    </thead>
                    <tbody>
                      {complexityTable.map((row) => (
                        <tr key={row.structure} className="border-b hover:bg-muted/50">
                          <td className="p-3 font-medium">{row.structure}</td>
                          <td className="p-3 text-center font-mono text-blue-600">{row.access}</td>
                          <td className="p-3 text-center font-mono text-green-600">{row.search}</td>
                          <td className="p-3 text-center font-mono text-orange-600">{row.insert}</td>
                          <td className="p-3 text-center font-mono text-red-600">{row.delete}</td>
                        </tr>
                      ))}
                    </tbody>
                  </table>
                </div>
              </CardContent>
            </Card>

            {/* Big-O Cheatsheet */}
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <TrendingUp className="h-5 w-5" />
                  Big-O Complexity Hierarchy
                </CardTitle>
              </CardHeader>
              <CardContent>
                <div className="grid grid-cols-2 md:grid-cols-4 lg:grid-cols-8 gap-3">
                  {[
                    { complexity: "O(1)", label: "Constant", color: "bg-green-100 text-green-700" },
                    { complexity: "O(log n)", label: "Logarithmic", color: "bg-green-100 text-green-700" },
                    { complexity: "O(n)", label: "Linear", color: "bg-yellow-100 text-yellow-700" },
                    { complexity: "O(n log n)", label: "Linearithmic", color: "bg-yellow-100 text-yellow-700" },
                    { complexity: "O(n²)", label: "Quadratic", color: "bg-orange-100 text-orange-700" },
                    { complexity: "O(n³)", label: "Cubic", color: "bg-orange-100 text-orange-700" },
                    { complexity: "O(2ⁿ)", label: "Exponential", color: "bg-red-100 text-red-700" },
                    { complexity: "O(n!)", label: "Factorial", color: "bg-red-100 text-red-700" },
                  ].map((item) => (
                    <div key={item.complexity} className={`p-3 rounded-lg text-center ${item.color}`}>
                      <div className="font-mono font-bold">{item.complexity}</div>
                      <div className="text-xs mt-1">{item.label}</div>
                    </div>
                  ))}
                </div>
              </CardContent>
            </Card>

            {/* Input Size Guide */}
            <Card>
              <CardHeader>
                <CardTitle>Input Size → Algorithm Guide</CardTitle>
                <CardDescription>Choose algorithm based on constraints</CardDescription>
              </CardHeader>
              <CardContent>
                <div className="space-y-3">
                  {[
                    { size: "n ≤ 20", approach: "Backtracking, recursion, brute force" },
                    { size: "n ≤ 100", approach: "O(n²) is acceptable" },
                    { size: "n ≤ 1,000", approach: "O(n²) might work, prefer O(n log n)" },
                    { size: "n ≤ 100,000", approach: "O(n log n) or O(n) required" },
                    { size: "n ≤ 1,000,000", approach: "O(n) only, must be linear" },
                    { size: "n > 1,000,000", approach: "O(log n) or O(1), use math/properties" },
                  ].map((item) => (
                    <div key={item.size} className="flex items-center gap-4 p-3 rounded-lg bg-muted/50">
                      <code className="font-mono text-sm bg-background px-3 py-1 rounded">{item.size}</code>
                      <ChevronRight className="h-4 w-4 text-muted-foreground" />
                      <span className="text-sm">{item.approach}</span>
                    </div>
                  ))}
                </div>
              </CardContent>
            </Card>
          </TabsContent>

          {/* Behavioral Tab */}
          <TabsContent value="behavioral" className="space-y-6">
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <MessageSquare className="h-5 w-5" />
                  The STAR Method
                </CardTitle>
                <CardDescription>Structure your behavioral responses</CardDescription>
              </CardHeader>
              <CardContent>
                <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
                  {[
                    { letter: "S", word: "Situation", desc: "Set the context" },
                    { letter: "T", word: "Task", desc: "What you needed to do" },
                    { letter: "A", word: "Action", desc: "What you actually did" },
                    { letter: "R", word: "Result", desc: "The outcome (quantify!)" },
                  ].map((item) => (
                    <div key={item.letter} className="p-4 rounded-lg border bg-card">
                      <div className="flex items-center gap-3 mb-2">
                        <div className="w-10 h-10 rounded-full bg-primary text-primary-foreground flex items-center justify-center font-bold text-lg">
                          {item.letter}
                        </div>
                        <span className="font-semibold">{item.word}</span>
                      </div>
                      <p className="text-sm text-muted-foreground">{item.desc}</p>
                    </div>
                  ))}
                </div>
              </CardContent>
            </Card>

            <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
              <Card>
                <CardHeader>
                  <CardTitle className="flex items-center gap-2">
                    <Users className="h-5 w-5" />
                    Common Questions
                  </CardTitle>
                </CardHeader>
                <CardContent>
                  <div className="space-y-3">
                    {[
                      "Tell me about a challenging bug you fixed",
                      "Describe a time you had a conflict with a teammate",
                      "What's your biggest professional failure?",
                      "How do you handle tight deadlines?",
                      "Tell me about a time you led a project",
                    ].map((question, i) => (
                      <div key={i} className="p-3 rounded-lg bg-muted/50 text-sm">
                        {question}
                      </div>
                    ))}
                  </div>
                </CardContent>
              </Card>

              <Card>
                <CardHeader>
                  <CardTitle className="flex items-center gap-2">
                    <Lightbulb className="h-5 w-5" />
                    Example STAR Response
                  </CardTitle>
                </CardHeader>
                <CardContent>
                  <div className="space-y-4 text-sm">
                    <div className="p-3 rounded-lg border-l-4 border-blue-500 bg-blue-50/50">
                      <span className="font-semibold text-blue-700">Situation:</span>
                      <p className="text-muted-foreground mt-1">Our payment service was timing out under load</p>
                    </div>
                    <div className="p-3 rounded-lg border-l-4 border-green-500 bg-green-50/50">
                      <span className="font-semibold text-green-700">Task:</span>
                      <p className="text-muted-foreground mt-1">I needed to reduce latency by 50%</p>
                    </div>
                    <div className="p-3 rounded-lg border-l-4 border-orange-500 bg-orange-50/50">
                      <span className="font-semibold text-orange-700">Action:</span>
                      <p className="text-muted-foreground mt-1">Profiled code, found N+1 queries, implemented caching</p>
                    </div>
                    <div className="p-3 rounded-lg border-l-4 border-purple-500 bg-purple-50/50">
                      <span className="font-semibold text-purple-700">Result:</span>
                      <p className="text-muted-foreground mt-1">Latency dropped 70%, saved $50K/month</p>
                    </div>
                  </div>
                </CardContent>
              </Card>
            </div>
          </TabsContent>

          {/* Tips Tab */}
          <TabsContent value="tips" className="space-y-6">
            <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
              <Card>
                <CardHeader>
                  <CardTitle className="flex items-center gap-2">
                    <Target className="h-5 w-5" />
                    The REACTO Method
                  </CardTitle>
                  <CardDescription>Step-by-step problem solving</CardDescription>
                </CardHeader>
                <CardContent>
                  <div className="space-y-4">
                    {[
                      { step: "R", name: "Repeat", desc: "Restate the problem in your own words" },
                      { step: "E", name: "Examples", desc: "Work through examples, find edge cases" },
                      { step: "A", name: "Approach", desc: "Describe your approach before coding" },
                      { step: "C", name: "Code", desc: "Write clean, readable code" },
                      { step: "T", name: "Test", desc: "Walk through with test cases" },
                      { step: "O", name: "Optimize", desc: "Discuss possible improvements" },
                    ].map((item, i) => (
                      <div key={item.step} className="flex items-start gap-4">
                        <div className="w-8 h-8 rounded-full bg-primary text-primary-foreground flex items-center justify-center font-bold shrink-0">
                          {item.step}
                        </div>
                        <div>
                          <span className="font-semibold">{item.name}</span>
                          <p className="text-sm text-muted-foreground">{item.desc}</p>
                        </div>
                      </div>
                    ))}
                  </div>
                </CardContent>
              </Card>

              <Card>
                <CardHeader>
                  <CardTitle className="flex items-center gap-2">
                    <Zap className="h-5 w-5" />
                    Interview Day Tips
                  </CardTitle>
                </CardHeader>
                <CardContent>
                  <div className="space-y-4">
                    <div className="p-4 rounded-lg bg-green-50 border border-green-200">
                      <h4 className="font-semibold text-green-800 mb-2">✅ Do</h4>
                      <ul className="text-sm text-green-700 space-y-1">
                        <li>• Think out loud - silence is your enemy</li>
                        <li>• Ask clarifying questions first</li>
                        <li>• Start with brute force, then optimize</li>
                        <li>• Test your code with examples</li>
                        <li>• Be honest if you don't know something</li>
                      </ul>
                    </div>
                    <div className="p-4 rounded-lg bg-red-50 border border-red-200">
                      <h4 className="font-semibold text-red-800 mb-2">❌ Don't</h4>
                      <ul className="text-sm text-red-700 space-y-1">
                        <li>• Jump straight into coding</li>
                        <li>• Stay silent while thinking</li>
                        <li>• Ignore edge cases</li>
                        <li>• Give up if stuck - ask for hints</li>
                        <li>• Memorize solutions blindly</li>
                      </ul>
                    </div>
                  </div>
                </CardContent>
              </Card>
            </div>

            {/* Communication Templates */}
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <MessageSquare className="h-5 w-5" />
                  Communication Templates
                </CardTitle>
                <CardDescription>Phrases to use during interviews</CardDescription>
              </CardHeader>
              <CardContent>
                <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                  {[
                    { situation: "Understanding", phrase: "Let me make sure I understand. We're given [X] and need to [Y]. Is that correct?" },
                    { situation: "Clarifying", phrase: "Should I handle the case where [edge case]?" },
                    { situation: "Approaching", phrase: "Let me start with a brute force approach, then optimize it." },
                    { situation: "Coding", phrase: "I'll use a hash map here for O(1) lookups." },
                    { situation: "Testing", phrase: "Let me trace through with input [X]..." },
                    { situation: "Stuck", phrase: "I'm thinking about [approach]. Would a hint help point me in the right direction?" },
                  ].map((item) => (
                    <div key={item.situation} className="p-3 rounded-lg bg-muted/50">
                      <Badge variant="outline" className="mb-2">{item.situation}</Badge>
                      <p className="text-sm italic">"{item.phrase}"</p>
                    </div>
                  ))}
                </div>
              </CardContent>
            </Card>

            {/* 8-Week Plan */}
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <BookOpen className="h-5 w-5" />
                  8-Week Study Plan
                </CardTitle>
              </CardHeader>
              <CardContent>
                <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
                  {[
                    { week: "1-2", focus: "Data Structures", topics: "Arrays, Strings, Linked Lists, Stacks" },
                    { week: "3-4", focus: "Algorithms", topics: "Sorting, Trees, Graphs, DP Basics" },
                    { week: "5-6", focus: "Patterns", topics: "Two Pointers, Sliding Window, Backtracking" },
                    { week: "7-8", focus: "Final Prep", topics: "System Design, Mock Interviews, Review" },
                  ].map((item) => (
                    <div key={item.week} className="p-4 rounded-lg border bg-card">
                      <Badge className="mb-2">Week {item.week}</Badge>
                      <h4 className="font-semibold mb-1">{item.focus}</h4>
                      <p className="text-sm text-muted-foreground">{item.topics}</p>
                    </div>
                  ))}
                </div>
              </CardContent>
            </Card>
          </TabsContent>
        </Tabs>
      </div>
    </div>
  );
}
