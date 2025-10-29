"use client";

import { useState } from "react";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Progress } from "@/components/ui/progress";
import { 
  BookOpen, 
  ArrowRight, 
  Clock, 
  Target,
  Star,
  TrendingUp,
  Users,
  CheckCircle,
  Play,
  Lightbulb,
  Zap,
  Award
} from "lucide-react";
import Link from "next/link";

interface RelatedLesson {
  id: number;
  title: string;
  description: string;
  difficulty: "Beginner" | "Intermediate" | "Advanced";
  duration: string;
  progress?: number;
  completed?: boolean;
  rating?: number;
  enrolledCount?: number;
  category: string;
  tags: string[];
  prerequisite?: boolean;
  recommended?: boolean;
  nextInPath?: boolean;
}

interface RelatedLessonsProps {
  currentLessonId: number;
  lessons: RelatedLesson[];
  title?: string;
  maxItems?: number;
  showCategories?: boolean;
  className?: string;
}

export default function RelatedLessons({
  currentLessonId,
  lessons,
  title = "Related Lessons",
  maxItems = 6,
  showCategories = true,
  className = ""
}: RelatedLessonsProps) {
  const [selectedCategory, setSelectedCategory] = useState<string | null>(null);

  // Filter out current lesson and apply category filter
  const filteredLessons = lessons
    .filter(lesson => lesson.id !== currentLessonId)
    .filter(lesson => !selectedCategory || lesson.category === selectedCategory)
    .slice(0, maxItems);

  // Get unique categories
  const categories = Array.from(new Set(lessons.map(lesson => lesson.category)));

  const getDifficultyColor = (difficulty: string) => {
    switch (difficulty) {
      case 'Beginner': return 'text-green-600 bg-green-100 border-green-200 dark:bg-green-950 dark:border-green-800';
      case 'Intermediate': return 'text-yellow-600 bg-yellow-100 border-yellow-200 dark:bg-yellow-950 dark:border-yellow-800';
      case 'Advanced': return 'text-red-600 bg-red-100 border-red-200 dark:bg-red-950 dark:border-red-800';
      default: return 'text-blue-600 bg-blue-100 border-blue-200 dark:bg-blue-950 dark:border-blue-800';
    }
  };

  const getLessonTypeIcon = (lesson: RelatedLesson) => {
    if (lesson.prerequisite) return <Target className="h-4 w-4 text-orange-500" />;
    if (lesson.nextInPath) return <ArrowRight className="h-4 w-4 text-blue-500" />;
    if (lesson.recommended) return <Star className="h-4 w-4 text-yellow-500" />;
    return <BookOpen className="h-4 w-4 text-primary" />;
  };

  const getLessonTypeLabel = (lesson: RelatedLesson) => {
    if (lesson.prerequisite) return "Prerequisite";
    if (lesson.nextInPath) return "Next in Path";
    if (lesson.recommended) return "Recommended";
    return "Related";
  };

  const renderLessonCard = (lesson: RelatedLesson) => (
    <Card key={lesson.id} className="group hover:border-primary/50 transition-all hover-lift">
      <CardHeader className="pb-3">
        <div className="flex items-start justify-between mb-2">
          <div className="flex items-center space-x-2">
            {getLessonTypeIcon(lesson)}
            <Badge variant="outline" className="text-xs">
              {getLessonTypeLabel(lesson)}
            </Badge>
          </div>
          
          {lesson.completed && (
            <CheckCircle className="h-5 w-5 text-green-500" />
          )}
        </div>

        <CardTitle className="text-lg leading-tight group-hover:text-primary transition-colors">
          {lesson.title}
        </CardTitle>
        
        <CardDescription className="text-sm line-clamp-2">
          {lesson.description}
        </CardDescription>
      </CardHeader>

      <CardContent className="pt-0">
        <div className="space-y-3">
          {/* Progress Bar */}
          {lesson.progress !== undefined && (
            <div className="space-y-1">
              <div className="flex items-center justify-between text-xs">
                <span className="text-muted-foreground">Progress</span>
                <span className="font-medium">{lesson.progress}%</span>
              </div>
              <Progress value={lesson.progress} className="h-1.5" />
            </div>
          )}

          {/* Lesson Metadata */}
          <div className="flex items-center justify-between">
            <div className="flex items-center space-x-2">
              <Badge variant="outline" className={`text-xs ${getDifficultyColor(lesson.difficulty)}`}>
                {lesson.difficulty}
              </Badge>
              <Badge variant="outline" className="text-xs">
                <Clock className="mr-1 h-3 w-3" />
                {lesson.duration}
              </Badge>
            </div>

            {lesson.rating && (
              <div className="flex items-center space-x-1 text-xs text-muted-foreground">
                <Star className="h-3 w-3 fill-current text-yellow-500" />
                <span>{lesson.rating.toFixed(1)}</span>
              </div>
            )}
          </div>

          {/* Tags */}
          {lesson.tags.length > 0 && (
            <div className="flex flex-wrap gap-1">
              {lesson.tags.slice(0, 3).map(tag => (
                <Badge key={tag} variant="secondary" className="text-xs">
                  {tag}
                </Badge>
              ))}
              {lesson.tags.length > 3 && (
                <Badge variant="secondary" className="text-xs">
                  +{lesson.tags.length - 3}
                </Badge>
              )}
            </div>
          )}

          {/* Enrollment Count */}
          {lesson.enrolledCount && (
            <div className="flex items-center space-x-1 text-xs text-muted-foreground">
              <Users className="h-3 w-3" />
              <span>{lesson.enrolledCount.toLocaleString()} enrolled</span>
            </div>
          )}

          {/* Action Button */}
          <Link href={`/learn/${lesson.id}`} className="block">
            <Button 
              className="w-full mt-3 group-hover:bg-primary group-hover:text-primary-foreground transition-colors"
              variant={lesson.nextInPath ? "default" : "outline"}
            >
              {lesson.completed ? (
                <>
                  <CheckCircle className="mr-2 h-4 w-4" />
                  Review
                </>
              ) : lesson.progress && lesson.progress > 0 ? (
                <>
                  <Play className="mr-2 h-4 w-4" />
                  Continue
                </>
              ) : (
                <>
                  <BookOpen className="mr-2 h-4 w-4" />
                  Start Lesson
                </>
              )}
            </Button>
          </Link>
        </div>
      </CardContent>
    </Card>
  );

  return (
    <Card className={`glass-card border-2 ${className}`}>
      <CardHeader>
        <CardTitle className="flex items-center">
          <TrendingUp className="mr-2 h-5 w-5 text-primary" />
          {title}
        </CardTitle>
        
        {showCategories && categories.length > 1 && (
          <div className="flex flex-wrap gap-2 mt-3">
            <Button
              variant={selectedCategory === null ? "default" : "outline"}
              size="sm"
              onClick={() => setSelectedCategory(null)}
              className="text-xs"
            >
              All
            </Button>
            {categories.map(category => (
              <Button
                key={category}
                variant={selectedCategory === category ? "default" : "outline"}
                size="sm"
                onClick={() => setSelectedCategory(category)}
                className="text-xs"
              >
                {category}
              </Button>
            ))}
          </div>
        )}
      </CardHeader>

      <CardContent>
        {filteredLessons.length > 0 ? (
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
            {filteredLessons.map(renderLessonCard)}
          </div>
        ) : (
          <div className="text-center py-8 text-muted-foreground">
            <BookOpen className="mx-auto h-12 w-12 mb-4 opacity-50" />
            <p>No related lessons found</p>
            <p className="text-sm">Try selecting a different category</p>
          </div>
        )}

        {lessons.length > maxItems && (
          <div className="text-center mt-6">
            <Link href="/curriculum">
              <Button variant="outline" className="hover-glow">
                <BookOpen className="mr-2 h-4 w-4" />
                View All Lessons
              </Button>
            </Link>
          </div>
        )}
      </CardContent>
    </Card>
  );
}
