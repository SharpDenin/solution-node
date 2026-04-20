import { useEffect, useState } from 'react';
import { api } from '../api/client';
import type { Question } from '../types';

export const Questions = () => {
  const [questions, setQuestions] = useState<Question[]>([]);
  const [editingId, setEditingId] = useState<string | null>(null);
  const [editText, setEditText] = useState('');
  const [editOrder, setEditOrder] = useState(0);
  const [newText, setNewText] = useState('');
  const [newOrder, setNewOrder] = useState(1);
  const [loading, setLoading] = useState(false);

  const loadQuestions = async () => {
    setLoading(true);
    try {
      const res = await api.get('/api/questions');
      // Показываем только активные вопросы (is_active === true)
      const activeQuestions = res.data.filter((q: Question) => q.is_active === true);
      const sorted = activeQuestions.sort((a: Question, b: Question) => a.order_index - b.order_index);
      setQuestions(sorted);
    } catch (err) {
      console.error('Ошибка загрузки вопросов', err);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    loadQuestions();
  }, []);

  const addQuestion = async () => {
    if (!newText.trim()) return;
    try {
      await api.post('/api/questions', { text: newText, order_index: newOrder });
      setNewText('');
      setNewOrder(questions.length + 1);
      loadQuestions();
    } catch (err) {
      console.error('Ошибка добавления вопроса', err);
      alert('Не удалось добавить вопрос');
    }
  };

  const updateQuestion = async (id: string) => {
    try {
      await api.put(`/api/questions/${id}`, { text: editText, order_index: editOrder, is_active: true });
      setEditingId(null);
      loadQuestions();
    } catch (err) {
      console.error('Ошибка обновления вопроса', err);
      alert('Не удалось обновить вопрос');
    }
  };

  const deleteQuestion = async (id: string) => {
    if (!window.confirm('Удалить вопрос?')) return;
    try {
      await api.delete(`/api/questions/${id}`);
      loadQuestions();
    } catch (err) {
      console.error('Ошибка удаления вопроса', err);
      alert('Не удалось удалить вопрос');
    }
  };

  const startEdit = (q: Question) => {
    setEditingId(q.id);
    setEditText(q.text);
    setEditOrder(q.order_index);
  };

  return (
    <div>
      <h2>Управление вопросами</h2>
      <div style={{ marginBottom: 20, padding: 16, background: 'white', borderRadius: 12 }}>
        <h3>Новый вопрос</h3>
        <input
          value={newText}
          onChange={e => setNewText(e.target.value)}
          placeholder="Текст вопроса"
          style={{ width: '100%', padding: 8, marginBottom: 8, borderRadius: 6, border: '1px solid #ccc' }}
        />
        <input
          type="number"
          value={newOrder}
          onChange={e => setNewOrder(parseInt(e.target.value) || 0)}
          placeholder="Порядковый номер"
          style={{ width: 100, padding: 8, marginRight: 8, borderRadius: 6, border: '1px solid #ccc' }}
        />
        <button onClick={addQuestion} style={styles.addBtn}>Добавить</button>
      </div>

      {loading && <p>Загрузка...</p>}

      {questions.map(q => (
        <div key={q.id} style={styles.item}>
          {editingId === q.id ? (
            <div style={{ flex: 1 }}>
              <input
                value={editText}
                onChange={e => setEditText(e.target.value)}
                style={{ width: '100%', marginBottom: 8, padding: 6, borderRadius: 6, border: '1px solid #ccc' }}
              />
              <input
                type="number"
                value={editOrder}
                onChange={e => setEditOrder(parseInt(e.target.value) || 0)}
                style={{ width: 80, padding: 6, borderRadius: 6, border: '1px solid #ccc', marginRight: 8 }}
              />
              <button onClick={() => updateQuestion(q.id)} style={styles.saveBtn}>Сохранить</button>
              <button onClick={() => setEditingId(null)} style={styles.cancelBtn}>Отмена</button>
            </div>
          ) : (
            <>
              <div style={{ flex: 1 }}>
                <strong>{q.order_index}.</strong> {q.text}
              </div>
              <button onClick={() => startEdit(q)} style={styles.editBtn}>✏️</button>
              <button onClick={() => deleteQuestion(q.id)} style={styles.deleteBtn}>🗑️</button>
            </>
          )}
        </div>
      ))}
      {!loading && questions.length === 0 && <p>Нет активных вопросов. Создайте первый вопрос.</p>}
    </div>
  );
};

const styles = {
  addBtn: {
    background: '#16a34a',
    color: 'white',
    border: 'none',
    padding: '8px 16px',
    borderRadius: 8,
    cursor: 'pointer',
  },
  item: {
    background: 'white',
    padding: 12,
    borderRadius: 8,
    marginBottom: 8,
    display: 'flex',
    alignItems: 'center',
    gap: 12,
  },
  editBtn: {
    background: '#eab308',
    border: 'none',
    padding: '6px 12px',
    borderRadius: 6,
    cursor: 'pointer',
  },
  deleteBtn: {
    background: '#ef4444',
    border: 'none',
    padding: '6px 12px',
    borderRadius: 6,
    cursor: 'pointer',
    color: 'white',
  },
  saveBtn: {
    background: '#16a34a',
    border: 'none',
    padding: '6px 12px',
    borderRadius: 6,
    cursor: 'pointer',
    marginRight: 8,
    color: 'white',
  },
  cancelBtn: {
    background: '#6b7280',
    border: 'none',
    padding: '6px 12px',
    borderRadius: 6,
    cursor: 'pointer',
    color: 'white',
  },
};