import { useState } from 'react';
import { useNavigate, Link } from 'react-router-dom';
import { api } from '../api/client';
import { useAuth } from '../context/AuthContext';
import { decodeToken } from '../utils/token';

export const Login = () => {
  const [login, setLogin] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');
  const navigate = useNavigate();
  const { login: authLogin } = useAuth();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    try {
      const res = await api.post('/login', { login, password });
      const token = res.data.token;
      authLogin(token);
      // Декодируем роль прямо сейчас, не дожидаясь обновления контекста
      const payload = decodeToken();
      if (payload?.role === 'admin') {
        navigate('/dashboard');
      } else {
        navigate('/create-report');
      }
    } catch (err) {
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
          onChange={(e) => setLogin(e.target.value)}
          style={styles.input}
          required
        />
        <input
          type="password"
          placeholder="Пароль"
          value={password}
          onChange={(e) => setPassword(e.target.value)}
          style={styles.input}
          required
        />
        {error && <div style={styles.error}>{error}</div>}
        <button type="submit" style={styles.button}>Войти</button>
        <p style={styles.link}>
          Нет аккаунта? <Link to="/register">Зарегистрироваться</Link>
        </p>
      </form>
    </div>
  );
};

const styles = {
  page: {
    minHeight: '100vh',
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'center',
    background: '#f4f6f8',
  },
  card: {
    background: 'white',
    padding: '32px',
    borderRadius: '16px',
    width: '360px',
    boxShadow: '0 4px 12px rgba(0,0,0,0.05)',
  },
  input: {
    width: '100%',
    padding: '10px 12px',
    marginBottom: '16px',
    border: '1px solid #d1d5db',
    borderRadius: '8px',
    fontSize: '14px',
  },
  button: {
    width: '100%',
    background: '#16a34a',
    color: 'white',
    border: 'none',
    padding: '10px',
    borderRadius: '8px',
    cursor: 'pointer',
    fontWeight: 500,
  },
  error: {
    color: '#ef4444',
    marginBottom: '12px',
    fontSize: '14px',
  },
  link: {
    marginTop: '16px',
    textAlign: 'center' as const,
    fontSize: '14px',
  },
};