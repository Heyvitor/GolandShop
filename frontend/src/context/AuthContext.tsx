import { createContext, useContext, useState, useEffect } from 'react';
import type { ReactNode } from 'react';
import { fetchApi } from '../lib/api';

type User = {
  id: string;
  name?: string;
  email?: string;
  role: string;
};

type AuthContextType = {
  user: User | null;
  loading: boolean;
  login: (user: User) => void;
  logout: () => void;
};

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export function AuthProvider({ children }: { children: ReactNode }) {
  const [user, setUser] = useState<User | null>(null);
  const [loading, setLoading] = useState(true);

  // Tenta restaurar a sessão ao carregar a página
  useEffect(() => {
    async function restoreSession() {
      // Se já estivermos na página de login ou registro, não precisamos "bloquear"
      // Mas ainda é bom checar caso o usuário já esteja logado para redirecionar.
      try {
        const data = await fetchApi<User>('/auth/me');
        setUser(data);
      } catch (err) {
        setUser(null);
      } finally {
        setLoading(false);
      }
    }
    
    // Pequena otimização: se for rota pública de loja (ex: /loja01), 
    // ou rota de auth, poderíamos pular, mas o 'me' é rápido.
    // O problema do usuário é que ele "pede" o /me mesmo sem estar logado.
    // Isso é normal para checar sessão, mas vamos garantir que não quebre nada.
    restoreSession();
  }, []);

  const login = (userData: User) => {
    setUser(userData);
  };

  const logout = () => {
    setUser(null);
  };

  return (
    <AuthContext.Provider value={{ user, loading, login, logout }}>
      {children}
    </AuthContext.Provider>
  );
}

export function useAuth() {
  const context = useContext(AuthContext);
  if (context === undefined) {
    throw new Error('useAuth must be used within an AuthProvider');
  }
  return context;
}
