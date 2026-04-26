import { useEffect, useState } from 'react';
import { useParams, Navigate } from 'react-router-dom';
import { api } from '../api/client';
import type { Checklist } from '../types';

export const ChecklistEntry = () => {
  const { id } = useParams<{ id: string }>();
  const [redirectTo, setRedirectTo] = useState<string | null>(null);

  useEffect(() => {
    if (!id) {
      setRedirectTo('/');
      return;
    }
    api.get('/api/checklists')
      .then(res => {
        const checklists: Checklist[] = res.data;
        const cl = checklists.find(c => c.id === id);
        if (!cl) {
          setRedirectTo('/');
          return;
        }
        // Для чек-листа по фенофазам сначала выбираем сорт
        if (cl.code === 'sort_control') {
          setRedirectTo(`/checklist/${id}/variety`);
        } else {
          // Для остальных (default) сразу на заполнение
          setRedirectTo(`/checklist/${id}/fill`);
        }
      })
      .catch(() => setRedirectTo('/'));
  }, [id]);

  if (redirectTo) {
    return <Navigate to={redirectTo} replace />;
  }
  return <div>Загрузка чек-листа...</div>;
};