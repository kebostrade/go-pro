"use client";

import { useState } from "react";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Progress } from "@/components/ui/progress";
import { Textarea } from "@/components/ui/textarea";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import {
  Send,
  CheckCircle,
  AlertCircle,
  Clock,
  Trophy,
  Star,
  MessageSquare,
  ThumbsUp,
  ThumbsDown,
  Lightbulb,
  Target,
  Award,
  TrendingUp,
  Zap,
  Eye,
  History,
  HelpCircle,
  ChevronDown,
  ChevronUp,
  Sparkles,
  BookOpen,
  Code2,
  Cpu,
  HardDrive
} from "lucide-react";
import { executeGoCode, analyzeCode } from "@/lib/code-execution";

interface Feedback {
  score: number;
  maxScore: number;
  passed: boolean;
  comments: string[];
  suggestions: string[];
  strengths: string[];
}

interface Hint {
  id: string;
  level: number;
  title: string;
  content: string;
  unlocked: boolean;
}

interface ExerciseSubmissionProps {
  exerciseId: string;
  title: string;
  description: string;
  requirements: string[];
  code: string;
  onSubmit?: (code: string, notes?: string) => Promise<Feedback>;
  previousSubmissions?: Array<{
    id: string;
    timestamp: string;
    score: number;
    feedback: string;
  }>;
  hints?: Hint[];
  difficulty?: 'easy' | 'medium' | 'hard';
  estimatedTime?: string;
  points?: number;
}

const ExerciseSubmission = ({
  exerciseId,
  title,
  description,
  requirements,
  code,
  onSubmit,
  previousSubmissions = [],
  hints = [],
  difficulty = 'medium',
  estimatedTime = '10-15 min',
  points = 100
}: ExerciseSubmissionProps) => {
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [feedback, setFeedback] = useState<Feedback | null>(null);
  const [notes, setNotes] = useState("");
  const [showPreviousSubmissions, setShowPreviousSubmissions] = useState(false);
  const [unlockedHints, setUnlockedHints] = useState<Set<string>>(new Set());
  const [showHints, setShowHints] = useState(false);
  const [attempts, setAttempts] = useState(0);
  const [startTime] = useState(Date.now());
  const [timeSpent, setTimeSpent] = useState(0);
  const [showRequirements, setShowRequirements] = useState(true);

  const handleSubmit = async () => {
    setAttempts(prev => prev + 1);

    if (!onSubmit) {
      // Enhanced demo feedback with code analysis
      setIsSubmitting(true);

      try {
        // Analyze the code
        const analysis = analyzeCode(code);
        const executionResult = await executeGoCode(code);

        // Calculate score based on multiple factors
        let score = 50; // Base score

        // Code structure (20 points)
        if (analysis.hasMain || !code.includes('package main')) score += 10;
        if (analysis.functionCount > 0) score += 10;

        // Complexity (15 points)
        if (analysis.complexity <= 10) score += 15;
        else if (analysis.complexity <= 15) score += 10;
        else score += 5;

        // Execution (15 points)
        if (!executionResult.error) score += 15;

        const passed = score >= 75;

        const comments: string[] = [];
        const suggestions: string[] = [];
        const strengths: string[] = [];

        // Generate feedback based on analysis
        if (passed) {
          comments.push("Excellent work! Your solution meets the requirements.");
          strengths.push("Clean and well-structured code");
        } else {
          comments.push("Good effort! There are areas for improvement.");
        }

        if (analysis.complexity > 15) {
          suggestions.push(`High complexity detected (${analysis.complexity}). Consider refactoring.`);
        } else {
          strengths.push("Good code complexity management");
        }

        if (executionResult.error) {
          comments.push(`Execution error: ${executionResult.error}`);
          suggestions.push("Fix the runtime errors before resubmitting");
        } else {
          strengths.push("Code executes without errors");
        }

        if (executionResult.warnings && executionResult.warnings.length > 0) {
          executionResult.warnings.forEach(warning => {
            suggestions.push(warning);
          });
        }

        if (executionResult.suggestions && executionResult.suggestions.length > 0) {
          suggestions.push(...executionResult.suggestions);
        }

        if (analysis.functionCount === 0 && code.includes('func')) {
          suggestions.push("Consider breaking down your code into functions");
        }

        setFeedback({
          score: Math.round(score),
          maxScore: 100,
          passed,
          comments,
          suggestions,
          strengths
        });

        // Unlock hints after failed attempts
        if (!passed && attempts >= 2 && hints.length > 0) {
          const newUnlocked = new Set(unlockedHints);
          hints.slice(0, Math.min(attempts - 1, hints.length)).forEach(hint => {
            newUnlocked.add(hint.id);
          });
          setUnlockedHints(newUnlocked);
        }
      } catch (error) {
        console.error("Analysis failed:", error);
      } finally {
        setIsSubmitting(false);
      }
      return;
    }

    setIsSubmitting(true);
    try {
      const result = await onSubmit(code, notes);
      setFeedback(result);
    } catch (error) {
      console.error("Submission failed:", error);
    } finally {
      setIsSubmitting(false);
    }
  };

  const unlockHint = (hintId: string) => {
    const newUnlocked = new Set(unlockedHints);
    newUnlocked.add(hintId);
    setUnlockedHints(newUnlocked);
  };

  const getDifficultyColor = (diff: string) => {
    switch (diff) {
      case 'easy': return 'text-green-500 bg-green-500/10 border-green-500/20';
      case 'medium': return 'text-yellow-500 bg-yellow-500/10 border-yellow-500/20';
      case 'hard': return 'text-red-500 bg-red-500/10 border-red-500/20';
      default: return 'text-blue-500 bg-blue-500/10 border-blue-500/20';
    }
  };

  const getScoreColor = (score: number, maxScore: number) => {
    const percentage = (score / maxScore) * 100;
    if (percentage >= 90) return "text-green-500";
    if (percentage >= 70) return "text-blue-500";
    if (percentage >= 50) return "text-yellow-500";
    return "text-red-500";
  };

  const getScoreBadge = (score: number, maxScore: number) => {
    const percentage = (score / maxScore) * 100;
    if (percentage >= 90) return { text: "Excellent", color: "bg-green-500 text-white" };
    if (percentage >= 70) return { text: "Good", color: "bg-blue-500 text-white" };
    if (percentage >= 50) return { text: "Needs Work", color: "bg-yellow-500 text-white" };
    return { text: "Poor", color: "bg-red-500 text-white" };
  };

  return (
    <div className="space-y-6">
      {/* Enhanced Exercise Info */}
      <Card className="glass-card border-2">
        <CardHeader>
          <div className="flex items-start justify-between mb-4">
            <div className="flex-1">
              <CardTitle className="text-xl flex items-center">
                <Trophy className="mr-2 h-5 w-5 text-primary" />
                {title}
              </CardTitle>
              <CardDescription className="text-base mt-2">
                {description}
              </CardDescription>
            </div>

            {/* Exercise Stats */}
            <div className="flex items-center space-x-2">
              <Badge variant="outline" className={`${getDifficultyColor(difficulty)} border`}>
                {difficulty.charAt(0).toUpperCase() + difficulty.slice(1)}
              </Badge>
              <Badge variant="outline" className="text-xs">
                <Clock className="mr-1 h-3 w-3" />
                {estimatedTime}
              </Badge>
              <Badge variant="outline" className="text-xs">
                <Star className="mr-1 h-3 w-3" />
                {points} pts
              </Badge>
            </div>
          </div>

          {/* Progress Indicators */}
          <div className="grid grid-cols-3 gap-4 mb-4">
            <div className="text-center p-3 bg-muted/50 rounded-lg">
              <div className="text-lg font-bold">{attempts}</div>
              <div className="text-xs text-muted-foreground">Attempts</div>
            </div>
            <div className="text-center p-3 bg-muted/50 rounded-lg">
              <div className="text-lg font-bold">{Math.floor(timeSpent / 60)}m</div>
              <div className="text-xs text-muted-foreground">Time Spent</div>
            </div>
            <div className="text-center p-3 bg-muted/50 rounded-lg">
              <div className="text-lg font-bold">{unlockedHints.size}/{hints.length}</div>
              <div className="text-xs text-muted-foreground">Hints Used</div>
            </div>
          </div>
        </CardHeader>

        <CardContent>
          <div className="space-y-4">
            {/* Collapsible Requirements */}
            <div>
              <Button
                variant="ghost"
                onClick={() => setShowRequirements(!showRequirements)}
                className="w-full justify-between p-0 h-auto font-medium mb-2"
              >
                <span className="flex items-center">
                  <Target className="mr-2 h-4 w-4" />
                  Requirements ({requirements.length})
                </span>
                {showRequirements ? <ChevronUp className="h-4 w-4" /> : <ChevronDown className="h-4 w-4" />}
              </Button>

              {showRequirements && (
                <ul className="space-y-2 animate-in slide-in-from-top-2 duration-200">
                  {requirements.map((req, index) => (
                    <li key={`req-${index}`} className="flex items-start space-x-3 p-2 rounded-lg hover:bg-muted/50 transition-colors">
                      <div className="w-2 h-2 rounded-full bg-primary mt-2 flex-shrink-0" />
                      <span className="text-sm leading-relaxed">{req}</span>
                    </li>
                  ))}
                </ul>
              )}
            </div>

            {/* Hints System */}
            {hints.length > 0 && (
              <div>
                <Button
                  variant="ghost"
                  onClick={() => setShowHints(!showHints)}
                  className="w-full justify-between p-0 h-auto font-medium mb-2"
                >
                  <span className="flex items-center">
                    <Lightbulb className="mr-2 h-4 w-4" />
                    Hints ({unlockedHints.size}/{hints.length} unlocked)
                  </span>
                  {showHints ? <ChevronUp className="h-4 w-4" /> : <ChevronDown className="h-4 w-4" />}
                </Button>

                {showHints && (
                  <div className="space-y-2 animate-in slide-in-from-top-2 duration-200">
                    {hints.map((hint, index) => (
                      <div
                        key={hint.id}
                        className={`p-3 rounded-lg border transition-all ${
                          unlockedHints.has(hint.id)
                            ? 'bg-blue-50 border-blue-200 dark:bg-blue-950 dark:border-blue-800'
                            : 'bg-muted/50 border-muted'
                        }`}
                      >
                        <div className="flex items-center justify-between mb-2">
                          <span className="text-sm font-medium flex items-center">
                            <HelpCircle className="mr-2 h-4 w-4" />
                            Hint {index + 1}: {hint.title}
                          </span>
                          {!unlockedHints.has(hint.id) && (
                            <Button
                              size="sm"
                              variant="outline"
                              onClick={() => unlockHint(hint.id)}
                              className="text-xs"
                            >
                              Unlock
                            </Button>
                          )}
                        </div>
                        {unlockedHints.has(hint.id) ? (
                          <p className="text-sm text-muted-foreground">{hint.content}</p>
                        ) : (
                          <p className="text-sm text-muted-foreground italic">Click "Unlock" to reveal this hint</p>
                        )}
                      </div>
                    ))}
                  </div>
                )}
              </div>
            )}
          </div>
        </CardContent>
      </Card>

      {/* Submission Form */}
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center">
            <Send className="mr-2 h-5 w-5" />
            Submit Your Solution
          </CardTitle>
          <CardDescription>
            Submit your code for automated review and feedback
          </CardDescription>
        </CardHeader>
        <CardContent className="space-y-4">
          <div>
            <label className="text-sm font-medium mb-2 block">
              Additional Notes (Optional)
            </label>
            <Textarea
              placeholder="Add any notes about your implementation, challenges faced, or questions..."
              value={notes}
              onChange={(e) => setNotes(e.target.value)}
              className="min-h-[100px]"
            />
          </div>

          <div className="flex items-center justify-between">
            <div className="text-sm text-muted-foreground">
              Code will be automatically evaluated against test cases
            </div>
            <Button
              onClick={handleSubmit}
              disabled={isSubmitting || !code.trim()}
              className="go-gradient text-white"
            >
              {isSubmitting ? (
                <>
                  <Clock className="mr-2 h-4 w-4 animate-spin" />
                  Submitting...
                </>
              ) : (
                <>
                  <Send className="mr-2 h-4 w-4" />
                  Submit Solution
                </>
              )}
            </Button>
          </div>
        </CardContent>
      </Card>

      {/* Feedback */}
      {feedback && (
        <Card className="border-primary/20 bg-primary/5">
          <CardHeader>
            <div className="flex items-center justify-between">
              <CardTitle className="flex items-center">
                {feedback.passed ? (
                  <CheckCircle className="mr-2 h-5 w-5 text-green-500" />
                ) : (
                  <AlertCircle className="mr-2 h-5 w-5 text-red-500" />
                )}
                Submission Feedback
              </CardTitle>
              <div className="flex items-center space-x-2">
                <Badge className={getScoreBadge(feedback.score, feedback.maxScore).color}>
                  {getScoreBadge(feedback.score, feedback.maxScore).text}
                </Badge>
                <div className={`text-2xl font-bold ${getScoreColor(feedback.score, feedback.maxScore)}`}>
                  {feedback.score}/{feedback.maxScore}
                </div>
              </div>
            </div>
          </CardHeader>
          <CardContent className="space-y-6">
            {/* Score Progress */}
            <div className="space-y-2">
              <div className="flex justify-between text-sm">
                <span>Score</span>
                <span>{Math.round((feedback.score / feedback.maxScore) * 100)}%</span>
              </div>
              <Progress value={(feedback.score / feedback.maxScore) * 100} className="h-2" />
            </div>

            {/* Feedback Sections */}
            <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
              {/* Comments */}
              <div className="space-y-3">
                <h4 className="font-medium flex items-center">
                  <MessageSquare className="mr-2 h-4 w-4" />
                  Comments
                </h4>
                <div className="space-y-2">
                  {feedback.comments.map((comment, index) => (
                    <div key={index} className="text-sm p-3 bg-background rounded-lg border">
                      {comment}
                    </div>
                  ))}
                </div>
              </div>

              {/* Suggestions */}
              <div className="space-y-3">
                <h4 className="font-medium flex items-center">
                  <ThumbsUp className="mr-2 h-4 w-4" />
                  Suggestions
                </h4>
                <div className="space-y-2">
                  {feedback.suggestions.map((suggestion, index) => (
                    <div key={index} className="text-sm p-3 bg-yellow-50 dark:bg-yellow-950 rounded-lg border border-yellow-200 dark:border-yellow-800">
                      {suggestion}
                    </div>
                  ))}
                </div>
              </div>

              {/* Strengths */}
              <div className="space-y-3">
                <h4 className="font-medium flex items-center">
                  <Star className="mr-2 h-4 w-4" />
                  Strengths
                </h4>
                <div className="space-y-2">
                  {feedback.strengths.map((strength, index) => (
                    <div key={index} className="text-sm p-3 bg-green-50 dark:bg-green-950 rounded-lg border border-green-200 dark:border-green-800">
                      {strength}
                    </div>
                  ))}
                </div>
              </div>
            </div>

            {/* Action Buttons */}
            <div className="flex items-center justify-between pt-4 border-t">
              <div className="flex items-center space-x-2">
                {feedback.passed && (
                  <Badge className="bg-green-500 text-white">
                    <Trophy className="mr-1 h-3 w-3" />
                    Exercise Completed!
                  </Badge>
                )}
              </div>
              <div className="flex items-center space-x-2">
                <Button variant="outline" size="sm">
                  <ThumbsUp className="mr-2 h-3 w-3" />
                  Helpful
                </Button>
                <Button variant="outline" size="sm">
                  <ThumbsDown className="mr-2 h-3 w-3" />
                  Not Helpful
                </Button>
              </div>
            </div>
          </CardContent>
        </Card>
      )}

      {/* Previous Submissions */}
      {previousSubmissions.length > 0 && (
        <Card>
          <CardHeader>
            <div className="flex items-center justify-between">
              <CardTitle className="text-lg">Previous Submissions</CardTitle>
              <Button
                variant="outline"
                size="sm"
                onClick={() => setShowPreviousSubmissions(!showPreviousSubmissions)}
              >
                {showPreviousSubmissions ? "Hide" : "Show"} History
              </Button>
            </div>
          </CardHeader>
          {showPreviousSubmissions && (
            <CardContent>
              <div className="space-y-3">
                {previousSubmissions.map((submission) => (
                  <div
                    key={submission.id}
                    className="flex items-center justify-between p-3 bg-muted/50 rounded-lg"
                  >
                    <div>
                      <div className="font-medium text-sm">
                        Submission {submission.id}
                      </div>
                      <div className="text-xs text-muted-foreground">
                        {submission.timestamp}
                      </div>
                    </div>
                    <div className="flex items-center space-x-2">
                      <Badge variant="outline">
                        Score: {submission.score}
                      </Badge>
                    </div>
                  </div>
                ))}
              </div>
            </CardContent>
          )}
        </Card>
      )}
    </div>
  );
};

export default ExerciseSubmission;
