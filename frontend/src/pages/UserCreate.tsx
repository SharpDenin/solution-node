import { useState } from 'react';
import { api } from '../api/client';

export const UserCreate = () => {
  const [fullName, setFullName] = useState('');
  const [login, setLogin] = useState('');
  const [password, setPassword] = useState('');
  const [role, setRole] = useState<'worker' | 'admin'>('worker');
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
        role,
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
    <div>
      <h2>Создать пользователя</h2>
      <form onSubmit={handleSubmit} style={styles.form}>
        <input type="text" placeholder="ФИО" value={fullName} onChange={e => setFullName(e.target.value)} style={styles.input} required />
        <input type="text" placeholder="Логин" value={login} onChange={e => setLogin(e.target.value)} style={styles.input} required />
        <input type="password" placeholder="Пароль" value={password} onChange={e => setPassword(e.target.value)} style={styles.input} required />
        <select value={role} onChange={e => setRole(e.target.value as any)} style={styles.input}>
          <option value="worker">Рабочий</option>
          <option value="admin">Администратор</option>
        </select>
        <button type="submit" style={styles.button}>Создать</button>
        {message && <div style={styles.success}>{message}</div>}
        {error && <div style={styles.error}>{error}</div>}
      </form>
    </div>
  );
};

const styles = {
  form: { maxWidth: 400, background: 'white', padding: 24, borderRadius: 12, boxShadow: '0 2px 8px rgba(0,0,0,0.05)' },
  input: { width: '100%', padding: '10px 12px', marginBottom: 16, border: '1px solid #d1d5db', borderRadius: 8, fontSize: 14, boxSizing: 'border-box' as const },
  button: { width: '100%', background: '#16a34a', color: 'white', border: 'none', padding: 10, borderRadius: 8, cursor: 'pointer', fontWeight: 500 },
  success: { color: '#16a34a', marginTop: 8 },
  error: { color: '#ef4444', marginTop: 8 },
};