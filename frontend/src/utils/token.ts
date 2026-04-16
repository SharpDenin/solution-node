import { jwtDecode } from 'jwt-decode';

export interface TokenPayload {
  user_id: string;
  role: 'admin' | 'worker';
  exp: number;
  // возможно, full_name? если нет, добавим позже
}

export const setToken = (token: string) => {
  localStorage.setItem('token', token);
};

export const getToken = () => localStorage.getItem('token');

export const removeToken = () => {
  localStorage.removeItem('token');
};

export const decodeToken = (): TokenPayload | null => {
  const token = getToken();
  if (!token) return null;
  try {
    return jwtDecode<TokenPayload>(token);
  } catch {
    return null;
  }
};

export const isTokenExpired = (): boolean => {
  const payload = decodeToken();
  if (!payload) return true;
  return payload.exp * 1000 < Date.now();
};