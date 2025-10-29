"use client";

import { useState, useEffect, useRef } from "react";
import Editor from "@monaco-editor/react";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { Input } from "@/components/ui/input";
import {
  Play,
  RotateCcw,
  Settings,
  Zap,
  Code2,
  Terminal,
  Eye,
  EyeOff,
  Lightbulb,
  Target,
  CheckCircle,
  AlertCircle,
  GitCompare,
  BookOpen,
  TrendingUp,
  Activity,
  Cpu,
  HardDrive
} from "lucide-react";
import { executeGoCode, analyzeCode } from "@/lib/code-execution";
import CodeDiffViewer from "./code-diff-viewer";
import CodeSnippetsLibrary from "./code-snippets-library";

interface InteractiveExampleProps {
  title: string;
  description: string;
  initialCode: string;
  language?: string;
  inputs?: Array<{
    name: string;
    type: 'string' | 'number' | 'boolean';
    defaultValue: any;
    description?: string;
  }>;
  expectedOutput?: string;
  explanation?: string;
  concepts?: string[];
}

interface ExecutionResult {
  output: string;
  error?: string;
  executionTime: number;
  memoryUsed?: number;
  complexity?: {
    time: string;
    space: string;
  };
  warnings?: string[];
  suggestions?: string[];
}

export default function InteractiveExample({
  title,
  description,
  initialCode,
  language = "go",
  inputs = [],
  expectedOutput,
  explanation,
  concepts = []
}: InteractiveExampleProps) {
  const [code, setCode] = useState(initialCode);
  const [inputValues, setInputValues] = useState<Record<string, any>>({});
  const [result, setResult] = useState<ExecutionResult | null>(null);
  const [isRunning, setIsRunning] = useState(false);
  const [showExplanation, setShowExplanation] = useState(false);
  const [attempts, setAttempts] = useState(0);
  const [showDiff, setShowDiff] = useState(false);
  const [showSnippets, setShowSnippets] = useState(false);
  const [editorTheme, setEditorTheme] = useState<'vs-dark' | 'light'>('vs-dark');
  const editorRef = useRef<any>(null);

  // Initialize input values
  useEffect(() => {
    const initialInputs: Record<string, any> = {};
    inputs.forEach(input => {
      initialInputs[input.name] = input.defaultValue;
    });
    setInputValues(initialInputs);
  }, [inputs]);

  const handleEditorDidMount = (editor: any) => {
    editorRef.current = editor;
  };

  const handleCodeChange = (value: string | undefined) => {
    setCode(value || "");
  };

  const handleRun = async () => {
    setIsRunning(true);
    setAttempts(prev => prev + 1);

    try {
      const executionResult = await executeGoCode(code, inputValues);
      setResult(executionResult);
    } catch (error) {
      setResult({
        output: "",
        error: "Failed to execute code",
        executionTime: 0
      });
    } finally {
      setIsRunning(false);
    }
  };

  const handleReset = () => {
    setCode(initialCode);
    setResult(null);
    setAttempts(0);
    // Reset input values
    const resetInputs: Record<string, any> = {};
    inputs.forEach(input => {
      resetInputs[input.name] = input.defaultValue;
    });
    setInputValues(resetInputs);
  };

  const updateInputValue = (name: string, value: any) => {
    setInputValues(prev => ({ ...prev, [name]: value }));
  };

  const isOutputCorrect = result && expectedOutput && result.output === expectedOutput.trim();

  return (
    <Card className="glass-card border-2 hover-lift">
      <CardHeader>
        <CardTitle className="flex items-center">
          <Zap className="mr-2 h-5 w-5 text-primary" />
          {title}
        </CardTitle>
        <CardDescription>{description}</CardDescription>
        
        {concepts.length > 0 && (
          <div className="flex flex-wrap gap-2 mt-3">
            {concepts.map(concept => (
              <Badge key={concept} variant="outline" className="text-xs">
                <Target className="mr-1 h-3 w-3" />
                {concept}
              </Badge>
            ))}
          </div>
        )}
      </CardHeader>

      <CardContent>
        <Tabs defaultValue="code" className="space-y-4">
          <TabsList className="grid w-full grid-cols-5">
            <TabsTrigger value="code">
              <Code2 className="mr-2 h-4 w-4" />
              Code
            </TabsTrigger>
            <TabsTrigger value="inputs">
              <Settings className="mr-2 h-4 w-4" />
              Inputs
            </TabsTrigger>
            <TabsTrigger value="output">
              <Terminal className="mr-2 h-4 w-4" />
              Output
            </TabsTrigger>
            <TabsTrigger value="diff">
              <GitCompare className="mr-2 h-4 w-4" />
              Compare
            </TabsTrigger>
            <TabsTrigger value="snippets">
              <BookOpen className="mr-2 h-4 w-4" />
              Snippets
            </TabsTrigger>
          </TabsList>

          <TabsContent value="code" className="space-y-4">
            <div className="space-y-3">
              {/* Monaco Editor */}
              <div className="border rounded-lg overflow-hidden">
                <Editor
                  height="300px"
                  language={language}
                  theme={editorTheme}
                  value={code}
                  onChange={handleCodeChange}
                  onMount={handleEditorDidMount}
                  options={{
                    minimap: { enabled: false },
                    fontSize: 14,
                    lineHeight: 20,
                    padding: { top: 16, bottom: 16 },
                    scrollBeyondLastLine: false,
                    automaticLayout: true,
                    tabSize: 2,
                    insertSpaces: true,
                    wordWrap: "on",
                    lineNumbers: "on",
                    glyphMargin: false,
                    folding: true,
                    renderLineHighlight: "line",
                    cursorBlinking: "smooth",
                    smoothScrolling: true,
                    formatOnPaste: true,
                    formatOnType: true,
                    autoIndent: "full",
                    bracketPairColorization: { enabled: true },
                    guides: {
                      bracketPairs: true,
                      indentation: true,
                    },
                    suggest: {
                      enabled: true,
                      showKeywords: true,
                      showSnippets: true,
                    },
                    quickSuggestions: true,
                    scrollbar: {
                      vertical: "visible",
                      horizontal: "visible",
                      useShadows: false,
                      verticalScrollbarSize: 10,
                      horizontalScrollbarSize: 10,
                    },
                  }}
                />
              </div>

              <div className="flex items-center justify-between">
                <div className="flex items-center space-x-2">
                  <Button onClick={handleRun} disabled={isRunning} className="hover-glow">
                    {isRunning ? (
                      <>
                        <div className="animate-spin rounded-full h-4 w-4 border-b-2 border-white mr-2" />
                        Running...
                      </>
                    ) : (
                      <>
                        <Play className="mr-2 h-4 w-4" />
                        Run Code
                      </>
                    )}
                  </Button>

                  <Button variant="outline" onClick={handleReset}>
                    <RotateCcw className="mr-2 h-4 w-4" />
                    Reset
                  </Button>

                  <Button
                    variant="outline"
                    onClick={() => setEditorTheme(editorTheme === 'vs-dark' ? 'light' : 'vs-dark')}
                  >
                    {editorTheme === 'vs-dark' ? '☀️' : '🌙'}
                  </Button>
                </div>

                {result && (
                  <div className="flex items-center space-x-3 text-sm">
                    <Badge variant="outline" className="flex items-center space-x-1">
                      <Zap className="h-3 w-3" />
                      <span>{result.executionTime}ms</span>
                    </Badge>
                    {result.memoryUsed && (
                      <Badge variant="outline" className="flex items-center space-x-1">
                        <HardDrive className="h-3 w-3" />
                        <span>{result.memoryUsed}KB</span>
                      </Badge>
                    )}
                    {result.complexity && (
                      <Badge variant="outline" className="flex items-center space-x-1">
                        <Activity className="h-3 w-3" />
                        <span>{result.complexity.time}</span>
                      </Badge>
                    )}
                  </div>
                )}
              </div>
            </div>
          </TabsContent>

          <TabsContent value="inputs" className="space-y-4">
            {inputs.length > 0 ? (
              <div className="space-y-4">
                {inputs.map(input => (
                  <div key={input.name} className="space-y-2">
                    <label className="text-sm font-medium">
                      {input.name} ({input.type})
                    </label>
                    {input.description && (
                      <p className="text-xs text-muted-foreground">{input.description}</p>
                    )}
                    
                    {input.type === 'string' && (
                      <Input
                        value={inputValues[input.name] || ''}
                        onChange={(e) => updateInputValue(input.name, e.target.value)}
                        placeholder={`Enter ${input.name}...`}
                      />
                    )}
                    
                    {input.type === 'number' && (
                      <Input
                        type="number"
                        value={inputValues[input.name] || 0}
                        onChange={(e) => updateInputValue(input.name, parseFloat(e.target.value) || 0)}
                        placeholder={`Enter ${input.name}...`}
                      />
                    )}
                    
                    {input.type === 'boolean' && (
                      <div className="flex items-center space-x-2">
                        <input
                          type="checkbox"
                          checked={inputValues[input.name] || false}
                          onChange={(e) => updateInputValue(input.name, e.target.checked)}
                          className="rounded"
                        />
                        <span className="text-sm">{inputValues[input.name] ? 'true' : 'false'}</span>
                      </div>
                    )}
                  </div>
                ))}
              </div>
            ) : (
              <div className="text-center py-8 text-muted-foreground">
                <Settings className="mx-auto h-12 w-12 mb-4 opacity-50" />
                <p>No input parameters for this example</p>
              </div>
            )}
          </TabsContent>

          <TabsContent value="output" className="space-y-4">
            {result ? (
              <div className="space-y-4">
                {result.error ? (
                  <Card className="border-red-200 bg-red-50 dark:bg-red-950 dark:border-red-800">
                    <CardContent className="p-4">
                      <div className="flex items-start space-x-3">
                        <AlertCircle className="h-5 w-5 text-red-500 mt-0.5" />
                        <div className="flex-1">
                          <h4 className="font-medium text-red-700 dark:text-red-300">Error</h4>
                          <pre className="text-sm text-red-600 dark:text-red-400 mt-1 whitespace-pre-wrap">{result.error}</pre>
                        </div>
                      </div>
                    </CardContent>
                  </Card>
                ) : (
                  <>
                    <Card className={`border-2 ${isOutputCorrect ? 'border-green-200 bg-green-50 dark:bg-green-950 dark:border-green-800' : 'border-blue-200 bg-blue-50 dark:bg-blue-950 dark:border-blue-800'}`}>
                      <CardContent className="p-4">
                        <div className="flex items-start justify-between mb-2">
                          <h4 className="font-medium flex items-center">
                            <Terminal className="mr-2 h-4 w-4" />
                            Output
                          </h4>
                          {isOutputCorrect && (
                            <Badge variant="outline" className="border-green-500 text-green-600">
                              <CheckCircle className="mr-1 h-3 w-3" />
                              Correct!
                            </Badge>
                          )}
                        </div>
                        <pre className="text-sm bg-background p-3 rounded border overflow-x-auto whitespace-pre-wrap">
                          {result.output || '(no output)'}
                        </pre>
                      </CardContent>
                    </Card>

                    {/* Performance Metrics */}
                    {result.complexity && (
                      <Card className="border-purple-200 bg-purple-50 dark:bg-purple-950 dark:border-purple-800">
                        <CardContent className="p-4">
                          <h4 className="font-medium flex items-center mb-3">
                            <TrendingUp className="mr-2 h-4 w-4 text-purple-600" />
                            Performance Analysis
                          </h4>
                          <div className="grid grid-cols-2 gap-3">
                            <div className="flex items-center space-x-2">
                              <Cpu className="h-4 w-4 text-purple-600" />
                              <div>
                                <div className="text-xs text-muted-foreground">Time Complexity</div>
                                <div className="text-sm font-medium">{result.complexity.time}</div>
                              </div>
                            </div>
                            <div className="flex items-center space-x-2">
                              <HardDrive className="h-4 w-4 text-purple-600" />
                              <div>
                                <div className="text-xs text-muted-foreground">Space Complexity</div>
                                <div className="text-sm font-medium">{result.complexity.space}</div>
                              </div>
                            </div>
                          </div>
                        </CardContent>
                      </Card>
                    )}

                    {/* Warnings */}
                    {result.warnings && result.warnings.length > 0 && (
                      <Card className="border-yellow-200 bg-yellow-50 dark:bg-yellow-950 dark:border-yellow-800">
                        <CardContent className="p-4">
                          <h4 className="font-medium flex items-center mb-2 text-yellow-700 dark:text-yellow-300">
                            <AlertCircle className="mr-2 h-4 w-4" />
                            Warnings
                          </h4>
                          <ul className="space-y-1">
                            {result.warnings.map((warning, index) => (
                              <li key={index} className="text-sm text-yellow-600 dark:text-yellow-400">
                                • {warning}
                              </li>
                            ))}
                          </ul>
                        </CardContent>
                      </Card>
                    )}

                    {/* Suggestions */}
                    {result.suggestions && result.suggestions.length > 0 && (
                      <Card className="border-blue-200 bg-blue-50 dark:bg-blue-950 dark:border-blue-800">
                        <CardContent className="p-4">
                          <h4 className="font-medium flex items-center mb-2 text-blue-700 dark:text-blue-300">
                            <Lightbulb className="mr-2 h-4 w-4" />
                            Suggestions
                          </h4>
                          <ul className="space-y-1">
                            {result.suggestions.map((suggestion, index) => (
                              <li key={index} className="text-sm text-blue-600 dark:text-blue-400">
                                • {suggestion}
                              </li>
                            ))}
                          </ul>
                        </CardContent>
                      </Card>
                    )}
                  </>
                )}

                {explanation && (
                  <Card>
                    <CardContent className="p-4">
                      <Button
                        variant="ghost"
                        onClick={() => setShowExplanation(!showExplanation)}
                        className="w-full justify-between p-0 h-auto font-medium mb-2"
                      >
                        <span className="flex items-center">
                          <Lightbulb className="mr-2 h-4 w-4" />
                          Explanation
                        </span>
                        {showExplanation ? <EyeOff className="h-4 w-4" /> : <Eye className="h-4 w-4" />}
                      </Button>

                      {showExplanation && (
                        <div className="text-sm text-muted-foreground leading-relaxed animate-in slide-in-from-top-2 duration-200">
                          {explanation}
                        </div>
                      )}
                    </CardContent>
                  </Card>
                )}
              </div>
            ) : (
              <div className="text-center py-8 text-muted-foreground">
                <Terminal className="mx-auto h-12 w-12 mb-4 opacity-50" />
                <p>Run the code to see the output</p>
              </div>
            )}
          </TabsContent>

          {/* Diff Tab */}
          <TabsContent value="diff" className="space-y-4">
            <CodeDiffViewer
              originalCode={initialCode}
              modifiedCode={code}
              originalLabel="Initial Code"
              modifiedLabel="Your Code"
              language={language}
            />
          </TabsContent>

          {/* Snippets Tab */}
          <TabsContent value="snippets" className="space-y-4">
            <CodeSnippetsLibrary
              onInsert={(snippetCode) => {
                setCode(code + '\n\n' + snippetCode);
              }}
            />
          </TabsContent>
        </Tabs>
      </CardContent>
    </Card>
  );
}
