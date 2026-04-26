import { Outlet, useNavigate } from 'react-router-dom';
import { useAuth } from '../../context/AuthContext';

export const WorkerLayout = () => {
  const { fullName, position } = useAuth();
  const navigate = useNavigate();

  const handleLogout = () => {
    navigate('/');
  };

  return (
    <div style={styles.container}>
      <header style={styles.header}>
        <div style={styles.leftGroup}>
          <span style={styles.logo}>📋 Чек-лист</span>
          {fullName && (
            <div style={styles.userInfo}>
              <span style={styles.userName}>{fullName}</span>
              {position && <span style={styles.userPosition}>{position}</span>}
            </div>
          )}
        </div>
        <button onClick={handleLogout} style={styles.logoutBtn}>Выход</button>
      </header>
      <main style={styles.main}>
        <Outlet />
      </main>
    </div>
  );
};

const styles = {
  container: { minHeight: '100vh', background: '#f4f6f8' },
  header: {
    background: 'white',
    padding: '12px 24px',
    borderBottom: '1px solid #e5e7eb',
    display: 'flex',
    justifyContent: 'space-between',
    alignItems: 'center',
  },
  leftGroup: {
    display: 'flex',
    alignItems: 'center',
    gap: '20px',
  },
  logo: { fontSize: '20px', fontWeight: 700, color: '#16a34a' },
  userInfo: {
    display: 'flex',
    alignItems: 'center',
    gap: '8px',
  },
  userName: {
    fontWeight: 600,
    fontSize: '15px',
    color: '#111827',
  },
  userPosition: {
    fontSize: '13px',
    color: '#6b7280',
    backgroundColor: '#f3f4f6',
    padding: '2px 10px',
    borderRadius: '20px',
  },
  main: { padding: '24px', maxWidth: '800px', margin: '0 auto' },
  logoutBtn: {
    background: '#ef4444',
    color: 'white',
    border: 'none',
    padding: '8px 16px',
    borderRadius: '8px',
    cursor: 'pointer',
  },
};