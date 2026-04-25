import { useState } from 'react';
import { useNavigate, useSearchParams } from 'react-router-dom';
import { api } from '../api/client';
import { useAuth } from '../context/AuthContext';
import { decodeToken } from '../utils/token';

export const Login = () => {
  const [login, setLogin] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');
  const navigate = useNavigate();
  const { login: authLogin } = useAuth();
  const [searchParams] = useSearchParams();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    try {
      const res = await api.post('/api/login', { login, password });
      const token = res.data.token;
      authLogin(token);
      const payload = decodeToken();
      const returnUrl = searchParams.get('returnUrl');
      if (payload?.role === 'admin') {
        navigate(returnUrl || '/dashboard');
      } else {
        navigate(returnUrl || '/');
      }
    } catch {
      setError('Неверный логин или пароль');
    }
  };

  return (
    <div style={styles.page}>
      <form onSubmit={handleSubmit} style={styles.card}>
        <h2>Вход в систему</h2>
        <input
          type="text"
          placeholder="Логин"
          value={login}
          onChange={e => setLogin(e.target.value)}
          style={styles.input}
          required
        />
        <input
          type="password"
          placeholder="Пароль"
          value={password}
          onChange={e => setPassword(e.target.value)}
          style={styles.input}
          required
        />
        {error && <div style={styles.error}>{error}</div>}
        <button type="submit" style={styles.button}>Войти</button>
      </form>
    </div>
  );
};

const styles = {
  page: { minHeight: '100vh', display: 'flex', alignItems: 'center', justifyContent: 'center', background: '#f4f6f8' },
  card: { background: 'white', padding: 32, borderRadius: 16, width: 360, boxShadow: '0 4px 12px rgba(0,0,0,0.05)' },
  input: { width: '100%', padding: '10px 12px', marginBottom: 16, border: '1px solid #d1d5db', borderRadius: 8, fontSize: 14, boxSizing: 'border-box' as const },
  button: { width: '100%', background: '#16a34a', color: 'white', border: 'none', padding: 10, borderRadius: 8, cursor: 'pointer', fontWeight: 500 },
  error: { color: '#ef4444', marginBottom: 12, fontSize: 14 },
};