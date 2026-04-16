import { BrowserRouter, Routes, Route } from "react-router-dom";

import Login from "./pages/Login";
import Register from "./pages/Register";
import Dashboard from "./pages/Dashboard";
import Reports from "./pages/Reports";
import Questions from "./pages/Questions";
import CreateReport from "./pages/CreateReport";
import ReportDetail from "./pages/ReportDetail";
import Layout from "./components/Layout";
import ProtectedRoute from "./auth/ProtectedRoute";

export default function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/login" element={<Login />} />
        <Route path="/register" element={<Register />} />

        <Route
          element={
            <ProtectedRoute>
              <Layout />
            </ProtectedRoute>
          }
        >
          <Route path="/" element={<Dashboard />} />
          <Route path="/reports" element={<Reports />} />
          <Route path="/reports/:id" element={<ReportDetail />} />
          <Route path="/questions" element={<Questions />} />
          <Route path="/create-report" element={<CreateReport />} />
        </Route>
      </Routes>
    </BrowserRouter>
  );
}