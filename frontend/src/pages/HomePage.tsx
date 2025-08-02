import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';

const HomePage: React.FC = () => {
  const [termoBusca, setTermoBusca] = useState('');
  const navigate = useNavigate();

  const handleBusca = (e: React.FormEvent) => {
    e.preventDefault();
    if (termoBusca.trim()) {
      navigate(`/buscar?q=${termoBusca}`);
    }
  };

  // Categorias com IDs corretos do TMDB
  const categorias = [
    { nome: 'Ação', emoji: '💥', cor: 'from-red-500 to-orange-500', generoId: '28' },
    { nome: 'Comédia', emoji: '😄', cor: 'from-yellow-500 to-orange-500', generoId: '35' },
    { nome: 'Drama', emoji: '🎭', cor: 'from-purple-500 to-pink-500', generoId: '18' },
    { nome: 'Terror', emoji: '👻', cor: 'from-gray-700 to-black', generoId: '27' },
    { nome: 'Romance', emoji: '❤️', cor: 'from-pink-500 to-red-500', generoId: '10749' },
    { nome: 'Ficção Científica', emoji: '🚀', cor: 'from-blue-500 to-purple-500', generoId: '878' }
  ];

  return (
    <div className="min-h-screen flex flex-col">
      {/* Hero Section */}
      <section className="flex-1 flex flex-col items-center justify-center text-center py-20 px-4">
        <div className="fade-in-up">
          <h1 className="text-6xl md:text-7xl font-extrabold text-white mb-6 leading-tight">
            Explore o <span className="text-orange bg-clip-text text-transparent">Cinema</span>
          </h1>
          <p className="text-xl text-gray-300 mb-12 max-w-3xl leading-relaxed">
            Descubra filmes incríveis, salve seus favoritos e mergulhe no universo cinematográfico. 
            Sua próxima aventura começa aqui! 🎬
          </p>
        </div>

        {/* Search Bar */}
        <form onSubmit={handleBusca} className="w-full max-w-2xl mb-16 fade-in-up">
          <div className="relative glass-effect rounded-2xl p-2">
            <div className="flex items-center">
              <div className="flex-1 relative">
                <input
                  type="text"
                  value={termoBusca}
                  onChange={(e) => setTermoBusca(e.target.value)}
                  placeholder="Buscar por 'Matrix', 'Interestelar', 'Vingadores'..."
                  className="w-full bg-transparent py-4 px-6 text-white text-lg rounded-xl focus:outline-none placeholder-gray-400"
                />
                <div className="absolute left-4 top-1/2 transform -translate-y-1/2 text-orange">
                  🔍
                </div>
              </div>
              <button
                type="submit"
                className="btn-orange ml-2 px-8 py-4 text-lg font-bold pulse-orange"
              >
                Buscar
              </button>
            </div>
          </div>
        </form>

        {/* Quick Categories */}
        <div className="w-full max-w-6xl fade-in-up">
          <h2 className="text-2xl font-bold text-white mb-8 text-center">
            Explore por <span className="text-orange">Categoria</span>
          </h2>
          <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-6 gap-4">
            {categorias.map((categoria, index) => (
              <button
                key={index}
                onClick={() => navigate(`/buscar?generoId=${categoria.generoId}`)}
                className={`bg-gradient-to-br ${categoria.cor} p-6 rounded-xl text-white font-semibold transition-all duration-300 transform hover:scale-105 hover:shadow-lg`}
              >
                <div className="text-3xl mb-2">{categoria.emoji}</div>
                <div className="text-sm">{categoria.nome}</div>
              </button>
            ))}
          </div>
        </div>
      </section>

      {/* Features Section */}
      <section className="py-20 px-4">
        <div className="container mx-auto max-w-6xl">
          <h2 className="text-4xl font-bold text-center text-white mb-16">
            Por que escolher o <span className="text-orange">CineHub</span>?
          </h2>
          <div className="grid md:grid-cols-3 gap-8">
            <div className="text-center fade-in-up">
              <div className="w-16 h-16 bg-orange rounded-full flex items-center justify-center mx-auto mb-4 text-2xl">
                🎯
              </div>
              <h3 className="text-xl font-bold text-white mb-3">Busca Inteligente</h3>
              <p className="text-gray-300">
                Encontre exatamente o que procura com nossa busca avançada e recomendações personalizadas.
              </p>
            </div>
            <div className="text-center fade-in-up">
              <div className="w-16 h-16 bg-orange rounded-full flex items-center justify-center mx-auto mb-4 text-2xl">
                ❤️
              </div>
              <h3 className="text-xl font-bold text-white mb-3">Lista de Favoritos</h3>
              <p className="text-gray-300">
                Salve seus filmes favoritos e acesse-os a qualquer momento. Nunca mais esqueça um filme incrível!
              </p>
            </div>
            <div className="text-center fade-in-up">
              <div className="w-16 h-16 bg-orange rounded-full flex items-center justify-center mx-auto mb-4 text-2xl">
                🎲
              </div>
              <h3 className="text-xl font-bold text-white mb-3">Descoberta Aleatória</h3>
              <p className="text-gray-300">
                Não sabe o que assistir? Use nosso gerador aleatório e descubra filmes surpreendentes!
              </p>
            </div>
          </div>
        </div>
      </section>
    </div>
  );
};

export default HomePage;