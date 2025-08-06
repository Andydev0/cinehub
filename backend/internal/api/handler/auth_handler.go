package handler

import (
	"net/http"

	"github.com/Andydev0/filmes-backend/internal/servico"
	"github.com/gin-gonic/gin"
)

// AuthHandler gerencia as requisições HTTP de autenticação.
type AuthHandler struct {
	servico servico.AuthServico
}

// NovoAuthHandler cria a instância do handler de autenticação.
func NovoAuthHandler(s servico.AuthServico) *AuthHandler {
	return &AuthHandler{servico: s}
}

// Registrar manipula a requisição de registro de um novo usuário.
func (h *AuthHandler) Registrar(c *gin.Context) {
	var input servico.RegistroInput

	// Valida e extrai os dados do JSON da requisição.
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "Dados de entrada inválidos: " + err.Error()})
		return
	}

	// Chama o serviço para executar a lógica de registro.
	_, err := h.servico.Registrar(input)
	if err != nil {
		// Retorna um erro específico se o email já estiver em uso.
		if err == servico.ErrEmailJaExiste {
			c.JSON(http.StatusConflict, gin.H{"erro": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"erro": "Falha ao registrar usuário"})
		return
	}

	// Retorna sucesso se o usuário for criado.
	c.JSON(http.StatusCreated, gin.H{"mensagem": "Usuário registrado com sucesso!"})
}

// Login manipula a requisição de login.
func (h *AuthHandler) Login(c *gin.Context) {
	var input servico.LoginInput

	// Valida e extrai os dados do JSON da requisição.
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "Dados de entrada inválidos"})
		return
	}

	// Chama o serviço para validar as credenciais e gerar o token.
	token, err := h.servico.Login(input)
	if err != nil {
		// Retorna 401 Unauthorized se as credenciais estiverem erradas.
		if err == servico.ErrCredenciaisInvalidas {
			c.JSON(http.StatusUnauthorized, gin.H{"erro": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"erro": "Falha ao realizar login"})
		return
	}

	// Se o login for bem-sucedido, retorna o token.
	c.JSON(http.StatusOK, gin.H{"token": token})
}
