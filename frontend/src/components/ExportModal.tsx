import { useState } from 'react';
import { api } from '../api/client';
import type { ReportFilters } from '../types';

interface Props {
  onClose: () => void;
  initialFilters: ReportFilters;
}

export const ExportModal: React.FC<Props> = ({ onClose, initialFilters }) => {
  const [filters, setFilters] = useState(initialFilters);
  const [loading, setLoading] = useState(false);

  const handleExport = async () => {
    setLoading(true);
    try {
      const params = new URLSearchParams();
      if (filters.date_from) params.append('date_from', filters.date_from);
      if (filters.date_to) params.append('date_to', filters.date_to);
      if (filters.place) params.append('place', filters.place);
      if (filters.user_name) params.append('user_name', filters.user_name); // новое поле
      
      const response = await api.get('/api/reports/export', {
        params,
        responseType: 'blob',
      });
      
      const url = window.URL.createObjectURL(new Blob([response.data]));
      const link = document.createElement('a');
      link.href = url;
      // Имя файла генерируется на бэке, но здесь можно задать своё
      link.setAttribute('download', `export_${new Date().toISOString().slice(0,10)}.xlsx`);
      document.body.appendChild(link);
      link.click();
      link.remove();
      window.URL.revokeObjectURL(url);
      onClose();
    } catch (err) {
      alert('Ошибка при выгрузке отчёта');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div style={styles.backdrop} onClick={onClose}>
      <div style={styles.modal} onClick={(e) => e.stopPropagation()}>
        <h3 style={styles.title}>Выгрузка отчётов</h3>
        
        <div style={styles.field}>
          <label style={styles.label}>Дата от</label>
          <input
            type="date"
            value={filters.date_from || ''}
            onChange={(e) => setFilters(prev => ({ ...prev, date_from: e.target.value }))}
            style={styles.input}
          />
        </div>

        <div style={styles.field}>
          <label style={styles.label}>Дата до</label>
          <input
            type="date"
            value={filters.date_to || ''}
            onChange={(e) => setFilters(prev => ({ ...prev, date_to: e.target.value }))}
            style={styles.input}
          />
        </div>

        <div style={styles.field}>
          <label style={styles.label}>Место работ</label>
          <input
            type="text"
            placeholder="Например: Растворный узел"
            value={filters.place || ''}
            onChange={(e) => setFilters(prev => ({ ...prev, place: e.target.value }))}
            style={styles.input}
          />
        </div>

        <div style={styles.field}>
          <label style={styles.label}>Ответственный (ФИО)</label>
          <input
            type="text"
            placeholder="Введите часть ФИО"
            value={filters.user_name || ''}
            onChange={(e) => setFilters(prev => ({ ...prev, user_name: e.target.value }))}
            style={styles.input}
          />
        </div>

        <div style={styles.actions}>
          <button onClick={onClose} style={styles.cancelBtn}>
            Отмена
          </button>
          <button onClick={handleExport} disabled={loading} style={styles.exportBtn}>
            {loading ? 'Подготовка файла...' : 'Скачать Excel'}
          </button>
        </div>
      </div>
    </div>
  );
};

const styles = {
  backdrop: {
    position: 'fixed' as const,
    top: 0,
    left: 0,
    right: 0,
    bottom: 0,
    backgroundColor: 'rgba(0, 0, 0, 0.5)',
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'center',
    zIndex: 1000,
  },
  modal: {
    backgroundColor: '#fff',
    borderRadius: 16,
    padding: '24px',
    width: '90%',
    maxWidth: '480px',
    boxShadow: '0 20px 25px -5px rgba(0,0,0,0.1), 0 10px 10px -5px rgba(0,0,0,0.04)',
  },
  title: {
    margin: '0 0 20px 0',
    fontSize: '1.5rem',
    fontWeight: 600,
    color: '#111827',
  },
  field: {
    marginBottom: 16,
  },
  label: {
    display: 'block',
    marginBottom: 6,
    fontSize: '0.875rem',
    fontWeight: 500,
    color: '#374151',
  },
  input: {
    width: '100%',
    padding: '10px 12px',
    borderRadius: 8,
    border: '1px solid #d1d5db',
    fontSize: '0.875rem',
    transition: 'border-color 0.2s',
    outline: 'none',
    boxSizing: 'border-box' as const,
  },
  actions: {
    display: 'flex',
    justifyContent: 'flex-end',
    gap: '12px',
    marginTop: '24px',
  },
  cancelBtn: {
    padding: '8px 16px',
    backgroundColor: '#f3f4f6',
    border: '1px solid #e5e7eb',
    borderRadius: 8,
    fontSize: '0.875rem',
    fontWeight: 500,
    cursor: 'pointer',
    color: '#374151',
    transition: 'background-color 0.2s',
  },
  exportBtn: {
    padding: '8px 20px',
    backgroundColor: '#2563eb',
    border: 'none',
    borderRadius: 8,
    fontSize: '0.875rem',
    fontWeight: 500,
    cursor: 'pointer',
    color: 'white',
    transition: 'background-color 0.2s',
    disabled: {
      opacity: 0.6,
      cursor: 'not-allowed',
    },
  },
};