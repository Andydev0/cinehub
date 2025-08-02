import React, { useState, useEffect } from 'react';
import api from '../services/api';

interface Opcao {
  id: number;
  texto: string;
}

interface Pergunta {
  pergunta: string;
  opcoes: Opcao[];
  respostaCorretaId: number;
}

const QuizPage: React.FC = () => {
  const [pergunta, setPergunta] = useState<Pergunta | null>(null);
  const [respostaSelecionada, setRespostaSelecionada] = useState<number | null>(null);
  const [isRespostaCorreta, setIsRespostaCorreta] = useState<boolean | null>(null);
  const [pontuacao, setPontuacao] = useState(0);
  const [carregando, setCarregando] = useState(true);
  const [perguntasRespondidas, setPerguntasRespondidas] = useState(0);
  const [acertos, setAcertos] = useState(0);

  const buscarPergunta = () => {
    setCarregando(true);
    setRespostaSelecionada(null);
    setIsRespostaCorreta(null);
    api.get('/quiz/pergunta')
      .then(response => {
        setPergunta(response.data);
      })
      .catch(() => {
        // Em caso de erro, mantém a pergunta atual ou mostra mensagem
      })
      .finally(() => setCarregando(false));
  };

  useEffect(() => {
    buscarPergunta();
  }, []);

  const handleResponder = (opcaoId: number) => {
    if (respostaSelecionada !== null) return;
    
    setRespostaSelecionada(opcaoId);
    setPerguntasRespondidas(prev => prev + 1);
    
    if (opcaoId === pergunta?.respostaCorretaId) {
      setIsRespostaCorreta(true);
      setPontuacao(pontuacao + 10);
      setAcertos(prev => prev + 1);
    } else {
      setIsRespostaCorreta(false);
    }
  };

  const getEstiloBotao = (opcaoId: number) => {
    if (respostaSelecionada === null) {
      return 'glass-effect hover:bg-white/20 border border-white/20 text-white hover:border-orange/50 transition-all duration-300';
    }
    if (opcaoId === pergunta?.respostaCorretaId) {
      return 'bg-green-500/80 border-2 border-green-400 text-white shadow-lg';
    }
    if (opcaoId === respostaSelecionada) {
      return 'bg-red-500/80 border-2 border-red-400 text-white shadow-lg';
    }
    return 'glass-effect opacity-50 border border-white/10 text-gray-400';
  };

  // Loading state
  if (carregando) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="text-center fade-in-up">
          <div className="w-16 h-16 border-4 border-orange border-t-transparent rounded-full animate-spin mx-auto mb-4"></div>
          <p className="text-xl text-gray-300"> Preparando sua pergunta...</p>
        </div>
      </div>
    );
  }

  const porcentagemAcerto = perguntasRespondidas > 0 ? Math.round((acertos / perguntasRespondidas) * 100) : 0;

  return (
    <div className="min-h-screen py-8">
      <div className="container mx-auto px-4">
        {/* Header */}
        <div className="text-center mb-12 fade-in-up">
          <h1 className="text-5xl md:text-6xl font-bold text-white mb-4">
            Quiz de <span className="text-orange">Cinema</span>
          </h1>
          <p className="text-xl text-gray-300 mb-8">
            Teste seus conhecimentos sobre filmes e ganhe pontos!
          </p>
        </div>

        {/* Stats */}
        <div className="max-w-4xl mx-auto mb-8 fade-in-up">
          <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
            <div className="glass-effect rounded-2xl p-6 text-center">
              <div className="text-3xl font-bold text-orange mb-2">{pontuacao}</div>
              <div className="text-gray-300 text-sm"> Pontos</div>
            </div>
            <div className="glass-effect rounded-2xl p-6 text-center">
              <div className="text-3xl font-bold text-blue-400 mb-2">{perguntasRespondidas}</div>
              <div className="text-gray-300 text-sm"> Perguntas</div>
            </div>
            <div className="glass-effect rounded-2xl p-6 text-center">
              <div className="text-3xl font-bold text-green-400 mb-2">{acertos}</div>
              <div className="text-gray-300 text-sm"> Acertos</div>
            </div>
            <div className="glass-effect rounded-2xl p-6 text-center">
              <div className="text-3xl font-bold text-purple-400 mb-2">{porcentagemAcerto}%</div>
              <div className="text-gray-300 text-sm"> Precisão</div>
            </div>
          </div>
        </div>

        {/* Quiz Card */}
        {pergunta && (
          <div className="max-w-4xl mx-auto fade-in-up">
            <div className="glass-effect rounded-3xl p-8 md:p-12">
              {/* Pergunta */}
              <div className="text-center mb-8">
                <h2 className="text-2xl md:text-3xl font-bold text-white mb-4 leading-relaxed">
                  {pergunta.pergunta}
                </h2>
                <div className="w-24 h-1 bg-gradient-to-r from-orange to-yellow-400 mx-auto rounded-full"></div>
              </div>

              {/* Opções */}
              <div className="grid grid-cols-1 md:grid-cols-2 gap-4 mb-8">
                {pergunta.opcoes.map((opcao, index) => (
                  <button
                    key={opcao.id}
                    onClick={() => handleResponder(opcao.id)}
                    disabled={respostaSelecionada !== null}
                    className={`p-6 rounded-2xl font-semibold text-lg transition-all duration-300 transform hover:scale-105 disabled:hover:scale-100 ${getEstiloBotao(opcao.id)}`}
                  >
                    <div className="flex items-center space-x-3">
                      <div className="w-8 h-8 rounded-full bg-white/20 flex items-center justify-center text-sm font-bold">
                        {String.fromCharCode(65 + index)}
                      </div>
                      <span>{opcao.texto}</span>
                    </div>
                  </button>
                ))}
              </div>

              {/* Feedback */}
              {isRespostaCorreta === true && (
                <div className="text-center mb-6 fade-in-up">
                  <div className="inline-flex items-center space-x-2 bg-green-500/20 border border-green-400/50 rounded-2xl px-6 py-3">
                    <span className="text-2xl">🎉</span>
                    <span className="text-green-400 font-bold text-lg">Resposta Correta! +10 pontos</span>
                  </div>
                </div>
              )}
              
              {isRespostaCorreta === false && (
                <div className="text-center mb-6 fade-in-up">
                  <div className="inline-flex items-center space-x-2 bg-red-500/20 border border-red-400/50 rounded-2xl px-6 py-3">
                    <span className="text-2xl">😔</span>
                    <span className="text-red-400 font-bold text-lg">Resposta Incorreta! Tente novamente</span>
                  </div>
                </div>
              )}
              
              {/* Próxima Pergunta */}
              {respostaSelecionada !== null && (
                <div className="text-center fade-in-up">
                  <button 
                    onClick={buscarPergunta} 
                    className="btn-orange px-8 py-4 text-lg font-bold"
                  >
                    Próxima Pergunta
                  </button>
                </div>
              )}
            </div>
          </div>
        )}

        {/* Dicas */}
        <div className="max-w-4xl mx-auto mt-12 fade-in-up">
          <div className="glass-effect rounded-2xl p-6">
            <h3 className="text-xl font-bold text-white mb-4 flex items-center">
              <span className="mr-2">💡</span>
              Dicas para o Quiz
            </h3>
            <div className="grid md:grid-cols-2 gap-4 text-gray-300">
              <div className="flex items-start space-x-3">
                <span className="text-orange">•</span>
                <span>Cada resposta correta vale 10 pontos</span>
              </div>
              <div className="flex items-start space-x-3">
                <span className="text-orange">•</span>
                <span>As perguntas são sobre filmes populares</span>
              </div>
              <div className="flex items-start space-x-3">
                <span className="text-orange">•</span>
                <span>Não há limite de perguntas</span>
              </div>
              <div className="flex items-start space-x-3">
                <span className="text-orange">•</span>
                <span>Desafie-se a melhorar sua precisão!</span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default QuizPage;