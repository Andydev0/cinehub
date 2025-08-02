package database

import (
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3" // Importa o driver do SQLite
)

// InitDB inicializa a conexão com o banco e cria o schema se necessário.
func InitDB() *sqlx.DB {
	log.Println("Tentando conectar ao banco de dados SQLite...")

	db, err := sqlx.Connect("sqlite3", "./filmes.db")
	if err != nil {
		log.Fatalf("Falha ao conectar ao banco de dados: %v", err)
	}

	err = criarSchema(db)
	if err != nil {
		log.Fatalf("Falha ao criar o schema do banco de dados: %v", err)
	}

	return db
}

// criarSchema executa as instruções SQL para criar as tabelas da aplicação.
func criarSchema(db *sqlx.DB) error {
	schema := `
	CREATE TABLE IF NOT EXISTS usuarios (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		nome TEXT NOT NULL,
		email TEXT NOT NULL UNIQUE,
		senha_hash TEXT NOT NULL
	);

	CREATE TABLE IF NOT EXISTS filmes_favoritos (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		usuario_id INTEGER NOT NULL,
		filme_id INTEGER NOT NULL,
		titulo TEXT NOT NULL,
		caminho_poster TEXT,
		data_adicionado DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (usuario_id) REFERENCES usuarios(id),
		UNIQUE (usuario_id, filme_id)
	);

	-- Adiciona a nova tabela para as avaliações dos usuários.
	CREATE TABLE IF NOT EXISTS avaliacoes (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		usuario_id INTEGER NOT NULL,
		filme_id INTEGER NOT NULL,
		nota INTEGER NOT NULL CHECK(nota >= 1 AND nota <= 5),
		comentario TEXT,
		data_criacao DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (usuario_id) REFERENCES usuarios(id),
		-- Garante que um usuário só pode ter uma avaliação por filme.
		UNIQUE (usuario_id, filme_id)
	);
	`

	_, err := db.Exec(schema)
	return err
}
