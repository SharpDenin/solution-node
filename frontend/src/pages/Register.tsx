import { useState } from "react";
import { api } from "../api/client";
import { useNavigate } from "react-router-dom";

export default function Register() {
  const [fullName, setFullName] = useState("");
  const [login, setLogin] = useState("");
  const [password, setPassword] = useState("");

  const nav = useNavigate();

  const submit = async () => {
    await api.post("/register", {
      full_name: fullName,
      login,
      password,
    });

    alert("Registered!");
    nav("/login");
  };

  return (
    <div style={{ padding: 40 }}>
      <h2>Register</h2>

      <input placeholder="Full name" onChange={e => setFullName(e.target.value)} />
      <br /><br />

      <input placeholder="Login" onChange={e => setLogin(e.target.value)} />
      <br /><br />

      <input type="password" placeholder="Password" onChange={e => setPassword(e.target.value)} />
      <br /><br />

      <button onClick={submit}>Create account</button>
    </div>
  );
}