import { useEffect, useState } from 'react';
import { useParams } from 'react-router-dom';
import { api } from '../api/client';
import type { ReportDetail as ReportDetailType } from '../types';

export const ReportDetail = () => {
  const { id } = useParams();
  const [report, setReport] = useState<ReportDetailType | null>(null);

  useEffect(() => {
    if (id) {
      api.get(`/reports/${id}`).then(res => setReport(res.data));
    }
  }, [id]);

  if (!report) return <div>Загрузка...</div>;

  return (
    <div>
      <h2>Детали отчета</h2>
      <div style={styles.info}>
        <p><strong>Место:</strong> {report.place}</p>
        <p><strong>Дата:</strong> {new Date(report.report_date).toLocaleDateString('ru-RU')}</p>
        <p><strong>Ответственный:</strong> {report.responsible_name}</p>
      </div>
      {report.answers.map((a, idx) => (
        <div key={idx} style={styles.answerCard}>
          <b>{a.question_text}</b>
          <p>{a.answer_text}</p>
          {a.image_url && <img src={a.image_url} style={styles.image} />}
        </div>
      ))}
    </div>
  );
};

const styles = {
  info: { background: 'white', padding: 16, borderRadius: 12, marginBottom: 24 },
  answerCard: { background: 'white', padding: 16, borderRadius: 12, marginBottom: 16 },
  image: { maxWidth: 300, marginTop: 8 },
};