"use client";

import { useState } from "react";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Progress } from "@/components/ui/progress";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import {
  Code2,
  Trophy,
  Flame,
  Clock,
  Target,
  Calendar,
  TrendingUp,
  BookOpen,
  CheckCircle,
  Circle,
  LoaderCircle,
  Play,
  BarChart3,
  ListTodo,
  Lightbulb
} from "lucide-react";
import { Category, AlgoProgress, Session, Problem } from "@/types/algorithms";
import ProblemsList from "./problems-list";
import SessionLog from "./session-log";
import PatternNotes from "./pattern-notes";

interface AlgorithmsTrackerProps {
  categories: Category[];
  progress: AlgoProgress;
  sessions: Session[];
  problems: Problem[];
}

const AlgorithmsTracker = ({ categories, progress, sessions, problems }: AlgorithmsTrackerProps) => {
  const [selectedCategory, setSelectedCategory] = useState<string | null>(null);

  const getCategoryIcon = (iconName: string) => {
    const icons: Record<string, React.ElementType> = {
      Array: BarChart3,
      Link: Code2,
      Layers: BookOpen,
      GitBranch: Code2,
      ArrowUpCircle: TrendingUp,
      Hash: Target,
      Network: Code2,
      Zap: Flame,
      RotateCcw: Code2,
      ArrowUpDown: Code2,
    };
    return icons[iconName] || Code2;
  };

  const formatTime = (minutes: number) => {
    if (minutes < 60) return `${minutes}m`;
    const hours = Math.floor(minutes / 60);
    const mins = minutes % 60;
    return mins > 0 ? `${hours}h ${mins}m` : `${hours}h`;
  };

  const getCompletionColor = (percentage: number) => {
    if (percentage >= 80) return 'text-green-600';
    if (percentage >= 50) return 'text-yellow-600';
    if (percentage >= 25) return 'text-orange-600';
    return 'text-gray-600';
  };

  return (
    <div className="space-y-6">
      {/* Stats Overview */}
      <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
        <Card>
          <CardContent className="pt-6">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm text-muted-foreground">Solved</p>
                <p className="text-2xl font-bold">{progress.completedProblems}</p>
                <p className="text-xs text-muted-foreground">of {progress.totalProblems}</p>
              </div>
              <CheckCircle className="h-8 w-8 text-green-500 opacity-80" />
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardContent className="pt-6">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm text-muted-foreground">Streak</p>
                <p className="text-2xl font-bold">{progress.currentStreak} days</p>
                <p className="text-xs text-muted-foreground">best: {progress.longestStreak}</p>
              </div>
              <Flame className="h-8 w-8 text-orange-500 opacity-80" />
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardContent className="pt-6">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm text-muted-foreground">Time Spent</p>
                <p className="text-2xl font-bold">{formatTime(progress.totalTimeSpent)}</p>
                <p className="text-xs text-muted-foreground">avg: {formatTime(progress.averageTimePerProblem)}/prob</p>
              </div>
              <Clock className="h-8 w-8 text-blue-500 opacity-80" />
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardContent className="pt-6">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm text-muted-foreground">Progress</p>
                <p className="text-2xl font-bold">{Math.round((progress.completedProblems / progress.totalProblems) * 100)}%</p>
                <Progress value={(progress.completedProblems / progress.totalProblems) * 100} className="mt-2 h-1.5" />
              </div>
              <Target className="h-8 w-8 text-purple-500 opacity-80" />
            </div>
          </CardContent>
        </Card>
      </div>

      {/* Difficulty Breakdown */}
      <Card>
        <CardHeader className="pb-3">
          <CardTitle className="text-lg">Difficulty Progress</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="flex gap-6">
            <div className="flex items-center gap-2">
              <Badge className="bg-green-100 text-green-800 border-green-200">Easy</Badge>
              <span className="font-medium">{progress.easySolved}</span>
              <span className="text-sm text-muted-foreground">solved</span>
            </div>
            <div className="flex items-center gap-2">
              <Badge className="bg-yellow-100 text-yellow-800 border-yellow-200">Medium</Badge>
              <span className="font-medium">{progress.mediumSolved}</span>
              <span className="text-sm text-muted-foreground">solved</span>
            </div>
            <div className="flex items-center gap-2">
              <Badge className="bg-red-100 text-red-800 border-red-200">Hard</Badge>
              <span className="font-medium">{progress.hardSolved}</span>
              <span className="text-sm text-muted-foreground">solved</span>
            </div>
          </div>
        </CardContent>
      </Card>

      {/* Categories Grid */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-5 gap-4">
        {categories.map((category) => {
          const Icon = getCategoryIcon(category.icon);
          const percentage = category.totalProblems > 0
            ? Math.round((category.completedProblems / category.totalProblems) * 100)
            : 0;

          return (
            <Card
              key={category.id}
              className={`cursor-pointer transition-all hover:shadow-md ${
                selectedCategory === category.id ? 'ring-2 ring-primary' : ''
              }`}
              onClick={() => setSelectedCategory(selectedCategory === category.id ? null : category.id)}
            >
              <CardContent className="pt-4 pb-4">
                <div className="flex items-center gap-3 mb-3">
                  <div className={`p-2 rounded-lg bg-${category.color}-100`}>
                    <Icon className={`h-5 w-5 text-${category.color}-600`} />
                  </div>
                  <div className="flex-1 min-w-0">
                    <h3 className="font-medium text-sm truncate">{category.name}</h3>
                    <p className="text-xs text-muted-foreground">
                      {category.completedProblems}/{category.totalProblems}
                    </p>
                  </div>
                </div>
                <Progress value={percentage} className="h-1.5" />
                <div className="flex justify-between mt-2">
                  <span className={`text-xs font-medium ${getCompletionColor(percentage)}`}>
                    {percentage}%
                  </span>
                  {category.inProgressProblems > 0 && (
                    <span className="text-xs text-blue-600 flex items-center gap-1">
                      <LoaderCircle className="h-3 w-3" />
                      {category.inProgressProblems} active
                    </span>
                  )}
                </div>
              </CardContent>
            </Card>
          );
        })}
      </div>

      {/* Main Content Tabs */}
      <Tabs defaultValue="problems" className="space-y-4">
        <TabsList className="grid w-full grid-cols-3 lg:w-auto lg:inline-grid">
          <TabsTrigger value="problems" className="flex items-center gap-2">
            <ListTodo className="h-4 w-4" />
            Problems
          </TabsTrigger>
          <TabsTrigger value="sessions" className="flex items-center gap-2">
            <Calendar className="h-4 w-4" />
            Sessions
          </TabsTrigger>
          <TabsTrigger value="patterns" className="flex items-center gap-2">
            <Lightbulb className="h-4 w-4" />
            Patterns
          </TabsTrigger>
        </TabsList>

        <TabsContent value="problems">
          <ProblemsList
            problems={problems}
            selectedCategory={selectedCategory}
            categories={categories}
          />
        </TabsContent>

        <TabsContent value="sessions">
          <SessionLog sessions={sessions} />
        </TabsContent>

        <TabsContent value="patterns">
          <PatternNotes />
        </TabsContent>
      </Tabs>
    </div>
  );
};

export default AlgorithmsTracker;
