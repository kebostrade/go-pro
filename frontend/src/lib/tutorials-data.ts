// Complete Tutorial Data for Go Learning Platform
export interface Tutorial {
  id: string;
  number: number;
  title: string;
  description: string;
  category: 'basics' | 'intermediate' | 'advanced' | 'projects';
  difficulty: 'beginner' | 'intermediate' | 'advanced';
  duration: string;
  topics: string[];
  projectPath: string;
  icon: string;
  color: string;
  prerequisites: string[];
  learningOutcomes: string[];
  featured: boolean;
}

export const tutorials: Tutorial[] = [
  {
    id: 'websocket-realtime',
    number: 12,
    title: 'WebSocket Real-Time Communication',
    description: 'Build real-time chat applications with WebSocket, implementing the Hub pattern for concurrent client management.',
    category: 'intermediate',
    difficulty: 'intermediate',
    duration: '4 hours',
    topics: ['WebSocket', 'Concurrency', 'Real-time', 'Hub Pattern'],
    projectPath: 'basic/projects/websocket-chat',
    icon: '💬',
    color: 'from-blue-500 to-cyan-500',
    prerequisites: ['Go Basics', 'HTTP', 'Goroutines'],
    learningOutcomes: [
      'WebSocket protocol implementation',
      'Hub pattern for message broadcasting',
      'Concurrent client management',
      'Real-time bidirectional communication'
    ],
    featured: true
  },
  {
    id: 'microservices',
    number: 13,
    title: 'Microservices Architecture',
    description: 'Design and implement microservices with gRPC, service discovery, API gateway, and distributed systems patterns.',
    category: 'advanced',
    difficulty: 'advanced',
    duration: '8 hours',
    topics: ['Microservices', 'gRPC', 'Service Discovery', 'API Gateway'],
    projectPath: 'basic/projects/microservices-demo',
    icon: '🏗️',
    color: 'from-purple-500 to-pink-500',
    prerequisites: ['HTTP', 'REST APIs', 'Docker'],
    learningOutcomes: [
      'Microservices architecture patterns',
      'gRPC service communication',
      'Service discovery and registration',
      'API gateway implementation'
    ],
    featured: true
  },
  {
    id: 'design-patterns',
    number: 14,
    title: 'Design Patterns in Go',
    description: 'Master essential design patterns including Singleton, Factory, Observer, Strategy, and more with Go implementations.',
    category: 'intermediate',
    difficulty: 'intermediate',
    duration: '5 hours',
    topics: ['Design Patterns', 'OOP', 'Architecture', 'Best Practices'],
    projectPath: 'basic/projects/design-patterns',
    icon: '🎨',
    color: 'from-green-500 to-teal-500',
    prerequisites: ['Go Basics', 'Interfaces', 'Structs'],
    learningOutcomes: [
      'Creational patterns (Singleton, Factory, Builder)',
      'Structural patterns (Adapter, Decorator, Facade)',
      'Behavioral patterns (Observer, Strategy, Command)',
      'Go-specific pattern implementations'
    ],
    featured: false
  },
  {
    id: 'concurrency-patterns',
    number: 15,
    title: 'Concurrency Patterns in Go',
    description: 'Deep dive into Go concurrency with 12 essential patterns including Worker Pool, Pipeline, Fan-Out/Fan-In, and more.',
    category: 'advanced',
    difficulty: 'advanced',
    duration: '6 hours',
    topics: ['Concurrency', 'Goroutines', 'Channels', 'Sync Primitives'],
    projectPath: 'basic/projects/concurrency-patterns',
    icon: '⚡',
    color: 'from-yellow-500 to-orange-500',
    prerequisites: ['Goroutines', 'Channels', 'Select Statement'],
    learningOutcomes: [
      'Worker Pool pattern',
      'Pipeline and Fan-Out/Fan-In',
      'Context cancellation',
      'Rate limiting and circuit breakers'
    ],
    featured: true
  },
  {
    id: 'web-architecture',
    number: 16,
    title: 'Web Architecture with Go',
    description: 'Build production-ready web applications with Clean Architecture, repository pattern, JWT auth, and middleware.',
    category: 'advanced',
    difficulty: 'advanced',
    duration: '7 hours',
    topics: ['Clean Architecture', 'REST API', 'JWT', 'Middleware'],
    projectPath: 'basic/projects/web-architecture',
    icon: '🌐',
    color: 'from-indigo-500 to-purple-500',
    prerequisites: ['HTTP', 'REST APIs', 'Database'],
    learningOutcomes: [
      'Clean Architecture principles',
      'Repository pattern implementation',
      'JWT authentication and authorization',
      'Middleware and request handling'
    ],
    featured: true
  },
  {
    id: 'devops-docker-k8s',
    number: 17,
    title: 'DevOps with Go - Docker, Kubernetes, Terraform',
    description: 'Deploy Go applications using Docker multi-stage builds, Kubernetes orchestration, and Terraform infrastructure as code.',
    category: 'advanced',
    difficulty: 'advanced',
    duration: '8 hours',
    topics: ['Docker', 'Kubernetes', 'Terraform', 'DevOps', 'CI/CD'],
    projectPath: 'basic/projects/devops-with-go',
    icon: '🚀',
    color: 'from-blue-600 to-indigo-600',
    prerequisites: ['Docker Basics', 'Linux', 'Cloud Concepts'],
    learningOutcomes: [
      'Multi-stage Docker builds',
      'Kubernetes deployments and services',
      'Terraform infrastructure provisioning',
      'Prometheus and Grafana monitoring'
    ],
    featured: true
  },
  {
    id: 'messaging-kafka-rabbitmq',
    number: 18,
    title: 'Messaging with Go - Kafka and RabbitMQ',
    description: 'Build scalable messaging systems with Apache Kafka for event streaming and RabbitMQ for message queuing.',
    category: 'advanced',
    difficulty: 'advanced',
    duration: '6 hours',
    topics: ['Kafka', 'RabbitMQ', 'Event Streaming', 'Message Queues'],
    projectPath: 'basic/projects/messaging-with-go',
    icon: '📨',
    color: 'from-red-500 to-pink-500',
    prerequisites: ['Distributed Systems', 'Concurrency'],
    learningOutcomes: [
      'Kafka producers and consumers',
      'Consumer groups and partitions',
      'RabbitMQ exchanges and queues',
      'Messaging patterns (Pub/Sub, Work Queue)'
    ],
    featured: true
  },
  {
    id: 'ethical-hacking',
    number: 19,
    title: 'Ethical Hacking with Go',
    description: 'Learn network security and penetration testing by building security tools: port scanner, packet sniffer, vulnerability scanner.',
    category: 'advanced',
    difficulty: 'advanced',
    duration: '7 hours',
    topics: ['Security', 'Networking', 'Penetration Testing', 'Cryptography'],
    projectPath: 'basic/projects/ethical-hacking',
    icon: '🔒',
    color: 'from-gray-700 to-gray-900',
    prerequisites: ['Networking', 'TCP/IP', 'HTTP'],
    learningOutcomes: [
      'Port scanning and network mapping',
      'Packet capture and analysis',
      'Web vulnerability scanning',
      'Password security and cryptography'
    ],
    featured: true
  },
  {
    id: 'postgres-redis',
    number: 20,
    title: 'PostgreSQL & Redis with Go',
    description: 'Master database operations with PostgreSQL (pgx, GORM) and Redis caching, including connection pooling, transactions, pub/sub, and production-ready patterns.',
    category: 'advanced',
    difficulty: 'advanced',
    duration: '8 hours',
    topics: ['PostgreSQL', 'Redis', 'pgx', 'GORM', 'Caching', 'Pub/Sub', 'Connection Pooling', 'Transactions'],
    projectPath: 'basic/projects/postgres-redis-go',
    icon: '🐘',
    color: 'from-blue-600 to-red-500',
    prerequisites: ['Database basics', 'SQL knowledge', 'Caching concepts'],
    learningOutcomes: [
      'PostgreSQL with pgx and GORM',
      'Connection pooling and transactions',
      'Redis caching and session management',
      'Pub/sub messaging systems',
      'Distributed locks and rate limiting',
      'Cache-aside and write-through patterns'
    ],
    featured: true
  },
  {
    id: 'api-technologies',
    number: 21,
    title: 'RESTful APIs, gRPC & GraphQL',
    description: 'Master modern API technologies with Go - build RESTful APIs with multiple frameworks, high-performance gRPC services, and flexible GraphQL servers.',
    category: 'advanced',
    difficulty: 'intermediate',
    duration: '8 hours',
    topics: ['REST', 'gRPC', 'GraphQL', 'Protocol Buffers', 'Chi', 'Gin', 'gqlgen', 'API Gateway', 'Middleware'],
    projectPath: 'basic/projects/api-technologies-go',
    icon: '🌐',
    color: 'from-indigo-500 to-purple-600',
    prerequisites: ['HTTP basics', 'JSON', 'API concepts'],
    learningOutcomes: [
      'Build REST APIs with net/http, Chi, and Gin',
      'Implement gRPC with Protocol Buffers',
      'Create GraphQL servers with schema-first approach',
      'Implement all RPC streaming patterns',
      'Build unified API gateways',
      'Choose the right API technology for your use case'
    ],
    featured: true
  },
  {
    id: 'cloud-cicd',
    number: 22,
    title: 'Cloud Platforms & CI/CD',
    description: 'Deploy Go applications to Google Cloud Platform and AWS with production-grade CI/CD pipelines. Master serverless, containers, messaging, and infrastructure as code.',
    category: 'advanced',
    difficulty: 'intermediate',
    duration: '10 hours',
    topics: ['GCP', 'AWS', 'Cloud Run', 'Lambda', 'S3', 'Cloud Storage', 'Pub/Sub', 'SQS', 'Terraform', 'GitHub Actions', 'Docker', 'Kubernetes'],
    projectPath: 'basic/projects/cloud-cicd-go',
    icon: '☁️',
    color: 'from-blue-500 to-cyan-600',
    prerequisites: ['Docker basics', 'REST APIs', 'Cloud concepts'],
    learningOutcomes: [
      'Deploy serverless apps to Cloud Run and Lambda',
      'Use cloud storage (GCS, S3) and messaging (Pub/Sub, SQS)',
      'Work with NoSQL databases (Firestore, DynamoDB)',
      'Build CI/CD pipelines with GitHub Actions',
      'Manage infrastructure with Terraform',
      'Containerize applications with Docker',
      'Deploy to Kubernetes (GKE)',
      'Implement cloud best practices'
    ],
    featured: true
  }
];

export const tutorialCategories = {
  basics: {
    name: 'Basics',
    description: 'Fundamental Go concepts and syntax',
    icon: '📚',
    color: 'blue'
  },
  intermediate: {
    name: 'Intermediate',
    description: 'Advanced features and patterns',
    icon: '🎯',
    color: 'purple'
  },
  advanced: {
    name: 'Advanced',
    description: 'Production-ready applications',
    icon: '🚀',
    color: 'red'
  },
  projects: {
    name: 'Projects',
    description: 'Real-world applications',
    icon: '💼',
    color: 'green'
  }
};

export const getTutorialById = (id: string): Tutorial | undefined => {
  return tutorials.find(t => t.id === id);
};

export const getTutorialsByCategory = (category: string): Tutorial[] => {
  return tutorials.filter(t => t.category === category);
};

export const getFeaturedTutorials = (): Tutorial[] => {
  return tutorials.filter(t => t.featured);
};

export const getTutorialStats = () => {
  return {
    total: tutorials.length,
    byDifficulty: {
      beginner: tutorials.filter(t => t.difficulty === 'beginner').length,
      intermediate: tutorials.filter(t => t.difficulty === 'intermediate').length,
      advanced: tutorials.filter(t => t.difficulty === 'advanced').length,
    },
    byCategory: {
      basics: tutorials.filter(t => t.category === 'basics').length,
      intermediate: tutorials.filter(t => t.category === 'intermediate').length,
      advanced: tutorials.filter(t => t.category === 'advanced').length,
      projects: tutorials.filter(t => t.category === 'projects').length,
    }
  };
};

