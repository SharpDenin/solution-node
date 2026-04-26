import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom';
import { AuthProvider, useAuth } from './context/AuthContext';
import { AdminLayout } from './components/layouts/AdminLayout';
import { WorkerLayout } from './components/layouts/WorkerLayout';
import { HomePage } from './pages/HomePage';
import { Login } from './pages/Login';
import { Dashboard } from './pages/Dashboard';
import { ReportDetail } from './pages/ReportDetail';
import { Questions } from './pages/Questions';
import { ChecklistReport } from './pages/ChecklistReport';
import { ThankYou } from './pages/ThankYou';
import { UserCreate } from './pages/UserCreate';
import { ProtectedRoute } from './components/ProtectedRoute';

const AppRoutes = () => {
  const { isLoading } = useAuth();

  if (isLoading) return <div>Загрузка...</div>;

  return (
    <Routes>
      {/* Публичные */}
      <Route path="/" element={<HomePage />} />
      <Route path="/login" element={<Login />} />

      {/* Рабочий (worker) */}
      <Route element={<ProtectedRoute allowedRoles={['worker']}><WorkerLayout /></ProtectedRoute>}>
        <Route path="/checklist/:id" element={<ChecklistReport />} />
        <Route path="/thank-you" element={<ThankYou />} />
      </Route>

      {/* Администратор */}
      <Route element={<ProtectedRoute allowedRoles={['admin']}><AdminLayout /></ProtectedRoute>}>
        <Route path="/dashboard" element={<Dashboard />} />
        <Route path="/reports/:id" element={<ReportDetail />} />
        <Route path="/questions" element={<Questions />} />
        <Route path="/admin/users/create" element={<UserCreate />} />
      </Route>

      {/* Fallback */}
      <Route path="*" element={<Navigate to="/" replace />} />
    </Routes>
  );
};

function App() {
  return (
    <BrowserRouter>
      <AuthProvider>
        <AppRoutes />
      </AuthProvider>
    </BrowserRouter>
  );
}

export default App;