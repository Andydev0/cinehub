import React, { useState, useEffect } from 'react';
import type { Filme } from '../types/Filme';
import api from '../services/api';
import FilmeCard from '../components/FilmeCard';

interface FilmeFavoritoResponse {
  id: number;
  filmeId: number;
  titulo: string;
  caminhoPoster: string;
}

const FavoritosPage: React.FC = () => {
  const [favoritos, setFavoritos] = useState<Filme[]>([]);
  const [carregando, setCarregando] = useState(true);
  const [erro, setErro] = useState<string | null>(null);

  const buscarFavoritos = () => {
    setCarregando(true);
    api.get('/favoritos')
      .then(response => {
        const filmesMapeados = response.data.map((fav: FilmeFavoritoResponse) => ({
          id: fav.filmeId,
          titulo: fav.titulo,
          caminhoPoster: fav.caminhoPoster,
          sinopse: '',
          dataLancamento: '',
          notaMedia: 0,
        }));
        setFavoritos(filmesMapeados);
      })
      .catch(() => {
        setErro('Falha ao carregar seus favoritos. Tente fazer login novamente.');
      })
      .finally(() => {
        setCarregando(false);
      });
  };

  useEffect(() => {
    buscarFavoritos();
  }, []);

  const handleRemoverFavorito = async (filmeId: number) => {
    try {
      await api.delete(`/favoritos/${filmeId}`);
      setFavoritos(favoritos.filter(filme => filme.id !== filmeId));
    } catch (error) {
      console.error('Falha ao remover o filme:', error);
    }
  };

  if (carregando) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="text-center">
          <div className="w-16 h-16 border-4 border-orange border-t-transparent rounded-full animate-spin mx-auto mb-4"></div>
          <p className="text-xl text-gray-300">Carregando seus filmes favoritos...</p>
        </div>
      </div>
    );
  }

  if (erro) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="text-center">
          <div className="text-6xl mb-4">😢</div>
          <p className="text-xl text-red-400 mb-4">{erro}</p>
          <button 
            onClick={buscarFavoritos}
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
        {/* Header */}
        <div className="text-center mb-12 fade-in-up">
          <h1 className="text-5xl font-bold text-white mb-4">
            Meus <span className="text-orange">Favoritos</span> ❤️
          </h1>
          <p className="text-xl text-gray-300 max-w-2xl mx-auto">
            Aqui estão todos os filmes que você salvou. Sua coleção pessoal de obras-primas!
          </p>
        </div>

        {favoritos.length > 0 ? (
          <>
            {/* Stats */}
            <div className="glass-effect rounded-2xl p-6 mb-8 fade-in-up">
              <div className="flex items-center justify-center space-x-8">
                <div className="text-center">
                  <div className="text-3xl font-bold text-orange">{favoritos.length}</div>
                  <div className="text-gray-300">Filmes Salvos</div>
                </div>
                <div className="text-center">
                  <div className="text-3xl font-bold text-orange">🏆</div>
                  <div className="text-gray-300">Sua Coleção</div>
                </div>
              </div>
            </div>

            {/* Grid de filmes */}
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-8">
              {favoritos.map((filme, index) => (
                <div key={filme.id} className="fade-in-up" style={{ animationDelay: `${index * 0.1}s` }}>
                  <FilmeCard
                    filme={filme}
                    isLogado={true}
                    acao="remover"
                    onRemoverFavorito={handleRemoverFavorito}
                  />
                </div>
              ))}
            </div>
          </>
        ) : (
          <div className="text-center py-20 fade-in-up">
            <div className="text-8xl mb-6">📺</div>
            <h2 className="text-3xl font-bold text-white mb-4">
              Nenhum filme favoritado ainda
            </h2>
            <p className="text-xl text-gray-300 mb-8 max-w-md mx-auto">
              Que tal começar explorando alguns filmes incríveis? Use a busca para descobrir seus próximos favoritos!
            </p>
            <button 
              onClick={() => window.location.href = '/'}
              className="btn-orange text-lg"
            >
              Descobrir Filmes
            </button>
          </div>
        )}
      </div>
    </div>
  );
};

export default FavoritosPage;