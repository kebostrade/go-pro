"use client";

import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Progress } from "@/components/ui/progress";
import {
  Clock,
  Star,
  CheckCircle,
  Lock,
  Play,
  Target,
  TrendingUp,
  Users
} from "lucide-react";
import Link from "next/link";

interface ChallengeCardProps {
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
  onStart?: () => void;
  onContinue?: () => void;
  onReview?: () => void;
}

const ChallengeCard = ({
  id,
  title,
  description,
  difficulty,
  category,
  estimatedTime,
  points,
  completed,
  locked,
  tags,
  completionRate,
  attempts,
  onStart,
  onContinue,
  onReview
}: ChallengeCardProps) => {
  const getDifficultyColor = (difficulty: string) => {
    switch (difficulty) {
      case "Beginner": return "text-green-600 bg-green-50 border-green-200";
      case "Intermediate": return "text-yellow-600 bg-yellow-50 border-yellow-200";
      case "Advanced": return "text-red-600 bg-red-50 border-red-200";
      default: return "text-gray-600 bg-gray-50 border-gray-200";
    }
  };

  const handleAction = () => {
    if (locked) return;
    
    if (completed && onReview) {
      onReview();
    } else if (attempts > 0 && onContinue) {
      onContinue();
    } else if (onStart) {
      onStart();
    }
  };

  return (
    <Card className={`relative transition-all hover:shadow-md ${locked ? 'opacity-60' : ''}`}>
      {locked && (
        <div className="absolute top-4 right-4 z-10">
          <Lock className="h-5 w-5 text-muted-foreground" />
        </div>
      )}
      
      <CardHeader>
        <div className="flex items-start justify-between">
          <div className="flex-1">
            <CardTitle className="text-lg mb-2 line-clamp-2">{title}</CardTitle>
            <CardDescription className="text-sm line-clamp-2">
              {description}
            </CardDescription>
          </div>
        </div>
        
        <div className="flex items-center gap-2 mt-3">
          <Badge className={getDifficultyColor(difficulty)}>
            {difficulty}
          </Badge>
          <Badge variant="outline">{category}</Badge>
          {completed && (
            <Badge className="bg-green-100 text-green-800 border-green-200">
              <CheckCircle className="mr-1 h-3 w-3" />
              Completed
            </Badge>
          )}
        </div>
      </CardHeader>
      
      <CardContent>
        <div className="space-y-4">
          <div className="flex items-center justify-between text-sm">
            <div className="flex items-center space-x-4">
              <div className="flex items-center space-x-1">
                <Clock className="h-4 w-4 text-muted-foreground" />
                <span>{estimatedTime}</span>
              </div>
              <div className="flex items-center space-x-1">
                <Star className="h-4 w-4 text-yellow-500" />
                <span>{points} pts</span>
              </div>
            </div>
          </div>
          
          <div className="space-y-2">
            <div className="flex items-center justify-between text-sm">
              <span className="text-muted-foreground">Success Rate</span>
              <span className="font-medium">{completionRate}%</span>
            </div>
            <Progress value={completionRate} className="h-2" />
          </div>

          <div className="flex items-center justify-between text-sm text-muted-foreground">
            <div className="flex items-center space-x-1">
              <Users className="h-4 w-4" />
              <span>{Math.floor(Math.random() * 1000) + 500} attempts</span>
            </div>
            {attempts > 0 && (
              <div className="flex items-center space-x-1">
                <TrendingUp className="h-4 w-4" />
                <span>Your attempts: {attempts}</span>
              </div>
            )}
          </div>

          <div className="flex flex-wrap gap-1 mt-2">
            {tags.slice(0, 3).map((tag) => (
              <Badge key={tag} variant="secondary" className="text-xs">
                {tag}
              </Badge>
            ))}
            {tags.length > 3 && (
              <Badge variant="secondary" className="text-xs">
                +{tags.length - 3} more
              </Badge>
            )}
          </div>

          <Link href={`/practice/challenge/${id}`}>
            <Button 
              className="w-full mt-4" 
              disabled={locked}
              variant={completed ? "outline" : "default"}
              onClick={handleAction}
            >
              {locked ? (
                <>
                  <Lock className="mr-2 h-4 w-4" />
                  Locked
                </>
              ) : completed ? (
                <>
                  <CheckCircle className="mr-2 h-4 w-4" />
                  Review Solution
                </>
              ) : attempts > 0 ? (
                <>
                  <Play className="mr-2 h-4 w-4" />
                  Continue Challenge
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
  );
};

export default ChallengeCard;
