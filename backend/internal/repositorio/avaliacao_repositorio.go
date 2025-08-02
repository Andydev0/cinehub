package repositorio

import (
	"github.com/Andydev0/filmes-backend/internal/dominio"
	"github.com/jmoiron/sqlx"
)

type AvaliacaoRepositorio interface {
	Salvar(avaliacao *dominio.Avaliacao) error
	BuscarPorFilmeID(filmeID int64) ([]dominio.AvaliacaoComUsuario, error)
}

type avaliacaoRepoSqlx struct{ db *sqlx.DB }

func NovaAvaliacaoRepositorio(db *sqlx.DB) AvaliacaoRepositorio {
	return &avaliacaoRepoSqlx{db: db}
}

func (r *avaliacaoRepoSqlx) Salvar(a *dominio.Avaliacao) error {
	// Usamos "upsert" (INSERT OR REPLACE) para que um usuário só possa ter uma avaliação por filme.
	query := `INSERT OR REPLACE INTO avaliacoes (id, usuario_id, filme_id, nota, comentario) 
	           VALUES ((SELECT id FROM avaliacoes WHERE usuario_id = ? AND filme_id = ?), ?, ?, ?, ?)`
	_, err := r.db.Exec(query, a.UsuarioID, a.FilmeID, a.UsuarioID, a.FilmeID, a.Nota, a.Comentario)
	return err
}

func (r *avaliacaoRepoSqlx) BuscarPorFilmeID(filmeID int64) ([]dominio.AvaliacaoComUsuario, error) {
	var avaliacoes []dominio.AvaliacaoComUsuario
	query := `SELECT a.id, a.nota, a.comentario, a.data_criacao, u.nome 
	          FROM avaliacoes a JOIN usuarios u ON a.usuario_id = u.id 
	          WHERE a.filme_id = ? ORDER BY a.data_criacao DESC`
	err := r.db.Select(&avaliacoes, query, filmeID)
	return avaliacoes, err
}
