import { createContext, useContext, useState } from 'react';
import type { ReactNode } from 'react';

type User = {
  id: string;
  name: string;
  email: string;
};

type AuthContextType = {
  user: User | null;
  loading: boolean;
  login: (user: User) => void;
  logout: () => void;
};

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export function AuthProvider({ children }: { children: ReactNode }) {
  // Since we rely on HttpOnly cookies, we don't have the token in memory.
  // We store the user info here when they login/register.
  const [user, setUser] = useState<User | null>(null);
  const [loading] = useState(false);

  return (
    <AuthContext.Provider value={{ user, loading, login: setUser, logout: () => setUser(null) }}>
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
