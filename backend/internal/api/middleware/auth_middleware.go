package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// AuthMiddleware cria um middleware do Gin para validar o token JWT.
func AuthMiddleware(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Pega o header de autorização da requisição.
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"erro": "Header de autorização não encontrado"})
			return
		}

		// O header deve estar no formato "Bearer <token>".
		headerParts := strings.Split(authHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"erro": "Header de autorização mal formatado"})
			return
		}

		tokenString := headerParts[1]

		// Valida o token.
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Verifica se o método de assinatura é o esperado (HS256).
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("método de assinatura inesperado: %v", token.Header["alg"])
			}
			return []byte(jwtSecret), nil
		})

		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"erro": "Token inválido: " + err.Error()})
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			// Extrai o ID do usuário (subject) do token.
			usuarioID, ok := claims["sub"].(float64) // O parser JSON trata números como float64
			if !ok {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"erro": "Claim de usuário inválida"})
				return
			}

			// Adiciona o ID do usuário ao contexto da requisição.
			// As próximas funções (handlers) poderão acessar este valor.
			c.Set("usuarioID", int64(usuarioID))
			c.Next() // Passa a requisição para o próximo handler.
		} else {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"erro": "Token ou claims inválidos"})
		}
	}
}
