// Package handler contém os manipuladores HTTP para a API de filmes.
// Este pacote é responsável por receber requisições HTTP, validar parâmetros,
// chamar os serviços apropriados e formatar as respostas.
package handler

import (
	"net/http"
	"strconv"

	"github.com/Andydev0/filmes-backend/internal/servico"
	"github.com/gin-gonic/gin"
)

// FilmeHandler encapsula a lógica de manipulação de requisições relacionadas a filmes.
// Utiliza uma instância de FilmeServico para executar operações de negócio.
type FilmeHandler struct {
	servico servico.FilmeServico // Serviço responsável pelas operações com filmes
}

// NovoFilmeHandler cria uma nova instância de FilmeHandler.
// Recebe uma implementação de FilmeServico como dependência.
// Parâmetros:
//   - s: Implementação da interface FilmeServico
// Retorno:
//   - Ponteiro para a instância criada de FilmeHandler
func NovoFilmeHandler(s servico.FilmeServico) *FilmeHandler {
	return &FilmeHandler{
		servico: s,
	}
}

// BuscarFilmes processa requisições para buscar filmes por termo de busca.
// Endpoint: GET /filmes/buscar?termo={termo}
// Parâmetros de consulta:
//   - termo: Termo para busca de filmes (obrigatório)
// Respostas:
//   - 200 OK: Lista de filmes encontrados
//   - 400 Bad Request: Parâmetro 'termo' ausente
//   - 500 Internal Server Error: Erro ao processar a requisição
func (h *FilmeHandler) BuscarFilmes(c *gin.Context) {
	// Extrai o termo de busca da query string
	termoDeBusca := c.Query("termo")
	if termoDeBusca == "" {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "O parâmetro 'termo' é obrigatório"})
		return
	}

	// Chama o serviço para buscar filmes pelo termo
	filmes, err := h.servico.BuscarFilmes(termoDeBusca)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erro": "Falha ao buscar filmes"})
		return
	}

	// Retorna os filmes encontrados
	c.JSON(http.StatusOK, filmes)
}

// ListarGeneros processa requisições para listar todos os gêneros de filmes disponíveis.
// Endpoint: GET /generos
// Respostas:
//   - 200 OK: Lista de gêneros disponíveis
//   - 500 Internal Server Error: Erro ao processar a requisição
func (h *FilmeHandler) ListarGeneros(c *gin.Context) {
	// Chama o serviço para listar os gêneros
	generos, err := h.servico.ListarGeneros()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erro": "Falha ao buscar gêneros"})
		return
	}
	
	// Retorna a lista de gêneros
	c.JSON(http.StatusOK, generos)
}

// BuscarAleatorio processa requisições para buscar um filme aleatório.
// Endpoint: GET /filmes/aleatorio?generoId={id}&ano={ano}
// Parâmetros de consulta (opcionais):
//   - generoId: ID do gênero para filtrar
//   - ano: Ano de lançamento para filtrar
// Respostas:
//   - 200 OK: Filme aleatório encontrado
//   - 404 Not Found: Nenhum filme encontrado com os filtros fornecidos
//   - 500 Internal Server Error: Erro ao processar a requisição
func (h *FilmeHandler) BuscarAleatorio(c *gin.Context) {
	// Extrai os parâmetros de filtro da query string
	generoID := c.Query("generoId")
	ano := c.Query("ano")

	// Chama o serviço para buscar um filme aleatório com os filtros
	filme, err := h.servico.BuscarFilmeAleatorio(generoID, ano)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"erro": err.Error()})
		return
	}
	
	// Retorna o filme encontrado
	c.JSON(http.StatusOK, filme)
}

// BuscarDetalhes processa requisições para buscar detalhes completos de um filme específico.
// Endpoint: GET /filmes/:id
// Parâmetros de rota:
//   - id: ID do filme (obrigatório)
// Respostas:
//   - 200 OK: Detalhes completos do filme
//   - 400 Bad Request: ID de filme inválido
//   - 404 Not Found: Filme não encontrado
//   - 500 Internal Server Error: Erro ao processar a requisição
func (h *FilmeHandler) BuscarDetalhes(c *gin.Context) {
	// Extrai e converte o ID do filme do parâmetro da URL
	filmeID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "ID de filme inválido"})
		return
	}

	// Chama o serviço para buscar os detalhes completos do filme
	filme, err := h.servico.BuscarDetalhes(filmeID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"erro": "Filme não encontrado"})
		return
	}

	// Retorna os detalhes do filme
	c.JSON(http.StatusOK, filme)
}

// BuscarFilmesPorGenero processa requisições para buscar filmes de um gênero específico.
// Endpoint: GET /filmes/genero?generoId={id}
// Parâmetros de consulta:
//   - generoId: ID do gênero (obrigatório)
// Respostas:
//   - 200 OK: Lista de filmes do gênero especificado
//   - 400 Bad Request: Parâmetro 'generoId' ausente
//   - 500 Internal Server Error: Erro ao processar a requisição
func (h *FilmeHandler) BuscarFilmesPorGenero(c *gin.Context) {
	// Extrai o ID do gênero da query string
	generoID := c.Query("generoId")
	if generoID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "O parâmetro 'generoId' é obrigatório"})
		return
	}

	// Chama o serviço para buscar filmes do gênero especificado
	filmes, err := h.servico.BuscarFilmesPorGenero(generoID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erro": "Falha ao buscar filmes por gênero"})
		return
	}

	// Retorna os filmes encontrados
	c.JSON(http.StatusOK, filmes)
}
