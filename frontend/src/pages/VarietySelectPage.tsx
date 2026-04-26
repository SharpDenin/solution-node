import { useEffect, useState } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { api } from '../api/client';
import type { Variety } from '../types';

export const VarietySelectPage = () => {
  const { id: checklistId } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const [varieties, setVarieties] = useState<Variety[]>([]);
  const [expandedId, setExpandedId] = useState<string | null>(null);

  useEffect(() => {
    api.get('/api/varieties')
      .then(res => setVarieties(res.data))
      .catch(console.error);
  }, []);

  const handleSelect = (varietyId: string) => {
    navigate(`/checklist/${checklistId}/phenophase?varietyId=${varietyId}`);
  };

  const toggleDescription = (id: string) => {
    setExpandedId(expandedId === id ? null : id);
  };

  return (
    <div style={styles.container}>
      <h2>Выберите сорт</h2>
      <div style={styles.list}>
        {varieties.map(v => (
          <div key={v.id} style={styles.card}>
            <div onClick={() => toggleDescription(v.id)} style={styles.cardHeader}>
              <strong>{v.name}</strong>
              {v.priority === 'high' && <span style={styles.highBadge}>⚠️ Высокий</span>}
              {v.priority === 'low' && <span style={styles.lowBadge}>✅ Низкий</span>}
              {v.image_url && (
                <img src={v.image_url} style={styles.thumb} alt={v.name} />
              )}
            </div>
            {expandedId === v.id && v.description && (
              <div style={styles.description}>{v.description}</div>
            )}
            <button onClick={() => handleSelect(v.id)} style={styles.selectBtn}>
              Выбрать
            </button>
          </div>
        ))}
      </div>
    </div>
  );
};

const styles = {
  container: { maxWidth: '600px', margin: '0 auto', padding: '24px' },
  list: { display: 'flex', flexDirection: 'column' as const, gap: '16px' },
  card: {
    background: 'white',
    borderRadius: '12px',
    padding: '16px',
    boxShadow: '0 2px 6px rgba(0,0,0,0.08)',
  },
  cardHeader: {
    display: 'flex',
    alignItems: 'center',
    gap: '8px',
    cursor: 'pointer',
  },
  highBadge: { color: '#ef4444', fontSize: '14px', fontWeight: 600 },
  lowBadge: { color: '#16a34a', fontSize: '14px', fontWeight: 600 },
  thumb: { width: '40px', height: '40px', borderRadius: '6px', objectFit: 'cover' as const },
  description: {
    marginTop: '12px',
    padding: '8px',
    backgroundColor: '#f9fafb',
    borderRadius: '8px',
    fontSize: '14px',
    color: '#4b5563',
  },
  selectBtn: {
    marginTop: '12px',
    background: '#16a34a',
    color: 'white',
    border: 'none',
    padding: '8px 16px',
    borderRadius: '8px',
    cursor: 'pointer',
    fontWeight: 500,
  },
};