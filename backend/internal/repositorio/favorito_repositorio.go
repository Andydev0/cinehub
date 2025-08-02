package repositorio

import (
	"database/sql"

	"github.com/Andydev0/filmes-backend/internal/dominio"
	"github.com/jmoiron/sqlx"
)

type FavoritoRepositorio interface {
	Salvar(favorito *dominio.FilmeFavorito) error
	ListarPorUsuarioID(usuarioID int64) ([]dominio.FilmeFavorito, error)
	Deletar(usuarioID, filmeID int64) error
	VerificarExistencia(usuarioID, filmeID int64) (bool, error)
}

type favoritoRepositorioSqlx struct {
	db *sqlx.DB
}

func NovoFavoritoRepositorio(db *sqlx.DB) FavoritoRepositorio {
	return &favoritoRepositorioSqlx{db: db}
}

func (r *favoritoRepositorioSqlx) Salvar(favorito *dominio.FilmeFavorito) error {
	query := "INSERT INTO filmes_favoritos (usuario_id, filme_id, titulo, caminho_poster) VALUES (?, ?, ?, ?)"
	_, err := r.db.Exec(query, favorito.UsuarioID, favorito.FilmeID, favorito.Titulo, favorito.CaminhoPoster)
	return err
}

func (r *favoritoRepositorioSqlx) ListarPorUsuarioID(usuarioID int64) ([]dominio.FilmeFavorito, error) {
	var favoritos []dominio.FilmeFavorito
	query := "SELECT * FROM filmes_favoritos WHERE usuario_id = ?"
	err := r.db.Select(&favoritos, query, usuarioID)
	return favoritos, err
}

func (r *favoritoRepositorioSqlx) Deletar(usuarioID, filmeID int64) error {
	query := "DELETE FROM filmes_favoritos WHERE usuario_id = ? AND filme_id = ?"
	_, err := r.db.Exec(query, usuarioID, filmeID)
	return err
}

// Implementação do novo método para verificar a existência de um favorito.
func (r *favoritoRepositorioSqlx) VerificarExistencia(usuarioID, filmeID int64) (bool, error) {
	var existe bool
	query := "SELECT EXISTS(SELECT 1 FROM filmes_favoritos WHERE usuario_id = ? AND filme_id = ?)"
	err := r.db.Get(&existe, query, usuarioID, filmeID)
	if err != nil && err != sql.ErrNoRows {
		return false, err
	}
	return existe, nil
}
