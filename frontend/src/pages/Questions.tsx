import { useEffect, useState } from 'react';
import { api } from '../api/client';
import { ImageUploader } from '../components/ImageUploader';
import type { Question, Checklist, Phenophase, QuestionPhenophaseFormula } from '../types';

export const Questions = () => {
  const [checklists, setChecklists] = useState<Checklist[]>([]);
  const [phenophases, setPhenophases] = useState<Phenophase[]>([]);
  const [activeChecklistId, setActiveChecklistId] = useState<string>('');
  const [questions, setQuestions] = useState<Question[]>([]);

  // поля для нового вопроса
  const [newText, setNewText] = useState('');
  const [newImageUrl, setNewImageUrl] = useState<string | undefined>(undefined);
  const [newFormulaDefault, setNewFormulaDefault] = useState('');
  const [newFormulas, setNewFormulas] = useState<QuestionPhenophaseFormula[]>([]);

  // редактирование
  const [editingId, setEditingId] = useState<string | null>(null);
  const [editText, setEditText] = useState('');
  const [editChecklistId, setEditChecklistId] = useState('');
  const [editImageUrl, setEditImageUrl] = useState<string | undefined>(undefined);
  const [editFormulaDefault, setEditFormulaDefault] = useState('');
  const [editFormulas, setEditFormulas] = useState<QuestionPhenophaseFormula[]>([]);
  const [editOrder, setEditOrder] = useState(0); // не показываем, но храним для отправки

  const [loading, setLoading] = useState(false);

  useEffect(() => {
    api.get('/api/checklists').then(res => {
      setChecklists(res.data);
      if (res.data.length > 0 && !activeChecklistId) {
        setActiveChecklistId(res.data[0].id);
      }
    });
    api.get('/api/phenophases').then(res => setPhenophases(res.data));
  }, []);

  useEffect(() => {
    if (activeChecklistId) {
      loadQuestions(activeChecklistId);
      resetNewForm();
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

  const resetNewForm = () => {
    setNewText('');
    setNewImageUrl(undefined);
    setNewFormulaDefault('');
    if (activeChecklist?.code === 'sort_control') {
      setNewFormulas(phenophases.map(p => ({ phenophase_id: p.id, formula: '' })));
    } else {
      setNewFormulas([]);
    }
  };

  const activeChecklist = checklists.find(c => c.id === activeChecklistId);
  const isSortControl = activeChecklist?.code === 'sort_control';

  const uploadImage = async (file: File): Promise<string | undefined> => {
    const formData = new FormData();
    formData.append('file', file);
    try {
      const res = await api.post('/api/upload', formData, {
        headers: { 'Content-Type': 'multipart/form-data' },
      });
      return res.data.url;
    } catch {
      alert('Ошибка загрузки изображения');
      return undefined;
    }
  };

  const getMaxOrderIndex = (): number => {
    if (questions.length === 0) return 1;
    return Math.max(...questions.map(q => q.order_index)) + 1;
  };

  const addQuestion = async () => {
    if (!newText.trim()) return;
    const orderIndex = getMaxOrderIndex();
    const payload: any = {
      text: newText,
      order_index: orderIndex,
      is_active: true,
      checklist_id: activeChecklistId,
      image_url: newImageUrl || undefined,
    };
    if (isSortControl) {
      payload.formulas = newFormulas.filter(f => f.formula.trim() !== '');
      payload.formula = undefined;
    } else {
      payload.formula = newFormulaDefault || undefined;
      payload.formulas = [];
    }
    try {
      await api.post('/api/questions', payload);
      loadQuestions(activeChecklistId);
    } catch {
      alert('Не удалось добавить вопрос');
    }
  };

  const startEdit = (q: Question) => {
    setEditingId(q.id);
    setEditText(q.text);
    setEditChecklistId(q.checklist_id);
    setEditImageUrl(q.image_url);
    setEditFormulaDefault(q.formula || '');
    setEditOrder(q.order_index); // сохраняем текущий порядок
    if (isSortControl) {
      const existingFormulas = q.formulas || [];
      const merged = phenophases.map(p => {
        const existing = existingFormulas.find(f => f.phenophase_id === p.id);
        return {
          phenophase_id: p.id,
          formula: existing ? existing.formula : '',
        };
      });
      setEditFormulas(merged);
    } else {
      setEditFormulas([]);
    }
  };

  const updateQuestion = async (id: string) => {
    const payload: any = {
      text: editText,
      order_index: editOrder, // сохраняем текущий порядок без изменений
      is_active: true,
      checklist_id: editChecklistId,
      image_url: editImageUrl || undefined,
    };
    if (isSortControl) {
      payload.formulas = editFormulas.filter(f => f.formula.trim() !== '');
      payload.formula = undefined;
    } else {
      payload.formula = editFormulaDefault || undefined;
      payload.formulas = [];
    }
    try {
      await api.put(`/api/questions/${id}`, payload);
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

  return (
    <div>
      <h2>Управление вопросами</h2>

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

      <div style={styles.newForm}>
        <h3>Новый вопрос для «{activeChecklist?.name || '...'}»</h3>
        <input
          value={newText}
          onChange={e => setNewText(e.target.value)}
          placeholder="Текст вопроса"
          style={styles.input}
        />
        <ImageUploader
          imageUrl={newImageUrl}
          onUpload={async (file) => {
            const url = await uploadImage(file);
            if (url) setNewImageUrl(url);
          }}
          onRemove={() => setNewImageUrl(undefined)}
        />

        {!isSortControl ? (
          <input
            value={newFormulaDefault}
            onChange={e => setNewFormulaDefault(e.target.value)}
            placeholder="Формула (например >0.5)"
            style={styles.input}
          />
        ) : (
          <div style={styles.formulaList}>
            <label>Формулы по фенофазам:</label>
            {phenophases.map(p => {
              const idx = newFormulas.findIndex(f => f.phenophase_id === p.id);
              const value = idx >= 0 ? newFormulas[idx].formula : '';
              const handleChange = (val: string) => {
                setNewFormulas(prev => {
                  const copy = [...prev];
                  if (idx >= 0) {
                    copy[idx] = { ...copy[idx], formula: val };
                  } else {
                    copy.push({ phenophase_id: p.id, formula: val });
                  }
                  return copy;
                });
              };
              return (
                <div key={p.id} style={styles.formulaRow}>
                  <span style={styles.phenoName}>{p.name}:</span>
                  <input
                    value={value}
                    onChange={e => handleChange(e.target.value)}
                    placeholder="формула"
                    style={{ ...styles.input, flex: 1 }}
                  />
                </div>
              );
            })}
          </div>
        )}
        <button onClick={addQuestion} style={styles.addBtn}>Добавить</button>
      </div>

      {loading && <p>Загрузка...</p>}

      {questions.map(q => (
        <div key={q.id} style={styles.item}>
          {editingId === q.id ? (
            <div style={{ flex: 1 }}>
              <input value={editText} onChange={e => setEditText(e.target.value)} style={styles.input} />
              <select
                value={editChecklistId}
                onChange={e => setEditChecklistId(e.target.value)}
                style={styles.select}
              >
                {checklists.map(c => (
                  <option key={c.id} value={c.id}>{c.name}</option>
                ))}
              </select>
              <ImageUploader
                imageUrl={editImageUrl}
                onUpload={async (file) => {
                  const url = await uploadImage(file);
                  if (url) setEditImageUrl(url);
                }}
                onRemove={() => setEditImageUrl(undefined)}
              />

              {!isSortControl ? (
                <input
                  value={editFormulaDefault}
                  onChange={e => setEditFormulaDefault(e.target.value)}
                  placeholder="Формула"
                  style={styles.input}
                />
              ) : (
                <div style={styles.formulaList}>
                  {phenophases.map(p => {
                    const idx = editFormulas.findIndex(f => f.phenophase_id === p.id);
                    const value = idx >= 0 ? editFormulas[idx].formula : '';
                    const handleChange = (val: string) => {
                      setEditFormulas(prev => {
                        const copy = [...prev];
                        if (idx >= 0) {
                          copy[idx] = { ...copy[idx], formula: val };
                        } else {
                          copy.push({ phenophase_id: p.id, formula: val });
                        }
                        return copy;
                      });
                    };
                    return (
                      <div key={p.id} style={styles.formulaRow}>
                        <span style={styles.phenoName}>{p.name}:</span>
                        <input
                          value={value}
                          onChange={e => handleChange(e.target.value)}
                          style={{ ...styles.input, flex: 1 }}
                        />
                      </div>
                    );
                  })}
                </div>
              )}
              <button onClick={() => updateQuestion(q.id)} style={styles.saveBtn}>Сохранить</button>
              <button onClick={() => setEditingId(null)} style={styles.cancelBtn}>Отмена</button>
            </div>
          ) : (
            <>
              <div style={{ flex: 1 }}>
                <strong>{q.order_index}.</strong> {q.text}
                {q.image_url && <img src={q.image_url} style={styles.thumb} alt="preview" />}
                {q.formula && <span style={{ color: '#6b7280', marginLeft: 8 }}>[{q.formula}]</span>}
                {q.formulas && q.formulas.length > 0 && (
                  <span style={{ color: '#6b7280', marginLeft: 8 }}>
                    (формул: {q.formulas.length})
                  </span>
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
  select: { padding: '8px 12px', borderRadius: 6, border: '1px solid #ccc', marginBottom: 8 },
  addBtn: { background: '#16a34a', color: 'white', border: 'none', padding: '8px 16px', borderRadius: 8, cursor: 'pointer' },
  item: { background: 'white', padding: 12, borderRadius: 8, marginBottom: 8, display: 'flex', alignItems: 'center', gap: 12 },
  editBtn: { background: '#eab308', border: 'none', padding: '6px 12px', borderRadius: 6, cursor: 'pointer' },
  deleteBtn: { background: '#ef4444', border: 'none', padding: '6px 12px', borderRadius: 6, cursor: 'pointer', color: 'white' },
  saveBtn: { background: '#16a34a', border: 'none', padding: '6px 12px', borderRadius: 6, cursor: 'pointer', marginRight: 8, color: 'white' },
  cancelBtn: { background: '#6b7280', border: 'none', padding: '6px 12px', borderRadius: 6, cursor: 'pointer', color: 'white' },
  formulaList: { marginBottom: 12 },
  formulaRow: { display: 'flex', alignItems: 'center', gap: 8, marginBottom: 4 },
  phenoName: { minWidth: 120, fontWeight: 500 },
  thumb: { width: 30, height: 30, borderRadius: 4, marginLeft: 8, verticalAlign: 'middle' },
};