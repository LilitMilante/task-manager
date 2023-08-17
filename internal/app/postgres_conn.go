package app

import (
	"database/sql"
	"errors"

	"task-manager/migrations"

	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

func ConnectToPostgres(dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	err = upMigrations(db)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func upMigrations(db *sql.DB) error {
	goose.SetBaseFS(migrations.EmbedMigrations)
	goose.SetLogger(goose.NopLogger())

	err := goose.SetDialect("postgres")
	if err != nil {
		return err
	}

	err = goose.Up(db, ".")
	if err != nil && !errors.Is(err, goose.ErrNoNextVersion) {
		return err
	}

	return nil
}
