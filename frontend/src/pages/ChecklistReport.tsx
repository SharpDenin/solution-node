import { useEffect, useState, type CSSProperties } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { api } from '../api/client';
import type { Question, AnswerPayload, Checklist } from '../types';

export const ChecklistReport = () => {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();

  const [checklist, setChecklist] = useState<Checklist | null>(null);
  const [questions, setQuestions] = useState<Question[]>([]);
  const [answers, setAnswers] = useState<Record<string, AnswerPayload>>({});
  const [date, setDate] = useState(new Date().toISOString().split('T')[0]);
  const [responsible, setResponsible] = useState('');
  const [place, setPlace] = useState('');
  const [sort, setSort] = useState('');
  const [prioritySort, setPrioritySort] = useState('low');
  const [loading, setLoading] = useState(false);
  const [checkingAccess, setCheckingAccess] = useState(true);

  useEffect(() => {
    if (!id) return;

    api.get('/api/checklists/available')
      .then(res => {
        const accessibleIds = (res.data as Checklist[]).map(c => c.id);
        if (!accessibleIds.includes(id)) {
          alert('У вас нет доступа к этому чек-листу');
          navigate('/');
          return;
        }
        // Загружаем инфо о чек-листе и вопросы
        return Promise.all([
          api.get(`/api/checklists/${id}`),
          api.get(`/api/checklists/${id}/questions`)
        ]);
      })
      .then((results) => {
        if (!results) return;
        const [clRes, qRes] = results;
        setChecklist(clRes.data);
        const sorted = (qRes.data as Question[])
          .filter(q => q.is_active)
          .sort((a, b) => a.order_index - b.order_index);
        setQuestions(sorted);
      })
      .catch(err => {
        console.error(err);
        navigate('/');
      })
      .finally(() => setCheckingAccess(false));
  }, [id, navigate]);

  const updateAnswer = (qId: string, answer_text: string) => {
    setAnswers(prev => ({
      ...prev,
      [qId]: {
        ...prev[qId],
        question_id: qId,
        answer_text,
      },
    }));
  };

  const uploadImage = async (qId: string, file: File) => {
    const formData = new FormData();
    formData.append('file', file);
    try {
      const res = await api.post('/api/upload', formData, {
        headers: { 'Content-Type': 'multipart/form-data' },
      });
      setAnswers(prev => ({
        ...prev,
        [qId]: {
          ...prev[qId],
          question_id: qId,
          image_url: res.data.url,
          answer_text: prev[qId]?.answer_text || '',
        },
      }));
    } catch {
      alert('Ошибка загрузки изображения');
    }
  };

  const removeImage = (qId: string) => {
    setAnswers(prev => ({
      ...prev,
      [qId]: {
        ...prev[qId],
        image_url: undefined,
      },
    }));
  };

  const handleSubmit = async () => {
    if (!date || !responsible) {
      alert('Заполните обязательные поля (дата, ответственный)');
      return;
    }
    const payload: any = {
      checklist_id: id,
      report_date: date,
      responsible_name: responsible,
      answers: Object.values(answers).map(a => ({
        question_id: a.question_id,
        answer_text: a.answer_text,
        image_url: a.image_url || '',
      })),
    };

    if (checklist?.code === 'sort_priority') {
      payload.sort = sort;
      payload.priority_sort = prioritySort;
    } else {
      payload.place = place || 'Не указано';
    }

    setLoading(true);
    try {
      await api.post('/api/reports', payload);
      navigate('/thank-you');
    } catch {
      alert('Ошибка отправки отчёта');
    } finally {
      setLoading(false);
    }
  };

  if (checkingAccess) return <div>Проверка доступа...</div>;
  if (!checklist) return <div>Чек-лист не найден</div>;

  return (
    <div>
      <h2>{checklist.name}</h2>
      <div style={styles.formGroup}>
        <label>Дата *</label>
        <input
          type="date"
          value={date}
          onChange={e => setDate(e.target.value)}
          style={styles.input}
        />
        <label>Ответственный *</label>
        <input
          placeholder="ФИО ответственного"
          value={responsible}
          onChange={e => setResponsible(e.target.value)}
          style={styles.input}
        />

        {checklist.code === 'sort_priority' ? (
          <>
            <label>Сорт</label>
            <input
              value={sort}
              onChange={e => setSort(e.target.value)}
              style={styles.input}
            />
            <label>Приоритет сорта</label>
            <select
              value={prioritySort}
              onChange={e => setPrioritySort(e.target.value)}
              style={styles.input}
            >
              <option value="low">Низкий</option>
              <option value="high">Высокий</option>
            </select>
          </>
        ) : (
          <>
            <label>Место работ</label>
            <input
              placeholder="Например: Растворный узел №1"
              value={place}
              onChange={e => setPlace(e.target.value)}
              style={styles.input}
            />
          </>
        )}
      </div>

      <hr />

      {questions.map(q => (
        <QuestionCard
          key={q.id}
          question={q}
          answer={answers[q.id]}
          onAnswerChange={(text: string) => updateAnswer(q.id, text)}
          onImageUpload={(file: File) => uploadImage(q.id, file)}
          onImageRemove={() => removeImage(q.id)}
        />
      ))}

      <button onClick={handleSubmit} disabled={loading} style={styles.submitBtn}>
        {loading ? 'Отправка...' : 'Завершить смену'}
      </button>
    </div>
  );
};

// ---- QuestionCard ----
interface QuestionCardProps {
  question: Question;
  answer?: AnswerPayload;
  onAnswerChange: (text: string) => void;
  onImageUpload: (file: File) => void;
  onImageRemove: () => void;
}

const QuestionCard = ({
  question,
  answer,
  onAnswerChange,
  onImageUpload,
  onImageRemove,
}: QuestionCardProps) => {
  const [dragActive, setDragActive] = useState(false);

  const handleDrag = (e: React.DragEvent) => {
    e.preventDefault();
    e.stopPropagation();
    if (e.type === 'dragenter' || e.type === 'dragover') {
      setDragActive(true);
    } else if (e.type === 'dragleave') {
      setDragActive(false);
    }
  };

  const handleDrop = (e: React.DragEvent) => {
    e.preventDefault();
    e.stopPropagation();
    setDragActive(false);
    const files = e.dataTransfer.files;
    if (files && files[0]) {
      onImageUpload(files[0]);
    }
  };

  return (
    <div style={styles.questionCard}>
      <b>{question.text}</b>
      <textarea
        placeholder="Ваш ответ"
        value={answer?.answer_text || ''}
        onChange={e => onAnswerChange(e.target.value)}
        style={styles.textarea}
        rows={3}
      />
      <div
        onDragEnter={handleDrag}
        onDragLeave={handleDrag}
        onDragOver={handleDrag}
        onDrop={handleDrop}
        style={{
          ...styles.dropZone,
          borderColor: dragActive ? '#16a34a' : '#ccc',
          backgroundColor: dragActive ? '#f0fdf4' : '#fafafa',
        }}
      >
        <input
          type="file"
          accept="image/*"
          onChange={(e) => {
            const file = e.target.files?.[0];
            if (file) onImageUpload(file);
          }}
          style={{ display: 'none' }}
          id={`file-${question.id}`}
        />
        <label htmlFor={`file-${question.id}`} style={styles.uploadLabel}>
          📸 Нажмите для выбора или перетащите изображение
        </label>
      </div>
      {answer?.image_url && (
        <div style={styles.imagePreview}>
          <img src={answer.image_url} style={styles.image} alt="фото" />
          <button onClick={onImageRemove} style={styles.removeImage}>
            ✕
          </button>
        </div>
      )}
    </div>
  );
};

// ---- Стили ----
const styles: Record<string, CSSProperties> = {
  formGroup: { marginBottom: 24 },
  input: { width: '100%', padding: '10px', marginBottom: 12, borderRadius: 8, border: '1px solid #ccc', boxSizing: 'border-box' },
  questionCard: { background: 'white', padding: 16, borderRadius: 12, marginBottom: 16 },
  textarea: { width: '100%', marginTop: 8, padding: 8, borderRadius: 8, border: '1px solid #ccc', boxSizing: 'border-box' },
  dropZone: {
    marginTop: 12,
    border: '2px dashed #ccc',
    borderRadius: 8,
    padding: '16px',
    textAlign: 'center' as const,
    cursor: 'pointer',
    transition: 'all 0.2s',
  },
  uploadLabel: {
    cursor: 'pointer',
    color: '#555',
  },
  imagePreview: {
    position: 'relative' as const,
    marginTop: 8,
    display: 'inline-block',
  },
  image: { maxWidth: 200, maxHeight: 150, borderRadius: 8, display: 'block' },
  removeImage: {
    position: 'absolute' as const,
    top: -8,
    right: -8,
    background: '#ef4444',
    color: 'white',
    border: 'none',
    borderRadius: '50%',
    width: 24,
    height: 24,
    cursor: 'pointer',
    fontSize: 12,
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'center',
  },
  submitBtn: { background: '#16a34a', color: 'white', border: 'none', padding: '12px', borderRadius: 8, cursor: 'pointer', width: '100%', marginTop: 16 },
};