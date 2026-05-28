import { useEffect, useState } from 'react';
import { useRoute } from 'wouter';
import { fetchApi } from '../lib/api';

export default function StoreView() {
  const [, params] = useRoute('/:slug');
  const [store, setStore] = useState<any>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  useEffect(() => {
    if (params?.slug) {
      fetchApi<any>(`/stores/view?slug=${params.slug}`)
        .then(setStore)
        .catch(() => setError('Loja não encontrada.'))
        .finally(() => setLoading(false));
    }
  }, [params?.slug]);

  if (loading) return <div className="container">Carregando loja...</div>;
  if (error) return <div className="container"><div className="alert">{error}</div></div>;

  return (
    <div className="container">
      <div className="card text-center">
        <h1 style={{ fontSize: '2.5rem', marginBottom: '1rem' }}>{store.name}</h1>
        <p className="text-muted">Bem-vindo à nossa loja oficial!</p>
        <hr className="mt-6 mb-6" style={{ border: 'none', borderTop: '1px solid var(--border)' }} />
        <div className="grid">
          <p>Vitrine de produtos em breve...</p>
        </div>
      </div>
    </div>
  );
}
