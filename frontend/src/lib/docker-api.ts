// API client for Docker environment management
// Note: This module has a simplified auth-less implementation for Docker endpoints
// which are expected to be local operations not requiring Firebase auth.

const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || '';

interface DockerApiResponse<T> {
  success: boolean;
  data?: T;
  error?: string;
}

export interface ServiceStatus {
  name: string;
  status: string;
  health: 'healthy' | 'unhealthy' | 'starting' | 'stopped';
}

export interface DockerStatus {
  topic_id: string;
  status: 'running' | 'stopped' | 'not_created' | 'error';
  services: ServiceStatus[];
  ports: Record<string, string>;
  error?: string;
  last_update: string;
}

async function dockerRequest<T>(
  endpoint: string,
  options: RequestInit = {}
): Promise<T> {
  if (!API_BASE_URL) {
    throw new Error('Backend not configured');
  }

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

  const response = await fetch(url, config);
  const data: DockerApiResponse<T> = await response.json();

  if (!response.ok || !data.success) {
    throw new Error(data.error || `Request failed: ${response.status}`);
  }

  return data as T;
}

export const dockerApi = {
  /**
   * Start Docker environment for a topic
   * POST /api/docker/up
   */
  start: async (topicId: string): Promise<DockerStatus> => {
    return dockerRequest<DockerStatus>('/api/docker/up', {
      method: 'POST',
      body: JSON.stringify({ topic_id: topicId }),
    });
  },

  /**
   * Stop Docker environment for a topic
   * POST /api/docker/down
   */
  stop: async (topicId: string): Promise<DockerStatus> => {
    return dockerRequest<DockerStatus>('/api/docker/down', {
      method: 'POST',
      body: JSON.stringify({ topic_id: topicId }),
    });
  },

  /**
   * Get Docker environment status for a topic
   * GET /api/docker/status?topic_id=xxx
   */
  getStatus: async (topicId: string): Promise<DockerStatus> => {
    return dockerRequest<DockerStatus>(`/api/docker/status?topic_id=${encodeURIComponent(topicId)}`);
  },
};
