import { useState } from 'react';
import { api } from '../api/client';

export const UserCreate = () => {
  const [fullName, setFullName] = useState('');
  const [login, setLogin] = useState('');
  const [password, setPassword] = useState('');
  const [role, setRole] = useState('node'); // default to node specialist
  const [message, setMessage] = useState('');
  const [error, setError] = useState('');

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setMessage('');
    setError('');
    try {
      await api.post('/api/register', {
        full_name: fullName,
        login,
        password,
        role, // sends the string role directly (backend expects "admin", "node", "phenophase")
      });
      setMessage('Пользователь успешно создан');
      setFullName('');
      setLogin('');
      setPassword('');
    } catch (err: any) {
      setError(err.response?.data || 'Ошибка при создании');
    }
  };

  return (
    <div style={styles.wrapper}>
      <h2 style={styles.title}>Создать пользователя</h2>
      <form onSubmit={handleSubmit} style={styles.form}>
        <input
          type="text"
          placeholder="ФИО"
          value={fullName}
          onChange={e => setFullName(e.target.value)}
          style={styles.input}
          required
        />
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
        <select
          value={role}
          onChange={e => setRole(e.target.value)}
          style={styles.input}
        >
          <option value="admin">Администратор</option>
          <option value="node">Специалист по растворным узлам</option>
          <option value="phenophase">Специалист по фенофазам</option>
        </select>
        <button type="submit" style={styles.button}>Создать</button>
        {message && <div style={styles.success}>{message}</div>}
        {error && <div style={styles.error}>{error}</div>}
      </form>
    </div>
  );
};

const styles = {
  wrapper: {
    display: 'flex',
    flexDirection: 'column' as const,
    alignItems: 'center',
    justifyContent: 'center',
    minHeight: 'calc(100vh - 100px)', // vertical centering
    padding: '40px 16px 0',
  },
  title: {
    marginBottom: '32px',
    fontSize: '24px',
    fontWeight: 600,
    color: '#111827',
    textAlign: 'center' as const,
  },
  form: {
    width: '100%',
    maxWidth: '420px',
    background: 'white',
    padding: '28px',
    borderRadius: '12px',
    boxShadow: '0 2px 8px rgba(0,0,0,0.05)',
    display: 'flex',
    flexDirection: 'column' as const,
    gap: '16px',
  },
  input: {
    padding: '10px 14px',
    border: '1px solid #d1d5db',
    borderRadius: '8px',
    fontSize: '14px',
    outline: 'none',
    boxSizing: 'border-box' as const,
  },
  button: {
    background: '#16a34a',
    color: 'white',
    border: 'none',
    padding: '10px',
    borderRadius: '8px',
    cursor: 'pointer',
    fontWeight: 500,
    fontSize: '14px',
  },
  success: {
    color: '#16a34a',
    fontSize: '14px',
    textAlign: 'center' as const,
  },
  error: {
    color: '#ef4444',
    fontSize: '14px',
    textAlign: 'center' as const,
  },
};