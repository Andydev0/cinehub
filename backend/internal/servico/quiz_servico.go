// Package servico contém a lógica de negócio da aplicação.
// Este arquivo implementa o serviço de quiz, que gera perguntas personalizadas
// baseadas nos filmes favoritos dos usuários.
package servico

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/Andydev0/filmes-backend/internal/dominio"
	"github.com/Andydev0/filmes-backend/internal/repositorio"
)

// PessoaPopular representa uma pessoa (ator/diretor) popular da API TMDB.
// Usada para obter informações sobre pessoas populares para gerar opções
// de respostas para o quiz.
type PessoaPopular struct {
	ID           int     `json:"id"`           // ID da pessoa na API TMDB
	Nome         string  `json:"name"`         // Nome da pessoa
	Departamento string  `json:"known_for_department"` // Departamento (Acting, Directing, etc)
	Popularidade float64 `json:"popularity"`   // Pontuação de popularidade
}

// RespostaPessoasPopulares representa a resposta da API TMDB para busca de pessoas populares.
type RespostaPessoasPopulares struct {
	Resultados []PessoaPopular `json:"results"` // Lista de pessoas populares retornadas
}

// QuizServico define a interface para o serviço de quiz.
// Permite gerar perguntas personalizadas baseadas nos filmes favoritos do usuário.
type QuizServico interface {
	// GerarPergunta cria uma pergunta de quiz personalizada para um usuário.
	// Parâmetros:
	//   - usuarioID: ID do usuário para quem a pergunta será gerada
	// Retorno:
	//   - Pergunta de quiz personalizada ou erro se não for possível gerar
	GerarPergunta(usuarioID int64) (*dominio.PerguntaQuiz, error)
}

// quizServicoImpl implementa a interface QuizServico.
// Gera perguntas variadas sobre filmes favoritos do usuário.
type quizServicoImpl struct {
	apiKey       string                         // Chave da API TMDB
	favoritoRepo repositorio.FavoritoRepositorio // Repositório de filmes favoritos
	filmeServico FilmeServico                   // Serviço para buscar informações de filmes
	historicoQuiz map[int64][]int64             // Cache para evitar repetição de perguntas (usuarioID -> filmeIDs já usados)
}

// NovoQuizServico cria uma nova instância do serviço de quiz.
// Inicializa o serviço com as dependências necessárias e cria o cache de histórico.
//
// Parâmetros:
//   - apiKey: Chave de API para o serviço TMDB
//   - favoritoRepo: Repositório para acessar os filmes favoritos dos usuários
//   - filmeServico: Serviço para buscar detalhes de filmes
//
// Retorno:
//   - Uma implementação da interface QuizServico
func NovoQuizServico(apiKey string, favoritoRepo repositorio.FavoritoRepositorio, filmeServico FilmeServico) QuizServico {
	return &quizServicoImpl{
		apiKey:       apiKey,
		favoritoRepo: favoritoRepo,
		filmeServico: filmeServico,
		historicoQuiz: make(map[int64][]int64), // Inicializa o mapa de histórico vazio
	}
}

// GerarPergunta cria uma pergunta de quiz personalizada baseada nos filmes favoritos do usuário.
// Implementa um sistema de rotação de filmes para evitar repetições e oferece diferentes tipos
// de perguntas (ano de lançamento, diretor, ator, gênero) para maior variedade.
//
// Parâmetros:
//   - usuarioID: ID do usuário para quem a pergunta será gerada
//
// Retorno:
//   - Pergunta de quiz personalizada ou erro se não for possível gerar
func (s *quizServicoImpl) GerarPergunta(usuarioID int64) (*dominio.PerguntaQuiz, error) {
	// Busca os filmes favoritos do usuário no repositório
	favoritos, err := s.favoritoRepo.ListarPorUsuarioID(usuarioID)
	if err != nil || len(favoritos) < 1 {
		return nil, fmt.Errorf("usuário não tem filmes favoritos suficientes para o quiz")
	}

	// Filtra filmes que ainda não foram usados no quiz para este usuário
	historicoUsuario := s.historicoQuiz[usuarioID]
	filmesDisponiveis := []dominio.FilmeFavorito{}
	
	// Verifica cada filme favorito para ver se já foi usado recentemente
	for _, favorito := range favoritos {
		jáUsado := false
		for _, filmeUsado := range historicoUsuario {
			if favorito.FilmeID == filmeUsado {
				jáUsado = true
				break
			}
		}
		if !jáUsado {
			filmesDisponiveis = append(filmesDisponiveis, favorito)
		}
	}
	
	// Se todos os filmes já foram usados, reinicia o ciclo usando todos novamente
	if len(filmesDisponiveis) == 0 {
		s.historicoQuiz[usuarioID] = []int64{} // Limpa o histórico
		filmesDisponiveis = favoritos         // Usa todos os favoritos
	}

	// Escolhe um filme favorito aleatório dentre os disponíveis
	filmeCorretoFavorito := filmesDisponiveis[rand.Intn(len(filmesDisponiveis))]
	
	// Busca detalhes completos do filme para gerar a pergunta
	detalhesFilmeCorreto, err := s.filmeServico.BuscarDetalhes(filmeCorretoFavorito.FilmeID)
	if err != nil {
		return nil, fmt.Errorf("falha ao buscar detalhes do filme correto")
	}
	
	// Registra o filme no histórico do usuário para evitar repetição
	s.historicoQuiz[usuarioID] = append(s.historicoQuiz[usuarioID], filmeCorretoFavorito.FilmeID)

	// Escolhe aleatoriamente o tipo de pergunta para maior variedade
	// 1: Ano de lançamento, 2: Diretor, 3: Ator, 4: Gênero
	tipoPergunta := rand.Intn(4) + 1

	// Gera a pergunta de acordo com o tipo escolhido
	switch tipoPergunta {
	case 1:
		return s.gerarPerguntaAno(detalhesFilmeCorreto)
	case 2:
		return s.gerarPerguntaDiretor(detalhesFilmeCorreto)
	case 3:
		return s.gerarPerguntaAtor(detalhesFilmeCorreto)
	case 4:
		return s.gerarPerguntaGenero(detalhesFilmeCorreto)
	default:
		// Fallback para pergunta de ano (mais simples e sempre disponível)
		return s.gerarPerguntaAno(detalhesFilmeCorreto)
	}
}

// gerarPerguntaAno gera uma pergunta sobre o ano de lançamento do filme.
// Cria opções de anos próximos ao ano correto para desafiar o usuário.
//
// Parâmetros:
//   - filme: Detalhes completos do filme para o qual a pergunta será gerada
//
// Retorno:
//   - Pergunta formatada ou erro se não for possível gerar
func (s *quizServicoImpl) gerarPerguntaAno(filme *dominio.DetalhesFilmeCompleto) (*dominio.PerguntaQuiz, error) {
	// Verifica se o filme tem data de lançamento válida
	if filme.DataLancamento == "" {
		return nil, fmt.Errorf("filme não tem data de lançamento")
	}

	// Extrai o ano da data de lançamento (formato: "2023-05-15")
	anoCorreto, err := strconv.Atoi(filme.DataLancamento[:4])
	if err != nil {
		return nil, fmt.Errorf("erro ao extrair ano do filme: %v", err)
	}

	// Gera 3 anos aleatórios diferentes como opções incorretas
	anosIncorretos := make([]int, 0, 3)
	for len(anosIncorretos) < 3 {
		// Gera um ano aleatório entre anoCorreto-10 e anoCorreto+10
		// para criar opções plausíveis mas incorretas
		variacaoAno := rand.Intn(21) - 10 // -10 a +10
		anoAleatorio := anoCorreto + variacaoAno

		// Garante que o ano está em um intervalo realista (1900-2025)
		// Evita anos futuros ou muito antigos
		if anoAleatorio < 1900 {
			anoAleatorio = 1900 + rand.Intn(20) // 1900-1919
		} else if anoAleatorio > 2025 {
			anoAleatorio = 2005 + rand.Intn(20) // 2005-2024
		}

		// Verifica se o ano é diferente do correto e não está já nas opções
		// Isso evita duplicatas nas opções e garante que a resposta correta não esteja entre as incorretas
		if anoAleatorio != anoCorreto && !contemInt(anosIncorretos, anoAleatorio) {
			anosIncorretos = append(anosIncorretos, anoAleatorio)
		}
	}

	// Cria as opções para a pergunta
	opcoes := []dominio.OpcaoQuiz{
		{ID: 1, Texto: strconv.Itoa(anoCorreto)},
		{ID: 2, Texto: strconv.Itoa(anosIncorretos[0])},
		{ID: 3, Texto: strconv.Itoa(anosIncorretos[1])},
		{ID: 4, Texto: strconv.Itoa(anosIncorretos[2])},
	}

	// Embaralha as opções para que a resposta correta não esteja sempre na mesma posição
	rand.Shuffle(len(opcoes), func(i, j int) {
		opcoes[i], opcoes[j] = opcoes[j], opcoes[i]
	})

	// Encontra o ID da opção correta após embaralhar
	respostaCorretaID := 0
	for i, op := range opcoes {
		if op.Texto == strconv.Itoa(anoCorreto) {
			respostaCorretaID = opcoes[i].ID
			break
		}
	}

	// Cria a estrutura da pergunta com texto, opções e ID da resposta correta
	pergunta := &dominio.PerguntaQuiz{
		Pergunta:          fmt.Sprintf("Em que ano foi lançado o filme '%s'?", filme.Titulo),
		Opcoes:            opcoes,
		RespostaCorretaID: respostaCorretaID,
	}

	return pergunta, nil
}

// contemInt verifica se um valor inteiro está presente em uma slice de inteiros.
// É uma função auxiliar para evitar duplicação de opções no quiz.
//
// Parâmetros:
//   - slice: A slice de inteiros onde procurar
//   - valor: O valor inteiro a ser procurado
//
// Retorno:
//   - true se o valor estiver presente, false caso contrário
func contemInt(slice []int, valor int) bool {
	for _, v := range slice {
		if v == valor {
			return true
		}
	}
	return false
}

// gerarPerguntaDiretor gera uma pergunta sobre o diretor do filme.
// Utiliza uma lista curada de diretores famosos como opções incorretas.
//
// Parâmetros:
//   - detalhes: Detalhes completos do filme para o qual a pergunta será gerada
//
// Retorno:
//   - Pergunta formatada ou erro se não for possível gerar
//   - Em caso de filme sem diretor, faz fallback para pergunta sobre o ano
func (s *quizServicoImpl) gerarPerguntaDiretor(detalhes *dominio.DetalhesFilmeCompleto) (*dominio.PerguntaQuiz, error) {
	// Verifica se o filme tem informação de diretor disponível
	if detalhes.Diretor == "" {
		// Se não temos o diretor, voltamos para pergunta de ano que é mais simples
		return s.gerarPerguntaAno(detalhes)
	}

	// Busca diretores de filmes populares para opções incorretas
	diretoresIncorretos, err := s.buscarDiretoresPopulares(detalhes.Diretor)
	if err != nil || len(diretoresIncorretos) < 3 {
		return s.gerarPerguntaAno(detalhes) // Fallback para pergunta de ano
	}

	// Cria as opções para a pergunta com o diretor correto e 3 incorretos
	opcoes := []dominio.OpcaoQuiz{
		{ID: 1, Texto: detalhes.Diretor},
		{ID: 2, Texto: diretoresIncorretos[0]},
		{ID: 3, Texto: diretoresIncorretos[1]},
		{ID: 4, Texto: diretoresIncorretos[2]},
	}

	// Embaralha as opções para que a resposta correta não esteja sempre na mesma posição
	rand.Shuffle(len(opcoes), func(i, j int) { 
		opcoes[i], opcoes[j] = opcoes[j], opcoes[i] 
	})

	// Encontra o ID da resposta correta após embaralhar
	respostaCorretaID := 0
	for i, op := range opcoes {
		if op.Texto == detalhes.Diretor {
			respostaCorretaID = opcoes[i].ID
			break
		}
	}

	// Cria e retorna a estrutura completa da pergunta
	return &dominio.PerguntaQuiz{
		Pergunta:          fmt.Sprintf("Quem dirigiu o filme '%s'?", detalhes.Titulo),
		Opcoes:            opcoes,
		RespostaCorretaID: respostaCorretaID,
	}, nil
}

// gerarPerguntaAtor gera uma pergunta sobre um ator que participou do filme.
// Seleciona um ator do elenco principal e busca atores populares para opções incorretas.
//
// Parâmetros:
//   - detalhes: Detalhes completos do filme para o qual a pergunta será gerada
//
// Retorno:
//   - Pergunta formatada ou erro se não for possível gerar
//   - Em caso de filme sem elenco, faz fallback para pergunta sobre o ano
func (s *quizServicoImpl) gerarPerguntaAtor(detalhes *dominio.DetalhesFilmeCompleto) (*dominio.PerguntaQuiz, error) {
	// Verifica se o filme tem informação de elenco disponível
	if len(detalhes.Elenco) == 0 {
		// Se não temos elenco, voltamos para pergunta de ano que é mais simples
		return s.gerarPerguntaAno(detalhes)
	}

	// Escolhe um ator principal (um dos primeiros 3) para ser a resposta correta
	atorCorreto := detalhes.Elenco[0]
	if len(detalhes.Elenco) > 1 {
		atorCorreto = detalhes.Elenco[rand.Intn(min(3, len(detalhes.Elenco)))]
	}

	// Busca atores de filmes populares para opções incorretas
	atoresIncorretos, err := s.buscarAtoresPopulares(atorCorreto.Nome)
	if err != nil || len(atoresIncorretos) < 3 {
		// Se não conseguimos obter atores suficientes, voltamos para pergunta de ano
		return s.gerarPerguntaAno(detalhes)
	}

	// Cria as opções para a pergunta com o ator correto e 3 incorretos
	opcoes := []dominio.OpcaoQuiz{
		{ID: 1, Texto: atorCorreto.Nome},
		{ID: 2, Texto: atoresIncorretos[0]},
		{ID: 3, Texto: atoresIncorretos[1]},
		{ID: 4, Texto: atoresIncorretos[2]},
	}

	// Embaralha as opções para que a resposta correta não esteja sempre na mesma posição
	rand.Shuffle(len(opcoes), func(i, j int) { 
		opcoes[i], opcoes[j] = opcoes[j], opcoes[i] 
	})

	// Encontra o ID da resposta correta após embaralhar
	respostaCorretaID := 0
	for i, op := range opcoes {
		if op.Texto == atorCorreto.Nome {
			respostaCorretaID = opcoes[i].ID
			break
		}
	}

	// Cria e retorna a estrutura completa da pergunta
	return &dominio.PerguntaQuiz{
		Pergunta:          fmt.Sprintf("Qual destes atores participou do filme '%s'?", detalhes.Titulo),
		Opcoes:            opcoes,
		RespostaCorretaID: respostaCorretaID,
	}, nil
}
	
// buscarDiretoresPopulares busca diretores populares da API do TMDB para usar como opções incorretas.
// Utiliza a API de pessoas populares e filtra por departamento "Directing".
//
// Parâmetros:
//   - diretorCorreto: Nome do diretor que deve ser excluído da lista retornada
//
// Retorno:
//   - Lista de 3 nomes de diretores diferentes do diretor correto
//   - Erro se ocorrer algum problema na requisição
func (s *quizServicoImpl) buscarDiretoresPopulares(diretorCorreto string) ([]string, error) {
	// Busca pessoas populares que são diretores
	paginaAleatoria := rand.Intn(5) + 1 // Páginas 1-5 para ter variedade
	url := fmt.Sprintf("https://api.themoviedb.org/3/person/popular?api_key=%s&language=pt-BR&page=%d", s.apiKey, paginaAleatoria)
	
	resp, err := http.Get(url)
	if err != nil {
		return s.buscarDiretoresFallback(diretorCorreto), nil // Fallback para lista estática
	}
	defer resp.Body.Close()
	
	var respostaPessoas RespostaPessoasPopulares
	if err := json.NewDecoder(resp.Body).Decode(&respostaPessoas); err != nil {
		return s.buscarDiretoresFallback(diretorCorreto), nil // Fallback
	}
	
	// Filtra apenas diretores e remove o diretor correto
	diretoresDisponiveis := []string{}
	for _, pessoa := range respostaPessoas.Resultados {
		// Considera pessoas do departamento "Directing" ou com alta popularidade
		if (pessoa.Departamento == "Directing" || pessoa.Popularidade > 10) && 
		   pessoa.Nome != diretorCorreto && 
		   len(diretoresDisponiveis) < 10 { // Coleta até 10 para ter opções
			diretoresDisponiveis = append(diretoresDisponiveis, pessoa.Nome)
		}
	}
	
	// Se não encontrou diretores suficientes na API, usa fallback
	if len(diretoresDisponiveis) < 3 {
		return s.buscarDiretoresFallback(diretorCorreto), nil
	}
	
	// Embaralha e seleciona 3
	rand.Shuffle(len(diretoresDisponiveis), func(i, j int) {
		diretoresDisponiveis[i], diretoresDisponiveis[j] = diretoresDisponiveis[j], diretoresDisponiveis[i]
	})
	
	return diretoresDisponiveis[:3], nil
}

// buscarDiretoresFallback lista de fallback para quando a API falha
func (s *quizServicoImpl) buscarDiretoresFallback(diretorCorreto string) []string {
	diretoresComuns := []string{
		"Christopher Nolan", "Steven Spielberg", "Martin Scorsese", "Quentin Tarantino",
		"James Cameron", "Ridley Scott", "David Fincher", "Tim Burton",
		"Denis Villeneuve", "Jordan Peele", "Greta Gerwig", "Damien Chazelle",
	}
	
	diretoresDisponiveis := []string{}
	for _, diretor := range diretoresComuns {
		if diretor != diretorCorreto {
			diretoresDisponiveis = append(diretoresDisponiveis, diretor)
		}
	}
	
	rand.Shuffle(len(diretoresDisponiveis), func(i, j int) {
		diretoresDisponiveis[i], diretoresDisponiveis[j] = diretoresDisponiveis[j], diretoresDisponiveis[i]
	})
	
	quantidade := min(3, len(diretoresDisponiveis))
	return diretoresDisponiveis[:quantidade]
}

// buscarAtoresPopulares busca atores populares da API do TMDB
func (s *quizServicoImpl) buscarAtoresPopulares(atorCorreto string) ([]string, error) {
	// Busca pessoas populares que são atores
	paginaAleatoria := rand.Intn(10) + 1 // Páginas 1-10 para maior variedade de atores
	url := fmt.Sprintf("https://api.themoviedb.org/3/person/popular?api_key=%s&language=pt-BR&page=%d", s.apiKey, paginaAleatoria)
	
	resp, err := http.Get(url)
	if err != nil {
		return s.buscarAtoresFallback(atorCorreto), nil // Fallback para lista estática
	}
	defer resp.Body.Close()
	
	var respostaPessoas RespostaPessoasPopulares
	if err := json.NewDecoder(resp.Body).Decode(&respostaPessoas); err != nil {
		return s.buscarAtoresFallback(atorCorreto), nil // Fallback
	}
	
	// Filtra apenas atores e remove o ator correto
	atoresDisponiveis := []string{}
	for _, pessoa := range respostaPessoas.Resultados {
		// Considera pessoas do departamento "Acting" ou com alta popularidade
		if (pessoa.Departamento == "Acting" || pessoa.Popularidade > 15) && 
		   pessoa.Nome != atorCorreto && 
		   len(atoresDisponiveis) < 15 { // Coleta até 15 para ter mais opções
			atoresDisponiveis = append(atoresDisponiveis, pessoa.Nome)
		}
	}
	
	// Se não encontrou atores suficientes na API, usa fallback
	if len(atoresDisponiveis) < 3 {
		return s.buscarAtoresFallback(atorCorreto), nil
	}
	
	// Embaralha e seleciona 3
	rand.Shuffle(len(atoresDisponiveis), func(i, j int) {
		atoresDisponiveis[i], atoresDisponiveis[j] = atoresDisponiveis[j], atoresDisponiveis[i]
	})
	
	return atoresDisponiveis[:3], nil
}

// buscarAtoresFallback lista de fallback para quando a API falha
func (s *quizServicoImpl) buscarAtoresFallback(atorCorreto string) []string {
	atoresComuns := []string{
		"Leonardo DiCaprio", "Brad Pitt", "Tom Hanks", "Will Smith",
		"Robert Downey Jr.", "Scarlett Johansson", "Jennifer Lawrence", "Emma Stone",
		"Ryan Gosling", "Margot Robbie", "Timothée Chalamet", "Zendaya",
		"Chris Evans", "Gal Gadot", "Ryan Reynolds", "Sandra Bullock",
	}
	
	atoresDisponiveis := []string{}
	for _, ator := range atoresComuns {
		if ator != atorCorreto {
			atoresDisponiveis = append(atoresDisponiveis, ator)
		}
	}
	
	rand.Shuffle(len(atoresDisponiveis), func(i, j int) {
		atoresDisponiveis[i], atoresDisponiveis[j] = atoresDisponiveis[j], atoresDisponiveis[i]
	})
	
	quantidade := min(3, len(atoresDisponiveis))
	return atoresDisponiveis[:quantidade]
}

// gerarPerguntaGenero gera uma pergunta sobre um dos gêneros do filme.
// Seleciona um gênero real do filme e busca gêneros comuns para opções incorretas.
//
// Parâmetros:
//   - detalhes: Detalhes completos do filme para o qual a pergunta será gerada
//
// Retorno:
//   - Pergunta formatada ou erro se não for possível gerar
//   - Em caso de filme sem gêneros, faz fallback para pergunta sobre o ano
func (s *quizServicoImpl) gerarPerguntaGenero(detalhes *dominio.DetalhesFilmeCompleto) (*dominio.PerguntaQuiz, error) {
	// Verifica se o filme tem informação de gêneros disponível
	if len(detalhes.Generos) == 0 {
		// Se não temos gêneros, voltamos para pergunta de ano que é mais simples
		return s.gerarPerguntaAno(detalhes)
	}

	// Escolhe um gênero do filme aleatoriamente para ser a resposta correta
	generoCorreto := detalhes.Generos[0].Nome
	if len(detalhes.Generos) > 1 {
		generoCorreto = detalhes.Generos[rand.Intn(len(detalhes.Generos))].Nome
	}

	// Lista curada de gêneros cinematográficos comuns para usar como opções incorretas
	generosComuns := []string{
		"Ação", "Comédia", "Drama", "Terror", "Romance", "Ficção científica",
		"Aventura", "Thriller", "Animação", "Documentário", "Fantasia", "Crime",
		"Mistério", "Família", "Guerra", "História", "Música", "Faroeste",
	}
	
	// Filtra gêneros que não são o gênero correto para evitar duplicidade
	generosDisponiveis := []string{}
	for _, genero := range generosComuns {
		if genero != generoCorreto {
			generosDisponiveis = append(generosDisponiveis, genero)
		}
	}

	// Verifica se temos gêneros suficientes para criar a pergunta
	if len(generosDisponiveis) < 3 {
		// Se não temos gêneros suficientes, voltamos para pergunta de ano
		return s.gerarPerguntaAno(detalhes)
	}

	// Embaralha a lista de gêneros disponíveis para garantir aleatoriedade
	rand.Shuffle(len(generosDisponiveis), func(i, j int) {
		generosDisponiveis[i], generosDisponiveis[j] = generosDisponiveis[j], generosDisponiveis[i]
	})

	// Cria as opções para a pergunta com o gênero correto e 3 incorretos
	opcoes := []dominio.OpcaoQuiz{
		{ID: 1, Texto: generoCorreto},
		{ID: 2, Texto: generosDisponiveis[0]},
		{ID: 3, Texto: generosDisponiveis[1]},
		{ID: 4, Texto: generosDisponiveis[2]},
	}

	// Embaralha as opções para que a resposta correta não esteja sempre na mesma posição
	rand.Shuffle(len(opcoes), func(i, j int) {
		opcoes[i], opcoes[j] = opcoes[j], opcoes[i]
	})

	// Encontra o ID da resposta correta após embaralhar
	respostaCorretaID := 0
	for i, op := range opcoes {
		if op.Texto == generoCorreto {
			respostaCorretaID = opcoes[i].ID
			break
		}
	}

	// Cria e retorna a estrutura completa da pergunta
	return &dominio.PerguntaQuiz{
		Pergunta:          fmt.Sprintf("Qual é um dos gêneros do filme '%s'?", detalhes.Titulo),
		Opcoes:            opcoes,
		RespostaCorretaID: respostaCorretaID,
	}, nil
}

// min retorna o menor entre dois números
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
