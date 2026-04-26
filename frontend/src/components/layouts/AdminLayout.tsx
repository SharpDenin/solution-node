import { Outlet, NavLink, useNavigate } from 'react-router-dom';
import { useAuth } from '../../context/AuthContext';

const linkStyle = ({ isActive }: { isActive: boolean }) => ({
  textDecoration: 'none',
  color: isActive ? '#16a34a' : '#111827',
  backgroundColor: isActive ? '#e8f5ee' : 'transparent',
  padding: '10px 12px',
  borderRadius: '8px',
  transition: 'all 0.2s',
  fontSize: '16px',
  fontWeight: 500,
  display: 'block',
});

export const AdminLayout = () => {
  const { fullName, position } = useAuth();
  const navigate = useNavigate();

  const handleLogout = () => {
    navigate('/');
  };

  return (
    <div style={styles.container}>
      <aside style={styles.sidebar}>
        <div style={styles.logo}>Панель администратора</div>
        <nav style={styles.nav}>
          <NavLink to="/dashboard" style={linkStyle}>📊 Отчёты</NavLink>
          <NavLink to="/questions" style={linkStyle}>❓ Вопросы</NavLink>
          <NavLink to="/admin/users/create" style={linkStyle}>👤 Пользователи</NavLink>
        </nav>

        {/* Информация о пользователе – снизу слева */}
        <div style={styles.userBlock}>
          <div style={styles.userName}>{fullName || 'Пользователь'}</div>
          {position && <div style={styles.userPosition}>{position}</div>}
        </div>

        <button onClick={handleLogout} style={styles.logoutBtn}>Выход</button>
      </aside>
      <main style={styles.main}>
        <Outlet />
      </main>
    </div>
  );
};

const styles = {
  container: {
    display: 'flex',
    minHeight: '100vh',
    background: '#f4f6f8',
  },
  sidebar: {
    width: '260px',
    background: 'white',
    borderRight: '1px solid #e5e7eb',
    display: 'flex',
    flexDirection: 'column' as const,
    padding: '24px 16px',
    justifyContent: 'space-between',
  },
  logo: {
    fontSize: '20px',
    fontWeight: 700,
    color: '#16a34a',
    marginBottom: '32px',
    paddingLeft: '8px',
  },
  nav: {
    display: 'flex',
    flexDirection: 'column' as const,
    gap: '8px',
    flex: 1,
  },
  userBlock: {
    marginTop: 'auto',
    padding: '12px 8px',
    borderTop: '1px solid #e5e7eb',
    marginBottom: '12px',
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
    display: 'inline-block',
    marginTop: '4px',
  },
  main: {
    flex: 1,
    padding: '24px',
    overflowY: 'auto' as const,
  },
  logoutBtn: {
    background: '#ef4444',
    color: 'white',
    border: 'none',
    padding: '10px',
    borderRadius: '8px',
    cursor: 'pointer',
    fontWeight: 500,
  },
};