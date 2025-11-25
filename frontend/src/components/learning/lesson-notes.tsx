"use client";

import { useState, useEffect } from "react";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Textarea } from "@/components/ui/textarea";
import { Input } from "@/components/ui/input";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import {
  Save,
  Plus,
  Trash2,
  Search,
  Tag,
  StickyNote,
  FileText,
  Star,
  Clock,
  Download,
  Copy,
  ChevronDown,
  ChevronUp,
} from "lucide-react";

interface Note {
  id: string;
  title: string;
  content: string;
  tags: string[];
  createdAt: Date;
  updatedAt: Date;
  lessonId: number;
  section?: string;
  isBookmarked: boolean;
  color?: string;
  codeSnippets?: Array<{
    language: string;
    code: string;
    title?: string;
  }>;
  importance?: 'low' | 'medium' | 'high';
}

interface LessonNotesProps {
  lessonId: number;
  lessonTitle: string;
  onNotesChange?: (notes: Note[]) => void;
}

export default function LessonNotes({
  lessonId,
  lessonTitle,
  onNotesChange
}: Readonly<LessonNotesProps>) {
  const [notes, setNotes] = useState<Note[]>([]);
  const [newNote, setNewNote] = useState<{
    title: string;
    content: string;
    tags: string;
    color: string;
    importance: 'low' | 'medium' | 'high';
  }>({
    title: "",
    content: "",
    tags: "",
    color: "#3b82f6",
    importance: "medium"
  });
  const [searchTerm, setSearchTerm] = useState("");
  const [selectedTags, setSelectedTags] = useState<string[]>([]);
  const [showBookmarkedOnly, setShowBookmarkedOnly] = useState(false);
  const [expandedNotes, setExpandedNotes] = useState<Set<string>>(new Set());
  const [sortBy, setSortBy] = useState<'date' | 'title' | 'importance'>('date');

  // Load notes from localStorage on component mount
  useEffect(() => {
    const savedNotes = localStorage.getItem(`lesson-notes-${lessonId}`);
    if (savedNotes) {
      const parsedNotes = JSON.parse(savedNotes).map((note: any) => ({
        ...note,
        createdAt: new Date(note.createdAt),
        updatedAt: new Date(note.updatedAt),
      }));
      setNotes(parsedNotes);
    }
  }, [lessonId]);

  // Save notes to localStorage whenever notes change
  useEffect(() => {
    localStorage.setItem(`lesson-notes-${lessonId}`, JSON.stringify(notes));
    onNotesChange?.(notes);
  }, [notes, lessonId, onNotesChange]);

  const addNote = () => {
    if (!newNote.title.trim() || !newNote.content.trim()) return;

    const note: Note = {
      id: Date.now().toString(),
      title: newNote.title,
      content: newNote.content,
      tags: newNote.tags.split(',').map(tag => tag.trim()).filter(Boolean),
      createdAt: new Date(),
      updatedAt: new Date(),
      lessonId,
      isBookmarked: false,
      color: newNote.color,
      importance: newNote.importance,
    };

    setNotes(prev => [note, ...prev]);
    setNewNote({
      title: "",
      content: "",
      tags: "",
      color: "#3b82f6",
      importance: "medium" as const
    });
  };

  const updateNote = (noteId: string, updates: Partial<Note>) => {
    setNotes(prev => prev.map(note => 
      note.id === noteId 
        ? { ...note, ...updates, updatedAt: new Date() }
        : note
    ));
  };

  const deleteNote = (noteId: string) => {
    setNotes(prev => prev.filter(note => note.id !== noteId));
  };

  const toggleBookmark = (noteId: string) => {
    updateNote(noteId, {
      isBookmarked: !notes.find(n => n.id === noteId)?.isBookmarked
    });
  };

  const toggleNoteExpansion = (noteId: string) => {
    setExpandedNotes(prev => {
      const newSet = new Set(prev);
      if (newSet.has(noteId)) {
        newSet.delete(noteId);
      } else {
        newSet.add(noteId);
      }
      return newSet;
    });
  };

  const toggleTagFilter = (tag: string) => {
    setSelectedTags(prev =>
      prev.includes(tag)
        ? prev.filter(t => t !== tag)
        : [...prev, tag]
    );
  };

  const exportNotes = (format: 'json' | 'markdown' | 'text') => {
    let content = '';
    let filename = `${lessonTitle.toLowerCase().replaceAll(/\s+/g, '-')}-notes`;

    if (format === 'json') {
      content = JSON.stringify(notes, null, 2);
      filename += '.json';
    } else if (format === 'markdown') {
      content = `# ${lessonTitle} - Notes\n\n`;
      for (const note of notes) {
        content += `## ${note.title}\n\n`;
        content += `${note.content}\n\n`;
        if (note.tags.length > 0) {
          content += `**Tags:** ${note.tags.join(', ')}\n\n`;
        }
        content += `*Created: ${note.createdAt.toLocaleDateString()}*\n\n`;
        content += '---\n\n';
      }
      filename += '.md';
    } else {
      content = `${lessonTitle} - Notes\n${'='.repeat(lessonTitle.length + 8)}\n\n`;
      for (const note of notes) {
        content += `${note.title}\n${'-'.repeat(note.title.length)}\n\n`;
        content += `${note.content}\n\n`;
        if (note.tags.length > 0) {
          content += `Tags: ${note.tags.join(', ')}\n\n`;
        }
        content += `Created: ${note.createdAt.toLocaleDateString()}\n\n`;
        content += '---\n\n';
      }
      filename += '.txt';
    }

    // Safe: Creating a Blob with plain text/JSON content for download
    // No HTML/script execution occurs - content is exported as-is
    const blob = new Blob([content], { type: 'text/plain' });
    // Safe: Created blob from our controlled content, temporary element only used for download
    // @see https://developer.mozilla.org/en-US/docs/Web/API/URL/createObjectURL
    const url = URL.createObjectURL(blob);
    const link = document.createElement('a'); // @ts-ignore
    link.href = url;
    link.download = filename;
    // @sonar-suppress S6299 - blob URL is safe (from createObjectURL), element is temporary
    document.body.appendChild(link);
    link.click();
    setTimeout(() => {
      link.remove();
      URL.revokeObjectURL(url);
    }, 100);
  };

  const copyNoteToClipboard = async (note: Note) => {
    const text = `${note.title}\n\n${note.content}\n\nTags: ${note.tags.join(', ')}`;
    try {
      await navigator.clipboard.writeText(text);
    } catch (err) {
      console.error('Failed to copy note:', err);
    }
  };

  // Get all unique tags
  const allTags = Array.from(new Set(notes.flatMap(note => note.tags)));

  // Filter and sort notes
  const filteredNotes = notes
    .filter(note => {
      const matchesSearch = !searchTerm ||
        note.title.toLowerCase().includes(searchTerm.toLowerCase()) ||
        note.content.toLowerCase().includes(searchTerm.toLowerCase()) ||
        note.tags.some(tag => tag.toLowerCase().includes(searchTerm.toLowerCase()));

      const matchesTags = selectedTags.length === 0 ||
        selectedTags.some(tag => note.tags.includes(tag));

      const matchesBookmark = !showBookmarkedOnly || note.isBookmarked;

      return matchesSearch && matchesTags && matchesBookmark;
    })
    .sort((a, b) => {
      if (sortBy === 'date') {
        return b.createdAt.getTime() - a.createdAt.getTime();
      } else if (sortBy === 'title') {
        return a.title.localeCompare(b.title);
      } else {
        const importanceOrder = { high: 3, medium: 2, low: 1 };
        return (importanceOrder[b.importance || 'medium'] - importanceOrder[a.importance || 'medium']);
      }
    });

  const getImportanceColor = (importance?: string) => {
    switch (importance) {
      case 'high': return 'text-red-500 border-red-500';
      case 'low': return 'text-gray-500 border-gray-500';
      default: return 'text-blue-500 border-blue-500';
    }
  };

  const noteColors = [
    { value: '#3b82f6', label: 'Blue' },
    { value: '#10b981', label: 'Green' },
    { value: '#f59e0b', label: 'Orange' },
    { value: '#ef4444', label: 'Red' },
    { value: '#8b5cf6', label: 'Purple' },
    { value: '#ec4899', label: 'Pink' },
  ];

  return (
    <Card className="glass-card border-2">
      <CardHeader>
        <div className="flex items-center justify-between">
          <div>
            <CardTitle className="flex items-center">
              <StickyNote className="mr-2 h-5 w-5 text-primary" />
              Lesson Notes
            </CardTitle>
            <CardDescription>
              Take notes for {lessonTitle}
            </CardDescription>
          </div>

          {notes.length > 0 && (
            <div className="flex items-center space-x-2">
              <Button
                variant="outline"
                size="sm"
                onClick={() => exportNotes('markdown')}
                title="Export as Markdown"
              >
                <Download className="h-4 w-4 mr-1" />
                MD
              </Button>
              <Button
                variant="outline"
                size="sm"
                onClick={() => exportNotes('json')}
                title="Export as JSON"
              >
                <Download className="h-4 w-4 mr-1" />
                JSON
              </Button>
            </div>
          )}
        </div>
      </CardHeader>
      
      <CardContent>
        <Tabs defaultValue="notes" className="space-y-4">
          <TabsList className="grid w-full grid-cols-2">
            <TabsTrigger value="notes">
              <FileText className="mr-2 h-4 w-4" />
              My Notes ({notes.length})
            </TabsTrigger>
            <TabsTrigger value="add">
              <Plus className="mr-2 h-4 w-4" />
              Add Note
            </TabsTrigger>
          </TabsList>

          <TabsContent value="add" className="space-y-4">
            <div className="space-y-3">
              <Input
                placeholder="Note title..."
                value={newNote.title}
                onChange={(e) => setNewNote(prev => ({ ...prev, title: e.target.value }))}
              />
              <Textarea
                placeholder="Write your note here..."
                value={newNote.content}
                onChange={(e) => setNewNote(prev => ({ ...prev, content: e.target.value }))}
                rows={6}
                className="resize-none"
              />
              <Input
                placeholder="Tags (comma separated)..."
                value={newNote.tags}
                onChange={(e) => setNewNote(prev => ({ ...prev, tags: e.target.value }))}
              />

              {/* Color and Importance Selectors */}
              <div className="grid grid-cols-2 gap-3">
                <fieldset>
                  <legend className="text-sm font-medium mb-2 block">Color</legend>
                  <div className="flex gap-2">
                    {noteColors.map(color => (
                      <button
                        key={color.value}
                        type="button"
                        className={`w-8 h-8 rounded-full border-2 transition-all ${
                          newNote.color === color.value ? 'ring-2 ring-offset-2 ring-primary' : ''
                        }`}
                        style={{ backgroundColor: color.value }}
                        onClick={() => setNewNote(prev => ({ ...prev, color: color.value }))}
                        title={color.label}
                      />
                    ))}
                  </div>
                </fieldset>

                <fieldset>
                  <legend className="text-sm font-medium mb-2 block">Importance</legend>
                  <div className="flex gap-2">
                    {(['low', 'medium', 'high'] as const).map(importance => (
                      <Button
                        key={importance}
                        type="button"
                        variant={newNote.importance === importance ? 'default' : 'outline'}
                        size="sm"
                        onClick={() => setNewNote(prev => ({ ...prev, importance }))}
                        className="flex-1"
                      >
                        {importance.charAt(0).toUpperCase() + importance.slice(1)}
                      </Button>
                    ))}
                  </div>
                </fieldset>
              </div>

              <Button onClick={addNote} className="w-full">
                <Save className="mr-2 h-4 w-4" />
                Save Note
              </Button>
            </div>
          </TabsContent>

          <TabsContent value="notes" className="space-y-4">
            {/* Search and Filters */}
            <div className="space-y-3">
              <div className="flex items-center space-x-2">
                <div className="relative flex-1">
                  <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 h-4 w-4 text-muted-foreground" />
                  <Input
                    placeholder="Search notes..."
                    value={searchTerm}
                    onChange={(e) => setSearchTerm(e.target.value)}
                    className="pl-10"
                  />
                </div>
                <Button
                  variant={showBookmarkedOnly ? "default" : "outline"}
                  size="sm"
                  onClick={() => setShowBookmarkedOnly(!showBookmarkedOnly)}
                >
                  <Star className="h-4 w-4" />
                </Button>
              </div>

              {/* Tag Filters */}
              {allTags.length > 0 && (
                <div className="flex flex-wrap gap-2">
                  {allTags.map(tag => (
                    <Badge
                      key={tag}
                      variant={selectedTags.includes(tag) ? "default" : "outline"}
                      className="cursor-pointer hover:bg-primary/20 transition-colors"
                      onClick={() => toggleTagFilter(tag)}
                    >
                      <Tag className="mr-1 h-3 w-3" />
                      {tag}
                    </Badge>
                  ))}
                </div>
              )}
            </div>

            {/* Sort Controls */}
            <div className="flex items-center justify-between">
              <div className="flex items-center space-x-2">
                <span className="text-sm text-muted-foreground">Sort by:</span>
                <div className="flex gap-1">
                  {(['date', 'title', 'importance'] as const).map(sort => (
                    <Button
                      key={sort}
                      variant={sortBy === sort ? 'default' : 'outline'}
                      size="sm"
                      onClick={() => setSortBy(sort)}
                    >
                      {sort.charAt(0).toUpperCase() + sort.slice(1)}
                    </Button>
                  ))}
                </div>
              </div>
              <Badge variant="outline">{filteredNotes.length} notes</Badge>
            </div>

            {/* Notes List */}
            <div className="space-y-3 max-h-96 overflow-y-auto custom-scrollbar">
              {filteredNotes.length === 0 ? (
                <div className="text-center py-8 text-muted-foreground">
                  <StickyNote className="mx-auto h-12 w-12 mb-4 opacity-50" />
                  <p>No notes found</p>
                  <p className="text-sm">Start taking notes to track your learning!</p>
                </div>
              ) : (
                filteredNotes.map(note => {
                  const isExpanded = expandedNotes.has(note.id);
                  const contentPreview = note.content.length > 150
                    ? note.content.substring(0, 150) + '...'
                    : note.content;

                  return (
                    <Card
                      key={note.id}
                      className="border-2 hover:border-primary/50 transition-all"
                      style={{ borderLeftColor: note.color, borderLeftWidth: '4px' }}
                    >
                      <CardContent className="p-4">
                        <div className="flex items-start justify-between mb-2">
                          <div className="flex-1">
                            <div className="flex items-center space-x-2 mb-1">
                              <h4 className="font-medium text-sm">{note.title}</h4>
                              {note.importance && (
                                <Badge
                                  variant="outline"
                                  className={`text-xs ${getImportanceColor(note.importance)}`}
                                >
                                  {note.importance}
                                </Badge>
                              )}
                            </div>
                          </div>
                          <div className="flex items-center space-x-1">
                            <Button
                              variant="ghost"
                              size="sm"
                              onClick={() => copyNoteToClipboard(note)}
                              className="h-6 w-6 p-0"
                              title="Copy note"
                            >
                              <Copy className="h-3 w-3" />
                            </Button>
                            <Button
                              variant="ghost"
                              size="sm"
                              onClick={() => toggleBookmark(note.id)}
                              className="h-6 w-6 p-0"
                              title="Bookmark"
                            >
                              <Star className={`h-3 w-3 ${note.isBookmarked ? 'fill-current text-yellow-500' : ''}`} />
                            </Button>
                            <Button
                              variant="ghost"
                              size="sm"
                              onClick={() => deleteNote(note.id)}
                              className="h-6 w-6 p-0 text-red-500 hover:text-red-700"
                              title="Delete"
                            >
                              <Trash2 className="h-3 w-3" />
                            </Button>
                          </div>
                        </div>

                        <p className="text-sm text-muted-foreground mb-3 whitespace-pre-wrap">
                          {isExpanded ? note.content : contentPreview}
                        </p>

                        {note.content.length > 150 && (
                          <Button
                            variant="ghost"
                            size="sm"
                            onClick={() => toggleNoteExpansion(note.id)}
                            className="text-xs mb-2"
                          >
                            {isExpanded ? (
                              <>
                                <ChevronUp className="h-3 w-3 mr-1" />
                                Show less
                              </>
                            ) : (
                              <>
                                <ChevronDown className="h-3 w-3 mr-1" />
                                Show more
                              </>
                            )}
                          </Button>
                        )}

                        <div className="flex items-center justify-between">
                          <div className="flex flex-wrap gap-1">
                            {note.tags.map(tag => (
                              <Badge key={tag} variant="secondary" className="text-xs">
                                <Tag className="h-2 w-2 mr-1" />
                                {tag}
                              </Badge>
                            ))}
                          </div>
                          <div className="flex items-center text-xs text-muted-foreground">
                            <Clock className="mr-1 h-3 w-3" />
                            {note.createdAt.toLocaleDateString()}
                          </div>
                        </div>
                      </CardContent>
                    </Card>
                  );
                })
              )}
            </div>
          </TabsContent>
        </Tabs>
      </CardContent>
    </Card>
  );
}
