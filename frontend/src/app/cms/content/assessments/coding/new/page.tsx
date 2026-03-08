'use client';

import { useState } from 'react';
import { useRouter } from 'next/navigation';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Textarea } from '@/components/ui/textarea';
import { Card } from '@/components/ui/card';
import { TestCaseBuilder } from '@/components/assessments/test-case-builder';
import MonacoCodeEditor from '@/components/learning/monaco-code-editor';
import { CreateCodingExerciseRequest, ExecutionLimits, TestCase } from '@/types/assessment';
import { ArrowLeft, Save, Eye, Plus } from 'lucide-react';
import Link from 'next/link';

export default function NewCodingExercisePage() {
  const router = useRouter();
  const [isLoading, setIsLoading] = useState(false);
  const [showPreview, setShowPreview] = useState(false);

  const [title, setTitle] = useState('');
  const [description, setDescription] = useState('');
  const [starterCode, setStarterCode] = useState('package main\n\nimport "fmt"\n\nfunc main() {\n\t// Your code here\n\tfmt.Println("Hello, World!")\n}\n');
  const [testCases, setTestCases] = useState<Omit<TestCase, 'id'>[]>([]);
  const [limits, setLimits] = useState<ExecutionLimits>({
    timeLimitSeconds: 5,
    memoryLimitMB: 128,
  });
  const [hints, setHints] = useState<string[]>([]);
  const [solutionTemplate, setSolutionTemplate] = useState('');
  const [newHint, setNewHint] = useState('');

  const addHint = () => {
    if (newHint.trim()) {
      setHints([...hints, newHint.trim()]);
      setNewHint('');
    }
  };

  const removeHint = (index: number) => {
    setHints(hints.filter((_, i) => i !== index));
  };

  const handleSubmit = async (publish: boolean = false) => {
    if (!title.trim()) {
      alert('Please enter an exercise title');
      return;
    }

    if (testCases.length === 0) {
      alert('Please add at least one test case');
      return;
    }

    setIsLoading(true);

    try {
      const payload: CreateCodingExerciseRequest = {
        lessonId: '',
        title: title.trim(),
        description: description.trim(),
        starterCode,
        testCases,
        limits,
        hints,
        solutionTemplate: solutionTemplate || undefined,
        passingScore: 70,
      };

      console.log('Creating coding exercise:', payload);
      alert(publish ? 'Exercise published!' : 'Exercise saved as draft!');
      router.push('/cms/content/assessments');
    } catch (error) {
      console.error('Error creating exercise:', error);
      alert('Failed to create exercise');
    } finally {
      setIsLoading(false);
    }
  };

  if (showPreview) {
    return (
      <div className="container mx-auto py-8">
        <div className="flex items-center gap-4 mb-6">
          <Button variant="ghost" onClick={() => setShowPreview(false)}>
            <ArrowLeft className="mr-2 h-4 w-4" />
            Back
          </Button>
          <h1 className="text-2xl font-bold">Exercise Preview</h1>
        </div>

        <Card className="p-8 max-w-4xl mx-auto">
          <h2 className="text-3xl font-bold mb-4">{title || 'Untitled Exercise'}</h2>

          {description && (
            <div className="prose max-w-none mb-6">
              <p className="text-gray-600">{description}</p>
            </div>
          )}

          <div className="mb-6">
            <h3 className="text-lg font-semibold mb-2">Starter Code</h3>
            <pre className="bg-gray-900 text-gray-100 p-4 rounded-lg overflow-x-auto">
              <code>{starterCode}</code>
            </pre>
          </div>

          <div className="mb-6">
            <h3 className="text-lg font-semibold mb-2">Test Cases</h3>
            <div className="space-y-2">
              {testCases.filter(tc => tc.isVisible).map((tc, i) => (
                <Card key={i} className="p-4">
                  <p className="text-sm font-medium">Test Case {i + 1} ({tc.points} pts)</p>
                  <div className="grid grid-cols-2 gap-4 mt-2">
                    <div>
                      <p className="text-xs text-gray-500">Input:</p>
                      <pre className="text-xs bg-gray-100 p-2 rounded">{tc.input || '(empty)'}</pre>
                    </div>
                    <div>
                      <p className="text-xs text-gray-500">Expected:</p>
                      <pre className="text-xs bg-gray-100 p-2 rounded">{tc.expectedOutput || '(empty)'}</pre>
                    </div>
                  </div>
                </Card>
              ))}
            </div>
          </div>

          <Card className="p-4 bg-gray-50">
            <p className="text-sm"><strong>Time Limit:</strong> {limits.timeLimitSeconds}s</p>
            <p className="text-sm"><strong>Memory Limit:</strong> {limits.memoryLimitMB}MB</p>
            <p className="text-sm"><strong>Test Cases:</strong> {testCases.length} total</p>
          </Card>
        </Card>
      </div>
    );
  }

  return (
    <div className="container mx-auto py-8 max-w-6xl">
      <div className="flex items-center gap-4 mb-8">
        <Link href="/cms/content/assessments">
          <Button variant="ghost" size="sm">
            <ArrowLeft className="mr-2 h-4 w-4" />
            Back
          </Button>
        </Link>
        <div>
          <h1 className="text-3xl font-bold">Create Coding Exercise</h1>
          <p className="text-gray-600 mt-1">Define a programming problem with automated tests</p>
        </div>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-8">
        <div className="space-y-6">
          <Card className="p-6">
            <h2 className="text-xl font-semibold mb-4">Basic Information</h2>
            <div className="space-y-4">
              <div>
                <Label>Exercise Title *</Label>
                <Input
                  value={title}
                  onChange={(e) => setTitle(e.target.value)}
                  placeholder="e.g., Reverse a String"
                  className="mt-2"
                />
              </div>
              <div>
                <Label>Problem Description</Label>
                <Textarea
                  value={description}
                  onChange={(e) => setDescription(e.target.value)}
                  placeholder="Describe the problem, requirements, examples..."
                  className="mt-2"
                  rows={6}
                />
              </div>
            </div>
          </Card>

          <Card className="p-6">
            <h2 className="text-xl font-semibold mb-4">Execution Limits</h2>
            <div className="grid grid-cols-2 gap-4">
              <div>
                <Label>Time Limit (seconds)</Label>
                <Input
                  type="number"
                  min="1"
                  max="30"
                  value={limits.timeLimitSeconds}
                  onChange={(e) => setLimits({ ...limits, timeLimitSeconds: parseInt(e.target.value) || 5 })}
                  className="mt-2"
                />
                <p className="text-xs text-gray-500 mt-1">1-30 seconds</p>
              </div>
              <div>
                <Label>Memory Limit (MB)</Label>
                <Input
                  type="number"
                  min="64"
                  max="256"
                  value={limits.memoryLimitMB}
                  onChange={(e) => setLimits({ ...limits, memoryLimitMB: parseInt(e.target.value) || 128 })}
                  className="mt-2"
                />
                <p className="text-xs text-gray-500 mt-1">64-256 MB</p>
              </div>
            </div>
          </Card>

          <Card className="p-6">
            <h2 className="text-xl font-semibold mb-4">Hints (Optional)</h2>
            <div className="space-y-3">
              <div className="flex gap-2">
                <Input
                  value={newHint}
                  onChange={(e) => setNewHint(e.target.value)}
                  placeholder="Add a hint..."
                  onKeyPress={(e) => e.key === 'Enter' && addHint()}
                />
                <Button onClick={addHint} size="sm">
                  <Plus className="h-4 w-4" />
                </Button>
              </div>
              <div className="space-y-2">
                {hints.map((hint, i) => (
                  <div key={i} className="flex items-center justify-between bg-gray-50 p-2 rounded">
                    <span className="text-sm">{hint}</span>
                    <Button variant="ghost" size="sm" onClick={() => removeHint(i)}>
                      ×
                    </Button>
                  </div>
                ))}
              </div>
            </div>
          </Card>
        </div>

        <div className="space-y-6">
          <Card className="p-6">
            <h2 className="text-xl font-semibold mb-4">Starter Code</h2>
            <div className="border rounded-lg overflow-hidden">
              <MonacoCodeEditor
                code={starterCode}
                onChange={setStarterCode}
                language="go"
                height="400px"
              />
            </div>
          </Card>

          <Card className="p-6">
            <TestCaseBuilder testCases={testCases} onChange={setTestCases} />
          </Card>

          <Card className="p-6">
            <h2 className="text-xl font-semibold mb-4">Solution Template (Optional)</h2>
            <p className="text-sm text-gray-500 mb-4">Reference solution (not visible to students)</p>
            <Textarea
              value={solutionTemplate}
              onChange={(e) => setSolutionTemplate(e.target.value)}
              placeholder="Your solution code..."
              className="font-mono text-sm"
              rows={8}
            />
          </Card>

          <div className="space-y-3">
            <Button
              onClick={() => setShowPreview(true)}
              variant="outline"
              className="w-full"
              disabled={testCases.length === 0}
            >
              <Eye className="mr-2 h-4 w-4" />
              Preview
            </Button>
            <Button
              onClick={() => handleSubmit(false)}
              variant="outline"
              className="w-full"
              disabled={isLoading || !title.trim() || testCases.length === 0}
            >
              <Save className="mr-2 h-4 w-4" />
              Save Draft
            </Button>
            <Button
              onClick={() => handleSubmit(true)}
              className="w-full"
              disabled={isLoading || !title.trim() || testCases.length === 0}
            >
              Publish Exercise
            </Button>
          </div>
        </div>
      </div>
    </div>
  );
}
