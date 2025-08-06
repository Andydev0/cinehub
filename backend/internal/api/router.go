// Package api contém a configuração do roteador HTTP e middlewares da aplicação.
// Este pacote é responsável por definir todas as rotas da API, configurar middlewares
// e injetar as dependências necessárias para os handlers.
package api

import (
	"github.com/Andydev0/filmes-backend/internal/api/handler"
	"github.com/Andydev0/filmes-backend/internal/api/middleware"
	"github.com/Andydev0/filmes-backend/internal/repositorio"
	"github.com/Andydev0/filmes-backend/internal/servico"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

// SetupRouter configura todas as rotas da API e suas dependências.
// Inicializa todos os componentes necessários (repositórios, serviços e handlers)
// e configura o middleware CORS para permitir requisições do frontend.
//
// Parâmetros:
//   - chaveAPI: Chave de API para o serviço TMDB
//   - db: Conexão com o banco de dados
//   - jwtSecret: Chave secreta para assinatura de tokens JWT
//
// Retorno:
//   - Engine do Gin configurado com todas as rotas e middlewares
func SetupRouter(chaveAPI string, db *sqlx.DB, jwtSecret string) *gin.Engine {
	// Inicialização de todos os componentes da aplicação usando injeção de dependência
	
	// Componentes relacionados a filmes
	filmeServico := servico.NovoFilmeServico(chaveAPI)
	filmeHandler := handler.NovoFilmeHandler(filmeServico)
	
	// Componentes relacionados a usuários e autenticação
	usuarioRepo := repositorio.NovoUsuarioRepositorio(db)
	authServico := servico.NovoAuthServico(usuarioRepo, jwtSecret)
	authHandler := handler.NovoAuthHandler(authServico)
	
	// Componentes relacionados a favoritos
	favoritoRepo := repositorio.NovoFavoritoRepositorio(db)
	favoritoServico := servico.NovoFavoritoServico(favoritoRepo)
	favoritoHandler := handler.NovoFavoritoHandler(favoritoServico)
	
	// Componentes relacionados a recomendações
	recomendacaoServico := servico.NovoRecomendacaoServico(favoritoRepo, filmeServico, chaveAPI)
	recomendacaoHandler := handler.NovoRecomendacaoHandler(recomendacaoServico)
	
	// Componentes relacionados ao quiz
	quizServico := servico.NovoQuizServico(chaveAPI, favoritoRepo, filmeServico)
	quizHandler := handler.NovoQuizHandler(quizServico)
	
	// Componentes relacionados a avaliações
	avaliacaoRepo := repositorio.NovaAvaliacaoRepositorio(db)
	avaliacaoServico := servico.NovaAvaliacaoServico(avaliacaoRepo)
	avaliacaoHandler := handler.NovaAvaliacaoHandler(avaliacaoServico)

	// Inicialização do router Gin
	router := gin.Default()
	
	// Configuração do middleware CORS
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:5173"} // Origem do frontend
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Authorization"}
	router.Use(cors.New(config))

	// Grupo de rotas com prefixo /v1 (versionamento da API)
	apiV1 := router.Group("/v1")
	{
		// ===== ROTAS PÚBLICAS =====
		
		// Rotas de autenticação
		auth := apiV1.Group("/auth")
		{
			// POST /v1/auth/registrar - Registra um novo usuário
			auth.POST("/registrar", authHandler.Registrar)
			
			// POST /v1/auth/login - Autentica um usuário existente
			auth.POST("/login", authHandler.Login)
		}

		// Rotas de busca de filmes (públicas)
		
		// GET /v1/filmes/buscar?termo={termo} - Busca filmes por termo
		apiV1.GET("/filmes/buscar", filmeHandler.BuscarFilmes)
		
		// GET /v1/filmes/genero?generoId={id} - Busca filmes por gênero
		apiV1.GET("/filmes/genero", filmeHandler.BuscarFilmesPorGenero)
		
		// GET /v1/filmes/aleatorio?generoId={id}&ano={ano} - Busca filme aleatório
		apiV1.GET("/filmes/aleatorio", filmeHandler.BuscarAleatorio)
		
		// GET /v1/generos - Lista todos os gêneros disponíveis
		apiV1.GET("/generos", filmeHandler.ListarGeneros)

		// Grupo para rotas relacionadas a um filme específico pelo ID
		filmePorId := apiV1.Group("/filmes/:id")
		{
			// GET /v1/filmes/:id - Busca detalhes de um filme específico
			filmePorId.GET("", filmeHandler.BuscarDetalhes)
			
			// GET /v1/filmes/:id/avaliacoes - Lista avaliações de um filme específico
			filmePorId.GET("/avaliacoes", avaliacaoHandler.ListarPorFilme)
		}

		// ===== ROTAS PROTEGIDAS =====
		// Todas as rotas abaixo requerem autenticação via JWT
		autenticado := apiV1.Group("/")
		autenticado.Use(middleware.AuthMiddleware(jwtSecret))
		{
			// Rotas para gerenciamento de favoritos
			favoritos := autenticado.Group("/favoritos")
			{
				// POST /v1/favoritos - Adiciona um filme aos favoritos
				favoritos.POST("", favoritoHandler.Adicionar)
				
				// GET /v1/favoritos - Lista todos os filmes favoritos do usuário
				favoritos.GET("", favoritoHandler.Listar)
				
				// DELETE /v1/favoritos/:id - Remove um filme dos favoritos
				favoritos.DELETE("/:id", favoritoHandler.Remover)
			}
			
			// GET /v1/recomendacoes - Obtém recomendações personalizadas
			autenticado.GET("/recomendacoes", recomendacaoHandler.ObterRecomendacoes)
			
			// GET /v1/quiz/pergunta - Obtém uma pergunta para o quiz
			autenticado.GET("/quiz/pergunta", quizHandler.ObterPergunta)

			// POST /v1/filmes/:id/avaliacoes - Cria uma nova avaliação para um filme
			autenticado.POST("/filmes/:id/avaliacoes", avaliacaoHandler.Criar)
		}
	}
	
	return router
}
