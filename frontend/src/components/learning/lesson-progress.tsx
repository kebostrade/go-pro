"use client";

import { useState, useEffect } from "react";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Progress } from "@/components/ui/progress";
import { 
  CheckCircle, 
  Circle, 
  Clock, 
  Target, 
  Award, 
  TrendingUp,
  Calendar,
  Flame,
  Star,
  BookOpen,
  Code2,
  Trophy,
  Zap
} from "lucide-react";

interface LessonProgressData {
  lessonId: number;
  completedObjectives: Set<number>;
  timeSpent: number;
  lastAccessed: Date;
  tabsVisited: Set<string>;
  exercisesCompleted: Set<string>;
  notesCount: number;
  bookmarked: boolean;
  completionPercentage: number;
}

interface LessonProgressProps {
  lessonId: number;
  objectives: string[];
  totalExercises: number;
  onProgressUpdate?: (progress: LessonProgressData) => void;
}

export default function LessonProgress({ 
  lessonId, 
  objectives, 
  totalExercises,
  onProgressUpdate 
}: LessonProgressProps) {
  const [progress, setProgress] = useState<LessonProgressData>({
    lessonId,
    completedObjectives: new Set(),
    timeSpent: 0,
    lastAccessed: new Date(),
    tabsVisited: new Set(),
    exercisesCompleted: new Set(),
    notesCount: 0,
    bookmarked: false,
    completionPercentage: 0,
  });

  const [streak, setStreak] = useState(0);
  const [achievements, setAchievements] = useState<string[]>([]);

  // Load progress from localStorage
  useEffect(() => {
    const savedProgress = localStorage.getItem(`lesson-progress-${lessonId}`);
    if (savedProgress) {
      const parsed = JSON.parse(savedProgress);
      setProgress({
        ...parsed,
        completedObjectives: new Set(parsed.completedObjectives),
        tabsVisited: new Set(parsed.tabsVisited),
        exercisesCompleted: new Set(parsed.exercisesCompleted),
        lastAccessed: new Date(parsed.lastAccessed),
      });
    }
  }, [lessonId]);

  // Save progress to localStorage
  useEffect(() => {
    const progressToSave = {
      ...progress,
      completedObjectives: Array.from(progress.completedObjectives),
      tabsVisited: Array.from(progress.tabsVisited),
      exercisesCompleted: Array.from(progress.exercisesCompleted),
    };
    localStorage.setItem(`lesson-progress-${lessonId}`, JSON.stringify(progressToSave));
    onProgressUpdate?.(progress);
  }, [progress, lessonId, onProgressUpdate]);

  // Calculate completion percentage
  useEffect(() => {
    const objectiveWeight = 40;
    const exerciseWeight = 40;
    const engagementWeight = 20;

    const objectiveProgress = (progress.completedObjectives.size / objectives.length) * objectiveWeight;
    const exerciseProgress = totalExercises > 0 ? (progress.exercisesCompleted.size / totalExercises) * exerciseWeight : exerciseWeight;
    const engagementProgress = Math.min(progress.tabsVisited.size / 3, 1) * engagementWeight;

    const totalProgress = objectiveProgress + exerciseProgress + engagementProgress;
    
    setProgress(prev => ({ ...prev, completionPercentage: Math.round(totalProgress) }));
  }, [progress.completedObjectives, progress.exercisesCompleted, progress.tabsVisited, objectives.length, totalExercises]);

  const toggleObjective = (index: number) => {
    setProgress(prev => {
      const newCompleted = new Set(prev.completedObjectives);
      if (newCompleted.has(index)) {
        newCompleted.delete(index);
      } else {
        newCompleted.add(index);
        // Check for achievements
        if (newCompleted.size === objectives.length) {
          setAchievements(prev => [...prev, 'objectives-master']);
        }
      }
      return { ...prev, completedObjectives: newCompleted };
    });
  };

  const markTabVisited = (tab: string) => {
    setProgress(prev => ({
      ...prev,
      tabsVisited: new Set([...prev.tabsVisited, tab]),
    }));
  };

  const markExerciseCompleted = (exerciseId: string) => {
    setProgress(prev => ({
      ...prev,
      exercisesCompleted: new Set([...prev.exercisesCompleted, exerciseId]),
    }));
  };

  const updateTimeSpent = (seconds: number) => {
    setProgress(prev => ({ ...prev, timeSpent: seconds }));
  };

  const toggleBookmark = () => {
    setProgress(prev => ({ ...prev, bookmarked: !prev.bookmarked }));
  };

  const getProgressColor = (percentage: number) => {
    if (percentage >= 90) return "text-green-500";
    if (percentage >= 70) return "text-blue-500";
    if (percentage >= 50) return "text-yellow-500";
    return "text-gray-500";
  };

  const getProgressMessage = (percentage: number) => {
    if (percentage >= 100) return "Lesson completed! 🎉";
    if (percentage >= 90) return "Almost there! 🚀";
    if (percentage >= 70) return "Great progress! 💪";
    if (percentage >= 50) return "Halfway there! 📈";
    if (percentage >= 25) return "Good start! 👍";
    return "Just getting started! 🌱";
  };

  return (
    <div className="space-y-6">
      {/* Overall Progress Card */}
      <Card className="glass-card border-2">
        <CardHeader>
          <CardTitle className="flex items-center justify-between">
            <span className="flex items-center">
              <TrendingUp className="mr-2 h-5 w-5 text-primary" />
              Lesson Progress
            </span>
            <Badge variant="outline" className={getProgressColor(progress.completionPercentage)}>
              {progress.completionPercentage}%
            </Badge>
          </CardTitle>
          <CardDescription>
            {getProgressMessage(progress.completionPercentage)}
          </CardDescription>
        </CardHeader>
        
        <CardContent className="space-y-4">
          <Progress value={progress.completionPercentage} className="h-3" />
          
          {/* Stats Grid */}
          <div className="grid grid-cols-2 lg:grid-cols-4 gap-4">
            <div className="text-center p-3 bg-gradient-to-r from-blue-500/10 to-blue-500/5 rounded-lg border border-blue-500/20">
              <Target className="h-5 w-5 text-blue-500 mx-auto mb-1" />
              <div className="text-sm font-medium">{progress.completedObjectives.size}/{objectives.length}</div>
              <div className="text-xs text-muted-foreground">Objectives</div>
            </div>
            
            <div className="text-center p-3 bg-gradient-to-r from-green-500/10 to-green-500/5 rounded-lg border border-green-500/20">
              <Trophy className="h-5 w-5 text-green-500 mx-auto mb-1" />
              <div className="text-sm font-medium">{progress.exercisesCompleted.size}/{totalExercises}</div>
              <div className="text-xs text-muted-foreground">Exercises</div>
            </div>
            
            <div className="text-center p-3 bg-gradient-to-r from-orange-500/10 to-orange-500/5 rounded-lg border border-orange-500/20">
              <Clock className="h-5 w-5 text-orange-500 mx-auto mb-1" />
              <div className="text-sm font-medium">{Math.floor(progress.timeSpent / 60)}m</div>
              <div className="text-xs text-muted-foreground">Time Spent</div>
            </div>
            
            <div className="text-center p-3 bg-gradient-to-r from-purple-500/10 to-purple-500/5 rounded-lg border border-purple-500/20">
              <BookOpen className="h-5 w-5 text-purple-500 mx-auto mb-1" />
              <div className="text-sm font-medium">{progress.tabsVisited.size}/3</div>
              <div className="text-xs text-muted-foreground">Sections</div>
            </div>
          </div>
        </CardContent>
      </Card>

      {/* Objectives Checklist */}
      <Card className="glass-card">
        <CardHeader>
          <CardTitle className="flex items-center">
            <Target className="mr-2 h-5 w-5 text-primary" />
            Learning Objectives
          </CardTitle>
        </CardHeader>
        
        <CardContent>
          <div className="space-y-3">
            {objectives.map((objective, index) => (
              <div
                key={index}
                className="flex items-start space-x-3 p-3 rounded-lg hover:bg-muted/50 cursor-pointer transition-colors group"
                onClick={() => toggleObjective(index)}
              >
                <div className="mt-0.5">
                  {progress.completedObjectives.has(index) ? (
                    <CheckCircle className="h-5 w-5 text-green-500" />
                  ) : (
                    <Circle className="h-5 w-5 text-muted-foreground group-hover:text-primary transition-colors" />
                  )}
                </div>
                <span className={`text-sm leading-relaxed transition-all ${
                  progress.completedObjectives.has(index) 
                    ? 'line-through text-muted-foreground' 
                    : 'group-hover:text-foreground'
                }`}>
                  {objective}
                </span>
              </div>
            ))}
          </div>
        </CardContent>
      </Card>

      {/* Achievements */}
      {achievements.length > 0 && (
        <Card className="glass-card border-yellow-500/20 bg-gradient-to-r from-yellow-500/5 to-transparent">
          <CardHeader>
            <CardTitle className="flex items-center">
              <Award className="mr-2 h-5 w-5 text-yellow-500" />
              Achievements Unlocked!
            </CardTitle>
          </CardHeader>
          
          <CardContent>
            <div className="flex flex-wrap gap-2">
              {achievements.map(achievement => (
                <Badge key={achievement} variant="outline" className="border-yellow-500 text-yellow-600">
                  <Star className="mr-1 h-3 w-3" />
                  {achievement === 'objectives-master' && 'Objectives Master'}
                </Badge>
              ))}
            </div>
          </CardContent>
        </Card>
      )}
    </div>
  );
}
