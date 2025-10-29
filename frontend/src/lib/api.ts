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
};

export type {
  Curriculum,
  CurriculumPhase,
  CurriculumLesson,
  Project,
  LessonDetail,
  LessonExercise,
  APIResponse,
};

export { APIError };
