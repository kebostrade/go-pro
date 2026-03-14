"use client";

import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Progress } from "@/components/ui/progress";
import { ArrowLeft, Trophy, Target, Star, CheckCircle } from "lucide-react";
import Link from "next/link";

interface Exercise {
  id: string;
  title: string;
  difficulty: "Easy" | "Medium" | "Hard";
  description: string;
}

interface Module {
  id: number;
  title: string;
  exercises: Exercise[];
  completedCount: number;
}

const modules: Module[] = [
  {
    id: 1,
    title: "Foundations",
    completedCount: 0,
    exercises: [
      { id: "1.1", title: "Improve Prompts", difficulty: "Easy", description: "Rewrite prompts to be more effective" },
      { id: "1.2", title: "Token Estimation", difficulty: "Medium", description: "Estimate tokens for various content types" },
      { id: "1.3", title: "PIERS Analysis", difficulty: "Medium", description: "Analyze a prompt using PIERS framework" },
    ],
  },
  {
    id: 2,
    title: "Core Techniques",
    completedCount: 0,
    exercises: [
      { id: "2.1", title: "Zero-Shot Classification", difficulty: "Easy", description: "Classify support tickets without examples" },
      { id: "2.2", title: "Few-Shot Examples", difficulty: "Medium", description: "Create examples for voice conversion" },
      { id: "2.3", title: "Chain-of-Thought", difficulty: "Medium", description: "Write CoT prompt for logic puzzle" },
      { id: "2.4", title: "JSON Output", difficulty: "Medium", description: "Extract movie info as structured JSON" },
      { id: "2.5", title: "Role Design", difficulty: "Hard", description: "Design salary negotiation coach prompt" },
    ],
  },
  {
    id: 3,
    title: "Advanced Patterns",
    completedCount: 0,
    exercises: [
      { id: "3.1", title: "ReAct Trace", difficulty: "Medium", description: "Write ReAct trace for population query" },
      { id: "3.2", title: "Self-Consistency", difficulty: "Medium", description: "Implement majority voting in Python" },
      { id: "3.3", title: "Tree of Thoughts", difficulty: "Hard", description: "Create ToT trace for trip planning" },
      { id: "3.4", title: "Prompt Chain Design", difficulty: "Hard", description: "Design feedback analysis chain" },
    ],
  },
  {
    id: 4,
    title: "Task-Specific Prompting",
    completedCount: 0,
    exercises: [
      { id: "4.1", title: "Code Generation", difficulty: "Medium", description: "Generate Luhn algorithm validator" },
      { id: "4.2", title: "Data Extraction", difficulty: "Medium", description: "Extract structured data from resume" },
      { id: "4.3", title: "Email Sequence", difficulty: "Medium", description: "Create 3-email welcome sequence" },
      { id: "4.4", title: "Agent Tools", difficulty: "Hard", description: "Design customer support agent tools" },
    ],
  },
  {
    id: 5,
    title: "Production & Optimization",
    completedCount: 0,
    exercises: [
      { id: "5.1", title: "Test Suite", difficulty: "Medium", description: "Create test suite for summarization" },
      { id: "5.2", title: "Evaluation Pipeline", difficulty: "Medium", description: "Build automated evaluation in Python" },
      { id: "5.3", title: "Prompt Registry", difficulty: "Hard", description: "Design versioned prompt registry" },
      { id: "5.4", title: "Security Review", difficulty: "Hard", description: "Review and fix security vulnerabilities" },
    ],
  },
];

const challenges: Exercise[] = [
  { id: "c1", title: "Multi-Format Converter", difficulty: "Hard", description: "Build chain converting Markdown/HTML/JSON/YAML/CSV" },
  { id: "c2", title: "Code Review Agent", difficulty: "Hard", description: "Design agent reading PRs and posting reviews" },
  { id: "c3", title: "Research Assistant", difficulty: "Hard", description: "Build research agent with citations" },
  { id: "c4", title: "Content Pipeline", difficulty: "Hard", description: "Create end-to-end content generation pipeline" },
];

const difficultyColors = {
  Easy: "bg-green-100 text-green-800 border-green-200",
  Medium: "bg-yellow-100 text-yellow-800 border-yellow-200",
  Hard: "bg-red-100 text-red-800 border-red-200",
};

export default function ExercisesPage() {
  const totalExercises = modules.reduce((sum, m) => sum + m.exercises.length, 0);
  const totalChallenges = challenges.length;
  const completedExercises = modules.reduce((sum, m) => sum + m.completedCount, 0);
  const overallProgress = Math.round((completedExercises / totalExercises) * 100);

  return (
    <div className="min-h-screen bg-gray-50">
      {/* Header */}
      <div className="bg-white border-b border-gray-200">
        <div className="max-w-7xl mx-auto px-4 py-6">
          <Link
            href="/learn/prompt-engineering"
            className="inline-flex items-center gap-2 text-sm text-gray-600 hover:text-gray-900 mb-4"
          >
            <ArrowLeft className="w-4 h-4" />
            Back to Course
          </Link>
          <h1 className="text-3xl font-bold text-gray-900">Prompt Engineering Exercises</h1>
          <p className="text-gray-600 mt-2">Practice exercises organized by module and difficulty</p>
        </div>
      </div>

      <div className="max-w-7xl mx-auto px-4 py-8">
        {/* Stats Overview */}
        <div className="grid grid-cols-1 md:grid-cols-4 gap-4 mb-8">
          <Card>
            <CardContent className="pt-6">
              <div className="flex items-center gap-3">
                <div className="p-2 bg-blue-100 rounded-lg">
                  <Target className="w-5 h-5 text-blue-600" />
                </div>
                <div>
                  <p className="text-sm text-gray-600">Total Exercises</p>
                  <p className="text-2xl font-bold text-gray-900">{totalExercises}</p>
                </div>
              </div>
            </CardContent>
          </Card>
          <Card>
            <CardContent className="pt-6">
              <div className="flex items-center gap-3">
                <div className="p-2 bg-purple-100 rounded-lg">
                  <Trophy className="w-5 h-5 text-purple-600" />
                </div>
                <div>
                  <p className="text-sm text-gray-600">Challenges</p>
                  <p className="text-2xl font-bold text-gray-900">{totalChallenges}</p>
                </div>
              </div>
            </CardContent>
          </Card>
          <Card>
            <CardContent className="pt-6">
              <div className="flex items-center gap-3">
                <div className="p-2 bg-green-100 rounded-lg">
                  <CheckCircle className="w-5 h-5 text-green-600" />
                </div>
                <div>
                  <p className="text-sm text-gray-600">Completed</p>
                  <p className="text-2xl font-bold text-gray-900">{completedExercises}</p>
                </div>
              </div>
            </CardContent>
          </Card>
          <Card>
            <CardContent className="pt-6">
              <div className="flex items-center gap-3">
                <div className="p-2 bg-amber-100 rounded-lg">
                  <Star className="w-5 h-5 text-amber-600" />
                </div>
                <div>
                  <p className="text-sm text-gray-600">Progress</p>
                  <p className="text-2xl font-bold text-gray-900">{overallProgress}%</p>
                </div>
              </div>
            </CardContent>
          </Card>
        </div>

        {/* Modules */}
        <div className="space-y-6">
          {modules.map((module) => (
            <Card key={module.id}>
              <CardHeader className="pb-4">
                <div className="flex items-center justify-between">
                  <CardTitle className="text-xl">
                    Module {module.id}: {module.title}
                  </CardTitle>
                  <div className="flex items-center gap-3">
                    <span className="text-sm text-gray-600">
                      {module.completedCount}/{module.exercises.length} completed
                    </span>
                    <Progress
                      value={(module.completedCount / module.exercises.length) * 100}
                      className="w-24 h-2"
                    />
                  </div>
                </div>
              </CardHeader>
              <CardContent>
                <div className="space-y-3">
                  {module.exercises.map((exercise) => (
                    <div
                      key={exercise.id}
                      className="flex items-center justify-between p-4 bg-gray-50 rounded-lg hover:bg-gray-100 transition-colors"
                    >
                      <div className="flex items-center gap-4">
                        <Badge variant="outline" className={difficultyColors[exercise.difficulty]}>
                          {exercise.difficulty}
                        </Badge>
                        <div>
                          <h3 className="font-medium text-gray-900">
                            {module.id}.{exercise.id.split(".")[1]} {exercise.title}
                          </h3>
                          <p className="text-sm text-gray-600">{exercise.description}</p>
                        </div>
                      </div>
                      <Button variant="outline" size="sm">
                        Start
                      </Button>
                    </div>
                  ))}
                </div>
              </CardContent>
            </Card>
          ))}
        </div>

        {/* Challenges */}
        <Card className="mt-8 border-purple-200">
          <CardHeader className="bg-purple-50 rounded-t-lg">
            <div className="flex items-center gap-2">
              <Trophy className="w-5 h-5 text-purple-600" />
              <CardTitle className="text-xl text-purple-900">Challenge Problems</CardTitle>
            </div>
            <p className="text-sm text-purple-700 mt-1">
              Advanced multi-step challenges to test your mastery
            </p>
          </CardHeader>
          <CardContent className="pt-6">
            <div className="space-y-3">
              {challenges.map((challenge) => (
                <div
                  key={challenge.id}
                  className="flex items-center justify-between p-4 bg-purple-50 rounded-lg hover:bg-purple-100 transition-colors border border-purple-200"
                >
                  <div className="flex items-center gap-4">
                    <Badge variant="outline" className={difficultyColors[challenge.difficulty]}>
                      {challenge.difficulty}
                    </Badge>
                    <div>
                      <h3 className="font-medium text-gray-900">{challenge.title}</h3>
                      <p className="text-sm text-gray-600">{challenge.description}</p>
                    </div>
                  </div>
                  <Button variant="outline" size="sm" className="border-purple-300 text-purple-700">
                    Attempt
                  </Button>
                </div>
              ))}
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  );
}
