import { useState } from "react";
import { api } from "../api/client";
import { setToken } from "../api/client";
import { useNavigate } from "react-router-dom";

export default function Login() {
  const [login, setLogin] = useState("");
  const [password, setPassword] = useState("");

  const nav = useNavigate();

  const handleLogin = async () => {
    try {
      const res = await api.post("/login", { login, password });

      setToken(res.data.token);

      nav("/");
    } catch {
      alert("invalid credentials");
    }
  };

  return (
    <div style={{ padding: 40 }}>
      <h2>Login</h2>

      <input placeholder="login" onChange={e => setLogin(e.target.value)} />
      <br /><br />

      <input type="password" placeholder="password" onChange={e => setPassword(e.target.value)} />
      <br /><br />

      <button onClick={handleLogin}>Login</button>
    </div>
  );
}