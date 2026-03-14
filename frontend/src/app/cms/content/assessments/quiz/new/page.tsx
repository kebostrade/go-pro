'use client';

import { useState } from 'react';
import { useRouter } from 'next/navigation';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Textarea } from '@/components/ui/textarea';
import { Card } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Switch } from '@/components/ui/switch';
import { QuestionBuilder } from '@/components/assessments/question-builder';
import { CreateQuizRequest, QuizSettings, Question } from '@/types/assessment';
import { ArrowLeft, Save, Eye } from 'lucide-react';
import Link from 'next/link';

export default function NewQuizPage() {
  const router = useRouter();
  const [isLoading, setIsLoading] = useState(false);
  const [showPreview, setShowPreview] = useState(false);

  const [title, setTitle] = useState('');
  const [description, setDescription] = useState('');
  const [questions, setQuestions] = useState<Omit<Question, 'id'>[]>([]);
  const [settings, setSettings] = useState<QuizSettings>({
    shuffleQuestions: false,
    shuffleOptions: false,
    showExplanations: true,
    maxAttempts: 1,
    passingScore: 70,
  });
  const [timeLimitMinutes, setTimeLimitMinutes] = useState<number | undefined>(undefined);

  const handleSubmit = async (publish: boolean = false) => {
    if (!title.trim()) {
      alert('Please enter a quiz title');
      return;
    }

    if (questions.length === 0) {
      alert('Please add at least one question');
      return;
    }

    setIsLoading(true);

    try {
      const payload: CreateQuizRequest = {
        lessonId: '', // Will be set from context or route param
        title: title.trim(),
        description: description.trim(),
        questions,
        settings,
        passingScore: settings.passingScore,
        timeLimitMinutes,
      };

      // API call will be implemented
      console.log('Creating quiz:', payload);

      // const response = await fetch('/api/cms/assessments', {
      //   method: 'POST',
      //   headers: { 'Content-Type': 'application/json' },
      //   body: JSON.stringify({ type: 'quiz', ...payload }),
      // });

      // if (!response.ok) throw new Error('Failed to create quiz');

      alert(publish ? 'Quiz published successfully!' : 'Quiz saved as draft!');
      router.push('/cms/content/assessments');
    } catch (error) {
      console.error('Error creating quiz:', error);
      alert('Failed to create quiz. Please try again.');
    } finally {
      setIsLoading(false);
    }
  };

  const totalPoints = questions.reduce((sum, q) => sum + q.points, 0);

  if (showPreview) {
    return (
      <div className="container mx-auto py-8">
        <div className="flex items-center justify-between mb-6">
          <div className="flex items-center gap-4">
            <Button
              variant="ghost"
              onClick={() => setShowPreview(false)}
            >
              <ArrowLeft className="mr-2 h-4 w-4" />
              Back to Editor
            </Button>
            <h1 className="text-2xl font-bold">Quiz Preview</h1>
          </div>
          <Badge variant="outline">Preview Mode</Badge>
        </div>

        <Card className="p-8 max-w-3xl mx-auto">
          <h2 className="text-3xl font-bold mb-4">{title || 'Untitled Quiz'}</h2>
          {description && (
            <p className="text-gray-600 mb-6">{description}</p>
          )}

          {timeLimitMinutes && (
            <p className="text-sm text-gray-500 mb-6">
              Time Limit: {timeLimitMinutes} minutes
            </p>
          )}

          <div className="space-y-8">
            {questions.map((question, index) => (
              <div key={index} className="border-b pb-6">
                <div className="flex items-start justify-between mb-4">
                  <h3 className="text-lg font-semibold">
                    Question {index + 1}
                    <span className="text-sm font-normal text-gray-500 ml-2">
                      ({question.points} pts)
                    </span>
                  </h3>
                  <Badge variant="outline">{question.type}</Badge>
                </div>
                <p className="text-gray-800 mb-4">{question.questionText}</p>

                {question.type === 'multiple_choice' && question.options && (
                  <div className="space-y-2 ml-4">
                    {question.options.map((option, optIndex) => (
                      <div key={optIndex} className="flex items-center gap-2">
                        <input
                          type="radio"
                          name={`preview-q${index}`}
                          disabled
                          className="h-4 w-4"
                        />
                        <span>{option || '(empty option)'}</span>
                      </div>
                    ))}
                  </div>
                )}

                {question.type === 'true_false' && (
                  <div className="space-y-2 ml-4">
                    <div className="flex items-center gap-2">
                      <input type="radio" disabled className="h-4 w-4" />
                      <span>True</span>
                    </div>
                    <div className="flex items-center gap-2">
                      <input type="radio" disabled className="h-4 w-4" />
                      <span>False</span>
                    </div>
                  </div>
                )}

                {question.type === 'short_answer' && (
                  <Textarea
                    disabled
                    placeholder="Student's answer will appear here..."
                    className="ml-4"
                  />
                )}

                {question.type === 'code_completion' && (
                  <Textarea
                    disabled
                    placeholder="Student's code completion will appear here..."
                    className="ml-4 font-mono text-sm"
                  />
                )}
              </div>
            ))}
          </div>

          <div className="mt-8 p-4 bg-gray-50 rounded-lg">
            <p className="text-sm text-gray-600">
              <strong>Total Points:</strong> {totalPoints}
            </p>
            <p className="text-sm text-gray-600">
              <strong>Passing Score:</strong> {settings.passingScore}%
            </p>
            <p className="text-sm text-gray-600">
              <strong>Questions:</strong> {questions.length}
            </p>
          </div>
        </Card>
      </div>
    );
  }

  return (
    <div className="container mx-auto py-8 max-w-5xl">
      <div className="flex items-center gap-4 mb-8">
        <Link href="/cms/content/assessments">
          <Button variant="ghost" size="sm">
            <ArrowLeft className="mr-2 h-4 w-4" />
            Back
          </Button>
        </Link>
        <div>
          <h1 className="text-3xl font-bold">Create Quiz</h1>
          <p className="text-gray-600 mt-1">
            Build a quiz with multiple question types
          </p>
        </div>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
        <div className="lg:col-span-2 space-y-6">
          <Card className="p-6">
            <h2 className="text-xl font-semibold mb-4">Basic Information</h2>
            <div className="space-y-4">
              <div>
                <Label htmlFor="title">Quiz Title *</Label>
                <Input
                  id="title"
                  value={title}
                  onChange={(e) => setTitle(e.target.value)}
                  placeholder="e.g., Go Fundamentals Quiz"
                  className="mt-2"
                />
              </div>

              <div>
                <Label htmlFor="description">Description</Label>
                <Textarea
                  id="description"
                  value={description}
                  onChange={(e) => setDescription(e.target.value)}
                  placeholder="Describe what this quiz covers..."
                  className="mt-2"
                  rows={3}
                />
              </div>
            </div>
          </Card>

          <Card className="p-6">
            <h2 className="text-xl font-semibold mb-4">Questions</h2>
            <QuestionBuilder
              questions={questions}
              onChange={setQuestions}
            />
          </Card>
        </div>

        <div className="space-y-6">
          <Card className="p-6">
            <h2 className="text-xl font-semibold mb-4">Quiz Settings</h2>
            <div className="space-y-4">
              <div className="flex items-center justify-between">
                <Label htmlFor="shuffle-questions">Shuffle Questions</Label>
                <Switch
                  id="shuffle-questions"
                  checked={settings.shuffleQuestions}
                  onCheckedChange={(checked) =>
                    setSettings({ ...settings, shuffleQuestions: checked })
                  }
                />
              </div>

              <div className="flex items-center justify-between">
                <Label htmlFor="shuffle-options">Shuffle Answer Options</Label>
                <Switch
                  id="shuffle-options"
                  checked={settings.shuffleOptions}
                  onCheckedChange={(checked) =>
                    setSettings({ ...settings, shuffleOptions: checked })
                  }
                />
              </div>

              <div className="flex items-center justify-between">
                <Label htmlFor="show-explanations">Show Explanations</Label>
                <Switch
                  id="show-explanations"
                  checked={settings.showExplanations}
                  onCheckedChange={(checked) =>
                    setSettings({ ...settings, showExplanations: checked })
                  }
                />
              </div>

              <div>
                <Label htmlFor="passing-score">Passing Score (%)</Label>
                <Input
                  id="passing-score"
                  type="number"
                  min="0"
                  max="100"
                  value={settings.passingScore}
                  onChange={(e) =>
                    setSettings({ ...settings, passingScore: parseInt(e.target.value) || 0 })
                  }
                  className="mt-2"
                />
              </div>

              <div>
                <Label htmlFor="max-attempts">Maximum Attempts (Optional)</Label>
                <Input
                  id="max-attempts"
                  type="number"
                  min="1"
                  value={settings.maxAttempts || ''}
                  onChange={(e) =>
                    setSettings({
                      ...settings,
                      maxAttempts: e.target.value ? parseInt(e.target.value) : undefined,
                    })
                  }
                  placeholder="Unlimited"
                  className="mt-2"
                />
              </div>

              <div>
                <Label htmlFor="time-limit">Time Limit (minutes, optional)</Label>
                <Input
                  id="time-limit"
                  type="number"
                  min="1"
                  value={timeLimitMinutes || ''}
                  onChange={(e) =>
                    setTimeLimitMinutes(e.target.value ? parseInt(e.target.value) : undefined)
                  }
                  placeholder="No time limit"
                  className="mt-2"
                />
              </div>
            </div>
          </Card>

          <Card className="p-6 bg-blue-50 border-blue-200">
            <h3 className="font-semibold text-blue-900 mb-2">Quiz Summary</h3>
            <div className="space-y-2 text-sm text-blue-800">
              <p>Questions: {questions.length}</p>
              <p>Total Points: {totalPoints}</p>
              <p>Passing Score: {settings.passingScore}% ({Math.ceil(totalPoints * settings.passingScore / 100)} points)</p>
            </div>
          </Card>

          <div className="space-y-3">
            <Button
              onClick={() => setShowPreview(true)}
              variant="outline"
              className="w-full"
              disabled={questions.length === 0}
            >
              <Eye className="mr-2 h-4 w-4" />
              Preview Quiz
            </Button>
            <Button
              onClick={() => handleSubmit(false)}
              variant="outline"
              className="w-full"
              disabled={isLoading || !title.trim() || questions.length === 0}
            >
              <Save className="mr-2 h-4 w-4" />
              Save Draft
            </Button>
            <Button
              onClick={() => handleSubmit(true)}
              className="w-full"
              disabled={isLoading || !title.trim() || questions.length === 0}
            >
              Publish Quiz
            </Button>
          </div>
        </div>
      </div>
    </div>
  );
}
