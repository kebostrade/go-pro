'use client';

import React, { useState, useRef, useEffect, useCallback } from 'react';
import Editor, { OnMount } from '@monaco-editor/react';
import * as monaco from 'monaco-editor';
import { Play, Send, RotateCcw, Copy, Maximize2, Minimize2, Check, X, Clock } from 'lucide-react';

// Types
interface TestCase {
  name: string;
  input?: string;
  expected: string;
}

interface TestResult {
  name: string;
  passed: boolean;
  expected: string;
  actual: string;
  error?: string;
  execution_time_ms?: number;
}

interface ExerciseResult {
  success: boolean;
  passed: boolean;
  score: number;
  results: TestResult[];
  execution_time_ms: number;
  message: string;
}

interface MonacoCodeEditorProps {
  initialCode: string;
  exerciseId: string;
  language?: 'go' | 'javascript' | 'python';
  height?: string;
  readOnly?: boolean;
  onChange?: (code: string) => void;
  onSubmit?: (code: string) => Promise<ExerciseResult | void>;
  testCases?: TestCase[];
}

// Test Results Display Component
const TestResultsDisplay: React.FC<{ results: ExerciseResult }> = ({ results }) => {
  return (
    <div className="mt-4 border border-gray-200 rounded-lg overflow-hidden">
      {/* Header */}
      <div className={`px-4 py-3 ${results.passed ? 'bg-green-50' : 'bg-red-50'}`}>
        <div className="flex items-center justify-between">
          <div className="flex items-center gap-2">
            {results.passed ? (
              <Check className="w-5 h-5 text-green-600" />
            ) : (
              <X className="w-5 h-5 text-red-600" />
            )}
            <h3 className={`font-semibold ${results.passed ? 'text-green-900' : 'text-red-900'}`}>
              {results.message}
            </h3>
          </div>
          <div className="flex items-center gap-4 text-sm">
            <span className={`font-medium ${results.passed ? 'text-green-700' : 'text-red-700'}`}>
              Score: {results.score}%
            </span>
            <span className="flex items-center gap-1 text-gray-600">
              <Clock className="w-4 h-4" />
              {results.execution_time_ms}ms
            </span>
          </div>
        </div>
      </div>

      {/* Test Cases */}
      <div className="divide-y divide-gray-200">
        {results.results.map((test, index) => (
          <div key={index} className="px-4 py-3 bg-white">
            <div className="flex items-start gap-3">
              {/* Status Icon */}
              <div className="mt-1">
                {test.passed ? (
                  <Check className="w-5 h-5 text-green-600" />
                ) : (
                  <X className="w-5 h-5 text-red-600" />
                )}
              </div>

              {/* Test Details */}
              <div className="flex-1 min-w-0">
                <div className="flex items-center justify-between mb-2">
                  <h4 className="font-medium text-gray-900">{test.name}</h4>
                  {test.execution_time_ms && (
                    <span className="text-xs text-gray-500">{test.execution_time_ms}ms</span>
                  )}
                </div>

                {/* Expected vs Actual */}
                {!test.passed && (
                  <div className="space-y-2 text-sm">
                    <div>
                      <span className="font-medium text-gray-700">Expected:</span>
                      <pre className="mt-1 p-2 bg-gray-50 rounded border border-gray-200 overflow-x-auto">
                        <code className="text-xs">{test.expected}</code>
                      </pre>
                    </div>
                    <div>
                      <span className="font-medium text-gray-700">Actual:</span>
                      <pre className="mt-1 p-2 bg-red-50 rounded border border-red-200 overflow-x-auto">
                        <code className="text-xs text-red-900">{test.actual}</code>
                      </pre>
                    </div>
                    {test.error && (
                      <div>
                        <span className="font-medium text-red-700">Error:</span>
                        <pre className="mt-1 p-2 bg-red-50 rounded border border-red-200 overflow-x-auto">
                          <code className="text-xs text-red-900">{test.error}</code>
                        </pre>
                      </div>
                    )}
                  </div>
                )}
              </div>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
};

// Main Monaco Code Editor Component
const MonacoCodeEditor: React.FC<MonacoCodeEditorProps> = ({
  initialCode,
  exerciseId,
  language = 'go',
  height = '500px',
  readOnly = false,
  onChange,
  onSubmit,
  testCases,
}) => {
  const [code, setCode] = useState(initialCode);
  const [theme, setTheme] = useState<'vs-dark' | 'light'>('vs-dark');
  const [fontSize, setFontSize] = useState(14);
  const [isFullscreen, setIsFullscreen] = useState(false);
  const [isRunning, setIsRunning] = useState(false);
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [results, setResults] = useState<ExerciseResult | null>(null);
  const [copySuccess, setCopySuccess] = useState(false);

  const editorRef = useRef<monaco.editor.IStandaloneCodeEditor | null>(null);
  const containerRef = useRef<HTMLDivElement>(null);

  // Local storage key
  const storageKey = `exercise-${exerciseId}-code`;

  // Load code from localStorage on mount
  useEffect(() => {
    const savedCode = localStorage.getItem(storageKey);
    if (savedCode) {
      setCode(savedCode);
    }
  }, [storageKey]);

  // Auto-save to localStorage
  useEffect(() => {
    const timeoutId = setTimeout(() => {
      localStorage.setItem(storageKey, code);
    }, 1000);
    return () => clearTimeout(timeoutId);
  }, [code, storageKey]);

  // Editor mount handler
  const handleEditorDidMount: OnMount = (editor, monaco) => {
    editorRef.current = editor;

    // Keyboard shortcuts
    editor.addCommand(monaco.KeyMod.CtrlCmd | monaco.KeyCode.Enter, () => {
      handleRun();
    });

    editor.addCommand(monaco.KeyMod.CtrlCmd | monaco.KeyMod.Shift | monaco.KeyCode.Enter, () => {
      handleSubmit();
    });
  };

  // Handle code change
  const handleCodeChange = (value: string | undefined) => {
    if (value !== undefined) {
      setCode(value);
      onChange?.(value);
    }
  };

  // Reset code to initial
  const handleReset = useCallback(() => {
    setCode(initialCode);
    localStorage.removeItem(storageKey);
    setResults(null);
  }, [initialCode, storageKey]);

  // Copy code to clipboard
  const handleCopy = useCallback(async () => {
    try {
      await navigator.clipboard.writeText(code);
      setCopySuccess(true);
      setTimeout(() => setCopySuccess(false), 2000);
    } catch (error) {
      console.error('Failed to copy code:', error);
    }
  }, [code]);

  // Toggle fullscreen
  const toggleFullscreen = useCallback(() => {
    if (!document.fullscreenElement) {
      containerRef.current?.requestFullscreen();
      setIsFullscreen(true);
    } else {
      document.exitFullscreen();
      setIsFullscreen(false);
    }
  }, []);

  // Handle fullscreen change
  useEffect(() => {
    const handleFullscreenChange = () => {
      setIsFullscreen(!!document.fullscreenElement);
    };
    document.addEventListener('fullscreenchange', handleFullscreenChange);
    return () => document.removeEventListener('fullscreenchange', handleFullscreenChange);
  }, []);

  // Run code (test locally)
  const handleRun = useCallback(async () => {
    if (!onSubmit || isRunning) return;

    setIsRunning(true);
    setResults(null);

    try {
      const result = await onSubmit(code);
      if (result) setResults(result);
    } catch (error) {
      console.error('Error running code:', error);
      setResults({
        success: false,
        passed: false,
        score: 0,
        results: [],
        execution_time_ms: 0,
        message: error instanceof Error ? error.message : 'Failed to run code',
      });
    } finally {
      setIsRunning(false);
    }
  }, [code, onSubmit, isRunning]);

  // Submit code for grading
  const handleSubmit = useCallback(async () => {
    if (!onSubmit || isSubmitting) return;

    setIsSubmitting(true);
    setResults(null);

    try {
      const result = await onSubmit(code);
      if (result) {
        setResults(result);
        // Clear localStorage on successful submission
        if (result.passed) {
          localStorage.removeItem(storageKey);
        }
      }
    } catch (error) {
      console.error('Error submitting code:', error);
      setResults({
        success: false,
        passed: false,
        score: 0,
        results: [],
        execution_time_ms: 0,
        message: error instanceof Error ? error.message : 'Failed to submit code',
      });
    } finally {
      setIsSubmitting(false);
    }
  }, [code, onSubmit, isSubmitting, storageKey]);

  // Editor options
  const editorOptions: monaco.editor.IStandaloneEditorConstructionOptions = {
    minimap: { enabled: !isFullscreen },
    fontSize,
    lineNumbers: 'on',
    scrollBeyondLastLine: false,
    automaticLayout: true,
    tabSize: language === 'go' ? 4 : 2,
    wordWrap: 'on',
    formatOnPaste: true,
    formatOnType: true,
    readOnly,
    scrollbar: {
      vertical: 'visible',
      horizontal: 'visible',
    },
    suggestOnTriggerCharacters: true,
    quickSuggestions: true,
    folding: true,
    renderWhitespace: 'selection',
  };

  return (
    <div
      ref={containerRef}
      className={`monaco-editor-container ${isFullscreen ? 'fixed inset-0 z-50 bg-white' : ''}`}
    >
      {/* Toolbar */}
      <div className="flex items-center justify-between px-4 py-2 bg-gray-100 border-b border-gray-200">
        <div className="flex items-center gap-2">
          {/* Theme Toggle */}
          <select
            value={theme}
            onChange={(e) => setTheme(e.target.value as 'vs-dark' | 'light')}
            className="px-3 py-1 text-sm border border-gray-300 rounded hover:border-gray-400 focus:outline-none focus:ring-2 focus:ring-blue-500"
            aria-label="Select theme"
          >
            <option value="vs-dark">Dark Theme</option>
            <option value="light">Light Theme</option>
          </select>

          {/* Font Size */}
          <select
            value={fontSize}
            onChange={(e) => setFontSize(Number(e.target.value))}
            className="px-3 py-1 text-sm border border-gray-300 rounded hover:border-gray-400 focus:outline-none focus:ring-2 focus:ring-blue-500"
            aria-label="Select font size"
          >
            <option value={12}>12px</option>
            <option value={14}>14px</option>
            <option value={16}>16px</option>
            <option value={18}>18px</option>
            <option value={20}>20px</option>
          </select>
        </div>

        <div className="flex items-center gap-2">
          {/* Reset Button */}
          <button
            onClick={handleReset}
            className="flex items-center gap-1 px-3 py-1 text-sm text-gray-700 hover:bg-gray-200 rounded transition-colors"
            aria-label="Reset code"
          >
            <RotateCcw className="w-4 h-4" />
            Reset
          </button>

          {/* Copy Button */}
          <button
            onClick={handleCopy}
            className="flex items-center gap-1 px-3 py-1 text-sm text-gray-700 hover:bg-gray-200 rounded transition-colors"
            aria-label="Copy code"
          >
            <Copy className="w-4 h-4" />
            {copySuccess ? 'Copied!' : 'Copy'}
          </button>

          {/* Fullscreen Toggle */}
          <button
            onClick={toggleFullscreen}
            className="flex items-center gap-1 px-3 py-1 text-sm text-gray-700 hover:bg-gray-200 rounded transition-colors"
            aria-label="Toggle fullscreen"
          >
            {isFullscreen ? (
              <>
                <Minimize2 className="w-4 h-4" />
                Exit
              </>
            ) : (
              <>
                <Maximize2 className="w-4 h-4" />
                Fullscreen
              </>
            )}
          </button>
        </div>
      </div>

      {/* Monaco Editor */}
      <div className={isFullscreen ? 'h-[calc(100vh-120px)]' : ''}>
        <Editor
          height={isFullscreen ? '100%' : height}
          language={language}
          value={code}
          onChange={handleCodeChange}
          theme={theme}
          options={editorOptions}
          onMount={handleEditorDidMount}
          loading={
            <div className="flex items-center justify-center h-full">
              <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
            </div>
          }
        />
      </div>

      {/* Action Buttons */}
      <div className="flex items-center justify-between px-4 py-3 bg-gray-50 border-t border-gray-200">
        <div className="text-sm text-gray-600">
          Press <kbd className="px-2 py-1 bg-white border border-gray-300 rounded text-xs">Ctrl+Enter</kbd> to run,{' '}
          <kbd className="px-2 py-1 bg-white border border-gray-300 rounded text-xs">Ctrl+Shift+Enter</kbd> to submit
        </div>
        <div className="flex items-center gap-2">
          {/* Run Button */}
          <button
            onClick={handleRun}
            disabled={isRunning || isSubmitting || !onSubmit}
            className="flex items-center gap-2 px-4 py-2 text-sm font-medium text-white bg-blue-600 rounded hover:bg-blue-700 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
            aria-label="Run code"
          >
            {isRunning ? (
              <>
                <div className="animate-spin rounded-full h-4 w-4 border-2 border-white border-t-transparent"></div>
                Running...
              </>
            ) : (
              <>
                <Play className="w-4 h-4" />
                Run Code
              </>
            )}
          </button>

          {/* Submit Button */}
          <button
            onClick={handleSubmit}
            disabled={isRunning || isSubmitting || !onSubmit}
            className="flex items-center gap-2 px-4 py-2 text-sm font-medium text-white bg-green-600 rounded hover:bg-green-700 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
            aria-label="Submit code"
          >
            {isSubmitting ? (
              <>
                <div className="animate-spin rounded-full h-4 w-4 border-2 border-white border-t-transparent"></div>
                Submitting...
              </>
            ) : (
              <>
                <Send className="w-4 h-4" />
                Submit
              </>
            )}
          </button>
        </div>
      </div>

      {/* Test Results */}
      {results && (
        <div className="px-4 pb-4">
          <TestResultsDisplay results={results} />
        </div>
      )}
    </div>
  );
};

export default MonacoCodeEditor;
