import { ApiResponse, PaginatedResponse } from './types';

// API Configuration
const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api/v1';

// API Client Class
class ApiClient {
  private baseURL: string;
  private defaultHeaders: Record<string, string>;

  constructor(baseURL: string = API_BASE_URL) {
    this.baseURL = baseURL;
    this.defaultHeaders = {
      'Content-Type': 'application/json',
    };
  }

  // Set authentication token
  setAuthToken(token: string) {
    this.defaultHeaders['Authorization'] = `Bearer ${token}`;
  }

  // Remove authentication token
  removeAuthToken() {
    delete this.defaultHeaders['Authorization'];
  }

  // Generic request method
  private async request<T>(
    endpoint: string,
    options: RequestInit = {}
  ): Promise<ApiResponse<T>> {
    const url = `${this.baseURL}${endpoint}`;
    
    const config: RequestInit = {
      ...options,
      headers: {
        ...this.defaultHeaders,
        ...options.headers,
      },
    };

    try {
      const response = await fetch(url, config);
      
      if (!response.ok) {
        const errorData = await response.json().catch(() => ({})) as { message?: string };
        throw new Error(errorData.message || `HTTP ${response.status}: ${response.statusText}`);
      }

      const data = await response.json() as T;
      return {
        success: true,
        data,
      };
    } catch (error) {
      console.error(`API request failed: ${endpoint}`, error);
      return {
        success: false,
        error: error instanceof Error ? error.message : 'Unknown error occurred',
      };
    }
  }

  // GET request
  async get<T>(endpoint: string, params?: Record<string, any>): Promise<ApiResponse<T>> {
    const url = new URL(`${this.baseURL}${endpoint}`);
    
    if (params) {
      Object.entries(params).forEach(([key, value]) => {
        if (value !== undefined && value !== null) {
          if (Array.isArray(value)) {
            value.forEach(v => url.searchParams.append(key, v.toString()));
          } else {
            url.searchParams.append(key, value.toString());
          }
        }
      });
    }

    return this.request<T>(url.pathname + url.search);
  }

  // POST request
  async post<T>(endpoint: string, data?: any): Promise<ApiResponse<T>> {
    return this.request<T>(endpoint, {
      method: 'POST',
      body: data ? JSON.stringify(data) : undefined,
    });
  }

  // PUT request
  async put<T>(endpoint: string, data?: any): Promise<ApiResponse<T>> {
    return this.request<T>(endpoint, {
      method: 'PUT',
      body: data ? JSON.stringify(data) : undefined,
    });
  }

  // PATCH request
  async patch<T>(endpoint: string, data?: any): Promise<ApiResponse<T>> {
    return this.request<T>(endpoint, {
      method: 'PATCH',
      body: data ? JSON.stringify(data) : undefined,
    });
  }

  // DELETE request
  async delete<T>(endpoint: string): Promise<ApiResponse<T>> {
    return this.request<T>(endpoint, {
      method: 'DELETE',
    });
  }

  // Upload file
  async upload<T>(endpoint: string, file: File, additionalData?: Record<string, any>): Promise<ApiResponse<T>> {
    const formData = new FormData();
    formData.append('file', file);
    
    if (additionalData) {
      Object.entries(additionalData).forEach(([key, value]) => {
        formData.append(key, value.toString());
      });
    }

    const headers = { ...this.defaultHeaders };
    delete headers['Content-Type']; // Let browser set content-type for FormData

    return this.request<T>(endpoint, {
      method: 'POST',
      body: formData,
      headers,
    });
  }
}

// Create singleton instance
export const apiClient = new ApiClient();

// Utility functions for common API patterns
export const withPagination = <T>(
  apiCall: (page: number, limit: number) => Promise<ApiResponse<PaginatedResponse<T>>>
) => {
  return async (page: number = 1, limit: number = 20): Promise<ApiResponse<PaginatedResponse<T>>> => {
    return apiCall(page, limit);
  };
};

export const withRetry = <T>(
  apiCall: () => Promise<ApiResponse<T>>,
  maxRetries: number = 3,
  delay: number = 1000
) => {
  return async (): Promise<ApiResponse<T>> => {
    let lastError: Error | null = null;
    
    for (let attempt = 1; attempt <= maxRetries; attempt++) {
      try {
        const result = await apiCall();
        if (result.success) {
          return result;
        }
        lastError = new Error(result.error || 'API call failed');
      } catch (error) {
        lastError = error instanceof Error ? error : new Error('Unknown error');
      }
      
      if (attempt < maxRetries) {
        await new Promise(resolve => setTimeout(resolve, delay * attempt));
      }
    }
    
    return {
      success: false,
      error: lastError?.message || 'Max retries exceeded',
    };
  };
};

// Error handling utilities
export const isApiError = (response: ApiResponse<any>): response is ApiResponse<never> & { success: false } => {
  return !response.success;
};

export const getErrorMessage = (response: ApiResponse<any>): string => {
  if (isApiError(response)) {
    return response.error || 'An unknown error occurred';
  }
  return '';
};

// Response validation utilities
export const validateResponse = <T>(
  response: ApiResponse<T>,
  validator?: (data: T) => boolean
): response is ApiResponse<T> & { success: true; data: T } => {
  if (!response.success || !response.data) {
    return false;
  }
  
  if (validator && !validator(response.data)) {
    return false;
  }
  
  return true;
};

// Cache utilities for API responses
class ApiCache {
  private cache = new Map<string, { data: any; timestamp: number; ttl: number }>();

  set<T>(key: string, data: T, ttlMs: number = 5 * 60 * 1000): void {
    this.cache.set(key, {
      data,
      timestamp: Date.now(),
      ttl: ttlMs,
    });
  }

  get<T>(key: string): T | null {
    const cached = this.cache.get(key);
    
    if (!cached) {
      return null;
    }
    
    if (Date.now() - cached.timestamp > cached.ttl) {
      this.cache.delete(key);
      return null;
    }
    
    return cached.data as T;
  }

  clear(): void {
    this.cache.clear();
  }

  delete(key: string): void {
    this.cache.delete(key);
  }
}

export const apiCache = new ApiCache();

// Cached API call wrapper
export const withCache = <T>(
  apiCall: () => Promise<ApiResponse<T>>,
  cacheKey: string,
  ttlMs: number = 5 * 60 * 1000
) => {
  return async (): Promise<ApiResponse<T>> => {
    // Try to get from cache first
    const cached = apiCache.get<T>(cacheKey);
    if (cached) {
      return {
        success: true,
        data: cached,
      };
    }
    
    // Make API call
    const result = await apiCall();
    
    // Cache successful responses
    if (result.success && result.data) {
      apiCache.set(cacheKey, result.data, ttlMs);
    }
    
    return result;
  };
};

// Request interceptors
type RequestInterceptor = (config: RequestInit) => RequestInit | Promise<RequestInit>;
type ResponseInterceptor<T> = (response: ApiResponse<T>) => ApiResponse<T> | Promise<ApiResponse<T>>;

class InterceptorManager {
  private requestInterceptors: RequestInterceptor[] = [];
  private responseInterceptors: ResponseInterceptor<any>[] = [];

  addRequestInterceptor(interceptor: RequestInterceptor): void {
    this.requestInterceptors.push(interceptor);
  }

  addResponseInterceptor<T>(interceptor: ResponseInterceptor<T>): void {
    this.responseInterceptors.push(interceptor);
  }

  async processRequest(config: RequestInit): Promise<RequestInit> {
    let processedConfig = config;
    
    for (const interceptor of this.requestInterceptors) {
      processedConfig = await interceptor(processedConfig);
    }
    
    return processedConfig;
  }

  async processResponse<T>(response: ApiResponse<T>): Promise<ApiResponse<T>> {
    let processedResponse = response;
    
    for (const interceptor of this.responseInterceptors) {
      processedResponse = await interceptor(processedResponse);
    }
    
    return processedResponse;
  }
}

export const interceptors = new InterceptorManager();

// Add default interceptors
// interceptors.addRequestInterceptor((config) => {
//   // Add timestamp to prevent caching
//   if (config.method === 'GET') {
//     const url = new URL(config.url || '');
//     url.searchParams.set('_t', Date.now().toString());
//     config.url = url.toString();
//   }
//
//   return config;
// });

interceptors.addResponseInterceptor((response) => {
  // Log errors in development
  if (!response.success && process.env.NODE_ENV === 'development') {
    console.error('API Error:', response.error);
  }
  
  return response;
});

export default apiClient;
