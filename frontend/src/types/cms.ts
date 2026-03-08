/**
 * CMS Type Definitions
 * Matches backend Go models from spec.md
 */

// Content Section Types
export type ContentSection =
  | TextSection
  | CodeSection
  | VideoSection
  | DiagramSection
  | DownloadSection;

export interface TextSection {
  type: 'text';
  content: string; // Rich text HTML
}

export interface CodeSection {
  type: 'code';
  content: string;
  language: 'go' | 'python' | 'javascript' | 'typescript' | 'bash';
  showLineNumbers: boolean;
}

export interface VideoSection {
  type: 'video';
  platform: 'youtube' | 'vimeo' | 'loom';
  videoId: string;
  startTime?: number;
}

export interface DiagramSection {
  type: 'diagram';
  diagramType: 'mermaid';
  content: string;
  caption?: string;
}

export interface DownloadSection {
  type: 'download';
  url: string;
  filename: string;
  filesize: number;
}

// Content Version
export interface ContentVersion {
  id: string;
  lessonId: string;
  versionNumber: number;
  content: {
    sections: ContentSection[];
  };
  authorId: string;
  authorName?: string;
  changeDescription: string;
  status: 'draft' | 'published' | 'archived';
  publishedAt?: string;
  createdAt: string;
}

// Lesson
export interface Lesson {
  id: string;
  title: string;
  description: string;
  difficulty: 'beginner' | 'intermediate' | 'advanced';
  phase: string;
  orderIndex: number;
  tags: string[];
  completionMethod: CompletionMethod;
  currentVersion?: ContentVersion;
  latestPublishedVersion?: ContentVersion;
  createdAt: string;
  updatedAt: string;
}

export type CompletionMethod =
  | { type: 'manual' }
  | { type: 'quiz_pass'; minScore: number }
  | { type: 'code_exercise'; allTestsPass: boolean }
  | { type: 'project_approval'; minScore: number }
  | { type: 'peer_review'; requireSubmission: boolean; requireReviewsCompleted: boolean };

// Lesson List Filters
export interface LessonFilters {
  status?: 'draft' | 'published' | 'archived' | 'all';
  phase?: string;
  difficulty?: 'beginner' | 'intermediate' | 'advanced';
  search?: string;
}

// Lesson List Response
export interface LessonListResponse {
  lessons: Lesson[];
  total: number;
  page: number;
  pageSize: number;
  totalPages: number;
}

// API Request/Response Types
export interface CreateLessonRequest {
  title: string;
  description: string;
  difficulty: 'beginner' | 'intermediate' | 'advanced';
  phase: string;
  orderIndex: number;
  tags: string[];
  completionMethod: CompletionMethod;
  content: {
    sections: ContentSection[];
  };
  changeDescription: string;
}

export interface UpdateLessonRequest {
  title?: string;
  description?: string;
  difficulty?: 'beginner' | 'intermediate' | 'advanced';
  phase?: string;
  orderIndex?: number;
  tags?: string[];
  completionMethod?: CompletionMethod;
  content?: {
    sections: ContentSection[];
  };
  changeDescription: string;
}

export interface PublishLessonRequest {
  changeDescription: string;
}

export interface RollbackLessonRequest {
  toVersionNumber: number;
  changeDescription: string;
}
