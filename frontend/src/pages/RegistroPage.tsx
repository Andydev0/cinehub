import React, { useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import api from '../services/api';

const RegistroPage: React.FC = () => {
  const [nome, setNome] = useState('');
  const [email, setEmail] = useState('');
  const [senha, setSenha] = useState('');
  const [confirmSenha, setConfirmSenha] = useState('');
  const [erro, setErro] = useState<string | null>(null);
  const [carregando, setCarregando] = useState(false);
  const [mostrarSenha, setMostrarSenha] = useState(false);
  const [sucesso, setSucesso] = useState(false);
  const navigate = useNavigate();

  const validarSenha = (senha: string) => {
    return senha.length >= 6;
  };

  const handleRegistro = async (e: React.FormEvent) => {
    e.preventDefault();
    setCarregando(true);
    setErro(null);

    // Validações
    if (!validarSenha(senha)) {
      setErro('A senha deve ter pelo menos 6 caracteres.');
      setCarregando(false);
      return;
    }

    if (senha !== confirmSenha) {
      setErro('As senhas não coincidem.');
      setCarregando(false);
      return;
    }

    try {
      await api.post('/auth/registrar', { nome, email, senha });
      setSucesso(true);
      setTimeout(() => navigate('/login'), 2000);
    } catch (error: any) {
      const msgErro = error.response?.data?.erro || 'Falha ao registrar. Tente novamente.';
      setErro(msgErro);
    } finally {
      setCarregando(false);
    }
  };

  if (sucesso) {
    return (
      <div className="min-h-screen flex items-center justify-center py-12 px-4">
        <div className="w-full max-w-md text-center fade-in-up">
          <div className="glass-effect rounded-2xl p-8">
            <div className="text-6xl mb-4">🎉</div>
            <h1 className="text-3xl font-bold text-white mb-4">
              Conta criada com <span className="text-orange">sucesso</span>!
            </h1>
            <p className="text-gray-300 mb-6">
              Bem-vindo ao CineHub, {nome}! Redirecionando para o login...
            </p>
            <div className="w-8 h-8 border-4 border-orange border-t-transparent rounded-full animate-spin mx-auto"></div>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen flex items-center justify-center py-12 px-4">
      <div className="w-full max-w-md">
        {/* Logo/Header */}
        <div className="text-center mb-8 fade-in-up">
          <div className="text-6xl mb-4">🎆</div>
          <h1 className="text-4xl font-bold text-white mb-2">
            Junte-se ao <span className="text-orange">CineHub</span>!
          </h1>
          <p className="text-gray-300">
            Crie sua conta e comece a descobrir filmes incríveis
          </p>
        </div>

        {/* Form */}
        <div className="glass-effect rounded-2xl p-8 fade-in-up">
          <form onSubmit={handleRegistro} className="space-y-6">
            {/* Nome */}
            <div>
              <label htmlFor="nome" className="block text-sm font-medium text-gray-300 mb-2">
                👤 Nome Completo
              </label>
              <input
                id="nome"
                type="text"
                value={nome}
                onChange={(e) => setNome(e.target.value)}
                required
                placeholder="Seu nome completo"
                className="input-orange"
              />
            </div>

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
                  placeholder="Mínimo 6 caracteres"
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
              {senha && (
                <div className="mt-2">
                  <div className={`text-xs ${
                    validarSenha(senha) ? 'text-green-400' : 'text-red-400'
                  }`}>
                    {validarSenha(senha) ? '✓ Senha válida' : '⚠️ Mínimo 6 caracteres'}
                  </div>
                </div>
              )}
            </div>

            {/* Confirmar Senha */}
            <div>
              <label htmlFor="confirmSenha" className="block text-sm font-medium text-gray-300 mb-2">
                🔁 Confirmar Senha
              </label>
              <input
                id="confirmSenha"
                type="password"
                value={confirmSenha}
                onChange={(e) => setConfirmSenha(e.target.value)}
                required
                placeholder="Digite a senha novamente"
                className="input-orange"
              />
              {confirmSenha && (
                <div className="mt-2">
                  <div className={`text-xs ${
                    senha === confirmSenha ? 'text-green-400' : 'text-red-400'
                  }`}>
                    {senha === confirmSenha ? '✓ Senhas coincidem' : '⚠️ Senhas não coincidem'}
                  </div>
                </div>
              )}
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
                  Criando conta...
                </div>
              ) : (
                <div className="flex items-center justify-center">
                  <span className="mr-2">✨</span>
                  Criar Conta
                </div>
              )}
            </button>
          </form>

          {/* Links */}
          <div className="mt-6 text-center space-y-4">
            <p className="text-gray-400">
              Já tem uma conta?{' '}
              <Link 
                to="/login" 
                className="text-orange hover:text-orange-light font-semibold transition-colors"
              >
                Faça login
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

        {/* Benefits */}
        <div className="mt-8 grid grid-cols-3 gap-4 text-center fade-in-up">
          <div className="text-gray-400">
            <div className="text-2xl mb-1">🎯</div>
            <div className="text-xs">Busca Inteligente</div>
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

export default RegistroPage;