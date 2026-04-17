export type UserRole = 'worker' | 'admin';

export interface User {
  id: string;
  fullName: string;
  login: string;
  role: UserRole;
}

export interface Question {
  id: string;
  text: string;
  order_index: number;
  is_active: boolean;
}

export interface AnswerPayload {
  question_id: string;
  answer_text: string;
  image_url?: string;
}

export interface Report {
  id: string;
  user_id: string;
  place: string;
  report_date: string;
  responsible_name: string;
  created_at: string;
}

export interface ReportDetail extends Report {
  answers: Array<{
    question_id: string;
    question_text: string;
    answer_text: string;
    image_url?: string;
  }>;
}

export interface ReportFilters {
  date_from?: string;
  date_to?: string;
  place?: string;
  user_name?: string;
}