import { useEffect, useState } from 'react';
import { api } from '../api/client';
import { ImageUploader } from '../components/ImageUploader';
import type { Variety } from '../types';

export const VarietyManagePage = () => {
  const [varieties, setVarieties] = useState<Variety[]>([]);
  const [loading, setLoading] = useState(false);

  // создание / редактирование
  const [showForm, setShowForm] = useState(false);
  const [editingId, setEditingId] = useState<string | null>(null);
  const [name, setName] = useState('');
  const [description, setDescription] = useState('');
  const [priority, setPriority] = useState<'high' | 'low'>('high');
  const [imageUrl, setImageUrl] = useState<string | undefined>(undefined);

  const fetchVarieties = async () => {
    setLoading(true);
    try {
      const res = await api.get('/api/varieties');
      setVarieties(res.data);
    } catch {
      alert('Ошибка загрузки сортов');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchVarieties();
  }, []);

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

  const openCreateForm = () => {
    setEditingId(null);
    setName('');
    setDescription('');
    setPriority('high');
    setImageUrl(undefined);
    setShowForm(true);
  };

  const openEditForm = (v: Variety) => {
    setEditingId(v.id);
    setName(v.name);
    setDescription(v.description || '');
    setPriority(v.priority);
    setImageUrl(v.image_url || undefined);
    setShowForm(true);
  };

  const handleSave = async () => {
    if (!name.trim()) return;
    const payload = {
      name,
      description,
      priority,
      image_url: imageUrl || '',
    };
    try {
      if (editingId) {
        await api.put(`/api/varieties/${editingId}`, payload);
      } else {
        await api.post('/api/varieties', payload);
      }
      setShowForm(false);
      fetchVarieties();
    } catch {
      alert('Ошибка сохранения сорта');
    }
  };

  const handleDelete = async (id: string) => {
    if (!window.confirm('Удалить сорт?')) return;
    try {
      await api.delete(`/api/varieties/${id}`);
      fetchVarieties();
    } catch {
      alert('Ошибка удаления сорта');
    }
  };

  return (
    <div>
      <div style={styles.header}>
        <h2>Сорта</h2>
        <button onClick={openCreateForm} style={styles.addBtn}>Добавить сорт</button>
      </div>

      {loading && <p>Загрузка...</p>}

      <div style={styles.list}>
        {varieties.map(v => (
          <div key={v.id} style={styles.card}>
            <div style={{ display: 'flex', alignItems: 'center', gap: 12 }}>
              {v.image_url && <img src={v.image_url} style={styles.thumb} alt={v.name} />}
              <div style={{ flex: 1 }}>
                <strong>{v.name}</strong>
                <span style={{ marginLeft: 8, color: v.priority === 'high' ? '#ef4444' : '#16a34a' }}>
                  {v.priority === 'high' ? '⚠️ Высокий' : '✅ Низкий'}
                </span>
                {v.description && <p style={{ color: '#4b5563', margin: '4px 0 0' }}>{v.description}</p>}
              </div>
              <button onClick={() => openEditForm(v)} style={styles.editBtn}>✏️</button>
              <button onClick={() => handleDelete(v.id)} style={styles.deleteBtn}>🗑️</button>
            </div>
          </div>
        ))}
      </div>

      {showForm && (
        <div style={styles.modalOverlay} onClick={() => setShowForm(false)}>
          <div style={styles.modal} onClick={e => e.stopPropagation()}>
            <h3>{editingId ? 'Редактировать сорт' : 'Новый сорт'}</h3>
            <input
              value={name}
              onChange={e => setName(e.target.value)}
              placeholder="Название"
              style={styles.input}
            />
            <textarea
              value={description}
              onChange={e => setDescription(e.target.value)}
              placeholder="Описание"
              style={styles.textarea}
              rows={3}
            />
            <select value={priority} onChange={e => setPriority(e.target.value as 'high' | 'low')} style={styles.input}>
              <option value="high">Высокий</option>
              <option value="low">Низкий</option>
            </select>
            <ImageUploader
              imageUrl={imageUrl}
              onUpload={async (file) => {
                const url = await uploadImage(file);
                if (url) setImageUrl(url);
              }}
              onRemove={() => setImageUrl(undefined)}
            />
            <div style={styles.modalActions}>
              <button onClick={() => setShowForm(false)} style={styles.cancelBtn}>Отмена</button>
              <button onClick={handleSave} style={styles.saveBtn}>Сохранить</button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
};

const styles = {
  header: { display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 24 },
  addBtn: { background: '#16a34a', color: 'white', border: 'none', padding: '8px 16px', borderRadius: 8, cursor: 'pointer' },
  list: { display: 'flex', flexDirection: 'column' as const, gap: 12 },
  card: { background: 'white', padding: 16, borderRadius: 12, boxShadow: '0 1px 3px rgba(0,0,0,0.1)' },
  thumb: { width: 40, height: 40, borderRadius: 6, objectFit: 'cover' as const },
  editBtn: { background: '#eab308', border: 'none', padding: '6px 12px', borderRadius: 6, cursor: 'pointer', marginLeft: 8 },
  deleteBtn: { background: '#ef4444', border: 'none', padding: '6px 12px', borderRadius: 6, cursor: 'pointer', color: 'white', marginLeft: 4 },
  modalOverlay: {
    position: 'fixed' as const, top: 0, left: 0, right: 0, bottom: 0,
    backgroundColor: 'rgba(0,0,0,0.5)', display: 'flex', alignItems: 'center', justifyContent: 'center', zIndex: 1000
  },
  modal: { backgroundColor: '#fff', borderRadius: 16, padding: 24, width: '90%', maxWidth: 480, boxShadow: '0 20px 25px -5px rgba(0,0,0,0.1)' },
  input: { width: '100%', padding: 10, marginBottom: 12, borderRadius: 8, border: '1px solid #d1d5db', fontSize: 14, boxSizing: 'border-box' as const },
  textarea: { width: '100%', padding: 10, marginBottom: 12, borderRadius: 8, border: '1px solid #d1d5db', fontSize: 14, resize: 'vertical' as const, boxSizing: 'border-box' as const },
  modalActions: { display: 'flex', justifyContent: 'flex-end', gap: 12, marginTop: 16 },
  cancelBtn: { padding: '8px 16px', background: '#f3f4f6', border: '1px solid #e5e7eb', borderRadius: 8, cursor: 'pointer' },
  saveBtn: { background: '#16a34a', color: 'white', border: 'none', padding: '8px 16px', borderRadius: 8, cursor: 'pointer' },
};