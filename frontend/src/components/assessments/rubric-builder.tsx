'use client';

import { useState } from 'react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Textarea } from '@/components/ui/textarea';
import { Card } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Plus, Trash2, GripVertical } from 'lucide-react';
import { RubricCriterion, PerformanceLevel } from '@/types/assessment';
import clsx from 'clsx';

interface RubricBuilderProps {
  rubric: Omit<RubricCriterion, 'id'>[];
  onChange: (rubric: Omit<RubricCriterion, 'id'>[]) => void;
  readonly?: boolean;
}

export function RubricBuilder({ rubric, onChange, readonly = false }: RubricBuilderProps) {
  const addCriterion = () => {
    const newCriterion: Omit<RubricCriterion, 'id'> = {
      description: '',
      maxPoints: 10,
      levels: [
        { name: 'Excellent', description: '', points: 10 },
        { name: 'Good', description: '', points: 7 },
        { name: 'Fair', description: '', points: 4 },
        { name: 'Poor', description: '', points: 1 },
      ],
    };
    onChange([...rubric, newCriterion]);
  };

  const updateCriterion = (
    index: number,
    updates: Partial<Omit<RubricCriterion, 'id'>>
  ) => {
    const newRubric = [...rubric];
    newRubric[index] = { ...newRubric[index], ...updates };
    onChange(newRubric);
  };

  const deleteCriterion = (index: number) => {
    onChange(rubric.filter((_, i) => i !== index));
  };

  return (
    <div className="space-y-6">
      {!readonly && (
        <div className="flex justify-end">
          <Button onClick={addCriterion} size="sm">
            <Plus className="mr-2 h-4 w-4" />
            Add Criterion
          </Button>
        </div>
      )}

      <div className="space-y-4">
        {rubric.map((criterion, index) => (
          <CriterionCard
            key={index}
            criterion={criterion}
            index={index}
            onUpdate={(updates) => updateCriterion(index, updates)}
            onDelete={() => deleteCriterion(index)}
            readonly={readonly}
          />
        ))}
      </div>

      {rubric.length === 0 && (
        <div className="text-center py-12 text-gray-500">
          <p>No rubric criteria yet. Add criteria to define grading standards.</p>
        </div>
      )}

      {rubric.length > 0 && !readonly && (
        <Card className="p-4 bg-blue-50 border-blue-200">
          <p className="text-sm text-blue-900">
            <strong>Total Points:</strong> {rubric.reduce((sum, c) => sum + c.maxPoints, 0)}
          </p>
        </Card>
      )}
    </div>
  );
}

interface CriterionCardProps {
  criterion: Omit<RubricCriterion, 'id'>;
  index: number;
  onUpdate: (updates: Partial<Omit<RubricCriterion, 'id'>>) => void;
  onDelete: () => void;
  readonly?: boolean;
}

function CriterionCard({ criterion, index, onUpdate, onDelete, readonly }: CriterionCardProps) {
  const [isExpanded, setIsExpanded] = useState(true);

  const updateLevel = (levelIndex: number, updates: Partial<PerformanceLevel>) => {
    const newLevels = [...criterion.levels];
    newLevels[levelIndex] = { ...newLevels[levelIndex], ...updates };
    onUpdate({ levels: newLevels });
  };

  return (
    <Card className={clsx('p-6', !isExpanded && 'opacity-75')}>
      <div className="flex items-start gap-4">
        <div className="flex items-center gap-2 mt-2">
          <GripVertical className="h-5 w-5 text-gray-400 cursor-move" />
          <span className="text-sm font-medium text-gray-500">C{index + 1}</span>
        </div>

        <div className="flex-1 space-y-4">
          <div className="flex items-center justify-between">
            <Badge variant="outline">{criterion.maxPoints} points max</Badge>
            <div className="flex items-center gap-2">
              {!readonly && (
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
                <Label>Criterion Description</Label>
                <Textarea
                  value={criterion.description}
                  onChange={(e) => onUpdate({ description: e.target.value })}
                  placeholder="Describe what this criterion evaluates..."
                  className="mt-2"
                  rows={2}
                  disabled={readonly}
                />
              </div>

              <div>
                <Label>Maximum Points</Label>
                <Input
                  type="number"
                  min="1"
                  value={criterion.maxPoints}
                  onChange={(e) => onUpdate({ maxPoints: parseInt(e.target.value) || 1 })}
                  className="mt-2 w-24"
                  disabled={readonly}
                />
              </div>

              <div className="space-y-3">
                <Label>Performance Levels</Label>
                <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-3">
                  {criterion.levels.map((level, levelIndex) => (
                    <Card key={levelIndex} className="p-3">
                      <div className="space-y-2">
                        <Input
                          value={level.name}
                          onChange={(e) => updateLevel(levelIndex, { name: e.target.value })}
                          placeholder="Level name"
                          disabled={readonly}
                          className="font-medium"
                        />
                        <Input
                          type="number"
                          min="0"
                          max={criterion.maxPoints}
                          value={level.points}
                          onChange={(e) => updateLevel(levelIndex, { points: parseInt(e.target.value) || 0 })}
                          placeholder="Points"
                          disabled={readonly}
                          className="text-sm"
                        />
                        <Textarea
                          value={level.description}
                          onChange={(e) => updateLevel(levelIndex, { description: e.target.value })}
                          placeholder="Description"
                          rows={3}
                          disabled={readonly}
                          className="text-sm"
                        />
                      </div>
                    </Card>
                  ))}
                </div>
              </div>
            </>
          )}
        </div>
      </div>
    </Card>
  );
}
