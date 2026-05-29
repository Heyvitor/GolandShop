import { useEffect, useState } from 'react';
import { useRoute } from 'wouter';
import { fetchApi } from '../lib/api';

type Store = {
  id: string;
  name: string;
  slug: string;
};

type Product = {
  id: string;
  name: string;
  description: string;
  price: number;
  variant: string;
  variant_price: number;
  shipping_type: 'free' | 'consult';
};

type StoreCatalog = {
  store: Store;
  items: Product[];
};

export default function StoreView() {
  const [, params] = useRoute('/:slug');
  const [catalog, setCatalog] = useState<StoreCatalog | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  useEffect(() => {
    if (params?.slug) {
      setLoading(true);
      setError('');
      fetchApi<StoreCatalog>(`/stores/catalog/${encodeURIComponent(params.slug)}`)
        .then(setCatalog)
        .catch(() => setError('Loja não encontrada.'))
        .finally(() => setLoading(false));
    }
  }, [params?.slug]);

  if (loading) {
    return (
      <div className="container">
        <div className="storefront-shell">
          <div className="storefront-hero-skeleton shimmer" />
          <div className="storefront-grid">
            {[...Array(3)].map((_, index) => (
              <div key={index} className="store-product-card">
                <div className="skeleton-line shimmer" />
                <div className="skeleton-line shimmer short" />
                <div className="skeleton-line shimmer" />
              </div>
            ))}
          </div>
        </div>
      </div>
    );
  }
  if (error) return <div className="container"><div className="alert">{error}</div></div>;
  if (!catalog) return null;

  return (
    <div className="container">
      <div className="storefront-shell">
        <section className="storefront-hero">
          <span className="badge badge-user">Loja Oficial</span>
          <h1>{catalog.store.name}</h1>
          <p className="text-muted">Confira os produtos disponiveis da loja {catalog.store.slug}.</p>
        </section>

        <section className="storefront-grid">
          {catalog.items.length === 0 ? (
            <div className="card">
              <h3>Sem produtos publicados</h3>
              <p className="text-muted mt-4">Esta loja ainda nao cadastrou produtos na vitrine.</p>
            </div>
          ) : (
            catalog.items.map((product) => (
              <article key={product.id} className="store-product-card">
                <div className="product-chip">
                  {product.shipping_type === 'free' ? 'Frete gratis' : 'Frete a consultar'}
                </div>
                <h3>{product.name}</h3>
                <p className="text-muted">{product.description}</p>
                <div className="store-product-meta">
                  <strong>R$ {product.price.toFixed(2)}</strong>
                  <span>{product.variant || 'Sem variante'}</span>
                  {product.variant_price > 0 && <span>Variante por R$ {product.variant_price.toFixed(2)}</span>}
                </div>
              </article>
            ))
          )}
        </section>
      </div>
    </div>
  );
}
