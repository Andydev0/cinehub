import React from 'react';
import type { Filme } from '../types/Filme';
import { Link } from 'react-router-dom';

interface Props {
  filme: Filme;
  isLogado: boolean;
  acao?: 'adicionar' | 'remover';
  onAdicionarFavorito?: (filme: Filme) => void;
  onRemoverFavorito?: (filmeId: number) => void;
}

const FilmeCard: React.FC<Props> = ({ filme, isLogado, acao, onAdicionarFavorito, onRemoverFavorito }) => {
  const handleButtonClick = (e: React.MouseEvent, action: () => void) => {
    e.preventDefault();
    e.stopPropagation();
    action();
  };

  const renderStars = (rating: number) => {
    const stars = [];
    const fullStars = Math.floor(rating / 2);
    const hasHalfStar = (rating / 2) % 1 >= 0.5;
    
    for (let i = 0; i < 5; i++) {
      if (i < fullStars) {
        stars.push(<span key={i} className="text-orange">★</span>);
      } else if (i === fullStars && hasHalfStar) {
        stars.push(<span key={i} className="text-orange">☆</span>);
      } else {
        stars.push(<span key={i} className="text-gray-500">☆</span>);
      }
    }
    return stars;
  };

  return (
    <Link to={`/filme/${filme.id}`} className="block h-full group">
      <div className="movie-card flex flex-col h-full fade-in-up">
        {/* Poster com overlay de informações */}
        <div className="relative overflow-hidden">
          <img
            src={filme.caminhoPoster}
            alt={`Pôster de ${filme.titulo}`}
            className="w-full h-96 object-cover transition-transform duration-300 group-hover:scale-110"
          />
          
          {/* Overlay com nota */}
          {filme.notaMedia > 0 && (
            <div className="absolute top-3 right-3 glass-effect rounded-lg px-3 py-1">
              <div className="flex items-center space-x-1">
                {renderStars(filme.notaMedia)}
                <span className="text-white text-sm font-bold ml-1">
                  {filme.notaMedia.toFixed(1)}
                </span>
              </div>
            </div>
          )}
          
          {/* Gradient overlay no bottom */}
          <div className="absolute bottom-0 left-0 right-0 h-20 bg-gradient-to-t from-black/80 to-transparent" />
        </div>

        {/* Conteúdo do card */}
        <div className="p-5 flex flex-col flex-grow">
          <h3 className="text-xl font-bold mb-2 text-white group-hover:text-orange transition-colors duration-300">
            {filme.titulo}
          </h3>
          
          {filme.dataLancamento && (
            <p className="text-gray-400 text-sm mb-3">
              📅 {filme.dataLancamento}
            </p>
          )}
          
          {filme.sinopse && (
            <p className="text-gray-300 text-sm flex-grow mb-4 line-clamp-3 leading-relaxed">
              {filme.sinopse}
            </p>
          )}

          {/* Botões de ação */}
          <div className="mt-auto">
            {isLogado && acao === 'adicionar' && onAdicionarFavorito && (
              <button
                onClick={(e) => handleButtonClick(e, () => onAdicionarFavorito(filme))}
                className="w-full btn-orange text-sm flex items-center justify-center space-x-2"
              >
                <span>❤️</span>
                <span>Adicionar aos Favoritos</span>
              </button>
            )}
            
            {isLogado && acao === 'remover' && onRemoverFavorito && (
              <button
                onClick={(e) => handleButtonClick(e, () => onRemoverFavorito(filme.id))}
                className="w-full px-4 py-3 text-sm font-semibold text-white bg-red-600 rounded-lg hover:bg-red-700 transition-all duration-300 transform hover:scale-105 flex items-center justify-center space-x-2"
              >
                <span>🗑️</span>
                <span>Remover dos Favoritos</span>
              </button>
            )}
          </div>
        </div>
      </div>
    </Link>
  );
};

export default FilmeCard;