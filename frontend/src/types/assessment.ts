// Assessment types for the Curriculum Management System

export type AssessmentType = 'quiz' | 'coding_exercise' | 'project';

export type QuestionType = 'multiple_choice' | 'true_false' | 'short_answer' | 'code_completion';

export type SubmissionFormat = 'repo_link' | 'zip_upload' | 'text';

export interface Assessment {
  id: string;
  lessonId: string;
  type: AssessmentType;
  title: string;
  description: string;
  config: QuizConfig | CodingExerciseConfig | ProjectConfig;
  passingScore: number;
  timeLimitMinutes?: number;
  orderIndex: number;
  createdAt: Date;
  updatedAt: Date;
}

export interface QuizConfig {
  questions: Question[];
  settings: QuizSettings;
}

export interface Question {
  id: string;
  type: QuestionType;
  questionText: string;
  options?: string[]; // For multiple choice
  correctAnswer: number | boolean | string; // Index for MC, bool for T/F, text for others
  explanation?: string;
  points: number;
  orderIndex: number;
  hints?: string[];
  tags?: string[];
  readonly?: boolean; // Prevent editing when true
}

export interface QuizSettings {
  shuffleQuestions: boolean;
  shuffleOptions: boolean;
  showExplanations: boolean;
  maxAttempts?: number;
  passingScore: number; // 0-100
}

export interface CodingExerciseConfig {
  description: string; // Rich text
  starterCode: string;
  testCases: TestCase[];
  limits: ExecutionLimits;
  hints: string[];
  solutionTemplate?: string;
}

export interface TestCase {
  id: string;
  input: string;
  expectedOutput: string;
  isVisible: boolean;
  points: number;
}

export interface ExecutionLimits {
  timeLimitSeconds: number;
  memoryLimitMB: number;
}

export interface ProjectConfig {
  description: string; // Rich text
  deliverables: string[];
  submissionFormat: SubmissionFormat;
  starterCodeUrl?: string;
  rubric: RubricCriterion[];
  requirePeerReview: boolean;
  peerReviewCount?: number;
}

export interface RubricCriterion {
  id: string;
  description: string;
  maxPoints: number;
  levels: PerformanceLevel[];
}

export interface PerformanceLevel {
  name: string; // "Excellent", "Good", "Fair", "Poor"
  description: string;
  points: number;
}

export interface Submission {
  id: string;
  assessmentId: string;
  userId: string;
  content: QuizSubmissionContent | CodeSubmissionContent | ProjectSubmissionContent;
  score?: number;
  feedback?: string;
  gradedBy?: string;
  gradedAt?: Date;
  submittedAt: Date;
  status: 'submitted' | 'graded' | 'returned';
}

export interface QuizSubmissionContent {
  answers: Array<{
    questionId: string;
    answer: string | number | boolean;
  }>;
  timeSpent?: number; // seconds
}

export interface CodeSubmissionContent {
  code: string;
  language: string;
  testResults?: TestResult[];
}

export interface TestResult {
  testCaseId: string;
  passed: boolean;
  actualOutput?: string;
  error?: string;
}

export interface ProjectSubmissionContent {
  format: SubmissionFormat;
  repoUrl?: string;
  fileUrl?: string;
  textContent?: string;
}

export interface InlineComment {
  id: string;
  submissionId: string;
  lineNumber: number;
  comment: string;
  authorId: string;
  createdAt: Date;
}

export interface GradebookEntry {
  studentId: string;
  studentName: string;
  assessmentId: string;
  assessmentTitle: string;
  score?: number;
  maxScore: number;
  submittedAt?: Date;
  gradedAt?: Date;
  status: 'submitted' | 'graded' | 'pending';
}

// API request/response types
export interface CreateQuizRequest {
  lessonId: string;
  title: string;
  description: string;
  questions: Omit<Question, 'id'>[];
  settings: QuizSettings;
  passingScore: number;
  timeLimitMinutes?: number;
}

export interface CreateCodingExerciseRequest {
  lessonId: string;
  title: string;
  description: string;
  starterCode: string;
  testCases: Omit<TestCase, 'id'>[];
  limits: ExecutionLimits;
  hints: string[];
  solutionTemplate?: string;
  passingScore: number;
}

export interface CreateProjectAssignmentRequest {
  lessonId: string;
  title: string;
  description: string;
  deliverables: string[];
  submissionFormat: SubmissionFormat;
  starterCodeUrl?: string;
  rubric: Omit<RubricCriterion, 'id'>[];
  requirePeerReview: boolean;
  peerReviewCount?: number;
  passingScore: number;
}

export interface SubmitQuizRequest {
  answers: Array<{
    questionId: string;
    answer: string | number | boolean;
  }>;
}

export interface SubmitCodeExerciseRequest {
  code: string;
}

export interface SubmitProjectRequest {
  format: SubmissionFormat;
  repoUrl?: string;
  fileUrl?: string;
  textContent?: string;
}

export interface GradeSubmissionRequest {
  rubricScores: Record<string, number>; // criterionId -> score
  feedback: string;
  releaseImmediately: boolean;
}
