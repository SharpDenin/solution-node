import { useEffect, useState } from "react";
import { api } from "../api/client";
import { useNavigate } from "react-router-dom";

export default function Dashboard() {
  const [reports, setReports] = useState<any[]>([]);
  const nav = useNavigate();

  useEffect(() => {
    // админский список (если нет прав — просто проигнорится)
    api.get("/reports?limit=5")
      .then(res => setReports(res.data))
      .catch(() => {});
  }, []);

  return (
    <div>
      <h2>Dashboard</h2>

      {/* 🚀 быстрые действия */}
      <div style={{ marginBottom: 20 }}>
        <button onClick={() => nav("/create-report")}>
          ➕ Create Report
        </button>

        <button
          style={{ marginLeft: 10 }}
          onClick={() =>
            window.open("http://localhost:8090/reports/export", "_blank")
          }
        >
          📊 Export Excel
        </button>
      </div>

      <hr />

      {/* 📄 последние отчеты */}
      <h3>Recent Reports</h3>

      {reports.length === 0 && <p>No data or no access</p>}

      {reports.map(r => (
        <div
          key={r.id}
          style={{
            background: "white",
            padding: 10,
            marginBottom: 10,
            borderRadius: 6,
            cursor: "pointer",
          }}
          onClick={() => nav(`/reports/${r.id}`)}
        >
          <b>{r.place}</b> — {r.report_date}
          <div>{r.responsible_name}</div>
        </div>
      ))}
    </div>
  );
}