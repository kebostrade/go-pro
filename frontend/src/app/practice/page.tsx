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
  Target,
  Clock,
  Star,
  Zap,
  Brain,
  CheckCircle,
  Lock,
  Play,
  ArrowRight,
  Filter,
  Search,
  TrendingUp,
  Award,
  Users,
  BookOpen,
  Lightbulb
} from "lucide-react";
import Link from "next/link";

interface Challenge {
  id: string;
  title: string;
  description: string;
  difficulty: "Beginner" | "Intermediate" | "Advanced";
  category: string;
  estimatedTime: string;
  points: number;
  completed: boolean;
  locked: boolean;
  tags: string[];
  completionRate: number;
  attempts: number;
}

interface SkillAssessment {
  id: string;
  title: string;
  description: string;
  category: string;
  questions: number;
  duration: string;
  difficulty: "Beginner" | "Intermediate" | "Advanced";
  completed: boolean;
  score?: number;
  maxScore: number;
}

export default function PracticePage() {
  const [activeTab, setActiveTab] = useState("challenges");
  const [selectedDifficulty, setSelectedDifficulty] = useState<string>("all");
  const [selectedCategory, setSelectedCategory] = useState<string>("all");
  const [searchQuery, setSearchQuery] = useState("");

  // Mock data for challenges
  const challenges: Challenge[] = [
    {
      id: "go-basics-1",
      title: "Variables and Constants",
      description: "Practice declaring and using variables and constants in Go",
      difficulty: "Beginner",
      category: "Fundamentals",
      estimatedTime: "15 min",
      points: 50,
      completed: true,
      locked: false,
      tags: ["variables", "constants", "types"],
      completionRate: 85,
      attempts: 1
    },
    {
      id: "go-basics-2",
      title: "Control Flow Mastery",
      description: "Master if statements, loops, and switch cases",
      difficulty: "Beginner",
      category: "Fundamentals",
      estimatedTime: "20 min",
      points: 75,
      completed: true,
      locked: false,
      tags: ["if", "for", "switch"],
      completionRate: 78,
      attempts: 2
    },
    {
      id: "data-structures-1",
      title: "Slice Manipulation",
      description: "Work with slices: append, copy, and slice operations",
      difficulty: "Intermediate",
      category: "Data Structures",
      estimatedTime: "25 min",
      points: 100,
      completed: false,
      locked: false,
      tags: ["slices", "arrays", "manipulation"],
      completionRate: 65,
      attempts: 0
    },
    {
      id: "concurrency-1",
      title: "Goroutine Basics",
      description: "Create and manage goroutines for concurrent execution",
      difficulty: "Intermediate",
      category: "Concurrency",
      estimatedTime: "30 min",
      points: 150,
      completed: false,
      locked: false,
      tags: ["goroutines", "concurrency", "channels"],
      completionRate: 45,
      attempts: 0
    },
    {
      id: "advanced-patterns-1",
      title: "Interface Design Patterns",
      description: "Implement advanced interface patterns and polymorphism",
      difficulty: "Advanced",
      category: "Design Patterns",
      estimatedTime: "45 min",
      points: 200,
      completed: false,
      locked: true,
      tags: ["interfaces", "patterns", "polymorphism"],
      completionRate: 25,
      attempts: 0
    }
  ];

  // Mock data for skill assessments
  const skillAssessments: SkillAssessment[] = [
    {
      id: "fundamentals-assessment",
      title: "Go Fundamentals Assessment",
      description: "Test your knowledge of Go basics, syntax, and core concepts",
      category: "Fundamentals",
      questions: 20,
      duration: "30 min",
      difficulty: "Beginner",
      completed: true,
      score: 85,
      maxScore: 100
    },
    {
      id: "concurrency-assessment",
      title: "Concurrency Mastery Test",
      description: "Evaluate your understanding of goroutines, channels, and concurrent patterns",
      category: "Concurrency",
      questions: 15,
      duration: "25 min",
      difficulty: "Intermediate",
      completed: false,
      maxScore: 100
    },
    {
      id: "microservices-assessment",
      title: "Microservices Architecture Quiz",
      description: "Assess your knowledge of building scalable microservices with Go",
      category: "Architecture",
      questions: 25,
      duration: "40 min",
      difficulty: "Advanced",
      completed: false,
      maxScore: 100
    }
  ];

  const categories = ["all", "Fundamentals", "Data Structures", "Concurrency", "Web Development", "Design Patterns", "Architecture"];
  const difficulties = ["all", "Beginner", "Intermediate", "Advanced"];

  const filteredChallenges = challenges.filter(challenge => {
    const matchesDifficulty = selectedDifficulty === "all" || challenge.difficulty === selectedDifficulty;
    const matchesCategory = selectedCategory === "all" || challenge.category === selectedCategory;
    const matchesSearch = challenge.title.toLowerCase().includes(searchQuery.toLowerCase()) ||
                         challenge.description.toLowerCase().includes(searchQuery.toLowerCase()) ||
                         challenge.tags.some(tag => tag.toLowerCase().includes(searchQuery.toLowerCase()));
    
    return matchesDifficulty && matchesCategory && matchesSearch;
  });

  const getDifficultyColor = (difficulty: string) => {
    switch (difficulty) {
      case "Beginner": return "text-green-600 bg-green-50 border-green-200";
      case "Intermediate": return "text-yellow-600 bg-yellow-50 border-yellow-200";
      case "Advanced": return "text-red-600 bg-red-50 border-red-200";
      default: return "text-gray-600 bg-gray-50 border-gray-200";
    }
  };

  const userStats = {
    totalChallenges: challenges.length,
    completedChallenges: challenges.filter(c => c.completed).length,
    totalPoints: challenges.filter(c => c.completed).reduce((sum, c) => sum + c.points, 0),
    currentStreak: 5,
    rank: "Go Apprentice",
    nextRank: "Go Developer"
  };

  return (
    <div className="min-h-screen bg-gradient-to-br from-background via-background to-accent/5">
      <div className="container max-w-7xl mx-auto px-4 py-8 sm:px-6 sm:py-10 lg:px-8 lg:py-12">
        {/* Header */}
        <div className="mb-10 lg:mb-12">
          <div className="flex flex-col lg:flex-row lg:items-center lg:justify-between mb-6 lg:mb-8">
            <div className="mb-4 lg:mb-0">
              <h1 className="text-3xl lg:text-4xl xl:text-5xl font-bold tracking-tight mb-3 bg-gradient-to-r from-primary to-primary/70 bg-clip-text text-transparent">
                Practice Hub
              </h1>
              <p className="text-muted-foreground text-lg lg:text-xl max-w-2xl">
                Sharpen your Go skills with interactive challenges and assessments
              </p>
            </div>
            <div className="text-left lg:text-right">
              <div className="text-3xl lg:text-4xl font-bold text-primary">{userStats.totalPoints}</div>
              <div className="text-sm lg:text-base text-muted-foreground">Total Points</div>
            </div>
          </div>

          {/* Quick Stats */}
          <div className="grid grid-cols-2 md:grid-cols-4 gap-4 lg:gap-6 mb-8">
          <Card>
            <CardContent className="p-4 text-center">
              <Trophy className="h-6 w-6 text-yellow-500 mx-auto mb-2" />
              <div className="text-2xl font-bold">{userStats.completedChallenges}</div>
              <div className="text-sm text-muted-foreground">Completed</div>
            </CardContent>
          </Card>
          <Card>
            <CardContent className="p-4 text-center">
              <Target className="h-6 w-6 text-blue-500 mx-auto mb-2" />
              <div className="text-2xl font-bold">{userStats.totalChallenges - userStats.completedChallenges}</div>
              <div className="text-sm text-muted-foreground">Remaining</div>
            </CardContent>
          </Card>
          <Card>
            <CardContent className="p-4 text-center">
              <Zap className="h-6 w-6 text-orange-500 mx-auto mb-2" />
              <div className="text-2xl font-bold">{userStats.currentStreak}</div>
              <div className="text-sm text-muted-foreground">Day Streak</div>
            </CardContent>
          </Card>
          <Card>
            <CardContent className="p-4 text-center">
              <Star className="h-6 w-6 text-purple-500 mx-auto mb-2" />
              <div className="text-2xl font-bold">{userStats.rank}</div>
              <div className="text-sm text-muted-foreground">Current Rank</div>
            </CardContent>
          </Card>
        </div>

        <Progress 
          value={(userStats.completedChallenges / userStats.totalChallenges) * 100} 
          className="h-2" 
        />
      </div>

      {/* Main Content */}
      <Tabs value={activeTab} onValueChange={setActiveTab} className="space-y-6">
        <TabsList className="grid w-full grid-cols-3 lg:w-[400px]">
          <TabsTrigger value="challenges">Challenges</TabsTrigger>
          <TabsTrigger value="assessments">Assessments</TabsTrigger>
          <TabsTrigger value="leaderboard">Leaderboard</TabsTrigger>
        </TabsList>

        {/* Challenges Tab */}
        <TabsContent value="challenges" className="space-y-6">
          {/* Filters */}
          <Card>
            <CardContent className="p-4">
              <div className="flex flex-col md:flex-row gap-4">
                <div className="flex-1">
                  <div className="relative">
                    <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 h-4 w-4 text-muted-foreground" />
                    <input
                      type="text"
                      placeholder="Search challenges..."
                      value={searchQuery}
                      onChange={(e) => setSearchQuery(e.target.value)}
                      className="w-full pl-10 pr-4 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-primary"
                    />
                  </div>
                </div>
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
                  <select
                    value={selectedCategory}
                    onChange={(e) => setSelectedCategory(e.target.value)}
                    className="px-3 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-primary"
                  >
                    {categories.map(category => (
                      <option key={category} value={category}>
                        {category === "all" ? "All Categories" : category}
                      </option>
                    ))}
                  </select>
                </div>
              </div>
            </CardContent>
          </Card>

          {/* Challenge Grid */}
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 2xl:grid-cols-5 gap-6 lg:gap-8">
            {filteredChallenges.map((challenge) => (
              <Card key={challenge.id} className={`relative ${challenge.locked ? 'opacity-60' : ''}`}>
                {challenge.locked && (
                  <div className="absolute top-4 right-4 z-10">
                    <Lock className="h-5 w-5 text-muted-foreground" />
                  </div>
                )}
                <CardHeader>
                  <div className="flex items-start justify-between">
                    <div className="flex-1">
                      <CardTitle className="text-lg mb-2">{challenge.title}</CardTitle>
                      <CardDescription className="text-sm">
                        {challenge.description}
                      </CardDescription>
                    </div>
                  </div>
                  <div className="flex items-center gap-2 mt-3">
                    <Badge className={getDifficultyColor(challenge.difficulty)}>
                      {challenge.difficulty}
                    </Badge>
                    <Badge variant="outline">{challenge.category}</Badge>
                  </div>
                </CardHeader>
                <CardContent>
                  <div className="space-y-3">
                    <div className="flex items-center justify-between text-sm">
                      <div className="flex items-center space-x-4">
                        <div className="flex items-center space-x-1">
                          <Clock className="h-4 w-4 text-muted-foreground" />
                          <span>{challenge.estimatedTime}</span>
                        </div>
                        <div className="flex items-center space-x-1">
                          <Star className="h-4 w-4 text-yellow-500" />
                          <span>{challenge.points} pts</span>
                        </div>
                      </div>
                    </div>
                    
                    <div className="flex items-center justify-between">
                      <div className="text-sm text-muted-foreground">
                        {challenge.completionRate}% completion rate
                      </div>
                      {challenge.completed && (
                        <CheckCircle className="h-5 w-5 text-green-500" />
                      )}
                    </div>

                    <div className="flex flex-wrap gap-1 mt-2">
                      {challenge.tags.slice(0, 3).map((tag) => (
                        <Badge key={tag} variant="secondary" className="text-xs">
                          {tag}
                        </Badge>
                      ))}
                    </div>

                    <Link href={`/practice/challenge/${challenge.id}`}>
                      <Button 
                        className="w-full mt-4" 
                        disabled={challenge.locked}
                        variant={challenge.completed ? "outline" : "default"}
                      >
                        {challenge.locked ? (
                          <>
                            <Lock className="mr-2 h-4 w-4" />
                            Locked
                          </>
                        ) : challenge.completed ? (
                          <>
                            <CheckCircle className="mr-2 h-4 w-4" />
                            Review
                          </>
                        ) : (
                          <>
                            <Play className="mr-2 h-4 w-4" />
                            Start Challenge
                          </>
                        )}
                      </Button>
                    </Link>
                  </div>
                </CardContent>
              </Card>
            ))}
          </div>
        </TabsContent>

        {/* Assessments Tab */}
        <TabsContent value="assessments" className="space-y-8">
          <div className="grid grid-cols-1 lg:grid-cols-2 xl:grid-cols-3 gap-6 lg:gap-8">
            {skillAssessments.map((assessment) => (
              <Card key={assessment.id}>
                <CardHeader>
                  <div className="flex items-start justify-between">
                    <div className="flex-1">
                      <CardTitle className="text-xl mb-2">{assessment.title}</CardTitle>
                      <CardDescription>
                        {assessment.description}
                      </CardDescription>
                    </div>
                    {assessment.completed && assessment.score && (
                      <div className="text-right">
                        <div className="text-2xl font-bold text-primary">
                          {assessment.score}%
                        </div>
                        <div className="text-sm text-muted-foreground">Score</div>
                      </div>
                    )}
                  </div>
                  <div className="flex items-center gap-2 mt-3">
                    <Badge className={getDifficultyColor(assessment.difficulty)}>
                      {assessment.difficulty}
                    </Badge>
                    <Badge variant="outline">{assessment.category}</Badge>
                  </div>
                </CardHeader>
                <CardContent>
                  <div className="space-y-4">
                    <div className="grid grid-cols-2 gap-4 text-sm">
                      <div className="flex items-center space-x-2">
                        <BookOpen className="h-4 w-4 text-muted-foreground" />
                        <span>{assessment.questions} questions</span>
                      </div>
                      <div className="flex items-center space-x-2">
                        <Clock className="h-4 w-4 text-muted-foreground" />
                        <span>{assessment.duration}</span>
                      </div>
                    </div>

                    {assessment.completed && assessment.score && (
                      <div className="space-y-2">
                        <div className="flex justify-between text-sm">
                          <span>Your Score</span>
                          <span>{assessment.score}/{assessment.maxScore}</span>
                        </div>
                        <Progress value={assessment.score} className="h-2" />
                      </div>
                    )}

                    <Link href={`/practice/assessment/${assessment.id}`}>
                      <Button className="w-full" variant={assessment.completed ? "outline" : "default"}>
                        {assessment.completed ? (
                          <>
                            <TrendingUp className="mr-2 h-4 w-4" />
                            Retake Assessment
                          </>
                        ) : (
                          <>
                            <Brain className="mr-2 h-4 w-4" />
                            Start Assessment
                          </>
                        )}
                      </Button>
                    </Link>
                  </div>
                </CardContent>
              </Card>
            ))}
          </div>
        </TabsContent>

        {/* Leaderboard Tab */}
        <TabsContent value="leaderboard" className="space-y-8">
          <div className="grid grid-cols-1 lg:grid-cols-3 xl:grid-cols-4 gap-6 lg:gap-8">
            {/* Global Leaderboard */}
            <div className="lg:col-span-2">
              <Card>
                <CardHeader>
                  <CardTitle className="flex items-center">
                    <Trophy className="mr-2 h-5 w-5 text-yellow-500" />
                    Global Leaderboard
                  </CardTitle>
                  <CardDescription>
                    Top performers this month
                  </CardDescription>
                </CardHeader>
                <CardContent>
                  <div className="space-y-4">
                    {[
                      { rank: 1, name: "Alex Chen", points: 2450, avatar: "AC", streak: 15 },
                      { rank: 2, name: "Sarah Johnson", points: 2380, avatar: "SJ", streak: 12 },
                      { rank: 3, name: "Mike Rodriguez", points: 2290, avatar: "MR", streak: 8 },
                      { rank: 4, name: "You", points: userStats.totalPoints, avatar: "YU", streak: userStats.currentStreak },
                      { rank: 5, name: "Emma Wilson", points: 2150, avatar: "EW", streak: 6 },
                    ].map((user) => (
                      <div key={user.rank} className={`flex items-center justify-between p-3 rounded-lg ${user.name === "You" ? "bg-primary/10 border border-primary/20" : "bg-muted/50"}`}>
                        <div className="flex items-center space-x-3">
                          <div className={`w-8 h-8 rounded-full flex items-center justify-center text-sm font-bold ${
                            user.rank === 1 ? "bg-yellow-500 text-white" :
                            user.rank === 2 ? "bg-gray-400 text-white" :
                            user.rank === 3 ? "bg-amber-600 text-white" :
                            "bg-primary text-primary-foreground"
                          }`}>
                            {user.rank <= 3 ? user.rank : user.avatar}
                          </div>
                          <div>
                            <div className="font-medium">{user.name}</div>
                            <div className="text-sm text-muted-foreground">
                              {user.streak} day streak
                            </div>
                          </div>
                        </div>
                        <div className="text-right">
                          <div className="font-bold">{user.points}</div>
                          <div className="text-sm text-muted-foreground">points</div>
                        </div>
                      </div>
                    ))}
                  </div>
                </CardContent>
              </Card>
            </div>

            {/* Stats and Achievements */}
            <div className="space-y-6">
              <Card>
                <CardHeader>
                  <CardTitle className="flex items-center">
                    <Award className="mr-2 h-5 w-5 text-purple-500" />
                    Your Achievements
                  </CardTitle>
                </CardHeader>
                <CardContent>
                  <div className="space-y-3">
                    <div className="flex items-center space-x-3 p-2 rounded-lg bg-green-50 border border-green-200">
                      <CheckCircle className="h-6 w-6 text-green-600" />
                      <div>
                        <div className="font-medium text-green-800">First Steps</div>
                        <div className="text-sm text-green-600">Complete your first challenge</div>
                      </div>
                    </div>
                    <div className="flex items-center space-x-3 p-2 rounded-lg bg-blue-50 border border-blue-200">
                      <Zap className="h-6 w-6 text-blue-600" />
                      <div>
                        <div className="font-medium text-blue-800">Streak Master</div>
                        <div className="text-sm text-blue-600">5 day learning streak</div>
                      </div>
                    </div>
                    <div className="flex items-center space-x-3 p-2 rounded-lg bg-gray-50 border border-gray-200 opacity-60">
                      <Trophy className="h-6 w-6 text-gray-400" />
                      <div>
                        <div className="font-medium text-gray-600">Challenge Master</div>
                        <div className="text-sm text-gray-500">Complete 10 challenges</div>
                      </div>
                    </div>
                  </div>
                </CardContent>
              </Card>

              <Card>
                <CardHeader>
                  <CardTitle className="flex items-center">
                    <Users className="mr-2 h-5 w-5 text-blue-500" />
                    Community Stats
                  </CardTitle>
                </CardHeader>
                <CardContent>
                  <div className="space-y-3 text-sm">
                    <div className="flex justify-between">
                      <span>Active Learners</span>
                      <span className="font-bold">1,247</span>
                    </div>
                    <div className="flex justify-between">
                      <span>Challenges Completed Today</span>
                      <span className="font-bold">89</span>
                    </div>
                    <div className="flex justify-between">
                      <span>Average Success Rate</span>
                      <span className="font-bold">73%</span>
                    </div>
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
