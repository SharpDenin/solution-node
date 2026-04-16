import { useEffect, useState } from "react";
import { useParams } from "react-router-dom";
import { api } from "../api/client";

type Answer = {
  question_id: string;
  question_text: string;
  answer_text: string;
  image_url?: string;
};

type Report = {
  id: string;
  place: string;
  report_date: string;
  responsible_name: string;
  answers: Answer[];
};

export default function ReportDetail() {
  const { id } = useParams();
  const [report, setReport] = useState<Report | null>(null);

  useEffect(() => {
    if (!id) return;

    api.get(`/reports/${id}`).then(res => {
      setReport(res.data);
    });
  }, [id]);

  if (!report) return <div>Loading...</div>;

  return (
    <div>
      <h2>Report</h2>

      {/* 📄 Общая инфа */}
      <div style={{ marginBottom: 20 }}>
        <p><b>Place:</b> {report.place}</p>
        <p><b>Date:</b> {report.report_date}</p>
        <p><b>Responsible:</b> {report.responsible_name}</p>
      </div>

      <hr />

      {/* 📋 Ответы */}
      {report.answers.map((a, i) => (
        <div
          key={i}
          style={{
            marginBottom: 20,
            padding: 15,
            background: "white",
            borderRadius: 8,
          }}
        >
          <b>{a.question_text}</b>

          <p>{a.answer_text}</p>

          {a.image_url && (
            <img
              src={a.image_url}
              alt="answer"
              style={{ maxWidth: 300, marginTop: 10 }}
            />
          )}
        </div>
      ))}
    </div>
  );
}