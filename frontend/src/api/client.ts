import axios from "axios";

export const api = axios.create({
  baseURL: "http://localhost:8090",
});

api.interceptors.request.use((config) => {
  const token = localStorage.getItem("token");

  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }

  return config;
});

export const getToken = () => localStorage.getItem("token");

export const setToken = (token: string) => {
  localStorage.setItem("token", token);
};

export const logout = () => {
  localStorage.removeItem("token");
};