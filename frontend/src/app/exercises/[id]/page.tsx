'use client';

import { useState, useEffect } from 'react';
import { useParams } from 'next/navigation';
import MonacoCodeEditor from '@/components/learning/monaco-code-editor';
import EditorErrorBoundary from '@/components/learning/editor-error-boundary';
import { ArrowLeft, BookOpen, Target, Clock } from 'lucide-react';
import Link from 'next/link';

// Mock exercise data (replace with API call)
interface Exercise {
  id: string;
  title: string;
  description: string;
  difficulty: 'beginner' | 'intermediate' | 'advanced';
  estimated_time: string;
  initial_code: string;
  instructions: string[];
  hints: string[];
  test_cases: Array<{
    name: string;
    input?: string;
    expected: string;
  }>;
}

const mockExercise: Exercise = {
  id: '1',
  title: 'FizzBuzz',
  description: 'Write a function that prints numbers from 1 to n. For multiples of 3, print "Fizz" instead of the number. For multiples of 5, print "Buzz". For multiples of both 3 and 5, print "FizzBuzz".',
  difficulty: 'beginner',
  estimated_time: '15 minutes',
  initial_code: `package main

import "fmt"

// FizzBuzz prints numbers from 1 to n with special rules:
// - Multiples of 3: print "Fizz"
// - Multiples of 5: print "Buzz"
// - Multiples of both: print "FizzBuzz"
func FizzBuzz(n int) {
	// TODO: Implement your solution here
}

func main() {
	FizzBuzz(15)
}
`,
  instructions: [
    'Create a loop from 1 to n',
    'Check if the number is divisible by both 3 and 5 first',
    'Then check for divisibility by 3',
    'Then check for divisibility by 5',
    'Otherwise, print the number',
  ],
  hints: [
    'Use the modulo operator (%) to check divisibility',
    'Remember to check for both 3 AND 5 before checking individually',
    'fmt.Println() will print to standard output',
  ],
  test_cases: [
    {
      name: 'FizzBuzz(15) output',
      expected: '1\n2\nFizz\n4\nBuzz\nFizz\n7\n8\nFizz\nBuzz\n11\nFizz\n13\n14\nFizzBuzz\n',
    },
    {
      name: 'FizzBuzz(5) output',
      expected: '1\n2\nFizz\n4\nBuzz\n',
    },
  ],
};

// Mock API call (replace with real API)
const submitExercise = async (exerciseId: string, code: string) => {
  // Simulate API delay
  await new Promise((resolve) => setTimeout(resolve, 1500));

  // Mock response
  return {
    success: true,
    passed: true,
    score: 100,
    results: [
      {
        name: 'FizzBuzz(15) output',
        passed: true,
        expected: '1\n2\nFizz\n4\nBuzz\nFizz\n7\n8\nFizz\nBuzz\n11\nFizz\n13\n14\nFizzBuzz\n',
        actual: '1\n2\nFizz\n4\nBuzz\nFizz\n7\n8\nFizz\nBuzz\n11\nFizz\n13\n14\nFizzBuzz\n',
        execution_time_ms: 2,
      },
      {
        name: 'FizzBuzz(5) output',
        passed: true,
        expected: '1\n2\nFizz\n4\nBuzz\n',
        actual: '1\n2\nFizz\n4\nBuzz\n',
        execution_time_ms: 1,
      },
    ],
    execution_time_ms: 3,
    message: 'All tests passed! Great work!',
  };
};

export default function ExercisePage() {
  const params = useParams();
  const exerciseId = params?.id as string;
  const [exercise, setExercise] = useState<Exercise | null>(null);
  const [showHints, setShowHints] = useState(false);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    // Simulate loading exercise data
    setTimeout(() => {
      setExercise(mockExercise);
      setLoading(false);
    }, 500);
  }, [exerciseId]);

  if (loading) {
    return (
      <div className="flex items-center justify-center min-h-screen">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
      </div>
    );
  }

  if (!exercise) {
    return (
      <div className="flex items-center justify-center min-h-screen">
        <div className="text-center">
          <h1 className="text-2xl font-bold text-gray-900 mb-2">Exercise not found</h1>
          <Link href="/exercises" className="text-blue-600 hover:underline">
            Back to exercises
          </Link>
        </div>
      </div>
    );
  }

  const difficultyColors = {
    beginner: 'bg-green-100 text-green-800',
    intermediate: 'bg-yellow-100 text-yellow-800',
    advanced: 'bg-red-100 text-red-800',
  };

  return (
    <div className="min-h-screen bg-gray-50">
      {/* Header */}
      <div className="bg-white border-b border-gray-200">
        <div className="max-w-7xl mx-auto px-4 py-4">
          <Link
            href="/exercises"
            className="inline-flex items-center gap-2 text-sm text-gray-600 hover:text-gray-900 mb-4"
          >
            <ArrowLeft className="w-4 h-4" />
            Back to exercises
          </Link>
          <div className="flex items-start justify-between">
            <div>
              <h1 className="text-3xl font-bold text-gray-900 mb-2">{exercise.title}</h1>
              <div className="flex items-center gap-4 text-sm">
                <span className={`px-2 py-1 rounded-full font-medium ${difficultyColors[exercise.difficulty]}`}>
                  {exercise.difficulty.charAt(0).toUpperCase() + exercise.difficulty.slice(1)}
                </span>
                <span className="flex items-center gap-1 text-gray-600">
                  <Clock className="w-4 h-4" />
                  {exercise.estimated_time}
                </span>
              </div>
            </div>
          </div>
        </div>
      </div>

      {/* Main Content */}
      <div className="max-w-7xl mx-auto px-4 py-8">
        <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
          {/* Left Column - Instructions */}
          <div className="lg:col-span-1 space-y-6">
            {/* Description */}
            <div className="bg-white rounded-lg border border-gray-200 p-6">
              <div className="flex items-center gap-2 mb-4">
                <BookOpen className="w-5 h-5 text-blue-600" />
                <h2 className="text-lg font-semibold text-gray-900">Description</h2>
              </div>
              <p className="text-gray-700">{exercise.description}</p>
            </div>

            {/* Instructions */}
            <div className="bg-white rounded-lg border border-gray-200 p-6">
              <div className="flex items-center gap-2 mb-4">
                <Target className="w-5 h-5 text-green-600" />
                <h2 className="text-lg font-semibold text-gray-900">Instructions</h2>
              </div>
              <ol className="list-decimal list-inside space-y-2">
                {exercise.instructions.map((instruction, index) => (
                  <li key={index} className="text-gray-700">
                    {instruction}
                  </li>
                ))}
              </ol>
            </div>

            {/* Hints */}
            <div className="bg-white rounded-lg border border-gray-200 p-6">
              <button
                onClick={() => setShowHints(!showHints)}
                className="flex items-center justify-between w-full text-left"
              >
                <h2 className="text-lg font-semibold text-gray-900">Hints</h2>
                <span className="text-sm text-blue-600">{showHints ? 'Hide' : 'Show'}</span>
              </button>
              {showHints && (
                <ul className="mt-4 space-y-2">
                  {exercise.hints.map((hint, index) => (
                    <li key={index} className="text-sm text-gray-600 bg-blue-50 rounded p-3">
                      💡 {hint}
                    </li>
                  ))}
                </ul>
              )}
            </div>
          </div>

          {/* Right Column - Code Editor */}
          <div className="lg:col-span-2">
            <div className="bg-white rounded-lg border border-gray-200 overflow-hidden">
              <EditorErrorBoundary>
                <MonacoCodeEditor
                  initialCode={exercise.initial_code}
                  exerciseId={exercise.id}
                  language="go"
                  height="600px"
                  onSubmit={(code) => submitExercise(exercise.id, code)}
                  testCases={exercise.test_cases}
                />
              </EditorErrorBoundary>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
