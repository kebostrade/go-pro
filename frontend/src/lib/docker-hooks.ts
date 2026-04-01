import { useState, useEffect, useCallback, useRef } from 'react';
import { dockerApi, DockerStatus } from './docker-api';

interface UseDockerEnvironmentOptions {
  /** Poll interval in milliseconds when status is running */
  pollInterval?: number;
  /** Enable automatic polling when running */
  autoPoll?: boolean;
}

interface UseDockerEnvironmentReturn {
  /** Current Docker status */
  status: DockerStatus | null;
  /** Loading state for start/stop operations */
  loading: boolean;
  /** Error message if last operation failed */
  error: string | null;
  /** Start the Docker environment */
  start: () => Promise<void>;
  /** Stop the Docker environment */
  stop: () => Promise<void>;
  /** Manually refresh status */
  refresh: () => Promise<void>;
}

export function useDockerEnvironment(
  topicId: string | null,
  options: UseDockerEnvironmentOptions = {}
): UseDockerEnvironmentReturn {
  const { pollInterval = 5000, autoPoll = true } = options;

  const [status, setStatus] = useState<DockerStatus | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const pollTimerRef = useRef<NodeJS.Timeout | null>(null);
  const topicIdRef = useRef(topicId);

  // Keep topicIdRef updated
  useEffect(() => {
    topicIdRef.current = topicId;
  }, [topicId]);

  // Clear polling on unmount
  useEffect(() => {
    return () => {
      if (pollTimerRef.current) {
        clearInterval(pollTimerRef.current);
      }
    };
  }, []);

  // Start polling when status is running
  useEffect(() => {
    if (!autoPoll || status?.status !== 'running') {
      // Stop polling if not running
      if (pollTimerRef.current) {
        clearInterval(pollTimerRef.current);
        pollTimerRef.current = null;
      }
      return;
    }

    // Start polling
    pollTimerRef.current = setInterval(async () => {
      if (!topicIdRef.current) return;
      try {
        const currentStatus = await dockerApi.getStatus(topicIdRef.current);
        setStatus(currentStatus);

        // Stop polling if no longer running
        if (currentStatus.status !== 'running') {
          if (pollTimerRef.current) {
            clearInterval(pollTimerRef.current);
            pollTimerRef.current = null;
          }
        }
      } catch {
        // Silently ignore polling errors - status will be stale
      }
    }, pollInterval);

    return () => {
      if (pollTimerRef.current) {
        clearInterval(pollTimerRef.current);
        pollTimerRef.current = null;
      }
    };
  }, [status?.status, autoPoll, pollInterval]);

  const start = useCallback(async () => {
    if (!topicId) return;

    setLoading(true);
    setError(null);

    try {
      const result = await dockerApi.start(topicId);
      setStatus(result);
    } catch (err) {
      const message = err instanceof Error ? err.message : 'Failed to start';
      setError(message);
    } finally {
      setLoading(false);
    }
  }, [topicId]);

  const stop = useCallback(async () => {
    if (!topicId) return;

    setLoading(true);
    setError(null);

    try {
      const result = await dockerApi.stop(topicId);
      setStatus(result);
    } catch (err) {
      const message = err instanceof Error ? err.message : 'Failed to stop';
      setError(message);
    } finally {
      setLoading(false);
    }
  }, [topicId]);

  const refresh = useCallback(async () => {
    if (!topicId) return;

    try {
      const result = await dockerApi.getStatus(topicId);
      setStatus(result);
      setError(null);
    } catch (err) {
      const message = err instanceof Error ? err.message : 'Failed to refresh';
      setError(message);
    }
  }, [topicId]);

  return {
    status,
    loading,
    error,
    start,
    stop,
    refresh,
  };
}
