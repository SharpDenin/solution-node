import { useState, useEffect } from 'react';
import { api } from '../api/client';
import type { UserAdmin } from '../types';

// ---- Модалка редактирования (без изменений) ----
const EditUserModal = ({
  user,
  onClose,
  onSaved,
}: {
  user: UserAdmin;
  onClose: () => void;
  onSaved: () => void;
}) => {
  const [fullName, setFullName] = useState(user.full_name);
  const [login, setLogin] = useState(user.login);
  const [role, setRole] = useState(user.role);
  const [position, setPosition] = useState(user.position || '');
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');

  const handleSave = async () => {
    if (!fullName.trim() || !login.trim() || !role) return;
    setLoading(true);
    setError('');
    try {
      await api.put(`/api/admin/users/${user.id}`, {
        full_name: fullName.trim(),
        login: login.trim(),
        role,
        position: position.trim() || null,
      });
      onSaved();
      onClose();
    } catch (err: any) {
      setError(err.response?.data || 'Ошибка сохранения');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div style={modalStyles.backdrop} onClick={onClose}>
      <div style={modalStyles.modal} onClick={e => e.stopPropagation()}>
        <h3 style={{ marginTop: 0 }}>Редактировать пользователя</h3>
        <input type="text" placeholder="ФИО" value={fullName} onChange={e => setFullName(e.target.value)} style={modalStyles.input} />
        <input type="text" placeholder="Логин" value={login} onChange={e => setLogin(e.target.value)} style={modalStyles.input} />
        <select value={role} onChange={e => setRole(e.target.value)} style={modalStyles.input}>
          <option value="admin">Администратор</option>
          <option value="node">Специалист по растворным узлам</option>
          <option value="phenophase">Специалист по фенофазам</option>
        </select>
        <input type="text" placeholder="Должность" value={position} onChange={e => setPosition(e.target.value)} style={modalStyles.input} />
        {error && <div style={{ color: 'red', marginBottom: 8 }}>{error}</div>}
        <div style={{ display: 'flex', gap: 8, justifyContent: 'flex-end' }}>
          <button onClick={onClose} style={modalStyles.cancelBtn}>Отмена</button>
          <button onClick={handleSave} disabled={loading} style={modalStyles.saveBtn}>
            {loading ? 'Сохранение...' : 'Сохранить'}
          </button>
        </div>
      </div>
    </div>
  );
};

const modalStyles = {
  backdrop: { position: 'fixed' as const, top: 0, left: 0, right: 0, bottom: 0, backgroundColor: 'rgba(0,0,0,0.5)', display: 'flex', alignItems: 'center', justifyContent: 'center', zIndex: 1000 },
  modal: { backgroundColor: 'white', borderRadius: 12, padding: 24, width: 400, maxWidth: '90%' },
  input: { width: '100%', padding: 8, marginBottom: 12, borderRadius: 6, border: '1px solid #ccc', boxSizing: 'border-box' as const },
  cancelBtn: { background: '#f3f4f6', border: '1px solid #d1d5db', padding: '6px 12px', borderRadius: 6, cursor: 'pointer' },
  saveBtn: { background: '#16a34a', color: 'white', border: 'none', padding: '6px 16px', borderRadius: 6, cursor: 'pointer' },
};

// ---- Основной компонент ----
export const UserCreate = () => {
  const [fullName, setFullName] = useState('');
  const [login, setLogin] = useState('');
  const [password, setPassword] = useState('');
  const [role, setRole] = useState('node');
  const [position, setPosition] = useState('');
  const [message, setMessage] = useState('');
  const [error, setError] = useState('');

  const [showUsers, setShowUsers] = useState(false);
  const [users, setUsers] = useState<UserAdmin[]>([]);
  const [usersLoading, setUsersLoading] = useState(false);
  const [usersError, setUsersError] = useState('');

  const [editingUser, setEditingUser] = useState<UserAdmin | null>(null);

  const fetchUsers = async () => {
    setUsersLoading(true);
    setUsersError('');
    try {
      const res = await api.get('/api/admin/users');
      const data = Array.isArray(res.data) ? res.data : (res.data?.users ? res.data.users : []);
      setUsers(data);
    } catch (err: any) {
      setUsersError('Ошибка загрузки пользователей');
      setUsers([]);
    } finally {
      setUsersLoading(false);
    }
  };

  useEffect(() => {
    if (showUsers) {
      fetchUsers();
    }
  }, [showUsers]);

  const handleBlock = async (id: string) => {
    try {
      await api.patch(`/api/admin/users/${id}/block`);
      fetchUsers();
    } catch {
      alert('Не удалось заблокировать');
    }
  };

  const handleUnblock = async (id: string) => {
    try {
      await api.patch(`/api/admin/users/${id}/unblock`);
      fetchUsers();
    } catch {
      alert('Не удалось разблокировать');
    }
  };

  const handleDelete = async (id: string) => {
    if (!window.confirm('Удалить пользователя?')) return;
    try {
      await api.delete(`/api/admin/users/${id}`);
      fetchUsers();
    } catch {
      alert('Не удалось удалить');
    }
  };

  const handleRestore = async (id: string) => {
    try {
      await api.patch(`/api/admin/users/${id}/restore`);
      fetchUsers();
    } catch {
      alert('Не удалось восстановить');
    }
  };

  const handleCreate = async (e: React.FormEvent) => {
    e.preventDefault();
    setMessage('');
    setError('');
    try {
      await api.post('/api/register', {
        full_name: fullName,
        login,
        password,
        role,
        position,
      });
      setMessage('Пользователь успешно создан');
      setFullName('');
      setLogin('');
      setPassword('');
      setPosition('');
      if (showUsers) fetchUsers();
    } catch (err: any) {
      setError(err.response?.data || 'Ошибка при создании');
    }
  };

  return (
    <div style={styles.wrapper}>
      <h2 style={styles.title}>Управление пользователями</h2>

      <form onSubmit={handleCreate} style={styles.form}>
        <input type="text" placeholder="ФИО" value={fullName} onChange={e => setFullName(e.target.value)} style={styles.input} required />
        <input type="text" placeholder="Логин" value={login} onChange={e => setLogin(e.target.value)} style={styles.input} required />
        <input type="password" placeholder="Пароль" value={password} onChange={e => setPassword(e.target.value)} style={styles.input} required />
        <select value={role} onChange={e => setRole(e.target.value)} style={styles.input}>
          <option value="admin">Администратор</option>
          <option value="node">Специалист по растворным узлам</option>
          <option value="phenophase">Специалист по фенофазам</option>
        </select>
        <input type="text" placeholder="Должность (например: агроном)" value={position} onChange={e => setPosition(e.target.value)} style={styles.input} />
        <button type="submit" style={styles.button}>Создать</button>
        {message && <div style={styles.success}>{message}</div>}
        {error && <div style={styles.error}>{error}</div>}
      </form>

      <div style={{ width: '100%', maxWidth: 800, marginTop: 32 }}>
        <button
          onClick={() => setShowUsers(!showUsers)}
          style={styles.accordionHeader}
        >
          <span>👥 Список пользователей</span>
          <span style={{ transform: showUsers ? 'rotate(90deg)' : 'none', transition: '0.2s' }}>▶</span>
        </button>

        {showUsers && (
          <div style={styles.usersPanel}>
            {usersLoading && <p>Загрузка...</p>}
            {usersError && <div style={{ color: 'red' }}>{usersError}</div>}

            {!usersLoading && users.length === 0 && !usersError && (
              <p>Нет пользователей (кроме системного администратора).</p>
            )}

            {users.length > 0 && (
              <div style={{ overflowX: 'auto' }}>
                <table style={styles.table}>
                  <thead>
                    <tr>
                      <th style={styles.th}>ФИО</th>
                      <th style={styles.th}>Логин</th>
                      <th style={styles.th}>Роль</th>
                      <th style={styles.th}>Должность</th>
                      <th style={styles.th}>Статус</th>
                      <th style={styles.th}>Действия</th>
                    </tr>
                  </thead>
                  <tbody>
                    {users.map(u => {
                      let status = 'Активен';
                      if (u.is_deleted) status = 'Удалён';
                      else if (u.is_blocked) status = 'Заблокирован';
                      return (
                        <tr key={u.id}>
                          <td style={styles.td}>{u.full_name}</td>
                          <td style={styles.td}>{u.login}</td>
                          <td style={styles.td}>{u.role === 'admin' ? 'Администратор' : u.role === 'node' ? 'Узлы' : 'Фенофазы'}</td>
                          <td style={styles.td}>{u.position || '—'}</td>
                          <td style={styles.td}>
                            <span style={{
                              color: u.is_deleted ? '#9ca3af' : u.is_blocked ? '#ef4444' : '#16a34a',
                              fontWeight: 500,
                            }}>
                              {status}
                            </span>
                          </td>
                          <td style={styles.td}>
                            <div style={{ display: 'flex', gap: 4, flexWrap: 'wrap' }}>
                              <button onClick={() => setEditingUser(u)} style={styles.actionBtn} title="Редактировать">✏️</button>
                              {!u.is_deleted && !u.is_blocked && (
                                <button onClick={() => handleBlock(u.id)} style={{...styles.actionBtn, background: '#f59e0b'}} title="Заблокировать">⛔</button>
                              )}
                              {!u.is_deleted && u.is_blocked && (
                                <button onClick={() => handleUnblock(u.id)} style={{...styles.actionBtn, background: '#10b981'}} title="Разблокировать">✅</button>
                              )}
                              {!u.is_deleted && (
                                <button onClick={() => handleDelete(u.id)} style={{...styles.actionBtn, background: '#ef4444', color: 'white'}} title="Удалить">🗑️</button>
                              )}
                              {u.is_deleted && (
                                <button onClick={() => handleRestore(u.id)} style={{...styles.actionBtn, background: '#3b82f6', color: 'white'}} title="Восстановить">🔄</button>
                              )}
                            </div>
                          </td>
                        </tr>
                      );
                    })}
                  </tbody>
                </table>
              </div>
            )}
          </div>
        )}
      </div>

      {editingUser && (
        <EditUserModal
          user={editingUser}
          onClose={() => setEditingUser(null)}
          onSaved={() => { fetchUsers(); setEditingUser(null); }}
        />
      )}
    </div>
  );
};

const styles = {
  wrapper: {
    display: 'flex',
    flexDirection: 'column' as const,
    alignItems: 'center',
    padding: '40px 16px',
    minHeight: '100vh',
  },
  title: {
    marginBottom: '24px',
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
  success: { color: '#16a34a', fontSize: '14px', textAlign: 'center' as const },
  error: { color: '#ef4444', fontSize: '14px', textAlign: 'center' as const },

  accordionHeader: {
    width: '100%',
    padding: '12px 16px',
    background: '#f3f4f6',
    border: '1px solid #e5e7eb',
    borderRadius: 8,
    display: 'flex',
    justifyContent: 'space-between',
    alignItems: 'center',
    cursor: 'pointer',
    fontWeight: 600,
    fontSize: 15,
    color: '#111827',
    transition: 'background 0.2s',
  },
  usersPanel: {
    marginTop: 8,
    background: 'white',
    borderRadius: 8,
    padding: 16,
    boxShadow: '0 1px 4px rgba(0,0,0,0.05)',
  },
  table: { width: '100%', borderCollapse: 'collapse' as const },
  th: { textAlign: 'left' as const, padding: '8px 6px', borderBottom: '2px solid #e5e7eb', fontWeight: 600, fontSize: 13 },
  td: { padding: '8px 6px', borderBottom: '1px solid #e5e7eb', fontSize: 13 },
  actionBtn: { border: 'none', background: '#e5e7eb', borderRadius: 4, padding: '4px 8px', cursor: 'pointer', fontSize: 14 },
};