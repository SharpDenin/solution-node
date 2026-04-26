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
            <div onClick={() => toggleDescription(v.id)} style={styles.cardBody}>
              <div style={styles.titleRow}>
                <strong>{v.name}</strong>
                {v.priority === 'high' && <span style={styles.highBadge}>⚠️ Высокий</span>}
                {v.priority === 'low' && <span style={styles.lowBadge}>✅ Низкий</span>}
              </div>
              {v.image_url && (
                <div style={styles.imageWrapper}>
                  <img src={v.image_url} alt={v.name} style={styles.image} />
                </div>
              )}
              {expandedId === v.id && v.description && (
                <div style={styles.description}>{v.description}</div>
              )}
            </div>
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
    display: 'flex',
    flexDirection: 'column' as const,
    gap: '12px',
  },
  cardBody: {
    cursor: 'pointer',
  },
  titleRow: {
    display: 'flex',
    alignItems: 'center',
    gap: '8px',
    marginBottom: '8px',
  },
  highBadge: { color: '#ef4444', fontSize: '14px', fontWeight: 600 },
  lowBadge: { color: '#16a34a', fontSize: '14px', fontWeight: 600 },
  imageWrapper: {
    marginTop: '8px',
    marginBottom: '8px',
    borderRadius: '12px',
    overflow: 'hidden',
    width: '100%',
    display: 'flex',
    justifyContent: 'center',
  },
  image: {
    maxWidth: '100%',
    maxHeight: '200px',
    objectFit: 'contain' as const,
    borderRadius: '12px',
  },
  description: {
    marginTop: '8px',
    padding: '12px',
    backgroundColor: '#f9fafb',
    borderRadius: '8px',
    fontSize: '14px',
    color: '#4b5563',
  },
  selectBtn: {
    background: '#16a34a',
    color: 'white',
    border: 'none',
    padding: '10px 16px',
    borderRadius: '8px',
    cursor: 'pointer',
    fontWeight: 500,
    alignSelf: 'stretch',
  },
};