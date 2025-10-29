"use client";

import { useState, useEffect } from "react";
import { useParams, useRouter } from "next/navigation";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Progress } from "@/components/ui/progress";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import CodeEditor from "@/components/learning/code-editor";
import {
  ArrowLeft,
  Clock,
  Star,
  Target,
  CheckCircle,
  AlertCircle,
  Lightbulb,
  BookOpen,
  Code2,
  Play,
  RotateCcw,
  Home,
  ChevronRight,
  Trophy,
  Users,
  TrendingUp
} from "lucide-react";
import Link from "next/link";

interface ChallengeData {
  id: string;
  title: string;
  description: string;
  difficulty: "Beginner" | "Intermediate" | "Advanced";
  category: string;
  estimatedTime: string;
  points: number;
  tags: string[];
  instructions: string;
  initialCode: string;
  solution: string;
  testCases: TestCase[];
  hints: string[];
  completionRate: number;
  totalAttempts: number;
}

interface TestCase {
  id: string;
  input: string;
  expectedOutput: string;
  description: string;
  passed?: boolean;
}

export default function ChallengePage() {
  const params = useParams();
  const router = useRouter();
  const [activeTab, setActiveTab] = useState("challenge");
  const [challengeData, setChallengeData] = useState<ChallengeData | null>(null);
  const [loading, setLoading] = useState(true);
  const [currentCode, setCurrentCode] = useState("");
  const [testResults, setTestResults] = useState<TestCase[]>([]);
  const [showHints, setShowHints] = useState(false);
  const [attempts, setAttempts] = useState(0);
  const [completed, setCompleted] = useState(false);

  const challengeId = params.id as string;

  useEffect(() => {
    // Mock data loading
    const mockChallenge: ChallengeData = {
      id: challengeId,
      title: "Variables and Constants",
      description: "Practice declaring and using variables and constants in Go",
      difficulty: "Beginner",
      category: "Fundamentals",
      estimatedTime: "15 min",
      points: 50,
      tags: ["variables", "constants", "types"],
      instructions: `
# Variables and Constants Challenge

In this challenge, you'll practice working with Go's variable and constant declarations.

## Your Task
1. Declare a variable \`name\` of type string and assign it your name
2. Declare a constant \`pi\` with the value 3.14159
3. Create a variable \`age\` and assign it an integer value
4. Use type inference to create a variable \`isStudent\` with a boolean value
5. Print all variables using fmt.Printf

## Requirements
- Use both \`var\` and short variable declaration (\`:=\`)
- Include at least one constant declaration
- Use proper Go formatting and naming conventions
- Your code should compile and run without errors

## Example Output
\`\`\`
Name: John Doe
Pi: 3.14159
Age: 25
Is Student: true
\`\`\`
      `,
      initialCode: `package main

import "fmt"

func main() {
    // TODO: Declare your variables and constants here
    
    // TODO: Print the values using fmt.Printf
}`,
      solution: `package main

import "fmt"

func main() {
    // Variable declarations
    var name string = "John Doe"
    const pi = 3.14159
    var age int = 25
    isStudent := true
    
    // Print the values
    fmt.Printf("Name: %s\\n", name)
    fmt.Printf("Pi: %.5f\\n", pi)
    fmt.Printf("Age: %d\\n", age)
    fmt.Printf("Is Student: %t\\n", isStudent)
}`,
      testCases: [
        {
          id: "test1",
          input: "",
          expectedOutput: "Name: John Doe\nPi: 3.14159\nAge: 25\nIs Student: true",
          description: "Should declare and print all variables correctly"
        },
        {
          id: "test2",
          input: "",
          expectedOutput: "compilation success",
          description: "Code should compile without errors"
        }
      ],
      hints: [
        "Remember that constants are declared with the 'const' keyword",
        "You can use both 'var' keyword and short declaration ':=' for variables",
        "Use fmt.Printf with appropriate format specifiers: %s for strings, %d for integers, %t for booleans, %f for floats",
        "Make sure to import the 'fmt' package to use Printf"
      ],
      completionRate: 85,
      totalAttempts: 1247
    };

    setChallengeData(mockChallenge);
    setCurrentCode(mockChallenge.initialCode);
    setLoading(false);
  }, [challengeId]);

  const handleCodeChange = (code: string) => {
    setCurrentCode(code);
  };

  const handleRunCode = async (code: string) => {
    setAttempts(prev => prev + 1);

    // Mock test execution
    const mockResults: TestCase[] = challengeData?.testCases.map(test => ({
      ...test,
      passed: Math.random() > 0.3 // Mock random pass/fail
    })) || [];

    setTestResults(mockResults);

    const allPassed = mockResults.every(test => test.passed);
    if (allPassed) {
      setCompleted(true);
    }

    // Convert TestCase[] to TestResult[] for CodeEditor
    const testResults = mockResults.map(test => ({
      name: test.description,
      passed: test.passed || false,
      message: test.passed ? `Expected: ${test.expectedOutput}` : `Expected: ${test.expectedOutput}, but test failed`
    }));

    return {
      output: allPassed ? "All tests passed! Great job!" : "Some tests failed. Keep trying!",
      tests: testResults
    };
  };

  const handleReset = () => {
    if (challengeData) {
      setCurrentCode(challengeData.initialCode);
      setTestResults([]);
      setShowHints(false);
    }
  };

  const getDifficultyColor = (difficulty: string) => {
    switch (difficulty) {
      case "Beginner": return "text-green-600 bg-green-50 border-green-200";
      case "Intermediate": return "text-yellow-600 bg-yellow-50 border-yellow-200";
      case "Advanced": return "text-red-600 bg-red-50 border-red-200";
      default: return "text-gray-600 bg-gray-50 border-gray-200";
    }
  };

  if (loading || !challengeData) {
    return (
      <div className="min-h-screen bg-gradient-to-br from-background via-background to-accent/5">
        <div className="container-responsive padding-responsive-y">
          <div className="flex items-center justify-center min-h-[60vh]">
            <div className="text-center">
              <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-primary mx-auto mb-4"></div>
              <p className="text-responsive text-muted-foreground">Loading challenge...</p>
            </div>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gradient-to-br from-background via-background to-accent/5">
      <div className="container-responsive padding-responsive-y">
        {/* Breadcrumb Navigation */}
        <div className="flex items-center space-x-2 text-sm text-muted-foreground mb-6 lg:mb-8">
        <Link href="/" className="hover:text-primary">
          <Home className="h-4 w-4" />
        </Link>
        <ChevronRight className="h-4 w-4" />
        <Link href="/practice" className="hover:text-primary">
          Practice
        </Link>
        <ChevronRight className="h-4 w-4" />
        <span className="text-foreground">{challengeData.title}</span>
      </div>

      {/* Challenge Header */}
      <div className="mb-8">
        <div className="flex items-start justify-between mb-4">
          <div className="flex-1">
            <div className="flex items-center space-x-3 mb-3">
              <Badge className={getDifficultyColor(challengeData.difficulty)}>
                {challengeData.difficulty}
              </Badge>
              <Badge variant="outline">{challengeData.category}</Badge>
              {completed && (
                <Badge className="bg-green-100 text-green-800 border-green-200">
                  <CheckCircle className="mr-1 h-3 w-3" />
                  Completed
                </Badge>
              )}
            </div>
            <h1 className="text-3xl font-bold tracking-tight mb-2">{challengeData.title}</h1>
            <p className="text-muted-foreground text-lg mb-4">{challengeData.description}</p>
            
            <div className="flex items-center space-x-6 text-sm">
              <div className="flex items-center space-x-1">
                <Clock className="h-4 w-4" />
                <span>{challengeData.estimatedTime}</span>
              </div>
              <div className="flex items-center space-x-1">
                <Star className="h-4 w-4 text-yellow-500" />
                <span>{challengeData.points} points</span>
              </div>
              <div className="flex items-center space-x-1">
                <Target className="h-4 w-4" />
                <span>{challengeData.completionRate}% success rate</span>
              </div>
              <div className="flex items-center space-x-1">
                <Users className="h-4 w-4" />
                <span>{challengeData.totalAttempts} attempts</span>
              </div>
            </div>
          </div>
          <div className="ml-6">
            <Link href="/practice">
              <Button variant="outline" size="sm">
                <ArrowLeft className="mr-2 h-4 w-4" />
                Back to Practice
              </Button>
            </Link>
          </div>
        </div>

        {/* Progress indicator */}
        {attempts > 0 && (
          <div className="mb-4">
            <div className="flex items-center justify-between text-sm mb-2">
              <span>Attempts: {attempts}</span>
              <span>{testResults.filter(t => t.passed).length}/{testResults.length} tests passed</span>
            </div>
            <Progress 
              value={testResults.length > 0 ? (testResults.filter(t => t.passed).length / testResults.length) * 100 : 0} 
              className="h-2" 
            />
          </div>
        )}
      </div>

      {/* Main Content */}
      <Tabs value={activeTab} onValueChange={setActiveTab} className="space-y-6">
        <TabsList className="grid w-full grid-cols-3 lg:w-[400px]">
          <TabsTrigger value="challenge">Challenge</TabsTrigger>
          <TabsTrigger value="solution">Solution</TabsTrigger>
          <TabsTrigger value="discussion">Discussion</TabsTrigger>
        </TabsList>

        {/* Challenge Tab */}
        <TabsContent value="challenge" className="space-y-6">
          <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
            {/* Instructions */}
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center">
                  <BookOpen className="mr-2 h-5 w-5" />
                  Instructions
                </CardTitle>
              </CardHeader>
              <CardContent>
                <div className="prose dark:prose-invert max-w-none">
                  <div dangerouslySetInnerHTML={{ 
                    __html: challengeData.instructions.replace(/\n/g, '<br>').replace(/```/g, '<pre><code>').replace(/`([^`]+)`/g, '<code>$1</code>') 
                  }} />
                </div>
                
                {/* Tags */}
                <div className="flex flex-wrap gap-2 mt-4 pt-4 border-t">
                  {challengeData.tags.map((tag) => (
                    <Badge key={tag} variant="secondary" className="text-xs">
                      {tag}
                    </Badge>
                  ))}
                </div>

                {/* Hints */}
                <div className="mt-4 pt-4 border-t">
                  <Button
                    variant="outline"
                    size="sm"
                    onClick={() => setShowHints(!showHints)}
                    className="mb-3"
                  >
                    <Lightbulb className="mr-2 h-4 w-4" />
                    {showHints ? 'Hide Hints' : 'Show Hints'}
                  </Button>
                  
                  {showHints && (
                    <div className="space-y-2">
                      {challengeData.hints.map((hint, index) => (
                        <div key={index} className="p-3 bg-blue-50 border border-blue-200 rounded-lg">
                          <div className="flex items-start space-x-2">
                            <Lightbulb className="h-4 w-4 text-blue-600 mt-0.5 flex-shrink-0" />
                            <p className="text-sm text-blue-800">{hint}</p>
                          </div>
                        </div>
                      ))}
                    </div>
                  )}
                </div>
              </CardContent>
            </Card>

            {/* Code Editor */}
            <CodeEditor
              title="Your Solution"
              description="Write your Go code here. Click 'Run Code' to test your solution."
              initialCode={currentCode}
              solution={challengeData.solution}
              language="go"
              onCodeChange={handleCodeChange}
              onRun={handleRunCode}
              tests={testResults.map(test => ({
                name: test.description,
                passed: test.passed || false,
                message: test.passed ? `Expected: ${test.expectedOutput}` : `Expected: ${test.expectedOutput}, but test failed`
              }))}
            />
          </div>

          {/* Test Results */}
          {testResults.length > 0 && (
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center">
                  <Target className="mr-2 h-5 w-5" />
                  Test Results
                </CardTitle>
              </CardHeader>
              <CardContent>
                <div className="space-y-3">
                  {testResults.map((test) => (
                    <div key={test.id} className={`p-3 rounded-lg border ${
                      test.passed 
                        ? 'bg-green-50 border-green-200' 
                        : 'bg-red-50 border-red-200'
                    }`}>
                      <div className="flex items-center space-x-2 mb-2">
                        {test.passed ? (
                          <CheckCircle className="h-4 w-4 text-green-600" />
                        ) : (
                          <AlertCircle className="h-4 w-4 text-red-600" />
                        )}
                        <span className={`font-medium ${
                          test.passed ? 'text-green-800' : 'text-red-800'
                        }`}>
                          {test.passed ? 'Passed' : 'Failed'}
                        </span>
                      </div>
                      <p className="text-sm text-muted-foreground">{test.description}</p>
                    </div>
                  ))}
                </div>
              </CardContent>
            </Card>
          )}

          {/* Action Buttons */}
          <div className="flex items-center justify-between">
            <Button variant="outline" onClick={handleReset}>
              <RotateCcw className="mr-2 h-4 w-4" />
              Reset Code
            </Button>
            
            {completed && (
              <div className="flex items-center space-x-4">
                <div className="flex items-center space-x-2 text-green-600">
                  <Trophy className="h-5 w-5" />
                  <span className="font-medium">Challenge Completed! +{challengeData.points} points</span>
                </div>
                <Link href="/practice">
                  <Button className="go-gradient text-white">
                    Next Challenge
                    <ChevronRight className="ml-2 h-4 w-4" />
                  </Button>
                </Link>
              </div>
            )}
          </div>
        </TabsContent>
      </Tabs>
      </div>
    </div>
  );
}
