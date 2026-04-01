'use client';

import { type Exercise } from '@/lib/topics-data';
import { CheckCircle, ChevronRight, Terminal, Lightbulb, AlertCircle } from 'lucide-react';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';

interface ExerciseCardProps {
  exercise: Exercise;
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
  index,
  completed = false,
  onComplete,
}) => {
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
          {!completed && (
            <Button 
              size="sm" 
              onClick={() => onComplete?.(exercise.id)}
            >
              Mark Complete
            </Button>
          )}
        </div>
      </CardContent>
    </Card>
  );
};

export default ExerciseCard;
