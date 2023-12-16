package sqlite

import (
	"context"
	"database/sql"
	"embed"
	"errors"

	"github.com/golang-migrate/migrate/v4"
	// migrate tools
	_ "github.com/golang-migrate/migrate/v4/database/sqlite"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	"github.com/golang-migrate/migrate/v4/source/iofs"

	// sqlite3 driver
	_ "github.com/mattn/go-sqlite3"

	"github.com/sreway/gophkeeper/internal/config"
)

type storage struct {
	db *sql.DB
}

//go:embed migrations/*.sql
var migrations embed.FS

func (s *storage) migrate() error {
	source, err := iofs.New(migrations, "migrations")
	if err != nil {
		return err
	}

	driver, err := sqlite3.WithInstance(s.db, &sqlite3.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithInstance("iofs", source, "sqlite", driver)
	if err != nil {
		return err
	}

	return m.Up()
}

func New(ctx context.Context, config *config.SQLite) (*storage, error) {
	s := new(storage)

	db, err := sql.Open("sqlite3", config.DSN)
	if err != nil {
		return nil, err
	}

	s.db = db

	err = s.migrate()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return nil, err
	}

	go func() {
		<-ctx.Done()
		s.db.Close()
	}()

	return s, nil
}
