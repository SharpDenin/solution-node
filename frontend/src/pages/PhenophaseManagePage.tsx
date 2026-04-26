import { useEffect, useState } from 'react';
import { api } from '../api/client';
import { ImageUploader } from '../components/ImageUploader';
import type { Phenophase } from '../types';

export const PhenophaseManagePage = () => {
  const [phenophases, setPhenophases] = useState<Phenophase[]>([]);
  const [loading, setLoading] = useState(false);

  const [showForm, setShowForm] = useState(false);
  const [editingId, setEditingId] = useState<string | null>(null);
  const [name, setName] = useState('');
  const [description, setDescription] = useState('');
  const [imageUrl, setImageUrl] = useState<string | undefined>(undefined);
  const [orderIndex, setOrderIndex] = useState(1);

  const fetchPhenophases = async () => {
    setLoading(true);
    try {
      const res = await api.get('/api/phenophases');
      setPhenophases(res.data);
    } catch {
      alert('Ошибка загрузки фенофаз');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchPhenophases();
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
    setImageUrl(undefined);
    setOrderIndex(phenophases.length + 1);
    setShowForm(true);
  };

  const openEditForm = (p: Phenophase) => {
    setEditingId(p.id);
    setName(p.name);
    setDescription(p.description || '');
    setImageUrl(p.image_url || undefined);
    setOrderIndex(p.order_index);
    setShowForm(true);
  };

  const handleSave = async () => {
    if (!name.trim()) return;
    const payload = {
      name,
      description,
      image_url: imageUrl || '',
      order_index: orderIndex,
    };
    try {
      if (editingId) {
        await api.put(`/api/phenophases/${editingId}`, payload);
      } else {
        await api.post('/api/phenophases', payload);
      }
      setShowForm(false);
      fetchPhenophases();
    } catch {
      alert('Ошибка сохранения фенофазы');
    }
  };

  const handleDelete = async (id: string) => {
    if (!window.confirm('Удалить фенофазу?')) return;
    try {
      await api.delete(`/api/phenophases/${id}`);
      fetchPhenophases();
    } catch {
      alert('Ошибка удаления фенофазы');
    }
  };

  return (
    <div>
      <div style={styles.header}>
        <h2>Фенофазы</h2>
        <button onClick={openCreateForm} style={styles.addBtn}>Добавить фенофазу</button>
      </div>

      {loading && <p>Загрузка...</p>}

      <div style={styles.list}>
        {phenophases.map(p => (
          <div key={p.id} style={styles.card}>
            <div style={{ display: 'flex', alignItems: 'center', gap: 12 }}>
              <strong style={{ minWidth: 24 }}>{p.order_index}.</strong>
              {p.image_url && <img src={p.image_url} style={styles.thumb} alt={p.name} />}
              <div style={{ flex: 1 }}>
                <strong>{p.name}</strong>
                {p.description && <p style={{ color: '#4b5563', margin: '4px 0 0' }}>{p.description}</p>}
              </div>
              <button onClick={() => openEditForm(p)} style={styles.editBtn}>✏️</button>
              <button onClick={() => handleDelete(p.id)} style={styles.deleteBtn}>🗑️</button>
            </div>
          </div>
        ))}
      </div>

      {showForm && (
        <div style={styles.modalOverlay} onClick={() => setShowForm(false)}>
          <div style={styles.modal} onClick={e => e.stopPropagation()}>
            <h3>{editingId ? 'Редактировать фенофазу' : 'Новая фенофаза'}</h3>
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
            <div style={{ marginBottom: 12 }}>
              <label>Порядковый номер</label>
              <input
                type="number"
                value={orderIndex}
                onChange={e => setOrderIndex(parseInt(e.target.value) || 1)}
                style={styles.input}
              />
            </div>
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