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
        .then(res => setReport(res.data))
        .catch(console.error);
    }
  }, [id]);

  if (!report) return <div>Загрузка...</div>;

  const isPriorityHigh = report.priority_sort === 'high';

  return (
    <div>
      <h2 style={{ color: isPriorityHigh ? '#ef4444' : 'inherit' }}>
        Детали отчёта
      </h2>
      <div style={styles.info}>
        {report.place && <p><strong>Место:</strong> {report.place}</p>}
        {report.sort && (
          <p><strong>Сорт:</strong> {report.sort} (приоритет: {report.priority_sort})</p>
        )}
        <p><strong>Дата:</strong> {new Date(report.report_date).toLocaleDateString('ru-RU')}</p>
        <p><strong>Ответственный:</strong> {report.responsible_name}</p>
      </div>
      {report.answers.map((a, idx) => (
        <div
          key={idx}
          style={{
            ...styles.answerCard,
            borderLeft: a.evaluation === 'good'
              ? '4px solid #16a34a'
              : a.evaluation === 'bad'
              ? '4px solid #ef4444'
              : 'none',
          }}
        >
          <b>{a.question_text}</b>
          <p>{a.answer_text}</p>
          {a.image_url && <img src={a.image_url} style={styles.image} alt="фото" />}
        </div>
      ))}
    </div>
  );
};

const styles = {
  info: { background: 'white', padding: 16, borderRadius: 12, marginBottom: 24 },
  answerCard: { background: 'white', padding: 16, borderRadius: 12, marginBottom: 16 },
  image: { maxWidth: '100%', maxHeight: 500, marginTop: 8, borderRadius: 8, objectFit: 'contain' as const },
};