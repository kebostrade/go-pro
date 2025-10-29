"use client";

import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Progress } from "@/components/ui/progress";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { 
  BookOpen, 
  Clock, 
  Trophy, 
  Target,
  Play,
  CheckCircle,
  Lock
} from "lucide-react";

interface Lesson {
  id: string;
  title: string;
  duration: string;
  completed: boolean;
  locked: boolean;
  type: 'lesson' | 'exercise' | 'project';
}

interface CourseProgressProps {
  courseTitle: string;
  courseDescription: string;
  totalLessons: number;
  completedLessons: number;
  estimatedTime: string;
  difficulty: 'Beginner' | 'Intermediate' | 'Advanced';
  lessons: Lesson[];
}

const CourseProgress = ({
  courseTitle,
  courseDescription,
  totalLessons,
  completedLessons,
  estimatedTime,
  difficulty,
  lessons
}: CourseProgressProps) => {
  const progressPercentage = (completedLessons / totalLessons) * 100;
  
  const getDifficultyColor = (level: string) => {
    switch (level) {
      case 'Beginner': return 'bg-green-500';
      case 'Intermediate': return 'bg-yellow-500';
      case 'Advanced': return 'bg-red-500';
      default: return 'bg-gray-500';
    }
  };

  const getTypeIcon = (type: string) => {
    switch (type) {
      case 'lesson': return BookOpen;
      case 'exercise': return Target;
      case 'project': return Trophy;
      default: return BookOpen;
    }
  };

  const getTypeColor = (type: string) => {
    switch (type) {
      case 'lesson': return 'text-blue-500';
      case 'exercise': return 'text-green-500';
      case 'project': return 'text-purple-500';
      default: return 'text-gray-500';
    }
  };

  return (
    <div className="space-y-6">
      {/* Course Overview */}
      <Card className="lesson-card">
        <CardHeader>
          <div className="flex items-start justify-between">
            <div className="space-y-2">
              <CardTitle className="text-2xl">{courseTitle}</CardTitle>
              <CardDescription className="text-base max-w-2xl">
                {courseDescription}
              </CardDescription>
            </div>
            <Badge 
              variant="secondary" 
              className={`${getDifficultyColor(difficulty)} text-white`}
            >
              {difficulty}
            </Badge>
          </div>
        </CardHeader>
        <CardContent className="space-y-6">
          {/* Progress Stats */}
          <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
            <div className="flex items-center space-x-3">
              <div className="flex h-10 w-10 items-center justify-center rounded-lg bg-primary/10">
                <BookOpen className="h-5 w-5 text-primary" />
              </div>
              <div>
                <p className="text-sm text-muted-foreground">Progress</p>
                <p className="text-lg font-semibold">{completedLessons}/{totalLessons} lessons</p>
              </div>
            </div>
            
            <div className="flex items-center space-x-3">
              <div className="flex h-10 w-10 items-center justify-center rounded-lg bg-primary/10">
                <Clock className="h-5 w-5 text-primary" />
              </div>
              <div>
                <p className="text-sm text-muted-foreground">Time Left</p>
                <p className="text-lg font-semibold">{estimatedTime}</p>
              </div>
            </div>
            
            <div className="flex items-center space-x-3">
              <div className="flex h-10 w-10 items-center justify-center rounded-lg bg-primary/10">
                <Trophy className="h-5 w-5 text-primary" />
              </div>
              <div>
                <p className="text-sm text-muted-foreground">Completion</p>
                <p className="text-lg font-semibold">{Math.round(progressPercentage)}%</p>
              </div>
            </div>
          </div>
          
          {/* Progress Bar */}
          <div className="space-y-2">
            <div className="flex justify-between text-sm">
              <span className="text-muted-foreground">Course Progress</span>
              <span className="font-medium">{Math.round(progressPercentage)}%</span>
            </div>
            <Progress value={progressPercentage} className="h-2" />
          </div>
        </CardContent>
      </Card>

      {/* Lessons List */}
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center">
            <BookOpen className="mr-2 h-5 w-5" />
            Course Content
          </CardTitle>
          <CardDescription>
            {totalLessons} lessons • {estimatedTime} total
          </CardDescription>
        </CardHeader>
        <CardContent>
          <div className="space-y-3">
            {lessons.map((lesson, index) => {
              const TypeIcon = getTypeIcon(lesson.type);
              const isNext = !lesson.completed && !lesson.locked && 
                           lessons.slice(0, index).every(l => l.completed);
              
              return (
                <div
                  key={lesson.id}
                  className={`flex items-center justify-between p-4 rounded-lg border transition-colors ${
                    lesson.completed 
                      ? 'bg-green-50 border-green-200 dark:bg-green-950 dark:border-green-800' 
                      : lesson.locked 
                        ? 'bg-muted/50 border-border opacity-60' 
                        : isNext
                          ? 'bg-primary/5 border-primary/20'
                          : 'bg-background border-border hover:bg-accent/50'
                  }`}
                >
                  <div className="flex items-center space-x-3">
                    <div className={`flex h-8 w-8 items-center justify-center rounded-lg ${
                      lesson.completed 
                        ? 'bg-green-500 text-white' 
                        : lesson.locked 
                          ? 'bg-muted text-muted-foreground'
                          : 'bg-primary/10'
                    }`}>
                      {lesson.completed ? (
                        <CheckCircle className="h-4 w-4" />
                      ) : lesson.locked ? (
                        <Lock className="h-4 w-4" />
                      ) : (
                        <TypeIcon className={`h-4 w-4 ${getTypeColor(lesson.type)}`} />
                      )}
                    </div>
                    
                    <div>
                      <h4 className={`font-medium ${lesson.locked ? 'text-muted-foreground' : ''}`}>
                        {lesson.title}
                      </h4>
                      <div className="flex items-center space-x-2 text-sm text-muted-foreground">
                        <Clock className="h-3 w-3" />
                        <span>{lesson.duration}</span>
                        <Badge variant="outline" className="text-xs capitalize">
                          {lesson.type}
                        </Badge>
                      </div>
                    </div>
                  </div>
                  
                  <div className="flex items-center space-x-2">
                    {lesson.completed && (
                      <Badge variant="secondary" className="text-xs">
                        Completed
                      </Badge>
                    )}
                    {isNext && (
                      <Badge className="text-xs go-gradient text-white">
                        Next
                      </Badge>
                    )}
                    <Button
                      size="sm"
                      variant={lesson.completed ? "outline" : "default"}
                      disabled={lesson.locked}
                      className={lesson.completed ? "" : "go-gradient text-white"}
                    >
                      {lesson.completed ? (
                        "Review"
                      ) : lesson.locked ? (
                        <Lock className="h-4 w-4" />
                      ) : (
                        <>
                          <Play className="mr-1 h-3 w-3" />
                          Start
                        </>
                      )}
                    </Button>
                  </div>
                </div>
              );
            })}
          </div>
        </CardContent>
      </Card>
    </div>
  );
};

export default CourseProgress;
