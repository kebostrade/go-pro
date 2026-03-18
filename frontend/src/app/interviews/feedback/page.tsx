'use client';

import { Suspense, useState, useEffect } from 'react';
import { useSearchParams } from 'next/navigation';
import FeedbackDisplay from '@/components/interviews/FeedbackDisplay';
import { Loader2 } from 'lucide-react';
import { api } from '@/lib/api';
import type { InterviewFeedback } from '@/types/interview';

function FeedbackContent() {
  const searchParams = useSearchParams();
  const type = (searchParams.get('type') || 'coding') as 'coding' | 'behavioral' | 'system_design';
  const difficulty = (searchParams.get('difficulty') || 'beginner') as 'beginner' | 'intermediate' | 'advanced';

  const [feedback, setFeedback] = useState<InterviewFeedback | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const fetchFeedback = async () => {
      setIsLoading(true);
      setError(null);

      try {
        // Try to get sessions from API
        const sessions = await api.getInterviewSessions();

        if (sessions && sessions.length > 0) {
          // Get the most recent completed session
          const latestSession = sessions.find(s => s.status === 'completed') || sessions[0];

          if (latestSession) {
            const feedbackData = await api.getInterviewFeedback(latestSession.id);
            setFeedback(feedbackData);
          } else {
            // Use mock feedback if no sessions
            setFeedback(getMockFeedback(type, difficulty));
          }
        } else {
          // Use mock feedback if API unavailable or no sessions
          setFeedback(getMockFeedback(type, difficulty));
        }
      } catch (err) {
        console.log('Using mock feedback:', err);
        // Use mock feedback on error
        setFeedback(getMockFeedback(type, difficulty));
      } finally {
        setIsLoading(false);
      }
    };

    fetchFeedback();
  }, [type, difficulty]);

  const handleRetry = () => {
    window.location.href = `/interviews/session?type=${type}&difficulty=${difficulty}`;
  };

  if (isLoading) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="text-center space-y-4">
          <Loader2 className="h-12 w-12 animate-spin text-primary mx-auto" />
          <p className="text-muted-foreground">Loading feedback...</p>
        </div>
      </div>
    );
  }

  if (error || !feedback) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="text-center space-y-4">
          <p className="text-destructive">{error || 'Failed to load feedback'}</p>
          <button
            onClick={handleRetry}
            className="text-primary hover:underline"
          >
            Try Again
          </button>
        </div>
      </div>
    );
  }

  return <FeedbackDisplay feedback={feedback} onRetry={handleRetry} />;
}

// Mock feedback generator for demo mode
function getMockFeedback(type: string, difficulty: string): InterviewFeedback {
  const feedbackByType: Record<string, InterviewFeedback> = {
    coding: {
      session_id: `session-coding-${difficulty}`,
      overall_score: 78,
      strengths: [
        'Clear code structure and good naming conventions',
        'Proper use of Go idioms and best practices',
        'Good error handling demonstrated',
        'Efficient algorithm design',
      ],
      improvements: [
        'Consider edge cases in algorithm design',
        'Add more inline comments for complex logic',
        'Optimize time complexity where possible',
        'Include unit tests for your solution',
      ],
      detailed_feedback: [
        {
          question_id: 'q1',
          score: 85,
          feedback: 'Excellent solution with clean code structure. You demonstrated good understanding of Go slices and iteration patterns.',
          strengths: ['Efficient algorithm', 'Clean code', 'Proper error handling', 'Good variable naming'],
          missed: ['Could optimize memory usage'],
        },
        {
          question_id: 'q2',
          score: 72,
          feedback: 'Correct solution but missing some edge case handling. Consider empty inputs and boundary conditions.',
          strengths: ['Correct logic', 'Good approach'],
          missed: ['Did not handle empty input', 'Missing nil checks', 'No error handling'],
        },
        {
          question_id: 'q3',
          score: 78,
          feedback: 'Well-structured answer with clear explanation. Good use of standard library functions.',
          strengths: ['Clear explanation', 'Structured response', 'Appropriate use of stdlib'],
          missed: ['Could provide more examples', 'Missing complexity analysis'],
        },
      ],
    },
    behavioral: {
      session_id: `session-behavioral-${difficulty}`,
      overall_score: 82,
      strengths: [
        'Strong use of STAR method',
        'Clear and concise communication',
        'Good self-awareness and reflection',
        'Demonstrated growth mindset',
      ],
      improvements: [
        'Provide more specific metrics and outcomes',
        'Include more details about team dynamics',
        'Discuss lessons learned from challenges',
        'Quantify your impact more clearly',
      ],
      detailed_feedback: [
        {
          question_id: 'q1',
          score: 88,
          feedback: 'Excellent use of the STAR method. Your situation was well-contextualized and the action steps were clear.',
          strengths: ['Clear context', 'Specific actions', 'Good outcome', 'Showed growth'],
          missed: ['Could quantify impact more'],
        },
        {
          question_id: 'q2',
          score: 76,
          feedback: 'Good response but could benefit from more specific examples and measurable outcomes.',
          strengths: ['Relevant example', 'Clear communication'],
          missed: ['Missing metrics', 'Could elaborate on challenges'],
        },
      ],
    },
    system_design: {
      session_id: `session-system-design-${difficulty}`,
      overall_score: 75,
      strengths: [
        'Good understanding of core system components',
        'Clear explanation of data flow',
        'Considered scalability aspects',
        'Discussed trade-offs appropriately',
      ],
      improvements: [
        'Include more detail on database schema',
        'Discuss failure modes and recovery',
        'Add capacity estimation and planning',
        'Consider security implications',
      ],
      detailed_feedback: [
        {
          question_id: 'q1',
          score: 80,
          feedback: 'Solid design with good consideration of scalability. The API design was well-thought-out.',
          strengths: ['Clear architecture', 'Good API design', 'Scalability considerations'],
          missed: ['Missing caching strategy', 'No discussion of monitoring'],
        },
        {
          question_id: 'q2',
          score: 70,
          feedback: 'Good foundation but needs more depth on distributed systems challenges.',
          strengths: ['Understood core concepts', 'Reasonable approach'],
          missed: ['Missing CAP theorem discussion', 'No mention of data consistency'],
        },
      ],
    },
  };

  return feedbackByType[type] || feedbackByType.coding;
}

export default function FeedbackPage() {
  return (
    <Suspense
      fallback={
        <div className="min-h-screen flex items-center justify-center">
          <div className="text-center space-y-4">
            <Loader2 className="h-12 w-12 animate-spin text-primary mx-auto" />
            <p className="text-muted-foreground">Loading feedback...</p>
          </div>
        </div>
      }
    >
      <FeedbackContent />
    </Suspense>
  );
}
