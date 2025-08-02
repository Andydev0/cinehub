import React, { useState, useEffect, useCallback } from 'react';
import { useParams } from 'react-router-dom';
import api from '../services/api';
import { useAuth } from '../context/AuthContext';

interface FilmeDetalhes {
  id: number;
  title: string;
  overview: string;
  poster_path: string;
  release_date: string;
  vote_average: number;
  genres: { id: number; name: string }[];
  elenco: { name: string; character: string; profile_path: string | null }[];
  diretor: string;
  escritores: string[];
  trailerKey: string;
}

interface Avaliacao {
  id: number;
  nota: number;
  comentario: string;
  dataCriacao: string;
  nomeUsuario: string;
}

const PaginaDetalhesFilme: React.FC = () => {
  const { id } = useParams<{ id: string }>();
  const { isLogado } = useAuth();

  const [filme, setFilme] = useState<FilmeDetalhes | null>(null);
  const [avaliacoes, setAvaliacoes] = useState<Avaliacao[]>([]);
  const [carregando, setCarregando] = useState(true);
  const [nota, setNota] = useState(0);
  const [comentario, setComentario] = useState('');

  const buscarDados = useCallback(() => {
    setCarregando(true);
    Promise.all([
      api.get(`/filmes/${id}`),
      api.get(`/filmes/${id}/avaliacoes`)
    ]).then(([resFilme, resAvaliacoes]) => {
      setFilme(resFilme.data);
      setAvaliacoes(resAvaliacoes.data);
    }).catch(err => {
      console.error("Falha ao buscar dados", err);
    }).finally(() => {
      setCarregando(false);
    });
  }, [id]);

  useEffect(() => {
    buscarDados();
  }, [buscarDados]);

  const handleEnviarAvaliacao = async (e: React.FormEvent) => {
    e.preventDefault();
    if (nota < 1 || nota > 5) {
      alert("Por favor, selecione uma nota entre 1 e 5.");
      return;
    }
    try {
      await api.post(`/filmes/${id}/avaliacoes`, { nota, comentario });
      buscarDados();
      setNota(0);
      setComentario('');
      alert('Avaliação enviada com sucesso!');
    } catch (error) {
      alert('Falha ao enviar avaliação. Você só pode enviar uma por filme.');
    }
  };

  if (carregando || !filme) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="glass-card p-8 text-center">
          <div className="animate-spin w-12 h-12 border-4 border-orange-500 border-t-transparent rounded-full mx-auto mb-4"></div>
          <p className="text-xl text-orange-400">Carregando detalhes do filme...</p>
        </div>
      </div>
    );
  }

  const placeholderFoto = 'https://via.placeholder.com/185x278.png?text=Sem+Foto';

  return (
    <div className="min-h-screen">
      {/* Hero Section com Background */}
      <div 
        className="relative h-96 bg-cover bg-center bg-no-repeat"
        style={{
          backgroundImage: filme.poster_path 
            ? `linear-gradient(rgba(0, 0, 0, 0.7), rgba(0, 0, 0, 0.8)), url(https://image.tmdb.org/t/p/original${filme.poster_path})`
            : 'linear-gradient(135deg, var(--primary-orange), var(--secondary-orange))'
        }}
      >
        <div className="absolute inset-0 flex items-end">
          <div className="container mx-auto p-6">
            <div className="glass-card p-8 max-w-4xl">
              <div className="flex flex-col md:flex-row gap-6">
                <img 
                  src={`https://image.tmdb.org/t/p/w500${filme.poster_path}`} 
                  alt={`Pôster de ${filme.title}`} 
                  className="w-48 h-72 object-cover rounded-xl shadow-2xl mx-auto md:mx-0"
                />
                <div className="flex-1 text-center md:text-left">
                  <h1 className="text-4xl md:text-5xl font-bold text-white mb-4 drop-shadow-lg">
                    {filme.title}
                  </h1>
                  <div className="text-xl text-orange-200 mb-4">
                    {filme.release_date.substring(0,4)}
                  </div>
                  <div className="flex flex-wrap gap-3 mb-6 justify-center md:justify-start">
                    {filme.genres.map(g => (
                      <span 
                        key={g.id} 
                        className="px-4 py-2 bg-orange-500/20 backdrop-blur-sm border border-orange-400/30 rounded-full text-orange-200 text-sm font-medium"
                      >
                        {g.name}
                      </span>
                    ))}
                  </div>
                  <div className="flex items-center gap-6 justify-center md:justify-start">
                    <div className="flex items-center gap-2">
                      <span className="text-2xl">⭐</span>
                      <span className="text-2xl font-bold text-yellow-400">
                        {filme.vote_average.toFixed(1)}
                      </span>
                      <span className="text-orange-200">TMDB</span>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <div className="container mx-auto p-6 space-y-8">
        {/* Sinopse */}
        <div className="glass-card p-8">
          <h2 className="text-3xl font-bold text-orange-400 mb-6 flex items-center gap-3">
            📖 Sinopse
          </h2>
          <p className="text-lg leading-relaxed text-gray-300">
            {filme.overview || 'Sinopse não disponível.'}
          </p>
          
          {/* Informações Técnicas */}
          <div className="mt-8 grid grid-cols-1 md:grid-cols-2 gap-6">
            {filme.diretor && (
              <div className="flex items-center gap-3">
                <span className="text-2xl">🎬</span>
                <div>
                  <span className="font-bold text-orange-400">Diretor:</span>
                  <span className="ml-2 text-gray-300">{filme.diretor}</span>
                </div>
              </div>
            )}
            {filme.escritores?.length > 0 && (
              <div className="flex items-center gap-3">
                <span className="text-2xl">✍️</span>
                <div>
                  <span className="font-bold text-orange-400">Roteiro:</span>
                  <span className="ml-2 text-gray-300">{filme.escritores.join(', ')}</span>
                </div>
              </div>
            )}
          </div>
        </div>

        {/* Trailer e Elenco */}
        <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
          <div className="lg:col-span-2">
            <div className="glass-card p-8">
              <h2 className="text-3xl font-bold text-orange-400 mb-6 flex items-center gap-3">
                🎥 Trailer
              </h2>
              {filme.trailerKey ? (
                <div className="aspect-w-16 aspect-h-9 rounded-xl overflow-hidden">
                  <iframe
                    src={`https://www.youtube.com/embed/${filme.trailerKey}`}
                    title={`Trailer de ${filme.title}`}
                    allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture"
                    allowFullScreen
                    className="w-full h-64 md:h-80 rounded-xl"
                  ></iframe>
                </div>
              ) : (
                <div className="bg-gray-800/50 rounded-xl p-12 text-center">
                  <span className="text-6xl mb-4 block">🎬</span>
                  <p className="text-gray-400">Trailer não disponível</p>
                </div>
              )}
            </div>
          </div>
          
          <div className="glass-card p-8">
            <h2 className="text-3xl font-bold text-orange-400 mb-6 flex items-center gap-3">
              🎭 Elenco
            </h2>
            <div className="space-y-4 max-h-96 overflow-y-auto">
              {filme.elenco && filme.elenco.length > 0 ? filme.elenco.map((ator, index) => (
                <div key={`${ator.name}-${index}`} className="flex items-center gap-4 bg-gray-800/30 p-3 rounded-lg backdrop-blur-sm border border-gray-700/50">
                  <img
                    src={ator.profile_path ? `https://image.tmdb.org/t/p/w185${ator.profile_path}` : placeholderFoto}
                    alt={`Foto de ${ator.name}`}
                    className="w-16 h-20 object-cover rounded-lg"
                  />
                  <div className="text-sm">
                    <p className="font-semibold text-white">{ator.name}</p>
                    <p className="text-orange-300">como {ator.character}</p>
                  </div>
                </div>
              )) : (
                <p className="text-gray-400 text-center py-8">Elenco não disponível</p>
              )}
            </div>
          </div>
        </div>
        
        {/* Seção de Avaliações */}
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-8">
          {isLogado && (
            <div className="glass-card p-8">
              <h2 className="text-3xl font-bold text-orange-400 mb-6 flex items-center gap-3">
                ⭐ Sua Avaliação
              </h2>
              <form onSubmit={handleEnviarAvaliacao} className="space-y-6">
                <div>
                  <label className="block text-orange-300 font-medium mb-2">Sua Nota</label>
                  <div className="flex items-center gap-2 mb-4">
                    <div className="flex">
                      {[1, 2, 3, 4, 5].map((estrela) => (
                        <button
                          key={estrela}
                          type="button"
                          onClick={() => setNota(estrela)}
                          onMouseEnter={() => {
                            const starsPreview = document.querySelectorAll('.star-preview');
                            starsPreview.forEach((star, index) => {
                              if (index < estrela) {
                                star.classList.add('text-yellow-400', 'scale-110');
                                star.classList.remove('text-gray-500');
                              } else {
                                star.classList.add('text-gray-500');
                                star.classList.remove('text-yellow-400', 'scale-110');
                              }
                            });
                          }}
                          onMouseLeave={() => {
                            const starsPreview = document.querySelectorAll('.star-preview');
                            starsPreview.forEach((star, index) => {
                              if (index < nota) {
                                star.classList.add('text-yellow-400');
                                star.classList.remove('text-gray-500', 'scale-110');
                              } else {
                                star.classList.add('text-gray-500');
                                star.classList.remove('text-yellow-400', 'scale-110');
                              }
                            });
                          }}
                          className="text-2xl focus:outline-none transition-transform duration-200"
                          aria-label={`Avaliar com ${estrela} estrelas`}
                        >
                          {estrela <= nota ? (
                            <span className="star-preview text-yellow-400 transition-all duration-200">⭐</span>
                          ) : (
                            <span className="star-preview text-gray-500 transition-all duration-200">☆</span>
                          )}
                        </button>
                      ))}
                    </div>
                    {nota > 0 && (
                      <span className="text-yellow-400 ml-2">{nota}/5</span>
                    )}
                  </div>
                </div>
                <div>
                  <label className="block text-orange-300 font-medium mb-2">Seu Comentário (opcional)</label>
                  <textarea 
                    value={comentario} 
                    onChange={(e) => setComentario(e.target.value)} 
                    rows={4} 
                    className="input-orange resize-none"
                    placeholder="Compartilhe sua opinião sobre o filme..."
                  ></textarea>
                </div>
                <button type="submit" className="btn-orange w-full">
                  Enviar Avaliação
                </button>
              </form>
            </div>
          )}
          
          <div className="glass-card p-8">
            <h2 className="text-3xl font-bold text-orange-400 mb-6 flex items-center gap-3">
              💬 Avaliações
            </h2>
            {avaliacoes.length > 0 ? (
              <div className="space-y-4 max-h-96 overflow-y-auto">
                {avaliacoes.map(a => (
                  <div key={a.id} className="bg-gray-800/30 p-4 rounded-lg backdrop-blur-sm border border-gray-700/50">
                    <div className="flex justify-between items-center mb-2">
                      <p className="font-bold text-white">{a.nomeUsuario}</p>
                      <div className="flex items-center gap-1">
                        {[1, 2, 3, 4, 5].map((estrela) => (
                          <span key={estrela} className="text-xl">
                            {estrela <= a.nota ? (
                              <span className="text-yellow-400">⭐</span>
                            ) : (
                              <span className="text-gray-500">☆</span>
                            )}
                          </span>
                        ))}
                      </div>
                    </div>
                    {a.comentario && (
                      <p className="text-gray-300 mb-2">{a.comentario}</p>
                    )}
                    <div className="flex justify-between items-center mt-2 text-sm text-gray-400">
                      <p>{new Date(a.dataCriacao).toLocaleDateString('pt-BR', {
                        day: '2-digit',
                        month: '2-digit',
                        year: 'numeric',
                        hour: '2-digit',
                        minute: '2-digit'
                      })}</p>
                    </div>
                  </div>
                ))}
              </div>
            ) : (
              <div className="text-center py-8">
                <span className="text-6xl mb-4 block">💬</span>
                <p className="text-gray-400">Este filme ainda não tem avaliações.</p>
                <p className="text-orange-400">Seja o primeiro a avaliar!</p>
              </div>
            )}
          </div>
        </div>
      </div>
    </div>
  );
};

export default PaginaDetalhesFilme;