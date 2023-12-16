package postgres

import (
	"context"
	"errors"

	"github.com/golang-migrate/migrate/v4"
	// migrate tools
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/sreway/gophkeeper/internal/config"
)

type storage struct {
	pool *pgxpool.Pool
}

func (s *storage) migrate(sourceURL, dsn string) error {
	m, err := migrate.New(sourceURL, dsn)
	if err != nil {
		return err
	}

	return m.Up()
}

func New(ctx context.Context, config *config.Postgres) (*storage, error) {
	poolConfig, err := pgxpool.ParseConfig(config.DSN)
	if err != nil {
		return nil, err
	}

	s := new(storage)

	s.pool, err = pgxpool.ConnectConfig(ctx, poolConfig)
	if err != nil {
		return nil, err
	}

	if len(config.SourceMigrations) > 0 {
		err = s.migrate(config.SourceMigrations, config.DSN)
		if err != nil && !errors.Is(err, migrate.ErrNoChange) {
			return nil, err
		}
	}

	go func() {
		<-ctx.Done()
		s.pool.Close()
	}()

	return s, nil
}
