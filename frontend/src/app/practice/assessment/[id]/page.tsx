"use client";

import { useState, useEffect } from "react";
import { useParams, useRouter } from "next/navigation";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Progress } from "@/components/ui/progress";
import {
  ArrowLeft,
  Clock,
  Brain,
  CheckCircle,
  AlertCircle,
  Home,
  ChevronRight,
  Trophy,
  Target,
  BookOpen,
  ArrowRight,
  RotateCcw
} from "lucide-react";
import Link from "next/link";

interface Question {
  id: string;
  type: "multiple-choice" | "code-completion" | "true-false";
  question: string;
  options?: string[];
  correctAnswer: string | number;
  explanation: string;
  points: number;
}

interface AssessmentData {
  id: string;
  title: string;
  description: string;
  category: string;
  difficulty: "Beginner" | "Intermediate" | "Advanced";
  duration: string;
  totalQuestions: number;
  totalPoints: number;
  questions: Question[];
  timeLimit: number; // in minutes
}

export default function AssessmentPage() {
  const params = useParams();
  const router = useRouter();
  const [assessmentData, setAssessmentData] = useState<AssessmentData | null>(null);
  const [loading, setLoading] = useState(true);
  const [currentQuestionIndex, setCurrentQuestionIndex] = useState(0);
  const [answers, setAnswers] = useState<Record<string, string | number>>({});
  const [timeRemaining, setTimeRemaining] = useState(0);
  const [isStarted, setIsStarted] = useState(false);
  const [isCompleted, setIsCompleted] = useState(false);
  const [showResults, setShowResults] = useState(false);
  const [score, setScore] = useState(0);

  const assessmentId = params.id as string;

  useEffect(() => {
    // Mock data loading
    const mockAssessment: AssessmentData = {
      id: assessmentId,
      title: "Go Fundamentals Assessment",
      description: "Test your knowledge of Go basics, syntax, and core concepts",
      category: "Fundamentals",
      difficulty: "Beginner",
      duration: "30 min",
      totalQuestions: 10,
      totalPoints: 100,
      timeLimit: 30,
      questions: [
        {
          id: "q1",
          type: "multiple-choice",
          question: "Which of the following is the correct way to declare a variable in Go?",
          options: [
            "var x int = 5",
            "int x = 5",
            "x := 5",
            "Both A and C are correct"
          ],
          correctAnswer: 3,
          explanation: "In Go, you can declare variables using 'var' keyword or short declaration ':='. Both 'var x int = 5' and 'x := 5' are correct.",
          points: 10
        },
        {
          id: "q2",
          type: "multiple-choice",
          question: "What is the zero value of a string in Go?",
          options: [
            "null",
            "undefined",
            "\"\"",
            "0"
          ],
          correctAnswer: 2,
          explanation: "The zero value of a string in Go is an empty string \"\".",
          points: 10
        },
        {
          id: "q3",
          type: "true-false",
          question: "Go supports method overloading.",
          correctAnswer: "false",
          explanation: "Go does not support method overloading. Each method must have a unique name within its type.",
          points: 10
        },
        {
          id: "q4",
          type: "multiple-choice",
          question: "Which keyword is used to define a constant in Go?",
          options: [
            "const",
            "final",
            "static",
            "readonly"
          ],
          correctAnswer: 0,
          explanation: "The 'const' keyword is used to define constants in Go.",
          points: 10
        },
        {
          id: "q5",
          type: "multiple-choice",
          question: "What is the correct way to create a slice in Go?",
          options: [
            "var s []int",
            "s := make([]int, 0)",
            "s := []int{}",
            "All of the above"
          ],
          correctAnswer: 3,
          explanation: "All three methods are correct ways to create a slice in Go.",
          points: 10
        }
      ]
    };

    setAssessmentData(mockAssessment);
    setTimeRemaining(mockAssessment.timeLimit * 60); // Convert to seconds
    setLoading(false);
  }, [assessmentId]);

  // Timer effect
  useEffect(() => {
    if (isStarted && !isCompleted && timeRemaining > 0) {
      const timer = setInterval(() => {
        setTimeRemaining(prev => {
          if (prev <= 1) {
            handleSubmitAssessment();
            return 0;
          }
          return prev - 1;
        });
      }, 1000);

      return () => clearInterval(timer);
    }
  }, [isStarted, isCompleted, timeRemaining]);

  const handleStartAssessment = () => {
    setIsStarted(true);
  };

  const handleAnswerChange = (questionId: string, answer: string | number) => {
    setAnswers(prev => ({
      ...prev,
      [questionId]: answer
    }));
  };

  const handleNextQuestion = () => {
    if (assessmentData && currentQuestionIndex < assessmentData.questions.length - 1) {
      setCurrentQuestionIndex(prev => prev + 1);
    }
  };

  const handlePreviousQuestion = () => {
    if (currentQuestionIndex > 0) {
      setCurrentQuestionIndex(prev => prev - 1);
    }
  };

  const handleSubmitAssessment = () => {
    if (!assessmentData) return;

    let totalScore = 0;
    assessmentData.questions.forEach(question => {
      const userAnswer = answers[question.id];
      if (userAnswer === question.correctAnswer) {
        totalScore += question.points;
      }
    });

    setScore(totalScore);
    setIsCompleted(true);
    setShowResults(true);
  };

  const formatTime = (seconds: number) => {
    const minutes = Math.floor(seconds / 60);
    const remainingSeconds = seconds % 60;
    return `${minutes}:${remainingSeconds.toString().padStart(2, '0')}`;
  };

  const getDifficultyColor = (difficulty: string) => {
    switch (difficulty) {
      case "Beginner": return "text-green-600 bg-green-50 border-green-200";
      case "Intermediate": return "text-yellow-600 bg-yellow-50 border-yellow-200";
      case "Advanced": return "text-red-600 bg-red-50 border-red-200";
      default: return "text-gray-600 bg-gray-50 border-gray-200";
    }
  };

  if (loading || !assessmentData) {
    return (
      <div className="min-h-screen bg-gradient-to-br from-background via-background to-accent/5">
        <div className="container-responsive padding-responsive-y">
          <div className="flex items-center justify-center min-h-[60vh]">
            <div className="text-center">
              <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-primary mx-auto mb-4"></div>
              <p className="text-responsive text-muted-foreground">Loading assessment...</p>
            </div>
          </div>
        </div>
      </div>
    );
  }

  // Pre-start screen
  if (!isStarted) {
    return (
      <div className="min-h-screen bg-gradient-to-br from-background via-background to-accent/5">
        <div className="container-responsive padding-responsive-y">
        {/* Breadcrumb Navigation */}
        <div className="flex items-center space-x-2 text-sm text-muted-foreground mb-6">
          <Link href="/" className="hover:text-primary">
            <Home className="h-4 w-4" />
          </Link>
          <ChevronRight className="h-4 w-4" />
          <Link href="/practice" className="hover:text-primary">
            Practice
          </Link>
          <ChevronRight className="h-4 w-4" />
          <span className="text-foreground">{assessmentData.title}</span>
        </div>

        <div className="max-w-2xl mx-auto">
          <Card>
            <CardHeader className="text-center">
              <div className="flex justify-center mb-4">
                <div className="p-3 rounded-full bg-primary/10">
                  <Brain className="h-8 w-8 text-primary" />
                </div>
              </div>
              <CardTitle className="text-2xl mb-2">{assessmentData.title}</CardTitle>
              <CardDescription className="text-base">
                {assessmentData.description}
              </CardDescription>
              <div className="flex items-center justify-center gap-2 mt-4">
                <Badge className={getDifficultyColor(assessmentData.difficulty)}>
                  {assessmentData.difficulty}
                </Badge>
                <Badge variant="outline">{assessmentData.category}</Badge>
              </div>
            </CardHeader>
            <CardContent className="space-y-6">
              <div className="grid grid-cols-2 gap-4 text-center">
                <div className="p-4 bg-muted/50 rounded-lg">
                  <BookOpen className="h-6 w-6 mx-auto mb-2 text-blue-500" />
                  <div className="font-bold">{assessmentData.totalQuestions}</div>
                  <div className="text-sm text-muted-foreground">Questions</div>
                </div>
                <div className="p-4 bg-muted/50 rounded-lg">
                  <Clock className="h-6 w-6 mx-auto mb-2 text-orange-500" />
                  <div className="font-bold">{assessmentData.duration}</div>
                  <div className="text-sm text-muted-foreground">Time Limit</div>
                </div>
              </div>

              <div className="space-y-4">
                <h3 className="font-semibold">Instructions:</h3>
                <ul className="space-y-2 text-sm text-muted-foreground">
                  <li>• You have {assessmentData.duration} to complete this assessment</li>
                  <li>• Each question has a point value, total possible score is {assessmentData.totalPoints}</li>
                  <li>• You can navigate between questions using the Previous/Next buttons</li>
                  <li>• Make sure to submit your assessment before time runs out</li>
                  <li>• You can review your answers before submitting</li>
                </ul>
              </div>

              <div className="flex items-center justify-between pt-4">
                <Link href="/practice">
                  <Button variant="outline">
                    <ArrowLeft className="mr-2 h-4 w-4" />
                    Back to Practice
                  </Button>
                </Link>
                <Button onClick={handleStartAssessment} className="go-gradient text-white">
                  <Brain className="mr-2 h-4 w-4" />
                  Start Assessment
                </Button>
              </div>
            </CardContent>
          </Card>
        </div>
      </div>
      </div>
    );
  }

  // Results screen
  if (showResults) {
    const percentage = Math.round((score / assessmentData.totalPoints) * 100);
    const passed = percentage >= 70;

    return (
      <div className="min-h-screen bg-gradient-to-br from-background via-background to-accent/5">
        <div className="container-responsive padding-responsive-y">
          <div className="max-w-2xl mx-auto">
          <Card>
            <CardHeader className="text-center">
              <div className="flex justify-center mb-4">
                <div className={`p-3 rounded-full ${passed ? 'bg-green-100' : 'bg-red-100'}`}>
                  {passed ? (
                    <Trophy className="h-8 w-8 text-green-600" />
                  ) : (
                    <Target className="h-8 w-8 text-red-600" />
                  )}
                </div>
              </div>
              <CardTitle className="text-2xl mb-2">
                {passed ? 'Congratulations!' : 'Assessment Complete'}
              </CardTitle>
              <CardDescription>
                {passed 
                  ? 'You have successfully passed the assessment!' 
                  : 'You can retake this assessment to improve your score.'
                }
              </CardDescription>
            </CardHeader>
            <CardContent className="space-y-6">
              <div className="text-center">
                <div className="text-4xl font-bold text-primary mb-2">{percentage}%</div>
                <div className="text-muted-foreground">
                  {score} out of {assessmentData.totalPoints} points
                </div>
                <Progress value={percentage} className="h-3 mt-4" />
              </div>

              <div className="grid grid-cols-2 gap-4">
                <div className="text-center p-4 bg-muted/50 rounded-lg">
                  <CheckCircle className="h-6 w-6 mx-auto mb-2 text-green-500" />
                  <div className="font-bold">
                    {assessmentData.questions.filter(q => answers[q.id] === q.correctAnswer).length}
                  </div>
                  <div className="text-sm text-muted-foreground">Correct</div>
                </div>
                <div className="text-center p-4 bg-muted/50 rounded-lg">
                  <AlertCircle className="h-6 w-6 mx-auto mb-2 text-red-500" />
                  <div className="font-bold">
                    {assessmentData.questions.filter(q => answers[q.id] !== q.correctAnswer).length}
                  </div>
                  <div className="text-sm text-muted-foreground">Incorrect</div>
                </div>
              </div>

              <div className="flex items-center justify-between pt-4">
                <Link href="/practice">
                  <Button variant="outline">
                    <ArrowLeft className="mr-2 h-4 w-4" />
                    Back to Practice
                  </Button>
                </Link>
                <div className="flex gap-2">
                  <Button variant="outline" onClick={() => window.location.reload()}>
                    <RotateCcw className="mr-2 h-4 w-4" />
                    Retake
                  </Button>
                  <Link href="/practice">
                    <Button className="go-gradient text-white">
                      Continue Learning
                      <ArrowRight className="ml-2 h-4 w-4" />
                    </Button>
                  </Link>
                </div>
              </div>
            </CardContent>
          </Card>
          </div>
        </div>
      </div>
    );
  }

  const currentQuestion = assessmentData.questions[currentQuestionIndex];
  const progress = ((currentQuestionIndex + 1) / assessmentData.questions.length) * 100;

  return (
    <div className="min-h-screen bg-gradient-to-br from-background via-background to-accent/5">
      <div className="container-responsive padding-responsive-y">
        {/* Header with timer */}
        <div className="flex flex-col lg:flex-row lg:items-center lg:justify-between mb-6 lg:mb-8">
        <div className="flex items-center space-x-4">
          <h1 className="text-2xl font-bold">{assessmentData.title}</h1>
          <Badge variant="outline">
            Question {currentQuestionIndex + 1} of {assessmentData.questions.length}
          </Badge>
        </div>
        <div className="flex items-center space-x-2 text-lg font-mono">
          <Clock className="h-5 w-5 text-orange-500" />
          <span className={timeRemaining < 300 ? 'text-red-600' : 'text-foreground'}>
            {formatTime(timeRemaining)}
          </span>
        </div>
      </div>

      {/* Progress bar */}
      <div className="mb-6">
        <Progress value={progress} className="h-2" />
      </div>

      {/* Question */}
      <div className="max-w-3xl mx-auto">
        <Card>
          <CardHeader>
            <CardTitle className="text-xl">{currentQuestion.question}</CardTitle>
            <div className="flex items-center justify-between">
              <Badge variant="secondary">{currentQuestion.points} points</Badge>
              <Badge variant="outline">{currentQuestion.type}</Badge>
            </div>
          </CardHeader>
          <CardContent className="space-y-4">
            {currentQuestion.type === "multiple-choice" && currentQuestion.options && (
              <div className="space-y-3">
                {currentQuestion.options.map((option, index) => (
                  <label key={index} className="flex items-center space-x-3 p-3 border rounded-lg cursor-pointer hover:bg-muted/50">
                    <input
                      type="radio"
                      name={currentQuestion.id}
                      value={index}
                      checked={answers[currentQuestion.id] === index}
                      onChange={() => handleAnswerChange(currentQuestion.id, index)}
                      className="w-4 h-4 text-primary"
                    />
                    <span>{option}</span>
                  </label>
                ))}
              </div>
            )}

            {currentQuestion.type === "true-false" && (
              <div className="space-y-3">
                <label className="flex items-center space-x-3 p-3 border rounded-lg cursor-pointer hover:bg-muted/50">
                  <input
                    type="radio"
                    name={currentQuestion.id}
                    value="true"
                    checked={answers[currentQuestion.id] === "true"}
                    onChange={() => handleAnswerChange(currentQuestion.id, "true")}
                    className="w-4 h-4 text-primary"
                  />
                  <span>True</span>
                </label>
                <label className="flex items-center space-x-3 p-3 border rounded-lg cursor-pointer hover:bg-muted/50">
                  <input
                    type="radio"
                    name={currentQuestion.id}
                    value="false"
                    checked={answers[currentQuestion.id] === "false"}
                    onChange={() => handleAnswerChange(currentQuestion.id, "false")}
                    className="w-4 h-4 text-primary"
                  />
                  <span>False</span>
                </label>
              </div>
            )}
          </CardContent>
        </Card>

        {/* Navigation */}
        <div className="flex items-center justify-between mt-6">
          <Button
            variant="outline"
            onClick={handlePreviousQuestion}
            disabled={currentQuestionIndex === 0}
          >
            <ArrowLeft className="mr-2 h-4 w-4" />
            Previous
          </Button>

          <div className="flex gap-2">
            {currentQuestionIndex === assessmentData.questions.length - 1 ? (
              <Button onClick={handleSubmitAssessment} className="go-gradient text-white">
                Submit Assessment
                <CheckCircle className="ml-2 h-4 w-4" />
              </Button>
            ) : (
              <Button onClick={handleNextQuestion}>
                Next
                <ArrowRight className="ml-2 h-4 w-4" />
              </Button>
            )}
          </div>
        </div>
      </div>
      </div>
    </div>
  );
}
