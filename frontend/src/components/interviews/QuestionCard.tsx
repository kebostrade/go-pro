'use client';

import { useState } from 'react';
import { Card, CardContent, CardHeader } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Clock, Lightbulb, ChevronLeft, ChevronRight } from 'lucide-react';
import type { Question } from '@/types/interview';

interface QuestionCardProps {
  question: Question;
  currentIndex: number;
  totalQuestions: number;
  timeRemaining: number;
  hintsUsed: number;
  maxHints: number;
  onNext: () => void;
  onPrevious: () => void;
  onHint: () => void;
  onEnd: () => void;
}

export default function QuestionCard({
  question,
  currentIndex,
  totalQuestions,
  timeRemaining,
  hintsUsed,
  maxHints,
  onNext,
  onPrevious,
  onHint,
  onEnd,
}: QuestionCardProps) {
  const [showHint, setShowHint] = useState(false);
  const [currentHintIndex, setCurrentHintIndex] = useState(0);

  const formatTime = (seconds: number) => {
    const mins = Math.floor(seconds / 60);
    const secs = seconds % 60;
    return `${mins}:${secs.toString().padStart(2, '0')}`;
  };

  const handleHint = () => {
    if (hintsUsed < maxHints) {
      setShowHint(true);
      onHint();
    }
  };

  const hasHints = question.expected_points && question.expected_points.length > 0;
  const hints = question.expected_points || [];

  return (
    <Card className="border-none shadow-lg mb-6">
      <CardHeader className="flex items-center justify-between pb-4">
        <div className="flex items-center gap-3">
          <Badge variant="outline" className="text-sm">
            Question {currentIndex + 1} of {totalQuestions}
          </Badge>
          {question.type && (
            <Badge variant="secondary" className="text-sm capitalize">
              {question.type.replace('_', ' ')}
            </Badge>
          )}
        </div>
        <Button
          variant="ghost"
          size="sm"
          onClick={onEnd}
          className="text-muted-foreground hover:text-destructive"
        >
          End Interview
        </Button>
      </CardHeader>
      <CardContent className="space-y-4">
        {/* Timer */}
        <div
          className={`flex items-center justify-center gap-2 py-3 px-4 rounded-lg ${
            timeRemaining < 60 ? 'bg-red-500/10 text-red-400' : 'bg-muted/30'
          }`}
        >
          <Clock className={`h-5 w-5 ${timeRemaining < 60 ? 'animate-pulse' : ''}`} />
          <span className="text-xl font-mono font-semibold">
            {formatTime(timeRemaining)}
          </span>
        </div>

        {/* Question */}
        <div className="bg-card rounded-lg p-6 border border-border/50">
          <p className="text-lg leading-relaxed">{question.content}</p>
        </div>

        {/* Hints Section */}
        {hasHints && (
          <div className="space-y-2">
            <div className="flex items-center justify-between">
              <span className="text-sm text-muted-foreground">
                Hints: {hintsUsed}/{maxHints}
              </span>
              {hintsUsed < maxHints && (
                <Button
                  variant="outline"
                  size="sm"
                  onClick={handleHint}
                  className="flex items-center gap-1"
                >
                  <Lightbulb className="h-4 w-4" />
                  Get Hint
                </Button>
              )}
            </div>
            {showHint && currentHintIndex < hints.length && (
              <div className="bg-amber-500/10 border border-amber-500/30 rounded-lg p-4 text-sm">
                <div className="flex items-start gap-2">
                  <Lightbulb className="h-4 w-4 text-amber-500 shrink-0 mt-0.5" />
                  <div>
                    <div className="font-semibold text-amber-400">Hint {currentHintIndex + 1}</div>
                    <div className="text-amber-200">{hints[currentHintIndex]}</div>
                  </div>
                </div>
              </div>
            )}
          </div>
        )}

        {/* Navigation */}
        <div className="flex items-center justify-between pt-4 border-t border-border/50">
          <Button
            variant="outline"
            onClick={onPrevious}
            disabled={currentIndex === 0}
            className="flex items-center gap-2"
          >
            <ChevronLeft className="h-4 w-4" />
            Previous
          </Button>
          <Button onClick={onNext} className="flex items-center gap-2">
            Next
            <ChevronRight className="h-4 w-4" />
          </Button>
        </div>
      </CardContent>
    </Card>
  );
}
