'use client';

import { useState, useEffect, useCallback } from 'react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
import PlaygroundEditor from '@/components/playground/PlaygroundEditor';
import QuestionCard from './QuestionCard';
import {
  Clock,
  Send,
  ArrowLeft,
  ArrowRight,
  Save,
  AlertCircle,
  CheckCircle,
  Code2,
  MessageSquare,
  Network,
  Loader2,
} from 'lucide-react';
import { api } from '@/lib/api';
import type { InterviewType, Difficulty, Question, Answer } from '@/types/interview';

interface InterviewSessionProps {
  type: InterviewType;
  difficulty: Difficulty;
  onComplete?: (answers: Answer[]) => void;
}

// Mock questions for fallback when API unavailable
const getMockQuestions = (type: InterviewType, difficulty: Difficulty): Question[] => {
  const questions: Record<InterviewType, Record<Difficulty, Question[]>> = {
    coding: {
      beginner: [
        {
          id: 'c1',
          content: 'Write a function that takes a slice of integers and returns the sum of all even numbers.',
          type: 'coding',
          difficulty: 'beginner',
          time_limit: 300,
          expected_points: ['Use a loop to iterate through the slice', 'Check if each number is even using modulo operator', 'Accumulate the sum in a variable'],
        },
        {
          id: 'c2',
          content: 'Write a function that reverses a string without using any built-in reverse functions.',
          type: 'coding',
          difficulty: 'beginner',
          time_limit: 300,
          expected_points: ['Convert string to slice of runes', 'Build a new string in reverse order', 'Handle empty strings gracefully'],
        },
        {
          id: 'c3',
          content: 'Write a function that finds the maximum value in a slice of integers.',
          type: 'coding',
          difficulty: 'beginner',
          time_limit: 300,
          expected_points: ['Initialize max with first element', 'Compare each element with current max', 'Update max when larger value found'],
        },
      ],
      intermediate: [
        {
          id: 'c4',
          content: 'Implement a function that merges two sorted slices into one sorted slice.',
          type: 'coding',
          difficulty: 'intermediate',
          time_limit: 600,
          expected_points: ['Use two pointers technique', 'Compare elements from both slices', 'Handle edge cases like empty slices'],
        },
        {
          id: 'c5',
          content: 'Write a function that detects if a linked list has a cycle.',
          type: 'coding',
          difficulty: 'intermediate',
          time_limit: 600,
          expected_points: ['Use slow and fast pointers', 'Slow pointer moves one step', 'Fast pointer moves two steps'],
        },
      ],
      advanced: [
        {
          id: 'c6',
          content: 'Design and implement a thread-safe LRU cache in Go.',
          type: 'coding',
          difficulty: 'advanced',
          time_limit: 900,
          expected_points: ['Use sync.RWMutex for thread safety', 'Implement doubly-linked list for cache order', 'Use hash map for O(1) lookups'],
        },
      ],
    },
    behavioral: {
      beginner: [
        {
          id: 'b1',
          content: 'Tell me about a time when you had to work with a difficult team member. How did you handle the situation?',
          type: 'behavioral',
          difficulty: 'beginner',
          time_limit: 300,
          expected_points: ['Show empathy and understanding', 'Focus on solutions, not blame', 'Mention the outcome'],
        },
        {
          id: 'b2',
          content: 'Describe a project where you took initiative without being asked. What was the impact?',
          type: 'behavioral',
          difficulty: 'beginner',
          time_limit: 300,
          expected_points: ['Be specific about your role', 'Quantify the impact', 'Show proactive mindset'],
        },
      ],
      intermediate: [
        {
          id: 'b3',
          content: 'Tell me about a time you had to make a difficult technical decision with limited information. What was your process?',
          type: 'behavioral',
          difficulty: 'intermediate',
          time_limit: 600,
          expected_points: ['Explain the constraints', 'Show your decision-making framework', 'Reflect on the outcome'],
        },
      ],
      advanced: [
        {
          id: 'b4',
          content: 'Describe a situation where you had to influence stakeholders without authority. How did you approach it?',
          type: 'behavioral',
          difficulty: 'advanced',
          time_limit: 900,
          expected_points: ['Demonstrate influence skills', 'Show understanding of stakeholder needs', 'Highlight the outcome'],
        },
      ],
    },
    system_design: {
      beginner: [
        {
          id: 's1',
          content: 'Design a URL shortener service. Explain your approach for storing and retrieving URLs.',
          type: 'system_design',
          difficulty: 'beginner',
          time_limit: 300,
          expected_points: ['Consider database choices', 'Handle collisions', 'Discuss scalability'],
        },
      ],
      intermediate: [
        {
          id: 's2',
          content: 'Design a real-time chat application. How would you handle message delivery and persistence?',
          type: 'system_design',
          difficulty: 'intermediate',
          time_limit: 600,
          expected_points: ['Consider WebSocket connections', 'Message queue architecture', 'Scaling strategies'],
        },
      ],
      advanced: [
        {
          id: 's3',
          content: 'Design a distributed key-value store that can handle billions of requests per day.',
          type: 'system_design',
          difficulty: 'advanced',
          time_limit: 900,
          expected_points: ['Consistent vs eventual consistency', 'Sharding strategy', 'Fault tolerance mechanisms'],
        },
      ],
    },
  };

  return questions[type]?.[difficulty] || questions.coding.beginner;
};

export default function InterviewSession({ type, difficulty, onComplete }: InterviewSessionProps) {
  const [questions, setQuestions] = useState<Question[]>([]);
  const [currentIndex, setCurrentIndex] = useState(0);
  const [answers, setAnswers] = useState<Answer[]>([]);
  const [timeRemaining, setTimeRemaining] = useState(0);
  const [hintsUsed, setHintsUsed] = useState(0);
  const [isLoading, setIsLoading] = useState(true);
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [currentAnswer, setCurrentAnswer] = useState('');
  const [sessionId, setSessionId] = useState<string | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [useMockData, setUseMockData] = useState(false);

  const maxHints = 2;

  // Initialize interview session
  useEffect(() => {
    const initSession = async () => {
      setIsLoading(true);
      setError(null);

      try {
        // Try to start interview via API
        const response = await api.startInterview(type, difficulty);
        setSessionId(response.session.id);
        setQuestions(response.session.questions.map(q => ({
          ...q,
          type: q.type as InterviewType,
          difficulty: q.difficulty as Difficulty,
        })));
        setTimeRemaining(response.session.questions[0]?.time_limit || 600);
        setUseMockData(false);
      } catch (err) {
        // Fallback to mock data if API unavailable
        console.log('API unavailable, using mock data:', err);
        const mockQuestions = getMockQuestions(type, difficulty);
        setQuestions(mockQuestions);
        setTimeRemaining(mockQuestions[0]?.time_limit || 600);
        setUseMockData(true);
      } finally {
        setIsLoading(false);
      }
    };

    initSession();
  }, [type, difficulty]);

  // Timer effect
  useEffect(() => {
    if (timeRemaining <= 0 || isLoading) return;

    const timer = setInterval(() => {
      setTimeRemaining((prev) => {
        if (prev <= 1) {
          // Time's up - auto-submit current answer
          handleNext(true);
          return 0;
        }
        return prev - 1;
      });
    }, 1000);

    return () => clearInterval(timer);
  }, [isLoading]);

  const handleNext = useCallback(async (autoSubmit = false) => {
    if (isSubmitting) return;

    const answerContent = currentAnswer.trim();
    if (!answerContent && !autoSubmit) {
      setError('Please provide an answer before continuing.');
      return;
    }

    setIsSubmitting(true);
    setError(null);

    const newAnswer: Answer = {
      question_id: questions[currentIndex]?.id || '',
      content: answerContent,
      created_at: new Date().toISOString(),
    };

    const newAnswers = [...answers, newAnswer];
    setAnswers(newAnswers);

    // Try to submit to API if available
    if (sessionId && !useMockData) {
      try {
        await api.submitInterviewAnswer(sessionId, answerContent);
      } catch (err) {
        console.error('Failed to submit answer to API:', err);
      }
    }

    if (currentIndex < questions.length - 1) {
      setCurrentIndex(currentIndex + 1);
      setCurrentAnswer('');
      setTimeRemaining(questions[currentIndex + 1]?.time_limit || 600);
      setHintsUsed(0);
    } else {
      // Interview complete
      onComplete?.(newAnswers);
      window.location.href = `/interviews/feedback?type=${type}&difficulty=${difficulty}`;
    }

    setIsSubmitting(false);
  }, [currentIndex, questions, currentAnswer, answers, sessionId, useMockData, type, difficulty, onComplete, isSubmitting]);

  const handlePrevious = useCallback(() => {
    if (currentIndex > 0) {
      setCurrentIndex(currentIndex - 1);
      setCurrentAnswer(answers[currentIndex - 1]?.content || '');
      setTimeRemaining(questions[currentIndex - 1]?.time_limit || 600);
    }
  }, [currentIndex, answers, questions]);

  const handleHint = useCallback(() => {
    if (hintsUsed < maxHints) {
      setHintsUsed((prev) => prev + 1);
    }
  }, [hintsUsed, maxHints]);

  const handleComplete = useCallback(() => {
    const allAnswers = [...answers];
    if (currentAnswer.trim()) {
      allAnswers.push({
        question_id: questions[currentIndex]?.id || '',
        content: currentAnswer,
        created_at: new Date().toISOString(),
      });
    }
    onComplete?.(allAnswers);
    window.location.href = `/interviews/feedback?type=${type}&difficulty=${difficulty}`;
  }, [answers, currentAnswer, currentIndex, questions, type, difficulty, onComplete]);

  if (isLoading) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="text-center space-y-4">
          <Loader2 className="h-12 w-12 animate-spin text-primary mx-auto" />
          <p className="text-muted-foreground">Loading interview questions...</p>
        </div>
      </div>
    );
  }

  const currentQuestion = questions[currentIndex];
  if (!currentQuestion) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <Card className="border-none shadow-lg">
          <CardContent className="p-6 text-center space-y-4">
            <AlertCircle className="h-16 w-16 text-destructive mx-auto" />
            <h2 className="text-xl font-bold">No questions available</h2>
            <p className="text-muted-foreground">
              Unable to load interview questions. Please try again.
            </p>
            <Button onClick={() => window.location.href = '/interviews/practice'}>
              Back to Practice
            </Button>
          </CardContent>
        </Card>
      </div>
    );
  }

  const TypeIcon = {
    coding: Code2,
    behavioral: MessageSquare,
    system_design: Network,
  }[type];

  return (
    <div className="min-h-screen bg-gradient-to-br from-background via-background to-background">
      <div className="container mx-auto px-4 py-6 max-w-7xl">
        {/* API Mode Indicator */}
        {useMockData && (
          <div className="mb-4 p-3 bg-amber-500/10 border border-amber-500/30 rounded-lg text-amber-600 text-sm flex items-center gap-2">
            <AlertCircle className="h-4 w-4" />
            Running in demo mode (backend unavailable)
          </div>
        )}

        {/* Header */}
        <div className="mb-6">
          <div className="flex items-center justify-between">
            <div>
              <h1 className="text-3xl font-bold tracking-tight flex items-center gap-3">
                <TypeIcon className="h-8 w-8 text-primary" />
                {type.replace('_', ' ').replace(/\b\w/g, (match) => match.toUpperCase())} Interview
              </h1>
              <p className="text-muted-foreground mt-1 capitalize">
                {difficulty} Level
              </p>
            </div>
            <div className="flex items-center gap-2">
              <Badge variant="outline" className="text-sm">
                Question {currentIndex + 1} of {questions.length}
              </Badge>
            </div>
          </div>
        </div>

        {/* Error Message */}
        {error && (
          <div className="mb-4 p-3 bg-destructive/10 border border-destructive/30 rounded-lg text-destructive text-sm">
            {error}
          </div>
        )}

        <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
          {/* Question Card */}
          <div className="lg:col-span-2">
            <QuestionCard
              question={currentQuestion}
              currentIndex={currentIndex}
              totalQuestions={questions.length}
              timeRemaining={timeRemaining}
              hintsUsed={hintsUsed}
              maxHints={maxHints}
              onNext={() => handleNext(false)}
              onPrevious={handlePrevious}
              onHint={handleHint}
              onEnd={handleComplete}
            />
          </div>

          {/* Answer Input */}
          <div className="lg:col-span-1">
            <Card className="border-none shadow-lg h-full">
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <Send className="h-5 w-5 text-primary" />
                  Your Answer
                </CardTitle>
              </CardHeader>
              <CardContent className="space-y-4">
                {type === 'coding' ? (
                  <Tabs defaultValue="editor" className="w-full">
                    <TabsList className="grid w-full grid-cols-1">
                      <TabsTrigger value="editor">
                        <Code2 className="h-4 w-4 mr-2" />
                        Code Editor
                      </TabsTrigger>
                    </TabsList>
                    <TabsContent value="editor" className="mt-4">
                      <div className="rounded-lg overflow-hidden border border-border/50">
                        <PlaygroundEditor
                          value={currentAnswer}
                          onChange={(value) => setCurrentAnswer(value || '')}
                          language="go"
                          height="400px"
                        />
                      </div>
                    </TabsContent>
                  </Tabs>
                ) : (
                  <textarea
                    value={currentAnswer}
                    onChange={(e) => setCurrentAnswer(e.target.value)}
                    placeholder="Type your answer here..."
                    className="w-full min-h-[400px] rounded-lg border border-border/50 bg-card p-4 text-sm focus:outline-none focus:ring-2 focus:ring-primary/50 resize-y"
                  />
                )}

                {/* Action Buttons */}
                <div className="flex gap-2 pt-4">
                  <Button
                    variant="outline"
                    onClick={() => {
                      setAnswers((prev) => [
                        ...prev,
                        {
                          question_id: currentQuestion.id,
                          content: currentAnswer,
                          created_at: new Date().toISOString(),
                        },
                      ]);
                    }}
                  >
                    <Save className="h-4 w-4 mr-2" />
                    Save Draft
                  </Button>
                </div>

                {/* Progress */}
                <div className="pt-4 border-t border-border/50">
                  <div className="text-sm text-muted-foreground mb-2">
                    Progress
                  </div>
                  <div className="w-full bg-muted rounded-full h-2">
                    <div
                      className="bg-primary h-full rounded-full transition-all duration-500"
                      style={{ width: `${((currentIndex + 1) / questions.length) * 100}%` }}
                    />
                  </div>
                </div>
              </CardContent>
            </Card>
          </div>
        </div>
      </div>
    </div>
  );
}
