package servico

import (
	"errors"

	"github.com/Andydev0/filmes-backend/internal/dominio"
	"github.com/Andydev0/filmes-backend/internal/repositorio"
)

var ErrFavoritoJaExiste = errors.New("este filme já está na lista de favoritos")

type AdicionarFavoritoInput struct {
	FilmeID       int64  `json:"filmeId" binding:"required"`
	Titulo        string `json:"titulo" binding:"required"`
	CaminhoPoster string `json:"caminhoPoster"`
}

type FavoritoServico interface {
	AdicionarFavorito(usuarioID int64, input AdicionarFavoritoInput) error
	ListarFavoritos(usuarioID int64) ([]dominio.FilmeFavorito, error)
	RemoverFavorito(usuarioID, filmeID int64) error
}

type favoritoServicoImpl struct {
	repo repositorio.FavoritoRepositorio
}

func NovoFavoritoServico(repo repositorio.FavoritoRepositorio) FavoritoServico {
	return &favoritoServicoImpl{repo: repo}
}

// Lógica de AdicionarFavorito atualizada
func (s *favoritoServicoImpl) AdicionarFavorito(usuarioID int64, input AdicionarFavoritoInput) error {
	// 1. Verifica se o favorito já existe antes de tentar inserir.
	existe, err := s.repo.VerificarExistencia(usuarioID, input.FilmeID)
	if err != nil {
		return err // Retorna erro de banco de dados
	}
	if existe {
		return ErrFavoritoJaExiste // Retorna nosso erro customizado
	}

	// 2. Se não existir, cria e salva o novo favorito.
	favorito := &dominio.FilmeFavorito{
		UsuarioID:     usuarioID,
		FilmeID:       input.FilmeID,
		Titulo:        input.Titulo,
		CaminhoPoster: input.CaminhoPoster,
	}
	return s.repo.Salvar(favorito)
}

func (s *favoritoServicoImpl) ListarFavoritos(usuarioID int64) ([]dominio.FilmeFavorito, error) {
	return s.repo.ListarPorUsuarioID(usuarioID)
}

func (s *favoritoServicoImpl) RemoverFavorito(usuarioID, filmeID int64) error {
	return s.repo.Deletar(usuarioID, filmeID)
}
