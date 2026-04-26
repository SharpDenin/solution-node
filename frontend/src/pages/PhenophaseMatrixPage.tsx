import { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { api } from '../api/client';
import type {
  PhenophaseMatrixReportResponse,
} from '../types';

export const PhenophaseMatrixPage = () => {
  const navigate = useNavigate();
  const [varieties, setVarieties] = useState<{ id: string; name: string }[]>([]);
  const [selectedVarietyId, setSelectedVarietyId] = useState('');
  const [matrix, setMatrix] = useState<PhenophaseMatrixReportResponse | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');

  useEffect(() => {
    api
      .get('/api/varieties')
      .then(res => setVarieties(res.data))
      .catch(() => setError('Ошибка загрузки сортов'));
  }, []);

  const fetchMatrix = async (varietyId: string) => {
    if (!varietyId) return;
    setLoading(true);
    setError('');
    try {
      const res = await api.get<PhenophaseMatrixReportResponse>(
        `/api/reports/phenophase-matrix?variety_id=${varietyId}`
      );
      setMatrix(res.data);
    } catch (err: any) {
      const msg = err?.response?.data || 'Не удалось загрузить матрицу';
      setError(typeof msg === 'string' ? msg : 'Ошибка сервера');
      setMatrix(null);
    } finally {
      setLoading(false);
    }
  };

  const handleClose = () => {
    // возвращаемся на предыдущую страницу (дашборд или другой раздел)
    navigate(-1);
  };

  return (
    <div style={styles.backdrop} onClick={handleClose}>
      <div style={styles.modal} onClick={e => e.stopPropagation()}>
        {/* Крестик закрытия */}
        <button onClick={handleClose} style={styles.closeBtn}>
          ✕
        </button>

        <h2 style={{ marginTop: 0, marginBottom: 24, fontSize: 20 }}>
          🌸 Отчёт по фенофазам
        </h2>

        <div style={styles.selector}>
          <label style={styles.label}>Сорт:</label>
          <select
            value={selectedVarietyId}
            onChange={e => {
              setSelectedVarietyId(e.target.value);
              fetchMatrix(e.target.value);
            }}
            style={styles.select}
          >
            <option value="">-- Выберите сорт --</option>
            {varieties.map(v => (
              <option key={v.id} value={v.id}>
                {v.name}
              </option>
            ))}
          </select>
        </div>

        {loading && <p style={{ marginTop: 16 }}>Загрузка...</p>}

        {error && <div style={styles.errorBanner}>{error}</div>}

        {matrix && matrix.columns.length > 0 && matrix.rows.length > 0 && (
          <div style={{ overflowX: 'auto', marginTop: 16 }}>
            <table style={styles.table}>
              <thead>
                <tr>
                  <th style={styles.th}>Вопрос</th>
                  {matrix.columns.map(col => (
                    <th key={col.phenophase_id} style={styles.th}>
                      {col.name}
                    </th>
                  ))}
                </tr>
              </thead>
              <tbody>
                {matrix.rows.map(row => (
                  <tr key={row.question_id}>
                    <td style={styles.tdQuestion}>{row.text}</td>
                    {row.cells.map((cell, idx) => {
                      const hasAnswer = cell.answer_text !== null && cell.answer_text !== undefined;
                      const result = cell.result;
                      const bgColor =
                        result === 'good'
                          ? '#f0fdf4'
                          : result === 'bad'
                          ? '#fef2f2'
                          : 'transparent';
                      const borderColor =
                        result === 'good'
                          ? '#16a34a'
                          : result === 'bad'
                          ? '#ef4444'
                          : '#e5e7eb';

                      return (
                        <td
                          key={idx}
                          style={{
                            ...styles.td,
                            backgroundColor: bgColor,
                            borderLeft: `3px solid ${borderColor}`,
                          }}
                          title={
                            hasAnswer
                              ? `Ответ: ${cell.answer_text}\nРезультат: ${result ?? '—'}`
                              : 'Нет данных'
                          }
                        >
                          {hasAnswer ? cell.answer_text : '—'}
                        </td>
                      );
                    })}
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}

        {matrix && matrix.rows.length === 0 && !loading && (
          <p style={{ marginTop: 16, color: '#6b7280' }}>
            Нет данных. Создайте хотя бы один отчёт по чек-листу фенофаз для этого сорта.
          </p>
        )}
      </div>
    </div>
  );
};

const styles = {
  backdrop: {
    position: 'fixed' as const,
    top: 0,
    left: 0,
    width: '100vw',
    height: '100vh',
    backgroundColor: 'rgba(0,0,0,0.5)',
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'center',
    zIndex: 1000,
  },
  modal: {
    backgroundColor: 'white',
    borderRadius: 16,
    padding: 24,
    width: '95vw',
    maxHeight: '90vh',
    overflowY: 'auto' as const,
    position: 'relative' as const,
    boxShadow: '0 20px 25px -5px rgba(0,0,0,0.1)',
  },
  closeBtn: {
    position: 'absolute' as const,
    top: 16,
    right: 16,
    background: 'transparent',
    border: 'none',
    fontSize: 24,
    cursor: 'pointer',
    color: '#6b7280',
    lineHeight: 1,
    padding: 0,
  },
  selector: {
    display: 'flex',
    alignItems: 'center',
    gap: 12,
    flexWrap: 'wrap' as const,
    marginBottom: 8,
  },
  label: {
    fontWeight: 600,
    fontSize: 14,
    color: '#111827',
  },
  select: {
    padding: '8px 12px',
    borderRadius: 8,
    border: '1px solid #d1d5db',
    fontSize: 14,
    minWidth: 220,
    outline: 'none',
  },
  table: {
    borderCollapse: 'collapse' as const,
    width: '100%',
    minWidth: 600,
    background: 'white',
    borderRadius: 12,
    overflow: 'hidden',
    boxShadow: '0 2px 8px rgba(0,0,0,0.05)',
  },
  th: {
    padding: '10px 8px',
    textAlign: 'left' as const,
    background: '#f3f4f6',
    fontWeight: 600,
    fontSize: 12,
    borderBottom: '2px solid #e5e7eb',
  },
  td: {
    padding: '8px 8px',
    borderBottom: '1px solid #e5e7eb',
    fontSize: 12,
    transition: 'background 0.2s',
  },
  tdQuestion: {
    padding: '8px 8px',
    fontWeight: 500,
    fontSize: 12,
    borderBottom: '1px solid #e5e7eb',
    backgroundColor: '#f9fafb',
  },
  errorBanner: {
    marginTop: 16,
    backgroundColor: '#fee2e2',
    color: '#b91c1c',
    padding: '12px 16px',
    borderRadius: 8,
    fontSize: 13,
  },
};