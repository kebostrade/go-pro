'use client';

import { useState } from 'react';
import { type Exercise, type Topic } from '@/lib/topics-data';
import { CheckCircle, ChevronRight, Terminal, Lightbulb, AlertCircle, Loader2, Sparkles, X } from 'lucide-react';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { api } from '@/lib/api';

interface ExerciseCardProps {
  exercise: Exercise;
  topic: Topic;
  index: number;
  completed?: boolean;
  onComplete?: (id: string) => void;
}

const difficultyColors = {
  easy: 'bg-green-100 text-green-800 dark:bg-green-900/30 dark:text-green-400',
  medium: 'bg-yellow-100 text-yellow-800 dark:bg-yellow-900/30 dark:text-yellow-400',
  hard: 'bg-red-100 text-red-800 dark:bg-red-900/30 dark:text-red-400',
};

export const ExerciseCard: React.FC<ExerciseCardProps> = ({
  exercise,
  topic,
  index,
  completed = false,
  onComplete,
}) => {
  const [isReviewing, setIsReviewing] = useState(false);
  const [feedback, setFeedback] = useState<string | null>(null);
  const [showFeedback, setShowFeedback] = useState(false);

  const handleSubmitReview = async () => {
    setIsReviewing(true);
    try {
      const result = await api.submitReview({
        user_id: 'current-user',
        topic_id: topic.id,
        exercise_id: exercise.id,
        code: exercise.starterCode || '',
      });
      setFeedback(result.feedback);
      setShowFeedback(true);
    } catch (err) {
      console.error('Review failed:', err);
    } finally {
      setIsReviewing(false);
    }
  };
  return (
    <Card className={completed ? 'border-green-500 bg-green-50/50 dark:bg-green-900/10' : ''}>
      <CardHeader className="pb-3">
        <div className="flex items-start justify-between">
          <div className="flex-1">
            <div className="flex items-center gap-2 mb-2">
              <Badge variant="outline" className="text-xs">
                Exercise {index}
              </Badge>
              <Badge className={difficultyColors[exercise.difficulty]}>
                {exercise.difficulty}
              </Badge>
              <Badge variant="secondary" className="text-xs">
                {exercise.completionMethod.replace('_', ' ')}
              </Badge>
            </div>
            <CardTitle className="text-lg">{exercise.title}</CardTitle>
          </div>
          {completed && (
            <CheckCircle className="h-6 w-6 text-green-500" />
          )}
        </div>
      </CardHeader>
      <CardContent className="pt-0">
        <p className="text-muted-foreground mb-4">{exercise.description}</p>
        
        {/* Requirements */}
        <div className="mb-4">
          <h4 className="text-sm font-semibold mb-2 flex items-center gap-1">
            <AlertCircle className="w-4 h-4" />
            Requirements
          </h4>
          <ul className="space-y-1">
            {exercise.requirements.map((req, i) => (
              <li key={i} className="text-sm flex items-start gap-2">
                <ChevronRight className="w-4 h-4 mt-0.5 flex-shrink-0" />
                {req}
              </li>
            ))}
          </ul>
        </div>
        
        {/* Hint */}
        <div className="bg-blue-50 dark:bg-blue-900/20 border border-blue-200 dark:border-blue-800 rounded-lg p-3 mb-4">
          <div className="flex items-start gap-2">
            <Lightbulb className="w-4 h-4 text-blue-600 dark:text-blue-400 mt-0.5" />
            <div>
              <p className="text-sm font-medium text-blue-800 dark:text-blue-200">Hint</p>
              <p className="text-sm text-blue-700 dark:text-blue-300">{exercise.solutionHint}</p>
            </div>
          </div>
        </div>
        
        {/* Actions */}
        <div className="flex gap-2">
          <Button variant="outline" size="sm" className="flex-1">
            <Terminal className="w-4 h-4 mr-2" />
            View Starter Code
          </Button>
          <button
            onClick={handleSubmitReview}
            disabled={isReviewing}
            className="flex items-center gap-2 px-3 py-1.5 bg-purple-600 hover:bg-purple-700 disabled:bg-purple-400 text-white rounded-lg text-sm font-medium transition-colors"
          >
            {isReviewing ? (
              <>
                <Loader2 className="w-4 h-4 animate-spin" />
                Analyzing...
              </>
            ) : (
              <>
                <Sparkles className="w-4 h-4" />
                Submit for Review
              </>
            )}
          </button>
          {!completed && (
            <Button 
              size="sm" 
              onClick={() => onComplete?.(exercise.id)}
            >
              Mark Complete
            </Button>
          )}
        </div>

        {showFeedback && feedback && (
          <div className="mt-4 p-4 bg-gray-50 dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700">
            <div className="flex items-center justify-between mb-2">
              <h4 className="font-semibold text-sm">AI Feedback</h4>
              <button
                onClick={() => setShowFeedback(false)}
                className="text-gray-500 hover:text-gray-700"
              >
                <X className="w-4 h-4" />
              </button>
            </div>
            <p className="text-sm text-gray-700 dark:text-gray-300 whitespace-pre-wrap">
              {feedback}
            </p>
          </div>
        )}
      </CardContent>
    </Card>
  );
};

export default ExerciseCard;
