import React, { useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import api from '../services/api';
import { useAuth } from '../context/AuthContext';

const LoginPage: React.FC = () => {
  const [email, setEmail] = useState('');
  const [senha, setSenha] = useState('');
  const [erro, setErro] = useState<string | null>(null);
  const [carregando, setCarregando] = useState(false);
  const [mostrarSenha, setMostrarSenha] = useState(false);
  const { login } = useAuth();
  const navigate = useNavigate();

  const handleLogin = async (e: React.FormEvent) => {
    e.preventDefault();
    setCarregando(true);
    setErro(null);

    try {
      const response = await api.post('/auth/login', { email, senha });
      login(response.data.token);
      navigate('/');
    } catch (error: any) {
      const msgErro = error.response?.data?.erro || 'Falha ao fazer login. Tente novamente.';
      setErro(msgErro);
    } finally {
      setCarregando(false);
    }
  };

  return (
    <div className="min-h-screen flex items-center justify-center py-12 px-4">
      <div className="w-full max-w-md">
        {/* Logo/Header */}
        <div className="text-center mb-8 fade-in-up">
          <div className="text-6xl mb-4">🎬</div>
          <h1 className="text-4xl font-bold text-white mb-2">
            Bem-vindo de <span className="text-orange">volta</span>!
          </h1>
          <p className="text-gray-300">
            Entre na sua conta para continuar explorando filmes incríveis
          </p>
        </div>

        {/* Form */}
        <div className="glass-effect rounded-2xl p-8 fade-in-up">
          <form onSubmit={handleLogin} className="space-y-6">
            {/* Email */}
            <div>
              <label htmlFor="email" className="block text-sm font-medium text-gray-300 mb-2">
                📧 Email
              </label>
              <input
                id="email"
                type="email"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                required
                placeholder="seu@email.com"
                className="input-orange"
              />
            </div>

            {/* Senha */}
            <div>
              <label htmlFor="senha" className="block text-sm font-medium text-gray-300 mb-2">
                🔒 Senha
              </label>
              <div className="relative">
                <input
                  id="senha"
                  type={mostrarSenha ? 'text' : 'password'}
                  value={senha}
                  onChange={(e) => setSenha(e.target.value)}
                  required
                  placeholder="••••••••"
                  className="input-orange pr-12"
                />
                <button
                  type="button"
                  onClick={() => setMostrarSenha(!mostrarSenha)}
                  className="absolute right-3 top-1/2 transform -translate-y-1/2 text-gray-400 hover:text-orange transition-colors"
                >
                  {mostrarSenha ? '🙈' : '👁️'}
                </button>
              </div>
            </div>

            {/* Error Message */}
            {erro && (
              <div className="bg-red-500/20 border border-red-500/50 rounded-lg p-3 fade-in-up">
                <p className="text-sm text-red-300 text-center flex items-center justify-center">
                  <span className="mr-2">⚠️</span>
                  {erro}
                </p>
              </div>
            )}

            {/* Submit Button */}
            <button
              type="submit"
              disabled={carregando}
              className={`w-full btn-orange text-lg py-4 ${
                carregando ? 'opacity-50 cursor-not-allowed' : ''
              }`}
            >
              {carregando ? (
                <div className="flex items-center justify-center">
                  <div className="w-5 h-5 border-2 border-white border-t-transparent rounded-full animate-spin mr-2"></div>
                  Entrando...
                </div>
              ) : (
                <div className="flex items-center justify-center">
                  <span className="mr-2">🚀</span>
                  Entrar
                </div>
              )}
            </button>
          </form>

          {/* Links */}
          <div className="mt-6 text-center space-y-4">
            <p className="text-gray-400">
              Não tem uma conta?{' '}
              <Link 
                to="/registrar" 
                className="text-orange hover:text-orange-light font-semibold transition-colors"
              >
                Registre-se aqui
              </Link>
            </p>
            
            <Link 
              to="/" 
              className="inline-block text-gray-400 hover:text-white transition-colors"
            >
              ← Voltar ao início
            </Link>
          </div>
        </div>

        {/* Features */}
        <div className="mt-8 grid grid-cols-3 gap-4 text-center fade-in-up">
          <div className="text-gray-400">
            <div className="text-2xl mb-1">🎯</div>
            <div className="text-xs">Busca Avançada</div>
          </div>
          <div className="text-gray-400">
            <div className="text-2xl mb-1">❤️</div>
            <div className="text-xs">Lista Favoritos</div>
          </div>
          <div className="text-gray-400">
            <div className="text-2xl mb-1">🎲</div>
            <div className="text-xs">Descoberta</div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default LoginPage;