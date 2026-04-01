'use client';

import { useState, useEffect } from "react";
import { useParams, useRouter } from "next/navigation";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Progress } from "@/components/ui/progress";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import {
  BookOpen,
  Code2,
  Trophy,
  ArrowRight,
  ArrowLeft,
  Clock,
  Target,
  CheckCircle,
  Home,
  ChevronRight,
  Lightbulb,
  FileText,
  Sparkles,
  Zap,
  Brain,
  Rocket,
  Settings
} from "lucide-react";
import Link from "next/link";

const lessons = [
  { id: "PE-01", title: "Introduction to Prompt Engineering", module: "Foundations", moduleNumber: 1, duration: "15 min", difficulty: "beginner", description: "Learn the fundamentals of prompt engineering and why it matters for AI interactions.", summary: "Covers what prompt engineering is, its importance, and how effective prompts can dramatically improve AI outputs.", objectives: ["Understand what prompt engineering is and why it matters", "Learn the history and evolution of prompt engineering", "Identify key use cases for prompt engineering", "Recognize the impact of well-crafted prompts"] },
  { id: "PE-02", title: "LLM Basics & Capabilities", module: "Foundations", moduleNumber: 1, duration: "20 min", difficulty: "beginner", description: "Explore how Large Language Models work and their capabilities.", summary: "Deep dive into how LLMs process text, their strengths, limitations, and how to leverage them effectively.", objectives: ["Understand how LLMs process and generate text", "Learn about token limits and context windows", "Identify LLM strengths and limitations", "Explore different LLM capabilities"] },
  { id: "PE-03", title: "Anatomy of a Prompt", module: "Foundations", moduleNumber: 1, duration: "20 min", difficulty: "beginner", description: "Break down the essential components that make up an effective prompt.", summary: "Learn the structure of prompts including instructions, context, input data, and output format specifications.", objectives: ["Identify the core components of a prompt", "Learn how to structure instructions clearly", "Understand the role of context in prompts", "Practice specifying output formats"] },
  { id: "PE-04", title: "Zero-Shot Prompting", module: "Foundations", moduleNumber: 1, duration: "15 min", difficulty: "beginner", description: "Master the art of prompting without examples.", summary: "Learn to craft effective prompts that work without providing examples, leveraging the model's inherent knowledge.", objectives: ["Understand what zero-shot prompting means", "Learn when to use zero-shot vs other approaches", "Craft clear zero-shot instructions", "Identify tasks suitable for zero-shot"] },
  { id: "PE-05", title: "Few-Shot Prompting", module: "Core Techniques", moduleNumber: 2, duration: "25 min", difficulty: "intermediate", description: "Learn to guide AI behavior with strategic examples.", summary: "Master the technique of providing examples in prompts to improve output quality and consistency.", objectives: ["Understand few-shot prompting and its benefits", "Learn to select effective examples", "Determine optimal number of examples", "Avoid common few-shot pitfalls"] },
  { id: "PE-06", title: "Chain-of-Thought Prompting", module: "Core Techniques", moduleNumber: 2, duration: "30 min", difficulty: "intermediate", description: "Teach AI to reason step-by-step for complex problems.", summary: "Learn to elicit reasoning chains from AI, improving accuracy on complex tasks like math and logic.", objectives: ["Understand chain-of-thought reasoning", "Learn to prompt for step-by-step thinking", "Apply CoT to mathematical problems", "Use CoT for logical reasoning tasks"] },
  { id: "PE-07", title: "Structured Output Generation", module: "Core Techniques", moduleNumber: 2, duration: "25 min", difficulty: "intermediate", description: "Generate consistent, parseable outputs from AI models.", summary: "Techniques for getting JSON, tables, and other structured formats from LLMs reliably.", objectives: ["Learn to request specific output formats", "Generate valid JSON from prompts", "Create tables and structured data", "Handle parsing and validation"] },
  { id: "PE-08", title: "Role & Persona Prompting", module: "Core Techniques", moduleNumber: 2, duration: "20 min", difficulty: "intermediate", description: "Shape AI responses through persona assignment.", summary: "Learn to assign roles and personas to AI for more targeted, expert-level responses.", objectives: ["Understand role-based prompting", "Create effective persona definitions", "Combine roles with task specifications", "Leverage domain expertise in prompts"] },
  { id: "PE-09", title: "ReAct Pattern", module: "Advanced Patterns", moduleNumber: 3, duration: "30 min", difficulty: "advanced", description: "Combine reasoning and action for agentic behavior.", summary: "Learn the ReAct framework that enables AI to reason about actions and act on reasoning.", objectives: ["Understand the ReAct paradigm", "Learn reasoning-action loops", "Implement tool use with ReAct", "Build simple agentic workflows"] },
  { id: "PE-10", title: "Self-Consistency", module: "Advanced Patterns", moduleNumber: 3, duration: "25 min", difficulty: "advanced", description: "Improve reliability through multiple reasoning paths.", summary: "Generate multiple responses and aggregate them for more reliable outputs.", objectives: ["Understand self-consistency approach", "Generate multiple reasoning chains", "Aggregate responses effectively", "Know when self-consistency helps"] },
  { id: "PE-11", title: "Tree of Thoughts", module: "Advanced Patterns", moduleNumber: 3, duration: "30 min", difficulty: "advanced", description: "Explore multiple solution paths systematically.", summary: "Advanced reasoning technique that explores branching possibilities like a decision tree.", objectives: ["Understand tree-structured reasoning", "Generate and evaluate branches", "Implement search strategies", "Apply ToT to complex problems"] },
  { id: "PE-12", title: "Prompt Chaining", module: "Advanced Patterns", moduleNumber: 3, duration: "25 min", difficulty: "advanced", description: "Break complex tasks into sequential prompt steps.", summary: "Learn to decompose complex tasks into a series of simpler prompts for better results.", objectives: ["Design prompt chain architectures", "Handle intermediate outputs", "Manage chain state and context", "Debug and optimize chains"] },
  { id: "PE-13", title: "Code Generation Prompts", module: "Task-Specific", moduleNumber: 4, duration: "30 min", difficulty: "intermediate", description: "Craft effective prompts for code generation tasks.", summary: "Best practices for prompting AI to write, explain, debug, and refactor code.", objectives: ["Structure code generation prompts", "Specify language and framework context", "Request code explanations", "Debug and improve generated code"] },
  { id: "PE-14", title: "Data Analysis Prompts", module: "Task-Specific", moduleNumber: 4, duration: "25 min", difficulty: "intermediate", description: "Use prompts for data analysis and insights extraction.", summary: "Techniques for leveraging AI in data analysis, pattern recognition, and insight generation.", objectives: ["Structure data analysis prompts", "Extract patterns from datasets", "Generate analysis reports", "Handle data formats effectively"] },
  { id: "PE-15", title: "Creative Writing Prompts", module: "Task-Specific", moduleNumber: 4, duration: "20 min", difficulty: "intermediate", description: "Unlock creative potential through strategic prompting.", summary: "Learn to prompt for stories, marketing copy, creative content, and artistic text.", objectives: ["Set creative constraints effectively", "Define style and tone", "Iterate on creative outputs", "Balance creativity with consistency"] },
  { id: "PE-16", title: "Agentic Prompt Patterns", module: "Task-Specific", moduleNumber: 4, duration: "35 min", difficulty: "advanced", description: "Design prompts for autonomous AI agent behavior.", summary: "Advanced patterns for creating AI agents that can plan, execute, and reflect on tasks.", objectives: ["Design agent system prompts", "Implement planning loops", "Enable tool selection", "Build reflection mechanisms"] },
  { id: "PE-17", title: "Prompt Evaluation", module: "Production", moduleNumber: 5, duration: "25 min", difficulty: "advanced", description: "Measure and evaluate prompt effectiveness systematically.", summary: "Learn evaluation metrics, testing strategies, and quality assessment for prompts.", objectives: ["Define evaluation criteria", "Create test cases for prompts", "Measure output quality", "Compare prompt variations"] },
  { id: "PE-18", title: "Prompt Optimization", module: "Production", moduleNumber: 5, duration: "30 min", difficulty: "advanced", description: "Iterate and improve prompts for better performance.", summary: "Techniques for refining prompts, reducing token usage, and improving output quality.", objectives: ["Identify optimization opportunities", "Reduce token consumption", "Improve response consistency", "A/B test prompt variations"] },
  { id: "PE-19", title: "Prompt Security", module: "Production", moduleNumber: 5, duration: "25 min", difficulty: "advanced", description: "Protect against prompt injection and misuse.", summary: "Security considerations including prompt injection, data leakage prevention, and safe guards.", objectives: ["Understand prompt injection attacks", "Implement input sanitization", "Design defensive prompts", "Prevent data leakage"] },
  { id: "PE-20", title: "Production Prompt Engineering", module: "Production", moduleNumber: 5, duration: "30 min", difficulty: "expert", description: "Deploy and manage prompts in production systems.", summary: "Best practices for versioning, monitoring, and maintaining prompts at scale.", objectives: ["Version control for prompts", "Monitor prompt performance", "Handle prompt updates safely", "Scale prompt management"] },
];

const moduleIcons: Record<number, React.ReactNode> = {
  1: <Lightbulb className="h-4 w-4" />,
  2: <Zap className="h-4 w-4" />,
  3: <Brain className="h-4 w-4" />,
  4: <Rocket className="h-4 w-4" />,
  5: <Settings className="h-4 w-4" />
};

const moduleColors: Record<number, string> = {
  1: "from-yellow-500/20 to-yellow-500/5 border-yellow-500/30",
  2: "from-blue-500/20 to-blue-500/5 border-blue-500/30",
  3: "from-purple-500/20 to-purple-500/5 border-purple-500/30",
  4: "from-green-500/20 to-green-500/5 border-green-500/30",
  5: "from-red-500/20 to-red-500/5 border-red-500/30"
};

interface Props {
  params: Promise<{ id: string }>;
}

export default function PromptEngineeringLessonClient({ params }: Props) {
  const { id: lessonId } = React.use(params);
  const router = useRouter();
  const [activeTab, setActiveTab] = useState("theory");
  const [lessonProgress, setLessonProgress] = useState<Record<string, boolean>>({});
  const [completedObjectives, setCompletedObjectives] = useState<Set<number>>(new Set());
  const [timeSpent, setTimeSpent] = useState(0);
  const [startTime] = useState(Date.now());

  const getActualLessonId = (id: string): string => {
    if (id.startsWith("PE-")) return id;
    const num = parseInt(id, 10);
    if (!isNaN(num) && num >= 1 && num <= 20) {
      return `PE-${num.toString().padStart(2, "0")}`;
    }
    return id;
  };

  const actualId = getActualLessonId(lessonId);
  const currentLesson = lessons.find(l => l.id === actualId);
  const currentIndex = lessons.findIndex(l => l.id === actualId);
  const prevLesson = currentIndex > 0 ? lessons[currentIndex - 1] : null;
  const nextLesson = currentIndex < lessons.length - 1 ? lessons[currentIndex + 1] : null;

  useEffect(() => {
    const savedProgress = localStorage.getItem("pe-lesson-progress");
    if (savedProgress) {
      setLessonProgress(JSON.parse(savedProgress));
    }
    const savedObjectives = localStorage.getItem(`pe-lesson-${lessonId}-objectives`);
    if (savedObjectives) {
      setCompletedObjectives(new Set(JSON.parse(savedObjectives)));
    }
  }, [lessonId]);

  useEffect(() => {
    const interval = setInterval(() => {
      setTimeSpent(Math.floor((Date.now() - startTime) / 1000));
    }, 1000);
    return () => clearInterval(interval);
  }, [startTime]);

  const handleObjectiveToggle = (index: number) => {
    const newCompleted = new Set(completedObjectives);
    if (newCompleted.has(index)) {
      newCompleted.delete(index);
    } else {
      newCompleted.add(index);
    }
    setCompletedObjectives(newCompleted);
    localStorage.setItem(`pe-lesson-${lessonId}-objectives`, JSON.stringify([...newCompleted]));
  };

  const handleCompleteLesson = () => {
    const newProgress = { ...lessonProgress, [lessonId]: true };
    setLessonProgress(newProgress);
    localStorage.setItem("pe-lesson-progress", JSON.stringify(newProgress));
    const allObjectives = new Set(currentLesson?.objectives.map((_, i) => i) || []);
    setCompletedObjectives(allObjectives);
    localStorage.setItem(`pe-lesson-${lessonId}-objectives`, JSON.stringify([...allObjectives]));
  };

  const completedCount = Object.values(lessonProgress).filter(Boolean).length;
  const overallProgress = (completedCount / lessons.length) * 100;
  const lessonCompletion = currentLesson
    ? (completedObjectives.size / currentLesson.objectives.length) * 100
    : 0;

  if (!currentLesson) {
    return (
      <div className="min-h-screen bg-gradient-to-br from-background via-background to-accent/5">
        <div className="container mx-auto px-4 py-8">
          <div className="text-center py-16">
            <h1 className="text-3xl font-bold mb-4">Lesson Not Found</h1>
            <p className="text-muted-foreground mb-6">The lesson you&apos;re looking for doesn&apos;t exist.</p>
            <Link href="/learn/prompt-engineering">
              <Button size="lg"><ArrowLeft className="mr-2 h-4 w-4" />Back to Course</Button>
            </Link>
          </div>
        </div>
      </div>
    );
  }

  const difficultyColors: Record<string, string> = {
    beginner: "bg-green-500/10 text-green-600 border-green-500/20",
    intermediate: "bg-yellow-500/10 text-yellow-600 border-yellow-500/20",
    advanced: "bg-orange-500/10 text-orange-600 border-orange-500/20",
    expert: "bg-red-500/10 text-red-600 border-red-500/20"
  };

  return (
    <div className="min-h-screen bg-gradient-to-br from-background via-background to-accent/5">
      <div className="container mx-auto px-4 py-8">
        <div className="mb-6">
          <div className="flex items-center justify-between mb-4">
            <div className="flex items-center space-x-2 text-sm text-muted-foreground">
              <Link href="/" className="hover:text-primary transition-colors"><Home className="h-4 w-4" /></Link>
              <ChevronRight className="h-4 w-4" />
              <Link href="/learn" className="hover:text-primary transition-colors">Learn</Link>
              <ChevronRight className="h-4 w-4" />
              <Link href="/learn/prompt-engineering" className="hover:text-primary transition-colors">Prompt Engineering</Link>
              <ChevronRight className="h-4 w-4" />
              <span className="text-foreground font-medium">{currentLesson.id}</span>
            </div>
            <Badge variant="outline" className="flex items-center gap-1">
              <Sparkles className="h-3 w-3" />
              {completedCount}/{lessons.length} Complete
            </Badge>
          </div>
          <Card className="bg-card/50 backdrop-blur-sm border-border/50">
            <CardContent className="p-4">
              <div className="flex items-center justify-between mb-2">
                <div className="flex items-center space-x-2">
                  <Target className="h-4 w-4 text-primary" />
                  <span className="text-sm font-medium">Course Progress</span>
                </div>
                <span className="text-sm text-muted-foreground">{Math.round(overallProgress)}%</span>
              </div>
              <Progress value={overallProgress} className="h-2" />
            </CardContent>
          </Card>
        </div>

        <Card className={`mb-8 p-6 lg:p-8 border-2 relative overflow-hidden bg-gradient-to-br ${moduleColors[currentLesson.moduleNumber]}`}>
          <div className="absolute top-0 right-0 w-96 h-96 bg-gradient-to-br from-primary/20 to-transparent rounded-full blur-3xl -z-10" />
          <div className="absolute bottom-0 left-0 w-64 h-64 bg-gradient-to-tr from-secondary/15 to-transparent rounded-full blur-2xl -z-10" />
          <div className="flex items-start justify-between mb-6">
            <div className="flex-1">
              <div className="flex items-center flex-wrap gap-2 mb-4">
                <Badge variant="outline" className="flex items-center gap-1">{moduleIcons[currentLesson.moduleNumber]}{currentLesson.id}</Badge>
                <Badge variant="secondary" className="flex items-center gap-1">{moduleIcons[currentLesson.moduleNumber]}Module {currentLesson.moduleNumber}: {currentLesson.module}</Badge>
                <Badge className={difficultyColors[currentLesson.difficulty]}>{currentLesson.difficulty.charAt(0).toUpperCase() + currentLesson.difficulty.slice(1)}</Badge>
                {lessonProgress[lessonId] && <Badge className="bg-green-500/10 text-green-600 border-green-500/20"><CheckCircle className="mr-1 h-3 w-3" />Completed</Badge>}
              </div>
              <h1 className="text-3xl lg:text-4xl font-bold mb-4 bg-gradient-to-r from-primary to-blue-600 bg-clip-text text-transparent">{currentLesson.title}</h1>
              <p className="text-muted-foreground text-base lg:text-lg mb-6 leading-relaxed max-w-4xl">{currentLesson.description}</p>
              <div className="grid grid-cols-2 lg:grid-cols-4 gap-4">
                <div className="flex items-center space-x-2 bg-background/50 px-4 py-3 rounded-xl border border-border/50"><Clock className="h-5 w-5 text-primary" /><div><div className="text-sm font-medium">{currentLesson.duration}</div><div className="text-xs text-muted-foreground">Duration</div></div></div>
                <div className="flex items-center space-x-2 bg-background/50 px-4 py-3 rounded-xl border border-border/50"><Target className="h-5 w-5 text-blue-500" /><div><div className="text-sm font-medium">{completedObjectives.size}/{currentLesson.objectives.length}</div><div className="text-xs text-muted-foreground">Objectives</div></div></div>
                <div className="flex items-center space-x-2 bg-background/50 px-4 py-3 rounded-xl border border-border/50"><BookOpen className="h-5 w-5 text-green-500" /><div><div className="text-sm font-medium">{Math.round(lessonCompletion)}%</div><div className="text-xs text-muted-foreground">Complete</div></div></div>
                <div className="flex items-center space-x-2 bg-background/50 px-4 py-3 rounded-xl border border-border/50"><Clock className="h-5 w-5 text-orange-500" /><div><div className="text-sm font-medium">{Math.floor(timeSpent / 60)}m</div><div className="text-xs text-muted-foreground">Time Spent</div></div></div>
              </div>
            </div>
          </div>
        </Card>

        <div className="flex gap-6">
          <div className="hidden lg:block w-80 flex-shrink-0">
            <div className="sticky top-6 space-y-6">
              <Card>
                <CardHeader className="pb-3"><CardTitle className="text-lg flex items-center"><Target className="mr-2 h-5 w-5 text-primary" />Lesson Progress</CardTitle></CardHeader>
                <CardContent className="space-y-4">
                  <div className="space-y-2">
                    <div className="flex items-center justify-between text-sm"><span className="text-muted-foreground">Completion</span><span className="font-medium">{Math.round(lessonCompletion)}%</span></div>
                    <Progress value={lessonCompletion} className="h-2" />
                  </div>
                  <Button onClick={handleCompleteLesson} disabled={lessonProgress[lessonId]} className="w-full" variant={lessonProgress[lessonId] ? "outline" : "default"}>
                    {lessonProgress[lessonId] ? <><CheckCircle className="mr-2 h-4 w-4" />Completed</> : <><CheckCircle className="mr-2 h-4 w-4" />Mark as Complete</>}
                  </Button>
                </CardContent>
              </Card>
              <Card>
                <CardHeader className="pb-3"><CardTitle className="text-lg flex items-center"><Target className="mr-2 h-5 w-5 text-primary" />Objectives</CardTitle></CardHeader>
                <CardContent>
                  <div className="space-y-2">
                    {currentLesson.objectives.map((objective, index) => (
                      <div key={index} className="flex items-start space-x-3 p-2 rounded-lg hover:bg-muted/50 cursor-pointer transition-colors" onClick={() => handleObjectiveToggle(index)}>
                        <CheckCircle className={`h-4 w-4 mt-0.5 flex-shrink-0 ${completedObjectives.has(index) ? "text-green-500" : "text-muted-foreground"}`} />
                        <span className={`text-sm leading-relaxed ${completedObjectives.has(index) ? "line-through text-muted-foreground" : ""}`}>{objective}</span>
                      </div>
                    ))}
                  </div>
                </CardContent>
              </Card>
            </div>
          </div>

          <div className="flex-1 min-w-0">
            <Tabs value={activeTab} onValueChange={setActiveTab} className="space-y-6">
              <TabsList className="grid w-full grid-cols-3 lg:w-[500px]">
                <TabsTrigger value="theory"><Lightbulb className="mr-2 h-4 w-4" />Theory</TabsTrigger>
                <TabsTrigger value="examples"><Code2 className="mr-2 h-4 w-4" />Examples</TabsTrigger>
                <TabsTrigger value="exercises"><Trophy className="mr-2 h-4 w-4" />Exercises</TabsTrigger>
              </TabsList>
              <TabsContent value="theory" className="space-y-6">
                <Card className="border-2">
                  <CardHeader className="bg-gradient-to-r from-primary/10 to-primary/5 border-b"><CardTitle className="flex items-center text-xl"><div className="p-2 rounded-lg bg-primary/10 mr-3"><FileText className="h-5 w-5 text-primary" /></div>Overview</CardTitle></CardHeader>
                  <CardContent className="pt-6"><p className="text-muted-foreground leading-relaxed">{currentLesson.summary}</p></CardContent>
                </Card>
                <Card className="border-2">
                  <CardHeader className="bg-gradient-to-r from-primary/10 to-primary/5 border-b"><CardTitle className="flex items-center text-xl"><div className="p-2 rounded-lg bg-primary/10 mr-3"><Target className="h-5 w-5 text-primary" /></div>Learning Objectives</CardTitle></CardHeader>
                  <CardContent className="pt-6">
                    <ul className="space-y-3">
                      {currentLesson.objectives.map((objective, index) => (
                        <li key={index} className="flex items-start space-x-3 cursor-pointer" onClick={() => handleObjectiveToggle(index)}>
                          <CheckCircle className={`h-4 w-4 mt-0.5 flex-shrink-0 ${completedObjectives.has(index) ? "text-green-500" : "text-muted-foreground"}`} />
                          <span className={`text-sm leading-relaxed ${completedObjectives.has(index) ? "line-through text-muted-foreground" : ""}`}>{objective}</span>
                        </li>
                      ))}
                    </ul>
                  </CardContent>
                </Card>
              </TabsContent>
              <TabsContent value="examples" className="space-y-6">
                <Card className="border-2">
                  <CardHeader className="bg-gradient-to-r from-blue-500/10 to-blue-500/5 border-b border-blue-500/20"><CardTitle className="flex items-center text-xl"><div className="p-2 rounded-lg bg-blue-500/10 mr-3"><Code2 className="h-5 w-5 text-blue-500" /></div>Prompt Examples</CardTitle></CardHeader>
                  <CardContent className="pt-6"><p className="text-muted-foreground">Example prompts for {currentLesson.title.toLowerCase()} will be available in the full course.</p></CardContent>
                </Card>
              </TabsContent>
              <TabsContent value="exercises" className="space-y-6">
                <Card className="border-2">
                  <CardHeader className="bg-gradient-to-r from-green-500/10 to-green-500/5 border-b border-green-500/20"><CardTitle className="flex items-center text-xl"><div className="p-2 rounded-lg bg-green-500/10 mr-3"><Trophy className="h-5 w-5 text-green-500" /></div>Practice Exercises</CardTitle></CardHeader>
                  <CardContent className="pt-6"><p className="text-muted-foreground">Practice exercises for {currentLesson.title.toLowerCase()} will help reinforce your learning.</p></CardContent>
                </Card>
              </TabsContent>
            </Tabs>

            <Card className="p-6 mt-12 border-2">
              <div className="flex flex-col sm:flex-row items-center justify-between gap-4">
                <div className="w-full sm:w-auto">
                  {prevLesson && <Link href={`/learn/prompt-engineering/${prevLesson.id}`} className="block w-full sm:w-auto"><Button variant="outline" size="lg" className="w-full sm:w-auto"><ArrowLeft className="mr-2 h-4 w-4" /><div className="flex flex-col items-start"><span className="text-xs text-muted-foreground">{prevLesson.id}</span><span className="text-sm">Previous</span></div></Button></Link>}
                </div>
                <Link href="/learn/prompt-engineering" className="w-full sm:w-auto"><Button variant="outline" size="lg" className="w-full sm:w-auto">All Lessons</Button></Link>
                <div className="w-full sm:w-auto">
                  {nextLesson && <Link href={`/learn/prompt-engineering/${nextLesson.id}`} className="block w-full sm:w-auto"><Button size="lg" className="w-full sm:w-auto"><div className="flex flex-col items-end"><span className="text-xs opacity-70">{nextLesson.id}</span><span className="text-sm">Next</span></div><ArrowRight className="ml-2 h-4 w-4" /></Button></Link>}
                </div>
              </div>
            </Card>
          </div>
        </div>
      </div>
    </div>
  );
}
