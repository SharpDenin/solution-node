import { useEffect, useState } from "react";
import { api } from "../api/client";

export default function Questions() {
  const [questions, setQuestions] = useState<any[]>([]);
  const [text, setText] = useState("");

  const load = () => {
    api.get("/questions").then(res => setQuestions(res.data));
  };

  useEffect(load, []);

  const create = async () => {
    await api.post("/questions", { text, order_index: 1 });
    setText("");
    load();
  };

  const remove = async (id: string) => {
    await api.delete(`/questions/${id}`);
    load();
  };

  return (
    <div>
      <h2>Questions</h2>

      <input value={text} onChange={e => setText(e.target.value)} />
      <button onClick={create}>Add</button>

      <ul>
        {questions.map(q => (
          <li key={q.id}>
            {q.text}
            <button onClick={() => remove(q.id)}>Delete</button>
          </li>
        ))}
      </ul>
    </div>
  );
}