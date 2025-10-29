// Export API client and utilities
export { apiClient, apiCache, interceptors, withCache, withRetry, withPagination } from './client';

// Export types
export * from './types';

// Export services
export { practiceService } from './services/practice';
export { communityService } from './services/community';
export { projectService } from './services/projects';

// Re-export individual service APIs for convenience
export {
  challengeApi,
  assessmentApi,
  progressApi as practiceProgressApi,
  leaderboardApi,
  statsApi as practiceStatsApi
} from './services/practice';

export {
  forumApi,
  replyApi,
  userApi,
  eventsApi,
  categoriesApi,
  tagsApi,
  notificationsApi
} from './services/community';

export {
  projectsApi,
  chaptersApi,
  tasksApi,
  projectProgressApi,
  certificatesApi
} from './services/projects';

// Utility functions for common API operations
import { practiceService } from './services/practice';
import { communityService } from './services/community';
import { projectService } from './services/projects';
import { ApiResponse } from './types';

// Authentication utilities
export const authUtils = {
  // Set authentication token for all API calls
  setToken: (token: string) => {
    // apiClient.setAuthToken(token); // TODO: Implement setAuthToken method
    // Store in localStorage for persistence
    if (typeof window !== 'undefined') {
      localStorage.setItem('auth_token', token);
    }
  },

  // Remove authentication token
  removeToken: () => {
    // apiClient.removeAuthToken(); // TODO: Implement removeAuthToken method
    if (typeof window !== 'undefined') {
      localStorage.removeItem('auth_token');
    }
  },

  // Initialize auth from stored token
  initializeAuth: () => {
    if (typeof window !== 'undefined') {
      const token = localStorage.getItem('auth_token');
      if (token) {
        // apiClient.setAuthToken(token); // TODO: Implement setAuthToken method
      }
    }
  },

  // Check if user is authenticated
  isAuthenticated: (): boolean => {
    if (typeof window !== 'undefined') {
      return !!localStorage.getItem('auth_token');
    }
    return false;
  }
};

// Common API operations
export const commonApi = {
  // Get user dashboard data
  getDashboardData: async (): Promise<ApiResponse<{
    user: any;
    practiceStats: any;
    projectProgress: any;
    communityActivity: any;
    recentActivity: any[];
    achievements: any[];
    notifications: any[];
  }>> => {
    try {
      const [
        userResponse,
        practiceStatsResponse,
        projectProgressResponse
      ] = await Promise.all([
        communityService.users.getCurrentUser(),
        practiceService.stats.getPracticeStats(),
        projectService.progress.getOverallProgress(),
      ]);

      if (!userResponse.success) {
        return { success: false, error: 'Failed to load user data' };
      }

      return {
        success: true,
        data: {
          user: userResponse.data,
          practiceStats: practiceStatsResponse.data,
          projectProgress: projectProgressResponse.data,
          communityActivity: null, // Placeholder
          recentActivity: [],
          achievements: [],
          notifications: []
        }
      };
    } catch (error) {
      return {
        success: false,
        error: error instanceof Error ? error.message : 'Failed to load dashboard data'
      };
    }
  },

  // Search across all content types
  globalSearch: async (query: string, filters?: {
    type?: 'challenges' | 'projects' | 'posts' | 'users';
    category?: string;
    difficulty?: string;
  }): Promise<ApiResponse<{
    challenges: any[];
    projects: any[];
    posts: any[];
    users: any[];
    total: number;
  }>> => {
    try {
      const searchPromises = [];

      if (!filters?.type || filters.type === 'challenges') {
        searchPromises.push(
          practiceService.challenges.getAll({ search: query, ...filters })
        );
      }

      if (!filters?.type || filters.type === 'projects') {
        searchPromises.push(
          projectService.projects.getAll({ search: query, ...filters })
        );
      }

      if (!filters?.type || filters.type === 'posts') {
        searchPromises.push(
          communityService.forum.getPosts({ search: query, ...filters })
        );
      }

      if (!filters?.type || filters.type === 'users') {
        searchPromises.push(
          communityService.users.searchUsers(query)
        );
      }

      const results = await Promise.all(searchPromises);
      
      return {
        success: true,
        data: {
          challenges: results[0]?.data?.items || [],
          projects: results[1]?.data?.items || [],
          posts: results[2]?.data?.items || [],
          users: results[3]?.data?.items || [],
          total: results.reduce((sum, result) => sum + (result?.data?.total || result?.data?.items?.length || 0), 0)
        }
      };
    } catch (error) {
      return {
        success: false,
        error: error instanceof Error ? error.message : 'Search failed'
      };
    }
  },

  // Get user's complete activity feed
  getActivityFeed: async (days: number = 30): Promise<ApiResponse<Array<{
    id: string;
    type: 'challenge_completed' | 'project_started' | 'post_created' | 'reply_posted' | 'achievement_earned';
    title: string;
    description: string;
    timestamp: string;
    points?: number;
    url?: string;
    metadata?: any;
  }>>> => {
    try {
      // This would combine activity from different services
      // For now, return a placeholder structure
      return {
        success: true,
        data: []
      };
    } catch (error) {
      return {
        success: false,
        error: error instanceof Error ? error.message : 'Failed to load activity feed'
      };
    }
  },

  // Get comprehensive user statistics
  getUserStats: async (userId?: string): Promise<ApiResponse<{
    practice: any;
    projects: any;
    community: any;
    overall: {
      totalPoints: number;
      level: number;
      rank: number;
      joinedAt: string;
      lastActive: string;
      achievements: any[];
    };
  }>> => {
    try {
      const [practiceStats, projectProgress, userProfile] = await Promise.all([
        practiceService.stats.getPracticeStats(),
        projectService.progress.getOverallProgress(),
        userId ? communityService.users.getProfile(userId) : communityService.users.getCurrentUser()
      ]);

      return {
        success: true,
        data: {
          practice: practiceStats.data,
          projects: projectProgress.data,
          community: {}, // Placeholder
          overall: {
            totalPoints: (practiceStats.data as any)?.totalPoints || 0,
            level: (projectProgress.data as any)?.currentLevel || 1,
            rank: 0, // Would need to calculate
            joinedAt: (userProfile.data as any)?.joinedAt || '',
            lastActive: new Date().toISOString(),
            achievements: []
          }
        }
      };
    } catch (error) {
      return {
        success: false,
        error: error instanceof Error ? error.message : 'Failed to load user statistics'
      };
    }
  }
};

// Error handling utilities
export const handleApiError = (error: any): string => {
  if (error?.response?.data?.message) {
    return error.response.data.message;
  }
  if (error?.message) {
    return error.message;
  }
  return 'An unexpected error occurred';
};

// Loading state management
export const createLoadingState = () => {
  let loading = false;
  let error: string | null = null;

  return {
    setLoading: (state: boolean) => { loading = state; },
    setError: (err: string | null) => { error = err; },
    getState: () => ({ loading, error }),
    reset: () => { loading = false; error = null; }
  };
};

// Initialize authentication on import
if (typeof window !== 'undefined') {
  authUtils.initializeAuth();
}
