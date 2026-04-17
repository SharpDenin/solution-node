import { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { api } from '../api/client';
import type { Report, ReportFilters } from '../types';
import { ExportModal } from '../components/ExportModal';

// Возвращает начало и конец текущего года в формате YYYY-MM-DD
const getCurrentYearRange = () => {
  const now = new Date();
  const start = new Date(now.getFullYear(), 0, 1).toISOString().split('T')[0];
  const end = new Date(now.getFullYear(), 11, 31).toISOString().split('T')[0];
  return { start, end };
};

export const Dashboard = () => {
  const [reports, setReports] = useState<Report[]>([]);
  const [filters, setFilters] = useState<ReportFilters>(() => {
    const { start, end } = getCurrentYearRange();
    return { date_from: start, date_to: end, place: '', user_name: '' };
  });
  const [loading, setLoading] = useState(false);
  const [showExportModal, setShowExportModal] = useState(false);
  const navigate = useNavigate();

  const fetchReports = async () => {
    setLoading(true);
    try {
      const params: any = {};
      if (filters.date_from) params.date_from = filters.date_from;
      if (filters.date_to) params.date_to = filters.date_to;
      if (filters.place) params.place = filters.place;
      if (filters.user_name) params.user_name = filters.user_name;

      const res = await api.get('/reports', { params });
      setReports(Array.isArray(res.data) ? res.data : []);
    } catch (err) {
      console.error('Ошибка загрузки отчётов', err);
      setReports([]);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchReports();
  }, [filters]);

  const handleFilterChange = (key: keyof ReportFilters, value: string) => {
    setFilters(prev => ({ ...prev, [key]: value || undefined }));
  };

  const resetFilters = () => {
    const { start, end } = getCurrentYearRange();
    setFilters({ date_from: start, date_to: end, place: '', user_name: '' });
  };

  return (
    <div>
      <div style={styles.header}>
        <h2>Отчёты</h2>
        <button onClick={() => setShowExportModal(true)} style={styles.exportBtn}>
          Экспорт в Excel
        </button>
      </div>

      <div style={styles.filters}>
        <label>Дата с:</label>
        <input
          type="date"
          value={filters.date_from || ''}
          onChange={(e) => handleFilterChange('date_from', e.target.value)}
          style={styles.filterInput}
        />
        <label>Дата по:</label>
        <input
          type="date"
          value={filters.date_to || ''}
          onChange={(e) => handleFilterChange('date_to', e.target.value)}
          style={styles.filterInput}
        />
        <label>Место:</label>
        <input
          type="text"
          placeholder="Название места"
          value={filters.place || ''}
          onChange={(e) => handleFilterChange('place', e.target.value)}
          style={styles.filterInput}
        />
        <label>Ответственный:</label>
        <input
          type="text"
          placeholder="ФИО ответственного"
          value={filters.user_name || ''}
          onChange={(e) => handleFilterChange('user_name', e.target.value)}
          style={styles.filterInput}
        />
        <button onClick={resetFilters} style={styles.resetBtn}>Сбросить</button>
      </div>

      {loading && <p>Загрузка...</p>}

      <div style={styles.list}>
        {reports.map(report => (
          <div
            key={report.id}
            style={styles.card}
            onClick={() => navigate(`/reports/${report.id}`)}
          >
            <div style={styles.cardHeader}>
              <strong>{report.place}</strong>
              <span>{new Date(report.report_date).toLocaleDateString('ru-RU')}</span>
            </div>
            <div>Ответственный: {report.responsible_name}</div>
            <div style={styles.cardFooter}>
              Создан: {new Date(report.created_at).toLocaleDateString('ru-RU')}
            </div>
          </div>
        ))}
        {!loading && reports.length === 0 && (
          <p style={styles.noData}>Отчёты не найдены. Попробуйте изменить фильтры.</p>
        )}
      </div>

      {showExportModal && (
        <ExportModal
          onClose={() => setShowExportModal(false)}
          initialFilters={filters}
        />
      )}
    </div>
  );
};

const styles = {
  header: {
    display: 'flex',
    justifyContent: 'space-between',
    alignItems: 'center',
    marginBottom: '24px',
  },
  exportBtn: {
    background: '#2563eb',
    color: 'white',
    border: 'none',
    padding: '8px 16px',
    borderRadius: '8px',
    cursor: 'pointer',
  },
  filters: {
    display: 'flex',
    gap: '12px',
    marginBottom: '24px',
    flexWrap: 'wrap' as const,
    alignItems: 'center',
  },
  filterInput: {
    padding: '8px 12px',
    border: '1px solid #d1d5db',
    borderRadius: '8px',
    fontSize: '14px',
  },
  resetBtn: {
    background: '#6b7280',
    color: 'white',
    border: 'none',
    padding: '8px 16px',
    borderRadius: '8px',
    cursor: 'pointer',
  },
  list: {
    display: 'flex',
    flexDirection: 'column' as const,
    gap: '12px',
  },
  card: {
    background: 'white',
    padding: '16px',
    borderRadius: '12px',
    cursor: 'pointer',
    boxShadow: '0 1px 3px rgba(0,0,0,0.1)',
    transition: 'box-shadow 0.2s',
  },
  cardHeader: {
    display: 'flex',
    justifyContent: 'space-between',
    marginBottom: '8px',
  },
  cardFooter: {
    marginTop: '8px',
    fontSize: '12px',
    color: '#6b7280',
  },
  noData: {
    textAlign: 'center' as const,
    padding: '20px',
    color: '#6b7280',
  },
};