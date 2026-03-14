'use client';

import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Textarea } from '@/components/ui/textarea';
import { Card } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Switch } from '@/components/ui/switch';
import { Plus, Trash2, Eye, EyeOff } from 'lucide-react';
import { TestCase } from '@/types/assessment';

interface TestCaseBuilderProps {
  testCases: Omit<TestCase, 'id'>[];
  onChange: (testCases: Omit<TestCase, 'id'>[]) => void;
  readonly?: boolean;
}

export function TestCaseBuilder({ testCases, onChange, readonly = false }: TestCaseBuilderProps) {
  const addTestCase = () => {
    const newTestCase: Omit<TestCase, 'id'> = {
      input: '',
      expectedOutput: '',
      isVisible: true,
      points: 10,
    };
    onChange([...testCases, newTestCase]);
  };

  const updateTestCase = (index: number, updates: Partial<Omit<TestCase, 'id'>>) => {
    const newTestCases = [...testCases];
    newTestCases[index] = { ...newTestCases[index], ...updates };
    onChange(newTestCases);
  };

  const deleteTestCase = (index: number) => {
    onChange(testCases.filter((_, i) => i !== index));
  };

  const visibleCases = testCases.filter((tc) => tc.isVisible).length;
  const hiddenCases = testCases.length - visibleCases;
  const totalPoints = testCases.reduce((sum, tc) => sum + tc.points, 0);

  return (
    <div className="space-y-4">
      <div className="flex items-center justify-between">
        <div>
          <h3 className="text-lg font-semibold">Test Cases</h3>
          <p className="text-sm text-gray-500">
            {visibleCases} visible, {hiddenCases} hidden • {totalPoints} total points
          </p>
        </div>
        {!readonly && (
          <Button onClick={addTestCase} size="sm">
            <Plus className="mr-2 h-4 w-4" />
            Add Test Case
          </Button>
        )}
      </div>

      <div className="space-y-3">
        {testCases.map((testCase, index) => (
          <TestCaseCard
            key={index}
            testCase={testCase}
            index={index}
            onUpdate={(updates) => updateTestCase(index, updates)}
            onDelete={() => deleteTestCase(index)}
            readonly={readonly}
          />
        ))}
      </div>

      {testCases.length === 0 && (
        <div className="text-center py-12 text-gray-500 border-2 border-dashed rounded-lg">
          <p>No test cases yet. Add test cases to define how the code will be evaluated.</p>
        </div>
      )}

      {!readonly && testCases.length > 0 && (
        <Card className="p-4 bg-blue-50 border-blue-200">
          <p className="text-sm text-blue-900">
            <strong>Tip:</strong> Visible test cases help students understand what's expected.
            Hidden test cases are used for grading and prevent hardcoded solutions.
          </p>
        </Card>
      )}
    </div>
  );
}

interface TestCaseCardProps {
  testCase: Omit<TestCase, 'id'>;
  index: number;
  onUpdate: (updates: Partial<Omit<TestCase, 'id'>>) => void;
  onDelete: () => void;
  readonly?: boolean;
}

function TestCaseCard({ testCase, index, onUpdate, onDelete, readonly }: TestCaseCardProps) {
  return (
    <Card className={`p-4 ${testCase.isVisible ? 'border-gray-200' : 'border-yellow-200 bg-yellow-50'}`}>
      <div className="space-y-4">
        <div className="flex items-center justify-between">
          <div className="flex items-center gap-2">
            <Badge variant={testCase.isVisible ? 'outline' : 'warning'}>
              Test Case {index + 1}
            </Badge>
            {!testCase.isVisible && (
              <Badge variant="warning" className="text-xs">
                <EyeOff className="mr-1 h-3 w-3" />
                Hidden
              </Badge>
            )}
          </div>

          <div className="flex items-center gap-2">
            <span className="text-sm text-gray-500">{testCase.points} pts</span>
            {!readonly && (
              <Button variant="ghost" size="sm" onClick={onDelete}>
                <Trash2 className="h-4 w-4" />
              </Button>
            )}
          </div>
        </div>

        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div>
            <Label>Input</Label>
            <Textarea
              value={testCase.input}
              onChange={(e) => onUpdate({ input: e.target.value })}
              placeholder="Input values for the test case..."
              className="mt-2 font-mono text-sm"
              rows={3}
              disabled={readonly}
            />
            <p className="text-xs text-gray-500 mt-1">
              This will be provided to the program via stdin
            </p>
          </div>

          <div>
            <Label>Expected Output</Label>
            <Textarea
              value={testCase.expectedOutput}
              onChange={(e) => onUpdate({ expectedOutput: e.target.value })}
              placeholder="Expected output from the program..."
              className="mt-2 font-mono text-sm"
              rows={3}
              disabled={readonly}
            />
            <p className="text-xs text-gray-500 mt-1">
              This is what the program should output to stdout
            </p>
          </div>
        </div>

        <div className="flex items-center justify-between">
          <div className="flex items-center gap-2">
            <Switch
              checked={testCase.isVisible}
              onCheckedChange={(checked) => onUpdate({ isVisible: checked })}
              disabled={readonly}
            />
            <Label className="text-sm">
              {testCase.isVisible ? (
                <>
                  <Eye className="inline h-4 w-4 mr-1" />
                  Visible to students
                </>
              ) : (
                <>
                  <EyeOff className="inline h-4 w-4 mr-1" />
                  Hidden (graded only)
                </>
              )}
            </Label>
          </div>

          <div className="flex items-center gap-2">
            <Label className="text-sm">Points:</Label>
            <Input
              type="number"
              min="1"
              value={testCase.points}
              onChange={(e) => onUpdate({ points: parseInt(e.target.value) || 1 })}
              className="w-20"
              disabled={readonly}
            />
          </div>
        </div>
      </div>
    </Card>
  );
}
