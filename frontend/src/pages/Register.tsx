import { useState } from 'react';
import { useLocation, Link } from 'wouter';
import { fetchApi } from '../lib/api';
import { useAuth } from '../context/AuthContext';

export default function Register() {
  const [name, setName] = useState('');
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [role, setRole] = useState('client');
  const [error, setError] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const [, setLocation] = useLocation();
  const { login } = useAuth();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    setIsLoading(true);
    try {
      const data = await fetchApi<{ user: any }>('/auth/register', {
        method: 'POST',
        body: JSON.stringify({ name, email, password, role })
      });
      login(data.user);
      setLocation('/');
    } catch (err: any) {
      setError(err.message);
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="auth-wrapper">
      <div className="card auth-card">
        <div className="text-center mb-6">
          <h1 style={{ fontSize: '1.875rem', fontWeight: 800, marginBottom: '0.5rem' }}>Criar sua conta</h1>
          <p className="text-muted">Junte-se a milhares de lojas e clientes</p>
        </div>

        {error && (
          <div className="alert">
            <span>⚠️</span> {error}
          </div>
        )}

        <form onSubmit={handleSubmit}>
          <div className="form-group">
            <label>Nome Completo</label>
            <input 
              type="text" required 
              placeholder="Seu nome"
              value={name} 
              onChange={e => setName(e.target.value)} 
            />
          </div>

          <div className="form-group">
            <label>Endereço de E-mail</label>
            <input 
              type="email" required 
              placeholder="exemplo@email.com"
              value={email} 
              onChange={e => setEmail(e.target.value)} 
            />
          </div>
          
          <div className="form-group">
            <label>Senha</label>
            <input 
              type="password" required 
              placeholder="Mínimo 8 caracteres"
              value={password} 
              onChange={e => setPassword(e.target.value)} 
              minLength={8}
            />
          </div>

          <div className="form-group">
            <label>Tipo de Conta</label>
            <select value={role} onChange={e => setRole(e.target.value)}>
              <option value="client">Cliente (Quero comprar)</option>
              <option value="user">Vendedor (Tenho uma loja)</option>
              <option value="admin">Administrador (Gestão)</option>
            </select>
          </div>

          <button type="submit" className="btn-primary mt-4" style={{ width: '100%' }} disabled={isLoading}>
            {isLoading ? 'Criando conta...' : 'Registrar agora'}
          </button>
        </form>

        <div className="mt-6 text-center text-sm">
          <p className="text-muted">
            Já possui uma conta? <Link href="/login" className="font-bold" style={{ color: 'var(--primary)', textDecoration: 'none' }}>Entrar na conta</Link>
          </p>
        </div>
      </div>
    </div>
  );
}
