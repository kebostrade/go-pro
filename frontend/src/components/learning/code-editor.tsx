"use client";

import { useState, useRef, useEffect } from "react";
import Editor from "@monaco-editor/react";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { Progress } from "@/components/ui/progress";
import {
  Play,
  RotateCcw,
  Copy,
  Check,
  AlertCircle,
  CheckCircle,
  Clock,
  Terminal,
  Eye,
  EyeOff,
  Maximize2,
  Minimize2,
  Settings,
  Palette,
  Zap,
  Download,
  Share2,
  History,
  Code2,
  Lightbulb,
  TrendingUp,
  Cpu,
  HardDrive,
  Activity
} from "lucide-react";
import { executeGoCode, analyzeCode } from "@/lib/code-execution";
import CodeDiffViewer from "./code-diff-viewer";

interface TestResult {
  name: string;
  passed: boolean;
  message?: string;
}

interface CodeEditorProps {
  title: string;
  description: string;
  initialCode: string;
  solution?: string;
  tests?: TestResult[];
  language?: string;
  theme?: string;
  readOnly?: boolean;
  onCodeChange?: (code: string) => void;
  onRun?: (code: string) => Promise<{ output: string; error?: string; tests?: TestResult[] }>;
  showLineNumbers?: boolean;
  enableAutoComplete?: boolean;
  enableFormatting?: boolean;
  maxHeight?: number;
}

const CodeEditor = ({
  title,
  description,
  initialCode,
  solution,
  tests = [],
  language = "go",
  theme = "vs-dark",
  readOnly = false,
  onCodeChange,
  onRun,
  showLineNumbers = true,
  enableAutoComplete = true,
  enableFormatting = true,
  maxHeight = 400
}: CodeEditorProps) => {
  const [code, setCode] = useState(initialCode);
  const [output, setOutput] = useState("");
  const [error, setError] = useState("");
  const [testResults, setTestResults] = useState<TestResult[]>(tests);
  const [isRunning, setIsRunning] = useState(false);
  const [copied, setCopied] = useState(false);
  const [showSolution, setShowSolution] = useState(false);
  const [currentTheme, setCurrentTheme] = useState(theme);
  const [isFullscreen, setIsFullscreen] = useState(false);
  const [codeHistory, setCodeHistory] = useState<string[]>([initialCode]);
  const [historyIndex, setHistoryIndex] = useState(0);
  const [executionTime, setExecutionTime] = useState(0);
  const [showSettings, setShowSettings] = useState(false);
  const editorRef = useRef<any>(null);

  const handleEditorDidMount = (editor: any) => {
    editorRef.current = editor;
  };

  const handleCodeChange = (value: string | undefined) => {
    const newCode = value || "";
    setCode(newCode);
    onCodeChange?.(newCode);

    // Add to history if significantly different
    if (newCode !== codeHistory[historyIndex] && newCode.length > 0) {
      const newHistory = codeHistory.slice(0, historyIndex + 1);
      newHistory.push(newCode);
      setCodeHistory(newHistory);
      setHistoryIndex(newHistory.length - 1);
    }
  };

  const handleRun = async () => {
    const startTime = Date.now();

    if (!onRun) {
      // Use enhanced execution engine
      setIsRunning(true);
      setOutput("Running your Go code...\n");
      setError("");

      try {
        const result = await executeGoCode(code);
        const endTime = Date.now();
        setExecutionTime(endTime - startTime);

        if (result.error) {
          setError(result.error);
          setOutput("");
        } else {
          setOutput(result.output);
          setError("");
        }

        // Generate test results based on execution
        const analysis = analyzeCode(code);
        setTestResults([
          { name: "Syntax Check", passed: !result.error },
          { name: "Code Structure", passed: analysis.hasMain || !code.includes('package main') },
          { name: "Complexity Analysis", passed: analysis.complexity <= 15, message: analysis.complexity > 15 ? `High complexity: ${analysis.complexity}` : undefined },
          { name: "Best Practices", passed: code.includes('error') || !code.includes('func'), message: !code.includes('error') && code.includes('func') ? "Consider adding error handling" : undefined },
        ]);
      } catch (err) {
        setError(err instanceof Error ? err.message : "An error occurred");
      } finally {
        setIsRunning(false);
      }
      return;
    }

    setIsRunning(true);
    setError("");
    setOutput("");

    try {
      const result = await onRun(code);
      const endTime = Date.now();
      setExecutionTime(endTime - startTime);
      setOutput(result.output);
      if (result.error) {
        setError(result.error);
      }
      if (result.tests) {
        setTestResults(result.tests);
      }
    } catch (err) {
      setError(err instanceof Error ? err.message : "An error occurred");
    } finally {
      setIsRunning(false);
    }
  };

  const handleReset = () => {
    setCode(initialCode);
    setOutput("");
    setError("");
    setTestResults(tests);
    setShowSolution(false);
  };

  const handleCopy = async () => {
    try {
      await navigator.clipboard.writeText(code);
      setCopied(true);
      setTimeout(() => setCopied(false), 2000);
    } catch (err) {
      console.error("Failed to copy code:", err);
    }
  };

  const toggleSolution = () => {
    if (solution) {
      if (showSolution) {
        setCode(initialCode);
      } else {
        setCode(solution);
      }
      setShowSolution(!showSolution);
    }
  };

  const formatCode = () => {
    if (editorRef.current && enableFormatting) {
      editorRef.current.getAction('editor.action.formatDocument').run();
    }
  };

  const toggleTheme = () => {
    const newTheme = currentTheme === 'vs-dark' ? 'vs-light' : 'vs-dark';
    setCurrentTheme(newTheme);
  };

  const toggleFullscreen = () => {
    setIsFullscreen(!isFullscreen);
  };

  const undoCode = () => {
    if (historyIndex > 0) {
      const newIndex = historyIndex - 1;
      setHistoryIndex(newIndex);
      setCode(codeHistory[newIndex]);
    }
  };

  const redoCode = () => {
    if (historyIndex < codeHistory.length - 1) {
      const newIndex = historyIndex + 1;
      setHistoryIndex(newIndex);
      setCode(codeHistory[newIndex]);
    }
  };

  const downloadCode = () => {
    const blob = new Blob([code], { type: 'text/plain' });
    const url = URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = `${title.toLowerCase().replace(/\s+/g, '-')}.go`;
    document.body.appendChild(a);
    a.click();
    document.body.removeChild(a);
    URL.revokeObjectURL(url);
  };

  const shareCode = async () => {
    if (navigator.share) {
      try {
        await navigator.share({
          title: title,
          text: description,
          url: window.location.href
        });
      } catch (err) {
        console.log('Error sharing:', err);
      }
    } else {
      // Fallback to clipboard
      await handleCopy();
    }
  };

  const passedTests = testResults.filter(test => test.passed).length;
  const totalTests = testResults.length;

  return (
    <div className={`space-y-6 ${isFullscreen ? 'fixed inset-0 z-50 bg-background p-6' : ''}`}>
      {/* Enhanced Header */}
      <Card className="glass-card border-2">
        <CardHeader>
          <div className="flex items-start justify-between mb-4">
            <div className="flex-1">
              <CardTitle className="text-xl flex items-center">
                <Code2 className="mr-2 h-5 w-5 text-primary" />
                {title}
              </CardTitle>
              <CardDescription className="mt-2 text-base">
                {description}
              </CardDescription>
            </div>

            {/* Status Badges */}
            <div className="flex items-center space-x-2">
              <Badge variant="outline" className="capitalize hover:bg-primary/10 transition-colors">
                {language}
              </Badge>
              {totalTests > 0 && (
                <Badge
                  variant={passedTests === totalTests ? "default" : "secondary"}
                  className={`transition-all ${passedTests === totalTests ? "bg-green-500 text-white shadow-lg" : "hover:bg-muted"}`}
                >
                  {passedTests}/{totalTests} Tests
                </Badge>
              )}
              {executionTime > 0 && (
                <Badge variant="outline" className="text-xs">
                  <Zap className="mr-1 h-3 w-3" />
                  {executionTime}ms
                </Badge>
              )}
            </div>
          </div>

          {/* Enhanced Toolbar */}
          <div className="flex items-center justify-between">
            <div className="flex items-center space-x-2">
              {/* Primary Actions */}
              <Button
                onClick={handleRun}
                disabled={isRunning}
                className="bg-gradient-to-r from-green-500 to-green-600 hover:from-green-600 hover:to-green-700 text-white shadow-lg hover:shadow-xl transition-all"
              >
                {isRunning ? (
                  <div className="animate-spin rounded-full h-4 w-4 border-b-2 border-white mr-2" />
                ) : (
                  <Play className="mr-2 h-4 w-4" />
                )}
                {isRunning ? "Running..." : "Run Code"}
              </Button>

              <Button variant="outline" onClick={handleReset} className="hover:bg-muted transition-colors">
                <RotateCcw className="mr-2 h-4 w-4" />
                Reset
              </Button>

              {solution && (
                <Button
                  variant={showSolution ? "default" : "outline"}
                  onClick={toggleSolution}
                  className="hover:shadow-md transition-all"
                >
                  {showSolution ? <EyeOff className="mr-2 h-4 w-4" /> : <Eye className="mr-2 h-4 w-4" />}
                  {showSolution ? "Hide Solution" : "Show Solution"}
                </Button>
              )}
            </div>

            {/* Secondary Actions */}
            <div className="flex items-center space-x-1">
              <Button variant="ghost" size="sm" onClick={undoCode} disabled={historyIndex <= 0} title="Undo">
                <History className="h-4 w-4" />
              </Button>

              <Button variant="ghost" size="sm" onClick={redoCode} disabled={historyIndex >= codeHistory.length - 1} title="Redo">
                <History className="h-4 w-4 scale-x-[-1]" />
              </Button>

              <Button variant="ghost" size="sm" onClick={formatCode} title="Format Code">
                <Settings className="h-4 w-4" />
              </Button>

              <Button variant="ghost" size="sm" onClick={toggleTheme} title="Toggle Theme">
                <Palette className="h-4 w-4" />
              </Button>

              <Button variant="ghost" size="sm" onClick={handleCopy} title="Copy Code">
                {copied ? <Check className="h-4 w-4 text-green-500" /> : <Copy className="h-4 w-4" />}
              </Button>

              <Button variant="ghost" size="sm" onClick={downloadCode} title="Download">
                <Download className="h-4 w-4" />
              </Button>

              <Button variant="ghost" size="sm" onClick={shareCode} title="Share">
                <Share2 className="h-4 w-4" />
              </Button>

              <Button variant="ghost" size="sm" onClick={toggleFullscreen} title={isFullscreen ? "Exit Fullscreen" : "Fullscreen"}>
                {isFullscreen ? <Minimize2 className="h-4 w-4" /> : <Maximize2 className="h-4 w-4" />}
              </Button>
            </div>
          </div>
        </CardHeader>
      </Card>

      {/* Editor and Output */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        {/* Code Editor */}
        <Card>
          <CardHeader className="pb-3">
            <div className="flex items-center justify-between">
              <CardTitle className="text-lg">Code Editor</CardTitle>
              <div className="flex items-center space-x-2">
                <Button
                  variant="outline"
                  size="sm"
                  onClick={handleCopy}
                  className="h-8"
                >
                  {copied ? (
                    <Check className="h-3 w-3" />
                  ) : (
                    <Copy className="h-3 w-3" />
                  )}
                </Button>
                <Button
                  variant="outline"
                  size="sm"
                  onClick={handleReset}
                  className="h-8"
                >
                  <RotateCcw className="h-3 w-3" />
                </Button>
                {solution && (
                  <Button
                    variant="outline"
                    size="sm"
                    onClick={toggleSolution}
                    className="h-8"
                  >
                    {showSolution ? "Hide" : "Show"} Solution
                  </Button>
                )}
              </div>
            </div>
          </CardHeader>
          <CardContent className="p-0">
            <div className="border rounded-lg overflow-hidden">
              <Editor
                height={isFullscreen ? "calc(100vh - 300px)" : `${maxHeight}px`}
                language={language}
                theme={currentTheme}
                value={code}
                onChange={handleCodeChange}
                onMount={handleEditorDidMount}
                options={{
                  readOnly,
                  minimap: { enabled: !isFullscreen },
                  fontSize: isFullscreen ? 16 : 14,
                  lineHeight: isFullscreen ? 24 : 20,
                  padding: { top: 16, bottom: 16 },
                  scrollBeyondLastLine: false,
                  automaticLayout: true,
                  tabSize: 2,
                  insertSpaces: true,
                  wordWrap: "on",
                  lineNumbers: showLineNumbers ? "on" : "off",
                  glyphMargin: false,
                  folding: true,
                  lineDecorationsWidth: 0,
                  lineNumbersMinChars: 3,
                  renderLineHighlight: "line",
                  cursorBlinking: "smooth",
                  cursorSmoothCaretAnimation: "on",
                  smoothScrolling: true,
                  mouseWheelZoom: true,
                  formatOnPaste: enableFormatting,
                  formatOnType: enableFormatting,
                  autoIndent: "full",
                  bracketPairColorization: { enabled: true },
                  guides: {
                    bracketPairs: true,
                    indentation: true,
                  },
                  suggest: {
                    enabled: enableAutoComplete,
                    showKeywords: true,
                    showSnippets: true,
                  },
                  quickSuggestions: enableAutoComplete,
                  parameterHints: { enabled: enableAutoComplete },
                  scrollbar: {
                    vertical: "visible",
                    horizontal: "visible",
                    useShadows: false,
                    verticalHasArrows: false,
                    horizontalHasArrows: false,
                    verticalScrollbarSize: 10,
                    horizontalScrollbarSize: 10,
                  },
                  overviewRulerBorder: false,
                  hideCursorInOverviewRuler: true,
                  overviewRulerLanes: 0,
                  roundedSelection: false,
                }}
              />
            </div>
            <div className="p-4 border-t">
              <Button
                onClick={handleRun}
                disabled={isRunning || readOnly}
                className="go-gradient text-white"
              >
                {isRunning ? (
                  <>
                    <Clock className="mr-2 h-4 w-4 animate-spin" />
                    Running...
                  </>
                ) : (
                  <>
                    <Play className="mr-2 h-4 w-4" />
                    Run Code
                  </>
                )}
              </Button>
            </div>
          </CardContent>
        </Card>

        {/* Output and Tests */}
        <Card>
          <CardHeader className="pb-3">
            <CardTitle className="text-lg flex items-center">
              <Terminal className="mr-2 h-5 w-5" />
              Output & Tests
            </CardTitle>
          </CardHeader>
          <CardContent>
            <Tabs defaultValue="output" className="w-full">
              <TabsList className="grid w-full grid-cols-2">
                <TabsTrigger value="output">Output</TabsTrigger>
                <TabsTrigger value="tests">
                  Tests {totalTests > 0 && `(${passedTests}/${totalTests})`}
                </TabsTrigger>
              </TabsList>
              
              <TabsContent value="output" className="mt-4">
                <div className="code-block min-h-[300px] p-4 font-mono text-sm">
                  {error ? (
                    <div className="text-red-500">
                      <div className="flex items-center mb-2">
                        <AlertCircle className="mr-2 h-4 w-4" />
                        Error
                      </div>
                      <pre className="whitespace-pre-wrap">{error}</pre>
                    </div>
                  ) : output ? (
                    <pre className="whitespace-pre-wrap">{output}</pre>
                  ) : (
                    <div className="text-muted-foreground italic">
                      Click "Run Code" to see the output here...
                    </div>
                  )}
                </div>
              </TabsContent>
              
              <TabsContent value="tests" className="mt-4">
                <div className="space-y-3">
                  {testResults.length > 0 ? (
                    testResults.map((test, index) => (
                      <div
                        key={index}
                        className={`flex items-start space-x-3 p-3 rounded-lg border ${
                          test.passed
                            ? 'bg-green-50 border-green-200 dark:bg-green-950 dark:border-green-800'
                            : 'bg-red-50 border-red-200 dark:bg-red-950 dark:border-red-800'
                        }`}
                      >
                        {test.passed ? (
                          <CheckCircle className="h-4 w-4 text-green-500 mt-0.5" />
                        ) : (
                          <AlertCircle className="h-4 w-4 text-red-500 mt-0.5" />
                        )}
                        <div className="flex-1">
                          <div className="font-medium text-sm">{test.name}</div>
                          {test.message && (
                            <div className="text-sm text-muted-foreground mt-1">
                              {test.message}
                            </div>
                          )}
                        </div>
                      </div>
                    ))
                  ) : (
                    <div className="text-muted-foreground italic text-center py-8">
                      No tests available for this exercise
                    </div>
                  )}
                </div>
              </TabsContent>
            </Tabs>
          </CardContent>
        </Card>
      </div>
    </div>
  );
};

export default CodeEditor;
