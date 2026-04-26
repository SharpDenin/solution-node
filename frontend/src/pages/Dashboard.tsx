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

export const Dashboard = () => {
  const [reports, setReports] = useState<Report[]>([]);
  const [checklists, setChecklists] = useState<Checklist[]>([]);
  const [filters, setFilters] = useState<ReportFilters>(() => {
    const { start, end } = getCurrentYearRange();
    return { date_from: start, date_to: end, user_name: '' };
  });
  const [selectedChecklistId, setSelectedChecklistId] = useState<string>('');
  const [loading, setLoading] = useState(false);
  const [showExportModal, setShowExportModal] = useState(false);
  const navigate = useNavigate();

  useEffect(() => {
    api.get('/api/checklists').then(res => setChecklists(res.data));
  }, []);

  const fetchReports = async () => {
    setLoading(true);
    try {
      const params: any = {};
      if (filters.date_from) params.date_from = filters.date_from;
      if (filters.date_to) params.date_to = filters.date_to;
      if (filters.user_name) params.user_name = filters.user_name;
      if (selectedChecklistId) params.checklist_id = selectedChecklistId;

      if (selectedChecklistId) {
        const cl = checklists.find(c => c.id === selectedChecklistId);
        if (cl?.code === 'default' && filters.place) {
          params['metadata_place'] = filters.place;
        } else if (cl?.code === 'sort_control') {
          if (filters.sort) params['metadata_sort'] = filters.sort;
          if (filters.priority_sort) params['metadata_priority_sort'] = filters.priority_sort;
        }
      }

      const res = await api.get('/api/reports', { params });
      const raw = Array.isArray(res.data) ? res.data : [];
      const parsed = raw.map((r: any) => ({
        ...r,
        place: r.metadata?.place || '',
        sort: r.metadata?.sort || '',
        priority_sort: r.metadata?.priority_sort || '',
      }));
      setReports(parsed);
    } catch (err) {
      console.error('Ошибка загрузки отчётов', err);
      setReports([]);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchReports();
  }, [filters, selectedChecklistId, checklists]);

  const handleFilterChange = (key: keyof ReportFilters, value: string) => {
    setFilters(prev => ({ ...prev, [key]: value || undefined }));
  };

  const handleChecklistChange = (id: string) => {
    setSelectedChecklistId(id);
    setFilters(prev => ({ ...prev, place: '', sort: '', priority_sort: '' }));
  };

  const resetFilters = () => {
    const { start, end } = getCurrentYearRange();
    setFilters({ date_from: start, date_to: end, user_name: '' });
    setSelectedChecklistId('');
  };

  const selectedChecklist = checklists.find(c => c.id === selectedChecklistId);

  return (
    <div>
      <div style={styles.header}>
        <h2>Отчёты</h2>
        <button onClick={() => setShowExportModal(true)} style={styles.exportBtn}>
          Экспорт в Excel
        </button>
      </div>

      <div style={styles.filters}>
        <label>Тип чек-листа:</label>
        <select
          value={selectedChecklistId}
          onChange={e => handleChecklistChange(e.target.value)}
          style={styles.filterInput}
        >
          <option value="">Все</option>
          {checklists.map(c => (
            <option key={c.id} value={c.id}>{c.name}</option>
          ))}
        </select>

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

        {selectedChecklist && selectedChecklist.code === 'default' && (
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
        {selectedChecklist && selectedChecklist.code === 'sort_control' && (
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
              style={styles.card}
              onClick={() => navigate(`/reports/${report.id}`)}
            >
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
        {!loading && reports.length === 0 && (
          <p style={styles.noData}>Отчёты не найдены.</p>
        )}
      </div>

      {showExportModal && (
        <ExportModal
          onClose={() => setShowExportModal(false)}
          initialFilters={{ ...filters, checklist_id: selectedChecklistId }}
          checklists={checklists}
        />
      )}
    </div>
  );
};

const styles = {
  header: { display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 24 },
  exportBtn: { background: '#2563eb', color: 'white', border: 'none', padding: '8px 16px', borderRadius: 8, cursor: 'pointer' },
  filters: { display: 'flex', gap: 12, marginBottom: 24, flexWrap: 'wrap' as const, alignItems: 'center' },
  filterInput: { padding: '8px 12px', border: '1px solid #d1d5db', borderRadius: 8, fontSize: 14 },
  resetBtn: { background: '#6b7280', color: 'white', border: 'none', padding: '8px 16px', borderRadius: 8, cursor: 'pointer' },
  list: { display: 'flex', flexDirection: 'column' as const, gap: 12 },
  card: { background: 'white', padding: 16, borderRadius: 12, cursor: 'pointer', boxShadow: '0 1px 3px rgba(0,0,0,0.1)' },
  cardHeader: { display: 'flex', justifyContent: 'space-between', marginBottom: 8 },
  cardFooter: { marginTop: 8, fontSize: 12, color: '#6b7280' },
  noData: { textAlign: 'center' as const, padding: 20, color: '#6b7280' },
};