import { useEffect, useState } from 'react';
import { api } from '../api/client';
import type { Question, Checklist } from '../types';

export const Questions = () => {
  const [checklists, setChecklists] = useState<Checklist[]>([]);
  const [activeChecklistId, setActiveChecklistId] = useState<string>('');
  const [questions, setQuestions] = useState<Question[]>([]);

  // Новый вопрос
  const [newText, setNewText] = useState('');
  const [newOrder, setNewOrder] = useState(1);
  const [newFormula, setNewFormula] = useState('');

  // Редактирование
  const [editingId, setEditingId] = useState<string | null>(null);
  const [editText, setEditText] = useState('');
  const [editOrder, setEditOrder] = useState(0);
  const [editChecklistId, setEditChecklistId] = useState('');
  const [editFormula, setEditFormula] = useState('');

  const [loading, setLoading] = useState(false);

  // Загружаем чек-листы и устанавливаем первую вкладку активной
  useEffect(() => {
    api.get('/api/checklists').then(res => {
      setChecklists(res.data);
      if (res.data.length > 0 && !activeChecklistId) {
        setActiveChecklistId(res.data[0].id);
      }
    });
  }, []);

  // При смене активной вкладки загружаем вопросы
  useEffect(() => {
    if (activeChecklistId) {
      loadQuestions(activeChecklistId);
    }
  }, [activeChecklistId]);

  const loadQuestions = async (checklistId: string) => {
    setLoading(true);
    try {
      const res = await api.get(`/api/checklists/${checklistId}/questions`);
      const active = res.data.filter((q: Question) => q.is_active);
      setQuestions(active.sort((a: Question, b: Question) => a.order_index - b.order_index));
    } catch (err) {
      console.error('Ошибка загрузки вопросов', err);
    } finally {
      setLoading(false);
    }
  };

  const addQuestion = async () => {
    if (!newText.trim()) return;
    try {
      await api.post('/api/questions', {
        text: newText,
        order_index: newOrder,
        checklist_id: activeChecklistId,
        formula: newFormula || null,
      });
      setNewText('');
      setNewFormula('');
      setNewOrder(questions.length + 1);
      loadQuestions(activeChecklistId);
    } catch {
      alert('Не удалось добавить вопрос');
    }
  };

  const updateQuestion = async (id: string) => {
    try {
      await api.put(`/api/questions/${id}`, {
        text: editText,
        order_index: editOrder,
        is_active: true,
        checklist_id: editChecklistId,
        formula: editFormula || null,
      });
      setEditingId(null);
      loadQuestions(activeChecklistId);
    } catch {
      alert('Не удалось обновить вопрос');
    }
  };

  const deleteQuestion = async (id: string) => {
    if (!window.confirm('Удалить вопрос?')) return;
    try {
      await api.delete(`/api/questions/${id}`);
      loadQuestions(activeChecklistId);
    } catch {
      alert('Не удалось удалить вопрос');
    }
  };

  const startEdit = (q: Question) => {
    setEditingId(q.id);
    setEditText(q.text);
    setEditOrder(q.order_index);
    setEditChecklistId(q.checklist_id);
    setEditFormula(q.formula || '');
  };

  const activeChecklist = checklists.find(c => c.id === activeChecklistId);

  return (
    <div>
      <h2>Управление вопросами</h2>

      {/* Вкладки чек-листов */}
      <div style={styles.tabs}>
        {checklists.map(cl => (
          <button
            key={cl.id}
            onClick={() => setActiveChecklistId(cl.id)}
            style={{
              ...styles.tabButton,
              ...(activeChecklistId === cl.id ? styles.activeTab : {}),
            }}
          >
            {cl.name}
          </button>
        ))}
      </div>

      {/* Форма создания вопроса (всегда подставляется активный чек-лист) */}
      <div style={styles.newForm}>
        <h3>
          Новый вопрос для «{activeChecklist?.name || '...'}»
        </h3>
        <input
          value={newText}
          onChange={e => setNewText(e.target.value)}
          placeholder="Текст вопроса"
          style={styles.input}
        />
        <input
          type="number"
          value={newOrder}
          onChange={e => setNewOrder(parseInt(e.target.value) || 0)}
          placeholder="Порядок"
          style={styles.inputSmall}
        />
        <input
          value={newFormula}
          onChange={e => setNewFormula(e.target.value)}
          placeholder="Формула (напр. >0.5)"
          style={styles.input}
        />
        <button onClick={addQuestion} style={styles.addBtn}>
          Добавить
        </button>
      </div>

      {loading && <p>Загрузка...</p>}

      {questions.map(q => (
        <div key={q.id} style={styles.item}>
          {editingId === q.id ? (
            <div style={{ flex: 1 }}>
              <input value={editText} onChange={e => setEditText(e.target.value)} style={styles.input} />
              <input
                type="number"
                value={editOrder}
                onChange={e => setEditOrder(parseInt(e.target.value) || 0)}
                style={styles.inputSmall}
              />
              <select
                value={editChecklistId}
                onChange={e => setEditChecklistId(e.target.value)}
                style={styles.select}
              >
                {checklists.map(c => (
                  <option key={c.id} value={c.id}>{c.name}</option>
                ))}
              </select>
              <input
                value={editFormula}
                onChange={e => setEditFormula(e.target.value)}
                placeholder="Формула"
                style={styles.input}
              />
              <button onClick={() => updateQuestion(q.id)} style={styles.saveBtn}>Сохранить</button>
              <button onClick={() => setEditingId(null)} style={styles.cancelBtn}>Отмена</button>
            </div>
          ) : (
            <>
              <div style={{ flex: 1 }}>
                <strong>{q.order_index}.</strong> {q.text}
                {q.formula && (
                  <span style={{ color: '#6b7280', marginLeft: 8 }}>[{q.formula}]</span>
                )}
              </div>
              <button onClick={() => startEdit(q)} style={styles.editBtn}>✏️</button>
              <button onClick={() => deleteQuestion(q.id)} style={styles.deleteBtn}>🗑️</button>
            </>
          )}
        </div>
      ))}

      {!loading && questions.length === 0 && (
        <p>Нет активных вопросов для этого чек-листа.</p>
      )}
    </div>
  );
};

const styles = {
  filterRow: { marginBottom: 20 },
  tabs: {
    display: 'flex',
    gap: 0,
    marginBottom: 24,
    borderBottom: '2px solid #e5e7eb',
  },
  tabButton: {
    padding: '10px 20px',
    background: 'transparent',
    border: 'none',
    borderBottom: '2px solid transparent',
    marginBottom: '-2px',
    cursor: 'pointer',
    fontSize: '15px',
    fontWeight: 500,
    color: '#6b7280',
    transition: 'all 0.2s',
    outline: 'none',
  },
  activeTab: {
    color: '#16a34a',
    borderBottom: '2px solid #16a34a',
  },
  newForm: { background: 'white', padding: 16, borderRadius: 12, marginBottom: 24 },
  input: { width: '100%', padding: 8, marginBottom: 8, borderRadius: 6, border: '1px solid #ccc', boxSizing: 'border-box' as const },
  inputSmall: { width: 100, padding: 8, marginRight: 8, borderRadius: 6, border: '1px solid #ccc' },
  select: { padding: '8px 12px', borderRadius: 6, border: '1px solid #ccc', marginBottom: 8 },
  addBtn: { background: '#16a34a', color: 'white', border: 'none', padding: '8px 16px', borderRadius: 8, cursor: 'pointer' },
  item: { background: 'white', padding: 12, borderRadius: 8, marginBottom: 8, display: 'flex', alignItems: 'center', gap: 12 },
  editBtn: { background: '#eab308', border: 'none', padding: '6px 12px', borderRadius: 6, cursor: 'pointer' },
  deleteBtn: { background: '#ef4444', border: 'none', padding: '6px 12px', borderRadius: 6, cursor: 'pointer', color: 'white' },
  saveBtn: { background: '#16a34a', border: 'none', padding: '6px 12px', borderRadius: 6, cursor: 'pointer', marginRight: 8, color: 'white' },
  cancelBtn: { background: '#6b7280', border: 'none', padding: '6px 12px', borderRadius: 6, cursor: 'pointer', color: 'white' },
};