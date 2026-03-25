'use client';

import { Suspense } from 'react';
import { useSearchParams } from 'next/navigation';
import InterviewSession from '@/components/interviews/InterviewSession';
import type { InterviewType, Difficulty, Answer } from '@/types/interview';

function SessionContent() {
  const searchParams = useSearchParams();
  const type = (searchParams.get('type') || 'coding') as InterviewType;
  const difficulty = (searchParams.get('difficulty') || 'beginner') as Difficulty;

  const handleComplete = (answers: Answer[]) => {
    // Navigate to feedback page with answers
    const params = new URLSearchParams({
      type,
      difficulty,
      completed: 'true',
    });
    window.location.href = `/interviews/feedback?${params.toString()}`;
  };

  return (
    <InterviewSession
      type={type}
      difficulty={difficulty}
      onComplete={handleComplete}
    />
  );
}

export default function SessionPage() {
  return (
    <Suspense
      fallback={
        <div className="min-h-screen flex items-center justify-center">
          <div className="text-center space-y-4">
            <div className="animate-spin rounded-full h-12 w-12 border-4 border-primary border-t-transparent mx-auto" />
            <p className="text-muted-foreground">Loading interview session...</p>
          </div>
        </div>
      }
    >
      <SessionContent />
    </Suspense>
  );
}
