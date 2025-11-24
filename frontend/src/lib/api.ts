// API client for GO-PRO backend

const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';

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

async function apiRequest<T>(
  endpoint: string,
  options: RequestInit = {}
): Promise<T> {
  const url = `${API_BASE_URL}${endpoint}`;
  
  const defaultHeaders = {
    'Content-Type': 'application/json',
    'Accept': 'application/json',
  };

  const config: RequestInit = {
    ...options,
    headers: {
      ...defaultHeaders,
      ...options.headers,
    },
  };

  try {
    const response = await fetch(url, config);
    const data: APIResponse<T> = await response.json();

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

  // Progress
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
    return apiRequest(`/api/v1/progress/${userId}?page=${page}&page_size=${pageSize}`);
  },

  async updateProgress(
    userId: string,
    lessonId: string,
    data: { completed: boolean; score: number }
  ): Promise<any> {
    return apiRequest(`/api/v1/progress/${userId}/lesson/${lessonId}`, {
      method: 'POST',
      body: JSON.stringify(data),
    });
  },

  // Exercise submission
  async submitExercise(
    exerciseId: string,
    code: string,
    authToken?: string
  ): Promise<ExerciseResult> {
    const headers: Record<string, string> = {};
    if (authToken) {
      headers['Authorization'] = `Bearer ${authToken}`;
    }

    return apiRequest<ExerciseResult>(`/api/v1/exercises/${exerciseId}/submit`, {
      method: 'POST',
      headers,
      body: JSON.stringify({
        code,
        language: 'go',
      } as ExerciseSubmission),
    });
  },

  // Lesson completion
  async completeLesson(
    lessonId: string,
    authToken?: string
  ): Promise<void> {
    const headers: Record<string, string> = {};
    if (authToken) {
      headers['Authorization'] = `Bearer ${authToken}`;
    }

    return apiRequest<void>(`/api/v1/lessons/${lessonId}/complete`, {
      method: 'POST',
      headers,
    });
  },

  // User progress
  async getUserProgress(
    userId: string,
    page = 1,
    pageSize = 20,
    authToken?: string
  ): Promise<ProgressListResponse> {
    const headers: Record<string, string> = {};
    if (authToken) {
      headers['Authorization'] = `Bearer ${authToken}`;
    }

    return apiRequest<ProgressListResponse>(
      `/api/v1/users/${userId}/progress?page=${page}&pageSize=${pageSize}`,
      { headers }
    );
  },

  // Progress statistics
  async getProgressStats(
    userId: string,
    authToken?: string
  ): Promise<ProgressStats> {
    const headers: Record<string, string> = {};
    if (authToken) {
      headers['Authorization'] = `Bearer ${authToken}`;
    }

    return apiRequest<ProgressStats>(
      `/api/v1/users/${userId}/progress/stats`,
      { headers }
    );
  },

  // Update lesson progress
  async updateLessonProgress(
    userId: string,
    lessonId: string,
    status: 'not_started' | 'in_progress' | 'completed',
    score?: number,
    authToken?: string
  ): Promise<Progress> {
    const headers: Record<string, string> = {};
    if (authToken) {
      headers['Authorization'] = `Bearer ${authToken}`;
    }

    const payload: UpdateProgressPayload = { status };
    if (score !== undefined) {
      payload.score = score;
    }

    return apiRequest<Progress>(
      `/api/v1/users/${userId}/lessons/${lessonId}/progress`,
      {
        method: 'POST',
        headers,
        body: JSON.stringify(payload),
      }
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
};

export { APIError };
