package main

import (
	"log"
	"os"

	"github.com/Andydev0/filmes-backend/internal/api"
	"github.com/Andydev0/filmes-backend/internal/database"
	"github.com/joho/godotenv"
)

func main() {
	// Carrega as variáveis do arquivo .env.
	err := godotenv.Load()
	if err != nil {
		log.Println("Aviso: Não foi possível encontrar o arquivo .env.")
	}

	log.Println("Iniciando aplicação...")

	// Inicializa a conexão com o banco de dados.
	db := database.InitDB()
	defer db.Close()

	log.Println("Conexão com o banco de dados estabelecida com sucesso.")

	// Lê a chave da API do TMDB do ambiente.
	chaveAPI := os.Getenv("TMDB_API_KEY")
	if chaveAPI == "" {
		log.Fatal("A variável de ambiente TMDB_API_KEY é obrigatória.")
	}

	// Lê o segredo do JWT do ambiente.
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("A variável de ambiente JWT_SECRET é obrigatória.")
	}

	// Passa as configurações e a conexão com o banco para o roteador.
	roteador := api.SetupRouter(chaveAPI, db, jwtSecret)

	log.Println("Servidor iniciado na porta 8080")
	if err := roteador.Run(":8080"); err != nil {
		log.Fatalf("Falha ao iniciar o servidor: %v", err)
	}
}
