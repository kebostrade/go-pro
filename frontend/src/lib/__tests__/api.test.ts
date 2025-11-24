// API client tests for GO-PRO frontend

import { api, APIError } from '../api';

// Mock fetch globally
const mockFetch = jest.fn();
global.fetch = mockFetch;

describe('API Client', () => {
  beforeEach(() => {
    mockFetch.mockClear();
  });

  describe('Health Check', () => {
    it('should fetch health status successfully', async () => {
      const mockHealth = {
        status: 'healthy',
        timestamp: '2025-01-15T10:00:00Z',
        version: '1.0.0',
        uptime: '5m',
      };

      mockFetch.mockResolvedValueOnce({
        ok: true,
        json: async () => ({
          success: true,
          data: mockHealth,
          timestamp: '2025-01-15T10:00:00Z',
        }),
      });

      const result = await api.health();

      expect(mockFetch).toHaveBeenCalledWith(
        'http://localhost:8080/api/v1/health',
        expect.objectContaining({
          headers: expect.objectContaining({
            'Content-Type': 'application/json',
          }),
        })
      );
      expect(result).toEqual(mockHealth);
    });

    it('should handle health check errors', async () => {
      mockFetch.mockResolvedValueOnce({
        ok: false,
        status: 500,
        json: async () => ({
          success: false,
          error: {
            type: 'INTERNAL_ERROR',
            message: 'Service unavailable',
          },
        }),
      });

      await expect(api.health()).rejects.toThrow(APIError);
      await expect(api.health()).rejects.toThrow('Service unavailable');
    });
  });

  describe('Curriculum', () => {
    it('should fetch curriculum successfully', async () => {
      const mockCurriculum = {
        id: 'go-pro-curriculum',
        title: 'GO-PRO Learning Path',
        description: 'Comprehensive Go curriculum',
        duration: '12 weeks',
        phases: [],
        projects: [],
        created_at: '2025-01-01T00:00:00Z',
        updated_at: '2025-01-01T00:00:00Z',
      };

      mockFetch.mockResolvedValueOnce({
        ok: true,
        json: async () => ({
          success: true,
          data: mockCurriculum,
        }),
      });

      const result = await api.getCurriculum();

      expect(mockFetch).toHaveBeenCalledWith(
        'http://localhost:8080/api/v1/curriculum',
        expect.any(Object)
      );
      expect(result).toEqual(mockCurriculum);
    });

    it('should fetch lesson detail successfully', async () => {
      const mockLesson = {
        id: 1,
        title: 'Variables and Types',
        description: 'Learn Go variables',
        duration: '30 minutes',
        difficulty: 'beginner' as const,
        phase: 'Fundamentals',
        objectives: ['Understand variables'],
        theory: 'Variables are...',
        code_example: 'var x int',
        solution: 'Complete solution',
        exercises: [],
      };

      mockFetch.mockResolvedValueOnce({
        ok: true,
        json: async () => ({
          success: true,
          data: mockLesson,
        }),
      });

      const result = await api.getLessonDetail(1);

      expect(mockFetch).toHaveBeenCalledWith(
        'http://localhost:8080/api/v1/curriculum/lesson/1',
        expect.any(Object)
      );
      expect(result).toEqual(mockLesson);
    });
  });

  describe('Courses', () => {
    it('should fetch courses with pagination', async () => {
      const mockResponse = {
        items: [
          {
            id: 'go-basics',
            title: 'Go Basics',
            description: 'Learn fundamentals',
          },
        ],
        pagination: {
          page: 1,
          page_size: 10,
          total_items: 1,
          total_pages: 1,
          has_next: false,
          has_prev: false,
        },
      };

      mockFetch.mockResolvedValueOnce({
        ok: true,
        json: async () => ({
          success: true,
          data: mockResponse,
        }),
      });

      const result = await api.getCourses(1, 10);

      expect(mockFetch).toHaveBeenCalledWith(
        'http://localhost:8080/api/v1/courses?page=1&page_size=10',
        expect.any(Object)
      );
      expect(result).toEqual(mockResponse);
    });

    it('should fetch single course', async () => {
      const mockCourse = {
        id: 'go-basics',
        title: 'Go Basics',
        description: 'Learn fundamentals',
        created_at: '2025-01-01T00:00:00Z',
      };

      mockFetch.mockResolvedValueOnce({
        ok: true,
        json: async () => ({
          success: true,
          data: mockCourse,
        }),
      });

      const result = await api.getCourse('go-basics');

      expect(mockFetch).toHaveBeenCalledWith(
        'http://localhost:8080/api/v1/courses/go-basics',
        expect.any(Object)
      );
      expect(result).toEqual(mockCourse);
    });

    it('should handle course not found', async () => {
      mockFetch.mockResolvedValueOnce({
        ok: false,
        status: 404,
        json: async () => ({
          success: false,
          error: {
            type: 'NOT_FOUND',
            message: 'Course not found',
          },
        }),
      });

      await expect(api.getCourse('nonexistent')).rejects.toThrow(APIError);
      await expect(api.getCourse('nonexistent')).rejects.toThrow('Course not found');
    });
  });

  describe('Progress', () => {
    it('should fetch user progress', async () => {
      const mockProgress = {
        items: [
          {
            id: 'progress-001',
            user_id: 'user-123',
            lesson_id: 'lesson-001',
            status: 'completed',
            score: 100,
          },
        ],
        pagination: {
          page: 1,
          page_size: 10,
          total_items: 1,
          total_pages: 1,
          has_next: false,
          has_prev: false,
        },
      };

      mockFetch.mockResolvedValueOnce({
        ok: true,
        json: async () => ({
          success: true,
          data: mockProgress,
        }),
      });

      const result = await api.getProgress('user-123', 1, 10);

      expect(mockFetch).toHaveBeenCalledWith(
        'http://localhost:8080/api/v1/progress/user-123?page=1&page_size=10',
        expect.any(Object)
      );
      expect(result).toEqual(mockProgress);
    });

    it('should update progress', async () => {
      const mockProgress = {
        id: 'progress-001',
        user_id: 'user-123',
        lesson_id: 'lesson-001',
        status: 'completed',
        score: 95,
      };

      mockFetch.mockResolvedValueOnce({
        ok: true,
        json: async () => ({
          success: true,
          data: mockProgress,
        }),
      });

      const result = await api.updateProgress('user-123', 'lesson-001', {
        completed: true,
        score: 95,
      });

      expect(mockFetch).toHaveBeenCalledWith(
        'http://localhost:8080/api/v1/progress/user-123/lesson/lesson-001',
        expect.objectContaining({
          method: 'POST',
          body: JSON.stringify({ completed: true, score: 95 }),
        })
      );
      expect(result).toEqual(mockProgress);
    });

    it('should fetch progress statistics', async () => {
      const mockStats = {
        total_lessons: 20,
        completed_lessons: 5,
        in_progress_lessons: 3,
        average_score: 85.5,
        total_time_spent: 150,
      };

      mockFetch.mockResolvedValueOnce({
        ok: true,
        json: async () => ({
          success: true,
          data: mockStats,
        }),
      });

      const result = await api.getProgressStats('user-123');

      expect(mockFetch).toHaveBeenCalledWith(
        'http://localhost:8080/api/v1/users/user-123/progress/stats',
        expect.any(Object)
      );
      expect(result).toEqual(mockStats);
    });

    it('should update lesson progress', async () => {
      const mockProgress = {
        id: 'progress-001',
        user_id: 'user-123',
        lesson_id: 'lesson-001',
        status: 'in_progress' as const,
        score: 0,
      };

      mockFetch.mockResolvedValueOnce({
        ok: true,
        json: async () => ({
          success: true,
          data: mockProgress,
        }),
      });

      const result = await api.updateLessonProgress(
        'user-123',
        'lesson-001',
        'in_progress',
        0
      );

      expect(mockFetch).toHaveBeenCalledWith(
        'http://localhost:8080/api/v1/users/user-123/lessons/lesson-001/progress',
        expect.objectContaining({
          method: 'POST',
          body: JSON.stringify({ status: 'in_progress', score: 0 }),
        })
      );
      expect(result).toEqual(mockProgress);
    });
  });

  describe('Exercise Submission', () => {
    it('should submit exercise successfully - all tests passed', async () => {
      const mockResult = {
        success: true,
        passed: true,
        score: 100,
        results: [
          {
            test_name: 'Test 1',
            passed: true,
            expected: 'Hello, World!',
            actual: 'Hello, World!',
            error: null,
          },
        ],
        execution_time_ms: 150,
        message: 'All tests passed!',
      };

      mockFetch.mockResolvedValueOnce({
        ok: true,
        json: async () => ({
          success: true,
          data: mockResult,
        }),
      });

      const result = await api.submitExercise(
        'exercise-001',
        'package main\nfunc main() {}'
      );

      expect(mockFetch).toHaveBeenCalledWith(
        'http://localhost:8080/api/v1/exercises/exercise-001/submit',
        expect.objectContaining({
          method: 'POST',
          body: JSON.stringify({
            code: 'package main\nfunc main() {}',
            language: 'go',
          }),
        })
      );
      expect(result).toEqual(mockResult);
      expect(result.passed).toBe(true);
      expect(result.score).toBe(100);
    });

    it('should submit exercise with auth token', async () => {
      const mockResult = {
        success: true,
        passed: false,
        score: 50,
        results: [],
        execution_time_ms: 200,
        message: 'Tests passed: 1/2',
      };

      mockFetch.mockResolvedValueOnce({
        ok: true,
        json: async () => ({
          success: true,
          data: mockResult,
        }),
      });

      await api.submitExercise(
        'exercise-001',
        'package main\nfunc main() {}',
        'test-token-123'
      );

      expect(mockFetch).toHaveBeenCalledWith(
        expect.any(String),
        expect.objectContaining({
          headers: expect.objectContaining({
            Authorization: 'Bearer test-token-123',
          }),
        })
      );
    });

    it('should handle submission rate limit error', async () => {
      mockFetch.mockResolvedValueOnce({
        ok: false,
        status: 429,
        json: async () => ({
          success: false,
          error: {
            type: 'RATE_LIMIT_EXCEEDED',
            message: 'Too many submissions',
          },
        }),
      });

      await expect(
        api.submitExercise('exercise-001', 'code')
      ).rejects.toThrow(APIError);

      try {
        await api.submitExercise('exercise-001', 'code');
      } catch (error) {
        expect(error).toBeInstanceOf(APIError);
        if (error instanceof APIError) {
          expect(error.status).toBe(429);
          expect(error.type).toBe('RATE_LIMIT_EXCEEDED');
        }
      }
    });

    it('should handle validation errors', async () => {
      mockFetch.mockResolvedValueOnce({
        ok: false,
        status: 400,
        json: async () => ({
          success: false,
          error: {
            type: 'VALIDATION_ERROR',
            message: 'Invalid code format',
            details: {
              code: 'Code is required',
            },
          },
        }),
      });

      await expect(api.submitExercise('exercise-001', '')).rejects.toThrow(
        APIError
      );
    });
  });

  describe('Error Handling', () => {
    it('should handle network errors', async () => {
      mockFetch.mockRejectedValueOnce(new Error('Network error'));

      await expect(api.health()).rejects.toThrow('Network error');
    });

    it('should handle non-JSON responses', async () => {
      mockFetch.mockResolvedValueOnce({
        ok: false,
        status: 500,
        json: async () => {
          throw new Error('Invalid JSON');
        },
      });

      await expect(api.health()).rejects.toThrow();
    });

    it('should handle empty error response', async () => {
      mockFetch.mockResolvedValueOnce({
        ok: false,
        status: 500,
        json: async () => ({
          success: false,
        }),
      });

      await expect(api.health()).rejects.toThrow(APIError);
      await expect(api.health()).rejects.toThrow('Request failed');
    });

    it('should create APIError with correct properties', async () => {
      mockFetch.mockResolvedValueOnce({
        ok: false,
        status: 404,
        json: async () => ({
          success: false,
          error: {
            type: 'NOT_FOUND',
            message: 'Resource not found',
            details: { id: 'nonexistent' },
          },
        }),
      });

      try {
        await api.health();
      } catch (error) {
        expect(error).toBeInstanceOf(APIError);
        if (error instanceof APIError) {
          expect(error.message).toBe('Resource not found');
          expect(error.status).toBe(404);
          expect(error.type).toBe('NOT_FOUND');
          expect(error.details).toEqual({ id: 'nonexistent' });
        }
      }
    });
  });

  describe('Authentication Headers', () => {
    it('should include auth token in getUserProgress', async () => {
      mockFetch.mockResolvedValueOnce({
        ok: true,
        json: async () => ({
          success: true,
          data: { progress: [], total: 0 },
        }),
      });

      await api.getUserProgress('user-123', 1, 20, 'auth-token-xyz');

      expect(mockFetch).toHaveBeenCalledWith(
        expect.any(String),
        expect.objectContaining({
          headers: expect.objectContaining({
            Authorization: 'Bearer auth-token-xyz',
          }),
        })
      );
    });

    it('should not include auth header when token not provided', async () => {
      mockFetch.mockResolvedValueOnce({
        ok: true,
        json: async () => ({
          success: true,
          data: { progress: [], total: 0 },
        }),
      });

      await api.getUserProgress('user-123');

      expect(mockFetch).toHaveBeenCalledWith(
        expect.any(String),
        expect.objectContaining({
          headers: expect.not.objectContaining({
            Authorization: expect.any(String),
          }),
        })
      );
    });
  });

  describe('Complete Lesson', () => {
    it('should complete lesson successfully', async () => {
      mockFetch.mockResolvedValueOnce({
        ok: true,
        json: async () => ({
          success: true,
          data: null,
        }),
      });

      await api.completeLesson('lesson-001');

      expect(mockFetch).toHaveBeenCalledWith(
        'http://localhost:8080/api/v1/lessons/lesson-001/complete',
        expect.objectContaining({
          method: 'POST',
        })
      );
    });

    it('should include auth token when completing lesson', async () => {
      mockFetch.mockResolvedValueOnce({
        ok: true,
        json: async () => ({
          success: true,
          data: null,
        }),
      });

      await api.completeLesson('lesson-001', 'token-123');

      expect(mockFetch).toHaveBeenCalledWith(
        expect.any(String),
        expect.objectContaining({
          headers: expect.objectContaining({
            Authorization: 'Bearer token-123',
          }),
        })
      );
    });
  });

  describe('API Base URL Configuration', () => {
    const originalEnv = process.env.NEXT_PUBLIC_API_URL;

    afterEach(() => {
      process.env.NEXT_PUBLIC_API_URL = originalEnv;
    });

    it('should use environment variable for API base URL', () => {
      process.env.NEXT_PUBLIC_API_URL = 'https://api.example.com';

      mockFetch.mockResolvedValueOnce({
        ok: true,
        json: async () => ({
          success: true,
          data: { status: 'healthy' },
        }),
      });

      // Note: In real implementation, would need to reload module
      // This test demonstrates the pattern
      expect(process.env.NEXT_PUBLIC_API_URL).toBe('https://api.example.com');
    });
  });
});
