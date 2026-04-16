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
      const response = await api.get('/reports/export', {
        params,
        responseType: 'blob',
      });
      const url = window.URL.createObjectURL(new Blob([response.data]));
      const link = document.createElement('a');
      link.href = url;
      link.setAttribute('download', 'reports.xlsx');
      document.body.appendChild(link);
      link.click();
      link.remove();
      onClose();
    } catch (err) {
      alert('Export failed');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div style={styles.backdrop} onClick={onClose}>
      <div style={styles.modal} onClick={(e) => e.stopPropagation()}>
        <h3>Export Filters</h3>
        <input
          type="date"
          placeholder="Date from"
          value={filters.date_from || ''}
          onChange={(e) => setFilters(prev => ({ ...prev, date_from: e.target.value }))}
          style={styles.input}
        />
        <input
          type="date"
          placeholder="Date to"
          value={filters.date_to || ''}
          onChange={(e) => setFilters(prev => ({ ...prev, date_to: e.target.value }))}
          style={styles.input}
        />
        <input
          type="text"
          placeholder="Place"
          value={filters.place || ''}
          onChange={(e) => setFilters(prev => ({ ...prev, place: e.target.value }))}
          style={styles.input}
        />
        <div style={styles.actions}>
          <button onClick={onClose} style={styles.cancelBtn}>Отмена</button>
          <button onClick={handleExport} disabled={loading} style={styles.exportBtn}>
            {loading ? 'Идет скачивание...' : 'Скачать'}
          </button>
        </div>
      </div>
    </div>
  );
};

const styles = {
  backdrop: {
    position: 'fixed' as const,
    inset: 0,
    background: 'rgba(0,0,0,0.5)',
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'center',
    zIndex: 1000,
  },
  modal: {
    background: 'white',
    padding: '24px',
    borderRadius: '16px',
    width: '360px',
  },
  input: {
    width: '100%',
    padding: '8px 12px',
    marginBottom: '12px',
    border: '1px solid #d1d5db',
    borderRadius: '8px',
  },
  actions: {
    display: 'flex',
    justifyContent: 'flex-end',
    gap: '12px',
    marginTop: '16px',
  },
  cancelBtn: {
    padding: '8px 16px',
    background: '#e5e7eb',
    border: 'none',
    borderRadius: '8px',
    cursor: 'pointer',
  },
  exportBtn: {
    padding: '8px 16px',
    background: '#2563eb',
    color: 'white',
    border: 'none',
    borderRadius: '8px',
    cursor: 'pointer',
  },
};