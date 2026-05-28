import { useEffect, useState } from 'react';
import { useLocation } from 'wouter';
import { fetchApi } from '../lib/api';
import { useAuth } from '../context/AuthContext';

type Store = {
  id: string;
  name: string;
  slug: string;
};

type Product = {
  id: string;
  store_id: string;
  name: string;
  description: string;
  price: number;
  variant: string;
  variant_price: number;
  shipping_type: 'free' | 'consult';
};

export default function Dashboard() {
  const { user, loading, logout } = useAuth();
  const [, setLocation] = useLocation();

  if (loading) return <div className="container">Carregando painel...</div>;
  if (!user) {
    setLocation('/auth/login');
    return null;
  }

  // Renderiza o painel baseado na Role do usuário
  switch (user.role) {
    case 'admin':
      return <AdminPanel user={user} logout={logout} />;
    case 'user':
      return <StorePanel user={user} logout={logout} />;
    default:
      return <ClientPanel user={user} logout={logout} />;
  }
}

/* -------------------------------------------------------------------------- */
/* PAINEL ADMINISTRADOR                                                       */
/* -------------------------------------------------------------------------- */
function AdminPanel({ user, logout }: any) {
  return (
    <div className="container">
      <header className="flex justify-between items-center mb-6">
        <div>
          <h1>Painel Admin</h1>
          <span className="badge badge-admin">Sistema Geral</span>
        </div>
        <div className="flex items-center gap-4">
          <span className="text-muted">Logado como: <strong>{user.email}</strong></span>
          <button className="btn-ghost" onClick={logout}>Sair</button>
        </div>
      </header>
      <div className="grid">
        <div className="card">
          <h3>Monitoramento</h3>
          <p className="mt-4 text-muted">Você tem acesso total a todas as lojas e pedidos do sistema.</p>
        </div>
      </div>
    </div>
  );
}

/* -------------------------------------------------------------------------- */
/* PAINEL VENDEDOR (LOJA)                                                     */
/* -------------------------------------------------------------------------- */
function StorePanel({ user, logout }: any) {
  const [store, setStore] = useState<Store | null>(null);
  const [products, setProducts] = useState<Product[]>([]);
  const [loading, setLoading] = useState(true);
  const [productLoading, setProductLoading] = useState(false);
  const [name, setName] = useState('');
  const [slug, setSlug] = useState('');
  const [showProductForm, setShowProductForm] = useState(false);
  const [productName, setProductName] = useState('');
  const [description, setDescription] = useState('');
  const [price, setPrice] = useState('');
  const [variant, setVariant] = useState('');
  const [variantPrice, setVariantPrice] = useState('');
  const [shippingType, setShippingType] = useState<'free' | 'consult'>('consult');

  useEffect(() => {
    fetchApi<Store>('/stores/mine')
      .then(setStore)
      .catch(() => setStore(null))
      .finally(() => setLoading(false));
  }, []);

  useEffect(() => {
    if (!store) return;

    setProductLoading(true);
    fetchApi<{ items: Product[] }>('/items')
      .then((data) => setProducts(data.items.filter((item) => item.store_id === store.id)))
      .catch(() => setProducts([]))
      .finally(() => setProductLoading(false));
  }, [store]);

  const handleCreateStore = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      const newStore = await fetchApi<Store>('/stores', {
        method: 'POST',
        body: JSON.stringify({ name, slug })
      });
      setStore(newStore);
    } catch (err: any) {
      alert(err.message);
    }
  };

  const handleCreateProduct = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!store) return;

    try {
      const newProduct = await fetchApi<Product>('/items', {
        method: 'POST',
        body: JSON.stringify({
          store_id: store.id,
          name: productName,
          description,
          price: Number(price),
          variant,
          variant_price: Number(variantPrice || '0'),
          shipping_type: shippingType
        })
      });

      setProducts((current) => [newProduct, ...current]);
      setProductName('');
      setDescription('');
      setPrice('');
      setVariant('');
      setVariantPrice('');
      setShippingType('consult');
      setShowProductForm(false);
    } catch (err: any) {
      alert(err.message);
    }
  };

  if (loading) return <div className="container">Carregando loja...</div>;

  return (
    <div className="container">
      <header className="flex justify-between items-center mb-6">
        <div>
          <h1>Minha Loja</h1>
          <span className="badge badge-user">Vendedor</span>
        </div>
        <div className="flex items-center gap-4">
          <span className="text-muted">Usuário: <strong>{user.email}</strong></span>
          <button className="btn-ghost" onClick={logout}>Sair</button>
        </div>
      </header>

      {!store ? (
        <div className="card auth-card" style={{ margin: '0 auto' }}>
          <h2>Crie sua Loja</h2>
          <p className="text-muted mb-6">Parece que você ainda não tem uma loja ativa.</p>
          <form onSubmit={handleCreateStore} className="flex flex-col gap-4">
            <div className="form-group">
              <label>Nome da Loja</label>
              <input required value={name} onChange={e => setName(e.target.value)} placeholder="Ex: Minha Loja Top" />
            </div>
            <div className="form-group">
              <label>Slug da URL (Único)</label>
              <input required value={slug} onChange={e => setSlug(e.target.value)} placeholder="ex: minha-loja-01" />
            </div>
            <button type="submit" className="btn-primary">Ativar Loja Agora</button>
          </form>
        </div>
      ) : (
        <div className="grid">
          <div className="card">
            <h3>{store.name}</h3>
            <p className="text-muted">Slug: /{store.slug}</p>
            <div className="mt-6 flex gap-4">
              <button className="btn-primary" onClick={() => setShowProductForm((value) => !value)}>
                {showProductForm ? 'Fechar Cadastro' : 'Gerenciar Produtos'}
              </button>
            </div>
          </div>
          <div className="card">
            <div className="flex justify-between items-center mb-4">
              <div>
                <h3>Produtos</h3>
                <p className="text-muted">Cadastre nome, descricao, valor, variante e frete.</p>
              </div>
              <strong>{products.length}</strong>
            </div>

            {showProductForm && (
              <form onSubmit={handleCreateProduct} className="flex flex-col gap-4 mb-6">
                <div className="form-group">
                  <label>Nome do produto</label>
                  <input required value={productName} onChange={e => setProductName(e.target.value)} placeholder="Ex: Camiseta Premium" />
                </div>
                <div className="form-group">
                  <label>Descricao</label>
                  <textarea
                    required
                    value={description}
                    onChange={e => setDescription(e.target.value)}
                    placeholder="Detalhes, material, medidas, beneficios"
                  />
                </div>
                <div className="product-form-grid">
                  <div className="form-group">
                    <label>Valor</label>
                    <input required type="number" min="0" step="0.01" value={price} onChange={e => setPrice(e.target.value)} placeholder="99.90" />
                  </div>
                  <div className="form-group">
                    <label>Variante</label>
                    <input value={variant} onChange={e => setVariant(e.target.value)} placeholder="Ex: Tamanho P / Vermelho" />
                  </div>
                  <div className="form-group">
                    <label>Valor da variante</label>
                    <input type="number" min="0" step="0.01" value={variantPrice} onChange={e => setVariantPrice(e.target.value)} placeholder="0.00" />
                  </div>
                  <div className="form-group">
                    <label>Frete</label>
                    <select value={shippingType} onChange={e => setShippingType(e.target.value as 'free' | 'consult')}>
                      <option value="consult">A consultar</option>
                      <option value="free">Frete gratuito</option>
                    </select>
                  </div>
                </div>
                <button type="submit" className="btn-primary">Cadastrar produto</button>
              </form>
            )}

            {productLoading ? (
              <p className="text-muted">Carregando produtos...</p>
            ) : products.length === 0 ? (
              <p className="text-muted">Nenhum produto cadastrado ainda.</p>
            ) : (
              <div className="product-list">
                {products.map((product) => (
                  <div key={product.id} className="product-row">
                    <div>
                      <h4>{product.name}</h4>
                      <p className="text-muted">{product.description}</p>
                    </div>
                    <div className="product-meta">
                      <strong>R$ {product.price.toFixed(2)}</strong>
                      <span className="text-muted">{product.variant || 'Sem variante'}</span>
                      <span className="text-muted">
                        {product.variant_price > 0 ? `Variante: R$ ${product.variant_price.toFixed(2)}` : 'Variante sem acrescimo'}
                      </span>
                      <span className="text-muted">
                        {product.shipping_type === 'free' ? 'Frete gratuito' : 'Frete a consultar'}
                      </span>
                    </div>
                  </div>
                ))}
              </div>
            )}
          </div>
        </div>
      )}
    </div>
  );
}

/* -------------------------------------------------------------------------- */
/* PAINEL CLIENTE                                                             */
/* -------------------------------------------------------------------------- */
function ClientPanel({ user, logout }: any) {
  return (
    <div className="container">
      <header className="flex justify-between items-center mb-6">
        <div>
          <h1>Meus Pedidos</h1>
          <span className="badge badge-client">Cliente</span>
        </div>
        <div className="flex items-center gap-4">
          <span className="text-muted">Olá, <strong>{user.email}</strong></span>
          <button className="btn-ghost" onClick={logout}>Sair</button>
        </div>
      </header>
      <div className="card">
        <p className="text-center text-muted py-10">Você ainda não realizou nenhum pedido.</p>
        <div className="text-center mt-4">
          <button className="btn-primary">Explorar Lojas</button>
        </div>
      </div>
    </div>
  );
}
