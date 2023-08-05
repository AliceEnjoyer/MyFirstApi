package sqlite

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

type Storage struct {
	db *sql.DB
}

func NewDatabase(storagePath string) (*Storage, error) {
	const fn = "storage.sqlite.new" // название функции

	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	stmt, err := db.Prepare(`
	CREATE TABLE IF NOT EXISTS url(
		id INTEGER PRIMARY KEY,
		alias TEXT NOT NULL UNIQUE,
		url TEXT NOT NULL);
	CREATE INDEX IF NOT EXISTS idx_alias ON url(alias);
	`) /*url - полная ссылка, alias - это УникальныЙ текст,
	который будет идти после нашего домена, что бы потом из краткой ссылки
	вставить ссылку из url, то есть будет все так:
	https://www.youtube.com/watch?v=rCJvW2xgnk0&t=1825s ->
	-> https://www.urlShortner.io/restApiGolang  */

	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	_, err = stmt.Exec() // stmt - это просто указатель на запрос в нашу датабазу
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	return &Storage{db: db}, nil

	// нужно изучить миграции баз данных
	// (пригодиться в нормальных проектах)
}
