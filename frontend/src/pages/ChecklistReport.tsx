import { useEffect, useState, type CSSProperties } from 'react';
import { useParams, useNavigate, useSearchParams } from 'react-router-dom';
import { api } from '../api/client';
import { useAuth } from '../context/AuthContext';
import type { Question, AnswerPayload, Checklist, Variety, Phenophase } from '../types';

interface QuestionCardProps {
  question: Question;
  answer?: AnswerPayload;
  onAnswerChange: (text: string) => void;
  onImageUpload: (file: File) => void;
  onImageRemove: () => void;
}

export const ChecklistReport = () => {
  const { id: checklistId } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const [searchParams] = useSearchParams();
  const { fullName, role } = useAuth();

  const [checklist, setChecklist] = useState<Checklist | null>(null);
  const [questions, setQuestions] = useState<Question[]>([]);
  const [answers, setAnswers] = useState<Record<string, AnswerPayload>>({});
  const [date] = useState(new Date().toISOString().split('T')[0]);

  const [place, setPlace] = useState('');
  const [variety, setVariety] = useState<Variety | null>(null);
  const [phenophase, setPhenophase] = useState<Phenophase | null>(null);

  const [loading, setLoading] = useState(false);
  const [checkingAccess, setCheckingAccess] = useState(true);

  useEffect(() => {
    if (!checklistId) return;

    const loadChecklistById = () => {
      return api.get('/api/checklists').then(res => {
        const all = res.data as Checklist[];
        const found = all.find(c => c.id === checklistId);
        if (!found) throw new Error('Чек-лист не найден');
        setChecklist(found);
        return found;
      });
    };

    const loadQuestions = async (cl: Checklist) => {
      if (cl.code === 'sort_control') {
        const varietyId = searchParams.get('varietyId');
        const phenophaseId = searchParams.get('phenophaseId');
        if (!varietyId || !phenophaseId) {
          navigate(`/checklist/${checklistId}/variety`);
          return Promise.reject('need selection');
        }
        const [vRes, pRes] = await Promise.all([
          api.get(`/api/varieties/${varietyId}`),
          api.get(`/api/phenophases/${phenophaseId}`)
        ]);
        setVariety(vRes.data);
        setPhenophase(pRes.data);

        // Новый эндпоинт с автозаполнением
        const res = await api.get(
          `/api/checklists/${checklistId}/questions/defaults?phenophase_id=${phenophaseId}`
        );
        return res.data;
      } else {
        const res = await api.get(`/api/checklists/${checklistId}/questions`);
        return res.data;
      }
    };

    const processData = (questionsData: Question[]) => {
      const sorted = questionsData
        .filter(q => q.is_active)
        .sort((a, b) => a.order_index - b.order_index);
      setQuestions(sorted);

      // Автозаполнение ответов из default_answer
      const initialAnswers: Record<string, AnswerPayload> = {};
      sorted.forEach(q => {
        initialAnswers[q.id] = {
          question_id: q.id,
          answer_text: q.default_answer ?? '',
          image_url: undefined,
        };
      });
      setAnswers(initialAnswers);
    };

    if (role === 'admin') {
      loadChecklistById()
        .then(cl => loadQuestions(cl))
        .then(data => {
          if (Array.isArray(data)) {
            processData(data);
          } else if (data) {
            processData(data);
          }
        })
        .catch(err => {
          if (err === 'need selection') return;
          console.error(err);
          alert('Ошибка загрузки данных чек-листа');
          navigate('/');
        })
        .finally(() => setCheckingAccess(false));
    } else {
      api.get('/api/checklists/available')
        .then(res => {
          const accessible = res.data as Checklist[];
          const found = accessible.find(c => c.id === checklistId);
          if (!found) {
            alert('У вас нет доступа к этому чек-листу');
            navigate('/');
            return Promise.reject('no access');
          }
          setChecklist(found);
          return loadQuestions(found);
        })
        .then(data => {
          if (Array.isArray(data)) {
            processData(data);
          } else if (data) {
            processData(data);
          }
        })
        .catch(err => {
          if (err === 'no access' || err === 'need selection') return;
          console.error(err);
          navigate('/login?returnUrl=/checklist/' + checklistId);
        })
        .finally(() => setCheckingAccess(false));
    }
  }, [checklistId, navigate, searchParams, role]);

  const updateAnswer = (qId: string, answer_text: string) => {
    setAnswers(prev => ({
      ...prev,
      [qId]: { ...prev[qId], question_id: qId, answer_text },
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
      [qId]: { ...prev[qId], image_url: undefined },
    }));
  };

  const handleSubmit = async () => {
    if (!date || !fullName) {
      alert('Заполните обязательные поля');
      return;
    }

    const payload: any = {
      checklist_id: checklistId,
      report_date: date,
      responsible_name: fullName,
      metadata: {},
      answers: Object.values(answers).map(a => ({
        question_id: a.question_id,
        answer_text: a.answer_text,
        image_url: a.image_url || '',
      })),
    };

    if (checklist?.code === 'sort_control' && variety && phenophase) {
      payload.variety_id = variety.id;
      payload.phenophase_id = phenophase.id;
      payload.metadata = {
        sort: variety.name,
        priority_sort: variety.priority,
        phenophase: phenophase.name,
        variety_id: variety.id,
        phenophase_id: phenophase.id,
      };
    } else if (checklist?.code === 'default') {
      payload.metadata = {
        place: place || 'Не указано',
      };
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
        <label>Дата</label>
        <input type="date" value={date} disabled style={{ ...styles.input, backgroundColor: '#f3f4f6' }} />
        <label>Ответственный</label>
        <input value={fullName || ''} disabled style={{ ...styles.input, backgroundColor: '#f3f4f6' }} />

        {checklist.code === 'default' && (
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

        {checklist.code === 'sort_control' && variety && phenophase && (
          <>
            <label>Сорт</label>
            <input
              value={variety.name}
              disabled
              style={{
                ...styles.input,
                backgroundColor: '#f3f4f6',
                color: variety.priority === 'high' ? '#ef4444' : '#16a34a',
                fontWeight: 600,
              }}
            />
            <label>Приоритет сорта</label>
            <input
              value={variety.priority === 'high' ? 'Высокий' : 'Низкий'}
              disabled
              style={{
                ...styles.input,
                backgroundColor: '#f3f4f6',
                color: variety.priority === 'high' ? '#ef4444' : '#16a34a',
                fontWeight: 600,
              }}
            />
            <label>Фенофаза</label>
            <input
              value={phenophase.name}
              disabled
              style={{ ...styles.input, backgroundColor: '#f3f4f6' }}
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

const QuestionCard = ({
  question,
  answer,
  onAnswerChange,
  onImageUpload,
  onImageRemove,
}: QuestionCardProps) => {
  const [dragActive, setDragActive] = useState(false);

  // Если есть дефолтный ответ – блокируем ввод и фото
  const hasDefault = !!question.default_answer;

  const handleDrag = (e: React.DragEvent) => {
    if (hasDefault) return;
    e.preventDefault();
    e.stopPropagation();
    if (e.type === 'dragenter' || e.type === 'dragover') {
      setDragActive(true);
    } else if (e.type === 'dragleave') {
      setDragActive(false);
    }
  };

  const handleDrop = (e: React.DragEvent) => {
    if (hasDefault) return;
    e.preventDefault();
    e.stopPropagation();
    setDragActive(false);
    const files = e.dataTransfer.files;
    if (files && files[0]) {
      onImageUpload(files[0]);
    }
  };

  const handleFileSelect = (e: React.ChangeEvent<HTMLInputElement>) => {
    if (hasDefault) return;
    const file = e.target.files?.[0];
    if (file) onImageUpload(file);
  };

  return (
    <div style={styles.questionCard}>
      <div style={{ display: 'flex', alignItems: 'center', gap: 8 }}>
        <b>{question.text}</b>
        {question.image_url && (
          <img src={question.image_url} style={{ width: 32, height: 32, borderRadius: 4 }} alt="q" />
        )}
        {hasDefault && (
          <span style={{ fontSize: 12, color: '#6b7280', marginLeft: 8 }}>
            (автозаполнено)
          </span>
        )}
      </div>
      <textarea
        placeholder="Ваш ответ"
        value={answer?.answer_text || ''}
        onChange={e => onAnswerChange(e.target.value)}
        style={{
          ...styles.textarea,
          backgroundColor: hasDefault ? '#f3f4f6' : 'white',
          color: hasDefault ? '#6b7280' : 'inherit',
        }}
        rows={3}
        disabled={hasDefault}
      />
      {!hasDefault && (
        <>
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
              onChange={handleFileSelect}
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
              <button onClick={onImageRemove} style={styles.removeImage}>✕</button>
            </div>
          )}
        </>
      )}
    </div>
  );
};

const styles: Record<string, CSSProperties> = {
  formGroup: { marginBottom: 24 },
  input: {
    width: '100%',
    padding: '10px',
    marginBottom: 12,
    borderRadius: 8,
    border: '1px solid #ccc',
    boxSizing: 'border-box',
  },
  questionCard: {
    background: 'white',
    padding: 16,
    borderRadius: 12,
    marginBottom: 16,
  },
  textarea: {
    width: '100%',
    marginTop: 8,
    padding: 8,
    borderRadius: 8,
    border: '1px solid #ccc',
    boxSizing: 'border-box',
  },
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
  image: {
    maxWidth: 200,
    maxHeight: 150,
    borderRadius: 8,
    display: 'block',
  },
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
  submitBtn: {
    background: '#16a34a',
    color: 'white',
    border: 'none',
    padding: '12px',
    borderRadius: 8,
    cursor: 'pointer',
    width: '100%',
    marginTop: 16,
  },
};