import { apiClient } from '../client';
import {
  ApiResponse,
  PaginatedResponse,
  Project,
  ProjectChapter,
  ProjectTask,
  UserProgress,
  ProjectFilters,
  UpdateUserProgressRequest
} from '../types';

// Projects API
export const projectsApi = {
  // Get all projects with filtering
  getAll: async (filters?: ProjectFilters, page: number = 1, limit: number = 20): Promise<ApiResponse<PaginatedResponse<Project>>> => {
    return apiClient.get('/projects', { ...filters, page, limit });
  },

  // Get project by ID
  getById: async (id: string): Promise<ApiResponse<Project>> => {
    return apiClient.get(`/projects/${id}`);
  },

  // Get project with user progress
  getWithProgress: async (id: string): Promise<ApiResponse<Project & { 
    userProgress?: UserProgress;
    chaptersProgress: Record<string, UserProgress>;
    tasksProgress: Record<string, UserProgress>;
    overallProgress: number;
    completedChapters: number;
    totalChapters: number;
  }>> => {
    return apiClient.get(`/projects/${id}/progress`);
  },

  // Get featured projects
  getFeatured: async (limit: number = 6): Promise<ApiResponse<Project[]>> => {
    return apiClient.get('/projects/featured', { limit });
  },

  // Get recommended projects for user
  getRecommended: async (limit: number = 6): Promise<ApiResponse<Project[]>> => {
    return apiClient.get('/projects/recommended', { limit });
  },

  // Get project statistics
  getStats: async (projectId: string): Promise<ApiResponse<{
    totalEnrolled: number;
    completionRate: number;
    averageRating: number;
    totalRatings: number;
    averageCompletionTime: number;
    difficultyRating: number;
    skillsDistribution: Record<string, number>;
    userFeedback: Array<{
      rating: number;
      comment: string;
      user: string;
      completedAt: string;
    }>;
  }>> => {
    return apiClient.get(`/projects/${projectId}/stats`);
  },

  // Enroll in project
  enroll: async (projectId: string): Promise<ApiResponse<{ enrolled: boolean; startedAt: string }>> => {
    return apiClient.post(`/projects/${projectId}/enroll`);
  },

  // Unenroll from project
  unenroll: async (projectId: string): Promise<ApiResponse<{ unenrolled: boolean }>> => {
    return apiClient.delete(`/projects/${projectId}/enroll`);
  },

  // Rate project
  rateProject: async (projectId: string, rating: number, comment?: string): Promise<ApiResponse<{
    rating: number;
    comment?: string;
    submittedAt: string;
  }>> => {
    return apiClient.post(`/projects/${projectId}/rate`, { rating, comment });
  },

  // Get user's enrolled projects
  getUserProjects: async (status?: 'not-started' | 'in-progress' | 'completed', page: number = 1, limit: number = 20): Promise<ApiResponse<PaginatedResponse<Project & { userProgress: UserProgress }>>> => {
    return apiClient.get('/projects/my-projects', { status, page, limit });
  }
};

// Chapters API
export const chaptersApi = {
  // Get all chapters for a project
  getByProject: async (projectId: string): Promise<ApiResponse<ProjectChapter[]>> => {
    return apiClient.get(`/projects/${projectId}/chapters`);
  },

  // Get chapter by ID
  getById: async (chapterId: string): Promise<ApiResponse<ProjectChapter>> => {
    return apiClient.get(`/chapters/${chapterId}`);
  },

  // Get chapter with user progress
  getWithProgress: async (chapterId: string): Promise<ApiResponse<ProjectChapter & {
    userProgress?: UserProgress;
    tasksProgress: Record<string, UserProgress>;
    overallProgress: number;
    completedTasks: number;
    totalTasks: number;
    estimatedTimeRemaining: number;
  }>> => {
    return apiClient.get(`/chapters/${chapterId}/progress`);
  },

  // Start chapter
  startChapter: async (chapterId: string): Promise<ApiResponse<{ started: boolean; startedAt: string }>> => {
    return apiClient.post(`/chapters/${chapterId}/start`);
  },

  // Complete chapter
  completeChapter: async (chapterId: string, timeSpent: number): Promise<ApiResponse<{
    completed: boolean;
    completedAt: string;
    pointsEarned: number;
    nextChapter?: string;
  }>> => {
    return apiClient.post(`/chapters/${chapterId}/complete`, { timeSpent });
  },

  // Get chapter resources
  getResources: async (chapterId: string): Promise<ApiResponse<Array<{
    id: string;
    title: string;
    type: 'link' | 'file' | 'video' | 'documentation';
    url: string;
    description?: string;
  }>>> => {
    return apiClient.get(`/chapters/${chapterId}/resources`);
  },

  // Submit chapter exercise
  submitExercise: async (chapterId: string, exerciseId: string, solution: string): Promise<ApiResponse<{
    submitted: boolean;
    feedback?: string;
    score?: number;
    passed: boolean;
  }>> => {
    return apiClient.post(`/chapters/${chapterId}/exercises/${exerciseId}/submit`, { solution });
  }
};

// Tasks API
export const tasksApi = {
  // Get all tasks for a chapter
  getByChapter: async (chapterId: string): Promise<ApiResponse<ProjectTask[]>> => {
    return apiClient.get(`/chapters/${chapterId}/tasks`);
  },

  // Get task by ID
  getById: async (taskId: string): Promise<ApiResponse<ProjectTask>> => {
    return apiClient.get(`/tasks/${taskId}`);
  },

  // Get task with user progress
  getWithProgress: async (taskId: string): Promise<ApiResponse<ProjectTask & {
    userProgress?: UserProgress;
    completed: boolean;
    startedAt?: string;
    completedAt?: string;
  }>> => {
    return apiClient.get(`/tasks/${taskId}/progress`);
  },

  // Start task
  startTask: async (taskId: string): Promise<ApiResponse<{ started: boolean; startedAt: string }>> => {
    return apiClient.post(`/tasks/${taskId}/start`);
  },

  // Complete task
  completeTask: async (taskId: string, timeSpent: number, notes?: string): Promise<ApiResponse<{
    completed: boolean;
    completedAt: string;
    pointsEarned: number;
  }>> => {
    return apiClient.post(`/tasks/${taskId}/complete`, { timeSpent, notes });
  },

  // Skip task
  skipTask: async (taskId: string, reason?: string): Promise<ApiResponse<{ skipped: boolean; skippedAt: string }>> => {
    return apiClient.post(`/tasks/${taskId}/skip`, { reason });
  },

  // Get task hints
  getHints: async (taskId: string): Promise<ApiResponse<string[]>> => {
    return apiClient.get(`/tasks/${taskId}/hints`);
  },

  // Submit task solution
  submitSolution: async (taskId: string, solution: string, language?: string): Promise<ApiResponse<{
    submitted: boolean;
    feedback?: string;
    score?: number;
    passed: boolean;
    suggestions?: string[];
  }>> => {
    return apiClient.post(`/tasks/${taskId}/submit`, { solution, language });
  }
};

// Progress API
export const projectProgressApi = {
  // Get overall project progress for user
  getOverallProgress: async (): Promise<ApiResponse<{
    totalProjects: number;
    enrolledProjects: number;
    completedProjects: number;
    inProgressProjects: number;
    totalTimeSpent: number;
    averageCompletionRate: number;
    skillsAcquired: string[];
    certificatesEarned: number;
    currentLevel: number;
    experiencePoints: number;
    nextLevelXP: number;
  }>> => {
    return apiClient.get('/projects/progress/overall');
  },

  // Update progress for specific item
  updateProgress: async (data: UpdateUserProgressRequest): Promise<ApiResponse<UserProgress>> => {
    return apiClient.post('/projects/progress/update', data);
  },

  // Get project progress
  getProjectProgress: async (projectId: string): Promise<ApiResponse<UserProgress & {
    chaptersCompleted: number;
    totalChapters: number;
    tasksCompleted: number;
    totalTasks: number;
    skillsLearned: string[];
    timeSpent: number;
    lastActivity: string;
  }>> => {
    return apiClient.get(`/projects/${projectId}/progress`);
  },

  // Get chapter progress
  getChapterProgress: async (chapterId: string): Promise<ApiResponse<UserProgress & {
    tasksCompleted: number;
    totalTasks: number;
    exercisesCompleted: number;
    totalExercises: number;
    timeSpent: number;
  }>> => {
    return apiClient.get(`/chapters/${chapterId}/progress`);
  },

  // Get learning path progress
  getLearningPath: async (userId?: string): Promise<ApiResponse<{
    currentProject?: Project;
    nextRecommendedProjects: Project[];
    completedProjects: Project[];
    skillsProgress: Array<{
      skill: string;
      level: number;
      experience: number;
      projectsContributing: string[];
    }>;
    learningGoals: Array<{
      id: string;
      title: string;
      description: string;
      targetDate: string;
      progress: number;
      relatedProjects: string[];
    }>;
  }>> => {
    return apiClient.get('/projects/learning-path', userId ? { userId } : {});
  },

  // Set learning goals
  setLearningGoals: async (goals: Array<{
    title: string;
    description: string;
    targetDate: string;
    relatedProjects: string[];
  }>): Promise<ApiResponse<{
    goalsSet: number;
    estimatedCompletionTime: number;
  }>> => {
    return apiClient.post('/projects/learning-goals', { goals });
  }
};

// Certificates API
export const certificatesApi = {
  // Get user's certificates
  getUserCertificates: async (page: number = 1, limit: number = 20): Promise<ApiResponse<PaginatedResponse<{
    id: string;
    projectId: string;
    projectTitle: string;
    issuedAt: string;
    certificateUrl: string;
    verificationCode: string;
    skills: string[];
    grade: 'A' | 'B' | 'C';
    completionTime: number;
  }>>> => {
    return apiClient.get('/certificates', { page, limit });
  },

  // Generate certificate for completed project
  generateCertificate: async (projectId: string): Promise<ApiResponse<{
    certificateId: string;
    certificateUrl: string;
    verificationCode: string;
    issuedAt: string;
  }>> => {
    return apiClient.post(`/projects/${projectId}/certificate`);
  },

  // Verify certificate
  verifyCertificate: async (verificationCode: string): Promise<ApiResponse<{
    valid: boolean;
    certificate?: {
      id: string;
      projectTitle: string;
      userName: string;
      issuedAt: string;
      skills: string[];
    };
  }>> => {
    return apiClient.get(`/certificates/verify/${verificationCode}`);
  },

  // Share certificate
  shareCertificate: async (certificateId: string, platform: 'linkedin' | 'twitter' | 'facebook'): Promise<ApiResponse<{
    shareUrl: string;
    message: string;
  }>> => {
    return apiClient.post(`/certificates/${certificateId}/share`, { platform });
  }
};

// Export all project-related APIs
export const projectService = {
  projects: projectsApi,
  chapters: chaptersApi,
  tasks: tasksApi,
  progress: projectProgressApi,
  certificates: certificatesApi,
};
