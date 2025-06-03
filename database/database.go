package database

import (
	"database/sql"

	migrate "github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/mattn/go-sqlite3"
)

func InitializeDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "database/database.db")
	if err != nil {
		return nil, err
	}

	// Run the migrations
	m, err := migrate.New("file://database/migrations", "sqlite3://database/database.db")
	if err != nil {
		return nil, err
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return nil, err
	}

	return db, nil
}
