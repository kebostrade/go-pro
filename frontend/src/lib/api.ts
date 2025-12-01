// API client for GO-PRO backend
import { auth } from './firebase';

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
  const user = auth.currentUser;
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
