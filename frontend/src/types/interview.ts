export type InterviewType = 'coding' | 'behavioral' | 'system_design';
export type Difficulty = 'beginner' | 'intermediate' | 'advanced';

export interface Question {
  id: string;
  content: string;
  type: InterviewType;
  difficulty: Difficulty;
  expected_points?: string[];
  time_limit: number;
}

export interface Answer {
  question_id: string;
  content: string;
  score?: number;
  feedback?: string;
  created_at: string;
}

export interface InterviewSession {
  id: string;
  user_id?: string;
  type: InterviewType;
  difficulty: Difficulty;
  questions: Question[];
  current_index?: number;
  answers?: Answer[];
  status: 'in_progress' | 'completed' | 'abandoned';
  score?: number;
  created_at: string;
  completed_at?: string;
}

export interface InterviewFeedback {
  session_id: string;
  overall_score: number;
  strengths: string[];
  improvements: string[];
  detailed_feedback: QuestionFeedback[];
}

export interface QuestionFeedback {
  question_id: string;
  score: number;
  feedback: string;
  strengths: string[];
  missed: string[];
}
