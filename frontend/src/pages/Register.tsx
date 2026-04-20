import { useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { api } from '../api/client';

export const Register = () => {
  const [fullName, setFullName] = useState('');
  const [login, setLogin] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');
  const [success, setSuccess] = useState(false);
  const navigate = useNavigate();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    setSuccess(false);
    try {
      await api.post('/api/register', {
        full_name: fullName,
        login,
        password,
      });
      setSuccess(true);
      setTimeout(() => navigate('/login'), 2000);
    } catch (err) {
      setError('Ошибка регистрации. Возможно, логин уже занят.');
    }
  };

  return (
    <div style={styles.page}>
      <form onSubmit={handleSubmit} style={styles.card}>
        <h2>Регистрация</h2>
        <input
          type="text"
          placeholder="ФИО"
          value={fullName}
          onChange={(e) => setFullName(e.target.value)}
          style={styles.input}
          required
        />
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
        {success && <div style={styles.success}>Регистрация успешна! Перенаправление на вход...</div>}
        <button type="submit" style={styles.button}>Зарегистрироваться</button>
        <p style={styles.link}>
          Уже есть аккаунт? <Link to="/login">Войти</Link>
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
  success: {
    color: '#16a34a',
    marginBottom: '12px',
    fontSize: '14px',
  },
  link: {
    marginTop: '16px',
    textAlign: 'center' as const,
    fontSize: '14px',
  },
};