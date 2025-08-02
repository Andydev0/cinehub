import React, { createContext, useState, useContext, useEffect, ReactNode } from 'react';
import api from '../services/api';

// Define a "forma" do nosso contexto
interface AuthContextType {
  token: string | null;
  isLogado: boolean;
  login: (token: string) => void;
  logout: () => void;
}

// Cria o contexto com um valor inicial undefined
const AuthContext = createContext<AuthContextType | undefined>(undefined);

// Cria o Provedor do contexto
export const AuthProvider: React.FC<{ children: ReactNode }> = ({ children }) => {
  const [token, setToken] = useState<string | null>(null);

  // Ao carregar, verifica se já existe um token no localStorage
  useEffect(() => {
    const storedToken = localStorage.getItem('authToken');
    if (storedToken) {
      setToken(storedToken);
      api.defaults.headers.common['Authorization'] = `Bearer ${storedToken}`;
    }
  }, []);

  const login = (newToken: string) => {
    localStorage.setItem('authToken', newToken);
    setToken(newToken);
    api.defaults.headers.common['Authorization'] = `Bearer ${newToken}`;
  };

  const logout = () => {
    localStorage.removeItem('authToken');
    setToken(null);
    delete api.defaults.headers.common['Authorization'];
    window.location.href = '/login';
  };

  const value = {
    token,
    isLogado: !!token, // isLogado é true se o token não for nulo
    login,
    logout,
  };

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
};

// Cria um hook customizado para facilitar o uso do contexto
export const useAuth = () => {
  const context = useContext(AuthContext);
  if (context === undefined) {
    throw new Error('useAuth deve ser usado dentro de um AuthProvider');
  }
  return context;
};