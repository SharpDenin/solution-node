import { Outlet, useNavigate } from 'react-router-dom';

export const WorkerLayout = () => {
  const navigate = useNavigate();

  const handleLogout = () => {
    navigate('/');  // на главную, где сессия сбросится
  };

  return (
    <div style={styles.container}>
      <header style={styles.header}>
        <span style={styles.logo}>📋 Чек-лист работника</span>
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
  header: { background: 'white', padding: '16px 24px', borderBottom: '1px solid #e5e7eb', display: 'flex', justifyContent: 'space-between', alignItems: 'center' },
  logo: { fontSize: '20px', fontWeight: 700, color: '#16a34a' },
  main: { padding: '24px', maxWidth: '800px', margin: '0 auto' },
  logoutBtn: { background: '#ef4444', color: 'white', border: 'none', padding: '8px 16px', borderRadius: '8px', cursor: 'pointer' },
};