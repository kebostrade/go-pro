/**
 * Docker Compose template for cloud topics
 * These topics require simulation or cloud infrastructure
 */

export interface CloudTopicConfig {
  topicId: string;
  topicTitle: string;
}

/**
 * Topics that require cloud infrastructure or simulation
 */
export const CLOUD_TOPICS = [
  'kubernetes',
  'docker-kubernetes',
  'distributed-systems',
  'aws-lambda',
] as const;

export type CloudTopicId = typeof CLOUD_TOPICS[number];

/**
 * Cloud topic fallback configurations
 */
export const CLOUD_TOPIC_CONFIG: Record<CloudTopicId, {
  description: string;
  simulationUrl?: string;
  localCommand?: string;
  localCommandArgs?: string[];
}> = {
  'kubernetes': {
    description: 'Kubernetes deployment simulation',
    simulationUrl: '/topics/kubernetes/simulation',
    localCommand: 'minikube',
    localCommandArgs: ['start', '--profile', 'go-pro'],
  },
  'docker-kubernetes': {
    description: 'Kubernetes with Docker Desktop',
    simulationUrl: '/topics/docker-kubernetes/simulation',
    localCommand: 'minikube',
    localCommandArgs: ['start', '--profile', 'go-pro-docker'],
  },
  'distributed-systems': {
    description: 'Distributed systems simulation (Raft, consensus)',
    simulationUrl: '/topics/distributed-systems/simulation',
  },
  'aws-lambda': {
    description: 'AWS Lambda with LocalStack simulation',
    simulationUrl: '/topics/aws-lambda/simulation',
    localCommand: 'localstack',
    localCommandArgs: ['start'],
  },
};

/**
 * Generate informational compose for cloud topics
 * This is not actually run - it provides documentation
 */
export function generateCloudCompose(config: CloudTopicConfig): string {
  const { topicId, topicTitle } = config;
  const cloudConfig = CLOUD_TOPIC_CONFIG[topicId as CloudTopicId];
  
  if (!cloudConfig) {
    return `# Cloud simulation not available for ${topicId}
# This topic requires manual setup or cloud infrastructure
`;
  }
  
  let compose = `# Cloud Topic: ${topicTitle}
# This topic requires special infrastructure

version: '3.8'

# This is a placeholder compose file for reference only.
# Actual deployment requires:
`;
  
  if (cloudConfig.localCommand) {
    compose += `# - Local command: ${cloudConfig.localCommand} ${(cloudConfig.localCommandArgs || []).join(' ')}
`;
  }
  
  if (cloudConfig.simulationUrl) {
    compose += `# - Simulation available at: ${cloudConfig.simulationUrl}
`;
  }
  
  compose += `#
# For more information, see the topic documentation.
`;
  
  return compose;
}

/**
 * Check if minikube is available (for kubernetes topics)
 */
export async function checkMinikubeAvailable(): Promise<boolean> {
  try {
    const response = await fetch('/api/system/check-command', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ command: 'minikube', args: ['version'] }),
    });
    return response.ok;
  } catch {
    return false;
  }
}

/**
 * Get cloud topic status (simulation vs. real infrastructure)
 */
export async function getCloudTopicStatus(topicId: string): Promise<{
  available: 'local' | 'simulation' | 'unavailable';
  message: string;
}> {
  if (topicId === 'kubernetes' || topicId === 'docker-kubernetes') {
    const hasMinikube = await checkMinikubeAvailable();
    if (hasMinikube) {
      return {
        available: 'local',
        message: 'Minikube is available. You can start a local Kubernetes cluster.',
      };
    }
    return {
      available: 'simulation',
      message: 'Minikube not found. Use the simulation mode instead.',
    };
  }
  
  return {
    available: 'simulation',
    message: CLOUD_TOPIC_CONFIG[topicId as CloudTopicId]?.description || 
             'Cloud simulation not available for this topic.',
  };
}
