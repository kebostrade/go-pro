/**
 * Topic category classification for Docker template generation
 */
export type TopicCategory = 'simple' | 'database' | 'messaging' | 'cloud';

/**
 * Category configuration
 */
export interface CategoryConfig {
  id: TopicCategory;
  label: string;
  description: string;
  /** Topics that require cloud/simulation instead of local Docker */
  isCloudTopic: boolean;
}

/**
 * All topic categories
 */
export const CATEGORIES: Record<TopicCategory, CategoryConfig> = {
  simple: {
    id: 'simple',
    label: 'Simple',
    description: 'Single Go service with optional health check',
    isCloudTopic: false,
  },
  database: {
    id: 'database',
    label: 'Database',
    description: 'Go service with PostgreSQL and/or Redis',
    isCloudTopic: false,
  },
  messaging: {
    id: 'messaging',
    label: 'Messaging',
    description: 'Go service with NATS or Kafka',
    isCloudTopic: false,
  },
  cloud: {
    id: 'cloud',
    label: 'Cloud/Simulation',
    description: 'Topics requiring cloud infrastructure or simulation',
    isCloudTopic: true,
  },
};

/**
 * Map of topic IDs to their categories
 * Based on existing docker-compose.yml files and D-10/D-11 decisions
 */
export const TOPIC_CATEGORIES: Record<string, TopicCategory> = {
  // Simple topics (single Go service)
  'rest-api': 'simple',
  'cli-tools': 'simple',
  'concurrent-patterns': 'simple',
  'error-handling': 'simple',
  'testing-patterns': 'simple',
  'gin-web': 'simple',
  'websocket-chat': 'simple',
  'grpc-services': 'simple',
  'grpc-service': 'simple',
  'message-queues': 'simple',
  'graphql-api': 'simple',
  'graphql': 'simple',
  'blockchain': 'simple',
  'iot-mqtt': 'simple',
  'ml-gorgonia': 'simple',
  'ml': 'simple',
  'data-processing': 'simple',
  'observability': 'simple',
  'security': 'simple',
  
  // Database topics (Go + PostgreSQL + Redis)
  'postgres-redis-go': 'database',
  'microservices': 'database',
  
  // Messaging topics (Go + NATS)
  'nats-events': 'messaging',
  
  // Cloud topics (require simulation or cloud infrastructure)
  'kubernetes': 'cloud',
  'docker-kubernetes': 'cloud',
  'distributed-systems': 'cloud',
  'aws-lambda': 'cloud',
};

/**
 * Get category for a topic
 */
export function getTopicCategory(topicId: string): TopicCategory {
  return TOPIC_CATEGORIES[topicId] ?? 'simple';
}

/**
 * Check if a topic requires cloud/simulation
 */
export function isCloudTopic(topicId: string): boolean {
  const category = getTopicCategory(topicId);
  return CATEGORIES[category]?.isCloudTopic ?? false;
}
