import { useState } from 'react';
import { useLocation, Link } from 'wouter';
import { fetchApi } from '../lib/api';
import { useAuth } from '../context/AuthContext';

export default function Login() {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');
  const [, setLocation] = useLocation();
  const { login } = useAuth();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    try {
      const data = await fetchApi<{ user: any }>('/auth/login', {
        method: 'POST',
        body: JSON.stringify({ email, password })
      });
      login(data.user);
      setLocation('/');
    } catch (err: any) {
      setError(err.message);
    }
  };

  return (
    <div className="container" style={{ maxWidth: '400px', marginTop: '4rem' }}>
      <div className="card">
        <h2 className="text-center mb-4">Login</h2>
        {error && <div className="alert">{error}</div>}
        <form onSubmit={handleSubmit} className="flex flex-col gap-4">
          <div>
            <label>Email</label>
            <input type="email" required value={email} onChange={e => setEmail(e.target.value)} />
          </div>
          <div>
            <label>Password</label>
            <input type="password" required value={password} onChange={e => setPassword(e.target.value)} />
          </div>
          <button type="submit" className="btn-primary mt-4">Entrar</button>
        </form>
        <p className="text-center mt-4 text-muted">
          Não tem conta? <Link href="/register">Cadastre-se</Link>
        </p>
      </div>
    </div>
  );
}
