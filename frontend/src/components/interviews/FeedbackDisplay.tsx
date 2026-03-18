'use client';

import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import {
  CheckCircle,
  XCircle,
  TrendingUp,
  AlertTriangle,
  RotateCcw,
  ArrowLeft,
  Share2,
  Download,
} from 'lucide-react';
import type { InterviewFeedback } from '@/types/interview';

interface FeedbackDisplayProps {
  feedback: InterviewFeedback;
  onRetry?: () => void;
}

// Mock feedback for demo
const mockFeedback: InterviewFeedback = {
  session_id: 'demo-session',
  overall_score: 75,
  strengths: [
    'Clear communication of your thought process',
    'Good use of appropriate Go idioms',
    'Proper error handling demonstrated',
  ],
  improvements: [
    'Consider edge cases in algorithm design',
    'Optimize time complexity where possible',
    'Add more detailed comments to complex logic',
  ],
  detailed_feedback: [
    {
      question_id: 'q1',
      score: 80,
      feedback: 'Strong solution with clean code structure. Good use of slices.',
      strengths: ['Efficient algorithm', 'Clean code', 'Proper error handling'],
      missed: ['Could optimize memory usage'],
    },
    {
      question_id: 'q2',
      score: 70,
      feedback: 'Correct solution but missing some edge case handling.',
      strengths: ['Correct logic', 'Good variable naming'],
      missed: ['Did not handle empty input', 'Missing nil checks'],
    },
    {
      question_id: 'q3',
      score: 75,
      feedback: 'Well-structured answer with clear explanation.',
      strengths: ['Clear explanation', 'Structured response'],
      missed: ['Could provide more examples'],
    },
  ],
};

export default function FeedbackDisplay({ feedback = mockFeedback, onRetry }: FeedbackDisplayProps) {
  const getScoreColor = (score: number) => {
    if (score >= 90) return 'text-green-500';
    if (score >= 70) return 'text-blue-500';
    if (score >= 50) return 'text-yellow-500';
    return 'text-red-500';
  };

  const getScoreLabel = (score: number) => {
    if (score >= 90) return 'Excellent';
    if (score >= 70) return 'Good';
    if (score >= 50) return 'Needs Improvement';
    return 'Keep Practicing';
  };

  return (
    <div className="min-h-screen bg-gradient-to-br from-background via-background to-background">
      <div className="container mx-auto px-4 py-8 max-w-5xl">
        {/* Header */}
        <div className="mb-8">
          <Button
            variant="ghost"
            onClick={() => window.location.href = '/interviews/practice'}
            className="mb-4"
          >
            <ArrowLeft className="h-4 w-4 mr-2" />
            Back to Practice
          </Button>
          <h1 className="text-4xl font-bold tracking-tight flex items-center gap-3">
            <CheckCircle className="h-10 w-10 text-primary" />
            Interview Results
          </h1>
          <p className="text-muted-foreground mt-2 text-lg">
            Here's how you performed in your interview
          </p>
        </div>

        {/* Overall Score Card */}
        <Card className="border-none shadow-lg mb-6 bg-gradient-to-br from-primary/10 to-primary/5">
          <CardContent className="p-8">
            <div className="flex flex-col md:flex-row items-center justify-between gap-6">
              <div className="text-center md:text-left">
                <div className={`text-6xl font-bold ${getScoreColor(feedback.overall_score)}`}>
                  {feedback.overall_score}%
                </div>
                <div className="text-xl text-muted-foreground capitalize">
                  {getScoreLabel(feedback.overall_score)}
                </div>
              </div>
              <div className="flex gap-3">
                <Button variant="outline" className="flex items-center gap-2">
                  <Share2 className="h-4 w-4" />
                  Share
                </Button>
                <Button variant="outline" className="flex items-center gap-2">
                  <Download className="h-4 w-4" />
                  Export
                </Button>
              </div>
            </div>
          </CardContent>
        </Card>

        <div className="grid grid-cols-1 md:grid-cols-2 gap-6 mb-6">
          {/* Strengths */}
          <Card className="border-none shadow-lg">
            <CardHeader>
              <CardTitle className="flex items-center gap-2 text-green-500">
                <CheckCircle className="h-5 w-5" />
                Strengths
              </CardTitle>
            </CardHeader>
            <CardContent>
              <ul className="space-y-3">
                {feedback.strengths.map((strength, index) => (
                  <li key={index} className="flex items-start gap-3">
                    <CheckCircle className="h-5 w-5 text-green-500 shrink-0 mt-0.5" />
                    <span className="text-sm">{strength}</span>
                  </li>
                ))}
              </ul>
            </CardContent>
          </Card>

          {/* Areas to Improve */}
          <Card className="border-none shadow-lg">
            <CardHeader>
              <CardTitle className="flex items-center gap-2 text-amber-500">
                <AlertTriangle className="h-5 w-5" />
                Areas to Improve
              </CardTitle>
            </CardHeader>
            <CardContent>
              <ul className="space-y-3">
                {feedback.improvements.map((improvement, index) => (
                  <li key={index} className="flex items-start gap-3">
                    <XCircle className="h-5 w-5 text-amber-500 shrink-0 mt-0.5" />
                    <span className="text-sm">{improvement}</span>
                  </li>
                ))}
              </ul>
            </CardContent>
          </Card>
        </div>

        {/* Detailed Feedback by Question */}
        <Card className="border-none shadow-lg mb-6">
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <TrendingUp className="h-5 w-5 text-primary" />
              Question-by-Question Breakdown
            </CardTitle>
          </CardHeader>
          <CardContent className="space-y-4">
            {feedback.detailed_feedback.map((item, index) => (
              <div key={item.question_id} className="bg-muted/30 rounded-lg p-4 border border-border/50">
                <div className="flex items-center justify-between mb-3">
                  <Badge variant="outline" className="text-sm">
                    Question {index + 1}
                  </Badge>
                  <Badge
                    variant="outline"
                    className={`${getScoreColor(item.score)} border-current`}
                  >
                    Score: {item.score}%
                  </Badge>
                </div>
                <p className="text-sm text-muted-foreground mb-3">{item.feedback}</p>

                {item.strengths && item.strengths.length > 0 && (
                  <div className="space-y-2">
                    <div className="text-sm font-semibold text-green-500">What went well:</div>
                    <ul className="space-y-1 ml-4">
                      {item.strengths.map((strength, i) => (
                        <li key={i} className="text-sm text-muted-foreground flex items-start gap-2">
                          <CheckCircle className="h-3 w-3 text-green-500 shrink-0 mt-0.5" />
                          <span>{strength}</span>
                        </li>
                      ))}
                    </ul>
                  </div>
                )}

                {item.missed && item.missed.length > 0 && (
                  <div className="space-y-2 mt-3">
                    <div className="text-sm font-semibold text-amber-500">To improve:</div>
                    <ul className="space-y-1 ml-4">
                      {item.missed.map((missed, i) => (
                        <li key={i} className="text-sm text-muted-foreground flex items-start gap-2">
                          <XCircle className="h-3 w-3 text-amber-500 shrink-0 mt-0.5" />
                          <span>{missed}</span>
                        </li>
                      ))}
                    </ul>
                  </div>
                )}
              </div>
            ))}
          </CardContent>
        </Card>

        {/* Actions */}
        <div className="flex flex-col md:flex-row items-center justify-center gap-4 mt-8">
          <Button
            variant="outline"
            size="lg"
            onClick={onRetry}
            className="flex items-center gap-2"
          >
            <RotateCcw className="h-5 w-5" />
            Try Another Interview
          </Button>
          <Button
            size="lg"
            onClick={() => window.location.href = '/interviews/practice'}
            className="flex items-center gap-2 bg-gradient-to-r from-primary to-blue-600"
          >
            <TrendingUp className="h-5 w-5" />
            Continue Practicing
          </Button>
        </div>

        {/* Tips */}
        <Card className="mt-8 border-none bg-gradient-to-r from-blue-500/5 to-purple-500/5">
          <CardContent className="p-6 text-center">
            <h3 className="text-xl font-bold mb-3">Keep Practicing!</h3>
            <p className="text-muted-foreground">
              Regular practice helps you build confidence and improve your interview skills.
              Try different interview types and difficulty levels to challenge yourself.
            </p>
          </CardContent>
        </Card>
      </div>
    </div>
  );
}
