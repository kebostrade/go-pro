'use client';

import { useState, useEffect, useCallback } from 'react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
import PlaygroundEditor from '@/components/playground/PlaygroundEditor';
import { api } from '@/lib/api';
import {
  Code2,
  Play,
  Settings,
  Terminal,
  FileCode,
  Sparkles,
  BookOpen,
  Zap,
  Copy,
  CheckCircle,
  Clock,
  ArrowRight,
  AlertCircle,
  Wand2,
  Bug,
  TestTube,
  Brain,
  Loader2,
} from 'lucide-react';

// Go code templates
const codeTemplates = {
  helloWorld: `package main

import "fmt"

func main() {
    fmt.Println("Hello, Go!")
}`,
  variables: `package main

import "fmt"

func main() {
    // Variable declarations
    var name string = "Go Learner"
    age := 25 // Short declaration

    // Constants
    const greeting = "Welcome"

    fmt.Printf("%s, %s! You are %d years old.\\n", greeting, name, age)
}`,
  functions: `package main

import "fmt"

// Function with multiple return values
func divide(a, b float64) (float64, error) {
    if b == 0 {
        return 0, fmt.Errorf("division by zero")
    }
    return a / b, nil
}

// Variadic function
func sum(nums ...int) int {
    total := 0
    for _, n := range nums {
        total += n
    }
    return total
}

func main() {
    result, err := divide(10, 2)
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    fmt.Printf("10 / 2 = %.2f\\n", result)
    fmt.Printf("Sum: %d\\n", sum(1, 2, 3, 4, 5))
}`,
  structs: `package main

import "fmt"

// Struct definition
type Person struct {
    Name string
    Age  int
    City string
}

// Method on struct
func (p *Person) Greet() {
    fmt.Printf("Hi, I'm %s from %s!\\n", p.Name, p.City)
}

func main() {
    // Create struct instance
    person := Person{
        Name: "Alice",
        Age:  30,
        City: "Boston",
    }

    person.Greet()
    fmt.Printf("%s is %d years old\\n", person.Name, person.Age)
}`,
  slices: `package main

import "fmt"

func main() {
    // Create a slice
    numbers := []int{1, 2, 3, 4, 5}

    // Append elements
    numbers = append(numbers, 6, 7, 8)

    // Slice operations
    subset := numbers[2:5] // Elements from index 2 to 4
    fmt.Printf("Subset: %v\\n", subset)

    // Iterate over slice
    for i, num := range numbers {
        fmt.Printf("Index %d: %d\\n", i, num)
    }

    // Length and capacity
    fmt.Printf("Length: %d, Capacity: %d\\n", len(numbers), cap(numbers))
}`,
  maps: `package main

import "fmt"

func main() {
    // Create a map
    scores := make(map[string]int)

    // Add entries
    scores["Alice"] = 95
    scores["Bob"] = 87
    scores["Charlie"] = 92

    // Access value
    fmt.Printf("Alice's score: %d\\n", scores["Alice"])

    // Check if key exists
    if score, exists := scores["David"]; exists {
        fmt.Printf("David's score: %d\\n", score)
    } else {
        fmt.Println("David not found")
    }

    // Iterate over map
    for name, score := range scores {
        fmt.Printf("%s: %d\\n", name, score)
    }

    // Delete entry
    delete(scores, "Bob")
    fmt.Printf("Map length after delete: %d\\n", len(scores))
}`,
  goroutines: `package main

import (
    "fmt"
    "time"
)

func sayHello(name string) {
    fmt.Printf("Hello, %s!\\n", name)
}

func main() {
    // Launch goroutine
    go sayHello("World")
    go sayHello("Go")

    // Anonymous goroutine
    go func(msg string) {
        fmt.Println(msg)
    }("Goroutines are fun!")

    // Wait for goroutines to finish
    time.Sleep(100 * time.Millisecond)
    fmt.Println("Done!")
}`,
  channels: `package main

import "fmt"

func main() {
    // Create a channel
    messages := make(chan string)

    // Send data (in goroutine)
    go func() {
        messages <- "ping"
    }()

    // Receive data
    msg := <-messages
    fmt.Println(msg)

    // Buffered channel
    buffered := make(chan int, 2)
    buffered <- 1
    buffered <- 2
    // buffered <- 3  // This would block!

    fmt.Println(<-buffered)
    fmt.Println(<-buffered)
}`,
  interfaces: `package main

import (
    "fmt"
    "math"
)

// Interface definition
type Shape interface {
    Area() float64
    Perimeter() float64
}

// Circle implements Shape
type Circle struct {
    Radius float64
}

func (c Circle) Area() float64 {
    return math.Pi * c.Radius * c.Radius
}

func (c Circle) Perimeter() float64 {
    return 2 * math.Pi * c.Radius
}

// Rectangle implements Shape
type Rectangle struct {
    Width, Height float64
}

func (r Rectangle) Area() float64 {
    return r.Width * r.Height
}

func (r Rectangle) Perimeter() float64 {
    return 2 * (r.Width + r.Height)
}

func printShapeInfo(s Shape) {
    fmt.Printf("Area: %.2f\\n", s.Area())
    fmt.Printf("Perimeter: %.2f\\n", s.Perimeter())
}

func main() {
    c := Circle{Radius: 5}
    r := Rectangle{Width: 4, Height: 3}

    fmt.Println("Circle:")
    printShapeInfo(c)

    fmt.Println("\\nRectangle:")
    printShapeInfo(r)
}`,
  errorHandling: `package main

import (
    "errors"
    "fmt"
)

// Custom error type
type ValidationError struct {
    Field string
    Msg  string
}

func (e *ValidationError) Error() string {
    return fmt.Sprintf("validation error on %s: %s", e.Field, e.Msg)
}

func validateAge(age int) error {
    if age < 0 {
        return &ValidationError{
            Field: "age",
            Msg:  "cannot be negative",
        }
    }
    if age > 150 {
        return &ValidationError{
            Field: "age",
            Msg:  "seems unrealistic",
        }
    }
    return nil
}

func main() {
    // Using errors.Is to check error type
    err := validateAge(-5)
    if err != nil {
        var valErr *ValidationError
        if errors.As(err, &valErr) {
            fmt.Printf("Field %q is invalid: %s\\n", valErr.Field, valErr.Msg)
        }
    }

    // Wrapping errors
    original := errors.New("original error")
    wrapped := fmt.Errorf("operation failed: %w", original)
    fmt.Println(wrapped)
}`,
};

// Example categories
const categories = [
  { id: 'basics', name: 'Basics', icon: BookOpen, count: 3 },
  { id: 'data', name: 'Data Structures', icon: FileCode, count: 3 },
  { id: 'concurrency', name: 'Concurrency', icon: Zap, count: 2 },
  { id: 'advanced', name: 'Advanced', icon: Sparkles, count: 2 },
];

export default function PlaygroundPage() {
  const [code, setCode] = useState(codeTemplates.helloWorld);
  const [output, setOutput] = useState<string>('');
  const [error, setError] = useState<string>('');
  const [isRunning, setIsRunning] = useState(false);
  const [executionTime, setExecutionTime] = useState<number | null>(null);
  const [selectedCategory, setSelectedCategory] = useState('basics');
  const [copied, setCopied] = useState(false);

  // AI-powered features state
  const [aiAnalysis, setAiAnalysis] = useState<any>(null);
  const [isAnalyzing, setIsAnalyzing] = useState(false);
  const [aiError, setAiError] = useState<string>('');
  const [sessionId, setSessionId] = useState<string>('');
  const [showAiPanel, setShowAiPanel] = useState(false);

  const handleCodeChange = (value: string | undefined) => {
    setCode(value || '');
  };

  const handleRun = useCallback(async () => {
    setIsRunning(true);
    setOutput('');
    setError('');
    setExecutionTime(null);

    try {
      const result = await api.executePlaygroundCode(code);

      setExecutionTime(result.execution_time_ms);

      if (result.success) {
        setOutput(result.output || '(no output)');
        setError('');
      } else {
        setOutput(result.output || '');
        setError(result.error || 'Execution failed');
      }
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to execute code');
      setOutput('');
    } finally {
      setIsRunning(false);
    }
  }, [code]);

  const handleCopy = useCallback(() => {
    navigator.clipboard.writeText(code);
    setCopied(true);
    setTimeout(() => setCopied(false), 2000);
  }, [code]);

  const handleClear = useCallback(() => {
    setCode('');
    setOutput('');
    setError('');
    setExecutionTime(null);
    setAiAnalysis(null);
    setAiError('');
  }, []);

  // AI-powered analysis handler
  const handleAnalyze = useCallback(async () => {
    if (!code.trim()) return;

    setIsAnalyzing(true);
    setAiError('');

    try {
      // Create session if not exists
      let sid = sessionId;
      if (!sid) {
        const session = await api.createPlaygroundSession();
        sid = session.session_id;
        setSessionId(sid);
      }

      // Execute with AI analysis
      const result = await api.executeWithAI(code, 'go', sid);

      setOutput(result.output || '');
      setError(result.error || '');
      setExecutionTime(result.execution_time_ms);

      if (result.ai_analysis) {
        setAiAnalysis(result.ai_analysis);
        setShowAiPanel(true); // Show AI panel to display results
      }
    } catch (err) {
      setAiError(err instanceof Error ? err.message : 'AI analysis failed');
    } finally {
      setIsAnalyzing(false);
    }
  }, [code, sessionId]);

  // Generate test cases
  const handleGenerateTests = useCallback(async () => {
    if (!code.trim()) return;

    setIsAnalyzing(true);
    setAiError('');

    try {
      const result = await api.generateTests(code);
      setAiAnalysis((prev: any) => ({
        ...prev,
        testCases: result.test_cases,
      }));
      setShowAiPanel(true);
    } catch (err) {
      setAiError(err instanceof Error ? err.message : 'Test generation failed');
    } finally {
      setIsAnalyzing(false);
    }
  }, [code]);

  // Explain error
  const handleExplainError = useCallback(async () => {
    if (!error || !code.trim()) return;

    setIsAnalyzing(true);
    try {
      const result = await api.explainError(code, error);
      setAiAnalysis((prev: any) => ({
        ...prev,
        errorExplanations: result.explanations,
      }));
      setShowAiPanel(true);
    } catch (err) {
      setAiError(err instanceof Error ? err.message : 'Error explanation failed');
    } finally {
      setIsAnalyzing(false);
    }
  }, [code, error]);

  // Keyboard shortcuts
  useEffect(() => {
    const handleKeyDown = (e: KeyboardEvent) => {
      // Ctrl+Enter or Cmd+Enter to run
      if ((e.ctrlKey || e.metaKey) && e.key === 'Enter') {
        e.preventDefault();
        handleRun();
      }
      // Ctrl+S or Cmd+S to save (prevent default, show feedback)
      if ((e.ctrlKey || e.metaKey) && e.key === 's') {
        e.preventDefault();
        handleCopy();
      }
    };

    window.addEventListener('keydown', handleKeyDown);
    return () => window.removeEventListener('keydown', handleKeyDown);
  }, [handleRun, handleCopy]);

  const templatesInCategory = Object.entries(codeTemplates).filter(([key]) => {
    if (selectedCategory === 'basics') return ['helloWorld', 'variables', 'functions'].includes(key);
    if (selectedCategory === 'data') return ['structs', 'slices', 'maps'].includes(key);
    if (selectedCategory === 'concurrency') return ['goroutines', 'channels'].includes(key);
    if (selectedCategory === 'advanced') return ['interfaces', 'errorHandling'].includes(key);
    return false;
  });

  return (
    <div className="min-h-screen bg-gradient-to-br from-background via-background to-background">
      <div className="container mx-auto px-4 py-6 max-w-7xl">
        {/* Header */}
        <div className="mb-6">
          <div className="flex flex-col md:flex-row md:items-center md:justify-between gap-4">
            <div>
              <h1 className="text-3xl font-bold tracking-tight flex items-center gap-2">
                <Code2 className="h-8 w-8 text-primary" />
                Go Playground
              </h1>
              <p className="text-muted-foreground mt-1">
                Write, run, and experiment with Go code
              </p>
            </div>
            <div className="flex items-center gap-2">
              <Button variant="outline" size="sm" onClick={handleCopy}>
                {copied ? (
                  <>
                    <CheckCircle className="h-4 w-4 mr-1 text-green-500" />
                    Copied!
                  </>
                ) : (
                  <>
                    <Copy className="h-4 w-4 mr-1" />
                    Copy
                  </>
                )}
              </Button>
              <Button variant="outline" size="sm" onClick={handleClear}>
                Clear
              </Button>
            </div>
          </div>
        </div>

        <div className="grid grid-cols-1 lg:grid-cols-4 gap-6">
          {/* Sidebar - Code Templates */}
          <div className="lg:col-span-1 space-y-4">
            <Card className="border-none shadow-lg">
              <CardHeader className="pb-2">
                <CardTitle className="text-lg flex items-center gap-2">
                  <FileCode className="h-5 w-5 text-primary" />
                  Templates
                </CardTitle>
              </CardHeader>
              <CardContent className="space-y-3">
                {/* Category tabs */}
                <div className="flex flex-wrap gap-2">
                  {categories.map((cat) => (
                    <button
                      key={cat.id}
                      onClick={() => setSelectedCategory(cat.id)}
                      className={`flex items-center gap-1 px-3 py-1.5 rounded-md text-sm transition-all ${
                        selectedCategory === cat.id
                          ? 'bg-primary text-primary-foreground'
                          : 'bg-muted hover:bg-muted/80'
                      }`}
                    >
                      <cat.icon className="h-4 w-4" />
                      {cat.name}
                    </button>
                  ))}
                </div>

                {/* Template buttons */}
                <div className="space-y-2 mt-4">
                  {templatesInCategory.map(([key, template]) => (
                    <button
                      key={key}
                      onClick={() => setCode(template)}
                      className={`w-full text-left px-3 py-2 rounded-lg text-sm transition-all hover:bg-accent ${
                        code === template ? 'bg-primary/10 border border-primary/50' : ''
                      }`}
                    >
                      <div className="flex items-center justify-between">
                        <span className="capitalize">
                          {key.replace(/([A-Z])/g, ' $1').trim()}
                        </span>
                        <Badge variant="outline" className="text-xs">Go</Badge>
                      </div>
                    </button>
                  ))}
                </div>
              </CardContent>
            </Card>

            {/* Quick Tips */}
            <Card className="border-none shadow-lg bg-gradient-to-br from-blue-500/5 to-purple-500/5">
              <CardHeader className="pb-2">
                <CardTitle className="text-lg flex items-center gap-2">
                  <Sparkles className="h-5 w-5 text-yellow-500" />
                  Quick Tips
                </CardTitle>
              </CardHeader>
              <CardContent className="space-y-2 text-sm text-muted-foreground">
                <div className="flex items-start gap-2">
                  <CheckCircle className="h-4 w-4 text-green-500 mt-0.5 shrink-0" />
                  <span>Press <kbd className="px-1.5 py-0.5 bg-muted rounded text-xs">Ctrl+Enter</kbd> to run</span>
                </div>
                <div className="flex items-start gap-2">
                  <CheckCircle className="h-4 w-4 text-green-500 mt-0.5 shrink-0" />
                  <span>Press <kbd className="px-1.5 py-0.5 bg-muted rounded text-xs">Ctrl+S</kbd> to copy code</span>
                </div>
                <div className="flex items-start gap-2">
                  <CheckCircle className="h-4 w-4 text-green-500 mt-0.5 shrink-0" />
                  <span>Execution timeout: 5 seconds</span>
                </div>
              </CardContent>
            </Card>
          </div>

          {/* Main Editor Area */}
          <div className="lg:col-span-3 space-y-4">
            <Tabs defaultValue="editor" className="w-full">
              <TabsList className="grid w-full grid-cols-2">
                <TabsTrigger value="editor">
                  <Code2 className="h-4 w-4 mr-2" />
                  Editor
                </TabsTrigger>
                <TabsTrigger value="output">
                  <Terminal className="h-4 w-4 mr-2" />
                  Output
                </TabsTrigger>
              </TabsList>

              <TabsContent value="editor" className="mt-4">
                <Card className="border-none shadow-lg overflow-hidden">
                  <CardContent className="p-0">
                    <PlaygroundEditor
                      value={code}
                      onChange={handleCodeChange}
                      language="go"
                      height="500px"
                    />
                  </CardContent>
                </Card>

                {/* Action buttons */}
                <div className="flex items-center justify-between mt-4">
                  <div className="flex items-center gap-4 text-sm text-muted-foreground">
                    {executionTime !== null && (
                      <div className="flex items-center gap-1">
                        <Clock className="h-4 w-4" />
                        <span>{executionTime}ms</span>
                      </div>
                    )}
                    {!executionTime && (
                        <span>Run your code to see the output</span>
                      )}
                  </div>
                  <Button
                    onClick={handleRun}
                    disabled={isRunning}
                    className="bg-gradient-to-r from-primary to-blue-600"
                    size="lg"
                  >
                    {isRunning ? (
                      <>
                        <div className="animate-spin rounded-full h-4 w-4 border-2 border-white border-t-transparent mr-2" />
                        Running...
                      </>
                    ) : (
                      <>
                        <Play className="h-4 w-4 mr-2" />
                        Run Code
                      </>
                    )}
                  </Button>
                </div>
              </TabsContent>

              <TabsContent value="output" className="mt-4">
                <Card className="border-none shadow-lg">
                  <CardHeader className="pb-2">
                    <div className="flex items-center justify-between">
                      <CardTitle className="text-lg flex items-center gap-2">
                        <Terminal className="h-5 w-5 text-green-500" />
                        Console Output
                      </CardTitle>
                      <div className="flex items-center gap-2">
                        {executionTime !== null && (
                          <Badge variant="outline" className="text-xs">
                            {executionTime}ms
                          </Badge>
                        )}
                        <Button
                          variant="ghost"
                          size="sm"
                          onClick={() => {
                            setOutput('');
                            setError('');
                            setExecutionTime(null);
                          }}
                        >
                          Clear
                        </Button>
                      </div>
                    </div>
                  </CardHeader>
                  <CardContent>
                    <div className="bg-gray-950 rounded-lg p-4 font-mono text-sm min-h-[300px] max-h-[400px] overflow-auto">
                      {!output && !error ? (
                        <div className="text-gray-500 flex items-center justify-center h-full">
                          <div className="text-center">
                            <Terminal className="h-12 w-12 mx-auto mb-3 opacity-20" />
                            <p>No output yet</p>
                            <p className="text-xs mt-1">Run your code to see results</p>
                          </div>
                        </div>
                      ) : (
                        <div className="space-y-1">
                          {error && (
                            <div className="flex items-start gap-2 text-red-400 mb-2 pb-2 border-b border-red-400/20">
                              <AlertCircle className="h-4 w-4 mt-0.5 shrink-0" />
                              <pre className="whitespace-pre-wrap break-all">{error}</pre>
                            </div>
                          )}
                          {output && (
                            <pre className="text-gray-300 whitespace-pre-wrap">{output}</pre>
                          )}
                          {executionTime !== null && (
                            <div className="text-gray-500 text-xs mt-2 pt-2 border-t border-gray-700">
                                Program exited in {executionTime}ms
                              </div>
                          )}
                        </div>
                      )}
                    </div>
                  </CardContent>
                </Card>
              </TabsContent>
            </Tabs>
          </div>
        </div>

        {/* Learning Resources */}
        <div className="mt-8">
          <Card className="border-none shadow-lg bg-gradient-to-r from-primary/5 to-transparent">
            <CardContent className="p-6">
              <div className="flex flex-col md:flex-row items-center justify-between gap-4">
                <div className="flex items-center gap-4">
                  <BookOpen className="h-10 w-10 text-primary" />
                  <div>
                    <h3 className="text-xl font-bold">Ready to Learn More?</h3>
                    <p className="text-muted-foreground">
                      Explore our structured Go curriculum with hands-on exercises
                    </p>
                  </div>
                </div>
                <Button className="bg-gradient-to-r from-primary to-blue-600">
                  View Curriculum
                  <ArrowRight className="h-4 w-4 ml-2" />
                </Button>
              </div>
            </CardContent>
          </Card>
        </div>
      </div>
    </div>
  );
}
