import React, { useState, useEffect } from 'react';
import { useSearchParams, useNavigate } from 'react-router-dom';
import type { Filme } from '../types/Filme';
import api from '../services/api';
import FilmeCard from '../components/FilmeCard';
import { useAuth } from '../context/AuthContext';

const ResultadosBuscaPage: React.FC = () => {
  const [searchParams] = useSearchParams();
  const navigate = useNavigate();
  const termoBusca = searchParams.get('q') || '';
  const generoId = searchParams.get('generoId') || '';
  const [novoTermo, setNovoTermo] = useState(termoBusca);

  const [filmes, setFilmes] = useState<Filme[]>([]);
  const [carregando, setCarregando] = useState(false);
  const [erro, setErro] = useState<string | null>(null);
  const [tipoConsulta, setTipoConsulta] = useState<'busca' | 'genero'>('busca');
  
  const { isLogado } = useAuth();

  // Mapeamento de IDs de gênero para nomes
  const generoNomes: { [key: string]: string } = {
    '28': 'Ação',
    '35': 'Comédia',
    '18': 'Drama',
    '27': 'Terror',
    '10749': 'Romance',
    '878': 'Ficção Científica'
  };

  const nomeGenero = generoId ? generoNomes[generoId] || 'Gênero Desconhecido' : '';

  useEffect(() => {
    if (termoBusca || generoId) {
      setCarregando(true);
      setErro(null);
      setFilmes([]);

      let endpoint = '';
      if (generoId) {
        endpoint = `/filmes/genero?generoId=${generoId}`;
        setTipoConsulta('genero');
      } else {
        endpoint = `/filmes/buscar?termo=${termoBusca}`;
        setTipoConsulta('busca');
      }

      api.get(endpoint)
        .then(response => setFilmes(response.data))
        .catch(() => setErro('Falha ao buscar filmes.'))
        .finally(() => setCarregando(false));
    }
  }, [termoBusca, generoId]);

  const handleAdicionarFavorito = async (filme: Filme) => {
    try {
      await api.post('/favoritos', {
        filmeId: filme.id,
        titulo: filme.titulo,
        caminhoPoster: filme.caminhoPoster,
      });
    } catch (error: any) {
      if (error.response?.status === 409) {
        console.log('Filme já está nos favoritos');
      } else {
        console.error('Falha ao adicionar filme aos favoritos');
      }
    }
  };

  const handleNovaBusca = (e: React.FormEvent) => {
    e.preventDefault();
    if (novoTermo.trim()) {
      navigate(`/buscar?q=${novoTermo}`);
    }
  };

  if (carregando) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="text-center">
          <div className="w-16 h-16 border-4 border-orange border-t-transparent rounded-full animate-spin mx-auto mb-4"></div>
          <p className="text-xl text-gray-300">Buscando filmes...</p>
        </div>
      </div>
    );
  }

  if (erro) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="text-center">
          <div className="text-6xl mb-4">😔</div>
          <p className="text-xl text-red-400 mb-4">{erro}</p>
          <button 
            onClick={() => window.location.reload()}
            className="btn-orange"
          >
            Tentar Novamente
          </button>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen py-8">
      <div className="container mx-auto px-4">
        {/* Header com busca */}
        <div className="text-center mb-12 fade-in-up">
          <h1 className="text-4xl md:text-5xl font-bold text-white mb-4">
            {tipoConsulta === 'genero' ? (
              <>Filmes de <span className="text-orange">{nomeGenero}</span></>
            ) : (
              <>Resultados para: <span className="text-orange">"{termoBusca}"</span></>
            )}
          </h1>
          
          {/* Nova busca */}
          <form onSubmit={handleNovaBusca} className="max-w-2xl mx-auto mt-8">
            <div className="glass-effect rounded-2xl p-2">
              <div className="flex items-center">
                <input
                  type="text"
                  value={novoTermo}
                  onChange={(e) => setNovoTermo(e.target.value)}
                  placeholder="Refinar busca ou tentar novo termo..."
                  className="flex-1 bg-transparent py-3 px-6 text-white rounded-xl focus:outline-none placeholder-gray-400"
                />
                <button
                  type="submit"
                  className="btn-orange ml-2 px-6 py-3"
                >
                  🔍 Buscar
                </button>
              </div>
            </div>
          </form>
        </div>

        {filmes.length > 0 ? (
          <>
            {/* Stats */}
            <div className="glass-effect rounded-2xl p-6 mb-8 fade-in-up">
              <div className="flex items-center justify-center space-x-8">
                <div className="text-center">
                  <div className="text-3xl font-bold text-orange">{filmes.length}</div>
                  <div className="text-gray-300">Filmes Encontrados</div>
                </div>
                <div className="text-center">
                  <div className="text-3xl font-bold text-orange">🎬</div>
                  <div className="text-gray-300">Resultados</div>
                </div>
              </div>
            </div>

            {/* Grid de filmes */}
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-8">
              {filmes.map((filme, index) => (
                <div key={filme.id} className="fade-in-up" style={{ animationDelay: `${index * 0.1}s` }}>
                  <FilmeCard
                    filme={filme}
                    isLogado={isLogado}
                    acao="adicionar"
                    onAdicionarFavorito={handleAdicionarFavorito}
                  />
                </div>
              ))}
            </div>
          </>
        ) : (
          <div className="text-center py-20 fade-in-up">
            <div className="text-8xl mb-6">🔍</div>
            <h2 className="text-3xl font-bold text-white mb-4">
              Nenhum resultado encontrado
            </h2>
            <p className="text-xl text-gray-300 mb-8 max-w-md mx-auto">
              Não encontramos filmes para "{termoBusca}". Que tal tentar uma busca diferente?
            </p>
            <div className="flex flex-col sm:flex-row gap-4 justify-center">
              <button 
                onClick={() => navigate('/')}
                className="btn-orange-outline"
              >
                Voltar ao Início
              </button>
              <button 
                onClick={() => navigate('/aleatorio')}
                className="btn-orange"
              >
                Filme Aleatório
              </button>
            </div>
          </div>
        )}
      </div>
    </div>
  );
};

export default ResultadosBuscaPage;