import { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { api } from '../api/client';
import { useAuth } from '../context/AuthContext';
import type { Checklist } from '../types';

export const HomePage = () => {
  const [checklists, setChecklists] = useState<Checklist[]>([]);
  const { token, role } = useAuth();
  const navigate = useNavigate();

  useEffect(() => {
    api.get('/api/checklists')
      .then(res => setChecklists(res.data))
      .catch(console.error);
  }, []);

  const handleChecklistClick = (id: string) => {
    if (!token) {
      navigate(`/login?returnUrl=/checklist/${id}`);
    } else {
      navigate(`/checklist/${id}`);
    }
  };

  return (
    <div style={styles.container}>
      <div style={styles.header}>
        <h1>Доступные чек-листы</h1>
        <div style={styles.buttons}>
          {!token && (
            <button onClick={() => navigate('/login')} style={styles.loginBtn}>Войти</button>
          )}
          {role === 'admin' && (
            <button onClick={() => navigate('/dashboard')} style={styles.reportsBtn}>
              Отчёты
            </button>
          )}
        </div>
      </div>
      <div style={styles.list}>
        {checklists.map(cl => (
          <div
            key={cl.id}
            style={styles.card}
            onClick={() => handleChecklistClick(cl.id)}
          >
            <strong>{cl.name}</strong>
            <div style={styles.code}>{cl.code}</div>
          </div>
        ))}
      </div>
    </div>
  );
};

const styles = {
  container: {
    maxWidth: '600px',
    margin: '40px auto',
    padding: '0 16px',
  },
  header: {
    display: 'flex',
    justifyContent: 'space-between',
    alignItems: 'center',
  },
  buttons: {
    display: 'flex',
    gap: '12px',
  },
  loginBtn: {
    background: '#2563eb',
    color: 'white',
    border: 'none',
    padding: '8px 16px',
    borderRadius: '8px',
    cursor: 'pointer',
  },
  reportsBtn: {
    background: '#16a34a',
    color: 'white',
    border: 'none',
    padding: '8px 16px',
    borderRadius: '8px',
    cursor: 'pointer',
  },
  list: {
    display: 'flex',
    flexDirection: 'column' as const,
    gap: '12px',
    marginTop: '24px',
  },
  card: {
    background: 'white',
    padding: '16px',
    borderRadius: '12px',
    cursor: 'pointer',
    boxShadow: '0 2px 6px rgba(0,0,0,0.08)',
    transition: 'box-shadow 0.2s',
  },
  code: {
    color: '#6b7280',
    fontSize: '14px',
  },
};