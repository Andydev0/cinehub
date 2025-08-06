package handler

import (
	"net/http"
	"strconv"

	"github.com/Andydev0/filmes-backend/internal/dominio"
	"github.com/Andydev0/filmes-backend/internal/servico"
	"github.com/gin-gonic/gin"
)

type AvaliacaoHandler struct{ servico servico.AvaliacaoServico }

func NovaAvaliacaoHandler(s servico.AvaliacaoServico) *AvaliacaoHandler {
	return &AvaliacaoHandler{servico: s}
}

func (h *AvaliacaoHandler) Criar(c *gin.Context) {
	filmeID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "ID de filme inválido"})
		return
	}
	usuarioID := c.MustGet("usuarioID").(int64)

	var input servico.AvaliacaoInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": err.Error()})
		return
	}

	if err := h.servico.Criar(usuarioID, filmeID, input); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erro": "Falha ao salvar avaliação"})
		return
	}
	c.Status(http.StatusCreated)
}

func (h *AvaliacaoHandler) ListarPorFilme(c *gin.Context) {
	filmeID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "ID de filme inválido"})
		return
	}

	avaliacoes, err := h.servico.ListarPorFilme(filmeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erro": "Falha ao buscar avaliações"})
		return
	}

	if avaliacoes == nil {
		avaliacoes = make([]dominio.AvaliacaoComUsuario, 0)
	}
	c.JSON(http.StatusOK, avaliacoes)
}
