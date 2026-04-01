'use client';

import React from 'react';
import { CheckCircle, XCircle, Terminal, AlertCircle } from 'lucide-react';

interface TestResult {
  test_name: string;
  passed: boolean;
  expected: string;
  actual: string;
  error?: string;
}

interface ExecuteResult {
  passed: boolean;
  score: number;
  results: TestResult[];
  execution_time: string;
}

interface OutputConsoleProps {
  result?: ExecuteResult | null;
  error?: string | null;
}

export default function OutputConsole({ result, error }: OutputConsoleProps) {
  if (error) {
    return (
      <div className="bg-red-900/20 border border-red-500/50 rounded-lg p-4">
        <div className="flex items-center gap-2 text-red-400 mb-2">
          <AlertCircle className="w-5 h-5" />
          <span className="font-semibold">Execution Error</span>
        </div>
        <pre className="font-mono text-sm text-red-300 whitespace-pre-wrap">
          {error}
        </pre>
      </div>
    );
  }

  if (!result) return null;

  return (
    <div className="bg-gray-900 border border-gray-700 rounded-lg overflow-hidden">
      {/* Header */}
      <div className="flex items-center justify-between px-4 py-2 bg-gray-800 border-b border-gray-700">
        <div className="flex items-center gap-2">
          <Terminal className="w-4 h-4 text-gray-400" />
          <span className="font-mono text-sm text-gray-300">Output</span>
        </div>
        <div className="flex items-center gap-4">
          <span className={`text-sm font-medium ${result.passed ? 'text-green-400' : 'text-red-400'}`}>
            {result.score}% passed
          </span>
          <span className="text-xs text-gray-500">
            {result.execution_time}
          </span>
        </div>
      </div>

      {/* Results */}
      <div className="p-4 space-y-3">
        {result.results.map((r, i) => (
          <div
            key={i}
            className={`p-3 rounded-lg border ${
              r.passed
                ? 'bg-green-900/20 border-green-500/30'
                : 'bg-red-900/20 border-red-500/30'
            }`}
          >
            <div className="flex items-center gap-2 mb-2">
              {r.passed ? (
                <CheckCircle className="w-4 h-4 text-green-400" />
              ) : (
                <XCircle className="w-4 h-4 text-red-400" />
              )}
              <span className="font-medium text-sm">{r.test_name}</span>
            </div>

            {!r.passed && (
              <div className="ml-6 space-y-1 text-sm">
                <div>
                  <span className="text-gray-400">Expected: </span>
                  <code className="text-green-300">{r.expected}</code>
                </div>
                <div>
                  <span className="text-gray-400">Actual: </span>
                  <code className="text-red-300">{r.actual}</code>
                </div>
                {r.error && (
                  <div className="text-red-400 mt-1">{r.error}</div>
                )}
              </div>
            )}
          </div>
        ))}
      </div>
    </div>
  );
}
