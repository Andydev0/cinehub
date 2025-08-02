package servico

import (
	"github.com/Andydev0/filmes-backend/internal/dominio"
	"github.com/Andydev0/filmes-backend/internal/repositorio"
)

type AvaliacaoInput struct {
	Nota       int    `json:"nota" binding:"required,min=1,max=5"`
	Comentario string `json:"comentario"`
}

type AvaliacaoServico interface {
	Criar(usuarioID, filmeID int64, input AvaliacaoInput) error
	ListarPorFilme(filmeID int64) ([]dominio.AvaliacaoComUsuario, error)
}

type avaliacaoServicoImpl struct {
	repo repositorio.AvaliacaoRepositorio
}

func NovaAvaliacaoServico(repo repositorio.AvaliacaoRepositorio) AvaliacaoServico {
	return &avaliacaoServicoImpl{repo: repo}
}

func (s *avaliacaoServicoImpl) Criar(usuarioID, filmeID int64, input AvaliacaoInput) error {
	avaliacao := &dominio.Avaliacao{
		UsuarioID:  usuarioID,
		FilmeID:    filmeID,
		Nota:       input.Nota,
		Comentario: input.Comentario,
	}
	return s.repo.Salvar(avaliacao)
}

func (s *avaliacaoServicoImpl) ListarPorFilme(filmeID int64) ([]dominio.AvaliacaoComUsuario, error) {
	return s.repo.BuscarPorFilmeID(filmeID)
}
