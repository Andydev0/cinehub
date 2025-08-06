package handler

import (
	"net/http"
	"strconv"

	"github.com/Andydev0/filmes-backend/internal/dominio"
	"github.com/Andydev0/filmes-backend/internal/servico"
	"github.com/gin-gonic/gin"
)

type FavoritoHandler struct {
	servico servico.FavoritoServico
}

func NovoFavoritoHandler(s servico.FavoritoServico) *FavoritoHandler {
	return &FavoritoHandler{servico: s}
}

// Lógica de Adicionar atualizada
func (h *FavoritoHandler) Adicionar(c *gin.Context) {
	var input servico.AdicionarFavoritoInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "Input inválido"})
		return
	}

	usuarioID := c.MustGet("usuarioID").(int64)

	err := h.servico.AdicionarFavorito(usuarioID, input)
	if err != nil {
		// Verifica se o erro é de duplicata e retorna o status correto.
		if err == servico.ErrFavoritoJaExiste {
			c.JSON(http.StatusConflict, gin.H{"erro": err.Error()})
			return
		}
		// Para outros erros, retorna 500.
		c.JSON(http.StatusInternalServerError, gin.H{"erro": "Falha ao adicionar favorito"})
		return
	}

	c.Status(http.StatusCreated)
}

// ... (Listar e Remover continuam os mesmos)
func (h *FavoritoHandler) Listar(c *gin.Context) {
	usuarioID := c.MustGet("usuarioID").(int64)

	favoritos, err := h.servico.ListarFavoritos(usuarioID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erro": "Falha ao listar favoritos"})
		return
	}

	if favoritos == nil {
		favoritos = make([]dominio.FilmeFavorito, 0)
	}

	c.JSON(http.StatusOK, favoritos)
}

func (h *FavoritoHandler) Remover(c *gin.Context) {
	usuarioID := c.MustGet("usuarioID").(int64)

	filmeID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "ID de filme inválido"})
		return
	}

	if err := h.servico.RemoverFavorito(usuarioID, filmeID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erro": "Falha ao remover favorito"})
		return
	}

	c.Status(http.StatusNoContent)
}
