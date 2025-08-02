package repositorio

import (
	"github.com/Andydev0/filmes-backend/internal/dominio"
	"github.com/jmoiron/sqlx"
)

// UsuarioRepositorio define a interface para operações de persistência de usuário.
type UsuarioRepositorio interface {
	Salvar(usuario *dominio.Usuario) error
	BuscarPorEmail(email string) (*dominio.Usuario, error)
}

// usuarioRepositorioSqlx é a implementação da interface usando sqlx.
type usuarioRepositorioSqlx struct {
	db *sqlx.DB
}

// NovoUsuarioRepositorio cria uma nova instância do repositório de usuário.
func NovoUsuarioRepositorio(db *sqlx.DB) UsuarioRepositorio {
	return &usuarioRepositorioSqlx{db: db}
}

// Salvar insere um novo usuário no banco de dados.
func (r *usuarioRepositorioSqlx) Salvar(usuario *dominio.Usuario) error {
	query := "INSERT INTO usuarios (nome, email, senha_hash) VALUES (?, ?, ?)"
	_, err := r.db.Exec(query, usuario.Nome, usuario.Email, usuario.SenhaHash)
	return err
}

// BuscarPorEmail encontra um usuário pelo seu email.
func (r *usuarioRepositorioSqlx) BuscarPorEmail(email string) (*dominio.Usuario, error) {
	var usuario dominio.Usuario
	query := "SELECT * FROM usuarios WHERE email = ?"
	err := r.db.Get(&usuario, query, email)
	if err != nil {
		return nil, err
	}
	return &usuario, nil
}
