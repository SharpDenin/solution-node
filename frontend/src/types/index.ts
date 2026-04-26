export type UserRole = 'admin' | 'worker';

export interface Checklist {
  id: string;
  name: string;
  code: string;             // 'default' | 'sort_control' ...
  allowed_role_id?: string;
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
  result?: 'good' | 'neutral' | 'bad';
  image_url?: string;
}

export interface Report {
  id: string;
  user_id: string;
  checklist_id: string;
  report_date: string;
  responsible_name: string;
  metadata?: string;        // JSON-строка с дополнительными данными
  created_at: string;
  // поля, извлечённые из metadata (заполняются после парсинга)
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
    result?: 'good' | 'neutral' | 'bad';
  }>;
  metadata?: string;
  place?: string;
  sort?: string;
  priority_sort?: string;
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