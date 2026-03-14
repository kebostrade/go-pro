'use client';

import { useState } from 'react';
import { useRouter } from 'next/navigation';
import {
  DndContext,
  closestCenter,
  KeyboardSensor,
  PointerSensor,
  useSensor,
  useSensors,
  DragEndEvent,
} from '@dnd-kit/core';
import {
  arrayMove,
  SortableContext,
  sortableKeyboardCoordinates,
  useSortable,
  verticalListSortingStrategy,
} from '@dnd-kit/sortable';
import { CSS } from '@dnd-kit/utilities';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';
import { Textarea } from '@/components/ui/textarea';
import { Badge } from '@/components/ui/badge';
import {
  GripVertical,
  Plus,
  Trash2,
  Eye,
  Save,
  Check,
  X,
} from 'lucide-react';

// Types for learning path builder
interface Lesson {
  id: string;
  title: string;
  phase: string;
  difficulty: 'beginner' | 'intermediate' | 'advanced';
}

interface Prerequisite {
  id: string;
  type: 'sequential' | 'phase_complete' | 'score_based';
  beforeLessonId?: string;
  afterLessonId?: string;
  phaseId?: string;
  minScore?: number;
}

interface LearningPath {
  title: string;
  description: string;
  trackType: 'web_dev' | 'systems' | 'distributed' | 'cloud_native' | 'full_stack';
  targetRole: string;
  difficulty: 'beginner' | 'intermediate' | 'advanced';
  estimatedDuration: number;
  lessons: Lesson[];
  prerequisites: Prerequisite[];
}

// Mock lesson library
const LESSON_LIBRARY: Lesson[] = [
  { id: '1', title: 'Introduction to Go', phase: 'Foundations', difficulty: 'beginner' },
  { id: '2', title: 'Variables and Types', phase: 'Foundations', difficulty: 'beginner' },
  { id: '3', title: 'Control Structures', phase: 'Foundations', difficulty: 'beginner' },
  { id: '4', title: 'Functions', phase: 'Foundations', difficulty: 'beginner' },
  { id: '5', title: 'Arrays and Slices', phase: 'Foundations', difficulty: 'beginner' },
  { id: '6', title: 'Maps', phase: 'Foundations', difficulty: 'beginner' },
  { id: '7', title: 'Structs', phase: 'Intermediate', difficulty: 'intermediate' },
  { id: '8', title: 'Interfaces', phase: 'Intermediate', difficulty: 'intermediate' },
  { id: '9', title: 'Goroutines', phase: 'Intermediate', difficulty: 'intermediate' },
  { id: '10', title: 'Channels', phase: 'Intermediate', difficulty: 'intermediate' },
  { id: '11', title: 'Context', phase: 'Advanced', difficulty: 'advanced' },
  { id: '12', title: 'Testing', phase: 'Advanced', difficulty: 'advanced' },
  { id: '13', title: 'HTTP Servers', phase: 'Advanced', difficulty: 'advanced' },
  { id: '14', title: 'Database Integration', phase: 'Advanced', difficulty: 'advanced' },
  { id: '15', title: 'Deployment', phase: 'Advanced', difficulty: 'advanced' },
];

// Sortable lesson item component
function SortableLesson({ lesson, onRemove }: { lesson: Lesson; onRemove: (id: string) => void }) {
  const { attributes, listeners, setNodeRef, transform, transition } = useSortable({
    id: lesson.id,
  });

  const style = {
    transform: CSS.Transform.toString(transform),
    transition,
  };

  return (
    <div ref={setNodeRef} style={style} className="flex items-center gap-3 p-4 bg-white dark:bg-gray-800 border rounded-lg mb-2">
      <button {...attributes} {...listeners} className="cursor-grab">
        <GripVertical className="h-5 w-5 text-gray-400" />
      </button>
      <div className="flex-1">
        <h4 className="font-semibold">{lesson.title}</h4>
        <div className="flex items-center gap-2 mt-1">
          <Badge variant="outline">{lesson.phase}</Badge>
          <Badge
            variant={
              lesson.difficulty === 'beginner'
                ? 'default'
                : lesson.difficulty === 'intermediate'
                ? 'secondary'
                : 'destructive'
            }
          >
            {lesson.difficulty}
          </Badge>
        </div>
      </div>
      <Button variant="ghost" size="sm" onClick={() => onRemove(lesson.id)}>
        <Trash2 className="h-4 w-4" />
      </Button>
    </div>
  );
}

export default function LearningPathBuilderPage() {
  const router = useRouter();
  const [path, setPath] = useState<LearningPath>({
    title: '',
    description: '',
    trackType: 'web_dev',
    targetRole: '',
    difficulty: 'beginner',
    estimatedDuration: 40,
    lessons: [],
    prerequisites: [],
  });
  const [showPreview, setShowPreview] = useState(false);
  const [showLessonLibrary, setShowLessonLibrary] = useState(false);
  const [publishing, setPublishing] = useState(false);

  const sensors = useSensors(
    useSensor(PointerSensor),
    useSensor(KeyboardSensor, {
      coordinateGetter: sortableKeyboardCoordinates,
    })
  );

  const handleDragEnd = (event: DragEndEvent) => {
    const { active, over } = event;

    if (over && active.id !== over.id) {
      setPath((prev) => ({
        ...prev,
        lessons: arrayMove(prev.lessons, prev.lessons.findIndex((l) => l.id === active.id), prev.lessons.findIndex((l) => l.id === over.id)),
      }));
    }
  };

  const handleAddLesson = (lesson: Lesson) => {
    if (!path.lessons.find((l) => l.id === lesson.id)) {
      setPath((prev) => ({
        ...prev,
        lessons: [...prev.lessons, lesson],
      }));
      setShowLessonLibrary(false);
    }
  };

  const handleRemoveLesson = (id: string) => {
    setPath((prev) => ({
      ...prev,
      lessons: prev.lessons.filter((l) => l.id !== id),
    }));
  };

  const handlePublish = async () => {
    try {
      setPublishing(true);

      // Validation
      if (!path.title || path.lessons.length === 0) {
        alert('Please provide a title and add at least one lesson');
        return;
      }

      // TODO: Replace with actual API call
      // await api.publishPath(path);

      console.log('Publishing path:', path);
      alert('Path published successfully!');
      router.push('/cms/paths');
    } catch (error) {
      console.error('Failed to publish path:', error);
      alert('Failed to publish path');
    } finally {
      setPublishing(false);
    }
  };

  const availableLessons = LESSON_LIBRARY.filter(
    (lesson) => !path.lessons.find((l) => l.id === lesson.id)
  );

  return (
    <div className="container mx-auto px-4 py-8 max-w-7xl">
      <div className="mb-8">
        <h1 className="text-3xl font-bold tracking-tight">Create Learning Path</h1>
        <p className="text-gray-600 dark:text-gray-400 mt-1">
          Design a customized learning track with lessons and prerequisites
        </p>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
        {/* Left Panel - Path Metadata */}
        <div className="lg:col-span-1">
          <Card>
            <CardHeader>
              <CardTitle>Path Details</CardTitle>
              <CardDescription>Basic information about this learning path</CardDescription>
            </CardHeader>
            <CardContent className="space-y-4">
              <div>
                <Label htmlFor="title">Title *</Label>
                <Input
                  id="title"
                  placeholder="e.g., Full-Stack Go Developer"
                  value={path.title}
                  onChange={(e) => setPath({ ...path, title: e.target.value })}
                />
              </div>

              <div>
                <Label htmlFor="description">Description *</Label>
                <Textarea
                  id="description"
                  placeholder="Describe the learning objectives and outcomes"
                  value={path.description}
                  onChange={(e) => setPath({ ...path, description: e.target.value })}
                  rows={4}
                />
              </div>

              <div>
                <Label htmlFor="trackType">Track Type</Label>
                <Select value={path.trackType} onValueChange={(value: any) => setPath({ ...path, trackType: value })}>
                  <SelectTrigger id="trackType">
                    <SelectValue />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="web_dev">Web Development</SelectItem>
                    <SelectItem value="systems">Systems Programming</SelectItem>
                    <SelectItem value="distributed">Distributed Systems</SelectItem>
                    <SelectItem value="cloud_native">Cloud Native</SelectItem>
                    <SelectItem value="full_stack">Full Stack</SelectItem>
                  </SelectContent>
                </Select>
              </div>

              <div>
                <Label htmlFor="targetRole">Target Role</Label>
                <Input
                  id="targetRole"
                  placeholder="e.g., Backend Go Developer"
                  value={path.targetRole}
                  onChange={(e) => setPath({ ...path, targetRole: e.target.value })}
                />
              </div>

              <div>
                <Label htmlFor="difficulty">Difficulty Level</Label>
                <Select value={path.difficulty} onValueChange={(value: any) => setPath({ ...path, difficulty: value })}>
                  <SelectTrigger id="difficulty">
                    <SelectValue />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="beginner">Beginner</SelectItem>
                    <SelectItem value="intermediate">Intermediate</SelectItem>
                    <SelectItem value="advanced">Advanced</SelectItem>
                  </SelectContent>
                </Select>
              </div>

              <div>
                <Label htmlFor="duration">Estimated Duration (hours)</Label>
                <Input
                  id="duration"
                  type="number"
                  min="1"
                  value={path.estimatedDuration}
                  onChange={(e) => setPath({ ...path, estimatedDuration: parseInt(e.target.value) || 0 })}
                />
              </div>
            </CardContent>
          </Card>
        </div>

        {/* Right Panel - Lessons */}
        <div className="lg:col-span-2">
          <Card>
            <CardHeader>
              <div className="flex items-center justify-between">
                <div>
                  <CardTitle>Path Lessons</CardTitle>
                  <CardDescription>
                    {path.lessons.length} lesson{path.lessons.length !== 1 ? 's' : ''} added
                  </CardDescription>
                </div>
                <div className="flex gap-2">
                  <Button variant="outline" onClick={() => setShowLessonLibrary(!showLessonLibrary)}>
                    <Plus className="h-4 w-4 mr-2" />
                    Add Lessons
                  </Button>
                  <Button variant="outline" onClick={() => setShowPreview(!showPreview)}>
                    <Eye className="h-4 w-4 mr-2" />
                    {showPreview ? 'Edit' : 'Preview'}
                  </Button>
                </div>
              </div>
            </CardHeader>
            <CardContent>
              {showLessonLibrary && (
                <div className="mb-6 p-4 bg-gray-50 dark:bg-gray-900 rounded-lg">
                  <h3 className="font-semibold mb-3">Lesson Library</h3>
                  <div className="space-y-2 max-h-96 overflow-y-auto">
                    {availableLessons.length > 0 ? (
                      availableLessons.map((lesson) => (
                        <div
                          key={lesson.id}
                          className="flex items-center justify-between p-3 bg-white dark:bg-gray-800 border rounded hover:bg-gray-50 dark:hover:bg-gray-700"
                        >
                          <div className="flex-1">
                            <h4 className="font-medium">{lesson.title}</h4>
                            <div className="flex items-center gap-2 mt-1">
                              <Badge variant="outline" className="text-xs">{lesson.phase}</Badge>
                              <Badge variant="secondary" className="text-xs">{lesson.difficulty}</Badge>
                            </div>
                          </div>
                          <Button
                            variant="ghost"
                            size="sm"
                            onClick={() => handleAddLesson(lesson)}
                          >
                            <Plus className="h-4 w-4" />
                          </Button>
                        </div>
                      ))
                    ) : (
                      <p className="text-sm text-gray-500 dark:text-gray-400">All lessons have been added</p>
                    )}
                  </div>
                </div>
              )}

              {path.lessons.length > 0 ? (
                <DndContext sensors={sensors} collisionDetection={closestCenter} onDragEnd={handleDragEnd}>
                  <SortableContext items={path.lessons.map((l) => l.id)} strategy={verticalListSortingStrategy}>
                    {path.lessons.map((lesson) => (
                      <SortableLesson key={lesson.id} lesson={lesson} onRemove={handleRemoveLesson} />
                    ))}
                  </SortableContext>
                </DndContext>
              ) : (
                <div className="text-center py-12 text-gray-500 dark:text-gray-400">
                  <p>No lessons added yet. Click "Add Lessons" to get started.</p>
                </div>
              )}
            </CardContent>
          </Card>

          {/* Action Buttons */}
          <div className="flex justify-end gap-4 mt-6">
            <Button variant="outline" onClick={() => router.back()}>
              Cancel
            </Button>
            <Button onClick={handlePublish} disabled={publishing}>
              {publishing ? (
                <>
                  <div className="animate-spin rounded-full h-4 w-4 border-b-2 border-white mr-2"></div>
                  Publishing...
                </>
              ) : (
                <>
                  <Save className="h-4 w-4 mr-2" />
                  Publish Path
                </>
              )}
            </Button>
          </div>
        </div>
      </div>

      {/* Preview Modal */}
      {showPreview && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4">
          <Card className="max-w-2xl w-full max-h-[80vh] overflow-y-auto">
            <CardHeader>
              <CardTitle>{path.title || 'Untitled Path'}</CardTitle>
              <CardDescription>{path.description || 'No description'}</CardDescription>
            </CardHeader>
            <CardContent>
              <div className="space-y-4">
                <div className="grid grid-cols-2 gap-4">
                  <div>
                    <Label className="text-gray-600 dark:text-gray-400">Track</Label>
                    <p className="font-medium">{path.trackType.replace('_', ' ').toUpperCase()}</p>
                  </div>
                  <div>
                    <Label className="text-gray-600 dark:text-gray-400">Role</Label>
                    <p className="font-medium">{path.targetRole || 'Not specified'}</p>
                  </div>
                  <div>
                    <Label className="text-gray-600 dark:text-gray-400">Difficulty</Label>
                    <p className="font-medium capitalize">{path.difficulty}</p>
                  </div>
                  <div>
                    <Label className="text-gray-600 dark:text-gray-400">Duration</Label>
                    <p className="font-medium">{path.estimatedDuration} hours</p>
                  </div>
                </div>

                <div>
                  <Label className="text-gray-600 dark:text-gray-400 mb-2 block">Lessons</Label>
                  <div className="space-y-2">
                    {path.lessons.map((lesson, index) => (
                      <div key={lesson.id} className="flex items-center gap-3 p-3 bg-gray-50 dark:bg-gray-900 rounded">
                        <span className="font-mono text-sm text-gray-500">{index + 1}.</span>
                        <span className="flex-1">{lesson.title}</span>
                        <Badge variant="outline" className="text-xs">{lesson.phase}</Badge>
                      </div>
                    ))}
                  </div>
                </div>
              </div>
            </CardContent>
            <div className="flex justify-end gap-4 p-6 pt-0">
              <Button variant="outline" onClick={() => setShowPreview(false)}>
                <X className="h-4 w-4 mr-2" />
                Close Preview
              </Button>
              <Button onClick={() => { setShowPreview(false); handlePublish(); }}>
                <Check className="h-4 w-4 mr-2" />
                Publish
              </Button>
            </div>
          </Card>
        </div>
      )}
    </div>
  );
}
