"use client";

import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Progress } from "@/components/ui/progress";
import { 
  BookOpen, 
  Clock, 
  Target, 
  Trophy,
  Play,
  CheckCircle,
  Lock,
  Star,
  Users,
  ArrowRight
} from "lucide-react";

interface LessonCardProps {
  id: string;
  title: string;
  description: string;
  duration: string;
  difficulty: 'Beginner' | 'Intermediate' | 'Advanced';
  type: 'lesson' | 'exercise' | 'project';
  completed: boolean;
  locked: boolean;
  progress?: number;
  rating?: number;
  enrolledCount?: number;
  tags?: string[];
  onStart?: () => void;
  onContinue?: () => void;
  onReview?: () => void;
}

const LessonCard = ({
  id,
  title,
  description,
  duration,
  difficulty,
  type,
  completed,
  locked,
  progress = 0,
  rating,
  enrolledCount,
  tags = [],
  onStart,
  onContinue,
  onReview
}: LessonCardProps) => {
  const getDifficultyColor = (level: string) => {
    switch (level) {
      case 'Beginner': return 'bg-green-500 text-white';
      case 'Intermediate': return 'bg-yellow-500 text-white';
      case 'Advanced': return 'bg-red-500 text-white';
      default: return 'bg-gray-500 text-white';
    }
  };

  const getTypeIcon = (lessonType: string) => {
    switch (lessonType) {
      case 'lesson': return BookOpen;
      case 'exercise': return Target;
      case 'project': return Trophy;
      default: return BookOpen;
    }
  };

  const getTypeColor = (lessonType: string) => {
    switch (lessonType) {
      case 'lesson': return 'text-blue-500 bg-blue-50 dark:bg-blue-950';
      case 'exercise': return 'text-green-500 bg-green-50 dark:bg-green-950';
      case 'project': return 'text-purple-500 bg-purple-50 dark:bg-purple-950';
      default: return 'text-gray-500 bg-gray-50 dark:bg-gray-950';
    }
  };

  const TypeIcon = getTypeIcon(type);

  const handleAction = () => {
    if (locked) return;
    
    if (completed) {
      onReview?.();
    } else if (progress > 0) {
      onContinue?.();
    } else {
      onStart?.();
    }
  };

  const getActionText = () => {
    if (locked) return 'Locked';
    if (completed) return 'Review';
    if (progress > 0) return 'Continue';
    return 'Start';
  };

  return (
    <Card className={`lesson-card transition-all duration-200 ${
      locked ? 'opacity-60' : 'hover:shadow-lg'
    } ${completed ? 'ring-1 ring-green-200 dark:ring-green-800' : ''}`}>
      <CardHeader className="space-y-3">
        <div className="flex items-start justify-between">
          <div className="flex items-center space-x-3">
            <div className={`flex h-10 w-10 items-center justify-center rounded-lg ${
              locked 
                ? 'bg-muted text-muted-foreground' 
                : getTypeColor(type)
            }`}>
              {locked ? (
                <Lock className="h-5 w-5" />
              ) : completed ? (
                <CheckCircle className="h-5 w-5 text-green-500" />
              ) : (
                <TypeIcon className="h-5 w-5" />
              )}
            </div>
            <div>
              <Badge 
                variant="secondary" 
                className={`text-xs ${getDifficultyColor(difficulty)} mb-1`}
              >
                {difficulty}
              </Badge>
              <Badge variant="outline" className="text-xs ml-2 capitalize">
                {type}
              </Badge>
            </div>
          </div>
          
          {completed && (
            <Badge className="bg-green-500 text-white">
              <CheckCircle className="mr-1 h-3 w-3" />
              Completed
            </Badge>
          )}
        </div>

        <div>
          <CardTitle className={`text-lg mb-2 ${locked ? 'text-muted-foreground' : ''}`}>
            {title}
          </CardTitle>
          <CardDescription className="text-sm line-clamp-2">
            {description}
          </CardDescription>
        </div>
      </CardHeader>

      <CardContent className="space-y-4">
        {/* Progress Bar (if in progress) */}
        {progress > 0 && progress < 100 && (
          <div className="space-y-2">
            <div className="flex justify-between text-sm">
              <span className="text-muted-foreground">Progress</span>
              <span className="font-medium">{Math.round(progress)}%</span>
            </div>
            <Progress value={progress} className="h-2" />
          </div>
        )}

        {/* Lesson Meta Info */}
        <div className="flex items-center justify-between text-sm text-muted-foreground">
          <div className="flex items-center space-x-4">
            <div className="flex items-center space-x-1">
              <Clock className="h-3 w-3" />
              <span>{duration}</span>
            </div>
            
            {rating && (
              <div className="flex items-center space-x-1">
                <Star className="h-3 w-3 fill-yellow-400 text-yellow-400" />
                <span>{rating.toFixed(1)}</span>
              </div>
            )}
            
            {enrolledCount && (
              <div className="flex items-center space-x-1">
                <Users className="h-3 w-3" />
                <span>{enrolledCount.toLocaleString()}</span>
              </div>
            )}
          </div>
        </div>

        {/* Tags */}
        {tags.length > 0 && (
          <div className="flex flex-wrap gap-1">
            {tags.slice(0, 3).map((tag, index) => (
              <Badge key={index} variant="outline" className="text-xs">
                {tag}
              </Badge>
            ))}
            {tags.length > 3 && (
              <Badge variant="outline" className="text-xs">
                +{tags.length - 3} more
              </Badge>
            )}
          </div>
        )}

        {/* Action Button */}
        <Button
          onClick={handleAction}
          disabled={locked}
          className={`w-full ${
            completed 
              ? 'variant-outline' 
              : locked 
                ? '' 
                : 'go-gradient text-white'
          }`}
          variant={completed ? 'outline' : 'default'}
        >
          {locked ? (
            <>
              <Lock className="mr-2 h-4 w-4" />
              {getActionText()}
            </>
          ) : completed ? (
            <>
              <BookOpen className="mr-2 h-4 w-4" />
              {getActionText()}
            </>
          ) : progress > 0 ? (
            <>
              <Play className="mr-2 h-4 w-4" />
              {getActionText()}
              <ArrowRight className="ml-2 h-4 w-4" />
            </>
          ) : (
            <>
              <Play className="mr-2 h-4 w-4" />
              {getActionText()}
              <ArrowRight className="ml-2 h-4 w-4" />
            </>
          )}
        </Button>
      </CardContent>
    </Card>
  );
};

export default LessonCard;
