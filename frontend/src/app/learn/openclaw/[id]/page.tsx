"use client";

import { useState, useEffect } from 'react';
import { useRouter, useParams } from 'next/navigation';
import Link from 'next/link';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { Card, CardContent } from '@/components/ui/card';
import {
  ArrowLeft,
  ArrowRight,
  BookOpen,
  CheckCircle,
  Clock,
  Menu,
  X,
  FileText,
  Code,
  Terminal,
  Shield,
  Server,
  MessageSquare,
  Wrench,
  Bot,
  Sparkles
} from 'lucide-react';

// Lesson data - in production this would come from the markdown files
const lessons = [
  {
    id: 'OC-01',
    title: 'What is OpenClaw? Architecture Deep Dive',
    module: 'Module 1: Getting Started',
    duration: '2 hours',
    description: 'Understand the Gateway, Nodes, and Skills architecture',
    topics: ['Architecture overview', 'Gateway components', 'Agent loop', 'Skills system'],
    icon: BookOpen,
    color: 'from-emerald-500 to-teal-500',
    prev: null,
    next: 'OC-02',
    markdownFile: 'course/openclaw/lessons/M1-getting-started/01-introduction/README.md'
  },
  {
    id: 'OC-02',
    title: 'Installation: Homebrew, Docker, Source',
    module: 'Module 1: Getting Started',
    duration: '2 hours',
    description: 'Install OpenClaw on macOS, Linux, or Windows',
    topics: ['Homebrew installation', 'Docker setup', 'Source installation', 'Verification'],
    icon: Terminal,
    color: 'from-emerald-500 to-teal-500',
    prev: 'OC-01',
    next: 'OC-03',
    markdownFile: 'course/openclaw/lessons/M1-getting-started/02-installation/README.md'
  },
  {
    id: 'OC-03',
    title: 'LLM Providers & Initial Configuration',
    module: 'Module 1: Getting Started',
    duration: '2 hours',
    description: 'Configure OpenAI, Anthropic, or Ollama providers',
    topics: ['Provider configuration', 'Environment variables', 'Gateway settings', 'Testing'],
    icon: Bot,
    color: 'from-emerald-500 to-teal-500',
    prev: 'OC-02',
    next: 'OC-04',
    markdownFile: 'course/openclaw/lessons/M1-getting-started/03-configuration/README.md'
  },
  {
    id: 'OC-04',
    title: 'Telegram Bot Integration',
    module: 'Module 2: Channel Integration',
    duration: '2 hours',
    description: 'Connect OpenClaw to Telegram for instant messaging',
    topics: ['BotFather setup', 'Configuration', 'Chat IDs', 'Access control'],
    icon: MessageSquare,
    color: 'from-blue-500 to-cyan-500',
    prev: 'OC-03',
    next: 'OC-05',
    markdownFile: 'course/openclaw/lessons/M2-channels/04-telegram/README.md'
  },
  {
    id: 'OC-05',
    title: 'Discord Bot Integration',
    module: 'Module 2: Channel Integration',
    duration: '2 hours',
    description: 'Add OpenClaw as a Discord bot',
    topics: ['Developer portal', 'OAuth setup', 'Permissions', 'Channel config'],
    icon: MessageSquare,
    color: 'from-blue-500 to-cyan-500',
    prev: 'OC-04',
    next: 'OC-06',
    markdownFile: 'course/openclaw/lessons/M2-channels/05-discord/README.md'
  },
  {
    id: 'OC-06',
    title: 'WhatsApp & Signal Integration',
    module: 'Module 2: Channel Integration',
    duration: '2 hours',
    description: 'Connect to WhatsApp and Signal',
    topics: ['QR pairing', 'Session management', 'Signal setup', 'Best practices'],
    icon: MessageSquare,
    color: 'from-blue-500 to-cyan-500',
    prev: 'OC-05',
    next: 'OC-07',
    markdownFile: 'course/openclaw/lessons/M2-channels/06-whatsapp-signal/README.md'
  },
  {
    id: 'OC-07',
    title: 'Skills System: Finding & Installing',
    module: 'Module 3: Skills & Extensibility',
    duration: '2 hours',
    description: 'Use community Skills to extend OpenClaw',
    topics: ['Finding skills', 'Installation', 'Configuration', 'Built-in skills'],
    icon: Wrench,
    color: 'from-purple-500 to-pink-500',
    prev: 'OC-06',
    next: 'OC-08',
    markdownFile: 'course/openclaw/lessons/M3-skills/07-skills-basics/README.md'
  },
  {
    id: 'OC-08',
    title: 'Building Custom Skills',
    module: 'Module 3: Skills & Extensibility',
    duration: '3 hours',
    description: 'Create your own custom Skills',
    topics: ['Skill manifest', 'Tool implementation', 'Registration', 'Testing'],
    icon: Code,
    color: 'from-purple-500 to-pink-500',
    prev: 'OC-07',
    next: 'OC-09',
    markdownFile: 'course/openclaw/lessons/M3-skills/08-custom-skills/README.md'
  },
  {
    id: 'OC-09',
    title: 'Docker & Docker Compose Deployment',
    module: 'Module 4: Production Deployment',
    duration: '2 hours',
    description: 'Deploy OpenClaw using containers',
    topics: ['docker-compose', 'Volume management', 'Operations', 'Backups'],
    icon: Server,
    color: 'from-orange-500 to-red-500',
    prev: 'OC-08',
    next: 'OC-10',
    markdownFile: 'course/openclaw/lessons/M4-production/09-docker-deployment/README.md'
  },
  {
    id: 'OC-10',
    title: 'VPS Production Deployment',
    module: 'Module 4: Production Deployment',
    duration: '3 hours',
    description: 'Deploy on cloud VPS with HTTPS',
    topics: ['VPS setup', 'Nginx reverse proxy', 'SSL/TLS', 'Automated backups'],
    icon: Server,
    color: 'from-orange-500 to-red-500',
    prev: 'OC-09',
    next: 'OC-11',
    markdownFile: 'course/openclaw/lessons/M4-production/10-vps-production/README.md'
  },
  {
    id: 'OC-11',
    title: 'Security Best Practices',
    module: 'Module 4: Production Deployment',
    duration: '2 hours',
    description: 'Secure your OpenClaw deployment',
    topics: ['Access control', 'API key management', 'Firewall', 'Incident response'],
    icon: Shield,
    color: 'from-orange-500 to-red-500',
    prev: 'OC-10',
    next: null,
    markdownFile: 'course/openclaw/lessons/M4-production/11-security-hardening/README.md'
  },
];

export default function OpenClawLessonPage() {
  const params = useParams();
  const router = useRouter();
  const [sidebarOpen, setSidebarOpen] = useState(false);
  const [content, setContent] = useState<string>('');
  const [loading, setLoading] = useState(true);

  const lessonId = params?.id as string | undefined;
  const currentLesson = lessons.find(l => l.id === lessonId);
  const currentIndex = lessons.findIndex(l => l.id === lessonId);

  useEffect(() => {
    // In a real app, we would fetch the markdown content here
    setLoading(false);
  }, [lessonId]);

  if (loading) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-emerald-500" />
      </div>
    );
  }

  if (!currentLesson) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="text-center">
          <h1 className="text-2xl font-bold mb-4">Lesson not found</h1>
          <Link href="/learn/openclaw">
            <Button>Back to OpenClaw Course</Button>
          </Link>
        </div>
      </div>
    );
  }

  const Icon = currentLesson.icon;

  return (
    <div className="min-h-screen bg-gradient-to-br from-background via-background to-accent/5">
      {/* Mobile sidebar toggle */}
      <div className="lg:hidden fixed top-4 left-4 z-50">
        <Button
          variant="outline"
          size="icon"
          onClick={() => setSidebarOpen(!sidebarOpen)}
          className="bg-background/80 backdrop-blur-sm"
        >
          {sidebarOpen ? <X className="h-5 w-5" /> : <Menu className="h-5 w-5" />}
        </Button>
      </div>

      {/* Sidebar */}
      <aside className={`
        fixed inset-y-0 left-0 z-40 w-72 bg-card/95 backdrop-blur-sm border-r transform transition-transform duration-300
        ${sidebarOpen ? 'translate-x-0' : '-translate-x-full lg:translate-x-0'}
      `}>
        <div className="h-full overflow-y-auto p-4">
          <Link href="/learn/openclaw" className="flex items-center gap-2 text-sm text-muted-foreground hover:text-foreground mb-4">
            <ArrowLeft className="h-4 w-4" />
            Back to Course
          </Link>
          
          <div className="space-y-2">
            {lessons.map((lesson, index) => (
              <Link
                key={lesson.id}
                href={`/learn/openclaw/${lesson.id}`}
                onClick={() => setSidebarOpen(false)}
                className={`
                  block p-3 rounded-lg transition-all duration-200
                  ${lesson.id === lessonId 
                    ? 'bg-gradient-to-r from-emerald-500/20 to-teal-500/10 border border-emerald-500/30' 
                    : 'hover:bg-muted/50'}
                `}
              >
                <div className="flex items-center gap-2 mb-1">
                  <span className="text-xs font-medium text-muted-foreground">{index + 1}.</span>
                  <span className={`text-sm font-medium ${lesson.id === lessonId ? 'text-emerald-500' : ''}`}>
                    {lesson.id}
                  </span>
                </div>
                <p className="text-xs text-muted-foreground line-clamp-2">{lesson.title}</p>
              </Link>
            ))}
          </div>
        </div>
      </aside>

      {/* Main content */}
      <div className="lg:ml-72 min-h-screen">
        <div className="container-responsive padding-responsive-y">
          {/* Header */}
          <div className="margin-responsive mb-8">
            <div className="flex items-center gap-2 mb-2">
              <Badge variant="outline" className="text-emerald-500 border-emerald-500/30">
                {currentLesson.module}
              </Badge>
              <Badge variant="secondary">
                <Clock className="h-3 w-3 mr-1" />
                {currentLesson.duration}
              </Badge>
            </div>
            
            <h1 className="text-2xl sm:text-3xl lg:text-4xl font-bold mb-4 flex items-center gap-3">
              <div className={`p-2 rounded-xl bg-gradient-to-br ${currentLesson.color}`}>
                <Icon className="h-6 w-6 text-white" />
              </div>
              {currentLesson.title}
            </h1>
            
            <p className="text-muted-foreground text-lg">{currentLesson.description}</p>
          </div>

          {/* Topics */}
          <Card className="margin-responsive mb-8 border-0 bg-card/50 backdrop-blur-sm">
            <CardContent className="p-5">
              <h3 className="font-semibold mb-3 flex items-center gap-2">
                <Sparkles className="h-4 w-4 text-emerald-500" />
                What you will learn
              </h3>
              <div className="grid sm:grid-cols-2 gap-2">
                {currentLesson.topics.map((topic, i) => (
                  <div key={i} className="flex items-center gap-2 text-sm">
                    <CheckCircle className="h-4 w-4 text-emerald-500 flex-shrink-0" />
                    <span>{topic}</span>
                  </div>
                ))}
              </div>
            </CardContent>
          </Card>

          {/* Content placeholder */}
          <div className="margin-responsive mb-8">
            <Card className="border-0 bg-card/50 backdrop-blur-sm">
              <CardContent className="p-6 sm:p-8">
                {loading ? (
                  <div className="flex items-center justify-center py-12">
                    <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-emerald-500" />
                  </div>
                ) : (
                  <div className="prose prose-sm sm:prose dark:prose-invert max-w-none">
                    <div className="text-center py-12 text-muted-foreground">
                      <FileText className="h-12 w-12 mx-auto mb-4 opacity-50" />
                      <p>Lesson content would be loaded from:</p>
                      <code className="text-xs bg-muted px-2 py-1 rounded mt-2 inline-block">
                        {currentLesson.markdownFile}
                      </code>
                    </div>
                  </div>
                )}
              </CardContent>
            </Card>
          </div>

          {/* Navigation */}
          <div className="flex flex-col sm:flex-row justify-between gap-4 margin-responsive">
            {currentLesson.prev ? (
              <Link href={`/learn/openclaw/${currentLesson.prev}`} className="flex-1">
                <Button variant="outline" className="w-full justify-start">
                  <ArrowLeft className="mr-2 h-4 w-4" />
                  {currentLesson.prev}
                </Button>
              </Link>
            ) : (
              <div className="flex-1" />
            )}
            
            {currentLesson.next ? (
              <Link href={`/learn/openclaw/${currentLesson.next}`} className="flex-1">
                <Button className="w-full justify-end bg-gradient-to-r from-emerald-500 to-teal-500 hover:from-emerald-600 hover:to-teal-600 text-white border-0">
                  {currentLesson.next}
                  <ArrowRight className="ml-2 h-4 w-4" />
                </Button>
              </Link>
            ) : (
              <Link href="/learn/openclaw" className="flex-1">
                <Button className="w-full bg-gradient-to-r from-emerald-500 to-teal-500 hover:from-emerald-600 hover:to-teal-600 text-white border-0">
                  Complete Course
                  <CheckCircle className="ml-2 h-4 w-4" />
                </Button>
              </Link>
            )}
          </div>
        </div>
      </div>
    </div>
  );
}
