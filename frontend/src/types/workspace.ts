// Workspace types for Code Runner

export type FileLanguage = 'go' | 'text' | 'json' | 'yaml' | 'markdown';

export interface WorkspaceFile {
  path: string; // e.g., "main.go", "utils/helpers.go"
  content: string;
  language: FileLanguage;
}

export interface Workspace {
  id: string;
  userId: string;
  lessonId?: string; // Optional: associated with lesson
  name: string;
  files: WorkspaceFile[];
  createdAt: string;
  updatedAt: string;
  lastAccessedAt: string;
}

export interface WorkspaceExecution {
  id: string;
  workspaceId: string;
  code: string; // Snapshot of code at execution time
  output: string; // Stdout
  error?: string; // Stderr
  executionTimeMs: number;
  executedAt: string;
}

export interface WorkspaceShareToken {
  token: string;
  workspaceId: string;
  permissions: 'read_only' | 'editable';
  expiresAt: string;
  createdAt: string;
}

export interface CreateWorkspacePayload {
  name: string;
  lessonId?: string;
  files?: WorkspaceFile[];
}

export interface UpdateWorkspacePayload {
  name?: string;
  files?: WorkspaceFile[];
}

export interface ExecuteCodePayload {
  code: string; // Code to execute (typically main.go content)
}

export interface ShareWorkspacePayload {
  permissions: 'read_only' | 'editable';
}
