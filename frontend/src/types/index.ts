export type UserRole = 'admin' | 'node' | 'phenophase';

export interface Checklist {
  id: string;
  name: string;
  code: string; // 'default' | 'sort_control'
}

export interface Question {
  id: string;
  text: string;
  order_index: number;
  is_active: boolean;
  checklist_id: string;
  formula?: string;
  image_url?: string;
  formulas?: QuestionPhenophaseFormula[];
}

export interface QuestionPhenophaseFormula {
  id?: string;
  question_id?: string;
  phenophase_id: string;
  formula: string;
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
  metadata?: any;
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
    result?: 'good' | 'neutral' | 'bad';
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

export interface CurrentUser {
  id: string;
  full_name: string;
  login: string;
  role: UserRole;
  position?: string;
}

export interface Variety {
  id: string;
  name: string;
  description?: string;
  priority: 'high' | 'low';
  image_url?: string;
  created_at?: string;
}

export interface Phenophase {
  id: string;
  name: string;
  description?: string;
  image_url?: string;
  order_index: number;
  created_at?: string;
}

export interface PhenophaseMatrixReportResponse {
  variety_id: string;
  columns: PhenophaseMatrixColumn[];
  rows: PhenophaseMatrixRow[];
}

export interface PhenophaseMatrixColumn {
  phenophase_id: string;
  name: string;
  order_index: number;
}

export interface PhenophaseMatrixRow {
  question_id: string;
  text: string;
  order_index: number;
  cells: PhenophaseMatrixCell[];
}

export interface PhenophaseMatrixCell {
  phenophase_id: string;
  answer_text: string | null;
  result: string | null;
  image_url: string | null;
  report_id: string | null;
}

export interface UserAdmin {
  id: string;
  full_name: string;
  login: string;
  role: string;
  position?: string;
  is_blocked: boolean;
  is_deleted: boolean;
  created_at: string;
  updated_at: string;
}