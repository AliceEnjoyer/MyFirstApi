package sqlite

import (
	"database/sql"
	"fmt"

	"github.com/AliceEnjoyer/MyFirstApi/internal/storage"
	"github.com/mattn/go-sqlite3"
	_ "github.com/mattn/go-sqlite3"
)

// 13)

type Storage struct {
	db *sql.DB
}

func NewDatabase(storagePath string) (*Storage, error) {
	const fn = "storage.sqlite.new" // название функции для возврата ошибок

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

func (s *Storage) SaveUrl(urlToSave, alias string) (int64, error) {
	const fn = "storage.sqlite.SaveUrl"

	stmt, err := s.db.Prepare("INSERT INTO url(url, alias) VALUES (?, ?)")
	if err != nil {
		return 0, fmt.Errorf("%s: %w", fn, err)
	}

	res, err := stmt.Exec(urlToSave, alias)
	if err != nil {
		// TODO: refactor this
		// сначало идет преобразование ошибки в ошибку sqlite3, после чего проверяеться
		// нормально ли преобразовалось и sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique
		if sqliteErr, ok := err.(sqlite3.Error); ok && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return 0, fmt.Errorf("%s: %w", fn, storage.ErrURLExists)
		}
		return 0, fmt.Errorf("%s: %w", fn, err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: %w", fn, err)
	}

	return id, nil
}

func (s *Storage) GetUrl(alias string) (string, error) {
	const fn = "storage.sqlite.GetUrl"

	stmt, err := s.db.Prepare("SELECT url FROM url WHERE alias = ?")
	if err != nil {
		return "", fmt.Errorf("%s: %w", fn, err)
	}

	var res string

	err = stmt.QueryRow(alias).Scan(&res)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("%s: %w", fn, storage.ErrURLNotFound)
		}
		return "", fmt.Errorf("%s: %w", fn, err)
	}

	return res, nil
}

func (s *Storage) DeleteUrl(alias string) error {
	const fn = "storage.sqlite.DeleteUrl"

	stmt, err := s.db.Prepare("DELETE FROM url WHERE alias = ?")
	if err != nil {
		return fmt.Errorf("%s: %w", fn, err)
	}

	_, err = stmt.Exec(alias)
	if err != nil {
		return fmt.Errorf("%s: %w", fn, err)
	}

	return nil
}

func (s *Storage) IfExists(alias string) (bool, error) {
	q := `SELECT COUNT(*) FROM url WHERE alias = ?`
	var res int

	if err := s.db.QueryRow(q, alias).Scan(&res); err != nil {
		return false, fmt.Errorf("can not check if user exists: %s", err.Error())
	}
	return res > 0, nil
}
