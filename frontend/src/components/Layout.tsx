import { Outlet, Link, useNavigate } from "react-router-dom";
import { logout } from "../api/client";

export default function Layout() {
  const nav = useNavigate();

  const handleLogout = () => {
    logout();
    nav("/login");
  };

  return (
    <div style={{ display: "flex", height: "100vh", fontFamily: "Arial" }}>
      <aside
        style={{
          width: 220,
          background: "#1f2937",
          color: "white",
          padding: 20,
        }}
      >
        <h3>Reports</h3>

        <nav style={{ display: "flex", flexDirection: "column", gap: 10 }}>
          <Link to="/" style={{ color: "white" }}>Dashboard</Link>
          <Link to="/reports" style={{ color: "white" }}>Reports</Link>
          <Link to="/create-report" style={{ color: "white" }}>New Report</Link>
          <Link to="/questions" style={{ color: "white" }}>Questions</Link>
        </nav>

        <button
          onClick={handleLogout}
          style={{ marginTop: 20, padding: 8 }}
        >
          Logout
        </button>
      </aside>

      <main style={{ flex: 1, padding: 20, background: "#f3f4f6" }}>
        <Outlet />
      </main>
    </div>
  );
}