export type Difficulty = 'Easy' | 'Medium' | 'Hard';
export type ProblemStatus = 'pending' | 'in_progress' | 'completed';

export interface Problem {
  id: string;
  number: number;
  title: string;
  difficulty: Difficulty;
  category: string;
  status: ProblemStatus;
  dateCompleted?: string;
  attempts: number;
  timeSpent: number; // in minutes
  notes: string;
  patterns: string[];
  reviewCount: number;
  nextReview?: string;
  leetcodeNumber?: number;
}

export interface Category {
  id: string;
  name: string;
  icon: string;
  totalProblems: number;
  completedProblems: number;
  inProgressProblems: number;
  color: string;
}

export interface Session {
  id: string;
  date: string;
  sessionNumber: number;
  topic: string;
  problemsSolved: number;
  timeSpent: number; // in minutes
  notes: string;
  problems: string[];
}

export interface SessionNote {
  id: string;
  sessionId: string;
  problemId: string;
  approach: string;
  timeComplexity: string;
  spaceComplexity: string;
  keyInsight: string;
  struggles: string[];
  code?: string;
}

export interface Pattern {
  id: string;
  name: string;
  description: string;
  template: string;
  problems: string[];
  whenToUse: string[];
  confidence: number; // 1-5
}

export interface AlgoProgress {
  totalProblems: number;
  completedProblems: number;
  inProgressProblems: number;
  currentStreak: number;
  longestStreak: number;
  totalTimeSpent: number; // in minutes
  averageTimePerProblem: number;
  easySolved: number;
  mediumSolved: number;
  hardSolved: number;
  lastSolvedDate?: string;
}

export interface ReviewItem {
  problemId: string;
  problemTitle: string;
  category: string;
  difficulty: Difficulty;
  lastReview: string;
  interval: number;
  due: boolean;
}

export const CATEGORIES: Category[] = [
  { id: 'arrays', name: 'Arrays & Strings', icon: 'Array', totalProblems: 25, completedProblems: 0, inProgressProblems: 0, color: 'blue' },
  { id: 'linkedlists', name: 'Linked Lists', icon: 'Link', totalProblems: 15, completedProblems: 0, inProgressProblems: 0, color: 'purple' },
  { id: 'stacks', name: 'Stacks & Queues', icon: 'Layers', totalProblems: 12, completedProblems: 0, inProgressProblems: 0, color: 'green' },
  { id: 'trees', name: 'Trees & BST', icon: 'GitBranch', totalProblems: 20, completedProblems: 0, inProgressProblems: 0, color: 'emerald' },
  { id: 'heaps', name: 'Heaps & Priority Q', icon: 'ArrowUpCircle', totalProblems: 10, completedProblems: 0, inProgressProblems: 0, color: 'orange' },
  { id: 'hash', name: 'Hash Tables', icon: 'Hash', totalProblems: 12, completedProblems: 0, inProgressProblems: 0, color: 'cyan' },
  { id: 'graphs', name: 'Graphs', icon: 'Network', totalProblems: 25, completedProblems: 0, inProgressProblems: 0, color: 'rose' },
  { id: 'dp', name: 'Dynamic Programming', icon: 'Zap', totalProblems: 30, completedProblems: 0, inProgressProblems: 0, color: 'yellow' },
  { id: 'backtracking', name: 'Backtracking', icon: 'RotateCcw', totalProblems: 15, completedProblems: 0, inProgressProblems: 0, color: 'pink' },
  { id: 'sorting', name: 'Sorting & Searching', icon: 'ArrowUpDown', totalProblems: 18, completedProblems: 0, inProgressProblems: 0, color: 'indigo' },
];

export const DIFFICULTY_COLORS: Record<Difficulty, string> = {
  Easy: 'bg-green-100 text-green-800 border-green-200',
  Medium: 'bg-yellow-100 text-yellow-800 border-yellow-200',
  Hard: 'bg-red-100 text-red-800 border-red-200',
};

export const STATUS_CONFIG: Record<ProblemStatus, { label: string; color: string; icon: string }> = {
  pending: { label: 'Not Started', color: 'bg-gray-100 text-gray-600', icon: 'Circle' },
  in_progress: { label: 'In Progress', color: 'bg-blue-100 text-blue-600', icon: 'LoaderCircle' },
  completed: { label: 'Completed', color: 'bg-green-100 text-green-600', icon: 'CheckCircle' },
};
