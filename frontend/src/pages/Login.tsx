import { useState, useEffect } from 'react';
import { useLocation, Link } from 'wouter';
import { fetchApi } from '../lib/api';
import { useAuth } from '../context/AuthContext';

export default function Login() {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const [, setLocation] = useLocation();
  const { user, login } = useAuth();

  useEffect(() => {
    if (user) {
      setLocation('/');
    }
  }, [user, setLocation]);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    setIsLoading(true);
    try {
      const data = await fetchApi<{ user: any }>('/auth/login', {
        method: 'POST',
        body: JSON.stringify({ email, password })
      });
      login(data.user);
      setLocation('/');
    } catch (err: any) {
      setError(err.message === 'invalid_credentials' ? 'E-mail ou senha incorretos.' : err.message);
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="auth-wrapper">
      <div className="card auth-card">
        <div className="text-center mb-6">
          <h1 style={{ fontSize: '1.875rem', fontWeight: 800, marginBottom: '0.5rem' }}>Bem-vindo de volta</h1>
          <p className="text-muted">Acesse sua conta para gerenciar seu negócio</p>
        </div>

        {error && (
          <div className="alert">
            <span>⚠️</span> {error}
          </div>
        )}

        <form onSubmit={handleSubmit}>
          <div className="form-group">
            <label>Endereço de E-mail</label>
            <input 
              type="email" 
              required 
              placeholder="exemplo@email.com"
              value={email} 
              onChange={e => setEmail(e.target.value)} 
            />
          </div>
          
          <div className="form-group">
            <label>Senha</label>
            <input 
              type="password" 
              required 
              placeholder="••••••••"
              value={password} 
              onChange={e => setPassword(e.target.value)} 
            />
          </div>

          <button type="submit" className="btn-primary mt-4" style={{ width: '100%' }} disabled={isLoading}>
            {isLoading ? 'Entrando...' : 'Entrar na conta'}
          </button>
        </form>

        <div className="mt-6 text-center text-sm">
          <p className="text-muted">
            Não tem uma conta ainda? <Link href="/auth/register" className="font-bold" style={{ color: 'var(--primary)', textDecoration: 'none' }}>Cadastre-se gratuitamente</Link>
          </p>
        </div>
      </div>
    </div>
  );
}
