import { apiClient } from '../client';
import {
  ApiResponse,
  PaginatedResponse,
  Challenge,
  Assessment,
  ChallengeSubmission,
  AssessmentSubmission,
  TestResult,
  UserProgress,
  ChallengeFilters,
  CreateChallengeSubmissionRequest,
  CreateAssessmentSubmissionRequest,
  UpdateUserProgressRequest,
  Leaderboard
} from '../types';

// Challenge API
export const challengeApi = {
  // Get all challenges with filtering
  getAll: async (filters?: ChallengeFilters, page: number = 1, limit: number = 20): Promise<ApiResponse<PaginatedResponse<Challenge>>> => {
    return apiClient.get('/challenges', { ...filters, page, limit });
  },

  // Get challenge by ID
  getById: async (id: string): Promise<ApiResponse<Challenge>> => {
    return apiClient.get(`/challenges/${id}`);
  },

  // Get challenge with user progress
  getWithProgress: async (id: string): Promise<ApiResponse<Challenge & { userProgress?: UserProgress }>> => {
    return apiClient.get(`/challenges/${id}/progress`);
  },

  // Submit challenge solution
  submit: async (data: CreateChallengeSubmissionRequest): Promise<ApiResponse<ChallengeSubmission>> => {
    return apiClient.post('/challenges/submit', data);
  },

  // Get submission results
  getSubmission: async (submissionId: string): Promise<ApiResponse<ChallengeSubmission>> => {
    return apiClient.get(`/challenges/submissions/${submissionId}`);
  },

  // Get user's submissions for a challenge
  getUserSubmissions: async (challengeId: string, page: number = 1, limit: number = 10): Promise<ApiResponse<PaginatedResponse<ChallengeSubmission>>> => {
    return apiClient.get(`/challenges/${challengeId}/submissions`, { page, limit });
  },

  // Run code without submitting
  runCode: async (challengeId: string, code: string, language: string): Promise<ApiResponse<TestResult[]>> => {
    return apiClient.post(`/challenges/${challengeId}/run`, { code, language });
  },

  // Get challenge hints
  getHints: async (challengeId: string): Promise<ApiResponse<string[]>> => {
    return apiClient.get(`/challenges/${challengeId}/hints`);
  },

  // Get challenge solution (if completed)
  getSolution: async (challengeId: string): Promise<ApiResponse<{ code: string; explanation: string }>> => {
    return apiClient.get(`/challenges/${challengeId}/solution`);
  },

  // Get challenge statistics
  getStats: async (challengeId: string): Promise<ApiResponse<{
    totalAttempts: number;
    successRate: number;
    averageAttempts: number;
    languageDistribution: Record<string, number>;
    difficultyRating: number;
  }>> => {
    return apiClient.get(`/challenges/${challengeId}/stats`);
  }
};

// Assessment API
export const assessmentApi = {
  // Get all assessments
  getAll: async (category?: string, difficulty?: string, page: number = 1, limit: number = 20): Promise<ApiResponse<PaginatedResponse<Assessment>>> => {
    return apiClient.get('/assessments', { category, difficulty, page, limit });
  },

  // Get assessment by ID
  getById: async (id: string): Promise<ApiResponse<Assessment>> => {
    return apiClient.get(`/assessments/${id}`);
  },

  // Start assessment (get questions without answers)
  start: async (id: string): Promise<ApiResponse<{
    assessmentId: string;
    sessionId: string;
    questions: Array<{
      id: string;
      type: string;
      question: string;
      options?: string[];
      points: number;
    }>;
    duration: number;
    startedAt: string;
  }>> => {
    return apiClient.post(`/assessments/${id}/start`);
  },

  // Submit assessment
  submit: async (data: CreateAssessmentSubmissionRequest): Promise<ApiResponse<AssessmentSubmission>> => {
    return apiClient.post('/assessments/submit', data);
  },

  // Get assessment results
  getSubmission: async (submissionId: string): Promise<ApiResponse<AssessmentSubmission & {
    detailedResults: Array<{
      questionId: string;
      question: string;
      userAnswer: string | number;
      correctAnswer: string | number;
      isCorrect: boolean;
      explanation?: string;
      points: number;
    }>;
  }>> => {
    return apiClient.get(`/assessments/submissions/${submissionId}`);
  },

  // Get user's assessment history
  getUserSubmissions: async (page: number = 1, limit: number = 10): Promise<ApiResponse<PaginatedResponse<AssessmentSubmission>>> => {
    return apiClient.get('/assessments/submissions', { page, limit });
  },

  // Get assessment statistics
  getStats: async (assessmentId: string): Promise<ApiResponse<{
    totalAttempts: number;
    averageScore: number;
    passRate: number;
    averageTime: number;
    scoreDistribution: Record<string, number>;
  }>> => {
    return apiClient.get(`/assessments/${assessmentId}/stats`);
  }
};

// Progress API
export const progressApi = {
  // Get user's overall progress
  getOverall: async (): Promise<ApiResponse<{
    challengesCompleted: number;
    assessmentsPassed: number;
    totalPoints: number;
    currentStreak: number;
    rank: number;
    level: number;
    experiencePoints: number;
    nextLevelXP: number;
  }>> => {
    return apiClient.get('/progress/overall');
  },

  // Update progress for specific item
  update: async (data: UpdateUserProgressRequest): Promise<ApiResponse<UserProgress>> => {
    return apiClient.post('/progress/update', data);
  },

  // Get progress for specific challenge
  getChallengeProgress: async (challengeId: string): Promise<ApiResponse<UserProgress>> => {
    return apiClient.get(`/progress/challenges/${challengeId}`);
  },

  // Get progress for specific assessment
  getAssessmentProgress: async (assessmentId: string): Promise<ApiResponse<UserProgress>> => {
    return apiClient.get(`/progress/assessments/${assessmentId}`);
  },

  // Get user's activity history
  getActivity: async (days: number = 30): Promise<ApiResponse<Array<{
    date: string;
    challengesCompleted: number;
    assessmentsCompleted: number;
    pointsEarned: number;
    timeSpent: number;
  }>>> => {
    return apiClient.get('/progress/activity', { days });
  },

  // Get skill levels
  getSkillLevels: async (): Promise<ApiResponse<Array<{
    skill: string;
    level: number;
    experience: number;
    nextLevelExperience: number;
    category: string;
    recentProgress: number;
  }>>> => {
    return apiClient.get('/progress/skills');
  }
};

// Leaderboard API
export const leaderboardApi = {
  // Get practice leaderboard
  get: async (period: 'daily' | 'weekly' | 'monthly' | 'all-time' = 'weekly', limit: number = 50): Promise<ApiResponse<Leaderboard>> => {
    return apiClient.get('/leaderboard/practice', { period, limit });
  },

  // Get challenge-specific leaderboard
  getChallenge: async (challengeId: string, limit: number = 20): Promise<ApiResponse<{
    challengeId: string;
    challengeTitle: string;
    entries: Array<{
      rank: number;
      user: {
        id: string;
        name: string;
        avatar: string;
        badge: string;
      };
      score: number;
      attempts: number;
      completedAt: string;
      language: string;
    }>;
  }>> => {
    return apiClient.get(`/leaderboard/challenges/${challengeId}`, { limit });
  },

  // Get user's rank
  getUserRank: async (period: 'daily' | 'weekly' | 'monthly' | 'all-time' = 'weekly'): Promise<ApiResponse<{
    rank: number;
    totalParticipants: number;
    points: number;
    percentile: number;
  }>> => {
    return apiClient.get('/leaderboard/rank', { period });
  }
};

// Statistics API
export const statsApi = {
  // Get practice statistics
  getPracticeStats: async (): Promise<ApiResponse<{
    totalChallenges: number;
    completedChallenges: number;
    totalAssessments: number;
    passedAssessments: number;
    averageScore: number;
    totalTimeSpent: number;
    favoriteCategory: string;
    strongestSkills: string[];
    improvementAreas: string[];
    recentActivity: Array<{
      date: string;
      activity: string;
      points: number;
    }>;
  }>> => {
    return apiClient.get('/stats/practice');
  },

  // Get detailed challenge statistics
  getChallengeStats: async (timeframe: 'week' | 'month' | 'year' = 'month'): Promise<ApiResponse<{
    completionRate: number;
    averageAttempts: number;
    categoryBreakdown: Record<string, {
      completed: number;
      total: number;
      averageScore: number;
    }>;
    difficultyBreakdown: Record<string, {
      completed: number;
      total: number;
      successRate: number;
    }>;
    languageUsage: Record<string, number>;
    timeSpentByCategory: Record<string, number>;
  }>> => {
    return apiClient.get('/stats/challenges', { timeframe });
  },

  // Get assessment performance
  getAssessmentStats: async (timeframe: 'week' | 'month' | 'year' = 'month'): Promise<ApiResponse<{
    totalTaken: number;
    totalPassed: number;
    averageScore: number;
    averageTime: number;
    categoryPerformance: Record<string, {
      taken: number;
      passed: number;
      averageScore: number;
    }>;
    improvementTrend: Array<{
      date: string;
      averageScore: number;
    }>;
    strongAreas: string[];
    weakAreas: string[];
  }>> => {
    return apiClient.get('/stats/assessments', { timeframe });
  }
};

// Export all practice-related APIs
export const practiceService = {
  challenges: challengeApi,
  assessments: assessmentApi,
  progress: progressApi,
  leaderboard: leaderboardApi,
  stats: statsApi,
};
