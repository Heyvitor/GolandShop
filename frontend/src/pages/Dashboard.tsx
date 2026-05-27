import { useEffect, useState } from 'react';
import { useLocation } from 'wouter';
import { fetchApi } from '../lib/api';
import { useAuth } from '../context/AuthContext';

type Item = {
  id: string;
  title: string;
  body: string;
  created_at: string;
};

export default function Dashboard() {
  const [items, setItems] = useState<Item[]>([]);
  const [title, setTitle] = useState('');
  const [body, setBody] = useState('');
  const { user, logout } = useAuth();
  const [, setLocation] = useLocation();

  useEffect(() => {
    loadItems();
  }, []);

  const loadItems = async () => {
    try {
      const data = await fetchApi<{ items: Item[] }>('/items');
      setItems(data.items || []);
    } catch (err: any) {
      if (err.message === 'missing_token' || err.message === 'invalid_token') {
        logout();
        setLocation('/login');
      }
    }
  };

  const handleCreate = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      await fetchApi('/items', {
        method: 'POST',
        body: JSON.stringify({ title, body })
      });
      setTitle('');
      setBody('');
      loadItems();
    } catch (err) {
      console.error(err);
    }
  };

  const handleLogout = async () => {
    await fetchApi('/auth/logout', { method: 'POST' }).catch(() => {});
    logout();
    setLocation('/login');
  };

  return (
    <div className="container">
      <header className="flex justify-between items-center mb-4">
        <h2>Dashboard</h2>
        <div className="flex gap-4 items-center">
          <span className="text-muted">Olá, {user?.name || 'Usuário'}</span>
          <button className="btn-danger" onClick={handleLogout}>Sair</button>
        </div>
      </header>

      <div className="card mb-4">
        <h3>Criar novo item</h3>
        <form onSubmit={handleCreate} className="flex flex-col gap-4 mt-4">
          <div>
            <label>Título</label>
            <input required value={title} onChange={e => setTitle(e.target.value)} />
          </div>
          <div>
            <label>Descrição</label>
            <input required value={body} onChange={e => setBody(e.target.value)} />
          </div>
          <button type="submit" className="btn-primary" style={{ alignSelf: 'flex-start' }}>Adicionar</button>
        </form>
      </div>

      <div className="flex flex-col gap-4">
        {items.map(item => (
          <div key={item.id} className="card">
            <h4>{item.title}</h4>
            <p className="mt-4 text-muted">{item.body}</p>
          </div>
        ))}
        {items.length === 0 && <p className="text-center text-muted">Nenhum item encontrado.</p>}
      </div>
    </div>
  );
}
