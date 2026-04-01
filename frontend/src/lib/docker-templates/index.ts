/**
 * Docker Template Generator
 * Generates docker-compose.yml files for all 15 Go learning topics
 */

import { getTopicCategory, isCloudTopic, TopicCategory } from './categories';
import { generateSimpleCompose, SIMPLE_TOPIC_PORTS, SimpleTopicConfig } from './topics/simple';
import { generateDatabaseCompose, DATABASE_TOPICS, DatabaseTopicConfig } from './topics/database';
import { generateMessagingCompose, MESSAGING_TOPICS, MessagingTopicConfig } from './topics/messaging';
import { generateCloudCompose, CLOUD_TOPICS, CloudTopicConfig } from './topics/cloud';

// Re-export everything for convenience
export type { TopicCategory } from './categories';
export { getTopicCategory, isCloudTopic, CATEGORIES } from './categories';
export type { CLOUD_TOPICS } from './topics/cloud';
export { SIMPLE_TOPIC_PORTS } from './topics/simple';
export { DATABASE_TOPICS } from './topics/database';
export { MESSAGING_TOPICS } from './topics/messaging';

/**
 * Topic metadata for template generation
 */
export interface TopicMetadata {
  id: string;
  title: string;
}

/**
 * Generate docker-compose.yml content for a topic
 * 
 * @param topic - Topic ID or full topic metadata
 * @returns The generated docker-compose.yml content, or null for cloud topics
 */
export function generateCompose(topic: string | TopicMetadata): string | null {
  const topicId = typeof topic === 'string' ? topic : topic.id;
  const topicTitle = typeof topic === 'string' ? topicId : topic.title;
  
  // Handle cloud topics (they don't actually run docker compose)
  if (isCloudTopic(topicId)) {
    return generateCloudCompose({ topicId, topicTitle });
  }
  
  const category = getTopicCategory(topicId);
  
  switch (category) {
    case 'simple': {
      const port = SIMPLE_TOPIC_PORTS[topicId] ?? 8080;
      return generateSimpleCompose({ topicId, topicTitle, port });
    }
    
    case 'database': {
      const dbConfig = DATABASE_TOPICS[topicId] || {
        postgresPort: 5432,
        redisPort: 6379,
        useRedis: true,
      };
      return generateDatabaseCompose({
        topicId,
        topicTitle,
        ...dbConfig,
      });
    }
    
    case 'messaging': {
      const msgConfig = MESSAGING_TOPICS[topicId] || {
        natsPort: 4222,
        goServicePort: 8080,
      };
      return generateMessagingCompose({
        topicId,
        topicTitle,
        ...msgConfig,
      });
    }
    
    default:
      // Unknown category, fall back to simple
      return generateSimpleCompose({ topicId, topicTitle });
  }
}

/**
 * Get topic metadata for template generation
 * This can be used to get the appropriate compose content for any topic
 */
export function getTopicCompose(topicId: string): {
  content: string | null;
  category: TopicCategory;
  isCloud: boolean;
} {
  const category = getTopicCategory(topicId);
  const isCloud = isCloudTopic(topicId);
  const content = isCloud ? null : generateCompose(topicId);
  
  return { content, category, isCloud };
}

/**
 * All supported topic IDs
 */
export const ALL_TOPICS = [
  // Simple topics
  'rest-api',
  'cli-tools',
  'concurrent-patterns',
  'error-handling',
  'testing-patterns',
  'gin-web',
  'websocket-chat',
  'grpc-services',
  'message-queues',
  'graphql-api',
  'blockchain',
  'iot-mqtt',
  'ml-gorgonia',
  'data-processing',
  'observability',
  'security',
  // Database topics
  'postgres-redis-go',
  'microservices',
  // Messaging topics
  'nats-events',
  // Cloud topics
  'kubernetes',
  'docker-kubernetes',
  'distributed-systems',
  'aws-lambda',
] as const;

export type TopicId = typeof ALL_TOPICS[number];
