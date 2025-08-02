package dominio

import "time"

// Genero representa a estrutura de um gênero retornada pela API do TMDB.
type Genero struct {
	ID   int    `json:"id"`
	Nome string `json:"name"`
}

// Filme é a nossa estrutura de domínio principal para a API.
type Filme struct {
	ID             int     `json:"id"`
	Titulo         string  `json:"titulo"`
	Sinopse        string  `json:"sinopse"`
	DataLancamento string  `json:"dataLancamento"`
	CaminhoPoster  string  `json:"caminhoPoster"`
	NotaMedia      float64 `json:"notaMedia"`
}

// RespostaBuscaTMDB espelha a resposta da busca da API externa.
type RespostaBuscaTMDB struct {
	Resultados []TMDBMovieResult `json:"results"`
}

// TMDBMovieResult representa um filme na resposta da API externa.
type TMDBMovieResult struct {
	ID             int      `json:"id"`
	Titulo         string   `json:"title"`
	Sinopse        string   `json:"overview"`
	DataLancamento string   `json:"release_date"`
	CaminhoPoster  string   `json:"poster_path"`
	NotaMedia      float64  `json:"vote_average"`
	Generos        []Genero `json:"genres"`
}

// FilmeFavorito representa a tabela 'filmes_favoritos' no nosso banco de dados.
type FilmeFavorito struct {
	ID             int64     `db:"id" json:"id"`
	UsuarioID      int64     `db:"usuario_id" json:"usuarioId"`
	FilmeID        int64     `db:"filme_id" json:"filmeId"`
	Titulo         string    `db:"titulo" json:"titulo"`
	CaminhoPoster  string    `db:"caminho_poster" json:"caminhoPoster"`
	DataAdicionado time.Time `db:"data_adicionado" json:"dataAdicionado"`
}

// Usuario representa a tabela 'usuarios' no nosso banco de dados.
type Usuario struct {
	ID        int64  `db:"id"`
	Nome      string `db:"nome"`
	Email     string `db:"email"`
	SenhaHash string `db:"senha_hash"`
}

// OpcaoQuiz representa uma única opção de resposta em uma pergunta.
type OpcaoQuiz struct {
	ID    int    `json:"id"`
	Texto string `json:"texto"`
}

// PerguntaQuiz representa a estrutura completa de uma pergunta do quiz.
type PerguntaQuiz struct {
	Pergunta          string      `json:"pergunta"`
	Opcoes            []OpcaoQuiz `json:"opcoes"`
	RespostaCorretaID int         `json:"respostaCorretaId"`
}

// Avaliacao representa a tabela 'avaliacoes' no nosso banco de dados.
type Avaliacao struct {
	ID          int64     `db:"id" json:"id"`
	UsuarioID   int64     `db:"usuario_id" json:"usuarioId"`
	FilmeID     int64     `db:"filme_id" json:"filmeId"`
	Nota        int       `db:"nota" json:"nota"`
	Comentario  string    `db:"comentario" json:"comentario"`
	DataCriacao time.Time `db:"data_criacao" json:"dataCriacao"`
}

// AvaliacaoComUsuario é uma struct para enviar uma avaliação junto com o nome de quem a fez.
type AvaliacaoComUsuario struct {
	ID          int64     `db:"id" json:"id"`
	Nota        int       `db:"nota" json:"nota"`
	Comentario  string    `db:"comentario" json:"comentario"`
	DataCriacao time.Time `db:"data_criacao" json:"dataCriacao"`
	NomeUsuario string    `db:"nome" json:"nomeUsuario"` // Vem da tabela 'usuarios'
}

// MembroElenco representa um ator/atriz no elenco de um filme.
type MembroElenco struct {
	Nome        string `json:"name"`
	Personagem  string `json:"character"`
	CaminhoFoto string `json:"profile_path"`
}

// MembroEquipe representa uma pessoa da equipe técnica (diretor, escritor).
type MembroEquipe struct {
	Nome string `json:"name"`
	Job  string `json:"job"`
}

// CreditosTMDB armazena as listas de elenco e equipe da API.
type CreditosTMDB struct {
	Elenco []MembroElenco `json:"cast"`
	Equipe []MembroEquipe `json:"crew"`
}

// Video representa um vídeo (trailer, teaser) associado a um filme.
type Video struct {
	Key  string `json:"key"`
	Site string `json:"site"`
	Tipo string `json:"type"`
}

// VideosTMDB armazena a lista de vídeos da API.
type VideosTMDB struct {
	Resultados []Video `json:"results"`
}

// DetalhesFilmeCompleto é a nossa nova struct de resposta, combinando tudo.
type DetalhesFilmeCompleto struct {
	*TMDBMovieResult                // Inclui todos os campos de detalhes básicos
	Elenco           []MembroElenco `json:"elenco"`
	Diretor          string         `json:"diretor"`
	Escritores       []string       `json:"escritores"`
	TrailerKey       string         `json:"trailerKey"`
}
