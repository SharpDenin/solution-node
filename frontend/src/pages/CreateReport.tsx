import { useEffect, useState } from "react";
import { api } from "../api/client";

type Question = {
  id: string;
  text: string;
};

type Answer = {
  question_id: string;
  text: string;
  image_url?: string;
};

export default function CreateReport() {
  const [questions, setQuestions] = useState<Question[]>([]);
  const [answers, setAnswers] = useState<Record<string, Answer>>({});

  const [place, setPlace] = useState("");
  const [date, setDate] = useState("");
  const [responsible, setResponsible] = useState("");

  const isValid =
    place.trim() !== "" &&
    date.trim() !== "" &&
    Object.values(answers).length > 0;

  useEffect(() => {
    api.get("/questions").then(res => {
      setQuestions(res.data);
    });
  }, []);

  const updateAnswer = (qId: string, value: Partial<Answer>) => {
    setAnswers(prev => ({
      ...prev,
      [qId]: {
        ...prev[qId],
        question_id: qId,
        text: prev[qId]?.text || "",
        ...value,
      },
    }));
  };

  const uploadImage = async (qId: string, file: File) => {
    const formData = new FormData();
    formData.append("file", file);

    const res = await api.post("/upload", formData, {
      headers: { "Content-Type": "multipart/form-data" },
    });

    updateAnswer(qId, { image_url: res.data.url });
  };

  const submit = async () => {
    if (!isValid) {
      alert("Fill required fields");
      return;
    }

    const payload = {
      place,
      report_date: date,
      responsible_name: responsible,
      answers: Object.values(answers),
    };

    await api.post("/reports", payload);

    alert("Report submitted!");
  };

  return (
    <div style={{ maxWidth: 800, margin: "0 auto" }}>
      <h2>New Report</h2>

      <div style={{ marginBottom: 20 }}>
        <input placeholder="Place" value={place} onChange={e => setPlace(e.target.value)} />
        <br /><br />

        <input type="date" value={date} onChange={e => setDate(e.target.value)} />
        <br /><br />

        <input
          placeholder="Responsible name"
          value={responsible}
          onChange={e => setResponsible(e.target.value)}
        />
      </div>

      <hr />

      {questions.map(q => (
        <div
          key={q.id}
          style={{
            marginBottom: 20,
            padding: 15,
            background: "#fff",
            borderRadius: 10,
            boxShadow: "0 1px 3px rgba(0,0,0,0.1)",
          }}
        >
          <b>{q.text}</b>

          <br /><br />

          <textarea
            placeholder="Answer"
            onChange={e => updateAnswer(q.id, { text: e.target.value })}
          />

          <br /><br />

          <input
            type="file"
            onChange={e => {
              const file = e.target.files?.[0];
              if (file) uploadImage(q.id, file);
            }}
          />

          {answers[q.id]?.image_url && (
            <div>
              <img
                src={answers[q.id].image_url}
                style={{ maxWidth: 200, marginTop: 10 }}
              />
            </div>
          )}
        </div>
      ))}

      <button disabled={!isValid} onClick={submit}>
        Close shift
      </button>
    </div>
  );
}