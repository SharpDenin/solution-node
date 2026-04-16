import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom';
import { AuthProvider, useAuth } from './context/AuthContext';
import { AdminLayout } from './components/layouts/AdminLayout';
import { WorkerLayout } from './components/layouts/WorkerLayout';
import { Login } from './pages/Login';
import { Register } from './pages/Register';
import { Dashboard } from './pages/Dashboard';
import { CreateReport } from './pages/CreateReport';
import { ThankYou } from './pages/ThankYou';
import { ReportDetail } from './pages/ReportDetail';
import { Questions } from './pages/Questions';
import { ProtectedRoute } from './components/ProtectedRoute';

const AppRoutes = () => {
  const { role, isLoading } = useAuth();

  if (isLoading) return <div>Загрузка...</div>;

  return (
    <Routes>
      <Route path="/login" element={<Login />} />
      <Route path="/register" element={<Register />} />

      <Route element={<ProtectedRoute allowedRoles={['admin']}><AdminLayout /></ProtectedRoute>}>
        <Route path="/dashboard" element={<Dashboard />} />
        <Route path="/reports/:id" element={<ReportDetail />} />
        <Route path="/questions" element={<Questions />} />
        <Route path="/" element={<Navigate to="/dashboard" replace />} />
      </Route>

      <Route element={<ProtectedRoute allowedRoles={['worker']}><WorkerLayout /></ProtectedRoute>}>
        <Route path="/create-report" element={<CreateReport />} />
        <Route path="/thank-you" element={<ThankYou />} />
        <Route path="/" element={<Navigate to="/create-report" replace />} />
      </Route>
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