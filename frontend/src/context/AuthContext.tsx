import React, { createContext, useState, useEffect, type ReactNode } from 'react';
import type { UserRole } from '../types';
import { getToken, setToken, decodeToken, removeToken, isTokenExpired } from '../utils/token';
import { api } from '../api/client';

interface AuthContextType {
  token: string | null;
  role: UserRole | null;
  userId: string | null;
  fullName: string | null;
  position: string | null;
  login: (token: string) => void;
  logout: () => void;
  isLoading: boolean;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export const AuthProvider: React.FC<{ children: ReactNode }> = ({ children }) => {
  const [token, setTokenState] = useState<string | null>(null);
  const [role, setRole] = useState<UserRole | null>(null);
  const [userId, setUserId] = useState<string | null>(null);
  const [fullName, setFullName] = useState<string | null>(null);
  const [position, setPosition] = useState<string | null>(null);
  const [isLoading, setIsLoading] = useState(true);

  const fetchCurrentUser = async (token: string) => {
    try {
      const res = await api.get('/api/me', {
        headers: { Authorization: `Bearer ${token}` },
      });
      const user = res.data;
      setFullName(user.full_name);
      setPosition(user.position || null);
    } catch {
      // ignore
    }
  };

  useEffect(() => {
    const storedToken = getToken();
    if (storedToken && !isTokenExpired()) {
      const payload = decodeToken();
      if (payload) {
        setTokenState(storedToken);
        setRole(payload.role);
        setUserId(payload.user_id);
        fetchCurrentUser(storedToken);
      } else {
        removeToken();
      }
    }
    setIsLoading(false);
  }, []);

  const login = async (newToken: string) => {
    setTokenState(newToken);
    const payload = decodeToken();
    if (payload) {
      setRole(payload.role);
      setUserId(payload.user_id);
    }
    setToken(newToken);
    await fetchCurrentUser(newToken);
  };

  const logout = () => {
    setTokenState(null);
    setRole(null);
    setUserId(null);
    setFullName(null);
    setPosition(null);
    removeToken();
  };

  return (
    <AuthContext.Provider value={{ token, role, userId, fullName, position, login, logout, isLoading }}>
      {children}
    </AuthContext.Provider>
  );
};

export const useAuth = () => {
  const context = React.useContext(AuthContext);
  if (!context) throw new Error('useAuth must be used within AuthProvider');
  return context;
};