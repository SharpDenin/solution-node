import { useEffect, useState } from 'react';
import { useParams } from 'react-router-dom';
import { api } from '../api/client';
import type { ReportDetail as ReportDetailType } from '../types';

export const ReportDetail = () => {
  const { id } = useParams();
  const [report, setReport] = useState<ReportDetailType | null>(null);

  useEffect(() => {
    if (id) {
      api.get(`/api/reports/${id}`)
        .then(res => {
          const data = res.data;
          // если metadata пришёл строкой — парсим, если объект — оставляем
          if (typeof data.metadata === 'string') {
            try {
              data.metadata = JSON.parse(data.metadata);
            } catch (e) {
              data.metadata = {};
            }
          }
          setReport(data);
        })
        .catch(console.error);
    }
  }, [id]);

  if (!report) return <div style={styles.loading}>Загрузка...</div>;

  // явное приведение к any, чтобы TypeScript не ругался на поля
  const meta = (report.metadata || {}) as any;
  const place = meta.place || '';
  const sort = meta.sort || '';
  const prioritySort = meta.priority_sort || '';

  const hasSort = !!sort;
  const isPriorityHigh = prioritySort === 'high';
  const titleColor = hasSort ? (isPriorityHigh ? '#ef4444' : '#16a34a') : 'inherit';

  const headerTitle = hasSort
    ? `Сорт: ${sort} (приоритет: ${prioritySort === 'high' ? 'высокий' : 'низкий'})`
    : place || 'Без места';

  return (
    <div style={styles.wrapper}>
      <h2 style={{ ...styles.mainTitle, color: titleColor }}>
        {headerTitle}
      </h2>

      <div style={styles.infoBlock}>
        <div style={styles.infoRow}>
          <span style={styles.label}>Дата:</span>
          <span style={styles.value}>
            {new Date(report.report_date).toLocaleDateString('ru-RU')}
          </span>
        </div>
        <div style={styles.infoRow}>
          <span style={styles.label}>Ответственный:</span>
          <span style={styles.value}>{report.responsible_name}</span>
        </div>
        {hasSort && (
          <div style={styles.infoRow}>
            <span style={styles.label}>Приоритет сорта:</span>
            <span style={{
              ...styles.value,
              color: isPriorityHigh ? '#ef4444' : '#16a34a',
              fontWeight: 600,
            }}>
              {isPriorityHigh ? '⚠️ Высокий' : '✅ Низкий'}
            </span>
          </div>
        )}
      </div>

      <div style={styles.answersContainer}>
        {report.answers.map((a, idx) => {
          let borderColor = 'transparent';
          let bgColor = 'white';
          if (a.result === 'good') {
            borderColor = '#16a34a';
            bgColor = '#f0fdf4';
          } else if (a.result === 'bad') {
            borderColor = '#ef4444';
            bgColor = '#fef2f2';
          } else if (a.result === 'neutral') {
            borderColor = '#e5e7eb';
          }

          return (
            <div
              key={idx}
              style={{
                ...styles.answerCard,
                borderLeft: `4px solid ${borderColor}`,
                backgroundColor: bgColor,
              }}
            >
              <div style={styles.questionHeader}>
                <span style={styles.questionNumber}>{idx + 1}.</span>
                <span style={styles.questionText}>{a.question_text}</span>
              </div>
              <div style={styles.answerContent}>
                {a.answer_text || <span style={styles.noAnswer}>Нет ответа</span>}
              </div>
              {a.image_url && (
                <div style={styles.imageWrapper}>
                  <img src={a.image_url} style={styles.image} alt="фото" />
                </div>
              )}
              {a.result && a.result !== 'neutral' && (
                <div style={{
                  ...styles.resultBadge,
                  backgroundColor: a.result === 'good' ? '#16a34a' : '#ef4444',
                }}>
                  {a.result === 'good' ? '✓ Хорошо' : '✗ Плохо'}
                </div>
              )}
            </div>
          );
        })}
      </div>
    </div>
  );
};

const styles = {
  wrapper: {
    maxWidth: '800px',
    margin: '0 auto',
    padding: '16px',
  },
  loading: {
    textAlign: 'center' as const,
    padding: '40px',
    color: '#6b7280',
  },
  mainTitle: {
    fontSize: '24px',
    fontWeight: 600,
    marginBottom: '24px',
    textAlign: 'center' as const,
    padding: '12px',
    backgroundColor: 'white',
    borderRadius: '12px',
    boxShadow: '0 2px 4px rgba(0,0,0,0.05)',
  },
  infoBlock: {
    backgroundColor: 'white',
    padding: '16px 20px',
    borderRadius: '12px',
    marginBottom: '24px',
    boxShadow: '0 2px 4px rgba(0,0,0,0.05)',
  },
  infoRow: {
    display: 'flex',
    marginBottom: '8px',
    fontSize: '15px',
  },
  label: {
    fontWeight: 600,
    color: '#374151',
    minWidth: '160px',
  },
  value: {
    color: '#111827',
  },
  answersContainer: {
    display: 'flex',
    flexDirection: 'column' as const,
    gap: '16px',
  },
  answerCard: {
    backgroundColor: 'white',
    borderRadius: '12px',
    padding: '16px 20px',
    boxShadow: '0 2px 4px rgba(0,0,0,0.05)',
    position: 'relative' as const,
  },
  questionHeader: {
    display: 'flex',
    alignItems: 'baseline',
    marginBottom: '12px',
    borderBottom: '1px solid #e5e7eb',
    paddingBottom: '8px',
  },
  questionNumber: {
    fontWeight: 700,
    color: '#16a34a',
    marginRight: '8px',
    fontSize: '16px',
  },
  questionText: {
    fontWeight: 600,
    fontSize: '16px',
    color: '#111827',
  },
  answerContent: {
    fontSize: '15px',
    color: '#4b5563',
    padding: '8px 0',
    textAlign: 'center' as const,
    fontStyle: 'italic',
  },
  noAnswer: {
    color: '#9ca3af',
  },
  imageWrapper: {
    textAlign: 'center' as const,
    marginTop: '12px',
  },
  image: {
    maxWidth: '100%',
    maxHeight: '400px',
    borderRadius: '8px',
    objectFit: 'contain' as const,
    border: '1px solid #e5e7eb',
  },
  resultBadge: {
    position: 'absolute' as const,
    top: '12px',
    right: '12px',
    color: 'white',
    fontSize: '12px',
    fontWeight: 600,
    padding: '4px 10px',
    borderRadius: '20px',
    letterSpacing: '0.5px',
  },
};