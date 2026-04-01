/**
 * Docker Compose template for database topics
 * Go service + PostgreSQL + optional Redis
 */

export interface DatabaseTopicConfig {
  topicId: string;
  topicTitle: string;
  postgresPort?: number;  // Host port for PostgreSQL
  redisPort?: number;     // Host port for Redis
  goServicePort?: number; // Host port for Go service
  useRedis?: boolean;     // Whether to include Redis
}

/**
 * Generate docker-compose.yml for database topics
 * Based on: basic/projects/microservices/docker-compose.yml
 */
export function generateDatabaseCompose(config: DatabaseTopicConfig): string {
  const {
    topicId,
    postgresPort = 5432,
    redisPort = 6379,
    goServicePort = 8080,
    useRedis = true,
  } = config;
  
  const serviceName = topicId.toLowerCase().replace(/[^a-z0-9-]/g, '-');
  
  let servicesYaml = `  postgres:
    image: postgres:15-alpine
    container_name: ${serviceName}-postgres
    environment:
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: password
      POSTGRES_DB: ${serviceName}
    ports:
      - "${postgresPort}:5432"
    volumes:
      - ${serviceName}-postgres-data:/var/lib/postgresql/data
    networks:
      - ${serviceName}-network
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U postgres"]
      interval: 10s
      timeout: 5s
      retries: 5`;

  if (useRedis) {
    servicesYaml += `

  redis:
    image: redis:7-alpine
    container_name: ${serviceName}-redis
    ports:
      - "${redisPort}:6379"
    volumes:
      - ${serviceName}-redis-data:/data
    networks:
      - ${serviceName}-network
    healthcheck:
      test: ["CMD", "redis-cli", "ping"]
      interval: 10s
      timeout: 5s
      retries: 5`;
  }

  servicesYaml += `

  ${serviceName}:
    build: .
    container_name: ${serviceName}
    ports:
      - "${goServicePort}:8080"
    environment:
      - GIN_MODE=release
      - DB_HOST=postgres
      - DB_PORT=5432
      - DB_USER=postgres
      - DB_PASSWORD=password
      - DB_NAME=${serviceName}`;
  
  if (useRedis) {
    servicesYaml += `
      - REDIS_HOST=redis
      - REDIS_PORT=6379`;
  }

  servicesYaml += `
    depends_on:
      postgres:
        condition: service_healthy`;
  
  if (useRedis) {
    servicesYaml += `
      redis:
        condition: service_healthy`;
  }

  servicesYaml += `
    networks:
      - ${serviceName}-network
    restart: unless-stopped
    healthcheck:
      test: ["CMD", "wget", "--no-verbose", "--tries=1", "--spider", "http://localhost:8080/health"]
      interval: 30s
      timeout: 10s
      retries: 3
      start_period: 40s`;

  return `version: '3.8'

services:
${servicesYaml}

networks:
  ${serviceName}-network:
    driver: bridge

volumes:
  ${serviceName}-postgres-data:
`;
}

/**
 * Database topic configurations
 */
export const DATABASE_TOPICS: Record<string, { postgresPort: number; redisPort: number; useRedis: boolean }> = {
  'postgres-redis-go': { postgresPort: 5432, redisPort: 6379, useRedis: true },
  'microservices': { postgresPort: 5433, redisPort: 6379, useRedis: true },
};
