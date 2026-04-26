import { useEffect, useState } from 'react';
import { useParams, useNavigate, useSearchParams } from 'react-router-dom';
import { api } from '../api/client';
import type { Phenophase } from '../types';

export const PhenophaseSelectPage = () => {
  const { id: checklistId } = useParams<{ id: string }>();
  const [searchParams] = useSearchParams();
  const varietyId = searchParams.get('varietyId');
  const navigate = useNavigate();

  const [phenophases, setPhenophases] = useState<Phenophase[]>([]);
  const [expandedId, setExpandedId] = useState<string | null>(null);

  useEffect(() => {
    api.get('/api/phenophases')
      .then(res => setPhenophases(res.data))
      .catch(console.error);
  }, []);

  const handleSelect = (phenophaseId: string) => {
    navigate(`/checklist/${checklistId}/fill?varietyId=${varietyId}&phenophaseId=${phenophaseId}`);
  };

  const toggleDescription = (id: string) => {
    setExpandedId(expandedId === id ? null : id);
  };

  return (
    <div style={styles.container}>
      <h2>Выберите фенофазу</h2>
      <div style={styles.list}>
        {phenophases.map(p => (
          <div key={p.id} style={styles.card}>
            <div onClick={() => toggleDescription(p.id)} style={styles.cardBody}>
              <strong style={{ display: 'block', marginBottom: 8 }}>{p.name}</strong>
              {p.image_url && (
                <div style={styles.imageWrapper}>
                  <img src={p.image_url} alt={p.name} style={styles.image} />
                </div>
              )}
              {expandedId === p.id && p.description && (
                <div style={styles.description}>{p.description}</div>
              )}
            </div>
            <button onClick={() => handleSelect(p.id)} style={styles.selectBtn}>
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
  imageWrapper: {
    marginTop: '8px',
    marginBottom: '8px',
    borderRadius: '12px',
    overflow: 'hidden',
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