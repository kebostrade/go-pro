// API Response Types
export interface ApiResponse<T> {
  success: boolean;
  data?: T;
  error?: string;
  message?: string;
}

export interface PaginatedResponse<T> {
  items: T[];
  total: number;
  page: number;
  limit: number;
  hasNext: boolean;
  hasPrev: boolean;
}

// User Types
export interface User {
  id: string;
  name: string;
  email: string;
  avatar: string;
  badge: string;
  reputation: number;
  joinedAt: string;
  location?: string;
  bio?: string;
  website?: string;
  github?: string;
  twitter?: string;
  specialties: string[];
  stats: {
    posts: number;
    answers: number;
    helpfulAnswers: number;
    followers: number;
    following: number;
    challengesCompleted: number;
    projectsCompleted: number;
    streak: number;
    points: number;
  };
}

// Practice Types
export interface Challenge {
  id: string;
  title: string;
  description: string;
  difficulty: "Beginner" | "Intermediate" | "Advanced";
  category: string;
  estimatedTime: string;
  points: number;
  tags: string[];
  instructions: string;
  hints: string[];
  starterCode: string;
  solution: string;
  testCases: TestCase[];
  completionRate: number;
  totalAttempts: number;
  createdAt: string;
  updatedAt: string;
}

export interface TestCase {
  id: string;
  input: string;
  expectedOutput: string;
  description: string;
  passed?: boolean;
}

export interface ChallengeSubmission {
  id: string;
  challengeId: string;
  userId: string;
  code: string;
  language: string;
  status: "pending" | "running" | "completed" | "failed";
  results?: TestResult[];
  score?: number;
  completedAt?: string;
  submittedAt: string;
}

export interface TestResult {
  testCaseId: string;
  passed: boolean;
  actualOutput?: string;
  error?: string;
  executionTime?: number;
}

export interface Assessment {
  id: string;
  title: string;
  description: string;
  category: string;
  difficulty: "Beginner" | "Intermediate" | "Advanced";
  questions: AssessmentQuestion[];
  duration: number; // in minutes
  maxScore: number;
  passingScore: number;
  createdAt: string;
  updatedAt: string;
}

export interface AssessmentQuestion {
  id: string;
  type: "multiple-choice" | "true-false" | "code-completion";
  question: string;
  options?: string[];
  correctAnswer: string | number;
  explanation?: string;
  points: number;
}

export interface AssessmentSubmission {
  id: string;
  assessmentId: string;
  userId: string;
  answers: AssessmentAnswer[];
  score: number;
  maxScore: number;
  percentage: number;
  passed: boolean;
  startedAt: string;
  completedAt: string;
  timeSpent: number; // in seconds
}

export interface AssessmentAnswer {
  questionId: string;
  answer: string | number;
  isCorrect: boolean;
  points: number;
}

// Community Types
export interface ForumPost {
  id: string;
  title: string;
  content: string;
  author: User;
  category: string;
  tags: string[];
  createdAt: string;
  updatedAt: string;
  replies: ForumReply[];
  views: number;
  likes: number;
  dislikes: number;
  solved: boolean;
  pinned: boolean;
  locked: boolean;
}

export interface ForumReply {
  id: string;
  postId: string;
  content: string;
  author: User;
  createdAt: string;
  updatedAt: string;
  likes: number;
  dislikes: number;
  isAccepted: boolean;
  hasCode: boolean;
  codeContent?: string;
  parentReplyId?: string;
  replies?: ForumReply[];
}

export interface CommunityEvent {
  id: string;
  title: string;
  description: string;
  type: "workshop" | "webinar" | "challenge" | "meetup";
  startDate: string;
  endDate: string;
  location?: string;
  isOnline: boolean;
  maxParticipants?: number;
  currentParticipants: number;
  organizer: User;
  tags: string[];
  registrationRequired: boolean;
  registrationDeadline?: string;
  createdAt: string;
}

// Project Types
export interface Project {
  id: string;
  title: string;
  description: string;
  longDescription: string;
  difficulty: "Beginner" | "Intermediate" | "Advanced";
  estimatedTime: string;
  technologies: string[];
  prerequisites: string[];
  learningOutcomes: string[];
  chapters: ProjectChapter[];
  category: string;
  githubRepo?: string;
  liveDemo?: string;
  createdAt: string;
  updatedAt: string;
}

export interface ProjectChapter {
  id: string;
  projectId: string;
  title: string;
  description: string;
  content: string;
  estimatedTime: string;
  order: number;
  codeExample?: string;
  tasks: ProjectTask[];
  prerequisites: string[];
  learningObjectives: string[];
  createdAt: string;
  updatedAt: string;
}

export interface ProjectTask {
  id: string;
  chapterId: string;
  title: string;
  description: string;
  type: "reading" | "coding" | "quiz" | "exercise";
  order: number;
  required: boolean;
  estimatedTime: string;
  resources?: string[];
  createdAt: string;
  updatedAt: string;
}

export interface UserProgress {
  userId: string;
  projectId?: string;
  challengeId?: string;
  assessmentId?: string;
  chapterId?: string;
  taskId?: string;
  status: "not-started" | "in-progress" | "completed" | "skipped";
  progress: number; // 0-100
  score?: number;
  attempts: number;
  timeSpent: number; // in seconds
  startedAt?: string;
  completedAt?: string;
  lastAccessedAt: string;
}

// Leaderboard Types
export interface LeaderboardEntry {
  rank: number;
  user: User;
  points: number;
  change: number; // position change from previous period
  streak: number;
  achievements: string[];
}

export interface Leaderboard {
  period: "daily" | "weekly" | "monthly" | "all-time";
  entries: LeaderboardEntry[];
  userRank?: number;
  totalParticipants: number;
  lastUpdated: string;
}

// Statistics Types
export interface UserStats {
  userId: string;
  totalPoints: number;
  currentStreak: number;
  longestStreak: number;
  challengesCompleted: number;
  assessmentsPassed: number;
  projectsCompleted: number;
  forumPosts: number;
  helpfulAnswers: number;
  reputation: number;
  rank: number;
  achievements: Achievement[];
  weeklyActivity: ActivityData[];
  skillLevels: SkillLevel[];
}

export interface Achievement {
  id: string;
  title: string;
  description: string;
  icon: string;
  category: string;
  rarity: "common" | "rare" | "epic" | "legendary";
  unlockedAt: string;
  progress?: number;
  maxProgress?: number;
}

export interface ActivityData {
  date: string;
  points: number;
  challenges: number;
  projects: number;
  forumActivity: number;
}

export interface SkillLevel {
  skill: string;
  level: number;
  experience: number;
  nextLevelExperience: number;
  category: string;
}

// API Request Types
export interface CreateChallengeSubmissionRequest {
  challengeId: string;
  code: string;
  language: string;
}

export interface CreateAssessmentSubmissionRequest {
  assessmentId: string;
  answers: AssessmentAnswer[];
  timeSpent: number;
}

export interface CreateForumPostRequest {
  title: string;
  content: string;
  category: string;
  tags: string[];
  hasCode?: boolean;
  codeContent?: string;
}

export interface CreateForumReplyRequest {
  postId: string;
  content: string;
  parentReplyId?: string;
  hasCode?: boolean;
  codeContent?: string;
}

export interface UpdateUserProgressRequest {
  projectId?: string;
  challengeId?: string;
  assessmentId?: string;
  chapterId?: string;
  taskId?: string;
  status: "not-started" | "in-progress" | "completed" | "skipped";
  progress: number;
  score?: number;
  timeSpent: number;
}

// Filter and Search Types
export interface ChallengeFilters {
  difficulty?: string;
  category?: string;
  tags?: string[];
  completed?: boolean;
  search?: string;
  sortBy?: "difficulty" | "points" | "completion-rate" | "created-at";
  sortOrder?: "asc" | "desc";
}

export interface ForumPostFilters {
  category?: string;
  tags?: string[];
  solved?: boolean;
  pinned?: boolean;
  author?: string;
  search?: string;
  sortBy?: "created-at" | "replies" | "views" | "likes";
  sortOrder?: "asc" | "desc";
}

export interface ProjectFilters {
  difficulty?: string;
  category?: string;
  technologies?: string[];
  completed?: boolean;
  search?: string;
  sortBy?: "difficulty" | "estimated-time" | "created-at";
  sortOrder?: "asc" | "desc";
}
