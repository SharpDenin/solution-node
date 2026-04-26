import { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { api } from '../api/client';
import { useAuth } from '../context/AuthContext';
import type { Checklist } from '../types';

export const HomePage = () => {
  const [checklists, setChecklists] = useState<Checklist[]>([]);
  const { token, logout } = useAuth();
  const navigate = useNavigate();

  // При заходе на главную сбрасываем токен — чистая публичная страница
  useEffect(() => {
    logout();
  }, [logout]);

  useEffect(() => {
    api.get('/api/checklists')
      .then(res => setChecklists(res.data))
      .catch(console.error);
  }, []);

  const handleChecklistClick = (id: string) => {
    // Токена нет – всегда редирект на логин
    if (!token) {
      navigate(`/login?returnUrl=/checklist/${id}`);
    } else {
      navigate(`/checklist/${id}`);
    }
  };

  return (
    <div style={styles.page}>
      {/* Шапка с названием компании — слева сверху */}
      <div style={styles.header}>
        <div style={styles.companyName}>ООО "Агроном-сад"</div>
      </div>

      {/* Основной контент — по центру */}
      <div style={styles.content}>
        <h1 style={styles.title}>Выберите необходимый функционал:</h1>

        <div style={styles.checklistPanel}>
          {checklists.map(cl => (
            <div
              key={cl.id}
              style={styles.card}
              onClick={() => handleChecklistClick(cl.id)}
            >
              <strong>{cl.name}</strong>
            </div>
          ))}
        </div>

        <button
          onClick={() => navigate('/dashboard')}
          style={styles.reportsBtn}
        >
          Отчёты
        </button>
      </div>
    </div>
  );
};

const styles = {
  page: {
    minHeight: '100vh',
    background: '#f4f6f8', // фон всей страницы
  },
  header: {
    width: '100%',
    padding: '20px 24px',
    boxSizing: 'border-box' as const,
  },
  companyName: {
    fontSize: '18px',
    fontWeight: 700 as const,
    color: '#111827',
    textAlign: 'left' as const,
  },
  content: {
    maxWidth: '600px',
    margin: '0 auto',
    padding: '0 16px 40px',
    textAlign: 'center' as const,
  },
  title: {
    marginBottom: '36px',
    fontSize: '24px',
    fontWeight: 600 as const,
    color: '#111827',
  },
  checklistPanel: {
    background: '#f9fafb',
    border: '1px solid #e5e7eb',
    borderRadius: '16px',
    padding: '20px',
    marginBottom: '32px',
    display: 'flex',
    flexDirection: 'column' as const,
    gap: '16px',
  },
  card: {
    background: 'white',
    padding: '18px 16px',
    borderRadius: '12px',
    cursor: 'pointer',
    boxShadow: '0 2px 6px rgba(0,0,0,0.06)',
    transition: 'box-shadow 0.2s',
    textAlign: 'left' as const,
    fontSize: '16px',
  },
  reportsBtn: {
    background: '#16a34a',
    color: 'white',
    border: 'none',
    padding: '12px 24px',
    borderRadius: '8px',
    fontSize: '16px',
    cursor: 'pointer',
    marginTop: '8px',
  },
};