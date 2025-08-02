import React, { useState, useEffect } from 'react';
import type { Filme } from '../types/Filme';
import api from '../services/api';
import FilmeCard from '../components/FilmeCard';
import { useAuth } from '../context/AuthContext';

interface Genero {
  id: number;
  name: string;
}

const GeradorAleatorioPage: React.FC = () => {
  const [generos, setGeneros] = useState<Genero[]>([]);
  const [generoSelecionado, setGeneroSelecionado] = useState('');
  const [ano, setAno] = useState('');
  const [filmeSugerido, setFilmeSugerido] = useState<Filme | null>(null);
  const [carregando, setCarregando] = useState(false);
  const [erro, setErro] = useState<string | null>(null);
  const [historico, setHistorico] = useState<Filme[]>([]);
  const { isLogado } = useAuth();

  useEffect(() => {
    api.get('/generos')
      .then(response => setGeneros(response.data))
      .catch(() => setErro('Falha ao carregar lista de gêneros.'));
  }, []);

  const handleSugerirFilme = async (e: React.FormEvent) => {
    e.preventDefault();
    setCarregando(true);
    setErro(null);

    try {
      const response = await api.get('/filmes/aleatorio', {
        params: {
          generoId: generoSelecionado,
          ano: ano,
        },
      });
      
      const novoFilme = response.data;
      setFilmeSugerido(novoFilme);
      
      // Adiciona ao histórico
      if (novoFilme) {
        setHistorico(prev => [novoFilme, ...prev.slice(0, 4)]);
      }
    } catch (error) {
      setErro('Nenhum filme encontrado com esses filtros. Tente novamente!');
    } finally {
      setCarregando(false);
    }
  };
  
  const handleAdicionarFavorito = async (filme: Filme) => {
    try {
      await api.post('/favoritos', { 
        filmeId: filme.id, 
        titulo: filme.titulo, 
        caminhoPoster: filme.caminhoPoster 
      });
    } catch (error: any) {
      console.error('Erro ao adicionar favorito:', error);
    }
  };

  const limparFiltros = () => {
    setGeneroSelecionado('');
    setAno('');
    setFilmeSugerido(null);
    setErro(null);
  };

  return (
    <div className="min-h-screen py-8">
      <div className="container mx-auto px-4">
        {/* Header */}
        <div className="text-center mb-12 fade-in-up">
          <div className="text-6xl mb-4">🎲</div>
          <h1 className="text-5xl font-bold text-white mb-4">
            Gerador <span className="text-orange">Aleatório</span>
          </h1>
          <p className="text-xl text-gray-300 max-w-3xl mx-auto leading-relaxed">
            Não sabe o que assistir? Configure seus filtros e deixe a sorte escolher 
            o filme perfeito para você! 🎬✨
          </p>
        </div>

        {/* Filtros */}
        <div className="max-w-4xl mx-auto mb-12 fade-in-up">
          <div className="glass-effect rounded-2xl p-8">
            <form onSubmit={handleSugerirFilme} className="space-y-6">
              <div className="grid md:grid-cols-2 gap-6">
                {/* Gênero */}
                <div>
                  <label className="block text-sm font-medium text-gray-300 mb-3">
                    🎭 Gênero
                  </label>
                  <select
                    value={generoSelecionado}
                    onChange={(e) => setGeneroSelecionado(e.target.value)}
                    className="input-orange"
                  >
                    <option value="">Qualquer Gênero</option>
                    {generos.map(g => (
                      <option key={g.id} value={g.id}>{g.name}</option>
                    ))}
                  </select>
                </div>

                {/* Ano */}
                <div>
                  <label className="block text-sm font-medium text-gray-300 mb-3">
                    📅 Ano
                  </label>
                  <input
                    type="number"
                    value={ano}
                    onChange={(e) => setAno(e.target.value)}
                    placeholder="Ex: 2020, 1995..."
                    min="1900"
                    max={new Date().getFullYear()}
                    className="input-orange"
                  />
                </div>
              </div>

              {/* Botões */}
              <div className="flex flex-col sm:flex-row gap-4 justify-center">
                <button
                  type="submit"
                  disabled={carregando}
                  className={`btn-orange text-lg px-8 py-4 pulse-orange ${
                    carregando ? 'opacity-50 cursor-not-allowed' : ''
                  }`}
                >
                  {carregando ? (
                    <div className="flex items-center justify-center">
                      <div className="w-5 h-5 border-2 border-white border-t-transparent rounded-full animate-spin mr-2"></div>
                      Sorteando...
                    </div>
                  ) : (
                    <div className="flex items-center justify-center">
                      <span className="mr-2">🎲</span>
                      Sortear Filme!
                    </div>
                  )}
                </button>
                
                <button
                  type="button"
                  onClick={limparFiltros}
                  className="btn-orange-outline px-8 py-4"
                >
                  🧩 Limpar Filtros
                </button>
              </div>
            </form>
          </div>
        </div>

        {/* Erro */}
        {erro && (
          <div className="max-w-2xl mx-auto mb-8 fade-in-up">
            <div className="bg-red-500/20 border border-red-500/50 rounded-2xl p-6 text-center">
              <div className="text-4xl mb-3">😔</div>
              <p className="text-red-300 text-lg">{erro}</p>
            </div>
          </div>
        )}

        {/* Filme Sugerido */}
        {filmeSugerido && (
          <div className="max-w-2xl mx-auto mb-12 fade-in-up">
            <div className="text-center mb-6">
              <h2 className="text-3xl font-bold text-white mb-2">
                🎆 Sua <span className="text-orange">Sugestão</span>!
              </h2>
              <p className="text-gray-300">Que tal assistir este filme hoje?</p>
            </div>
            <FilmeCard
              filme={filmeSugerido}
              isLogado={isLogado}
              acao="adicionar"
              onAdicionarFavorito={handleAdicionarFavorito}
            />
          </div>
        )}

        {/* Histórico */}
        {historico.length > 0 && (
          <div className="max-w-6xl mx-auto fade-in-up">
            <h2 className="text-3xl font-bold text-center text-white mb-8">
              📜 Histórico de <span className="text-orange">Sugestões</span>
            </h2>
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-5 gap-6">
              {historico.map((filme, index) => (
                <div key={`${filme.id}-${index}`} className="fade-in-up" style={{ animationDelay: `${index * 0.1}s` }}>
                  <FilmeCard
                    filme={filme}
                    isLogado={isLogado}
                    acao="adicionar"
                    onAdicionarFavorito={handleAdicionarFavorito}
                  />
                </div>
              ))}
            </div>
          </div>
        )}

        {/* Dicas */}
        {!filmeSugerido && !carregando && (
          <div className="max-w-4xl mx-auto text-center fade-in-up">
            <div className="glass-effect rounded-2xl p-8">
              <h3 className="text-2xl font-bold text-white mb-6">
                💡 Dicas para uma <span className="text-orange">boa escolha</span>
              </h3>
              <div className="grid md:grid-cols-3 gap-6">
                <div className="text-center">
                  <div className="text-3xl mb-3">🎭</div>
                  <h4 className="font-semibold text-white mb-2">Experimente Gêneros</h4>
                  <p className="text-gray-300 text-sm">
                    Deixe em branco para descobrir gêneros que você nunca assistiu!
                  </p>
                </div>
                <div className="text-center">
                  <div className="text-3xl mb-3">📅</div>
                  <h4 className="font-semibold text-white mb-2">Explore Épocas</h4>
                  <p className="text-gray-300 text-sm">
                    Filmes clássicos dos anos 80-90 ou lançamentos recentes?
                  </p>
                </div>
                <div className="text-center">
                  <div className="text-3xl mb-3">🎲</div>
                  <h4 className="font-semibold text-white mb-2">Confie na Sorte</h4>
                  <p className="text-gray-300 text-sm">
                    Às vezes os melhores filmes são os que encontramos por acaso!
                  </p>
                </div>
              </div>
            </div>
          </div>
        )}
      </div>
    </div>
  );
};

export default GeradorAleatorioPage;