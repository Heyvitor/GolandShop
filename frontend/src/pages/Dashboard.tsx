import { useEffect, useState } from 'react';
import { useLocation } from 'wouter';
import { fetchApi } from '../lib/api';
import { useAuth } from '../context/AuthContext';

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
  const [store, setStore] = useState<any>(null);
  const [loading, setLoading] = useState(true);
  const [name, setName] = useState('');
  const [slug, setSlug] = useState('');

  useEffect(() => {
    fetchApi('/stores/view?owner=true') // Endpoint fictício por enquanto
      .then(setStore)
      .catch(() => setStore(null))
      .finally(() => setLoading(false));
  }, []);

  const handleCreateStore = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      const newStore = await fetchApi<any>('/stores', {
        method: 'POST',
        body: JSON.stringify({ name, slug })
      });
      setStore(newStore);
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
            <div className="mt-6">
              <button className="btn-primary">Gerenciar Produtos</button>
            </div>
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
