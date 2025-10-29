import { apiClient } from '../client';
import {
  ApiResponse,
  PaginatedResponse,
  ForumPost,
  ForumReply,
  User,
  CommunityEvent,
  ForumPostFilters,
  CreateForumPostRequest,
  CreateForumReplyRequest
} from '../types';

// Forum API
export const forumApi = {
  // Get all forum posts with filtering
  getPosts: async (filters?: ForumPostFilters, page: number = 1, limit: number = 20): Promise<ApiResponse<PaginatedResponse<ForumPost>>> => {
    return apiClient.get('/forum/posts', { ...filters, page, limit });
  },

  // Get post by ID with replies
  getPost: async (id: string): Promise<ApiResponse<ForumPost>> => {
    return apiClient.get(`/forum/posts/${id}`);
  },

  // Create new forum post
  createPost: async (data: CreateForumPostRequest): Promise<ApiResponse<ForumPost>> => {
    return apiClient.post('/forum/posts', data);
  },

  // Update forum post
  updatePost: async (id: string, data: Partial<CreateForumPostRequest>): Promise<ApiResponse<ForumPost>> => {
    return apiClient.patch(`/forum/posts/${id}`, data);
  },

  // Delete forum post
  deletePost: async (id: string): Promise<ApiResponse<void>> => {
    return apiClient.delete(`/forum/posts/${id}`);
  },

  // Like/unlike a post
  toggleLike: async (id: string): Promise<ApiResponse<{ liked: boolean; totalLikes: number }>> => {
    return apiClient.post(`/forum/posts/${id}/like`);
  },

  // Dislike/undislike a post
  toggleDislike: async (id: string): Promise<ApiResponse<{ disliked: boolean; totalDislikes: number }>> => {
    return apiClient.post(`/forum/posts/${id}/dislike`);
  },

  // Mark post as solved
  markSolved: async (id: string): Promise<ApiResponse<ForumPost>> => {
    return apiClient.patch(`/forum/posts/${id}/solve`);
  },

  // Pin/unpin post (moderator only)
  togglePin: async (id: string): Promise<ApiResponse<{ pinned: boolean }>> => {
    return apiClient.patch(`/forum/posts/${id}/pin`);
  },

  // Lock/unlock post (moderator only)
  toggleLock: async (id: string): Promise<ApiResponse<{ locked: boolean }>> => {
    return apiClient.patch(`/forum/posts/${id}/lock`);
  },

  // Report post
  reportPost: async (id: string, reason: string): Promise<ApiResponse<void>> => {
    return apiClient.post(`/forum/posts/${id}/report`, { reason });
  },

  // Get post views
  incrementViews: async (id: string): Promise<ApiResponse<{ views: number }>> => {
    return apiClient.post(`/forum/posts/${id}/view`);
  },

  // Get trending posts
  getTrending: async (timeframe: 'day' | 'week' | 'month' = 'week', limit: number = 10): Promise<ApiResponse<ForumPost[]>> => {
    return apiClient.get('/forum/posts/trending', { timeframe, limit });
  },

  // Get user's posts
  getUserPosts: async (userId: string, page: number = 1, limit: number = 20): Promise<ApiResponse<PaginatedResponse<ForumPost>>> => {
    return apiClient.get(`/forum/users/${userId}/posts`, { page, limit });
  }
};

// Reply API
export const replyApi = {
  // Get replies for a post
  getReplies: async (postId: string, page: number = 1, limit: number = 20): Promise<ApiResponse<PaginatedResponse<ForumReply>>> => {
    return apiClient.get(`/forum/posts/${postId}/replies`, { page, limit });
  },

  // Create new reply
  createReply: async (data: CreateForumReplyRequest): Promise<ApiResponse<ForumReply>> => {
    return apiClient.post('/forum/replies', data);
  },

  // Update reply
  updateReply: async (id: string, content: string): Promise<ApiResponse<ForumReply>> => {
    return apiClient.patch(`/forum/replies/${id}`, { content });
  },

  // Delete reply
  deleteReply: async (id: string): Promise<ApiResponse<void>> => {
    return apiClient.delete(`/forum/replies/${id}`);
  },

  // Like/unlike a reply
  toggleLike: async (id: string): Promise<ApiResponse<{ liked: boolean; totalLikes: number }>> => {
    return apiClient.post(`/forum/replies/${id}/like`);
  },

  // Dislike/undislike a reply
  toggleDislike: async (id: string): Promise<ApiResponse<{ disliked: boolean; totalDislikes: number }>> => {
    return apiClient.post(`/forum/replies/${id}/dislike`);
  },

  // Accept reply as answer
  acceptAnswer: async (id: string): Promise<ApiResponse<ForumReply>> => {
    return apiClient.patch(`/forum/replies/${id}/accept`);
  },

  // Report reply
  reportReply: async (id: string, reason: string): Promise<ApiResponse<void>> => {
    return apiClient.post(`/forum/replies/${id}/report`, { reason });
  },

  // Get user's replies
  getUserReplies: async (userId: string, page: number = 1, limit: number = 20): Promise<ApiResponse<PaginatedResponse<ForumReply>>> => {
    return apiClient.get(`/forum/users/${userId}/replies`, { page, limit });
  }
};

// User API
export const userApi = {
  // Get user profile
  getProfile: async (id: string): Promise<ApiResponse<User>> => {
    return apiClient.get(`/users/${id}`);
  },

  // Get current user profile
  getCurrentUser: async (): Promise<ApiResponse<User>> => {
    return apiClient.get('/users/me');
  },

  // Update user profile
  updateProfile: async (data: Partial<User>): Promise<ApiResponse<User>> => {
    return apiClient.patch('/users/me', data);
  },

  // Upload avatar
  uploadAvatar: async (file: File): Promise<ApiResponse<{ avatarUrl: string }>> => {
    return apiClient.upload('/users/me/avatar', file);
  },

  // Follow/unfollow user
  toggleFollow: async (userId: string): Promise<ApiResponse<{ following: boolean; followerCount: number }>> => {
    return apiClient.post(`/users/${userId}/follow`);
  },

  // Get user's followers
  getFollowers: async (userId: string, page: number = 1, limit: number = 20): Promise<ApiResponse<PaginatedResponse<User>>> => {
    return apiClient.get(`/users/${userId}/followers`, { page, limit });
  },

  // Get user's following
  getFollowing: async (userId: string, page: number = 1, limit: number = 20): Promise<ApiResponse<PaginatedResponse<User>>> => {
    return apiClient.get(`/users/${userId}/following`, { page, limit });
  },

  // Search users
  searchUsers: async (query: string, page: number = 1, limit: number = 20): Promise<ApiResponse<PaginatedResponse<User>>> => {
    return apiClient.get('/users/search', { q: query, page, limit });
  },

  // Get top contributors
  getTopContributors: async (period: 'week' | 'month' | 'year' = 'month', limit: number = 20): Promise<ApiResponse<User[]>> => {
    return apiClient.get('/users/top-contributors', { period, limit });
  },

  // Get user activity
  getUserActivity: async (userId: string, days: number = 30): Promise<ApiResponse<Array<{
    date: string;
    type: 'post' | 'reply' | 'like' | 'accepted_answer';
    title: string;
    url: string;
    points: number;
  }>>> => {
    return apiClient.get(`/users/${userId}/activity`, { days });
  },

  // Get user statistics
  getUserStats: async (userId: string): Promise<ApiResponse<{
    totalPosts: number;
    totalReplies: number;
    acceptedAnswers: number;
    reputation: number;
    joinedAt: string;
    lastActive: string;
    badges: string[];
    achievements: Array<{
      id: string;
      title: string;
      description: string;
      unlockedAt: string;
    }>;
  }>> => {
    return apiClient.get(`/users/${userId}/stats`);
  }
};

// Events API
export const eventsApi = {
  // Get all events
  getEvents: async (
    type?: 'workshop' | 'webinar' | 'challenge' | 'meetup',
    upcoming?: boolean,
    page: number = 1,
    limit: number = 20
  ): Promise<ApiResponse<PaginatedResponse<CommunityEvent>>> => {
    return apiClient.get('/events', { type, upcoming, page, limit });
  },

  // Get event by ID
  getEvent: async (id: string): Promise<ApiResponse<CommunityEvent>> => {
    return apiClient.get(`/events/${id}`);
  },

  // Register for event
  registerForEvent: async (id: string): Promise<ApiResponse<{ registered: boolean; waitlisted?: boolean }>> => {
    return apiClient.post(`/events/${id}/register`);
  },

  // Unregister from event
  unregisterFromEvent: async (id: string): Promise<ApiResponse<{ unregistered: boolean }>> => {
    return apiClient.delete(`/events/${id}/register`);
  },

  // Get user's registered events
  getUserEvents: async (upcoming?: boolean, page: number = 1, limit: number = 20): Promise<ApiResponse<PaginatedResponse<CommunityEvent>>> => {
    return apiClient.get('/events/my-events', { upcoming, page, limit });
  },

  // Get event participants
  getParticipants: async (eventId: string, page: number = 1, limit: number = 50): Promise<ApiResponse<PaginatedResponse<User>>> => {
    return apiClient.get(`/events/${eventId}/participants`, { page, limit });
  }
};

// Categories API
export const categoriesApi = {
  // Get all forum categories
  getCategories: async (): Promise<ApiResponse<Array<{
    id: string;
    name: string;
    description: string;
    postCount: number;
    color: string;
    icon: string;
  }>>> => {
    return apiClient.get('/forum/categories');
  },

  // Get category statistics
  getCategoryStats: async (categoryId: string): Promise<ApiResponse<{
    totalPosts: number;
    totalReplies: number;
    activeUsers: number;
    topContributors: User[];
    recentActivity: Array<{
      type: 'post' | 'reply';
      title: string;
      author: string;
      createdAt: string;
    }>;
  }>> => {
    return apiClient.get(`/forum/categories/${categoryId}/stats`);
  }
};

// Tags API
export const tagsApi = {
  // Get popular tags
  getPopularTags: async (limit: number = 50): Promise<ApiResponse<Array<{
    name: string;
    count: number;
    category?: string;
  }>>> => {
    return apiClient.get('/forum/tags/popular', { limit });
  },

  // Search tags
  searchTags: async (query: string, limit: number = 20): Promise<ApiResponse<Array<{
    name: string;
    count: number;
  }>>> => {
    return apiClient.get('/forum/tags/search', { q: query, limit });
  },

  // Get posts by tag
  getPostsByTag: async (tag: string, page: number = 1, limit: number = 20): Promise<ApiResponse<PaginatedResponse<ForumPost>>> => {
    return apiClient.get(`/forum/tags/${tag}/posts`, { page, limit });
  }
};

// Notifications API
export const notificationsApi = {
  // Get user notifications
  getNotifications: async (unreadOnly?: boolean, page: number = 1, limit: number = 20): Promise<ApiResponse<PaginatedResponse<{
    id: string;
    type: 'reply' | 'like' | 'mention' | 'follow' | 'accepted_answer';
    title: string;
    message: string;
    read: boolean;
    createdAt: string;
    relatedUrl?: string;
    actor?: {
      id: string;
      name: string;
      avatar: string;
    };
  }>>> => {
    return apiClient.get('/notifications', { unreadOnly, page, limit });
  },

  // Mark notification as read
  markAsRead: async (id: string): Promise<ApiResponse<void>> => {
    return apiClient.patch(`/notifications/${id}/read`);
  },

  // Mark all notifications as read
  markAllAsRead: async (): Promise<ApiResponse<void>> => {
    return apiClient.patch('/notifications/read-all');
  },

  // Get unread count
  getUnreadCount: async (): Promise<ApiResponse<{ count: number }>> => {
    return apiClient.get('/notifications/unread-count');
  }
};

// Export all community-related APIs
export const communityService = {
  forum: forumApi,
  replies: replyApi,
  users: userApi,
  events: eventsApi,
  categories: categoriesApi,
  tags: tagsApi,
  notifications: notificationsApi,
};
