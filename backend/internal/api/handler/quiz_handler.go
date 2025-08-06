package handler

import (
	"net/http"

	"github.com/Andydev0/filmes-backend/internal/servico"
	"github.com/gin-gonic/gin"
)

type QuizHandler struct {
	servico servico.QuizServico
}

func NovoQuizHandler(s servico.QuizServico) *QuizHandler {
	return &QuizHandler{servico: s}
}

// ObterPergunta agora usa o ID do usuário para gerar uma pergunta personalizada.
func (h *QuizHandler) ObterPergunta(c *gin.Context) {
	// Pega o ID do usuário logado que foi injetado pelo middleware.
	usuarioID := c.MustGet("usuarioID").(int64)

	// Passa o ID para o serviço.
	pergunta, err := h.servico.GerarPergunta(usuarioID)
	if err != nil {
		// Retorna um erro amigável se o usuário não tiver favoritos.
		c.JSON(http.StatusNotFound, gin.H{"erro": err.Error()})
		return
	}
	c.JSON(http.StatusOK, pergunta)
}
