import { useNavigate } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';

export const ThankYou = () => {
  const { logout } = useAuth();
  const navigate = useNavigate();

  const handleExit = () => {
    logout();
    navigate('/');
  };

  return (
    <div style={styles.container}>
      <div style={styles.card}>
        <h2>✅ Спасибо за отчёт!</h2>
        <p>Ваш отчёт успешно отправлен.</p>
        <button onClick={handleExit} style={styles.exitBtn}>Выход</button>
      </div>
    </div>
  );
};

const styles = {
  container: {
    display: 'flex',
    justifyContent: 'center',
    alignItems: 'center',
    minHeight: '60vh',
  },
  card: {
    background: 'white',
    padding: '32px',
    borderRadius: '16px',
    textAlign: 'center' as const,
    boxShadow: '0 4px 12px rgba(0,0,0,0.05)',
  },
  exitBtn: {
    display: 'inline-block',
    marginTop: '16px',
    background: '#ef4444',
    color: 'white',
    padding: '10px 20px',
    borderRadius: '8px',
    border: 'none',
    cursor: 'pointer',
  },
};