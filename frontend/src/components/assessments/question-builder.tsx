'use client';

import { useState } from 'react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Textarea } from '@/components/ui/textarea';
import { Card } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
import { Plus, Trash2, GripVertical } from 'lucide-react';
import { Question, QuestionType } from '@/types/assessment';
import clsx from 'clsx';

interface QuestionBuilderProps {
  questions: Omit<Question, 'id'>[];
  onChange: (questions: Omit<Question, 'id'>[]) => void;
  readonly?: boolean;
}

export function QuestionBuilder({ questions, onChange, readonly = false }: QuestionBuilderProps) {
  const [selectedType, setSelectedType] = useState<QuestionType>('multiple_choice');

  const addQuestion = () => {
    const newQuestion: Omit<Question, 'id'> = {
      type: selectedType,
      questionText: '',
      points: 1,
      orderIndex: questions.length,
      ...(selectedType === 'multiple_choice' && {
        options: ['', '', '', ''],
        correctAnswer: 0,
      }),
      ...(selectedType === 'true_false' && {
        correctAnswer: false,
      }),
      ...(selectedType === 'short_answer' && {
        correctAnswer: '',
      }),
      ...(selectedType === 'code_completion' && {
        correctAnswer: '',
      }),
    };
    onChange([...questions, newQuestion]);
  };

  const updateQuestion = (index: number, updates: Partial<Omit<Question, 'id'>>) => {
    const newQuestions = [...questions];
    newQuestions[index] = { ...newQuestions[index], ...updates };
    onChange(newQuestions);
  };

  const deleteQuestion = (index: number) => {
    onChange(questions.filter((_, i) => i !== index));
  };

  return (
    <div className="space-y-6">
      {!readonly && (
        <div className="flex items-center justify-between">
          <Tabs value={selectedType} onValueChange={(v) => setSelectedType(v as QuestionType)}>
            <TabsList>
              <TabsTrigger value="multiple_choice">Multiple Choice</TabsTrigger>
              <TabsTrigger value="true_false">True/False</TabsTrigger>
              <TabsTrigger value="short_answer">Short Answer</TabsTrigger>
              <TabsTrigger value="code_completion">Code Completion</TabsTrigger>
            </TabsList>
          </Tabs>
          <Button onClick={addQuestion} size="sm">
            <Plus className="mr-2 h-4 w-4" />
            Add {selectedType.replace('_', ' ')}
          </Button>
        </div>
      )}

      <div className="space-y-4">
        {questions.map((question, index) => (
          <QuestionCard
            key={index}
            question={question}
            index={index}
            onUpdate={(updates) => updateQuestion(index, updates)}
            onDelete={() => deleteQuestion(index)}
            readonly={readonly}
          />
        ))}
      </div>

      {questions.length === 0 && (
        <div className="text-center py-12 text-gray-500">
          <p>No questions yet. Add your first question to get started.</p>
        </div>
      )}
    </div>
  );
}

interface QuestionCardProps {
  question: Omit<Question, 'id'>;
  index: number;
  onUpdate: (updates: Partial<Omit<Question, 'id'>>) => void;
  onDelete: () => void;
  readonly?: boolean;
}

function QuestionCard({ question, index, onUpdate, onDelete }: QuestionCardProps) {
  const [isExpanded, setIsExpanded] = useState(true);

  return (
    <Card className={clsx('p-6', !isExpanded && 'opacity-75')}>
      <div className="flex items-start gap-4">
        <div className="flex items-center gap-2 mt-2">
          <GripVertical className="h-5 w-5 text-gray-400 cursor-move" />
          <span className="text-sm font-medium text-gray-500">Q{index + 1}</span>
        </div>

        <div className="flex-1 space-y-4">
          <div className="flex items-center justify-between">
            <Badge variant="outline">{question.type.replace('_', ' ')}</Badge>
            <div className="flex items-center gap-2">
              <span className="text-sm text-gray-500">{question.points} pts</span>
              {!question.readonly && (
                <>
                  <Button
                    variant="ghost"
                    size="sm"
                    onClick={() => setIsExpanded(!isExpanded)}
                  >
                    {isExpanded ? 'Collapse' : 'Expand'}
                  </Button>
                  <Button
                    variant="ghost"
                    size="sm"
                    onClick={onDelete}
                  >
                    <Trash2 className="h-4 w-4" />
                  </Button>
                </>
              )}
            </div>
          </div>

          {isExpanded && (
            <>
              <div>
                <Label>Question</Label>
                <Textarea
                  value={question.questionText}
                  onChange={(e) => onUpdate({ questionText: e.target.value })}
                  placeholder="Enter your question..."
                  className="mt-2"
                  rows={3}
                  disabled={question.readonly}
                />
              </div>

              {question.type === 'multiple_choice' && question.options && (
                <MultipleChoiceEditor
                  options={question.options}
                  correctAnswer={question.correctAnswer as number}
                  onUpdate={(updates) => onUpdate(updates)}
                  readonly={question.readonly}
                />
              )}

              {question.type === 'true_false' && (
                <TrueFalseEditor
                  correctAnswer={question.correctAnswer as boolean}
                  onUpdate={(updates) => onUpdate(updates)}
                  readonly={question.readonly}
                />
              )}

              {question.type === 'code_completion' && (
                <CodeCompletionEditor
                  correctAnswer={question.correctAnswer as string}
                  onUpdate={(updates) => onUpdate(updates)}
                  readonly={question.readonly}
                />
              )}

              <div>
                <Label>Explanation (Optional)</Label>
                <Textarea
                  value={question.explanation || ''}
                  onChange={(e) => onUpdate({ explanation: e.target.value })}
                  placeholder="Explanation shown after answering..."
                  className="mt-2"
                  rows={2}
                  disabled={question.readonly}
                />
              </div>

              <div>
                <Label>Points</Label>
                <Input
                  type="number"
                  min="1"
                  value={question.points}
                  onChange={(e) => onUpdate({ points: parseInt(e.target.value) || 1 })}
                  className="mt-2 w-24"
                  disabled={question.readonly}
                />
              </div>
            </>
          )}
        </div>
      </div>
    </Card>
  );
}

function MultipleChoiceEditor({
  options,
  correctAnswer,
  onUpdate,
  readonly,
}: {
  options: string[];
  correctAnswer: number;
  onUpdate: (updates: { options?: string[]; correctAnswer?: number }) => void;
  readonly?: boolean;
}) {
  const updateOption = (index: number, value: string) => {
    const newOptions = [...options];
    newOptions[index] = value;
    onUpdate({ options: newOptions });
  };

  const addOption = () => {
    onUpdate({ options: [...options, ''] });
  };

  const removeOption = (index: number) => {
    if (options.length <= 2) return;
    const newOptions = options.filter((_, i) => i !== index);
    const newCorrectAnswer = correctAnswer === index ? 0 : correctAnswer > index ? correctAnswer - 1 : correctAnswer;
    onUpdate({ options: newOptions, correctAnswer: newCorrectAnswer });
  };

  return (
    <div className="space-y-2">
      <Label>Options (select correct answer)</Label>
      {options.map((option, index) => (
        <div key={index} className="flex items-center gap-2">
          <input
            type="radio"
            name={`correct-${Math.random()}`}
            checked={correctAnswer === index}
            onChange={() => onUpdate({ correctAnswer: index })}
            disabled={readonly}
            className="h-4 w-4"
          />
          <Input
            value={option}
            onChange={(e) => updateOption(index, e.target.value)}
            placeholder={`Option ${index + 1}`}
            disabled={readonly}
          />
          {!readonly && options.length > 2 && (
            <Button
              variant="ghost"
              size="sm"
              onClick={() => removeOption(index)}
            >
              <Trash2 className="h-4 w-4" />
            </Button>
          )}
        </div>
      ))}
      {!readonly && options.length < 6 && (
        <Button variant="outline" size="sm" onClick={addOption} className="mt-2">
          <Plus className="mr-2 h-4 w-4" />
          Add Option
        </Button>
      )}
    </div>
  );
}

function TrueFalseEditor({
  correctAnswer,
  onUpdate,
  readonly,
}: {
  correctAnswer: boolean;
  onUpdate: (updates: { correctAnswer?: boolean }) => void;
  readonly?: boolean;
}) {
  return (
    <div className="space-y-2">
      <Label>Correct Answer</Label>
      <div className="flex gap-4">
        <label className="flex items-center gap-2">
          <input
            type="radio"
            name={`tf-${Math.random()}`}
            checked={correctAnswer === true}
            onChange={() => onUpdate({ correctAnswer: true })}
            disabled={readonly}
            className="h-4 w-4"
          />
          <span>True</span>
        </label>
        <label className="flex items-center gap-2">
          <input
            type="radio"
            name={`tf-${Math.random()}`}
            checked={correctAnswer === false}
            onChange={() => onUpdate({ correctAnswer: false })}
            disabled={readonly}
            className="h-4 w-4"
          />
          <span>False</span>
        </label>
      </div>
    </div>
  );
}

function CodeCompletionEditor({
  correctAnswer,
  onUpdate,
  readonly,
}: {
  correctAnswer: string;
  onUpdate: (updates: { correctAnswer?: string }) => void;
  readonly?: boolean;
}) {
  return (
    <div>
      <Label>Correct Answer (Code)</Label>
      <Textarea
        value={correctAnswer}
        onChange={(e) => onUpdate({ correctAnswer: e.target.value })}
        placeholder="Enter the complete code..."
        className="mt-2 font-mono text-sm"
        rows={4}
        disabled={readonly}
      />
    </div>
  );
}
