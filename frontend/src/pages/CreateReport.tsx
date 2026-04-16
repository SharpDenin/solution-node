import { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { api } from '../api/client';
import type { Question, AnswerPayload } from '../types';

export const CreateReport = () => {
  const [questions, setQuestions] = useState<Question[]>([]);
  const [answers, setAnswers] = useState<Record<string, AnswerPayload>>({});
  const [place, setPlace] = useState('');
  const [date, setDate] = useState(new Date().toISOString().split('T')[0]);
  const [responsible, setResponsible] = useState('');
  const [loading, setLoading] = useState(false);
  const navigate = useNavigate();

  useEffect(() => {
    api.get('/questions')
      .then(res => {
        // Фильтруем только активные вопросы
        const active = res.data.filter((q: Question) => q.is_active === true);
        const sorted = active.sort((a: Question, b: Question) => a.order_index - b.order_index);
        setQuestions(sorted);
      })
      .catch(err => console.error('Ошибка загрузки вопросов', err));
  }, []);

  const updateAnswer = (qId: string, answer_text: string) => {
    setAnswers(prev => ({
      ...prev,
      [qId]: { question_id: qId, answer_text },
    }));
  };

  const uploadImage = async (qId: string, file: File) => {
    const formData = new FormData();
    formData.append('file', file);
    try {
      const res = await api.post('/upload', formData, {
        headers: { 'Content-Type': 'multipart/form-data' },
      });
      setAnswers(prev => ({
        ...prev,
        [qId]: {
          ...prev[qId],
          question_id: qId,
          answer_text: prev[qId]?.answer_text || '',
          image_url: res.data.url,
        },
      }));
    } catch (err) {
      alert('Ошибка загрузки изображения');
    }
  };

  const handleSubmit = async () => {
    if (!place || !date || Object.keys(answers).length === 0) {
      alert('Заполните место, дату и ответьте хотя бы на один вопрос');
      return;
    }
    setLoading(true);
    try {
      const payload = {
        place,
        report_date: date,
        responsible_name: responsible || 'Не указан',
        answers: Object.values(answers).map(a => ({
          question_id: a.question_id,
          answer_text: a.answer_text,
          image_url: a.image_url || '',
        })),
      };
      console.log('Отправляемые данные:', payload); // <- посмотрите в консоли браузера
      await api.post('/reports', payload);
      navigate('/thank-you');
    } catch (err: any) {
      console.error('Ошибка:', err);
      const message = err.response?.data || 'Ошибка при отправке отчёта';
      alert(message);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div>
      <h2>Новый отчёт</h2>
      <div style={styles.formGroup}>
        <label>Место работ *</label>
        <input
          placeholder="Например: Растворный узел №1"
          value={place}
          onChange={e => setPlace(e.target.value)}
          style={styles.input}
        />
        <label>Дата *</label>
        <input
          type="date"
          value={date}
          onChange={e => setDate(e.target.value)}
          style={styles.input}
        />
        <label>Ответственный</label>
        <input
          placeholder="ФИО ответственного"
          value={responsible}
          onChange={e => setResponsible(e.target.value)}
          style={styles.input}
        />
      </div>

      <hr />

      {questions.map(q => (
        <div key={q.id} style={styles.questionCard}>
          <b>{q.text}</b>
          <textarea
            placeholder="Ваш ответ"
            onChange={e => updateAnswer(q.id, e.target.value)}
            style={styles.textarea}
            rows={3}
          />
          <input type="file" accept="image/*" onChange={e => {
            const file = e.target.files?.[0];
            if (file) uploadImage(q.id, file);
          }} />
          {answers[q.id]?.image_url && (
            <img src={answers[q.id].image_url} style={styles.image} alt="фото" />
          )}
        </div>
      ))}

      <button onClick={handleSubmit} disabled={loading} style={styles.submitBtn}>
        {loading ? 'Отправка...' : 'Завершить смену'}
      </button>
    </div>
  );
};

const styles = {
  formGroup: { marginBottom: 24 },
  input: { width: '100%', padding: '10px', marginBottom: 12, borderRadius: 8, border: '1px solid #ccc' },
  questionCard: { background: 'white', padding: 16, borderRadius: 12, marginBottom: 16 },
  textarea: { width: '100%', marginTop: 8, padding: 8, borderRadius: 8, border: '1px solid #ccc' },
  image: { maxWidth: 200, marginTop: 8, display: 'block' },
  submitBtn: { background: '#16a34a', color: 'white', border: 'none', padding: '12px', borderRadius: 8, cursor: 'pointer', width: '100%' },
};