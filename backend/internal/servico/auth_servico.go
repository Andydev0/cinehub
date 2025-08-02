package servico

import (
	"database/sql"
	"errors"
	"time"

	"github.com/Andydev0/filmes-backend/internal/dominio"
	"github.com/Andydev0/filmes-backend/internal/repositorio"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// Define os erros que o serviço pode retornar.
var (
	ErrEmailJaExiste        = errors.New("o e-mail fornecido já está em uso")
	ErrCredenciaisInvalidas = errors.New("credenciais inválidas")
)

// RegistroInput define os campos para o body do request de registro.
type RegistroInput struct {
	Nome  string `json:"nome" binding:"required"`
	Email string `json:"email" binding:"required,email"`
	Senha string `json:"senha" binding:"required,min=6"`
}

// LoginInput define os campos para o body do request de login.
type LoginInput struct {
	Email string `json:"email" binding:"required,email"`
	Senha string `json:"senha" binding:"required"`
}

// AuthServico é a interface que define os contratos do nosso serviço de autenticação.
type AuthServico interface {
	Registrar(input RegistroInput) (*dominio.Usuario, error)
	Login(input LoginInput) (string, error)
}

// authServicoImpl é a implementação da interface AuthServico.
type authServicoImpl struct {
	repo      repositorio.UsuarioRepositorio
	jwtSecret string
}

// NovoAuthServico cria a instância do serviço de autenticação com suas dependências.
func NovoAuthServico(repo repositorio.UsuarioRepositorio, jwtSecret string) AuthServico {
	return &authServicoImpl{
		repo:      repo,
		jwtSecret: jwtSecret,
	}
}

// Registrar executa a lógica de criar um novo usuário.
func (s *authServicoImpl) Registrar(input RegistroInput) (*dominio.Usuario, error) {
	// Busca o usuário para ver se o email já existe.
	usuarioExistente, err := s.repo.BuscarPorEmail(input.Email)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	if usuarioExistente != nil {
		return nil, ErrEmailJaExiste
	}

	// Gera o hash da senha para armazenamento seguro.
	senhaHash, err := bcrypt.GenerateFromPassword([]byte(input.Senha), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Cria o objeto de domínio do usuário.
	novoUsuario := &dominio.Usuario{
		Nome:      input.Nome,
		Email:     input.Email,
		SenhaHash: string(senhaHash),
	}

	// Salva o usuário no banco através do repositório.
	if err := s.repo.Salvar(novoUsuario); err != nil {
		return nil, err
	}

	return novoUsuario, nil
}

// Login executa a lógica de autenticação e retorna um token JWT.
func (s *authServicoImpl) Login(input LoginInput) (string, error) {
	// Busca o usuário pelo email.
	usuario, err := s.repo.BuscarPorEmail(input.Email)
	if err != nil {
		// Retorna erro genérico se o usuário não for encontrado.
		if err == sql.ErrNoRows {
			return "", ErrCredenciaisInvalidas
		}
		return "", err
	}

	// Compara a senha enviada com o hash salvo no banco.
	err = bcrypt.CompareHashAndPassword([]byte(usuario.SenhaHash), []byte(input.Senha))
	if err != nil {
		// Se a senha não bate, retorna o mesmo erro genérico.
		return "", ErrCredenciaisInvalidas
	}

	// Define as informações (claims) que irão no payload do token.
	claims := jwt.MapClaims{
		"sub": usuario.ID,                           // Subject (ID do usuário)
		"exp": time.Now().Add(time.Hour * 8).Unix(), // Expiração do token
	}

	// Cria o token com o método de assinatura e as claims.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Assina o token com a chave secreta e o retorna como uma string.
	tokenString, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
