/**
 * Docker Compose template for simple Go topics
 * Single service with optional port exposure
 */

export interface SimpleTopicConfig {
  topicId: string;
  topicTitle: string;
  port?: number;  // Host port, defaults to topic-specific default
}

/**
 * Generate docker-compose.yml for simple Go topics
 * Based on: basic/projects/rest-api/docker-compose.yml
 */
export function generateSimpleCompose(config: SimpleTopicConfig): string {
  const { topicId, topicTitle, port = 8080 } = config;
  
  // Convert topic ID to service name (e.g., rest-api -> rest-api)
  const serviceName = topicId.toLowerCase().replace(/[^a-z0-9-]/g, '-');
  
  return `version: '3.8'

services:
  ${serviceName}:
    build: .
    container_name: ${serviceName}
    ports:
      - "${port}:8080"
    environment:
      - GIN_MODE=release
    restart: unless-stopped
    healthcheck:
      test: ["CMD", "wget", "-q", "--spider", "http://localhost:8080/health"]
      interval: 30s
      timeout: 10s
      retries: 3
    networks:
      - ${serviceName}-network

networks:
  ${serviceName}-network:
    driver: bridge
`;
}

/**
 * Default ports for simple topics
 * Avoids conflicts with other services
 */
export const SIMPLE_TOPIC_PORTS: Record<string, number> = {
  'rest-api': 8080,
  'gin-web': 8081,
  'grpc-service': 8082,
  'grpc-services': 8082,
  'graphql': 8083,
  'graphql-api': 8083,
  'websocket-chat': 8084,
  'cli-tools': 8085,
  'testing-patterns': 8086,
  'concurrent-patterns': 8087,
  'error-handling': 8088,
  'blockchain': 8089,
  'iot-mqtt': 8090,
  'ml': 8091,
  'ml-gorgonia': 8091,
  'data-processing': 8092,
  'observability': 8093,
  'security': 8094,
  'message-queues': 8095,
};
