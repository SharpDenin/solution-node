import { Outlet, NavLink, useNavigate } from 'react-router-dom';
import { useAuth } from '../../context/AuthContext';

const linkStyle = ({ isActive }: { isActive: boolean }) => ({
  textDecoration: 'none',
  color: isActive ? '#16a34a' : '#111827',
  backgroundColor: isActive ? '#e8f5ee' : 'transparent',
  padding: '8px 12px',
  borderRadius: '8px',
  transition: 'all 0.2s',
});

export const AdminLayout = () => {
  const navigate = useNavigate();

  const handleLogout = () => {
    navigate('/');  // уходим на главную, токен cбросится там
  };

  return (
    <div style={styles.container}>
      <aside style={styles.sidebar}>
        <div style={styles.logo}>📋 Растворные узлы</div>
        <nav style={styles.nav}>
          <NavLink to="/dashboard" style={linkStyle}>📊 Отчёты</NavLink>
          <NavLink to="/questions" style={linkStyle}>❓ Вопросы</NavLink>
          <NavLink to="/admin/users/create" style={linkStyle}>👤 Пользователи</NavLink>
        </nav>
        <button onClick={handleLogout} style={styles.logoutBtn}>Выход</button>
      </aside>
      <main style={styles.main}>
        <Outlet />
      </main>
    </div>
  );
};

const styles = {
  container: { display: 'flex', minHeight: '100vh', background: '#f4f6f8' },
  sidebar: { width: '260px', background: 'white', borderRight: '1px solid #e5e7eb', display: 'flex', flexDirection: 'column' as const, padding: '24px 16px', justifyContent: 'space-between' },
  logo: { fontSize: '20px', fontWeight: 700, color: '#16a34a', marginBottom: '32px', paddingLeft: '8px' },
  nav: { display: 'flex', flexDirection: 'column' as const, gap: '8px' },
  main: { flex: 1, padding: '24px', overflowY: 'auto' as const },
  logoutBtn: { background: '#ef4444', color: 'white', border: 'none', padding: '10px', borderRadius: '8px', cursor: 'pointer', fontWeight: 500, marginTop: 'auto' },
};