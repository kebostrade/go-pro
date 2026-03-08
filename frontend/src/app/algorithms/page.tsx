"use client";

import AlgorithmsTracker from "@/components/algorithms/algorithms-tracker";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import ProblemsList from "@/components/algorithms/problems-list";
import SessionLog from "@/components/algorithms/session-log";
import { AlgoProgress, Problem, Category, Session } from "@/types/algorithms";

// Sample data
const sampleCategories: Category[] = [
  { id: 'arrays', name: 'Arrays & Strings', icon: 'Array', totalProblems: 25, completedProblems: 5, inProgressProblems: 2, color: 'blue' },
  { id: 'linkedlists', name: 'Linked Lists', icon: 'Link', totalProblems: 15, completedProblems: 2, inProgressProblems: 1, color: 'purple' },
  { id: 'stacks', name: 'Stacks & Queues', icon: 'Layers', totalProblems: 12, completedProblems: 3, inProgressProblems: 1, color: 'green' },
  { id: 'trees', name: 'Trees & BST', icon: 'GitBranch', totalProblems: 20, completedProblems: 4, inProgressProblems: 2, color: 'emerald' },
  { id: 'dp', name: 'Dynamic Programming', icon: 'Zap', totalProblems: 30, completedProblems: 6, inProgressProblems: 3, color: 'yellow' },
  { id: 'graphs', name: 'Graphs', icon: 'Network', totalProblems: 25, completedProblems: 3, inProgressProblems: 2, color: 'rose' },
];

const sampleProblems: Problem[] = [
  { id: '1', number: 1, title: 'Two Sum', difficulty: 'Easy', category: 'arrays', status: 'completed', dateCompleted: '2024-01-15', attempts: 1, timeSpent: 15, notes: 'Use hash map for O(n) complexity', patterns: ['Hash Map'], reviewCount: 2, nextReview: '2026-01-22', leetcodeNumber: 1 },
  { id: '2', number: 2, title: 'Contains Duplicate', difficulty: 'Easy', category: 'arrays', status: 'completed', dateCompleted: '2026-01-16', attempts: 2, timeSpent: 12, notes: 'Use set for duplicate detection', patterns: ['Hash Set'], reviewCount: 1, nextReview: '2026-01-23', leetcodeNumber: 217 },
  { id: '3', number: 3, title: 'Best Time to Buy/Sell Stock', difficulty: 'Easy', category: 'arrays', status: 'completed', dateCompleted: '2026-01-17', attempts: 1, timeSpent: 18, notes: 'Track min price, calculate max profit at each step', patterns: ['Single Pass'], reviewCount: 0, nextReview: '2026-01-20', leetcodeNumber: 121 },
  { id: '4', number: 4, title: 'Product of Array Except Self', difficulty: 'Medium', category: 'arrays', status: 'in_progress', attempts: 1, timeSpent: 0, notes: '', patterns: [], reviewCount: 0 },
  { id: '5', number: 5, title: 'Maximum Subarray', difficulty: 'Medium', category: 'arrays', status: 'pending', attempts: 0, timeSpent: 0, notes: '', patterns: [], reviewCount: 0 },
  { id: '6', number: 6, title: '3Sum', difficulty: 'Medium', category: 'arrays', status: 'pending', attempts: 0, timeSpent: 0, notes: '', patterns: [], reviewCount: 0 },
  { id: '7', number: 7, title: 'Reverse Linked List', difficulty: 'Easy', category: 'linkedlists', status: 'pending', attempts: 0, timeSpent: 0, notes: '', patterns: [], reviewCount: 0 },
  { id: '8', number: 8, title: 'Merge Two Sorted Lists', difficulty: 'Easy', category: 'linkedlists', status: 'pending', attempts: 0, timeSpent: 0, notes: '', patterns: [], reviewCount: 0 },
  { id: '9', number: 9, title: 'Linked List Cycle', difficulty: 'Easy', category: 'linkedlists', status: 'pending', attempts: 0, timeSpent: 0, notes: '', patterns: [], reviewCount: 0 },
  { id: '10', number: 10, title: 'Valid Parentheses', difficulty: 'Easy', category: 'stacks', status: 'pending', attempts: 0, timeSpent: 0, notes: '', patterns: [], reviewCount: 0 },
  { id: '11', number: 11, title: 'Min Stack', difficulty: 'Medium', category: 'stacks', status: 'pending', attempts: 0, timeSpent: 0, notes: '', patterns: [], reviewCount: 0 },
  { id: '12', number: 12, title: 'Maximum Depth of Binary Tree', difficulty: 'Easy', category: 'trees', status: 'pending', attempts: 0, timeSpent: 0, notes: '', patterns: [], reviewCount: 0 },
  { id: '13', number: 13, title: 'Number of Islands', difficulty: 'Medium', category: 'graphs', status: 'pending', attempts: 0, timeSpent: 0, notes: '', patterns: [], reviewCount: 0 },
  { id: '14', number: 14, title: 'Climbing Stairs', difficulty: 'Easy', category: 'dp', status: 'pending', attempts: 0, timeSpent: 0, notes: '', patterns: [], reviewCount: 0 },
  { id: '15', number: 15, title: 'Coin Change', difficulty: 'Medium', category: 'dp', status: 'pending', attempts: 0, timeSpent: 0, notes: '', patterns: [], reviewCount: 0 },
];

const sampleSessions: Session[] = [
  { id: '1', date: '2026-01-15', sessionNumber: 1, topic: 'Arrays & Hash Maps', problemsSolved: 3, timeSpent: 45, notes: 'Focus on Two Sum, Contains Duplicate patterns', problems: ['1', '2', '3'] },
  { id: '2', date: '2026-01-16', sessionNumber: 2, topic: 'Arrays & Two Pointers', problemsSolved: 2, timeSpent: 32, notes: 'Two Sum II (optimized) - already did easier!', problems: ['4', '5'] },
];

const sampleProgress: AlgoProgress = {
  totalProblems: 182,
  completedProblems: 3,
  inProgressProblems: 1,
  currentStreak: 2,
  longestStreak: 5,
  totalTimeSpent: 73, // minutes
  averageTimePerProblem: 24,
  easySolved: 1,
  mediumSolved: 1,
  hardSolved: 1,
  lastSolvedDate: '2026-01-17',
};

export default function AlgorithmsPage() {
  const progress = sampleProgress;
  const problems = sampleProblems
  const sessions = sampleSessions
  const categories = sampleCategories

  return (
    <div className="container mx-auto py-8 px-4 md:px-8 lg:px-12">
      <h1 className="text-3xl font-bold mb-6">Algorithms Practice</h1>
      <p className="text-muted-foreground mb-4">
        Track your progress through structured practice sessions
      </p>

      <div className="mb-6">
        <Tabs defaultValue="tracker">
          <TabsList>
            <TabsTrigger value="tracker">Tracker</TabsTrigger>
            <TabsTrigger value="sessions">Sessions</TabsTrigger>
            <TabsTrigger value="problems">Problems</TabsTrigger>
          </TabsList>

          <TabsContent value="tracker">
            <AlgorithmsTracker
              categories={categories}
              progress={progress}
              sessions={sessions}
              problems={problems}
            />
          </TabsContent>

          <TabsContent value="sessions">
            <SessionLog sessions={sessions} />
          </TabsContent>

          <TabsContent value="problems">
            <ProblemsList problems={problems} selectedCategory={null} categories={categories} />
          </TabsContent>
        </Tabs>
      </div>
    </div>
  );
}
