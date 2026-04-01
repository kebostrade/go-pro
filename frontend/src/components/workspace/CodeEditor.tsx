'use client';

import React, { useState } from 'react';
import MonacoEditor from './MonacoEditor';
import OutputConsole from './OutputConsole';
import { Play, RotateCcw } from 'lucide-react';

interface CodeEditorProps {
  topicId: string;
  initialCode?: string;
  testCases?: TestCase[];
}

interface TestCase {
  name: string;
  input: string;
  expected: string;
}

interface ExecuteResult {
  passed: boolean;
  score: number;
  results: TestResult[];
  execution_time: string;
}

interface TestResult {
  test_name: string;
  passed: boolean;
  expected: string;
  actual: string;
  error?: string;
}

const DEFAULT_CODE = `package main

import "fmt"

func main() {
	fmt.Println("Hello, World!")
}
`;

export default function CodeEditor({
  topicId,
  initialCode,
  testCases = [],
}: CodeEditorProps) {
  const [code, setCode] = useState(initialCode || DEFAULT_CODE);
  const [output, setOutput] = useState<ExecuteResult | null>(null);
  const [isRunning, setIsRunning] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const handleExecute = async () => {
    setIsRunning(true);
    setError(null);
    setOutput(null);

    try {
      const response = await fetch('/api/execute', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          code,
          topic_id: topicId,
          test_cases: testCases,
        }),
      });

      if (!response.ok) {
        const errorData = await response.text();
        throw new Error(errorData || `HTTP ${response.status}`);
      }

      const result: ExecuteResult = await response.json();
      setOutput(result);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Execution failed');
    } finally {
      setIsRunning(false);
    }
  };

  const handleReset = () => {
    setCode(DEFAULT_CODE);
    setOutput(null);
    setError(null);
  };

  return (
    <div className="space-y-4">
      {/* Editor Header */}
      <div className="flex items-center justify-between px-4 py-2 bg-gray-100 dark:bg-gray-800 rounded-t-lg border border-b-0 border-gray-200 dark:border-gray-700">
        <span className="font-mono text-sm text-gray-600 dark:text-gray-400">main.go</span>
        <div className="flex items-center gap-2">
          <button
            onClick={handleReset}
            className="p-2 text-gray-500 hover:text-gray-700 dark:text-gray-400 dark:hover:text-gray-200"
            title="Reset to default"
          >
            <RotateCcw className="w-4 h-4" />
          </button>
          <button
            onClick={handleExecute}
            disabled={isRunning}
            className="flex items-center gap-2 px-4 py-2 bg-blue-600 hover:bg-blue-700 disabled:bg-blue-400 text-white rounded-lg font-medium transition-colors"
          >
            <Play className="w-4 h-4" />
            {isRunning ? 'Running...' : 'Run'}
          </button>
        </div>
      </div>

      {/* Editor */}
      <div className="h-80 border border-gray-200 dark:border-gray-700 rounded-lg overflow-hidden">
        <MonacoEditor
          value={code}
          onChange={(v) => setCode(v || '')}
          language="go"
          height="100%"
          theme="vs-dark"
        />
      </div>

      {/* Output */}
      {(output || error) && (
        <OutputConsole result={output} error={error} />
      )}
    </div>
  );
}
