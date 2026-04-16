import { Navigate } from "react-router-dom";

export default function ProtectedRoute({ children }: { children: any }) {
  const token = localStorage.getItem("token");

  if (!token) {
    return <Navigate to="/login" />;
  }

  return children;
}