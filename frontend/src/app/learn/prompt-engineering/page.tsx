"use client";

import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import {
  BookOpen,
  Wrench,
  GitBranch,
  Target,
  Rocket,
  ArrowRight,
  Clock,
  CheckCircle,
  Sparkles
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
    title: "Foundations",
    description: "LLM basics, prompt anatomy, and zero-shot prompting fundamentals",
    lessons: 4,
    duration: "2-3 hours",
    icon: BookOpen,
    color: "text-blue-600 dark:text-blue-400",
    gradientFrom: "from-blue-500/20",
    gradientTo: "to-blue-500/10",
  },
  {
    id: 2,
    title: "Core Techniques",
    description: "Few-shot, chain-of-thought, and structured output patterns",
    lessons: 4,
    duration: "3-4 hours",
    icon: Wrench,
    color: "text-green-600 dark:text-green-400",
    gradientFrom: "from-green-500/20",
    gradientTo: "to-green-500/10",
  },
  {
    id: 3,
    title: "Advanced Patterns",
    description: "ReAct, self-consistency, tree-of-thoughts, and prompt chaining",
    lessons: 4,
    duration: "4-5 hours",
    icon: GitBranch,
    color: "text-purple-600 dark:text-purple-400",
    gradientFrom: "from-purple-500/20",
    gradientTo: "to-purple-500/10",
  },
  {
    id: 4,
    title: "Task-Specific",
    description: "Code generation, data analysis, and creative writing prompts",
    lessons: 4,
    duration: "3-4 hours",
    icon: Target,
    color: "text-orange-600 dark:text-orange-400",
    gradientFrom: "from-orange-500/20",
    gradientTo: "to-orange-500/10",
  },
  {
    id: 5,
    title: "Production",
    description: "Evaluation, optimization, security, and deployment strategies",
    lessons: 4,
    duration: "4-5 hours",
    icon: Rocket,
    color: "text-red-600 dark:text-red-400",
    gradientFrom: "from-red-500/20",
    gradientTo: "to-red-500/10",
  },
];

export default function PromptEngineeringPage() {
  const totalLessons = modules.reduce((acc, m) => acc + m.lessons, 0);

  return (
    <div className="min-h-screen bg-gradient-to-br from-background via-background to-accent/5">
      <div className="container-responsive padding-responsive-y">
        {/* Hero Section */}
        <div className="margin-responsive mb-12">
          <div className="flex flex-col items-center text-center space-y-6">
            <div className="inline-flex items-center gap-2 px-4 py-2 rounded-full bg-gradient-to-r from-blue-500/10 to-blue-500/5 border border-blue-500/20 mb-4">
              <Sparkles className="h-4 w-4 text-blue-500" />
              <span className="text-sm font-medium text-blue-600 dark:text-blue-400">20 Lessons</span>
              <span className="text-muted-foreground">·</span>
              <span className="text-sm text-muted-foreground">16-21 Hours</span>
            </div>

            <h1 className="text-4xl md:text-5xl lg:text-6xl font-bold tracking-tight bg-gradient-to-r from-primary via-primary/80 to-primary/60 bg-clip-text text-transparent animate-in fade-in slide-in-from-bottom-4 duration-700">
              Prompt Engineering
            </h1>

            <p className="text-lg md:text-xl text-muted-foreground max-w-2xl animate-in fade-in slide-in-from-bottom-5 duration-700 delay-100">
              Master the art of communicating with AI. Learn techniques from basic prompting to advanced agentic patterns used in production systems.
            </p>

            <div className="flex flex-wrap items-center justify-center gap-4 pt-4">
              <Link href="/learn/prompt-engineering/PE-01">
                <Button size="lg" className="group">
                  Start Learning
                  <ArrowRight className="ml-2 h-4 w-4 group-hover:translate-x-1 transition-transform" />
                </Button>
              </Link>
              <Button variant="outline" size="lg">
                View Curriculum
              </Button>
            </div>
          </div>
        </div>

        {/* Stats Overview */}
        <div className="grid grid-cols-1 sm:grid-cols-3 gap-4 mb-12 max-w-3xl mx-auto">
          <Card suppressHydrationWarning className="border-blue-200/50 dark:border-blue-800/50">
            <CardContent className="p-6 text-center">
              <div className="text-3xl font-bold bg-gradient-to-br from-blue-600 to-blue-500 bg-clip-text text-transparent">5</div>
              <div className="text-sm font-medium text-muted-foreground mt-1">Modules</div>
            </CardContent>
          </Card>
          <Card suppressHydrationWarning className="border-green-200/50 dark:border-green-800/50">
            <CardContent className="p-6 text-center">
              <div className="text-3xl font-bold bg-gradient-to-br from-green-600 to-green-500 bg-clip-text text-transparent">{totalLessons}</div>
              <div className="text-sm font-medium text-muted-foreground mt-1">Lessons</div>
            </CardContent>
          </Card>
          <Card suppressHydrationWarning className="border-orange-200/50 dark:border-orange-800/50">
            <CardContent className="p-6 text-center">
              <div className="text-3xl font-bold bg-gradient-to-br from-orange-600 to-orange-500 bg-clip-text text-transparent">16-21h</div>
              <div className="text-sm font-medium text-muted-foreground mt-1">Duration</div>
            </CardContent>
          </Card>
        </div>

        {/* Module Cards */}
        <div className="space-y-6">
          <h2 className="text-2xl font-bold text-center mb-8">Course Modules</h2>

          <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
            {modules.map((module, index) => (
              <Link key={module.id} href={`/learn/prompt-engineering/${module.id}`}>
                <Card
                  suppressHydrationWarning
                  className={`group relative overflow-hidden border-border/50 hover:border-primary/40 hover:shadow-xl hover:shadow-primary/5 transition-all duration-500 hover:-translate-y-1 cursor-pointer ${
                    index === 0 ? "lg:col-span-2" : ""
                  }`}
                >
                  <div className="absolute inset-0 bg-gradient-to-br from-primary/5 via-transparent to-transparent opacity-0 group-hover:opacity-100 transition-opacity duration-500" />

                  <CardContent className="relative p-6">
                    <div className="flex items-start gap-5">
                      <div className={`p-4 rounded-2xl bg-gradient-to-br ${module.gradientFrom} ${module.gradientTo} group-hover:scale-110 transition-all duration-300 shrink-0`}>
                        <module.icon className={`h-7 w-7 ${module.color}`} />
                      </div>

                      <div className="flex-1 space-y-3">
                        <div className="flex items-center gap-3 flex-wrap">
                          <Badge variant="outline" className="text-xs font-bold tracking-wider uppercase border-primary/30 text-primary">
                            Module {module.id}
                          </Badge>
                          <Badge variant="secondary" className="text-xs">
                            {module.lessons} Lessons
                          </Badge>
                          <div className="flex items-center gap-1.5 text-xs text-muted-foreground">
                            <Clock className="h-3.5 w-3.5" />
                            <span>{module.duration}</span>
                          </div>
                        </div>

                        <h3 className="text-xl font-bold group-hover:text-primary transition-colors duration-300">
                          {module.title}
                        </h3>

                        <p className="text-sm text-muted-foreground leading-relaxed">
                          {module.description}
                        </p>

                        <div className="flex items-center gap-2 pt-2">
                          <span className="text-sm font-medium text-primary group-hover:underline">
                            Start Module
                          </span>
                          <ArrowRight className="h-4 w-4 text-primary group-hover:translate-x-1 transition-transform" />
                        </div>
                      </div>
                    </div>
                  </CardContent>
                </Card>
              </Link>
            ))}
          </div>
        </div>

        {/* CTA Section */}
        <div className="mt-16 text-center">
          <Card suppressHydrationWarning className="relative overflow-hidden border-primary/20 bg-gradient-to-br from-primary/5 via-background to-background max-w-2xl mx-auto">
            <div className="absolute top-0 right-0 w-64 h-64 bg-gradient-to-br from-primary/10 to-transparent rounded-full blur-3xl -z-10" />
            <CardContent className="p-8">
              <CheckCircle className="h-12 w-12 text-primary mx-auto mb-4" />
              <h3 className="text-2xl font-bold mb-2">Ready to Master Prompts?</h3>
              <p className="text-muted-foreground mb-6">
                Start with the foundations and work your way up to production-ready prompting techniques.
              </p>
              <Link href="/learn/prompt-engineering/PE-01">
                <Button size="lg" className="group">
                  Begin with Module 1
                  <ArrowRight className="ml-2 h-4 w-4 group-hover:translate-x-1 transition-transform" />
                </Button>
              </Link>
            </CardContent>
          </Card>
        </div>
      </div>
    </div>
  );
}
