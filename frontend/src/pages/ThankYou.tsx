import { Link } from 'react-router-dom';

export const ThankYou = () => {
  return (
    <div style={styles.container}>
      <div style={styles.card}>
        <h2>✅ Спасибо за отчёт!</h2>
        <p>Ваш отчёт успешно отправлен.</p>
        <Link to="/create-report" style={styles.button}>Создать новый отчёт</Link>
      </div>
    </div>
  );
};

const styles = {
  container: {
    display: 'flex',
    justifyContent: 'center',
    alignItems: 'center',
    minHeight: '60vh',
  },
  card: {
    background: 'white',
    padding: '32px',
    borderRadius: '16px',
    textAlign: 'center' as const,
    boxShadow: '0 4px 12px rgba(0,0,0,0.05)',
  },
  button: {
    display: 'inline-block',
    marginTop: '16px',
    background: '#16a34a',
    color: 'white',
    padding: '10px 20px',
    borderRadius: '8px',
    textDecoration: 'none',
  },
};