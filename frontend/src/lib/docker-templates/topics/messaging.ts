/**
 * Docker Compose template for messaging topics
 * Go service + NATS with JetStream
 */

export interface MessagingTopicConfig {
  topicId: string;
  topicTitle: string;
  natsPort?: number;     // Client port for NATS
  natsMonitorPort?: number; // Monitoring port
  goServicePort?: number; // Host port for Go service
}

/**
 * Generate docker-compose.yml for messaging topics
 * Based on: basic/projects/nats-events/docker-compose.yml
 */
export function generateMessagingCompose(config: MessagingTopicConfig): string {
  const {
    topicId,
    natsPort = 4222,
    natsMonitorPort = 8222,
    goServicePort = 8080,
  } = config;
  
  const serviceName = topicId.toLowerCase().replace(/[^a-z0-9-]/g, '-');
  
  return `version: '3.8'

services:
  nats:
    image: nats:2.10-alpine
    container_name: ${serviceName}-nats
    ports:
      - "${natsPort}:4222"  # Client connections
      - "${natsMonitorPort}:8222"  # Monitoring
    command: ["-js"]  # Enable JetStream
    networks:
      - ${serviceName}-network
    healthcheck:
      test: ["CMD", "wget", "-q", "-O", "-", "http://localhost:8222/healthz"]
      interval: 10s
      timeout: 5s
      retries: 5

  ${serviceName}:
    build: .
    container_name: ${serviceName}
    ports:
      - "${goServicePort}:8080"
    environment:
      - GIN_MODE=release
      - NATS_URL=nats://nats:4222
    depends_on:
      nats:
        condition: service_healthy
    networks:
      - ${serviceName}-network
    restart: unless-stopped
    healthcheck:
      test: ["CMD", "wget", "--no-verbose", "--tries=1", "--spider", "http://localhost:8080/health"]
      interval: 30s
      timeout: 10s
      retries: 3

networks:
  ${serviceName}-network:
    driver: bridge
`;
}

/**
 * Messaging topic configurations
 */
export const MESSAGING_TOPICS: Record<string, { natsPort: number; goServicePort: number }> = {
  'nats-events': { natsPort: 4222, goServicePort: 8080 },
};
