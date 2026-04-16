import { useEffect, useState } from "react";
import { api } from "../api/client";
import { Link } from "react-router-dom";

export default function Reports() {
  const [reports, setReports] = useState<any[]>([]);

  useEffect(() => {
    api.get("/reports").then(res => setReports(res.data));
  }, []);

  return (
    <div>
      <h2>Reports</h2>

      <button
        onClick={() =>
          window.open("http://localhost:8090/reports/export", "_blank")
        }
      >
        Export Excel
      </button>

      <table>
        <thead>
          <tr>
            <th>Place</th>
            <th>Date</th>
            <th>Responsible</th>
          </tr>
        </thead>
        <tbody>
          {reports.map(r => (
            <tr key={r.id}>
              <td>{r.place}</td>
              <td>{r.report_date}</td>
              <td>{r.responsible_name}</td>
              <td>
                <Link to={`/reports/${r.id}`}>Open</Link>
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}