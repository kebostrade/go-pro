"use client";

import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import {
  BookOpen,
  MessageSquare,
  Bot,
  Shield,
  Server,
  ArrowRight,
  Clock,
  CheckCircle,
  Sparkles,
  Terminal,
  Code,
  Lock,
  Cloud,
  Cpu,
  Link2,
  Wrench,
  Layers,
  Rocket,
  Users,
  Globe,
  Zap
} from "lucide-react";
import Link from "next/link";

interface Module {
  id: number;
  title: string;
  description: string;
  lessons: number;
  duration: string;
  icon: React.ElementType;
  color: string;
  gradientFrom: string;
  gradientTo: string;
}

const modules: Module[] = [
  {
    id: 1,
    title: "Getting Started",
    description: "Installation, configuration, and understanding OpenClaw architecture",
    lessons: 3,
    duration: "2-3 hours",
    icon: BookOpen,
    color: "text-emerald-600 dark:text-emerald-400",
    gradientFrom: "from-emerald-500/20",
    gradientTo: "to-emerald-500/10",
  },
  {
    id: 2,
    title: "Channel Integration",
    description: "Connect to Telegram, Discord, WhatsApp, and more",
    lessons: 3,
    duration: "2-3 hours",
    icon: MessageSquare,
    color: "text-blue-600 dark:text-blue-400",
    gradientFrom: "from-blue-500/20",
    gradientTo: "to-blue-500/10",
  },
  {
    id: 3,
    title: "Skills & Extensibility",
    description: "Extend OpenClaw with community and custom Skills",
    lessons: 2,
    duration: "3-4 hours",
    icon: Wrench,
    color: "text-purple-600 dark:text-purple-400",
    gradientFrom: "from-purple-500/20",
    gradientTo: "to-purple-500/10",
  },
  {
    id: 4,
    title: "Production Deployment",
    description: "Docker, VPS, security hardening, and HTTPS setup",
    lessons: 3,
    duration: "4-5 hours",
    icon: Server,
    color: "text-orange-600 dark:text-orange-400",
    gradientFrom: "from-orange-500/20",
    gradientTo: "to-orange-500/10",
  },
];

const features = [
  {
    icon: Bot,
    title: "Self-Hosted AI Agent",
    description: "Run your own AI assistant with complete data privacy",
    color: "from-emerald-500/20 to-emerald-500/5",
  },
  {
    icon: Link2,
    title: "Multi-Channel Support",
    description: "Telegram, Discord, WhatsApp, Slack, Signal, and more",
    color: "from-blue-500/20 to-blue-500/5",
  },
  {
    icon: Layers,
    title: "Skills System",
    description: "Extend capabilities with community and custom Skills",
    color: "from-purple-500/20 to-purple-500/5",
  },
  {
    icon: Shield,
    title: "Privacy First",
    description: "Your data stays on your infrastructure",
    color: "from-red-500/20 to-red-500/5",
  },
  {
    icon: Cpu,
    title: "Local Models",
    description: "Run with Ollama for complete offline capability",
    color: "from-cyan-500/20 to-cyan-500/5",
  },
  {
    icon: Globe,
    title: "Open Source",
    description: "150K+ GitHub stars, community-driven development",
    color: "from-amber-500/20 to-amber-500/5",
  },
];

const quickStartSteps = [
  {
    step: "01",
    title: "Install OpenClaw",
    description: "Use Homebrew, Docker, or install from source",
    code: "brew install openclaw",
  },
  {
    step: "02",
    title: "Configure API Key",
    description: "Add your OpenAI or Anthropic API key",
    code: "ANTHROPIC_API_KEY=sk-ant-xxx",
  },
  {
    step: "03",
    title: "Connect a Channel",
    description: "Set up Telegram, Discord, or another platform",
    code: 'channels.telegram.enabled = true',
  },
  {
    step: "04",
    title: "Start Chatting",
    description: "Message your bot and ask questions!",
    code: "openclaw start",
  },
];

export default function OpenClawPage() {
  const totalLessons = modules.reduce((acc, m) => acc + m.lessons, 0);

  return (
    <div className="min-h-screen bg-gradient-to-br from-background via-background to-accent/5">
      <div className="container-responsive padding-responsive-y">
        {/* Hero Section */}
        <div className="margin-responsive mb-12">
          <div className="flex flex-col items-center text-center space-y-6">
            <div className="inline-flex items-center gap-2 px-4 py-2 rounded-full bg-gradient-to-r from-emerald-500/10 to-emerald-500/5 border border-emerald-500/20 mb-4">
              <Sparkles className="h-4 w-4 text-emerald-500" />
              <span className="text-sm font-medium text-emerald-600 dark:text-emerald-400">11 Lessons</span>
              <span className="text-sm text-muted-foreground">•</span>
              <span className="text-sm font-medium text-emerald-600 dark:text-emerald-400">4 Projects</span>
            </div>

            <h1 className="text-3xl sm:text-4xl lg:text-5xl font-bold tracking-tight">
              Build Your Own{" "}
              <span className="bg-gradient-to-r from-emerald-500 via-teal-500 to-cyan-500 bg-clip-text text-transparent">
                AI Agent
              </span>
              <br />
              with OpenClaw
            </h1>

            <p className="text-muted-foreground text-sm sm:text-base lg:text-lg max-w-2xl leading-relaxed">
              Master the fastest-growing open-source AI agent platform (150K+ GitHub stars). 
              Learn to build self-hosted AI assistants that connect to Telegram, Discord, WhatsApp, 
              and more—while keeping your data private.
            </p>

            <div className="flex flex-col sm:flex-row gap-3 pt-4">
              <Link href="/learn/openclaw/OC-01">
                <Button size="lg" className="bg-gradient-to-r from-emerald-500 to-teal-500 hover:from-emerald-600 hover:to-teal-600 text-white border-0">
                  Start Learning
                  <ArrowRight className="ml-2 h-4 w-4" />
                </Button>
              </Link>
              <Link href="/docs/OPENCLAW_GUIDE.md" target="_blank">
                <Button size="lg" variant="outline">
                  <BookOpen className="mr-2 h-4 w-4" />
                  Read Guide
                </Button>
              </Link>
            </div>
          </div>
        </div>

        {/* Stats */}
        <div className="grid grid-cols-2 md:grid-cols-4 gap-4 mb-12">
          {[
            { label: "Lessons", value: "11", icon: BookOpen },
            { label: "Projects", value: "4", icon: Rocket },
            { label: "Hours", value: "12+", icon: Clock },
            { label: "Channels", value: "8+", icon: Link2 },
          ].map((stat, i) => (
            <Card key={i} className="border-0 bg-card/50 backdrop-blur-sm">
              <CardContent className="flex items-center gap-3 p-4">
                <div className="p-2 rounded-lg bg-gradient-to-br from-emerald-500/20 to-teal-500/10">
                  <stat.icon className="h-4 w-4 text-emerald-500" />
                </div>
                <div>
                  <p className="text-2xl font-bold">{stat.value}</p>
                  <p className="text-xs text-muted-foreground">{stat.label}</p>
                </div>
              </CardContent>
            </Card>
          ))}
        </div>

        {/* Features Grid */}
        <div className="margin-responsive mb-12">
          <h2 className="text-xl sm:text-2xl font-bold text-center mb-8">
            Why <span className="text-emerald-500">OpenClaw</span>?
          </h2>
          <div className="grid sm:grid-cols-2 lg:grid-cols-3 gap-4">
            {features.map((feature, i) => (
              <Card key={i} className="group hover:shadow-lg transition-all duration-300 border-0 bg-card/50 backdrop-blur-sm overflow-hidden">
                <div className={`absolute inset-0 bg-gradient-to-br ${feature.color} opacity-0 group-hover:opacity-100 transition-opacity duration-500`} />
                <CardContent className="p-5 relative z-10">
                  <div className="flex items-start gap-3">
                    <div className="p-2 rounded-lg bg-background/80">
                      <feature.icon className="h-5 w-5 text-emerald-500" />
                    </div>
                    <div>
                      <h3 className="font-semibold mb-1">{feature.title}</h3>
                      <p className="text-xs text-muted-foreground">{feature.description}</p>
                    </div>
                  </div>
                </CardContent>
              </Card>
            ))}
          </div>
        </div>

        {/* Quick Start */}
        <div className="margin-responsive mb-12">
          <h2 className="text-xl sm:text-2xl font-bold text-center mb-8">
            Quick Start
          </h2>
          <div className="grid sm:grid-cols-2 lg:grid-cols-4 gap-4">
            {quickStartSteps.map((item, i) => (
              <div key={i} className="relative">
                <div className="absolute -top-2 -left-2 text-6xl font-bold text-emerald-500/10">{item.step}</div>
                <Card className="h-full border-0 bg-card/50 backdrop-blur-sm">
                  <CardContent className="p-5">
                    <h3 className="font-semibold mb-2 flex items-center gap-2">
                      <Zap className="h-4 w-4 text-emerald-500" />
                      {item.title}
                    </h3>
                    <p className="text-xs text-muted-foreground mb-3">{item.description}</p>
                    <code className="text-xs bg-muted px-2 py-1 rounded block overflow-x-auto">
                      {item.code}
                    </code>
                  </CardContent>
                </Card>
              </div>
            ))}
          </div>
        </div>

        {/* Course Modules */}
        <div className="margin-responsive mb-12">
          <h2 className="text-xl sm:text-2xl font-bold text-center mb-8">
            Course Modules
          </h2>
          <div className="grid gap-4">
            {modules.map((module, i) => (
              <Link key={module.id} href={`/learn/openclaw/OC-0${module.id}`}>
                <Card className="group hover:shadow-xl transition-all duration-300 hover:-translate-y-1 cursor-pointer border-0 bg-card/50 backdrop-blur-sm">
                  <CardContent className="p-5 sm:p-6">
                    <div className="flex flex-col sm:flex-row sm:items-center gap-4">
                      <div className={`p-4 rounded-xl bg-gradient-to-br ${module.gradientFrom} ${module.gradientTo}`}>
                        <module.icon className={`h-6 w-6 ${module.color}`} />
                      </div>
                      <div className="flex-1">
                        <div className="flex items-center gap-2 mb-1">
                          <span className="text-sm font-medium text-muted-foreground">Module {module.id}</span>
                          <Badge variant="secondary" className="text-xs">
                            {module.lessons} lessons
                          </Badge>
                        </div>
                        <h3 className="text-lg font-semibold mb-1">{module.title}</h3>
                        <p className="text-sm text-muted-foreground">{module.description}</p>
                      </div>
                      <div className="flex items-center gap-2 text-muted-foreground">
                        <Clock className="h-4 w-4" />
                        <span className="text-sm">{module.duration}</span>
                        <ArrowRight className="h-4 w-4 group-hover:translate-x-1 transition-transform" />
                      </div>
                    </div>
                  </CardContent>
                </Card>
              </Link>
            ))}
          </div>
        </div>

        {/* Prerequisites */}
        <Card className="margin-responsive mb-12 border-0 bg-gradient-to-r from-amber-500/10 to-orange-500/10">
          <CardContent className="p-6">
            <div className="flex items-start gap-4">
              <div className="p-3 rounded-lg bg-amber-500/20">
                <CheckCircle className="h-6 w-6 text-amber-500" />
              </div>
              <div>
                <h3 className="font-semibold mb-2">Prerequisites</h3>
                <ul className="text-sm text-muted-foreground space-y-1">
                  <li>• API key from OpenAI or Anthropic</li>
                  <li>• Node.js 18+ (or Docker for container approach)</li>
                  <li>• Basic command line knowledge</li>
                  <li>• 4GB RAM minimum (8GB recommended)</li>
                </ul>
              </div>
            </div>
          </CardContent>
        </Card>

        {/* CTA */}
        <div className="text-center">
          <h3 className="text-lg font-semibold mb-4">Ready to build your AI agent?</h3>
          <Link href="/learn/openclaw/OC-01">
            <Button size="lg" className="bg-gradient-to-r from-emerald-500 to-teal-500 hover:from-emerald-600 hover:to-teal-600 text-white border-0">
              Start Module 1: Getting Started
              <ArrowRight className="ml-2 h-4 w-4" />
            </Button>
          </Link>
        </div>
      </div>
    </div>
  );
}
