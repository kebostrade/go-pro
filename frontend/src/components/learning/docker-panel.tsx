'use client';

import React from 'react';
import { useDockerEnvironment } from '@/lib/docker-hooks';
import { Loader2, Play, Square, RefreshCw, AlertCircle, CheckCircle, XCircle } from 'lucide-react';

interface DockerPanelProps {
  /** Topic ID for the Docker environment */
  topicId: string;
  /** Optional className for styling */
  className?: string;
}

export function DockerPanel({ topicId, className = '' }: DockerPanelProps) {
  const { status, loading, error, start, stop, refresh } = useDockerEnvironment(topicId);

  const isRunning = status?.status === 'running';
  const isStopped = status?.status === 'stopped';

  const getStatusIcon = () => {
    if (loading) return <Loader2 className="w-4 h-4 animate-spin" />;
    if (error || status?.status === 'error') return <XCircle className="w-4 h-4 text-red-500" />;
    if (isRunning) return <CheckCircle className="w-4 h-4 text-green-500" />;
    if (isStopped) return <AlertCircle className="w-4 h-4 text-yellow-500" />;
    return <AlertCircle className="w-4 h-4 text-gray-400" />;
  };

  const getStatusText = () => {
    if (loading) return 'Loading...';
    if (error) return 'Error';
    if (status?.status === 'error') return 'Error';
    if (isRunning) return 'Running';
    if (isStopped) return 'Stopped';
    return 'Not Created';
  };

  const getStatusColor = () => {
    if (loading) return 'text-gray-500';
    if (error || status?.status === 'error') return 'text-red-500';
    if (isRunning) return 'text-green-500';
    return 'text-yellow-500';
  };

  return (
    <div className={`border rounded-lg p-4 bg-slate-50 dark:bg-slate-900 ${className}`}>
      <div className="flex items-center justify-between mb-3">
        <div className="flex items-center gap-2">
          <div className={getStatusColor()}>
            {getStatusIcon()}
          </div>
          <span className={`font-medium ${getStatusColor()}`}>
            Docker: {getStatusText()}
          </span>
        </div>

        <button
          onClick={refresh}
          className="p-1.5 rounded-md hover:bg-slate-200 dark:hover:bg-slate-700 transition-colors"
          title="Refresh status"
          disabled={loading}
        >
          <RefreshCw className="w-4 h-4" />
        </button>
      </div>

      {/* Services list when running */}
      {isRunning && status?.services && status.services.length > 0 && (
        <div className="mb-3 pl-6">
          <ul className="text-sm space-y-1">
            {status.services.map((service) => (
              <li key={service.name} className="flex items-center gap-2">
                <span className={`w-2 h-2 rounded-full ${
                  service.health === 'healthy' ? 'bg-green-500' :
                  service.health === 'starting' ? 'bg-yellow-500' : 'bg-red-500'
                }`} />
                <span className="text-slate-700 dark:text-slate-300">
                  {service.name}
                </span>
                <span className="text-slate-400 text-xs">
                  ({service.status})
                </span>
              </li>
            ))}
          </ul>
        </div>
      )}

      {/* Error message */}
      {error && (
        <div className="mb-3 p-2 bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded text-sm text-red-600 dark:text-red-400">
          {error}
        </div>
      )}

      {/* Action buttons */}
      <div className="flex gap-2">
        {!isRunning ? (
          <button
            onClick={start}
            disabled={loading}
            className="flex-1 flex items-center justify-center gap-2 px-3 py-2 bg-green-600 hover:bg-green-700 disabled:bg-green-400 text-white rounded-md transition-colors text-sm font-medium"
          >
            {loading ? (
              <Loader2 className="w-4 h-4 animate-spin" />
            ) : (
              <Play className="w-4 h-4" />
            )}
            Start Environment
          </button>
        ) : (
          <button
            onClick={stop}
            disabled={loading}
            className="flex-1 flex items-center justify-center gap-2 px-3 py-2 bg-red-600 hover:bg-red-700 disabled:bg-red-400 text-white rounded-md transition-colors text-sm font-medium"
          >
            {loading ? (
              <Loader2 className="w-4 h-4 animate-spin" />
            ) : (
              <Square className="w-4 h-4" />
            )}
            Stop Environment
          </button>
        )}
      </div>

      {/* Info about Docker Desktop requirement */}
      {!isRunning && !loading && (
        <p className="mt-2 text-xs text-slate-500 dark:text-slate-400">
          Requires Docker Desktop running locally
        </p>
      )}
    </div>
  );
}

export default DockerPanel;
