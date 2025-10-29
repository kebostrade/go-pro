"use client";

import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Progress } from "@/components/ui/progress";
import {
  Clock,
  BookOpen,
  Brain,
  TrendingUp,
  CheckCircle,
  Target,
  Award,
  RotateCcw
} from "lucide-react";
import Link from "next/link";

interface AssessmentCardProps {
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
  onStart?: () => void;
  onRetake?: () => void;
}

const AssessmentCard = ({
  id,
  title,
  description,
  category,
  questions,
  duration,
  difficulty,
  completed,
  score,
  maxScore,
  onStart,
  onRetake
}: AssessmentCardProps) => {
  const getDifficultyColor = (difficulty: string) => {
    switch (difficulty) {
      case "Beginner": return "text-green-600 bg-green-50 border-green-200";
      case "Intermediate": return "text-yellow-600 bg-yellow-50 border-yellow-200";
      case "Advanced": return "text-red-600 bg-red-50 border-red-200";
      default: return "text-gray-600 bg-gray-50 border-gray-200";
    }
  };

  const getScoreColor = (score: number, maxScore: number) => {
    const percentage = (score / maxScore) * 100;
    if (percentage >= 90) return "text-green-600";
    if (percentage >= 70) return "text-yellow-600";
    return "text-red-600";
  };

  const getPerformanceBadge = (score: number, maxScore: number) => {
    const percentage = (score / maxScore) * 100;
    if (percentage >= 90) return { text: "Excellent", color: "bg-green-100 text-green-800 border-green-200" };
    if (percentage >= 80) return { text: "Good", color: "bg-blue-100 text-blue-800 border-blue-200" };
    if (percentage >= 70) return { text: "Pass", color: "bg-yellow-100 text-yellow-800 border-yellow-200" };
    return { text: "Needs Improvement", color: "bg-red-100 text-red-800 border-red-200" };
  };

  const handleAction = () => {
    if (completed && onRetake) {
      onRetake();
    } else if (onStart) {
      onStart();
    }
  };

  return (
    <Card className="transition-all hover:shadow-md">
      <CardHeader>
        <div className="flex items-start justify-between">
          <div className="flex-1">
            <CardTitle className="text-xl mb-2">{title}</CardTitle>
            <CardDescription>
              {description}
            </CardDescription>
          </div>
          {completed && score !== undefined && (
            <div className="text-right ml-4">
              <div className={`text-2xl font-bold ${getScoreColor(score, maxScore)}`}>
                {Math.round((score / maxScore) * 100)}%
              </div>
              <div className="text-sm text-muted-foreground">Score</div>
            </div>
          )}
        </div>
        
        <div className="flex items-center gap-2 mt-3">
          <Badge className={getDifficultyColor(difficulty)}>
            {difficulty}
          </Badge>
          <Badge variant="outline">{category}</Badge>
          {completed && score !== undefined && (
            <Badge className={getPerformanceBadge(score, maxScore).color}>
              {getPerformanceBadge(score, maxScore).text}
            </Badge>
          )}
        </div>
      </CardHeader>
      
      <CardContent>
        <div className="space-y-4">
          <div className="grid grid-cols-2 gap-4 text-sm">
            <div className="flex items-center space-x-2">
              <BookOpen className="h-4 w-4 text-muted-foreground" />
              <span>{questions} questions</span>
            </div>
            <div className="flex items-center space-x-2">
              <Clock className="h-4 w-4 text-muted-foreground" />
              <span>{duration}</span>
            </div>
          </div>

          {completed && score !== undefined && (
            <div className="space-y-3">
              <div className="space-y-2">
                <div className="flex justify-between text-sm">
                  <span>Your Score</span>
                  <span className="font-medium">{score}/{maxScore}</span>
                </div>
                <Progress value={(score / maxScore) * 100} className="h-2" />
              </div>
              
              <div className="grid grid-cols-2 gap-4 text-sm">
                <div className="flex items-center space-x-2">
                  <Target className="h-4 w-4 text-blue-500" />
                  <span>Accuracy: {Math.round((score / maxScore) * 100)}%</span>
                </div>
                <div className="flex items-center space-x-2">
                  <Award className="h-4 w-4 text-purple-500" />
                  <span>Rank: Top {Math.floor(Math.random() * 30) + 10}%</span>
                </div>
              </div>
            </div>
          )}

          {!completed && (
            <div className="space-y-3">
              <div className="p-3 bg-blue-50 border border-blue-200 rounded-lg">
                <div className="flex items-center space-x-2 text-blue-800">
                  <Brain className="h-4 w-4" />
                  <span className="text-sm font-medium">Assessment Overview</span>
                </div>
                <p className="text-sm text-blue-700 mt-1">
                  Test your knowledge and get personalized feedback on your progress.
                </p>
              </div>
            </div>
          )}

          <Link href={`/practice/assessment/${id}`}>
            <Button 
              className="w-full" 
              variant={completed ? "outline" : "default"}
              onClick={handleAction}
            >
              {completed ? (
                <>
                  <RotateCcw className="mr-2 h-4 w-4" />
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

          {completed && (
            <div className="flex items-center justify-center space-x-4 text-sm text-muted-foreground pt-2">
              <div className="flex items-center space-x-1">
                <TrendingUp className="h-4 w-4" />
                <span>Completed {Math.floor(Math.random() * 30) + 1} days ago</span>
              </div>
            </div>
          )}
        </div>
      </CardContent>
    </Card>
  );
};

export default AssessmentCard;
