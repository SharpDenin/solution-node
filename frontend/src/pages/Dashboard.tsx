import { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { api } from '../api/client';
import type { Report, ReportFilters, Checklist } from '../types';
import { ExportModal } from '../components/ExportModal';

const getCurrentYearRange = () => {
  const now = new Date();
  const start = new Date(now.getFullYear(), 0, 1).toISOString().split('T')[0];
  const end = new Date(now.getFullYear(), 11, 31).toISOString().split('T')[0];
  return { start, end };
};

const parseMetadata = (report: any) => {
  if (typeof report.metadata === 'string') {
    try {
      report.metadata = JSON.parse(report.metadata);
    } catch {
      report.metadata = {};
    }
  }
  return {
    ...report,
    place: report.place || report.metadata?.place || '',
    sort: report.metadata?.sort || '',
    priority_sort: report.metadata?.priority_sort || '',
  };
};

export const Dashboard = () => {
  const [reports, setReports] = useState<Report[]>([]);
  const [checklists, setChecklists] = useState<Checklist[]>([]);
  const [activeChecklistId, setActiveChecklistId] = useState<string>('');
  const [filters, setFilters] = useState<ReportFilters>(() => {
    const { start, end } = getCurrentYearRange();
    return { date_from: start, date_to: end, user_name: '' };
  });
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [showExportModal, setShowExportModal] = useState(false);
  const navigate = useNavigate();

  useEffect(() => {
    api.get('/api/checklists').then(res => {
      setChecklists(res.data);
      if (res.data.length > 0 && !activeChecklistId) {
        setActiveChecklistId(res.data[0].id);
      }
    });
  }, []);

  const fetchReports = async () => {
    if (!activeChecklistId) return;
    setLoading(true);
    setError(null);
    try {
      const params: any = {
        checklist_id: activeChecklistId,
        limit: 1000,
        offset: 0,
      };
      if (filters.date_from) params.date_from = filters.date_from;
      if (filters.date_to) params.date_to = filters.date_to;
      if (filters.user_name) params.user_name = filters.user_name;

      const cl = checklists.find(c => c.id === activeChecklistId);
      if (cl?.code === 'default' && filters.place) {
        params['metadata_place'] = filters.place;
      } else if (cl?.code === 'sort_control') {
        if (filters.sort) params['metadata_sort'] = filters.sort;
        if (filters.priority_sort) params['metadata_priority_sort'] = filters.priority_sort;
      }

      const res = await api.get('/api/reports', { params });
      const raw = Array.isArray(res.data) ? res.data : [];
      const parsed = raw.map(parseMetadata);
      setReports(parsed);
    } catch (err: any) {
      console.error('Ошибка загрузки отчётов', err);
      const message = err?.response?.data || 'Не удалось загрузить отчёты';
      setError(typeof message === 'string' ? message : 'Ошибка сервера');
      setReports([]);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchReports();
  }, [filters, activeChecklistId, checklists]);

  const handleFilterChange = (key: keyof ReportFilters, value: string) => {
    setFilters(prev => ({ ...prev, [key]: value || undefined }));
  };

  const switchTab = (checklistId: string) => {
    setActiveChecklistId(checklistId);
    setFilters(prev => ({
      ...prev,
      place: '',
      sort: '',
      priority_sort: '',
    }));
  };

  const resetFilters = () => {
    const { start, end } = getCurrentYearRange();
    setFilters({ date_from: start, date_to: end, user_name: '' });
  };

  // Удаление отчёта
  const deleteReport = async (id: string) => {
    if (!window.confirm('Вы уверены, что хотите удалить этот отчёт?')) return;
    try {
      await api.delete(`/api/reports/${id}`);
      fetchReports(); // обновить список после удаления
    } catch (err: any) {
      alert(err?.response?.data || 'Ошибка удаления отчёта');
    }
  };

  const activeChecklist = checklists.find(c => c.id === activeChecklistId);

  return (
    <div>
      <div style={styles.header}>
        <h2>Отчёты</h2>
        <button onClick={() => setShowExportModal(true)} style={styles.exportBtn}>
          Экспорт в Excel
        </button>
      </div>

      <div style={styles.tabs}>
        {checklists.map(cl => (
          <button
            key={cl.id}
            onClick={() => switchTab(cl.id)}
            style={{
              ...styles.tabButton,
              ...(activeChecklistId === cl.id ? styles.activeTab : {}),
            }}
          >
            {cl.name}
          </button>
        ))}
      </div>

      <div style={styles.filters}>
        <label>Дата с:</label>
        <input
          type="date"
          value={filters.date_from || ''}
          onChange={e => handleFilterChange('date_from', e.target.value)}
          style={styles.filterInput}
        />
        <label>Дата по:</label>
        <input
          type="date"
          value={filters.date_to || ''}
          onChange={e => handleFilterChange('date_to', e.target.value)}
          style={styles.filterInput}
        />
        <label>Ответственный:</label>
        <input
          type="text"
          placeholder="ФИО"
          value={filters.user_name || ''}
          onChange={e => handleFilterChange('user_name', e.target.value)}
          style={styles.filterInput}
        />

        {activeChecklist?.code === 'default' && (
          <>
            <label>Место:</label>
            <input
              type="text"
              placeholder="Название места"
              value={filters.place || ''}
              onChange={e => handleFilterChange('place', e.target.value)}
              style={styles.filterInput}
            />
          </>
        )}
        {activeChecklist?.code === 'sort_control' && (
          <>
            <label>Сорт:</label>
            <input
              type="text"
              value={filters.sort || ''}
              onChange={e => handleFilterChange('sort', e.target.value)}
              style={styles.filterInput}
            />
            <label>Приоритет:</label>
            <select
              value={filters.priority_sort || ''}
              onChange={e => handleFilterChange('priority_sort', e.target.value)}
              style={styles.filterInput}
            >
              <option value="">Все</option>
              <option value="high">Высокий</option>
              <option value="low">Низкий</option>
            </select>
          </>
        )}

        <button onClick={resetFilters} style={styles.resetBtn}>Сбросить</button>
      </div>

      {error && <div style={styles.error}>{error}</div>}
      {loading && <p>Загрузка...</p>}

      <div style={styles.list}>
        {reports.map(report => {
          const title = report.sort
            ? `Сорт: ${report.sort}`
            : report.place || 'Без названия';
          const priorityColor =
            report.priority_sort === 'high' ? '#ef4444' :
            report.priority_sort === 'low' ? '#16a34a' :
            'inherit';

          return (
            <div
              key={report.id}
              style={{ ...styles.card, position: 'relative' }}
              onClick={() => navigate(`/reports/${report.id}`)}
            >
              {/* Кнопка удаления */}
              <button
                onClick={(e) => {
                  e.stopPropagation(); // чтобы не переходить на страницу отчёта
                  deleteReport(report.id);
                }}
                style={styles.deleteBtn}
                title="Удалить отчёт"
              >
                🗑️
              </button>

              <div style={styles.cardHeader}>
                <strong style={{ color: priorityColor }}>{title}</strong>
                <span>{new Date(report.report_date).toLocaleDateString('ru-RU')}</span>
              </div>
              <div>Ответственный: {report.responsible_name}</div>
              <div style={styles.cardFooter}>
                Создан: {new Date(report.created_at).toLocaleDateString('ru-RU')}
              </div>
            </div>
          );
        })}
        {!loading && !error && reports.length === 0 && (
          <p style={styles.noData}>Отчёты не найдены.</p>
        )}
      </div>

      {showExportModal && (
        <ExportModal
          onClose={() => setShowExportModal(false)}
          initialFilters={{ ...filters, checklist_id: activeChecklistId }}
          checklists={checklists}
        />
      )}
    </div>
  );
};

const styles = {
  header: { display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 24 },
  exportBtn: { background: '#2563eb', color: 'white', border: 'none', padding: '8px 16px', borderRadius: 8, cursor: 'pointer' },
  tabs: {
    display: 'flex',
    gap: 0,
    marginBottom: 24,
    borderBottom: '2px solid #e5e7eb',
  },
  tabButton: {
    padding: '10px 20px',
    background: 'transparent',
    border: 'none',
    borderBottom: '2px solid transparent',
    marginBottom: '-2px',
    cursor: 'pointer',
    fontSize: '15px',
    fontWeight: 500,
    color: '#6b7280',
    transition: 'all 0.2s',
    outline: 'none',
  },
  activeTab: {
    color: '#16a34a',
    borderBottom: '2px solid #16a34a',
  },
  filters: { display: 'flex', gap: 12, marginBottom: 24, flexWrap: 'wrap' as const, alignItems: 'center' },
  filterInput: { padding: '8px 12px', border: '1px solid #d1d5db', borderRadius: 8, fontSize: 14 },
  resetBtn: { background: '#6b7280', color: 'white', border: 'none', padding: '8px 16px', borderRadius: 8, cursor: 'pointer' },
  list: { display: 'flex', flexDirection: 'column' as const, gap: 12 },
  card: {
    background: 'white',
    padding: 16,
    borderRadius: 12,
    cursor: 'pointer',
    boxShadow: '0 1px 3px rgba(0,0,0,0.1)',
    position: 'relative', // важно для позиционирования кнопки удаления
  },
  cardHeader: { display: 'flex', justifyContent: 'space-between', marginBottom: 8 },
  cardFooter: { marginTop: 8, fontSize: 12, color: '#6b7280' },
  noData: { textAlign: 'center' as const, padding: 20, color: '#6b7280' },
  error: { background: '#fee2e2', color: '#b91c1c', padding: '12px', borderRadius: 8, marginBottom: 16 },
  deleteBtn: {
    position: 'absolute' as const,
    bottom: 8,
    right: 8,
    background: 'transparent',
    border: 'none',
    fontSize: 18,
    cursor: 'pointer',
    color: '#ef4444',
    zIndex: 2,
    padding: 4,
    lineHeight: 1,
  },
};