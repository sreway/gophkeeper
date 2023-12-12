package postgres

import (
	"errors"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"

	"github.com/jackc/pgerrcode"

	"github.com/sreway/gophkeeper/internal/domain/models"
)

func errHandle(err error) error {
	var e *pgconn.PgError

	if errors.Is(err, pgx.ErrNoRows) {
		return models.ErrNotFound
	}

	if errors.As(err, &e) {
		switch e.Code {
		case pgerrcode.UniqueViolation:
			return models.ErrAlreadyExists
		default:
			return err
		}
	}

	return err
}
