export type UserRole = 'admin' | 'worker';

export interface Checklist {
  id: string;
  name: string;
  code: string;             // 'default' | 'sort_priority' ...
  allowed_role_id: string;  // UUID роли, которой разрешён доступ
}

export interface Question {
  id: string;
  text: string;
  order_index: number;
  is_active: boolean;
  checklist_id: string;
  formula?: string;
}

export interface AnswerPayload {
  question_id: string;
  answer_text: string;
  image_url?: string;
}

export interface Report {
  id: string;
  user_id: string;
  checklist_id: string;
  report_date: string;
  responsible_name: string;
  created_at: string;
  place?: string;
  sort?: string;
  priority_sort?: string;
}

export interface ReportDetail extends Report {
  answers: Array<{
    question_id: string;
    question_text: string;
    answer_text: string;
    image_url?: string;
    evaluation?: 'good' | 'bad' | null;
  }>;
}

export interface ReportFilters {
  date_from?: string;
  date_to?: string;
  user_name?: string;
  checklist_id?: string;
  place?: string;
  sort?: string;
  priority_sort?: string;
}