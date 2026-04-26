import { useState } from 'react';
import { api } from '../api/client';
import type { ReportFilters, Checklist } from '../types';

interface Props {
  onClose: () => void;
  initialFilters: ReportFilters;
  checklists: Checklist[];
}

export const ExportModal: React.FC<Props> = ({ onClose, initialFilters, checklists }) => {
  const [filters, setFilters] = useState(initialFilters);
  const [loading, setLoading] = useState(false);

  const selectedChecklist = checklists.find(c => c.id === filters.checklist_id);

  const handleExport = async () => {
    setLoading(true);
    try {
      const params: any = {};
      if (filters.date_from) params.date_from = filters.date_from;
      if (filters.date_to) params.date_to = filters.date_to;
      if (filters.user_name) params.user_name = filters.user_name;
      if (filters.checklist_id) params.checklist_id = filters.checklist_id;

      // Метаданные-фильтры
      if (selectedChecklist?.code === 'default' && filters.place) {
        params['metadata_place'] = filters.place;
      } else if (selectedChecklist?.code === 'sort_control') {
        if (filters.sort) params['metadata_sort'] = filters.sort;
        if (filters.priority_sort) params['metadata_priority_sort'] = filters.priority_sort;
      }

      const response = await api.get('/api/reports/export', {
        params,
        responseType: 'blob',
      });

      const url = window.URL.createObjectURL(new Blob([response.data]));
      const link = document.createElement('a');
      link.href = url;
      link.setAttribute('download', `export_${new Date().toISOString().slice(0, 10)}.xlsx`);
      document.body.appendChild(link);
      link.click();
      link.remove();
      window.URL.revokeObjectURL(url);
      onClose();
    } catch (err: any) {
      let message = 'Ошибка при выгрузке отчёта';

      // Если сервер ответил 404 – значит, нет данных
      if (err.response?.status === 404) {
        message = 'Нет отчётов по заданным фильтрам';
      } else if (err.response?.data) {
        // Пробуем извлечь текстовое сообщение из ответа
        try {
          if (typeof err.response.data === 'string') {
            message = err.response.data;
          } else if (err.response.data instanceof Blob) {
            const text = await err.response.data.text();
            if (text) message = text;
          }
        } catch (e) {
          // останется дефолтное сообщение
        }
      }

      alert(message);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div style={styles.backdrop} onClick={onClose}>
      <div style={styles.modal} onClick={e => e.stopPropagation()}>
        <h3 style={styles.title}>Выгрузка отчётов</h3>

        <div style={styles.field}>
          <label style={styles.label}>Тип чек-листа</label>
          <select
            value={filters.checklist_id || ''}
            onChange={e =>
              setFilters(prev => ({
                ...prev,
                checklist_id: e.target.value || undefined,
                place: '',
                sort: '',
                priority_sort: ''
              }))
            }
            style={styles.input}
          >
            <option value="">Все</option>
            {checklists.map(c => (
              <option key={c.id} value={c.id}>{c.name}</option>
            ))}
          </select>
        </div>

        <div style={styles.field}>
          <label style={styles.label}>Дата от</label>
          <input type="date" value={filters.date_from || ''} onChange={e => setFilters(prev => ({ ...prev, date_from: e.target.value }))} style={styles.input} />
        </div>
        <div style={styles.field}>
          <label style={styles.label}>Дата до</label>
          <input type="date" value={filters.date_to || ''} onChange={e => setFilters(prev => ({ ...prev, date_to: e.target.value }))} style={styles.input} />
        </div>
        <div style={styles.field}>
          <label style={styles.label}>Ответственный (ФИО)</label>
          <input type="text" placeholder="Введите часть ФИО" value={filters.user_name || ''} onChange={e => setFilters(prev => ({ ...prev, user_name: e.target.value }))} style={styles.input} />
        </div>

        {selectedChecklist?.code === 'default' && (
          <div style={styles.field}>
            <label style={styles.label}>Место работ</label>
            <input type="text" placeholder="Например: Растворный узел" value={filters.place || ''} onChange={e => setFilters(prev => ({ ...prev, place: e.target.value }))} style={styles.input} />
          </div>
        )}
        {selectedChecklist?.code === 'sort_control' && (
          <>
            <div style={styles.field}>
              <label style={styles.label}>Сорт</label>
              <input type="text" placeholder="Сорт" value={filters.sort || ''} onChange={e => setFilters(prev => ({ ...prev, sort: e.target.value }))} style={styles.input} />
            </div>
            <div style={styles.field}>
              <label style={styles.label}>Приоритет</label>
              <select value={filters.priority_sort || ''} onChange={e => setFilters(prev => ({ ...prev, priority_sort: e.target.value }))} style={styles.input}>
                <option value="">Все</option>
                <option value="high">Высокий</option>
                <option value="low">Низкий</option>
              </select>
            </div>
          </>
        )}

        <div style={styles.actions}>
          <button onClick={onClose} style={styles.cancelBtn}>Отмена</button>
          <button onClick={handleExport} disabled={loading} style={styles.exportBtn}>
            {loading ? 'Подготовка файла...' : 'Скачать Excel'}
          </button>
        </div>
      </div>
    </div>
  );
};

const styles = {
  backdrop: { position: 'fixed' as const, top: 0, left: 0, right: 0, bottom: 0, backgroundColor: 'rgba(0,0,0,0.5)', display: 'flex', alignItems: 'center', justifyContent: 'center', zIndex: 1000 },
  modal: { backgroundColor: '#fff', borderRadius: 16, padding: 24, width: '90%', maxWidth: 480, boxShadow: '0 20px 25px -5px rgba(0,0,0,0.1)' },
  title: { margin: '0 0 20px 0', fontSize: '1.5rem', fontWeight: 600, color: '#111827' },
  field: { marginBottom: 16 },
  label: { display: 'block', marginBottom: 6, fontSize: '0.875rem', fontWeight: 500, color: '#374151' },
  input: { width: '100%', padding: '10px 12px', borderRadius: 8, border: '1px solid #d1d5db', fontSize: '0.875rem', outline: 'none', boxSizing: 'border-box' as const },
  actions: { display: 'flex', justifyContent: 'flex-end', gap: 12, marginTop: 24 },
  cancelBtn: { padding: '8px 16px', backgroundColor: '#f3f4f6', border: '1px solid #e5e7eb', borderRadius: 8, fontSize: '0.875rem', fontWeight: 500, cursor: 'pointer', color: '#374151' },
  exportBtn: { padding: '8px 20px', backgroundColor: '#2563eb', border: 'none', borderRadius: 8, fontSize: '0.875rem', fontWeight: 500, cursor: 'pointer', color: 'white' },
};