'use client';

import { useState } from 'react';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { Progress } from '@/components/ui/progress';
import {
  Accordion,
  AccordionContent,
  AccordionItem,
  AccordionTrigger,
} from '@/components/ui/accordion';
import {
  Server,
  Database,
  Cloud,
  Zap,
  Shield,
  Network,
  Layers,
  HardDrive,
  ArrowRightLeft,
  GitBranch,
  Globe,
  Lock,
  Activity,
  Cpu,
  MemoryStick,
  Timer,
  Scale,
  RefreshCw,
  AlertTriangle,
  CheckCircle,
  ChevronRight,
  BookOpen,
  Play,
  ExternalLink,
} from 'lucide-react';
import Link from 'next/link';

interface SystemDesignTopic {
  id: string;
  title: string;
  description: string;
  icon: React.ReactNode;
  category: 'fundamentals' | 'data' | 'communication' | 'patterns' | 'scalability';
  difficulty: 'beginner' | 'intermediate' | 'advanced';
  keyPoints: string[];
  pros: string[];
  cons: string[];
  realWorld: string[];
  questions: string[];
}

const systemDesignTopics: SystemDesignTopic[] = [
  // Fundamentals
  {
    id: 'scalability',
    title: 'Scalability Fundamentals',
    description: 'Horizontal vs Vertical scaling, scaling strategies, and trade-offs',
    icon: <Scale className="h-6 w-6" />,
    category: 'fundamentals',
    difficulty: 'beginner',
    keyPoints: [
      'Vertical Scaling (Scale Up): Add more resources to existing server',
      'Horizontal Scaling (Scale Out): Add more servers to the pool',
      'Auto-scaling: Automatically adjust resources based on demand',
      'Elasticity: Ability to scale up and down dynamically',
    ],
    pros: ['Handles increased load', 'Improved performance', 'Better resource utilization'],
    cons: ['Increased complexity', 'Higher costs', 'Potential consistency issues'],
    realWorld: ['Netflix auto-scaling during peak hours', 'AWS EC2 Auto Scaling', 'Kubernetes HPA'],
    questions: ['When would you choose vertical over horizontal scaling?', 'How do you handle state in horizontally scaled systems?'],
  },
  {
    id: 'load-balancing',
    title: 'Load Balancing',
    description: 'Distribute traffic across multiple servers for reliability and performance',
    icon: <Network className="h-6 w-6" />,
    category: 'fundamentals',
    difficulty: 'beginner',
    keyPoints: [
      'Round Robin: Distribute requests sequentially',
      'Least Connections: Route to server with fewest active connections',
      'IP Hash: Route based on client IP for session persistence',
      'Weighted Round Robin: Distribute based on server capacity',
      'Health Checks: Remove unhealthy servers from rotation',
    ],
    pros: ['High availability', 'Improved performance', 'Fault tolerance', 'Flexible scaling'],
    cons: ['Single point of failure (if not redundant)', 'Added latency', 'Cost'],
    realWorld: ['NGINX', 'HAProxy', 'AWS ALB/NLB', 'Cloudflare Load Balancing'],
    questions: ['How does a load balancer detect a failed server?', 'Compare L4 vs L7 load balancers'],
  },
  {
    id: 'availability',
    title: 'High Availability & Fault Tolerance',
    description: 'Design systems that remain operational despite failures',
    icon: <Activity className="h-6 w-6" />,
    category: 'fundamentals',
    difficulty: 'intermediate',
    keyPoints: [
      'Availability = Uptime / (Uptime + Downtime)',
      '99.9% (three 9s) = 8.76 hours downtime/year',
      '99.99% (four 9s) = 52.6 minutes downtime/year',
      'Redundancy: Multiple copies of critical components',
      'Failover: Automatic switching to standby system',
      'Circuit Breaker: Prevent cascade failures',
    ],
    pros: ['Business continuity', 'User trust', 'Reduced data loss'],
    cons: ['Higher infrastructure cost', 'Increased complexity', 'More maintenance'],
    realWorld: ['Netflix Chaos Monkey', 'AWS Multi-AZ deployments', 'Google\'s global load balancing'],
    questions: ['How do you calculate availability of a system with dependencies?', 'What\'s the difference between MTBF and MTTR?'],
  },

  // Data Layer
  {
    id: 'database-sharding',
    title: 'Database Sharding',
    description: 'Horizontal partitioning of data across multiple database instances',
    icon: <Database className="h-6 w-6" />,
    category: 'data',
    difficulty: 'advanced',
    keyPoints: [
      'Horizontal Partitioning: Split rows across servers',
      'Shard Key: Column used to determine shard placement',
      'Range-based Sharding: Partition by value ranges',
      'Hash-based Sharding: Partition by hash of shard key',
      'Directory-based Sharding: Lookup table for shard mapping',
    ],
    pros: ['Scales write operations', 'Improved query performance', 'Geographic distribution'],
    cons: ['Complex cross-shard queries', 'Rebalancing difficulty', 'Potential hotspots'],
    realWorld: ['Instagram sharding with PostgreSQL', 'Uber\'s schemaless storage', 'Discord\'s Cassandra clusters'],
    questions: ['How do you choose a good shard key?', 'How do you handle resharding with minimal downtime?'],
  },
  {
    id: 'caching',
    title: 'Caching Strategies',
    description: 'Store frequently accessed data for faster retrieval',
    icon: <Zap className="h-6 w-6" />,
    category: 'data',
    difficulty: 'intermediate',
    keyPoints: [
      'Cache-Aside: Application manages cache explicitly',
      'Write-Through: Write to cache and DB simultaneously',
      'Write-Back: Write to cache, async to DB',
      'Write-Around: Write directly to DB, cache on read miss',
      'TTL (Time To Live): Automatic cache invalidation',
      'LRU/LFU: Cache eviction policies',
    ],
    pros: ['Reduced latency', 'Lower database load', 'Better user experience'],
    cons: ['Stale data risk', 'Memory costs', 'Cache invalidation complexity'],
    realWorld: ['Redis at Twitter', 'Memcached at Facebook', 'CDN edge caching'],
    questions: ['How do you handle cache stampede?', 'When would you use write-through vs write-back?'],
  },
  {
    id: 'replication',
    title: 'Database Replication',
    description: 'Copy data across multiple database servers for redundancy and performance',
    icon: <RefreshCw className="h-6 w-6" />,
    category: 'data',
    difficulty: 'intermediate',
    keyPoints: [
      'Master-Slave (Primary-Replica): Writes to master, reads from replicas',
      'Master-Master: Both nodes can accept writes',
      'Synchronous Replication: Wait for confirmation before commit',
      'Asynchronous Replication: Fire and forget',
      'Read Replicas: Scale read operations horizontally',
    ],
    pros: ['High availability', 'Read scalability', 'Disaster recovery', 'Geographic distribution'],
    cons: ['Replication lag', 'Conflict resolution', 'Increased complexity'],
    realWorld: ['PostgreSQL streaming replication', 'MySQL Group Replication', 'MongoDB Replica Sets'],
    questions: ['How do you handle split-brain in master-master replication?', 'What\'s the impact of replication lag on consistency?'],
  },
  {
    id: 'cap-theorem',
    title: 'CAP Theorem',
    description: 'Trade-offs between Consistency, Availability, and Partition tolerance',
    icon: <AlertTriangle className="h-6 w-6" />,
    category: 'data',
    difficulty: 'intermediate',
    keyPoints: [
      'Consistency: All nodes see same data at same time',
      'Availability: Every request receives a response',
      'Partition Tolerance: System works despite network failures',
      'CP: Consistent + Partition tolerant (e.g., HBase, MongoDB)',
      'AP: Available + Partition tolerant (e.g., Cassandra, DynamoDB)',
      'CA: Not possible in distributed systems (network failures happen)',
    ],
    pros: ['Clear framework for trade-off decisions', 'Helps choose right database'],
    cons: ['Oversimplification', 'Spectrum rather than binary', 'Doesn\'t account for latency'],
    realWorld: ['Cassandra (AP)', 'MongoDB (CP)', 'Google Spanner (CP with external consistency)'],
    questions: ['Why can\'t you have all three in a distributed system?', 'How does eventual consistency fit into CAP?'],
  },

  // Communication
  {
    id: 'api-design',
    title: 'API Design Patterns',
    description: 'REST, GraphQL, gRPC, and WebSocket communication patterns',
    icon: <ArrowRightLeft className="h-6 w-6" />,
    category: 'communication',
    difficulty: 'intermediate',
    keyPoints: [
      'REST: Stateless, resource-based, HTTP methods',
      'GraphQL: Query language, single endpoint, flexible responses',
      'gRPC: Protocol buffers, HTTP/2, streaming support',
      'WebSocket: Full-duplex, real-time bidirectional',
      'Versioning: URL path, query param, header, or content negotiation',
    ],
    pros: ['Standardized communication', 'Platform independence', 'Easier integration'],
    cons: ['Over-fetching (REST)', 'Complexity (GraphQL)', 'Learning curve'],
    realWorld: ['GitHub REST API', 'Shopify GraphQL', 'Google gRPC services', 'Slack WebSocket'],
    questions: ['When would you choose GraphQL over REST?', 'How do you handle breaking changes in API versioning?'],
  },
  {
    id: 'message-queues',
    title: 'Message Queues & Event Streaming',
    description: 'Asynchronous communication between services',
    icon: <Layers className="h-6 w-6" />,
    category: 'communication',
    difficulty: 'intermediate',
    keyPoints: [
      'Point-to-Point: One producer, one consumer',
      'Pub/Sub: One producer, multiple consumers',
      'Message Ordering: FIFO vs best-effort',
      'At-least-once vs Exactly-once delivery',
      'Dead Letter Queue: Handle failed messages',
      'Backpressure: Handle slow consumers',
    ],
    pros: ['Decoupled services', 'Reliability', 'Load leveling', 'Async processing'],
    cons: ['Added complexity', 'Message ordering challenges', 'Debugging difficulty'],
    realWorld: ['RabbitMQ', 'Apache Kafka', 'AWS SQS/SNS', 'Google Pub/Sub'],
    questions: ['How do you ensure exactly-once processing?', 'Compare Kafka vs RabbitMQ use cases'],
  },
  {
    id: 'microservices',
    title: 'Microservices Architecture',
    description: 'Decompose applications into small, independent services',
    icon: <GitBranch className="h-6 w-6" />,
    category: 'patterns',
    difficulty: 'advanced',
    keyPoints: [
      'Single Responsibility: Each service does one thing well',
      'Independent Deployment: Deploy services without affecting others',
      'Service Discovery: Locate services dynamically',
      'API Gateway: Single entry point for all services',
      'Circuit Breaker: Prevent cascade failures',
      'Saga Pattern: Distributed transactions',
    ],
    pros: ['Independent scaling', 'Technology diversity', 'Fault isolation', 'Team autonomy'],
    cons: ['Distributed system complexity', 'Network latency', 'Data consistency challenges', 'Operational overhead'],
    realWorld: ['Netflix microservices', 'Uber\'s service mesh', 'Amazon\'s two-pizza teams'],
    questions: ['How do you handle distributed transactions in microservices?', 'What\'s the role of a service mesh?'],
  },

  // Patterns
  {
    id: 'cdn',
    title: 'Content Delivery Networks (CDN)',
    description: 'Distribute content from edge locations closer to users',
    icon: <Globe className="h-6 w-6" />,
    category: 'patterns',
    difficulty: 'beginner',
    keyPoints: [
      'Edge Servers: Content cached at geographic locations',
      'Origin Server: Source of truth for content',
      'Cache Hit/Miss: Content served from edge vs origin',
      'TTL: How long content stays cached',
      'Invalidation: Force cache refresh',
      'Dynamic Content: CDN for API responses',
    ],
    pros: ['Reduced latency', 'Lower origin load', 'DDoS protection', 'SSL termination'],
    cons: ['Cache staleness', 'Cost at scale', 'Complex invalidation'],
    realWorld: ['Cloudflare', 'AWS CloudFront', 'Akamai', 'Fastly'],
    questions: ['How do you handle cache invalidation for user-specific content?', 'Compare push vs pull CDN models'],
  },
  {
    id: 'rate-limiting',
    title: 'Rate Limiting & Throttling',
    description: 'Control the rate of requests to protect services',
    icon: <Timer className="h-6 w-6" />,
    category: 'patterns',
    difficulty: 'intermediate',
    keyPoints: [
      'Token Bucket: Tokens replenished at fixed rate',
      'Leaky Bucket: Requests processed at constant rate',
      'Fixed Window: Count requests in time window',
      'Sliding Window: More accurate than fixed window',
      'Distributed Rate Limiting: Across multiple servers',
      'Per-user vs Global limits',
    ],
    pros: ['Protect from abuse', 'Ensure fair usage', 'Prevent resource exhaustion'],
    cons: ['User experience impact', 'Distributed state complexity', 'Clock synchronization'],
    realWorld: ['GitHub API rate limits', 'Twitter API tiers', 'Stripe rate limiting'],
    questions: ['How do you implement rate limiting in a distributed system?', 'Compare token bucket vs leaky bucket algorithms'],
  },
  {
    id: 'id-generation',
    title: 'Distributed ID Generation',
    description: 'Generate unique identifiers across distributed systems',
    icon: <Cpu className="h-6 w-6" />,
    category: 'patterns',
    difficulty: 'intermediate',
    keyPoints: [
      'UUID: 128-bit, collision-resistant, no coordination',
      'Snowflake: Twitter\'s 64-bit, time-ordered IDs',
      'Database Sequences: Centralized counter',
      'Range Allocation: Assign ID ranges to servers',
      'ULID: Time-sortable, 128-bit IDs',
      'KSUID: K-Sortable Unique Identifier',
    ],
    pros: ['No single point of failure', 'Time-ordered', 'Distributed generation'],
    cons: ['Clock dependency', 'Coordination overhead', 'ID length'],
    realWorld: ['Twitter Snowflake', 'Instagram ID generation', 'MongoDB ObjectID'],
    questions: ['How does Snowflake handle clock drift?', 'What are the trade-offs between UUID and Snowflake?'],
  },

  // Scalability
  {
    id: 'consistent-hashing',
    title: 'Consistent Hashing',
    description: 'Distribute data evenly with minimal redistribution on changes',
    icon: <RefreshCw className="h-6 w-6" />,
    category: 'scalability',
    difficulty: 'advanced',
    keyPoints: [
      'Hash Ring: Virtual circle of hash values',
      'Virtual Nodes: Multiple positions per server',
      'Minimal Rebalancing: Only K/N items move on server change',
      'Replication: Multiple copies across ring',
      'Hot Spot Mitigation: Distribute load evenly',
    ],
    pros: ['Even distribution', 'Minimal data movement', 'Horizontal scalability'],
    cons: ['Complexity', 'Virtual node overhead', 'Uneven distribution without virtual nodes'],
    realWorld: ['DynamoDB', 'Cassandra', 'Memcached client libraries', 'Discord'],
    questions: ['Why are virtual nodes important in consistent hashing?', 'How do you handle server addition/removal?'],
  },
  {
    id: 'event-sourcing',
    title: 'Event Sourcing & CQRS',
    description: 'Store state changes as events, separate read/write models',
    icon: <HardDrive className="h-6 w-6" />,
    category: 'scalability',
    difficulty: 'advanced',
    keyPoints: [
      'Event Store: Append-only log of all events',
      'Event Replay: Reconstruct state from events',
      'CQRS: Command Query Responsibility Segregation',
      'Projections: Read models built from events',
      'Event Versioning: Handle schema evolution',
      'Snapshots: Optimize replay performance',
    ],
    pros: ['Complete audit trail', 'Temporal queries', 'Scalable reads', 'Event replay'],
    cons: ['Complexity', 'Eventual consistency', 'Event versioning challenges'],
    realWorld: ['EventStoreDB', 'Axon Framework', 'Kafka-based event sourcing'],
    questions: ['How do you handle event schema changes?', 'When is CQRS overkill?'],
  },
];

const categoryConfig = {
  fundamentals: { label: 'Fundamentals', color: 'from-blue-500 to-cyan-500', bg: 'bg-blue-500/10' },
  data: { label: 'Data Layer', color: 'from-green-500 to-emerald-500', bg: 'bg-green-500/10' },
  communication: { label: 'Communication', color: 'from-purple-500 to-pink-500', bg: 'bg-purple-500/10' },
  patterns: { label: 'Patterns', color: 'from-orange-500 to-red-500', bg: 'bg-orange-500/10' },
  scalability: { label: 'Scalability', color: 'from-indigo-500 to-violet-500', bg: 'bg-indigo-500/10' },
};

const difficultyConfig = {
  beginner: { label: 'Beginner', color: 'bg-green-100 text-green-700 border-green-200' },
  intermediate: { label: 'Intermediate', color: 'bg-yellow-100 text-yellow-700 border-yellow-200' },
  advanced: { label: 'Advanced', color: 'bg-red-100 text-red-700 border-red-200' },
};

export default function SystemDesignPage() {
  const [selectedCategory, setSelectedCategory] = useState<string>('all');
  const [searchQuery, setSearchQuery] = useState('');
  const [viewMode, setViewMode] = useState<'cards' | 'accordion'>('cards');

  const filteredTopics = systemDesignTopics.filter(topic => {
    if (selectedCategory !== 'all' && topic.category !== selectedCategory) return false;
    if (searchQuery && !topic.title.toLowerCase().includes(searchQuery.toLowerCase()) &&
        !topic.description.toLowerCase().includes(searchQuery.toLowerCase())) return false;
    return true;
  });

  const progress = {
    completed: 4,
    total: systemDesignTopics.length,
    byCategory: {
      fundamentals: { completed: 2, total: systemDesignTopics.filter(t => t.category === 'fundamentals').length },
      data: { completed: 1, total: systemDesignTopics.filter(t => t.category === 'data').length },
      communication: { completed: 0, total: systemDesignTopics.filter(t => t.category === 'communication').length },
      patterns: { completed: 1, total: systemDesignTopics.filter(t => t.category === 'patterns').length },
      scalability: { completed: 0, total: systemDesignTopics.filter(t => t.category === 'scalability').length },
    },
  };

  return (
    <div className="min-h-screen bg-gray-50">
      {/* Header */}
      <div className="bg-gradient-to-r from-indigo-600 via-purple-600 to-pink-600 text-white">
        <div className="max-w-7xl mx-auto px-4 py-12">
          <div className="flex flex-col lg:flex-row lg:items-center lg:justify-between gap-6">
            <div>
              <div className="flex items-center gap-3 mb-2">
                <div className="p-2 rounded-xl bg-white/20 backdrop-blur">
                  <Layers className="h-8 w-8" />
                </div>
                <div>
                  <h1 className="text-3xl lg:text-4xl font-bold">System Design</h1>
                  <p className="text-white/80">Master distributed systems architecture</p>
                </div>
              </div>
            </div>
            <div className="flex gap-4">
              <div className="text-center px-6 py-3 rounded-xl bg-white/10 backdrop-blur">
                <div className="text-2xl font-bold">{progress.completed}/{progress.total}</div>
                <div className="text-sm text-white/70">Completed</div>
              </div>
              <Link href="/interviews/session?type=system_design&difficulty=beginner">
                <Button className="h-full bg-white text-purple-600 hover:bg-white/90">
                  <Play className="h-4 w-4 mr-2" />
                  Practice
                </Button>
              </Link>
            </div>
          </div>
        </div>
      </div>

      <div className="max-w-7xl mx-auto px-4 py-8">
        {/* Progress by Category */}
        <div className="grid grid-cols-2 md:grid-cols-5 gap-4 mb-8">
          {Object.entries(categoryConfig).map(([key, config]) => (
            <Card
              key={key}
              className={`cursor-pointer transition-all hover:shadow-lg ${
                selectedCategory === key ? 'ring-2 ring-primary' : ''
              }`}
              onClick={() => setSelectedCategory(selectedCategory === key ? 'all' : key)}
            >
              <CardContent className="p-4">
                <div className="flex items-center gap-2 mb-2">
                  <div className={`w-3 h-3 rounded-full bg-gradient-to-r ${config.color}`} />
                  <span className="font-medium text-sm">{config.label}</span>
                </div>
                <div className="text-2xl font-bold">
                  {progress.byCategory[key as keyof typeof progress.byCategory].completed}/
                  {progress.byCategory[key as keyof typeof progress.byCategory].total}
                </div>
                <Progress
                  value={(progress.byCategory[key as keyof typeof progress.byCategory].completed /
                    progress.byCategory[key as keyof typeof progress.byCategory].total) * 100}
                  className="h-1.5 mt-2"
                />
              </CardContent>
            </Card>
          ))}
        </div>

        {/* Filters & View Toggle */}
        <div className="flex flex-col md:flex-row gap-4 mb-6">
          <div className="flex-1">
            <input
              type="text"
              placeholder="Search topics..."
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
              className="w-full px-4 py-2 border rounded-lg bg-white focus:outline-none focus:ring-2 focus:ring-primary/50"
            />
          </div>
          <div className="flex gap-2">
            <Button
              variant={viewMode === 'cards' ? 'default' : 'outline'}
              onClick={() => setViewMode('cards')}
            >
              Cards
            </Button>
            <Button
              variant={viewMode === 'accordion' ? 'default' : 'outline'}
              onClick={() => setViewMode('accordion')}
            >
              Accordion
            </Button>
          </div>
        </div>

        {/* Content */}
        {viewMode === 'cards' ? (
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            {filteredTopics.map((topic) => {
              const config = categoryConfig[topic.category];
              const diffConfig = difficultyConfig[topic.difficulty];

              return (
                <Card key={topic.id} className="overflow-hidden hover:shadow-xl transition-all group">
                  <div className={`h-2 bg-gradient-to-r ${config.color}`} />
                  <CardHeader>
                    <div className="flex items-center justify-between mb-2">
                      <div className={`p-2 rounded-lg ${config.bg}`}>
                        {topic.icon}
                      </div>
                      <Badge variant="outline" className={diffConfig.color}>
                        {diffConfig.label}
                      </Badge>
                    </div>
                    <CardTitle className="text-lg">{topic.title}</CardTitle>
                    <CardDescription>{topic.description}</CardDescription>
                  </CardHeader>
                  <CardContent>
                    <div className="space-y-3">
                      <div>
                        <p className="text-sm font-medium text-muted-foreground mb-1">Key Points:</p>
                        <ul className="text-sm space-y-1">
                          {topic.keyPoints.slice(0, 2).map((point, i) => (
                            <li key={i} className="flex items-start gap-2">
                              <CheckCircle className="h-4 w-4 text-green-500 shrink-0 mt-0.5" />
                              <span className="text-muted-foreground">{point}</span>
                            </li>
                          ))}
                        </ul>
                      </div>
                      <div className="flex flex-wrap gap-1">
                        {topic.realWorld.slice(0, 3).map((example) => (
                          <Badge key={example} variant="secondary" className="text-xs">
                            {example}
                          </Badge>
                        ))}
                      </div>
                    </div>
                  </CardContent>
                </Card>
              );
            })}
          </div>
        ) : (
          <Accordion type="single" collapsible className="space-y-4">
            {filteredTopics.map((topic) => {
              const config = categoryConfig[topic.category];
              const diffConfig = difficultyConfig[topic.difficulty];

              return (
                <AccordionItem
                  key={topic.id}
                  value={topic.id}
                  className="bg-white rounded-lg border px-6"
                >
                  <AccordionTrigger className="hover:no-underline">
                    <div className="flex items-center gap-4 w-full">
                      <div className={`p-2 rounded-lg ${config.bg}`}>
                        {topic.icon}
                      </div>
                      <div className="flex-1 text-left">
                        <div className="flex items-center gap-2">
                          <span className="font-semibold">{topic.title}</span>
                          <Badge variant="outline" className={diffConfig.color}>
                            {diffConfig.label}
                          </Badge>
                        </div>
                        <p className="text-sm text-muted-foreground">{topic.description}</p>
                      </div>
                    </div>
                  </AccordionTrigger>
                  <AccordionContent className="pt-4 pb-6">
                    <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
                      {/* Key Points */}
                      <div>
                        <h4 className="font-semibold mb-3 flex items-center gap-2">
                          <BookOpen className="h-4 w-4 text-blue-500" />
                          Key Concepts
                        </h4>
                        <ul className="space-y-2">
                          {topic.keyPoints.map((point, i) => (
                            <li key={i} className="flex items-start gap-2 text-sm">
                              <CheckCircle className="h-4 w-4 text-green-500 shrink-0 mt-0.5" />
                              <span>{point}</span>
                            </li>
                          ))}
                        </ul>
                      </div>

                      {/* Pros & Cons */}
                      <div className="space-y-4">
                        <div>
                          <h4 className="font-semibold mb-2 text-green-600">Pros</h4>
                          <ul className="space-y-1">
                            {topic.pros.map((pro, i) => (
                              <li key={i} className="text-sm flex items-center gap-2">
                                <CheckCircle className="h-3 w-3 text-green-500" />
                                {pro}
                              </li>
                            ))}
                          </ul>
                        </div>
                        <div>
                          <h4 className="font-semibold mb-2 text-red-600">Cons</h4>
                          <ul className="space-y-1">
                            {topic.cons.map((con, i) => (
                              <li key={i} className="text-sm flex items-center gap-2">
                                <AlertTriangle className="h-3 w-3 text-red-500" />
                                {con}
                              </li>
                            ))}
                          </ul>
                        </div>
                      </div>

                      {/* Real World Examples */}
                      <div>
                        <h4 className="font-semibold mb-3 flex items-center gap-2">
                          <Globe className="h-4 w-4 text-purple-500" />
                          Real World Examples
                        </h4>
                        <div className="flex flex-wrap gap-2">
                          {topic.realWorld.map((example) => (
                            <Badge key={example} variant="secondary">
                              {example}
                            </Badge>
                          ))}
                        </div>
                      </div>

                      {/* Interview Questions */}
                      <div>
                        <h4 className="font-semibold mb-3 flex items-center gap-2">
                          <ExternalLink className="h-4 w-4 text-orange-500" />
                          Common Interview Questions
                        </h4>
                        <ul className="space-y-2">
                          {topic.questions.map((question, i) => (
                            <li key={i} className="text-sm p-2 rounded bg-muted/50">
                              {question}
                            </li>
                          ))}
                        </ul>
                      </div>
                    </div>

                    <div className="mt-6 pt-4 border-t">
                      <Link href={`/interviews/session?type=system_design&difficulty=${topic.difficulty}`}>
                        <Button className="w-full sm:w-auto">
                          <Play className="h-4 w-4 mr-2" />
                          Practice This Topic
                        </Button>
                      </Link>
                    </div>
                  </AccordionContent>
                </AccordionItem>
              );
            })}
          </Accordion>
        )}

        {/* Quick Reference Card */}
        <Card className="mt-8 border-0 shadow-lg bg-gradient-to-r from-indigo-50 to-purple-50">
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <Shield className="h-5 w-5 text-indigo-500" />
              System Design Interview Cheat Sheet
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
              <div>
                <h4 className="font-semibold mb-2 text-indigo-600">Requirements</h4>
                <ul className="text-sm space-y-1 text-muted-foreground">
                  <li>• Functional requirements</li>
                  <li>• Non-functional requirements</li>
                  <li>• Capacity estimation</li>
                  <li>• Constraints & assumptions</li>
                </ul>
              </div>
              <div>
                <h4 className="font-semibold mb-2 text-purple-600">High-Level Design</h4>
                <ul className="text-sm space-y-1 text-muted-foreground">
                  <li>• API design (REST/GraphQL)</li>
                  <li>• Data model</li>
                  <li>• System components</li>
                  <li>• Data flow diagrams</li>
                </ul>
              </div>
              <div>
                <h4 className="font-semibold mb-2 text-pink-600">Deep Dive</h4>
                <ul className="text-sm space-y-1 text-muted-foreground">
                  <li>• Database schema</li>
                  <li>• Scaling strategy</li>
                  <li>• Caching layer</li>
                  <li>• Load balancing</li>
                </ul>
              </div>
              <div>
                <h4 className="font-semibold mb-2 text-orange-600">Trade-offs</h4>
                <ul className="text-sm space-y-1 text-muted-foreground">
                  <li>• Consistency vs Availability</li>
                  <li>• Latency vs Throughput</li>
                  <li>• Cost vs Performance</li>
                  <li>• Complexity vs Simplicity</li>
                </ul>
              </div>
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  );
}
