import React, { useState, useEffect } from 'react';
import type { Filme } from '../types/Filme';
import api from '../services/api';
import FilmeCard from '../components/FilmeCard';

const RecomendacoesPage: React.FC = () => {
  const [recomendacoes, setRecomendacoes] = useState<Filme[]>([]);
  const [carregando, setCarregando] = useState(true);
  const [erro, setErro] = useState<string | null>(null);

  useEffect(() => {
    api.get('/recomendacoes')
      .then(response => {
        setRecomendacoes(response.data);
      })
      .catch(() => {
        setErro('Não foi possível carregar as recomendações. Adicione mais filmes aos seus favoritos!');
      })
      .finally(() => {
        setCarregando(false);
      });
  }, []);

  // A função para adicionar aos favoritos foi atualizada.
  const handleAdicionarFavorito = async (filme: Filme) => {
    try {
      await api.post('/favoritos', {
        filmeId: filme.id,
        titulo: filme.titulo,
        caminhoPoster: filme.caminhoPoster,
      });
      alert(`'${filme.titulo}' foi adicionado aos seus favoritos!`);
    } catch (error: any) {
      // Verifica o código de status do erro que a API retornou.
      if (error.response?.status === 409) {
        alert('Este filme já está na sua lista de favoritos.');
      } else {
        alert('Falha ao adicionar o filme aos favoritos.');
      }
    }
  };

  if (carregando) {
    return <p className="text-center p-10">Analisando seus gostos e buscando recomendações...</p>;
  }

  if (erro) {
    return <p className="text-center text-red-500 p-10">{erro}</p>;
  }

  return (
    <div className="container mx-auto p-4">
      <h1 className="text-3xl font-bold mb-6">Filmes Recomendados para Você</h1>
      {recomendacoes.length > 0 ? (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-6">
          {recomendacoes.map((filme) => (
            <FilmeCard
              key={filme.id}
              filme={filme}
              isLogado={true}
              acao="adicionar"
              onAdicionarFavorito={handleAdicionarFavorito}
            />
          ))}
        </div>
      ) : (
        <p className="text-center text-gray-400">Não encontramos recomendações. Tente adicionar mais filmes aos favoritos.</p>
      )}
    </div>
  );
};

export default RecomendacoesPage;