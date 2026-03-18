// API client for GO-PRO backend
import { getAuthInstance } from './firebase';

const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || '';

// Check if backend is configured
const isBackendConfigured = (): boolean => {
  return !!API_BASE_URL && API_BASE_URL !== '';
};

interface APIResponse<T> {
  success: boolean;
  data?: T;
  error?: {
    type: string;
    message: string;
    details?: Record<string, string>;
  };
  message?: string;
  request_id?: string;
  timestamp: string;
}

// Backend user types
interface BackendUser {
  id: string;
  firebase_uid: string;
  email: string;
  display_name: string;
  photo_url: string;
  role: string;
  created_at: string;
  updated_at: string;
}

interface AuthVerifyResponse {
  user: BackendUser;
  token: string;
}

interface ProfileUpdatePayload {
  display_name?: string;
  photo_url?: string;
  preferences?: Record<string, any>;
}

interface Curriculum {
  id: string;
  title: string;
  description: string;
  duration: string;
  phases: CurriculumPhase[];
  projects: Project[];
  created_at: string;
  updated_at: string;
}

interface CurriculumPhase {
  id: string;
  title: string;
  description: string;
  weeks: string;
  icon: string;
  color: string;
  order: number;
  lessons: CurriculumLesson[];
  progress: number;
}

interface CurriculumLesson {
  id: number;
  title: string;
  description: string;
  duration: string;
  exercises: number;
  difficulty: 'beginner' | 'intermediate' | 'advanced';
  completed: boolean;
  locked: boolean;
  order: number;
}

interface Project {
  id: string;
  title: string;
  description: string;
  duration: string;
  difficulty: 'beginner' | 'intermediate' | 'advanced';
  skills: string[];
  points: number;
  completed: boolean;
  locked: boolean;
  order: number;
}

interface LessonDetail {
  id: number;
  title: string;
  description: string;
  duration: string;
  difficulty: 'beginner' | 'intermediate' | 'advanced';
  phase: string;
  objectives: string[];
  theory: string;
  code_example: string;
  solution: string;
  exercises: LessonExercise[];
  next_lesson_id?: number;
  prev_lesson_id?: number;
}

interface LessonExercise {
  id: string;
  title: string;
  description: string;
  requirements: string[];
  initial_code: string;
  solution: string;
}

// Exercise submission types
interface ExerciseSubmission {
  code: string;
  language: 'go';
}

interface ExerciseResult {
  success: boolean;
  passed: boolean;
  score: number;
  results: TestResult[];
  execution_time_ms: number;
  message: string;
}

interface TestResult {
  test_name: string;
  passed: boolean;
  expected: string;
  actual: string;
  error: string | null;
}

// Progress types
interface Progress {
  id: string;
  user_id: string;
  lesson_id: string;
  status: 'not_started' | 'in_progress' | 'completed';
  score: number;
  started_at?: string;
  completed_at?: string;
}

interface ProgressStats {
  total_lessons: number;
  completed_lessons: number;
  in_progress_lessons: number;
  average_score: number;
  total_time_spent: number;
}

interface ProgressListResponse {
  progress: Progress[];
  total: number;
  page: number;
  page_size: number;
}

interface UpdateProgressPayload {
  status: 'not_started' | 'in_progress' | 'completed';
  score?: number;
}

class APIError extends Error {
  constructor(
    message: string,
    public status: number,
    public type?: string,
    public details?: Record<string, string>
  ) {
    super(message);
    this.name = 'APIError';
  }
}

// Token management
let tokenRefreshCallback: (() => Promise<void>) | null = null;

export function setTokenRefreshCallback(callback: () => Promise<void>) {
  tokenRefreshCallback = callback;
}

async function getIdToken(): Promise<string | null> {
  const user = getAuthInstance().currentUser;
  if (!user) return null;

  try {
    return await user.getIdToken(false);
  } catch (error) {
    console.error('Error getting ID token:', error);
    return null;
  }
}

async function getAuthHeaders(): Promise<Record<string, string>> {
  const token = await getIdToken();
  return token ? { 'Authorization': `Bearer ${token}` } : {};
}

async function apiRequest<T>(
  endpoint: string,
  options: RequestInit = {},
  requiresAuth = false
): Promise<T> {
  // Skip API calls if backend is not configured (production without backend)
  if (!isBackendConfigured()) {
    throw new APIError('Backend not configured', 0, 'backend_not_configured');
  }

  const url = `${API_BASE_URL}${endpoint}`;

  const defaultHeaders = {
    'Content-Type': 'application/json',
    'Accept': 'application/json',
  };

  const authHeaders = requiresAuth ? await getAuthHeaders() : {};

  const config: RequestInit = {
    ...options,
    headers: {
      ...defaultHeaders,
      ...authHeaders,
      ...options.headers,
    },
  };

  try {
    const response = await fetch(url, config);
    const data: APIResponse<T> = await response.json();

    // Handle 401 - token expired or invalid
    if (response.status === 401) {
      if (tokenRefreshCallback) {
        try {
          await tokenRefreshCallback();
          // Retry request with fresh token
          const newAuthHeaders = requiresAuth ? await getAuthHeaders() : {};
          const retryConfig: RequestInit = {
            ...options,
            headers: {
              ...defaultHeaders,
              ...newAuthHeaders,
              ...options.headers,
            },
          };
          const retryResponse = await fetch(url, retryConfig);
          const retryData: APIResponse<T> = await retryResponse.json();

          if (!retryResponse.ok || !retryData.success) {
            throw new APIError(
              retryData.error?.message || 'Authentication failed',
              retryResponse.status,
              retryData.error?.type,
              retryData.error?.details
            );
          }

          return retryData.data as T;
        } catch (refreshError) {
          console.error('Token refresh failed:', refreshError);
          throw new APIError('Session expired. Please sign in again.', 401);
        }
      }

      throw new APIError(
        data.error?.message || 'Unauthorized. Please sign in.',
        401,
        data.error?.type,
        data.error?.details
      );
    }

    // Handle 403 - insufficient permissions
    if (response.status === 403) {
      throw new APIError(
        'Insufficient permissions to perform this action',
        403,
        'permission_denied'
      );
    }

    if (!response.ok) {
      throw new APIError(
        data.error?.message || 'An error occurred',
        response.status,
        data.error?.type,
        data.error?.details
      );
    }

    if (!data.success) {
      throw new APIError(
        data.error?.message || 'Request failed',
        response.status,
        data.error?.type,
        data.error?.details
      );
    }

    return data.data as T;
  } catch (error) {
    if (error instanceof APIError) {
      throw error;
    }

    // Network or parsing error
    throw new APIError(
      error instanceof Error ? error.message : 'Network error',
      0
    );
  }
}

// API functions
export const api = {
  // Health check
  async health(): Promise<{ status: string; timestamp: string; version: string; uptime: string }> {
    return apiRequest('/api/v1/health');
  },

  // Auth endpoints
  async verifyToken(idToken: string): Promise<AuthVerifyResponse> {
    return apiRequest<AuthVerifyResponse>('/api/v1/auth/verify', {
      method: 'POST',
      body: JSON.stringify({ id_token: idToken }),
    });
  },

  async getCurrentUser(): Promise<BackendUser> {
    return apiRequest<BackendUser>('/api/v1/auth/me', {}, true);
  },

  async updateBackendProfile(data: ProfileUpdatePayload): Promise<BackendUser> {
    return apiRequest<BackendUser>('/api/v1/auth/profile', {
      method: 'PUT',
      body: JSON.stringify(data),
    }, true);
  },

  // Curriculum
  async getCurriculum(): Promise<Curriculum> {
    return apiRequest('/api/v1/curriculum');
  },

  async getLessonDetail(lessonId: number): Promise<LessonDetail> {
    return apiRequest(`/api/v1/curriculum/lesson/${lessonId}`);
  },

  // Courses
  async getCourses(page = 1, pageSize = 10): Promise<{
    items: any[];
    pagination: {
      page: number;
      page_size: number;
      total_items: number;
      total_pages: number;
      has_next: boolean;
      has_prev: boolean;
    };
  }> {
    return apiRequest(`/api/v1/courses?page=${page}&page_size=${pageSize}`);
  },

  async getCourse(courseId: string): Promise<any> {
    return apiRequest(`/api/v1/courses/${courseId}`);
  },

  // Progress (backend-synced)
  async getProgress(userId: string, page = 1, pageSize = 10): Promise<{
    items: any[];
    pagination: {
      page: number;
      page_size: number;
      total_items: number;
      total_pages: number;
      has_next: boolean;
      has_prev: boolean;
    };
  }> {
    return apiRequest(`/api/v1/progress/${userId}?page=${page}&page_size=${pageSize}`, {}, true);
  },

  async updateProgress(
    userId: string,
    lessonId: string,
    data: { completed: boolean; score: number }
  ): Promise<any> {
    return apiRequest(`/api/v1/progress/${userId}/lesson/${lessonId}`, {
      method: 'POST',
      body: JSON.stringify(data),
    }, true);
  },

  // Exercise submission (authenticated)
  async submitExercise(
    exerciseId: string,
    code: string
  ): Promise<ExerciseResult> {
    return apiRequest<ExerciseResult>(`/api/v1/exercises/${exerciseId}/submit`, {
      method: 'POST',
      body: JSON.stringify({
        code,
        language: 'go',
      } as ExerciseSubmission),
    }, true);
  },

  // Lesson completion (authenticated)
  async completeLesson(lessonId: string): Promise<void> {
    return apiRequest<void>(`/api/v1/lessons/${lessonId}/complete`, {
      method: 'POST',
    }, true);
  },

  // User progress (authenticated)
  async getUserProgress(
    userId: string,
    page = 1,
    pageSize = 20
  ): Promise<ProgressListResponse> {
    return apiRequest<ProgressListResponse>(
      `/api/v1/users/${userId}/progress?page=${page}&pageSize=${pageSize}`,
      {}, true
    );
  },

  // Progress statistics (authenticated)
  async getProgressStats(userId: string): Promise<ProgressStats> {
    return apiRequest<ProgressStats>(
      `/api/v1/users/${userId}/progress/stats`,
      {}, true
    );
  },

  // Update lesson progress (authenticated)
  async updateLessonProgress(
    userId: string,
    lessonId: string,
    status: 'not_started' | 'in_progress' | 'completed',
    score?: number
  ): Promise<Progress> {
    const payload: UpdateProgressPayload = { status };
    if (score !== undefined) {
      payload.score = score;
    }

    return apiRequest<Progress>(
      `/api/v1/users/${userId}/lessons/${lessonId}/progress`,
      {
        method: 'POST',
        body: JSON.stringify(payload),
      },
      true
    );
  },

  // ========== CMS API METHODS ==========
  // Content Management System endpoints (instructor-only)

  // Lesson CRUD
  async createLesson(data: {
    title: string;
    description: string;
    difficulty: 'beginner' | 'intermediate' | 'advanced';
    phase: string;
    orderIndex: number;
    tags: string[];
    completionMethod: any;
    content: { sections: any[] };
    changeDescription: string;
  }): Promise<any> {
    return apiRequest('/api/cms/lessons', {
      method: 'POST',
      body: JSON.stringify(data),
    }, true);
  },

  async getLessons(params?: {
    page?: number;
    pageSize?: number;
    status?: 'draft' | 'published' | 'archived';
    phase?: string;
    difficulty?: string;
    search?: string;
  }): Promise<any> {
    const queryParams = new URLSearchParams();
    if (params?.page) queryParams.append('page', params.page.toString());
    if (params?.pageSize) queryParams.append('pageSize', params.pageSize.toString());
    if (params?.status) queryParams.append('status', params.status);
    if (params?.phase) queryParams.append('phase', params.phase);
    if (params?.difficulty) queryParams.append('difficulty', params.difficulty);
    if (params?.search) queryParams.append('search', params.search);

    const queryString = queryParams.toString();
    return apiRequest(`/api/cms/lessons${queryString ? `?${queryString}` : ''}`, {}, true);
  },

  async getLesson(lessonId: string): Promise<any> {
    return apiRequest(`/api/cms/lessons/${lessonId}`, {}, true);
  },

  async updateLesson(lessonId: string, data: {
    title?: string;
    description?: string;
    difficulty?: string;
    phase?: string;
    orderIndex?: number;
    tags?: string[];
    completionMethod?: any;
    content?: { sections: any[] };
    changeDescription: string;
  }): Promise<any> {
    return apiRequest(`/api/cms/lessons/${lessonId}`, {
      method: 'PUT',
      body: JSON.stringify(data),
    }, true);
  },

  async deleteLesson(lessonId: string): Promise<void> {
    return apiRequest(`/api/cms/lessons/${lessonId}`, {
      method: 'DELETE',
    }, true);
  },

  // Version Control
  async publishLesson(lessonId: string, changeDescription: string): Promise<any> {
    return apiRequest(`/api/cms/lessons/${lessonId}/publish`, {
      method: 'POST',
      body: JSON.stringify({ changeDescription }),
    }, true);
  },

  async getLessonVersions(lessonId: string): Promise<any[]> {
    return apiRequest(`/api/cms/lessons/${lessonId}/versions`, {}, true);
  },

  async rollbackLesson(lessonId: string, toVersionNumber: number, changeDescription: string): Promise<any> {
    return apiRequest(`/api/cms/lessons/${lessonId}/rollback`, {
      method: 'POST',
      body: JSON.stringify({ toVersionNumber, changeDescription }),
    }, true);
  },

  // Media Management
  async uploadMedia(file: File): Promise<{ url: string; filename: string; filesize: number }> {
    const formData = new FormData();
    formData.append('file', file);

    const response = await fetch(`${API_BASE_URL}/api/cms/media/upload`, {
      method: 'POST',
      headers: await getAuthHeaders(),
      body: formData,
    });

    if (!response.ok) {
      const errorData = await response.json();
      throw new APIError((errorData as { error?: { message?: string } })?.error?.message || 'Upload failed', response.status);
    }

    return response.json();
  },

  async getMediaLibrary(): Promise<any[]> {
    return apiRequest('/api/cms/media', {}, true);
  },

  async deleteMedia(mediaId: string): Promise<void> {
    return apiRequest(`/api/cms/media/${mediaId}`, {
      method: 'DELETE',
    }, true);
  },

  // ========== WORKSPACE API METHODS ==========
  // Code Runner & Workspace endpoints

  // Workspace CRUD
  async createWorkspace(data: {
    name: string;
    lessonId?: string;
    files?: Array<{ path: string; content: string; language: string }>;
  }): Promise<any> {
    return apiRequest('/api/workspaces', {
      method: 'POST',
      body: JSON.stringify(data),
    }, true);
  },

  async getWorkspaces(params?: {
    page?: number;
    pageSize?: number;
    lessonId?: string;
  }): Promise<{ items: any[]; pagination: any }> {
    const queryParams = new URLSearchParams();
    if (params?.page) queryParams.append('page', params.page.toString());
    if (params?.pageSize) queryParams.append('pageSize', params.pageSize.toString());
    if (params?.lessonId) queryParams.append('lessonId', params.lessonId);

    const queryString = queryParams.toString();
    return apiRequest(`/api/workspaces${queryString ? `?${queryString}` : ''}`, {}, true);
  },

  async getWorkspace(workspaceId: string): Promise<any> {
    return apiRequest(`/api/workspaces/${workspaceId}`, {}, true);
  },

  async updateWorkspace(
    workspaceId: string,
    data: {
      name?: string;
      files?: Array<{ path: string; content: string; language: string }>;
    }
  ): Promise<any> {
    return apiRequest(`/api/workspaces/${workspaceId}`, {
      method: 'PUT',
      body: JSON.stringify(data),
    }, true);
  },

  async deleteWorkspace(workspaceId: string): Promise<void> {
    return apiRequest(`/api/workspaces/${workspaceId}`, {
      method: 'DELETE',
    }, true);
  },

  // Code Execution
  async executeCode(workspaceId: string, code: string): Promise<{
    output: string;
    error?: string;
    executionTimeMs: number;
    executedAt: string;
  }> {
    return apiRequest(`/api/workspaces/${workspaceId}/execute`, {
      method: 'POST',
      body: JSON.stringify({ code }),
    }, true);
  },

  async getExecutionHistory(workspaceId: string): Promise<any[]> {
    return apiRequest(`/api/workspaces/${workspaceId}/executions`, {}, true);
  },

  // Workspace Sharing
  async shareWorkspace(
    workspaceId: string,
    permissions: 'read_only' | 'editable'
  ): Promise<{ token: string; shareUrl: string; expiresAt: string }> {
    return apiRequest(`/api/workspaces/${workspaceId}/share`, {
      method: 'POST',
      body: JSON.stringify({ permissions }),
    }, true);
  },

  async getSharedWorkspace(token: string): Promise<any> {
    return apiRequest(`/api/workspaces/shared/${token}`, {}, false);
  },

  // ========== PLAYGROUND API METHODS ==========
  // Code execution in playground

  async executePlaygroundCode(code: string): Promise<{
    output: string;
    error?: string;
    execution_time_ms: number;
    success: boolean;
  }> {
    return apiRequest('/api/v1/playground/execute', {
      method: 'POST',
      body: JSON.stringify({ code }),
    }, false); // No auth required for playground
  },

  // AI-powered code analysis
  async analyzeCode(code: string, language: string = 'go'): Promise<{
    complexity: number;
    issues: Array<{ severity: string; message: string; line: number }>;
    patterns: string[];
    strengths: string[];
    suggestions: Array<{ description: string; fix?: string }>;
    functions: Array<{ name: string; params: string[]; returnType: string }>;
    variables: Array<{ name: string; type: string; mutable: boolean }>;
  }> {
    return apiRequest('/api/v1/playground/analyze', {
      method: 'POST',
      body: JSON.stringify({ code, language }),
    }, false);
  },

  // AI code completion
  async completeCode(code: string, position: number): Promise<{
    completions: Array<{ text: string; type: string; documentation?: string }>;
  }> {
    return apiRequest('/api/v1/playground/complete', {
      method: 'POST',
      body: JSON.stringify({ code, position }),
    }, false);
  },

  // AI error explanation
  async explainError(code: string, error: string): Promise<{
    explanations: Array<{ message: string; line?: number; fix?: string; learnMore?: string }>;
  }> {
    return apiRequest('/api/v1/playground/explain', {
      method: 'POST',
      body: JSON.stringify({ code, error }),
    }, false);
  },

  // AI test generation
  async generateTests(code: string): Promise<{
    test_cases: Array<{ name: string; input: string; expected: string; description: string }>;
  }> {
    return apiRequest('/api/v1/playground/generate-tests', {
      method: 'POST',
      body: JSON.stringify({ code }),
    }, false);
  },

  // Execute code with AI analysis
  async executeWithAI(code: string, language: string = 'go', sessionId?: string): Promise<{
    output: string;
    error?: string;
    execution_time_ms: number;
    success: boolean;
    ai_analysis?: {
      complexity: number;
      issues: Array<{ severity: string; message: string; line: number }>;
      patterns: string[];
      strengths: string[];
      suggestions: Array<{ description: string; fix?: string }>;
    };
    test_results?: Array<{ testName: string; passed: boolean; expected: string; actual: string; error: string }>;
  }> {
    return apiRequest('/api/v1/playground/execute-ai', {
      method: 'POST',
      body: JSON.stringify({ code, language, session_id: sessionId }),
    }, false);
  },

  // Session management
  async createPlaygroundSession(): Promise<{ session_id: string }> {
    return apiRequest('/api/v1/playground/sessions', {
      method: 'POST',
      body: JSON.stringify({}),
    }, false);
  },

  async getPlaygroundSession(sessionId: string): Promise<{
    session_id: string;
    history?: Array<{ code: string; language: string; output: string; timestamp: string }>;
  }> {
    return apiRequest(`/api/v1/playground/sessions/${sessionId}`, {}, false);
  },

  async getSessionHistory(sessionId: string): Promise<{
    history: Array<{ code: string; language: string; output: string; timestamp: string }>;
  }> {
    return apiRequest(`/api/v1/playground/sessions/${sessionId}/history`, {}, false);
  },

  // ========== INTERVIEW API METHODS ==========
  // AI-powered mock interviews

  async startInterview(type: 'coding' | 'behavioral' | 'system_design', difficulty: 'beginner' | 'intermediate' | 'advanced'): Promise<{
    session: {
      id: string;
      user_id: string;
      type: string;
      difficulty: string;
      questions: Array<{
        id: string;
        content: string;
        type: string;
        difficulty: string;
        expected_points?: string[];
        time_limit: number;
      }>;
      current_index: number;
      answers: Array<{
        question_id: string;
        content: string;
        score?: number;
        feedback?: string;
        created_at: string;
      }>;
      status: string;
      score?: number;
      created_at: string;
      completed_at?: string;
    };
    first_question: {
      id: string;
      content: string;
      type: string;
      difficulty: string;
      expected_points?: string[];
      time_limit: number;
    };
  }> {
    return apiRequest('/api/v1/interview/start', {
      method: 'POST',
      body: JSON.stringify({ type, difficulty }),
    }, true);
  },

  async submitInterviewAnswer(sessionId: string, answer: string): Promise<{
    answer: {
      question_id: string;
      content: string;
      score?: number;
      feedback?: string;
      created_at: string;
    };
    completed: boolean;
    next_question?: {
      id: string;
      content: string;
      type: string;
      difficulty: string;
      expected_points?: string[];
      time_limit: number;
    };
    session?: {
      id: string;
      status: string;
      score?: number;
    };
  }> {
    return apiRequest('/api/v1/interview/answer', {
      method: 'POST',
      body: JSON.stringify({ session_id: sessionId, answer }),
    }, true);
  },

  async getInterviewFeedback(sessionId: string): Promise<{
    session_id: string;
    overall_score: number;
    strengths: string[];
    improvements: string[];
    detailed_feedback: Array<{
      question_id: string;
      score: number;
      feedback: string;
      strengths: string[];
      missed: string[];
    }>;
  }> {
    return apiRequest(`/api/v1/interview/feedback/${sessionId}`, {}, true);
  },

  async getInterviewSessions(): Promise<Array<{
    id: string;
    user_id?: string;
    type: string;
    difficulty: string;
    questions?: Array<{
      id: string;
      content: string;
      type: string;
      difficulty: string;
      time_limit: number;
    }>;
    status: string;
    score?: number;
    created_at: string;
    completed_at?: string;
  }>> {
    return apiRequest('/api/v1/interview/sessions', {}, true);
  },

  async getInterviewSession(sessionId: string): Promise<{
    id: string;
    user_id: string;
    type: string;
    difficulty: string;
    questions: Array<{
      id: string;
      content: string;
      type: string;
      difficulty: string;
      expected_points?: string[];
      time_limit: number;
    }>;
    current_index: number;
    answers: Array<{
      question_id: string;
      content: string;
      score?: number;
      feedback?: string;
      created_at: string;
    }>;
    status: string;
    score?: number;
    created_at: string;
    completed_at?: string;
  }> {
    return apiRequest(`/api/v1/interview/sessions/${sessionId}`, {}, true);
  },
};

export type {
  Curriculum,
  CurriculumPhase,
  CurriculumLesson,
  Project,
  LessonDetail,
  LessonExercise,
  APIResponse,
  ExerciseSubmission,
  ExerciseResult,
  TestResult,
  Progress,
  ProgressStats,
  ProgressListResponse,
  UpdateProgressPayload,
  BackendUser,
  AuthVerifyResponse,
  ProfileUpdatePayload,
};

export { APIError };
