package handler

import (
	"net/http"

	"github.com/Andydev0/filmes-backend/internal/dominio"
	"github.com/Andydev0/filmes-backend/internal/servico"
	"github.com/gin-gonic/gin"
)

// RecomendacaoHandler gerencia as requisições de recomendação.
type RecomendacaoHandler struct {
	servico servico.RecomendacaoServico
}

// NovoRecomendacaoHandler cria a instância do handler.
func NovoRecomendacaoHandler(s servico.RecomendacaoServico) *RecomendacaoHandler {
	return &RecomendacaoHandler{servico: s}
}

// ObterRecomendacoes é o método que lida com a rota GET /recomendacoes.
func (h *RecomendacaoHandler) ObterRecomendacoes(c *gin.Context) {
	usuarioID := c.MustGet("usuarioID").(int64)

	recomendacoes, err := h.servico.RecomendarFilmes(usuarioID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erro": "Falha ao gerar recomendações"})
		return
	}

	// Garante que a resposta seja sempre um array vazio e não null.
	if recomendacoes == nil {
		recomendacoes = make([]dominio.Filme, 0)
	}

	c.JSON(http.StatusOK, recomendacoes)
}
